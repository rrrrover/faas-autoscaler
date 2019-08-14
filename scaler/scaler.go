package scaler

import (
	"github.com/openfaas/faas/gateway/handlers"
	"github.com/openfaas/faas/gateway/plugin"
	"github.com/rrrrover/faas-autoscaler/types"
	"log"
	"math"
	"net/url"
	"strconv"
)

type scaler struct {
	FunctionsProviderURL url.URL
	authInjector         handlers.AuthInjector
}

func NewAutoScaler(funcProviderURL url.URL, injector handlers.AuthInjector) AutoScaler {
	return &scaler{
		FunctionsProviderURL: funcProviderURL,
		authInjector:         injector,
	}
}

func (s *scaler) ScaleUp(alert types.PrometheusInnerAlert) error {
	if alert.Status != "firing" {
		return nil
	}
	functionName := alert.Labels.FunctionName
	service := plugin.NewExternalServiceQuery(s.FunctionsProviderURL, s.authInjector)
	queryResponse, getErr := service.GetReplicas(functionName)
	if getErr == nil {
		newReplicas := queryResponse.MaxReplicas*queryResponse.ScalingFactor/100 + queryResponse.AvailableReplicas
		if newReplicas >= queryResponse.MaxReplicas {
			newReplicas = queryResponse.MaxReplicas
		}
		if newReplicas == queryResponse.AvailableReplicas {
			return nil
		}
		log.Printf("scale-up function: [%s] %d=>%d", alert.Labels.FunctionName, queryResponse.AvailableReplicas, newReplicas)
		return service.SetReplicas(functionName, newReplicas)
	}
	return getErr
}

func (s *scaler) ScaleDown(alert types.PrometheusInnerAlert) error {
	if alert.Status != "firing" {
		return nil
	}
	functionName := alert.Labels.FunctionName
	service := plugin.NewExternalServiceQuery(s.FunctionsProviderURL, s.authInjector)
	queryResponse, getErr := service.GetReplicas(functionName)
	if getErr == nil {
		deltaReplicas := queryResponse.MaxReplicas * queryResponse.ScalingFactor / 100
		newReplicas := queryResponse.AvailableReplicas - deltaReplicas
		if queryResponse.AvailableReplicas <= queryResponse.MinReplicas+deltaReplicas {
			newReplicas = queryResponse.MinReplicas
		}
		if newReplicas == queryResponse.AvailableReplicas {
			return nil
		}
		log.Printf("scale-down function: [%s] %d=>%d", alert.Labels.FunctionName, queryResponse.AvailableReplicas, newReplicas)
		return service.SetReplicas(functionName, newReplicas)
	}
	return getErr
}

func (s *scaler) AutoScale(alert types.PrometheusInnerAlert) error {
	log.Printf("%+v", alert)
	functionName := alert.Labels.FunctionName
	service := plugin.NewExternalServiceQuery(s.FunctionsProviderURL, s.authInjector)
	queryResponse, getErr := service.GetReplicas(functionName)
	if getErr == nil {
		value, _ := strconv.ParseFloat(alert.Labels.Value, 64)
		targetValue, _ := strconv.ParseFloat(alert.Labels.Target, 64)
		newReplicas := uint64(math.Ceil(float64(queryResponse.AvailableReplicas) * value / targetValue))
		if newReplicas > queryResponse.MaxReplicas {
			newReplicas = queryResponse.MaxReplicas
		}
		if newReplicas < queryResponse.MinReplicas {
			newReplicas = queryResponse.MinReplicas
		}
		log.Printf("auto-scale function: [%s] %d=>%d", alert.Labels.FunctionName, queryResponse.AvailableReplicas, newReplicas)
		return service.SetReplicas(functionName, newReplicas)
	}
	return getErr
}
