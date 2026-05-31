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
const csvFileButtonEl = document.querySelector('label[for="csvFile"]');
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
const masterSchoolFilterEl = document.querySelector("#masterSchoolFilter");
const masterSchoolYearFilterEl = document.querySelector("#masterSchoolYearFilter");
const masterGradeFilterEl = document.querySelector("#masterGradeFilter");
const masterClassFilterEl = document.querySelector("#masterClassFilter");
const masterSearchEl = document.querySelector("#masterSearch");
const refreshMasterDataBtn = document.querySelector("#refreshMasterData");
const openSchoolTreeSchoolDialogBtn = document.querySelector("#openSchoolTreeSchoolDialog");
const openSchoolTreeYearDialogBtn = document.querySelector("#openSchoolTreeYearDialog");
const openSchoolTreeClassDialogBtn = document.querySelector("#openSchoolTreeClassDialog");
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
const masterStudentDetailEl = document.querySelector("#masterStudentDetail");
const masterStudentEditorPanelEl = document.querySelector(".master-editor-panel");
const masterStudentEditorContentEl = document.querySelector("#masterStudentEditorContent");
const masterStudentIdEl = document.querySelector("#masterStudentId");
const masterStudentCodeEl = document.querySelector("#masterStudentCode");
const masterStudentNameEl = document.querySelector("#masterStudentName");
const masterStudentClassEl = document.querySelector("#masterStudentClass");
const masterStudentStatusEl = document.querySelector("#masterStudentStatus");
const masterParentEditorRowsEl = document.querySelector("#masterParentEditorRows");
const newMasterStudentBtn = document.querySelector("#newMasterStudent");
const editMasterStudentBtn = document.querySelector("#editMasterStudent");
const addMasterParentBtn = document.querySelector("#addMasterParent");
const saveMasterStudentBtn = document.querySelector("#saveMasterStudent");
const masterConflictPanelEl = document.querySelector("#masterConflictPanel");
const masterConflictCountEl = document.querySelector("#masterConflictCount");
const masterConflictListEl = document.querySelector("#masterConflictList");
const masterImportPanelEl = document.querySelector(".master-import-panel");
const schoolTreeCountEl = document.querySelector("#schoolTreeCount");
const schoolTreeListEl = document.querySelector("#schoolTreeList");
const schoolTreeDetailEl = document.querySelector("#schoolTreeDetail");
const schoolTreeSchoolEditorEl = document.querySelector("#schoolTreeSchoolEditor");
const schoolTreeSchoolIdEl = document.querySelector("#schoolTreeSchoolId");
const schoolTreeSchoolCodeEl = document.querySelector("#schoolTreeSchoolCode");
const schoolTreeSchoolNameEl = document.querySelector("#schoolTreeSchoolName");
const schoolTreeSchoolStatusEl = document.querySelector("#schoolTreeSchoolStatus");
const saveSchoolTreeSchoolBtn = document.querySelector("#saveSchoolTreeSchool");
const newSchoolTreeSchoolBtn = document.querySelector("#newSchoolTreeSchool");
const schoolTreeYearEditorEl = document.querySelector("#schoolTreeYearEditor");
const schoolTreeYearIdEl = document.querySelector("#schoolTreeYearId");
const schoolTreeYearSchoolEl = document.querySelector("#schoolTreeYearSchool");
const schoolTreeYearCodeEl = document.querySelector("#schoolTreeYearCode");
const schoolTreeYearNameEl = document.querySelector("#schoolTreeYearName");
const schoolTreeYearStatusEl = document.querySelector("#schoolTreeYearStatus");
const saveSchoolTreeYearBtn = document.querySelector("#saveSchoolTreeYear");
const newSchoolTreeYearBtn = document.querySelector("#newSchoolTreeYear");
const schoolTreeClassEditorEl = document.querySelector("#schoolTreeClassEditor");
const schoolTreeClassIdEl = document.querySelector("#schoolTreeClassId");
const schoolTreeClassYearEl = document.querySelector("#schoolTreeClassYear");
const schoolTreeClassGradeEl = document.querySelector("#schoolTreeClassGrade");
const schoolTreeClassNameEl = document.querySelector("#schoolTreeClassName");
const schoolTreeClassStatusEl = document.querySelector("#schoolTreeClassStatus");
const saveSchoolTreeClassBtn = document.querySelector("#saveSchoolTreeClass");
const newSchoolTreeClassBtn = document.querySelector("#newSchoolTreeClass");
const feeScheduleLoadStatusEl = document.querySelector("#feeScheduleLoadStatus");
const refreshFeeSchedulesBtn = document.querySelector("#refreshFeeSchedules");
const openFeeScheduleDialogBtn = document.querySelector("#openFeeScheduleDialog");
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
const feeScheduleOperatorEl = document.querySelector("#feeScheduleOperator");
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
const openInvoiceDialogBtn = document.querySelector("#openInvoiceDialog");
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
const invoiceDetailSummaryEl = document.querySelector("#invoiceDetailSummary");
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
const openNotificationDialogBtn = document.querySelector("#openNotificationDialog");
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
const exportAdminReportClassesBtn = document.querySelector("#exportAdminReportClasses");
const exportAdminReportInvoicesBtn = document.querySelector("#exportAdminReportInvoices");
const exportAdminReportTransactionsBtn = document.querySelector("#exportAdminReportTransactions");
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
const operationsStatusEl = document.querySelector("#operationsStatus");
const refreshOperationsBtn = document.querySelector("#refreshOperations");
const operationSourceFilterEl = document.querySelector("#operationSourceFilter");
const operationLevelFilterEl = document.querySelector("#operationLevelFilter");
const operationLimitEl = document.querySelector("#operationLimit");
const operationLogCountEl = document.querySelector("#operationLogCount");
const operationLogRowsEl = document.querySelector("#operationLogRows");
const auditLogCountEl = document.querySelector("#auditLogCount");
const auditLogRowsEl = document.querySelector("#auditLogRows");
const adminUsersStatusEl = document.querySelector("#adminUsersStatus");
const refreshAdminUsersBtn = document.querySelector("#refreshAdminUsers");
const newAdminUserBtn = document.querySelector("#newAdminUser");
const adminUserIdEl = document.querySelector("#adminUserId");
const adminUserEmailEl = document.querySelector("#adminUserEmail");
const adminUserPhoneEl = document.querySelector("#adminUserPhone");
const adminUserDisplayNameEl = document.querySelector("#adminUserDisplayName");
const adminUserStatusEl = document.querySelector("#adminUserStatus");
const adminUserPasswordEl = document.querySelector("#adminUserPassword");
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
const currentSectionKickerEl = document.querySelector("#currentSectionKicker");
const currentSectionTitleEl = document.querySelector("#currentSectionTitle");
const currentSectionDescriptionEl = document.querySelector("#currentSectionDescription");
const loginScreenEl = document.querySelector("#loginScreen");
const appShellEl = document.querySelector("#appShell");
const loginFormEl = document.querySelector("#loginForm");
const loginEmailEl = document.querySelector("#loginEmail");
const loginPasswordEl = document.querySelector("#loginPassword");
const loginSubmitBtn = document.querySelector("#loginSubmit");
const loginStatusEl = document.querySelector("#loginStatus");
const bootstrapFormEl = document.querySelector("#bootstrapForm");
const bootstrapDisplayNameEl = document.querySelector("#bootstrapDisplayName");
const bootstrapEmailEl = document.querySelector("#bootstrapEmail");
const bootstrapPhoneEl = document.querySelector("#bootstrapPhone");
const bootstrapPasswordEl = document.querySelector("#bootstrapPassword");
const bootstrapPasswordConfirmEl = document.querySelector("#bootstrapPasswordConfirm");
const bootstrapSubmitBtn = document.querySelector("#bootstrapSubmit");
const bootstrapStatusEl = document.querySelector("#bootstrapStatus");
const authUserBadgeEl = document.querySelector("#authUserBadge");
const logoutButton = document.querySelector("#logoutButton");
const appDialogEl = document.querySelector("#appDialog");
const appDialogKickerEl = document.querySelector("#appDialogKicker");
const appDialogTitleEl = document.querySelector("#appDialogTitle");
const appDialogBodyEl = document.querySelector("#appDialogBody");
const appDialogErrorEl = document.querySelector("#appDialogError");
const appDialogActionsEl = document.querySelector("#appDialogActions");
const appDialogCloseBtn = document.querySelector("#appDialogClose");
const openEmailConfigDialogBtn = document.querySelector("#openEmailConfigDialog");
const openCronConfigDialogBtn = document.querySelector("#openCronConfigDialog");

let banks = [];
let currentItems = [];
let selectedId = "";
let feeColumnCollapsed = false;
let savedEmailConfig = {};
let masterDataOptions = { schools: [], schoolYears: [], classes: [] };
let masterDataLoaded = false;
let masterStudentsData = [];
let selectedMasterStudentKey = "";
let masterStudentParentDrafts = [];
let schoolTreeData = { schools: [] };
let selectedSchoolTreeNode = { type: "", id: "" };
let paymentImportState = null;
let masterImportState = null;
let feeScheduleOptions = { feeTypes: [], schoolYears: [], classes: [] };
let feeSchedulesLoaded = false;
let invoiceOptions = { schedules: [], schoolYears: [], classes: [] };
let invoicesLoaded = false;
let invoicesData = [];
let selectedInvoiceId = "";
let paymentReconciliationLoaded = false;
let paymentReconciliationData = { providers: [], invoices: [], transactions: [], intents: {}, summary: {} };
let paymentReconSelection = { type: "", id: "" };
let notificationLoaded = false;
let notificationOptions = { templates: [], campaigns: [], schoolYears: [], classes: [] };
let notificationPreviewData = { recipients: [], summary: {}, campaign: null, logs: [] };
let currentNotificationCampaignId = "";
let adminOptions = { schoolYears: [], classes: [] };
let adminDashboardLoaded = false;
let adminReportsLoaded = false;
let operationsLoaded = false;
let adminUsersLoaded = false;
let adminUsersData = { users: [], roles: [], permissions: [] };
let authSession = null;
let refreshAuthPromise = null;
let activeDialogRestore = null;
let activeDialogOnClose = null;

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

const tabMetadata = {
  dashboardTab: {
    kicker: "Tổng quan",
    title: "Dashboard thu học phí",
    description: "Theo dõi thu học phí, công nợ và các hóa đơn cần xử lý.",
  },
  masterDataTab: {
    kicker: "Trường & học sinh",
    title: "Học sinh, phụ huynh, lớp",
    description: "Quản lý dữ liệu nền cho học sinh, phụ huynh, lớp và năm học.",
  },
  feeTemplateTab: {
    kicker: "Học phí",
    title: "Bảng phí theo kỳ",
    description: "Thiết lập biểu phí, phụ phí và preview trước khi sinh hóa đơn.",
  },
  invoiceTab: {
    kicker: "Học phí",
    title: "Hóa đơn",
    description: "Sinh, kiểm tra và xuất hóa đơn/PDF receipt từ bảng phí đã lưu.",
  },
  reconciliationTab: {
    kicker: "Thanh toán",
    title: "Đối soát thanh toán",
    description: "Theo dõi intent, giao dịch, tiền mặt và trạng thái đối soát hóa đơn.",
  },
  paymentsTab: {
    kicker: "Thanh toán",
    title: "Thanh toán VietQR",
    description: "Import batch thanh toán legacy, sinh QR và kiểm tra payload thanh toán.",
  },
  notificationTab: {
    kicker: "Liên lạc",
    title: "Thông báo học phí",
    description: "Tạo campaign, preview người nhận và theo dõi log gửi thông báo.",
  },
  emailTab: {
    kicker: "Liên lạc",
    title: "Email & Cron",
    description: "Cấu hình provider email, preview/dry-run và quản lý lịch gửi cục bộ.",
  },
  reportsTab: {
    kicker: "Quản trị",
    title: "Báo cáo công nợ",
    description: "Xem và export báo cáo lớp, hóa đơn và giao dịch thanh toán.",
  },
  operationsTab: {
    kicker: "Quản trị",
    title: "Vận hành",
    description: "Kiểm tra operational logs, audit logs và các lỗi nền cần xử lý.",
  },
  usersTab: {
    kicker: "Quản trị",
    title: "Người dùng và quyền",
    description: "Quản lý user, role và permission trước khi bật enforcement đầy đủ.",
  },
};

