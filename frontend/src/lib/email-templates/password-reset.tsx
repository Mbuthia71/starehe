import { EmailLayout, emailButton, emailLink } from "./Layout";

export interface PasswordResetEmailProps {
  memberName: string;
  resetUrl: string;
  expiresInMinutes: number;
  requestedFromIp?: string;
  logoUrl?: string;
}

export function PasswordResetEmail({
  memberName,
  resetUrl,
  expiresInMinutes,
  requestedFromIp,
  logoUrl,
}: PasswordResetEmailProps) {
  return (
    <EmailLayout
      headline="RESET PIN"
      eyebrow="Security request"
      preheader="A request to reset your Siohioma PIN was received."
      logoUrl={logoUrl}
    >
      <p style={{ margin: "0 0 18px", fontWeight: 500, color: "#1a2620" }}>
        Hi {memberName.split(" ")[0]},
      </p>
      <p style={{ margin: "0 0 18px" }}>
        Someone — hopefully you — asked to reset the PIN for your Siohioma
        account. Tap the button below within{" "}
        <strong>{expiresInMinutes} minutes</strong> to choose a new one.
      </p>
      <p style={{ margin: "0 0 28px", textAlign: "center" as const }}>
        <a href={resetUrl} style={emailButton}>
          Reset my PIN
        </a>
      </p>
      <p style={{ margin: "0 0 18px", fontSize: 13, color: "#8b948f" }}>
        Or paste this link into your browser:
        <br />
        <a href={resetUrl} style={{ ...emailLink, wordBreak: "break-all" }}>
          {resetUrl}
        </a>
      </p>
      {requestedFromIp && (
        <p style={{ margin: 0, fontSize: 12, color: "#8b948f" }}>
          Request came from IP <code>{requestedFromIp}</code>. If that wasn't
          you, ignore this email and your PIN stays the same.
        </p>
      )}
    </EmailLayout>
  );
}
