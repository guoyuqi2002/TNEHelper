'use strict';

const element = id => document.getElementById(id);
const money = new Intl.NumberFormat('en', { style: 'currency', currency: 'CNY', currencyDisplay: 'code' });
let results = [];
let reviews = new Map();
let selected;
let importVersion = 0;

function formatMoney(value) {
  return money.format(value / 100);
}

function displayResults() {
  const matched = results.filter(result => result.status === 'matched').length;
  element('batch-summary').textContent =
    results.length + ' claims · ' + matched + ' matched · ' + (results.length - matched) +
    ' exceptions · ' + reviews.size + ' reviews saved';
  const rows = results.map(result => {
    const row = document.createElement('tr');
    const review = reviews.get(result.id);
    const values = [result.id + ' · ' + result.description, result.invoice || '—', formatMoney(result.claimed), reviewGuidance[result.status].label];
    for (const value of values) {
      const cell = row.insertCell();
      cell.textContent = value;
    }
    const action = document.createElement('button');
    action.type = 'button';
    action.textContent = review ? 'Edit review' : 'Open details';
    action.setAttribute('aria-label', action.textContent + ' for ' + result.id);
    action.addEventListener('click', () => showReview(result));
    row.insertCell().append(action);
    if (review) {
      const state = document.createElement('span');
      state.className = 'review-state';
      state.textContent = review.decision;
      row.lastElementChild.append(state);
    }
    return row;
  });
  element('result-rows').replaceChildren(...rows);
}

function showReview(result) {
  selected = result;
  const review = reviews.get(result.id);
  const guide = reviewGuidance[result.status];
  element('review-title').textContent = result.id + ' · ' + guide.label;
  const evidence = [
    ['Invoice reference', result.invoice || 'Missing'],
    ['Claimed amount', formatMoney(result.claimed)],
    ['Invoice records', result.candidates.length ? result.candidates.map(formatMoney).join(' / ') : 'No matching record'],
  ];
  element('evidence').replaceChildren(...evidence.flatMap(([label, value]) => {
    const name = document.createElement('dt');
    const description = document.createElement('dd');
    name.textContent = label;
    description.textContent = value;
    return [name, description];
  }));
  element('guidance').textContent = guide.next;
  element('decision').value = review?.decision ?? '';
  element('review-note').value = review?.note ?? '';
  element('review-note').setCustomValidity('');
  element('review-status').textContent = '';
  element('review').hidden = false;
  element('review-title').focus();
}

function loadBatch(batch, label) {
  const next = reconcile(batch);
  results = next;
  reviews = new Map();
  selected = undefined;
  element('review').hidden = true;
  element('results').hidden = false;
  element('export').disabled = false;
  element('import-error').hidden = true;
  element('status').textContent = label + ' loaded. Export a copy before replacing this batch or closing the page.';
  displayResults();
}

element('load-sample').addEventListener('click', () => {
  importVersion++;
  loadBatch(sampleBatch, 'Synthetic sample');
});
element('batch-file').addEventListener('change', async event => {
  const input = event.currentTarget;
  const file = input.files[0];
  if (!file) return;
  const version = ++importVersion;
  input.value = '';
  try {
    const batch = JSON.parse(await file.text());
    if (version === importVersion) loadBatch(batch, file.name);
  } catch (error) {
    if (version === importVersion) {
      element('import-error').textContent = error instanceof SyntaxError ? 'The selected file is not valid JSON.' : error.message;
      element('import-error').hidden = false;
    }
  }
});
element('review-note').addEventListener('input', event => event.target.setCustomValidity(''));
element('review-form').addEventListener('submit', event => {
  event.preventDefault();
  if (!selected) return;
  const noteField = element('review-note');
  const note = noteField.value.trim();
  if (!note) {
    noteField.setCustomValidity('Enter a reason for this decision.');
    noteField.reportValidity();
    return;
  }
  const decision = element('decision').value;
  reviews.set(selected.id, { decision, note });
  displayResults();
  element('review-status').textContent = 'Review saved for ' + selected.id + '.';
});
element('export').addEventListener('click', () => {
  const blob = new Blob(['\ufeff', toCsv(results, reviews)], { type: 'text/csv;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  const download = document.createElement('a');
  download.href = url;
  download.download = 'invoice-reconciliation.csv';
  document.body.append(download);
  download.click();
  download.remove();
  setTimeout(() => URL.revokeObjectURL(url), 0);
});
