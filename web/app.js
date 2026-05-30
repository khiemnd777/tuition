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
const paymentMappingPanelEl = document.querySelector("#paymentMappingPanel");
const paymentMappingCountEl = document.querySelector("#paymentMappingCount");
const paymentMappingSummaryEl = document.querySelector("#paymentMappingSummary");
const paymentMappingRowsEl = document.querySelector("#paymentMappingRows");
const applyPaymentImportBtn = document.querySelector("#applyPaymentImport");
const cancelPaymentImportBtn = document.querySelector("#cancelPaymentImport");
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
const masterMappingPanelEl = document.querySelector("#masterMappingPanel");
const masterMappingCountEl = document.querySelector("#masterMappingCount");
const masterMappingSummaryEl = document.querySelector("#masterMappingSummary");
const masterMappingRowsEl = document.querySelector("#masterMappingRows");
const clearMasterImportMappingBtn = document.querySelector("#clearMasterImportMapping");
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
const paymentReconStatusEl = document.querySelector("#paymentReconStatus");
const refreshPaymentReconBtn = document.querySelector("#refreshPaymentRecon");
const paymentProviderFilterEl = document.querySelector("#paymentProviderFilter");
const paymentInvoiceStatusFilterEl = document.querySelector("#paymentInvoiceStatusFilter");
const paymentTransactionStatusFilterEl = document.querySelector("#paymentTransactionStatusFilter");
const paymentReconSummaryEl = document.querySelector("#paymentReconSummary");
const paymentReconInvoiceCountEl = document.querySelector("#paymentReconInvoiceCount");
const paymentReconInvoiceRowsEl = document.querySelector("#paymentReconInvoiceRows");
const paymentReconTransactionCountEl = document.querySelector("#paymentReconTransactionCount");
const paymentReconTransactionRowsEl = document.querySelector("#paymentReconTransactionRows");
const paymentReconDetailEl = document.querySelector("#paymentReconDetail");
const notificationStatusEl = document.querySelector("#notificationStatus");
const refreshNotificationsBtn = document.querySelector("#refreshNotifications");
const notificationCampaignNameEl = document.querySelector("#notificationCampaignName");
const notificationCampaignTypeEl = document.querySelector("#notificationCampaignType");
const notificationTemplateEl = document.querySelector("#notificationTemplate");
const notificationSchoolYearEl = document.querySelector("#notificationSchoolYear");
const notificationGradeEl = document.querySelector("#notificationGrade");
const notificationClassEl = document.querySelector("#notificationClass");
const notificationPeriodEl = document.querySelector("#notificationPeriod");
const notificationInvoiceStatusEl = document.querySelector("#notificationInvoiceStatus");
const notificationDueBeforeEl = document.querySelector("#notificationDueBefore");
const previewNotificationsBtn = document.querySelector("#previewNotifications");
const saveNotificationCampaignBtn = document.querySelector("#saveNotificationCampaign");
const sendNotificationCampaignBtn = document.querySelector("#sendNotificationCampaign");
const notificationSummaryEl = document.querySelector("#notificationSummary");
const notificationRecipientCountEl = document.querySelector("#notificationRecipientCount");
const notificationRecipientsEl = document.querySelector("#notificationRecipients");
const notificationCampaignCountEl = document.querySelector("#notificationCampaignCount");
const notificationCampaignRowsEl = document.querySelector("#notificationCampaignRows");
const notificationLogCountEl = document.querySelector("#notificationLogCount");
const notificationLogsEl = document.querySelector("#notificationLogs");
const adminDashboardStatusEl = document.querySelector("#adminDashboardStatus");
const refreshAdminDashboardBtn = document.querySelector("#refreshAdminDashboard");
const adminDashboardYearEl = document.querySelector("#adminDashboardYear");
const adminDashboardGradeEl = document.querySelector("#adminDashboardGrade");
const adminDashboardClassEl = document.querySelector("#adminDashboardClass");
const adminDashboardPeriodEl = document.querySelector("#adminDashboardPeriod");
const adminDashboardMonthEl = document.querySelector("#adminDashboardMonth");
const adminDashboardInvoiceStatusEl = document.querySelector("#adminDashboardInvoiceStatus");
const adminDashboardMetricsEl = document.querySelector("#adminDashboardMetrics");
const adminTopClassCountEl = document.querySelector("#adminTopClassCount");
const adminTopClassRowsEl = document.querySelector("#adminTopClassRows");
const adminAttentionCountEl = document.querySelector("#adminAttentionCount");
const adminAttentionRowsEl = document.querySelector("#adminAttentionRows");
const adminReportsStatusEl = document.querySelector("#adminReportsStatus");
const refreshAdminReportsBtn = document.querySelector("#refreshAdminReports");
const adminReportsYearEl = document.querySelector("#adminReportsYear");
const adminReportsGradeEl = document.querySelector("#adminReportsGrade");
const adminReportsClassEl = document.querySelector("#adminReportsClass");
const adminReportsPeriodEl = document.querySelector("#adminReportsPeriod");
const adminReportsMonthEl = document.querySelector("#adminReportsMonth");
const adminReportsInvoiceStatusEl = document.querySelector("#adminReportsInvoiceStatus");
const adminReportsSummaryEl = document.querySelector("#adminReportsSummary");
const adminReportClassCountEl = document.querySelector("#adminReportClassCount");
const adminReportClassRowsEl = document.querySelector("#adminReportClassRows");
const adminReportInvoiceCountEl = document.querySelector("#adminReportInvoiceCount");
const adminReportInvoiceRowsEl = document.querySelector("#adminReportInvoiceRows");
const adminUsersStatusEl = document.querySelector("#adminUsersStatus");
const refreshAdminUsersBtn = document.querySelector("#refreshAdminUsers");
const adminUserIdEl = document.querySelector("#adminUserId");
const adminUserEmailEl = document.querySelector("#adminUserEmail");
const adminUserDisplayNameEl = document.querySelector("#adminUserDisplayName");
const adminUserStatusEl = document.querySelector("#adminUserStatus");
const adminUserRolesEl = document.querySelector("#adminUserRoles");
const clearAdminUserBtn = document.querySelector("#clearAdminUser");
const saveAdminUserBtn = document.querySelector("#saveAdminUser");
const assignAdminUserRolesBtn = document.querySelector("#assignAdminUserRoles");
const adminUserCountEl = document.querySelector("#adminUserCount");
const adminUserRowsEl = document.querySelector("#adminUserRows");
const adminRoleCountEl = document.querySelector("#adminRoleCount");
const adminRoleListEl = document.querySelector("#adminRoleList");
const tabButtons = [...document.querySelectorAll(".tab-button")];
const tabPanels = [...document.querySelectorAll(".tab-panel")];

