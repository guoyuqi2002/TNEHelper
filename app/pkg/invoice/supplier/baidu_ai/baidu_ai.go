package baidu_ai

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	apiUrl            = "https://aip.baidubce.com/rest/2.0/ocr/v1/multiple_invoice"
	apiAccessInterval = time.Second / 10
)

// credential returns the API key and name configured via environment variables.
func credential() (apiKey string, apiName string) {
	return os.Getenv("BAIDU_OCR_API_KEY"), os.Getenv("BAIDU_OCR_API_NAME")
}

const (
	MaxBase64EncodedSize = 4*1024*1024 - 256 // slightly under 4MB
)

type baiduAi struct {
	client        http.Client
	lastAccess    time.Time
	postFailCount int
}

func New() *baiduAi {
	b := &baiduAi{}
	return b
}

// maximum number of retries per request
const postFailCountMax = 3

type BaiduAiError struct {
	Code        int
	Msg         string
	Description string
}

func (e BaiduAiError) Error() string {
	return fmt.Sprintf("[%d/%s] %s", e.Code, e.Msg, e.Description)
}

var (
	ErrImageSizeError  = BaiduAiError{Code: 216202, Msg: "image size error", Description: "invalid image size: base64-encoded size must be under 4MB and resolution no higher than 4096x4096; please re-upload the image"}
	ErrRecognizeError  = BaiduAiError{Code: 216630, Msg: "recognize error", Description: "recognition error; please retry and make sure the image contains the corresponding invoice or certificate"}
	ErrFileFormatError = BaiduAiError{Code: 216201, Msg: "file format error", Description: "unsupported image format: currently supported formats are PNG, JPG, JPEG and BMP; please transcode or replace the image"}
)

func (b *baiduAi) Request(req MultipleInvoiceRequest) (res MultipleInvoiceResponse, err error) {
	apiKey, apiName := credential()
	if apiKey == "" {
		return res, fmt.Errorf("BAIDU_OCR_API_KEY not set; see README")
	}
	_ = apiName // optional, used only for credential bookkeeping

	data := url.Values{}
	switch {
	case req.Image != "":
		data.Set("image", req.Image)
	case req.Url != "":
		data.Set("url", req.Url)
	case req.PdfFile != "":
		data.Set("pdf_file", req.PdfFile)
		if req.PdfFileNum != 0 {
			data.Set("pdf_file_num", strconv.Itoa(req.PdfFileNum))
		}
	case req.OfdFile != "":
		data.Set("ofd_file", req.OfdFile)
		if req.OfdFileNum != 0 {
			data.Set("ofd_file_num", strconv.Itoa(req.OfdFileNum))
		}
	default:
		return res, fmt.Errorf("one of image/url/pdf_file/ofd_file must be provided")
	}
	if req.VerifyParameter {
		data.Set("verify_parameter", "true")
	}
	if req.Probability {
		data.Set("probability", "true")
	}
	if req.Location {
		data.Set("location", "true")
	}

	now := time.Now()
	// fmt.Println("now:", now, ", gap:", now.Sub(baiduAiLastAccess))
	if gap := now.Sub(b.lastAccess); gap < apiAccessInterval {
		time.Sleep(apiAccessInterval - gap)
	}
	defer func() { b.lastAccess = time.Now() }()

	if resp, err := b.postForm(apiUrl, data, apiKey); nil != err {
		return res, err
	} else {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&res); nil != err {
			return res, err
		}

		// handle known errors
		switch res.ErrorCode {
		case 216202:
			return res, ErrImageSizeError
		case 216630:
			return res, ErrRecognizeError
		case 216201:
			return res, ErrFileFormatError
		}

		if res.ErrorCode != 0 || res.ErrorMsg != "" {
			return res, fmt.Errorf("Baidu OCR error: [%d] %s", res.ErrorCode, res.ErrorMsg)
		}

		// fmt.Printf("result: %+v\n", response)

		return res, nil
	}
}