const tabAccess = {
  dashboardTab: ["dashboard.view"],
  masterDataTab: ["student.view"],
  feeTemplateTab: ["fee.view"],
  invoiceTab: ["invoice.view"],
  reconciliationTab: ["payment.view"],
  paymentsTab: ["payment.create"],
  notificationTab: ["notification.view"],
  emailTab: { any: ["email_config.view", "email_config.update", "notification.send", "email_cron.view", "email_cron.update"] },
  reportsTab: ["report.view"],
  operationsTab: ["operation_log.view"],
  usersTab: ["user.view"],
};

const permissionAliases = {
  "dashboard.view": ["admin.dashboard.read"],
  "student.view": ["master_data.read"],
  "student.create": ["master_data.write"],
  "student.update": ["master_data.write"],
  "school_tree.view": ["master_data.read"],
  "school_tree.update": ["master_data.write"],
  "fee.view": ["fee_schedules.read"],
  "fee.update": ["fee_schedules.write"],
  "invoice.view": ["invoices.read"],
  "invoice.create": ["invoices.write"],
  "payment.view": ["payments.read"],
  "payment.create": ["payments.write"],
  "payment.reconcile": ["payments.reconcile"],
  "notification.view": ["notifications.read"],
  "notification.create": ["notifications.write"],
  "notification.send": ["notifications.send", "email.send"],
  "email_config.view": ["email.config.read"],
  "email_config.update": ["email.config.write"],
  "email_cron.view": ["email.cron.manage"],
  "email_cron.update": ["email.cron.manage"],
  "report.view": ["admin.reports.read"],
  "report.export": ["admin.reports.export"],
  "operation_log.view": ["operations.read"],
  "audit_log.view": ["audit.read"],
  "user.view": ["system.users.read"],
  "user.create": ["system.users.write"],
  "user.update": ["system.users.write"],
  "user.assign_role": ["system.users.assign_roles"],
};

const nativeFetch = window.fetch.bind(window);
window.fetch = authAwareFetch;

init();

function muiIcon(name) {
  return `<span class="mui-icon" aria-hidden="true">${escapeHtml(name)}</span>`;
}

function openAppDialog({ title, kicker = "Dialog", icon = "", nodes = [], content = null, actions = [], size = "md", onClose = null } = {}) {
  if (!appDialogEl) return;
  if (appDialogEl.open) {
    const previousOnClose = activeDialogOnClose;
    activeDialogOnClose = null;
    appDialogEl.close("replace");
    restoreDialogContent();
    if (previousOnClose) previousOnClose();
  } else {
    restoreDialogContent();
  }
  appDialogEl.className = `app-dialog app-dialog-${size}`;
  appDialogKickerEl.textContent = kicker;
  appDialogTitleEl.innerHTML = `${icon ? muiIcon(icon) : ""}<span>${escapeHtml(title || "")}</span>`;
  appDialogBodyEl.innerHTML = "";
  appDialogActionsEl.innerHTML = "";
  appDialogErrorEl.hidden = true;
  appDialogErrorEl.textContent = "";

  const records = [];
  nodes.filter(Boolean).forEach((node) => {
    const placeholder = document.createComment(`dialog:${node.id || node.className || "node"}`);
    const parent = node.parentNode;
    const nextSibling = node.nextSibling;
    parent.insertBefore(placeholder, node);
    node.hidden = false;
    appDialogBodyEl.appendChild(node);
    records.push({ node, parent, nextSibling, placeholder });
  });
  if (content instanceof HTMLElement) {
    appDialogBodyEl.appendChild(content);
  } else if (typeof content === "string" && content) {
    appDialogBodyEl.innerHTML += content;
  }

  activeDialogRestore = () => {
    records.reverse().forEach(({ node, parent, nextSibling, placeholder }) => {
      node.hidden = true;
      if (nextSibling && nextSibling.parentNode === parent) {
        parent.insertBefore(node, nextSibling);
      } else {
        parent.appendChild(node);
      }
      placeholder.remove();
    });
    appDialogBodyEl.innerHTML = "";
    appDialogActionsEl.innerHTML = "";
  };
  activeDialogOnClose = onClose;

  actions.forEach((action) => appDialogActionsEl.appendChild(dialogActionButton(action)));
  if (!actions.length) {
    appDialogActionsEl.appendChild(dialogActionButton({ label: "Đóng", icon: "close", onClick: closeAppDialog }));
  }
  appDialogEl.showModal();
  window.setTimeout(() => {
    const first = appDialogBodyEl.querySelector("input:not([type='hidden']), select, textarea, button:not(:disabled)") || appDialogCloseBtn;
    first?.focus();
  }, 0);
}

function dialogActionButton({ label, icon = "", variant = "", onClick = null, closeOnSuccess = false } = {}) {
  const button = document.createElement("button");
  button.type = "button";
  button.className = variant;
  button.innerHTML = `${icon ? muiIcon(icon) : ""}<span>${escapeHtml(label || "OK")}</span>`;
  button.addEventListener("click", async () => {
    clearDialogError();
    const previous = button.innerHTML;
    button.disabled = true;
    if (onClick) {
      try {
        const result = await onClick();
        if (result === false) {
          button.disabled = false;
          button.innerHTML = previous;
          return;
        }
      } catch (error) {
        showDialogError(error?.message || "Không xử lý được thao tác");
        button.disabled = false;
        button.innerHTML = previous;
        return;
      }
    }
    button.disabled = false;
    button.innerHTML = previous;
    if (closeOnSuccess) closeAppDialog();
  });
  return button;
}

function closeAppDialog() {
  if (appDialogEl?.open) {
    appDialogEl.close("close");
  }
}

function restoreDialogContent() {
  if (activeDialogRestore) {
    activeDialogRestore();
    activeDialogRestore = null;
  }
}

function showDialogError(message) {
  appDialogErrorEl.textContent = message;
  appDialogErrorEl.hidden = false;
}

function clearDialogError() {
  appDialogErrorEl.hidden = true;
  appDialogErrorEl.textContent = "";
}

function confirmDialog({ title, message, confirmLabel = "Xác nhận", confirmIcon = "check", danger = false } = {}) {
  return new Promise((resolve) => {
    let settled = false;
    openAppDialog({
      title,
      kicker: "Confirm",
      icon: danger ? "warning" : "help",
      size: "sm",
      content: `<div class="dialog-message">${escapeHtml(message || "")}</div>`,
      onClose: () => {
        if (!settled) resolve(false);
      },
      actions: [
        {
          label: "Hủy",
          icon: "close",
          onClick: () => {
            settled = true;
            resolve(false);
          },
          closeOnSuccess: true,
        },
        {
          label: confirmLabel,
          icon: confirmIcon,
          variant: danger ? "danger" : "primary",
          onClick: () => {
            settled = true;
            resolve(true);
          },
          closeOnSuccess: true,
        },
      ],
    });
  });
}

async function init() {
  showLogin("Đang kiểm tra phiên");
  const bootstrapStatus = await loadAuthBootstrapStatus();
  if (!bootstrapStatus) {
    return;
  }
  if (bootstrapStatus.needsBootstrap) {
    showBootstrap("Tạo Admin để bắt đầu sử dụng hệ thống");
    return;
  }
  const session = await loadCurrentAuthSession();
  if (!session) {
    showLogin(loginStatusEl.textContent || "Vui lòng đăng nhập");
    return;
  }
  showApp(session);
  await initializeAppData();
}

async function initializeAppData() {
  status("Đang tải", "busy");
  await loadBanks();
  if (hasPermission("email_config.view")) {
    await loadEmailConfig();
  }
  if (hasPermission("email_cron.view") || hasPermission("email_cron.update")) {
    await loadEmailCron();
  }
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
  renderOperations(null);
  renderAdminUsers(null);
  renderFeeTemplate(defaultPaymentItems);
  renderRows(hasPermission("payment.create") ? sampleRows : []);
  if (hasPermission("dashboard.view")) {
    await loadAdminDashboard();
  }
  if (hasPermission("student.view")) {
    await loadMasterData();
  }
  if (hasPermission("payment.create")) {
    await generate();
  } else {
    status("Không đủ quyền", "error");
  }
  if (hasPermission("notification.send") && hasPermission("payment.create")) {
    await previewEmail();
  }
  await loadActiveTabData(tabPanels.find((panel) => panel.classList.contains("active"))?.id || "");
}

async function loadCurrentAuthSession() {
  const sessionRes = await nativeFetch("/api/v1/auth/session");
  if (sessionRes.ok) {
    return sessionRes.json();
  }
  if (sessionRes.status === 401 && (await refreshAuthSession())) {
    const retryRes = await nativeFetch("/api/v1/auth/session");
    if (retryRes.ok) {
      return retryRes.json();
    }
  }
  const text = await sessionRes.text();
  if (text && sessionRes.status !== 401) {
    setLoginStatus(text, "error");
  }
  return null;
}

async function loadAuthBootstrapStatus() {
  const res = await nativeFetch("/api/v1/auth/bootstrap");
  const text = await res.text();
  if (!res.ok) {
    showLogin(text || "Không kiểm tra được trạng thái Admin");
    setLoginStatus(text || "Không kiểm tra được trạng thái Admin", "error");
    return null;
  }
  return JSON.parse(text);
}

async function authAwareFetch(input, options = {}) {
  const path = requestPath(input);
  const res = await nativeFetch(input, options);
  if (res.status !== 401 || path.startsWith("/api/v1/auth/")) {
    return res;
  }
  if (!(await refreshAuthSession())) {
    showLogin("Phiên đăng nhập đã hết hạn");
    return res;
  }
  return nativeFetch(input, options);
}

function requestPath(input) {
  const value = typeof input === "string" ? input : input?.url || "";
  try {
    return new URL(value, window.location.origin).pathname;
  } catch {
    return value;
  }
}

async function refreshAuthSession() {
  if (refreshAuthPromise) return refreshAuthPromise;
  refreshAuthPromise = nativeFetch("/api/v1/auth/refresh", { method: "POST" })
    .then(async (res) => {
      if (!res.ok) return false;
      authSession = await res.json();
      updateAuthBadge(authSession);
      return true;
    })
    .catch(() => false)
    .finally(() => {
      refreshAuthPromise = null;
    });
  return refreshAuthPromise;
}

function showLogin(message = "") {
  authSession = null;
  loginScreenEl.hidden = false;
  appShellEl.hidden = true;
  loginFormEl.hidden = false;
  bootstrapFormEl.hidden = true;
  if (message) setLoginStatus(message, message.includes("hết hạn") ? "error" : "");
  loginSubmitBtn.disabled = false;
}

function showBootstrap(message = "") {
  authSession = null;
  loginScreenEl.hidden = false;
  appShellEl.hidden = true;
  loginFormEl.hidden = true;
  bootstrapFormEl.hidden = false;
  setBootstrapStatus(message, "");
  bootstrapSubmitBtn.disabled = false;
}

function showApp(session) {
  authSession = session;
  loginScreenEl.hidden = true;
  appShellEl.hidden = false;
  loginPasswordEl.value = "";
  setLoginStatus("", "");
  updateAuthBadge(session);
  applyPermissionUI();
  activateInitialAllowedTab();
}

function updateAuthBadge(session) {
  const user = session?.user || {};
  authUserBadgeEl.textContent = user.displayName || user.email || user.phone || "Đã đăng nhập";
}

