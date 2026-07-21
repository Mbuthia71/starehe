// All Fineract self-service API calls + M-Pesa bridge for the member app
import { authHeaders, getClientId } from "./auth";
import { config } from "./config";
import { logger, logAPIError } from "./logger";

const BASE = config.fineract.baseUrl;
const MPESA_BASE = config.mpesa.baseUrl;

async function selfGet(path: string) {
  const res = await fetch(`${BASE}${path}`, { headers: authHeaders() });
  if (res.status === 401) {
    logger.warn("Session expired", { path });
    throw new Error("Session expired. Please log in again.");
  }
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    logAPIError(path, { status: res.status, error: err });
    throw new Error(err.defaultUserMessage || `Error ${res.status}`);
  }
  return res.json();
}

async function selfPost(path: string, body: object) {
  const res = await fetch(`${BASE}${path}`, {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify(body),
  });
  if (res.status === 401) {
    logger.warn("Session expired", { path });
    throw new Error("Session expired. Please log in again.");
  }
  const data = await res.json();
  if (!res.ok) {
    logAPIError(path, { status: res.status, error: data, body });
    throw new Error(data.defaultUserMessage || data.errors?.[0]?.defaultUserMessage || JSON.stringify(data.errors) || `Error ${res.status}`);
  }
  return data;
}

// ─── Member / Client ────────────────────────────────────────────────────────

export async function getSelfClient() {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  const data = await selfGet(`/clients/${clientId}`);
  return {
    id: String(data.id),
    firstName: data.firstname || data.displayName?.split(" ")[0] || "",
    lastName: data.lastname || data.displayName?.split(" ").slice(1).join(" ") || "",
    memberNo: data.accountNo || `TJS-${data.id}`,
    kycStatus: data.active ? "Verified" : "Pending",
    phone: data.mobileNo || "",
    msisdn: data.mobileNo?.replace(/\s/g, "") || "",
    joinedDate: data.activationDate?.join?.("-") || "",
    branch: data.officeName || "",
    biometricEnabled: false,
    nextOfKin: { name: "", relation: "", phone: "" },
  };
}

// ─── Savings Accounts ────────────────────────────────────────────────────────

export async function getSelfAccounts() {
  const clientId = getClientId();
  console.log("Fetching accounts for clientId:", clientId);
  if (!clientId) throw new Error("Not authenticated");
  const data = await selfGet(`/savingsaccounts`);
  const items: any[] = data.pageItems ?? data.savingsAccounts ?? [];
  console.log("Total accounts returned:", items.length);
  // Filter by clientId on client side since API doesn't support the parameter
  const filtered = items.filter((a: any) => a.clientId === clientId);
  console.log("Filtered accounts for clientId", clientId, ":", filtered.length, filtered.map((a: any) => ({ id: a.id, accountNo: a.accountNo, clientId: a.clientId })));
  return filtered.map((a: any) => ({
    id: String(a.id),
    accountNo: a.accountNo,
    type: mapSavingsType(a.depositType?.value || a.savingsProductName || "Savings"),
    balance: a.summary?.accountBalance ?? 0,
    available: a.summary?.availableBalance ?? a.summary?.accountBalance ?? 0,
    interestRate: a.nominalAnnualInterestRate ?? 0,
    status: a.status?.value === "Active" ? "Active" : a.status?.value === "Dormant" ? "Dormant" : "Closed",
    openedOn: a.timeline?.activatedOnDate?.join?.("-") || "",
  }));
}

function mapSavingsType(raw: string): string {
  if (/share/i.test(raw)) return "Share Capital";
  if (/fixed|deposit/i.test(raw)) return "Fixed Deposit";
  if (/junior|target/i.test(raw)) return "Junior Savings";
  return "Savings";
}

// ─── Transactions ────────────────────────────────────────────────────────────

