package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rrrrover/faas-autoscaler/scaler"
	"github.com/rrrrover/faas-autoscaler/types"
	"io/ioutil"
	"log"
	"net/http"
)

func NewAutoScaleHandlerFunc(scaler scaler.AutoScaler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, readErr := ioutil.ReadAll(request.Body)
		if readErr != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte("Unable to read alert."))

			log.Println(readErr)
			return
		}

		var req types.PrometheusAlert
		err := json.Unmarshal(body, &req)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte("Unable to parse alert, bad format."))
			log.Println(err)
			return
		}

		var errors []error
		for _, alert := range req.Alerts {
			err = scaler.AutoScale(alert)
			if err != nil {
				errors = append(errors, err)
			}
		}

		if len(errors) > 0 {
			log.Println(errors)
			var errorOutput string
			for d, err := range errors {
				errorOutput += fmt.Sprintf("[%d] %s\n", d, err)
			}
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(errorOutput))
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}
