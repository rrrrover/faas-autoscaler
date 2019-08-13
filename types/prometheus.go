package types

// PrometheusInnerAlertLabel PrometheusInnerAlertLabel
type PrometheusInnerAlertLabel struct {
	AlertName    string `json:"alertname"`
	FunctionName string `json:"function_name"`
	Value        string `json:"value"`
	Target       string `json:"target"`
}

// PrometheusInnerAlert PrometheusInnerAlert
type PrometheusInnerAlert struct {
	Status string                    `json:"status"`
	Labels PrometheusInnerAlertLabel `json:"labels"`
}

// PrometheusAlert as produced by AlertManager
type PrometheusAlert struct {
	Status   string                 `json:"status"`
	Receiver string                 `json:"receiver"`
	Alerts   []PrometheusInnerAlert `json:"alerts"`
}