let banks = [];
let currentItems = [];
let selectedId = "";
let feeColumnCollapsed = false;
let savedEmailConfig = {};
let masterDataOptions = { schoolYears: [], classes: [] };
let masterDataLoaded = false;
let paymentImportState = null;
let masterImportState = null;
let feeScheduleOptions = { feeTypes: [], schoolYears: [], classes: [] };
let feeSchedulesLoaded = false;
let invoiceOptions = { schedules: [], schoolYears: [], classes: [] };
let invoicesLoaded = false;
let paymentReconciliationLoaded = false;
let paymentReconciliationData = { providers: [], invoices: [], transactions: [], intents: {}, summary: {} };
let notificationLoaded = false;
let notificationOptions = { templates: [], campaigns: [], schoolYears: [], classes: [] };
let notificationPreviewData = { recipients: [], summary: {}, campaign: null, logs: [] };
let currentNotificationCampaignId = "";
let adminOptions = { schoolYears: [], classes: [] };
let adminDashboardLoaded = false;
let adminReportsLoaded = false;
let adminUsersLoaded = false;
let adminUsersData = { users: [], roles: [], permissions: [] };

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
  renderPaymentReconciliation(null);
  renderNotificationControls();
  renderNotificationPreview(null);
  renderNotificationCampaigns([]);
  renderNotificationLogs([]);
  renderAdminFilters("dashboard");
  renderAdminFilters("reports");
  renderAdminDashboard(null);
  renderAdminReports(null);
  renderAdminUsers(null);
  renderFeeTemplate(defaultPaymentItems);
  renderRows(sampleRows);
  await loadAdminDashboard();
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
  if (targetId === "dashboardTab") {
    await loadAdminDashboard();
  }
  if (targetId === "masterDataTab") {
    await loadMasterData();
  }
  if (targetId === "feeTemplateTab") {
    await loadFeeSchedules();
  }
  if (targetId === "invoiceTab") {
    await loadInvoices();
  }
  if (targetId === "reconciliationTab") {
    await loadPaymentReconciliation();
  }
  if (targetId === "notificationTab") {
    await loadNotifications();
  }
  if (targetId === "reportsTab") {
    await loadAdminReports();
  }
  if (targetId === "usersTab") {
    await loadAdminUsers();
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

async function prepareImportMapping(target, file) {
  if (!file) return;
  if (target === "payments") {
    await activateTab("paymentsTab");
    status("Đang đọc file", "busy");
  } else {
    setMasterStatus("Đang đọc file", "busy");
  }

  const form = new FormData();
  form.append("file", file);
  const res = await fetch(`/api/v1/import/fields?target=${encodeURIComponent(target)}`, {
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
  if (!res.ok || !data) {
    if (target === "payments") {
      status("Lỗi", "error");
      previewEl.className = "preview-empty";
      previewEl.textContent = text || "Không đọc được file import";
    } else {
      setMasterStatus(text || "Không đọc được file import", "error");
      masterImportSummaryEl.textContent = text || "Không đọc được file import";
    }
    return;
  }

  const state = {
    target,
    file,
    headers: data.headers || [],
    fields: data.fields || [],
    mapping: data.suggestedMapping || {},
    preview: data.preview || [],
  };
  if (target === "payments") {
    paymentImportState = state;
    renderImportMapping(state);
    status("Đã đọc file", "ready");
  } else {
    masterImportState = state;
    renderImportMapping(state);
    masterImportSummaryEl.className = "master-import-summary";
    masterImportSummaryEl.textContent = `${file.name} · kiểm tra mapping trước khi import`;
    setMasterStatus("Đã đọc file", "ready");
  }
}

function renderImportMapping(state) {
  const els = importMappingElements(state.target);
  els.panel.hidden = false;
  els.count.textContent = `${state.headers.length} cột`;
  els.summary.textContent = `${state.file.name} · ${state.headers.length} cột · preview ${state.preview.length} dòng`;
  els.rows.innerHTML = state.headers
    .map((header, idx) => {
      const selected = state.mapping[header] || "";
      return `
        <tr>
          <td><strong>${escapeHtml(header || `Cột ${idx + 1}`)}</strong></td>
          <td>
            <select data-import-source="${escapeAttr(header)}">
              ${importMappingOptions(state.fields, selected)}
            </select>
          </td>
          <td>${escapeHtml(importSampleValues(state, header) || "-")}</td>
        </tr>
      `;
    })
    .join("");
  els.rows.querySelectorAll("select[data-import-source]").forEach((select) => {
    select.addEventListener("change", () => {
      state.mapping[select.dataset.importSource] = select.value;
    });
  });
}

function importMappingElements(target) {
  if (target === "master_data") {
    return {
      panel: masterMappingPanelEl,
      count: masterMappingCountEl,
      summary: masterMappingSummaryEl,
      rows: masterMappingRowsEl,
    };
  }
  return {
    panel: paymentMappingPanelEl,
    count: paymentMappingCountEl,
    summary: paymentMappingSummaryEl,
    rows: paymentMappingRowsEl,
  };
}

function importMappingOptions(fields, selected) {
  return [
    `<option value="">Bỏ qua</option>`,
    ...fields.map((field) => {
      const label = `${field.label}${field.required ? " *" : ""}`;
      const isSelected = field.key === selected ? "selected" : "";
      return `<option value="${escapeAttr(field.key)}" ${isSelected}>${escapeHtml(label)}</option>`;
    }),
  ].join("");
}

function importSampleValues(state, header) {
  return state.preview
    .map((row) => row.values?.[header] || "")
    .filter(Boolean)
    .slice(0, 3)
    .join(" · ");
}

function collectImportMapping(state) {
  if (!state) return {};
  const els = importMappingElements(state.target);
  const mapping = {};
  els.rows.querySelectorAll("select[data-import-source]").forEach((select) => {
    mapping[select.dataset.importSource] = select.value;
  });
  return mapping;
}

async function submitPaymentImport() {
  if (!paymentImportState?.file) {
    status("Chưa chọn file", "error");
    return;
  }
  status("Đang import", "busy");
  const form = new FormData();
  form.append("file", paymentImportState.file);
  form.append("mapping", JSON.stringify(collectImportMapping(paymentImportState)));
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
  clearPaymentImport();
  await activateTab("paymentsTab");
  await generate();
}

function clearPaymentImport() {
  paymentImportState = null;
  paymentMappingPanelEl.hidden = true;
  paymentMappingRowsEl.innerHTML = "";
  paymentMappingSummaryEl.textContent = "Chưa có file import";
  paymentMappingCountEl.textContent = "0 cột";
  csvFileEl.value = "";
}

function clearMasterImport() {
  masterImportState = null;
  masterMappingPanelEl.hidden = true;
  masterMappingRowsEl.innerHTML = "";
  masterMappingSummaryEl.textContent = "Chưa có file import";
  masterMappingCountEl.textContent = "0 cột";
  masterImportSummaryEl.className = "master-import-summary";
  masterImportSummaryEl.textContent = "Chưa có file import";
  masterCsvFileEl.value = "";
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

function adminFilterElements(kind) {
  if (kind === "reports") {
    return {
      year: adminReportsYearEl,
      grade: adminReportsGradeEl,
      classEl: adminReportsClassEl,
      period: adminReportsPeriodEl,
      month: adminReportsMonthEl,
      status: adminReportsInvoiceStatusEl,
    };
  }
  return {
    year: adminDashboardYearEl,
    grade: adminDashboardGradeEl,
    classEl: adminDashboardClassEl,
    period: adminDashboardPeriodEl,
    month: adminDashboardMonthEl,
    status: adminDashboardInvoiceStatusEl,
  };
}

function setAdminStatus(el, message, tone = "ready") {
  el.textContent = message;
  el.dataset.tone = tone;
}

function renderAdminFilters(kind) {
  const elements = adminFilterElements(kind);
  const selectedYear = elements.year.value;
  const selectedGrade = elements.grade.value;
  const selectedClass = elements.classEl.value;
  elements.year.innerHTML = [
    `<option value="">Tất cả năm học</option>`,
    ...(adminOptions.schoolYears || []).map((item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.code)}</option>`),
  ].join("");
  elements.year.value = optionValueOrEmpty(elements.year, selectedYear);

  const grades = [...new Set(
    (adminOptions.classes || [])
      .filter((item) => !elements.year.value || item.schoolYearId === elements.year.value)
      .map((item) => item.grade)
      .filter(Boolean),
  )].sort((a, b) => a.localeCompare(b, "vi", { numeric: true }));
  elements.grade.innerHTML = [
    `<option value="">Tất cả khối</option>`,
    ...grades.map((grade) => `<option value="${escapeAttr(grade)}">${escapeHtml(grade)}</option>`),
  ].join("");
  elements.grade.value = optionValueOrEmpty(elements.grade, selectedGrade);

  const classes = (adminOptions.classes || []).filter((item) => {
    if (elements.year.value && item.schoolYearId !== elements.year.value) return false;
    if (elements.grade.value && item.grade !== elements.grade.value) return false;
    return true;
  });
  elements.classEl.innerHTML = [
    `<option value="">Tất cả lớp</option>`,
    ...classes.map((item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.schoolYearCode)} · ${escapeHtml(item.name)}</option>`),
  ].join("");
  elements.classEl.value = optionValueOrEmpty(elements.classEl, selectedClass);
}

function adminFilterParams(kind) {
  const elements = adminFilterElements(kind);
  const params = new URLSearchParams();
  if (elements.year.value) params.set("schoolYearId", elements.year.value);
  if (elements.grade.value) params.set("grade", elements.grade.value);
  if (elements.classEl.value) params.set("classId", elements.classEl.value);
  if (elements.period.value.trim()) params.set("periodCode", elements.period.value.trim());
  if (elements.month.value) params.set("month", elements.month.value);
  if (elements.status.value) params.set("status", elements.status.value);
  return params;
}

async function loadAdminDashboard(force = false) {
  if (!force && adminDashboardLoaded) {
    return;
  }
  setAdminStatus(adminDashboardStatusEl, "Đang tải", "busy");
  const params = adminFilterParams("dashboard");
  const query = params.toString();
  const res = await fetch(`/api/v1/admin/dashboard${query ? `?${query}` : ""}`);
  const text = await res.text();
  if (!res.ok) {
    adminDashboardLoaded = false;
    renderAdminDashboard(null);
    setAdminStatus(adminDashboardStatusEl, text || "Chưa cấu hình DB", "error");
    return;
  }
  const data = JSON.parse(text);
  adminOptions = data.options || adminOptions;
  adminDashboardLoaded = true;
  renderAdminFilters("dashboard");
  renderAdminDashboard(data);
  setAdminStatus(adminDashboardStatusEl, "Sẵn sàng", "ready");
}

function renderAdminDashboard(data) {
  renderAdminMetrics(adminDashboardMetricsEl, data?.summary || null);
  renderAdminTopClasses(data?.topClasses || []);
  renderAdminAttentionInvoices(data?.attentionInvoices || []);
}

function renderAdminMetrics(root, summary) {
  if (!summary) {
    root.textContent = "Chưa có dữ liệu";
    return;
  }
  root.innerHTML = `
    <div><strong>${formatMoney(summary.totalReceivable || 0)}</strong><span>Cần thu</span></div>
    <div><strong>${formatMoney(summary.totalCollected || 0)}</strong><span>Đã thu</span></div>
    <div><strong>${formatMoney(summary.outstandingAmount || 0)}</strong><span>Còn thiếu</span></div>
    <div><strong>${formatPercent(summary.collectionRate || 0)}</strong><span>Tỷ lệ thu</span></div>
    <div><strong>${Number(summary.unpaidStudentCount || 0)}</strong><span>HS unpaid</span></div>
    <div><strong>${Number(summary.partialPaymentCount || 0)}</strong><span>Partial</span></div>
    <div><strong>${Number(summary.overpaidManualReviewCount || 0)}</strong><span>Overpaid/review</span></div>
    <div><strong>${Number(summary.unmatchedTransactionCount || 0) + Number(summary.manualReviewCount || 0)}</strong><span>Giao dịch cần xử lý</span></div>
  `;
}

function renderAdminTopClasses(rows) {
  adminTopClassCountEl.textContent = `${rows.length} lớp`;
  adminTopClassRowsEl.innerHTML = rows
    .map(
      (row) => `
        <tr>
          <td><strong>${escapeHtml(row.className || "-")}</strong><small>${escapeHtml(row.schoolYearCode || "")} · Khối ${escapeHtml(row.grade || "-")}</small></td>
          <td>${Number(row.studentCount || 0)}</td>
          <td class="money">${formatMoney(row.totalAmount || 0)}</td>
          <td class="money">${formatMoney(row.paidAmount || 0)}</td>
          <td class="money">${formatMoney(row.outstandingAmount || 0)}</td>
          <td>${formatPercent(row.collectionRate || 0)}</td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    adminTopClassRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Chưa có dữ liệu lớp</td></tr>`;
  }
}

function renderAdminAttentionInvoices(rows) {
  adminAttentionCountEl.textContent = `${rows.length} hóa đơn`;
  adminAttentionRowsEl.innerHTML = rows
    .map(
      (invoice) => `
        <tr>
          <td><strong>${escapeHtml(invoice.invoiceCode || "")}</strong><small>${escapeHtml(invoice.dueDate || "")}</small></td>
          <td>${escapeHtml(invoice.studentCode || "")}<small>${escapeHtml(invoice.studentName || "")}</small></td>
          <td>${escapeHtml(invoice.className || "")}<small>${escapeHtml(invoice.periodCode || "")}</small></td>
          <td class="money">${formatMoney(invoice.totalAmount || 0)}</td>
          <td class="money">${formatMoney(invoice.outstandingAmount || 0)}</td>
          <td><span class="tag">${escapeHtml(invoice.status || "")}</span></td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    adminAttentionRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Không có hóa đơn cần xử lý</td></tr>`;
  }
}

async function loadAdminReports(force = false) {
  if (!force && adminReportsLoaded) {
    return;
  }
  setAdminStatus(adminReportsStatusEl, "Đang tải", "busy");
  const params = adminFilterParams("reports");
  const query = params.toString();
  const res = await fetch(`/api/v1/admin/reports${query ? `?${query}` : ""}`);
  const text = await res.text();
  if (!res.ok) {
    adminReportsLoaded = false;
    renderAdminReports(null);
    setAdminStatus(adminReportsStatusEl, text || "Chưa cấu hình DB", "error");
    return;
  }
  const data = JSON.parse(text);
  adminOptions = data.options || adminOptions;
  adminReportsLoaded = true;
  renderAdminFilters("reports");
  renderAdminReports(data);
  setAdminStatus(adminReportsStatusEl, "Sẵn sàng", "ready");
}

function renderAdminReports(data) {
  renderAdminMetrics(adminReportsSummaryEl, data?.summary || null);
  renderAdminReportClasses(data?.classRows || []);
  renderAdminReportInvoices(data?.invoiceRows || []);
}

function renderAdminReportClasses(rows) {
  adminReportClassCountEl.textContent = `${rows.length} lớp`;
  adminReportClassRowsEl.innerHTML = rows
    .map(
      (row) => `
        <tr>
          <td><strong>${escapeHtml(row.className || "-")}</strong><small>${escapeHtml(row.schoolYearCode || "")} · Khối ${escapeHtml(row.grade || "-")}</small></td>
          <td>${Number(row.invoiceCount || 0)}<small>${Number(row.studentCount || 0)} học sinh</small></td>
          <td class="money">${formatMoney(row.totalAmount || 0)}</td>
          <td class="money">${formatMoney(row.paidAmount || 0)}</td>
          <td class="money">${formatMoney(row.outstandingAmount || 0)}</td>
          <td><span class="tag">U${Number(row.unpaidCount || 0)} P${Number(row.partialCount || 0)} R${Number(row.manualReviewCount || 0)}</span></td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    adminReportClassRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Chưa có báo cáo theo lớp</td></tr>`;
  }
}

function renderAdminReportInvoices(rows) {
  adminReportInvoiceCountEl.textContent = `${rows.length} hóa đơn`;
  adminReportInvoiceRowsEl.innerHTML = rows
    .map(
      (invoice) => `
        <tr>
          <td><strong>${escapeHtml(invoice.invoiceCode || "")}</strong><small>${escapeHtml(invoice.dueDate || "")}</small></td>
          <td>${escapeHtml(invoice.studentCode || "")}<small>${escapeHtml(invoice.studentName || "")}</small></td>
          <td>${escapeHtml(invoice.className || "")}<small>${escapeHtml(invoice.periodCode || "")}</small></td>
          <td class="money">${formatMoney(invoice.totalAmount || 0)}</td>
          <td class="money">${formatMoney(invoice.paidAmount || 0)}</td>
          <td><span class="tag">${escapeHtml(invoice.status || "")}</span></td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    adminReportInvoiceRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Chưa có hóa đơn trong bộ lọc</td></tr>`;
  }
}

async function loadAdminUsers(force = false) {
  if (!force && adminUsersLoaded) {
    return;
  }
  setAdminStatus(adminUsersStatusEl, "Đang tải", "busy");
  const res = await fetch("/api/v1/admin/users");
  const text = await res.text();
  if (!res.ok) {
    adminUsersLoaded = false;
    adminUsersData = { users: [], roles: [], permissions: [] };
    renderAdminUsers(null);
    setAdminStatus(adminUsersStatusEl, text || "Chưa cấu hình DB", "error");
    return;
  }
  adminUsersData = JSON.parse(text);
  adminUsersLoaded = true;
  renderAdminUsers(adminUsersData);
  setAdminStatus(adminUsersStatusEl, "Sẵn sàng", "ready");
}

function renderAdminUsers(data) {
  const users = data?.users || [];
  const roles = data?.roles || [];
  renderAdminUserRoleSelect(roles);
  renderAdminUserRows(users);
  renderAdminRoleList(roles);
}

function renderAdminUserRoleSelect(roles) {
  const selected = selectedOptionValues(adminUserRolesEl);
  adminUserRolesEl.innerHTML = roles
    .map((role) => `<option value="${escapeAttr(role.code || "")}">${escapeHtml(role.name || role.code || "")}</option>`)
    .join("");
  [...adminUserRolesEl.options].forEach((option) => {
    option.selected = selected.includes(option.value);
  });
}

function renderAdminUserRows(users) {
  adminUserCountEl.textContent = `${users.length} user`;
  adminUserRowsEl.innerHTML = users
    .map((user) => {
      const roleText = (user.roles || []).map((role) => role.code).join(", ");
      return `
        <tr data-admin-user-id="${escapeAttr(user.id || "")}">
          <td><strong>${escapeHtml(user.email || "")}</strong></td>
          <td>${escapeHtml(user.displayName || "")}</td>
          <td><span class="tag">${escapeHtml(user.status || "")}</span></td>
          <td>${escapeHtml(roleText || "-")}</td>
        </tr>
      `;
    })
    .join("");
  if (!users.length) {
    adminUserRowsEl.innerHTML = `<tr><td colspan="4" class="empty-cell">Chưa có user</td></tr>`;
  }
  adminUserRowsEl.querySelectorAll("[data-admin-user-id]").forEach((row) => {
    row.addEventListener("click", () => selectAdminUser(row.dataset.adminUserId));
  });
}

function renderAdminRoleList(roles) {
  adminRoleCountEl.textContent = `${roles.length} role`;
  adminRoleListEl.innerHTML = roles
    .map((role) => {
      const permissions = (role.permissions || []).map((permission) => `<span class="tag">${escapeHtml(permission.code || "")}</span>`).join("");
      return `
        <div class="admin-role-item">
          <div>
            <strong>${escapeHtml(role.name || role.code || "")}</strong>
            <small>${escapeHtml(role.code || "")}${role.isSystem ? " · system" : ""}</small>
            <p>${escapeHtml(role.description || "")}</p>
          </div>
          <div class="admin-permission-list">${permissions || `<span class="tag">no permissions</span>`}</div>
        </div>
      `;
    })
    .join("");
  if (!roles.length) {
    adminRoleListEl.textContent = "Chưa có role";
  }
}

function selectAdminUser(userId) {
  const user = (adminUsersData.users || []).find((item) => item.id === userId);
  if (!user) return;
  adminUserIdEl.value = user.id || "";
  adminUserEmailEl.value = user.email || "";
  adminUserDisplayNameEl.value = user.displayName || "";
  adminUserStatusEl.value = optionValueOrEmpty(adminUserStatusEl, user.status || "active") || "active";
  const roleCodes = (user.roles || []).map((role) => role.code);
  [...adminUserRolesEl.options].forEach((option) => {
    option.selected = roleCodes.includes(option.value);
  });
}

function clearAdminUserForm() {
  adminUserIdEl.value = "";
  adminUserEmailEl.value = "";
  adminUserDisplayNameEl.value = "";
  adminUserStatusEl.value = "active";
  [...adminUserRolesEl.options].forEach((option) => {
    option.selected = false;
  });
}

function selectedOptionValues(selectEl) {
  return [...selectEl.selectedOptions].map((option) => option.value).filter(Boolean);
}

async function saveAdminUser() {
  setAdminStatus(adminUsersStatusEl, "Đang lưu", "busy");
  const res = await fetch("/api/v1/admin/users/save", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-ABC-Admin-Permission": "system.users.write",
    },
    body: JSON.stringify({
      id: adminUserIdEl.value,
      email: adminUserEmailEl.value,
      displayName: adminUserDisplayNameEl.value,
      status: adminUserStatusEl.value,
    }),
  });
  const text = await res.text();
  if (!res.ok) {
    setAdminStatus(adminUsersStatusEl, text || "Không lưu được user", "error");
    return;
  }
  const data = JSON.parse(text);
  if (data.user?.id) {
    adminUserIdEl.value = data.user.id;
  }
  adminUsersLoaded = false;
  await loadAdminUsers(true);
  if (data.user?.id) {
    selectAdminUser(data.user.id);
  }
  setAdminStatus(adminUsersStatusEl, "Đã lưu user", "ready");
}

async function assignAdminUserRoles() {
  if (!adminUserIdEl.value) {
    setAdminStatus(adminUsersStatusEl, "Chọn hoặc lưu user trước", "error");
    return;
  }
  setAdminStatus(adminUsersStatusEl, "Đang lưu roles", "busy");
  const res = await fetch("/api/v1/admin/users/roles", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-ABC-Admin-Permission": "system.users.assign_roles",
    },
    body: JSON.stringify({
      userId: adminUserIdEl.value,
      roleCodes: selectedOptionValues(adminUserRolesEl),
    }),
  });
  const text = await res.text();
  if (!res.ok) {
    setAdminStatus(adminUsersStatusEl, text || "Không lưu được roles", "error");
    return;
  }
  adminUsersLoaded = false;
  await loadAdminUsers(true);
  selectAdminUser(adminUserIdEl.value);
  setAdminStatus(adminUsersStatusEl, "Đã lưu roles", "ready");
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
  const file = masterImportState?.file || masterCsvFileEl.files[0];
  if (!file) {
    setMasterStatus("Chưa chọn file", "error");
    return;
  }
  if (apply && !window.confirm("Áp dụng import sẽ ghi dữ liệu học sinh, phụ huynh và lớp vào database. Tiếp tục?")) {
    return;
  }

  setMasterStatus(apply ? "Đang áp dụng" : "Đang kiểm tra", "busy");
  const form = new FormData();
  form.append("file", file);
  if (masterImportState) {
    form.append("mapping", JSON.stringify(collectImportMapping(masterImportState)));
  }
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
    masterImportSummaryEl.className = "master-import-summary";
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