export async function getSelfTransactions(accountId: string) {
  const data = await selfGet(`/savingsaccounts/${accountId}?associations=transactions`);
  const items: any[] = data.pageItems ?? data.transactions ?? [];
  return items.map((t: any) => ({
    id: String(t.id),
    date: t.date?.join?.("-") || "",
    description: t.transactionType?.value || "Transaction",
    accountId,
    category: mapTxCategory(t.transactionType?.value || ""),
    amount: t.transactionType?.withdrawal ? -t.amount : t.amount,
    status: t.reversed ? "Failed" : "Completed",
    reference: t.paymentDetailData?.receiptNumber || undefined,
  }));
}

function mapTxCategory(type: string): string {
  if (/deposit/i.test(type)) return "Deposit";
  if (/withdrawal/i.test(type)) return "Withdrawal";
  if (/repayment/i.test(type)) return "Loan Repayment";
  if (/disburse/i.test(type)) return "Disbursement";
  if (/interest/i.test(type)) return "Interest";
  return "Deposit";
}

// ─── Loans ───────────────────────────────────────────────────────────────────

export async function getSelfLoans() {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  const data = await selfGet(`/loans`);
  const items: any[] = data.pageItems ?? data.loans ?? [];
  // Filter by clientId on client side since API doesn't support the parameter
  const filtered = items.filter((l: any) => l.clientId === clientId);
  return filtered.map((l: any) => ({
    id: String(l.id),
    accountNo: l.accountNo,
    product: l.loanProductName,
    principal: l.principal ?? l.approvedPrincipal ?? 0,
    outstanding: l.summary?.principalOutstanding ?? l.principal ?? 0,
    paid: l.summary?.principalPaid ?? 0,
    nextPaymentAmount: l.repaymentSchedule?.periods?.find((p: any) => !p.complete)?.totalDueForPeriod ?? 0,
    nextPaymentDate: l.repaymentSchedule?.periods?.find((p: any) => !p.complete)?.dueDate?.join?.("-") ?? "",
    stage: mapLoanStage(l.status?.value || ""),
    daysOverdue: l.summary?.overdueSinceDate ? Math.floor((Date.now() - new Date(l.summary.overdueSinceDate.join("-")).getTime()) / 86400000) : 0,
    guarantors: [],
    disbursedOn: l.timeline?.actualDisbursementDate?.join?.("-"),
    termMonths: l.numberOfRepayments ?? 12,
    monthsPaid: l.summary?.numberOfRepaymentsMade ?? 0,
    rate: l.annualInterestRate ?? l.interestRatePerPeriod ?? 0,
  }));
}

function mapLoanStage(status: string): string {
  if (/pending/i.test(status)) return "Applied";
  if (/approved/i.test(status)) return "Approved";
  if (/active/i.test(status)) return "Active";
  if (/arrears/i.test(status) || /overdue/i.test(status)) return "Arrears";
  if (/close/i.test(status)) return "Closed";
  if (/disburse/i.test(status)) return "Disbursed";
  return "Applied";
}

// ─── Loan Products ────────────────────────────────────────────────────────────

export async function getSelfLoanProducts() {
  const data = await selfGet("/loanproducts");
  return (Array.isArray(data) ? data : []).map((p: any) => ({
    id: String(p.id),
    name: p.name,
    rate: p.annualInterestRate ?? p.interestRatePerPeriod ?? 0,
    minTermMonths: p.minNumberOfRepayments ?? 1,
    maxTermMonths: p.maxNumberOfRepayments ?? 60,
    minAmount: p.minPrincipal ?? 5000,
    maxAmount: p.maxPrincipal ?? 3000000,
    eligibility: p.description || "Contact branch for eligibility details",
    description: p.description || p.name,
  }));
}

// ─── Apply for a Loan ────────────────────────────────────────────────────────

