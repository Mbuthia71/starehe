// Payment Gateway Integrations
// Pesalink, RTGS, and Bank Transfer APIs

import { authHeaders, getClientId } from "./auth";
import { config } from "./config";

const BASE = config.fineract.baseUrl;

export interface Bank {
  code: string;
  name: string;
  pesalinkCode?: string;
  rtgsCode?: string;
}

export interface BankAccount {
  accountNumber: string;
  accountName: string;
  bankCode: string;
}

export interface TransferRequest {
  type: "pesalink" | "rtgs" | "bank";
  fromAccount: string;
  toAccount: string;
  toBankCode: string;
  toAccountName: string;
  amount: number;
  reference?: string;
  toBankName?: string;
}

export interface TransferResponse {
  success: boolean;
  transactionId?: string;
  message: string;
  status?: string;
}

// Kenyan Banks List (simplified)
const KENYAN_BANKS: Bank[] = [
  { code: "01", name: "Equity Bank", pesalinkCode: "EQBK", rtgsCode: "EQBK" },
  { code: "07", name: "KCB Bank", pesalinkCode: "KCBL", rtgsCode: "KCBL" },
  { code: "23", name: "Cooperative Bank", pesalinkCode: "COOP", rtgsCode: "COOP" },
  { code: "11", name: "Standard Chartered", pesalinkCode: "SCBL", rtgsCode: "SCBL" },
  { code: "63", name: "Absa Bank", pesalinkCode: "ABSA", rtgsCode: "ABSA" },
  { code: "33", name: "Diamond Trust Bank", pesalinkCode: "DTBK", rtgsCode: "DTBK" },
  { code: "58", name: "Family Bank", pesalinkCode: "FMBK", rtgsCode: "FMBK" },
  { code: "17", name: "National Bank", pesalinkCode: "NBKE", rtgsCode: "NBKE" },
  { code: "75", name: "NCBA Bank", pesalinkCode: "NCBA", rtgsCode: "NCBA" },
  { code: "54", name: "I&M Bank", pesalinkCode: "IMBK", rtgsCode: "IMBK" },
  { code: "18", name: "NIC Bank", pesalinkCode: "NICB", rtgsCode: "NICB" },
  { code: "44", name: "Prime Bank", pesalinkCode: "PRMB", rtgsCode: "PRMB" },
  { code: "70", name: "Sidian Bank", pesalinkCode: "SDNB", rtgsCode: "SDNB" },
  { code: "82", name: "Spire Bank", pesalinkCode: "SPRB", rtgsCode: "SPRB" },
  { code: "08", name: "Barclays Bank", pesalinkCode: "BARC", rtgsCode: "BARC" },
];

// Get all banks
export function getBanks(): Bank[] {
  return KENYAN_BANKS;
}

// Get bank by code
export function getBankByCode(code: string): Bank | undefined {
  return KENYAN_BANKS.find(bank => bank.code === code || bank.pesalinkCode === code || bank.rtgsCode === code);
}

// Validate account number format (basic validation)
export function validateAccountNumber(accountNumber: string, bankCode: string): boolean {
  const bank = getBankByCode(bankCode);
  if (!bank) return false;
  
  // Most Kenyan banks use 10-14 digit account numbers
  const cleaned = accountNumber.replace(/\s/g, "");
  return cleaned.length >= 10 && cleaned.length <= 14 && /^\d+$/.test(cleaned);
}

// Simulate account name lookup (in production, this would call the bank's API)
export async function lookupAccountName(accountNumber: string, bankCode: string): Promise<string> {
  // Simulate API call delay
  await new Promise(resolve => setTimeout(resolve, 1000));
  
  // In production, this would call the actual bank's API or Pesalink/RTGS gateway
  // For now, return a placeholder
  return "Account Holder Name (Verify with Bank)";
}

// Pesalink Transfer (Real-time, up to 1M KES)
export async function processPesalinkTransfer(request: TransferRequest): Promise<TransferResponse> {
  try {
    // Validate
    if (!validateAccountNumber(request.toAccount, request.toBankCode)) {
      throw new Error("Invalid account number format");
    }

    if (request.amount > 1000000) {
      throw new Error("Pesalink transfers limited to 1,000,000 KES. Use RTGS for larger amounts.");
    }

    // In production, this would call the Pesalink API
    // For now, simulate the transfer
    await new Promise(resolve => setTimeout(resolve, 2000));

    const transactionId = `PSLK${Date.now()}${Math.random().toString(36).slice(2, 7).toUpperCase()}`;

    // Record the transfer in Fineract
    const bank = getBankByCode(request.toBankCode);
    await recordExternalTransfer(
      request.fromAccount,
      request.amount,
      "Pesalink",
      transactionId,
      bank?.name || request.toBankCode,
      request.toAccount
    );

    return {
      success: true,
      transactionId,
      message: "Pesalink transfer initiated successfully",
      status: "PENDING",
    };
  } catch (error: any) {
    return {
      success: false,
      message: error.message || "Pesalink transfer failed",
    };
  }
}

