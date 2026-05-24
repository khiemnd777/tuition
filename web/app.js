const rowsEl = document.querySelector("#rows");
const rowCountEl = document.querySelector("#rowCount");
const resultCountEl = document.querySelector("#resultCount");
const resultsEl = document.querySelector("#results");
const previewEl = document.querySelector("#preview");
const statusEl = document.querySelector("#status");
const dataPanelEl = document.querySelector(".data-panel");
const generateBtn = document.querySelector("#generate");
const addRowBtn = document.querySelector("#addRow");
const loadSampleBtn = document.querySelector("#loadSample");
const csvFileEl = document.querySelector("#csvFile");
const toggleFeeColumnBtn = document.querySelector("#toggleFeeColumn");
const feeTemplateEl = document.querySelector("#feeTemplate");
const applyFeeTemplateBtn = document.querySelector("#applyFeeTemplate");
const emailProviderEl = document.querySelector("#emailProvider");
const gmailAddressEl = document.querySelector("#gmailAddress");
const gmailAppPasswordEl = document.querySelector("#gmailAppPassword");
const gmailAddressFieldEl = document.querySelector("#gmailAddressField");
const gmailAppPasswordFieldEl = document.querySelector("#gmailAppPasswordField");
const resendApiKeyFieldEl = document.querySelector("#resendApiKeyField");
const resendApiKeyEl = document.querySelector("#resendApiKey");
const emailFromEl = document.querySelector("#emailFrom");
const emailReplyToEl = document.querySelector("#emailReplyTo");
const emailSubjectEl = document.querySelector("#emailSubject");
const schoolNameEl = document.querySelector("#schoolName");
const schoolNameEnEl = document.querySelector("#schoolNameEn");
const paymentPeriodEl = document.querySelector("#paymentPeriod");
const publicBaseUrlEl = document.querySelector("#publicBaseUrl");
const emailTemplateEl = document.querySelector("#emailTemplate");
const saveEmailConfigBtn = document.querySelector("#saveEmailConfig");
const previewEmailBtn = document.querySelector("#previewEmail");
const dryRunEmailBtn = document.querySelector("#dryRunEmail");
const sendEmailBtn = document.querySelector("#sendEmail");
const emailConfigStatusEl = document.querySelector("#emailConfigStatus");
const emailSendStatusEl = document.querySelector("#emailSendStatus");
const emailPreviewSubjectEl = document.querySelector("#emailPreviewSubject");
const emailPreviewFrameEl = document.querySelector("#emailPreviewFrame");
const cronEnabledEl = document.querySelector("#cronEnabled");
const cronSendTimeEl = document.querySelector("#cronSendTime");
const cronDailyLimitEl = document.querySelector("#cronDailyLimit");
const cronQueueSummaryEl = document.querySelector("#cronQueueSummary");
const saveCronBtn = document.querySelector("#saveCron");
const disableCronBtn = document.querySelector("#disableCron");
const runCronNowBtn = document.querySelector("#runCronNow");
const cronStatusEl = document.querySelector("#cronStatus");
const masterDataStatusEl = document.querySelector("#masterDataStatus");
const masterSchoolYearFilterEl = document.querySelector("#masterSchoolYearFilter");
const masterGradeFilterEl = document.querySelector("#masterGradeFilter");
const masterClassFilterEl = document.querySelector("#masterClassFilter");
const masterSearchEl = document.querySelector("#masterSearch");
const refreshMasterDataBtn = document.querySelector("#refreshMasterData");
const masterCsvFileEl = document.querySelector("#masterCsvFile");
const checkMasterImportBtn = document.querySelector("#checkMasterImport");
const applyMasterImportBtn = document.querySelector("#applyMasterImport");
const masterImportCountEl = document.querySelector("#masterImportCount");
const masterImportSummaryEl = document.querySelector("#masterImportSummary");
const masterStudentsEl = document.querySelector("#masterStudents");
const masterStudentCountEl = document.querySelector("#masterStudentCount");
const masterConflictPanelEl = document.querySelector("#masterConflictPanel");
const masterConflictCountEl = document.querySelector("#masterConflictCount");
const masterConflictListEl = document.querySelector("#masterConflictList");
const feeScheduleLoadStatusEl = document.querySelector("#feeScheduleLoadStatus");
const refreshFeeSchedulesBtn = document.querySelector("#refreshFeeSchedules");
const previewFeeScheduleBtn = document.querySelector("#previewFeeSchedule");
const saveFeeScheduleBtn = document.querySelector("#saveFeeSchedule");
const feeScheduleYearEl = document.querySelector("#feeScheduleYear");
const feeScheduleGradeEl = document.querySelector("#feeScheduleGrade");
const feeScheduleClassEl = document.querySelector("#feeScheduleClass");
const feeSchedulePeriodEl = document.querySelector("#feeSchedulePeriod");
const feeScheduleMonthEl = document.querySelector("#feeScheduleMonth");
const feeScheduleStatusEl = document.querySelector("#feeScheduleStatus");
const feeScheduleNameEl = document.querySelector("#feeScheduleName");
const feeScheduleNotesEl = document.querySelector("#feeScheduleNotes");
const feeScheduleItemsEl = document.querySelector("#feeScheduleItems");
const feeScheduleItemTotalEl = document.querySelector("#feeScheduleItemTotal");
const feeAdjustmentsCsvEl = document.querySelector("#feeAdjustmentsCsv");
const feeAdjustmentCountEl = document.querySelector("#feeAdjustmentCount");
const feeScheduleSummaryEl = document.querySelector("#feeScheduleSummary");
const feeSchedulePreviewCountEl = document.querySelector("#feeSchedulePreviewCount");
const feeSchedulePreviewRowsEl = document.querySelector("#feeSchedulePreviewRows");
const feeSchedulesEl = document.querySelector("#feeSchedules");
const feeScheduleCountEl = document.querySelector("#feeScheduleCount");
const invoiceStatusEl = document.querySelector("#invoiceStatus");
const refreshInvoicesBtn = document.querySelector("#refreshInvoices");
const previewInvoicesBtn = document.querySelector("#previewInvoices");
const generateInvoicesBtn = document.querySelector("#generateInvoices");
const invoiceScheduleEl = document.querySelector("#invoiceSchedule");
const invoiceBankBinEl = document.querySelector("#invoiceBankBin");
const invoiceBankAccountEl = document.querySelector("#invoiceBankAccount");
const invoiceIssueDateEl = document.querySelector("#invoiceIssueDate");
const invoiceDueDateEl = document.querySelector("#invoiceDueDate");
const invoiceRegenerateEl = document.querySelector("#invoiceRegenerate");
const invoicePreviewCountEl = document.querySelector("#invoicePreviewCount");
const invoicePreviewSummaryEl = document.querySelector("#invoicePreviewSummary");
const invoicePreviewRowsEl = document.querySelector("#invoicePreviewRows");
const invoiceRowsEl = document.querySelector("#invoiceRows");
const invoiceCountEl = document.querySelector("#invoiceCount");
const invoicePaymentStatusEl = document.querySelector("#invoicePaymentStatus");
const invoicePaymentPreviewEl = document.querySelector("#invoicePaymentPreview");
const tabButtons = [...document.querySelectorAll(".tab-button")];
const tabPanels = [...document.querySelectorAll(".tab-panel")];

let banks = [];
let currentItems = [];
let selectedId = "";
let feeColumnCollapsed = false;
let savedEmailConfig = {};
let masterDataOptions = { schoolYears: [], classes: [] };
let masterDataLoaded = false;
let feeScheduleOptions = { feeTypes: [], schoolYears: [], classes: [] };
let feeSchedulesLoaded = false;
let invoiceOptions = { schedules: [], schoolYears: [], classes: [] };
let invoicesLoaded = false;

const defaultPaymentItems = [
  { label: "Tiền học phí Tháng 04", labelEn: "Tuition fees for April", amount: 3950000 },
  { label: "Phí xe đưa rước Tháng 04", labelEn: "Shuttle fees for April", amount: 3030000 },
  { label: "Tiền học phí Tháng 05", labelEn: "Tuition fees for May", amount: 3950000 },
  { label: "Bảo hiểm y tế", labelEn: "Health Insurance", amount: 0 },
  { label: "Đồng phục", labelEn: "Uniform fee", amount: 0 },
  { label: "Sách CTQT", labelEn: "International material", amount: 0 },
  { label: "Các khoản phí tháng trước", labelEn: "Previous month's fees", amount: 0 },
];

