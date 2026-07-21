// Audit logging for tracking admin actions
// Critical for banking compliance and security

import { getRole, getStoredMember } from "./auth";

export type AuditAction = 
  | "login"
  | "logout"
  | "loan_approve"
  | "loan_disburse"
  | "loan_reject"
  | "transfer_initiate"
  | "transfer_complete"
  | "card_issue"
  | "card_revoke"
  | "user_create"
  | "user_update"
  | "user_delete"
  | "data_view";

export interface AuditLog {
  id: string;
  timestamp: string;
  userId: number;
  username: string;
  role: string;
  action: AuditAction;
  resource: string;
  resourceId?: string;
  details: Record<string, any>;
  ipAddress?: string;
  userAgent?: string;
}

// In-memory storage for audit logs (in production, use database)
const auditLogs: AuditLog[] = [];

export function logAudit(action: AuditAction, resource: string, details: Record<string, any> = {}, resourceId?: string) {
  const member = getStoredMember();
  if (!member) return;

  const log: AuditLog = {
    id: crypto.randomUUID(),
    timestamp: new Date().toISOString(),
    userId: member.id,
    username: member.displayName,
    role: getRole(),
    action,
    resource,
    resourceId,
    details,
    ipAddress: details.ipAddress || "unknown",
    userAgent: navigator.userAgent,
  };

  auditLogs.push(log);
  console.log('[AUDIT]', log);

  // In production, send to backend API
  // await fetch('/api/audit', { method: 'POST', body: JSON.stringify(log) });
}

export function getAuditLogs(filters?: { userId?: number; action?: AuditAction; resource?: string }): AuditLog[] {
  let filtered = auditLogs;

  if (filters?.userId) {
    filtered = filtered.filter(log => log.userId === filters.userId);
  }
  if (filters?.action) {
    filtered = filtered.filter(log => log.action === filters.action);
  }
  if (filters?.resource) {
    filtered = filtered.filter(log => log.resource === filters.resource);
  }

  return filtered.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());
}

// Convenience functions for common audit actions
export function auditLogin() {
  logAudit("login", "auth", {});
}

export function auditLogout() {
  logAudit("logout", "auth", {});
}

export function auditLoanApprove(loanId: string, amount: number) {
  logAudit("loan_approve", "loan", { amount }, loanId);
}

export function auditLoanDisburse(loanId: string, amount: number, disbursementAmount: number) {
  logAudit("loan_disburse", "loan", { amount, disbursementAmount }, loanId);
}

export function auditLoanReject(loanId: string, reason: string) {
  logAudit("loan_reject", "loan", { reason }, loanId);
}

export function auditTransferInitiate(fromAccount: string, toAccount: string, amount: number) {
  logAudit("transfer_initiate", "transfer", { fromAccount, toAccount, amount });
}

export function auditCardIssue(clientId: number, cardNumber: string) {
  logAudit("card_issue", "card", { clientId, cardNumber });
}

export function auditDataView(resource: string, resourceId?: string) {
  logAudit("data_view", resource, {}, resourceId);
}
