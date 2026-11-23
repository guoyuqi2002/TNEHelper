# Offline reconciliation demo

Open index.html directly and load the synthetic sample. It contains a match, an amount mismatch, a missing invoice, a missing invoice number, and two claims sharing one invoice reference. Open a result to inspect its evidence, record a decision and export the CSV. Reviews are kept in the current page session.

The web entry bypasses the Workday login, tenant navigation and OCR stages. It runs the local demonstration from structured invoice and expense records. The desktop Go/Wails application remains the integration entry point.

## Import a local JSON batch

Use this structure. Amounts are decimal strings in CNY. Expense IDs must be unique and the invoice value may be an empty string when recognition did not supply a number.

    {
      "expenses": [
        { "id": "EXP-101", "invoice": "INV-1101", "description": "Airport transfer", "amount": "86.50" }
      ],
      "invoices": [
        { "invoice": "INV-1101", "amount": "86.50" }
      ]
    }

For static hosting, run npm run build and use the dist folder. There are no runtime dependencies or API keys.
