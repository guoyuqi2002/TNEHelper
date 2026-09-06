package invoice_db

import (
	"fmt"
	"testing"
)

func TestMatch(t *testing.T) {
	{
		index0 := Index("033002300111", "89186289", "")
		found := Match(index0)
		fmt.Printf("found: %+v\n", found)
	}

	{
		index1 := Index("", "", "24342000000027096064")
		found := Match(index1)
		fmt.Printf("found: %+v\n", found)
	}
}
