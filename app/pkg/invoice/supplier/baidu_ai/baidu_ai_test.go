package baidu_ai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestRequest(t *testing.T) {
	// This test performs a live call to the Baidu OCR API and therefore requires
	// credentials. It is skipped when BAIDU_OCR_API_KEY is not configured.
	if os.Getenv("BAIDU_OCR_API_KEY") == "" {
		t.Skip("BAIDU_OCR_API_KEY not set; skipping live OCR test")
	}

	pdfPath := os.Getenv("BAIDU_OCR_TEST_PDF")
	if pdfPath == "" {
		t.Skip("BAIDU_OCR_TEST_PDF (path to a sample invoice PDF) not set; skipping live OCR test")
	}

	bytes, err := os.ReadFile(pdfPath)
	if err != nil {
		t.Fatal(err.Error())
	}
	req := MultipleInvoiceRequest{
		PdfFile:    base64.StdEncoding.EncodeToString(bytes),
		PdfFileNum: 1,
	}
	res, err := New().Request(req)
	fmt.Printf("err: %+v\n", err)
	bytes, _ = json.MarshalIndent(res, "", "  ")
	fmt.Println("res:\n", string(bytes))
}
