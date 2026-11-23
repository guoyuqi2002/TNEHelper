# Workday Expense & Invoice Assistant

Desktop automation suite that reconciles travel & expense reports against VAT
invoices — built during a tax-operations internship at a global
professional-services firm for a pharmaceutical-client ERP workflow.

The tool replaces a fully manual, error-prone reconciliation process: it logs
into a Workday tenant, batch-downloads every expense-report attachment, parses
or OCRs the invoices, matches them against a SQLite database of
tax-bureau-issued invoices, and exports ready-to-review Excel workbooks.

## Features

- **Workday browser automation** — drives a real Chrome session (chromedp) to
  navigate a Workday tenant, wait for user SSO login, and harvest
  expense-report attachments (both report-level and per-line) via injected JS.
- **Attachment pipeline** — MD5 deduplication of downloaded files, per-report
  progress tracking, and crash-safe state persisted with `encoding/gob`.
- **Invoice recognition** — local PDF/OFD parsing (page splitting, embedded
  image extraction) plus Baidu AI OCR for VAT invoices, train tickets, taxi
  receipts, air itineraries, toll invoices and more.
- **Invoice matching** — imports tax-bureau invoice exports into SQLite
  (GORM), then matches recognized invoices by invoice code/number/e-invoice
  number and verifies amounts exactly with `math/big.Rat`.
- **Excel generation** — typed, code-generated sheet definitions
  (`pkg/sheets`) produce matched/unmatched workbooks with per-type sheets.
- **Offline demo web app** — an interactive reconciliation and review workspace for
  synthetic examples and local JSON import (`demo-web/`).

## Architecture

```
download (chromedp)  ->  deduplicate (MD5)  ->  recognize (PDF/OFD + OCR)
                                                             |
   SQLite invoice DB (GORM)  <------------------------------+
                                        |
                        match (invoice no. + exact amounts)  ->  export (Excel)
```

| Layer | Tech |
| --- | --- |
| Desktop shell | Go + Wails v2 |
| Frontend | Vue 3, Pinia, Vite |
| Browser automation | chromedp |
| Storage | SQLite via GORM |
| OCR | Baidu AI OCR (multiple-invoice endpoint) |
| Invoice files | pdfcpu (PDF), custom OFD page reader |
| Excel | tealeg/xlsx + generated typed sheet definitions |

## Repository layout

```
app/        Wails v2 desktop application (Go backend + Vue 3 frontend)
demo-web/   Standalone offline browser demo from synthetic or local JSON records
```

## Build & run (app)

Prerequisites: Go 1.24+, Node.js, the [Wails CLI](https://wails.io) v2, and a
Chrome/Chromium installation.

```bash
cd app
go mod download        # vendor/ is intentionally not shipped; go.sum is included
cd frontend && npm install && cd ..
wails dev              # development
wails build            # production binary
```

### Environment variables

Copy `app/.env.example` and export the variables before running:

| Variable | Purpose |
| --- | --- |
| `BAIDU_OCR_API_KEY` | Baidu AI OCR API key (required for recognition) |
| `BAIDU_OCR_API_NAME` | Optional label for the API key |
| `WORKDAY_HOME_URL` | Workday tenant home URL for login |
| `WORKDAY_TASK_URL` | Workday task URL used to look up expense reports |
| `KNOWN_PURCHASER_KEYWORD` | Optional buyer-name keyword for classifying unmatched VAT invoices |

## Run the offline demo

Open demo-web/index.html directly in a browser and select **Load sample batch**. The demo starts from structured invoice and expense records. It can match by invoice number and exact amount, present the evidence for exceptions, save a review decision for the current session and export the results as CSV.

The demo requires no credentials or external services. The desktop integration above is the path for real tenant access and invoice recognition. See demo-web/README.md for the local JSON import format.
## Data & confidentiality

This repository contains **no client data**. Credentials, tenant URLs, client
logos, invoice samples and databases were removed, and all identifying
configuration was moved to environment variables (see `.env.example`). The
Workday integration expects your own tenant URLs; the invoice sheet
definitions use neutralized English column headers.