export async function applySelfLoan(productId: string, amount: number, termMonths: number, interestRate: number) {
  console.log('[Loan Application] Starting application:', { productId, amount, termMonths, interestRate });
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  
  // Deduct processing fee from disbursement amount
  // Fee structure: 500 KES total (200 to microfinance, 300 to SACCO)
  // Applicant receives (amount - 500) but repays full amount
  const disbursementAmount = amount - config.loan.processingFee;
  
  if (disbursementAmount <= 0) {
    throw new Error(`Loan amount must be at least ${config.loan.processingFee} KES to cover processing fee`);
  }
  
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "2-digit", year: "numeric" }).replace(/\//g, " ");
  const payload = {
    clientId,
    productId: Number(productId),
    principal: amount, // Full amount for repayment
    loanTermFrequency: termMonths,
    loanTermFrequencyType: 2,
    numberOfRepayments: termMonths,
    repaymentEvery: 1,
    repaymentFrequencyType: 2,
    interestRatePerPeriod: interestRate,
    amortizationType: 1,
    interestType: 0,
    interestCalculationPeriodType: 1,
    transactionProcessingStrategyCode: "mifos-standard-strategy",
    submittedOnDate: today,
    expectedDisbursementDate: today,
    dateFormat: "dd MM yyyy",
    locale: "en",
    loanType: 1, // 1 = Individual loan (mandatory parameter in Fineract)
    // Custom fields for fee tracking
    disbursementAmount, // Actual amount to be disbursed
    processingFee: config.loan.processingFee,
    microfinanceCut: config.loan.microfinanceCut,
    saccoCut: config.loan.saccoCut,
  };
  console.log('[Loan Application] Sending to Fineract:', payload);
  console.log('[Loan Application] Fee breakdown:', {
    requestedAmount: amount,
    processingFee: config.loan.processingFee,
    disbursementAmount,
    microfinanceCut: config.loan.microfinanceCut,
    saccoCut: config.loan.saccoCut,
  });
  const result = await selfPost("/loans", payload);
  console.log('[Loan Application] Success:', result);
  return result;
}

// ─── Account Lookup (for P2P transfers) ───────────────────────────────────────────

export async function lookupAccountByNumber(accountNo: string) {
  const data = await selfGet(`/savingsaccounts`);
  const items: any[] = data.pageItems ?? data.savingsAccounts ?? [];
  const account = items.find((a: any) => a.accountNo === accountNo);
  if (!account) throw new Error("Account not found");
  return {
    id: String(account.id),
    accountNo: account.accountNo,
    clientName: account.clientName,
    clientId: account.clientId ? String(account.clientId) : null,
    type: mapSavingsType(account.depositType?.value || account.savingsProductName || "Savings"),
  };
}

// ─── Real-Time Balance (Bypass Cache) ────────────────────────────────────────────

export async function getRealTimeBalance(accountId: string) {
  // Direct fetch to bypass React Query cache for critical balance checks
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  const res = await fetch(`${BASE}/savingsaccounts/${accountId}`, {
    headers: authHeaders(),
    cache: "no-store", // Bypass browser cache
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.defaultUserMessage || `Error ${res.status}`);
  }
  const data = await res.json();
  return {
    balance: data.summary?.accountBalance ?? 0,
    available: data.summary?.availableBalance ?? data.summary?.accountBalance ?? 0,
  };
}

// ─── P2P Transfers (Internal Ledger) ───────────────────────────────────────────

export async function transferToMember(fromAccountId: string, toAccountId: string, amount: number, description?: string, toClientId?: string | null) {
  try {
    console.log('[P2P Transfer] Starting transfer:', { fromAccountId, toAccountId, amount, description });
    
    const clientId = getClientId();
    if (!clientId) throw new Error("Not authenticated");
    
    // Validate inputs
    if (!fromAccountId || !toAccountId) {
      throw new Error('From account and To account are required');
    }
    if (!amount || amount <= 0) {
      throw new Error('Transfer amount must be greater than 0');
    }
    if (fromAccountId === toAccountId) {
      throw new Error('Cannot transfer to the same account');
    }
    
    const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "long", year: "numeric" });
    // Produces e.g. "22 June 2026" — matches Fineract dateFormat "dd MMMM yyyy"
    
    console.log('[P2P Transfer] Sending request to Fineract:', { 
      fromAccountId: Number(fromAccountId), 
      toAccountId: Number(toAccountId), 
      fromClientId: Number(clientId),
      toClientId,
      amount, 
      transferDate: today 
    });
    
    const result = await selfPost("/accounttransfers", {
      fromOfficeId: 1,
      fromClientId: Number(clientId),
      fromAccountType: 2, // 2 = Savings Portfolio
      fromAccountId: Number(fromAccountId),
      toOfficeId: 1,
      toClientId: toClientId ? Number(toClientId) : undefined,
      toAccountType: 2, // 2 = Savings Portfolio
      toAccountId: Number(toAccountId),
      transferDate: today,
      transferAmount: amount,
      transferDescription: description || "P2P Transfer",
      dateFormat: "dd MMMM yyyy",
      locale: "en",
    });
    
    console.log('[P2P Transfer] Transfer successful:', result);
    return result;
  } catch (error: any) {
    console.error('[P2P Transfer] Error:', error);
    if (error.response) {
      console.error('[P2P Transfer] API response:', error.response.status, error.response.data);
    }
    throw new Error(error.message || 'Transfer failed. Please try again.');
  }
}

