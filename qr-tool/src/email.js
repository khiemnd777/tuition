import { formatVND } from "./vietqr.js";

export const EMAIL_TEMPLATE_VERSION = 1;

export const DEFAULT_EMAIL_TEMPLATE = {
  schemaVersion: EMAIL_TEMPLATE_VERSION,
  name: "Thông báo thanh toán",
  subject: "Thông báo thanh toán {{student_name}} - {{bill_number}}",
  html: `
    <p>Kính gửi <strong>{{parent_name}}</strong>,</p>
    <p>Thông tin thanh toán của <strong>{{student_name}}</strong>, lớp {{class_name}}:</p>
    {{payment_items}}
    <p><strong>Tổng cộng: {{amount}}</strong></p>
    <p>Ngân hàng: {{bank_name}}<br>Số tài khoản: {{bank_account}}<br>Mã tham chiếu: {{bill_number}}<br>Nội dung: {{payment_note}}</p>
    {{qr_image}}
    <p>Vui lòng kiểm tra thông tin trước khi thực hiện thanh toán.</p>
  `.trim(),
};

export const MERGE_FIELDS = [
  ["{{recipient_email}}", "Email"],
  ["{{student_name}}", "Học sinh"],
  ["{{parent_name}}", "Phụ huynh"],
  ["{{class_name}}", "Lớp"],
  ["{{amount}}", "Số tiền"],
  ["{{payment_items}}", "Khoản phí"],
  ["{{bank_name}}", "Ngân hàng"],
  ["{{bank_account}}", "Số tài khoản"],
  ["{{bill_number}}", "Mã tham chiếu"],
  ["{{payment_note}}", "Nội dung"],
  ["{{qr_image}}", "Ảnh QR"],
];

function escapeHTML(value = "") {
  return String(value)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

function paymentItemsHTML(item) {
  const items = item.paymentItems?.length
    ? item.paymentItems
    : [{ label: "Tổng phí cần thanh toán", labelEn: "Total fees due", amount: item.amount }];
  const rows = items.map((paymentItem) => `
    <tr>
      <td style="border:1px solid #cbd5e1;padding:8px;">${escapeHTML(paymentItem.label)}</td>
      <td style="border:1px solid #cbd5e1;padding:8px;text-align:right;">${escapeHTML(formatVND(paymentItem.amount))}</td>
    </tr>`).join("");
  return `<table role="presentation" cellpadding="0" cellspacing="0" style="width:100%;border-collapse:collapse;margin:16px 0;font-family:Arial,sans-serif;font-size:14px;"><thead><tr><th style="border:1px solid #cbd5e1;background:#f1f5f9;padding:8px;text-align:left;">Khoản phí</th><th style="border:1px solid #cbd5e1;background:#f1f5f9;padding:8px;text-align:right;">Số tiền</th></tr></thead><tbody>${rows}</tbody></table>`;
}

function tokenValues(item) {
  return {
    "{{recipient_email}}": item.email || "",
    "{{student_name}}": item.studentName || "",
    "{{parent_name}}": item.parentName || "Quý phụ huynh",
    "{{class_name}}": item.className || "-",
    "{{amount}}": formatVND(item.amount),
    "{{bank_name}}": item.bankName || "",
    "{{bank_account}}": item.bankAccount || "",
    "{{bill_number}}": item.billNumber || "",
    "{{payment_note}}": item.note || "",
  };
}

function replaceTextTokens(value, item, htmlMode) {
  let result = String(value || "");
  for (const [token, raw] of Object.entries(tokenValues(item))) {
    result = result.split(token).join(htmlMode ? escapeHTML(raw) : String(raw));
  }
  return result;
}

function htmlToText(html) {
  return String(html)
    .replace(/<style[\s\S]*?<\/style>/gi, "")
    .replace(/<br\s*\/?>/gi, "\n")
    .replace(/<\/(p|div|tr|table|h[1-6]|li)>/gi, "\n")
    .replace(/<li[^>]*>/gi, "- ")
    .replace(/<[^>]+>/g, "")
    .replace(/&nbsp;/g, " ")
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/\n{3,}/g, "\n\n")
    .trim();
}