function setPaymentReconStatus(message, tone = "ready") {
  paymentReconStatusEl.textContent = message;
  paymentReconStatusEl.dataset.tone = tone;
}

async function loadPaymentReconciliation(force = false) {
  if (!force && paymentReconciliationLoaded) {
    return;
  }
  setPaymentReconStatus("Đang tải", "busy");
  const params = new URLSearchParams();
  if (paymentProviderFilterEl.value) {
    params.set("provider", paymentProviderFilterEl.value);
  }
  if (paymentInvoiceStatusFilterEl.value) {
    params.set("invoiceStatus", paymentInvoiceStatusFilterEl.value);
  }
  if (paymentTransactionStatusFilterEl.value) {
    params.set("transactionStatus", paymentTransactionStatusFilterEl.value);
  }
  const query = params.toString();
  const res = await fetch(`/api/v1/payments/reconciliation${query ? `?${query}` : ""}`);
  const text = await res.text();
  if (!res.ok) {
    paymentReconciliationLoaded = false;
    renderPaymentReconciliation(null);
    setPaymentReconStatus(text || "Chưa cấu hình DB", "error");
    return;
  }
  paymentReconciliationData = JSON.parse(text);
  paymentReconciliationLoaded = true;
  renderPaymentReconciliation(paymentReconciliationData);
  setPaymentReconStatus("Sẵn sàng", "ready");
}

