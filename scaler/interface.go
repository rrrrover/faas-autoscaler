package scaler

import (
	"github.com/rrrrover/faas-autoscaler/types"
)

type AutoScaler interface {
	ScaleUp(alert types.PrometheusInnerAlert) error
	ScaleDown(alert types.PrometheusInnerAlert) error
	AutoScale(alert types.PrometheusInnerAlert) error
}
