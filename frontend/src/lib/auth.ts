// Auth integration with Go backend for Starehian Society Platform

import { API_CONFIG } from "./api";
import { logger } from "./logger";

const TOKEN_KEY = "oss_token";
const REFRESH_TOKEN_KEY = "oss_refresh_token";
const USER_ID_KEY = "oss_user_id";
const ROLE_KEY = "oss_role";
const SESSION_TIMEOUT_KEY = "oss_last_activity";
const SESSION_TIMEOUT_MS = 60 * 60 * 1000; // 1 hour idle

export type UserRole = "super_admin" | "moderator" | "support" | "member";

export interface AuthUser {
  id: string;
  phone: string;
  email?: string;
  role: UserRole;
  status: string;
  created_at: string;
}

export interface AuthMember {
  id: string;
  displayName: string;
  firstName?: string;
  lastName?: string;
  emailAddress?: string;
  role: UserRole;
}

export interface AuthState {
  token: string;
  refreshToken: string;
  userId: string;
  role: UserRole;
}

const isBrowser = () => typeof window !== "undefined" && typeof localStorage !== "undefined";

export function getToken(): string | null {
  if (!isBrowser()) return null;
  return localStorage.getItem(TOKEN_KEY);
}

export function getRefreshToken(): string | null {
  if (!isBrowser()) return null;
  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function getUserId(): string | null {
  if (!isBrowser()) return null;
  return localStorage.getItem(USER_ID_KEY);
}

export function getRole(): UserRole {
  if (!isBrowser()) return "member";
  const v = localStorage.getItem(ROLE_KEY);
  return (v as UserRole) || "member";
}

export function isAuthenticated(): boolean {
  return !!getToken() && !!getUserId();
}

export function saveAuth(state: AuthState) {
  localStorage.setItem(TOKEN_KEY, state.token);
  localStorage.setItem(REFRESH_TOKEN_KEY, state.refreshToken);
  localStorage.setItem(USER_ID_KEY, state.userId);
  localStorage.setItem(ROLE_KEY, state.role);
}

export function getStoredMember(): AuthMember | null {
  if (!isBrowser()) return null;
  const userId = getUserId();
  const role = getRole();
  if (!userId) return null;
  
  // For now, return a minimal member object from localStorage
  // In the future, this should fetch from the API
  const stored = localStorage.getItem("oss_member_data");
  if (stored) {
    return JSON.parse(stored) as AuthMember;
  }
  
  // Fallback to basic info
  return {
    id: userId,
    displayName: "Old Starehian",
    role,
  };
}

export function clearAuth() {
  if (!isBrowser()) return;
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem(USER_ID_KEY);
  localStorage.removeItem(ROLE_KEY);
}

// ─── Role-based permissions ──────────────────────────────────────────────
export const permissions = {
  super_admin: ["read", "write", "delete", "approve", "manage_users", "view_all", "broadcast", "bulk"],
  moderator: ["read", "write", "approve", "moderate"],
  support: ["read", "write", "support"],
  member: ["read_own", "write_own"],
} as const;

export function hasPermission(action: string): boolean {
  const role = getRole();
  return (permissions[role as keyof typeof permissions] as readonly string[]).includes(action) || false;
}

export function canApproveConnections(): boolean { return hasPermission("approve"); }
export function canManageUsers(): boolean { return hasPermission("manage_users"); }
export function canViewAllData(): boolean { return hasPermission("view_all"); }
export function canBroadcast(): boolean { return hasPermission("broadcast"); }
export function canBulkOperations(): boolean { return hasPermission("bulk"); }

// ─── API calls to Go backend ──────────────────────────────────────────────
async function apiCall<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${API_CONFIG.baseUrl}${endpoint}`;
  const token = getToken();
  
  const response = await fetch(url, {
    ...options,
    headers: {
      ...API_CONFIG.headers(token || undefined),
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: "Request failed" }));
    throw new Error(error.error || `HTTP ${response.status}`);
  }

  return response.json();
}

// Request OTP for phone number
export async function requestOTP(phone: string): Promise<void> {
  const formData = new URLSearchParams();
  formData.append("phone", phone);
  
  await apiCall("/auth/request-otp", {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: formData.toString(),
  });
  logger.info(`OTP requested for: ${phone}`);
}

// Sign up with phone + password + full name
export async function signup(phone: string, password: string, full_name: string): Promise<AuthState> {
  const response = await apiCall<{ tokens: { access_token: string; refresh_token: string }; user: AuthUser }>("/auth/signup", {
    method: "POST",
    body: JSON.stringify({ phone, password, full_name }),
  });

  const state: AuthState = {
    token: response.tokens.access_token,
    refreshToken: response.tokens.refresh_token,
    userId: response.user.id,
    role: response.user.role as UserRole,
  };
  saveAuth(state);
  logger.info(`User signed up: ${response.user.id}`);
  return state;
}

// Login with OTP
export async function loginWithOTP(phone: string, otp: string): Promise<AuthState> {
  const response = await apiCall<{ tokens: { access_token: string; refresh_token: string }; user: AuthUser }>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ phone, otp }),
  });

  const state: AuthState = {
    token: response.tokens.access_token,
    refreshToken: response.tokens.refresh_token,
    userId: response.user.id,
    role: response.user.role as UserRole,
  };
  saveAuth(state);
  logger.info(`User logged in: ${response.user.id}`);
  return state;
}

// Admin login with email/password
export async function adminLogin(email: string, password: string): Promise<AuthState> {
  const response = await apiCall<{ tokens: { access_token: string; refresh_token: string }; user: AuthUser }>("/auth/admin/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });

  const state: AuthState = {
    token: response.tokens.access_token,
    refreshToken: response.tokens.refresh_token,
    userId: response.user.id,
    role: response.user.role as UserRole,
  };
  saveAuth(state);
  logger.info(`Admin logged in: ${response.user.id}`);
  return state;
}

// Refresh token
export async function refreshToken(): Promise<AuthState> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) throw new Error("No refresh token");

  const response = await apiCall<{ access_token: string; refresh_token: string }>("/auth/refresh", {
    method: "POST",
    body: JSON.stringify({ refresh_token: refreshToken }),
  });

  const state: AuthState = {
    token: response.access_token,
    refreshToken: response.refresh_token,
    userId: getUserId() || "",
    role: getRole(),
  };
  saveAuth(state);
  return state;
}

// Logout
export async function logout(): Promise<void> {
  try {
    await apiCall("/auth/logout", { method: "POST" });
  } catch (error: any) {
    logger.error("Logout API call failed:", error);
  }
  clearAuth();
  logger.info("User logged out");
}

// Login with phone or email + password
export async function login(identifier: string, password: string): Promise<AuthState> {
  // If identifier contains @, try admin login
  if (identifier.includes("@")) {
    return adminLogin(identifier, password);
  }
  
  // Otherwise, treat as phone and use regular login with password
  const response = await apiCall<{ tokens: { access_token: string; refresh_token: string }; user: AuthUser }>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ phone: identifier, password }),
  });

  const state: AuthState = {
    token: response.tokens.access_token,
    refreshToken: response.tokens.refresh_token,
    userId: response.user.id,
    role: response.user.role as UserRole,
  };
  saveAuth(state);
  logger.info(`User logged in: ${response.user.id}`);
  return state;
}

export function authHeaders(): Record<string, string> {
  const token = getToken();
  return {
    Authorization: `Bearer ${token ?? ""}`,
    "Content-Type": "application/json",
  };
}

// ─── Session / device helpers ────────────────────────────────────────────
export function updateActivity() {
  if (!isBrowser()) return;
  localStorage.setItem(SESSION_TIMEOUT_KEY, Date.now().toString());
}

export function isSessionExpired(): boolean {
  if (!isBrowser()) return false;
  const lastActivity = localStorage.getItem(SESSION_TIMEOUT_KEY);
  if (!lastActivity) return false;
  return Date.now() - Number(lastActivity) > SESSION_TIMEOUT_MS;
}

export function startSessionTimeoutCheck(callback: () => void) {
  const interval = setInterval(() => {
    if (isSessionExpired()) {
      clearInterval(interval);
      callback();
    }
  }, 30000);
  return interval;
}

export function getDeviceFingerprint(): string {
  if (typeof navigator === "undefined" || typeof window === "undefined") return "server";
  const ua = navigator.userAgent;
  const language = navigator.language;
  const platform = navigator.platform;
  const screenSize = `${window.screen.width}x${window.screen.height}`;
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  return btoa(`${ua}|${language}|${platform}|${screenSize}|${timezone}`);
}

export function isTrustedDevice(): boolean {
  if (!isBrowser()) return true;
  const fingerprint = getDeviceFingerprint();
  const trusted = localStorage.getItem("oss_trusted_devices");
  if (!trusted) return false;
  return (JSON.parse(trusted) as string[]).includes(fingerprint);
}

export function trustCurrentDevice() {
  if (!isBrowser()) return;
  const fingerprint = getDeviceFingerprint();
  const trusted = localStorage.getItem("oss_trusted_devices");
  const devices: string[] = trusted ? JSON.parse(trusted) : [];
  if (!devices.includes(fingerprint)) {
    devices.push(fingerprint);
    localStorage.setItem("oss_trusted_devices", JSON.stringify(devices));
  }
}

export function isBiometricAvailable(): boolean {
  return typeof window !== "undefined" &&
    typeof (window as any).PublicKeyCredential !== "undefined";
}

export function detectDeviceSecurity(): { isEmulator: boolean; isRooted: boolean; isJailbroken: boolean } {
  if (typeof navigator === "undefined") return { isEmulator: false, isRooted: false, isJailbroken: false };
  const ua = navigator.userAgent;
  const isEmulator = /Android SDK|Genymotion|Bluestacks|Nox|LDPlayer|Memu|Simulator|Emulator/i.test(ua);
  const isRooted = /Android/i.test(ua) && /root|superuser|magisk|test-keys/i.test(ua.toLowerCase());
  const isJailbroken = /iPhone|iPad|iPod/i.test(ua) && /Cydia|Sileo|checkra1n|unc0ver/i.test(ua);
  return { isEmulator, isRooted, isJailbroken };
}
