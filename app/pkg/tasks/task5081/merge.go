package task5081

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tne-helper/pkg/invoice"
	"tne-helper/pkg/invoice/supplier/baidu_ai"
	"tne-helper/pkg/invoice_db"

	"codeberg.org/tealeg/xlsx/v4"
)

type InvoiceType = baidu_ai.MultipleInvoiceType

type ResultMap = baidu_ai.MultipleInvoiceResponseResult

// custom train-ticket types: distinguish tickets with and without a ticket number
const (
	TrainTicket_New          InvoiceType = "train_ticket_new"          // train ticket without a ticket number, needs matching
	TrainTicket_Old          InvoiceType = "train_ticket_old"          // train ticket with a ticket number, no matching needed
	VatInvoice_namebatched   InvoiceType = "vat_invoice_namebatched"   // VAT invoice whose buyer matched
	VatInvoice_nameupbatched InvoiceType = "vat_invoice_nameupbatched" // VAT invoice whose buyer did not match
)

// drops the first column (Seq No)
func govInvoiceToSlice(gi *invoice_db.GovInvoice) []string {
	return []string{
		gi.InvoiceCode,
		gi.InvoiceNumber,
		gi.EInvoiceNumber,
		gi.SellerTaxID,
		gi.SellerName,
		gi.BuyerTaxID,
		gi.BuyerName,
		gi.IssueDate,
		gi.TaxCategoryCode,
		gi.SpecialBusinessType,
		gi.GoodsOrTaxableService,
		gi.Specification,
		gi.Unit,
		gi.Quantity,
		gi.UnitPrice,
		gi.Amount,
		gi.TaxRate,
		gi.TaxAmount,
		gi.TotalWithTax,
		gi.InvoiceSource,
		gi.InvoiceKind,
		gi.InvoiceStatus,
		gi.IsPositiveInvoice,
		gi.RiskLevel,
		gi.Issuer,
		gi.Remark,
	}
}

// extract the amount from a train ticket using a regular expression
var numberRegex = regexp.MustCompile(`-?\d+\.?\d*`)

func extractNumberString(s string) string {
	if s == "" {
		return ""
	}

	match := numberRegex.FindString(s)
	if match == "" {
		return ""
	}

	if _, err := strconv.ParseFloat(match, 64); err == nil {
		return match
	}

	return ""
}

// display names used as Excel sheet names
var invoiceTypeSheetNames = map[InvoiceType]string{
	baidu_ai.VatInvoice:       "VAT Invoice",
	VatInvoice_namebatched:    "VAT Invoice",
	VatInvoice_nameupbatched:  "VAT Invoice (Other)",
	baidu_ai.TaxiReceipt:      "Taxi Receipt",
	TrainTicket_Old:           "Train Ticket (Old)",
	TrainTicket_New:           "Train Ticket (New)",
	baidu_ai.QuotaInvoice:     "Quota Invoice",
	baidu_ai.AirTicket:        "Air Itinerary",
	baidu_ai.PrintedInvoice:   "Printed Invoice",
	baidu_ai.TollInvoice:      "Toll Invoice",
	baidu_ai.TaxiOnlineTicket: "Ride-hailing Itinerary",
	__UNKNOWN__:               "Unknown Type",
}

// knownPurchaserKeyword returns the purchaser-name keyword (configurable via the
// KNOWN_PURCHASER_KEYWORD environment variable) used to classify unmatched VAT
// invoices by buyer.
func knownPurchaserKeyword() string {
	if kw := os.Getenv("KNOWN_PURCHASER_KEYWORD"); kw != "" {
		return kw
	}
	return "YOUR_COMPANY_NAME"
}

