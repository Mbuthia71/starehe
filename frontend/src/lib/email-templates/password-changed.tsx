import { EmailLayout, emailLink } from "./Layout";

export interface PasswordChangedEmailProps {
  memberName: string;
  changedAt: string;
  deviceLabel?: string;
  logoUrl?: string;
}

export function PasswordChangedEmail({
  memberName,
  changedAt,
  deviceLabel,
  logoUrl,
}: PasswordChangedEmailProps) {
  return (
    <EmailLayout
      headline="PIN UPDATED"
      eyebrow="Account security"
      preheader="Your Siohioma PIN was just changed."
      logoUrl={logoUrl}
    >
      <p style={{ margin: "0 0 18px", fontWeight: 500, color: "#1a2620" }}>
        Hi {memberName.split(" ")[0]},
      </p>
      <p style={{ margin: "0 0 18px" }}>
        Your Siohioma PIN was successfully updated on{" "}
        <strong>{changedAt}</strong>
        {deviceLabel ? (
          <>
            {" "}
            from <strong>{deviceLabel}</strong>
          </>
        ) : null}
        .
      </p>
      <p style={{ margin: 0, fontSize: 13, color: "#8b948f" }}>
        Didn't do this? Call our 24-hour line{" "}
        <strong style={{ color: "#1f3a2e" }}>+254 711 000 000</strong> or email{" "}
        <a href="mailto:security@siohioma.co.ke" style={emailLink}>
          security@siohioma.co.ke
        </a>{" "}
        immediately. We'll freeze the account while we sort it out.
      </p>
    </EmailLayout>
  );
}
