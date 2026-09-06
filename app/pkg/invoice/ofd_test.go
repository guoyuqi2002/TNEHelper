package invoice

import (
	"fmt"
	"os"
	"testing"
)

// sample OFD file used by the tests; skipped when not present
const ofdSample = "testdata/sample.ofd"

func TestNPageOfOfd(t *testing.T) {
	if _, err := os.Stat(ofdSample); err != nil {
		t.Skipf("sample file %s not found, skipping", ofdSample)
	}
	n, err := NPageOfOfd(ofdSample)
	fmt.Println("n:", n, ", err:", err)
}

func TestNPageOfOfd2(t *testing.T) {
	if _, err := os.Stat(ofdSample); err != nil {
		t.Skipf("sample file %s not found, skipping", ofdSample)
	}
	if buf, err := os.ReadFile(ofdSample); err != nil {
		panic(err.Error())
	} else {
		n, err := NPageOfOfd2(buf)
		fmt.Println("n:", n, ", err:", err)
	}
}
