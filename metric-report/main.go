package main

import (
	"log"
	"net/http"

	"github.com/patrickmn/go-cache"

	"github.com/gorilla/mux"
	"github.com/sponnoly/metric-reporter/handler"
)

func main() {
	r := mux.NewRouter()
	h := handler.Handler{
		MetricsCache:                 cache.New(cache.NoExpiration, -1),
		InstrumentationTimeInSeconds: 3600, //Time in Seconds for instrumentation, 1 hour
	}

	r.HandleFunc("/metric/{key}", h.InsertMetric).
		Methods("POST")
	r.HandleFunc("/metric/{key}/sum", h.GetMetricSum).
		Methods("GET")

	//listen at port 8080
	log.Fatal(http.ListenAndServe(":8080", r))
}
