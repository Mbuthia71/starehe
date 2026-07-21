/**
 * Siohioma email templates registry.
 *
 * Backend wiring (Windsurf):
 *   import { renderToStaticMarkup } from "react-dom/server";
 *   const html = renderToStaticMarkup(<WelcomeEmail {...props} />);
 *   // send via SMTP / Resend / SES
 */
export { EmailLayout, emailButton, emailLink, emailColors, emailFonts } from "./Layout";
export { WelcomeEmail } from "./welcome";
export { PasswordResetEmail } from "./password-reset";
export { PasswordChangedEmail } from "./password-changed";
export { OtpCodeEmail } from "./otp-code";
export { LoginAlertEmail } from "./login-alert";
export { TransactionReceiptEmail } from "./transaction-receipt";
export { LoanStatusEmail } from "./loan-status";
export { StatementReadyEmail } from "./statement-ready";

import { WelcomeEmail } from "./welcome";
import { PasswordResetEmail } from "./password-reset";
import { PasswordChangedEmail } from "./password-changed";
import { OtpCodeEmail } from "./otp-code";
import { LoginAlertEmail } from "./login-alert";
import { TransactionReceiptEmail } from "./transaction-receipt";
import { LoanStatusEmail } from "./loan-status";
import { StatementReadyEmail } from "./statement-ready";

export const emailTemplates = [
  {
    id: "welcome",
    name: "Welcome",
    subject: "Welcome to Siohioma",
    Component: WelcomeEmail,
    sample: {
      memberName: "Achieng Otieno",
      memberNumber: "SHM-04821",
      msisdn: "+254 712 884 901",
      activationUrl: "https://siohioma.co.ke/activate?t=abc123",
    },
  },
  {
    id: "password-reset",
    name: "PIN reset request",
    subject: "Reset your Siohioma PIN",
    Component: PasswordResetEmail,
    sample: {
      memberName: "Achieng Otieno",
      resetUrl: "https://siohioma.co.ke/reset?t=abc123",
      expiresInMinutes: 15,
      requestedFromIp: "41.90.64.21",
    },
  },
  {
    id: "password-changed",
    name: "PIN updated",
    subject: "Your Siohioma PIN was changed",
    Component: PasswordChangedEmail,
    sample: {
      memberName: "Achieng Otieno",
      changedAt: "12 Nov 2026, 14:32 EAT",
      deviceLabel: "iPhone 14 · Safari",
    },
  },
  {
    id: "otp",
    name: "One-time code",
    subject: "Your Siohioma code",
    Component: OtpCodeEmail,
    sample: {
      memberName: "Achieng Otieno",
      code: "482910",
      expiresInMinutes: 5,
      purpose: "sign in",
    },
  },
  {
    id: "login-alert",
    name: "New sign-in alert",
    subject: "New sign-in to your account",
    Component: LoginAlertEmail,
    sample: {
      memberName: "Achieng Otieno",
      device: "Samsung Galaxy A54 · Chrome",
      location: "Nairobi, Kenya",
      when: "12 Nov 2026, 14:32 EAT",
    },
  },
  {
    id: "receipt",
    name: "Transaction receipt",
    subject: "Receipt · Siohioma",
    Component: TransactionReceiptEmail,
    sample: {
      memberName: "Achieng Otieno",
      type: "withdrawal" as const,
      amount: "KES 12,500",
      counterparty: "M-Pesa · +254 712 884 901",
      reference: "SHM-WD-09128",
      balanceAfter: "KES 184,320",
      when: "12 Nov 2026, 14:32 EAT",
    },
  },
  {
    id: "loan-status",
    name: "Loan status update",
    subject: "Your Siohioma loan update",
    Component: LoanStatusEmail,
    sample: {
      memberName: "Achieng Otieno",
      status: "approved" as const,
      productName: "Emergency Loan",
      amount: "KES 50,000",
      nextStep:
        "Funds will be sent to your M-Pesa (+254 712 884 901) within 24 hours.",
      ctaLabel: "View loan",
      ctaUrl: "https://siohioma.co.ke/loans",
    },
  },
  {
    id: "statement-ready",
    name: "Statement ready",
    subject: "Your statement is ready",
    Component: StatementReadyEmail,
    sample: {
      memberName: "Achieng Otieno",
      period: "October 2026",
      accountLabel: "Siohioma Savings · 0148-2210",
      downloadUrl: "https://siohioma.co.ke/statements/2026-10.pdf",
    },
  },
] as const;
