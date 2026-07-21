// Data masking for sensitive client information
// Masks PII for non-admin users to protect privacy

import { getRole, canViewAllData } from "./auth";

export type MaskableField = "phone" | "email" | "idNumber" | "accountNumber" | "address";

// Masking patterns
const maskPatterns = {
  phone: (value: string) => {
    if (!value) return "";
    const cleaned = value.replace(/\s/g, "");
    if (cleaned.length < 4) return value;
    return cleaned.slice(0, -4).replace(/./g, "*") + cleaned.slice(-4);
  },
  email: (value: string) => {
    if (!value) return "";
    const [local, domain] = value.split("@");
    if (!local || !domain) return value;
    const maskedLocal = local.slice(0, 2) + "*".repeat(Math.max(0, local.length - 2));
    return `${maskedLocal}@${domain}`;
  },
  idNumber: (value: string) => {
    if (!value) return "";
    if (value.length < 4) return value;
    return "*".repeat(value.length - 4) + value.slice(-4);
  },
  accountNumber: (value: string) => {
    if (!value) return "";
    if (value.length < 4) return value;
    return "*".repeat(value.length - 4) + value.slice(-4);
  },
  address: (value: string) => {
    if (!value) return "";
    const words = value.split(" ");
    if (words.length < 2) return value;
    return words[0] + " " + "*".repeat(value.length - words[0].length - 1);
  },
};

export function maskField(field: MaskableField, value: string): string {
  // Admins with view_all permission see unmasked data
  if (canViewAllData()) {
    return value;
  }

  const masker = maskPatterns[field];
  return masker ? masker(value) : value;
}

export function maskClientData(client: any): any {
  if (canViewAllData()) {
    return client;
  }

  return {
    ...client,
    mobileNo: maskField("phone", client.mobileNo || ""),
    emailAddress: maskField("email", client.emailAddress || ""),
    accountNo: maskField("accountNumber", client.accountNo || ""),
    // Add other sensitive fields as needed
  };
}

export function maskAccountData(account: any): any {
  if (canViewAllData()) {
    return account;
  }

  return {
    ...account,
    accountNo: maskField("accountNumber", account.accountNo || ""),
  };
}

export function maskTransactionData(transaction: any): any {
  if (canViewAllData()) {
    return transaction;
  }

  // Don't mask transaction amounts, but could mask account numbers
  return {
    ...transaction,
    accountId: maskField("accountNumber", transaction.accountId || ""),
  };
}