const defaultFeeTypes = [
  { code: "tuition", labelVi: "Học phí", labelEn: "Tuition", displayOrder: 10 },
  { code: "lunch", labelVi: "Tiền ăn", labelEn: "Lunch", displayOrder: 20 },
  { code: "shuttle", labelVi: "Phí xe đưa rước", labelEn: "Shuttle", displayOrder: 30 },
  { code: "uniform", labelVi: "Đồng phục", labelEn: "Uniform", displayOrder: 40 },
  { code: "insurance", labelVi: "Bảo hiểm", labelEn: "Insurance", displayOrder: 50 },
  { code: "materials", labelVi: "Học liệu", labelEn: "Learning materials", displayOrder: 60 },
  { code: "previous_fees", labelVi: "Phí kỳ trước", labelEn: "Previous fees", displayOrder: 70 },
  { code: "custom", labelVi: "Khoản phí khác", labelEn: "Custom fee", displayOrder: 100 },
];

const sampleRows = [
  {
    id: "real-vib",
    studentName: "Nguyễn Duy Khiêm",
    parentName: "Nguyễn Duy Khiêm",
    className: "3.02",
    bankBin: "970441",
    bankAccount: "625704060370690",
    email: "",
    amount: 10930000,
    paymentItems: clonePaymentItems(defaultPaymentItems),
    billNumber: "TESTVIB01",
    note: "Hoc phi thang 04+05",
  },
  {
    id: "real-vpbank",
    studentName: "Nguyễn Duy Khiêm",
    parentName: "Nguyễn Duy Khiêm",
    className: "3.02",
    bankBin: "970432",
    bankAccount: "0974322365",
    email: "",
    amount: 10930000,
    paymentItems: clonePaymentItems(defaultPaymentItems),
    billNumber: "TESTVPB01",
    note: "Hoc phi thang 04+05",
  },
];

init();

function muiIcon(name) {
  return `<span class="mui-icon" aria-hidden="true">${escapeHtml(name)}</span>`;
}

async function init() {
  status("Đang tải", "busy");
  await loadBanks();
  await loadEmailConfig();
  await loadEmailCron();
  renderMasterFilters();
  renderMasterStudents([]);
  renderFeeScheduleControls();
  renderFeeScheduleItems(defaultFeeTypes);
  renderFeeSchedulePreview(null);
  renderFeeSchedules([]);
  renderInvoiceControls();
  renderInvoicePreview(null);
  renderInvoices([]);
  renderFeeTemplate(defaultPaymentItems);
  renderRows(sampleRows);
  await loadMasterData();
  await generate();
  await previewEmail();
}

async function loadBanks() {
  const res = await fetch("/api/v1/banks");
  if (!res.ok) throw new Error("Không tải được danh sách ngân hàng");
  const data = await res.json();
  banks = data.banks || [];
}

function renderRows(rows) {
  rowsEl.innerHTML = "";
  rows.forEach((row) => rowsEl.appendChild(rowTemplate(row)));
  updateRowCount();
}

async function activateTab(targetId) {
  tabButtons.forEach((button) => {
    const isActive = button.dataset.tabTarget === targetId;
    button.classList.toggle("active", isActive);
    button.setAttribute("aria-selected", String(isActive));
  });
  tabPanels.forEach((panel) => {
    const isActive = panel.id === targetId;
    panel.hidden = !isActive;
    panel.classList.toggle("active", isActive);
  });
  if (targetId === "masterDataTab") {
    await loadMasterData();
  }
  if (targetId === "feeTemplateTab") {
    await loadFeeSchedules();
  }
  if (targetId === "invoiceTab") {
    await loadInvoices();
  }
  if (targetId === "emailTab") {
    await previewEmail();
  }
}

function rowTemplate(row = {}) {
  const tr = document.createElement("tr");
  tr.dataset.id = row.id || crypto.randomUUID();
  tr.innerHTML = `
    <td><input data-field="studentName" value="${escapeAttr(row.studentName || "")}" /></td>
    <td><input data-field="parentName" value="${escapeAttr(row.parentName || "")}" /></td>
    <td><input data-field="className" value="${escapeAttr(row.className || "")}" /></td>
    <td><select data-field="bankBin">${bankOptions(row.bankBin || "970415")}</select></td>
    <td><input data-field="bankAccount" value="${escapeAttr(row.bankAccount || "")}" inputmode="numeric" /></td>
    <td><input data-field="email" value="${escapeAttr(row.email || "")}" type="email" /></td>
    <td>${paymentItemsTemplate(row.paymentItems, row.amount)}</td>
    <td><input data-field="billNumber" value="${escapeAttr(row.billNumber || "")}" maxlength="25" /></td>
    <td><input data-field="note" value="${escapeAttr(row.note || "")}" maxlength="60" /></td>
    <td><button class="delete-row danger" type="button" title="Xóa dòng">${muiIcon("delete")}<span>Xóa</span></button></td>
  `;
  tr.querySelector(".delete-row").addEventListener("click", () => {
    tr.remove();
    updateRowCount();
  });
  initializeFeeEditor(tr);
  return tr;
}

function renderFeeTemplate(items) {
  feeTemplateEl.innerHTML = paymentItemsTemplate(items);
  initializeFeeEditor(feeTemplateEl);
}

function paymentItemsTemplate(items = [], fallbackAmount = 0) {
  const normalized = items.length
    ? items
    : [{ label: "Tổng phí cần thanh toán", labelEn: "Total fees due", amount: Number(fallbackAmount || 0) }];
  return `
    <div class="fee-editor">
      <div class="fee-items">
        ${normalized.map(paymentItemTemplate).join("")}
      </div>
      <div class="fee-toolbar">
        <button class="add-fee" type="button">${muiIcon("add_circle")}<span>Thêm khoản</span></button>
        <span class="fee-summary">
          <span class="fee-count">${normalized.length} khoản</span>
          <span class="fee-total">0 đ</span>
        </span>
      </div>
    </div>
  `;
}

function paymentItemTemplate(item = {}) {
  return `
    <div class="fee-item">
      <input class="fee-label" data-fee-field="label" value="${escapeAttr(item.label || "")}" placeholder="Diễn giải" />
      <input class="fee-label-en" data-fee-field="labelEn" value="${escapeAttr(item.labelEn || "")}" placeholder="Explanation" />
      <input class="fee-amount" data-fee-field="amount" value="${Number(item.amount || 0)}" type="number" min="0" step="1000" />
      <button class="remove-fee danger" type="button" title="Xóa khoản phí">${muiIcon("remove_circle")}<span>Xóa</span></button>
    </div>
  `;
}

function initializeFeeEditor(root) {
  root.querySelector(".add-fee").addEventListener("click", () => addPaymentItem(root));
  root.querySelectorAll(".remove-fee").forEach((button) => {
    button.addEventListener("click", () => {
      button.closest(".fee-item").remove();
      updateFeeTotal(root);
    });
  });
  root.querySelectorAll(".fee-amount").forEach((input) => {
    input.addEventListener("input", () => updateFeeTotal(root));
  });
  updateFeeTotal(root);
}

function addPaymentItem(root) {
  const list = root.querySelector(".fee-items");
  const wrapper = document.createElement("div");
  wrapper.innerHTML = paymentItemTemplate({
    label: "Khoản phí",
    labelEn: "Fee",
    amount: 0,
  }).trim();
  const item = wrapper.firstElementChild;
  list.appendChild(item);
  item.querySelector(".remove-fee").addEventListener("click", () => {
    item.remove();
    updateFeeTotal(root);
  });
  item.querySelector(".fee-amount").addEventListener("input", () => updateFeeTotal(root));
  updateFeeTotal(root);
}

function updateFeeTotal(root) {
  const items = collectPaymentItems(root);
  const total = items.reduce((sum, item) => sum + Number(item.amount || 0), 0);
  const countEl = root.querySelector(".fee-count");
  if (countEl) {
    countEl.textContent = `${items.length} khoản`;
  }
  root.querySelector(".fee-total").textContent = formatMoney(total);
}

function setFeeColumnCollapsed(collapsed) {
  feeColumnCollapsed = collapsed;
  dataPanelEl.classList.toggle("fees-collapsed", feeColumnCollapsed);
  toggleFeeColumnBtn.innerHTML = `${muiIcon(feeColumnCollapsed ? "view_week" : "view_column")}<span>${feeColumnCollapsed ? "Mở rộng khoản phí" : "Thu gọn khoản phí"}</span>`;
  toggleFeeColumnBtn.setAttribute("aria-expanded", String(!feeColumnCollapsed));
}

function currentFeeTemplateItems() {
  const items = collectPaymentItems(feeTemplateEl);
  if (items.length) {
    return clonePaymentItems(items);
  }
  return [{ label: "Tổng phí cần thanh toán", labelEn: "Total fees due", amount: 0 }];
}

function clonePaymentItems(items = []) {
  return items.map((item) => ({
    label: item.label || "",
    labelEn: item.labelEn || "",
    amount: Number(item.amount || 0),
  }));
}

