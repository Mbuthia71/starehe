import type { StatementDocumentProps } from "./StatementDocument";

export const sampleStatement: StatementDocumentProps = {
  member: {
    fullName: "Achieng Otieno",
    memberNumber: "SHM-04821",
    msisdn: "+254 712 884 901",
    address: ["P.O. Box 4821 – 00100", "Kilimani, Nairobi", "Kenya"],
  },
  account: {
    productName: "Siohioma Savings Account",
    accountNumber: "0148-2210-09",
    subtitle: "MEMBER: ACHIENG OTIENO",
  },
  period: { from: "01 Oct 2026", to: "31 Oct 2026", label: "October 2026" },
  summary: {
    beginningBalance: 132_480.0,
    deposits: 18_000.0,
    mpesaIn: 42_500.0,
    otherCredits: 1_240.5,
    withdrawals: 12_500.0,
    loanRepayments: 8_400.0,
    endingBalance: 173_320.5,
    averageBalance: 154_220.18,
    interestEarnedPeriod: 642.18,
    interestYtd: 5_120.4,
    dividendsCredited: 0.0,
    apy: 0.0625,
    daysInPeriod: 31,
  },
  sections: [
    {
      title: "Member contributions",
      rows: [
        { date: "10/02", description: "MONTHLY SAVINGS · STANDING ORDER", amount: 6_000 },
        { date: "10/12", description: "MONTHLY SAVINGS · STANDING ORDER", amount: 6_000 },
        { date: "10/22", description: "MONTHLY SAVINGS · STANDING ORDER", amount: 6_000 },
      ],
    },
    {
      title: "M-Pesa deposits",
      rows: [
        { date: "10/03", description: "M-PESA DEPOSIT · +254712884901 · REF QJK7H2A", amount: 5_000 },
        { date: "10/07", description: "M-PESA DEPOSIT · +254712884901 · REF QJL9P12", amount: 12_000 },
        { date: "10/14", description: "M-PESA DEPOSIT · +254712884901 · REF QJM1B43", amount: 8_500 },
        { date: "10/19", description: "M-PESA DEPOSIT · +254712884901 · REF QJN3C77", amount: 7_000 },
        { date: "10/27", description: "M-PESA DEPOSIT · +254712884901 · REF QJP4D02", amount: 10_000 },
      ],
    },
    {
      title: "Withdrawals",
      rows: [
        { date: "10/09", description: "M-PESA WITHDRAWAL · +254712884901 · REF WD09128", amount: 7_500 },
        { date: "10/24", description: "M-PESA WITHDRAWAL · +254712884901 · REF WD24910", amount: 5_000 },
      ],
    },
    {
      title: "Loan activity",
      rows: [
        { date: "10/05", description: "EMERGENCY LOAN REPAYMENT · INSTALMENT 4/12", amount: 4_200 },
        { date: "10/25", description: "EMERGENCY LOAN REPAYMENT · INSTALMENT 5/12", amount: 4_200 },
      ],
    },
    {
      title: "Interest & other credits",
      rows: [
        { date: "10/31", description: "INTEREST CREDIT · SAVINGS ACCOUNT", amount: 642.18 },
        { date: "10/31", description: "FEE REVERSAL · M-PESA TRANSFER", amount: 598.32 },
      ],
    },
  ],
  page: { current: 1, total: 1 },
};
