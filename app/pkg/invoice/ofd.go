package invoice

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
)

var (
	ErrBadOfd = errors.New("bad ofd file")
)

type OfdDocument struct {
	XMLName xml.Name `xml:"Document"`
	Pages   []struct {
		Page struct{} `xml:"Page"`
	} `xml:"Pages"`
}

func NPageOfOfd(fp string) (int, error) {
	if z, err := zip.OpenReader(fp); err != nil {
		return -1, err
	} else {
		defer z.Close()

		for _, f := range z.File {
			if f.Name == "Doc_0/Document.xml" {
				if r, err := f.Open(); err != nil {
					return -1, err
				} else {
					defer r.Close()
					d := OfdDocument{}

					if err := xml.NewDecoder(r).Decode(&d); err != nil {
						return -1, err
					} else {
						return len(d.Pages), nil
					}
				}
			}
		}
		return -1, ErrBadOfd
	}
}

func NPageOfOfd2(buf []byte) (int, error) {
	if z, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf))); err != nil {
		return -1, err
	} else {
		for _, f := range z.File {
			if f.Name == "Doc_0/Document.xml" {
				if r, err := f.Open(); err != nil {
					return -1, err
				} else {
					defer r.Close()
					d := OfdDocument{}

					if err := xml.NewDecoder(r).Decode(&d); err != nil {
						return -1, err
					} else {
						return len(d.Pages), nil
					}
				}
			}
		}
		return -1, ErrBadOfd
	}
}