// RTGS Transfer (High-value, same-day settlement)
export async function processRtgsTransfer(request: TransferRequest): Promise<TransferResponse> {
  try {
    // Validate
    if (!validateAccountNumber(request.toAccount, request.toBankCode)) {
      throw new Error("Invalid account number format");
    }

    if (request.amount < 1000000) {
      throw new Error("RTGS transfers require minimum 1,000,000 KES. Use Pesalink for smaller amounts.");
    }

    // In production, this would call the RTGS API via Central Bank
    // For now, simulate the transfer
    await new Promise(resolve => setTimeout(resolve, 3000));

    const transactionId = `RTGS${Date.now()}${Math.random().toString(36).slice(2, 7).toUpperCase()}`;

    // Record the transfer in Fineract
    const bank = getBankByCode(request.toBankCode);
    await recordExternalTransfer(
      request.fromAccount,
      request.amount,
      "RTGS",
      transactionId,
      bank?.name || request.toBankCode,
      request.toAccount
    );

    return {
      success: true,
      transactionId,
      message: "RTGS transfer submitted for processing",
      status: "PROCESSING",
    };
  } catch (error: any) {
    return {
      success: false,
      message: error.message || "RTGS transfer failed",
    };
  }
}

// Bank Transfer (Standard inter-bank transfer)
export async function processBankTransfer(request: TransferRequest): Promise<TransferResponse> {
  try {
    // Validate
    if (!validateAccountNumber(request.toAccount, request.toBankCode)) {
      throw new Error("Invalid account number format");
    }

    // In production, this would call the bank's API or clearing house
    // For now, simulate the transfer
    await new Promise(resolve => setTimeout(resolve, 2500));

    const transactionId = `BNKTR${Date.now()}${Math.random().toString(36).slice(2, 7).toUpperCase()}`;

    // Record the transfer in Fineract
    const bank = getBankByCode(request.toBankCode);
    await recordExternalTransfer(
      request.fromAccount,
      request.amount,
      "Bank Transfer",
      transactionId,
      bank?.name || request.toBankCode,
      request.toAccount
    );

    return {
      success: true,
      transactionId,
      message: "Bank transfer initiated successfully",
      status: "PENDING",
    };
  } catch (error: any) {
    return {
      success: false,
      message: error.message || "Bank transfer failed",
    };
  }
}

// Check transaction status
export async function checkTransferStatus(transactionId: string): Promise<TransferResponse> {
  // In production, this would query the payment gateway
  await new Promise(resolve => setTimeout(resolve, 500));
  
  return {
    success: true,
    transactionId,
    message: "Transaction completed successfully",
    status: "COMPLETED",
  };
}

// Get transfer limits
export function getTransferLimits(type: "pesalink" | "rtgs" | "bank") {
  switch (type) {
    case "pesalink":
      return { min: 10, max: 1000000, fee: 0, settlement: "Real-time" };
    case "rtgs":
      return { min: 1000000, max: 10000000000, fee: 500, settlement: "Same-day" };
    case "bank":
      return { min: 100, max: 50000000, fee: 100, settlement: "1-2 business days" };
  }
}

// Record external transfer in Fineract as a withdrawal transaction
async function recordExternalTransfer(accountId: string, amount: number, type: string, reference: string, toBank: string, toAccount: string) {
  try {
    const clientId = getClientId();
    if (!clientId) throw new Error("Not authenticated");

    const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });

    const res = await fetch(`${BASE}/savingsaccounts/${accountId}/transactions`, {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify({
        dateFormat: "dd MMM yyyy",
        locale: "en",
        transactionDate: today,
        transactionAmount: amount,
        paymentTypeId: 1, // Cash payment
        transactionType: "withdrawal",
        note: `${type} transfer to ${toBank} - ${toAccount}. Ref: ${reference}`,
      }),
    });

    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      console.error("Failed to record external transfer in Fineract:", err);
      // Don't throw - the transfer already succeeded, just log the error
    }

    console.log("External transfer recorded in Fineract");
  } catch (error) {
    console.error("Error recording external transfer:", error);
  }
}
