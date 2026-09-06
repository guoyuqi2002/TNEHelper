package task5081

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"tne-helper/pkg/common"
)

type Task struct {
	name  string
	state taskState
}

type taskState struct {
	SortedList []string
	Summary    taskSummary
	Reports    map[string]taskStateLine
}

type taskStateLine struct {
	Done1 bool // attachment download finished
	Done2 bool // deduplication finished
	Done3 bool // AI invoice recognition finished

	/**
	 * Done1
	 */
	// file information queried from the system; file names may be duplicated
	// and file contents may be duplicated as well
	// files are represented by file ids and the total file count is recorded
	// original file names are stored in the OriginalNames field
	Attachments   []int
	Lines         [][]int
	NFile         int
	OriginalNames []string // list of original file names
	DownloadUrls  []string // list of download URLs

	/**
	 * Done2
	 */
	UniqueAttachmentIDs []int // deduplicated list of attachment IDs

	/**
	 * Done3
	 */
}

// name is used to create the task directory, summary initializes the task state
func OpenTask(name string) *Task {
	if err := os.MkdirAll(common.Data+name, 0755); err != nil {
		panic(err.Error())
	}

	t := &Task{
		name: name,
	}
	t.loadState()

	return t
}

// get the current task state
func (t *Task) State() taskState {
	state := taskState{
		SortedList: []string{},
		Reports:    map[string]taskStateLine{},
		Summary:    taskSummary{},
	}
	if f, err := os.Open(t.stateFilePath()); err != nil {
		// file not found, return an empty state
	} else {
		defer f.Close()
		d := gob.NewDecoder(f)
		if err := d.Decode(&state); err != nil {
			panic(err.Error())
		}
	}

	for er := range t.state.Summary {
		if _, ok := state.Reports[er]; !ok {
			state.Reports[er] = taskStateLine{}
		}
	}
	for er := range state.Reports {
		if _, ok := t.state.Summary[er]; !ok {
			delete(state.Reports, er)
		}
	}

	return state
}

func (t *Task) loadState() {
	if f, err := os.Open(t.stateFilePath()); err != nil {
		fmt.Println("task state file not found:", t.stateFilePath())
		fmt.Println("initialize task state from scratch")

		// initialize the state
		t.state.SortedList = []string{}
		t.state.Summary = taskSummary{}
		t.state.Reports = map[string]taskStateLine{}

		legacySupport := true
		if legacySupport {
			// try to parse the initial JSON-format state file
			fmt.Println("using legacy JSON state file")
			legacy := common.Data + t.name + "/task5081.json"
			if bytes, err := os.ReadFile(legacy); err != nil {
				// file not found, return an empty state
			} else if err := json.Unmarshal(bytes, &t.state); err != nil {
				panic(err.Error())
			}
		}
	} else {
		defer f.Close()

		d := gob.NewDecoder(f)
		if err := d.Decode(&t.state); err != nil {
			panic(err.Error())
		}
	}
}

func (t *Task) saveState() {
	_ = os.Rename(t.stateFilePath(), t.stateFilePath()+".old")
	if f, err := os.Create(t.stateFilePath()); err != nil {
		panic(err.Error())
	} else {
		defer f.Close()

		e := gob.NewEncoder(f)
		if err := e.Encode(t.state); err != nil {
			panic(err.Error())
		}
	}
}

// tne-helper-data/$name/attachments/$taskLineName/$id
func (t *Task) attachmentDir(er string) string {
	er = strings.ReplaceAll(er, ":", "_")
	er = strings.ReplaceAll(er, " ", "_")
	return common.Data + t.name + "/attachments/" + er + "/"
}

func (t *Task) stateFilePath() string {
	return common.Data + t.name + "/task5081.gob"
}
