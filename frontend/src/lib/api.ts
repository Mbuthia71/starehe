// Central API client for the Starehian Society Platform.
// ──────────────────────────────────────────────────────────────────────────────
// BACKEND SWAP POINT
// ──────────────────────────────────────────────────────────────────────────────
// Go Backend API
//   Base:    http://localhost:3000/api (or VITE_API_BASE_URL from env)
//   Auth:    Bearer <JWT token>
//   Routes:  POST /auth/request-otp
//            POST /auth/signup
//            POST /auth/login
//            POST /auth/admin/login
//            POST /auth/refresh
//            POST /auth/logout
//            GET  /profiles/me
//            PUT  /profiles/me
//            GET  /profiles/:id
//            POST /profiles/search
//            POST /connections/
//            POST /posts/
//            GET  /posts/feed
//            POST /chat/direct/:id
//            etc.
//
export const API_CONFIG = {
  baseUrl: (import.meta as any).env?.VITE_API_BASE_URL || "/api",
  timeout: parseInt((import.meta as any).env?.VITE_API_TIMEOUT || "30000"),
  headers: (token?: string) => ({
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  }),
} as const;

// ──────────────────────────────────────────────────────────────────────────────
// Types
// ──────────────────────────────────────────────────────────────────────────────
export type KycStatus = "Verified" | "Pending" | "Flagged";
export type AccountType = "Savings" | "Share Capital" | "Fixed Deposit" | "Junior Savings";
export type LoanStage =
  | "Applied"
  | "Committee Review"
  | "Approved"
  | "Disbursed"
  | "Active"
  | "Arrears"
  | "Closed";
export type TxStatus = "Completed" | "Pending" | "Failed";
export type LoanApplicationStatus = "Submitted" | "Under Review" | "Approved" | "Disbursed" | "Rejected";

export interface Member {
  id: string;
  firstName: string;
  lastName: string;
  memberNo: string;
  kycStatus: KycStatus;
  phone: string;
  msisdn: string; // +254…
  joinedDate: string;
  branch: string;
  biometricEnabled: boolean;
  nextOfKin: { name: string; relation: string; phone: string };
}

export interface Account {
  id: string;
  accountNo: string;
  type: AccountType;
  balance: number;
  available: number;
  interestRate: number;
  status: "Active" | "Dormant" | "Closed";
  openedOn: string;
}

export interface Guarantor {
  memberName: string;
  pledgedAmount: number;
  status: "Accepted" | "Pending" | "Declined";
}

export interface Loan {
  id: string;
  accountNo: string;
  product: string;
  principal: number;
  outstanding: number;
  paid: number;
  nextPaymentAmount: number;
  nextPaymentDate: string;
  stage: LoanStage;
  daysOverdue: number;
  guarantors: Guarantor[];
  disbursedOn?: string;
  termMonths: number;
  monthsPaid: number;
  rate: number;
}

export interface Transaction {
  id: string;
  date: string;
  description: string;
  accountId: string;
  category:
    | "Deposit"
    | "Withdrawal"
    | "Loan Repayment"
    | "Disbursement"
    | "Interest"
    | "Share Capital"
    | "Transfer";
  amount: number; // signed: + credit to member, − debit
  status: TxStatus;
  reference?: string;
}

export interface LoanProduct {
  id: string;
  name: string;
  rate: number;
  minTermMonths: number;
  maxTermMonths: number;
  minAmount: number;
  maxAmount: number;
  eligibility: string;
  description: string;
}

export interface GuarantorRequest {
  id: string;
  fromMemberName: string;
  fromMemberNo: string;
  loanProduct: string;
  loanAmount: number;
  pledgeAmount: number;
  termMonths: number;
  requestedAt: string;
  status: "Pending" | "Accepted" | "Declined";
}

export interface LoanApplication {
  id: string;
  product: string;
  amount: number;
  termMonths: number;
  submittedAt: string;
  status: LoanApplicationStatus;
  timeline: { label: string; at?: string; state: "done" | "pending" | "fail" }[];
}

// ──────────────────────────────────────────────────────────────────────────────
// Current member (the only person this app cares about)
// ──────────────────────────────────────────────────────────────────────────────
// NOTE: This is now fetched from selfApi.ts via getSelfClient()
// The mock data has been removed - use getSelfClient() instead
export const currentMember: Member = {
  id: "",
  firstName: "",
  lastName: "",
  memberNo: "",
  kycStatus: "Pending",
  phone: "",
  msisdn: "",
  joinedDate: "",
  branch: "",
  biometricEnabled: false,
  nextOfKin: { name: "", relation: "", phone: "" },
};

// ──────────────────────────────────────────────────────────────────────────────
// My accounts
// ──────────────────────────────────────────────────────────────────────────────
// NOTE: This is now fetched from selfApi.ts via getSelfAccounts()
// The mock data has been removed - use getSelfAccounts() instead
export const myAccounts: Account[] = [];

export const primaryAccount: Account = {
  id: "",
  accountNo: "",
  type: "Savings",
  balance: 0,
  available: 0,
  interestRate: 0,
  status: "Active",
  openedOn: "",
};

// ──────────────────────────────────────────────────────────────────────────────
// My loans
// ──────────────────────────────────────────────────────────────────────────────
// NOTE: This is now fetched from selfApi.ts via getSelfLoans()
// The mock data has been removed - use getSelfLoans() instead
export const myLoans: Loan[] = [];

export const loanApplications: LoanApplication[] = [];

// ──────────────────────────────────────────────────────────────────────────────
// Transactions — this member's ledger, newest first
// Negative = debit from member, Positive = credit to member.
// ──────────────────────────────────────────────────────────────────────────────
// NOTE: This is now fetched from selfApi.ts via getSelfTransactions()
// The mock data has been removed - use getSelfTransactions() instead
export const myTransactions: Transaction[] = [];

// Latest few, for the home card
export const recentTransactions: Transaction[] = [];

// ──────────────────────────────────────────────────────────────────────────────
// Loan products the member can apply for
// ──────────────────────────────────────────────────────────────────────────────
// NOTE: This is now fetched from selfApi.ts via getSelfLoanProducts()
// The mock data has been removed - use getSelfLoanProducts() instead
export const loanProducts: LoanProduct[] = [];

// ──────────────────────────────────────────────────────────────────────────────
// Guarantor pledge requests sent to ME by other members
// ──────────────────────────────────────────────────────────────────────────────
// NOTE: This would be fetched from Fineract guarantor API
// The mock data has been removed - implement real API call
export const guarantorRequests: GuarantorRequest[] = [];

// ──────────────────────────────────────────────────────────────────────────────
// Statements — months for which a downloadable PDF exists
// ──────────────────────────────────────────────────────────────────────────────
// NOTE: This would be fetched from Fineract statement API
// The mock data has been removed - implement real API call
export const availableStatements: any[] = [];

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────
export const totalSavings = 0;

export const activeLoan: Loan | undefined = undefined;