// ─── Virtual Card Details (Member Only) ───────────────────────────────────────

export async function getMemberCardDetails(clientId: number) {
  const notes = await selfGet(`/clients/${clientId}/notes`);
  const items: any[] = Array.isArray(notes) ? notes : (notes.pageItems ?? []);
  const cardNote = items
    .map((n: any) => { try { return JSON.parse(n.note); } catch { return null; } })
    .filter((d: any) => d && d.cardNumber)
    .pop();
  if (!cardNote) return null;
  return {
    cardNumber: cardNote.cardNumber,
    expiryDate: cardNote.expiryDate,
    cvv: cardNote.cvv,
    cardType: cardNote.cardType || "Virtual",
    status: cardNote.status || "Active",
    issuedAt: cardNote.issuedAt,
    dailyLimit: cardNote.dailyLimit,
    monthlyLimit: cardNote.monthlyLimit,
  };
}

// ─── M-Pesa STK Push (loan repayment / deposit) ───────────────────────────────

export async function stkPush(msisdn: string, amount: number, accountRef: string) {
  const res = await fetch(`${MPESA_BASE}/stk/push`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ msisdn: msisdn.replace(/\s/g, ""), amount, account_ref: accountRef, description: `Payment for ${accountRef}` }),
  });
  if (!res.ok) throw new Error("M-Pesa STK push failed");
  return res.json();
}

// ─── M-Pesa B2C Withdrawal ────────────────────────────────────────────────────

export async function withdrawViaMpesa(accountId: string, amount: number, msisdn: string) {
  const clientId = getClientId();
  const res = await fetch(`${MPESA_BASE}/b2c/withdraw`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ account_id: accountId, client_id: String(clientId), msisdn: msisdn.replace(/\s/g, ""), amount }),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.message || "Withdrawal failed");

  // Record the withdrawal in Fineract as a transaction
  try {
    const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
    const transactionRes = await fetch(`${BASE}/savingsaccounts/${accountId}/transactions`, {
      method: "POST",
      headers: authHeaders(),
      body: JSON.stringify({
        dateFormat: "dd MMM yyyy",
        locale: "en",
        transactionDate: today,
        transactionAmount: amount,
        paymentTypeId: 1, // Cash payment
        transactionType: "withdrawal",
        note: `M-Pesa withdrawal to ${msisdn}. Ref: ${data.transaction_id || data.ConversationID || "N/A"}`,
      }),
    });

    if (!transactionRes.ok) {
      const err = await transactionRes.json().catch(() => ({}));
      console.error("Failed to record M-Pesa withdrawal in Fineract:", err);
      // Don't throw - the withdrawal already succeeded, just log the error
    } else {
      console.log("M-Pesa withdrawal recorded in Fineract");
    }
  } catch (error) {
    console.error("Error recording M-Pesa withdrawal:", error);
  }

  return data;
}

// ─── Transaction status ───────────────────────────────────────────────────────

export async function getTxStatus(txId: string) {
  const res = await fetch(`${MPESA_BASE}/tx/status/${txId}`);
  return res.json();
}

// ─── Loan Template (for dynamic form generation) ───────────────────────────────

export async function getLoanTemplate(productId?: string) {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  let params = `activeOnly=true&staffInSelectedOfficeOnly=true&clientId=${clientId}&templateType=individual`;
  if (productId) params += `&productId=${productId}`;
  return selfGet(`/loans/template?${params}`);
}

