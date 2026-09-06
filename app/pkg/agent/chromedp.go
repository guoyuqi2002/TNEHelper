package agent

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
)

const profilePath = "./profile/"

// HomeURL returns the Workday tenant home page URL, configurable via the
// WORKDAY_HOME_URL environment variable.
func HomeURL() string {
	if u := os.Getenv("WORKDAY_HOME_URL"); u != "" {
		return u
	}
	return "https://www.myworkday.com/YOUR_TENANT/d/home.htmld"
}

func debugger(a *chromedp.ExecAllocator) {
	chromedp.Flag("headless", false)(a)
	// Like in Puppeteer.
	chromedp.Flag("hide-scrollbars", false)(a)
	chromedp.Flag("mute-audio", true)(a)

}

var allocCtx, taskCtx context.Context
var cancel0, cancel1 context.CancelFunc

var ChDownload = make(chan browser.EventDownloadProgress)

func TaskCtx() context.Context {
	return taskCtx
}

func Initialize() {
	if err := os.MkdirAll(profilePath, 0755); nil != err {
		panic(err.Error())
	}
	profilePath, _ := filepath.Abs(profilePath)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(profilePath),
		debugger,
	)

	allocCtx, cancel0 = chromedp.NewExecAllocator(context.Background(), opts...)

	taskCtx, cancel1 = chromedp.NewContext(allocCtx)

	if err := chromedp.Run(taskCtx, chromedp.EmulateViewport(1920, 1080)); err != nil {
		panic(err.Error())
	}

	chromedp.ListenTarget(taskCtx, func(v interface{}) {
		switch ev := v.(type) {
		case *browser.EventDownloadProgress:
			ChDownload <- *ev
		default:
			// do nothing
			// fmt.Printf("unknown event: %T, %+v\n", ev, ev)
		}
	})

	chromedp.Run(taskCtx,
		browser.
			SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(downloadPath()).
			WithEventsEnabled(true))
}

func Login() {
	chromedp.Run(taskCtx, chromedp.Navigate(HomeURL()))

	chromedp.Run(taskCtx, chromedp.WaitReady("div[data-uxi-element-id='signout-button']"))
}

func downloadPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	if runtime.GOOS == "windows" {
		return wd + `\blob\`
	}

	return wd + "/blob/"
}
