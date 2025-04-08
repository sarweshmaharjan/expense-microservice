package services

import (
	"fmt"
	"math"

	"github.com/sarweshmaharjan/expenses_divider/data"
	"github.com/sarweshmaharjan/expenses_divider/model"
)

const (
	VARYING = "VARYING"
	FIXED   = "FIXED"
	FLOW    = "FLOW"
)

func GenerateMonthlyExpenseDivision(fconfig model.FinancialConfig) []model.MonthlyExpenseDivision {

	if fconfig.CurrentSalary < fconfig.CapIncomeLimit {
		fmt.Println("Current salary is less than cap limit")
		return nil
	}

	totalFixedAmount, totalVaryingAmount, fixedExpenses := calculateFixedExpenses(fconfig)

	remainingBalance := fconfig.CurrentSalary - (totalFixedAmount + totalVaryingAmount)

	adjustVaryingExpenses := allocateRemainingBalance(remainingBalance, totalVaryingAmount, fixedExpenses)

	if fconfig.SalaryCurrency == "NPR" {
		adjustVaryingExpenses = nepaleseExpenseAmountPrecision(adjustVaryingExpenses)
	}

	finalMonthlyExpenses := adjustForDiscrepancy(fconfig.CurrentSalary, adjustVaryingExpenses)
	return generateMonthlyExpenseDivision(finalMonthlyExpenses)
}

func calculateFixedExpenses(fconfig model.FinancialConfig) (float64, float64, model.FinancialConfig) {
	var totalFixedAmount float64
	var totalVaryingAmount float64

	for i, expense := range fconfig.Expenses {
		if !expense.Active {
			continue
		}
		if expense.IsFixed {
			fconfig.Expenses[i].Amount = expense.Max
			totalFixedAmount += expense.Max
		} else {
			totalVaryingAmount += expense.Min
			fconfig.Expenses[i].Amount = expense.Min
		}
	}

	return totalFixedAmount, totalVaryingAmount, fconfig
}

func allocateRemainingBalance(remainingBalance, totalFixedAmount float64, fconfig model.FinancialConfig) model.FinancialConfig {
	extraAmount := 0.0
	equalDivisionBy := 0.0
	remainingDivisionBy := 0.0

	if remainingBalance == 0 {
		return fconfig
	}

	// First loop: Allocate remaining balance, adjust when max limit is hit
	for i, expense := range fconfig.Expenses {
		if !expense.Active {
			continue
		}
		if !expense.IsFixed {
			// Calculate the addition based on the ratio of the current expense amount to the total fixed amount
			additionRatio := math.Floor(remainingBalance * (expense.Amount / totalFixedAmount))
			adjustedAmount := expense.Amount + additionRatio
			fconfig.Expenses[i].Amount = adjustedAmount
			equalDivisionBy++
		}
	}

	for i, expense := range fconfig.Expenses {
		if !expense.Active {
			continue
		}
		if expense.IsFixed {
			continue
		}
		if expense.Amount == 0 {
			continue
		}
		if expense.Max != 0 && expense.Amount > expense.Max {
			diff := expense.Amount - expense.Max
			extraAmount += diff
			fconfig.Expenses[i].Amount = expense.Max
			fconfig.Expenses[i].IsMaxReached = true
			equalDivisionBy--
		}
	}

	// Final allocation loop based on Type: Investment -> Saving -> Liabilities
	expenseTypes := []string{"Investment", "Saving", "Liabilities"}

	for ie, expenseType := range expenseTypes {
		if extraAmount <= 0 {
			return fconfig
		}
		for i, expense := range fconfig.Expenses {
			if !expense.Active {
				continue
			}
			if expense.IsFixed {
				continue
			}
			if expense.IsMaxReached {
				continue
			}
			if expense.Max == 0 && ie == 1 {
				remainingDivisionBy++
			}
			if expense.Type == expenseType {

				if expense.Max > 0 {

					diff := expense.Max - expense.Amount
					ratioAmount := math.Floor(extraAmount * (1 / equalDivisionBy))
					additionAmount := math.Min(diff, ratioAmount)

					if additionAmount == diff {
						fconfig.Expenses[i].IsMaxReached = true
					}

					fconfig.Expenses[i].Amount += additionAmount
					extraAmount -= additionAmount

					if extraAmount <= 0 {
						return fconfig
					}

					continue
				}
				ratioAmount := math.Floor(extraAmount * (1 / equalDivisionBy))
				fconfig.Expenses[i].Amount += ratioAmount
				extraAmount -= ratioAmount

				// Early exit if no remaining balance
				if extraAmount <= 0 {
					return fconfig
				}
			}
		}
	}

	if extraAmount > 0 {
		ratioAmount := math.Floor(extraAmount * (1 / remainingDivisionBy))
		for i, expense := range fconfig.Expenses {
			if !expense.Active {
				continue
			}
			if expense.IsFixed {
				continue
			}
			if expense.IsMaxReached {
				continue
			}
			if expense.Max == 0 {
				fconfig.Expenses[i].Amount += ratioAmount
				extraAmount -= ratioAmount
				if extraAmount <= 0 {
					return fconfig
				}
			}
		}
	}

	return fconfig
}

