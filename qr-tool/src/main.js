import DOMPurify from "dompurify";
import QRCode from "qrcode";

import { BANKS } from "./banks.js";
import { COFFEE_TRANSFER, generateCoffeeVietQR } from "./coffee.js";
import {
  DEFAULT_EMAIL_TEMPLATE,
  MERGE_FIELDS,
  buildEml,
  normalizeImportedTemplate,
  renderEmailTemplate,
} from "./email.js";
import {
  GMAIL_FREE_APPS_SCRIPT,
  createEmailBundle,
  createGmailFreeBundle,
  createGmailSheetData,
  createQRBundle,
  isGmailSheetTemplateURL,
  safeFilename,
} from "./exporter.js";
import {
  clearPersonalGmailSheetURL,
  loadPersonalGmailSheetURL,
  savePersonalGmailSheetURL,
} from "./gmail-sheet.js";
import {
  IMPORT_FIELD_GROUPS,
  buildPaymentRows,
  previewValues,
  readSpreadsheet,
  suggestMapping,
  validateMapping,
} from "./importer.js";
import { buildQRItem, formatVND } from "./vietqr.js";
import "./styles.css";

const elements = Object.fromEntries([
  "coffeeButton", "coffeeDialog", "coffeeQR", "coffeeAmountLabel", "resetButton", "resetDialog",
  "spreadsheetFile", "dropZone", "fileSummary", "recipientSection", "recipientStatus",
  "recipientSharedFields", "recipientPerRowHelp", "recipientBank", "recipientAccount", "recipientAccountName",
  "recipientSetupSummary", "recipientSetupHint", "confirmRecipientSetup", "mappingSection",
  "mappingSummary", "mappingRows", "mappingHint", "applyMapping",
  "resultsSection", "totalCount", "validCount", "errorCount", "resultSearch", "resultFilter", "resultRows",
  "recipientExportGate", "recipientExportTitle", "recipientExportSummary", "recipientExportAccounts", "confirmRecipientExport",
  "qrEmpty", "qrDetail", "previewStatus", "previewName", "previewMeta", "previewQR", "previewBank",
  "previewAccount", "previewAccountNameRow", "previewAccountName", "previewAmount", "previewNote", "previewErrors", "previewPayload", "exportQR",
  "emailSection", "templateName", "emailFrom", "templateFile", "exportTemplate", "emailSubject",
  "mergeFields", "emailEditor", "linkURL", "insertLink", "emailPreviewRecipient", "emailPreviewFrame",
  "copySubject", "copyEmail", "copyQR", "downloadEml", "gmailFirstUseFlow", "gmailReturningFlow",
  "gmailSetupTitle", "gmailSetupDescription", "createGmailSheet", "connectGmailSheet",
  "downloadGmailUpgrade", "downloadGmailUpgradeConnected", "toggleGmailUpgradeSteps", "gmailUpgradeSteps",
  "gmailStep1", "gmailStep2", "gmailDownloadedState", "exportGmailData", "openGmailSheet",
  "changeGmailSheet", "disconnectGmailSheet", "exportGmailFree", "exportEmailBundle", "gmailGuideDialog",
  "gmailConnectDialog", "gmailConnectForm", "gmailSheetURL", "gmailSheetURLError", "saveGmailSheet", "toast",
].map((id) => [id, document.getElementById(id)]));

const GMAIL_SHEET_TEMPLATE_URL = String(import.meta.env.VITE_GMAIL_SHEET_TEMPLATE_URL || "").trim();
const browserStorage = (() => {
  try { return window.localStorage; } catch { return null; }
})();

const state = {
  table: null,
  mapping: {},
  items: [],
  selectedId: "",
  template: structuredClone(DEFAULT_EMAIL_TEMPLATE),
  gmailDataFilename: "",
  gmailSheetURL: loadPersonalGmailSheetURL(browserStorage),
  recipientConfig: null,
  recipientSetupConfirmed: false,
  recipientExportConfirmed: false,
  savedRange: null,
  toastTimer: null,
};

const SANITIZE_CONFIG = {
  ALLOWED_TAGS: ["p", "div", "br", "strong", "b", "em", "i", "u", "ul", "ol", "li", "a", "h1", "h2", "h3", "table", "thead", "tbody", "tr", "th", "td", "hr", "span"],
  ALLOWED_ATTR: ["href", "target", "rel", "style", "colspan", "rowspan", "role", "cellpadding", "cellspacing"],
  ALLOW_DATA_ATTR: false,
};