// [STEP 5]
// merge the recognition results into the Microsoft Excel format requested by the client
// saves two Excel files and returns the paths they were saved to
func (t *Task) Merge(exportPath string) (matchedPath string, unmatchedPath string) {

	// create the two Excel files
	matchedFile := xlsx.NewFile()
	unmatchedFile := xlsx.NewFile()

	// all known invoice types (unknown excluded)
	knownTypes := []InvoiceType{
		baidu_ai.VatInvoice,
		baidu_ai.TaxiReceipt,
		TrainTicket_New,
		TrainTicket_Old,
		baidu_ai.QuotaInvoice,
		baidu_ai.AirTicket,
		baidu_ai.PrintedInvoice,
		baidu_ai.TollInvoice,
		baidu_ai.TaxiOnlineTicket,
	}

	// TrainTicket_Old, TaxiReceipt, QuotaInvoice, TollInvoice, TaxiOnlineTicket, PrintedInvoice all go directly into matched
	unmatchedTypes := []InvoiceType{
		VatInvoice_namebatched,   // needs gov_invoice matching
		VatInvoice_nameupbatched, //
		TrainTicket_New,          // needs gov_invoice matching
		baidu_ai.AirTicket,       // needs gov_invoice matching
		__UNKNOWN__,              // unknown type
	}

	// header line for gov_invoice columns
	govInvoiceHeaders := []string{
		"Invoice Code", "Invoice Number", "E-Invoice Number", "Seller Tax ID", "Seller Name",
		"Buyer Tax ID", "Buyer Name", "Issue Date", "Tax Category Code", "Special Business Type",
		"Goods or Taxable Service", "Specification", "Unit", "Quantity", "Unit Price", "Amount",
		"Tax Rate", "Tax Amount", "Total With Tax", "Invoice Source", "Invoice Type", "Invoice Status",
		"Is Positive Invoice", "Risk Level", "Issuer", "Remark",
	}

	// pre-create the sheets
	matchedSheets := make(map[InvoiceType]*xlsx.Sheet)
	for _, invoiceType := range knownTypes {
		sheet, _ := matchedFile.AddSheet(invoiceTypeSheetNames[invoiceType])
		m := invoiceMaps[invoiceType]
		headerRow := sheet.AddRow()
		headerRow.AddCell().SetString("Expense Report ID")

		// for types that need gov_invoice matching, the matched sheet only contains the gov_invoice headers
		if invoiceType == baidu_ai.VatInvoice || invoiceType == TrainTicket_New || invoiceType == baidu_ai.AirTicket {
			for _, h := range govInvoiceHeaders {
				headerRow.AddCell().SetString(h)
			}
		} else {
			for _, h := range m.Headers {
				headerRow.AddCell().SetString(h)
			}
		}
		matchedSheets[invoiceType] = sheet
	}

	unmatchedSheets := make(map[InvoiceType]*xlsx.Sheet)
	for _, invoiceType := range unmatchedTypes {
		sheet, _ := unmatchedFile.AddSheet(invoiceTypeSheetNames[invoiceType])
		headerRow := sheet.AddRow()
		headerRow.AddCell().SetString("Expense Report ID")

		switch invoiceType {
		case VatInvoice_namebatched, VatInvoice_nameupbatched, TrainTicket_New, baidu_ai.AirTicket:
			// for types that need gov_invoice matching, the unmatched sheet contains
			// the recognized-field headers + a separator column + the gov_invoice headers
			var m invoiceMap
			switch invoiceType {
			case VatInvoice_namebatched, VatInvoice_nameupbatched:
				m = invoiceMaps[baidu_ai.VatInvoice] // use the VatInvoice headers

			case TrainTicket_New:
				m = invoiceMaps[TrainTicket_New]

			case baidu_ai.AirTicket:
				m = invoiceMaps[baidu_ai.AirTicket]
			}
			for _, h := range m.Headers {
				headerRow.AddCell().SetString(h)
			}
			// add an empty column as a separator
			headerRow.AddCell().SetString("")
			for _, h := range govInvoiceHeaders {
				headerRow.AddCell().SetString(h)
			}

		case __UNKNOWN__:
			m := invoiceMaps[__UNKNOWN__]
			for _, h := range m.Headers {
				headerRow.AddCell().SetString(h)
			}

		default:
			panic("unrecognized invoice type: " + string(invoiceType))
		}
		unmatchedSheets[invoiceType] = sheet
	}

	// define the row-writing helper
	writeRow := func(sheet *xlsx.Sheet, rowData []string) {
		row := sheet.AddRow()
		for _, v := range rowData {
			row.AddCell().SetString(v)
		}
	}
	// the frontend progress bar relies on this output
	total := len(t.state.SortedList)
	for idx, er := range t.state.SortedList {
		// send progress to the frontend, used by the progress bar of the exporting overlay
		fmt.Printf("__EXPORT_PROGRESS__:%d:%d:excel\n", idx+1, total)

		// print the Expense Report currently being processed
		fmt.Printf("processing Expense Report: ER %s\n", er)

		// only process successfully recognized files
		if line := t.state.Reports[er]; !line.Done3 {
			continue
		}

		// read the recognition results of the previous step
		parseResults := []invoice.ParseResult{}
		if parseResult, err := os.ReadFile(t.ParseResultFile(er)); err != nil {
			panic(err.Error())
		} else if err := json.Unmarshal(parseResult, &parseResults); err != nil {
			panic(err.Error())
		}

		// one Expense Report may contain multiple invoice files
		for resultIdx, parseResult := range parseResults {
			// print the index of the invoice file being processed
			fmt.Println("processing invoice file", resultIdx+1)
			// one invoice file may contain multiple invoices
			for _, wr := range parseResult.Result.WordsResult {
				// special handling: for train tickets the ticket number decides which type is used
				actualType := wr.Type
				if wr.Type == baidu_ai.TrainTicket {
					ticketNum := wr.Result.GetField("ticket_num")
					if ticketNum != "" {
						actualType = TrainTicket_Old // has a ticket number
					} else {
						actualType = TrainTicket_New // no ticket number
					}
				}

				// check whether the invoice type is known
				if m, ok := invoiceMaps[actualType]; !ok {
					fmt.Println("unknown invoice type:", wr.Type)
					// include unknown invoice types in the final Excel as well, in a separate sheet
					data := invoiceMaps[__UNKNOWN__].Extractor(wr.Result)
					data.AllFields[0] = string(wr.Type)
					row := append([]string{er}, data.AllFields...)
					writeRow(unmatchedSheets[__UNKNOWN__], row)
				} else {
					data := m.Extractor(wr.Result)
					switch data.Type {

					// these three must be compared against GovInvoice
					case baidu_ai.VatInvoice:
						// VAT invoices use Amount for verification, with special unmatched classification
						index := invoice_db.Index(data.InvoiceCode, data.InvoiceNum, data.InvoiceNumDigit)
						matchedInvoices := invoice_db.Match(index)

						if len(matchedInvoices) > 0 {
							// sum the amounts of all matches
							sum := new(big.Rat)
							for _, inv := range matchedInvoices {
								if inv.Amount != "" {
									amount := new(big.Rat)
									if _, ok := amount.SetString(inv.Amount); ok {
										sum.Add(sum, amount)
									}
								}
							}

							expectedAmount := new(big.Rat)
							amountParsed := false
							if data.TotalAmount != "" {
								if _, ok := expectedAmount.SetString(data.TotalAmount); ok {
									amountParsed = true
								} else {
									panic(fmt.Sprintf("error: cannot parse amount '%s' (ER: %s, invoice type: %s)", data.TotalAmount, er, data.Type))
								}
							} else {
								fmt.Printf("amount is empty (ER: %s, invoice type: %s), will be treated as unmatched\n", er, data.Type)
							}

							if amountParsed && sum.Cmp(expectedAmount) == 0 {
								// matched, write to the matched sheet
								for _, inv := range matchedInvoices {
									row := append([]string{er}, govInvoiceToSlice(&inv)...)
									writeRow(matchedSheets[baidu_ai.VatInvoice], row)
								}
							} else {
								// amounts differ, fall back to matching by name
								for _, inv := range matchedInvoices {
									row := append([]string{er}, data.AllFields...)
									row = append(row, "") // separator column
									row = append(row, govInvoiceToSlice(&inv)...)
									writeRow(unmatchedSheets[VatInvoice_namebatched], row)
								}
							}
						} else {
							// no GovInvoice matched, classify by buyer name
							purchaserName := data.AllFields[5] // buyer name is the 6th column (index 5)
							if strings.Contains(purchaserName, knownPurchaserKeyword()) {
								// put into the known-purchaser sheet, pad 27 empty columns to keep column counts consistent
								row := append([]string{er}, data.AllFields...)
								row = append(row, make([]string, 27)...)
								writeRow(unmatchedSheets[VatInvoice_namebatched], row)
							} else {
								// put into the other-purchaser sheet
								row := append([]string{er}, data.AllFields...)
								row = append(row, make([]string, 27)...)
								writeRow(unmatchedSheets[VatInvoice_nameupbatched], row)
							}
						}

					case TrainTicket_New:
						// train ticket without ticket number: match gov_invoice using TotalWithTax
						matched, unmatched := processIndexAndMatching(data, er, "TotalWithTax")
						for _, row := range matched {
							writeRow(matchedSheets[data.Type], row)
						}
						for _, row := range unmatched {
							writeRow(unmatchedSheets[data.Type], row)
						}

					case baidu_ai.AirTicket:
						// air itinerary recognition is pending, TBD
						matched, unmatched := processIndexAndMatching(data, er, "TotalWithTax")
						for _, row := range matched {
							writeRow(matchedSheets[data.Type], row)
						}
						for _, row := range unmatched {
							writeRow(unmatchedSheets[data.Type], row)
						}
						// goes directly into matched
					case TrainTicket_Old, baidu_ai.TaxiReceipt, baidu_ai.QuotaInvoice, baidu_ai.TollInvoice, baidu_ai.TaxiOnlineTicket, baidu_ai.PrintedInvoice:
						row := append([]string{er}, data.AllFields...)
						writeRow(matchedSheets[data.Type], row)

					default:
						panic("unknown invoice type: " + string(data.Type))
					}
				}
			}
		}
	}

	// use the dup method to handle duplicate file names
	getUniquePath := func(baseName string) string {
		path := exportPath + "/" + baseName
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path
		}

		// file already exists, append a counter
		for counter := 1; ; counter++ {
			newName := fmt.Sprintf("%s(%d).xlsx", baseName[:len(baseName)-5], counter)
			path = exportPath + "/" + newName
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return path
			}
		}
	}

	// save the matched file
	matchedPath = getUniquePath("matched.xlsx")
	if err := matchedFile.Save(matchedPath); err != nil {
		panic(fmt.Sprintf("failed to save matched.xlsx: %v", err))
	}

	// save the unmatched file
	unmatchedPath = getUniquePath("unmatched.xlsx")
	if err := unmatchedFile.Save(unmatchedPath); err != nil {
		panic(fmt.Sprintf("failed to save unmatched.xlsx: %v", err))
	}

	return matchedPath, unmatchedPath
}

