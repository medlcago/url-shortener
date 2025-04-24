package server

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type MetricsServer struct {
	server *http.Server
}

func NewMetricsServer(addr string) *MetricsServer {
	return &MetricsServer{
		server: &http.Server{
			Addr:    addr,
			Handler: promhttp.Handler(),
		},
	}
}

func (m *MetricsServer) Listen() error {
	return m.server.ListenAndServe()
}

func (m *MetricsServer) Shutdown(ctx context.Context) error {
	return m.server.Shutdown(ctx)
}
