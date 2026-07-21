import { EmailLayout, emailButton } from "./Layout";

export interface StatementReadyEmailProps {
  memberName: string;
  period: string; // e.g. "October 2026"
  accountLabel: string;
  downloadUrl: string;
  logoUrl?: string;
}

export function StatementReadyEmail({
  memberName,
  period,
  accountLabel,
  downloadUrl,
  logoUrl,
}: StatementReadyEmailProps) {
  return (
    <EmailLayout
      headline="STATEMENT"
      eyebrow={`${period} · ready to download`}
      preheader={`Your ${period} statement is ready.`}
      logoUrl={logoUrl}
    >
      <p style={{ margin: "0 0 18px", fontWeight: 500, color: "#1a2620" }}>
        Hi {memberName.split(" ")[0]},
      </p>
      <p style={{ margin: "0 0 28px" }}>
        Your <strong>{period}</strong> statement for{" "}
        <strong>{accountLabel}</strong> is ready. The PDF link below is valid
        for 14 days.
      </p>
      <p style={{ margin: "0 0 28px", textAlign: "center" as const }}>
        <a href={downloadUrl} style={emailButton}>
          Download statement
        </a>
      </p>
      <p style={{ margin: 0, fontSize: 13, color: "#8b948f" }}>
        Need an older statement? Visit the Statements tab in the Siohioma app
        and request any month from the past 7 years.
      </p>
    </EmailLayout>
  );
}