type ExtractedData struct {
	// core matching fields
	// invoice type
	Type baidu_ai.MultipleInvoiceType
	// invoice number
	InvoiceNum string
	// invoice code
	InvoiceCode string
	// fully-digitalized (e-)invoice number
	InvoiceNumDigit string
	// total amount
	TotalAmount string
	// all fields
	AllFields []string
}

type invoiceMap struct {
	Extractor func(result ResultMap) ExtractedData
	Headers   []string
}

// processIndexAndMatching consolidates the index lookup and matching logic
// it returns the matched and the unmatched invoices
// govInvoiceSumField: which GovInvoice field to sum, e.g. "Amount" or "TotalWithTax"
func processIndexAndMatching(data ExtractedData, erFullName string, govInvoiceSumField string) (matched [][]string, unmatched [][]string) {
	// build the index from the invoice number information and look up matching GovInvoices
	index := invoice_db.Index(data.InvoiceCode, data.InvoiceNum, data.InvoiceNumDigit)
	matchedInvoices := invoice_db.Match(index)

	if len(matchedInvoices) > 0 {
		// sum the given field over all matches (big.Rat avoids floating-point precision issues)
		sum := new(big.Rat)
		for _, inv := range matchedInvoices {
			// pick the field value according to govInvoiceSumField
			var fieldValue string
			switch govInvoiceSumField {
			case "Amount":
				fieldValue = inv.Amount
			case "TotalWithTax":
				fieldValue = inv.TotalWithTax
			}

			if fieldValue != "" {
				amount := new(big.Rat)
				if _, ok := amount.SetString(fieldValue); ok {
					sum.Add(sum, amount)
				}
			}
		}
		// both amounts are compared exactly using big.Rat
		expectedAmount := new(big.Rat)
		amountParsed := false
		// check whether the amount can be parsed
		if data.TotalAmount != "" {
			if _, ok := expectedAmount.SetString(data.TotalAmount); ok {
				amountParsed = true
			} else {
				panic(fmt.Sprintf("error: cannot parse amount '%s' (ER: %s, invoice type: %s)", data.TotalAmount, erFullName, data.Type))
			}
		} else {
			fmt.Printf("amount is empty (ER: %s, invoice type: %s), will be treated as unmatched\n", erFullName, data.Type)
		}

		// check whether the sum equals the extracted total amount
		if amountParsed && sum.Cmp(expectedAmount) == 0 {
			// matched
			for _, inv := range matchedInvoices {
				row := append([]string{erFullName}, govInvoiceToSlice(&inv)...)
				matched = append(matched, row)
			}
		} else {
			// not matched, merge the output
			for _, inv := range matchedInvoices {
				row := append([]string{erFullName}, data.AllFields...)
				// add a separating empty column
				row = append(row, "") //
				row = append(row, govInvoiceToSlice(&inv)...)
				unmatched = append(unmatched, row)
			}
		}
	} else {
		// no GovInvoice matched, pad 27 empty columns to keep column counts consistent
		row := append([]string{erFullName}, data.AllFields...)
		row = append(row, make([]string, 27)...)
		unmatched = append(unmatched, row)
	}

	return matched, unmatched
}

