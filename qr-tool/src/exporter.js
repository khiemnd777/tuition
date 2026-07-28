import JSZip from "jszip";
import * as XLSX from "xlsx";

import { buildEml, renderEmailTemplate, templateForGmailMerge } from "./email.js";
import { ascii, formatVND } from "./vietqr.js";

function csvCell(value) {
  const text = String(value ?? "");
  return /[",\r\n]/.test(text) ? `"${text.replace(/"/g, '""')}"` : text;
}

export function toCSV(rows, columns) {
  const header = columns.map(([key, label]) => csvCell(label || key)).join(",");
  const body = rows.map((row) => columns.map(([key]) => csvCell(row[key])).join(",")).join("\r\n");
  return `\uFEFF${header}${body ? `\r\n${body}` : ""}\r\n`;
}

export function safeFilename(value) {
  return ascii(String(value || "file"))
    .replace(/[^A-Za-z0-9_-]+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^[-_.]+|[-_.]+$/g, "")
    .slice(0, 120) || "file";
}

export function isGmailSheetTemplateURL(value) {
  try {
    const url = new URL(String(value || ""));
    return url.protocol === "https:"
      && url.hostname === "docs.google.com"
      && /^\/spreadsheets\/d\/[^/]+\/copy\/?$/.test(url.pathname);
  } catch {
    return false;
  }
}

function allocateName(base, used) {
  const root = safeFilename(base);
  let name = root;
  let suffix = 2;
  while (used.has(name.toLowerCase())) {
    name = `${root}-${suffix}`;
    suffix += 1;
  }
  used.add(name.toLowerCase());
  return name;
}

function itemBaseName(item) {
  return [item.billNumber, item.studentName || item.id].filter(Boolean).join("-");
}

function qrBase64(item) {
  return String(item.qrData || "").replace(/^data:image\/png;base64,/, "").replace(/\s/g, "");
}

function manifestRow(item, filename = "") {
  return {
    sourceRow: item.sourceRow || "",
    id: item.id || "",
    studentCode: item.studentCode || "",
    studentName: item.studentName || "",
    schoolName: item.schoolName || "",
    cohort: item.cohort || "",
    year: item.year || "",
    parentName: item.parentName || "",
    className: item.className || "",
    email: item.email || "",
    bankBin: item.bankBin || "",
    bankName: item.bankName || "",
    bankAccount: item.bankAccount || "",
    amount: item.amount || 0,
    billNumber: item.billNumber || "",
    note: item.note || "",
    qrFilename: filename,
    status: item.errors?.length ? "error" : "ready",
    errors: (item.errors || []).join("; "),
    vietqr: item.vietqr || "",
  };
}

const MANIFEST_COLUMNS = [
  ["sourceRow", "Dòng nguồn"],
  ["id", "ID"],
  ["studentCode", "Mã học sinh"],
  ["studentName", "Tên học sinh"],
  ["schoolName", "Tên trường"],
  ["cohort", "Niên khóa"],
  ["year", "Năm"],
  ["parentName", "Tên phụ huynh"],
  ["className", "Lớp"],
  ["email", "Email"],
  ["bankBin", "BIN"],
  ["bankName", "Ngân hàng"],
  ["bankAccount", "Số tài khoản"],
  ["amount", "Số tiền"],
  ["billNumber", "Mã tham chiếu"],
  ["note", "Nội dung"],
  ["qrFilename", "File QR"],
  ["status", "Trạng thái"],
  ["errors", "Lỗi"],
  ["vietqr", "VietQR payload"],
];

export async function createQRBundle(items) {
  const zip = new JSZip();
  const used = new Set();
  const manifest = [];
  const errorRows = [];

  for (const item of items) {
    let filename = "";
    if (!item.errors?.length && item.qrData) {
      filename = `${allocateName(itemBaseName(item), used)}.png`;
      zip.file(`qr/${filename}`, qrBase64(item), { base64: true });
    } else {
      errorRows.push(manifestRow(item));
    }
    manifest.push(manifestRow(item, filename));
  }

  zip.file("manifest.csv", toCSV(manifest, MANIFEST_COLUMNS));
  if (errorRows.length > 0) zip.file("errors.csv", toCSV(errorRows, MANIFEST_COLUMNS));
  return zip.generateAsync({ type: "blob", mimeType: "application/zip", compression: "DEFLATE" });
}

function splitRecipientName(item) {
  const words = String(item.parentName || item.studentName || "").trim().split(/\s+/).filter(Boolean);
  if (words.length <= 1) return { firstName: words[0] || "", lastName: "" };
  return { firstName: words.at(-1), lastName: words.slice(0, -1).join(" ") };
}

function gmailRows(items, fileById) {
  return items.map((item) => {
    const { firstName, lastName } = splitRecipientName(item);
    const paymentItems = (item.paymentItems || []).map((entry) => `${entry.label}: ${formatVND(entry.amount)}`).join(" | ");
    return {
      Email: String(item.email || ""),
      FirstName: String(firstName),
      LastName: String(lastName),
      ParentName: String(item.parentName || ""),
      StudentCode: String(item.studentCode || ""),
      StudentName: String(item.studentName || ""),
      SchoolName: String(item.schoolName || ""),
      Cohort: String(item.cohort || ""),
      Year: String(item.year || ""),
      ClassName: String(item.className || ""),
      Amount: formatVND(item.amount),
      PaymentItems: paymentItems,
      BankName: String(item.bankName || ""),
      BankAccount: String(item.bankAccount || ""),
      BillNumber: String(item.billNumber || ""),
      PaymentNote: String(item.note || ""),
      QrFilename: fileById.get(item.id)?.qr || "",
    };
  });
}

function workbookArray(rows) {
  const workbook = XLSX.utils.book_new();
  const sheet = XLSX.utils.json_to_sheet(rows, { skipHeader: false });
  XLSX.utils.book_append_sheet(workbook, sheet, "Recipients");
  return XLSX.write(workbook, { type: "array", bookType: "xlsx" });
}

function emailIsValid(value) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(value || ""));
}

function gmailFreeRows(items, template) {
  return items.map((item) => {
    const errors = [...(item.errors || [])];
    if (!emailIsValid(item.email)) errors.push(item.email ? "Email không hợp lệ" : "Thiếu email người nhận");
    if (!item.qrData) errors.push("QR chưa sẵn sàng");
    const rendered = renderEmailTemplate(template, item, { qrSrc: "cid:dekisugi_qr" });
    return {
      SourceRow: String(item.sourceRow || ""),
      Send: errors.length ? "NO" : "YES",
      Email: String(item.email || ""),
      StudentCode: String(item.studentCode || ""),
      StudentName: String(item.studentName || ""),
      ParentName: String(item.parentName || ""),
      SchoolName: String(item.schoolName || ""),
      Cohort: String(item.cohort || ""),
      Year: String(item.year || ""),
      ClassName: String(item.className || ""),
      Amount: formatVND(item.amount),
      BillNumber: String(item.billNumber || ""),
      Subject: rendered.subject,
      HtmlBody: rendered.html,
      TextBody: rendered.text,
      QrBase64: qrBase64(item),
      QrFilename: `${safeFilename(itemBaseName(item))}.png`,
      Status: errors.length ? "ERROR" : "READY",
      SentAt: "",
      Error: [...new Set(errors)].join("; "),
    };
  });
}

export function createGmailSheetData(items, template, { exportedAt = new Date().toISOString() } = {}) {
  if (!Array.isArray(items)) throw new Error("Danh sách email không hợp lệ");
  if (items.length > 500) throw new Error("Google Sheet mẫu chỉ nhận tối đa 500 dòng mỗi lần import");
  const rows = gmailFreeRows(items, template);
  const ready = rows.filter((row) => row.Status === "READY").length;
  return {
    kind: "dekisugi.gmail-data",
    schemaVersion: 1,
    exportedAt,
    summary: {
      total: rows.length,
      ready,
      errors: rows.length - ready,
    },
    rows,
  };
}

function gmailFreeWorkbookArray(rows) {
  const workbook = XLSX.utils.book_new();
  workbook.Props = {
    Title: "DEKISUGI Gmail Free Sender",
    Subject: "Danh sách email thanh toán với QR riêng",
    Author: "DEKISUGI QR Tool",
  };

  const guideSheet = XLSX.utils.aoa_to_sheet([
    ["DEKISUGI EMAIL — BẮT ĐẦU TỪ ĐÂY", "", "", ""],
    ["Chặng", "Bước", "Bấm ở đâu / bấm gì", "Kết quả cần thấy"],
    ["1. Tạo Google Sheet", "1.1", "Upload workbook này lên Google Drive rồi nhấp đúp để mở bằng Google Sheets.", "File mở được trong trình chỉnh sửa Google Sheets."],
    ["1. Tạo Google Sheet", "1.2", "Trên thanh menu chọn File → Save as Google Sheets, sau đó chuyển sang bản mới.", "Tên bản mới không còn đuôi .xlsx."],
    ["2. Dán code", "2.1", "Trong bản Google Sheets mới, chọn Extensions → Apps Script.", "Tab Apps Script mở và có file Code.gs."],
    ["2. Dán code", "2.2", "Trong Code.gs nhấn ⌘A trên Mac hoặc Ctrl+A trên Windows, xoá code mẫu, rồi paste toàn bộ 02_CODE_GUI_EMAIL.gs.", "Code DEKISUGI xuất hiện trong Code.gs."],
    ["2. Dán code", "2.3", "Bấm biểu tượng Save project hình đĩa mềm ở thanh trên.", "Saving project… đổi thành Saved to Drive."],
    ["3. Chạy setup", "3.1", "Trên thanh phía trên code, bấm danh sách đang hiện onOpen và chọn setup.", "Danh sách hàm đang hiển thị setup."],
    ["3. Chạy setup", "3.2", "Bấm nút ▶ Run ở bên trái nút Debug.", "Hộp Authorization required xuất hiện ở lần chạy đầu."],
    ["3. Cấp quyền", "3.3", "Trong Authorization required, bấm Review permissions rồi chọn đúng tài khoản Gmail dùng để gửi.", "Màn hình xem quyền truy cập xuất hiện."],
    ["3. Cấp quyền", "3.4", "Nếu thấy Google hasn't verified this app: bấm Advanced → Go to <tên project> (unsafe). Chỉ tiếp tục khi URL là script.google.com và code không bị sửa.", "Màn hình Continue hoặc Allow xuất hiện."],
    ["3. Cấp quyền", "3.5", "Kiểm tra quyền liên quan đến Sheet hiện tại và gửi email, rồi bấm Continue hoặc Allow.", "Quay lại Apps Script; Execution log hiện Execution completed."],
    ["4. Kiểm tra", "4.1", "Quay lại Google Sheets và reload trang. Kiểm tra menu DEKISUGI Email có mục 0. Nhập dữ liệu mới.", "Sheet đã sẵn sàng dùng lại cho các đợt thu sau."],
    ["4. Kết nối", "4.2", "Copy link Sheet trên thanh địa chỉ. Quay lại QR Tool, bấm Tôi đã có Sheet Gmail và dán link.", "QR Tool hiện Đã kết nối Sheet Gmail."],
    ["4. Kiểm tra", "4.3", "Chọn DEKISUGI Email → 1. Kiểm tra danh sách.", "Popup báo số dòng READY, ERROR và SKIP."],
    ["5. Gửi thử", "5.1", "Trong sheet EMAILS, chọn một ô ở dòng READY rồi chọn DEKISUGI Email → 2. Gửi thử dòng đang chọn.", "Hộp nhập email nhận bản thử xuất hiện."],
    ["5. Gửi thử", "5.2", "Nhập email của chính bạn, bấm OK, rồi kiểm tra nội dung, QR, số tiền và học sinh.", "Email thử đến hộp thư của bạn; dòng vẫn chưa chuyển SENT."],
    ["5. Gửi thật", "5.3", "Chỉ khi email thử đúng, chọn DEKISUGI Email → 3. Gửi email chưa gửi; đọc số lượng rồi bấm YES.", "Dòng gửi thành công chuyển SENT và không tự gửi lại."],
    ["6. Đợt thu mới", "6.1", "Từ QR Tool, tải file DEKISUGI_GMAIL_DATA.json rồi mở lại Sheet này.", "Không dán code, không chạy setup và không cấp quyền lại."],
    ["6. Đợt thu mới", "6.2", "Chọn DEKISUGI Email → 0. Nhập dữ liệu mới, chọn file JSON và xác nhận thay danh sách.", "Sheet EMAILS chứa dữ liệu mới; code và quyền Gmail được giữ nguyên."],
    [],
    ["Trạng thái", "", "Ý nghĩa", "Bạn cần làm gì"],
    ["READY", "", "Đủ dữ liệu và đang chờ gửi.", "Có thể gửi thử hoặc gửi thật."],
    ["SENT", "", "Đã gửi.", "Không sửa về READY nếu không thực sự muốn gửi lại."],
    ["ERROR", "", "Không gửi được; nguyên nhân nằm ở cột Error.", "Sửa dữ liệu rồi chạy 1. Kiểm tra danh sách."],
    ["SENDING", "", "Đang xử lý hoặc lần chạy trước bị gián đoạn.", "Chạy 1. Kiểm tra danh sách để khôi phục trạng thái."],
    ["SKIP", "", "Dòng có Send = NO.", "Đổi Send thành YES nếu muốn gửi."],
    [],
    ["Bị kẹt?", "", "Mở BAT_DAU_TU_DAY.html trong ZIP để xem cách xử lý khi không thấy menu, popup bị chặn hoặc tài khoản bị Google/Workspace chặn.", "Không gửi hoặc attach ba file trong ZIP cho phụ huynh."],
  ]);
  guideSheet["!cols"] = [{ wch: 22 }, { wch: 10 }, { wch: 95 }, { wch: 58 }];
  guideSheet["!merges"] = [XLSX.utils.decode_range("A1:D1")];
  XLSX.utils.book_append_sheet(workbook, guideSheet, "BAT_DAU");

  const emailSheet = XLSX.utils.json_to_sheet(rows, { skipHeader: false });
  emailSheet["!autofilter"] = { ref: `A1:T${Math.max(2, rows.length + 1)}` };
  emailSheet["!cols"] = [
    { wch: 10 }, { wch: 8 }, { wch: 30 }, { wch: 16 }, { wch: 24 },
    { wch: 24 }, { wch: 24 }, { wch: 15 }, { wch: 12 }, { wch: 12 },
    { wch: 16 }, { wch: 18 }, { wch: 42 }, { hidden: true }, { hidden: true },
    { hidden: true }, { hidden: true }, { wch: 12 }, { wch: 22 }, { wch: 42 },
  ];
  XLSX.utils.book_append_sheet(workbook, emailSheet, "EMAILS");
  return XLSX.write(workbook, { type: "array", bookType: "xlsx", compression: true });
}

export const GMAIL_FREE_APPS_SCRIPT = String.raw`/**
 * DEKISUGI Gmail Free Sender
 * Script chỉ đọc workbook hiện tại và gửi email bằng tài khoản Google đang cấp quyền.
 * @OnlyCurrentDoc
 */

var DEKISUGI_CONFIG = Object.freeze({
  SHEET_NAME: "EMAILS",
  MAX_ROWS: 500,
  MAX_CELL_LENGTH: 49000,
  MAX_PER_BATCH: 90,
  RESERVED_QUOTA: 10,
  SEND_DELAY_MS: 500,
  QR_CONTENT_ID: "dekisugi_qr"
});

var DEKISUGI_HEADERS = [
  "SourceRow", "Send", "Email", "StudentCode", "StudentName", "ParentName",
  "SchoolName", "Cohort", "Year", "ClassName", "Amount", "BillNumber",
  "Subject", "HtmlBody", "TextBody", "QrBase64", "QrFilename", "Status",
  "SentAt", "Error"
];

function onOpen() {
  SpreadsheetApp.getUi()
    .createMenu("DEKISUGI Email")
    .addItem("0. Nhập dữ liệu mới", "showImportDialog")
    .addSeparator()
    .addItem("1. Kiểm tra danh sách", "validateList")
    .addItem("2. Gửi thử dòng đang chọn", "sendSelectedTest")
    .addSeparator()
    .addItem("3. Gửi email chưa gửi", "sendPendingEmails")
    .addItem("4. Xem các dòng bị lỗi", "showErrorRows")
    .addToUi();
}

function setup() {
  var result = validateListInternal_();
  formatSheet_();
  var quota = MailApp.getRemainingDailyQuota();
  onOpen();
  var message = "SETUP HOÀN TẤT | READY: " + result.ready +
    " | ERROR: " + result.errors +
    " | Quota Google còn lại: " + quota +
    " | Quay lại Google Sheets và reload trang.";
  console.log(message);
  SpreadsheetApp.getActive().toast(message, "DEKISUGI Email", 12);
  return message;
}

function showImportDialog() {
  var html = HtmlService.createHtmlOutput(
    '<!doctype html><html><head><base target="_top"><meta charset="utf-8"><style>' +
    'body{margin:0;padding:22px;font:14px/1.55 Arial,sans-serif;color:#172033;background:#f7f9fc}h2{margin:0 0 8px;font-size:20px}p{margin:7px 0;color:#475467}.warning{margin:14px 0;padding:11px 13px;border-left:4px solid #c9851d;background:#fff8e8;color:#6f4c12}.picker{display:block;margin-top:15px;padding:20px;border:2px dashed #9fb3df;border-radius:8px;background:#fff;text-align:center;color:#2b50a1;font-weight:700;cursor:pointer}.picker input{display:block;width:100%;margin-top:10px}button{width:100%;min-height:42px;margin-top:14px;border:0;border-radius:6px;background:#2b50a1;color:#fff;font:700 14px Arial;cursor:pointer}button:disabled{opacity:.45;cursor:not-allowed}.status{min-height:22px;margin-top:12px;font-weight:700}.error{color:#b42318}.success{color:#16794c}</style></head><body>' +
    '<h2>Nhập dữ liệu đợt thu mới</h2>' +
    '<p>Chọn file <strong>DEKISUGI_GMAIL_DATA.json</strong> vừa tải từ QR Tool.</p>' +
    '<p class="warning"><strong>Lưu ý:</strong> danh sách mới sẽ thay toàn bộ danh sách hiện tại trong sheet EMAILS. Code và quyền Gmail vẫn được giữ nguyên.</p>' +
    '<label class="picker">Chọn file JSON<input id="file" type="file" accept="application/json,.json"></label>' +
    '<button id="run" type="button" disabled>Thay bằng dữ liệu mới</button><p id="status" class="status">Chưa chọn file.</p>' +
    '<script>' +
    'var fileInput=document.getElementById("file");var runButton=document.getElementById("run");var statusNode=document.getElementById("status");' +
    'fileInput.onchange=function(){runButton.disabled=!fileInput.files[0];statusNode.className="status";statusNode.textContent=fileInput.files[0]?"Đã chọn: "+fileInput.files[0].name:"Chưa chọn file.";};' +
    'runButton.onclick=function(){var file=fileInput.files[0];if(!file)return;runButton.disabled=true;statusNode.className="status";statusNode.textContent="Đang nhập dữ liệu…";var reader=new FileReader();reader.onload=function(){try{var payload=JSON.parse(reader.result);google.script.run.withSuccessHandler(function(result){statusNode.className="status success";statusNode.textContent="Đã nhập "+result.total+" dòng · READY: "+result.ready+" · ERROR: "+result.errors;runButton.disabled=false;}).withFailureHandler(function(error){statusNode.className="status error";statusNode.textContent=error&&error.message?error.message:String(error);runButton.disabled=false;}).importDataset(payload);}catch(error){statusNode.className="status error";statusNode.textContent="Không đọc được JSON: "+error.message;runButton.disabled=false;}};reader.onerror=function(){statusNode.className="status error";statusNode.textContent="Không đọc được file đã chọn.";runButton.disabled=false;};reader.readAsText(file);};' +
    '<\/script></body></html>'
  ).setWidth(540).setHeight(430);
  SpreadsheetApp.getUi().showModalDialog(html, "DEKISUGI Email · Nhập dữ liệu mới");
}

function importDataset(payload) {
  var lock = LockService.getUserLock();
  if (!lock.tryLock(5000)) throw new Error("Một thao tác khác đang chạy. Vui lòng thử lại sau vài giây.");
  try {
    validateDataset_(payload);
    var sheet = SpreadsheetApp.getActive().getSheetByName(DEKISUGI_CONFIG.SHEET_NAME);
    if (!sheet) sheet = SpreadsheetApp.getActive().insertSheet(DEKISUGI_CONFIG.SHEET_NAME);
    if (sheet.getMaxColumns() < DEKISUGI_HEADERS.length) {
      sheet.insertColumnsAfter(sheet.getMaxColumns(), DEKISUGI_HEADERS.length - sheet.getMaxColumns());
    }
    if (sheet.getFilter()) sheet.getFilter().remove();
    sheet.showColumns(1, sheet.getMaxColumns());
    sheet.clear();
    sheet.getRange(1, 1, 1, DEKISUGI_HEADERS.length).setValues([DEKISUGI_HEADERS]);
    var values = payload.rows.map(function (row) {
      return DEKISUGI_HEADERS.map(function (header) { return safeSheetText_(row[header]); });
    });
    var range = sheet.getRange(2, 1, values.length, DEKISUGI_HEADERS.length);
    range.setNumberFormat("@");
    range.setValues(values);
    SpreadsheetApp.flush();
    var result = validateListInternal_();
    formatSheet_();
    return { total: values.length, ready: result.ready, errors: result.errors, skipped: result.skipped };
  } finally {
    lock.releaseLock();
  }
}

function validateDataset_(payload) {
  if (!payload || payload.kind !== "dekisugi.gmail-data" || Number(payload.schemaVersion) !== 1) {
    throw new Error("Đây không phải file dữ liệu Gmail do DEKISUGI xuất.");
  }
  if (!Array.isArray(payload.rows) || !payload.rows.length) throw new Error("File dữ liệu không có dòng nào để nhập.");
  if (payload.rows.length > DEKISUGI_CONFIG.MAX_ROWS) throw new Error("Chỉ nhận tối đa 500 dòng mỗi lần import.");
  payload.rows.forEach(function (row, rowIndex) {
    if (!row || typeof row !== "object" || Array.isArray(row)) throw new Error("Dòng " + (rowIndex + 1) + " không hợp lệ.");
    DEKISUGI_HEADERS.forEach(function (header) {
      var value = row[header] === undefined || row[header] === null ? "" : String(row[header]);
      if (value.length > DEKISUGI_CONFIG.MAX_CELL_LENGTH) {
        throw new Error("Dòng " + (rowIndex + 1) + ", cột " + header + " vượt giới hạn Google Sheets.");
      }
    });
  });
}

function safeSheetText_(value) {
  var text = value === undefined || value === null ? "" : String(value);
  return /^[=+\-@]/.test(text) ? "'" + text : text;
}

function validateList() {
  try {
    var result = validateListInternal_();
    SpreadsheetApp.getUi().alert(
      "Kiểm tra hoàn tất",
      "READY: " + result.ready + " | ERROR: " + result.errors + " | SKIP: " + result.skipped,
      SpreadsheetApp.getUi().ButtonSet.OK
    );
  } catch (error) {
    showFailure_(error);
  }
}

function sendSelectedTest() {
  try {
    var context = getContext_();
    var activeSheet = SpreadsheetApp.getActiveSheet();
    var activeRange = activeSheet.getActiveRange();
    if (activeSheet.getName() !== DEKISUGI_CONFIG.SHEET_NAME || !activeRange || activeRange.getRow() < 2) {
      throw new Error("Hãy chọn một ô trên dòng cần gửi thử trong sheet EMAILS.");
    }
    var rowIndex = activeRange.getRow() - 1;
    var record = recordForRow_(context, rowIndex);
    var errors = validateRecord_(record);
    if (errors.length) throw new Error(errors.join("; "));
    if (MailApp.getRemainingDailyQuota() < 1) throw new Error("Tài khoản đã hết quota gửi email hôm nay.");

    var ui = SpreadsheetApp.getUi();
    var prompt = ui.prompt(
      "Gửi thử",
      "Nhập email sẽ nhận bản thử. Email của phụ huynh không bị thay đổi:",
      ui.ButtonSet.OK_CANCEL
    );
    if (prompt.getSelectedButton() !== ui.Button.OK) return;
    var testEmail = String(prompt.getResponseText() || "").trim();
    if (!isEmail_(testEmail)) throw new Error("Email nhận bản thử không hợp lệ.");
    sendRecord_(record, testEmail);
    ui.alert("Đã gửi một email thử tới " + testEmail + ". Dòng này vẫn giữ trạng thái " + record.Status + ".");
  } catch (error) {
    showFailure_(error);
  }
}

function sendPendingEmails() {
  var lock = LockService.getUserLock();
  if (!lock.tryLock(2000)) {
    SpreadsheetApp.getUi().alert("Một lượt gửi khác đang chạy. Vui lòng chờ và thử lại.");
    return;
  }

  try {
    validateListInternal_();
    var context = getContext_();
    var remainingQuota = MailApp.getRemainingDailyQuota();
    var safeQuota = Math.max(0, Math.min(
      DEKISUGI_CONFIG.MAX_PER_BATCH,
      remainingQuota - DEKISUGI_CONFIG.RESERVED_QUOTA
    ));
    if (safeQuota < 1) throw new Error("Quota còn lại không đủ. Script luôn chừa 10 email cho nhu cầu gửi thông thường.");

    var pendingIndexes = [];
    for (var rowIndex = 1; rowIndex < context.values.length; rowIndex += 1) {
      var record = recordForRow_(context, rowIndex);
      if (record.Status === "READY" && String(record.Send).toUpperCase() !== "NO") pendingIndexes.push(rowIndex);
    }
    if (!pendingIndexes.length) throw new Error("Không có dòng READY nào đang chờ gửi.");

    var sendCount = Math.min(safeQuota, pendingIndexes.length);
    var ui = SpreadsheetApp.getUi();
    var decision = ui.alert(
      "Xác nhận gửi email thật",
      "Sẽ gửi " + sendCount + " email từ tài khoản Google đang đăng nhập.\n" +
        "Mỗi dòng SENT sẽ không được tự động gửi lại. Tiếp tục?",
      ui.ButtonSet.YES_NO
    );
    if (decision !== ui.Button.YES) return;

    var sent = 0;
    var failed = 0;
    for (var index = 0; index < sendCount; index += 1) {
      var currentRowIndex = pendingIndexes[index];
      var currentRecord = recordForRow_(context, currentRowIndex);
      writeCell_(context, currentRowIndex, "Status", "SENDING");
      writeCell_(context, currentRowIndex, "Error", "");
      SpreadsheetApp.flush();
      try {
        sendRecord_(currentRecord, currentRecord.Email);
        writeCell_(context, currentRowIndex, "Status", "SENT");
        writeCell_(context, currentRowIndex, "SentAt", new Date());
        sent += 1;
      } catch (error) {
        var message = cleanError_(error);
        writeCell_(context, currentRowIndex, "Status", "ERROR");
        writeCell_(context, currentRowIndex, "Error", message);
        failed += 1;
        if (/quota|limit|too many|service invoked/i.test(message)) break;
      }
      SpreadsheetApp.flush();
      Utilities.sleep(DEKISUGI_CONFIG.SEND_DELAY_MS);
    }

    ui.alert(
      "Lượt gửi hoàn tất",
      "Đã gửi: " + sent + " | Lỗi: " + failed +
        "\nCác dòng READY còn lại có thể gửi ở lượt hoặc ngày tiếp theo.",
      ui.ButtonSet.OK
    );
  } catch (error) {
    showFailure_(error);
  } finally {
    lock.releaseLock();
  }
}

function showErrorRows() {
  try {
    var context = getContext_();
    var lines = [];
    for (var rowIndex = 1; rowIndex < context.values.length; rowIndex += 1) {
      var record = recordForRow_(context, rowIndex);
      if (record.Status === "ERROR") {
        lines.push("Dòng " + (rowIndex + 1) + " · " + (record.Email || "thiếu email") + " · " + (record.Error || "Lỗi chưa xác định"));
      }
      if (lines.length >= 12) break;
    }
    SpreadsheetApp.getUi().alert(
      "Các dòng lỗi",
      lines.length ? lines.join("\n") : "Không có dòng ERROR.",
      SpreadsheetApp.getUi().ButtonSet.OK
    );
  } catch (error) {
    showFailure_(error);
  }
}

function validateListInternal_() {
  var context = getContext_();
  var statuses = [];
  var errorsOutput = [];
  var result = { ready: 0, errors: 0, skipped: 0 };
  for (var rowIndex = 1; rowIndex < context.values.length; rowIndex += 1) {
    var record = recordForRow_(context, rowIndex);
    if (record.Status === "SENT") {
      statuses.push(["SENT"]);
      errorsOutput.push([record.Error || ""]);
      continue;
    }
    var errors = validateRecord_(record);
    var skipped = String(record.Send || "").toUpperCase() === "NO" && !errors.length;
    var status = errors.length ? "ERROR" : (skipped ? "SKIP" : "READY");
    statuses.push([status]);
    errorsOutput.push([errors.join("; ")]);
    if (status === "READY") result.ready += 1;
    else if (status === "ERROR") result.errors += 1;
    else result.skipped += 1;
  }
  if (statuses.length) {
    context.sheet.getRange(2, context.index.Status + 1, statuses.length, 1).setValues(statuses);
    context.sheet.getRange(2, context.index.Error + 1, errorsOutput.length, 1).setValues(errorsOutput);
  }
  SpreadsheetApp.flush();
  return result;
}

function getContext_() {
  var sheet = SpreadsheetApp.getActive().getSheetByName(DEKISUGI_CONFIG.SHEET_NAME);
  if (!sheet) throw new Error("Không tìm thấy sheet EMAILS.");
  var values = sheet.getDataRange().getDisplayValues();
  if (values.length < 2) throw new Error("Sheet EMAILS chưa có dữ liệu.");
  var index = {};
  values[0].forEach(function (header, columnIndex) { index[String(header).trim()] = columnIndex; });
  ["Send", "Email", "Subject", "HtmlBody", "TextBody", "QrBase64", "QrFilename", "Status", "SentAt", "Error"].forEach(function (header) {
    if (index[header] === undefined) throw new Error("Thiếu cột bắt buộc: " + header);
  });
  return { sheet: sheet, values: values, index: index };
}

function recordForRow_(context, rowIndex) {
  var record = {};
  Object.keys(context.index).forEach(function (header) {
    record[header] = context.values[rowIndex][context.index[header]];
  });
  return record;
}

function validateRecord_(record) {
  var errors = [];
  if (!isEmail_(record.Email)) errors.push("Email không hợp lệ");
  if (!String(record.Subject || "").trim()) errors.push("Thiếu Subject");
  if (!String(record.HtmlBody || "").trim()) errors.push("Thiếu HtmlBody");
  if (!String(record.TextBody || "").trim()) errors.push("Thiếu TextBody");
  var base64 = String(record.QrBase64 || "").replace(/\s/g, "");
  if (!base64 || !/^[A-Za-z0-9+/=]+$/.test(base64)) errors.push("QR Base64 không hợp lệ");
  return errors;
}

function sendRecord_(record, recipient) {
  var cleanBase64 = String(record.QrBase64 || "").replace(/^data:image\/png;base64,/, "").replace(/\s/g, "");
  var qrBlob = Utilities.newBlob(
    Utilities.base64Decode(cleanBase64),
    "image/png",
    String(record.QrFilename || "vietqr.png")
  );
  var inlineImages = {};
  inlineImages[DEKISUGI_CONFIG.QR_CONTENT_ID] = qrBlob;
  MailApp.sendEmail({
    to: recipient,
    subject: String(record.Subject || ""),
    body: String(record.TextBody || ""),
    htmlBody: String(record.HtmlBody || ""),
    inlineImages: inlineImages
  });
}

function formatSheet_() {
  var context = getContext_();
  context.sheet.setFrozenRows(1);
  context.sheet.getRange(1, 1, 1, context.values[0].length)
    .setBackground("#e8eefb")
    .setFontColor("#172033")
    .setFontWeight("bold");
  if (!context.sheet.getFilter()) context.sheet.getDataRange().createFilter();
  ["HtmlBody", "TextBody", "QrBase64", "QrFilename"].forEach(function (header) {
    context.sheet.hideColumns(context.index[header] + 1);
  });
}

function writeCell_(context, rowIndex, header, value) {
  context.sheet.getRange(rowIndex + 1, context.index[header] + 1).setValue(value);
  context.values[rowIndex][context.index[header]] = String(value || "");
}

function isEmail_(value) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(value || "").trim());
}

function cleanError_(error) {
  return String(error && error.message ? error.message : error || "Lỗi chưa xác định").replace(/[\r\n]+/g, " ").slice(0, 500);
}

function showFailure_(error) {
  SpreadsheetApp.getUi().alert("Không thể thực hiện", cleanError_(error), SpreadsheetApp.getUi().ButtonSet.OK);
}`;

function gmailFreeGuideHTML({ total, ready, errors }) {
  return `<!doctype html>
<html lang="vi"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>DEKISUGI Gmail miễn phí — Bắt đầu từ đây</title>
<style>
:root{font-family:Arial,sans-serif;color:#172033;background:#f4f7fb;font-synthesis:none}*{box-sizing:border-box}body{margin:0}.wrap{max-width:940px;margin:36px auto;padding:0 18px 64px}.card{background:#fff;border:1px solid #d9dee8;border-radius:10px;padding:32px;box-shadow:0 12px 30px rgba(15,23,42,.08)}h1{margin:0;font-size:30px;letter-spacing:-.02em}h2{margin:0;font-size:18px;color:#24324a}p,li{line-height:1.6}.lead{margin:9px 0;color:#536077}.summary{display:flex;gap:8px;flex-wrap:wrap;margin:20px 0}.pill{padding:7px 11px;border-radius:999px;background:#eef3ff;color:#2b50a1;font-weight:700;font-size:13px}.pill.error{background:#fff1f0;color:#b42318}.notice{margin:18px 0;padding:13px 15px;border-left:4px solid #16794c;background:#edf9f3;color:#315b47}.notice.warning{border-left-color:#98620b;background:#fff8e8;color:#6f4c12}.phase{margin-top:18px;padding:20px;border:1px solid #d9dee8;border-radius:9px}.phase.emphasis{border-color:#9fb3df;background:#f8faff}.phase-head{display:flex;align-items:center;gap:11px}.phase-number{flex:0 0 32px;width:32px;height:32px;display:grid;place-items:center;border-radius:50%;background:#2b50a1;color:#fff;font-weight:800}.phase-head small{display:block;margin-top:3px;color:#667085}.steps{margin:14px 0 0;padding-left:23px}.steps li{margin:9px 0;padding-left:4px}.where{margin:15px 0 7px;font-size:13px;font-weight:700}.toolbar{display:flex;flex-wrap:wrap;gap:7px;padding:11px;border:1px solid #c7cfdd;border-radius:7px;background:#eef1f5;font:12px/1.3 Consolas,monospace}.toolbar span{padding:7px 10px;border:1px solid #c7cfdd;border-radius:5px;background:#fff}.toolbar .run{border-color:#8ea8de;color:#2b50a1;font-weight:800}.toolbar .function-selector{border:2px solid #cf7a13;background:#fff1b8;color:#633b00;font-size:14px;font-weight:900;box-shadow:0 0 0 2px rgba(207,122,19,.12)}.expected{margin:14px 0 0;padding:11px 13px;border-left:3px solid #6ea787;background:#f2f8f4;color:#385848}.auth-flow{margin-top:14px;padding:14px;border:1px solid #e8c98b;border-radius:7px;background:#fffaf0}.auth-flow h3{margin:0 0 8px;font-size:15px;color:#724e0e}.button-name{display:inline-block;padding:2px 7px;border:1px solid #b5c0d1;border-radius:4px;background:#fff;color:#24324a;font-weight:700;white-space:nowrap}code{padding:2px 5px;border-radius:4px;background:#f1f4f8;font-family:Consolas,monospace;font-size:.92em}.code-onopen,.code-setup{display:inline-block;padding:3px 7px;font-weight:900}.code-onopen{border:1px solid #cf7a13;background:#fff1b8;color:#633b00}.code-setup{border:1px solid #2b50a1;background:#dfe8ff;color:#183b88}.status-table{width:100%;margin-top:14px;border-collapse:collapse;font-size:13px}.status-table th,.status-table td{padding:10px;border:1px solid #d9dee8;text-align:left;vertical-align:top}.status-table th{background:#f3f6fa}.trouble{margin-top:24px}.trouble details{margin-top:8px;border:1px solid #d9dee8;border-radius:7px;background:#fff}.trouble summary{padding:12px 14px;font-weight:700;cursor:pointer}.trouble p{margin:0;padding:0 14px 14px;color:#536077}.files{margin-top:24px;padding-top:20px;border-top:1px solid #d9dee8}.footer{margin-top:24px;color:#667085;font-size:13px}@media(max-width:640px){.wrap{margin:12px auto;padding:0 10px 36px}.card{padding:20px}.phase{padding:15px}h1{font-size:24px}.status-table{display:block;overflow:auto}}
</style></head>
<body><main class="wrap"><article class="card">
<h1>Gửi email bằng Gmail miễn phí</h1>
<p class="lead">Hướng dẫn dành cho người chưa từng dùng Apps Script. Việc cài đặt chỉ làm một lần; các đợt thu sau sẽ nhập file dữ liệu mới vào chính Sheet này.</p>
<div class="summary"><span class="pill">Tổng ${total} dòng</span><span class="pill">${ready} dòng READY</span><span class="pill error">${errors} dòng cần kiểm tra</span></div>
<p class="notice"><strong>Hiện tại chưa có email nào được gửi.</strong> Gmail chỉ bắt đầu gửi sau khi chính bạn cấp quyền, gửi thử và xác nhận lượt gửi thật.</p>

<section class="phase"><div class="phase-head"><span class="phase-number">1</span><div><h2>Tạo bản Google Sheets</h2><small>Thực hiện trong Google Drive</small></div></div>
<ol class="steps"><li>Upload file <code>01_DANH_SACH_GUI.xlsx</code> lên Google Drive.</li><li>Nhấp đúp file và chọn mở bằng Google Sheets.</li><li>Trong thanh menu của Google Sheets, chọn <span class="button-name">File</span> → <span class="button-name">Save as Google Sheets</span>.</li><li>Google tạo một bản mới. Chuyển sang bản mới này và đóng tab file Excel cũ để tránh nhầm.</li></ol>
<p class="expected"><strong>Kết quả đúng:</strong> tên bản Google Sheets mới không còn đuôi <code>.xlsx</code>. Chỉ bản này mới dùng được Apps Script.</p></section>

<section class="phase"><div class="phase-head"><span class="phase-number">2</span><div><h2>Mở Apps Script, dán và lưu code</h2><small>Thực hiện từ bản Google Sheets mới</small></div></div>
<ol class="steps"><li>Trên thanh menu, chọn <span class="button-name">Extensions</span> → <span class="button-name">Apps Script</span>. Một tab mới sẽ mở.</li><li>Ở cột <strong>Files</strong> bên trái, bấm file <code>Code.gs</code>.</li><li>Bấm vào vùng code, nhấn <strong>⌘A</strong> trên Mac hoặc <strong>Ctrl+A</strong> trên Windows, rồi xoá code mẫu.</li><li>Mở file <code>02_CODE_GUI_EMAIL.gs</code> trong ZIP bằng trình soạn thảo text, copy toàn bộ và paste vào <code>Code.gs</code>.</li><li>Trên thanh phía trên code, bấm biểu tượng <strong>Save project</strong> hình đĩa mềm. Chờ chữ <strong>Saving project…</strong> đổi thành <strong>Saved to Drive</strong>.</li></ol>
<p class="expected"><strong>Kết quả đúng:</strong> file bên trái vẫn tên <code>Code.gs</code>, bên phải là code DEKISUGI và phía trên báo đã lưu vào Drive.</p></section>

<section class="phase emphasis"><div class="phase-head"><span class="phase-number">3</span><div><h2>Chạy setup và cấp quyền lần đầu</h2><small>Đây là bước chỉ làm một lần</small></div></div>
<p class="where">Nhìn thanh công cụ ngay phía trên vùng code:</p><div class="toolbar"><span>💾 Save</span><span class="run">▶ Run</span><span>Debug</span><span class="function-selector">onOpen ▼</span></div>
<ol class="steps"><li>Bấm danh sách tên hàm đang hiện <code class="code-onopen">onOpen</code>.</li><li>Trong danh sách vừa mở, chọn <code class="code-setup">setup</code>. Kiểm tra ô tên hàm đã đổi từ <code class="code-onopen">onOpen</code> thành <code class="code-setup">setup</code>.</li><li>Bấm nút <span class="button-name">▶ Run</span> nằm bên trái nút <strong>Debug</strong>.</li></ol>
<div class="auth-flow"><h3>Khi thấy “Authorization required”</h3><ol class="steps"><li>Bấm nút xanh <span class="button-name">Review permissions</span>, không bấm Cancel.</li><li>Chọn đúng tài khoản Gmail mà bạn muốn dùng để gửi email.</li><li>Nếu Google hiện <strong>Google hasn’t verified this app</strong>, bấm <span class="button-name">Advanced</span>, rồi bấm <strong>Go to &lt;tên project&gt; (unsafe)</strong>.</li><li>Chỉ làm bước trên khi thanh địa chỉ là <code>script.google.com</code>, project là project bạn vừa tạo và code vẫn nguyên bản từ file DEKISUGI. Từ “unsafe” ở đây báo project cá nhân chưa được Google xuất bản/xác minh; nó không tự chứng minh code an toàn.</li><li>Kiểm tra các quyền được hỏi chỉ liên quan đến Google Sheet hiện tại và gửi email, rồi bấm <span class="button-name">Continue</span> hoặc <span class="button-name">Allow</span>.</li><li>Quay lại tab Apps Script. Nhìn bảng <strong>Execution log</strong> ở phía dưới và chờ dòng <strong>Execution completed</strong>.</li></ol></div>
<p class="notice warning"><strong>Quy tắc an toàn:</strong> không nhập password ở trang nào ngoài <code>accounts.google.com</code>. Nếu không muốn cấp quyền, hãy bấm Cancel; sẽ không có email nào được gửi.</p>
<p class="expected"><strong>Kết quả đúng:</strong> hàm đang chọn là <code class="code-setup">setup</code>, Execution log kết thúc bằng <strong>Execution completed</strong>, không có dòng lỗi màu đỏ.</p></section>

<section class="phase"><div class="phase-head"><span class="phase-number">4</span><div><h2>Kiểm tra và kết nối Sheet</h2><small>Menu DEKISUGI xuất hiện sau khi reload</small></div></div>
<ol class="steps"><li>Quay lại tab Google Sheets và reload trang bằng <strong>⌘R</strong> trên Mac hoặc <strong>Ctrl+R</strong> trên Windows.</li><li>Chờ Sheet tải xong. Trên thanh menu sẽ có mục <span class="button-name">DEKISUGI Email</span>.</li><li>Kiểm tra menu có mục <strong>0. Nhập dữ liệu mới</strong>; đây là nút dùng cho các tháng sau.</li><li>Copy toàn bộ link Sheet trên thanh địa chỉ.</li><li>Quay lại QR Tool, bấm <strong>Tôi đã có Sheet Gmail</strong>, dán link rồi bấm <strong>Lưu Sheet của tôi</strong>.</li><li>Trở lại Sheet và bấm <strong>DEKISUGI Email → 1. Kiểm tra danh sách</strong>.</li><li>Đọc số dòng <code>READY</code>, <code>ERROR</code> và <code>SKIP</code>. Script chỉ gửi dòng <code>READY</code> có cột <code>Send</code> khác <code>NO</code>.</li></ol>
<p class="expected"><strong>Kết quả đúng:</strong> QR Tool hiện <strong>Đã kết nối Sheet Gmail</strong>; trong Sheet, popup báo số dòng sẵn sàng và dòng lỗi.</p></section>

<section class="phase"><div class="phase-head"><span class="phase-number">5</span><div><h2>Gửi thử trước, gửi thật sau</h2><small>Luôn gửi vào email của chính bạn trước</small></div></div>
<ol class="steps"><li>Mở sheet <code>EMAILS</code> ở phía dưới màn hình.</li><li>Bấm một ô bất kỳ trên một dòng có trạng thái <code>READY</code>.</li><li>Chọn <strong>DEKISUGI Email → 2. Gửi thử dòng đang chọn</strong>.</li><li>Trong hộp <strong>Gửi thử</strong>, nhập email của chính bạn rồi bấm <span class="button-name">OK</span>. Email phụ huynh trong Sheet không bị thay đổi.</li><li>Mở Gmail của bạn và kiểm tra đúng tên học sinh, số tiền, mã tham chiếu và ảnh QR.</li><li>Chỉ khi bản thử hoàn toàn đúng, quay lại Sheet và chọn <strong>DEKISUGI Email → 3. Gửi email chưa gửi</strong>.</li><li>Đọc số lượng email trong hộp <strong>Xác nhận gửi email thật</strong>. Bấm <span class="button-name">YES</span> để gửi hoặc <span class="button-name">NO</span> để dừng.</li></ol>
<p class="expected"><strong>Chống gửi trùng:</strong> dòng gửi thành công chuyển thành <code>SENT</code> và không tự gửi lại ở lượt sau. Dòng lỗi chuyển thành <code>ERROR</code> và ghi nguyên nhân ở cột <code>Error</code>.</p></section>

<section class="phase emphasis"><div class="phase-head"><span class="phase-number">✓</span><div><h2>Đợt thu tiếp theo: dùng lại Sheet này</h2><small>Không cài đặt lại</small></div></div>
<ol class="steps"><li>Import Excel mới vào QR Tool, map dữ liệu và kiểm tra QR/email như bình thường.</li><li>Bấm <strong>Tải dữ liệu tháng này</strong>, rồi bấm <strong>Mở Sheet Gmail của tôi</strong>.</li><li>Trong Sheet, chọn <strong>DEKISUGI Email → 0. Nhập dữ liệu mới</strong>.</li><li>Chọn file <code>DEKISUGI_GMAIL_DATA.json</code> vừa tải và bấm <strong>Thay bằng dữ liệu mới</strong>.</li><li>Gửi thử cho chính bạn rồi mới gửi hàng loạt.</li></ol>
<p class="expected"><strong>Không làm lại:</strong> không dán code, không chạy <code class="code-setup">setup</code> và không cấp quyền lại. Danh sách mới thay danh sách cũ; code và quyền Gmail được giữ nguyên.</p></section>

<section class="phase"><div class="phase-head"><span class="phase-number">✓</span><div><h2>Hiểu các trạng thái</h2><small>Không cần đoán dòng nào đã gửi</small></div></div>
<table class="status-table"><thead><tr><th>Trạng thái</th><th>Ý nghĩa</th><th>Bạn cần làm gì</th></tr></thead><tbody><tr><td><code>READY</code></td><td>Đủ dữ liệu và đang chờ gửi.</td><td>Có thể gửi thử hoặc gửi thật.</td></tr><tr><td><code>SENT</code></td><td>Đã gửi thành công.</td><td>Không sửa về READY nếu không thực sự muốn gửi lại.</td></tr><tr><td><code>ERROR</code></td><td>Không thể gửi; xem cột Error.</td><td>Sửa dữ liệu rồi chạy lại “1. Kiểm tra danh sách”.</td></tr><tr><td><code>SKIP</code></td><td>Dòng có Send = NO.</td><td>Đổi Send thành YES nếu muốn gửi.</td></tr><tr><td><code>SENDING</code></td><td>Đang xử lý hoặc lần trước bị gián đoạn.</td><td>Chạy lại “1. Kiểm tra danh sách” để khôi phục.</td></tr></tbody></table></section>

<section class="trouble"><h2>Nếu bị kẹt</h2><details><summary>Không thấy mục DEKISUGI Email trong Google Sheets</summary><p>Kiểm tra bạn đang ở bản Google Sheets mới, không phải file <code>.xlsx</code>. Sau đó xác nhận Apps Script đã báo <strong>Execution completed</strong> và reload lại Google Sheets.</p></details><details><summary>Bấm Run nhưng không thấy cửa sổ chọn tài khoản</summary><p>Kiểm tra biểu tượng popup bị chặn ở bên phải thanh địa chỉ Edge/Chrome. Cho phép popup từ <code>script.google.com</code>, rồi chọn lại hàm <code>setup</code> và bấm <strong>Run</strong>.</p></details><details><summary>Google báo “This app is blocked” hoặc chính sách quản trị không cho phép</summary><p>Đây là giới hạn của tài khoản Google/Workspace. Không cố vượt qua. Hãy dùng một tài khoản Gmail cá nhân thông thường hoặc liên hệ quản trị viên Workspace.</p></details><details><summary>Execution log có dòng lỗi màu đỏ</summary><p>Bấm vào dòng lỗi để đọc nội dung. Kiểm tra code đã được copy đầy đủ, Sheet có tab <code>EMAILS</code>, rồi bấm Save và chạy lại <code>setup</code>. Không gửi thật cho đến khi setup hoàn tất.</p></details><details><summary>Không nhận được email thử</summary><p>Kiểm tra Spam/Thư rác, email thử đã nhập và cột <code>Error</code>. Chạy <strong>1. Kiểm tra danh sách</strong> và xem quota còn lại trước khi thử lại.</p></details></section>

<section class="files"><h2>Ba file trong ZIP dùng để làm gì?</h2><ul><li><code>01_DANH_SACH_GUI.xlsx</code>: danh sách người nhận, nội dung email và QR riêng từng dòng.</li><li><code>02_CODE_GUI_EMAIL.gs</code>: code bạn tự paste và tự cấp quyền trong tài khoản Google của mình.</li><li><code>BAT_DAU_TU_DAY.html</code>: trang hướng dẫn đang đọc.</li></ul><p class="notice warning"><strong>Không gửi hoặc attach ba file này cho phụ huynh.</strong> Workbook chứa dữ liệu thanh toán và email người nhận.</p></section>
<p class="footer">Script gửi tối đa 90 email mỗi lượt/ngày và luôn chừa ít nhất 10 recipients trong quota Google báo còn lại. Google có thể thay đổi quota theo loại tài khoản.</p>
</article></main></body></html>`;
}

export async function createGmailFreeBundle(items, template) {
  const zip = new JSZip();
  const rows = gmailFreeRows(items, template);
  const ready = rows.filter((row) => row.Status === "READY").length;
  zip.file("01_DANH_SACH_GUI.xlsx", gmailFreeWorkbookArray(rows));
  zip.file("02_CODE_GUI_EMAIL.gs", GMAIL_FREE_APPS_SCRIPT);
  zip.file("BAT_DAU_TU_DAY.html", gmailFreeGuideHTML({ total: rows.length, ready, errors: rows.length - ready }));
  return zip.generateAsync({ type: "blob", mimeType: "application/zip", compression: "DEFLATE" });
}

export async function createGmailMergeBundle(items, template) {
  const zip = new JSZip();
  const eligible = items.filter((item) => !item.errors?.length && emailIsValid(item.email));
  const rejected = items
    .filter((item) => item.errors?.length || !emailIsValid(item.email))
    .map((item) => ({
      ...manifestRow(item),
      errors: item.errors?.length ? item.errors.join("; ") : (item.email ? "Email không hợp lệ" : "Thiếu email người nhận"),
    }));
  const fileById = new Map(eligible.map((item) => [item.id, { qr: `${safeFilename(itemBaseName(item))}.png` }]));
  const gmail = templateForGmailMerge(template);

  zip.file("gmail-mail-merge.xlsx", workbookArray(gmailRows(eligible, fileById)));
  zip.file("gmail-template.html", gmail.html);
  zip.file("template.json", JSON.stringify(template, null, 2));
  if (rejected.length > 0) zip.file("email-errors.csv", toCSV(rejected, MANIFEST_COLUMNS));
  zip.file("README.txt", [
    "DEKISUGI QR Tool - Gmail Mail Merge",
    "",
    "1. Upload gmail-mail-merge.xlsx lên Google Drive và mở bằng Google Sheets.",
    "2. Trong Gmail, bật Mail Merge và chọn Add from a spreadsheet.",
    "3. Chọn cột Email, FirstName và LastName; copy nội dung gmail-template.html vào email.",
    "4. Gmail không cho merge tag trong subject và không hỗ trợ attachment/QR khác nhau cho từng recipient.",
    "5. Nếu mỗi người nhận cần QR riêng, dùng email bulk ZIP (.eml) từ DEKISUGI QR Tool.",
    "",
    `Template có QR riêng: ${gmail.hasPerRecipientQR ? "Có - Gmail sẽ thay bằng cảnh báo text" : "Không"}`,
    "Gói này không gửi email và không kết nối tài khoản Google.",
  ].join("\r\n"));
  return zip.generateAsync({ type: "blob", mimeType: "application/zip", compression: "DEFLATE" });
}

export async function createEmailBundle(items, template, { from = "" } = {}) {
  const zip = new JSZip();
  const used = new Set();
  const eligible = [];
  const exportErrors = [];
  const providerRows = [];
  const providerJSON = [];
  const fileById = new Map();

  for (const item of items) {
    if (item.errors?.length || !item.qrData) {
      exportErrors.push({ ...manifestRow(item), errors: (item.errors || ["QR chưa sẵn sàng"]).join("; ") });
      continue;
    }
    if (!emailIsValid(item.email)) {
      exportErrors.push({ ...manifestRow(item), errors: item.email ? "Email không hợp lệ" : "Thiếu email người nhận" });
      continue;
    }

    const base = allocateName(itemBaseName(item), used);
    const qrFilename = `${base}.png`;
    const emlFilename = `${base}.eml`;
    const htmlFilename = `${base}.html`;
    const textFilename = `${base}.txt`;
    const contentId = `qr-${safeFilename(item.id || base)}`;
    const cidRendered = renderEmailTemplate(template, item, { qrSrc: `cid:${contentId}` });
    const fileRendered = renderEmailTemplate(template, item, { qrSrc: `../qr/${qrFilename}` });
    const eml = buildEml({
      to: item.email,
      from,
      subject: cidRendered.subject,
      html: cidRendered.html,
      text: cidRendered.text,
      qrBase64: qrBase64(item),
      contentId,
      qrFilename,
    });

    zip.file(`qr/${qrFilename}`, qrBase64(item), { base64: true });
    zip.file(`messages/${emlFilename}`, eml);
    zip.file(`emails/${htmlFilename}`, fileRendered.html);
    zip.file(`emails/${textFilename}`, fileRendered.text);
    fileById.set(item.id, { qr: qrFilename, eml: emlFilename, html: htmlFilename, text: textFilename });
    eligible.push(item);
    providerRows.push({
      to: item.email,
      subject: cidRendered.subject,
      studentCode: item.studentCode,
      studentName: item.studentName,
      schoolName: item.schoolName,
      cohort: item.cohort,
      year: item.year,
      className: item.className,
      parentName: item.parentName,
      htmlFile: `emails/${htmlFilename}`,
      textFile: `emails/${textFilename}`,
      qrFile: `qr/${qrFilename}`,
      emlFile: `messages/${emlFilename}`,
      billNumber: item.billNumber,
    });
    providerJSON.push(JSON.stringify({
      to: item.email,
      subject: cidRendered.subject,
      html: cidRendered.html,
      text: cidRendered.text,
      attachments: [{ filename: qrFilename, contentId, path: `qr/${qrFilename}`, inline: true }],
      metadata: {
        billNumber: item.billNumber,
        studentCode: item.studentCode,
        studentName: item.studentName,
        schoolName: item.schoolName,
        cohort: item.cohort,
        year: item.year,
        className: item.className,
        sourceRow: item.sourceRow,
      },
    }));
  }

  zip.file("template/template.json", JSON.stringify(template, null, 2));
  zip.file("template/template.html", template.html);
  zip.file("recipients/bulk-email.csv", toCSV(providerRows, [
    ["to", "to"],
    ["subject", "subject"],
    ["studentCode", "student_code"],
    ["studentName", "student_name"],
    ["schoolName", "school_name"],
    ["cohort", "cohort"],
    ["year", "year"],
    ["className", "class_name"],
    ["parentName", "parent_name"],
    ["htmlFile", "html_file"],
    ["textFile", "text_file"],
    ["qrFile", "qr_file"],
    ["emlFile", "eml_file"],
    ["billNumber", "bill_number"],
  ]));
  zip.file("recipients/bulk-email.jsonl", `${providerJSON.join("\n")}${providerJSON.length ? "\n" : ""}`);
  zip.file("manifest.csv", toCSV(items.map((item) => manifestRow(item, fileById.get(item.id)?.qr || "")), MANIFEST_COLUMNS));
  if (exportErrors.length > 0) zip.file("email-errors.csv", toCSV(exportErrors, MANIFEST_COLUMNS));
  zip.file("README.txt", [
    "DEKISUGI QR Tool - Email export bundle",
    "",
    "Gói này chỉ tạo nội dung và file nháp; nó không gửi email, không đăng nhập và không kết nối provider.",
    "messages/: email RFC 822 (.eml) có X-Unsent: 1 và QR inline CID.",
    "bulk-email.csv/jsonl: định dạng portable; cần map lại theo schema của email provider cụ thể.",
    "qr/: ảnh QR tương ứng từng người nhận.",
    "Dùng bundle Gmail miễn phí riêng nếu muốn gửi bằng Google Sheets + Apps Script.",
  ].join("\r\n"));

  return zip.generateAsync({ type: "blob", mimeType: "application/zip", compression: "DEFLATE" });
}
