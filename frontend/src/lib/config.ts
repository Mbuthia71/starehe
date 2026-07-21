// Centralized configuration for the Siohioma SACCO system
// All environment variables should be set in .env or environment

export const config = {
  // Fineract API Configuration
  fineract: {
    baseUrl: import.meta.env.VITE_FINERACT_BASE_URL || "https://banking.siohioma.com/fineract-provider/api/v1",
    tenantId: import.meta.env.VITE_FINERACT_TENANT_ID || "default",
  },

  // M-Pesa Configuration
  mpesa: {
    baseUrl: import.meta.env.VITE_MPESA_BASE_URL || "https://204.168.218.166:8081",
    shortcode: import.meta.env.VITE_MPESA_SHORTCODE || "",
    passkey: import.meta.env.VITE_MPESA_PASSKEY || "",
    consumerKey: import.meta.env.VITE_MPESA_CONSUMER_KEY || "",
    consumerSecret: import.meta.env.VITE_MPESA_CONSUMER_SECRET || "",
  },

  // Payment Gateway Configuration
  paymentGateway: {
    pesalinkEnabled: import.meta.env.VITE_PESALINK_ENABLED === "true",
    rtgsEnabled: import.meta.env.VITE_RTGS_ENABLED === "true",
    bankTransferEnabled: import.meta.env.VITE_BANK_TRANSFER_ENABLED === "true",
  },

  // Email Configuration
  email: {
    rendererUrl: import.meta.env.VITE_EMAIL_RENDERER_URL || "http://204.168.218.166:3002",
    from: import.meta.env.VITE_EMAIL_FROM || "noreply@siohioma.com",
    replyTo: import.meta.env.VITE_EMAIL_REPLY_TO || "support@siohioma.com",
    resendApiKey: import.meta.env.RESEND_API_KEY || "",
  },

  // Loan Fee Configuration
  loan: {
    processingFee: Number(import.meta.env.VITE_LOAN_PROCESSING_FEE) || 500,
    microfinanceCut: Number(import.meta.env.VITE_LOAN_PROCESSING_FEE) * 0.4 || 200,
    saccoCut: Number(import.meta.env.VITE_LOAN_PROCESSING_FEE) * 0.6 || 300,
    minAmount: Number(import.meta.env.VITE_LOAN_MIN_AMOUNT) || 5000,
    maxAmount: Number(import.meta.env.VITE_LOAN_MAX_AMOUNT) || 3000000,
  },

  // Transfer Configuration
  transfer: {
    minAmount: Number(import.meta.env.VITE_TRANSFER_MIN_AMOUNT) || 50,
    maxAmount: Number(import.meta.env.VITE_TRANSFER_MAX_AMOUNT) || 1000000,
  },

  // Application Configuration
  app: {
    name: import.meta.env.VITE_APP_NAME || "Siohioma SACCO",
    version: import.meta.env.VITE_APP_VERSION || "1.0.0",
    environment: import.meta.env.MODE || "development",
    sessionTimeoutMinutes: Number(import.meta.env.VITE_SESSION_TIMEOUT_MINUTES) || 30,
  },

  // Security Configuration
  security: {
    maxLoginAttempts: Number(import.meta.env.VITE_MAX_LOGIN_ATTEMPTS) || 5,
    otpExpiryMinutes: Number(import.meta.env.VITE_OTP_EXPIRY_MINUTES) || 10,
    pinMaxAttempts: Number(import.meta.env.VITE_PIN_MAX_ATTEMPTS) || 3,
  },

  // Account Number Sequences
  accountNumbers: {
    // Savings accounts: 1008100001, 1008100002, etc.
    savings: {
      prefix: "10081",
      startSequence: 1,
      format: (seq: number) => `10081${String(seq).padStart(4, '0')}`,
    },
    // Loan accounts: 2008200001, 2008200002, etc.
    loans: {
      prefix: "20082",
      startSequence: 1,
      format: (seq: number) => `20082${String(seq).padStart(4, '0')}`,
    },
    // Share capital accounts: 3008300001, 3008300002, etc.
    shareCapital: {
      prefix: "30083",
      startSequence: 1,
      format: (seq: number) => `30083${String(seq).padStart(4, '0')}`,
    },
    // Fixed deposit accounts: 4008400001, 4008400002, etc.
    fixedDeposit: {
      prefix: "40084",
      startSequence: 1,
      format: (seq: number) => `40084${String(seq).padStart(4, '0')}`,
    },
  },
} as const;

// Type-safe config access
export type Config = typeof config;