function currentPermissionSet() {
  return new Set((authSession?.user?.permissions || []).map((permission) => permission.code).filter(Boolean));
}

function hasPermission(permission) {
  if (!permission) return true;
  const permissions = currentPermissionSet();
  if (permissions.has(permission)) return true;
  return (permissionAliases[permission] || []).some((alias) => permissions.has(alias));
}

function hasAnyPermission(permissions) {
  return (permissions || []).some((permission) => hasPermission(permission));
}

function canUseTab(tabId) {
  const access = tabAccess[tabId];
  if (!access) return true;
  if (Array.isArray(access)) return access.every((permission) => hasPermission(permission));
  if (access.any) return hasAnyPermission(access.any);
  return true;
}

function setElementAllowed(el, allowed) {
  if (el) el.hidden = !allowed;
}

function applyPermissionUI() {
  tabButtons.forEach((button) => {
    button.hidden = !canUseTab(button.dataset.tabTarget);
  });
  tabPanels.forEach((panel) => {
    if (!canUseTab(panel.id)) {
      panel.hidden = true;
      panel.classList.remove("active");
    }
  });

  setElementAllowed(csvFileButtonEl, hasPermission("payment.create"));
  setElementAllowed(loadSampleBtn, hasPermission("payment.create"));
  setElementAllowed(generateBtn, hasPermission("payment.create"));
  setElementAllowed(addRowBtn, hasPermission("payment.create"));
  setElementAllowed(toggleFeeColumnBtn, hasPermission("payment.create"));
  setElementAllowed(masterStudentEditorPanelEl, hasPermission("student.update"));
  setElementAllowed(masterImportPanelEl, hasPermission("student.create"));
  setElementAllowed(openSchoolTreeSchoolDialogBtn, hasPermission("school_tree.update"));
  setElementAllowed(openSchoolTreeYearDialogBtn, hasPermission("school_tree.update"));
  setElementAllowed(openSchoolTreeClassDialogBtn, hasPermission("school_tree.update"));
  setElementAllowed(saveSchoolTreeSchoolBtn, hasPermission("school_tree.update"));
  setElementAllowed(newSchoolTreeSchoolBtn, hasPermission("school_tree.update"));
  setElementAllowed(saveSchoolTreeYearBtn, hasPermission("school_tree.update"));
  setElementAllowed(newSchoolTreeYearBtn, hasPermission("school_tree.update"));
  setElementAllowed(saveSchoolTreeClassBtn, hasPermission("school_tree.update"));
  setElementAllowed(newSchoolTreeClassBtn, hasPermission("school_tree.update"));
  setElementAllowed(openFeeScheduleDialogBtn, hasPermission("fee.view") || hasPermission("fee.update"));
  setElementAllowed(previewFeeScheduleBtn, hasPermission("fee.view"));
  setElementAllowed(saveFeeScheduleBtn, hasPermission("fee.update"));
  setElementAllowed(openInvoiceDialogBtn, hasPermission("invoice.view") || hasPermission("invoice.create"));
  setElementAllowed(previewInvoicesBtn, hasPermission("invoice.view"));
  setElementAllowed(generateInvoicesBtn, hasPermission("invoice.create"));
  setElementAllowed(openNotificationDialogBtn, hasPermission("notification.create") || hasPermission("notification.send"));
  setElementAllowed(previewNotificationsBtn, hasPermission("notification.send"));
  setElementAllowed(saveNotificationCampaignBtn, hasPermission("notification.create"));
  setElementAllowed(sendNotificationCampaignBtn, hasPermission("notification.send"));
  setElementAllowed(exportAdminReportClassesBtn, hasPermission("report.export"));
  setElementAllowed(exportAdminReportInvoicesBtn, hasPermission("report.export"));
  setElementAllowed(exportAdminReportTransactionsBtn, hasPermission("report.export"));
  setElementAllowed(newAdminUserBtn, hasPermission("user.create") || hasPermission("user.update"));
  setElementAllowed(saveAdminUserBtn, hasPermission("user.create") || hasPermission("user.update"));
  setElementAllowed(assignAdminUserRolesBtn, hasPermission("user.assign_role"));
  setElementAllowed(openEmailConfigDialogBtn, hasPermission("email_config.view") || hasPermission("email_config.update"));
  setElementAllowed(saveEmailConfigBtn, hasPermission("email_config.update"));
  setElementAllowed(previewEmailBtn, hasPermission("notification.send"));
  setElementAllowed(dryRunEmailBtn, hasPermission("notification.send"));
  setElementAllowed(sendEmailBtn, hasPermission("notification.send"));
  setElementAllowed(openCronConfigDialogBtn, hasPermission("email_cron.view") || hasPermission("email_cron.update"));
  setElementAllowed(saveCronBtn, hasPermission("email_cron.update"));
  setElementAllowed(disableCronBtn, hasPermission("email_cron.update"));
  setElementAllowed(runCronNowBtn, hasPermission("email_cron.update"));
}

function activateInitialAllowedTab() {
  const current = tabPanels.find((panel) => panel.classList.contains("active"))?.id || "dashboardTab";
  const targetId = canUseTab(current) ? current : tabButtons.find((button) => !button.hidden)?.dataset.tabTarget || "";
  if (!targetId) {
    currentSectionKickerEl.textContent = "Không đủ quyền";
    currentSectionTitleEl.textContent = "Chưa có màn hình được cấp quyền";
    currentSectionDescriptionEl.textContent = "Liên hệ quản trị viên để cập nhật role.";
    return;
  }
  updateCurrentSection(targetId);
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
}

function setLoginStatus(message, tone = "") {
  loginStatusEl.textContent = message;
  loginStatusEl.dataset.tone = tone;
}

function setBootstrapStatus(message, tone = "") {
  bootstrapStatusEl.textContent = message;
  bootstrapStatusEl.dataset.tone = tone;
}

async function submitLogin(event) {
  event.preventDefault();
  loginSubmitBtn.disabled = true;
  setLoginStatus("Đang đăng nhập", "busy");
  const res = await nativeFetch("/api/v1/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      identifier: loginEmailEl.value.trim(),
      password: loginPasswordEl.value,
    }),
  });
  const text = await res.text();
  if (!res.ok) {
    loginSubmitBtn.disabled = false;
    setLoginStatus(text || "Không đăng nhập được", "error");
    return;
  }
  const session = JSON.parse(text);
  showApp(session);
  await initializeAppData();
}

async function submitBootstrap(event) {
  event.preventDefault();
  const password = bootstrapPasswordEl.value;
  if (password !== bootstrapPasswordConfirmEl.value) {
    setBootstrapStatus("Password xác nhận không khớp", "error");
    return;
  }
  if (!bootstrapEmailEl.value.trim() && !bootstrapPhoneEl.value.trim()) {
    setBootstrapStatus("Nhập Email hoặc SĐT", "error");
    return;
  }
  bootstrapSubmitBtn.disabled = true;
  setBootstrapStatus("Đang tạo Admin", "busy");
  const res = await nativeFetch("/api/v1/auth/bootstrap", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      displayName: bootstrapDisplayNameEl.value.trim(),
      email: bootstrapEmailEl.value.trim(),
      phone: bootstrapPhoneEl.value.trim(),
      password,
    }),
  });
  const text = await res.text();
  if (!res.ok) {
    bootstrapSubmitBtn.disabled = false;
    setBootstrapStatus(text || "Không tạo được Admin", "error");
    return;
  }
  const session = JSON.parse(text);
  bootstrapPasswordEl.value = "";
  bootstrapPasswordConfirmEl.value = "";
  showApp(session);
  await initializeAppData();
}

async function logout() {
  await nativeFetch("/api/v1/auth/logout", { method: "POST" });
  showLogin("Đã đăng xuất");
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
  if (!canUseTab(targetId)) {
    return;
  }
  updateCurrentSection(targetId);
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
  await loadActiveTabData(targetId);
}

async function loadActiveTabData(targetId) {
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
  if (targetId === "operationsTab") {
    await loadOperations();
  }
  if (targetId === "usersTab") {
    await loadAdminUsers();
  }
  if (targetId === "emailTab") {
    if (hasPermission("notification.send")) {
      await previewEmail();
    }
  }
}

function updateCurrentSection(targetId) {
  const metadata = tabMetadata[targetId] || tabMetadata.dashboardTab;
  currentSectionKickerEl.textContent = metadata.kicker;
  currentSectionTitleEl.textContent = metadata.title;
  currentSectionDescriptionEl.textContent = metadata.description;
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
  if (!hasPermission("payment.create")) {
    status("Không đủ quyền", "error");
    return;
  }
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
    await loadSchoolTree();
    await loadMasterStudents();
  }
}