function bankOptions(selected) {
  return banks
    .map((bank) => {
      const label = `${bank.shortName || bank.code} - ${bank.bin}`;
      const isSelected = bank.bin === selected ? "selected" : "";
      return `<option value="${bank.bin}" ${isSelected}>${escapeHtml(label)}</option>`;
    })
    .join("");
}

function collectRows() {
  return [...rowsEl.querySelectorAll("tr")].map((tr, index) => {
    const row = { id: tr.dataset.id || `row-${index + 1}` };
    tr.querySelectorAll("[data-field]").forEach((field) => {
      const key = field.dataset.field;
      row[key] = field.value.trim();
    });
    row.paymentItems = collectPaymentItems(tr);
    row.amount = row.paymentItems.reduce((sum, item) => sum + Number(item.amount || 0), 0);
    return row;
  });
}

function collectPaymentItems(tr) {
  return [...tr.querySelectorAll(".fee-item")]
    .map((item) => ({
      label: item.querySelector('[data-fee-field="label"]').value.trim(),
      labelEn: item.querySelector('[data-fee-field="labelEn"]').value.trim(),
      amount: Number(item.querySelector('[data-fee-field="amount"]').value || 0),
    }))
    .filter((item) => item.label || item.labelEn || item.amount);
}

async function generate() {
  status("Đang sinh QR", "busy");
  const rows = collectRows();
  const res = await fetch("/api/v1/vietqr/batch", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ rows }),
  });

  if (!res.ok) {
    status("Lỗi", "error");
    previewEl.className = "preview-empty";
    previewEl.textContent = await res.text();
    return;
  }

  const data = await res.json();
  currentItems = data.items || [];
  selectedId = currentItems.find((item) => item.vietqr)?.id || currentItems[0]?.id || "";
  renderResults();
  renderPreview();
  status("Sẵn sàng", "ready");
}

async function loadEmailConfig() {
  const res = await fetch("/api/v1/email/config");
  if (!res.ok) return;
  const config = await res.json();
  savedEmailConfig = config;
  emailProviderEl.value = config.provider || "gmail";
  gmailAddressEl.value = config.gmailAddress || "";
  gmailAppPasswordEl.placeholder = config.hasGmailAppPassword ? config.gmailAppPasswordMasked : "16-character app password";
  emailFromEl.value = config.from || "";
  emailReplyToEl.value = config.replyTo || "";
  emailSubjectEl.value = config.subject || "";
  schoolNameEl.value = config.schoolName || "";
  schoolNameEnEl.value = config.schoolNameEn || "";
  paymentPeriodEl.value = config.paymentPeriod || "";
  publicBaseUrlEl.value = config.publicBaseUrl || "";
  resendApiKeyEl.placeholder = config.hasApiKey ? config.apiKeyMasked : "re_...";
  updateEmailProviderFields();
  renderEmailConfigStatus(config);
}

async function loadEmailCron() {
  const res = await fetch("/api/v1/email/cron");
  if (!res.ok) return;
  renderCronStatus(await res.json());
}

async function loadMasterData(force = false) {
  if (!force && masterDataLoaded) {
    return;
  }
  const loaded = await loadMasterDataOptions();
  if (loaded) {
    await loadMasterStudents();
  }
}

async function loadMasterDataOptions() {
  setMasterStatus("Đang tải", "busy");
  const res = await fetch("/api/v1/master-data/options");
  const text = await res.text();
  if (!res.ok) {
    masterDataLoaded = false;
    masterDataOptions = { schoolYears: [], classes: [] };
    renderMasterFilters();
    renderMasterStudents([]);
    renderMasterImport(null);
    setMasterStatus(text || "Chưa cấu hình DB", "error");
    return false;
  }
  masterDataOptions = JSON.parse(text);
  masterDataLoaded = true;
  renderMasterFilters();
  setMasterStatus("Sẵn sàng", "ready");
  return true;
}

function renderMasterFilters() {
  const selectedYear = masterSchoolYearFilterEl.value;
  const selectedGrade = masterGradeFilterEl.value;
  const selectedClass = masterClassFilterEl.value;

  masterSchoolYearFilterEl.innerHTML = [
    `<option value="">Tất cả năm học</option>`,
    ...(masterDataOptions.schoolYears || []).map(
      (item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.code)}</option>`,
    ),
  ].join("");
  masterSchoolYearFilterEl.value = optionValueOrEmpty(masterSchoolYearFilterEl, selectedYear);

  const grades = [...new Set(
    (masterDataOptions.classes || [])
      .filter((item) => !masterSchoolYearFilterEl.value || item.schoolYearId === masterSchoolYearFilterEl.value)
      .map((item) => item.grade)
      .filter(Boolean),
  )].sort((a, b) => a.localeCompare(b, "vi", { numeric: true }));
  masterGradeFilterEl.innerHTML = [
    `<option value="">Tất cả khối</option>`,
    ...grades.map((grade) => `<option value="${escapeAttr(grade)}">${escapeHtml(grade)}</option>`),
  ].join("");
  masterGradeFilterEl.value = optionValueOrEmpty(masterGradeFilterEl, selectedGrade);

  renderMasterClassFilter(selectedClass);
}

function renderMasterClassFilter(selectedClass = masterClassFilterEl.value) {
  const classes = (masterDataOptions.classes || []).filter((item) => {
    if (masterSchoolYearFilterEl.value && item.schoolYearId !== masterSchoolYearFilterEl.value) return false;
    if (masterGradeFilterEl.value && item.grade !== masterGradeFilterEl.value) return false;
    return true;
  });
  masterClassFilterEl.innerHTML = [
    `<option value="">Tất cả lớp</option>`,
    ...classes.map(
      (item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.schoolYearCode)} · ${escapeHtml(item.name)}</option>`,
    ),
  ].join("");
  masterClassFilterEl.value = optionValueOrEmpty(masterClassFilterEl, selectedClass);
}

function optionValueOrEmpty(selectEl, value) {
  return [...selectEl.options].some((option) => option.value === value) ? value : "";
}

async function loadMasterStudents() {
  if (!masterDataLoaded) {
    return;
  }
  setMasterStatus("Đang tải", "busy");
  const params = new URLSearchParams();
  if (masterSchoolYearFilterEl.value) params.set("schoolYearId", masterSchoolYearFilterEl.value);
  if (masterGradeFilterEl.value) params.set("grade", masterGradeFilterEl.value);
  if (masterClassFilterEl.value) params.set("classId", masterClassFilterEl.value);
  if (masterSearchEl.value.trim()) params.set("q", masterSearchEl.value.trim());

  const res = await fetch(`/api/v1/master-data/students?${params.toString()}`);
  const text = await res.text();
  if (!res.ok) {
    renderMasterStudents([]);
    setMasterStatus(text || "Không tải được học sinh", "error");
    return;
  }
  const data = JSON.parse(text);
  renderMasterStudents(data.students || []);
  setMasterStatus("Sẵn sàng", "ready");
}

function renderMasterStudents(students) {
  masterStudentCountEl.textContent = `${students.length} học sinh`;
  masterStudentsEl.innerHTML = students
    .map((student) => {
      const primary = (student.parents || []).find((parent) => parent.isPrimary) || (student.parents || [])[0] || {};
      const parentNames = (student.parents || []).map((parent) => parent.parentName).filter(Boolean).join(", ");
      const billingEmails = (student.parents || [])
        .filter((parent) => parent.receivesBillingEmail && parent.isActive && parent.emailActive && parent.email)
        .map((parent) => parent.email)
        .join(", ");
      return `
        <tr>
          <td><strong>${escapeHtml(student.studentCode || "")}</strong></td>
          <td>${escapeHtml(student.studentName || "")}</td>
          <td>${escapeHtml(student.schoolYearCode || "")}</td>
          <td>${escapeHtml(student.grade || "")}</td>
          <td>${escapeHtml(student.className || "")}</td>
          <td>${escapeHtml(parentNames || primary.parentName || "-")}</td>
          <td>${escapeHtml(billingEmails || "-")}</td>
          <td><span class="tag">${escapeHtml(student.status || "active")}</span></td>
        </tr>
      `;
    })
    .join("");
  if (!students.length) {
    masterStudentsEl.innerHTML = `<tr><td colspan="8" class="empty-cell">Chưa có dữ liệu học sinh</td></tr>`;
  }
}

