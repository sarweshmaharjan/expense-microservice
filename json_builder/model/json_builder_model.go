package model

type FinancialConfig struct {
	SalaryCurrency string    `yaml:"salary_currency"`
	CurrentSalary  float64   `yaml:"current_salary"`
	CapIncomeLimit float64   `yaml:"cap_income_limit"`
	Expenses       []Expense `yaml:"expenses"`
}
type Expense struct {
	Name           string  `yaml:"name"`
	IsFixed        bool    `yaml:"is_fixed"`
	Min            float64 `yaml:"min"`
	Max            float64 `yaml:"max"`
	Type           string  `yaml:"type"`
	ExpectedAmount float64 `yaml:"expected_amount"`
	IsMaxReached   bool    `yaml:"-"`
	IsMinReached   bool    `yaml:"-"`
	Active         bool    `yaml:"active"`
}