async function loadMasterDataOptions() {
  setMasterStatus("Đang tải", "busy");
  const res = await fetch("/api/v1/master-data/options");
  const text = await res.text();
  if (!res.ok) {
    masterDataLoaded = false;
    masterDataOptions = { schools: [], schoolYears: [], classes: [] };
    schoolTreeData = { schools: [] };
    renderMasterFilters();
    renderSchoolTree();
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
  const selectedSchool = masterSchoolFilterEl.value;
  const selectedYear = masterSchoolYearFilterEl.value;
  const selectedGrade = masterGradeFilterEl.value;
  const selectedClass = masterClassFilterEl.value;

  masterSchoolFilterEl.innerHTML = [
    `<option value="">Tất cả trường</option>`,
    ...(masterDataOptions.schools || []).map(
      (item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.code)}</option>`,
    ),
  ].join("");
  masterSchoolFilterEl.value = optionValueOrEmpty(masterSchoolFilterEl, selectedSchool);

  masterSchoolYearFilterEl.innerHTML = [
    `<option value="">Tất cả năm học</option>`,
    ...(masterDataOptions.schoolYears || [])
      .filter((item) => !masterSchoolFilterEl.value || item.schoolId === masterSchoolFilterEl.value)
      .map((item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.schoolCode || "")} · ${escapeHtml(item.code)}</option>`),
  ].join("");
  masterSchoolYearFilterEl.value = optionValueOrEmpty(masterSchoolYearFilterEl, selectedYear);

  const grades = [...new Set(
    (masterDataOptions.classes || [])
      .filter((item) => !masterSchoolFilterEl.value || item.schoolId === masterSchoolFilterEl.value)
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
    if (masterSchoolFilterEl.value && item.schoolId !== masterSchoolFilterEl.value) return false;
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
  renderMasterStudentClassSelect();
}

function renderMasterStudentClassSelect(selectedClass = masterStudentClassEl.value) {
  const classes = [...(masterDataOptions.classes || [])].sort((a, b) => {
    const left = [a.schoolCode, a.schoolYearCode, a.grade, a.name].filter(Boolean).join(" ");
    const right = [b.schoolCode, b.schoolYearCode, b.grade, b.name].filter(Boolean).join(" ");
    return left.localeCompare(right, "vi", { numeric: true });
  });
  masterStudentClassEl.innerHTML = [
    `<option value="">Chọn lớp</option>`,
    ...classes.map((item) => {
      const label = [item.schoolCode, item.schoolYearCode, item.grade ? `Khối ${item.grade}` : "", item.name].filter(Boolean).join(" · ");
      return `<option value="${escapeAttr(item.id)}">${escapeHtml(label)}</option>`;
    }),
  ].join("");
  masterStudentClassEl.value = optionValueOrEmpty(masterStudentClassEl, selectedClass || masterClassFilterEl.value);
}

function optionValueOrEmpty(selectEl, value) {
  return [...selectEl.options].some((option) => option.value === value) ? value : "";
}

async function loadSchoolTree() {
  const res = await fetch("/api/v1/school-tree");
  const text = await res.text();
  if (!res.ok) {
    schoolTreeData = { schools: [] };
    renderSchoolTree();
    setMasterStatus(text || "Không tải được cây trường", "error");
    return false;
  }
  schoolTreeData = JSON.parse(text);
  renderSchoolTree();
  return true;
}

function renderSchoolTree() {
  const schools = schoolTreeData.schools || [];
  renderSchoolTreeSelects();
  const nodeCount = countSchoolTreeNodes(schools);
  schoolTreeCountEl.textContent = `${nodeCount} node`;
  if (!schools.length) {
    schoolTreeListEl.innerHTML = `<div class="school-tree-empty">${muiIcon("account_tree")}<span>Chưa có cây trường</span></div>`;
    renderSchoolTreeDetail(null);
    return;
  }
  schoolTreeListEl.innerHTML = schools.map((school) => renderSchoolTreeSchool(school)).join("");
  schoolTreeListEl.querySelectorAll("[data-tree-type]").forEach((button) => {
    button.addEventListener("click", () => selectSchoolTreeNode(button.dataset.treeType, button.dataset.treeId));
  });
  renderSchoolTreeDetail(findSchoolTreeNode(selectedSchoolTreeNode.type, selectedSchoolTreeNode.id));
}

function renderSchoolTreeSchool(school) {
  const active = selectedSchoolTreeNode.type === "school" && selectedSchoolTreeNode.id === school.id;
  return `
    <div class="school-tree-school">
      <button class="school-tree-node ${active ? "is-selected" : ""}" type="button" data-tree-type="school" data-tree-id="${escapeAttr(school.id)}">
        ${muiIcon("apartment")}
        <span><strong>${escapeHtml(school.name || school.code || "-")}</strong><small>${escapeHtml(school.code || "")} · ${Number(school.studentCount || 0)} HS · ${Number(school.feeScheduleCount || 0)} bảng phí</small></span>
      </button>
      <div class="school-tree-children">
        ${(school.schoolYears || []).map((year) => renderSchoolTreeYear(year)).join("")}
      </div>
    </div>
  `;
}

function renderSchoolTreeYear(year) {
  const active = selectedSchoolTreeNode.type === "year" && selectedSchoolTreeNode.id === year.id;
  return `
    <div class="school-tree-year">
      <button class="school-tree-node ${active ? "is-selected" : ""}" type="button" data-tree-type="year" data-tree-id="${escapeAttr(year.id)}">
        ${muiIcon("event")}
        <span><strong>${escapeHtml(year.code || "-")}</strong><small>${Number(year.classCount || 0)} lớp · ${Number(year.studentCount || 0)} HS · ${Number(year.adjustmentCount || 0)} điều chỉnh</small></span>
      </button>
      <div class="school-tree-children">
        ${(year.grades || []).map((grade) => renderSchoolTreeGrade(year, grade)).join("")}
      </div>
    </div>
  `;
}

function renderSchoolTreeGrade(year, grade) {
  const gradeId = `${year.id}|${grade.grade || ""}`;
  const active = selectedSchoolTreeNode.type === "grade" && selectedSchoolTreeNode.id === gradeId;
  return `
    <div class="school-tree-grade">
      <button class="school-tree-node ${active ? "is-selected" : ""}" type="button" data-tree-type="grade" data-tree-id="${escapeAttr(gradeId)}">
        ${muiIcon("stacked_line_chart")}
        <span><strong>Khối ${escapeHtml(grade.grade || "-")}</strong><small>${Number(grade.classCount || 0)} lớp · ${Number(grade.studentCount || 0)} HS</small></span>
      </button>
      <div class="school-tree-children">
        ${(grade.classes || []).map((item) => renderSchoolTreeClass(item)).join("")}
      </div>
    </div>
  `;
}

function renderSchoolTreeClass(item) {
  const active = selectedSchoolTreeNode.type === "class" && selectedSchoolTreeNode.id === item.id;
  return `
    <button class="school-tree-node school-tree-class ${active ? "is-selected" : ""}" type="button" data-tree-type="class" data-tree-id="${escapeAttr(item.id)}">
      ${muiIcon("school")}
      <span><strong>${escapeHtml(item.name || "-")}</strong><small>${Number(item.studentCount || 0)} HS · ${Number(item.activeFeeScheduleCount || 0)} active · ${Number(item.adjustmentCount || 0)} điều chỉnh</small></span>
    </button>
  `;
}

function countSchoolTreeNodes(schools) {
  return schools.reduce((sum, school) => {
    const years = school.schoolYears || [];
    return sum + 1 + years.reduce((yearSum, year) => {
      const grades = year.grades || [];
      return yearSum + 1 + grades.reduce((gradeSum, grade) => gradeSum + 1 + (grade.classes || []).length, 0);
    }, 0);
  }, 0);
}

function renderSchoolTreeSelects() {
  const schools = masterDataOptions.schools || [];
  const years = masterDataOptions.schoolYears || [];
  const selectedYearSchool = schoolTreeYearSchoolEl.value;
  const selectedClassYear = schoolTreeClassYearEl.value;
  schoolTreeYearSchoolEl.innerHTML = [
    `<option value="">Chọn trường</option>`,
    ...schools.map((item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.code)} · ${escapeHtml(item.name || "")}</option>`),
  ].join("");
  schoolTreeYearSchoolEl.value = optionValueOrEmpty(schoolTreeYearSchoolEl, selectedYearSchool);

  schoolTreeClassYearEl.innerHTML = [
    `<option value="">Chọn năm học</option>`,
    ...years
      .filter((item) => !schoolTreeYearSchoolEl.value || item.schoolId === schoolTreeYearSchoolEl.value)
      .map((item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.schoolCode || "")} · ${escapeHtml(item.code || "")}</option>`),
  ].join("");
  schoolTreeClassYearEl.value = optionValueOrEmpty(schoolTreeClassYearEl, selectedClassYear);
}

function findSchoolTreeNode(type, id) {
  if (!type || !id) return null;
  for (const school of schoolTreeData.schools || []) {
    if (type === "school" && school.id === id) return { type, school };
    for (const year of school.schoolYears || []) {
      if (type === "year" && year.id === id) return { type, school, year };
      for (const grade of year.grades || []) {
        const gradeId = `${year.id}|${grade.grade || ""}`;
        if (type === "grade" && gradeId === id) return { type, school, year, grade };
        for (const item of grade.classes || []) {
          if (type === "class" && item.id === id) return { type, school, year, grade, classItem: item };
        }
      }
    }
  }
  return null;
}

async function selectSchoolTreeNode(type, id) {
  selectedSchoolTreeNode = { type, id };
  const node = findSchoolTreeNode(type, id);
  fillSchoolTreeForms(node);
  applySchoolTreeFilters(node);
  renderSchoolTree();
  await loadMasterStudents();
}

function applySchoolTreeFilters(node) {
  if (!node) return;
  masterSchoolFilterEl.value = optionValueOrEmpty(masterSchoolFilterEl, node.school?.id || "");
  renderMasterFilters();
  if (node.year) {
    masterSchoolYearFilterEl.value = optionValueOrEmpty(masterSchoolYearFilterEl, node.year.id || "");
  }
  renderMasterFilters();
  if (node.grade) {
    masterGradeFilterEl.value = optionValueOrEmpty(masterGradeFilterEl, node.grade.grade || "");
  }
  renderMasterClassFilter();
  if (node.classItem) {
    masterClassFilterEl.value = optionValueOrEmpty(masterClassFilterEl, node.classItem.id || "");
  }
}

function fillSchoolTreeForms(node) {
  if (!node) return;
  if (node.school) {
    schoolTreeSchoolIdEl.value = node.school.id || "";
    schoolTreeSchoolCodeEl.value = node.school.code || "";
    schoolTreeSchoolNameEl.value = node.school.name || "";
    schoolTreeSchoolStatusEl.value = optionValueOrEmpty(schoolTreeSchoolStatusEl, node.school.status || "active");
    schoolTreeYearSchoolEl.value = optionValueOrEmpty(schoolTreeYearSchoolEl, node.school.id || "");
  }
  if (node.year) {
    schoolTreeYearIdEl.value = node.year.id || "";
    schoolTreeYearSchoolEl.value = optionValueOrEmpty(schoolTreeYearSchoolEl, node.year.schoolId || node.school?.id || "");
    schoolTreeYearCodeEl.value = node.year.code || "";
    schoolTreeYearNameEl.value = node.year.name || "";
    schoolTreeYearStatusEl.value = optionValueOrEmpty(schoolTreeYearStatusEl, node.year.status || "active");
    renderSchoolTreeSelects();
    schoolTreeClassYearEl.value = optionValueOrEmpty(schoolTreeClassYearEl, node.year.id || "");
  }
  if (node.grade) {
    schoolTreeClassGradeEl.value = node.grade.grade || "";
  }
  if (node.classItem) {
    schoolTreeClassIdEl.value = node.classItem.id || "";
    schoolTreeClassYearEl.value = optionValueOrEmpty(schoolTreeClassYearEl, node.classItem.schoolYearId || "");
    schoolTreeClassGradeEl.value = node.classItem.grade || "";
    schoolTreeClassNameEl.value = node.classItem.name || "";
    schoolTreeClassStatusEl.value = optionValueOrEmpty(schoolTreeClassStatusEl, node.classItem.status || "active");
  }
}

function renderSchoolTreeDetail(node) {
  if (!node) {
    schoolTreeDetailEl.innerHTML = `<div class="detail-placeholder">${muiIcon("account_tree")}<span>Chọn trường, năm học, khối hoặc lớp.</span></div>`;
    return;
  }
  const subject = node.classItem || node.grade || node.year || node.school;
  const title = node.classItem?.name || (node.grade ? `Khối ${node.grade.grade}` : node.year?.code || node.school?.name || "-");
  const lines = [
    ["Trường", node.school?.code || "-"],
    ["Năm học", node.year?.code || "-"],
    ["Khối", node.grade?.grade || node.classItem?.grade || "-"],
    ["Sĩ số", subject?.studentCount ?? "-"],
    ["Bảng phí", subject?.feeScheduleCount ?? "-"],
    ["Điều chỉnh", subject?.adjustmentCount ?? "-"],
  ];
  const latest = node.classItem?.latestFeeScheduleId
    ? `<div class="detail-section"><p class="detail-section-title">Bảng phí gần nhất</p><div class="detail-grid"><span>Kỳ</span><strong>${escapeHtml(node.classItem.latestPeriodCode || "-")}</strong><span>Tên</span><strong>${escapeHtml(node.classItem.latestFeeScheduleName || "-")}</strong><span>Trạng thái</span><strong>${escapeHtml(node.classItem.latestScheduleStatus || "-")}</strong></div></div>`
    : "";
  const openFees = hasPermission("fee.view") && (node.year || node.grade || node.classItem)
    ? `<button data-tree-open-fees="true" type="button"><span class="mui-icon" aria-hidden="true">format_list_bulleted</span><span>Mở bảng phí</span></button>`
    : "";
  schoolTreeDetailEl.innerHTML = `
    <div class="detail-hero">
      ${muiIcon(node.classItem ? "school" : node.grade ? "stacked_line_chart" : node.year ? "event" : "apartment")}
      <div>
        <strong>${escapeHtml(title)}</strong>
        <span>${escapeHtml(node.type)}</span>
      </div>
    </div>
    <div class="detail-grid">
      ${lines.map(([label, value]) => `<span>${escapeHtml(label)}</span><strong>${escapeHtml(String(value))}</strong>`).join("")}
    </div>
    ${latest}
    <div class="detail-actions">${openFees}</div>
  `;
  const openButton = schoolTreeDetailEl.querySelector("[data-tree-open-fees]");
  if (openButton) {
    openButton.addEventListener("click", () => openSchoolTreeFeeScope(node));
  }
}

async function openSchoolTreeFeeScope(node) {
  if (!hasPermission("fee.view")) {
    setMasterStatus("Không đủ quyền mở bảng phí", "error");
    return;
  }
  await activateTab("feeTemplateTab");
  await loadFeeSchedules(true);
  feeScheduleYearEl.value = optionValueOrEmpty(feeScheduleYearEl, node.year?.id || node.classItem?.schoolYearId || "");
  renderFeeScheduleControls();
  feeScheduleGradeEl.value = optionValueOrEmpty(feeScheduleGradeEl, node.grade?.grade || node.classItem?.grade || "");
  renderFeeScheduleClassFilter();
  feeScheduleClassEl.value = optionValueOrEmpty(feeScheduleClassEl, node.classItem?.id || "");
  await loadFeeScheduleList();
}

function clearSchoolTreeSchoolForm() {
  schoolTreeSchoolIdEl.value = "";
  schoolTreeSchoolCodeEl.value = "";
  schoolTreeSchoolNameEl.value = "";
  schoolTreeSchoolStatusEl.value = "active";
}

function clearSchoolTreeYearForm() {
  schoolTreeYearIdEl.value = "";
  schoolTreeYearCodeEl.value = "";
  schoolTreeYearNameEl.value = "";
  schoolTreeYearStatusEl.value = "active";
}

function clearSchoolTreeClassForm() {
  schoolTreeClassIdEl.value = "";
  schoolTreeClassGradeEl.value = "";
  schoolTreeClassNameEl.value = "";
  schoolTreeClassStatusEl.value = "active";
}

function openSchoolTreeEntityDialog(kind) {
  const config = {
    school: {
      title: schoolTreeSchoolIdEl.value ? "Sửa trường" : "Tạo trường",
      icon: "apartment",
      node: schoolTreeSchoolEditorEl,
      save: saveSchoolTreeSchool,
    },
    year: {
      title: schoolTreeYearIdEl.value ? "Sửa năm học" : "Tạo năm học",
      icon: "event",
      node: schoolTreeYearEditorEl,
      save: saveSchoolTreeYear,
    },
    class: {
      title: schoolTreeClassIdEl.value ? "Sửa lớp" : "Tạo lớp",
      icon: "school",
      node: schoolTreeClassEditorEl,
      save: saveSchoolTreeClass,
    },
  }[kind];
  if (!config?.node) return;
  openAppDialog({
    title: config.title,
    kicker: "School tree",
    icon: config.icon,
    nodes: [config.node],
    size: "md",
    actions: [
      { label: "Đóng", icon: "close", onClick: closeAppDialog },
      { label: "Lưu", icon: "save", variant: "primary", onClick: config.save, closeOnSuccess: true },
    ],
  });
}

async function saveSchoolTreeSchool() {
  return saveSchoolTreeEntity("/api/v1/school-tree/schools/save", {
    id: schoolTreeSchoolIdEl.value,
    code: schoolTreeSchoolCodeEl.value,
    name: schoolTreeSchoolNameEl.value,
    status: schoolTreeSchoolStatusEl.value,
  }, "Đã lưu trường");
}

async function saveSchoolTreeYear() {
  return saveSchoolTreeEntity("/api/v1/school-tree/school-years/save", {
    id: schoolTreeYearIdEl.value,
    schoolId: schoolTreeYearSchoolEl.value,
    code: schoolTreeYearCodeEl.value,
    name: schoolTreeYearNameEl.value,
    status: schoolTreeYearStatusEl.value,
  }, "Đã lưu năm học");
}

async function saveSchoolTreeClass() {
  return saveSchoolTreeEntity("/api/v1/school-tree/classes/save", {
    id: schoolTreeClassIdEl.value,
    schoolYearId: schoolTreeClassYearEl.value,
    grade: schoolTreeClassGradeEl.value,
    name: schoolTreeClassNameEl.value,
    status: schoolTreeClassStatusEl.value,
  }, "Đã lưu lớp");
}

async function saveSchoolTreeEntity(url, payload, successMessage) {
  if (!hasPermission("school_tree.update")) {
    setMasterStatus("Không đủ quyền", "error");
    return false;
  }
  setMasterStatus("Đang lưu cây trường", "busy");
  const res = await fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  const text = await res.text();
  if (!res.ok) {
    setMasterStatus(text || "Không lưu được cây trường", "error");
    return false;
  }
  masterDataLoaded = false;
  feeSchedulesLoaded = false;
  invoiceOptions = { schedules: [], schoolYears: [], classes: [] };
  await loadMasterData(true);
  setMasterStatus(successMessage, "ready");
  return true;
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

function exportAdminReport(dataset) {
  if (!hasPermission("report.export")) {
    setAdminStatus(adminReportsStatusEl, "Không đủ quyền", "error");
    return;
  }
  const params = adminFilterParams("reports");
  params.set("dataset", dataset);
  const link = document.createElement("a");
  link.href = `/api/v1/admin/reports/export?${params.toString()}`;
  link.download = "";
  document.body.appendChild(link);
  link.click();
  link.remove();
}

async function loadOperations(force = false) {
  if (!force && operationsLoaded) {
    return;
  }
  setAdminStatus(operationsStatusEl, "Đang tải", "busy");
  const limit = Math.min(Math.max(Number(operationLimitEl.value || 100), 10), 500);
  const operationParams = new URLSearchParams();
  if (operationSourceFilterEl.value) operationParams.set("source", operationSourceFilterEl.value);
  if (operationLevelFilterEl.value) operationParams.set("level", operationLevelFilterEl.value);
  operationParams.set("limit", String(limit));
  const auditParams = new URLSearchParams({ limit: String(limit) });
  const [operationRes, auditRes] = await Promise.all([
    fetch(`/api/v1/admin/operation-logs?${operationParams.toString()}`),
    fetch(`/api/v1/admin/audit-logs?${auditParams.toString()}`),
  ]);
  const [operationText, auditText] = await Promise.all([operationRes.text(), auditRes.text()]);
  if (!operationRes.ok || !auditRes.ok) {
    operationsLoaded = false;
    renderOperations(null);
    setAdminStatus(operationsStatusEl, operationText || auditText || "Chưa cấu hình DB", "error");
    return;
  }
  const operationData = JSON.parse(operationText);
  const auditData = JSON.parse(auditText);
  operationsLoaded = true;
  renderOperations({ operationLogs: operationData.logs || [], auditLogs: auditData.logs || [] });
  setAdminStatus(operationsStatusEl, "Sẵn sàng", "ready");
}

function renderOperations(data) {
  renderOperationLogs(data?.operationLogs || []);
  renderAuditLogs(data?.auditLogs || []);
}

function renderOperationLogs(rows) {
  operationLogCountEl.textContent = `${rows.length} log`;
  operationLogRowsEl.innerHTML = rows
    .map(
      (row) => `
        <tr>
          <td><strong>${escapeHtml(formatDateTime(row.occurredAt))}</strong><small>${escapeHtml(row.operation || "")}</small></td>
          <td><span class="tag">${escapeHtml(row.source || "")}</span><small>${escapeHtml(row.level || "")}</small></td>
          <td><span class="tag">${escapeHtml(row.status || "")}</span></td>
          <td>${escapeHtml(row.message || "")}</td>
          <td>${escapeHtml(row.entityType || "-")}<small>${escapeHtml(row.entityId || "")}</small></td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    operationLogRowsEl.innerHTML = `<tr><td colspan="5" class="empty-cell">Chưa có operational log</td></tr>`;
  }
}

function renderAuditLogs(rows) {
  auditLogCountEl.textContent = `${rows.length} log`;
  auditLogRowsEl.innerHTML = rows
    .map(
      (row) => `
        <tr>
          <td><strong>${escapeHtml(formatDateTime(row.occurredAt))}</strong><small>${escapeHtml(row.requestId || "")}</small></td>
          <td>${escapeHtml(row.actorName || row.actorUserId || "-")}<small>${escapeHtml(row.ipAddress || "")}</small></td>
          <td><span class="tag">${escapeHtml(row.action || "")}</span></td>
          <td>${escapeHtml(row.reason || metadataReason(row.metadata) || "")}</td>
          <td>${escapeHtml(row.entityType || "-")}<small>${escapeHtml(row.entityId || "")}</small></td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    auditLogRowsEl.innerHTML = `<tr><td colspan="5" class="empty-cell">Chưa có audit log</td></tr>`;
  }
}

function metadataReason(metadata) {
  if (!metadata || typeof metadata !== "object") return "";
  return String(metadata.reason || "");
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
      const roleText = (user.roles || []).map((role) => role.name || role.code).join(", ");
      return `
        <tr data-admin-user-id="${escapeAttr(user.id || "")}">
          <td><strong>${escapeHtml(user.displayName || "")}</strong></td>
          <td>${escapeHtml(user.email || "-")}</td>
          <td>${escapeHtml(user.phone || "-")}</td>
          <td><span class="tag">${escapeHtml(user.status || "")}</span></td>
          <td><span class="tag">${user.hasPassword ? "Set" : "Missing"}</span></td>
          <td>${escapeHtml(roleText || "-")}</td>
        </tr>
      `;
    })
    .join("");
  if (!users.length) {
    adminUserRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Chưa có user</td></tr>`;
  }
  adminUserRowsEl.querySelectorAll("[data-admin-user-id]").forEach((row) => {
    row.addEventListener("click", () => {
      selectAdminUser(row.dataset.adminUserId);
      openAdminUserDialog();
    });
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
  adminUserPhoneEl.value = user.phone || "";
  adminUserDisplayNameEl.value = user.displayName || "";
  adminUserStatusEl.value = optionValueOrEmpty(adminUserStatusEl, user.status || "active") || "active";
  adminUserPasswordEl.value = "";
  const roleCodes = (user.roles || []).map((role) => role.code);
  [...adminUserRolesEl.options].forEach((option) => {
    option.selected = roleCodes.includes(option.value);
  });
}

function openAdminUserDialog() {
  openAppDialog({
    title: adminUserIdEl.value ? "Sửa user" : "Tạo user",
    kicker: "User and role admin",
    icon: "admin_panel_settings",
    nodes: [document.querySelector(".admin-user-form")],
    size: "lg",
    actions: [
      { label: "Đóng", icon: "close", onClick: closeAppDialog },
      { label: "Clear", icon: "backspace", onClick: clearAdminUserForm },
      { label: "Lưu roles", icon: "assignment_ind", onClick: assignAdminUserRoles },
      { label: "Lưu user", icon: "save", variant: "primary", onClick: saveAdminUser, closeOnSuccess: true },
    ],
  });
}

function clearAdminUserForm() {
  adminUserIdEl.value = "";
  adminUserEmailEl.value = "";
  adminUserPhoneEl.value = "";
  adminUserDisplayNameEl.value = "";
  adminUserStatusEl.value = "active";
  adminUserPasswordEl.value = "";
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
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      id: adminUserIdEl.value,
      email: adminUserEmailEl.value,
      phone: adminUserPhoneEl.value,
      displayName: adminUserDisplayNameEl.value,
      status: adminUserStatusEl.value,
      password: adminUserPasswordEl.value,
    }),
  });
  const text = await res.text();
  if (!res.ok) {
    setAdminStatus(adminUsersStatusEl, text || "Không lưu được user", "error");
    return false;
  }
  const data = JSON.parse(text);
  if (data.user?.id) {
    adminUserIdEl.value = data.user.id;
  }
  adminUserPasswordEl.value = "";
  adminUsersLoaded = false;
  await loadAdminUsers(true);
  if (data.user?.id) {
    selectAdminUser(data.user.id);
  }
  setAdminStatus(adminUsersStatusEl, "Đã lưu user", "ready");
  return true;
}

async function assignAdminUserRoles() {
  if (!adminUserIdEl.value) {
    setAdminStatus(adminUsersStatusEl, "Chọn hoặc lưu user trước", "error");
    return false;
  }
  setAdminStatus(adminUsersStatusEl, "Đang lưu roles", "busy");
  const res = await fetch("/api/v1/admin/users/roles", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      userId: adminUserIdEl.value,
      roleCodes: selectedOptionValues(adminUserRolesEl),
    }),
  });
  const text = await res.text();
  if (!res.ok) {
    setAdminStatus(adminUsersStatusEl, text || "Không lưu được roles", "error");
    return false;
  }
  adminUsersLoaded = false;
  await loadAdminUsers(true);
  selectAdminUser(adminUserIdEl.value);
  setAdminStatus(adminUsersStatusEl, "Đã lưu roles", "ready");
  return true;
}

async function loadMasterStudents() {
  if (!masterDataLoaded) {
    return;
  }
  setMasterStatus("Đang tải", "busy");
  const params = new URLSearchParams();
  if (masterSchoolFilterEl.value) params.set("schoolId", masterSchoolFilterEl.value);
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

function renderMasterStudents(students = []) {
  masterStudentsData = students || [];
  if (!masterStudentsData.some((student) => masterStudentKey(student) === selectedMasterStudentKey)) {
    selectedMasterStudentKey = masterStudentsData[0] ? masterStudentKey(masterStudentsData[0]) : "";
  }
  masterStudentCountEl.textContent = `${students.length} học sinh`;
  masterStudentsEl.innerHTML = students
    .map((student) => {
      const key = masterStudentKey(student);
      const primary = (student.parents || []).find((parent) => parent.isPrimary) || (student.parents || [])[0] || {};
      const parentNames = (student.parents || []).map((parent) => parent.parentName).filter(Boolean).join(", ");
      const billingEmails = (student.parents || [])
        .filter((parent) => parent.receivesBillingEmail && parent.isActive && parent.emailActive && parent.email)
        .map((parent) => parent.email)
        .join(", ");
      return `
        <tr data-master-student-row="${escapeAttr(key)}" class="${key === selectedMasterStudentKey ? "is-selected" : ""}">
          <td><strong>${escapeHtml(student.studentCode || "")}</strong></td>
          <td>${escapeHtml(student.studentName || "")}</td>
          <td>${escapeHtml([student.schoolCode, student.schoolYearCode].filter(Boolean).join(" · "))}</td>
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
  masterStudentsEl.querySelectorAll("[data-master-student-row]").forEach((row) => {
    row.addEventListener("click", () => selectMasterStudent(row.dataset.masterStudentRow));
  });
  const selected = masterStudentsData.find((student) => masterStudentKey(student) === selectedMasterStudentKey);
  renderMasterStudentDetail(selected);
  populateMasterStudentForm(selected);
}

function masterStudentKey(student) {
  return student?.id || student?.studentId || student?.studentCode || student?.studentName || "";
}

function selectMasterStudent(key) {
  selectedMasterStudentKey = key || "";
  masterStudentsEl.querySelectorAll("[data-master-student-row]").forEach((row) => {
    row.classList.toggle("is-selected", row.dataset.masterStudentRow === selectedMasterStudentKey);
  });
  const selected = masterStudentsData.find((student) => masterStudentKey(student) === selectedMasterStudentKey);
  renderMasterStudentDetail(selected);
  populateMasterStudentForm(selected);
}

function openMasterStudentDialog() {
  openAppDialog({
    title: masterStudentIdEl.value ? "Sửa học sinh" : "Tạo học sinh",
    kicker: "Manual upsert",
    icon: "edit_note",
    nodes: [masterStudentEditorContentEl],
    size: "lg",
    actions: [
      { label: "Hủy", icon: "close", onClick: closeAppDialog },
      { label: "Thêm phụ huynh", icon: "group_add", onClick: addMasterParentDraft },
      { label: "Lưu học sinh", icon: "save", variant: "primary", onClick: saveMasterStudent, closeOnSuccess: true },
    ],
  });
}

function renderMasterStudentDetail(student) {
  if (!student) {
    masterStudentDetailEl.innerHTML = `<div class="detail-placeholder">${muiIcon("person_search")}<span>Chọn một học sinh để xem lớp, năm học và thông tin phụ huynh.</span></div>`;
    return;
  }
  const parents = student.parents || [];
  const activeBillingEmails = parents
    .filter((parent) => parent.receivesBillingEmail && parent.isActive && parent.emailActive && parent.email)
    .map((parent) => parent.email);
  const parentList = parents.length
    ? parents
        .map((parent) => {
          const flags = [
            parent.isPrimary ? "Chính" : "",
            parent.receivesBillingEmail ? "Nhận phí" : "",
            parent.isActive === false || parent.emailActive === false ? "Tạm dừng" : "",
          ]
            .filter(Boolean)
            .join(" · ");
          return `
            <li>
              <strong>${escapeHtml(parent.parentName || "-")}</strong>
              <span>${escapeHtml(parent.relationship || "Phụ huynh")}${flags ? ` · ${escapeHtml(flags)}` : ""}</span>
              <small>${escapeHtml(parent.email || "Chưa có email")}</small>
            </li>
          `;
        })
        .join("")
    : `<li><strong>Chưa có phụ huynh</strong><span>Import hoặc cập nhật master data</span></li>`;
  masterStudentDetailEl.innerHTML = `
    <div class="detail-hero">
      ${muiIcon("person")}
      <div>
        <strong>${escapeHtml(student.studentName || "-")}</strong>
        <span>${escapeHtml(student.studentCode || "Chưa có mã HS")}</span>
      </div>
    </div>
    <div class="detail-grid">
      <span>Năm học</span><strong>${escapeHtml(student.schoolYearCode || "-")}</strong>
      <span>Trường</span><strong>${escapeHtml(student.schoolCode || "-")}</strong>
      <span>Khối / lớp</span><strong>${escapeHtml([student.grade ? `Khối ${student.grade}` : "", student.className || ""].filter(Boolean).join(" · ") || "-")}</strong>
      <span>Trạng thái</span><strong>${escapeHtml(student.status || "active")}</strong>
      <span>Email nhận phí</span><strong>${escapeHtml(activeBillingEmails.join(", ") || "-")}</strong>
    </div>
    <div class="detail-section">
      <p class="detail-section-title">Phụ huynh</p>
      <ul class="detail-list">${parentList}</ul>
    </div>
    <div class="detail-actions">
      <button data-edit-master-student="true" type="button">${muiIcon("edit")}<span>Sửa học sinh</span></button>
    </div>
  `;
  masterStudentDetailEl.querySelector("[data-edit-master-student]")?.addEventListener("click", openMasterStudentDialog);
}

function defaultMasterParentDraft(isPrimary = false) {
  return {
    id: "",
    parentName: "",
    email: "",
    emailActive: true,
    isPrimary,
    isActive: true,
    receivesBillingEmail: true,
  };
}

function populateMasterStudentForm(student) {
  if (!student) {
    clearMasterStudentForm();
    return;
  }
  masterStudentIdEl.value = student.id || "";
  masterStudentCodeEl.value = student.studentCode || "";
  masterStudentNameEl.value = student.studentName || "";
  renderMasterStudentClassSelect(student.classId || "");
  masterStudentStatusEl.value = optionValueOrEmpty(masterStudentStatusEl, student.status || "active") || "active";
  masterStudentParentDrafts = (student.parents || []).map((parent) => ({
    id: parent.id || "",
    parentName: parent.parentName || "",
    email: parent.email || "",
    emailActive: parent.emailActive !== false,
    isPrimary: !!parent.isPrimary,
    isActive: parent.isActive !== false,
    receivesBillingEmail: !!parent.receivesBillingEmail,
  }));
  if (!masterStudentParentDrafts.length) {
    masterStudentParentDrafts = [defaultMasterParentDraft(true)];
  }
  renderMasterParentEditorRows();
}

function clearMasterStudentForm() {
  masterStudentIdEl.value = "";
  masterStudentCodeEl.value = "";
  masterStudentNameEl.value = "";
  renderMasterStudentClassSelect(masterClassFilterEl.value || "");
  masterStudentStatusEl.value = "active";
  masterStudentParentDrafts = [defaultMasterParentDraft(true)];
  renderMasterParentEditorRows();
}

function renderMasterParentEditorRows() {
  if (!masterStudentParentDrafts.length) {
    masterStudentParentDrafts = [defaultMasterParentDraft(true)];
  }
  masterParentEditorRowsEl.innerHTML = masterStudentParentDrafts
    .map((parent, index) => `
      <div class="master-parent-row" data-master-parent-index="${index}">
        <input type="hidden" data-master-parent-field="id" value="${escapeAttr(parent.id || "")}" />
        <label>
          <span>Phụ huynh</span>
          <input data-master-parent-field="parentName" value="${escapeAttr(parent.parentName || "")}" />
        </label>
        <label>
          <span>Email</span>
          <input data-master-parent-field="email" type="email" value="${escapeAttr(parent.email || "")}" />
        </label>
        <label class="master-parent-check">
          <input data-master-parent-field="isPrimary" name="masterParentPrimary" type="radio" ${parent.isPrimary ? "checked" : ""} />
          <span>Chính</span>
        </label>
        <label class="master-parent-check">
          <input data-master-parent-field="isActive" type="checkbox" ${parent.isActive !== false ? "checked" : ""} />
          <span>Active</span>
        </label>
        <label class="master-parent-check">
          <input data-master-parent-field="receivesBillingEmail" type="checkbox" ${parent.receivesBillingEmail !== false ? "checked" : ""} />
          <span>Nhận phí</span>
        </label>
        <label class="master-parent-check">
          <input data-master-parent-field="emailActive" type="checkbox" ${parent.emailActive !== false ? "checked" : ""} />
          <span>Email active</span>
        </label>
        <button data-remove-master-parent="${index}" type="button" ${masterStudentParentDrafts.length === 1 ? "disabled" : ""}>
          ${muiIcon("delete")}<span>Xóa</span>
        </button>
      </div>
    `)
    .join("");
  masterParentEditorRowsEl.querySelectorAll("[data-remove-master-parent]").forEach((button) => {
    button.addEventListener("click", () => {
      masterStudentParentDrafts = collectMasterParentDrafts();
      masterStudentParentDrafts.splice(Number(button.dataset.removeMasterParent), 1);
      renderMasterParentEditorRows();
    });
  });
}

function collectMasterParentDrafts() {
  return [...masterParentEditorRowsEl.querySelectorAll("[data-master-parent-index]")].map((row) => {
    const field = (name) => row.querySelector(`[data-master-parent-field="${name}"]`);
    return {
      id: field("id")?.value.trim() || "",
      parentName: field("parentName")?.value.trim() || "",
      email: field("email")?.value.trim() || "",
      isPrimary: !!field("isPrimary")?.checked,
      isActive: !!field("isActive")?.checked,
      receivesBillingEmail: !!field("receivesBillingEmail")?.checked,
      emailActive: !!field("emailActive")?.checked,
    };
  });
}

function addMasterParentDraft() {
  masterStudentParentDrafts = collectMasterParentDrafts();
  masterStudentParentDrafts.push(defaultMasterParentDraft(!masterStudentParentDrafts.some((parent) => parent.isPrimary && parent.isActive)));
  renderMasterParentEditorRows();
}

async function saveMasterStudent() {
  const payload = {
    id: masterStudentIdEl.value.trim(),
    studentCode: masterStudentCodeEl.value.trim(),
    studentName: masterStudentNameEl.value.trim(),
    classId: masterStudentClassEl.value,
    status: masterStudentStatusEl.value,
    parents: collectMasterParentDrafts(),
  };
  setMasterStatus("Đang lưu", "busy");
  const res = await fetch("/api/v1/master-data/students/save", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  const text = await res.text();
  let data = null;
  try {
    data = JSON.parse(text);
  } catch {
    data = null;
  }
  if (!res.ok) {
    setMasterStatus(text || "Không lưu được học sinh", "error");
    return false;
  }
  const saved = data?.student;
  if (saved) {
    selectedMasterStudentKey = masterStudentKey(saved);
    masterSchoolFilterEl.value = optionValueOrEmpty(masterSchoolFilterEl, saved.schoolId || "") || masterSchoolFilterEl.value;
    renderMasterFilters();
    masterSchoolYearFilterEl.value = optionValueOrEmpty(masterSchoolYearFilterEl, saved.schoolYearId || "") || masterSchoolYearFilterEl.value;
    renderMasterFilters();
    masterGradeFilterEl.value = optionValueOrEmpty(masterGradeFilterEl, saved.grade || "") || masterGradeFilterEl.value;
    renderMasterClassFilter(saved.classId || "");
  }
  await loadMasterData(true);
  setMasterStatus("Đã lưu học sinh", "ready");
  return true;
}

async function submitMasterImport(apply) {
  const file = masterImportState?.file || masterCsvFileEl.files[0];
  if (!file) {
    setMasterStatus("Chưa chọn file", "error");
    return;
  }
  if (apply) {
    const confirmed = await confirmDialog({
      title: "Áp dụng import master data?",
      message: "Thao tác này sẽ ghi dữ liệu học sinh, phụ huynh và lớp vào database.",
      confirmLabel: "Áp dụng import",
      confirmIcon: "publish",
      danger: true,
    });
    if (!confirmed) return;
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
    operatorName: feeScheduleOperatorEl.value.trim(),
    items: collectFeeScheduleItems(),
    adjustments: parseFeeAdjustmentsCsv(),
  };
}

function openFeeScheduleDialog() {
  openAppDialog({
    title: "Tạo/sửa bảng phí",
    kicker: "Production fees",
    icon: "price_change",
    nodes: [document.querySelector(".fee-schedule-grid"), document.querySelector(".fee-schedule-body")],
    size: "xl",
    actions: [
      { label: "Đóng", icon: "close", onClick: closeAppDialog },
      { label: "Preview", icon: "visibility", onClick: previewFeeSchedule },
      { label: "Lưu bảng phí", icon: "save", variant: "primary", onClick: saveFeeSchedule, closeOnSuccess: true },
    ],
  });
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
    return false;
  }
  renderFeeSchedulePreview(data);
  setFeeScheduleStatus(data.issues?.length ? "Có lỗi" : "Đã preview", data.issues?.length ? "error" : "ready");
  return !data.issues?.length;
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
    return false;
  }
  renderFeeSchedulePreview(data.preview || null);
  renderFeeSchedules(data.schedules || []);
  setFeeScheduleStatus("Đã lưu", "ready");
  return true;
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

function openInvoiceDialog() {
  openAppDialog({
    title: "Cấu hình sinh hóa đơn",
    kicker: "Production invoices",
    icon: "request_quote",
    nodes: [document.querySelector(".invoice-toolbar")],
    size: "lg",
    actions: [
      { label: "Đóng", icon: "close", onClick: closeAppDialog },
      { label: "Preview", icon: "visibility", onClick: previewInvoices },
    ],
  });
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
    return false;
  }
  renderInvoicePreview(data);
  setInvoiceStatus(data.issues?.length ? "Có lỗi" : "Đã preview", data.issues?.length ? "error" : "ready");
  return !data.issues?.length;
}

