package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin" // Ensure the gin package is installed via `go get -u github.com/gin-gonic/gin`
	"github.com/sarweshmaharjan/hub/model"
)

func main() {

	router := gin.Default()

	router.GET("/", generateFinance)
	router.Run("hub:8000") // Replace 'localhost' with the service name 'hub'
}

func generateFinance(c *gin.Context) {
	// Step 1: Call json_builder container
	resp1, err := http.Get("http://json_builder:8002/")
	if err != nil {
		http.Error(c.Writer, "Error calling json_builder: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp1.Body.Close()

	payload1, err := io.ReadAll(resp1.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Error reading response from json_builder: %v; resp1:%v", err, payload1)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}

	// Unmarshal payload1 into model.FinancialConfig
	var res model.FinancialConfig
	if err := json.Unmarshal(payload1, &res); err != nil {
		errMsg := fmt.Sprintf("Error unmarshalling response from json_builder: %v; resp1:%v", err, payload1)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}

	// Marshal res back to JSON
	payload1, err = json.Marshal(res)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshalling FinancialConfig: %v; resp1:%v", err, payload1)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}

	// Step 2: Call expenses_divider with payload from Step 1
	resp2, err := http.Post("http://expenses_divider:8001/", "application/json", bytes.NewReader(payload1))
	if err != nil {
		errMsg := fmt.Sprintf("Error calling expenses_divider: %v -- resp2:%v -- resp1:%v", err, resp2, payload1)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}
	defer resp2.Body.Close()

	payload2, err := io.ReadAll(resp2.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Error reading response from expenses_divider: %v -- resp2:%v -- resp1:%v", err, payload2, payload2)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}

	// Unmarshal payload2 into []model.MonthlyExpenseDivision
	var newExpenses []model.MonthlyExpenseDivision
	if err := json.Unmarshal(payload2, &newExpenses); err != nil {
		errMsg := fmt.Sprintf("Error unmarshalling response from expenses_divider: %v -- resp2:%v -- resp1:%v", err, payload2, payload1)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}

	// Marshal newExpenses back to JSON
	payload2, err = json.Marshal(newExpenses)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshalling MonthlyExpenseDivision: %v; resp2:%v", err, payload2)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}

	// Step 3: Call report_generator with payload from Step 2
	resp3, err := http.Post("http://report_generator:8003/", "application/json", bytes.NewReader(payload2))
	if err != nil {
		errMsg := fmt.Sprintf("Error calling report_generator: %v; resp3:%v -- payload2:%v", err, resp3, payload2)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}
	defer resp3.Body.Close()

	// Decode the base64-encoded PDF content from resp3.Body
	payload3, err := io.ReadAll(resp3.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Error reading response from report_generator: %v; resp3:%v", err, resp3)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}

	pdfBytes, err := base64.StdEncoding.DecodeString(string(payload3))
	if err != nil {
		errMsg := fmt.Sprintf("Error decoding base64 PDF content: %v; payload3:%v", err, payload3)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}

	// Write the decoded PDF response to the client
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.Writer.Header().Set("Content-Disposition", "inline; filename=report.pdf")
	if _, err := c.Writer.Write(pdfBytes); err != nil {
		errMsg := fmt.Sprintf("Error writing PDF response to client: %v", err)
		http.Error(c.Writer, errMsg, http.StatusInternalServerError)
		return
	}
}