async function submitMasterImport(apply) {
  const file = masterCsvFileEl.files[0];
  if (!file) {
    setMasterStatus("Chưa chọn CSV", "error");
    return;
  }
  if (apply && !window.confirm("Áp dụng import sẽ ghi dữ liệu học sinh, phụ huynh và lớp vào database. Tiếp tục?")) {
    return;
  }

  setMasterStatus(apply ? "Đang áp dụng" : "Đang kiểm tra", "busy");
  const form = new FormData();
  form.append("file", file);
  const res = await fetch(`/api/v1/master-data/import/csv?apply=${apply ? "true" : "false"}`, {
    method: "POST",
    body: form,
  });
  const text = await res.text();
  let data = null;
  try {
    data = JSON.parse(text);
  } catch {
    data = null;
  }
  if (!data) {
    setMasterStatus(text || "Import failed", "error");
    return;
  }
  renderMasterImport(data);
  if (!res.ok) {
    setMasterStatus(data.issues?.length ? "Có conflict" : "Import failed", "error");
    return;
  }
  if (data.options) {
    masterDataOptions = data.options;
    renderMasterFilters();
  } else {
    await loadMasterDataOptions();
  }
  await loadMasterStudents();
  setMasterStatus(data.applied ? "Đã áp dụng" : "Đã kiểm tra", data.issues?.length ? "error" : "ready");
}

function renderMasterImport(data) {
  if (!data) {
    masterImportCountEl.textContent = "0 dòng";
    masterImportSummaryEl.textContent = "Chưa có file import";
    renderMasterConflicts([]);
    return;
  }
  const summary = data.summary || {};
  masterImportCountEl.textContent = `${summary.totalRows || 0} dòng`;
  masterImportSummaryEl.className = `master-import-summary${summary.conflicts ? " error" : ""}`;
  masterImportSummaryEl.textContent = [
    data.applied ? "APPLIED" : "PREVIEW",
    `ready ${summary.ready || 0}`,
    `create/link ${summary.created || 0}`,
    `unchanged ${summary.unchanged || 0}`,
    `conflicts ${summary.conflicts || 0}`,
  ].join(" · ");
  renderMasterConflicts(data.issues || []);
}

function renderMasterConflicts(issues) {
  masterConflictPanelEl.hidden = !issues.length;
  masterConflictCountEl.textContent = `${issues.length} lỗi`;
  masterConflictListEl.innerHTML = issues
    .map(
      (issue) => `
        <div class="master-conflict-item">
          <strong>Dòng ${escapeHtml(issue.rowNumber || "-")} · ${escapeHtml(issue.studentCode || "-")}</strong>
          <span>${escapeHtml(issue.type || "")}</span>
          <p>${escapeHtml(issue.message || "")}</p>
          ${issue.existing || issue.incoming ? `<small>Existing: ${escapeHtml(issue.existing || "-")} · Incoming: ${escapeHtml(issue.incoming || "-")}</small>` : ""}
        </div>
      `,
    )
    .join("");
}

async function loadFeeSchedules(force = false) {
  if (!force && feeSchedulesLoaded) {
    return;
  }
  const loaded = await loadFeeScheduleOptions();
  if (loaded) {
    await loadFeeScheduleList();
  }
}

async function loadFeeScheduleOptions() {
  setFeeScheduleStatus("Đang tải", "busy");
  const res = await fetch("/api/v1/fee-schedules/options");
  const text = await res.text();
  if (!res.ok) {
    feeSchedulesLoaded = false;
    renderFeeSchedules([]);
    renderFeeSchedulePreview(null);
    setFeeScheduleStatus(text || "Chưa cấu hình DB", "error");
    return false;
  }
  const data = JSON.parse(text);
  feeScheduleOptions = {
    feeTypes: data.feeTypes || defaultFeeTypes,
    schoolYears: data.schoolYears || [],
    classes: data.classes || [],
  };
  feeSchedulesLoaded = true;
  renderFeeScheduleControls();
  renderFeeScheduleItems(feeScheduleOptions.feeTypes);
  setFeeScheduleStatus("Sẵn sàng", "ready");
  return true;
}

async function loadFeeScheduleList() {
  if (!feeSchedulesLoaded) {
    return;
  }
  const params = new URLSearchParams();
  if (feeScheduleYearEl.value) params.set("schoolYearId", feeScheduleYearEl.value);
  if (feeScheduleClassEl.value) params.set("classId", feeScheduleClassEl.value);
  if (!feeScheduleClassEl.value && feeScheduleGradeEl.value) params.set("grade", feeScheduleGradeEl.value);

  const res = await fetch(`/api/v1/fee-schedules?${params.toString()}`);
  const text = await res.text();
  if (!res.ok) {
    renderFeeSchedules([]);
    setFeeScheduleStatus(text || "Không tải được bảng phí", "error");
    return;
  }
  const data = JSON.parse(text);
  renderFeeSchedules(data.schedules || []);
  setFeeScheduleStatus("Sẵn sàng", "ready");
}

function renderFeeScheduleControls() {
  const selectedYear = feeScheduleYearEl.value;
  const selectedGrade = feeScheduleGradeEl.value;
  const selectedClass = feeScheduleClassEl.value;

  feeScheduleYearEl.innerHTML = [
    `<option value="">Chọn năm học</option>`,
    ...(feeScheduleOptions.schoolYears || []).map(
      (item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.code)}</option>`,
    ),
  ].join("");
  feeScheduleYearEl.value = optionValueOrEmpty(feeScheduleYearEl, selectedYear);

  const grades = [...new Set(
    (feeScheduleOptions.classes || [])
      .filter((item) => !feeScheduleYearEl.value || item.schoolYearId === feeScheduleYearEl.value)
      .map((item) => item.grade)
      .filter(Boolean),
  )].sort((a, b) => a.localeCompare(b, "vi", { numeric: true }));
  feeScheduleGradeEl.innerHTML = [
    `<option value="">Toàn năm học</option>`,
    ...grades.map((grade) => `<option value="${escapeAttr(grade)}">Khối ${escapeHtml(grade)}</option>`),
  ].join("");
  feeScheduleGradeEl.value = optionValueOrEmpty(feeScheduleGradeEl, selectedGrade);

  renderFeeScheduleClassFilter(selectedClass);
}

function renderFeeScheduleClassFilter(selectedClass = feeScheduleClassEl.value) {
  const classes = (feeScheduleOptions.classes || []).filter((item) => {
    if (feeScheduleYearEl.value && item.schoolYearId !== feeScheduleYearEl.value) return false;
    if (feeScheduleGradeEl.value && item.grade !== feeScheduleGradeEl.value) return false;
    return true;
  });
  feeScheduleClassEl.innerHTML = [
    `<option value="">Theo năm/khối</option>`,
    ...classes.map(
      (item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.schoolYearCode)} · ${escapeHtml(item.name)}</option>`,
    ),
  ].join("");
  feeScheduleClassEl.value = optionValueOrEmpty(feeScheduleClassEl, selectedClass);
}

function renderFeeScheduleItems(feeTypes) {
  const types = feeTypes?.length ? feeTypes : defaultFeeTypes;
  feeScheduleItemsEl.innerHTML = types
    .map(
      (item) => `
        <tr data-fee-type-id="${escapeAttr(item.id || "")}" data-fee-type-code="${escapeAttr(item.code || "")}" data-display-order="${Number(item.displayOrder || 0)}">
          <td><span class="tag">${escapeHtml(item.code || "")}</span></td>
          <td><input data-fee-schedule-field="labelVi" value="${escapeAttr(item.labelVi || "")}" /></td>
          <td><input data-fee-schedule-field="labelEn" value="${escapeAttr(item.labelEn || "")}" /></td>
          <td><input data-fee-schedule-field="amount" type="number" min="0" step="1000" value="0" /></td>
        </tr>
      `,
    )
    .join("");
  feeScheduleItemsEl.querySelectorAll('[data-fee-schedule-field="amount"]').forEach((input) => {
    input.addEventListener("input", updateFeeScheduleItemTotal);
  });
  updateFeeScheduleItemTotal();
}

function updateFeeScheduleItemTotal() {
  const total = collectFeeScheduleItems().reduce((sum, item) => sum + Number(item.amount || 0), 0);
  feeScheduleItemTotalEl.textContent = formatMoney(total);
}

function collectFeeScheduleItems() {
  return [...feeScheduleItemsEl.querySelectorAll("tr")]
    .map((row) => ({
      feeTypeId: row.dataset.feeTypeId || "",
      feeTypeCode: row.dataset.feeTypeCode || "",
      labelVi: row.querySelector('[data-fee-schedule-field="labelVi"]').value.trim(),
      labelEn: row.querySelector('[data-fee-schedule-field="labelEn"]').value.trim(),
      amount: Number(row.querySelector('[data-fee-schedule-field="amount"]').value || 0),
      displayOrder: Number(row.dataset.displayOrder || 0),
    }))
    .filter((item) => item.amount > 0);
}

