package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %v", hits)))
}

func (cfg *apiConfig) handleMetricsReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Metrics reset successfully"))
}
