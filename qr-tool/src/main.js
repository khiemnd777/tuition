import DOMPurify from "dompurify";
import QRCode from "qrcode";

import { BANKS } from "./banks.js";
import {
  DEFAULT_EMAIL_TEMPLATE,
  MERGE_FIELDS,
  buildEml,
  normalizeImportedTemplate,
  renderEmailTemplate,
} from "./email.js";
import { createEmailBundle, createGmailMergeBundle, createQRBundle, safeFilename } from "./exporter.js";
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
  "resetButton", "resetDialog", "spreadsheetFile", "dropZone", "fileSummary", "mappingSection",
  "mappingSummary", "defaultBank", "defaultAccount", "mappingRows", "mappingHint", "applyMapping",
  "resultsSection", "totalCount", "validCount", "errorCount", "resultSearch", "resultFilter", "resultRows",
  "qrEmpty", "qrDetail", "previewStatus", "previewName", "previewMeta", "previewQR", "previewBank",
  "previewAccount", "previewAmount", "previewNote", "previewErrors", "previewPayload", "exportQR",
  "emailSection", "templateName", "emailFrom", "templateFile", "exportTemplate", "emailSubject",
  "mergeFields", "emailEditor", "linkURL", "insertLink", "emailPreviewRecipient", "emailPreviewFrame",
  "copySubject", "copyEmail", "copyQR", "downloadEml", "exportGmail", "exportEmailBundle", "toast",
].map((id) => [id, document.getElementById(id)]));

