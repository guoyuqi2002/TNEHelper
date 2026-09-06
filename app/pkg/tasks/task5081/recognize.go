package task5081

import (
	"encoding/json"
	"fmt"
	"os"
	"tne-helper/pkg/invoice"
)

// [STEP 4]
// recognize invoices via the Baidu AI API
// the recognition results of all attachments of an Expense Report are saved in
// t.ParseResultFile(er), i.e.
// tne-helper-data/$name/attachments/$taskLineName/$id/parse_result.json
func (t *Task) Recognize() {
	fmt.Println("recognizing invoices...")
	failureCount := 0

	for idx, er := range t.state.SortedList {
		line := t.state.Reports[er]
		fmt.Printf("[%d/%d] %s\n", idx+1, len(t.state.SortedList), er)
		if !line.Done2 || line.Done3 {
			continue
		}

	AGAIN:
		if res, err := invoice.ParseFiles(t.attachmentDir(er), line.UniqueAttachmentIDs, line.OriginalNames); err != nil {
			fmt.Println("error while calling the recognition API:", err.Error())
			failureCount += 1
			if failureCount <= MAX_TRY {
				goto AGAIN
			} else {
				panic("too many errors, aborting the task")
			}
		} else {
			failureCount = 0
			if bytes, err := json.MarshalIndent(res, "", "  "); err != nil {
				panic(err.Error())
			} else if err := os.WriteFile(t.ParseResultFile(er), bytes, 0644); err != nil {
				panic(err.Error())
			}

			line.Done3 = true
			t.state.Reports[er] = line
			t.saveState()
		}
	}
}

// path of the Baidu attachment recognition result file
// tne-helper-data/$name/attachments/$taskLineName/$id/parse_result.json
func (t *Task) ParseResultFile(er string) string {
	return t.attachmentDir(er) + "parse_result.json"
}
