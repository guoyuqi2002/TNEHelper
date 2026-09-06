package common

import (
	"io"
	"mime"
	"net/http"
)

func DownloadWithCookies(url string, cookies []http.Cookie) (data []byte, filename string, err error) {
	// create a new http.Client each time
	// and assign cookies manually
	client := &http.Client{}

	if req, err := http.NewRequest("GET", url, nil); err != nil {
		return nil, "", err
	} else {
		for _, c := range cookies {
			req.AddCookie(&c)
		}

		if res, err := client.Do(req); err != nil {
			return nil, "", err
		} else {
			defer res.Body.Close()

			_, params, _ := mime.ParseMediaType(res.Header.Get("Content-Disposition"))
			filename = params["filename"]

			data, err = io.ReadAll(res.Body)

			return data, filename, err
		}
	}
}
