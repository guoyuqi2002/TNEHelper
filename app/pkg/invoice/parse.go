package invoice

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tne-helper/pkg/invoice/supplier/baidu_ai"
)

type ParseResults = []ParseResult
type ParseResult struct {
	AttID   int // attachment ID
	IdxPage int // page index, starting from 0; valid for PDF/OFD files

	// error flags
	ErrComment            string // error description
	IsUnsupportedFileType bool   // unsupported file type
	IsFileTooLarge        bool   // file too large (even after page splitting) for the OCR API to process
	IsBadFile             bool   // file could not be recognized

	Result baidu_ai.MultipleInvoiceResponse
}

var client = baidu_ai.New()

// known errors that cannot be resolved by retrying are returned to the caller directly
func IsKnownBaiduAiError(err error) bool {
	switch err {
	case baidu_ai.ErrImageSizeError, baidu_ai.ErrRecognizeError, baidu_ai.ErrFileFormatError:
		return true
	default:
		return false
	}
}

// process all files under an Expense Report directory
// dir is the path of the raw attachments folder
// uniqueIDs is the deduplicated list of attachment IDs
// originalNames is the list of original attachment file names
// files that are themselves problematic are flagged in the ParseResult
// other cases are unlikely to error; an error is returned only for issues such as
// network problems or the cloud OCR account running out of quota
func ParseFiles(dir string, uniqueIDs []int, originalNames []string) (ParseResults, error) {
	res := []ParseResult{}
	for _, id := range uniqueIDs {
		fp := dir + strconv.Itoa(id)
		if bytes, err := os.ReadFile(fp); err != nil {
			res = append(res, ParseResult{AttID: id, IsBadFile: true, ErrComment: err.Error()})
		} else {
			fn := originalNames[id]
			fmt.Printf(" - %s\n", fn)
			ext := strings.ToLower(filepath.Ext(fn))
			switch ext {
			case ".jpg", ".jpeg", ".png", ".bmp":
				encoded := base64.StdEncoding.EncodeToString(bytes)
				if len(encoded) > baidu_ai.MaxBase64EncodedSize {
					res = append(res, ParseResult{AttID: id, IsFileTooLarge: true, ErrComment: "file too large"})
				} else if r, err := client.Request(baidu_ai.MultipleInvoiceRequest{Image: encoded}); err == nil {
					res = append(res, ParseResult{AttID: id, Result: r})
				} else if IsKnownBaiduAiError(err) {
					res = append(res, ParseResult{AttID: id, IsBadFile: true, ErrComment: err.Error()})
				} else {
					return nil, err
				}
			case ".ofd":
				encoded := base64.StdEncoding.EncodeToString(bytes)
				if len(encoded) > baidu_ai.MaxBase64EncodedSize {
					res = append(res, ParseResult{AttID: id, IsFileTooLarge: true, ErrComment: "file too large"})
				} else if nPage, err := NPageOfOfd2(bytes); err != nil {
					res = append(res, ParseResult{AttID: id, IsBadFile: true, ErrComment: err.Error()})
				} else {
					for idx := range nPage {
						if r, err := client.Request(baidu_ai.MultipleInvoiceRequest{OfdFile: encoded, OfdFileNum: idx + 1}); err == nil {
							res = append(res, ParseResult{AttID: id, IdxPage: idx, Result: r})
						} else if IsKnownBaiduAiError(err) {
							res = append(res, ParseResult{AttID: id, IdxPage: idx, IsBadFile: true, ErrComment: err.Error()})
						} else {
							return nil, err
						}
					}
				}
			case ".pdf":
				if nPage, err := NPageOfPdf2(bytes); err != nil {
					res = append(res, ParseResult{AttID: id, IsBadFile: true, ErrComment: err.Error()})
				} else if splitBytes, err := SplitPdf(bytes, nPage); err != nil {
					res = append(res, ParseResult{AttID: id, IsBadFile: true, ErrComment: err.Error()})
				} else {
					fmt.Printf("   PDF has %d pages\n", nPage)
					for idx, singlePageBytes := range splitBytes {
						fmt.Printf("    - page %d/%d\n", idx+1, nPage)
						encoded := base64.StdEncoding.EncodeToString(singlePageBytes.Payload)
						req := baidu_ai.MultipleInvoiceRequest{PdfFile: encoded}
						if singlePageBytes.IsImage {
							fmt.Println("       this page is a large image, recognizing the image directly")
							req = baidu_ai.MultipleInvoiceRequest{Image: encoded}
						}
						if len(encoded) > baidu_ai.MaxBase64EncodedSize {
							res = append(res, ParseResult{AttID: id, IdxPage: idx, IsFileTooLarge: true, ErrComment: "PDF file still too large after splitting"})
						} else if r, err := client.Request(req); err == nil {
							res = append(res, ParseResult{AttID: id, IdxPage: idx, Result: r})
						} else if IsKnownBaiduAiError(err) {
							res = append(res, ParseResult{AttID: id, IdxPage: idx, IsBadFile: true, ErrComment: err.Error()})
						} else {
							return nil, err
						}
					}
				}

			default:
				res = append(res, ParseResult{AttID: id, IsUnsupportedFileType: true, ErrComment: "unexpected file extension"})
			}
		}
	}

	return res, nil
}
