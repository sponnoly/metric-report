package handler

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

// Handler carries the metrics cache and the instrumentation time in seconds
type Handler struct {
	MetricsCache                 *cache.Cache
	InstrumentationTimeInSeconds time.Duration
}

// request is the request json for the InsertMetric end point
type request struct {
	Value int `json:"value"`
}

// response is the response json for the GetMetricSum end point
type response struct {
	Value int `json:"value"`
}

// metricData is the data stored in the cache
type metricData struct {
	value int
	time  time.Time
}

// GetMetricSum gets the metric sum upto to the last **ValidTimeInSeconds**
func (h *Handler) GetMetricSum(w http.ResponseWriter, r *http.Request) {
	var b []byte
	key := mux.Vars(r)["key"]

	resp := response{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	cachedVal, ok := h.MetricsCache.Get(key)
	if !ok {
		b, _ = json.Marshal(resp)
		w.Write(b)
		return
	}

	cachedList := cachedVal.(*list.List)
	newList := list.New()

	for element := cachedList.Front(); element != nil; element = element.Next() {
		data := element.Value.(metricData)
		metricTime := data.time
		validMetricTime := metricTime.Add(h.InstrumentationTimeInSeconds * time.Second)

		if validMetricTime.After(time.Now()) {
			resp.Value = resp.Value + data.value
			data := metricData{value: data.value, time: data.time}

			newList.PushBack(data)
		} else {
			h.MetricsCache.Set(key, newList, cache.NoExpiration)
			break
		}
	}

	b, _ = json.Marshal(resp)
	w.Write(b)

	return
}

// InsertMetric sets the metric value in the cache
func (h *Handler) InsertMetric(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	var cachedList *list.List
	b, _ := ioutil.ReadAll(r.Body)
	req := request{}

	err := json.Unmarshal(b, &req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cachedQueue, ok := h.MetricsCache.Get(key)
	if !ok {
		cachedList = list.New()
	} else {
		cachedList = cachedQueue.(*list.List)
	}
	data := metricData{value: req.Value, time: time.Now()}
	cachedList.PushFront(data)

	h.MetricsCache.Set(key, cachedList, cache.NoExpiration)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return
}