async function generateInvoices() {
  const confirmed = await confirmDialog({
    title: "Sinh hóa đơn?",
    message: "Thao tác này sẽ ghi dữ liệu invoice vào database.",
    confirmLabel: "Sinh hóa đơn",
    confirmIcon: "post_add",
    danger: true,
  });
  if (!confirmed) return false;
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
    return false;
  }
  renderInvoicePreview(data.preview || null);
  renderInvoices(data.invoices || []);
  setInvoiceStatus("Đã sinh hóa đơn", "ready");
  return true;
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

function renderInvoices(invoices = []) {
  invoicesData = invoices || [];
  if (!invoicesData.some((invoice) => invoice.id === selectedInvoiceId)) {
    selectedInvoiceId = invoicesData[0]?.id || "";
  }
  invoiceCountEl.textContent = `${invoices.length} hóa đơn`;
  invoiceRowsEl.innerHTML = invoices
    .map(
      (invoice) => `
        <tr data-invoice-row="${escapeAttr(invoice.id || "")}" class="${invoice.id === selectedInvoiceId ? "is-selected" : ""}">
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
  invoiceRowsEl.querySelectorAll("[data-invoice-row]").forEach((row) => {
    row.addEventListener("click", (event) => {
      if (event.target.closest("button, a")) return;
      selectInvoice(row.dataset.invoiceRow);
    });
  });
  invoiceRowsEl.querySelectorAll("[data-invoice-qr]").forEach((button) => {
    button.addEventListener("click", () => loadInvoicePayment(button.dataset.invoiceQr));
  });
  if (selectedInvoiceId) {
    selectInvoice(selectedInvoiceId);
  } else {
    renderInvoiceDetail(null);
    invoicePaymentStatusEl.textContent = "Chưa chọn";
    invoicePaymentStatusEl.dataset.tone = "";
    invoicePaymentPreviewEl.className = "preview-empty";
    invoicePaymentPreviewEl.textContent = "Chưa chọn hóa đơn";
  }
}

async function loadInvoicePayment(invoiceId) {
  if (!invoiceId) {
    return;
  }
  selectInvoice(invoiceId, { keepQr: true });
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

function selectInvoice(invoiceId, options = {}) {
  selectedInvoiceId = invoiceId || "";
  invoiceRowsEl.querySelectorAll("[data-invoice-row]").forEach((row) => {
    row.classList.toggle("is-selected", row.dataset.invoiceRow === selectedInvoiceId);
  });
  renderInvoiceDetail(invoicesData.find((invoice) => invoice.id === selectedInvoiceId));
  if (!options.keepQr) {
    invoicePaymentStatusEl.textContent = selectedInvoiceId ? "Đã chọn" : "Chưa chọn";
    invoicePaymentStatusEl.dataset.tone = selectedInvoiceId ? "ready" : "";
    invoicePaymentPreviewEl.className = "preview-empty";
    invoicePaymentPreviewEl.textContent = selectedInvoiceId ? "Bấm QR để xem payload thanh toán của hóa đơn này" : "Chưa chọn hóa đơn";
  }
}

function renderInvoiceDetail(invoice) {
  if (!invoice) {
    invoiceDetailSummaryEl.innerHTML = `<div class="detail-placeholder">${muiIcon("receipt_long")}<span>Chọn một hóa đơn để xem tổng tiền, trạng thái và thao tác QR/PDF.</span></div>`;
    return;
  }
  const outstanding = Math.max(Number(invoice.totalAmount || 0) - Number(invoice.paidAmount || 0), 0);
  invoiceDetailSummaryEl.innerHTML = `
    <div class="detail-hero">
      ${muiIcon("receipt_long")}
      <div>
        <strong>${escapeHtml(invoice.invoiceCode || "-")}</strong>
        <span>${escapeHtml(invoice.studentCode || "")} · ${escapeHtml(invoice.studentName || "")}</span>
      </div>
    </div>
    <div class="detail-grid">
      <span>Lớp / kỳ</span><strong>${escapeHtml([invoice.className || "", invoice.periodCode || ""].filter(Boolean).join(" · ") || "-")}</strong>
      <span>Tổng tiền</span><strong>${formatMoney(invoice.totalAmount || 0)}</strong>
      <span>Đã thu</span><strong>${formatMoney(invoice.paidAmount || 0)}</strong>
      <span>Còn thiếu</span><strong>${formatMoney(outstanding)}</strong>
      <span>Status</span><strong>${escapeHtml(invoice.status || "unpaid")}</strong>
    </div>
    <div class="detail-actions">
      <button type="button" data-detail-invoice-qr="${escapeAttr(invoice.id || "")}">${muiIcon("qr_code")}<span>Xem QR</span></button>
      <a class="button-link" href="/api/v1/invoices/pdf?id=${encodeURIComponent(invoice.id || "")}" target="_blank" rel="noreferrer">${muiIcon("picture_as_pdf")}<span>Mở PDF</span></a>
    </div>
  `;
  invoiceDetailSummaryEl.querySelectorAll("[data-detail-invoice-qr]").forEach((button) => {
    button.addEventListener("click", () => loadInvoicePayment(button.dataset.detailInvoiceQr));
  });
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
  const canWritePayments = hasPermission("payment.create");
  paymentReconInvoiceRowsEl.innerHTML = invoices
    .map((invoice) => {
      const paid = Number(invoice.paidAmount || 0);
      const total = Number(invoice.totalAmount || 0);
      const outstanding = Math.max(total - paid, 0);
      const intent = intents?.[invoice.id];
      const intentLabel = intent?.provider ? `${intent.provider}: ${intent.status}` : "";
      const actions = canWritePayments
        ? `
          <div class="invoice-actions">
            <button type="button" data-recon-intent="${escapeAttr(invoice.id || "")}" data-recon-provider="manual_vietqr">${muiIcon("qr_code")}<span>QR</span></button>
            ${hasPayOS ? `<button type="button" data-recon-intent="${escapeAttr(invoice.id || "")}" data-recon-provider="payos">${muiIcon("link")}<span>payOS</span></button>` : ""}
            <button type="button" data-recon-cash="${escapeAttr(invoice.id || "")}" data-recon-default-amount="${escapeAttr(outstanding || total)}">${muiIcon("payments")}<span>Tiền mặt</span></button>
          </div>
        `
        : "";
      return `
        <tr data-recon-invoice-row="${escapeAttr(invoice.id || "")}">
          <td><strong>${escapeHtml(invoice.invoiceCode || "")}</strong>${intentLabel ? `<small>${escapeHtml(intentLabel)}</small>` : ""}</td>
          <td>${escapeHtml(invoice.studentCode || "")} · ${escapeHtml(invoice.studentName || "")}</td>
          <td>${escapeHtml(invoice.className || "")}</td>
          <td class="money">${formatMoney(total)}</td>
          <td class="money">${formatMoney(paid)}</td>
          <td><span class="tag">${escapeHtml(invoice.status || "unpaid")}</span></td>
          <td>${actions}</td>
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
    row.classList.toggle("is-selected", paymentReconSelection.type === "invoice" && paymentReconSelection.id === row.dataset.reconInvoiceRow);
    row.addEventListener("click", (event) => {
      if (event.target.closest("button, a")) return;
      const invoice = (paymentReconciliationData.invoices || []).find((item) => item.id === row.dataset.reconInvoiceRow);
      renderPaymentReconDetail(invoiceDetailTemplate(invoice), { type: "invoice", id: row.dataset.reconInvoiceRow });
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
    row.classList.toggle("is-selected", paymentReconSelection.type === "transaction" && paymentReconSelection.id === row.dataset.reconTransactionRow);
    row.addEventListener("click", () => {
      const transaction = (paymentReconciliationData.transactions || []).find((item) => item.id === row.dataset.reconTransactionRow);
      renderPaymentReconDetail(transactionDetailTemplate(transaction), { type: "transaction", id: row.dataset.reconTransactionRow });
    });
  });
}

async function createPaymentIntent(invoiceId, provider) {
  if (!invoiceId) return;
  paymentReconSelection = { type: "invoice", id: invoiceId };
  updatePaymentReconActiveRows();
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
    renderPaymentReconDetail(`<div class="reconciliation-error">${escapeHtml(text || "Không tạo được payment intent")}</div>`, { type: "invoice", id: invoiceId });
    return;
  }
  renderPaymentReconDetail(paymentIntentDetailTemplate(data), { type: "invoice", id: invoiceId });
  paymentReconciliationLoaded = false;
  await loadPaymentReconciliation(true);
}

