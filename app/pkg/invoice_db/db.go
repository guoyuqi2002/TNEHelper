package invoice_db

import (
	"context"
	"errors"
	"tne-helper/pkg/common"

	"gorm.io/gorm"
)

var db = common.DB
var table = gorm.G[GovInvoice](db)
var ctx = context.Background()

func init() {
	db.AutoMigrate(&GovInvoice{})
}

var ErrBadInvoiceNumber = errors.New("bad invoice number")

// surprisingly, cases exist where all three identifiers are present at once
func Index(invoiceCode, invoiceNumber, eInvoiceNumber string) string {
	return "A:" + eInvoiceNumber + ";B:" + invoiceCode + ";C:" + invoiceNumber
}

func (gi *GovInvoice) BeforeSave(tx *gorm.DB) error {
	index := Index(gi.InvoiceCode, gi.InvoiceNumber, gi.EInvoiceNumber)
	gi.Index = index
	return nil
}

type GovInvoice struct {
	ID int `gorm:"primary_key"`

	// name of the imported file
	// used for undoing an import: if a file was imported multiple times,
	// its records can be deleted via this field
	FileName string `gorm:"index"`

	// non-unique index on the (invoice code, invoice number, e-invoice number) triple
	// either (invoice code, invoice number) or (e-invoice number) is used:
	// either A: e-invoice number
	// or B: invoice code; C: invoice number
	Index string `gorm:"index"`

	SeqNo                 string
	InvoiceCode           string
	InvoiceNumber         string
	EInvoiceNumber        string
	SellerTaxID           string
	SellerName            string
	BuyerTaxID            string
	BuyerName             string
	IssueDate             string
	TaxCategoryCode       string
	SpecialBusinessType   string
	GoodsOrTaxableService string
	Specification         string
	Unit                  string
	Quantity              string
	UnitPrice             string
	Amount                string
	TaxRate               string
	TaxAmount             string
	TotalWithTax          string
	InvoiceSource         string
	InvoiceKind           string
	InvoiceStatus         string
	IsPositiveInvoice     string
	RiskLevel             string
	Issuer                string
	Remark                string
}
