package main

import (
	"encoding/json"
	"net/http"

	"github.com/sarweshmaharjan/json_builder/services"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		res := services.Load()
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(res)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})

	http.ListenAndServe("json_builder:8002", nil)
}