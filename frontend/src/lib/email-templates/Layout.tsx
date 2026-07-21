import type { ReactNode, CSSProperties } from "react";

/**
 * Email-safe layout scaffold for all Siohioma transactional emails.
 *
 * Design language:
 *  - Centered logo (top)
 *  - Massive ultra-thin display headline (Cormorant Garamond 300)
 *  - Short gold rule underneath
 *  - Quiet body text (Jost 400)
 *  - Script signoff (Dancing Script)
 *  - Hairline + footer address
 *
 * NOTE for backend (Windsurf): render this with @react-email/render or
 * renderToStaticMarkup. All styling is inline so it survives Gmail / Outlook.
 * Replace `logoUrl` with a publicly hosted absolute URL before sending.
 */

export interface EmailLayoutProps {
  /** Big all-caps display word, e.g. "WELCOME!" */
  headline: string;
  /** Short text under the gold rule. Optional. */
  eyebrow?: string;
  /** Email body content */
  children: ReactNode;
  /** Closing line above the signature */
  signoff?: string;
  /** Signature line, rendered in script */
  signature?: string;
  /** Absolute URL to the Siohioma logo */
  logoUrl?: string;
  /** Preheader (preview text shown in inbox list) */
  preheader?: string;
}

const colors = {
  ink: "#1a2620",
  body: "#3a4a42",
  muted: "#8b948f",
  hairline: "#e6e2d8",
  cream: "#faf7ef",
  gold: "#d4a04a",
  green: "#1f3a2e",
};

const fontDisplay =
  '"Cormorant Garamond", "Forum", Georgia, "Times New Roman", serif';
const fontBody =
  '"Jost", "Helvetica Neue", Arial, sans-serif';
const fontScript = '"Dancing Script", "Apple Chancery", cursive';

