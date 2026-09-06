package invoice_db

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestFileNameImported(t *testing.T) {
	fileName := "invoice-export-beijing-2024.1-3.xlsx"
	imported := FileNameImported(fileName)
	fmt.Println(fileName, "imported:", imported)
}

func TestAllFileNames(t *testing.T) {
	fileNames := AllFileNames()
	fmt.Printf("all file names: %+v\n", fileNames)
}

func TestDeleteByFileName(t *testing.T) {
	fileName := "invoice-export-beijing-2025.11-12.xlsx"
	rowsAffected := DeleteByFileName(fileName)
	fmt.Println("deleted rows affected:", rowsAffected)
}

// TestImportGovInvoice imports every .xlsx file found in the directory given by
// GOV_INVOICE_DIR; it is skipped when the variable is not set.
func TestImportGovInvoice(t *testing.T) {
	dir := os.Getenv("GOV_INVOICE_DIR")
	if dir == "" {
		t.Skip("GOV_INVOICE_DIR not set; skipping import test")
	}
	if entries, err := os.ReadDir(dir); err != nil {
		panic(err.Error())
	} else {
		for _, entry := range entries {
			if name := entry.Name(); strings.HasSuffix(strings.ToLower(name), ".xlsx") {
				fp := dir + "/" + entry.Name()
				Import(fp, name)
			}
		}
	}
}