export async function getCollateralTemplate() {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  return selfGet(`/clients/${clientId}/collaterals/template`);
}

// ─── Loan Officers ───────────────────────────────────────────────────────────────

export async function getLoanOfficers() {
  const data = await selfGet("/staff?officeId=1");
  return (Array.isArray(data) ? data : (data.pageItems ?? [])).map((s: any) => ({
    id: String(s.id),
    name: `${s.firstname} ${s.lastname}`,
    office: s.officeName,
  }));
}

// ─── Fund Options ───────────────────────────────────────────────────────────────

export async function getFundOptions() {
  const data = await selfGet("/funds");
  return (Array.isArray(data) ? data : []).map((f: any) => ({
    id: String(f.id),
    name: f.name,
  }));
}

// ─── Interest Calculation Methods ───────────────────────────────────────────────

export async function getInterestCalculationMethods() {
  return [
    { id: 1, name: "Declining Balance" },
    { id: 2, name: "Flat" },
    { id: 3, name: "Declining Balance (Equal Installments)" },
    { id: 4, name: "Declining Balance (Equal Principal Payments)" },
  ];
}

// ─── Amortization Types ─────────────────────────────────────────────────────────

export async function getAmortizationTypes() {
  return [
    { id: 1, name: "Equal Principal Payments" },
    { id: 0, name: "Equal Installments" },
  ];
}

// ─── Repayment Strategies ───────────────────────────────────────────────────────

export async function getRepaymentStrategies() {
  return [
    { id: 1, name: "Whole Term" },
    { id: 2, name: "Same as Repayment Period" },
  ];
}

// ─── Client Documents (KYC) ─────────────────────────────────────────────────────

export async function getClientDocuments() {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  const data = await selfGet(`/clients/${clientId}/documents`);
  return (Array.isArray(data) ? data : []).map((d: any) => ({
    id: String(d.id),
    name: d.name,
    type: d.documentType?.name || "Document",
    description: d.description,
    fileName: d.fileName,
    uploadedOn: d.createdDate?.join?.("-") || "",
    status: d.status?.value || "Active",
  }));
}

export async function uploadClientDocument(file: File, documentTypeId: string, description: string) {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  
  const formData = new FormData();
  formData.append("file", file);
  formData.append("name", file.name);
  formData.append("documentTypeId", documentTypeId);
  formData.append("description", description);
  formData.append("locale", "en");
  formData.append("dateFormat", "dd MMM yyyy");
  formData.append("submittedOnDate", new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" }));

  const res = await fetch(`${BASE}/clients/${clientId}/documents`, {
    method: "POST",
    headers: authHeaders(),
    body: formData,
  });
  
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.defaultUserMessage || `Error ${res.status}`);
  }
  return res.json();
}

// ─── Document Types ─────────────────────────────────────────────────────────────

export async function getDocumentTypes() {
  const data = await selfGet("/documents");
  return (Array.isArray(data) ? data : []).map((d: any) => ({
    id: String(d.id),
    name: d.name,
    description: d.description,
  }));
}

// ─── Family Members ─────────────────────────────────────────────────────────────

export async function getFamilyMembers() {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  const data = await selfGet(`/clients/${clientId}/familymembers`);
  return (Array.isArray(data) ? data : []).map((f: any) => ({
    id: String(f.id),
    firstName: f.firstname,
    lastName: f.lastname,
    relation: f.relation?.name || "Family Member",
    profession: f.profession,
    age: f.age,
    isDependent: f.isDependent,
  }));
}

