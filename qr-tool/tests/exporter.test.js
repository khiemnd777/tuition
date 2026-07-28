import JSZip from "jszip";
import * as XLSX from "xlsx";
import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

import { DEFAULT_EMAIL_TEMPLATE } from "../src/email.js";
import {
  GMAIL_FREE_APPS_SCRIPT,
  createEmailBundle,
  createGmailFreeBundle,
  createGmailSheetData,
  createQRBundle,
  isGmailSheetTemplateURL,
  safeFilename,
} from "../src/exporter.js";

const validItem = {
  id: "row-001",
  sourceRow: 2,
  studentCode: "HS001",
  studentName: "Nguyễn An",
  schoolName: "DEKISUGI School",
  cohort: "2024–2028",
  year: "Năm 3",
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

  it("only accepts a Google Sheets copy URL for the published template", () => {
    expect(isGmailSheetTemplateURL("https://docs.google.com/spreadsheets/d/abc123/copy")).toBe(true);
    expect(isGmailSheetTemplateURL("https://docs.google.com/spreadsheets/d/abc123/edit")).toBe(false);
    expect(isGmailSheetTemplateURL("https://evil.example/spreadsheets/d/abc123/copy")).toBe(false);
  });

  it("exports valid QR images plus manifest and row errors", async () => {
    const blob = await createQRBundle([
      validItem,
      { ...validItem, id: "row-002", billNumber: "SUN002", qrData: "", errors: ["Thiếu tài khoản"] },
    ]);
    const zip = await JSZip.loadAsync(await blob.arrayBuffer());
    expect(Object.keys(zip.files)).toContain("qr/SUN001-Nguyen-An.png");
    const manifest = await zip.file("manifest.csv").async("string");
    expect(manifest).toContain("Mã học sinh");
    expect(manifest).toContain("HS001");
    expect(manifest).toContain("DEKISUGI School");
    expect(manifest).toContain("SUN001");
    expect(await zip.file("errors.csv").async("string")).toContain("Thiếu tài khoản");
  });

  it("exports EML and provider data without legacy Gmail Mail Merge artifacts", async () => {
    const blob = await createEmailBundle([validItem], DEFAULT_EMAIL_TEMPLATE);
    const zip = await JSZip.loadAsync(await blob.arrayBuffer());
    const names = Object.keys(zip.files);
    expect(names).toContain("messages/SUN001-Nguyen-An.eml");
    expect(names).not.toContain("recipients/gmail-mail-merge.xlsx");
    expect(names).not.toContain("template/gmail-template.html");
    expect(names).toContain("recipients/bulk-email.jsonl");
    expect(await zip.file("messages/SUN001-Nguyen-An.eml").async("string")).toContain("X-Unsent: 1");
    expect(await zip.file("recipients/bulk-email.csv").async("string")).toContain("student_code");
    expect(await zip.file("recipients/bulk-email.jsonl").async("string")).toContain('"schoolName":"DEKISUGI School"');
    expect(await zip.file("README.txt").async("string")).toContain("không gửi email");
  });

  it("exports a self-contained Gmail Free workbook, Apps Script and offline guide", async () => {
    const invalidItem = {
      ...validItem,
      id: "row-002",
      sourceRow: 3,
      email: "email-sai",
      qrData: "",
      errors: ["Thiếu số tài khoản"],
    };
    const blob = await createGmailFreeBundle([validItem, invalidItem], DEFAULT_EMAIL_TEMPLATE);
    const zip = await JSZip.loadAsync(await blob.arrayBuffer());

    expect(Object.keys(zip.files).sort()).toEqual([
      "01_DANH_SACH_GUI.xlsx",
      "02_CODE_GUI_EMAIL.gs",
      "BAT_DAU_TU_DAY.html",
    ]);

    const workbookData = await zip.file("01_DANH_SACH_GUI.xlsx").async("arraybuffer");
    const workbook = XLSX.read(workbookData, { type: "array" });
    expect(workbook.SheetNames).toEqual(["BAT_DAU", "EMAILS"]);
    const guideRows = XLSX.utils.sheet_to_json(workbook.Sheets.BAT_DAU, { header: 1, defval: "" });
    const workbookGuide = guideRows.flat().join(" ");
    expect(workbookGuide).toContain("Authorization required");
    expect(workbookGuide).toContain("Review permissions");
    expect(workbookGuide).toContain("đang hiện onOpen và chọn setup");
    expect(workbookGuide).toContain("Execution completed");
    expect(workbookGuide).toContain("Gửi thử dòng đang chọn");
    expect(workbookGuide).toContain("0. Nhập dữ liệu mới");
    expect(workbookGuide).toContain("Không dán code, không chạy setup");
    const rows = XLSX.utils.sheet_to_json(workbook.Sheets.EMAILS, { defval: "" });
    expect(rows).toHaveLength(2);
    expect(rows[0]).toMatchObject({
      Send: "YES",
      Email: "parent@example.com",
      StudentCode: "HS001",
      Status: "READY",
    });
    expect(rows[0].HtmlBody).toContain("cid:dekisugi_qr");
    expect(rows[0].QrBase64).not.toContain("data:image/png;base64,");
    expect(rows[1].Send).toBe("NO");
    expect(rows[1].Status).toBe("ERROR");
    expect(rows[1].Error).toContain("Email không hợp lệ");

    const script = await zip.file("02_CODE_GUI_EMAIL.gs").async("string");
    expect(script).toBe(GMAIL_FREE_APPS_SCRIPT);
    expect(() => new Function(script)).not.toThrow();
    expect(script).toContain("MailApp.getRemainingDailyQuota");
    expect(script).toContain("MAX_PER_BATCH: 90");
    expect(script).toContain("LockService.getUserLock");
    expect(script).toContain("inlineImages");
    expect(script).toContain("Utilities.base64Decode");
    expect(script).toContain("hideColumns");
    expect(script).toContain('addItem("0. Nhập dữ liệu mới", "showImportDialog")');
    expect(script).toContain("function importDataset(payload)");
    expect(script).toContain('payload.kind !== "dekisugi.gmail-data"');
    expect(script).toContain("safeSheetText_");
    const setupSource = script.match(/function setup\(\) \{([\s\S]*?)\n\}/)?.[1] || "";
    expect(setupSource).toContain("SpreadsheetApp.getActive().toast");
    expect(setupSource).toContain("console.log(message)");
    expect(setupSource).not.toContain("getUi().alert");
    expect(script).toContain("Xác nhận gửi email thật");
    expect(script).not.toContain("UrlFetchApp");
    expect(script).not.toContain("ScriptApp.newTrigger");

    const guide = await zip.file("BAT_DAU_TU_DAY.html").async("string");
    expect(guide).toContain("Google Sheets");
    expect(guide).toContain("Apps Script");
    expect(guide).toContain("Save as Google Sheets");
    expect(guide).toContain("script.google.com");
    expect(guide).toContain("Review permissions");
    expect(guide).toContain('class="code-onopen">onOpen</code> thành <code class="code-setup">setup');
    expect(guide).toContain("Execution completed");
    expect(guide).toContain("popup bị chặn");
    expect(guide).toContain("This app is blocked");
    expect(guide).toContain("email của chính bạn");
    expect(guide).toContain("tối đa 90");
    expect(guide).toContain("1 dòng READY");
    expect(guide).toContain("Đợt thu tiếp theo: dùng lại Sheet này");
    expect(guide).toContain("0. Nhập dữ liệu mới");
  });

  it("exports one versioned data file for the Google Sheet template", () => {
    const payload = createGmailSheetData([validItem], DEFAULT_EMAIL_TEMPLATE, {
      exportedAt: "2026-07-27T01:02:03.000Z",
    });

    expect(payload).toMatchObject({
      kind: "dekisugi.gmail-data",
      schemaVersion: 1,
      exportedAt: "2026-07-27T01:02:03.000Z",
      summary: { total: 1, ready: 1, errors: 0 },
    });
    expect(payload.rows).toHaveLength(1);
    expect(payload.rows[0]).toMatchObject({
      Email: "parent@example.com",
      StudentCode: "HS001",
      Status: "READY",
      Send: "YES",
    });
    expect(payload.rows[0].HtmlBody).toContain("cid:dekisugi_qr");
    expect(payload.rows[0].QrBase64).not.toContain("data:image/png;base64,");
  });

  it("ships a user-owned Sheet sidebar with guarded sending and opt-in scheduling", () => {
    const code = readFileSync(new URL("../google-sheet-template/Code.gs", import.meta.url), "utf8");
    const sidebar = readFileSync(new URL("../google-sheet-template/Sidebar.html", import.meta.url), "utf8");
    const manifest = JSON.parse(readFileSync(new URL("../google-sheet-template/appsscript.json", import.meta.url), "utf8"));

    expect(code).toContain('kind !== "dekisugi.gmail-data"');
    expect(code).toContain("safeSheetText_");
    expect(code).toContain("LockService.getDocumentLock");
    expect(code).toContain("TEST_SENT_AT");
    expect(code).toContain("ScriptApp.newTrigger");
    expect(code).toContain("runScheduledBatch");
    expect(code).toContain('getSheetByName("BAT_DAU")');
    expect(code).toContain("DEKISUGI EMAIL — BẮT ĐẦU TỪ ĐÂY");
    expect(code).toContain("MailApp.getRemainingDailyQuota");
    expect(code).toContain("RESERVED_QUOTA: 10");
    expect(code).toContain("inlineImages");
    expect(code).not.toContain("UrlFetchApp");
    expect(code).not.toContain("getUi().alert");
    expect(code).not.toContain("getUi().prompt");

    expect(sidebar).toContain("DEKISUGI_GMAIL_DATA.json");
    expect(sidebar).toContain("Gửi 1 email thử");
    expect(sidebar).toContain("Tạm dừng lịch gửi");
    expect(sidebar).toContain("confirmModal");
    expect(sidebar).not.toContain("window.confirm");
    expect(sidebar).not.toContain("window.alert");

    expect(manifest.oauthScopes).toEqual(expect.arrayContaining([
      "https://www.googleapis.com/auth/spreadsheets.currentonly",
      "https://www.googleapis.com/auth/script.send_mail",
      "https://www.googleapis.com/auth/script.scriptapp",
    ]));
  });
});