function renderPaymentReconciliation(data) {
  const providers = data?.providers || paymentReconciliationData.providers || [];
  renderPaymentProviderFilter(providers);
  renderPaymentReconSummary(data?.summary || null);
  renderPaymentReconInvoices(data?.invoices || [], data?.intents || {});
  renderPaymentReconTransactions(data?.transactions || []);
}

function renderPaymentProviderFilter(providers) {
  const selected = paymentProviderFilterEl.value;
  paymentProviderFilterEl.innerHTML = [
    `<option value="">Tất cả</option>`,
    ...providers.map((provider) => {
      const suffix = provider.configured ? "" : " · thiếu cấu hình";
      return `<option value="${escapeAttr(provider.code)}">${escapeHtml(provider.displayName || provider.code)}${suffix}</option>`;
    }),
  ].join("");
  paymentProviderFilterEl.value = optionValueOrEmpty(paymentProviderFilterEl, selected);
}

function renderPaymentReconSummary(summary) {
  if (!summary) {
    paymentReconSummaryEl.textContent = "Chưa có dữ liệu đối soát";
    return;
  }
  const outstanding = summary.outstandingAmount || 0;
  paymentReconSummaryEl.innerHTML = `
    <div><strong>${formatMoney(summary.totalReceivable || 0)}</strong><span>Cần thu</span></div>
    <div><strong>${formatMoney(summary.totalCollected || 0)}</strong><span>Đã nhận</span></div>
    <div><strong>${formatMoney(outstanding)}</strong><span>Còn thiếu</span></div>
    <div><strong>${Number(summary.partialCount || 0)}</strong><span>Partial</span></div>
    <div><strong>${Number(summary.overpaidCount || 0)}</strong><span>Overpaid</span></div>
    <div><strong>${Number(summary.unmatchedCount || 0) + Number(summary.manualReviewCount || 0)}</strong><span>Cần xử lý</span></div>
  `;
}

