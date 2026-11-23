'use strict';

const fs = require('node:fs');
const path = require('node:path');

const out = path.join(__dirname, 'dist');
fs.mkdirSync(path.join(out, 'src'), { recursive: true });
for (const file of ['index.html', 'src/reconcile.js', 'src/main.js', 'src/style.css']) {
  fs.copyFileSync(path.join(__dirname, file), path.join(out, file));
}
console.log('Built ' + path.join(out, 'index.html'));
