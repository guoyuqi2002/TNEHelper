package task5081

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// [STEP 3]
// deduplicate
func (t *Task) Deduplicate() {
	fmt.Println("deduplicating...")
	extensions := map[string]int{}

	for idx, er := range t.state.SortedList {
		line := t.state.Reports[er]
		fmt.Printf("[%d/%d] %s\n", idx+1, len(t.state.SortedList), er)
		if !line.Done1 || line.Done2 {
			continue
		}

		unique := []int{}
		m := map[[16]byte]bool{}

		for fid := range line.NFile {
			fp := t.attachmentDir(er) + strconv.Itoa(fid)
			if bytes, err := os.ReadFile(fp); err != nil {
				panic(err.Error())
			} else {
				sum := md5.Sum(bytes)
				if !m[sum] {
					// not seen before
					m[sum] = true
					unique = append(unique, fid)
				}
			}
		}
		line.Done2 = true
		line.UniqueAttachmentIDs = unique

		// update
		t.state.Reports[er] = line
	}

	for _, er := range t.state.SortedList {
		line := t.state.Reports[er]
		for _, fid := range line.UniqueAttachmentIDs {
			n := line.OriginalNames[fid]
			ext := filepath.Ext(n)
			ext = strings.ToLower(ext)
			extensions[ext] += 1
		}
	}
	fmt.Printf("file extension statistics: %+v\n", extensions)

	t.saveState()
}