func nepaleseExpenseAmountPrecision(fconfig model.FinancialConfig) model.FinancialConfig {
	for i, expense := range fconfig.Expenses {
		fconfig.Expenses[i].Amount = math.Floor(expense.Amount)
	}
	return fconfig
}

func adjustForDiscrepancy(currentSalary float64, fconfig model.FinancialConfig) model.FinancialConfig {
	var total float64
	var discrepancy float64

	for _, expense := range fconfig.Expenses {
		total += expense.Amount
	}

	discrepancy = total - currentSalary

	for i, expense := range fconfig.Expenses {
		if !expense.Active {
			continue
		}
		if expense.IsFixed {
			continue
		}
		if !expense.IsFixed {
			if discrepancy > 0 {
				if expense.Amount-discrepancy < expense.Min {
					continue
				}
				fconfig.Expenses[i].Amount -= math.Abs(discrepancy)
				break
			}
			if discrepancy < 0 {
				if expense.Max != 0 {
					continue
				}
				fconfig.Expenses[i].Amount += math.Abs(discrepancy)
				break
			}
		}
	}

	total = 0
	for _, expense := range fconfig.Expenses {
		total += expense.Amount
	}

	discrepancy = total - currentSalary
	if discrepancy != 0 {
		fmt.Printf("Discrepancy not resolved: %f", discrepancy)
	}

	return fconfig
}

func generateMonthlyExpenseDivision(fconfig model.FinancialConfig) []model.MonthlyExpenseDivision {
	var monthlyExpenses []model.MonthlyExpenseDivision
	var totalInvestmentAmount, totalSavingAmount, totalLiabilitiesAmount float64
	var totalInvestmentRatio, totalSavingRatio, totalLiabilitiesRatio float64

	for _, expense := range fconfig.Expenses {
		if !expense.Active {
			continue
		}
		if expense.Type == data.Investment {
			monthlyExpenses = append(monthlyExpenses, model.MonthlyExpenseDivision{
				Name:   expense.Name,
				Amount: expense.Amount,
				Type:   expense.Type,
				Ratio:  math.Round((expense.Amount / fconfig.CurrentSalary) * 100),
			})
			totalInvestmentAmount += expense.Amount
			totalInvestmentRatio += math.Round((expense.Amount / fconfig.CurrentSalary) * 100)
		}
	}

	monthlyExpenses = append(monthlyExpenses, model.MonthlyExpenseDivision{
		Name:   data.TotalInvestment,
		Amount: totalInvestmentAmount,
		Type:   data.LiabilitiesShortHand,
		Ratio:  totalInvestmentRatio,
	})

	for _, expense := range fconfig.Expenses {
		if !expense.Active {
			continue
		}
		if expense.Type == data.Saving {
			monthlyExpenses = append(monthlyExpenses, model.MonthlyExpenseDivision{
				Name:   expense.Name,
				Amount: expense.Amount,
				Type:   expense.Type,
				Ratio:  math.Round((expense.Amount / fconfig.CurrentSalary) * 100),
			})
			totalSavingAmount += expense.Amount
			totalSavingRatio += math.Round((expense.Amount / fconfig.CurrentSalary) * 100)
		}
	}

	monthlyExpenses = append(monthlyExpenses, model.MonthlyExpenseDivision{
		Name:   data.TotalSaving,
		Amount: totalSavingAmount,
		Type:   data.SavingShortHand,
		Ratio:  totalSavingRatio,
	})

	for _, expense := range fconfig.Expenses {
		if !expense.Active {
			continue
		}
		if expense.Type == data.Liabilities {
			monthlyExpenses = append(monthlyExpenses, model.MonthlyExpenseDivision{
				Name:   expense.Name,
				Amount: expense.Amount,
				Type:   expense.Type,
				Ratio:  math.Round((expense.Amount / fconfig.CurrentSalary) * 100),
			})
			totalLiabilitiesAmount += expense.Amount
			totalLiabilitiesRatio += math.Round((expense.Amount / fconfig.CurrentSalary) * 100)
		}
	}

	monthlyExpenses = append(monthlyExpenses, model.MonthlyExpenseDivision{
		Name:   data.TotalLiabilities,
		Amount: totalLiabilitiesAmount,
		Type:   data.InvestmentShortHand,
		Ratio:  totalLiabilitiesRatio,
	})

	monthlyExpenses = append(monthlyExpenses, model.MonthlyExpenseDivision{
		Name:   data.TotalSalary,
		Amount: fconfig.CurrentSalary,
		Type:   "T",
		Ratio:  100,
	})

	return monthlyExpenses
}
