package task5081

import (
	"os"
	"slices"
	"tne-helper/pkg/common"
	"tne-helper/pkg/sheets/generated/s5081"

	"codeberg.org/tealeg/xlsx/v4"
)

type taskSummary map[string]int

func (t *Task) InputFileDumpPath() string {
	return common.Data + t.name + "/input5081.xlsx.bin"
}

// set the input 5081 file
func (t *Task) SetInput(fp string) {
	// back up the old file first
	_ = os.Rename(t.InputFileDumpPath(), t.InputFileDumpPath()+".old")

	// copy it into the internal storage area
	common.CopyFile(t.InputFileDumpPath(), fp)

	// open the file
	if inFile, err := xlsx.OpenFile(fp, xlsx.ColLimit(s5081.NKey)); nil != err {
		panic("failed to open input file: " + err.Error())
	} else {
		// select Sheet1
		in, ok := inFile.Sheet["Sheet1"]
		if !ok {
			panic("Sheet1 not found")
		}

		// verify the header line
		s5081.CheckHeads(in)

		// count the sub-entries
		t.state.Summary = taskSummary{}
		if err := in.ForEachRow(func(_r *xlsx.Row) error {
			rowIn := s5081.TypedRow{Row: _r}
			if idxRow := rowIn.GetCoordinate(); idxRow <= s5081.HeadRowIndex {
				return nil // Skip header row
			}

			if er := rowIn.GetCell(s5081.ꝎExpense_Report).String(); er != "" {
				t.state.Summary[er] += 1
			}

			return nil
		}); nil != err {
			panic("error processing rows: " + err.Error())
		}

		// build the sorted list
		t.state.SortedList = make([]string, 0, len(t.state.Summary))
		for er := range t.state.Summary {
			t.state.SortedList = append(t.state.SortedList, er)
		}
		slices.Sort(t.state.SortedList)

		// save the state
		t.saveState()
	}
}