function collectFeeScheduleDraft() {
  const classId = feeScheduleClassEl.value;
  return {
    schoolYearId: feeScheduleYearEl.value,
    classId,
    grade: classId ? "" : feeScheduleGradeEl.value,
    periodCode: feeSchedulePeriodEl.value.trim(),
    month: Number(feeScheduleMonthEl.value || 0),
    name: feeScheduleNameEl.value.trim(),
    notes: feeScheduleNotesEl.value.trim(),
    status: feeScheduleStatusEl.value || "draft",
    items: collectFeeScheduleItems(),
    adjustments: parseFeeAdjustmentsCsv(),
  };
}

function parseFeeAdjustmentsCsv() {
  const lines = feeAdjustmentsCsvEl.value
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter(Boolean);
  if (!lines.length) {
    return [];
  }

  let header = ["student_code", "adjustment_type", "fee_type_code", "amount", "reason"];
  const first = lines[0].split(",").map((part) => headerKeyClient(part));
  const hasHeader = first.includes("student_code") || first.includes("student_id");
  const dataLines = hasHeader ? lines.slice(1) : lines;
  if (hasHeader) {
    header = first;
  }

  const adjustments = dataLines
    .map((line) => {
      const parts = line.split(",").map((part) => part.trim());
      const field = (name, fallbackIndex) => {
        const index = header.indexOf(name);
        if (index >= 0) return parts[index] || "";
        return parts[fallbackIndex] || "";
      };
      const reasonIndex = header.indexOf("reason");
      const reason = reasonIndex >= 0 ? parts.slice(reasonIndex).join(",").trim() : parts.slice(4).join(",").trim();
      return {
        studentCode: field("student_code", 0).toUpperCase(),
        adjustmentType: headerKeyClient(field("adjustment_type", 1)),
        feeTypeCode: headerKeyClient(field("fee_type_code", 2)),
        amount: parseMoneyInput(field("amount", 3)),
        reason,
      };
    })
    .filter((item) => item.studentCode || item.adjustmentType || item.amount || item.reason);
  feeAdjustmentCountEl.textContent = `${adjustments.length} điều chỉnh`;
  return adjustments;
}

async function previewFeeSchedule() {
  setFeeScheduleStatus("Đang preview", "busy");
  const res = await fetch("/api/v1/fee-schedules/preview", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(collectFeeScheduleDraft()),
  });
  const text = await res.text();
  let data = null;
  try {
    data = JSON.parse(text);
  } catch {
    data = null;
  }
  if (!res.ok || !data) {
    renderFeeSchedulePreview({ rows: [], issues: [{ message: text || "Preview failed", type: "preview_failed" }] });
    setFeeScheduleStatus("Lỗi", "error");
    return;
  }
  renderFeeSchedulePreview(data);
  setFeeScheduleStatus(data.issues?.length ? "Có lỗi" : "Đã preview", data.issues?.length ? "error" : "ready");
}

async function saveFeeSchedule() {
  setFeeScheduleStatus("Đang lưu", "busy");
  const res = await fetch("/api/v1/fee-schedules/save", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(collectFeeScheduleDraft()),
  });
  const text = await res.text();
  let data = null;
  try {
    data = JSON.parse(text);
  } catch {
    data = null;
  }
  if (!res.ok) {
    const issues = data?.issues || [{ type: "save_failed", message: text || "Save failed" }];
    renderFeeSchedulePreview({ rows: [], issues });
    setFeeScheduleStatus("Lỗi", "error");
    return;
  }
  renderFeeSchedulePreview(data.preview || null);
  renderFeeSchedules(data.schedules || []);
  setFeeScheduleStatus("Đã lưu", "ready");
}

function renderFeeSchedulePreview(data) {
  const rows = data?.rows || [];
  const issues = data?.issues || [];
  feeSchedulePreviewCountEl.textContent = `${rows.length} học sinh`;
  if (!data) {
    feeScheduleSummaryEl.className = "fee-schedule-summary";
    feeScheduleSummaryEl.textContent = "Chưa có preview";
    feeSchedulePreviewRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Chưa có preview</td></tr>`;
    return;
  }

  const summary = data.summary || {};
  feeScheduleSummaryEl.className = `fee-schedule-summary${issues.length ? " error" : ""}`;
  const summaryText = [
    `${summary.studentCount || rows.length || 0} học sinh`,
    `mặc định ${formatMoney(summary.baseAmount || 0)}`,
    `điều chỉnh ${formatMoney(summary.adjustments || 0)}`,
    `phải thu ${formatMoney(summary.totalAmount || 0)}`,
  ].join(" · ");
  const issueList = issues.length
    ? `<div class="fee-issue-list">${issues
        .map((issue) => `<div><strong>${escapeHtml(issue.type || "issue")}</strong> ${escapeHtml(issue.studentCode || "")} ${escapeHtml(issue.message || "")}</div>`)
        .join("")}</div>`
    : "";
  feeScheduleSummaryEl.innerHTML = `<div>${escapeHtml(summaryText)}</div>${issueList}`;

  feeSchedulePreviewRowsEl.innerHTML = rows
    .map(
      (row) => `
        <tr>
          <td><strong>${escapeHtml(row.studentCode || "")}</strong></td>
          <td>${escapeHtml(row.studentName || "")}</td>
          <td>${escapeHtml(row.className || "")}</td>
          <td>${formatMoney(row.baseAmount || 0)}</td>
          <td class="${Number(row.adjustmentAmount || 0) < 0 ? "fee-negative" : ""}">${formatMoney(row.adjustmentAmount || 0)}</td>
          <td class="money">${formatMoney(row.totalAmount || 0)}</td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    feeSchedulePreviewRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Không có học sinh trong phạm vi này</td></tr>`;
  }
}

function renderFeeSchedules(schedules) {
  feeScheduleCountEl.textContent = `${schedules.length} bảng phí`;
  feeSchedulesEl.innerHTML = schedules
    .map((schedule) => {
      const scope = schedule.className || (schedule.grade ? `Khối ${schedule.grade}` : "Toàn năm học");
      const period = [schedule.periodCode, schedule.month ? `T${schedule.month}` : ""].filter(Boolean).join(" · ");
      return `
        <div class="fee-schedule-list-item">
          <div>
            <strong>${escapeHtml(schedule.name || schedule.periodCode || "Bảng phí")}</strong>
            <span>${escapeHtml(schedule.schoolYearCode || "")} · ${escapeHtml(scope)} · ${escapeHtml(period)}</span>
          </div>
          <div>
            <span class="tag">${escapeHtml(schedule.status || "draft")}</span>
            <span class="money">${formatMoney(schedule.itemTotal || 0)}</span>
            <span>${Number(schedule.adjustmentCount || 0)} điều chỉnh</span>
          </div>
        </div>
      `;
    })
    .join("");
  if (!schedules.length) {
    feeSchedulesEl.innerHTML = `<div class="empty-cell fee-list-empty">Chưa có bảng phí đã lưu</div>`;
  }
}

async function loadInvoices(force = false) {
  if (!force && invoicesLoaded) {
    return;
  }
  const loaded = await loadInvoiceOptions();
  if (loaded) {
    await loadInvoiceList();
  }
}

async function loadInvoiceOptions() {
  setInvoiceStatus("Đang tải", "busy");
  const res = await fetch("/api/v1/invoices/options");
  const text = await res.text();
  if (!res.ok) {
    invoicesLoaded = false;
    invoiceOptions = { schedules: [], schoolYears: [], classes: [] };
    renderInvoiceControls();
    renderInvoicePreview(null);
    renderInvoices([]);
    setInvoiceStatus(text || "Chưa cấu hình DB", "error");
    return false;
  }
  invoiceOptions = JSON.parse(text);
  invoicesLoaded = true;
  renderInvoiceControls();
  setInvoiceStatus("Sẵn sàng", "ready");
  return true;
}

async function loadInvoiceList() {
  if (!invoicesLoaded) {
    return;
  }
  const res = await fetch("/api/v1/invoices");
  const text = await res.text();
  if (!res.ok) {
    renderInvoices([]);
    setInvoiceStatus(text || "Không tải được hóa đơn", "error");
    return;
  }
  const data = JSON.parse(text);
  renderInvoices(data.invoices || []);
  setInvoiceStatus("Sẵn sàng", "ready");
}

function renderInvoiceControls() {
  const selectedSchedule = invoiceScheduleEl.value;
  invoiceScheduleEl.innerHTML = [
    `<option value="">Chọn bảng phí</option>`,
    ...(invoiceOptions.schedules || []).map((schedule) => {
      const scope = schedule.className || (schedule.grade ? `Khối ${schedule.grade}` : "Toàn năm học");
      const period = [schedule.periodCode, schedule.month ? `T${schedule.month}` : ""].filter(Boolean).join(" · ");
      const label = `${schedule.schoolYearCode || ""} · ${scope} · ${period || schedule.name || schedule.id}`;
      return `<option value="${escapeAttr(schedule.id)}">${escapeHtml(label)}</option>`;
    }),
  ].join("");
  invoiceScheduleEl.value = optionValueOrEmpty(invoiceScheduleEl, selectedSchedule);

  const selectedBank = invoiceBankBinEl.value || "970415";
  invoiceBankBinEl.innerHTML = bankOptions(selectedBank);
  invoiceBankBinEl.value = optionValueOrEmpty(invoiceBankBinEl, selectedBank);
  if (!invoiceIssueDateEl.value) {
    invoiceIssueDateEl.value = new Date().toISOString().slice(0, 10);
  }
}

function collectInvoiceRequest() {
  return {
    feeScheduleId: invoiceScheduleEl.value,
    bankBin: invoiceBankBinEl.value,
    bankAccount: invoiceBankAccountEl.value.trim(),
    issueDate: invoiceIssueDateEl.value,
    dueDate: invoiceDueDateEl.value,
    regenerate: invoiceRegenerateEl.value === "true",
  };
}

async function previewInvoices() {
  setInvoiceStatus("Đang preview", "busy");
  const res = await fetch("/api/v1/invoices/preview", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(collectInvoiceRequest()),
  });
  const text = await res.text();
  let data = null;
  try {
    data = JSON.parse(text);
  } catch {
    data = null;
  }
  if (!res.ok || !data) {
    renderInvoicePreview({ rows: [], issues: [{ type: "preview_failed", message: text || "Preview failed" }] });
    setInvoiceStatus("Lỗi", "error");
    return;
  }
  renderInvoicePreview(data);
  setInvoiceStatus(data.issues?.length ? "Có lỗi" : "Đã preview", data.issues?.length ? "error" : "ready");
}

async function generateInvoices() {
  if (!window.confirm("Sinh hóa đơn sẽ ghi dữ liệu invoice vào database. Tiếp tục?")) {
    return;
  }
  setInvoiceStatus("Đang sinh hóa đơn", "busy");
  const res = await fetch("/api/v1/invoices/generate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(collectInvoiceRequest()),
  });
  const text = await res.text();
  let data = null;
  try {
    data = JSON.parse(text);
  } catch {
    data = null;
  }
  if (!res.ok || !data) {
    renderInvoicePreview(data?.preview || { rows: [], issues: data?.issues || [{ type: "generate_failed", message: text || "Generate failed" }] });
    setInvoiceStatus("Lỗi", "error");
    return;
  }
  renderInvoicePreview(data.preview || null);
  renderInvoices(data.invoices || []);
  setInvoiceStatus("Đã sinh hóa đơn", "ready");
}

function renderInvoicePreview(data) {
  const rows = data?.rows || [];
  const issues = data?.issues || [];
  invoicePreviewCountEl.textContent = `${rows.length} hóa đơn`;
  if (!data) {
    invoicePreviewSummaryEl.className = "invoice-summary";
    invoicePreviewSummaryEl.textContent = "Chưa có preview";
    invoicePreviewRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Chưa có preview</td></tr>`;
    return;
  }

  const summary = data.summary || {};
  invoicePreviewSummaryEl.className = `invoice-summary${issues.length ? " error" : ""}`;
  const summaryText = [
    `${summary.studentCount || rows.length || 0} học sinh`,
    `đã có ${summary.existingCount || 0}`,
    `mặc định ${formatMoney(summary.baseAmount || 0)}`,
    `điều chỉnh ${formatMoney(summary.adjustmentAmount || 0)}`,
    `phải thu ${formatMoney(summary.totalAmount || 0)}`,
  ].join(" · ");
  const issueList = issues.length
    ? `<div class="fee-issue-list">${issues
        .map((issue) => `<div><strong>${escapeHtml(issue.type || "issue")}</strong> ${escapeHtml(issue.studentCode || "")} ${escapeHtml(issue.message || "")}</div>`)
        .join("")}</div>`
    : "";
  invoicePreviewSummaryEl.innerHTML = `<div>${escapeHtml(summaryText)}</div>${issueList}`;

  invoicePreviewRowsEl.innerHTML = rows
    .map(
      (row) => `
        <tr>
          <td><strong>${escapeHtml(row.invoiceCode || "")}</strong></td>
          <td>${escapeHtml(row.studentCode || "")} · ${escapeHtml(row.studentName || "")}</td>
          <td>${escapeHtml(row.className || "")}</td>
          <td>${escapeHtml(row.periodCode || "")}</td>
          <td class="money">${formatMoney(row.totalAmount || 0)}</td>
          <td><span class="tag">${escapeHtml(row.existing ? "existing" : row.status || "unpaid")}</span></td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    invoicePreviewRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Không có hóa đơn trong preview</td></tr>`;
  }
}

