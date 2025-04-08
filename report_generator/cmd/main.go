package main

import (
	"encoding/json"
	"net/http"

	"github.com/sarweshmaharjan/report_generator/model"
	"github.com/sarweshmaharjan/report_generator/services"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var expenses []model.MonthlyExpenseDivision
		if err := json.NewDecoder(r.Body).Decode(&expenses); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		pdfBytes, err := services.GenerateMonthlyFinanceReport(expenses)
		if err != nil {
			http.Error(w, "Error generating report", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(pdfBytes))
	})

	http.ListenAndServe("report_generator:8003", nil)
}
