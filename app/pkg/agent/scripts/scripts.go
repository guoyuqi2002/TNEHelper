package scripts

import (
	_ "embed"
	"strings"
)

//go:embed fill_in_company_with_china.js
var fillInCompanyWithChina string

func FillInCompanyWithChina(dr string) string {
	template := fillInCompanyWithChina
	template = strings.ReplaceAll(template, "${DocumentNumber}", dr)
	return template
}

//go:embed click_tab_attachments.js
var ClickTabAttachments string

//go:embed click_tab_expense_lines.js
var ClickTabExpenseLines string

//go:embed collect_attachments.js
var CollectAttachments string

//go:embed collect_expense_line_attachments.js
var CollectExpenseLineAttachments string

//go:embed set_temp_id_for_each_expense_line.js
var SetTempIDForEachExpenseLine string
