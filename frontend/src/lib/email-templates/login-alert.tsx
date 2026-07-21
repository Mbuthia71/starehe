import { EmailLayout, emailLink } from "./Layout";

export interface LoginAlertEmailProps {
  memberName: string;
  device: string;
  location: string;
  when: string;
  logoUrl?: string;
}

export function LoginAlertEmail({
  memberName,
  device,
  location,
  when,
  logoUrl,
}: LoginAlertEmailProps) {
  return (
    <EmailLayout
      headline="NEW SIGN-IN"
      eyebrow="Device alert"
      preheader={`New sign-in to your Siohioma account from ${device}.`}
      logoUrl={logoUrl}
    >
      <p style={{ margin: "0 0 18px", fontWeight: 500, color: "#1a2620" }}>
        Hi {memberName.split(" ")[0]},
      </p>
      <p style={{ margin: "0 0 24px" }}>
        We noticed a new sign-in to your Siohioma account.
      </p>
      <table
        role="presentation"
        cellPadding={0}
        cellSpacing={0}
        width="100%"
        style={{
          marginBottom: 28,
          borderTop: "1px solid #e6e2d8",
        }}
      >
        <tbody>
          {[
            ["Device", device],
            ["Where", location],
            ["When", when],
          ].map(([k, v]) => (
            <tr key={k}>
              <td
                style={{
                  padding: "12px 0",
                  borderBottom: "1px solid #e6e2d8",
                  fontSize: 12,
                  textTransform: "uppercase",
                  letterSpacing: "0.14em",
                  color: "#8b948f",
                  width: 120,
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
                }}
              >
                {v}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <p style={{ margin: 0, fontSize: 13, color: "#8b948f" }}>
        If that was you, no action needed. If not, reset your PIN now and call{" "}
        <a href="tel:+254711000000" style={emailLink}>
          +254 711 000 000
        </a>
        .
      </p>
    </EmailLayout>
  );
}