function cashReceiptDialog(defaultAmount) {
  return new Promise((resolve) => {
    let settled = false;
    const body = document.createElement("div");
    body.className = "dialog-form-grid";
    body.innerHTML = `
      <label>
        <span>Số tiền thu tiền mặt</span>
        <input data-cash-field="amount" inputmode="numeric" value="${escapeAttr(String(defaultAmount || ""))}" />
      </label>
      <label>
        <span>Người thu tiền</span>
        <input data-cash-field="collectorName" />
      </label>
      <label>
        <span>Mã phiếu thu</span>
        <input data-cash-field="receiptReference" value="${escapeAttr(`CASH${Date.now()}`)}" />
      </label>
      <label>
        <span>Lý do ghi nhận</span>
        <input data-cash-field="reason" value="Thu tiền mặt học phí" />
      </label>
    `;
    const field = (name) => body.querySelector(`[data-cash-field="${name}"]`);
    openAppDialog({
      title: "Ghi nhận tiền mặt",
      kicker: "Đối soát",
      icon: "payments",
      content: body,
      size: "md",
      onClose: () => {
        if (!settled) resolve(null);
      },
      actions: [
        {
          label: "Hủy",
          icon: "close",
          onClick: () => {
            settled = true;
            resolve(null);
          },
          closeOnSuccess: true,
        },
        {
          label: "Ghi nhận",
          icon: "payments",
          variant: "primary",
          onClick: () => {
            const amount = parseMoneyInput(field("amount").value);
            if (!amount || amount <= 0) {
              showDialogError("Số tiền phải lớn hơn 0");
              return false;
            }
            settled = true;
            resolve({
              amount,
              collectorName: field("collectorName").value.trim(),
              receiptReference: field("receiptReference").value.trim(),
              reason: field("reason").value.trim(),
            });
          },
          closeOnSuccess: true,
        },
      ],
    });
  });
}