function renderInvoices(invoices) {
  invoiceCountEl.textContent = `${invoices.length} hóa đơn`;
  invoiceRowsEl.innerHTML = invoices
    .map(
      (invoice) => `
        <tr>
          <td><strong>${escapeHtml(invoice.invoiceCode || "")}</strong></td>
          <td>${escapeHtml(invoice.studentCode || "")} · ${escapeHtml(invoice.studentName || "")}</td>
          <td>${escapeHtml(invoice.className || "")}</td>
          <td>${escapeHtml(invoice.periodCode || "")}</td>
          <td class="money">${formatMoney(invoice.totalAmount || 0)}</td>
          <td><span class="tag">${escapeHtml(invoice.status || "unpaid")}</span></td>
          <td>
            <div class="invoice-actions">
              <button type="button" data-invoice-qr="${escapeAttr(invoice.id || "")}">${muiIcon("qr_code")}<span>QR</span></button>
              <a class="button-link" href="/api/v1/invoices/pdf?id=${encodeURIComponent(invoice.id || "")}" target="_blank" rel="noreferrer">${muiIcon("picture_as_pdf")}<span>PDF</span></a>
            </div>
          </td>
        </tr>
      `,
    )
    .join("");
  if (!invoices.length) {
    invoiceRowsEl.innerHTML = `<tr><td colspan="7" class="empty-cell">Chưa có hóa đơn</td></tr>`;
  }
  invoiceRowsEl.querySelectorAll("[data-invoice-qr]").forEach((button) => {
    button.addEventListener("click", () => loadInvoicePayment(button.dataset.invoiceQr));
  });
}

async function loadInvoicePayment(invoiceId) {
  if (!invoiceId) {
    return;
  }
  invoicePaymentStatusEl.textContent = "Đang tải";
  invoicePaymentStatusEl.dataset.tone = "busy";
  const res = await fetch(`/api/v1/invoices/payment?id=${encodeURIComponent(invoiceId)}`);
  const text = await res.text();
  let item = null;
  try {
    item = JSON.parse(text);
  } catch {
    item = null;
  }
  if (!res.ok || !item || !item.vietqr) {
    invoicePaymentStatusEl.textContent = "Lỗi";
    invoicePaymentStatusEl.dataset.tone = "error";
    invoicePaymentPreviewEl.className = "preview-empty";
    invoicePaymentPreviewEl.textContent = item?.errors?.join(", ") || text || "Không tải được QR hóa đơn";
    return;
  }
  invoicePaymentStatusEl.textContent = "Sẵn sàng";
  invoicePaymentStatusEl.dataset.tone = "ready";
  invoicePaymentPreviewEl.className = "preview-content";
  invoicePaymentPreviewEl.innerHTML = `
    <img src="${item.qrData}" alt="QR hóa đơn" />
    <div class="preview-meta">
      <div class="meta-row">${muiIcon("person")}<strong>${escapeHtml(item.studentName || "")}</strong></div>
      <div class="meta-row">${muiIcon("school")}<span>${escapeHtml(item.className || "")}</span></div>
      <div class="meta-row">${muiIcon("confirmation_number")}<span>${escapeHtml(item.billNumber || "")}</span></div>
      <div class="meta-row">${muiIcon("account_balance")}<span>${escapeHtml(item.bankName || item.bankBin)} / ${escapeHtml(item.bankAccount)}</span></div>
      <div class="meta-row money">${muiIcon("payments")}<span>${formatMoney(item.amount || 0)}</span></div>
    </div>
    <textarea class="payload" readonly>${escapeHtml(item.vietqr || "")}</textarea>
  `;
}

function setMasterStatus(message, tone = "ready") {
  masterDataStatusEl.textContent = message;
  masterDataStatusEl.dataset.tone = tone;
}

function setFeeScheduleStatus(message, tone = "ready") {
  feeScheduleLoadStatusEl.textContent = message;
  feeScheduleLoadStatusEl.dataset.tone = tone;
}

