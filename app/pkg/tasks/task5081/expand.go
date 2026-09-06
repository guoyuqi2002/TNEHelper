package task5081

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"tne-helper/pkg/common"

	"codeberg.org/tealeg/xlsx/v4"
)

// [STEP 2]
// Expand exports the expanded results to a folder
// overview location: <exportPath>/$TaskName/overview.xlsx
// attachment location: <exportPath>/$TaskName/$ExpenseReport/$filename
func (t *Task) Expand(exportPath string) {
	fmt.Println("exporting expanded results to folder...") // create the export folder
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		panic(err.Error())
	}
	container := exportPath + "/"

	xf := xlsx.NewFile()
	sh, _ := xf.AddSheet("Sheet1")
	head := sh.AddRow()
	for _, col := range []string{"Expense Report", "Current Attachment Index", "Attachment Count", "Attachment Name", "Download URL"} {
		head.AddCell().SetString(col)
	}

	// export to the folder
	total := len(t.state.SortedList)
	for idx, er := range t.state.SortedList {
		// send a progress event to the frontend
		fmt.Printf("__EXPORT_PROGRESS__:%d:%d:attachment\n", idx+1, total)

		// print the Expense Report being processed
		fmt.Printf("exporting deduplicated attachment files: ER %s\n", er)

		if line := t.state.Reports[er]; !line.Done1 {
			row := sh.AddRow()
			row.AddCell().SetString(er)
			row.AddCell().SetInt(-1)
			row.AddCell().SetInt(-1)
			row.AddCell().SetString("download task not finished")
			row.AddCell().SetString("")
		} else if line.NFile == 0 {
			row := sh.AddRow()
			row.AddCell().SetString(er)
			row.AddCell().SetInt(0)
			row.AddCell().SetInt(0)
			row.AddCell().SetString("no attachments")
			row.AddCell().SetString("")
		} else {
			erSanitized := strings.ReplaceAll(er, ":", "_")
			erSanitized = strings.ReplaceAll(erSanitized, " ", "_")
			ead := exportPath + "/" + erSanitized + "/"
			if err := os.MkdirAll(ead, 0755); err != nil {
				panic(err.Error())
			}

			ad := t.attachmentDir(er)

			dup := map[string]int{}

			for idx := range line.NFile {
				src := ad + strconv.Itoa(idx)
				fn := line.OriginalNames[idx]
				row := sh.AddRow()
				row.AddCell().SetString(er)
				row.AddCell().SetInt(idx + 1)
				row.AddCell().SetInt(line.NFile)
				row.AddCell().SetString(fn)
				row.AddCell().SetString(line.DownloadUrls[idx])

				dup[fn] += 1
				if dup[fn] > 1 {
					fn = strconv.Itoa(dup[fn]) + "_" + fn
				}
				dst := ead + fn

				fmt.Println("┏", dst)
				fmt.Println("┗", src)
				common.CopyFile(dst, src)
			}
		}
	}

	overview := container + "overview.xlsx"
	fmt.Println("saving overview to", overview)
	if err := xf.Save(overview); err != nil {
		panic(err.Error())
	}
}
