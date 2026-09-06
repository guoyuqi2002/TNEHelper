package agent

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"tne-helper/pkg/agent/scripts"
	"tne-helper/pkg/common"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func WithAwaitPromise(p *runtime.EvaluateParams) *runtime.EvaluateParams {
	return p.WithAwaitPromise(true)
}

// Task5081URL returns the Workday task URL, configurable via the
// WORKDAY_TASK_URL environment variable.
func Task5081URL() string {
	if u := os.Getenv("WORKDAY_TASK_URL"); u != "" {
		return u
	}
	return "https://www.myworkday.com/YOUR_TENANT/d/task/YOUR_TASK.htmld"
}

func Task5081Line(er string, tmpAttachmentsDir string) (attachments []int, lines [][]int, originalNames, downloadUrls []string, err error) {
	// Document Number
	dr := strings.TrimPrefix(er, "Expense Report: ")
	// attachment counter
	idxAttachment := 0
	// used for downloading attachments
	cookies := toHttpCookies()

	//
	// First screen
	//
	fmt.Println("entering first screen")
	entry := Task5081URL()

	selButtonOK := `button[data-uxi-actionbutton-action="bpf-submit"]`
	chromedp.Run(taskCtx,
		chromedp.Navigate("about:blank"),
		chromedp.Navigate(entry),
		chromedp.WaitReady(selButtonOK),
	)

	fmt.Println("filling in search criteria")
	selCompanyHierarchy := `div[data-automation-label="Company Hierarchy"]`
	selCompanyChina := `div[data-automation-label="China"]`
	chromedp.Run(taskCtx,
		chromedp.Evaluate(scripts.FillInCompanyWithChina(dr), nil),
		chromedp.WaitReady(selCompanyHierarchy),
		chromedp.Click(selCompanyHierarchy),
		chromedp.WaitReady(selCompanyChina),
		chromedp.Click(selCompanyChina),

		// first OK click closes the company dropdown
		chromedp.Sleep(time.Millisecond*100),
		chromedp.Click(selButtonOK),
		// second OK click submits the query
		chromedp.Sleep(time.Millisecond*100),
		chromedp.Click(selButtonOK),
	)

	//
	// Second screen
	//
	fmt.Println("entering second screen")
	selDetail := `a[title="` + er + `"]`
	chromedp.Run(taskCtx,
		chromedp.WaitReady(selDetail),
		chromedp.Sleep(time.Millisecond*200),
		chromedp.Click(selDetail),
	)

	//
	// Third screen
	//
	fmt.Println("entering third screen")
	// switch to the attachments tab
	selSpinnerContainer := `#spinnerContainer`
	selTabs := `div[data-automation-id="tabLabel"]`
	fmt.Println("switching to the main attachments view")
	chromedp.Run(taskCtx,
		chromedp.WaitReady(selTabs),
		chromedp.WaitNotPresent(selSpinnerContainer),
		chromedp.Evaluate(scripts.ClickTabAttachments, nil),
	)

	result0 := []struct {
		Title string // Title = Filename + Date + Person
		Url   string
	}{}
	chromedp.Run(taskCtx,
		chromedp.Evaluate(scripts.CollectAttachments, &result0, WithAwaitPromise),
	)
	fmt.Println("main attachments:")
	for _, r := range result0 {
		url := r.Url
		if len(url) > 64 {
			url = url[:64] + "..."
		}
		fmt.Printf(" - %s:\n   %s\n", r.Title, url)
	}
	for _, a := range result0 {
		fp := tmpAttachmentsDir + strconv.Itoa(idxAttachment)
		fmt.Println("downloading:", a.Title)
		if bytes, filename, err := common.DownloadWithCookies(a.Url, cookies); err != nil {
			return nil, nil, nil, nil, err
		} else if err := os.WriteFile(fp, bytes, 0644); err != nil {
			return nil, nil, nil, nil, err
		} else {
			attachments = append(attachments, idxAttachment)
			originalNames = append(originalNames, filename)
			downloadUrls = append(downloadUrls, a.Url)

			idxAttachment += 1
		}
	}
	fmt.Println("all main attachments downloaded")

	//
	// Fourth screen
	//
	selLoadingPanel := `div[data-automation-id="wd-LoadingPanel"]`
	chromedp.Run(taskCtx,
		chromedp.WaitReady(selTabs),
		chromedp.Evaluate(scripts.ClickTabExpenseLines, nil),
		chromedp.Sleep(time.Millisecond*500),
		chromedp.WaitNotVisible(selLoadingPanel),
	)

	fmt.Println("counting expense lines")
	lineHasAttachment := []bool{}
	chromedp.Run(taskCtx,
		chromedp.Evaluate(scripts.SetTempIDForEachExpenseLine, &lineHasAttachment),
	)

	fmt.Println("expense line count:", len(lineHasAttachment))
	result1 := [][]struct {
		Title string // Title = Filename + Date + Person
		Url   string
	}{}
	for idx, hasAttachment := range lineHasAttachment {
		if !hasAttachment {
			fmt.Printf(" - skipping line %d, no attachments\n", idx)
			result1 = append(result1, nil)
		} else {
			fmt.Printf(" - processing line %d, has attachments, clicking\n", idx)

			id := "#TNE_HELPER_TMP_ID_" + strconv.Itoa(idx)
			partial := []struct {
				Title string // Title = Filename + Date + Person
				Url   string
			}{}
			chromedp.Run(taskCtx,
				chromedp.Click(id),
				chromedp.WaitReady(id+`[data-automation-selected="true"]`),
				chromedp.Sleep(time.Millisecond*500),
				chromedp.WaitNotVisible(selLoadingPanel),
				chromedp.Evaluate(scripts.CollectExpenseLineAttachments, &partial, WithAwaitPromise),
			)
			if len(partial) == 0 {
				fmt.Println("[warning] no attachments found after clicking, skipping")
			}

			fmt.Println("attachment info:")
			for _, r := range partial {
				url := r.Url
				if len(url) > 64 {
					url = url[:64] + "..."
				}
				fmt.Printf(" - %s:\n   %s\n", r.Title, url)
			}
			result1 = append(result1, partial)
		}
	}
	fmt.Println("preparing to download line attachments")
	for _, b := range result1 {
		p := []int{}
		for _, a := range b {
			fp := tmpAttachmentsDir + strconv.Itoa(idxAttachment)
			fmt.Println("downloading:", a.Title)
			if bytes, filename, err := common.DownloadWithCookies(a.Url, cookies); err != nil {
				return nil, nil, nil, nil, err
			} else if err := os.WriteFile(fp, bytes, 0644); err != nil {
				return nil, nil, nil, nil, err
			} else {
				p = append(p, idxAttachment)
				originalNames = append(originalNames, filename)
				downloadUrls = append(downloadUrls, a.Url)

				idxAttachment += 1
			}
		}
		lines = append(lines, p)
	}
	fmt.Println("all line attachments downloaded")

	time.Sleep(time.Millisecond * 1000)
	fmt.Printf("main attachments: %v, line attachments: %+v, original attachment names: %v\n", attachments, lines, originalNames)
	return attachments, lines, originalNames, downloadUrls, nil
}