function setInvoiceStatus(message, tone = "ready") {
  invoiceStatusEl.textContent = message;
  invoiceStatusEl.dataset.tone = tone;
}


function collectEmailConfig() {
  return {
    provider: emailProviderEl.value || "gmail",
    apiKey: resendApiKeyEl.value.trim(),
    gmailAddress: gmailAddressEl.value.trim(),
    gmailAppPassword: gmailAppPasswordEl.value.trim(),
    from: emailFromEl.value.trim(),
    replyTo: emailReplyToEl.value.trim(),
    subject: emailSubjectEl.value.trim(),
    schoolName: schoolNameEl.value.trim(),
    schoolNameEn: schoolNameEnEl.value.trim(),
    paymentPeriod: paymentPeriodEl.value.trim(),
    publicBaseUrl: publicBaseUrlEl.value.trim(),
  };
}

async function saveEmailConfig() {
  setEmailStatus("Saving");
  const res = await fetch("/api/v1/email/config", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(collectEmailConfig()),
  });
  const text = await res.text();
  if (!res.ok) {
    setEmailStatus(text || "Save failed", true);
    return false;
  }
  const config = JSON.parse(text);
  savedEmailConfig = config;
  emailProviderEl.value = config.provider || emailProviderEl.value || "gmail";
  gmailAppPasswordEl.value = "";
  gmailAppPasswordEl.placeholder = config.hasGmailAppPassword ? config.gmailAppPasswordMasked : "16-character app password";
  resendApiKeyEl.value = "";
  resendApiKeyEl.placeholder = config.hasApiKey ? config.apiKeyMasked : "re_...";
  renderEmailConfigStatus(config);
  setEmailStatus("Đã lưu cấu hình");
  return true;
}

async function previewEmail() {
  const rows = collectRows();
  if (!rows.length) return;
  const res = await fetch("/api/v1/email/preview", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      rows: [rows[0]],
      template: emailTemplateEl.value,
      config: collectEmailConfig(),
    }),
  });
  const text = await res.text();
  if (!res.ok) {
    emailPreviewSubjectEl.textContent = "Lỗi";
    setEmailStatus(text || "Preview failed", true);
    return;
  }
  const data = JSON.parse(text);
  emailPreviewSubjectEl.textContent = data.subject || "Preview";
  emailPreviewFrameEl.srcdoc = data.html || "";
}

async function sendEmails(dryRun) {
  const saved = await saveEmailConfig();
  if (!saved) return;
  await generate();
  const rows = collectRows();
  setEmailStatus(dryRun ? "Checking" : "Sending");
  const res = await fetch("/api/v1/email/send", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      rows,
      template: emailTemplateEl.value,
      dryRun,
    }),
  });
  const text = await res.text();
  if (!res.ok) {
    setEmailStatus(text || "Send failed", true);
    return;
  }
  const data = JSON.parse(text);
  renderEmailResults(data.results || []);
}

function renderEmailResults(results) {
  const sent = results.filter((item) => item.status === "sent").length;
  const dry = results.filter((item) => item.status === "dry_run").length;
  const skipped = results.filter((item) => item.status === "skipped").length;
  const errors = results.filter((item) => item.status === "error").length;
  const rows = results
    .map((item) => {
      const message = item.error || item.resendId || item.messageId || item.status;
      return `<div><strong>${escapeHtml(item.status)}</strong> ${escapeHtml(item.email || "-")} · ${escapeHtml(item.studentName || "-")} · ${escapeHtml(message)}</div>`;
    })
    .join("");
  emailSendStatusEl.innerHTML = `
    <div class="email-summary">sent ${sent} · dry ${dry} · skipped ${skipped} · errors ${errors}</div>
    <div class="email-result-list">${rows}</div>
  `;
}

function setEmailStatus(message, isError = false) {
  emailSendStatusEl.className = `email-status${isError ? " error" : ""}`;
  emailSendStatusEl.textContent = message;
}

async function saveEmailCron(enabled) {
  const saved = await saveEmailConfig();
  if (!saved) return;
  await generate();
  setCronStatus("Saving");
  const res = await fetch("/api/v1/email/cron", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      rows: collectRows(),
      template: emailTemplateEl.value,
      enabled,
      dailyLimit: Number(cronDailyLimitEl.value || 500),
      sendTime: cronSendTimeEl.value || "08:00",
    }),
  });
  const text = await res.text();
  if (!res.ok) {
    setCronStatus(text || "Save cron failed", true);
    return;
  }
  renderCronStatus(JSON.parse(text));
}

async function disableEmailCron() {
  setCronStatus("Saving");
  const res = await fetch("/api/v1/email/cron", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ enabled: false }),
  });
  const text = await res.text();
  if (!res.ok) {
    setCronStatus(text || "Disable cron failed", true);
    return;
  }
  renderCronStatus(JSON.parse(text));
}

async function runEmailCronNow() {
  if (!window.confirm("Chạy cron sẽ gửi email thật qua provider hiện tại theo giới hạn còn lại. Tiếp tục?")) {
    return;
  }
  setCronStatus("Running");
  const res = await fetch("/api/v1/email/cron/run", { method: "POST" });
  const text = await res.text();
  if (!res.ok) {
    setCronStatus(text || "Run cron failed", true);
    return;
  }
  renderCronStatus(JSON.parse(text));
}

function renderCronStatus(data) {
  cronEnabledEl.value = String(Boolean(data.enabled));
  cronSendTimeEl.value = data.sendTime || "08:00";
  cronDailyLimitEl.value = data.dailyLimit || 500;
  cronQueueSummaryEl.value = `${data.queued || 0} chờ / ${data.sent || 0} gửi / ${data.errors || 0} lỗi`;

  const nextRun = data.nextRunAt ? formatDateTime(data.nextRunAt) : "-";
  const lastRun = data.lastRunAt ? formatDateTime(data.lastRunAt) : "-";
  const sentLast24h = data.sentLast24h ?? data.sentToday ?? 0;
  const summary = [
    data.enabled ? "ACTIVE" : "PAUSED",
    `limit ${data.dailyLimit || 500}/24h`,
    `sent 24h ${sentLast24h}`,
    `next ${nextRun}`,
    `last ${lastRun}`,
  ].join(" · ");
  const recent = (data.lastResults || [])
    .slice(0, 5)
    .map((item) => {
      const message = item.error || item.resendId || item.messageId || item.status;
      return `<div><strong>${escapeHtml(item.status)}</strong> ${escapeHtml(item.email || "-")} · ${escapeHtml(item.studentName || "-")} · ${escapeHtml(message || "")}</div>`;
    })
    .join("");
  cronStatusEl.className = "email-status";
  cronStatusEl.innerHTML = `<div class="email-summary">${escapeHtml(summary)}</div><div class="email-result-list">${recent}</div>`;
}

function setCronStatus(message, isError = false) {
  cronStatusEl.className = `email-status${isError ? " error" : ""}`;
  cronStatusEl.textContent = message;
}

function updateEmailProviderFields() {
  const provider = emailProviderEl.value || "gmail";
  const isGmail = provider === "gmail";
  gmailAddressFieldEl.hidden = !isGmail;
  gmailAppPasswordFieldEl.hidden = !isGmail;
  resendApiKeyFieldEl.hidden = isGmail;
}

function currentEmailConfigStatusData() {
  return {
    ...savedEmailConfig,
    provider: emailProviderEl.value || savedEmailConfig.provider || "gmail",
    gmailAddress: gmailAddressEl.value.trim() || savedEmailConfig.gmailAddress || "",
    hasGmailAppPassword:
      Boolean(gmailAppPasswordEl.value.trim()) ||
      Boolean(savedEmailConfig.hasGmailAppPassword) ||
      gmailAppPasswordEl.placeholder !== "16-character app password",
    hasApiKey:
      Boolean(resendApiKeyEl.value.trim()) ||
      Boolean(savedEmailConfig.hasApiKey) ||
      resendApiKeyEl.placeholder !== "re_...",
  };
}

function renderEmailConfigStatus(config) {
  const provider = config.provider || emailProviderEl.value || "gmail";
  if (provider === "resend") {
    emailConfigStatusEl.textContent = config.hasApiKey ? "Resend đã cấu hình" : "Thiếu Resend key";
    return;
  }
  emailConfigStatusEl.textContent =
    config.gmailAddress && config.hasGmailAppPassword ? "Gmail đã cấu hình" : "Thiếu Gmail SMTP";
}

