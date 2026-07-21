## What's in the system today

Searched the codebase — **no PDF generator and no email sending exist yet**. The Statements route (`src/routes/_app.statements.tsx`) has a "Generate PDF" button with no implementation, and `auth.ts` has an `emailAddress` field but zero send logic. So both features need to be created from scratch.

You said design-only and Windsurf will wire the backend. So I'll build presentation-layer templates that render exactly what the recipient will see, with clear data-prop hooks the backend can populate.

---

## 1. Email templates — Kid & Coe inspired, Siohioma branded

Visual language from image 1, translated to Siohioma:
- Centered logo at top (use existing `src/assets/siohioma-logo.png`)
- Massive ultra-thin display headline ("WELCOME!", "RESET PASSWORD", "CONFIRM", "RECEIPT"), letter-spaced, light weight
- Short gold/amber underline accent under the headline (matches the yellow stroke in the reference)
- Tight body copy in a clean grotesk, generous white space
- Sign-off in a flowing script font ("The Siohioma Team")
- Hairline divider + small circular icons at the bottom
- Footer with address block, centered, muted

### Files

Create `src/lib/email-templates/` with one `.tsx` per template, plus a shared `Layout.tsx`:

- `Layout.tsx` — header (logo), headline + gold rule, body slot, signoff, hairline, social icons, footer address. Inline styles only (email-safe), Google Font links for `Cormorant Garamond` / `Jost` (display thin) and `Dancing Script` (signoff).
- `welcome.tsx` — "WELCOME!" — new member, with activation link
- `password-reset.tsx` — "RESET PASSWORD" — secure link + expiry note
- `password-changed.tsx` — "PASSWORD UPDATED" — confirmation + "wasn't you?" line
- `otp-code.tsx` — "YOUR CODE" — big spaced 6-digit OTP block
- `login-alert.tsx` — "NEW SIGN-IN" — device, location, time
- `transaction-receipt.tsx` — "RECEIPT" — amount, ref, counterparty, balance after
- `loan-status.tsx` — "LOAN UPDATE" — status pill, amount, next step
- `statement-ready.tsx` — "STATEMENT READY" — period + download CTA

Each template is a pure React component with typed props (e.g. `{ memberName, resetUrl, expiresInMinutes }`) so Windsurf only has to import and pass data.

### Preview route

`src/routes/_app.dev.emails.tsx` — internal-only page that renders each template in an iframe with sample data, so you can review them in the preview without sending anything. Stripped from production by living under an obvious dev path.

---

## 2. PDF statement template — TD Bank inspired, Siohioma branded

Visual language from image 2, translated to a SACCO member statement:
- Top-left: Siohioma logo + tagline ("Member-owned. Member-run.")
- Top-right: "STATEMENT OF ACCOUNT" eyebrow, page X of Y, statement period, member ref, primary account #
- Member address block (left), account meta (right)
- Bold section title: "Siohioma Savings Account" + "MEMBER: <Name>" + "Account # ……"
- **ACCOUNT SUMMARY** band (green hairline, uppercase eyebrow) — two columns: balances vs interest/period stats
- **DAILY ACCOUNT ACTIVITY** band with subsections:
  - Deposits (member contributions, M-Pesa in)
  - Withdrawals (M-Pesa out, transfers)
  - Loan activity (disbursements, repayments)
  - Dividends / interest credited
- Each row: posting date · description · amount, right-aligned subtotals
- Footer bar: "Call +254 … for 24-hour member support or visit siohioma.co.ke" + regulatory line ("Regulated by SASRA")

### Implementation (design-only, browser-rendered)

Build it as an HTML/CSS template that mirrors the PDF layout 1:1, so Windsurf can either:
(a) print it to PDF via headless Chromium server-side, or
(b) hand the same component to `@react-pdf/renderer` later.

Files:
- `src/lib/pdf-templates/StatementDocument.tsx` — the full statement, A4-sized, paginated by section, typed props (`member`, `account`, `period`, `summary`, `transactions[]`).
- `src/lib/pdf-templates/statement-styles.ts` — print-ready CSS (A4 page size, mm units, Siohioma green replacing TD green, monospaced numerics via `.num`).
- `src/routes/_app.dev.statement.tsx` — preview route that renders `StatementDocument` with realistic mock data (12 months of mixed deposits, M-Pesa, loan repayments, dividends) so you can see and screenshot it.

The existing `Generate PDF` button in `_app.statements.tsx` will be re-pointed to navigate to/print this template. Actual PDF byte generation is left for Windsurf — the template is the contract.

---

## Design tokens (shared)

- Brand green: reuse the Siohioma green already in `src/styles.css` (replacing TD's `#00853f`)
- Accent gold: existing amber token (matches the yellow underline in image 1)
- Display: `Cormorant Garamond` light (300) for big headlines
- Body: `Jost` / `Inter` 400 for paragraphs
- Signoff script: `Dancing Script`
- Numerics: tabular-nums everywhere in the PDF
- All currency via the existing `formatKES` in `src/lib/format.ts`

---

## Out of scope (Windsurf will handle)

- SMTP / Resend / SES wiring
- Server route that renders template → HTML → email body
- Headless Chromium / `@react-pdf/renderer` install + PDF byte stream
- Storing generated PDFs in object storage
- Triggering emails from auth events

Confirm and I'll build it.