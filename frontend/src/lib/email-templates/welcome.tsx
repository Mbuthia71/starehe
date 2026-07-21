import { EmailLayout, emailButton, emailLink } from "./Layout";

export interface WelcomeEmailProps {
  memberName: string;
  memberNumber: string;
  msisdn: string;
  activationUrl: string;
  logoUrl?: string;
}

export function WelcomeEmail({
  memberName,
  memberNumber,
  msisdn,
  activationUrl,
  logoUrl,
}: WelcomeEmailProps) {
  return (
    <EmailLayout
      headline="WELCOME!"
      preheader={`Your Siohioma membership is ready, ${memberName.split(" ")[0]}.`}
      logoUrl={logoUrl}
    >
      <p style={{ margin: "0 0 18px", fontWeight: 500, color: "#1a2620" }}>
        Hi {memberName.split(" ")[0]},
      </p>
      <p style={{ margin: "0 0 18px" }}>
        You've successfully joined Siohioma SACCO under member number{" "}
        <strong style={{ color: "#1a2620" }}>{memberNumber}</strong>. Your
        phone <strong>{msisdn}</strong> is now linked to your account.
      </p>
      <p style={{ margin: "0 0 28px" }}>
        Activate your account to set your PIN and start saving, borrowing, and
        moving money with M-Pesa.
      </p>
      <p style={{ margin: "0 0 28px", textAlign: "center" as const }}>
        <a href={activationUrl} style={emailButton}>
          Activate account
        </a>
      </p>
      <p style={{ margin: 0, fontSize: 13, color: "#8b948f" }}>
        If you didn't sign up, just ignore this email or write us at{" "}
        <a href="mailto:hello@siohioma.co.ke" style={emailLink}>
          hello@siohioma.co.ke
        </a>
        .
      </p>
    </EmailLayout>
  );
}
