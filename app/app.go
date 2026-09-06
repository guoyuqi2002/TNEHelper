package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"tne-helper/pkg/agent"
	"tne-helper/pkg/common"
	"tne-helper/pkg/invoice_db"
	"tne-helper/pkg/tasks/task5081"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the Wails application struct
type App struct {
	ctx             context.Context
	currentTask     *task5081.Task
	currentTaskName string
	Initialize1     bool
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{}
}

// Startup is called when the application starts
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	common.StdoutRedirect(func(line string) {
		a.emitLog(line)
	})
}

// emitLog emits a log event
func (a *App) emitLog(msg string) {
	runtime.EventsEmit(a.ctx, "process:log", msg)
}

// emitError emits an error event
func (a *App) emitError(err error) {
	runtime.EventsEmit(a.ctx, "process:error", err.Error())
}

// SelectFile is a generic file-selection function for the frontend.
// It only opens a dialog and returns the path of the selected file.
func (a *App) SelectFile(title string) (string, error) {
	selectedPath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   title,
		Filters: []runtime.FileFilter{{DisplayName: "Excel Files (*.xlsx)", Pattern: "*.xlsx"}},
	})
	if err != nil {
		return "", err
	}
	return selectedPath, nil
}

// SelectDirectory is a generic directory-selection function for the frontend.
// It only opens a dialog and returns the path of the selected directory.
func (a *App) SelectDirectory(title string) (string, error) {
	selectedPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	if err != nil {
		return "", err
	}
	return selectedPath, nil
}

func (a *App) CreateTask(taskName string, filepath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	a.currentTask = task5081.OpenTask(taskName)
	a.currentTaskName = taskName

	// let the task load the input 5081 file
	a.currentTask.SetInput(filepath)
	a.emitLog("unique expense reports: " + strconv.Itoa(len(a.currentTask.State().Summary)))
	return nil
}

func (a *App) OpenTask(taskName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	a.currentTask = task5081.OpenTask(taskName)
	a.currentTaskName = taskName
	return nil
}

// enumerate all non-tmp folder names under the tne-helper-data directory as the task list
func (a *App) Tasks() ([]string, error) {
	if entries, err := os.ReadDir(common.Data); err != nil {
		return nil, err
	} else {
		ret := []string{}
		for _, entry := range entries {
			isDir := entry.IsDir()
			name := entry.Name()

			if isDir && name != "tmp" {
				ret = append(ret, name)
			}
		}
		return ret, nil
	}
}

type TaskState struct {
	NDone1 int
	NDone2 int
	NDone3 int
	Total  int
}

// this call is relatively expensive; call it about once every three seconds
func (a *App) GetTaskState() (TaskState, error) {
	s := a.currentTask.State()
	ret := TaskState{}
	ret.Total = len(s.Reports)
	for _, line := range s.Reports {
		if line.Done1 {
			ret.NDone1++
		}
		if line.Done2 {
			ret.NDone2++
		}
		if line.Done3 {
			ret.NDone3++
		}
	}
	return ret, nil
}

// export the download deduplication results to the given folder
func (a *App) Expand5081(exportPath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	if a.currentTask == nil {
		return fmt.Errorf("task not initialized")
	}

	// perform the export
	a.currentTask.Expand(exportPath)
	return nil
}

// RunTask1 runs the main task pipeline
func (a *App) RunTask1() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	if !a.Initialize1 {
		go agent.ForceWake()
		agent.Initialize()

		a.emitLog("loading login credentials")
		agent.LoadCookies()

		a.emitLog("waiting for the user to log in")
		agent.Login()

		a.emitLog("saving login credentials")
		agent.SaveCookies()

		a.emitLog("browser environment ready")

		a.Initialize1 = true
	}

	a.currentTask.Download()

	a.emitLog("task finished")
	return nil
}

func (a *App) RunTask2() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	a.currentTask.Deduplicate()
	return nil
}

func (a *App) RunTask3() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	a.currentTask.Recognize()
	return nil
}

type GovInvoiceIndex struct {
	Year   string
	Month  string
	Branch string
	Raw    string
}

// parse file names with the structure Year-Month-Branch,
// e.g. for 2025-01-Beijing the Year is 2025, the Month is 01 and the Branch is Beijing;
// Year is 4 digits and Month is 2 digits.
// results are sorted in reverse chronological order, newest first
func (a *App) GovInvoiceIndices() []GovInvoiceIndex {
	AllFileNames := invoice_db.AllFileNames()
	parsedFileNames := []GovInvoiceIndex{}

	for _, fileName := range AllFileNames {
		parts := strings.SplitN(fileName, "-", 3)
		if len(parts) == 3 {
			year := parts[0]
			month := parts[1]
			branch := parts[2]
			raw := fileName
			parsedFileNames = append(parsedFileNames, GovInvoiceIndex{Year: year, Month: month, Branch: branch, Raw: raw})
		} else {
			parsedFileNames = append(parsedFileNames, GovInvoiceIndex{Year: "9999", Month: "99", Branch: "error: " + fileName, Raw: fileName})
		}
	}
	slices.SortFunc(parsedFileNames, func(a, b GovInvoiceIndex) int {
		return strings.Compare(a.Year+a.Month+a.Branch, b.Year+b.Month+b.Branch) * -1
	})
	return parsedFileNames
}

// rawIndex corresponds to the Raw field of GovInvoiceIndex
func (a *App) DeleteGovInvoiceIndex(rawIndex string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	invoice_db.DeleteByFileName(rawIndex)
	return nil
}

// fp is the file path, rawIndex corresponds to the Raw field of GovInvoiceIndex
func (a *App) ImportGovInvoices(fp string, rawIndex string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	invoice_db.Import(fp, rawIndex)
	return nil
}

// generate the two Excel files and save them to the given path
func (a *App) ExportExcel(exportPath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	// generate the two Excel matching files
	matchedPath, unmatchedPath := a.currentTask.Merge(exportPath)

	a.emitLog("saved to: " + matchedPath)
	a.emitLog("saved to: " + unmatchedPath)

	a.emitLog("Excel export complete")
	return nil
}