export function EmailLayout({
  headline,
  eyebrow,
  children,
  signoff = "Cheers,",
  signature = "The Siohioma Team",
  logoUrl = "https://siohioma.co.ke/logo.png",
  preheader,
}: EmailLayoutProps) {
  const outer: CSSProperties = {
    margin: 0,
    padding: "40px 16px",
    backgroundColor: "#ffffff",
    fontFamily: fontBody,
    color: colors.body,
    WebkitFontSmoothing: "antialiased",
  };

  const wrap: CSSProperties = {
    maxWidth: 560,
    margin: "0 auto",
  };

  return (
    <div style={outer}>
      {preheader && (
        <div
          style={{
            display: "none",
            overflow: "hidden",
            lineHeight: 1,
            opacity: 0,
            maxHeight: 0,
            maxWidth: 0,
            color: "transparent",
          }}
        >
          {preheader}
        </div>
      )}

      <table
        role="presentation"
        cellPadding={0}
        cellSpacing={0}
        width="100%"
        style={wrap}
      >
        <tbody>
          {/* Logo */}
          <tr>
            <td align="center" style={{ paddingTop: 8, paddingBottom: 36 }}>
              <img
                src={logoUrl}
                alt="Siohioma"
                width={68}
                height={68}
                style={{
                  display: "block",
                  width: 68,
                  height: 68,
                  objectFit: "contain",
                }}
              />
            </td>
          </tr>

          {/* Headline */}
          <tr>
            <td align="center" style={{ paddingBottom: 8 }}>
              <h1
                style={{
                  margin: 0,
                  fontFamily: fontDisplay,
                  fontWeight: 300,
                  fontSize: 52,
                  lineHeight: 1.05,
                  letterSpacing: "0.04em",
                  color: colors.ink,
                }}
              >
                {headline}
              </h1>
            </td>
          </tr>

          {/* Gold rule */}
          <tr>
            <td align="center" style={{ paddingTop: 14, paddingBottom: 8 }}>
              <div
                style={{
                  width: 36,
                  height: 3,
                  backgroundColor: colors.gold,
                  borderRadius: 2,
                  margin: "0 auto",
                }}
              />
            </td>
          </tr>

          {eyebrow && (
            <tr>
              <td
                align="center"
                style={{
                  paddingTop: 6,
                  paddingBottom: 24,
                  fontSize: 11,
                  fontWeight: 500,
                  letterSpacing: "0.18em",
                  textTransform: "uppercase",
                  color: colors.muted,
                }}
              >
                {eyebrow}
              </td>
            </tr>
          )}

          {/* Body */}
          <tr>
            <td
              style={{
                paddingTop: eyebrow ? 8 : 32,
                paddingBottom: 32,
                fontSize: 15,
                lineHeight: 1.7,
                color: colors.body,
              }}
            >
              {children}
            </td>
          </tr>

          {/* Signoff */}
          <tr>
            <td style={{ paddingBottom: 4, fontSize: 15, color: colors.body }}>
              {signoff}
            </td>
          </tr>
          <tr>
            <td
              style={{
                paddingBottom: 40,
                fontFamily: fontScript,
                fontSize: 26,
                color: colors.green,
              }}
            >
              {signature}
            </td>
          </tr>

          {/* Hairline */}
          <tr>
            <td>
              <div
                style={{
                  height: 1,
                  backgroundColor: colors.hairline,
                  margin: "0 0 28px 0",
                }}
              />
            </td>
          </tr>

          {/* Socials */}
          <tr>
            <td align="center" style={{ paddingBottom: 20 }}>
              {[
                { label: "fb", bg: "#2f6f8f" },
                { label: "tw", bg: "#4ab3e0" },
                { label: "ig", bg: "#e0744a" },
                { label: "wa", bg: "#3fae5a" },
              ].map((s) => (
                <span
                  key={s.label}
                  style={{
                    display: "inline-block",
                    width: 28,
                    height: 28,
                    lineHeight: "28px",
                    margin: "0 4px",
                    borderRadius: 999,
                    backgroundColor: s.bg,
                    color: "#ffffff",
                    fontSize: 10,
                    fontWeight: 600,
                    textTransform: "uppercase",
                    letterSpacing: "0.05em",
                    textAlign: "center",
                  }}
                >
                  {s.label}
                </span>
              ))}
            </td>
          </tr>

          {/* Footer */}
          <tr>
            <td
              align="center"
              style={{
                fontSize: 11,
                lineHeight: 1.7,
                color: colors.muted,
                letterSpacing: "0.02em",
                paddingBottom: 8,
              }}
            >
              Siohioma SACCO Society Ltd
              <br />
              Kimathi Street, 4th Floor · Nairobi, Kenya
              <br />
              Regulated by SASRA · License No. CS/1820
            </td>
          </tr>
          <tr>
            <td
              align="center"
              style={{
                fontSize: 11,
                color: colors.muted,
                paddingTop: 6,
              }}
            >
              <a
                href="https://siohioma.co.ke"
                style={{ color: colors.muted, textDecoration: "underline" }}
              >
                siohioma.co.ke
              </a>
              {"  ·  "}
              <a
                href="mailto:hello@siohioma.co.ke"
                style={{ color: colors.muted, textDecoration: "underline" }}
              >
                hello@siohioma.co.ke
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}

/** Shared button style for CTAs inside email bodies. */
export const emailButton: CSSProperties = {
  display: "inline-block",
  padding: "14px 28px",
  backgroundColor: colors.green,
  color: "#ffffff",
  fontFamily: fontBody,
  fontSize: 13,
  fontWeight: 500,
  letterSpacing: "0.14em",
  textTransform: "uppercase",
  textDecoration: "none",
  borderRadius: 999,
};

export const emailLink: CSSProperties = {
  color: colors.green,
  textDecoration: "underline",
  fontWeight: 500,
};

export const emailColors = colors;
export const emailFonts = { fontDisplay, fontBody, fontScript };