export function renderEmailTemplate(template, item, { qrSrc = "" } = {}) {
  const subject = replaceTextTokens(template.subject, item, false).replace(/[\r\n]+/g, " ").trim();
  let body = replaceTextTokens(template.html, item, true);
  body = body.split("{{payment_items}}").join(paymentItemsHTML(item));
  const qr = qrSrc
    ? `<p style="text-align:center;margin:20px 0;"><img src="${escapeHTML(qrSrc)}" alt="QR thanh toán" style="display:block;width:240px;max-width:100%;height:auto;margin:0 auto;"></p>`
    : `<p><strong>[QR thanh toán]</strong></p>`;
  body = body.split("{{qr_image}}").join(qr);
  const html = `<!doctype html><html><body style="margin:0;padding:0;background:#ffffff;color:#172033;font-family:Arial,Helvetica,sans-serif;font-size:15px;line-height:1.55;"><div style="max-width:680px;margin:0 auto;padding:24px 18px;">${body}</div></body></html>`;
  return { subject, html, text: htmlToText(html) };
}

function utf8Base64(value) {
  const bytes = new TextEncoder().encode(String(value));
  let binary = "";
  for (let index = 0; index < bytes.length; index += 0x8000) {
    binary += String.fromCharCode(...bytes.subarray(index, index + 0x8000));
  }
  if (typeof btoa === "function") return btoa(binary);
  return Buffer.from(bytes).toString("base64");
}

function wrapBase64(value) {
  return String(value || "").replace(/\s/g, "").match(/.{1,76}/g)?.join("\r\n") || "";
}

function encodedHeader(value) {
  return `=?UTF-8?B?${utf8Base64(String(value).replace(/[\r\n]+/g, " "))}?=`;
}

export function buildEml({ to, from = "", subject, html, text, qrBase64, contentId, qrFilename }) {
  const token = String(contentId || "payment").replace(/[^A-Za-z0-9_-]/g, "") || "payment";
  const relatedBoundary = `dekisugi-related-${token}`;
  const alternativeBoundary = `dekisugi-alternative-${token}`;
  const headers = ["X-Unsent: 1"];
  if (from) headers.push(`From: ${String(from).replace(/[\r\n]+/g, " ")}`);
  headers.push(`To: ${String(to || "").replace(/[\r\n]+/g, " ")}`);
  headers.push(`Subject: ${encodedHeader(subject)}`);
  headers.push("MIME-Version: 1.0");
  headers.push(`Content-Type: multipart/related; boundary="${relatedBoundary}"`);

  const parts = [
    headers.join("\r\n"),
    "",
    `--${relatedBoundary}`,
    `Content-Type: multipart/alternative; boundary="${alternativeBoundary}"`,
    "",
    `--${alternativeBoundary}`,
    "Content-Type: text/plain; charset=UTF-8",
    "Content-Transfer-Encoding: base64",
    "",
    wrapBase64(utf8Base64(text)),
    `--${alternativeBoundary}`,
    "Content-Type: text/html; charset=UTF-8",
    "Content-Transfer-Encoding: base64",
    "",
    wrapBase64(utf8Base64(html)),
    `--${alternativeBoundary}--`,
  ];

  if (qrBase64) {
    const filename = String(qrFilename || "vietqr.png").replace(/["\r\n]/g, "");
    parts.push(
      `--${relatedBoundary}`,
      `Content-Type: image/png; name="${filename}"`,
      "Content-Transfer-Encoding: base64",
      `Content-ID: <${token}>`,
      `Content-Disposition: inline; filename="${filename}"`,
      "",
      wrapBase64(qrBase64),
    );
  }
  parts.push(`--${relatedBoundary}--`, "");
  return parts.join("\r\n");
}

const GMAIL_TOKENS = {
  "{{recipient_email}}": "@Email",
  "{{student_name}}": "@StudentName",
  "{{parent_name}}": "@ParentName",
  "{{class_name}}": "@ClassName",
  "{{amount}}": "@Amount",
  "{{bank_name}}": "@BankName",
  "{{bank_account}}": "@BankAccount",
  "{{bill_number}}": "@BillNumber",
  "{{payment_note}}": "@PaymentNote",
};

export function templateForGmailMerge(template) {
  let html = String(template.html || "");
  for (const [token, gmailToken] of Object.entries(GMAIL_TOKENS)) html = html.split(token).join(gmailToken);
  html = html.split("{{payment_items}}").join("@PaymentItems");
  const hasPerRecipientQR = html.includes("{{qr_image}}");
  html = html.split("{{qr_image}}").join("[Gmail Mail Merge không thể chèn QR riêng theo từng người nhận]");
  return { html, hasPerRecipientQR };
}

export function normalizeImportedTemplate(value) {
  if (!value || typeof value !== "object") throw new Error("Template JSON không hợp lệ");
  const name = String(value.name || "Template email").trim().slice(0, 100);
  const subject = String(value.subject || "").trim().slice(0, 300);
  const html = String(value.html || "").trim();
  if (!subject || !html) throw new Error("Template cần có tiêu đề và nội dung");
  return { schemaVersion: EMAIL_TEMPLATE_VERSION, name, subject, html };
}
