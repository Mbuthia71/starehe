// Account Number Generation System
// Generates unique account numbers using configurable sequences

import { config } from "./config";

type AccountType = "savings" | "loans" | "shareCapital" | "fixedDeposit";

// In-memory sequence storage (in production, this should be stored in database)
const sequences: Record<AccountType, number> = {
  savings: config.accountNumbers.savings.startSequence,
  loans: config.accountNumbers.loans.startSequence,
  shareCapital: config.accountNumbers.shareCapital.startSequence,
  fixedDeposit: config.accountNumbers.fixedDeposit.startSequence,
};

// Generated account numbers to prevent duplicates
const generatedNumbers = new Set<string>();

/**
 * Generate a unique account number for a given account type
 * @param type - The type of account to generate a number for
 * @returns A unique account number
 */
export function generateAccountNumber(type: AccountType): string {
  const configType = config.accountNumbers[type];
  let accountNumber: string;
  let attempts = 0;
  const maxAttempts = 1000;

  do {
    // Increment sequence
    sequences[type]++;
    
    // Format account number
    accountNumber = configType.format(sequences[type]);
    
    attempts++;
    
    if (attempts >= maxAttempts) {
      throw new Error(`Failed to generate unique account number for ${type} after ${maxAttempts} attempts`);
    }
  } while (generatedNumbers.has(accountNumber));

  // Mark as generated
  generatedNumbers.add(accountNumber);

  return accountNumber;
}

/**
 * Get the current sequence number for an account type
 * @param type - The type of account
 * @returns The current sequence number
 */
export function getCurrentSequence(type: AccountType): number {
  return sequences[type];
}

/**
 * Set the sequence number for an account type (useful for initialization)
 * @param type - The type of account
 * @param sequence - The sequence number to set
 */
export function setSequence(type: AccountType, sequence: number): void {
  sequences[type] = sequence;
}

/**
 * Reset all sequences to their starting values (use with caution)
 */
export function resetSequences(): void {
  sequences.savings = config.accountNumbers.savings.startSequence;
  sequences.loans = config.accountNumbers.loans.startSequence;
  sequences.shareCapital = config.accountNumbers.shareCapital.startSequence;
  sequences.fixedDeposit = config.accountNumbers.fixedDeposit.startSequence;
  generatedNumbers.clear();
}

/**
 * Validate if an account number matches the expected format for a type
 * @param accountNumber - The account number to validate
 * @param type - The expected account type
 * @returns True if the account number matches the format
 */
export function validateAccountNumberFormat(accountNumber: string, type: AccountType): boolean {
  const configType = config.accountNumbers[type];
  return accountNumber.startsWith(configType.prefix) && accountNumber.length === 9;
}

/**
 * Extract the account type from an account number
 * @param accountNumber - The account number
 * @returns The account type or null if invalid
 */
export function getAccountTypeFromNumber(accountNumber: string): AccountType | null {
  for (const [type, configType] of Object.entries(config.accountNumbers) as [AccountType, any][]) {
    if (accountNumber.startsWith(configType.prefix)) {
      return type;
    }
  }
  return null;
}

/**
 * Get the sequence number from an account number
 * @param accountNumber - The account number
 * @returns The sequence number or null if invalid
 */
export function getSequenceFromNumber(accountNumber: string): number | null {
  const type = getAccountTypeFromNumber(accountNumber);
  if (!type) return null;
  
  const configType = config.accountNumbers[type];
  const sequenceStr = accountNumber.replace(configType.prefix, "");
  const sequence = parseInt(sequenceStr, 10);
  
  return isNaN(sequence) ? null : sequence;
}
