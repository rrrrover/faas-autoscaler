package types

import (
	"net/http"
)

type AutoScaleHandlerSet struct {
	ScaleUpHandlerFunc   http.HandlerFunc
	ScaleDownHandlerFunc http.HandlerFunc
	AutoScaleHandlerFunc http.HandlerFunc
}
