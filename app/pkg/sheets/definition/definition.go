package definition

import (
	"errors"
)

var ErrInvalidHeadLine = errors.New("invalid head line")

type Definition struct {
	Name         string
	Comment      string
	SheetName    string
	Heads        []string
	HeadRowIndex int
}

var Definitions = []Definition{
	{
		Name:         "s5081",
		Comment:      "Find Expense Report Lines for Organization - Tax Team",
		SheetName:    "^[Ss]heet1$",
		Heads:        []string{"Expense Report Line", "Expense Report", "Expense Report Accounting Date", "Expense Report Line Date", "Receipt", "Expense Item", "Spend Category", "Spend Category ID", "Cost Center ID", "Line Amount in Local Currency", "Currency", "Local Currency", "Line Amount in USD", "Extended Amount in Ledger Base Currency", "Expense Report Line Extended Amount in Reporting Currency", "Memo", "Worktags", "Destination Location", "Destination Country", "Project", "Project Phase", "Project Phase ID", "Billable", "Personal Expense", "Credit Card Transaction?", "Location", "Location ID", "Business Line", "Business Line ID", "Company", "Company ID", "Expense Report Status", "Pay To", "Employee ID", "Expense Report Memo"},
		HeadRowIndex: 8,
	},
	{
		Name:         "gov_invoices",
		Comment:      "All Invoices From Government Tax Bureau",
		SheetName:    `^Summary Sheet\d*$`,
		Heads:        []string{"Seq No", "Invoice Code", "Invoice Number", "E-Invoice Number", "Seller Tax ID", "Seller Name", "Buyer Tax ID", "Buyer Name", "Issue Date", "Tax Category Code", "Special Business Type", "Goods or Taxable Service", "Specification", "Unit", "Quantity", "Unit Price", "Amount", "Tax Rate", "Tax Amount", "Total With Tax", "Invoice Source", "Invoice Type", "Invoice Status", "Is Positive Invoice", "Risk Level", "Issuer", "Remark"},
		HeadRowIndex: 0,
	},
}
