package task5081

import (
	"fmt"
	"os"
	"testing"
)

// TestParseList parses the 5081 export given by TASK5081_INPUT;
// it is skipped when the variable is not set.
func TestParseList(t *testing.T) {
	fp := os.Getenv("TASK5081_INPUT")
	if fp == "" {
		t.Skip("TASK5081_INPUT not set; skipping parse-list test")
	}
	task := OpenTask("example_task")
	task.SetInput(fp)
	fmt.Println(task.state.Summary)
}
