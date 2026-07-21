// Email sending via Resend + local renderer
import { config } from "./config";
import { logger } from "./logger";

const RENDERER_URL = config.email.rendererUrl;
const RESEND_API_KEY = config.email.resendApiKey;
const FROM_EMAIL = config.email.from;
const REPLY_TO = config.email.replyTo;

export interface EmailTemplate {
  template: string;
  props: Record<string, unknown>;
}

export interface SendEmailOptions {
  to: string | string[];
  subject: string;
  template: string;
  props: Record<string, unknown>;
}

/**
 * Render email template via the renderer service
 */
async function renderTemplate(template: string, props: Record<string, unknown>): Promise<string> {
  const res = await fetch(`${RENDERER_URL}/render/${template}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(props),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.error || `Renderer error ${res.status}`);
  }

  const data = await res.json();
  return data.html;
}

/**
 * Send email via Resend API
 */
async function sendViaResend(to: string | string[], subject: string, html: string): Promise<void> {
  if (!RESEND_API_KEY) {
    throw new Error("RESEND_API_KEY not configured");
  }

  const recipients = Array.isArray(to) ? to : [to];

  const res = await fetch("https://api.resend.com/emails", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${RESEND_API_KEY}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      from: FROM_EMAIL,
      to: recipients,
      subject,
      html,
      reply_to: REPLY_TO,
    }),
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || `Resend error ${res.status}`);
  }

  logger.info("Email sent via Resend", { to: recipients, subject });
}

/**
 * Render and send an email
 */
export async function sendEmail({ to, subject, template, props }: SendEmailOptions): Promise<void> {
  try {
    logger.info("Rendering email template", { template, to });
    const html = await renderTemplate(template, props);
    await sendViaResend(to, subject, html);
    logger.info("Email sent successfully", { template, to });
  } catch (error) {
    logger.error("Failed to send email", { template, to, error: (error as Error).message });
    throw error;
  }
}

/**
 * Quick helpers for common email types
 */
export const emailHelpers = {
  sendOtp: (to: string, code: string, memberName: string) =>
    sendEmail({
      to,
      subject: "Your Siohioma code",
      template: "otp",
      props: { memberName, code, expiresInMinutes: 5, purpose: "sign in" },
    }),

  sendWelcome: (to: string, memberName: string, memberNumber: string, activationUrl: string) =>
    sendEmail({
      to,
      subject: "Welcome to Siohioma",
      template: "welcome",
      props: { memberName, memberNumber, msisdn: "", activationUrl },
    }),

  sendTransactionReceipt: (to: string, memberName: string, type: string, amount: string, reference: string, balanceAfter: string) =>
    sendEmail({
      to,
      subject: `Receipt · Siohioma`,
      template: "transaction-receipt",
      props: { memberName, type, amount, counterparty: "Siohioma SACCO", reference, balanceAfter, when: new Date().toLocaleString() },
    }),

  sendLoanStatus: (to: string, memberName: string, status: string, productName: string, amount: string) =>
    sendEmail({
      to,
      subject: "Your Siohioma loan update",
      template: "loan-status",
      props: { memberName, status, productName, amount, nextStep: "Check your member portal for details", ctaLabel: "View loan", ctaUrl: "https://banking.siohioma.com/loans" },
    }),

  sendPasswordReset: (to: string, memberName: string, resetUrl: string, requestedFromIp?: string) =>
    sendEmail({
      to,
      subject: "Reset your Siohioma PIN",
      template: "password-reset",
      props: { memberName, resetUrl, expiresInMinutes: 15, requestedFromIp, logoUrl: "https://banking.siohioma.com/siohioma-logo.png" },
    }),

  sendPasswordChanged: (to: string, memberName: string, deviceLabel?: string) =>
    sendEmail({
      to,
      subject: "Your Siohioma PIN was updated",
      template: "password-changed",
      props: { memberName, changedAt: new Date().toLocaleString(), deviceLabel, logoUrl: "https://banking.siohioma.com/siohioma-logo.png" },
    }),

  sendLoginAlert: (to: string, memberName: string, device: string, location: string) =>
    sendEmail({
      to,
      subject: "New sign-in to your Siohioma account",
      template: "login-alert",
      props: { memberName, device, location, when: new Date().toLocaleString(), logoUrl: "https://banking.siohioma.com/siohioma-logo.png" },
    }),

  sendStatementReady: (to: string, memberName: string, period: string, accountLabel: string, downloadUrl: string) =>
    sendEmail({
      to,
      subject: `Your ${period} statement is ready`,
      template: "statement-ready",
      props: { memberName, period, accountLabel, downloadUrl, logoUrl: "https://banking.siohioma.com/siohioma-logo.png" },
    }),
};
