import JSZip from "jszip";
import { describe, expect, it } from "vitest";

import { DEFAULT_EMAIL_TEMPLATE } from "../src/email.js";
import { createEmailBundle, createQRBundle, safeFilename } from "../src/exporter.js";

const validItem = {
  id: "row-001",
  sourceRow: 2,
  studentName: "Nguyễn An",
  parentName: "Nguyễn Văn Bình",
  className: "3.02",
  bankBin: "970415",
  bankName: "VietinBank",
  bankAccount: "0011001932418",
  email: "parent@example.com",
  amount: 120000,
  paymentItems: [{ label: "Học phí", labelEn: "Tuition", amount: 120000 }],
  billNumber: "SUN001",
  note: "HP NGUYEN AN",
  vietqr: "000201-test",
  qrData: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAAB",
  errors: [],
};

describe("bulk export bundles", () => {
  it("sanitizes filenames deterministically", () => {
    expect(safeFilename("SUN001 Nguyễn/An")).toBe("SUN001-Nguyen-An");
  });

  it("exports valid QR images plus manifest and row errors", async () => {
    const blob = await createQRBundle([
      validItem,
      { ...validItem, id: "row-002", billNumber: "SUN002", qrData: "", errors: ["Thiếu tài khoản"] },
    ]);
    const zip = await JSZip.loadAsync(await blob.arrayBuffer());
    expect(Object.keys(zip.files)).toContain("qr/SUN001-Nguyen-An.png");
    expect(await zip.file("manifest.csv").async("string")).toContain("SUN001");
    expect(await zip.file("errors.csv").async("string")).toContain("Thiếu tài khoản");
  });

  it("exports EML, provider JSONL and Gmail workbook without sending", async () => {
    const blob = await createEmailBundle([validItem], DEFAULT_EMAIL_TEMPLATE);
    const zip = await JSZip.loadAsync(await blob.arrayBuffer());
    const names = Object.keys(zip.files);
    expect(names).toContain("messages/SUN001-Nguyen-An.eml");
    expect(names).toContain("recipients/gmail-mail-merge.xlsx");
    expect(names).toContain("recipients/bulk-email.jsonl");
    expect(await zip.file("messages/SUN001-Nguyen-An.eml").async("string")).toContain("X-Unsent: 1");
    expect(await zip.file("README.txt").async("string")).toContain("không gửi email");
  });
});
