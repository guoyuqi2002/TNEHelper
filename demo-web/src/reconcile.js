'use strict';

const sampleBatch = {
  expenses: [
    { id: 'EXP-101', invoice: 'INV-1101', description: 'Airport transfer', amount: '86.50' },
    { id: 'EXP-102', invoice: 'INV-1102', description: 'Client site rail fare', amount: '240.00' },
    { id: 'EXP-103', invoice: 'INV-1103', description: 'Hotel stay', amount: '520.00' },
    { id: 'EXP-104', invoice: 'INV-1104', description: 'Meeting supplies', amount: '112.00' },
    { id: 'EXP-105', invoice: '', description: 'Taxi receipt with an unreadable number', amount: '42.00' },
    { id: 'EXP-106', invoice: 'INV-1104', description: 'Second claim for meeting supplies', amount: '112.00' },
  ],
  invoices: [
    { invoice: 'INV-1101', amount: '86.50' },
    { invoice: 'INV-1102', amount: '225.00' },
    { invoice: 'INV-1104', amount: '112.00' },
  ],
};

const reviewGuidance = {
  matched: { label: 'Matched', next: 'The invoice number and exact amount match.' },
  amount: { label: 'Amount mismatch', next: 'Compare the claimed amount with the invoice and check for a partial claim or data entry error.' },
  missing: { label: 'Invoice missing', next: 'Check the invoice number or obtain the missing invoice record.' },
  incomplete: { label: 'Number missing', next: 'Review the receipt and record the invoice number if it can be verified.' },
  duplicate: { label: 'Duplicate reference', next: 'Compare all claims and invoice records for this number before approving any claim.' },
};

function cents(value) {
  if (typeof value !== 'string' || !/^\d{1,12}(?:\.\d{1,2})?$/.test(value)) {
    throw new Error('Enter amounts as decimal strings with up to two places, such as "86.50".');
  }
  const [units, fraction = ''] = value.split('.');
  return Number(units) * 100 + Number(fraction.padEnd(2, '0'));
}

function field(value, name, optional = false) {
  if (typeof value !== 'string' || (!optional && value.trim() === '')) {
    throw new Error(name + ' must be a text field.');
  }
  return value.trim();
}

function reconcile(batch) {
  if (!batch || !Array.isArray(batch.expenses) || !Array.isArray(batch.invoices) || !batch.expenses.length) {
    throw new Error('The JSON file must contain a nonempty expenses array and an invoices array.');
  }
  const byInvoice = new Map();
  for (const row of batch.invoices) {
    const invoice = field(row?.invoice, 'Invoice number').toUpperCase();
    const amounts = byInvoice.get(invoice) ?? [];
    amounts.push(cents(row.amount));
    byInvoice.set(invoice, amounts);
  }

  const ids = new Set();
  const claimCount = new Map();
  const claims = batch.expenses.map(row => {
    const id = field(row?.id, 'Expense ID');
    if (ids.has(id)) throw new Error('Expense IDs must be unique.');
    ids.add(id);
    const invoice = field(row.invoice, 'Invoice number', true).toUpperCase();
    if (invoice) claimCount.set(invoice, (claimCount.get(invoice) ?? 0) + 1);
    return { id, invoice, description: field(row.description, 'Description'), claimed: cents(row.amount) };
  });

  return claims.map(claim => {
    const candidates = byInvoice.get(claim.invoice) ?? [];
    let status;
    if (!claim.invoice) status = 'incomplete';
    else if (candidates.length > 1 || claimCount.get(claim.invoice) > 1) status = 'duplicate';
    else if (candidates.length === 0) status = 'missing';
    else status = candidates[0] === claim.claimed ? 'matched' : 'amount';
    return { ...claim, candidates, status };
  });
}

function toCsv(results, reviews) {
  const fields = ['Expense ID', 'Description', 'Invoice', 'Claimed CNY', 'Invoice amounts CNY', 'Match result', 'Review decision', 'Review note'];
  const format = value => (value / 100).toFixed(2);
  const rows = results.map(result => {
    const review = reviews.get(result.id);
    return [result.id, result.description, result.invoice, format(result.claimed),
      result.candidates.map(format).join(' / '), reviewGuidance[result.status].label,
      review?.decision ?? '', review?.note ?? ''];
  });
  const cell = value => {
    const safe = /^[\s]*[=+\-@]|^[\t\r\n]/.test(value) ? "'" + value : value;
    return '"' + safe.replaceAll('"', '""') + '"';
  };
  return [fields, ...rows].map(row => row.map(cell).join(',')).join('\r\n') + '\r\n';
}

if (typeof module !== 'undefined') {
  module.exports = { sampleBatch, reviewGuidance, reconcile, toCsv };
}
