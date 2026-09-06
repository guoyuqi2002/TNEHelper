package invoice_db

import (
	"fmt"
	"tne-helper/pkg/sheets/generated/gov_invoices"

	"codeberg.org/tealeg/xlsx/v4"
)

// check whether a sheet with the given name has already been imported
func FileNameImported(fileName string) bool {
	if count, err := table.Where("`FileName` = ?", fileName).Count(ctx, "Index"); err != nil {
		panic(err.Error())
	} else {
		return count > 0
	}
}

// list all imported sheet file names
func AllFileNames() []string {
	fileNames := []string{}
	if err := table.Select("DISTINCT `FileName`").Scan(ctx, &fileNames); err != nil {
		panic(err.Error())
	} else {
		return fileNames
	}
}

// delete all invoice records imported from the given sheet file name
func DeleteByFileName(fileName string) int {
	if rowsAffected, err := table.Where("`FileName` = ?", fileName).Delete(ctx); err != nil {
		panic(err.Error())
	} else {
		return rowsAffected
	}
}

// import a tax-bureau invoice detail sheet
// fp is the path of the Excel file
// fileName identifies the imported file and can later be used to delete the data it imported
// before each import, check whether fileName has already been imported;
// if so, prompt the user whether to delete the existing data first
func Import(fp string, fileName string) {
	fmt.Println("importing tax-bureau invoice detail:", fp)
	if f, err := xlsx.OpenFile(fp); err != nil {
		panic(err.Error())
	} else {
		atLeastOneSheetImported := false
		for name, sheet := range f.Sheet {
			fmt.Println(" - processing sheet:", name)
			if gov_invoices.ReMatchTargetSheet.MatchString(name) {
				atLeastOneSheetImported = true
				buf := []GovInvoice{}
				gov_invoices.CheckHeads(sheet)

				if err := sheet.ForEachRow(func(_r *xlsx.Row) error {
					rowIn := gov_invoices.TypedRow{Row: _r}
					if idxRow := rowIn.GetCoordinate(); idxRow <= gov_invoices.HeadRowIndex {
						return nil // Skip header row
					}

					inv := GovInvoice{
						FileName: fileName,

						SeqNo:                 rowIn.GetCell(gov_invoices.ꝎSeq_No).Value,
						InvoiceCode:           rowIn.GetCell(gov_invoices.ꝎInvoice_Code).Value,
						InvoiceNumber:         rowIn.GetCell(gov_invoices.ꝎInvoice_Number).Value,
						EInvoiceNumber:        rowIn.GetCell(gov_invoices.ꝎE_Invoice_Number).Value,
						SellerTaxID:           rowIn.GetCell(gov_invoices.ꝎSeller_Tax_ID).Value,
						SellerName:            rowIn.GetCell(gov_invoices.ꝎSeller_Name).Value,
						BuyerTaxID:            rowIn.GetCell(gov_invoices.ꝎBuyer_Tax_ID).Value,
						BuyerName:             rowIn.GetCell(gov_invoices.ꝎBuyer_Name).Value,
						IssueDate:             rowIn.GetCell(gov_invoices.ꝎIssue_Date).Value,
						TaxCategoryCode:       rowIn.GetCell(gov_invoices.ꝎTax_Category_Code).Value,
						SpecialBusinessType:   rowIn.GetCell(gov_invoices.ꝎSpecial_Business_Type).Value,
						GoodsOrTaxableService: rowIn.GetCell(gov_invoices.ꝎGoods_or_Taxable_Service).Value,
						Specification:         rowIn.GetCell(gov_invoices.ꝎSpecification).Value,
						Unit:                  rowIn.GetCell(gov_invoices.ꝎUnit).Value,
						Quantity:              rowIn.GetCell(gov_invoices.ꝎQuantity).Value,
						UnitPrice:             rowIn.GetCell(gov_invoices.ꝎUnit_Price).Value,
						Amount:                rowIn.GetCell(gov_invoices.ꝎAmount).Value,
						TaxRate:               rowIn.GetCell(gov_invoices.ꝎTax_Rate).Value,
						TaxAmount:             rowIn.GetCell(gov_invoices.ꝎTax_Amount).Value,
						TotalWithTax:          rowIn.GetCell(gov_invoices.ꝎTotal_With_Tax).Value,
						InvoiceSource:         rowIn.GetCell(gov_invoices.ꝎInvoice_Source).Value,
						InvoiceKind:           rowIn.GetCell(gov_invoices.ꝎInvoice_Type).Value,
						InvoiceStatus:         rowIn.GetCell(gov_invoices.ꝎInvoice_Status).Value,
						IsPositiveInvoice:     rowIn.GetCell(gov_invoices.ꝎIs_Positive_Invoice).Value,
						RiskLevel:             rowIn.GetCell(gov_invoices.ꝎRisk_Level).Value,
						Issuer:                rowIn.GetCell(gov_invoices.ꝎIssuer).Value,
						Remark:                rowIn.GetCell(gov_invoices.ꝎRemark).Value,
					}

					// in sheets exported by the tax bureau, empty e-invoice
					// numbers are filled in as "--"
					// restore them to empty strings here
					if inv.EInvoiceNumber == "--" {
						inv.EInvoiceNumber = ""
					} else if len(inv.EInvoiceNumber) <= 5 {
						panic("bad e-invoice number format: " + inv.EInvoiceNumber)
					}

					buf = append(buf, inv)
					return nil
				}); nil != err {
					panic("error processing rows: " + err.Error())
				}

				const Step = 200
				for idx := 0; idx < len(buf); idx += Step {
					end := idx + Step
					if end > len(buf) {
						end = len(buf)
					}
					fmt.Printf("    - [%d-%d/%d]\n", idx, end, len(buf))
					r := buf[idx:end]
					if err := table.CreateInBatches(ctx, &r, Step); err != nil {
						panic(err.Error())
					}
				}

			}
		}

		if !atLeastOneSheetImported {
			panic("no summary sheet found in the file")
		}
	}
}