function renderResults() {
  resultsEl.innerHTML = "";
  resultCountEl.textContent = `${currentItems.filter((item) => item.vietqr).length} mã QR`;

  currentItems.forEach((item) => {
    const card = document.createElement("article");
    card.className = `qr-card${item.id === selectedId ? " active" : ""}`;
    card.tabIndex = 0;
    card.innerHTML = item.vietqr
      ? `
        <img src="${item.qrData}" alt="QR ${escapeAttr(item.studentName)}" />
        <div>
          <h3>${muiIcon("person")}<span>${escapeHtml(item.studentName || "Chưa có tên")}</span></h3>
          <p class="card-meta">${muiIcon("family_restroom")}<span>${escapeHtml(item.parentName || "")}</span></p>
          <p><span class="tag">${muiIcon("account_balance")}<span>${escapeHtml(item.bankName || item.bankBin)}</span></span></p>
          <p class="money">${formatMoney(item.amount)}</p>
        </div>
      `
      : `
        <div>
          <h3>${muiIcon("error")}<span>${escapeHtml(item.studentName || "Dòng lỗi")}</span></h3>
          <p class="error">${escapeHtml((item.errors || []).join(", "))}</p>
        </div>
      `;
    card.addEventListener("click", () => selectItem(item.id));
    card.addEventListener("keydown", (event) => {
      if (event.key === "Enter" || event.key === " ") {
        event.preventDefault();
        selectItem(item.id);
      }
    });
    resultsEl.appendChild(card);
  });
}

function renderPreview() {
  const item = currentItems.find((entry) => entry.id === selectedId);
  if (!item || !item.vietqr) {
    previewEl.className = "preview-empty";
    previewEl.textContent = item?.errors?.join(", ") || "Chưa có QR";
    return;
  }

  const params = new URLSearchParams({
    bankBin: item.bankBin,
    account: item.bankAccount,
    amount: String(item.amount),
    billNumber: item.billNumber || "",
    note: item.note || "",
  });

  previewEl.className = "preview-content";
  previewEl.innerHTML = `
    <img src="${item.qrData}" alt="QR preview" />
    <div class="preview-meta">
      <div class="meta-row">${muiIcon("person")}<strong>${escapeHtml(item.studentName || "Chưa có tên")}</strong></div>
      <div class="meta-row">${muiIcon("family_restroom")}<span>${escapeHtml(item.parentName || "")}</span></div>
      <div class="meta-row">${muiIcon("account_balance")}<span>${escapeHtml(item.bankName || item.bankBin)} / ${escapeHtml(item.bankAccount)}</span></div>
      <div class="meta-row">${muiIcon("confirmation_number")}<span>Bill: ${escapeHtml(item.billNumber || "")}</span></div>
      <div class="meta-row money">${muiIcon("payments")}<span>${formatMoney(item.amount)}</span></div>
      <a class="preview-link" href="/api/v1/qr.png?${params.toString()}" target="_blank" rel="noreferrer">${muiIcon("open_in_new")}<span>Mở PNG scan test</span></a>
    </div>
    <textarea class="payload" readonly>${escapeHtml(item.vietqr)}</textarea>
  `;
}

function selectItem(id) {
  selectedId = id;
  renderResults();
  renderPreview();
}

function updateRowCount() {
  rowCountEl.textContent = `${rowsEl.querySelectorAll("tr").length} bản ghi`;
}

function status(text, tone = "ready") {
  statusEl.textContent = text;
  statusEl.dataset.tone = tone;
}

function formatMoney(amount) {
  return new Intl.NumberFormat("vi-VN", {
    style: "currency",
    currency: "VND",
    maximumFractionDigits: 0,
  }).format(Number(amount || 0));
}

function parseMoneyInput(value) {
  const normalized = String(value || "").replace(/[.,\s_]/g, "");
  const amount = Number(normalized || 0);
  return Number.isFinite(amount) ? amount : 0;
}

function headerKeyClient(value) {
  return String(value || "")
    .trim()
    .toLowerCase()
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/đ/g, "d")
    .replace(/[^a-z0-9]+/g, "_")
    .replace(/^_+|_+$/g, "");
}

function formatDateTime(value) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return "-";
  }
  return new Intl.DateTimeFormat("vi-VN", {
    dateStyle: "short",
    timeStyle: "short",
  }).format(date);
}

function escapeHtml(value) {
  return String(value).replace(/[&<>"']/g, (char) => {
    return {
      "&": "&amp;",
      "<": "&lt;",
      ">": "&gt;",
      '"': "&quot;",
      "'": "&#039;",
    }[char];
  });
}

function escapeAttr(value) {
  return escapeHtml(value).replace(/`/g, "&#096;");
}

generateBtn.addEventListener("click", generate);

toggleFeeColumnBtn.addEventListener("click", () => {
  setFeeColumnCollapsed(!feeColumnCollapsed);
});

addRowBtn.addEventListener("click", async () => {
  rowsEl.appendChild(
    rowTemplate({
      id: crypto.randomUUID(),
      bankBin: "970415",
      paymentItems: currentFeeTemplateItems(),
    }),
  );
  updateRowCount();
  await activateTab("paymentsTab");
});

loadSampleBtn.addEventListener("click", async () => {
  renderRows(sampleRows);
  await activateTab("paymentsTab");
  await generate();
});

applyFeeTemplateBtn.addEventListener("click", async () => {
  const templateItems = currentFeeTemplateItems();
  rowsEl.querySelectorAll("tr").forEach((tr) => {
    const feeCell = tr.children[6];
    feeCell.innerHTML = paymentItemsTemplate(clonePaymentItems(templateItems));
    initializeFeeEditor(feeCell);
  });
  await generate();
  await previewEmail();
});

csvFileEl.addEventListener("change", async () => {
  const file = csvFileEl.files[0];
  if (!file) return;

  status("Đang import", "busy");
  const form = new FormData();
  form.append("file", file);

  const res = await fetch("/api/v1/import/csv", {
    method: "POST",
    body: form,
  });

  if (!res.ok) {
    status("Lỗi", "error");
    previewEl.className = "preview-empty";
    previewEl.textContent = await res.text();
    return;
  }

  const data = await res.json();
  renderRows(data.rows || []);
  await activateTab("paymentsTab");
  await generate();
  csvFileEl.value = "";
});

tabButtons.forEach((button) => {
  button.addEventListener("click", () => activateTab(button.dataset.tabTarget));
});

masterSchoolYearFilterEl.addEventListener("change", async () => {
  renderMasterFilters();
  await loadMasterStudents();
});
masterGradeFilterEl.addEventListener("change", async () => {
  renderMasterClassFilter();
  await loadMasterStudents();
});
masterClassFilterEl.addEventListener("change", loadMasterStudents);
masterSearchEl.addEventListener("input", () => {
  window.clearTimeout(masterSearchEl.dataset.timer);
  masterSearchEl.dataset.timer = window.setTimeout(loadMasterStudents, 250);
});
refreshMasterDataBtn.addEventListener("click", () => loadMasterData(true));
masterCsvFileEl.addEventListener("change", () => {
  const file = masterCsvFileEl.files[0];
  masterImportSummaryEl.textContent = file ? file.name : "Chưa có file import";
});
checkMasterImportBtn.addEventListener("click", () => submitMasterImport(false));
applyMasterImportBtn.addEventListener("click", () => submitMasterImport(true));

refreshFeeSchedulesBtn.addEventListener("click", () => loadFeeSchedules(true));
feeScheduleYearEl.addEventListener("change", async () => {
  renderFeeScheduleControls();
  await loadFeeScheduleList();
});
feeScheduleGradeEl.addEventListener("change", async () => {
  renderFeeScheduleClassFilter();
  await loadFeeScheduleList();
});
feeScheduleClassEl.addEventListener("change", loadFeeScheduleList);
feeAdjustmentsCsvEl.addEventListener("input", () => {
  parseFeeAdjustmentsCsv();
});
previewFeeScheduleBtn.addEventListener("click", previewFeeSchedule);
saveFeeScheduleBtn.addEventListener("click", saveFeeSchedule);

refreshInvoicesBtn.addEventListener("click", () => loadInvoices(true));
previewInvoicesBtn.addEventListener("click", previewInvoices);
generateInvoicesBtn.addEventListener("click", generateInvoices);

saveEmailConfigBtn.addEventListener("click", saveEmailConfig);
previewEmailBtn.addEventListener("click", previewEmail);
dryRunEmailBtn.addEventListener("click", () => sendEmails(true));
sendEmailBtn.addEventListener("click", () => sendEmails(false));
emailProviderEl.addEventListener("change", () => {
  updateEmailProviderFields();
  renderEmailConfigStatus(currentEmailConfigStatusData());
});
emailTemplateEl.addEventListener("change", previewEmail);
saveCronBtn.addEventListener("click", () => saveEmailCron(cronEnabledEl.value === "true"));
disableCronBtn.addEventListener("click", disableEmailCron);
runCronNowBtn.addEventListener("click", runEmailCronNow);
