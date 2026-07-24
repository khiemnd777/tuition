import { describe, expect, it } from "vitest";

import {
  IMPORT_FIELD_GROUPS,
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
  it("groups concise fields and keeps legacy fee presets out of the visible list", () => {
    expect(IMPORT_FIELD_GROUPS.map((group) => group.label)).toEqual([
      "Học sinh & trường",
      "Phụ huynh",
      "Thanh toán",
      "Khoản phí / nâng cao",
    ]);
    const visibleKeys = IMPORT_FIELD_GROUPS.flatMap((group) => group.fields.map((field) => field.key));
    expect(visibleKeys).toEqual(expect.arrayContaining([
      "student_code", "student_name", "school_name", "cohort", "year", "class_name",
      "parent_name", "email", "bank_bin", "bank_account", "amount", "bill_number", "note",
      "fee_item", "payment_items",
    ]));
    expect(visibleKeys).not.toContain("tuition_april");
  });

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

  it("suggests the new school metadata and treats legacy fee headers as generic fees", () => {
    const mapping = suggestMapping([
      "Mã học sinh", "Tên trường", "Niên khóa", "Năm", "Lớp", "Học phí tháng 04",
    ]);
    expect(mapping).toEqual({
      "Mã học sinh": "student_code",
      "Tên trường": "school_name",
      "Niên khóa": "cohort",
      "Năm": "year",
      "Lớp": "class_name",
      "Học phí tháng 04": "fee_item",
    });
  });

  it("keeps concise student and school metadata on each payment row", () => {
    const metadataTable = {
      headers: ["Mã học sinh", "Họ và tên", "Tên trường", "Niên khóa", "Năm", "Lớp"],
      records: [["HS001", "Nguyễn An", "DEKISUGI", "2024–2028", "Năm 3", "3.02"]],
    };
    const rows = buildPaymentRows(metadataTable, suggestMapping(metadataTable.headers), {
      bankBin: "970415",
      bankAccount: "0011001932418",
    });
    expect(rows[0]).toMatchObject({
      studentCode: "HS001",
      studentName: "Nguyễn An",
      schoolName: "DEKISUGI",
      cohort: "2024–2028",
      year: "Năm 3",
      className: "3.02",
    });
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
