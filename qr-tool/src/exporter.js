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
    studentName: item.studentName || "",
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
  ["studentName", "Học sinh"],
  ["parentName", "Phụ huynh"],
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
      StudentName: String(item.studentName || ""),
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
      metadata: { billNumber: item.billNumber, studentName: item.studentName, sourceRow: item.sourceRow },
    }));
  }

  const gmail = templateForGmailMerge(template);
  zip.file("template/template.json", JSON.stringify(template, null, 2));
  zip.file("template/template.html", template.html);
  zip.file("template/gmail-template.html", gmail.html);
  zip.file("recipients/gmail-mail-merge.xlsx", workbookArray(gmailRows(eligible, fileById)));
  zip.file("recipients/bulk-email.csv", toCSV(providerRows, [
    ["to", "to"],
    ["subject", "subject"],
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
    "recipients/gmail-mail-merge.xlsx: dữ liệu text cho Gmail Mail Merge.",
    "template/gmail-template.html: nội dung có @MergeTag để copy vào Gmail.",
    "Gmail Mail Merge không hỗ trợ QR/attachment khác nhau theo từng recipient; dùng messages/*.eml hoặc provider JSONL khi cần QR riêng.",
    "bulk-email.csv/jsonl: định dạng portable; cần map lại theo schema của email provider cụ thể.",
    "qr/: ảnh QR tương ứng từng người nhận.",
    "",
    `Template có QR riêng: ${gmail.hasPerRecipientQR ? "Có" : "Không"}`,
  ].join("\r\n"));

  return zip.generateAsync({ type: "blob", mimeType: "application/zip", compression: "DEFLATE" });
}
