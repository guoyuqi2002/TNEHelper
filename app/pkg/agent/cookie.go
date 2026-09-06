package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"tne-helper/pkg/common"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const CookiesFile = common.Data + "cookies.json"

var lastCookie = int64(0)

func SaveCookies() {
	if now := time.Now().UnixMilli(); now-lastCookie <= 1000 {
		// skip saves that are too frequent
		return
	} else {
		lastCookie = now
	}

	chromedp.Run(taskCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		_ = os.Rename(CookiesFile, CookiesFile+".old")
		if cookies, err := network.GetCookies().Do(ctx); err != nil {
			panic(err.Error())
		} else if bytes, err := json.MarshalIndent(cookies, "", "  "); err != nil {
			panic(err.Error())
		} else if err := os.WriteFile(CookiesFile, bytes, 0644); err != nil {
			panic(err.Error())
		}
		return nil
	}))
}

func toHttpCookies() []http.Cookie {
	res := []http.Cookie{}
	chromedp.Run(taskCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		if cookies, err := network.GetCookies().Do(ctx); err != nil {
			panic(err.Error())
		} else {
			for _, c := range cookies {
				res = append(res, http.Cookie{
					Name:  c.Name,
					Value: c.Value,
				})
			}
		}
		return nil
	}))
	return res
}

func LoadCookies() {
	chromedp.Run(taskCtx, chromedp.ActionFunc(func(ctx context.Context) error {

		cookies := []*network.Cookie{}
		if bytes, err := os.ReadFile(CookiesFile); err != nil {
			// File not found or other error, skip loading
			fmt.Println("cookies file read error, skip loading")
		} else if err := json.Unmarshal(bytes, &cookies); err != nil {
			fmt.Println("cookies file unmarshal error, skip loading")
		} else {
			var params []*network.CookieParam
			for _, c := range cookies {
				// Set cookie expiry to one year from now
				exp := cdp.TimeSinceEpoch(time.Now().AddDate(1, 0, 0))
				cp := &network.CookieParam{
					Name:         c.Name,
					Value:        c.Value,
					Domain:       c.Domain,
					Path:         c.Path,
					Secure:       c.Secure,
					HTTPOnly:     c.HTTPOnly,
					SameSite:     c.SameSite,
					Priority:     c.Priority,
					SourceScheme: c.SourceScheme,
					SourcePort:   c.SourcePort,
					PartitionKey: c.PartitionKey,
					Expires:      &exp,
				}
				params = append(params, cp)
			}
			if err := network.SetCookies(params).Do(ctx); err != nil {
				panic(err.Error())
			}
		}

		return nil
	}))
}
