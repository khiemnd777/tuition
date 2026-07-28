import { describe, expect, it } from "vitest";

import { COFFEE_TRANSFER, generateCoffeeVietQR } from "../src/coffee.js";

describe("coffee support VietQR", () => {
  it("uses the approved VPBank recipient and 30,000 VND default", () => {
    expect(COFFEE_TRANSFER).toMatchObject({
      bankBin: "970432",
      bankName: "VPBank",
      accountNumber: "0974322365",
      accountName: "NGUYỄN DUY KHIÊM",
      purpose: "KHAO KHIEM CAFE",
      defaultAmount: 30000,
    });
    expect(generateCoffeeVietQR()).toEqual({
      amount: 30000,
      payload: "00020101021238540010A00000072701240006970432011009743223650208QRIBFTTA53037045405300005802VN62190815KHAO KHIEM CAFE63041BB7",
    });
  });

  it("creates a static QR without an amount for the voluntary option", () => {
    expect(generateCoffeeVietQR(0)).toEqual({
      amount: 0,
      payload: "00020101021138540010A00000072701240006970432011009743223650208QRIBFTTA53037045802VN62190815KHAO KHIEM CAFE6304BF41",
    });
  });
});