DOMPurify.addHook("afterSanitizeAttributes", (node) => {
  if (node.hasAttribute?.("href")) {
    const href = node.getAttribute("href").trim();
    if (!/^(https?:|mailto:)/i.test(href)) node.removeAttribute("href");
    else {
      node.setAttribute("target", "_blank");
      node.setAttribute("rel", "noopener noreferrer");
    }
  }
  if (node.hasAttribute?.("style") && /(url\s*\(|expression\s*\(|@import)/i.test(node.getAttribute("style"))) {
    node.removeAttribute("style");
  }
});

function sanitizeTemplateHTML(value) {
  return DOMPurify.sanitize(String(value || ""), SANITIZE_CONFIG);
}

function escapeHTML(value = "") {
  return String(value)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

function currentTemplate() {
  return {
    schemaVersion: 1,
    name: elements.templateName.value.trim() || "Template email",
    subject: elements.emailSubject.value.trim(),
    html: sanitizeTemplateHTML(elements.emailEditor.innerHTML),
  };
}

function showToast(message, type = "success") {
  clearTimeout(state.toastTimer);
  elements.toast.textContent = message;
  elements.toast.className = `toast${type === "error" ? " error" : ""}`;
  elements.toast.hidden = false;
  state.toastTimer = setTimeout(() => { elements.toast.hidden = true; }, 3600);
}

function setStatus(message, type = "neutral") {
  elements.fileSummary.textContent = message;
  elements.fileSummary.className = `status-pill ${type}`;
}

function setStep(step) {
  document.querySelectorAll(".stepper li").forEach((item) => {
    const value = Number(item.dataset.step);
    item.classList.toggle("active", value === step);
    item.classList.toggle("done", value < step);
  });
}

function setBusy(button, busy, busyLabel) {
  if (!button.dataset.label) button.dataset.label = button.textContent;
  button.disabled = busy;
  button.textContent = busy ? busyLabel : button.dataset.label;
}

function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = filename;
  document.body.append(anchor);
  anchor.click();
  anchor.remove();
  setTimeout(() => URL.revokeObjectURL(url), 2000);
}

function selectedItem() {
  return state.items.find((item) => item.id === state.selectedId) || null;
}

function initializeBanks() {
  elements.recipientBank.insertAdjacentHTML("beforeend", BANKS
    .slice()
    .sort((a, b) => a.shortName.localeCompare(b.shortName, "vi"))
    .map((bank) => `<option value="${bank.bin}">${escapeHTML(bank.shortName)} · ${bank.bin}</option>`)
    .join(""));
}

function currentRecipientMode() {
  return document.querySelector('input[name="recipientMode"]:checked')?.value === "per_row" ? "per_row" : "shared";
}

function bankByBIN(bin) {
  return BANKS.find((bank) => bank.bin === String(bin || "")) || null;
}

function readRecipientConfig() {
  return {
    mode: currentRecipientMode(),
    bankBin: elements.recipientBank.value,
    bankAccount: elements.recipientAccount.value.trim(),
    accountName: elements.recipientAccountName.value.trim(),
  };
}

function renderRecipientMode() {
  const perRow = currentRecipientMode() === "per_row";
  elements.recipientSharedFields.hidden = perRow;
  elements.recipientPerRowHelp.hidden = !perRow;
  elements.confirmRecipientSetup.textContent = perRow ? "Tiếp tục map tài khoản từ Excel" : "Kiểm tra và tiếp tục";
}

function invalidateRecipientSetup() {
  state.recipientConfig = null;
  state.recipientSetupConfirmed = false;
  state.recipientExportConfirmed = false;
  state.items = [];
  state.selectedId = "";
  elements.recipientStatus.textContent = "Chưa xác nhận";
  elements.recipientStatus.className = "status-pill neutral";
  elements.recipientSetupSummary.hidden = true;
  elements.mappingSection.hidden = true;
  elements.resultsSection.hidden = true;
  elements.emailSection.hidden = true;
  renderRecipientMode();
  renderExportAvailability();
  if (state.table) setStep(2);
}

function validateSharedRecipient(config) {
  const errors = [];
  if (!bankByBIN(config.bankBin)) errors.push("Hãy chọn ngân hàng nhận tiền");
  if (!config.bankAccount) errors.push("Hãy nhập số tài khoản nhận tiền");
  else if (!/^[A-Za-z0-9]+$/.test(config.bankAccount)) errors.push("Số tài khoản chỉ được chứa chữ và số");
  else if (config.bankAccount.length > 19) errors.push("Số tài khoản tối đa 19 ký tự");
  return errors;
}

function restorePerRowRecipientSuggestions() {
  const suggestions = suggestMapping(state.table?.headers || []);
  for (const header of state.table?.headers || []) {
    if (["bank_bin", "bank_account"].includes(suggestions[header]) && !state.mapping[header]) {
      state.mapping[header] = suggestions[header];
    }
  }
}

function confirmRecipientSetup() {
  if (!state.table) return showToast("Hãy chọn file Excel hoặc CSV trước", "error");
  const config = readRecipientConfig();
  const errors = config.mode === "shared" ? validateSharedRecipient(config) : [];
  if (errors.length) return showToast(errors.join(" · "), "error");

  state.recipientConfig = config;
  state.recipientSetupConfirmed = true;
  state.recipientExportConfirmed = false;
  if (config.mode === "per_row") restorePerRowRecipientSuggestions();

  elements.recipientStatus.textContent = "Đã nhập thông tin";
  elements.recipientStatus.className = "status-pill success";
  elements.recipientSetupSummary.hidden = false;
  if (config.mode === "shared") {
    const bank = bankByBIN(config.bankBin);
    elements.recipientSetupSummary.innerHTML = `<strong>Sẽ dùng một tài khoản cho toàn bộ QR:</strong><span>${escapeHTML(bank?.shortName || config.bankBin)} · ${escapeHTML(config.bankAccount)}${config.accountName ? ` · ${escapeHTML(config.accountName)}` : ""}</span>`;
  } else {
    elements.recipientSetupSummary.innerHTML = "<strong>Đã chọn chế độ nâng cao:</strong><span>Tài khoản nhận tiền sẽ được đọc riêng từ từng dòng Excel.</span>";
  }

  elements.mappingSection.hidden = false;
  renderMapping();
  setStep(3);
  elements.mappingSection.scrollIntoView({ behavior: "smooth", block: "start" });
}

function recipientFieldGroups() {
  if (state.recipientConfig?.mode === "per_row") return IMPORT_FIELD_GROUPS;
  return IMPORT_FIELD_GROUPS.map((group) => ({
    ...group,
    fields: group.fields.filter((field) => !["bank_bin", "bank_account"].includes(field.key)),
  })).filter((group) => group.fields.length > 0);
}

function renderRecipientExportGate() {
  const validItems = state.items.filter((item) => !item.errors.length);
  const accounts = new Map();
  validItems.forEach((item) => {
    const key = `${item.bankBin}|${item.bankAccount}`;
    if (!accounts.has(key)) accounts.set(key, item);
  });

  const confirmed = state.recipientExportConfirmed && validItems.length > 0;
  elements.recipientExportGate.classList.toggle("is-confirmed", confirmed);
  elements.recipientExportTitle.textContent = confirmed ? "Đã xác nhận tài khoản nhận tiền" : "Xác nhận tài khoản trước khi xuất";
  elements.recipientExportSummary.textContent = state.recipientConfig?.mode === "per_row"
    ? `${validItems.length} mã QR hợp lệ đang dùng ${accounts.size} tài khoản nhận tiền từ Excel.`
    : `${validItems.length} mã QR hợp lệ sẽ chuyển về một tài khoản duy nhất.`;

  const visibleAccounts = [...accounts.values()].slice(0, 6);
  elements.recipientExportAccounts.innerHTML = visibleAccounts.map((item) => {
    const bank = bankByBIN(item.bankBin);
    const owner = item.accountName ? ` · ${escapeHTML(item.accountName)}` : "";
    return `<li><strong>${escapeHTML(bank?.shortName || item.bankBin)}</strong><span>${escapeHTML(item.bankAccount)}${owner}</span></li>`;
  }).join("") + (accounts.size > visibleAccounts.length ? `<li><strong>Và ${accounts.size - visibleAccounts.length} tài khoản khác</strong><span>Kiểm tra từng dòng trước khi xác nhận.</span></li>` : "");

  elements.confirmRecipientExport.disabled = validItems.length === 0 || confirmed;
  elements.confirmRecipientExport.textContent = confirmed ? "✓ Đã cho phép xuất" : "Tôi đã kiểm tra, cho phép xuất";
}

function confirmRecipientExportAccounts() {
  if (!state.recipientSetupConfirmed) return showToast("Hãy xác nhận cách nhận tiền trước", "error");
  if (!state.items.some((item) => !item.errors.length)) return showToast("Chưa có mã QR hợp lệ để xác nhận", "error");
  state.recipientExportConfirmed = true;
  renderRecipientExportGate();
  renderExportAvailability();
  showToast("Đã xác nhận tài khoản nhận tiền. Bạn có thể xuất QR và email.");
}

function requireRecipientExportConfirmation() {
  if (state.recipientExportConfirmed) return true;
  showToast("Hãy kiểm tra và xác nhận tài khoản nhận tiền trước khi xuất", "error");
  elements.recipientExportGate?.scrollIntoView({ behavior: "smooth", block: "center" });
  return false;
}

function renderExportAvailability() {
  const allowed = state.recipientExportConfirmed;
  const hasSelected = Boolean(selectedItem());
  const validCount = state.items.filter((item) => !item.errors.length).length;
  elements.exportQR.disabled = !allowed || validCount === 0;
  elements.exportEmailBundle.disabled = !allowed;
  elements.exportGmailData.disabled = !allowed;
  elements.copySubject.disabled = !allowed || !hasSelected;
  elements.copyEmail.disabled = !allowed || !hasSelected;
  elements.copyQR.disabled = !allowed || !selectedItem()?.qrData;
  elements.downloadEml.disabled = !allowed || !selectedItem()?.qrData;
}

function initializeTemplate() {
  elements.templateName.value = state.template.name;
  elements.emailSubject.value = state.template.subject;
  elements.emailEditor.innerHTML = sanitizeTemplateHTML(state.template.html);
  elements.mergeFields.innerHTML = MERGE_FIELDS.map(([token, label]) =>
    `<button class="token-button" type="button" data-token="${escapeHTML(token)}" title="${escapeHTML(token)}">+ ${escapeHTML(label)}</button>`,
  ).join("");
}

function renderMapping() {
  if (!state.table) return;
  const fieldGroups = recipientFieldGroups();
  elements.mappingSummary.textContent = `${state.table.headers.length} cột · ${state.table.records.length} dòng`;
  elements.mappingRows.innerHTML = state.table.headers.map((header) => {
    const selected = state.mapping[header] || "";
    const options = [
      `<option value="">Bỏ qua</option>`,
      ...fieldGroups.map((group) => `<optgroup label="${escapeHTML(group.label)}">${group.fields.map((field) =>
        `<option value="${field.key}" ${selected === field.key ? "selected" : ""}>${escapeHTML(field.label)}${field.required ? " *" : ""}</option>`,
      ).join("")}</optgroup>`),
    ].join("");
    return `<tr><td><strong>${escapeHTML(header)}</strong></td><td><select data-source="${escapeHTML(header)}">${options}</select></td><td>${escapeHTML(previewValues(state.table, header) || "-")}</td></tr>`;
  }).join("");
  updateMappingHint();
}

function collectMapping() {
  const mapping = {};
  elements.mappingRows.querySelectorAll("select[data-source]").forEach((select) => { mapping[select.dataset.source] = select.value; });
  state.mapping = mapping;
  return mapping;
}

function updateMappingHint() {
  const mapping = collectMapping();
  const mapped = Object.values(mapping).filter(Boolean).length;
  const fees = Object.values(mapping).filter((value) => value === "fee_item").length;
  const errors = validateMapping(mapping);
  const recipientNote = state.recipientConfig?.mode === "per_row"
    ? " · Tài khoản lấy riêng theo từng dòng Excel"
    : " · Cột tài khoản trong Excel được bỏ qua";
  elements.mappingHint.textContent = errors.length
    ? errors.join(" · ")
    : `${mapped}/${state.table?.headers.length || 0} cột đã map${fees ? ` · ${fees} khoản phí custom` : ""}${recipientNote}.`;
  elements.mappingHint.style.color = errors.length ? "var(--danger)" : "";
}

function handleMappingChange() {
  updateMappingHint();
  if (!state.items.length) return;
  state.items = [];
  state.selectedId = "";
  state.recipientExportConfirmed = false;
  elements.resultsSection.hidden = true;
  elements.emailSection.hidden = true;
  renderExportAvailability();
  setStep(3);
  showToast("Dữ liệu map đã thay đổi. Hãy bấm “Kiểm tra và sinh QR” lại.");
}

async function handleSpreadsheet(file) {
  if (!file) return;
  setStatus("Đang đọc file…", "busy");
  try {
    state.table = await readSpreadsheet(file);
    state.mapping = suggestMapping(state.table.headers);
    state.items = [];
    state.selectedId = "";
    state.recipientConfig = null;
    state.recipientSetupConfirmed = false;
    state.recipientExportConfirmed = false;
    elements.recipientSection.hidden = false;
    elements.recipientStatus.textContent = "Chưa xác nhận";
    elements.recipientStatus.className = "status-pill neutral";
    elements.recipientSetupSummary.hidden = true;
    elements.mappingSection.hidden = true;
    elements.resultsSection.hidden = true;
    elements.emailSection.hidden = true;
    renderRecipientMode();
    renderExportAvailability();
    setStatus(`${file.name} · ${state.table.records.length} dòng`, "success");
    setStep(2);
    elements.recipientSection.scrollIntoView({ behavior: "smooth", block: "start" });
  } catch (error) {
    setStatus("Không đọc được file", "danger");
    showToast(error.message || "Không đọc được bảng dữ liệu", "error");
  }
}

async function generateItems() {
  if (!state.recipientSetupConfirmed || !state.recipientConfig) {
    showToast("Hãy xác nhận tài khoản nhận tiền trước khi map dữ liệu", "error");
    elements.recipientSection.scrollIntoView({ behavior: "smooth", block: "start" });
    return;
  }
  const mapping = collectMapping();
  const errors = validateMapping(mapping);
  if (state.recipientConfig.mode === "per_row") {
    const targets = Object.values(mapping);
    if (!targets.includes("bank_bin")) errors.push("Chế độ tài khoản theo dòng cần map field BIN ngân hàng");
    if (!targets.includes("bank_account")) errors.push("Chế độ tài khoản theo dòng cần map field Số tài khoản");
  }
  if (errors.length > 0) {
    showToast(errors.join(" · "), "error");
    return;
  }
  setBusy(elements.applyMapping, true, "Đang sinh QR…");
  try {
    const rows = buildPaymentRows(state.table, mapping, state.recipientConfig);
    const items = [];
    for (let index = 0; index < rows.length; index += 1) {
      const item = buildQRItem(rows[index]);
      if (!item.errors.length) {
        try {
          item.qrData = await QRCode.toDataURL(item.vietqr, { width: 512, margin: 2, errorCorrectionLevel: "M", type: "image/png" });
        } catch {
          item.errors.push("Không render được ảnh QR");
        }
      }
      items.push(item);
      if (index > 0 && index % 25 === 0) elements.applyMapping.textContent = `Đang sinh QR ${index}/${rows.length}…`;
    }
    state.items = items;
    state.selectedId = items.find((item) => !item.errors.length)?.id || items[0]?.id || "";
    state.recipientExportConfirmed = false;
    elements.resultsSection.hidden = false;
    elements.emailSection.hidden = false;
    renderRecipientExportGate();
    renderResults();
    renderSelected();
    renderExportAvailability();
    setStep(4);
    elements.resultsSection.scrollIntoView({ behavior: "smooth", block: "start" });
    showToast(`Đã xử lý ${items.length} dòng hoàn toàn trên thiết bị`);
  } catch (error) {
    showToast(error.message || "Không thể xử lý dữ liệu", "error");
  } finally {
    setBusy(elements.applyMapping, false);
  }
}

function filteredItems() {
  const query = elements.resultSearch.value.trim().toLocaleLowerCase("vi");
  const filter = elements.resultFilter.value;
  return state.items.filter((item) => {
    const ready = !item.errors.length;
    if (filter === "ready" && !ready) return false;
    if (filter === "error" && ready) return false;
    if (!query) return true;
    return [
      item.studentCode,
      item.studentName,
      item.schoolName,
      item.cohort,
      item.year,
      item.className,
      item.parentName,
      item.email,
      item.billNumber,
    ]
      .some((value) => String(value || "").toLocaleLowerCase("vi").includes(query));
  });
}

function renderResults() {
  const valid = state.items.filter((item) => !item.errors.length).length;
  const errors = state.items.length - valid;
  elements.totalCount.textContent = `${state.items.length} dòng`;
  elements.validCount.textContent = `${valid} hợp lệ`;
  elements.errorCount.textContent = `${errors} lỗi`;
  renderExportAvailability();
  const items = filteredItems();
  elements.resultRows.innerHTML = items.length ? items.map((item) => {
    const ready = !item.errors.length;
    const context = [
      item.studentCode ? `Mã: ${item.studentCode}` : "",
      item.schoolName,
      item.cohort ? `Niên khóa: ${item.cohort}` : "",
      item.year ? `Năm: ${item.year}` : "",
      item.className ? `Lớp: ${item.className}` : "",
    ].filter(Boolean).join(" · ");
    return `<tr data-id="${escapeHTML(item.id)}" class="${item.id === state.selectedId ? "selected" : ""}">
      <td>${item.sourceRow || "-"}</td>
      <td><span class="status-pill ${ready ? "success" : "danger"}">${ready ? "Sẵn sàng" : "Có lỗi"}</span></td>
      <td><span class="row-name">${escapeHTML(item.studentName || "Không có tên")}</span><span class="row-sub">${escapeHTML(context || "-")}</span></td>
      <td>${escapeHTML(item.email || "-")}</td>
      <td>${escapeHTML(formatVND(item.amount))}</td>
      <td>${escapeHTML(item.billNumber)}</td>
    </tr>`;
  }).join("") : `<tr><td colspan="6">Không có dòng phù hợp bộ lọc.</td></tr>`;
}

function renderSelected() {
  const item = selectedItem();
  if (!item) {
    elements.qrEmpty.hidden = false;
    elements.qrDetail.hidden = true;
    renderEmailPreview();
    renderExportAvailability();
    return;
  }
  const ready = !item.errors.length;
  elements.qrEmpty.hidden = true;
  elements.qrDetail.hidden = false;
  elements.previewStatus.textContent = ready ? "Sẵn sàng" : "Có lỗi";
  elements.previewStatus.className = `status-pill ${ready ? "success" : "danger"}`;
  elements.previewName.textContent = item.studentName || "Không có tên";
  elements.previewMeta.textContent = [
    item.studentCode ? `Mã ${item.studentCode}` : "",
    item.schoolName,
    item.cohort ? `Niên khóa ${item.cohort}` : "",
    item.year ? `Năm ${item.year}` : "",
    item.className ? `Lớp ${item.className}` : "",
    item.email,
    `Dòng ${item.sourceRow || "-"}`,
  ].filter(Boolean).join(" · ");
  elements.previewQR.hidden = !item.qrData;
  elements.previewQR.src = item.qrData || "";
  elements.previewBank.textContent = item.bankName ? `${item.bankName} · ${item.bankBin}` : item.bankBin || "-";
  elements.previewAccount.textContent = item.bankAccount || "-";
  elements.previewAccountNameRow.hidden = !item.accountName;
  elements.previewAccountName.textContent = item.accountName || "";
  elements.previewAmount.textContent = formatVND(item.amount);
  elements.previewNote.textContent = item.note || "-";
  elements.previewErrors.hidden = ready;
  elements.previewErrors.textContent = item.errors.join(" · ");
  elements.previewPayload.textContent = item.vietqr || "Chưa có payload";
  elements.emailPreviewRecipient.textContent = item.email || "Thiếu email";
  renderEmailPreview();
  renderExportAvailability();
}

function emailPreviewStatusDocument(message, isError = false) {
  const color = isError ? "#b42318" : "#667085";
  return `<!doctype html><html lang="vi"><body style="margin:0;padding:24px;background:#ffffff;color:${color};font-family:Arial,Helvetica,sans-serif;font-size:14px;line-height:1.5;">${escapeHTML(message)}</body></html>`;
}

function setEmailPreviewDocument(html) {
  elements.emailPreviewFrame.srcdoc = html;
}

function renderEmailPreview() {
  const item = selectedItem();
  if (!item) {
    setEmailPreviewDocument(emailPreviewStatusDocument("Chọn một dòng để preview email."));
    return;
  }
  try {
    const template = currentTemplate();
    state.template = template;
    const rendered = renderEmailTemplate(template, item, { qrSrc: item.qrData || "" });
    const previewHTML = DOMPurify.sanitize(rendered.html, {
      WHOLE_DOCUMENT: true,
      ADD_TAGS: ["html", "head", "body", "img"],
      ADD_ATTR: ["src", "alt", "width", "height"],
      ALLOW_DATA_ATTR: false,
    });
    if (!previewHTML.trim()) throw new Error("Empty preview document");
    setEmailPreviewDocument(previewHTML);
  } catch {
    setEmailPreviewDocument(emailPreviewStatusDocument("Preview không thể hiển thị. Vui lòng kiểm tra lại template.", true));
  }
}

function saveEditorRange() {
  const selection = window.getSelection();
  if (selection?.rangeCount && elements.emailEditor.contains(selection.anchorNode)) state.savedRange = selection.getRangeAt(0).cloneRange();
}

function restoreEditorRange() {
  elements.emailEditor.focus();
  if (!state.savedRange) return;
  const selection = window.getSelection();
  selection.removeAllRanges();
  selection.addRange(state.savedRange);
}

function insertToken(token) {
  restoreEditorRange();
  document.execCommand("insertText", false, token);
  saveEditorRange();
  renderEmailPreview();
}

async function copyText(value, successMessage) {
  try {
    await navigator.clipboard.writeText(value);
  } catch {
    const input = document.createElement("textarea");
    input.value = value;
    input.style.position = "fixed";
    input.style.opacity = "0";
    document.body.append(input);
    input.select();
    document.execCommand("copy");
    input.remove();
  }
  showToast(successMessage);
}

async function copyRichEmail() {
  if (!requireRecipientExportConfirmation()) return;
  const item = selectedItem();
  if (!item) return showToast("Chưa chọn recipient", "error");
  const rendered = renderEmailTemplate(currentTemplate(), item, { qrSrc: item.qrData || "" });
  try {
    if (!window.ClipboardItem || !navigator.clipboard?.write) throw new Error("unsupported");
    await navigator.clipboard.write([new ClipboardItem({
      "text/html": new Blob([rendered.html], { type: "text/html" }),
      "text/plain": new Blob([rendered.text], { type: "text/plain" }),
    })]);
    showToast("Đã copy rich email; nếu client bỏ ảnh, dùng nút Copy QR");
  } catch {
    await copyText(rendered.text, "Trình duyệt không hỗ trợ rich clipboard; đã copy bản text");
  }
}

function pngBlobFromDataURL(dataURL) {
  const base64 = String(dataURL).split(",")[1] || "";
  const binary = atob(base64);
  const bytes = new Uint8Array(binary.length);
  for (let index = 0; index < binary.length; index += 1) bytes[index] = binary.charCodeAt(index);
  return new Blob([bytes], { type: "image/png" });
}

async function copySelectedQR() {
  if (!requireRecipientExportConfirmation()) return;
  const item = selectedItem();
  if (!item?.qrData) return showToast("Dòng đang chọn chưa có QR hợp lệ", "error");
  try {
    await navigator.clipboard.write([new ClipboardItem({ "image/png": pngBlobFromDataURL(item.qrData) })]);
    showToast("Đã copy ảnh QR");
  } catch {
    downloadBlob(pngBlobFromDataURL(item.qrData), `${safeFilename(item.billNumber || item.id)}.png`);
    showToast("Clipboard ảnh không khả dụng; đã tải QR xuống máy");
  }
}

function downloadSelectedEml() {
  if (!requireRecipientExportConfirmation()) return;
  const item = selectedItem();
  if (!item?.qrData) return showToast("Dòng đang chọn chưa có QR hợp lệ", "error");
  if (!item.email) return showToast("Dòng đang chọn thiếu email người nhận", "error");
  const contentId = `qr-${safeFilename(item.id || item.billNumber)}`;
  const rendered = renderEmailTemplate(currentTemplate(), item, { qrSrc: `cid:${contentId}` });
  const eml = buildEml({
    to: item.email,
    from: elements.emailFrom.value.trim(),
    subject: rendered.subject,
    html: rendered.html,
    text: rendered.text,
    qrBase64: String(item.qrData).replace(/^data:image\/png;base64,/, ""),
    contentId,
    qrFilename: `${safeFilename(item.billNumber || item.id)}.png`,
  });
  downloadBlob(new Blob([eml], { type: "message/rfc822;charset=utf-8" }), `${safeFilename([item.billNumber, item.studentName].filter(Boolean).join("-"))}.eml`);
  showToast("Đã tạo email nháp .eml; không có email nào được gửi");
}

async function exportQRZip() {
  if (!requireRecipientExportConfirmation()) return;
  setBusy(elements.exportQR, true, "Đang đóng gói…");
  try {
    const blob = await createQRBundle(state.items);
    downloadBlob(blob, `dekisugi-qr-${new Date().toISOString().slice(0, 10)}.zip`);
    showToast("Đã tạo ZIP QR và manifest");
  } catch (error) {
    showToast(error.message || "Không tạo được ZIP QR", "error");
  } finally { setBusy(elements.exportQR, false); }
}

async function exportEmailZip() {
  if (!requireRecipientExportConfirmation()) return;
  setBusy(elements.exportEmailBundle, true, "Đang đóng gói…");
  try {
    const blob = await createEmailBundle(state.items, currentTemplate(), { from: elements.emailFrom.value.trim() });
    downloadBlob(blob, `dekisugi-email-bulk-${new Date().toISOString().slice(0, 10)}.zip`);
    showToast("Đã tạo email bulk ZIP; không có email nào được gửi");
  } catch (error) {
    showToast(error.message || "Không tạo được email bundle", "error");
  } finally { setBusy(elements.exportEmailBundle, false); }
}

async function exportGmailFreeZip(button = elements.exportGmailFree) {
  if (!requireRecipientExportConfirmation()) return;
  setBusy(button, true, "Đang tạo bộ Gmail…");
  try {
    const blob = await createGmailFreeBundle(state.items, currentTemplate());
    downloadBlob(blob, `dekisugi-gmail-free-${new Date().toISOString().slice(0, 10)}.zip`);
    showToast("Đã tạo bộ gửi Gmail miễn phí; chưa có email nào được gửi");
    elements.gmailGuideDialog.showModal();
  } catch (error) {
    showToast(error.message || "Không tạo được bộ gửi Gmail", "error");
  } finally { setBusy(button, false); }
}

function exportGmailSheetDataFile() {
  if (!requireRecipientExportConfirmation()) return;
  try {
    const payload = createGmailSheetData(state.items, currentTemplate());
    const blob = new Blob([JSON.stringify(payload)], { type: "application/json;charset=utf-8" });
    const filename = `DEKISUGI_GMAIL_DATA_${new Date().toISOString().slice(0, 10)}.json`;
    downloadBlob(blob, filename);
    state.gmailDataFilename = filename;
    renderGmailSheetFlow();
    showToast(`Đã tải ${payload.summary.total} dòng. Bây giờ tiếp tục bước 2: mở Google Sheet gửi Gmail.`);
  } catch (error) {
    showToast(error.message || "Không tạo được dữ liệu Gmail", "error");
  }
}

function renderGmailSheetFlow() {
  const connected = Boolean(state.gmailSheetURL);
  const templateAvailable = isGmailSheetTemplateURL(GMAIL_SHEET_TEMPLATE_URL);
  elements.gmailFirstUseFlow.hidden = connected;
  elements.gmailReturningFlow.hidden = !connected;

  if (!connected) {
    elements.gmailSetupTitle.textContent = templateAvailable
      ? "Tạo công cụ Gmail của tôi"
      : "Thiết lập Gmail một lần";
    elements.gmailSetupDescription.textContent = templateAvailable
      ? "Google sẽ tạo một bản Sheet riêng thuộc tài khoản của bạn. Sau lần này, các tháng sau bạn chỉ nhập dữ liệu mới."
      : "Tải bộ cài để tạo một Google Sheet dùng lâu dài. Sau lần này, các tháng sau bạn không phải dán code hoặc chạy setup lại.";
    elements.createGmailSheet.textContent = templateAvailable
      ? "Tạo Sheet Gmail của tôi"
      : "Tải bộ cài Gmail lần đầu";
    return;
  }

  const downloaded = Boolean(state.gmailDataFilename);
  elements.gmailStep1.classList.toggle("is-current", !downloaded);
  elements.gmailStep1.classList.toggle("is-complete", downloaded);
  elements.gmailStep2.classList.toggle("is-pending", !downloaded);
  elements.gmailStep2.classList.toggle("is-current", downloaded);
  elements.exportGmailData.textContent = downloaded ? "Tải lại dữ liệu" : "Tải dữ liệu tháng này";
  elements.gmailDownloadedState.hidden = !downloaded;
  if (downloaded) {
    elements.gmailDownloadedState.textContent = `Đã tải: ${state.gmailDataFilename}`;
  }
}

function startGmailSetup() {
  if (isGmailSheetTemplateURL(GMAIL_SHEET_TEMPLATE_URL)) {
    window.open(GMAIL_SHEET_TEMPLATE_URL, "_blank", "noopener,noreferrer");
    showToast("Google Sheets sẽ yêu cầu tạo bản sao. Sau khi thiết lập, quay lại bấm “Tôi đã có Sheet Gmail”.");
    return;
  }
  exportGmailFreeZip(elements.createGmailSheet);
}

function openGmailConnectDialog() {
  elements.gmailSheetURL.value = state.gmailSheetURL;
  elements.gmailSheetURLError.hidden = true;
  elements.gmailSheetURLError.textContent = "";
  elements.gmailConnectDialog.showModal();
  setTimeout(() => elements.gmailSheetURL.focus(), 0);
}

function saveGmailSheetConnection() {
  try {
    state.gmailSheetURL = savePersonalGmailSheetURL(browserStorage, elements.gmailSheetURL.value);
    elements.gmailConnectDialog.close();
    renderGmailSheetFlow();
    showToast("Đã ghi nhớ Sheet Gmail. Những lần sau bạn không cần chạy setup lại.");
  } catch (error) {
    elements.gmailSheetURLError.textContent = error.message || "Không lưu được link Google Sheet";
    elements.gmailSheetURLError.hidden = false;
  }
}

function disconnectGmailSheet() {
  clearPersonalGmailSheetURL(browserStorage);
  state.gmailSheetURL = "";
  state.gmailDataFilename = "";
  renderGmailSheetFlow();
  showToast("Đã ngắt kết nối trên trình duyệt này. Google Sheet của bạn không bị xoá.");
}

function downloadGmailUpgradeScript() {
  downloadBlob(
    new Blob([GMAIL_FREE_APPS_SCRIPT], { type: "text/plain;charset=utf-8" }),
    "DEKISUGI_NANG_CAP_CODE.gs",
  );
  showToast("Đã tải code mới. Thay nội dung Code.gs trong Sheet cũ rồi reload Google Sheets.");
}

function toggleGmailUpgradeInstructions() {
  const shouldShow = elements.gmailUpgradeSteps.hidden;
  elements.gmailUpgradeSteps.hidden = !shouldShow;
  elements.toggleGmailUpgradeSteps.setAttribute("aria-expanded", String(shouldShow));
  elements.toggleGmailUpgradeSteps.textContent = shouldShow ? "Ẩn hướng dẫn" : "Xem hướng dẫn 4 bước";
}

function openPersonalGmailSheet() {
  if (!state.gmailSheetURL) {
    openGmailConnectDialog();
    return;
  }
  window.open(state.gmailSheetURL, "_blank", "noopener,noreferrer");
  showToast(state.gmailDataFilename
    ? "Trong Sheet, chọn DEKISUGI Email → 0. Nhập dữ liệu mới."
    : "Đã mở Sheet Gmail của bạn. Hãy tải dữ liệu mới trước khi gửi đợt thu mới.");
}

function exportTemplateFile() {
  const template = currentTemplate();
  downloadBlob(new Blob([JSON.stringify(template, null, 2)], { type: "application/json" }), `${safeFilename(template.name)}.json`);
  showToast("Đã export template; ứng dụng không lưu bản sao trong browser");
}

async function importTemplateFile(file) {
  if (!file) return;
  try {
    const imported = normalizeImportedTemplate(JSON.parse(await file.text()));
    imported.html = sanitizeTemplateHTML(imported.html);
    state.template = imported;
    initializeTemplate();
    renderEmailPreview();
    showToast("Đã import và sanitize template");
  } catch (error) {
    showToast(error.message || "Không import được template", "error");
  } finally {
    elements.templateFile.value = "";
  }
}

function resetSession() {
  state.table = null;
  state.mapping = {};
  state.items = [];
  state.selectedId = "";
  state.template = structuredClone(DEFAULT_EMAIL_TEMPLATE);
  state.gmailDataFilename = "";
  state.recipientConfig = null;
  state.recipientSetupConfirmed = false;
  state.recipientExportConfirmed = false;
  elements.spreadsheetFile.value = "";
  elements.recipientBank.value = "";
  elements.recipientAccount.value = "";
  elements.recipientAccountName.value = "";
  document.querySelector('input[name="recipientMode"][value="shared"]').checked = true;
  elements.resultSearch.value = "";
  elements.resultFilter.value = "all";
  elements.emailFrom.value = "";
  elements.recipientSection.hidden = true;
  elements.mappingSection.hidden = true;
  elements.resultsSection.hidden = true;
  elements.emailSection.hidden = true;
  initializeTemplate();
  renderRecipientMode();
  renderExportAvailability();
  renderGmailSheetFlow();
  setStatus("Chưa có file", "neutral");
  setStep(1);
  window.scrollTo({ top: 0, behavior: "smooth" });
  showToast("Đã xoá dữ liệu phiên khỏi bộ nhớ");
}

async function renderCoffeeQR(amount = COFFEE_TRANSFER.defaultAmount) {
  const { amount: selectedAmount, payload } = generateCoffeeVietQR(amount);
  elements.coffeeQR.src = await QRCode.toDataURL(payload, {
    width: 500,
    margin: 2,
    errorCorrectionLevel: "M",
    type: "image/png",
  });
  const amountLabel = selectedAmount > 0 ? formatVND(selectedAmount) : "Người gửi tự nhập số tiền";
  elements.coffeeAmountLabel.textContent = amountLabel;
  elements.coffeeQR.alt = selectedAmount > 0
    ? `VietQR ${amountLabel} gửi ${COFFEE_TRANSFER.accountName}`
    : `VietQR tùy tâm gửi ${COFFEE_TRANSFER.accountName}`;
  document.querySelectorAll("[data-coffee-amount]").forEach((button) => {
    button.setAttribute("aria-pressed", String(Number(button.dataset.coffeeAmount) === selectedAmount));
  });
}

async function openCoffeeDialog() {
  try {
    await renderCoffeeQR(COFFEE_TRANSFER.defaultAmount);
    elements.coffeeDialog.showModal();
  } catch (error) {
    showToast(error.message || "Không tạo được VietQR cà phê", "error");
  }
}

elements.spreadsheetFile.addEventListener("change", () => handleSpreadsheet(elements.spreadsheetFile.files[0]));
["dragenter", "dragover"].forEach((eventName) => elements.dropZone.addEventListener(eventName, (event) => {
  event.preventDefault(); elements.dropZone.classList.add("dragging");
}));
["dragleave", "drop"].forEach((eventName) => elements.dropZone.addEventListener(eventName, (event) => {
  event.preventDefault(); elements.dropZone.classList.remove("dragging");
}));
elements.dropZone.addEventListener("drop", (event) => handleSpreadsheet(event.dataTransfer.files[0]));
document.querySelectorAll('input[name="recipientMode"]').forEach((input) => input.addEventListener("change", invalidateRecipientSetup));
[elements.recipientBank, elements.recipientAccount, elements.recipientAccountName].forEach((input) => {
  input.addEventListener("input", () => {
    if (state.recipientSetupConfirmed) invalidateRecipientSetup();
  });
});
elements.confirmRecipientSetup.addEventListener("click", confirmRecipientSetup);
elements.confirmRecipientExport.addEventListener("click", confirmRecipientExportAccounts);
elements.mappingRows.addEventListener("change", handleMappingChange);
elements.applyMapping.addEventListener("click", generateItems);
elements.resultSearch.addEventListener("input", renderResults);
elements.resultFilter.addEventListener("change", renderResults);
elements.resultRows.addEventListener("click", (event) => {
  const row = event.target.closest("tr[data-id]");
  if (!row) return;
  state.selectedId = row.dataset.id;
  renderResults();
  renderSelected();
  setStep(5);
});
elements.exportQR.addEventListener("click", exportQRZip);
elements.templateName.addEventListener("input", renderEmailPreview);
elements.emailSubject.addEventListener("input", renderEmailPreview);
elements.emailEditor.addEventListener("input", renderEmailPreview);
elements.emailEditor.addEventListener("keyup", saveEditorRange);
elements.emailEditor.addEventListener("mouseup", saveEditorRange);
elements.mergeFields.addEventListener("click", (event) => {
  const button = event.target.closest("button[data-token]");
  if (button) insertToken(button.dataset.token);
});
document.querySelectorAll("[data-command]").forEach((button) => button.addEventListener("click", () => {
  restoreEditorRange();
  document.execCommand(button.dataset.command, false);
  saveEditorRange();
  renderEmailPreview();
}));
elements.insertLink.addEventListener("click", () => {
  const url = elements.linkURL.value.trim();
  if (!/^https?:\/\//i.test(url)) return showToast("Link cần bắt đầu bằng http:// hoặc https://", "error");
  restoreEditorRange();
  document.execCommand("createLink", false, url);
  elements.linkURL.value = "";
  saveEditorRange();
  renderEmailPreview();
});
elements.copySubject.addEventListener("click", () => {
  if (!requireRecipientExportConfirmation()) return;
  const item = selectedItem();
  if (!item) return showToast("Chưa chọn recipient", "error");
  copyText(renderEmailTemplate(currentTemplate(), item).subject, "Đã copy tiêu đề");
});
elements.copyEmail.addEventListener("click", copyRichEmail);
elements.copyQR.addEventListener("click", copySelectedQR);
elements.downloadEml.addEventListener("click", downloadSelectedEml);
elements.exportEmailBundle.addEventListener("click", exportEmailZip);
elements.exportGmailData.addEventListener("click", exportGmailSheetDataFile);
elements.openGmailSheet.addEventListener("click", openPersonalGmailSheet);
elements.createGmailSheet.addEventListener("click", startGmailSetup);
elements.connectGmailSheet.addEventListener("click", openGmailConnectDialog);
elements.downloadGmailUpgrade.addEventListener("click", downloadGmailUpgradeScript);
elements.downloadGmailUpgradeConnected.addEventListener("click", downloadGmailUpgradeScript);
elements.toggleGmailUpgradeSteps.addEventListener("click", toggleGmailUpgradeInstructions);
elements.changeGmailSheet.addEventListener("click", openGmailConnectDialog);
elements.disconnectGmailSheet.addEventListener("click", disconnectGmailSheet);
elements.saveGmailSheet.addEventListener("click", saveGmailSheetConnection);
elements.gmailConnectForm.addEventListener("submit", (event) => {
  event.preventDefault();
  saveGmailSheetConnection();
});
elements.exportGmailFree.addEventListener("click", () => exportGmailFreeZip(elements.exportGmailFree));
elements.exportTemplate.addEventListener("click", exportTemplateFile);
elements.templateFile.addEventListener("change", () => importTemplateFile(elements.templateFile.files[0]));
elements.coffeeButton.addEventListener("click", openCoffeeDialog);
document.querySelectorAll("[data-coffee-amount]").forEach((button) => button.addEventListener("click", async () => {
  try {
    await renderCoffeeQR(button.dataset.coffeeAmount);
  } catch (error) {
    showToast(error.message || "Không cập nhật được VietQR cà phê", "error");
  }
}));
elements.resetButton.addEventListener("click", () => elements.resetDialog.showModal());
elements.resetDialog.addEventListener("close", () => { if (elements.resetDialog.returnValue === "confirm") resetSession(); });

initializeBanks();
initializeTemplate();
renderRecipientMode();
renderGmailSheetFlow();
renderExportAvailability();
renderEmailPreview();
