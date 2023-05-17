package api

import (
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

func (api *API) MetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		metrics.WritePrometheus(w, true)
	}
}