const state = {
  table: null,
  mapping: {},
  items: [],
  selectedId: "",
  template: structuredClone(DEFAULT_EMAIL_TEMPLATE),
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
  elements.defaultBank.insertAdjacentHTML("beforeend", BANKS
    .slice()
    .sort((a, b) => a.shortName.localeCompare(b.shortName, "vi"))
    .map((bank) => `<option value="${bank.bin}">${escapeHTML(bank.shortName)} · ${bank.bin}</option>`)
    .join(""));
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
  elements.mappingSummary.textContent = `${state.table.headers.length} cột · ${state.table.records.length} dòng`;
  elements.mappingRows.innerHTML = state.table.headers.map((header) => {
    const selected = state.mapping[header] || "";
    const options = [
      `<option value="">Bỏ qua</option>`,
      ...IMPORT_FIELD_GROUPS.map((group) => `<optgroup label="${escapeHTML(group.label)}">${group.fields.map((field) =>
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
  elements.mappingHint.textContent = errors.length
    ? errors.join(" · ")
    : `${mapped}/${state.table?.headers.length || 0} cột đã map${fees ? ` · ${fees} khoản phí custom` : ""}.`;
  elements.mappingHint.style.color = errors.length ? "var(--danger)" : "";
}

async function handleSpreadsheet(file) {
  if (!file) return;
  setStatus("Đang đọc file…", "busy");
  try {
    state.table = await readSpreadsheet(file);
    state.mapping = suggestMapping(state.table.headers);
    state.items = [];
    state.selectedId = "";
    elements.mappingSection.hidden = false;
    elements.resultsSection.hidden = true;
    elements.emailSection.hidden = true;
    renderMapping();
    setStatus(`${file.name} · ${state.table.records.length} dòng`, "success");
    setStep(2);
    elements.mappingSection.scrollIntoView({ behavior: "smooth", block: "start" });
  } catch (error) {
    setStatus("Không đọc được file", "danger");
    showToast(error.message || "Không đọc được bảng dữ liệu", "error");
  }
}

async function generateItems() {
  const mapping = collectMapping();
  const errors = validateMapping(mapping);
  if (errors.length > 0) {
    showToast(errors.join(" · "), "error");
    return;
  }
  setBusy(elements.applyMapping, true, "Đang sinh QR…");
  try {
    const rows = buildPaymentRows(state.table, mapping, {
      bankBin: elements.defaultBank.value,
      bankAccount: elements.defaultAccount.value,
    });
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
    elements.resultsSection.hidden = false;
    elements.emailSection.hidden = false;
    renderResults();
    renderSelected();
    setStep(3);
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
  elements.exportQR.disabled = valid === 0;
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
  elements.previewAmount.textContent = formatVND(item.amount);
  elements.previewNote.textContent = item.note || "-";
  elements.previewErrors.hidden = ready;
  elements.previewErrors.textContent = item.errors.join(" · ");
  elements.previewPayload.textContent = item.vietqr || "Chưa có payload";
  elements.emailPreviewRecipient.textContent = item.email || "Thiếu email";
  renderEmailPreview();
}

function renderEmailPreview() {
  const item = selectedItem();
  if (!item) {
    elements.emailPreviewFrame.srcdoc = "<p style='font-family:Arial;padding:24px;color:#667085'>Chọn một dòng để preview email.</p>";
    return;
  }
  const template = currentTemplate();
  state.template = template;
  const rendered = renderEmailTemplate(template, item, { qrSrc: item.qrData || "" });
  elements.emailPreviewFrame.srcdoc = DOMPurify.sanitize(rendered.html, {
    WHOLE_DOCUMENT: true,
    ADD_TAGS: ["html", "head", "body", "img"],
    ADD_ATTR: ["src", "alt", "width", "height"],
    ALLOW_DATA_ATTR: false,
  });
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
  setBusy(elements.exportEmailBundle, true, "Đang đóng gói…");
  try {
    const blob = await createEmailBundle(state.items, currentTemplate(), { from: elements.emailFrom.value.trim() });
    downloadBlob(blob, `dekisugi-email-bulk-${new Date().toISOString().slice(0, 10)}.zip`);
    showToast("Đã tạo email bulk ZIP; không có email nào được gửi");
  } catch (error) {
    showToast(error.message || "Không tạo được email bundle", "error");
  } finally { setBusy(elements.exportEmailBundle, false); }
}

async function exportGmailZip() {
  setBusy(elements.exportGmail, true, "Đang tạo workbook…");
  try {
    const blob = await createGmailMergeBundle(state.items, currentTemplate());
    downloadBlob(blob, `dekisugi-gmail-mail-merge-${new Date().toISOString().slice(0, 10)}.zip`);
    showToast("Đã tạo Gmail Mail Merge workbook và hướng dẫn");
  } catch (error) {
    showToast(error.message || "Không tạo được Gmail Mail Merge", "error");
  } finally { setBusy(elements.exportGmail, false); }
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
  elements.spreadsheetFile.value = "";
  elements.defaultBank.value = "";
  elements.defaultAccount.value = "";
  elements.resultSearch.value = "";
  elements.resultFilter.value = "all";
  elements.emailFrom.value = "";
  elements.mappingSection.hidden = true;
  elements.resultsSection.hidden = true;
  elements.emailSection.hidden = true;
  initializeTemplate();
  setStatus("Chưa có file", "neutral");
  setStep(1);
  window.scrollTo({ top: 0, behavior: "smooth" });
  showToast("Đã xoá dữ liệu phiên khỏi bộ nhớ");
}

elements.spreadsheetFile.addEventListener("change", () => handleSpreadsheet(elements.spreadsheetFile.files[0]));
["dragenter", "dragover"].forEach((eventName) => elements.dropZone.addEventListener(eventName, (event) => {
  event.preventDefault(); elements.dropZone.classList.add("dragging");
}));
["dragleave", "drop"].forEach((eventName) => elements.dropZone.addEventListener(eventName, (event) => {
  event.preventDefault(); elements.dropZone.classList.remove("dragging");
}));
elements.dropZone.addEventListener("drop", (event) => handleSpreadsheet(event.dataTransfer.files[0]));
elements.mappingRows.addEventListener("change", updateMappingHint);
elements.applyMapping.addEventListener("click", generateItems);
elements.resultSearch.addEventListener("input", renderResults);
elements.resultFilter.addEventListener("change", renderResults);
elements.resultRows.addEventListener("click", (event) => {
  const row = event.target.closest("tr[data-id]");
  if (!row) return;
  state.selectedId = row.dataset.id;
  renderResults();
  renderSelected();
  setStep(4);
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
  const item = selectedItem();
  if (!item) return showToast("Chưa chọn recipient", "error");
  copyText(renderEmailTemplate(currentTemplate(), item).subject, "Đã copy tiêu đề");
});
elements.copyEmail.addEventListener("click", copyRichEmail);
elements.copyQR.addEventListener("click", copySelectedQR);
elements.downloadEml.addEventListener("click", downloadSelectedEml);
elements.exportEmailBundle.addEventListener("click", exportEmailZip);
elements.exportGmail.addEventListener("click", exportGmailZip);
elements.exportTemplate.addEventListener("click", exportTemplateFile);
elements.templateFile.addEventListener("change", () => importTemplateFile(elements.templateFile.files[0]));
elements.resetButton.addEventListener("click", () => elements.resetDialog.showModal());
elements.resetDialog.addEventListener("close", () => { if (elements.resetDialog.returnValue === "confirm") resetSession(); });

initializeBanks();
initializeTemplate();
renderEmailPreview();
