/**
 * Print-ready styles for the Siohioma statement PDF.
 *
 * Mounted as a single <style> tag inside <StatementDocument /> so the
 * markup is self-contained — Windsurf can pipe the same HTML to headless
 * Chromium for PDF generation without touching Tailwind.
 */
export const statementCss = `
@page {
  size: A4;
  margin: 14mm 12mm 16mm 12mm;
}

.sio-doc {
  font-family: "Inter", "Helvetica Neue", Arial, sans-serif;
  color: #1a2620;
  font-size: 9.5pt;
  line-height: 1.4;
  background: #ffffff;
}

.sio-doc * { box-sizing: border-box; }

.sio-page {
  width: 186mm;
  min-height: 268mm;
  margin: 0 auto;
  padding: 0;
  page-break-after: always;
  position: relative;
}
.sio-page:last-child { page-break-after: auto; }

/* ── HEADER ──────────────────────────────────────────────── */
.sio-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding-bottom: 10mm;
}
.sio-brand {
  display: flex;
  align-items: center;
  gap: 4mm;
}
.sio-brand img {
  width: 16mm;
  height: 16mm;
  object-fit: contain;
}
.sio-brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.1;
}
.sio-brand-name {
  font-family: "Cormorant Garamond", "Forum", Georgia, serif;
  font-weight: 500;
  font-size: 20pt;
  color: #1f3a2e;
  letter-spacing: 0.01em;
}
.sio-brand-tag {
  font-size: 7.5pt;
  color: #5e6b65;
  margin-top: 1mm;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}
.sio-doc-type {
  text-align: right;
  font-size: 8pt;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: #1f3a2e;
  font-weight: 600;
}

.sio-meta {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8mm;
  margin-top: 4mm;
  font-size: 8.5pt;
}
.sio-meta .label {
  color: #5e6b65;
}
.sio-meta .value {
  color: #1a2620;
  font-weight: 500;
  font-variant-numeric: tabular-nums;
}
.sio-meta-row {
  display: flex;
  justify-content: space-between;
  gap: 6mm;
  padding: 1mm 0;
}
.sio-meta-block {
  display: flex;
  flex-direction: column;
}
.sio-addr {
  font-size: 8.5pt;
  line-height: 1.5;
  color: #1a2620;
}

/* ── BAND TITLES ─────────────────────────────────────────── */
.sio-band {
  margin-top: 6mm;
  padding-bottom: 1.5mm;
  border-bottom: 0.5pt solid #1f3a2e;
  font-size: 8pt;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #1f3a2e;
}

.sio-section-title {
  margin-top: 5mm;
  font-family: "Cormorant Garamond", "Forum", Georgia, serif;
  font-size: 14pt;
  font-weight: 500;
  color: #1a2620;
}
.sio-section-sub {
  font-size: 8.5pt;
  color: #5e6b65;
  margin-top: 0.5mm;
}
.sio-section-acct {
  float: right;
  font-size: 9pt;
  font-weight: 600;
  color: #1a2620;
  margin-top: -8mm;
  font-variant-numeric: tabular-nums;
  letter-spacing: 0.02em;
}

/* ── SUMMARY ─────────────────────────────────────────────── */
.sio-summary {
  margin-top: 3mm;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8mm;
  font-size: 9pt;
}
.sio-summary .row {
  display: flex;
  justify-content: space-between;
  padding: 1.4mm 0;
}
.sio-summary .row .k { color: #1a2620; }
.sio-summary .row .v {
  font-variant-numeric: tabular-nums;
  font-weight: 500;
}
.sio-summary .row.total {
  border-top: 0.5pt solid #c8c2b0;
  padding-top: 2.2mm;
  margin-top: 1mm;
  font-weight: 600;
}

/* ── TABLES ──────────────────────────────────────────────── */
.sio-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 2mm;
  font-size: 8.5pt;
}
.sio-table thead th {
  text-align: left;
  font-weight: 700;
  font-size: 7.5pt;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #5e6b65;
  padding: 1.5mm 1mm;
  border-bottom: 0.4pt solid #d8d2c0;
}
.sio-table thead th.amount,
.sio-table tbody td.amount {
  text-align: right;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}
.sio-table tbody td {
  padding: 1.5mm 1mm;
  border-bottom: 0.25pt solid #ece8da;
  color: #1a2620;
  vertical-align: top;
}
.sio-table tbody td.date {
  width: 18mm;
  color: #5e6b65;
  font-variant-numeric: tabular-nums;
}
.sio-table tbody td.desc {
  font-weight: 500;
}
.sio-table .subtotal-row td {
  padding-top: 2mm;
  border-bottom: none;
  border-top: 0.4pt solid #d8d2c0;
  font-weight: 700;
  color: #1f3a2e;
  font-size: 8.5pt;
}

.sio-sub-title {
  margin-top: 5mm;
  margin-bottom: 0;
  font-size: 10pt;
  font-weight: 700;
  color: #1a2620;
}

/* ── FOOTER ──────────────────────────────────────────────── */
.sio-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  border-top: 1pt solid #1f3a2e;
  padding-top: 2mm;
  text-align: center;
}
.sio-footer .call {
  font-family: "Cormorant Garamond", "Forum", Georgia, serif;
  font-size: 11pt;
  color: #1f3a2e;
  letter-spacing: 0.02em;
}
.sio-footer .reg {
  margin-top: 1mm;
  font-size: 6.5pt;
  letter-spacing: 0.06em;
  color: #5e6b65;
  text-transform: uppercase;
}

/* Screen preview only */
.sio-preview-shell {
  background: #e9e6dc;
  padding: 24px 12px;
  min-height: 100vh;
}
.sio-preview-shell .sio-page {
  background: #ffffff;
  box-shadow: 0 8px 30px rgba(0,0,0,0.10);
  padding: 14mm 12mm 16mm 12mm;
  margin-bottom: 16px;
  width: 210mm;
  min-height: 297mm;
}
@media print {
  .sio-preview-shell { background: #ffffff; padding: 0; }
  .sio-preview-shell .sio-page { box-shadow: none; margin: 0; padding: 0; width: 186mm; min-height: 268mm; }
}
`;
