package main

import (
	"encoding/json"
	"net/http"

	"github.com/sarweshmaharjan/expenses_divider/model"
	"github.com/sarweshmaharjan/expenses_divider/services"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var financialConfig model.FinancialConfig

		if err := json.NewDecoder(r.Body).Decode(&financialConfig); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		newExpenses := services.GenerateMonthlyExpenseDivision(financialConfig)
		if newExpenses == nil {
			http.Error(w, "Failed to generate monthly expenses", http.StatusInternalServerError)
			return
		}



		response, err := json.Marshal(newExpenses)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	http.ListenAndServe("expenses_divider:8001", nil)
}