import { EmailLayout, emailFonts } from "./Layout";

export interface TransactionReceiptEmailProps {
  memberName: string;
  type: "deposit" | "withdrawal" | "transfer" | "loan_repayment";
  amount: string; // pre-formatted, e.g. "KES 12,500"
  counterparty?: string;
  reference: string;
  balanceAfter: string;
  when: string;
  logoUrl?: string;
}

const typeCopy: Record<TransactionReceiptEmailProps["type"], string> = {
  deposit: "Deposit received",
  withdrawal: "Withdrawal to M-Pesa",
  transfer: "Transfer sent",
  loan_repayment: "Loan repayment",
};

export function TransactionReceiptEmail({
  memberName,
  type,
  amount,
  counterparty,
  reference,
  balanceAfter,
  when,
  logoUrl,
}: TransactionReceiptEmailProps) {
  return (
    <EmailLayout
      headline="RECEIPT"
      eyebrow={typeCopy[type]}
      preheader={`${typeCopy[type]} · ${amount}`}
      logoUrl={logoUrl}
    >
      <p style={{ margin: "0 0 18px", fontWeight: 500, color: "#1a2620" }}>
        Hi {memberName.split(" ")[0]},
      </p>
      <p style={{ margin: "0 0 28px" }}>
        Your {typeCopy[type].toLowerCase()} has cleared. Here are the details.
      </p>

      <div
        style={{
          textAlign: "center",
          padding: "28px 12px",
          marginBottom: 24,
          backgroundColor: "#faf7ef",
          borderRadius: 16,
        }}
      >
        <div
          style={{
            fontSize: 11,
            letterSpacing: "0.18em",
            textTransform: "uppercase",
            color: "#8b948f",
            marginBottom: 8,
          }}
        >
          Amount
        </div>
        <div
          style={{
            fontFamily: emailFonts.fontDisplay,
            fontWeight: 400,
            fontSize: 42,
            color: "#1f3a2e",
            lineHeight: 1,
          }}
        >
          {amount}
        </div>
      </div>

      <table
        role="presentation"
        cellPadding={0}
        cellSpacing={0}
        width="100%"
        style={{ marginBottom: 28, borderTop: "1px solid #e6e2d8" }}
      >
        <tbody>
          {([
            counterparty ? ["To / From", counterparty] : null,
            ["Reference", reference],
            ["Date", when],
            ["Balance after", balanceAfter],
          ].filter(Boolean) as [string, string][]).map(([k, v]) => (
            <tr key={k}>
              <td
                style={{
                  padding: "12px 0",
                  borderBottom: "1px solid #e6e2d8",
                  fontSize: 12,
                  textTransform: "uppercase",
                  letterSpacing: "0.14em",
                  color: "#8b948f",
                  width: 140,
                }}
              >
                {k}
              </td>
              <td
                style={{
                  padding: "12px 0",
                  borderBottom: "1px solid #e6e2d8",
                  fontSize: 14,
                  color: "#1a2620",
                  fontWeight: 500,
                  textAlign: "right" as const,
                  fontVariantNumeric: "tabular-nums",
                }}
              >
                {v}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <p style={{ margin: 0, fontSize: 13, color: "#8b948f" }}>
        Keep this email as proof of the transaction.
      </p>
    </EmailLayout>
  );
}