function renderPaymentReconInvoices(invoices, intents) {
  paymentReconInvoiceCountEl.textContent = `${invoices.length} hóa đơn`;
  const hasPayOS = (paymentReconciliationData.providers || []).some((provider) => provider.code === "payos");
  paymentReconInvoiceRowsEl.innerHTML = invoices
    .map((invoice) => {
      const paid = Number(invoice.paidAmount || 0);
      const total = Number(invoice.totalAmount || 0);
      const outstanding = Math.max(total - paid, 0);
      const intent = intents?.[invoice.id];
      const intentLabel = intent?.provider ? `${intent.provider}: ${intent.status}` : "";
      return `
        <tr data-recon-invoice-row="${escapeAttr(invoice.id || "")}">
          <td><strong>${escapeHtml(invoice.invoiceCode || "")}</strong>${intentLabel ? `<small>${escapeHtml(intentLabel)}</small>` : ""}</td>
          <td>${escapeHtml(invoice.studentCode || "")} · ${escapeHtml(invoice.studentName || "")}</td>
          <td>${escapeHtml(invoice.className || "")}</td>
          <td class="money">${formatMoney(total)}</td>
          <td class="money">${formatMoney(paid)}</td>
          <td><span class="tag">${escapeHtml(invoice.status || "unpaid")}</span></td>
          <td>
            <div class="invoice-actions">
              <button type="button" data-recon-intent="${escapeAttr(invoice.id || "")}" data-recon-provider="manual_vietqr">${muiIcon("qr_code")}<span>QR</span></button>
              ${hasPayOS ? `<button type="button" data-recon-intent="${escapeAttr(invoice.id || "")}" data-recon-provider="payos">${muiIcon("link")}<span>payOS</span></button>` : ""}
              <button type="button" data-recon-cash="${escapeAttr(invoice.id || "")}" data-recon-default-amount="${escapeAttr(outstanding || total)}">${muiIcon("payments")}<span>Tiền mặt</span></button>
            </div>
          </td>
        </tr>
      `;
    })
    .join("");
  if (!invoices.length) {
    paymentReconInvoiceRowsEl.innerHTML = `<tr><td colspan="7" class="empty-cell">Chưa có hóa đơn để đối soát</td></tr>`;
  }
  paymentReconInvoiceRowsEl.querySelectorAll("[data-recon-intent]").forEach((button) => {
    button.addEventListener("click", () => createPaymentIntent(button.dataset.reconIntent, button.dataset.reconProvider));
  });
  paymentReconInvoiceRowsEl.querySelectorAll("[data-recon-cash]").forEach((button) => {
    button.addEventListener("click", () => recordManualCashReceipt(button.dataset.reconCash, Number(button.dataset.reconDefaultAmount || 0)));
  });
  paymentReconInvoiceRowsEl.querySelectorAll("[data-recon-invoice-row]").forEach((row) => {
    row.addEventListener("click", (event) => {
      if (event.target.closest("button, a")) return;
      const invoice = (paymentReconciliationData.invoices || []).find((item) => item.id === row.dataset.reconInvoiceRow);
      renderPaymentReconDetail(invoiceDetailTemplate(invoice));
    });
  });
}

