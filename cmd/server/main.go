package main

import (
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		t := r.PathValue("type")
		n := r.PathValue("name")
		v := r.PathValue("value")

		if n == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		switch t {
		case "gauge":
			_, err := strconv.ParseFloat(v, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case "counter":
			_, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return

		}

		w.WriteHeader(http.StatusOK)
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

type MemStorage struct {
	storageGauge   map[string]float64
	storageCounter map[string]int64
}