async function recordManualCashReceipt(invoiceId, defaultAmount) {
  if (!invoiceId) return;
  paymentReconSelection = { type: "invoice", id: invoiceId };
  updatePaymentReconActiveRows();
  const receipt = await cashReceiptDialog(defaultAmount);
  if (!receipt) return false;
  setPaymentReconStatus("Đang ghi nhận", "busy");
  const res = await fetch("/api/v1/payments/cash-receipts", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ invoiceId, ...receipt }),
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
    renderPaymentReconDetail(`<div class="reconciliation-error">${escapeHtml(text || "Không ghi nhận được tiền mặt")}</div>`, { type: "invoice", id: invoiceId });
    return false;
  }
  paymentReconciliationLoaded = false;
  await loadPaymentReconciliation(true);
  renderPaymentReconDetail(transactionDetailTemplate(data.transaction), { type: "transaction", id: data.transaction?.id || "" });
  setPaymentReconStatus("Đã ghi nhận", "ready");
  return true;
}

function renderPaymentReconDetail(html, selection = null) {
  if (selection) {
    paymentReconSelection = selection;
    updatePaymentReconActiveRows();
  }
  paymentReconDetailEl.innerHTML = html || "Chưa chọn hóa đơn hoặc giao dịch";
}

function updatePaymentReconActiveRows() {
  paymentReconInvoiceRowsEl.querySelectorAll("[data-recon-invoice-row]").forEach((row) => {
    row.classList.toggle("is-selected", paymentReconSelection.type === "invoice" && paymentReconSelection.id === row.dataset.reconInvoiceRow);
  });
  paymentReconTransactionRowsEl.querySelectorAll("[data-recon-transaction-row]").forEach((row) => {
    row.classList.toggle("is-selected", paymentReconSelection.type === "transaction" && paymentReconSelection.id === row.dataset.reconTransactionRow);
  });
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

function openNotificationDialog() {
  openAppDialog({
    title: "Cấu hình campaign",
    kicker: "Invoice campaigns",
    icon: "campaign",
    nodes: [document.querySelector(".notification-toolbar")],
    size: "lg",
    actions: [
      { label: "Đóng", icon: "close", onClick: closeAppDialog },
      { label: "Preview", icon: "visibility", onClick: previewNotifications },
      { label: "Lưu campaign", icon: "save", variant: "primary", onClick: saveNotificationCampaign },
    ],
  });
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
  const confirmed = await confirmDialog({
    title: "Gửi campaign?",
    message: "Thao tác này sẽ gửi email thật qua provider hiện tại và ghi log theo từng invoice/recipient.",
    confirmLabel: "Gửi campaign",
    confirmIcon: "send",
    danger: true,
  });
  if (!confirmed) return false;
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
    return false;
  }
  const data = JSON.parse(text);
  currentNotificationCampaignId = data.campaign?.id || currentNotificationCampaignId;
  renderNotificationResults(data);
  notificationLoaded = false;
  setNotificationStatus("Đã xử lý gửi", "ready");
  return true;
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
    row.addEventListener("click", async () => {
      const campaign = notificationOptions.campaigns.find((item) => item.id === row.dataset.campaignId);
      if (campaign) {
        await selectNotificationCampaign(campaign);
        openNotificationDialog();
      }
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

function openEmailConfigDialog() {
  openAppDialog({
    title: "Cấu hình email",
    kicker: "Email settings",
    icon: "settings",
    nodes: [document.querySelector(".form-grid")],
    size: "lg",
    actions: [
      { label: "Đóng", icon: "close", onClick: closeAppDialog },
      { label: "Preview email", icon: "visibility", onClick: previewEmail },
      { label: "Lưu cấu hình", icon: "save", variant: "primary", onClick: saveEmailConfig, closeOnSuccess: true },
    ],
  });
}

function openCronConfigDialog() {
  openAppDialog({
    title: "Cấu hình cron",
    kicker: "Email cron",
    icon: "schedule",
    nodes: [document.querySelector(".cron-grid")],
    size: "md",
    actions: [
      { label: "Đóng", icon: "close", onClick: closeAppDialog },
      { label: "Tắt cron", icon: "event_busy", onClick: disableEmailCron },
      { label: "Lưu cron", icon: "event_available", variant: "primary", onClick: () => saveEmailCron(cronEnabledEl.value === "true"), closeOnSuccess: true },
    ],
  });
}

async function saveEmailConfig() {
  if (!hasPermission("email_config.update")) {
    setEmailStatus("Không đủ quyền", true);
    return false;
  }
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
  if (!hasPermission("notification.send")) {
    setEmailStatus("Không đủ quyền", true);
    return;
  }
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
  if (!dryRun) {
    const confirmed = await confirmDialog({
      title: "Gửi email thật?",
      message: "Thao tác này sẽ gửi email thật qua provider hiện tại cho các dòng đang có trong bảng.",
      confirmLabel: "Gửi email",
      confirmIcon: "send",
      danger: true,
    });
    if (!confirmed) return false;
  }
  const saved = hasPermission("email_config.update") ? await saveEmailConfig() : true;
  if (!saved) return false;
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
    return false;
  }
  const data = JSON.parse(text);
  renderEmailResults(data.results || []);
  return true;
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
  const saved = hasPermission("email_config.update") ? await saveEmailConfig() : true;
  if (!saved) return false;
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
    return false;
  }
  renderCronStatus(JSON.parse(text));
  return true;
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
    return false;
  }
  renderCronStatus(JSON.parse(text));
  return true;
}

