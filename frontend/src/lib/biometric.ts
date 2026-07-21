import { BiometricAuth } from '@aparajita/capacitor-biometric-auth';

export async function isBiometricAvailable(): Promise<boolean> {
  try {
    const result = await BiometricAuth.checkBiometry();
    return result.isAvailable;
  } catch {
    return false;
  }
}

export async function authenticateWithBiometric(reason: string): Promise<boolean> {
  try {
    await BiometricAuth.authenticate({
      reason,
      cancelTitle: 'Cancel',
      allowDeviceCredential: true,
    });
    return true;
  } catch {
    return false;
  }
}

const secretKey = (key: string) => `biometric_secret_${key}`;

export async function setBiometricSecret(key: string, value: string): Promise<void> {
  try {
    if (typeof localStorage !== 'undefined') {
      localStorage.setItem(secretKey(key), value);
    }
  } catch (error) {
    console.error('Failed to set biometric secret:', error);
  }
}

export async function getBiometricSecret(key: string): Promise<string | null> {
  try {
    if (typeof localStorage !== 'undefined') {
      return localStorage.getItem(secretKey(key));
    }
    return null;
  } catch {
    return null;
  }
}

export async function removeBiometricSecret(key: string): Promise<void> {
  try {
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem(secretKey(key));
    }
  } catch (error) {
    console.error('Failed to remove biometric secret:', error);
  }
}
