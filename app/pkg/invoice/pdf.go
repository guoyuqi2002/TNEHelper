package invoice

import (
	"bytes"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// number of PDF pages, in the range [1, N]
func NPageOfPdf(fp string) (int, error) {
	if f, err := os.Open(fp); err != nil {
		return -1, err
	} else {
		defer f.Close()
		conf := model.NewDefaultConfiguration()
		conf.Cmd = model.EXTRACTPAGES

		if ctx, err := api.ReadValidateAndOptimize(f, conf); err != nil {
			return -1, err
		} else {
			return ctx.PageCount, nil
		}
	}
}

func NPageOfPdf2(buf []byte) (int, error) {
	rs := bytes.NewReader(buf)
	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.EXTRACTPAGES

	if ctx, err := api.ReadValidateAndOptimize(rs, conf); err != nil {
		return -1, err
	} else {
		return ctx.PageCount, nil
	}

}

func SplitPdfLegacy(buf []byte, nPage int) ([][]byte, error) {
	rs := bytes.NewReader(buf)
	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.EXTRACTPAGES

	ctx, err := api.ReadValidateAndOptimize(rs, conf)
	if err != nil {
		return nil, err
	}

	ret := [][]byte{}
	for idx := range nPage {
		pageNr := idx + 1 // starting from 1
		if r, err := api.ExtractPage(ctx, pageNr); err != nil {
			panic(err.Error())
		} else if bytes, err := io.ReadAll(r); err != nil {
			panic(err.Error())
		} else {
			ret = append(ret, bytes)
		}
	}

	return ret, nil
}

type splitPdf struct {
	Payload []byte
	IsImage bool
}

func SplitPdf(buf []byte, nPage int) ([]splitPdf, error) {
	rs := bytes.NewReader(buf)
	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.EXTRACTPAGES

	ctx, err := api.ReadValidateAndOptimize(rs, conf)
	if err != nil {
		return nil, err
	}

	ret := []splitPdf{}
	for idx := range nPage {
		pageNr := idx + 1 // starting from 1

		//fmt.Println(">>>>>>>>>>>>>>>> page:", pageNr)
		if ctxNew, err := pdfcpu.ExtractPages(ctx, []int{pageNr}, false); err != nil {
			panic(err.Error())
		} else {

			// first extract images larger than 1000x1000 from this PDF page
			bypass := [][]byte{}
			for _, obj := range ctxNew.XRefTable.Table {
				if obj, ok := obj.Object.(types.StreamDict); ok {
					if obj.Dict["Type"] == types.Name("XObject") && obj.Dict["Subtype"] == types.Name("Image") {
						width, height := obj.Dict.IntEntry("Width"), obj.Dict.IntEntry("Height")
						if width != nil && height != nil {
							width := *width
							height := *height
							if width >= 1000 && height >= 1000 && len(obj.Raw) > 0 {
								//fmt.Printf("%d: [%dx%d] %v\n", idx, width, height, obj)
								bypass = append(bypass, obj.Raw)
							}
						}
					}
				}
			}

			if len(bypass) == 1 {
				// if there is exactly one such image, recognize the image directly
				ret = append(ret, splitPdf{Payload: bypass[0], IsImage: true})
			} else {
				// otherwise return this PDF page
				b := bytes.NewBuffer(nil)
				if err := api.WriteContext(ctxNew, b); err != nil {
					panic(err.Error())
				}
				ret = append(ret, splitPdf{Payload: b.Bytes(), IsImage: false})
			}
		}
	}

	return ret, nil
}