function renderPaymentReconTransactions(transactions) {
  paymentReconTransactionCountEl.textContent = `${transactions.length} giao dịch`;
  paymentReconTransactionRowsEl.innerHTML = transactions
    .map(
      (transaction) => `
        <tr data-recon-transaction-row="${escapeAttr(transaction.id || "")}">
          <td><strong>${escapeHtml(transaction.provider || "")}</strong><small>${escapeHtml(transaction.providerTransactionId || "")}</small></td>
          <td>${escapeHtml(formatDateTime(transaction.transactionTime))}</td>
          <td class="money">${formatMoney(transaction.amount || 0)}</td>
          <td>${escapeHtml(transaction.accountNumber || "")}</td>
          <td>${escapeHtml(transaction.description || transaction.referenceCode || "")}</td>
          <td><span class="tag">${escapeHtml(transaction.invoiceCode || transaction.status || "unmatched")}</span></td>
        </tr>
      `,
    )
    .join("");
  if (!transactions.length) {
    paymentReconTransactionRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Chưa có giao dịch vào</td></tr>`;
  }
  paymentReconTransactionRowsEl.querySelectorAll("[data-recon-transaction-row]").forEach((row) => {
    row.addEventListener("click", () => {
      const transaction = (paymentReconciliationData.transactions || []).find((item) => item.id === row.dataset.reconTransactionRow);
      renderPaymentReconDetail(transactionDetailTemplate(transaction));
    });
  });
}

async function createPaymentIntent(invoiceId, provider) {
  if (!invoiceId) return;
  setPaymentReconStatus("Đang tạo intent", "busy");
  const res = await fetch("/api/v1/payments/intents", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ invoiceId, provider }),
  });
  const text = await res.text();
  let data = null;
  try {
    data = JSON.parse(text);
  } catch {
    data = null;
  }
  if (!res.ok || !data?.intent) {
    setPaymentReconStatus("Lỗi", "error");
    renderPaymentReconDetail(`<div class="reconciliation-error">${escapeHtml(text || "Không tạo được payment intent")}</div>`);
    return;
  }
  renderPaymentReconDetail(paymentIntentDetailTemplate(data));
  paymentReconciliationLoaded = false;
  await loadPaymentReconciliation(true);
}

async function recordManualCashReceipt(invoiceId, defaultAmount) {
  if (!invoiceId) return;
  const amountValue = window.prompt("Số tiền thu tiền mặt", String(defaultAmount || ""));
  if (amountValue === null) return;
  const amount = parseMoneyInput(amountValue);
  const collectorName = window.prompt("Người thu tiền", "");
  if (collectorName === null) return;
  const receiptReference = window.prompt("Mã phiếu thu", `CASH${Date.now()}`);
  if (receiptReference === null) return;
  if (!window.confirm("Ghi nhận khoản thu tiền mặt vào ledger đối soát?")) {
    return;
  }
  setPaymentReconStatus("Đang ghi nhận", "busy");
  const res = await fetch("/api/v1/payments/cash-receipts", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ invoiceId, amount, collectorName, receiptReference }),
  });
  const text = await res.text();
  let data = null;
  try {
    data = JSON.parse(text);
  } catch {
    data = null;
  }
  if (!res.ok || !data) {
    setPaymentReconStatus("Lỗi", "error");
    renderPaymentReconDetail(`<div class="reconciliation-error">${escapeHtml(text || "Không ghi nhận được tiền mặt")}</div>`);
    return;
  }
  paymentReconciliationLoaded = false;
  await loadPaymentReconciliation(true);
  renderPaymentReconDetail(transactionDetailTemplate(data.transaction));
  setPaymentReconStatus("Đã ghi nhận", "ready");
}

function renderPaymentReconDetail(html) {
  paymentReconDetailEl.innerHTML = html || "Chưa chọn hóa đơn hoặc giao dịch";
}

function invoiceDetailTemplate(invoice) {
  if (!invoice) return "Không tìm thấy hóa đơn";
  return `
    <div class="reconciliation-detail-grid">
      <span>Mã hóa đơn</span><strong>${escapeHtml(invoice.invoiceCode || "")}</strong>
      <span>Học sinh</span><strong>${escapeHtml(invoice.studentCode || "")} · ${escapeHtml(invoice.studentName || "")}</strong>
      <span>Lớp / kỳ</span><strong>${escapeHtml(invoice.className || "")} · ${escapeHtml(invoice.periodCode || "")}</strong>
      <span>Phải thu</span><strong>${formatMoney(invoice.totalAmount || 0)}</strong>
      <span>Đã thu</span><strong>${formatMoney(invoice.paidAmount || 0)}</strong>
      <span>Status</span><strong>${escapeHtml(invoice.status || "")}</strong>
    </div>
  `;
}

function transactionDetailTemplate(transaction) {
  if (!transaction) return "Không tìm thấy giao dịch";
  return `
    <div class="reconciliation-detail-grid">
      <span>Provider</span><strong>${escapeHtml(transaction.provider || "")}</strong>
      <span>Reference</span><strong>${escapeHtml(transaction.providerTransactionId || transaction.referenceCode || "")}</strong>
      <span>Hóa đơn</span><strong>${escapeHtml(transaction.invoiceCode || "Chưa match")}</strong>
      <span>Số tiền</span><strong>${formatMoney(transaction.amount || 0)}</strong>
      <span>Thời gian</span><strong>${escapeHtml(formatDateTime(transaction.transactionTime))}</strong>
      <span>Nội dung</span><strong>${escapeHtml(transaction.description || "")}</strong>
      <span>Status</span><strong>${escapeHtml(transaction.status || "")}</strong>
    </div>
  `;
}

function paymentIntentDetailTemplate(data) {
  const intent = data.intent || {};
  const qr = data.qr || {};
  const link = intent.paymentUrl ? `<a class="button-link" href="${escapeAttr(intent.paymentUrl)}" target="_blank" rel="noreferrer">${muiIcon("open_in_new")}<span>Mở link</span></a>` : "";
  const image = qr.qrData ? `<img class="reconciliation-qr" src="${escapeAttr(qr.qrData)}" alt="QR thanh toán" />` : "";
  return `
    <div class="reconciliation-intent-detail">
      ${image}
      <div class="reconciliation-detail-grid">
        <span>Provider</span><strong>${escapeHtml(intent.provider || "")}</strong>
        <span>Intent</span><strong>${escapeHtml(intent.intentCode || "")}</strong>
        <span>Reference</span><strong>${escapeHtml(intent.providerReference || "")}</strong>
        <span>Số tiền</span><strong>${formatMoney(intent.amount || 0)}</strong>
        <span>Status</span><strong>${escapeHtml(intent.status || "")}</strong>
      </div>
      ${link}
      ${qr.vietqr ? `<textarea class="payload" readonly>${escapeHtml(qr.vietqr)}</textarea>` : ""}
    </div>
  `;
}

async function loadNotifications(force = false) {
  if (notificationLoaded && !force) return;
  setNotificationStatus("Đang tải", "busy");
  const [optionsRes, logsRes] = await Promise.all([
    fetch("/api/v1/notifications/options"),
    fetch("/api/v1/notifications/logs?limit=50"),
  ]);
  const optionsText = await optionsRes.text();
  if (!optionsRes.ok) {
    setNotificationStatus(optionsText || "Không tải được notification", "error");
    renderNotificationControls();
    renderNotificationPreview(null);
    renderNotificationCampaigns([]);
    renderNotificationLogs([]);
    return;
  }
  notificationOptions = JSON.parse(optionsText);
  renderNotificationControls();
  renderNotificationCampaigns(notificationOptions.campaigns || []);
  if (logsRes.ok) {
    const logsData = await logsRes.json();
    renderNotificationLogs(logsData.logs || []);
  } else {
    renderNotificationLogs([]);
  }
  notificationLoaded = true;
  setNotificationStatus("Sẵn sàng", "ready");
}

function renderNotificationControls() {
  const templates = notificationOptions.templates || [];
  const selectedTemplate = notificationTemplateEl.value;
  notificationTemplateEl.innerHTML = templates.length
    ? templates
        .map((item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.name || item.code)} v${Number(item.version || 1)}</option>`)
        .join("")
    : `<option value="">Chưa có template</option>`;
  if (templates.some((item) => item.id === selectedTemplate)) {
    notificationTemplateEl.value = selectedTemplate;
  } else {
    const type = notificationCampaignTypeEl.value || "first_notice";
    const match = templates.find((item) => item.code === type) || templates[0];
    if (match) notificationTemplateEl.value = match.id;
  }

  const currentYear = notificationSchoolYearEl.value;
  notificationSchoolYearEl.innerHTML = `<option value="">Tất cả</option>${(notificationOptions.schoolYears || [])
    .map((item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.code || item.name || "")}</option>`)
    .join("")}`;
  if ((notificationOptions.schoolYears || []).some((item) => item.id === currentYear)) {
    notificationSchoolYearEl.value = currentYear;
  }
  renderNotificationGradeOptions();
  renderNotificationClassOptions();
}

function renderNotificationGradeOptions() {
  const current = notificationGradeEl.value;
  const yearId = notificationSchoolYearEl.value;
  const grades = [
    ...new Set(
      (notificationOptions.classes || [])
        .filter((item) => !yearId || item.schoolYearId === yearId)
        .map((item) => item.grade)
        .filter(Boolean),
    ),
  ].sort((a, b) => String(a).localeCompare(String(b), "vi", { numeric: true }));
  notificationGradeEl.innerHTML = `<option value="">Tất cả</option>${grades
    .map((grade) => `<option value="${escapeAttr(grade)}">${escapeHtml(grade)}</option>`)
    .join("")}`;
  if (grades.includes(current)) {
    notificationGradeEl.value = current;
  }
}

function renderNotificationClassOptions() {
  const current = notificationClassEl.value;
  const yearId = notificationSchoolYearEl.value;
  const grade = notificationGradeEl.value;
  const classes = (notificationOptions.classes || []).filter((item) => {
    return (!yearId || item.schoolYearId === yearId) && (!grade || item.grade === grade);
  });
  notificationClassEl.innerHTML = `<option value="">Tất cả</option>${classes
    .map((item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.name)} · Khối ${escapeHtml(item.grade || "-")}</option>`)
    .join("")}`;
  if (classes.some((item) => item.id === current)) {
    notificationClassEl.value = current;
  }
}