type MultipleInvoiceRequest struct {
	// Image data, base64-encoded and then urlencoded; after encoding the size must not exceed 4MB,
	// the shortest side must be at least 15px and the longest side at most 4096px; jpg/jpeg/png/bmp formats are supported.
	// Priority: image > url > pdf_file > ofd_file; when the image field is present, the url/pdf_file/ofd_file fields are ignored
	Image string `json:"image"` // choose one of image/url/pdf_file/ofd_file

	// Full image URL, at most 1024 bytes; the corresponding image must be under 4MB after base64 encoding,
	// the shortest side at least 15px and the longest side at most 4096px; jpg/jpeg/png/bmp formats are supported.
	// Priority: image > url > pdf_file > ofd_file; when the image field is present, the url field is ignored.
	// Make sure hotlink protection is disabled for the URL.
	Url string `json:"url"` // choose one of image/pdf_file/ofd_file

	// PDF file, base64-encoded and then urlencoded; after encoding the size must not exceed 4MB,
	// the shortest side must be at least 15px and the longest side at most 4096px.
	// Priority: image > url > pdf_file > ofd_file; when the image or url field is present, the pdf_file field is ignored
	PdfFile string `json:"pdf_file"` // choose one of image/url/ofd_file

	// Page number of the PDF file to recognize; valid when pdf_file is provided.
	// If omitted, page 1 is recognized by default.
	PdfFileNum int `json:"pdf_file_num"`

	// OFD file, base64-encoded and then urlencoded; after encoding the size must not exceed 4MB,
	// the shortest side must be at least 15px and the longest side at most 4096px.
	// Priority: image > url > pdf_file > ofd_file; when the image, url or pdf_file field is present, the ofd_file field is ignored
	OfdFile string `json:"ofd_file"` // choose one of image/url/pdf_file

	// Page number of the OFD file to recognize; valid when ofd_file is provided.
	// If omitted, page 1 is recognized by default.
	OfdFileNum int `json:"ofd_file_num"`

	// Whether to enable verification, default false. When true, the six fields required by the
	// invoice-verification API are returned; see the API documentation for details.
	VerifyParameter bool `json:"verify_parameter"`

	// Whether to return field-level confidence, default false (not returned)
	Probability bool `json:"probability"`

	// Whether to return field position coordinates, default false (not returned)
	Location bool `json:"location"`
}

type ErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type MultipleInvoiceResponseResult map[string][]struct {
	Row  string `json:"row"`  // row number
	Word string `json:"word"` // content
}

type MultipleInvoiceResponse struct {
	ErrorResponse // error case

	LogId          uint64 `json:"log_id"`           // unique log id, used for troubleshooting
	PdfFileSize    int    `json:"pdf_file_size"`    // total page count of the uploaded PDF, returned when pdf_file is provided
	WordsResultNum uint32 `json:"words_result_num"` // number of recognition results, i.e. element count of words_result
	WordsResult    []struct {
		Probability float64                       `json:"probability"` // confidence of the classification of a single invoice
		Left        int                           `json:"left"`        // horizontal coordinate of the top-left corner of the invoice bounding box
		Top         int                           `json:"top"`         // vertical coordinate of the top-left corner of the invoice bounding box
		Width       int                           `json:"width"`       // width of the invoice bounding box
		Height      int                           `json:"height"`      // height of the invoice bounding box
		Type        MultipleInvoiceType           `json:"type"`        // category of each invoice
		Result      MultipleInvoiceResponseResult `json:"result"`      // recognition result for a single invoice

	} `json:"words_result"` // recognition results
}

type MultipleInvoiceType string

const (
	VatInvoice          MultipleInvoiceType = "vat_invoice"           // VAT invoice
	TaxiReceipt         MultipleInvoiceType = "taxi_receipt"          // taxi receipt
	TrainTicket         MultipleInvoiceType = "train_ticket"          // train ticket
	QuotaInvoice        MultipleInvoiceType = "quota_invoice"         // fixed-amount invoice
	AirTicket           MultipleInvoiceType = "air_ticket"            // air travel itinerary
	RollNormalInvoice   MultipleInvoiceType = "roll_normal_invoice"   // roll invoice
	PrintedInvoice      MultipleInvoiceType = "printed_invoice"       // machine-printed invoice
	PrintedElecInvoice  MultipleInvoiceType = "printed_elec_invoice"  // machine-printed electronic invoice
	BusTicket           MultipleInvoiceType = "bus_ticket"            // bus ticket
	TollInvoice         MultipleInvoiceType = "toll_invoice"          // road/bridge toll invoice
	FerryTicket         MultipleInvoiceType = "ferry_ticket"          // ferry ticket
	MotorVehicleInvoice MultipleInvoiceType = "motor_vehicle_invoice" // motor vehicle sales invoice
	UsedVehicleInvoice  MultipleInvoiceType = "used_vehicle_invoice"  // used vehicle sales invoice
	TaxiOnlineTicket    MultipleInvoiceType = "taxi_online_ticket"    // ride-hailing itinerary
	LimitInvoice        MultipleInvoiceType = "limit_invoice"         // capped invoice
	ShoppingReceipt     MultipleInvoiceType = "shopping_receipt"      // shopping receipt
	PosInvoice          MultipleInvoiceType = "pos_invoice"           // POS receipt
	Others              MultipleInvoiceType = "others"                // others
)

func (b *baiduAi) postForm(url string, data url.Values, apiKey string) (*http.Response, error) {
	encoded := data.Encode()

again:
	req, err := http.NewRequest("POST", url, strings.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := b.client.Do(req)
	if err != nil {
		b.postFailCount += 1
		if b.postFailCount < postFailCountMax {
			log.Println("Baidu OCR request failed, retrying...")
			goto again
		}
	} else {
		b.postFailCount = 0
	}

	return resp, err

}

func (r MultipleInvoiceResponseResult) GetField(field string) string {
	if items, exists := r[field]; exists && len(items) > 0 {
		return items[0].Word
	}
	return ""
}
