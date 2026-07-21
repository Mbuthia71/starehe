// OTP Authentication System
import { authHeaders } from "./auth";
import { config } from "./config";
import { emailHelpers } from "./email";

const BASE = config.fineract.baseUrl;
const TENANT = config.fineract.tenantId;

export interface OtpRequest {
  msisdn: string;
}

export interface OtpVerify {
  msisdn: string;
  otp: string;
  deviceId: string;
}

export interface OtpResponse {
  success: boolean;
  message: string;
  tempToken?: string;
}

// Request OTP for login
export async function requestOtp(msisdn: string): Promise<OtpResponse> {
  try {
    const res = await fetch(`${BASE}/authentication`, {
      method: "POST",
      headers: {
        "Fineract-Platform-TenantId": TENANT,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: msisdn,
        password: "", // Will be set after OTP verification
      }),
    });

    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.defaultUserMessage || "Failed to request OTP");
    }

    // In a real implementation, this would call an SMS service
    // For now, we'll simulate OTP generation and send via email
    const otp = Math.floor(100000 + Math.random() * 900000).toString();
    console.log(`[OTP] Generated OTP for ${msisdn}: ${otp}`);

    // Send OTP via email
    try {
      await emailHelpers.sendOtp(msisdn, otp, "Member");
    } catch (emailError) {
      console.error("[OTP] Failed to send email:", emailError);
      // Continue anyway - OTP is still valid
    }

    // Store OTP in localStorage for verification (in production, this would be server-side)
    localStorage.setItem(`otp_${msisdn}`, JSON.stringify({
      otp,
      expiresAt: Date.now() + 5 * 60 * 1000, // 5 minutes
    }));

    return {
      success: true,
      message: "OTP sent successfully",
    };
  } catch (error: any) {
    console.error("[OTP] Request error:", error);
    throw new Error(error.message || "Failed to request OTP");
  }
}

// Verify OTP
export async function verifyOtp(msisdn: string, otp: string, deviceId: string): Promise<OtpResponse> {
  try {
    const stored = localStorage.getItem(`otp_${msisdn}`);
    if (!stored) {
      throw new Error("OTP expired or not found. Please request a new OTP.");
    }

    const { otp: storedOtp, expiresAt } = JSON.parse(stored);
    
    if (Date.now() > expiresAt) {
      localStorage.removeItem(`otp_${msisdn}`);
      throw new Error("OTP expired. Please request a new OTP.");
    }

    if (otp !== storedOtp) {
      throw new Error("Invalid OTP. Please try again.");
    }

    // Clear OTP after successful verification
    localStorage.removeItem(`otp_${msisdn}`);

    // Generate temporary token
    const tempToken = btoa(`${msisdn}:${Date.now()}`);

    return {
      success: true,
      message: "OTP verified successfully",
      tempToken,
    };
  } catch (error: any) {
    console.error("[OTP] Verify error:", error);
    throw new Error(error.message || "Failed to verify OTP");
  }
}

// Resend OTP
export async function resendOtp(msisdn: string): Promise<OtpResponse> {
  // Check if there's a recent OTP request to prevent spam
  const lastRequest = localStorage.getItem(`otp_last_${msisdn}`);
  if (lastRequest && Date.now() - Number(lastRequest) < 30000) {
    throw new Error("Please wait 30 seconds before requesting another OTP");
  }

  localStorage.setItem(`otp_last_${msisdn}`, Date.now().toString());
  return requestOtp(msisdn);
}