async function runEmailCronNow() {
  const confirmed = await confirmDialog({
    title: "Chạy cron ngay?",
    message: "Thao tác này sẽ gửi email thật qua provider hiện tại theo giới hạn còn lại.",
    confirmLabel: "Chạy cron",
    confirmIcon: "play_arrow",
    danger: true,
  });
  if (!confirmed) return false;
  setCronStatus("Running");
  const res = await fetch("/api/v1/email/cron/run", { method: "POST" });
  const text = await res.text();
  if (!res.ok) {
    setCronStatus(text || "Run cron failed", true);
    return false;
  }
  renderCronStatus(JSON.parse(text));
  return true;
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

loginFormEl.addEventListener("submit", submitLogin);
bootstrapFormEl.addEventListener("submit", submitBootstrap);
logoutButton.addEventListener("click", logout);

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

appDialogCloseBtn.addEventListener("click", closeAppDialog);
appDialogEl.addEventListener("click", (event) => {
  if (event.target === appDialogEl) {
    closeAppDialog();
  }
});
appDialogEl.addEventListener("close", () => {
  const onClose = activeDialogOnClose;
  activeDialogOnClose = null;
  restoreDialogContent();
  if (onClose) onClose();
});

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
exportAdminReportClassesBtn.addEventListener("click", () => exportAdminReport("classes"));
exportAdminReportInvoicesBtn.addEventListener("click", () => exportAdminReport("invoices"));
exportAdminReportTransactionsBtn.addEventListener("click", () => exportAdminReport("transactions"));
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
refreshOperationsBtn.addEventListener("click", () => loadOperations(true));
operationSourceFilterEl.addEventListener("change", () => loadOperations(true));
operationLevelFilterEl.addEventListener("change", () => loadOperations(true));
operationLimitEl.addEventListener("change", () => loadOperations(true));

refreshAdminUsersBtn.addEventListener("click", () => loadAdminUsers(true));
newAdminUserBtn.addEventListener("click", () => {
  clearAdminUserForm();
  openAdminUserDialog();
});
clearAdminUserBtn.addEventListener("click", clearAdminUserForm);
saveAdminUserBtn.addEventListener("click", saveAdminUser);
assignAdminUserRolesBtn.addEventListener("click", assignAdminUserRoles);

masterSchoolFilterEl.addEventListener("change", async () => {
  renderMasterFilters();
  await loadMasterStudents();
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
newMasterStudentBtn.addEventListener("click", () => {
  clearMasterStudentForm();
  openMasterStudentDialog();
});
editMasterStudentBtn.addEventListener("click", openMasterStudentDialog);
addMasterParentBtn.addEventListener("click", addMasterParentDraft);
saveMasterStudentBtn.addEventListener("click", saveMasterStudent);
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
openSchoolTreeSchoolDialogBtn.addEventListener("click", () => openSchoolTreeEntityDialog("school"));
openSchoolTreeYearDialogBtn.addEventListener("click", () => openSchoolTreeEntityDialog("year"));
openSchoolTreeClassDialogBtn.addEventListener("click", () => openSchoolTreeEntityDialog("class"));
saveSchoolTreeSchoolBtn.addEventListener("click", saveSchoolTreeSchool);
newSchoolTreeSchoolBtn.addEventListener("click", clearSchoolTreeSchoolForm);
schoolTreeYearSchoolEl.addEventListener("change", renderSchoolTreeSelects);
saveSchoolTreeYearBtn.addEventListener("click", saveSchoolTreeYear);
newSchoolTreeYearBtn.addEventListener("click", clearSchoolTreeYearForm);
saveSchoolTreeClassBtn.addEventListener("click", saveSchoolTreeClass);
newSchoolTreeClassBtn.addEventListener("click", clearSchoolTreeClassForm);

refreshFeeSchedulesBtn.addEventListener("click", () => loadFeeSchedules(true));
openFeeScheduleDialogBtn.addEventListener("click", openFeeScheduleDialog);
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
openInvoiceDialogBtn.addEventListener("click", openInvoiceDialog);
previewInvoicesBtn.addEventListener("click", previewInvoices);
generateInvoicesBtn.addEventListener("click", generateInvoices);
refreshPaymentReconBtn.addEventListener("click", () => loadPaymentReconciliation(true));
paymentProviderFilterEl.addEventListener("change", () => loadPaymentReconciliation(true));
paymentInvoiceStatusFilterEl.addEventListener("change", () => loadPaymentReconciliation(true));
paymentTransactionStatusFilterEl.addEventListener("change", () => loadPaymentReconciliation(true));
refreshNotificationsBtn.addEventListener("click", () => loadNotifications(true));
openNotificationDialogBtn.addEventListener("click", openNotificationDialog);
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
openEmailConfigDialogBtn.addEventListener("click", openEmailConfigDialog);
openCronConfigDialogBtn.addEventListener("click", openCronConfigDialog);
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
