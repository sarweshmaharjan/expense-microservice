package services

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/go-pdf/fpdf"
	"github.com/sarweshmaharjan/report_generator/data"
	"github.com/sarweshmaharjan/report_generator/model"
)

func GenerateMonthlyFinanceReport(newExpenses []model.MonthlyExpenseDivision) (string, error) {
	pdf := fpdf.New("P", "mm", "A4", "")

	pdf.AddPage()

	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(254, 180, 120) // default color is white

	pdf.CellFormat(40, 10, "Expenses List", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Amount", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Types", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "% of Total Salary", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 10, "Done", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 10, "Dates", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	for _, v := range newExpenses {
		pdf.SetFillColor(255, 255, 255) // default color is white
		switch v.Name {
		case data.TotalInvestment:
			pdf.SetFillColor(133, 255, 149) // green
		case data.TotalSaving:
			pdf.SetFillColor(145, 255, 255) // blue
		case data.TotalLiabilities:
			pdf.SetFillColor(254, 85, 78) // red
		case data.TotalSalary:
			pdf.SetFillColor(220, 151, 255) // purple
		default:
			pdf.SetFillColor(255, 255, 255) // fallback to white
		}

		pdf.CellFormat(40, 10, v.Name, "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", v.Amount), "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 10, v.Type, "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f%%", v.Ratio), "1", 0, "C", true, 0, "")
		pdf.CellFormat(20, 10, "", "1", 0, "C", true, 0, "")
		pdf.CellFormat(25, 10, "", "1", 0, "C", true, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return "", err
	}

	// Encode the PDF content to base64
	base64PDF := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Return the base64-encoded PDF as a string
	return base64PDF, nil
}
