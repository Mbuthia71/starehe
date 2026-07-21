import { EmailLayout, emailButton } from "./Layout";

export interface LoanStatusEmailProps {
  memberName: string;
  status: "submitted" | "approved" | "disbursed" | "rejected";
  productName: string;
  amount: string;
  nextStep?: string;
  ctaLabel?: string;
  ctaUrl?: string;
  logoUrl?: string;
}

const headlineByStatus: Record<LoanStatusEmailProps["status"], string> = {
  submitted: "RECEIVED",
  approved: "APPROVED",
  disbursed: "DISBURSED",
  rejected: "ON HOLD",
};

const eyebrowByStatus: Record<LoanStatusEmailProps["status"], string> = {
  submitted: "Loan application received",
  approved: "Loan approved",
  disbursed: "Funds on the way",
  rejected: "Loan application update",
};

export function LoanStatusEmail({
  memberName,
  status,
  productName,
  amount,
  nextStep,
  ctaLabel,
  ctaUrl,
  logoUrl,
}: LoanStatusEmailProps) {
  return (
    <EmailLayout
      headline={headlineByStatus[status]}
      eyebrow={eyebrowByStatus[status]}
      preheader={`${productName} · ${amount}`}
      logoUrl={logoUrl}
    >
      <p style={{ margin: "0 0 18px", fontWeight: 500, color: "#1a2620" }}>
        Hi {memberName.split(" ")[0]},
      </p>
      <p style={{ margin: "0 0 24px" }}>
        Your <strong>{productName}</strong> application for{" "}
        <strong>{amount}</strong> is{" "}
        <span style={{ color: "#1f3a2e", fontWeight: 600 }}>{status}</span>.
      </p>
      {nextStep && (
        <p style={{ margin: "0 0 28px" }}>{nextStep}</p>
      )}
      {ctaLabel && ctaUrl && (
        <p style={{ margin: "0 0 28px", textAlign: "center" as const }}>
          <a href={ctaUrl} style={emailButton}>
            {ctaLabel}
          </a>
        </p>
      )}
      <p style={{ margin: 0, fontSize: 13, color: "#8b948f" }}>
        Questions? Reply to this email and a loans officer will get back to
        you.
      </p>
    </EmailLayout>
  );
}
