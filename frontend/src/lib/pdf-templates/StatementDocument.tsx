import { statementCss } from "./statement-styles";

export interface StatementMember {
  fullName: string;
  memberNumber: string;
  msisdn: string;
  address: string[]; // multi-line
}

export interface StatementAccount {
  productName: string; // e.g. "Siohioma Savings Account"
  accountNumber: string; // e.g. "0148-2210-09"
  subtitle?: string; // e.g. "MEMBER: ACHIENG OTIENO"
}

export interface StatementSummary {
  beginningBalance: number;
  deposits: number;
  mpesaIn: number;
  otherCredits: number;
  withdrawals: number;
  loanRepayments: number;
  endingBalance: number;
  averageBalance: number;
  interestEarnedPeriod: number;
  interestYtd: number;
  dividendsCredited: number;
  apy: number; // e.g. 0.0625
  daysInPeriod: number;
}

export interface StatementTxn {
  date: string; // "07/02"
  description: string;
  amount: number;
}

export interface StatementSection {
  title: string; // "Deposits" | "M-Pesa deposits" | "Withdrawals" | "Loan activity" | "Dividends & interest"
  rows: StatementTxn[];
}

export interface StatementDocumentProps {
  member: StatementMember;
  account: StatementAccount;
  period: { from: string; to: string; label: string };
  summary: StatementSummary;
  sections: StatementSection[];
  page?: { current: number; total: number };
  logoUrl?: string;
  /** When true wraps in a screen preview shell with grey background */
  preview?: boolean;
}

const fmt = (n: number) =>
  n.toLocaleString("en-KE", { minimumFractionDigits: 2, maximumFractionDigits: 2 });

const fmtPct = (n: number) => `${(n * 100).toFixed(2)}%`;

export function StatementDocument(props: StatementDocumentProps) {
  const {
    member,
    account,
    period,
    summary,
    sections,
    page = { current: 1, total: 1 },
    logoUrl = "https://siohioma.co.ke/logo.png",
    preview = false,
  } = props;

  const Page = (
    <div className="sio-page">
      {/* HEADER */}
      <div className="sio-header">
        <div className="sio-brand">
          <img src={logoUrl} alt="Siohioma" />
          <div className="sio-brand-text">
            <div className="sio-brand-name">Siohioma</div>
            <div className="sio-brand-tag">Member-owned · Member-run</div>
          </div>
        </div>
        <div>
          <div className="sio-doc-type">Statement of Account</div>
          <div className="sio-meta" style={{ gridTemplateColumns: "auto auto", marginTop: "3mm" }}>
            <div className="sio-meta-block">
              <div className="sio-meta-row">
                <span className="label">Page:</span>
                <span className="value">{page.current} of {page.total}</span>
              </div>
              <div className="sio-meta-row">
                <span className="label">Period:</span>
                <span className="value">{period.from} – {period.to}</span>
              </div>
              <div className="sio-meta-row">
                <span className="label">Member #:</span>
                <span className="value">{member.memberNumber}</span>
              </div>
              <div className="sio-meta-row">
                <span className="label">Primary A/C:</span>
                <span className="value">{account.accountNumber}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* ADDRESS */}
      <div className="sio-addr">
        <div style={{ fontWeight: 600, textTransform: "uppercase", letterSpacing: "0.04em" }}>
          {member.fullName}
        </div>
        {member.address.map((l, i) => (
          <div key={i}>{l}</div>
        ))}
        <div>{member.msisdn}</div>
      </div>

      {/* SECTION TITLE */}
      <div className="sio-section-acct">Account # {account.accountNumber}</div>
      <h2 className="sio-section-title">{account.productName}</h2>
      <div className="sio-section-sub">
        {account.subtitle ?? `MEMBER: ${member.fullName.toUpperCase()}`}
      </div>

      {/* ACCOUNT SUMMARY */}
      <div className="sio-band">Account summary</div>
      <div className="sio-summary">
        <div>
          <div className="row"><span className="k">Beginning balance</span><span className="v">{fmt(summary.beginningBalance)}</span></div>
          <div className="row"><span className="k">Deposits</span><span className="v">{fmt(summary.deposits)}</span></div>
          <div className="row"><span className="k">M-Pesa deposits</span><span className="v">{fmt(summary.mpesaIn)}</span></div>
          <div className="row"><span className="k">Other credits</span><span className="v">{fmt(summary.otherCredits)}</span></div>
          <div className="row"><span className="k">Withdrawals</span><span className="v">({fmt(summary.withdrawals)})</span></div>
          <div className="row"><span className="k">Loan repayments</span><span className="v">({fmt(summary.loanRepayments)})</span></div>
          <div className="row total"><span className="k">Ending balance</span><span className="v">{fmt(summary.endingBalance)}</span></div>
        </div>
        <div>
          <div className="row"><span className="k">Average collected balance</span><span className="v">{fmt(summary.averageBalance)}</span></div>
          <div className="row"><span className="k">Interest earned this period</span><span className="v">{fmt(summary.interestEarnedPeriod)}</span></div>
          <div className="row"><span className="k">Interest paid year-to-date</span><span className="v">{fmt(summary.interestYtd)}</span></div>
          <div className="row"><span className="k">Dividends credited</span><span className="v">{fmt(summary.dividendsCredited)}</span></div>
          <div className="row"><span className="k">Annual percentage yield</span><span className="v">{fmtPct(summary.apy)}</span></div>
          <div className="row"><span className="k">Days in period</span><span className="v">{summary.daysInPeriod}</span></div>
        </div>
      </div>

      {/* DAILY ACTIVITY */}
      <div className="sio-band" style={{ marginTop: "8mm" }}>Daily account activity</div>

      {sections.map((s) => {
        const subtotal = s.rows.reduce((a, r) => a + r.amount, 0);
        return (
          <div key={s.title}>
            <h3 className="sio-sub-title">{s.title}</h3>
            <table className="sio-table">
              <thead>
                <tr>
                  <th>Posting date</th>
                  <th>Description</th>
                  <th className="amount">Amount (KES)</th>
                </tr>
              </thead>
              <tbody>
                {s.rows.map((r, i) => (
                  <tr key={i}>
                    <td className="date">{r.date}</td>
                    <td className="desc">{r.description}</td>
                    <td className="amount">{fmt(r.amount)}</td>
                  </tr>
                ))}
                <tr className="subtotal-row">
                  <td />
                  <td style={{ textAlign: "right" }}>Subtotal</td>
                  <td className="amount">{fmt(subtotal)}</td>
                </tr>
              </tbody>
            </table>
          </div>
        );
      })}

      {/* FOOTER */}
      <div className="sio-footer">
        <div className="call">
          Call <strong>+254 711 000 000</strong> for 24-hour member support or visit{" "}
          <strong>siohioma.co.ke</strong>
        </div>
        <div className="reg">
          Member deposits protected · Regulated by SASRA · License No. CS/1820
        </div>
      </div>
    </div>
  );

  return (
    <div className={preview ? "sio-preview-shell" : "sio-doc"}>
      <style dangerouslySetInnerHTML={{ __html: statementCss }} />
      <div className={preview ? "sio-doc" : ""}>{Page}</div>
    </div>
  );
}