export async function addFamilyMember(member: any) {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  return selfPost(`/clients/${clientId}/familymembers`, {
    ...member,
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

// ─── Client Addresses ───────────────────────────────────────────────────────────

export async function getClientAddresses() {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  const data = await selfGet(`/clients/${clientId}/addresses`);
  return (Array.isArray(data) ? data : []).map((a: any) => ({
    id: String(a.id),
    addressType: a.addressType?.name || "Address",
    street: a.street,
    city: a.city,
    region: a.stateProvince,
    postalCode: a.postalCode,
    country: a.country?.name || "Kenya",
    isActive: a.isActive,
  }));
}

export async function addAddress(address: any) {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  return selfPost(`/clients/${clientId}/addresses`, {
    ...address,
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

// ─── Share Capital Management ─────────────────────────────────────────────────────

export async function getShareAccounts() {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  const data = await selfGet(`/savingsaccounts?clientId=${clientId}`);
  const items: any[] = data.pageItems ?? data.savingsAccounts ?? [];
  return items.filter((a: any) => a.savingsProductId === 1 || a.savingsProductName?.toLowerCase().includes("share")).map((a: any) => ({
    id: String(a.id),
    accountNo: a.accountNo,
    shares: a.summary?.accountBalance || 0,
    status: a.status?.value || "Active",
    openedOn: a.timeline?.activatedOnDate?.join?.("-") || "",
  }));
}

export async function purchaseShares(accountId: string, amount: number) {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  
  return selfPost(`/savingsaccounts/${accountId}/transactions`, {
    dateFormat: "dd MMM yyyy",
    locale: "en",
    transactionDate: today,
    transactionAmount: amount,
    paymentTypeId: 1, // Cash payment
    note: "Share purchase",
  });
}

export async function getDividends(accountId: string) {
  const data = await selfGet(`/savingsaccounts/${accountId}/transactions`);
  const items: any[] = data.pageItems ?? data ?? [];
  return items.filter((t: any) => t.transactionType?.value?.toLowerCase().includes("dividend")).map((t: any) => ({
    id: String(t.id),
    amount: t.amount,
    date: t.date?.join?.("-") || "",
    description: t.transactionType?.value || "Dividend",
  }));
}

// ─── Charges & Fees ─────────────────────────────────────────────────────────────

export async function getLoanCharges(loanId: string) {
  const data = await selfGet(`/loans/${loanId}/charges`);
  return (Array.isArray(data) ? data : []).map((c: any) => ({
    id: String(c.id),
    name: c.name,
    type: c.chargeCalculationType?.value || "Fixed",
    amount: c.amount,
    paid: c.amountPaid || 0,
    waived: c.amountWaived || 0,
    dueDate: c.dueDate?.join?.("-") || "",
    status: c.status?.value || "Pending",
  }));
}

export async function waiveCharge(chargeId: string, reason: string) {
  return selfPost(`/charges/${chargeId}/waive`, {
    reason,
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

export async function payCharge(chargeId: string, amount: number) {
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  return selfPost(`/charges/${chargeId}/payments`, {
    amount,
    paymentDate: today,
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

// ─── Reporting & Analytics ───────────────────────────────────────────────────────

export async function getLoanPortfolioReport() {
  const data = await selfGet("/loans");
  const items: any[] = data.pageItems ?? data.loans ?? [];
  return {
    totalLoans: items.length,
    totalPrincipal: items.reduce((sum: number, l: any) => sum + (l.principal || 0), 0),
    totalOutstanding: items.reduce((sum: number, l: any) => sum + (l.summary?.principalOutstanding || 0), 0),
    activeLoans: items.filter((l: any) => l.status?.value === "Active").length,
    overdueLoans: items.filter((l: any) => l.status?.value === "Overdue").length,
    byProduct: items.reduce((acc: any, l: any) => {
      const product = l.loanProductName || "Unknown";
      acc[product] = (acc[product] || 0) + 1;
      return acc;
    }, {}),
  };
}

export async function getDelinquencyReport() {
  const data = await selfGet("/loans");
  const items: any[] = data.pageItems ?? data.loans ?? [];
  return items.filter((l: any) => l.status?.value === "Overdue" || l.status?.value === "Arrears").map((l: any) => ({
    id: String(l.id),
    accountNo: l.accountNo,
    clientName: l.clientName,
    product: l.loanProductName,
    outstanding: l.summary?.principalOutstanding || 0,
    daysOverdue: l.summary?.overdueSinceDate ? Math.floor((Date.now() - new Date(l.summary.overdueSinceDate.join("-")).getTime()) / 86400000) : 0,
    nextPaymentDate: l.repaymentSchedule?.periods?.find((p: any) => !p.complete)?.dueDate?.join?.("-") || "",
  }));
}

export async function getFinancialStatements() {
  const data = await selfGet("/savingsaccounts");
  const items: any[] = data.pageItems ?? data.savingsAccounts ?? [];
  return {
    totalSavings: items.reduce((sum: number, a: any) => sum + (a.summary?.accountBalance || 0), 0),
    totalAccounts: items.length,
    activeAccounts: items.filter((a: any) => a.status?.value === "Active").length,
    byProduct: items.reduce((acc: any, a: any) => {
      const product = a.savingsProductName || "Unknown";
      acc[product] = (acc[product] || 0) + (a.summary?.accountBalance || 0);
      return acc;
    }, {}),
  };
}

// ─── Audit Trails ───────────────────────────────────────────────────────────────

export async function getAuditLogs(limit: number = 50) {
  const data = await selfGet(`/audit?limit=${limit}`);
  return (Array.isArray(data) ? data : []).map((log: any) => ({
    id: String(log.id),
    action: log.actionName,
    entity: log.entityName,
    resourceId: log.resourceId,
    madeBy: log.madeByUserName,
    madeOn: log.madeOnDate?.join?.("-") || "",
    details: log.details,
  }));
}

// ─── Loan Workflow (Approval, Disbursement, Repayment, Closure) ─────────────────

export async function approveLoan(loanId: string, note: string = "Approved") {
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  return selfPost(`/loans/${loanId}?command=approve`, {
    approvedOnDate: today,
    expectedDisbursementDate: today,
    note,
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

export async function disburseLoan(loanId: string, amount: number, paymentTypeId: number = 1) {
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  return selfPost(`/loans/${loanId}/transactions?command=disburse`, {
    transactionDate: today,
    transactionAmount: amount,
    paymentTypeId,
    note: "Loan disbursement",
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

export async function makeRepayment(loanId: string, amount: number, paymentTypeId: number = 1) {
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  return selfPost(`/loans/${loanId}/transactions?command=repayment`, {
    transactionDate: today,
    transactionAmount: amount,
    paymentTypeId,
    note: "Loan repayment",
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

export async function closeLoan(loanId: string, note: string = "Loan closed") {
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  return selfPost(`/loans/${loanId}?command=close`, {
    closedOnDate: today,
    note,
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

export async function getLoanStatus(loanId: string) {
  const data = await selfGet(`/loans/${loanId}`);
  return {
    id: String(data.id),
    status: data.status?.value || "Unknown",
    stage: mapLoanStage(data.status?.value || ""),
    principal: data.principal || 0,
    outstanding: data.summary?.principalOutstanding || 0,
    paid: data.summary?.principalPaid || 0,
    approvedDate: data.timeline?.approvedOnDate?.join?.("-") || "",
    disbursedDate: data.timeline?.actualDisbursementDate?.join?.("-") || "",
    closedDate: data.timeline?.closedOnDate?.join?.("-") || "",
  };
}

// ─── Savings Products (Fixed Deposits, Recurring Deposits, Target Savings) ─────────

export async function getSavingsProducts() {
  const data = await selfGet("/savingsproducts");
  return (Array.isArray(data) ? data : []).map((p: any) => ({
    id: String(p.id),
    name: p.name,
    type: p.depositType?.value || "Savings",
    minBalance: p.minRequiredOpeningBalance || 0,
    interestRate: p.nominalAnnualInterestRate || 0,
    currency: p.currency?.code || "KES",
    description: p.description || "",
  }));
}

// ─── Statement Generation ───────────────────────────────────────────────────────────

export async function getStatementData(accountId: string, startDate: string, endDate: string) {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  
  // Fetch account details + transactions in one call.
  // Fineract's /savingsaccounts/{id}/transactions endpoint rejects GET (405),
  // so transactions must be loaded via the associations query param instead.
  const accountData = await selfGet(`/savingsaccounts/${accountId}?associations=transactions`);
  const allTx: any[] = accountData.transactions ?? [];
  
  const start = new Date(startDate);
  const end = new Date(endDate);
  
  // Filter transactions by date range on client side
  const periodTransactions = allTx.filter((t: any) => {
    const txDate = new Date(t.date?.join?.("-") || "");
    return txDate >= start && txDate <= end;
  });
  
  // Calculate opening and closing balances
  const openingBalance = allTx
    .filter((t: any) => {
      const txDate = new Date(t.date?.join?.("-") || "");
      return txDate < start;
    })
    .reduce((sum: number, t: any) => sum + (t.transactionType?.withdrawal ? -t.amount : t.amount), 0);
  
  const closingBalance = allTx
    .filter((t: any) => {
      const txDate = new Date(t.date?.join?.("-") || "");
      return txDate <= end;
    })
    .reduce((sum: number, t: any) => sum + (t.transactionType?.withdrawal ? -t.amount : t.amount), 0);
  
  // Calculate totals for the period
  const totalDeposits = periodTransactions
    .filter((t: any) => !t.transactionType?.withdrawal)
    .reduce((sum: number, t: any) => sum + t.amount, 0);
  
  const totalWithdrawals = periodTransactions
    .filter((t: any) => t.transactionType?.withdrawal)
    .reduce((sum: number, t: any) => sum + t.amount, 0);
  
  return {
    account: {
      accountNo: accountData.accountNo,
      accountName: accountData.savingsProductName,
      balance: accountData.summary?.accountBalance || 0,
      currency: accountData.currency?.code || "KES",
    },
    period: {
      startDate,
      endDate,
    },
    summary: {
      openingBalance,
      closingBalance,
      totalDeposits,
      totalWithdrawals,
      netChange: totalDeposits - totalWithdrawals,
    },
    transactions: periodTransactions.map((t: any) => ({
      id: String(t.id),
      date: t.date?.join?.("-") || "",
      description: t.transactionType?.value || "Transaction",
      amount: t.transactionType?.withdrawal ? -t.amount : t.amount,
      balance: 0, // Would need running balance calculation
      reference: t.paymentDetailData?.receiptNumber || "",
    })),
  };
}

export async function getAvailableStatements(accountId: string) {
  // Generate list of available monthly statements for the last 12 months
  const statements = [];
  const now = new Date();
  
  for (let i = 0; i < 12; i++) {
    const date = new Date(now.getFullYear(), now.getMonth() - i, 1);
    const year = date.getFullYear();
    const month = date.getMonth() + 1;
    const monthName = date.toLocaleDateString("en-US", { month: "long" });
    
    const startDate = `${year}-${String(month).padStart(2, "0")}-01`;
    const endDate = `${year}-${String(month).padStart(2, "0")}-${new Date(year, month, 0).getDate()}`;
    
    statements.push({
      id: `${year}-${month}`,
      label: `${monthName} ${year}`,
      startDate,
      endDate,
      accountId,
    });
  }
  
  return statements;
}

export async function openSavingsAccount(productId: string, initialDeposit: number) {
  const clientId = getClientId();
  if (!clientId) throw new Error("Not authenticated");
  
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  
  return selfPost("/savingsaccounts", {
    clientId: Number(clientId),
    productId: Number(productId),
    submittedOnDate: today,
    externalId: `SA-${Date.now()}`,
    nominalAnnualInterestRate: 5,
    minRequiredOpeningBalance: initialDeposit,
    accountNo: null,
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

export async function makeDeposit(accountId: string, amount: number, paymentTypeId: number = 1) {
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  return selfPost(`/savingsaccounts/${accountId}/transactions?command=deposit`, {
    transactionDate: today,
    transactionAmount: amount,
    paymentTypeId,
    note: "Deposit",
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

export async function makeWithdrawal(accountId: string, amount: number, paymentTypeId: number = 1) {
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  return selfPost(`/savingsaccounts/${accountId}/transactions?command=withdrawal`, {
    transactionDate: today,
    transactionAmount: amount,
    paymentTypeId,
    note: "Withdrawal",
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}

export async function closeSavingsAccount(accountId: string) {
  const today = new Date().toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
  return selfPost(`/savingsaccounts/${accountId}?command=close`, {
    closedOnDate: today,
    note: "Account closed",
    locale: "en",
    dateFormat: "dd MMM yyyy",
  });
}
