import { generateVietQR, parseAmount } from "./vietqr.js";

export const COFFEE_TRANSFER = Object.freeze({
  bankBin: "970432",
  bankName: "VPBank",
  accountNumber: "0974322365",
  accountName: "NGUYỄN DUY KHIÊM",
  purpose: "KHAO KHIEM CAFE",
  defaultAmount: 30000,
});

export function generateCoffeeVietQR(amount = COFFEE_TRANSFER.defaultAmount) {
  const selectedAmount = parseAmount(amount);
  return {
    amount: selectedAmount,
    payload: generateVietQR({
      bankBin: COFFEE_TRANSFER.bankBin,
      accountNumber: COFFEE_TRANSFER.accountNumber,
      amount: selectedAmount,
      purpose: COFFEE_TRANSFER.purpose,
    }),
  };
}