function collectNotificationInput() {
  return {
    campaignId: currentNotificationCampaignId,
    name: notificationCampaignNameEl.value.trim(),
    campaignType: notificationCampaignTypeEl.value || "first_notice",
    templateId: notificationTemplateEl.value || "",
    schoolYearId: notificationSchoolYearEl.value || "",
    classId: notificationClassEl.value || "",
    grade: notificationGradeEl.value || "",
    periodCode: notificationPeriodEl.value.trim(),
    invoiceStatus: notificationInvoiceStatusEl.value || "",
    dueOnOrBefore: notificationDueBeforeEl.value || "",
  };
}

async function previewNotifications() {
  setNotificationStatus("Đang preview", "busy");
  const res = await fetch("/api/v1/notifications/campaigns/preview", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(collectNotificationInput()),
  });
  const text = await res.text();
  if (!res.ok) {
    setNotificationStatus(text || "Preview failed", "error");
    return null;
  }
  const data = JSON.parse(text);
  notificationPreviewData = data;
  renderNotificationPreview(data);
  if (data.campaign?.id) currentNotificationCampaignId = data.campaign.id;
  setNotificationStatus((data.issues || []).length ? "Cần xử lý" : "Preview xong", (data.issues || []).length ? "error" : "ready");
  return data;
}

async function saveNotificationCampaign() {
  setNotificationStatus("Đang lưu", "busy");
  const res = await fetch("/api/v1/notifications/campaigns/save", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(collectNotificationInput()),
  });
  const text = await res.text();
  if (!res.ok) {
    setNotificationStatus(text || "Không lưu được campaign", "error");
    return null;
  }
  const data = JSON.parse(text);
  currentNotificationCampaignId = data.campaign?.id || currentNotificationCampaignId;
  notificationPreviewData = { ...notificationPreviewData, recipients: data.recipients || [], summary: data.summary || {}, campaign: data.campaign };
  renderNotificationPreview(notificationPreviewData);
  renderNotificationCampaigns(data.campaigns || notificationOptions.campaigns || []);
  notificationLoaded = false;
  setNotificationStatus("Đã lưu", "ready");
  return data.campaign;
}

