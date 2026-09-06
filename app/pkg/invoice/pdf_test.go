package invoice

import (
	"fmt"
	"os"
	"testing"
)

// sample PDF file used by the tests; skipped when not present
const pdfSample = "testdata/sample.pdf"

func TestNPageOfPdf(t *testing.T) {
	if _, err := os.Stat(pdfSample); err != nil {
		t.Skipf("sample file %s not found, skipping", pdfSample)
	}
	n, err := NPageOfPdf(pdfSample)
	fmt.Println("n:", n, ", err:", err)
}

func TestNPageOfPdf2(t *testing.T) {
	if _, err := os.Stat(pdfSample); err != nil {
		t.Skipf("sample file %s not found, skipping", pdfSample)
	}
	if buf, err := os.ReadFile(pdfSample); err != nil {
		panic(err.Error())
	} else {
		n, err := NPageOfPdf2(buf)
		fmt.Println("n:", n, ", err:", err)
	}
}

func TestSplitPdf(t *testing.T) {
	if _, err := os.Stat(pdfSample); err != nil {
		t.Skipf("sample file %s not found, skipping", pdfSample)
	}
	if buf, err := os.ReadFile(pdfSample); err != nil {
		panic(err.Error())
	} else if n, err := NPageOfPdf2(buf); err != nil {
		panic(err.Error())
	} else if arr, err := SplitPdf(buf, n); err != nil {
		panic(err.Error())
	} else {
		outDir := t.TempDir()
		for idx, b := range arr {
			fp := fmt.Sprintf("%s/page_%03d.pdf", outDir, idx+1)
			if b.IsImage {
				fp = fmt.Sprintf("%s/page_%03d.jpg", outDir, idx+1)
			}
			if err := os.WriteFile(fp, b.Payload, 0644); err != nil {
				panic(err.Error())
			}
		}
	}
}
