package task5081

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"tne-helper/pkg/agent"
	"tne-helper/pkg/common"
)

const MAX_TRY = 5

// [STEP 1]
// iterate over all entries whose attachments have not been downloaded yet;
// each callback calls Commit to report that the download is complete
func (t *Task) Download() {
	t.iterate(func(idx, total int, er string) error {
		defer agent.SaveCookies()

		tmpContainer := common.Data + "tmp/attachments/"
		tmpDir := tmpContainer + strconv.FormatInt(time.Now().UnixMilli(), 10) + "/"
		// prepare the temporary folder
		if err := os.RemoveAll(tmpContainer); err != nil {
			panic(err.Error())
		} else if err := os.MkdirAll(tmpDir, 0755); err != nil {
			panic(err.Error())
		}
		if attachments, lines, originalNames, downloadUrls, err := agent.Task5081Line(er, tmpDir); err != nil {
			return err
		} else {
			t.Commit(er, tmpDir, attachments, lines, originalNames, downloadUrls)
		}

		return nil
	})
}

func (t *Task) iterate(foo func(idx, total int, er string) error) {
	fmt.Println("downloading attachments...")
	total := len(t.state.Summary)
	failureCount := 0

	for idx, er := range t.state.SortedList {
		fmt.Printf("[%d/%d] %s\n", idx+1, total, er)

		// skip entries that are already done
		if t.state.Reports[er].Done1 {
			continue
		}
	AGAIN:
		if err := foo(idx, total, er); err != nil {
			fmt.Println("error while downloading attachments:", err.Error())
			failureCount += 1
			if failureCount <= MAX_TRY {
				goto AGAIN
			} else {
				panic("too many errors, aborting the task")
			}
		} else {
			failureCount = 0
		}
	}
}

// report the attachment download result for one ER
// er is the ER name
// tmpDir is the temporary folder; attachments, lines, originalNames and downloadUrls
// are respectively the downloaded attachment ID list, the per-line ID lists,
// the original file name list and the download URL list
func (t *Task) Commit(er string, tmpDir string, attachments []int, lines [][]int, originalNames, downloadUrls []string) {
	// verify the counts
	n := 0
	n += len(attachments)
	for _, line := range lines {
		n += len(line)
	}
	if n != len(originalNames) {
		panic("mismatched number of original names")
	}
	if n != len(downloadUrls) {
		panic("mismatched number of download URLs")
	}

	// verify the contents of the attachments folder
	if entries, err := os.ReadDir(tmpDir); err != nil {
		panic(err.Error())
	} else if len(entries) != n {
		panic("bad attachments directory: mismatched number of files")
	} else {
		m := map[int]bool{}
		for _, e := range entries {
			if id, err := strconv.Atoi(e.Name()); err != nil {
				panic(err.Error())
			} else {
				m[id] = true
			}
		}
		for x := range n {
			if !m[x] {
				panic("bad attachments directory: missing file " + strconv.Itoa(x))
			}
		}
	}

	if err := os.MkdirAll(t.attachmentDir(er), 0755); err != nil {
		panic(err.Error())
	}

	// move the files into the task directory
	for x := range n {
		if err := os.Rename(tmpDir+strconv.Itoa(x), t.attachmentDir(er)+strconv.Itoa(x)); err != nil {
			panic(err.Error())
		}
	}

	// update the state
	stateLine := taskStateLine{
		Done1:         true,
		Attachments:   attachments,
		Lines:         lines,
		NFile:         n,
		OriginalNames: originalNames,
		DownloadUrls:  downloadUrls,
	}

	t.state.Reports[er] = stateLine
	t.saveState()
}
