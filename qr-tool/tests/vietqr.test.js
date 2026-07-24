import { describe, expect, it } from "vitest";

import {
  buildQRItem,
  cleanPaymentRow,
  crc16CCITTFalse,
  generateVietQR,
} from "../src/vietqr.js";

describe("VietQR contract", () => {
  it("matches the documented CRC fixture", () => {
    const input = "00020101021138570010A00000072701270006970403011200110123456780208QRIBFTTA53037045802VN6304";
    expect(crc16CCITTFalse(input)).toBe("F4E5");
  });

  it("matches the existing static IBFT payload", () => {
    expect(generateVietQR({ bankBin: "970423", accountNumber: "0099999999" })).toBe(
      "00020101021138540010A00000072701240006970423011000999999990208QRIBFTTA53037045802VN6304CBB4",
    );
  });

  it("matches the existing dynamic IBFT payload", () => {
    expect(generateVietQR({
      bankBin: "970403",
      accountNumber: "0011012345678",
      amount: 180000,
      billNumber: "NPS6869",
      purpose: "thanh toan don hang",
      dynamic: true,
    })).toBe(
      "00020101021238570010A00000072701270006970403011300110123456780208QRIBFTTA530370454061800005802VN62340107NPS68690819thanh toan don hang63042E2E",
    );
  });

  it("uses payment item totals instead of raw amount", () => {
    const row = cleanPaymentRow({
      studentCode: "HS001",
      studentName: "Nguyen An",
      schoolName: "DEKISUGI School",
      cohort: "2024-2028",
      year: "Nam 3",
      bankBin: "970415",
      bankAccount: "0011001932418",
      amount: 1,
      paymentItems: [
        { label: "Học phí", amount: 3_950_000 },
        { label: "Phí xe", amount: 3_030_000 },
      ],
      billNumber: "SUN001",
      note: "Hoc phi",
    });
    expect(row.amount).toBe(6_980_000);
    expect(row).toMatchObject({
      studentCode: "HS001",
      schoolName: "DEKISUGI School",
      cohort: "2024-2028",
      year: "Nam 3",
    });
    expect(buildQRItem(row).vietqr).toContain("54076980000");
  });

  it("rejects unknown bank BINs without silently generating a QR", () => {
    const item = buildQRItem({ bankBin: "999999", bankAccount: "123", amount: 1000 });
    expect(item.errors).toContain("BIN ngân hàng không nằm trong danh sách VietQR");
    expect(item.vietqr).toBe("");
  });
});
