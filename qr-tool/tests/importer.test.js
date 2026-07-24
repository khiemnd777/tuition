import { describe, expect, it } from "vitest";

import {
  buildPaymentRows,
  normalizeHeader,
  suggestMapping,
  validateMapping,
} from "../src/importer.js";

const table = {
  headers: ["Họ và tên", "Email phụ huynh", "Tổng phí", "Học phí T7", "Phí xe T7"],
  records: [
    ["Nguyễn An", "parent@example.com", "120.000", "3.950.000", "3.030.000"],
  ],
};

describe("browser import mapping", () => {
  it("normalizes Vietnamese spreadsheet headers", () => {
    expect(normalizeHeader("  Số tài khoản  ")).toBe("so_tai_khoan");
  });

  it("suggests known aliases and leaves custom fee columns explicit", () => {
    const mapping = suggestMapping(table.headers);
    expect(mapping["Họ và tên"]).toBe("student_name");
    expect(mapping["Email phụ huynh"]).toBe("email");
    expect(mapping["Tổng phí"]).toBe("amount");
    expect(mapping["Học phí T7"]).toBe("");
  });

  it("maps multiple custom fee columns and lets their total override amount", () => {
    const rows = buildPaymentRows(table, {
      "Họ và tên": "student_name",
      "Email phụ huynh": "email",
      "Tổng phí": "amount",
      "Học phí T7": "fee_item",
      "Phí xe T7": "fee_item",
    }, {
      bankBin: "970415",
      bankAccount: "0011001932418",
    });

    expect(rows).toHaveLength(1);
    expect(rows[0].studentName).toBe("Nguyễn An");
    expect(rows[0].paymentItems).toEqual([
      { label: "Học phí T7", labelEn: "Học phí T7", amount: 3_950_000 },
      { label: "Phí xe T7", labelEn: "Phí xe T7", amount: 3_030_000 },
    ]);
    expect(rows[0].amount).toBe(6_980_000);
  });

  it("prevents duplicate singleton targets", () => {
    expect(validateMapping({ A: "email", B: "email", C: "fee_item", D: "fee_item" })).toEqual([
      "Trường Email đang được map từ nhiều cột",
    ]);
  });
});