const __UNKNOWN__ InvoiceType = "__UNKNOWN__"

var invoiceMaps = map[InvoiceType]invoiceMap{
	// handle other unknown invoice types
	__UNKNOWN__: {Headers: []string{"Invoice Type", "Raw Data"},
		Extractor: func(result ResultMap) ExtractedData {
			if unknownInvoices, err := json.Marshal(result); err != nil {
				panic(err.Error())
			} else {
				return ExtractedData{
					Type:      __UNKNOWN__,
					AllFields: []string{"[Invoice Type]", string(unknownInvoices)},
				}
			}
		}},

	// newly added: covers invoices that carry both code+number and the e-invoice number
	baidu_ai.VatInvoice: {Headers: []string{"Service Type", "Invoice Code", "Invoice Number", "E-Invoice Number", "Issue Date", "Buyer Name", "Buyer Unified Social Credit Code",
		"Seller Name", "Seller Tax ID", "Goods Name", "Specification", "Quantity", "Unit Price", "Amount", "Tax Rate", "Total Amount", "Total Tax", "Total Amount (Figures)"},
		Extractor: func(result ResultMap) ExtractedData {
			invoicecode := result.GetField("InvoiceCode")
			invoicenum := result.GetField("InvoiceNum")

			var extractedFields []string
			if len(invoicenum) == 20 {
				extractedFields = []string{
					result.GetField("ServiceType"),
					"",
					"",
					invoicenum,
					result.GetField("InvoiceDate"),
					result.GetField("PurchaserName"),
					result.GetField("PurchaserRegisterNum"),
					result.GetField("SellerName"),
					result.GetField("SellerRegisterNum"),
					result.GetField("CommodityName"),
					result.GetField("CommodityType"),
					result.GetField("CommodityNum"),
					result.GetField("CommodityPrice"),
					result.GetField("CommodityAmount"),
					result.GetField("CommodityTaxRate"),
					result.GetField("TotalAmount"),
					result.GetField("TotalTax"),
					result.GetField("AmountInFiguers"),
				}
			} else if len(invoicecode) == 20 {
				extractedFields = []string{
					result.GetField("ServiceType"),
					"",
					"",
					invoicecode,
					result.GetField("InvoiceDate"),
					result.GetField("PurchaserName"),
					result.GetField("PurchaserRegisterNum"),
					result.GetField("SellerName"),
					result.GetField("SellerRegisterNum"),
					result.GetField("CommodityName"),
					result.GetField("CommodityType"),
					result.GetField("CommodityNum"),
					result.GetField("CommodityPrice"),
					result.GetField("CommodityAmount"),
					result.GetField("CommodityTaxRate"),
					result.GetField("TotalAmount"),
					result.GetField("TotalTax"),
					result.GetField("AmountInFiguers"),
				}
			} else {
				extractedFields = []string{
					result.GetField("ServiceType"),
					invoicecode,
					invoicenum,
					result.GetField("InvoiceNumDigit"),
					result.GetField("InvoiceDate"),
					result.GetField("PurchaserName"),
					result.GetField("PurchaserRegisterNum"),
					result.GetField("SellerName"),
					result.GetField("SellerRegisterNum"),
					result.GetField("CommodityName"),
					result.GetField("CommodityType"),
					result.GetField("CommodityNum"),
					result.GetField("CommodityPrice"),
					result.GetField("CommodityAmount"),
					result.GetField("CommodityTaxRate"),
					result.GetField("TotalAmount"),
					result.GetField("TotalTax"),
					result.GetField("AmountInFiguers"),
				}
			}

			return ExtractedData{
				Type:            baidu_ai.VatInvoice,
				InvoiceCode:     extractedFields[1],
				InvoiceNum:      extractedFields[2],
				InvoiceNumDigit: extractedFields[3],
				TotalAmount:     extractedFields[15],
				AllFields:       extractedFields,
			}
		}},

	baidu_ai.TaxiReceipt: {Headers: []string{"Service Type", "Invoice Number", "Invoice Code", "Date", "Amount", "Fuel Surcharge", "Booking Service Fee", "Total Fare", "Price Per km", "Distance"},
		Extractor: func(result ResultMap) ExtractedData {
			extractedFields := []string{
				result.GetField("ServiceType"),
				result.GetField("InvoiceCode"),
				result.GetField("InvoiceNum"),
				result.GetField("Date"),
				result.GetField("Fare"),
				result.GetField("FuelOilSurcharge"),
				result.GetField("CallServiceSurcharge"),
				result.GetField("TotalFare"),
				result.GetField("PricePerkm"),
				result.GetField("Distance"),
			}
			return ExtractedData{
				Type:      baidu_ai.TaxiReceipt,
				AllFields: extractedFields,
			}
		}},

	// e-invoice number column added, invoice number column removed
	// TrainTicket_Old: train ticket with a ticket number, goes directly into matched
	TrainTicket_Old: {Headers: []string{"Service Type", "Buyer Name", "Buyer Unified Social Credit Code", "Ticket Number", "E-Invoice Number", "Electronic Ticket Number", "Issue Date", "Ticket fare; refund fee when the ticket is a refund", "Amount Excl. Tax", "Tax Rate", "Tax Amount", "Passenger Name", "Train Number", "Departure Date"},
		Extractor: func(result ResultMap) ExtractedData {
			extractedFields := []string{
				result.GetField("ServiceType"),
				result.GetField("purchaser_name"),
				result.GetField("purchaser_register_num"),
				result.GetField("ticket_num"),
				result.GetField("invoice_num"),
				result.GetField("elec_ticket_num"),
				result.GetField("invoice_date"),
				result.GetField("ticket_rates"),
				result.GetField("fare"),
				result.GetField("tax_rate"),
				result.GetField("tax"),
				result.GetField("name"),
				result.GetField("train_num"),
				result.GetField("date"),
			}
			return ExtractedData{
				Type:      TrainTicket_Old,
				AllFields: extractedFields,
			}
		}},

	// TrainTicket_New: train ticket without a ticket number (new-style e-ticket), needs gov_invoice matching
	TrainTicket_New: {Headers: []string{"Service Type", "Buyer Name", "Buyer Unified Social Credit Code", "Ticket Number", "E-Invoice Number", "Electronic Ticket Number", "Issue Date", "Ticket fare; refund fee when the ticket is a refund", "Amount Excl. Tax", "Tax Rate", "Tax Amount", "Passenger Name", "Train Number", "Departure Date"},
		Extractor: func(result ResultMap) ExtractedData {
			ticketRates := result.GetField("ticket_rates")
			ticketRatesNumber := extractNumberString(ticketRates)

			extractedFields := []string{
				result.GetField("ServiceType"),
				result.GetField("purchaser_name"),
				result.GetField("purchaser_register_num"),
				result.GetField("ticket_num"),
				result.GetField("invoice_num"),
				result.GetField("elec_ticket_num"),
				result.GetField("invoice_date"),
				ticketRates,
				result.GetField("fare"),
				result.GetField("tax_rate"),
				result.GetField("tax"),
				result.GetField("name"),
				result.GetField("train_num"),
				result.GetField("date"),
			}
			return ExtractedData{
				Type:            TrainTicket_New,
				InvoiceCode:     "",
				InvoiceNum:      "",
				InvoiceNumDigit: extractedFields[4],
				TotalAmount:     ticketRatesNumber,
				AllFields:       extractedFields,
			}
		}},

	baidu_ai.QuotaInvoice: {Headers: []string{"Service Type", "Invoice Code", "Invoice Number", "Amount", "Invoice Location", "Invoice Name"},
		Extractor: func(result ResultMap) ExtractedData {
			extractedFields := []string{
				result.GetField("ServiceType"),
				result.GetField("invoice_code"),
				result.GetField("invoice_number"),
				result.GetField("Invoice_rate"),
				result.GetField("Location"),
				result.GetField("invoice_type"),
			}
			return ExtractedData{
				Type:      baidu_ai.QuotaInvoice,
				AllFields: extractedFields,
			}
		}},

	// air itineraries have no invoice number or code, related columns removed
	baidu_ai.AirTicket: {Headers: []string{"Service Type", "E-Invoice Number", "Total Amount"},
		Extractor: func(result ResultMap) ExtractedData {
			extractedFields := []string{
				result.GetField("ServiceType"),
				result.GetField("invoice_num"),
				result.GetField("total_amount"),
			}
			return ExtractedData{
				Type:            baidu_ai.AirTicket,
				InvoiceCode:     "",
				InvoiceNum:      "",
				InvoiceNumDigit: extractedFields[1],
				TotalAmount:     extractedFields[2],
				AllFields:       extractedFields,
			}
		}},

	baidu_ai.PrintedInvoice: {Headers: []string{"Service Type", "Invoice Code", "Invoice Number", "Issue Date", "Total Amount (Figures)"},
		Extractor: func(result ResultMap) ExtractedData {
			extractedFields := []string{
				result.GetField("ServiceType"),
				result.GetField("InvoiceCode"),
				result.GetField("InvoiceNum"),
				result.GetField("InvoiceDate"),
				result.GetField("AmountInFiguers"),
			}
			return ExtractedData{
				Type:            baidu_ai.PrintedInvoice,
				InvoiceCode:     extractedFields[1],
				InvoiceNum:      extractedFields[2],
				InvoiceNumDigit: "",
				TotalAmount:     extractedFields[4],
				AllFields:       extractedFields,
			}
		}},

	baidu_ai.TollInvoice: {Headers: []string{"Service Type", "Invoice Code", "Invoice Number", "Date", "Time", "Amount", "Entrance", "Exit"},
		Extractor: func(result ResultMap) ExtractedData {
			extractedFields := []string{
				result.GetField("ServiceType"),
				result.GetField("InvoiceCode"),
				result.GetField("InvoiceNum"),
				result.GetField("OutDate"),
				result.GetField("OutTime"),
				result.GetField("TotalAmount"),
				result.GetField("Entrance"),
				result.GetField("Exit"),
			}
			return ExtractedData{
				Type:      baidu_ai.TollInvoice,
				AllFields: extractedFields,
			}
		}},

	baidu_ai.TaxiOnlineTicket: {Headers: []string{"Service Type", "Service Provider", "Trip Start Time", "Trip End Time", "Passenger Phone", "Application Date", "Total Fare", "Trip Count", "Trip Details", "Pickup Time", "Pickup Date", "Car Type", "Distance", "Start Place", "Destination", "City", "Fare"},
		Extractor: func(result ResultMap) ExtractedData {
			extractedFields := []string{
				result.GetField("ServiceType"),
				result.GetField("service_provider"),
				result.GetField("start_time"),
				result.GetField("destination_time"),
				result.GetField("phone"),
				result.GetField("application_date"),
				result.GetField("total_fare"),
				result.GetField("item_num"),
				result.GetField("items"),
				result.GetField("pickup_time"),
				result.GetField("pickup_date"),
				result.GetField("car_type"),
				result.GetField("distance"),
				result.GetField("start_place"),
				result.GetField("destination_place"),
				result.GetField("city"),
				result.GetField("fare"),
			}
			return ExtractedData{
				Type:      baidu_ai.TaxiOnlineTicket,
				AllFields: extractedFields,
			}
		}},
}
