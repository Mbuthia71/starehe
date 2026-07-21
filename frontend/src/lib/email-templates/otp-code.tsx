import { EmailLayout, emailFonts } from "./Layout";

export interface OtpCodeEmailProps {
  memberName: string;
  code: string;
  expiresInMinutes: number;
  purpose?: string;
  logoUrl?: string;
}

export function OtpCodeEmail({
  memberName,
  code,
  expiresInMinutes,
  purpose = "sign in",
  logoUrl,
}: OtpCodeEmailProps) {
  return (
    <EmailLayout
      headline="YOUR CODE"
      eyebrow="One-time passcode"
      preheader={`Your Siohioma code is ${code}`}
      logoUrl={logoUrl}
    >
      <p style={{ margin: "0 0 18px", fontWeight: 500, color: "#1a2620" }}>
        Hi {memberName.split(" ")[0]},
      </p>
      <p style={{ margin: "0 0 24px" }}>
        Use the code below to {purpose}. It expires in{" "}
        <strong>{expiresInMinutes} minutes</strong>.
      </p>
      <div
        style={{
          margin: "0 0 28px",
          padding: "28px 12px",
          textAlign: "center",
          backgroundColor: "#faf7ef",
          borderRadius: 16,
          border: "1px solid #efeada",
        }}
      >
        <div
          style={{
            fontFamily: emailFonts.fontDisplay,
            fontWeight: 300,
            fontSize: 46,
            letterSpacing: "0.32em",
            color: "#1f3a2e",
            paddingLeft: "0.32em",
          }}
        >
          {code}
        </div>
      </div>
      <p style={{ margin: 0, fontSize: 13, color: "#8b948f" }}>
        Never share this code. Siohioma staff will never ask you for it.
      </p>
    </EmailLayout>
  );
}