async function sendNotificationCampaign() {
  if (!window.confirm("Gửi campaign sẽ gửi email thật qua provider hiện tại và ghi log theo từng invoice/recipient. Tiếp tục?")) {
    return;
  }
  setNotificationStatus("Đang gửi", "busy");
  const input = { ...collectNotificationInput(), confirmSend: true, dryRun: false };
  const res = await fetch("/api/v1/notifications/campaigns/send", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
  const text = await res.text();
  if (!res.ok) {
    setNotificationStatus(text || "Không gửi được campaign", "error");
    return;
  }
  const data = JSON.parse(text);
  currentNotificationCampaignId = data.campaign?.id || currentNotificationCampaignId;
  renderNotificationResults(data);
  notificationLoaded = false;
  setNotificationStatus("Đã xử lý gửi", "ready");
}

function renderNotificationPreview(data) {
  const issues = data?.issues || [];
  const recipients = data?.recipients || [];
  const summary = data?.summary || {};
  notificationRecipientCountEl.textContent = `${recipients.length} email`;
  if (!data) {
    notificationSummaryEl.textContent = "Chưa có preview";
    notificationRecipientsEl.innerHTML = `<tr><td colspan="6">Chưa có recipient</td></tr>`;
    return;
  }
  if (issues.length) {
    notificationSummaryEl.innerHTML = issues.map((item) => `<div><strong>${escapeHtml(item.type)}</strong><span>${escapeHtml(item.message)}</span></div>`).join("");
  } else {
    notificationSummaryEl.innerHTML = `
      <div><strong>${Number(summary.invoiceCount || 0)}</strong><span>invoice</span></div>
      <div><strong>${Number(summary.recipientCount || 0)}</strong><span>recipient</span></div>
      <div><strong>${formatMoney(summary.totalAmount || 0)}</strong><span>phải thu</span></div>
      <div><strong>${formatMoney(summary.unpaidAmount || 0)}</strong><span>còn phải thu</span></div>
      <div><strong>${Number(summary.alreadySent || 0)}</strong><span>đã gửi trước</span></div>
    `;
  }
  notificationRecipientsEl.innerHTML = recipients.length
    ? recipients.map(notificationRecipientRowTemplate).join("")
    : `<tr><td colspan="6">Không có recipient phù hợp</td></tr>`;
}

function notificationRecipientRowTemplate(item) {
  const status = item.alreadySent ? "already sent" : item.status || item.invoiceStatus || "pending";
  return `
    <tr>
      <td><strong>${escapeHtml(item.invoiceCode || "")}</strong><small>${escapeHtml(item.periodCode || "")}</small></td>
      <td>${escapeHtml(item.studentCode || "")}<small>${escapeHtml(item.studentName || "")}</small></td>
      <td>${escapeHtml(item.className || "")}<small>${escapeHtml(item.dueDate || "")}</small></td>
      <td>${escapeHtml(item.recipientEmail || "")}<small>${escapeHtml(item.recipientName || "")}</small></td>
      <td>${formatMoney(item.amount || 0)}<small>Đã thu ${formatMoney(item.paidAmount || 0)}</small></td>
      <td><span class="status-pill">${escapeHtml(status)}</span></td>
    </tr>
  `;
}

function renderNotificationCampaigns(campaigns) {
  notificationOptions.campaigns = campaigns || [];
  notificationCampaignCountEl.textContent = `${notificationOptions.campaigns.length} campaign`;
  notificationCampaignRowsEl.innerHTML = notificationOptions.campaigns.length
    ? notificationOptions.campaigns.map(notificationCampaignRowTemplate).join("")
    : `<tr><td colspan="5">Chưa có campaign</td></tr>`;
  notificationCampaignRowsEl.querySelectorAll("tr[data-campaign-id]").forEach((row) => {
    row.addEventListener("click", () => {
      const campaign = notificationOptions.campaigns.find((item) => item.id === row.dataset.campaignId);
      if (campaign) selectNotificationCampaign(campaign);
    });
  });
}

function notificationCampaignRowTemplate(item) {
  return `
    <tr data-campaign-id="${escapeAttr(item.id || "")}">
      <td><strong>${escapeHtml(item.name || item.code || "")}</strong><small>${escapeHtml(item.campaignType || "")}</small></td>
      <td>${escapeHtml(item.periodCode || "-")}<small>${escapeHtml(item.className || item.grade || item.schoolYearCode || "Tất cả")}</small></td>
      <td>${escapeHtml(item.template || "")}<small>v${Number(item.templateVersion || 1)}</small></td>
      <td>${Number(item.recipientCount || 0)}<small>sent ${Number(item.sentCount || 0)} · errors ${Number(item.errorCount || 0)}</small></td>
      <td><span class="status-pill">${escapeHtml(item.status || "")}</span></td>
    </tr>
  `;
}

async function selectNotificationCampaign(campaign) {
  currentNotificationCampaignId = campaign.id || "";
  notificationCampaignNameEl.value = campaign.name || "";
  notificationCampaignTypeEl.value = campaign.campaignType || "first_notice";
  notificationSchoolYearEl.value = campaign.schoolYearId || "";
  renderNotificationGradeOptions();
  notificationGradeEl.value = campaign.grade || "";
  renderNotificationClassOptions();
  notificationClassEl.value = campaign.classId || "";
  notificationPeriodEl.value = campaign.periodCode || "";
  notificationInvoiceStatusEl.value = campaign.invoiceStatus || "";
  notificationDueBeforeEl.value = campaign.dueOnOrBefore || "";
  if ((notificationOptions.templates || []).some((item) => item.id === campaign.templateId)) {
    notificationTemplateEl.value = campaign.templateId;
  }
  await loadNotificationLogs(campaign.id);
}

async function loadNotificationLogs(campaignId = "") {
  const params = new URLSearchParams({ limit: "50" });
  if (campaignId) params.set("campaignId", campaignId);
  const res = await fetch(`/api/v1/notifications/logs?${params.toString()}`);
  if (!res.ok) {
    renderNotificationLogs([]);
    return;
  }
  const data = await res.json();
  renderNotificationLogs(data.logs || []);
}

function renderNotificationLogs(logs) {
  notificationLogCountEl.textContent = `${logs.length} log`;
  notificationLogsEl.innerHTML = logs.length
    ? logs.map(notificationLogTemplate).join("")
    : "Chưa có log gửi";
}

function notificationLogTemplate(item) {
  const message = item.error || item.providerMessageId || item.status || "";
  return `
    <div class="notification-log-item">
      <strong>${escapeHtml(item.status || "")} · ${escapeHtml(item.invoiceCode || "")}</strong>
      <span>${escapeHtml(item.recipientEmail || "")}</span>
      <small>${escapeHtml(item.campaignName || "")} · ${escapeHtml(item.templateCode || "")} v${Number(item.templateVersion || 1)}</small>
      <small>${escapeHtml(formatDateTime(item.sentAt))} · ${escapeHtml(item.provider || (item.dryRun ? "dry-run" : "-"))}</small>
      <small>${escapeHtml(message)}</small>
    </div>
  `;
}

function renderNotificationResults(data) {
  const results = data.results || [];
  const logs = data.logs || [];
  notificationPreviewData = {
    ...notificationPreviewData,
    campaign: data.campaign,
    summary: data.summary || {},
    recipients: (notificationPreviewData.recipients || []).map((recipient) => {
      const result = results.find((item) => item.id === recipient.id || item.email === recipient.recipientEmail);
      return result ? { ...recipient, status: result.status, lastError: result.error || "" } : recipient;
    }),
  };
  renderNotificationPreview(notificationPreviewData);
  renderNotificationLogs(logs);
  if (data.campaign) {
    const campaigns = [data.campaign, ...(notificationOptions.campaigns || []).filter((item) => item.id !== data.campaign.id)];
    renderNotificationCampaigns(campaigns);
  }
}

function setNotificationStatus(message, tone = "ready") {
  notificationStatusEl.textContent = message;
  notificationStatusEl.dataset.tone = tone;
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

function formatPercent(value) {
  return `${Math.round(Number(value || 0) * 1000) / 10}%`;
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
  await prepareImportMapping("payments", file);
});
applyPaymentImportBtn.addEventListener("click", submitPaymentImport);
cancelPaymentImportBtn.addEventListener("click", clearPaymentImport);

tabButtons.forEach((button) => {
  button.addEventListener("click", () => activateTab(button.dataset.tabTarget));
});

refreshAdminDashboardBtn.addEventListener("click", () => loadAdminDashboard(true));
adminDashboardYearEl.addEventListener("change", async () => {
  renderAdminFilters("dashboard");
  await loadAdminDashboard(true);
});
adminDashboardGradeEl.addEventListener("change", async () => {
  renderAdminFilters("dashboard");
  await loadAdminDashboard(true);
});
adminDashboardClassEl.addEventListener("change", () => loadAdminDashboard(true));
adminDashboardPeriodEl.addEventListener("change", () => loadAdminDashboard(true));
adminDashboardMonthEl.addEventListener("change", () => loadAdminDashboard(true));
adminDashboardInvoiceStatusEl.addEventListener("change", () => loadAdminDashboard(true));

refreshAdminReportsBtn.addEventListener("click", () => loadAdminReports(true));
adminReportsYearEl.addEventListener("change", async () => {
  renderAdminFilters("reports");
  await loadAdminReports(true);
});
adminReportsGradeEl.addEventListener("change", async () => {
  renderAdminFilters("reports");
  await loadAdminReports(true);
});
adminReportsClassEl.addEventListener("change", () => loadAdminReports(true));
adminReportsPeriodEl.addEventListener("change", () => loadAdminReports(true));
adminReportsMonthEl.addEventListener("change", () => loadAdminReports(true));
adminReportsInvoiceStatusEl.addEventListener("change", () => loadAdminReports(true));

refreshAdminUsersBtn.addEventListener("click", () => loadAdminUsers(true));
clearAdminUserBtn.addEventListener("click", clearAdminUserForm);
saveAdminUserBtn.addEventListener("click", saveAdminUser);
assignAdminUserRolesBtn.addEventListener("click", assignAdminUserRoles);

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
masterCsvFileEl.addEventListener("change", async () => {
  const file = masterCsvFileEl.files[0];
  if (!file) {
    clearMasterImport();
    return;
  }
  masterImportSummaryEl.textContent = file.name;
  await prepareImportMapping("master_data", file);
});
clearMasterImportMappingBtn.addEventListener("click", clearMasterImport);
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
refreshPaymentReconBtn.addEventListener("click", () => loadPaymentReconciliation(true));
paymentProviderFilterEl.addEventListener("change", () => loadPaymentReconciliation(true));
paymentInvoiceStatusFilterEl.addEventListener("change", () => loadPaymentReconciliation(true));
paymentTransactionStatusFilterEl.addEventListener("change", () => loadPaymentReconciliation(true));
refreshNotificationsBtn.addEventListener("click", () => loadNotifications(true));
notificationCampaignTypeEl.addEventListener("change", () => {
  currentNotificationCampaignId = "";
  const match = (notificationOptions.templates || []).find((item) => item.code === notificationCampaignTypeEl.value);
  if (match) notificationTemplateEl.value = match.id;
});
notificationSchoolYearEl.addEventListener("change", () => {
  renderNotificationGradeOptions();
  renderNotificationClassOptions();
});
notificationGradeEl.addEventListener("change", renderNotificationClassOptions);
previewNotificationsBtn.addEventListener("click", previewNotifications);
saveNotificationCampaignBtn.addEventListener("click", saveNotificationCampaign);
sendNotificationCampaignBtn.addEventListener("click", sendNotificationCampaign);

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
