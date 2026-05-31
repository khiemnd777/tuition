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
const masterBillingFilterEl = document.querySelector("#masterBillingFilter");
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
const feeGuideStepsEl = document.querySelector("#feeGuideSteps");
const feeScheduleSchoolEl = document.querySelector("#feeScheduleSchool");
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
const feeAdjustmentRowsEl = document.querySelector("#feeAdjustmentRows");
const addFeeAdjustmentRowBtn = document.querySelector("#addFeeAdjustmentRow");
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
const invoiceWorkbenchStepsEl = document.querySelector("#invoiceWorkbenchSteps");
const invoiceScheduleEl = document.querySelector("#invoiceSchedule");
const invoiceBankBinEl = document.querySelector("#invoiceBankBin");
const invoiceBankAccountEl = document.querySelector("#invoiceBankAccount");
const invoiceIssueDateEl = document.querySelector("#invoiceIssueDate");
const invoiceDueDateEl = document.querySelector("#invoiceDueDate");
const invoiceRegenerateEl = document.querySelector("#invoiceRegenerate");
const invoicePreviewCountEl = document.querySelector("#invoicePreviewCount");
const invoicePreviewSummaryEl = document.querySelector("#invoicePreviewSummary");
const invoiceIssuePanelEl = document.querySelector("#invoiceIssuePanel");
const invoicePreviewRowsEl = document.querySelector("#invoicePreviewRows");
const invoiceRowsEl = document.querySelector("#invoiceRows");
const invoiceCountEl = document.querySelector("#invoiceCount");
const exportInvoiceCsvBtn = document.querySelector("#exportInvoiceCsv");
const invoiceDetailSummaryEl = document.querySelector("#invoiceDetailSummary");
const invoicePaymentStatusEl = document.querySelector("#invoicePaymentStatus");
const invoicePaymentPreviewEl = document.querySelector("#invoicePaymentPreview");
const paymentReconStatusEl = document.querySelector("#paymentReconStatus");
const refreshPaymentReconBtn = document.querySelector("#refreshPaymentRecon");
const paymentReconStepsEl = document.querySelector("#paymentReconSteps");
const paymentReconSchoolFilterEl = document.querySelector("#paymentReconSchoolFilter");
const paymentReconYearFilterEl = document.querySelector("#paymentReconYearFilter");
const paymentReconGradeFilterEl = document.querySelector("#paymentReconGradeFilter");
const paymentReconClassFilterEl = document.querySelector("#paymentReconClassFilter");
const paymentReconPeriodFilterEl = document.querySelector("#paymentReconPeriodFilter");
const paymentProviderFilterEl = document.querySelector("#paymentProviderFilter");
const paymentInvoiceStatusFilterEl = document.querySelector("#paymentInvoiceStatusFilter");
const paymentTransactionStatusFilterEl = document.querySelector("#paymentTransactionStatusFilter");
const paymentReconSummaryEl = document.querySelector("#paymentReconSummary");
const paymentReconInvoiceCountEl = document.querySelector("#paymentReconInvoiceCount");
const paymentReconInvoiceRowsEl = document.querySelector("#paymentReconInvoiceRows");
const paymentReconTransactionCountEl = document.querySelector("#paymentReconTransactionCount");
const paymentReconTransactionRowsEl = document.querySelector("#paymentReconTransactionRows");
const paymentReconReviewCountEl = document.querySelector("#paymentReconReviewCount");
const paymentReconReviewRowsEl = document.querySelector("#paymentReconReviewRows");
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
const notificationForceResendEl = document.querySelector("#notificationForceResend");
const notificationWorkbenchStepsEl = document.querySelector("#notificationWorkbenchSteps");
const previewNotificationsBtn = document.querySelector("#previewNotifications");
const saveNotificationCampaignBtn = document.querySelector("#saveNotificationCampaign");
const retryNotificationRecipientsBtn = document.querySelector("#retryNotificationRecipients");
const sendNotificationCampaignBtn = document.querySelector("#sendNotificationCampaign");
const notificationSummaryEl = document.querySelector("#notificationSummary");
const notificationRecipientCountEl = document.querySelector("#notificationRecipientCount");
const notificationRecipientsEl = document.querySelector("#notificationRecipients");
const notificationEmailPreviewStatusEl = document.querySelector("#notificationEmailPreviewStatus");
const notificationEmailPreviewFrameEl = document.querySelector("#notificationEmailPreviewFrame");
const notificationCampaignCountEl = document.querySelector("#notificationCampaignCount");
const notificationCampaignRowsEl = document.querySelector("#notificationCampaignRows");
const notificationLogCountEl = document.querySelector("#notificationLogCount");
const notificationLogsEl = document.querySelector("#notificationLogs");
const notificationCronSnapshotEl = document.querySelector("#notificationCronSnapshot");
const openNotificationCronConfigBtn = document.querySelector("#openNotificationCronConfig");
const adminDashboardStatusEl = document.querySelector("#adminDashboardStatus");
const refreshAdminDashboardBtn = document.querySelector("#refreshAdminDashboard");
const adminDashboardSchoolEl = document.querySelector("#adminDashboardSchool");
const adminDashboardYearEl = document.querySelector("#adminDashboardYear");
const adminDashboardGradeEl = document.querySelector("#adminDashboardGrade");
const adminDashboardClassEl = document.querySelector("#adminDashboardClass");
const adminDashboardPeriodEl = document.querySelector("#adminDashboardPeriod");
const adminDashboardMonthEl = document.querySelector("#adminDashboardMonth");
const adminDashboardInvoiceStatusEl = document.querySelector("#adminDashboardInvoiceStatus");
const adminDashboardMetricsEl = document.querySelector("#adminDashboardMetrics");
const operatorOnboardingEl = document.querySelector("#operatorOnboarding");
const adminWorkQueueEl = document.querySelector("#adminWorkQueue");
const adminQuickActionsEl = document.querySelector("#adminQuickActions");
const adminReadinessSeverityEl = document.querySelector("#adminReadinessSeverity");
const adminReadinessTypeEl = document.querySelector("#adminReadinessType");
const adminReadinessCenterEl = document.querySelector("#adminReadinessCenter");
const adminTopClassCountEl = document.querySelector("#adminTopClassCount");
const adminTopClassRowsEl = document.querySelector("#adminTopClassRows");
const adminAttentionCountEl = document.querySelector("#adminAttentionCount");
const adminAttentionRowsEl = document.querySelector("#adminAttentionRows");
const adminReportsStatusEl = document.querySelector("#adminReportsStatus");
const refreshAdminReportsBtn = document.querySelector("#refreshAdminReports");
const exportAdminReportClassesBtn = document.querySelector("#exportAdminReportClasses");
const exportAdminReportInvoicesBtn = document.querySelector("#exportAdminReportInvoices");
const exportAdminReportTransactionsBtn = document.querySelector("#exportAdminReportTransactions");
const adminReportsSchoolEl = document.querySelector("#adminReportsSchool");
const adminReportsYearEl = document.querySelector("#adminReportsYear");
const adminReportsGradeEl = document.querySelector("#adminReportsGrade");
const adminReportsClassEl = document.querySelector("#adminReportsClass");
const adminReportsPeriodEl = document.querySelector("#adminReportsPeriod");
const adminReportsMonthEl = document.querySelector("#adminReportsMonth");
const adminReportsInvoiceStatusEl = document.querySelector("#adminReportsInvoiceStatus");
const adminReportsProviderEl = document.querySelector("#adminReportsProvider");
const adminReportsSummaryEl = document.querySelector("#adminReportsSummary");
const adminReportClassCountEl = document.querySelector("#adminReportClassCount");
const adminReportClassRowsEl = document.querySelector("#adminReportClassRows");
const adminReportInvoiceCountEl = document.querySelector("#adminReportInvoiceCount");
const adminReportInvoiceRowsEl = document.querySelector("#adminReportInvoiceRows");
const adminReportTransactionCountEl = document.querySelector("#adminReportTransactionCount");
const adminReportTransactionRowsEl = document.querySelector("#adminReportTransactionRows");
const operationsStatusEl = document.querySelector("#operationsStatus");
const refreshOperationsBtn = document.querySelector("#refreshOperations");
const operationSourceFilterEl = document.querySelector("#operationSourceFilter");
const operationLevelFilterEl = document.querySelector("#operationLevelFilter");
const operationNameFilterEl = document.querySelector("#operationNameFilter");
const operationStatusFilterEl = document.querySelector("#operationStatusFilter");
const auditActionFilterEl = document.querySelector("#auditActionFilter");
const operationEntityTypeFilterEl = document.querySelector("#operationEntityTypeFilter");
const operationLimitEl = document.querySelector("#operationLimit");
const operationsSummaryEl = document.querySelector("#operationsSummary");
const operationLogCountEl = document.querySelector("#operationLogCount");
const operationLogRowsEl = document.querySelector("#operationLogRows");
const auditLogCountEl = document.querySelector("#auditLogCount");
const auditLogRowsEl = document.querySelector("#auditLogRows");
const operationLogDetailEl = document.querySelector("#operationLogDetail");
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
const legacyPaymentActionsEl = document.querySelector(".legacy-payment-actions");
const appBreadcrumbsEl = document.querySelector("#appBreadcrumbs");
const currentSectionKickerEl = document.querySelector("#currentSectionKicker");
const currentSectionTitleEl = document.querySelector("#currentSectionTitle");
const currentSectionDescriptionEl = document.querySelector("#currentSectionDescription");
const appContextSchoolEl = document.querySelector("#appContextSchool");
const appContextYearEl = document.querySelector("#appContextYear");
const appContextPeriodEl = document.querySelector("#appContextPeriod");
const appContextMonthEl = document.querySelector("#appContextMonth");
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
let masterStudentsRawData = [];
let masterStudentsData = [];
let selectedMasterStudentKey = "";
let masterStudentParentDrafts = [];
let schoolTreeData = { schools: [] };
let selectedSchoolTreeNode = { type: "", id: "" };
let paymentImportState = null;
let masterImportState = null;
let feeScheduleOptions = { schools: [], feeTypes: [], schoolYears: [], classes: [] };
let feeSchedulesData = [];
let feeSchedulesLoaded = false;
let invoiceOptions = { schedules: [], schoolYears: [], classes: [] };
let invoicesLoaded = false;
let invoicesData = [];
let invoiceDetailCache = new Map();
let selectedInvoiceId = "";
let paymentReconciliationLoaded = false;
let paymentReconciliationData = { providers: [], schools: [], schoolYears: [], classes: [], invoices: [], transactions: [], intents: {}, matches: {}, summary: {} };
let paymentReconSelection = { type: "", id: "" };
let notificationLoaded = false;
let notificationOptions = { templates: [], campaigns: [], schoolYears: [], classes: [] };
let notificationPreviewData = { recipients: [], summary: {}, campaign: null, logs: [] };
let currentNotificationCampaignId = "";
let selectedNotificationRecipientKey = "";
let selectedNotificationRecipientIds = new Set();
let notificationCronData = null;
let adminOptions = { schools: [], schoolYears: [], classes: [] };
let adminDashboardLoaded = false;
let adminDashboardData = null;
let adminReportsLoaded = false;
let adminReportProviders = [];
let adminReportsData = null;
let operationsLoaded = false;
let operationsData = { operationLogs: [], auditLogs: [], operationSummary: {}, auditSummary: {} };
let selectedOperationLog = { type: "", id: "" };
let adminUsersLoaded = false;
let adminUsersData = { users: [], roles: [], permissions: [] };
let authSession = null;
let refreshAuthPromise = null;
let appContext = { schoolId: "", schoolYearId: "", periodCode: "", month: "" };
let appContextApplyTimer = 0;
let activeDialogRestore = null;
let activeDialogOnClose = null;
let lastDialogTrigger = null;
let interactiveRowsScheduled = false;
const interactiveRowSelector = [
  "tr[data-master-student-row]",
  "tr[data-invoice-row]",
  "tr[data-recon-invoice-row]",
  "tr[data-recon-transaction-row]",
  "tr[data-recipient-key]",
  "tr[data-operation-log-row]",
  "tr[data-audit-log-row]",
].join(", ");

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
    title: "Việc cần xử lý",
    description: "Theo dõi công nợ, giao dịch cần đối soát và các bước vận hành hằng ngày.",
    breadcrumbs: ["Tổng quan", "Việc cần xử lý"],
  },
  masterDataTab: {
    kicker: "Thiết lập",
    title: "Học sinh & phụ huynh",
    description: "Quản lý trường, năm học, lớp, học sinh, phụ huynh và dữ liệu nhận phí.",
    breadcrumbs: ["Thiết lập", "Học sinh & phụ huynh"],
  },
  feeTemplateTab: {
    kicker: "Học phí",
    title: "Bảng phí theo kỳ",
    description: "Thiết lập biểu phí, phụ phí và preview trước khi sinh hóa đơn.",
    breadcrumbs: ["Học phí", "Bảng phí"],
  },
  invoiceTab: {
    kicker: "Học phí",
    title: "Hóa đơn",
    description: "Sinh, kiểm tra và xuất hóa đơn/PDF receipt từ bảng phí đã lưu.",
    breadcrumbs: ["Học phí", "Hóa đơn"],
  },
  reconciliationTab: {
    kicker: "Thu tiền",
    title: "Đối soát thanh toán",
    description: "Theo dõi intent, giao dịch, tiền mặt và trạng thái đối soát hóa đơn.",
    breadcrumbs: ["Thu tiền", "Đối soát"],
  },
  paymentsTab: {
    kicker: "Thu tiền",
    title: "Công cụ QR/import",
    description: "Công cụ phụ cho batch thanh toán legacy, sinh QR và kiểm tra payload.",
    breadcrumbs: ["Thu tiền", "Công cụ QR/import"],
  },
  notificationTab: {
    kicker: "Liên lạc",
    title: "Thông báo học phí",
    description: "Tạo campaign, preview người nhận và theo dõi log gửi thông báo.",
    breadcrumbs: ["Liên lạc", "Thông báo"],
  },
  emailTab: {
    kicker: "Liên lạc",
    title: "Email & Cron",
    description: "Cấu hình provider email, preview/dry-run và quản lý lịch gửi cục bộ.",
    breadcrumbs: ["Liên lạc", "Email & Cron"],
  },
  reportsTab: {
    kicker: "Báo cáo & vận hành",
    title: "Báo cáo công nợ",
    description: "Xem và export báo cáo lớp, hóa đơn và giao dịch thanh toán.",
    breadcrumbs: ["Báo cáo & vận hành", "Báo cáo"],
  },
  operationsTab: {
    kicker: "Báo cáo & vận hành",
    title: "Vận hành",
    description: "Kiểm tra operational logs, audit logs và các lỗi nền cần xử lý.",
    breadcrumbs: ["Báo cáo & vận hành", "Vận hành"],
  },
  usersTab: {
    kicker: "Thiết lập",
    title: "Người dùng & quyền",
    description: "Quản lý user, role và permission trước khi bật enforcement đầy đủ.",
    breadcrumbs: ["Thiết lập", "Người dùng & quyền"],
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

const sessionRecoveryStorageKey = "abcsun.sessionRecovery.v1";
const permissionSummaryGroups = [
  { key: "view", label: "View", verbs: ["view", "read"] },
  { key: "create", label: "Create", verbs: ["create"] },
  { key: "update", label: "Update", verbs: ["update", "write", "manage"] },
  { key: "send", label: "Send", verbs: ["send"] },
  { key: "reconcile", label: "Reconcile", verbs: ["reconcile"] },
  { key: "export", label: "Export", verbs: ["export"] },
  { key: "administer", label: "Administer", verbs: ["assign_role", "administer"] },
];

const nativeFetch = window.fetch.bind(window);
window.fetch = authAwareFetch;

init();

function muiIcon(name) {
  return `<span class="mui-icon" aria-hidden="true">${escapeHtml(name)}</span>`;
}

function openAppDialog({ title, kicker = "Dialog", icon = "", nodes = [], content = null, actions = [], size = "md", onClose = null } = {}) {
  if (!appDialogEl) return;
  lastDialogTrigger = document.activeElement instanceof HTMLElement ? document.activeElement : null;
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
  appDialogEl.removeAttribute("aria-busy");
  appDialogEl.removeAttribute("aria-describedby");
  delete appDialogEl.dataset.busy;

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
  const actionLabel = label || "OK";
  button.type = "button";
  button.className = variant;
  button.innerHTML = `${icon ? muiIcon(icon) : ""}<span>${escapeHtml(actionLabel)}</span>`;
  button.addEventListener("click", async () => {
    clearDialogError();
    const previous = button.innerHTML;
    const previousMinWidth = button.style.minWidth;
    let shouldClose = closeOnSuccess;
    setDialogActionsBusy(button, true);
    if (onClick) {
      try {
        const result = await onClick();
        if (result === false) {
          shouldClose = false;
        }
      } catch (error) {
        showDialogError(error?.message || "Không xử lý được thao tác");
        shouldClose = false;
      }
    }
    setDialogActionsBusy(button, false);
    button.innerHTML = previous;
    button.style.minWidth = previousMinWidth;
    if (shouldClose) closeAppDialog();
  });
  return button;
}

function setDialogActionsBusy(activeButton, busy) {
  if (!appDialogEl || !appDialogActionsEl) return;
  const controls = [appDialogCloseBtn, ...appDialogActionsEl.querySelectorAll("button")].filter(Boolean);
  if (busy) {
    appDialogEl.dataset.busy = "true";
    appDialogEl.setAttribute("aria-busy", "true");
    controls.forEach((control) => {
      control.dataset.dialogWasDisabled = control.disabled ? "true" : "false";
      control.disabled = true;
    });
    if (activeButton) {
      const width = Math.ceil(activeButton.getBoundingClientRect().width);
      if (width) activeButton.style.minWidth = `${width}px`;
      activeButton.disabled = true;
      activeButton.setAttribute("aria-busy", "true");
      activeButton.innerHTML = `${muiIcon("progress_activity")}<span>Đang xử lý</span>`;
    }
    return;
  }
  delete appDialogEl.dataset.busy;
  appDialogEl.removeAttribute("aria-busy");
  controls.forEach((control) => {
    control.disabled = control.dataset.dialogWasDisabled === "true";
    delete control.dataset.dialogWasDisabled;
  });
  activeButton?.removeAttribute("aria-busy");
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
  appDialogEl?.setAttribute("aria-describedby", "appDialogError");
  window.setTimeout(() => appDialogErrorEl.focus(), 0);
}

function clearDialogError() {
  appDialogErrorEl.hidden = true;
  appDialogErrorEl.textContent = "";
  appDialogEl?.removeAttribute("aria-describedby");
}

function scheduleEnhanceInteractiveRows() {
  if (interactiveRowsScheduled) return;
  interactiveRowsScheduled = true;
  window.requestAnimationFrame(() => {
    interactiveRowsScheduled = false;
    enhanceInteractiveRows();
  });
}

function enhanceInteractiveRows(root = document) {
  root.querySelectorAll(interactiveRowSelector).forEach((row) => {
    row.tabIndex = 0;
    row.setAttribute("aria-selected", row.classList.contains("is-selected") ? "true" : "false");
  });
}

function activateFocusedInteractiveRow(event) {
  if (!["Enter", " "].includes(event.key)) return;
  const row = event.target.closest?.(interactiveRowSelector);
  if (!row || event.target !== row) return;
  event.preventDefault();
  row.click();
}

function confirmDialog({
  title,
  message,
  confirmLabel = "Xác nhận",
  confirmIcon = "check",
  danger = false,
  details = [],
  actor = false,
  auditNote = "",
} = {}) {
  return new Promise((resolve) => {
    let settled = false;
    const detailRows = [
      ...(actor ? [{ label: "Actor", value: currentActorLabel() }] : []),
      ...details.filter((item) => item?.label || item?.value),
      ...(auditNote ? [{ label: "Audit", value: auditNote }] : []),
    ];
    const detailsHtml = detailRows.length
      ? `<dl class="dialog-risk-list">${detailRows
          .map(
            (item) => `
              <div>
                <dt>${escapeHtml(item.label || "")}</dt>
                <dd>${escapeHtml(item.value || "-")}</dd>
              </div>
            `,
          )
          .join("")}</dl>`
      : "";
    openAppDialog({
      title,
      kicker: "Confirm",
      icon: danger ? "warning" : "help",
      size: "sm",
      content: `<div class="dialog-message${danger ? " dialog-message-danger" : ""}">${escapeHtml(message || "")}${detailsHtml}</div>`,
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

function currentActorLabel() {
  const user = authSession?.user || {};
  return user.displayName || user.email || user.phone || "Current session";
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
  const recovery = readSessionRecoveryState();
  if (recovery?.context) {
    appContext = {
      schoolId: recovery.context.schoolId || "",
      schoolYearId: recovery.context.schoolYearId || "",
      periodCode: recovery.context.periodCode || "",
      month: recovery.context.month || "",
    };
  }
  loginScreenEl.hidden = true;
  appShellEl.hidden = false;
  loginPasswordEl.value = "";
  setLoginStatus("", "");
  updateAuthBadge(session);
  applyPermissionUI();
  activateInitialAllowedTab(recovery?.activeTabId || "");
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

function activeTabId() {
  return tabPanels.find((panel) => panel.classList.contains("active"))?.id || "dashboardTab";
}

function updateVisibleMenuGroups() {
  document.querySelectorAll(".menu-group").forEach((group) => {
    group.hidden = !group.querySelector(".tab-button:not([hidden])");
  });
}

function updateLegacyPaymentActions(targetId = activeTabId()) {
  setElementAllowed(legacyPaymentActionsEl, targetId === "paymentsTab" && hasPermission("payment.create"));
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
  setElementAllowed(exportInvoiceCsvBtn, hasPermission("report.export"));
  setElementAllowed(openNotificationDialogBtn, hasPermission("notification.create") || hasPermission("notification.send"));
  setElementAllowed(previewNotificationsBtn, hasPermission("notification.send"));
  setElementAllowed(saveNotificationCampaignBtn, hasPermission("notification.create"));
  setElementAllowed(retryNotificationRecipientsBtn, hasPermission("notification.send"));
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
  setElementAllowed(openNotificationCronConfigBtn, hasPermission("email_cron.view") || hasPermission("email_cron.update"));
  setElementAllowed(saveCronBtn, hasPermission("email_cron.update"));
  setElementAllowed(disableCronBtn, hasPermission("email_cron.update"));
  setElementAllowed(runCronNowBtn, hasPermission("email_cron.update"));
  updateVisibleMenuGroups();
  updateLegacyPaymentActions();
  renderAdminQuickActions();
}

function activateInitialAllowedTab(preferredTabId = "") {
  const current = tabPanels.find((panel) => panel.classList.contains("active"))?.id || "dashboardTab";
  const targetId = canUseTab(preferredTabId)
    ? preferredTabId
    : canUseTab(current)
      ? current
      : tabButtons.find((button) => !button.hidden)?.dataset.tabTarget || "";
  if (!targetId) {
    currentSectionKickerEl.textContent = "Không đủ quyền";
    currentSectionTitleEl.textContent = "Chưa có màn hình được cấp quyền";
    currentSectionDescriptionEl.textContent = "Liên hệ quản trị viên để cập nhật role.";
    if (appBreadcrumbsEl) appBreadcrumbsEl.textContent = "Không đủ quyền";
    updateLegacyPaymentActions("");
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
  persistSessionRecoveryState(targetId);
}

function readSessionRecoveryState() {
  try {
    return JSON.parse(window.localStorage.getItem(sessionRecoveryStorageKey) || "null");
  } catch {
    return null;
  }
}

function persistSessionRecoveryState(tabId = activeTabId()) {
  try {
    window.localStorage.setItem(
      sessionRecoveryStorageKey,
      JSON.stringify({
        activeTabId: tabId || "dashboardTab",
        context: appContext,
        savedAt: new Date().toISOString(),
      }),
    );
  } catch {
    // Session recovery is best-effort only.
  }
}

function clearSessionRecoveryState() {
  try {
    window.localStorage.removeItem(sessionRecoveryStorageKey);
  } catch {
    // Ignore storage failures.
  }
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
  clearSessionRecoveryState();
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
  await applyAppContextToTab(targetId, true);
  persistSessionRecoveryState(targetId);
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
  renderBreadcrumbs(metadata);
  updateLegacyPaymentActions(targetId);
  renderAppContextControls();
}

function renderBreadcrumbs(metadata) {
  if (!appBreadcrumbsEl) return;
  const parts = metadata.breadcrumbs?.length ? metadata.breadcrumbs : [metadata.kicker, metadata.title].filter(Boolean);
  appBreadcrumbsEl.innerHTML = parts
    .map((part, index) => {
      const current = index === parts.length - 1 ? ` aria-current="page"` : "";
      return `<span${current}>${escapeHtml(part)}</span>`;
    })
    .join(`<span class="breadcrumb-separator">/</span>`);
}

function uniqueContextOptions(items) {
  const byId = new Map();
  items.forEach((item) => {
    const id = item?.id || "";
    if (!id || byId.has(id)) return;
    byId.set(id, item);
  });
  return [...byId.values()];
}

function contextSchoolOptions() {
  return uniqueContextOptions([
    ...(masterDataOptions.schools || []),
    ...(schoolTreeData.schools || []),
    ...(feeScheduleOptions.schools || []),
  ]).sort((a, b) => String(a.code || a.name || "").localeCompare(String(b.code || b.name || ""), "vi", { numeric: true }));
}

function contextYearOptions() {
  const explicitYears = [
    ...(masterDataOptions.schoolYears || []),
    ...(adminOptions.schoolYears || []),
    ...(feeScheduleOptions.schoolYears || []),
    ...(invoiceOptions.schoolYears || []),
    ...(notificationOptions.schoolYears || []),
  ].map((item) => ({
    id: item.id || item.schoolYearId || "",
    code: item.code || item.schoolYearCode || item.name || "",
    name: item.name || "",
    schoolId: item.schoolId || "",
    schoolCode: item.schoolCode || "",
  }));
  const yearsFromClasses = [
    ...(masterDataOptions.classes || []),
    ...(adminOptions.classes || []),
    ...(feeScheduleOptions.classes || []),
    ...(invoiceOptions.classes || []),
    ...(notificationOptions.classes || []),
  ].map((item) => ({
    id: item.schoolYearId || "",
    code: item.schoolYearCode || "",
    name: item.schoolYearCode || "",
    schoolId: item.schoolId || "",
    schoolCode: item.schoolCode || "",
  }));
  return uniqueContextOptions([...explicitYears, ...yearsFromClasses])
    .filter((item) => !appContext.schoolId || !item.schoolId || item.schoolId === appContext.schoolId)
    .sort((a, b) => String(a.code || a.name || "").localeCompare(String(b.code || b.name || ""), "vi", { numeric: true }));
}

function renderAppContextControls() {
  if (!appContextSchoolEl || !appContextYearEl) return;
  const schools = contextSchoolOptions();
  appContextSchoolEl.innerHTML = [
    `<option value="">${schools.length ? "Tất cả trường" : "ABC SUN"}</option>`,
    ...schools.map((item) => {
      const label = [item.code, item.name && item.name !== item.code ? item.name : ""].filter(Boolean).join(" · ");
      return `<option value="${escapeAttr(item.id)}">${escapeHtml(label || item.id)}</option>`;
    }),
  ].join("");
  appContextSchoolEl.value = optionValueOrEmpty(appContextSchoolEl, appContext.schoolId);

  const years = contextYearOptions();
  appContextYearEl.innerHTML = [
    `<option value="">Tất cả năm học</option>`,
    ...years.map((item) => {
      const label = [item.schoolCode, item.code || item.name].filter(Boolean).join(" · ");
      return `<option value="${escapeAttr(item.id)}">${escapeHtml(label || item.id)}</option>`;
    }),
  ].join("");
  appContextYearEl.value = optionValueOrEmpty(appContextYearEl, appContext.schoolYearId);
  appContextPeriodEl.value = appContext.periodCode || "";
  appContextMonthEl.value = appContext.month || "";
}

function syncAppContextFromActiveTab(targetId = activeTabId()) {
  const next = readTabContext(targetId);
  if (Object.prototype.hasOwnProperty.call(next, "schoolId")) appContext.schoolId = next.schoolId || "";
  if (Object.prototype.hasOwnProperty.call(next, "schoolYearId")) appContext.schoolYearId = next.schoolYearId || "";
  if (Object.prototype.hasOwnProperty.call(next, "periodCode")) appContext.periodCode = next.periodCode || "";
  if (Object.prototype.hasOwnProperty.call(next, "month")) appContext.month = next.month || "";
  renderAppContextControls();
  persistSessionRecoveryState(targetId);
}

function readTabContext(targetId) {
  if (targetId === "dashboardTab") {
    return {
      schoolId: adminDashboardSchoolEl.value,
      schoolYearId: adminDashboardYearEl.value,
      periodCode: adminDashboardPeriodEl.value.trim(),
      month: adminDashboardMonthEl.value,
    };
  }
  if (targetId === "reportsTab") {
    return {
      schoolId: adminReportsSchoolEl.value,
      schoolYearId: adminReportsYearEl.value,
      periodCode: adminReportsPeriodEl.value.trim(),
      month: adminReportsMonthEl.value,
    };
  }
  if (targetId === "masterDataTab") {
    return {
      schoolId: masterSchoolFilterEl.value,
      schoolYearId: masterSchoolYearFilterEl.value,
    };
  }
  if (targetId === "feeTemplateTab") {
    return {
      schoolId: feeScheduleSchoolEl.value,
      schoolYearId: feeScheduleYearEl.value,
      periodCode: feeSchedulePeriodEl.value.trim(),
      month: feeScheduleMonthEl.value,
    };
  }
  if (targetId === "notificationTab") {
    return {
      schoolYearId: notificationSchoolYearEl.value,
      periodCode: notificationPeriodEl.value.trim(),
    };
  }
  if (targetId === "reconciliationTab") {
    return {
      schoolId: paymentReconSchoolFilterEl.value,
      schoolYearId: paymentReconYearFilterEl.value,
      periodCode: paymentReconPeriodFilterEl.value.trim(),
    };
  }
  if (targetId === "paymentsTab") {
    return { periodCode: paymentPeriodEl.value.trim() };
  }
  return {};
}

async function applyAppContextToActiveTab() {
  appContext = {
    schoolId: appContextSchoolEl.value || "",
    schoolYearId: appContextYearEl.value || "",
    periodCode: appContextPeriodEl.value.trim(),
    month: appContextMonthEl.value || "",
  };
  await applyAppContextToTab(activeTabId(), true);
  persistSessionRecoveryState();
}

function scheduleApplyAppContext() {
  window.clearTimeout(appContextApplyTimer);
  appContextApplyTimer = window.setTimeout(() => {
    applyAppContextToActiveTab();
  }, 250);
}

async function applyAppContextToTab(targetId, reloadData = false) {
  if (targetId === "dashboardTab") {
    adminDashboardSchoolEl.value = optionValueOrEmpty(adminDashboardSchoolEl, appContext.schoolId);
    adminDashboardYearEl.value = optionValueOrEmpty(adminDashboardYearEl, appContext.schoolYearId);
    adminDashboardPeriodEl.value = appContext.periodCode;
    adminDashboardMonthEl.value = appContext.month;
    renderAdminFilters("dashboard");
    if (reloadData) await loadAdminDashboard(true);
  } else if (targetId === "reportsTab") {
    adminReportsSchoolEl.value = optionValueOrEmpty(adminReportsSchoolEl, appContext.schoolId);
    adminReportsYearEl.value = optionValueOrEmpty(adminReportsYearEl, appContext.schoolYearId);
    adminReportsPeriodEl.value = appContext.periodCode;
    adminReportsMonthEl.value = appContext.month;
    renderAdminFilters("reports");
    if (reloadData) await loadAdminReports(true);
  } else if (targetId === "masterDataTab") {
    masterSchoolFilterEl.value = optionValueOrEmpty(masterSchoolFilterEl, appContext.schoolId);
    masterSchoolYearFilterEl.value = optionValueOrEmpty(masterSchoolYearFilterEl, appContext.schoolYearId);
    renderMasterFilters();
    if (reloadData) {
      await loadSchoolTree();
      await loadMasterStudents();
    }
  } else if (targetId === "feeTemplateTab") {
    feeScheduleSchoolEl.value = optionValueOrEmpty(feeScheduleSchoolEl, appContext.schoolId);
    feeScheduleYearEl.value = optionValueOrEmpty(feeScheduleYearEl, appContext.schoolYearId);
    feeSchedulePeriodEl.value = appContext.periodCode;
    feeScheduleMonthEl.value = appContext.month;
    renderFeeScheduleControls();
    if (reloadData) await loadFeeScheduleList();
  } else if (targetId === "notificationTab") {
    notificationSchoolYearEl.value = optionValueOrEmpty(notificationSchoolYearEl, appContext.schoolYearId);
    notificationPeriodEl.value = appContext.periodCode;
    renderNotificationGradeOptions();
    renderNotificationClassOptions();
  } else if (targetId === "reconciliationTab") {
    paymentReconSchoolFilterEl.value = optionValueOrEmpty(paymentReconSchoolFilterEl, appContext.schoolId);
    paymentReconYearFilterEl.value = optionValueOrEmpty(paymentReconYearFilterEl, appContext.schoolYearId);
    paymentReconPeriodFilterEl.value = appContext.periodCode;
    renderPaymentReconFilters(paymentReconciliationData);
    if (reloadData) await loadPaymentReconciliation(true);
  } else if (targetId === "paymentsTab") {
    paymentPeriodEl.value = appContext.periodCode;
  }
  renderAppContextControls();
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
      <button class="remove-fee icon-only danger" type="button" aria-label="Xóa khoản phí" title="Xóa khoản phí">${muiIcon("remove_circle")}<span class="sr-only">Xóa</span></button>
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
  const confirmed = await confirmDialog({
    title: "Import bảng thanh toán?",
    message: "Thao tác này sẽ thay thế các dòng QR/import hiện tại bằng dữ liệu từ file đã map.",
    confirmLabel: "Import vào bảng",
    confirmIcon: "publish",
    danger: true,
    actor: true,
    details: [
      { label: "File", value: paymentImportState.file?.name || "-" },
      { label: "Preview", value: `${paymentImportState.preview?.length || 0} dòng` },
      { label: "Scope", value: "Legacy QR/import table, không ghi production invoice" },
    ],
  });
  if (!confirmed) return;
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
  renderAppContextControls();
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
  const params = new URLSearchParams();
  if (appContext.periodCode) params.set("periodCode", appContext.periodCode);
  if (appContext.month) params.set("month", appContext.month);
  const query = params.toString();
  const res = await fetch(`/api/v1/school-tree${query ? `?${query}` : ""}`);
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
  renderAppContextControls();
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
        <span class="school-tree-node-body">
          <strong>${escapeHtml(school.name || school.code || "-")}</strong>
          <small>${escapeHtml(school.code || "")} · ${Number(school.studentCount || 0)} HS · ${Number(school.classCount || 0)} lớp</small>
          ${renderSchoolTreeBadges(school)}
        </span>
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
        <span class="school-tree-node-body">
          <strong>${escapeHtml(year.code || "-")}</strong>
          <small>${Number(year.classCount || 0)} lớp · ${Number(year.studentCount || 0)} HS · ${Number(year.adjustmentCount || 0)} điều chỉnh</small>
          ${renderSchoolTreeBadges(year)}
        </span>
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
        <span class="school-tree-node-body">
          <strong>Khối ${escapeHtml(grade.grade || "-")}</strong>
          <small>${Number(grade.classCount || 0)} lớp · ${Number(grade.studentCount || 0)} HS</small>
          ${renderSchoolTreeBadges(grade)}
        </span>
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
      <span class="school-tree-node-body">
        <strong>${escapeHtml(item.name || "-")}</strong>
        <small>${Number(item.studentCount || 0)} HS · ${Number(item.activeFeeScheduleCount || 0)} active · ${Number(item.adjustmentCount || 0)} điều chỉnh</small>
        ${renderSchoolTreeBadges(item)}
      </span>
    </button>
  `;
}

function renderSchoolTreeBadges(subject = {}) {
  const studentCount = Number(subject.studentCount || 0);
  const readyCount = Number(subject.billingReadyStudentCount || 0);
  const activeScheduleCount = Number(subject.currentActiveScheduleCount || 0);
  const scheduleCount = Number(subject.currentFeeScheduleCount || 0);
  const invoiceCount = Number(subject.currentInvoiceCount || 0);
  const issueCount = Number(subject.issueCount || 0);
  const issueBadge = issueCount > 0
    ? `<span class="tree-badge tree-badge-warning">${muiIcon("priority_high")}<span>${issueCount}</span></span>`
    : `<span class="tree-badge tree-badge-ready">${muiIcon("check")}<span>OK</span></span>`;
  return `
    <span class="school-tree-badges">
      ${issueBadge}
      <span class="tree-badge">${muiIcon("alternate_email")}<span>${readyCount}/${studentCount}</span></span>
      <span class="tree-badge">${muiIcon("price_change")}<span>${activeScheduleCount || scheduleCount}</span></span>
      <span class="tree-badge">${muiIcon("request_quote")}<span>${invoiceCount}</span></span>
    </span>
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
  syncAppContextFromActiveTab("masterDataTab");
  renderSchoolTree();
  await loadMasterStudents();
  renderSchoolTreeDetail(findSchoolTreeNode(type, id));
}

function applySchoolTreeFilters(node) {
  if (!node) return;
  masterSchoolFilterEl.value = optionValueOrEmpty(masterSchoolFilterEl, node.school?.id || "");
  renderMasterFilters();
  masterSchoolYearFilterEl.value = optionValueOrEmpty(masterSchoolYearFilterEl, node.year?.id || node.classItem?.schoolYearId || "");
  renderMasterFilters();
  masterGradeFilterEl.value = optionValueOrEmpty(masterGradeFilterEl, node.grade?.grade || node.classItem?.grade || "");
  renderMasterClassFilter();
  masterClassFilterEl.value = optionValueOrEmpty(masterClassFilterEl, node.classItem?.id || "");
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

function schoolTreeSubject(node) {
  return node?.classItem || node?.grade || node?.year || node?.school || {};
}

function schoolTreeTypeLabel(type) {
  return {
    school: "Trường",
    year: "Năm học",
    grade: "Khối",
    class: "Lớp",
  }[type] || "School tree";
}

function schoolTreePeriodLabel() {
  return [appContext.periodCode, appContext.month ? `T${appContext.month}` : ""].filter(Boolean).join(" · ") || "Tất cả kỳ";
}

function schoolTreeScopeLabel(node) {
  if (!node) return "";
  return [
    node.school?.code || node.school?.name || "",
    node.year?.code || "",
    node.grade?.grade ? `Khối ${node.grade.grade}` : node.classItem?.grade ? `Khối ${node.classItem.grade}` : "",
    node.classItem?.name || "",
  ].filter(Boolean).join(" · ");
}

function renderSchoolTreeMetrics(subject = {}) {
  const studentCount = Number(subject.studentCount || 0);
  const billingReady = Number(subject.billingReadyStudentCount || 0);
  const activeSchedules = Number(subject.currentActiveScheduleCount || 0);
  const schedules = Number(subject.currentFeeScheduleCount || 0);
  const invoices = Number(subject.currentInvoiceCount || 0);
  const openInvoices = Number(subject.openInvoiceCount || 0);
  const issues = Number(subject.issueCount || 0);
  const items = [
    { icon: "groups", label: "Học sinh", value: studentCount },
    { icon: "alternate_email", label: "Nhận phí", value: `${billingReady}/${studentCount}` },
    { icon: "price_change", label: "Bảng phí kỳ", value: `${activeSchedules}/${schedules}` },
    { icon: "request_quote", label: "Hóa đơn kỳ", value: `${invoices}/${studentCount}` },
    { icon: "pending_actions", label: "Đang mở", value: openInvoices },
    { icon: issues ? "priority_high" : "check_circle", label: "Cần xử lý", value: issues },
  ];
  return `
    <div class="school-tree-metrics">
      ${items.map((item) => `
        <div class="school-tree-metric">
          ${muiIcon(item.icon)}
          <span>${escapeHtml(item.label)}</span>
          <strong>${escapeHtml(String(item.value))}</strong>
        </div>
      `).join("")}
    </div>
  `;
}

function schoolTreeReadinessItems(subject = {}) {
  const studentCount = Number(subject.studentCount || 0);
  const billingReady = Number(subject.billingReadyStudentCount || 0);
  const missingBilling = Number(subject.missingBillingRecipientCount || 0);
  const activeSchedules = Number(subject.currentActiveScheduleCount || 0);
  const schedules = Number(subject.currentFeeScheduleCount || 0);
  const invoices = Number(subject.currentInvoiceCount || 0);
  const missingInvoices = Math.max(0, studentCount - invoices);
  const openInvoices = Number(subject.openInvoiceCount || 0);
  return [
    {
      tone: missingBilling > 0 ? "warning" : "ready",
      icon: missingBilling > 0 ? "mark_email_unread" : "mark_email_read",
      label: "Người nhận phí",
      value: missingBilling > 0 ? `${missingBilling} HS thiếu email nhận phí` : `${billingReady}/${studentCount} HS sẵn sàng`,
    },
    {
      tone: studentCount > 0 && activeSchedules === 0 ? "warning" : "ready",
      icon: activeSchedules > 0 ? "price_check" : "price_change",
      label: "Bảng phí theo kỳ",
      value: activeSchedules > 0 ? `${activeSchedules} active / ${schedules} phù hợp` : "Chưa có bảng phí active",
    },
    {
      tone: missingInvoices > 0 ? "warning" : "ready",
      icon: missingInvoices > 0 ? "request_quote" : "fact_check",
      label: "Hóa đơn theo kỳ",
      value: missingInvoices > 0 ? `${missingInvoices} HS chưa có hóa đơn` : `${invoices}/${studentCount} hóa đơn`,
    },
    {
      tone: openInvoices > 0 ? "info" : "ready",
      icon: openInvoices > 0 ? "pending_actions" : "done_all",
      label: "Công nợ đang mở",
      value: openInvoices > 0 ? `${openInvoices} hóa đơn unpaid/partial/review` : "Không có hóa đơn đang mở",
    },
  ];
}

function renderSchoolTreeReadiness(subject = {}) {
  return `
    <ul class="school-tree-readiness-list">
      ${schoolTreeReadinessItems(subject).map((item) => `
        <li class="readiness-item readiness-${escapeAttr(item.tone)}">
          ${muiIcon(item.icon)}
          <span>${escapeHtml(item.label)}</span>
          <strong>${escapeHtml(item.value)}</strong>
        </li>
      `).join("")}
    </ul>
  `;
}

function studentHasBillingRecipient(student) {
  return (student?.parents || []).some((parent) => parent.receivesBillingEmail && parent.isActive && parent.emailActive && parent.email);
}

function renderSchoolTreeRosterPreview() {
  const students = masterStudentsData || [];
  const visible = students.slice(0, 6);
  if (!visible.length) {
    return `<div class="school-tree-roster-empty">Không có học sinh trong phạm vi đang chọn</div>`;
  }
  return `
    <div class="school-tree-roster">
      ${visible.map((student) => {
        const billingReady = studentHasBillingRecipient(student);
        return `
          <div class="school-tree-roster-row">
            <strong>${escapeHtml(student.studentName || "-")}</strong>
            <span>${escapeHtml([student.studentCode, student.className].filter(Boolean).join(" · ") || "-")}</span>
            <small class="${billingReady ? "ready" : "warning"}">${billingReady ? "Có email nhận phí" : "Thiếu email nhận phí"}</small>
          </div>
        `;
      }).join("")}
      ${students.length > visible.length ? `<div class="school-tree-roster-more">+${students.length - visible.length} học sinh khác</div>` : ""}
    </div>
  `;
}

function renderSchoolTreeActions(node) {
  const studentAction = hasPermission("student.view")
    ? `<button data-tree-open-students="true" type="button">${muiIcon("groups")}<span>Xem học sinh</span></button>`
    : "";
  const feeAction = hasPermission("fee.view") && (node.year || node.grade || node.classItem)
    ? `<button data-tree-open-fees="true" type="button">${muiIcon("price_change")}<span>Thiết lập bảng phí</span></button>`
    : "";
  const invoiceAction = (hasPermission("invoice.view") || hasPermission("invoice.create")) && (node.year || node.grade || node.classItem)
    ? `<button data-tree-open-invoices="true" type="button">${muiIcon("request_quote")}<span>Sinh hóa đơn</span></button>`
    : "";
  return `<div class="detail-actions school-tree-action-row">${studentAction}${feeAction}${invoiceAction}</div>`;
}

function renderSchoolTreeDetail(node) {
  if (!node) {
    schoolTreeDetailEl.innerHTML = `<div class="detail-placeholder">${muiIcon("account_tree")}<span>Chọn trường, năm học, khối hoặc lớp.</span></div>`;
    return;
  }
  const subject = schoolTreeSubject(node);
  const title = node.classItem?.name || (node.grade ? `Khối ${node.grade.grade}` : node.year?.code || node.school?.name || "-");
  const scope = schoolTreeScopeLabel(node);
  const latest = node.classItem?.latestFeeScheduleId
    ? `<div class="detail-section"><p class="detail-section-title">Bảng phí gần nhất</p><div class="detail-grid"><span>Kỳ</span><strong>${escapeHtml(node.classItem.latestPeriodCode || "-")}</strong><span>Tên</span><strong>${escapeHtml(node.classItem.latestFeeScheduleName || "-")}</strong><span>Trạng thái</span><strong>${escapeHtml(node.classItem.latestScheduleStatus || "-")}</strong></div></div>`
    : "";
  schoolTreeDetailEl.innerHTML = `
    <div class="detail-hero">
      ${muiIcon(node.classItem ? "school" : node.grade ? "stacked_line_chart" : node.year ? "event" : "apartment")}
      <div>
        <strong>${escapeHtml(title)}</strong>
        <span>${escapeHtml([schoolTreeTypeLabel(node.type), scope, schoolTreePeriodLabel()].filter(Boolean).join(" · "))}</span>
      </div>
    </div>
    ${renderSchoolTreeMetrics(subject)}
    <div class="detail-section">
      <p class="detail-section-title">Readiness kỳ đang chọn</p>
      ${renderSchoolTreeReadiness(subject)}
    </div>
    ${latest}
    <div class="detail-section">
      <p class="detail-section-title">Học sinh trong phạm vi</p>
      ${renderSchoolTreeRosterPreview()}
    </div>
    ${renderSchoolTreeActions(node)}
  `;
  schoolTreeDetailEl.querySelector("[data-tree-open-students]")?.addEventListener("click", () => openSchoolTreeStudentScope(node));
  schoolTreeDetailEl.querySelector("[data-tree-open-fees]")?.addEventListener("click", () => openSchoolTreeFeeScope(node));
  schoolTreeDetailEl.querySelector("[data-tree-open-invoices]")?.addEventListener("click", () => openSchoolTreeInvoiceScope(node));
}

async function openSchoolTreeStudentScope(node) {
  if (!hasPermission("student.view")) {
    setMasterStatus("Không đủ quyền xem học sinh", "error");
    return;
  }
  applySchoolTreeFilters(node);
  syncAppContextFromActiveTab("masterDataTab");
  await loadMasterStudents();
  document.querySelector(".master-list-panel")?.scrollIntoView({ block: "start", behavior: "smooth" });
  setMasterStatus("Đã lọc học sinh theo cây trường", "ready");
}

async function openSchoolTreeFeeScope(node) {
  if (!hasPermission("fee.view")) {
    setMasterStatus("Không đủ quyền mở bảng phí", "error");
    return;
  }
  await activateTab("feeTemplateTab");
  await loadFeeSchedules(true);
  feeScheduleSchoolEl.value = optionValueOrEmpty(feeScheduleSchoolEl, node.school?.id || "");
  feeScheduleYearEl.value = optionValueOrEmpty(feeScheduleYearEl, node.year?.id || node.classItem?.schoolYearId || "");
  feeSchedulePeriodEl.value = appContext.periodCode || feeSchedulePeriodEl.value;
  feeScheduleMonthEl.value = appContext.month || feeScheduleMonthEl.value;
  renderFeeScheduleControls();
  feeScheduleGradeEl.value = optionValueOrEmpty(feeScheduleGradeEl, node.grade?.grade || node.classItem?.grade || "");
  renderFeeScheduleClassFilter();
  feeScheduleClassEl.value = optionValueOrEmpty(feeScheduleClassEl, node.classItem?.id || "");
  await loadFeeScheduleList();
  if (hasPermission("fee.update") && !openFeeScheduleDialogBtn.hidden) {
    openFeeScheduleDialog();
  }
}

function findSchoolTreeInvoiceSchedule(node) {
  const targetYearId = node.year?.id || node.classItem?.schoolYearId || "";
  const targetGrade = node.grade?.grade || node.classItem?.grade || "";
  const targetClassId = node.classItem?.id || "";
  return [...(invoiceOptions.schedules || [])]
    .map((schedule) => {
      if (targetYearId && schedule.schoolYearId !== targetYearId) return { schedule, score: -1 };
      if (appContext.periodCode && String(schedule.periodCode || "").toLowerCase() !== appContext.periodCode.toLowerCase()) return { schedule, score: -1 };
      if (appContext.month && Number(schedule.month || 0) !== Number(appContext.month)) return { schedule, score: -1 };
      let score = -1;
      if (targetClassId && schedule.classId === targetClassId) {
        score = 30;
      } else if (targetGrade && schedule.scopeType === "grade" && String(schedule.grade || "").toLowerCase() === targetGrade.toLowerCase()) {
        score = 20;
      } else if (schedule.scopeType === "school_year") {
        score = 10;
      }
      if (score < 0) return { schedule, score };
      if (schedule.status === "active") score += 5;
      if (schedule.periodCode) score += 1;
      return { schedule, score };
    })
    .filter((item) => item.score >= 0)
    .sort((a, b) => b.score - a.score)[0]?.schedule || null;
}

async function openSchoolTreeInvoiceScope(node) {
  if (!hasPermission("invoice.view") && !hasPermission("invoice.create")) {
    setMasterStatus("Không đủ quyền mở hóa đơn", "error");
    return;
  }
  await activateTab("invoiceTab");
  await loadInvoices(true);
  const schedule = findSchoolTreeInvoiceSchedule(node);
  if (!schedule) {
    setInvoiceStatus("Chưa có bảng phí phù hợp với phạm vi và kỳ đang chọn", "error");
    return;
  }
  invoiceScheduleEl.value = optionValueOrEmpty(invoiceScheduleEl, schedule.id);
  if (!openInvoiceDialogBtn.hidden) {
    openInvoiceDialog();
  }
  setInvoiceStatus("Đã chọn bảng phí phù hợp", "ready");
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
      school: adminReportsSchoolEl,
      year: adminReportsYearEl,
      grade: adminReportsGradeEl,
      classEl: adminReportsClassEl,
      period: adminReportsPeriodEl,
      month: adminReportsMonthEl,
      status: adminReportsInvoiceStatusEl,
      provider: adminReportsProviderEl,
    };
  }
  return {
    school: adminDashboardSchoolEl,
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
  const selectedSchool = elements.school.value;
  const selectedYear = elements.year.value;
  const selectedGrade = elements.grade.value;
  const selectedClass = elements.classEl.value;
  elements.school.innerHTML = [
    `<option value="">Tất cả trường</option>`,
    ...(adminOptions.schools || []).map((item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.code || item.name || item.id)}</option>`),
  ].join("");
  elements.school.value = optionValueOrEmpty(elements.school, selectedSchool);

  elements.year.innerHTML = [
    `<option value="">Tất cả năm học</option>`,
    ...(adminOptions.schoolYears || [])
      .filter((item) => !elements.school.value || item.schoolId === elements.school.value)
      .map((item) => {
        const label = [item.schoolCode, item.code].filter(Boolean).join(" · ");
        return `<option value="${escapeAttr(item.id)}">${escapeHtml(label || item.code || item.id)}</option>`;
      }),
  ].join("");
  elements.year.value = optionValueOrEmpty(elements.year, selectedYear);

  const grades = [...new Set(
    (adminOptions.classes || [])
      .filter((item) => !elements.school.value || item.schoolId === elements.school.value)
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
    if (elements.school.value && item.schoolId !== elements.school.value) return false;
    if (elements.year.value && item.schoolYearId !== elements.year.value) return false;
    if (elements.grade.value && item.grade !== elements.grade.value) return false;
    return true;
  });
  elements.classEl.innerHTML = [
    `<option value="">Tất cả lớp</option>`,
    ...classes.map((item) => `<option value="${escapeAttr(item.id)}">${escapeHtml(item.schoolYearCode)} · ${escapeHtml(item.name)}</option>`),
  ].join("");
  elements.classEl.value = optionValueOrEmpty(elements.classEl, selectedClass);
  if (elements.provider) {
    renderAdminReportProviderFilter();
  }
  renderAppContextControls();
}

function renderAdminReportProviderFilter() {
  if (!adminReportsProviderEl) return;
  const selected = adminReportsProviderEl.value;
  adminReportsProviderEl.innerHTML = [
    `<option value="">Tất cả provider</option>`,
    ...adminReportProviders.map((provider) => {
      const suffix = provider.configured ? "" : " · thiếu cấu hình";
      return `<option value="${escapeAttr(provider.code)}">${escapeHtml(provider.displayName || provider.code)}${suffix}</option>`;
    }),
  ].join("");
  adminReportsProviderEl.value = optionValueOrEmpty(adminReportsProviderEl, selected);
}

function adminFilterParams(kind) {
  const elements = adminFilterElements(kind);
  const params = new URLSearchParams();
  if (elements.school.value) params.set("schoolId", elements.school.value);
  if (elements.year.value) params.set("schoolYearId", elements.year.value);
  if (elements.grade.value) params.set("grade", elements.grade.value);
  if (elements.classEl.value) params.set("classId", elements.classEl.value);
  if (elements.period.value.trim()) params.set("periodCode", elements.period.value.trim());
  if (elements.month.value) params.set("month", elements.month.value);
  if (elements.status.value) params.set("status", elements.status.value);
  if (elements.provider?.value) params.set("provider", elements.provider.value);
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
  adminDashboardData = data || null;
  renderAdminMetrics(adminDashboardMetricsEl, data?.summary || null);
  renderOperatorOnboarding(data);
  renderAdminWorkQueue(data);
  renderAdminQuickActions();
  renderAdminReadiness(data?.readiness || null);
  renderAdminTopClasses(data?.topClasses || []);
  renderAdminAttentionInvoices(data?.attentionInvoices || []);
}

function renderOperatorOnboarding(data) {
  if (!operatorOnboardingEl) return;
  const issues = data?.readiness?.issues || [];
  const issueTypes = new Set(issues.map((issue) => issue.type).filter(Boolean));
  const hasDashboardData = !!data?.summary;
  const hasSchools = (adminOptions.schools || []).length > 0 || (masterDataOptions.schools || []).length > 0 || (schoolTreeData.schools || []).length > 0;
  const hasYears = (adminOptions.schoolYears || []).length > 0 || (masterDataOptions.schoolYears || []).length > 0 || (feeScheduleOptions.schoolYears || []).length > 0;
  const hasClasses = (adminOptions.classes || []).length > 0 || (masterDataOptions.classes || []).length > 0 || (feeScheduleOptions.classes || []).length > 0;
  const billingIssueTypes = [
    "student_missing_billing_recipient",
    "billing_recipient_missing_email",
    "billing_recipient_email_inactive",
    "inactive_parent_selected_for_billing",
  ];
  const studentIssueTypes = ["student_missing_parent", "student_missing_class", "class_has_no_students"];
  const feeIssueTypes = ["class_missing_fee_schedule", "fee_schedule_empty_items", "fee_schedule_zero_amount_items"];
  const steps = [
    {
      id: "school",
      label: "Tạo hoặc kiểm tra trường",
      detail: "Trường là scope gốc cho năm học, lớp, học sinh và báo cáo.",
      status: hasSchools ? "ready" : "warning",
      statusLabel: hasSchools ? "Ready" : "Cần tạo",
      action: "students",
    },
    {
      id: "classes",
      label: "Tạo năm học và lớp",
      detail: "Năm học/cohort và lớp cần sẵn sàng trước import hoặc lập bảng phí.",
      status: hasYears && hasClasses ? "ready" : "warning",
      statusLabel: hasYears && hasClasses ? "Ready" : "Thiếu dữ liệu",
      action: "students",
    },
    {
      id: "students",
      label: "Import học sinh và phụ huynh",
      detail: "Student code là định danh bền vững; phụ huynh cần được liên kết rõ.",
      status: hasAnyIssue(issueTypes, studentIssueTypes) ? "warning" : hasDashboardData ? "ready" : "info",
      statusLabel: hasAnyIssue(issueTypes, studentIssueTypes) ? "Cần xử lý" : hasDashboardData ? "Ready" : "Chưa tải",
      action: "students",
    },
    {
      id: "billing",
      label: "Resolve billing recipients",
      detail: "Mỗi học sinh cần người nhận phí active, nhận billing, email active và có email.",
      status: hasAnyIssue(issueTypes, billingIssueTypes) ? "error" : hasDashboardData ? "ready" : "info",
      statusLabel: hasAnyIssue(issueTypes, billingIssueTypes) ? "Blocking" : hasDashboardData ? "Ready" : "Chưa tải",
      action: "students",
    },
    {
      id: "email",
      label: "Cấu hình email provider",
      detail: "Chỉ lưu masked state trong UI; không hiển thị app password/API key thật.",
      status: issueTypes.has("email_provider_not_configured") ? "warning" : hasDashboardData ? "ready" : "info",
      statusLabel: issueTypes.has("email_provider_not_configured") ? "Chưa cấu hình" : hasDashboardData ? "Ready" : "Chưa tải",
      action: "email_config",
    },
    {
      id: "fees",
      label: "Tạo bảng phí đầu tiên",
      detail: "Preview scope, item, adjustment và billing readiness trước khi lưu active.",
      status: hasAnyIssue(issueTypes, feeIssueTypes) ? "warning" : hasDashboardData ? "ready" : "info",
      statusLabel: hasAnyIssue(issueTypes, feeIssueTypes) ? "Cần bảng phí" : hasDashboardData ? "Ready" : "Chưa tải",
      action: "fees",
    },
    {
      id: "invoice",
      label: "Preview invoice batch đầu tiên",
      detail: "Preview giúp thấy blocking issue trước khi ghi invoice vào database.",
      status: Number(data?.summary?.receivableAmount || 0) > 0 ? "ready" : hasDashboardData ? "info" : "warning",
      statusLabel: Number(data?.summary?.receivableAmount || 0) > 0 ? "Ready" : hasDashboardData ? "Next" : "Cần tải",
      action: "invoices",
    },
  ];
  operatorOnboardingEl.innerHTML = steps.map(renderOperatorOnboardingStep).join("");
}

function hasAnyIssue(issueTypes, types) {
  return types.some((type) => issueTypes.has(type));
}

function renderOperatorOnboardingStep(step, index) {
  const action = dashboardActionDefinitions().find((item) => item.id === step.action);
  const canOpen = action && dashboardActionAllowed(action);
  const tone = step.status === "error" ? "error" : step.status === "warning" ? "warning" : step.status === "ready" ? "ready" : "busy";
  return `
    <div class="operator-onboarding-item" data-status="${escapeAttr(step.status)}">
      <span class="operator-onboarding-index">${index + 1}</span>
      <span class="operator-onboarding-copy">
        <strong>${escapeHtml(step.label)}</strong>
        <small>${escapeHtml(step.detail)}</small>
      </span>
      <span class="status-pill" data-tone="${escapeAttr(tone)}">${escapeHtml(step.statusLabel)}</span>
      ${
        canOpen
          ? `<button type="button" data-dashboard-action="${escapeAttr(step.action)}">${muiIcon("open_in_new")}<span>Mở</span></button>`
          : `<span class="status-pill">Không đủ quyền</span>`
      }
    </div>
  `;
}

function dashboardActionDefinitions() {
  return [
    {
      id: "students",
      label: "Cập nhật học sinh",
      detail: "Kiểm tra lớp, phụ huynh và người nhận phí",
      icon: "groups",
      tabId: "masterDataTab",
      permission: "student.view",
    },
    {
      id: "fees",
      label: "Lập bảng phí",
      detail: "Mở cấu hình phí theo năm học, khối, lớp và kỳ thu",
      icon: "format_list_bulleted",
      tabId: "feeTemplateTab",
      permission: "fee.update",
      afterOpen: () => {
        if (!openFeeScheduleDialogBtn.hidden) openFeeScheduleDialog();
      },
    },
    {
      id: "invoices",
      label: "Sinh hóa đơn",
      detail: "Preview và sinh hóa đơn từ bảng phí đã lưu",
      icon: "request_quote",
      tabId: "invoiceTab",
      permission: "invoice.create",
      afterOpen: () => {
        if (!openInvoiceDialogBtn.hidden) openInvoiceDialog();
      },
    },
    {
      id: "notify",
      label: "Gửi thông báo",
      detail: "Dry-run campaign trước khi gửi email thật",
      icon: "campaign",
      tabId: "notificationTab",
      permission: "notification.send",
      afterOpen: () => {
        if (!openNotificationDialogBtn.hidden) openNotificationDialog();
      },
    },
    {
      id: "reconcile",
      label: "Đối soát giao dịch",
      detail: "Kiểm tra tiền vào, tiền mặt và trạng thái invoice",
      icon: "account_balance_wallet",
      tabId: "reconciliationTab",
      permission: "payment.view",
    },
    {
      id: "email_config",
      label: "Email & cron",
      detail: "Kiểm tra provider, quota và queue gửi email",
      icon: "mark_email_read",
      tabId: "emailTab",
      permission: "email_config.view",
    },
    {
      id: "operations",
      label: "Log vận hành",
      detail: "Mở operation logs để xem lỗi nền gần đây",
      icon: "monitor_heart",
      tabId: "operationsTab",
      permission: "operation_log.view",
      secondary: true,
    },
    {
      id: "qr_tool",
      label: "Công cụ QR/import",
      detail: "Mở batch QR legacy khi cần kiểm tra payload riêng lẻ",
      icon: "qr_code_scanner",
      tabId: "paymentsTab",
      permission: "payment.create",
      secondary: true,
    },
  ];
}

function dashboardActionAllowed(action) {
  return canUseTab(action.tabId) && (!action.permission || hasPermission(action.permission));
}

function renderAdminQuickActions() {
  if (!adminQuickActionsEl) return;
  const actions = dashboardActionDefinitions().filter(dashboardActionAllowed);
  if (!actions.length) {
    adminQuickActionsEl.innerHTML = `<div class="empty-task">Chưa có tác vụ được cấp quyền</div>`;
    return;
  }
  adminQuickActionsEl.innerHTML = actions
    .map(
      (action) => `
        <button class="${action.secondary ? "secondary-task-action" : ""}" type="button" data-dashboard-action="${escapeAttr(action.id)}">
          ${muiIcon(action.icon)}
          <span>
            <strong>${escapeHtml(action.label)}</strong>
            <small>${escapeHtml(action.detail)}</small>
          </span>
        </button>
      `,
    )
    .join("");
}

function renderAdminWorkQueue(data) {
  if (!adminWorkQueueEl) return;
  if (!data?.summary) {
    adminWorkQueueEl.innerHTML = `<div class="empty-task">Chưa có dữ liệu dashboard</div>`;
    return;
  }
  const summary = data.summary || {};
  const items = [
    {
      count: Number(summary.unmatchedTransactionCount || 0) + Number(summary.manualReviewCount || 0),
      label: "Giao dịch cần đối soát",
      detail: "Giao dịch chưa khớp hoặc invoice đang manual review",
      action: "reconcile",
      permission: "payment.view",
    },
    {
      count: Number(summary.unpaidStudentCount || 0),
      label: "Hóa đơn chưa thu",
      detail: "Học sinh còn invoice unpaid trong bộ lọc hiện tại",
      action: "invoices",
      permission: "invoice.view",
    },
    {
      count: Number(summary.partialPaymentCount || 0),
      label: "Thanh toán một phần",
      detail: "Invoice cần kiểm tra tiếp trước khi chốt thu",
      action: "reconcile",
      permission: "payment.view",
    },
    {
      count: Number(summary.overpaidManualReviewCount || 0),
      label: "Overpaid hoặc cần review",
      detail: "Khoản thu cần xử lý thủ công",
      action: "reconcile",
      permission: "payment.view",
    },
  ].filter((item) => item.count > 0);

  if (!items.length) {
    adminWorkQueueEl.innerHTML = `<div class="empty-task">Không có việc cần xử lý trong bộ lọc hiện tại</div>`;
    return;
  }
  adminWorkQueueEl.innerHTML = items
    .map((item) => {
      const canOpen = hasPermission(item.permission);
      return `
        <div class="work-queue-item">
          <span class="work-queue-count">${Number(item.count)}</span>
          <span class="work-queue-copy">
            <strong>${escapeHtml(item.label)}</strong>
            <small>${escapeHtml(item.detail)}</small>
          </span>
          ${
            canOpen
              ? `<button type="button" data-dashboard-action="${escapeAttr(item.action)}">${muiIcon("open_in_new")}<span>Mở</span></button>`
              : `<span class="status-pill">Không đủ quyền</span>`
          }
        </div>
      `;
    })
    .join("");
}

function renderAdminReadiness(readiness) {
  if (!adminReadinessCenterEl) return;
  if (!readiness?.summary) {
    adminReadinessCenterEl.innerHTML = `<div class="empty-task">Chưa có dữ liệu readiness</div>`;
    renderAdminReadinessTypeFilter([]);
    return;
  }
  const issues = readiness.issues || [];
  renderAdminReadinessTypeFilter(issues);
  const severityFilter = adminReadinessSeverityEl?.value || "";
  const typeFilter = adminReadinessTypeEl?.value || "";
  const filtered = issues.filter((issue) => {
    if (severityFilter && issue.severity !== severityFilter) return false;
    if (typeFilter && issue.type !== typeFilter) return false;
    return true;
  });
  const summary = readiness.summary || {};
  const summaryHtml = `
    <div class="readiness-summary-grid">
      ${renderReadinessSummaryCard("blocking", "Blocking", summary.blockingCount || 0)}
      ${renderReadinessSummaryCard("warning", "Warning", summary.warningCount || 0)}
      ${renderReadinessSummaryCard("info", "Info", summary.infoCount || 0)}
      ${renderReadinessSummaryCard("total", "Tổng", summary.totalCount || 0)}
    </div>
  `;
  if (!filtered.length) {
    adminReadinessCenterEl.innerHTML = `${summaryHtml}<div class="empty-task">Không có readiness issue theo bộ lọc hiện tại</div>`;
    return;
  }
  const visible = filtered.slice(0, 120);
  const hiddenCount = Math.max(0, filtered.length - visible.length);
  const groups = ["blocking", "warning", "info"]
    .map((severity) => {
      const rows = visible.filter((issue) => issue.severity === severity);
      if (!rows.length) return "";
      return `
        <div class="readiness-group" data-severity="${escapeAttr(severity)}">
          <div class="readiness-group-title">
            <span>${escapeHtml(readinessSeverityLabel(severity))}</span>
            <strong>${rows.length}</strong>
          </div>
          <div class="readiness-issue-list">
            ${rows.map(renderReadinessIssue).join("")}
          </div>
        </div>
      `;
    })
    .join("");
  const moreHtml = hiddenCount
    ? `<div class="readiness-more">Đang hiển thị 120 issue đầu tiên, còn ${hiddenCount} issue trong bộ lọc.</div>`
    : "";
  adminReadinessCenterEl.innerHTML = `${summaryHtml}${groups}${moreHtml}`;
}

function renderAdminReadinessTypeFilter(issues) {
  if (!adminReadinessTypeEl) return;
  const selected = adminReadinessTypeEl.value;
  const types = [...new Set((issues || []).map((issue) => issue.type).filter(Boolean))].sort((a, b) =>
    readinessTypeLabel(a).localeCompare(readinessTypeLabel(b), "vi", { numeric: true }),
  );
  adminReadinessTypeEl.innerHTML = [
    `<option value="">Tất cả issue</option>`,
    ...types.map((type) => `<option value="${escapeAttr(type)}">${escapeHtml(readinessTypeLabel(type))}</option>`),
  ].join("");
  adminReadinessTypeEl.value = optionValueOrEmpty(adminReadinessTypeEl, selected);
}

function renderReadinessSummaryCard(severity, label, count) {
  return `
    <div class="readiness-summary-card" data-severity="${escapeAttr(severity)}">
      <strong>${Number(count || 0)}</strong>
      <span>${escapeHtml(label)}</span>
    </div>
  `;
}

function renderReadinessIssue(issue) {
  const action = dashboardActionDefinitions().find((item) => item.id === issue.action);
  const canOpen = action && dashboardActionAllowed(action);
  const count = Number(issue.referenceCount || 0);
  const countHtml = count > 1 ? `<span class="readiness-count">x${count}</span>` : "";
  return `
    <div class="readiness-issue" data-severity="${escapeAttr(issue.severity || "warning")}">
      <span class="readiness-severity-chip">${escapeHtml(readinessSeverityLabel(issue.severity))}</span>
      <span class="readiness-copy">
        <strong>${escapeHtml(issue.entityLabel || "-")}</strong>
        <small>${escapeHtml([issue.scope, issue.message].filter(Boolean).join(" · "))}</small>
      </span>
      <span class="readiness-type">${escapeHtml(readinessTypeLabel(issue.type))}</span>
      ${countHtml}
      ${
        canOpen
          ? `<button type="button" data-dashboard-action="${escapeAttr(issue.action || "")}" data-readiness-issue="${escapeAttr(issue.id || "")}">${muiIcon("open_in_new")}<span>Mở</span></button>`
          : `<span class="status-pill">Không đủ quyền</span>`
      }
    </div>
  `;
}

function readinessSeverityLabel(severity) {
  switch (severity) {
    case "blocking":
      return "Blocking";
    case "info":
      return "Info";
    default:
      return "Warning";
  }
}

function readinessTypeLabel(type) {
  const labels = {
    billing_recipient_email_inactive: "Email inactive",
    billing_recipient_missing_email: "Thiếu email nhận phí",
    class_has_no_students: "Lớp chưa có học sinh",
    class_missing_fee_schedule: "Thiếu bảng phí",
    cron_over_limit_sends: "Cron vượt quota",
    cron_queue_errors: "Cron có lỗi",
    cron_recent_errors: "Lỗi cron gần đây",
    email_provider_not_configured: "Email chưa cấu hình",
    fee_schedule_empty_items: "Bảng phí trống",
    fee_schedule_zero_amount_items: "Dòng phí bằng 0",
    inactive_parent_selected_for_billing: "Quan hệ inactive nhận phí",
    incoming_transaction_unmatched: "Giao dịch chưa khớp",
    invoice_manual_review: "Invoice manual review",
    invoice_missing_payment_data: "Invoice thiếu QR",
    invoice_overpaid: "Invoice overpaid",
    invoice_partial: "Invoice partial",
    invoice_unpaid: "Invoice unpaid",
    notification_failed_recipients: "Thông báo gửi lỗi",
    student_adjustment_missing_reason: "Điều chỉnh thiếu lý do",
    student_missing_billing_recipient: "Thiếu người nhận phí",
    student_missing_class: "Học sinh thiếu lớp",
    student_missing_parent: "Thiếu phụ huynh",
    transaction_manual_review: "Giao dịch cần review",
  };
  return labels[type] || String(type || "Issue").replaceAll("_", " ");
}

function readinessIssueById(issueId) {
  return (adminDashboardData?.readiness?.issues || []).find((issue) => issue.id === issueId) || null;
}

async function runDashboardAction(actionId, readinessIssueId = "") {
  const action = dashboardActionDefinitions().find((item) => item.id === actionId);
  if (!action || !dashboardActionAllowed(action)) return;
  const issue = readinessIssueById(readinessIssueId);
  applyReadinessIssueContext(issue);
  await activateTab(action.tabId);
  await applyReadinessIssueTargetFilters(action.tabId, issue);
  if (action.afterOpen) {
    window.setTimeout(action.afterOpen, 0);
  }
}

function applyReadinessIssueContext(issue) {
  if (!issue) return;
  appContext = {
    ...appContext,
    schoolId: issue.schoolId || appContext.schoolId,
    schoolYearId: issue.schoolYearId || appContext.schoolYearId,
    periodCode: issue.periodCode || appContext.periodCode,
    month: issue.month ? String(issue.month) : appContext.month,
  };
  renderAppContextControls();
}

async function applyReadinessIssueTargetFilters(tabId, issue) {
  if (!issue) return;
  if (tabId === "masterDataTab") {
    masterSchoolFilterEl.value = optionValueOrEmpty(masterSchoolFilterEl, issue.schoolId || appContext.schoolId);
    masterSchoolYearFilterEl.value = optionValueOrEmpty(masterSchoolYearFilterEl, issue.schoolYearId || appContext.schoolYearId);
    masterGradeFilterEl.value = optionValueOrEmpty(masterGradeFilterEl, issue.grade || "");
    renderMasterFilters();
    masterClassFilterEl.value = optionValueOrEmpty(masterClassFilterEl, issue.classId || "");
    if (issue.entityType === "student" && issue.referenceCode) {
      masterSearchEl.value = issue.referenceCode;
    }
    await loadMasterStudents();
    if (issue.entityType === "student" && issue.entityId) {
      selectMasterStudent(issue.entityId);
    }
  } else if (tabId === "feeTemplateTab") {
    feeScheduleSchoolEl.value = optionValueOrEmpty(feeScheduleSchoolEl, issue.schoolId || appContext.schoolId);
    feeScheduleYearEl.value = optionValueOrEmpty(feeScheduleYearEl, issue.schoolYearId || appContext.schoolYearId);
    feeSchedulePeriodEl.value = issue.periodCode || appContext.periodCode || feeSchedulePeriodEl.value;
    feeScheduleMonthEl.value = issue.month ? String(issue.month) : appContext.month || feeScheduleMonthEl.value;
    renderFeeScheduleControls();
    feeScheduleGradeEl.value = optionValueOrEmpty(feeScheduleGradeEl, issue.grade || "");
    renderFeeScheduleClassFilter(issue.classId || "");
    await loadFeeScheduleList();
  } else if (tabId === "reconciliationTab") {
    paymentReconSchoolFilterEl.value = optionValueOrEmpty(paymentReconSchoolFilterEl, issue.schoolId || appContext.schoolId);
    paymentReconYearFilterEl.value = optionValueOrEmpty(paymentReconYearFilterEl, issue.schoolYearId || appContext.schoolYearId);
    paymentReconPeriodFilterEl.value = issue.periodCode || appContext.periodCode || paymentReconPeriodFilterEl.value;
    renderPaymentReconFilters(paymentReconciliationData);
    paymentReconGradeFilterEl.value = optionValueOrEmpty(paymentReconGradeFilterEl, issue.grade || "");
    renderPaymentReconFilters(paymentReconciliationData);
    paymentReconClassFilterEl.value = optionValueOrEmpty(paymentReconClassFilterEl, issue.classId || "");
    if (issue.type === "incoming_transaction_unmatched") {
      paymentTransactionStatusFilterEl.value = optionValueOrEmpty(paymentTransactionStatusFilterEl, "unmatched");
    } else if (issue.type === "transaction_manual_review") {
      paymentTransactionStatusFilterEl.value = optionValueOrEmpty(paymentTransactionStatusFilterEl, "manual_review");
    }
    await loadPaymentReconciliation(true);
  } else if (tabId === "notificationTab") {
    notificationSchoolYearEl.value = optionValueOrEmpty(notificationSchoolYearEl, issue.schoolYearId || appContext.schoolYearId);
    notificationPeriodEl.value = issue.periodCode || appContext.periodCode || notificationPeriodEl.value;
    renderNotificationGradeOptions();
    notificationGradeEl.value = optionValueOrEmpty(notificationGradeEl, issue.grade || "");
    renderNotificationClassOptions();
    notificationClassEl.value = optionValueOrEmpty(notificationClassEl, issue.classId || "");
  } else if (tabId === "operationsTab") {
    operationSourceFilterEl.value = optionValueOrEmpty(operationSourceFilterEl, "background_job");
    operationLevelFilterEl.value = optionValueOrEmpty(operationLevelFilterEl, "error");
    await loadOperations(true);
  }
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
  adminReportProviders = data.providers || adminReportProviders;
  adminReportsData = data;
  adminReportsLoaded = true;
  renderAdminFilters("reports");
  renderAdminReports(data);
  setAdminStatus(adminReportsStatusEl, "Sẵn sàng", "ready");
}

function renderAdminReports(data) {
  adminReportsData = data || null;
  renderAdminMetrics(adminReportsSummaryEl, data?.summary || null);
  renderAdminReportClasses(data?.classRows || []);
  renderAdminReportInvoices(data?.invoiceRows || []);
  renderAdminReportTransactions(data?.transactions || []);
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
          <td class="money">${formatMoney(invoice.paidAmount || 0)}<small>Còn ${formatMoney(invoice.outstandingAmount || 0)}</small></td>
          <td><span class="tag ${paymentReconStatusTone(invoice.status)}">${escapeHtml(invoice.status || "")}</span><small>${Number(invoice.sentCount || 0)} sent · ${Number(invoice.paymentIntentCount || 0)} intent · ${Number(invoice.matchedPaymentCount || 0)} match</small></td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    adminReportInvoiceRowsEl.innerHTML = `<tr><td colspan="6" class="empty-cell">Chưa có hóa đơn trong bộ lọc</td></tr>`;
  }
}

function renderAdminReportTransactions(rows) {
  adminReportTransactionCountEl.textContent = `${rows.length} giao dịch`;
  adminReportTransactionRowsEl.innerHTML = rows
    .map(
      (transaction) => `
        <tr>
          <td><strong>${escapeHtml(transaction.provider || transaction.providerCode || "")}</strong><small>${escapeHtml(transaction.providerTransactionId || transaction.referenceCode || "")}</small></td>
          <td class="money">${formatMoney(transaction.amount || 0)}<small>${escapeHtml(formatDateTime(transaction.transactionTime))}</small></td>
          <td>${escapeHtml(transaction.matchReason || transaction.description || "-")}<small>${escapeHtml(paymentMatchMeta(transaction) || transaction.accountNumber || "")}</small></td>
          <td>${escapeHtml(transaction.invoiceCode || "Chưa match")}<small>${escapeHtml([transaction.studentCode || "", transaction.studentName || ""].filter(Boolean).join(" · "))}</small></td>
          <td><span class="tag ${paymentReconStatusTone(transaction.matchStatus || transaction.status)}">${escapeHtml(transaction.status || "")}</span></td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    adminReportTransactionRowsEl.innerHTML = `<tr><td colspan="5" class="empty-cell">Chưa có giao dịch trong bộ lọc</td></tr>`;
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
  if (operationNameFilterEl.value.trim()) operationParams.set("operation", operationNameFilterEl.value.trim());
  if (operationStatusFilterEl.value.trim()) operationParams.set("status", operationStatusFilterEl.value.trim());
  if (operationEntityTypeFilterEl.value.trim()) operationParams.set("entityType", operationEntityTypeFilterEl.value.trim());
  operationParams.set("limit", String(limit));
  const auditParams = new URLSearchParams({ limit: String(limit) });
  if (auditActionFilterEl.value.trim()) auditParams.set("action", auditActionFilterEl.value.trim());
  if (operationEntityTypeFilterEl.value.trim()) auditParams.set("entityType", operationEntityTypeFilterEl.value.trim());
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
  renderOperations({
    operationLogs: operationData.logs || [],
    auditLogs: auditData.logs || [],
    operationSummary: operationData.summary || {},
    auditSummary: auditData.summary || {},
  });
  setAdminStatus(operationsStatusEl, "Sẵn sàng", "ready");
}

function renderOperations(data) {
  operationsData = data || { operationLogs: [], auditLogs: [], operationSummary: {}, auditSummary: {} };
  renderOperationsSummary(operationsData.operationSummary || {}, operationsData.auditSummary || {});
  renderOperationLogs(data?.operationLogs || []);
  renderAuditLogs(data?.auditLogs || []);
  if (selectedOperationLog.id) {
    const rows = selectedOperationLog.type === "audit" ? operationsData.auditLogs : operationsData.operationLogs;
    const selected = (rows || []).find((row) => row.id === selectedOperationLog.id);
    if (selected) {
      renderOperationLogDetail(selectedOperationLog.type, selected);
      return;
    }
  }
  renderOperationLogDetail("", null);
}

function renderOperationsSummary(operationSummary = {}, auditSummary = {}) {
  if (!operationsSummaryEl) return;
  const cards = [
    { label: "Webhook errors", value: Number(operationSummary.webhookErrorCount || 0) },
    { label: "Email errors", value: Number(operationSummary.emailErrorCount || 0) },
    { label: "Cron errors", value: Number(operationSummary.cronErrorCount || 0) },
    { label: "Background failures", value: Number(operationSummary.backgroundJobErrorCount || 0) },
    { label: "Audit actions", value: Number(auditSummary.totalCount || 0) },
    { label: "Money actions", value: Number(auditSummary.moneyActionCount || 0) },
    { label: "Fee actions", value: Number(auditSummary.feeActionCount || 0) },
    { label: "User actions", value: Number(auditSummary.userActionCount || 0) },
  ];
  operationsSummaryEl.innerHTML = cards.map((card) => `<div><strong>${escapeHtml(card.value)}</strong><span>${escapeHtml(card.label)}</span></div>`).join("");
}

function renderOperationLogs(rows) {
  operationLogCountEl.textContent = `${rows.length} log`;
  operationLogRowsEl.innerHTML = rows
    .map(
      (row) => `
        <tr data-operation-log-row="${escapeAttr(row.id || "")}">
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
  operationLogRowsEl.querySelectorAll("[data-operation-log-row]").forEach((row) => {
    row.classList.toggle("is-selected", selectedOperationLog.type === "operation" && selectedOperationLog.id === row.dataset.operationLogRow);
    row.addEventListener("click", () => {
      const log = (operationsData.operationLogs || []).find((item) => item.id === row.dataset.operationLogRow);
      selectedOperationLog = { type: "operation", id: row.dataset.operationLogRow };
      renderOperationLogDetail("operation", log);
      updateOperationActiveRows();
    });
  });
}

function renderAuditLogs(rows) {
  auditLogCountEl.textContent = `${rows.length} log`;
  auditLogRowsEl.innerHTML = rows
    .map(
      (row) => `
        <tr data-audit-log-row="${escapeAttr(row.id || "")}">
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
  auditLogRowsEl.querySelectorAll("[data-audit-log-row]").forEach((row) => {
    row.classList.toggle("is-selected", selectedOperationLog.type === "audit" && selectedOperationLog.id === row.dataset.auditLogRow);
    row.addEventListener("click", () => {
      const log = (operationsData.auditLogs || []).find((item) => item.id === row.dataset.auditLogRow);
      selectedOperationLog = { type: "audit", id: row.dataset.auditLogRow };
      renderOperationLogDetail("audit", log);
      updateOperationActiveRows();
    });
  });
}

function metadataReason(metadata) {
  if (!metadata || typeof metadata !== "object") return "";
  return String(metadata.reason || "");
}

function updateOperationActiveRows() {
  operationLogRowsEl.querySelectorAll("[data-operation-log-row]").forEach((row) => {
    row.classList.toggle("is-selected", selectedOperationLog.type === "operation" && selectedOperationLog.id === row.dataset.operationLogRow);
  });
  auditLogRowsEl.querySelectorAll("[data-audit-log-row]").forEach((row) => {
    row.classList.toggle("is-selected", selectedOperationLog.type === "audit" && selectedOperationLog.id === row.dataset.auditLogRow);
  });
}

function renderOperationLogDetail(type, log) {
  if (!operationLogDetailEl) return;
  if (!log) {
    operationLogDetailEl.innerHTML = `<div class="detail-placeholder">${muiIcon("manage_search")}<span>Chọn một operation hoặc audit log để xem metadata đã lọc secret và mở workflow liên quan.</span></div>`;
    return;
  }
  const isAudit = type === "audit";
  const title = isAudit ? log.action || "Audit" : log.operation || "Operation";
  const status = isAudit ? log.entityType || "audit" : log.status || log.level || "operation";
  const message = isAudit ? log.reason || metadataReason(log.metadata) || "-" : log.message || "-";
  const metadata = JSON.stringify(log.metadata || {}, null, 2);
  const drilldown = operationDrilldownTarget(log);
  operationLogDetailEl.innerHTML = `
    <div class="detail-hero">
      ${muiIcon(isAudit ? "policy" : "running_with_errors")}
      <div>
        <strong>${escapeHtml(title)}</strong>
        <span>${escapeHtml(formatDateTime(log.occurredAt))}</span>
      </div>
    </div>
    <div class="reconciliation-detail-grid">
      <span>Type</span><strong>${escapeHtml(isAudit ? "Audit log" : "Operation log")}</strong>
      <span>Status</span><strong>${escapeHtml(status)}</strong>
      <span>Actor/source</span><strong>${escapeHtml(isAudit ? log.actorName || log.actorUserId || "-" : log.source || "-")}</strong>
      <span>Entity</span><strong>${escapeHtml([log.entityType || "", log.entityId || ""].filter(Boolean).join(" · ") || "-")}</strong>
      <span>Request</span><strong>${escapeHtml(log.requestId || "-")}</strong>
      <span>Message</span><strong>${escapeHtml(message)}</strong>
    </div>
    <div class="detail-section">
      <h3 class="detail-section-title">Sanitized metadata</h3>
      <textarea class="payload" readonly>${escapeHtml(metadata)}</textarea>
    </div>
    <div class="detail-actions">
      ${drilldown ? `<button type="button" data-operation-drilldown="${escapeAttr(drilldown.tabId)}">${muiIcon(drilldown.icon)}<span>${escapeHtml(drilldown.label)}</span></button>` : ""}
    </div>
  `;
  operationLogDetailEl.querySelectorAll("[data-operation-drilldown]").forEach((button) => {
    button.addEventListener("click", () => openOperationDrilldown(log));
  });
}

function operationDrilldownTarget(log) {
  const entityType = log?.entityType || "";
  if (["invoice", "manual_cash_receipt"].includes(entityType)) {
    return { tabId: "invoiceTab", label: "Open invoice", icon: "receipt_long" };
  }
  if (["payment_transaction", "provider_event", "reconciliation_match"].includes(entityType)) {
    return { tabId: "reconciliationTab", label: "Open reconciliation", icon: "account_balance_wallet" };
  }
  if (["notification_campaign", "notification_recipient", "notification_log"].includes(entityType)) {
    return { tabId: "notificationTab", label: "Open notifications", icon: "campaign" };
  }
  if (["student", "parent", "class"].includes(entityType)) {
    return { tabId: "masterDataTab", label: "Open master data", icon: "account_tree" };
  }
  if (["fee_schedule", "student_fee_adjustment"].includes(entityType)) {
    return { tabId: "feeTemplateTab", label: "Open fees", icon: "format_list_bulleted" };
  }
  return null;
}

async function openOperationDrilldown(log) {
  const target = operationDrilldownTarget(log);
  if (!target) return;
  await activateTab(target.tabId);
  if (target.tabId === "invoiceTab" && log.entityId) {
    selectedInvoiceId = log.entityId;
    await loadInvoices(true);
    selectInvoice(log.entityId);
  } else if (target.tabId === "reconciliationTab") {
    await loadPaymentReconciliation(true);
  } else if (target.tabId === "notificationTab") {
    await loadNotifications(true);
  } else if (target.tabId === "masterDataTab") {
    await loadMasterData(true);
  } else if (target.tabId === "feeTemplateTab") {
    await loadFeeSchedules(true);
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
      const summaries = rolePermissionSummaries(role.permissions || []);
      return `
        <div class="admin-role-item">
          <div>
            <strong>${escapeHtml(role.name || role.code || "")}</strong>
            <small>${escapeHtml(role.code || "")}${role.isSystem ? " · system" : ""}</small>
            <p>${escapeHtml(role.description || "")}</p>
            <p>${escapeHtml(roleTemplateSummary(role.code))}</p>
          </div>
          <div class="admin-permission-summary">
            ${summaries.map(renderPermissionSummaryGroup).join("") || `<span class="tag">no permissions</span>`}
          </div>
        </div>
      `;
    })
    .join("");
  if (!roles.length) {
    adminRoleListEl.textContent = "Chưa có role";
  }
}

function roleTemplateSummary(code = "") {
  switch (code) {
    case "admin":
      return "Admin: cấu hình hệ thống, user/role, tất cả workflow học phí và vận hành.";
    case "staff":
      return "Staff: quản lý học sinh/phụ huynh, lớp, thông báo và các bước chăm sóc vận hành.";
    case "accountant":
      return "Accountant: lập phí, hóa đơn, thu tiền, đối soát, báo cáo, email/cron và audit.";
    default:
      return "Custom role: kiểm tra nhóm permission trước khi gán cho operator.";
  }
}

function rolePermissionSummaries(permissions) {
  const byGroup = new Map(permissionSummaryGroups.map((group) => [group.key, { ...group, modules: new Set(), raw: [] }]));
  permissions.forEach((permission) => {
    const code = permission.code || "";
    const group = permissionSummaryGroup(code);
    if (!byGroup.has(group.key)) {
      byGroup.set(group.key, { ...group, modules: new Set(), raw: [] });
    }
    const entry = byGroup.get(group.key);
    entry.modules.add(permissionModuleLabel(code));
    entry.raw.push(code);
  });
  return [...byGroup.values()].filter((group) => group.modules.size || group.raw.length);
}

function permissionSummaryGroup(code) {
  const verb = String(code || "").split(".").pop() || "";
  return permissionSummaryGroups.find((group) => group.verbs.includes(verb)) || { key: "other", label: "Other", verbs: [] };
}

function permissionModuleLabel(code) {
  const module = String(code || "").split(".")[0] || "system";
  const labels = {
    audit_log: "audit",
    dashboard: "dashboard",
    email_config: "email config",
    email_cron: "email cron",
    fee: "fees",
    invoice: "invoices",
    notification: "notifications",
    operation_log: "operations",
    payment: "payments",
    report: "reports",
    school_tree: "school tree",
    student: "students",
    user: "users",
  };
  return labels[module] || module.replaceAll("_", " ");
}

function renderPermissionSummaryGroup(group) {
  const modules = [...group.modules].sort((a, b) => a.localeCompare(b, "vi", { numeric: true }));
  return `
    <div class="admin-permission-group">
      <span>${escapeHtml(group.label)}</span>
      <strong>${escapeHtml(modules.join(", ") || "-")}</strong>
    </div>
  `;
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
  renderSchoolTreeDetail(findSchoolTreeNode(selectedSchoolTreeNode.type, selectedSchoolTreeNode.id));
  setMasterStatus("Sẵn sàng", "ready");
}

function renderMasterStudents(students = []) {
  masterStudentsRawData = students || [];
  masterStudentsData = filterMasterStudentsForRelationshipView(masterStudentsRawData);
  if (!masterStudentsData.some((student) => masterStudentKey(student) === selectedMasterStudentKey)) {
    selectedMasterStudentKey = masterStudentsData[0] ? masterStudentKey(masterStudentsData[0]) : "";
  }
  masterStudentCountEl.textContent = masterStudentsData.length === masterStudentsRawData.length
    ? `${masterStudentsData.length} học sinh`
    : `${masterStudentsData.length}/${masterStudentsRawData.length} học sinh`;
  masterStudentsEl.innerHTML = masterStudentsData
    .map((student) => {
      const key = masterStudentKey(student);
      const warning = masterStudentContactWarning(student);
      const billingCount = masterStudentBillingCount(student);
      const parentCount = masterStudentParentCount(student);
      const invoiceAttention = Number(student.invoiceAttentionCount || 0);
      return `
        <tr data-master-student-row="${escapeAttr(key)}" class="${key === selectedMasterStudentKey ? "is-selected" : ""}">
          <td><strong>${escapeHtml(student.studentCode || "")}</strong></td>
          <td>${escapeHtml(student.studentName || "")}</td>
          <td>${escapeHtml([student.schoolCode, student.schoolYearCode].filter(Boolean).join(" · "))}</td>
          <td>${escapeHtml(student.grade || "")}</td>
          <td>${escapeHtml(student.className || "")}</td>
          <td>${renderMasterParentSummary(student)}</td>
          <td><span class="relationship-pill ${billingCount > 0 ? "is-ready" : "is-warning"}">${muiIcon(billingCount > 0 ? "mark_email_read" : "mark_email_unread")}<span>${billingCount}/${parentCount}</span></span></td>
          <td><span class="relationship-pill ${warning ? "is-warning" : "is-ready"}">${muiIcon(warning ? "priority_high" : "check_circle")}<span>${escapeHtml(masterStudentContactWarningLabel(warning))}</span></span></td>
          <td>${invoiceAttention > 0 ? `<span class="relationship-pill is-info">${muiIcon("request_quote")}<span>${invoiceAttention}</span></span>` : `<span class="muted-cell">-</span>`}</td>
          <td><span class="tag">${escapeHtml(student.status || "active")}</span></td>
        </tr>
      `;
    })
    .join("");
  if (!masterStudentsData.length) {
    masterStudentsEl.innerHTML = `<tr><td colspan="10" class="empty-cell">Chưa có dữ liệu học sinh</td></tr>`;
  }
  masterStudentsEl.querySelectorAll("[data-master-student-row]").forEach((row) => {
    row.addEventListener("click", () => selectMasterStudent(row.dataset.masterStudentRow));
  });
  const selected = masterStudentsData.find((student) => masterStudentKey(student) === selectedMasterStudentKey);
  renderMasterStudentDetail(selected);
  populateMasterStudentForm(selected);
}

function filterMasterStudentsForRelationshipView(students = []) {
  const filter = masterBillingFilterEl?.value || "";
  return (students || []).filter((student) => {
    const warning = masterStudentContactWarning(student);
    const billingCount = masterStudentBillingCount(student);
    const invoiceAttention = Number(student.invoiceAttentionCount || 0);
    if (filter === "ready") return billingCount > 0 && !warning;
    if (filter === "missing") return billingCount === 0;
    if (filter === "warning") return !!warning;
    if (filter === "attention") return invoiceAttention > 0;
    return true;
  });
}

function masterStudentParentCount(student) {
  return Number(student?.parentCount ?? (student?.parents || []).length);
}

function masterStudentBillingCount(student) {
  if (student && Object.prototype.hasOwnProperty.call(student, "billingRecipientCount")) {
    return Number(student.billingRecipientCount || 0);
  }
  return (student?.parents || []).filter(masterParentBillingReady).length;
}

function masterStudentContactWarning(student) {
  if (student?.contactWarning) return student.contactWarning;
  if (masterStudentParentCount(student) === 0) return "missing_parent";
  if (masterStudentBillingCount(student) === 0) return "missing_billing_recipient";
  return "";
}

function masterStudentContactWarningLabel(warning) {
  return {
    missing_parent: "Thiếu PH",
    missing_billing_recipient: "Thiếu billing",
  }[warning] || "OK";
}

function renderMasterParentSummary(student) {
  const parents = student?.parents || [];
  if (!parents.length) return `<span class="muted-cell">Chưa có</span>`;
  const primary = parents.find((parent) => parent.isPrimary) || parents[0];
  const extra = parents.length > 1 ? ` +${parents.length - 1}` : "";
  return `
    <div class="relationship-summary">
      <strong>${escapeHtml(primary.parentName || "-")}${escapeHtml(extra)}</strong>
      <span>${escapeHtml(masterRelationshipLabel(primary.relationship))}${primary.phone ? ` · ${escapeHtml(primary.phone)}` : ""}</span>
    </div>
  `;
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
    size: "xl",
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
  const parentCount = masterStudentParentCount(student);
  const billingCount = masterStudentBillingCount(student);
  const warning = masterStudentContactWarning(student);
  const invoiceAttention = Number(student.invoiceAttentionCount || 0);
  masterStudentDetailEl.innerHTML = `
    <div class="detail-hero">
      ${muiIcon("person")}
      <div>
        <strong>${escapeHtml(student.studentName || "-")}</strong>
        <span>${escapeHtml(student.studentCode || "Chưa có mã HS")}</span>
      </div>
    </div>
    <div class="student-relationship-metrics">
      <div>${muiIcon("supervisor_account")}<span>Phụ huynh</span><strong>${parentCount}</strong></div>
      <div>${muiIcon("alternate_email")}<span>Billing</span><strong>${billingCount}/${parentCount}</strong></div>
      <div>${muiIcon(warning ? "priority_high" : "check_circle")}<span>Contact</span><strong>${escapeHtml(masterStudentContactWarningLabel(warning))}</strong></div>
      <div>${muiIcon("request_quote")}<span>HĐ cần xử lý</span><strong>${invoiceAttention}</strong></div>
    </div>
    <div class="detail-grid">
      <span>Năm học</span><strong>${escapeHtml(student.schoolYearCode || "-")}</strong>
      <span>Trường</span><strong>${escapeHtml(student.schoolCode || "-")}</strong>
      <span>Khối / lớp</span><strong>${escapeHtml([student.grade ? `Khối ${student.grade}` : "", student.className || ""].filter(Boolean).join(" · ") || "-")}</strong>
      <span>Trạng thái</span><strong>${escapeHtml(student.status || "active")}</strong>
      <span>Quy tắc billing</span><strong>Active + nhận phí + email active + có email</strong>
    </div>
    <div class="detail-section">
      <p class="detail-section-title">Quan hệ phụ huynh</p>
      ${renderMasterParentRelationshipTable(student)}
    </div>
    <div class="detail-section">
      <p class="detail-section-title">Anh/chị em</p>
      ${renderMasterStudentSiblings(student)}
    </div>
    <div class="detail-actions">
      <button data-edit-master-student="true" type="button">${muiIcon("edit")}<span>Sửa học sinh</span></button>
      ${hasPermission("invoice.view") ? `<button data-master-student-invoices="true" type="button">${muiIcon("request_quote")}<span>Hóa đơn</span></button>` : ""}
      ${hasPermission("notification.view") || hasPermission("notification.send") ? `<button data-master-student-notifications="true" type="button">${muiIcon("campaign")}<span>Thông báo</span></button>` : ""}
      <button data-master-student-class="true" type="button">${muiIcon("account_tree")}<span>Lớp</span></button>
    </div>
  `;
  masterStudentDetailEl.querySelector("[data-edit-master-student]")?.addEventListener("click", openMasterStudentDialog);
  masterStudentDetailEl.querySelector("[data-master-student-invoices]")?.addEventListener("click", () => openMasterStudentInvoices(student));
  masterStudentDetailEl.querySelector("[data-master-student-notifications]")?.addEventListener("click", () => openMasterStudentNotifications(student));
  masterStudentDetailEl.querySelector("[data-master-student-class]")?.addEventListener("click", () => openMasterStudentClassScope(student));
}

function renderMasterParentRelationshipTable(student) {
  const parents = student.parents || [];
  if (!parents.length) {
    return `<div class="relationship-empty">Chưa có phụ huynh. Import hoặc cập nhật master data.</div>`;
  }
  return `
    <div class="relationship-table">
      <div class="relationship-table-head"><span>Quan hệ</span><span>Liên hệ</span><span>Flags</span><span>Billing</span></div>
      ${parents.map((parent) => {
        const ready = parent.billingReady ?? masterParentBillingReady(parent);
        return `
          <div class="relationship-row">
            <div><strong>${escapeHtml(parent.parentName || "-")}</strong><span>${escapeHtml(masterRelationshipLabel(parent.relationship))}</span></div>
            <div><strong>${escapeHtml(parent.email || "Chưa có email")}</strong><span>${escapeHtml(parent.phone || "Chưa có SĐT")}</span></div>
            <div class="relationship-flags">${masterParentFlagTags(parent)}</div>
            <div><span class="relationship-pill ${ready ? "is-ready" : "is-warning"}">${muiIcon(ready ? "check_circle" : "priority_high")}<span>${ready ? "Valid" : "Blocked"}</span></span></div>
          </div>
        `;
      }).join("")}
    </div>
  `;
}

function renderMasterStudentSiblings(student) {
  const siblings = student.siblings || [];
  if (!siblings.length) {
    return `<div class="relationship-empty">Chưa phát hiện học sinh dùng chung phụ huynh.</div>`;
  }
  return `
    <div class="sibling-list">
      ${siblings.map((sibling) => `
        <div class="sibling-item">
          ${muiIcon("family_restroom")}
          <div>
            <strong>${escapeHtml(sibling.studentCode || "")} · ${escapeHtml(sibling.studentName || "")}</strong>
            <span>${escapeHtml([sibling.schoolYearCode, sibling.grade ? `Khối ${sibling.grade}` : "", sibling.className].filter(Boolean).join(" · ") || "-")}</span>
            <small>${escapeHtml((sibling.sharedParentNames || []).join(", ") || "Chung phụ huynh")}</small>
          </div>
        </div>
      `).join("")}
    </div>
  `;
}

function masterParentBillingReady(parent) {
  return !!(parent?.isActive && parent?.receivesBillingEmail && parent?.emailActive && parent?.email);
}

function masterRelationshipLabel(value) {
  return {
    guardian: "Phụ huynh",
    mother: "Mẹ",
    me: "Mẹ",
    mom: "Mẹ",
    father: "Bố",
    bo: "Bố",
    dad: "Bố",
    grandparent: "Ông/bà",
    ong_ba: "Ông/bà",
    other: "Khác",
  }[value || ""] || value || "Phụ huynh";
}

function masterParentFlagTags(parent) {
  const tags = [
    parent.isPrimary ? ["Chính", "is-ready"] : ["Phụ", ""],
    parent.isActive ? ["Active", "is-ready"] : ["Inactive", "is-warning"],
    parent.receivesBillingEmail ? ["Nhận phí", "is-info"] : ["Không gửi", ""],
    parent.emailActive ? ["Email active", "is-ready"] : ["Email off", "is-warning"],
  ];
  return tags.map(([label, tone]) => `<span class="relationship-chip ${tone}">${escapeHtml(label)}</span>`).join("");
}

function masterRelationshipOptions(selected = "guardian") {
  const options = [
    ["guardian", "Phụ huynh"],
    ["mother", "Mẹ"],
    ["father", "Bố"],
    ["grandparent", "Ông/bà"],
    ["other", "Khác"],
  ];
  return options.map(([value, label]) => `<option value="${escapeAttr(value)}" ${value === selected ? "selected" : ""}>${escapeHtml(label)}</option>`).join("");
}

async function openMasterStudentInvoices(student) {
  if (!hasPermission("invoice.view")) {
    setMasterStatus("Không đủ quyền xem hóa đơn", "error");
    return;
  }
  await activateTab("invoiceTab");
  await loadInvoices(true);
  const match = invoicesData.find((invoice) => invoice.studentCode === student.studentCode);
  if (match) {
    selectInvoice(match.id);
    setInvoiceStatus("Đã chọn hóa đơn của học sinh", "ready");
  } else {
    setInvoiceStatus("Chưa có hóa đơn cho học sinh này trong danh sách hiện tại", "error");
  }
}

async function openMasterStudentNotifications(student) {
  if (!hasPermission("notification.view") && !hasPermission("notification.send")) {
    setMasterStatus("Không đủ quyền mở thông báo", "error");
    return;
  }
  await activateTab("notificationTab");
  await loadNotifications(true);
  notificationSchoolYearEl.value = optionValueOrEmpty(notificationSchoolYearEl, student.schoolYearId || "");
  renderNotificationGradeOptions();
  notificationGradeEl.value = optionValueOrEmpty(notificationGradeEl, student.grade || "");
  renderNotificationClassOptions();
  notificationClassEl.value = optionValueOrEmpty(notificationClassEl, student.classId || "");
  if (hasPermission("notification.send") && !openNotificationDialogBtn.hidden) {
    openNotificationDialog();
  }
}

function openMasterStudentClassScope(student) {
  const node = findSchoolTreeNode("class", student.classId || "");
  if (node) {
    selectedSchoolTreeNode = { type: "class", id: student.classId };
    fillSchoolTreeForms(node);
    applySchoolTreeFilters(node);
    renderSchoolTree();
  }
  document.querySelector(".master-tree-panel")?.scrollIntoView({ block: "start", behavior: "smooth" });
}

function defaultMasterParentDraft(isPrimary = false) {
  return {
    id: "",
    parentName: "",
    email: "",
    phone: "",
    relationship: "guardian",
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
    phone: parent.phone || "",
    relationship: parent.relationship || "guardian",
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
        <label>
          <span>SĐT</span>
          <input data-master-parent-field="phone" value="${escapeAttr(parent.phone || "")}" inputmode="tel" />
        </label>
        <label>
          <span>Quan hệ</span>
          <select data-master-parent-field="relationship">${masterRelationshipOptions(parent.relationship || "guardian")}</select>
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
    .join("") + `<div class="master-parent-warning" data-master-parent-warning></div>`;
  masterParentEditorRowsEl.querySelectorAll("[data-remove-master-parent]").forEach((button) => {
    button.addEventListener("click", () => {
      masterStudentParentDrafts = collectMasterParentDrafts();
      masterStudentParentDrafts.splice(Number(button.dataset.removeMasterParent), 1);
      renderMasterParentEditorRows();
    });
  });
  masterParentEditorRowsEl.querySelectorAll("input, select").forEach((field) => {
    field.addEventListener("input", updateMasterParentEditorWarning);
    field.addEventListener("change", updateMasterParentEditorWarning);
  });
  updateMasterParentEditorWarning();
}

function collectMasterParentDrafts() {
  return [...masterParentEditorRowsEl.querySelectorAll("[data-master-parent-index]")].map((row) => {
    const field = (name) => row.querySelector(`[data-master-parent-field="${name}"]`);
    return {
      id: field("id")?.value.trim() || "",
      parentName: field("parentName")?.value.trim() || "",
      email: field("email")?.value.trim() || "",
      phone: field("phone")?.value.trim() || "",
      relationship: field("relationship")?.value || "guardian",
      isPrimary: !!field("isPrimary")?.checked,
      isActive: !!field("isActive")?.checked,
      receivesBillingEmail: !!field("receivesBillingEmail")?.checked,
      emailActive: !!field("emailActive")?.checked,
    };
  });
}

function masterParentDraftBillingReady(parent) {
  return !!(parent.isActive && parent.receivesBillingEmail && parent.emailActive && parent.email);
}

function updateMasterParentEditorWarning() {
  const warningEl = masterParentEditorRowsEl.querySelector("[data-master-parent-warning]");
  if (!warningEl) return;
  const drafts = collectMasterParentDrafts();
  const readyCount = drafts.filter(masterParentDraftBillingReady).length;
  const parentCount = drafts.length;
  if (readyCount > 0) {
    warningEl.dataset.tone = "ready";
    warningEl.innerHTML = `${muiIcon("check_circle")}<span>${readyCount}/${parentCount} phụ huynh hợp lệ để nhận email học phí.</span>`;
    return;
  }
  warningEl.dataset.tone = "warning";
  warningEl.innerHTML = `${muiIcon("priority_high")}<span>Chưa có billing recipient hợp lệ: cần Active + Nhận phí + Email active + có email.</span>`;
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
      actor: true,
      details: [
        { label: "File", value: file.name || "-" },
        { label: "Preview", value: `${masterImportState?.preview?.length || 0} dòng` },
        { label: "Audit", value: "Conflict sẽ bị báo lỗi, không silent overwrite dữ liệu lệch" },
      ],
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
    schools: data.schools || [],
    feeTypes: data.feeTypes || defaultFeeTypes,
    schoolYears: data.schoolYears || [],
    classes: data.classes || [],
  };
  feeSchedulesLoaded = true;
  renderFeeScheduleControls();
  renderFeeScheduleItems(feeScheduleOptions.feeTypes);
  renderFeeAdjustmentRows();
  setFeeScheduleStatus("Sẵn sàng", "ready");
  return true;
}

async function loadFeeScheduleList() {
  if (!feeSchedulesLoaded) {
    return;
  }
  const params = new URLSearchParams();
  if (feeScheduleSchoolEl.value) params.set("schoolId", feeScheduleSchoolEl.value);
  if (feeScheduleYearEl.value) params.set("schoolYearId", feeScheduleYearEl.value);
  if (feeScheduleClassEl.value) params.set("classId", feeScheduleClassEl.value);
  if (!feeScheduleClassEl.value && feeScheduleGradeEl.value) params.set("grade", feeScheduleGradeEl.value);
  if (feeSchedulePeriodEl.value.trim()) params.set("periodCode", feeSchedulePeriodEl.value.trim());
  if (feeScheduleMonthEl.value) params.set("month", feeScheduleMonthEl.value);

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
  const selectedSchool = feeScheduleSchoolEl.value;
  const selectedYear = feeScheduleYearEl.value;
  const selectedGrade = feeScheduleGradeEl.value;
  const selectedClass = feeScheduleClassEl.value;

  feeScheduleSchoolEl.innerHTML = [
    `<option value="">Tất cả trường</option>`,
    ...(feeScheduleOptions.schools || []).map((item) => {
      const label = [item.code, item.name && item.name !== item.code ? item.name : ""].filter(Boolean).join(" · ");
      return `<option value="${escapeAttr(item.id)}">${escapeHtml(label || item.id)}</option>`;
    }),
  ].join("");
  feeScheduleSchoolEl.value = optionValueOrEmpty(feeScheduleSchoolEl, selectedSchool);

  feeScheduleYearEl.innerHTML = [
    `<option value="">Chọn năm học</option>`,
    ...(feeScheduleOptions.schoolYears || [])
      .filter((item) => !feeScheduleSchoolEl.value || item.schoolId === feeScheduleSchoolEl.value)
      .map((item) => {
        const label = [item.schoolCode, item.code].filter(Boolean).join(" · ");
        return `<option value="${escapeAttr(item.id)}">${escapeHtml(label || item.id)}</option>`;
      }),
  ].join("");
  feeScheduleYearEl.value = optionValueOrEmpty(feeScheduleYearEl, selectedYear);

  const grades = [...new Set(
    (feeScheduleOptions.classes || [])
      .filter((item) => !feeScheduleSchoolEl.value || item.schoolId === feeScheduleSchoolEl.value)
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
  renderAppContextControls();
}

function renderFeeScheduleClassFilter(selectedClass = feeScheduleClassEl.value) {
  const classes = (feeScheduleOptions.classes || []).filter((item) => {
    if (feeScheduleSchoolEl.value && item.schoolId !== feeScheduleSchoolEl.value) return false;
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
    .filter((item) => item.amount !== 0);
}

function defaultFeeAdjustmentDraft() {
  return {
    studentCode: "",
    adjustmentType: "discount",
    feeTypeCode: "",
    amount: 0,
    reason: "",
  };
}

function feeAdjustmentTypeOptions(selected) {
  const options = [
    ["discount", "Giảm trừ"],
    ["surcharge", "Phụ thu"],
    ["waiver", "Miễn giảm"],
    ["carry_over", "Chuyển kỳ trước"],
  ];
  return options
    .map(([value, label]) => `<option value="${escapeAttr(value)}" ${value === selected ? "selected" : ""}>${escapeHtml(label)}</option>`)
    .join("");
}

function feeAdjustmentFeeTypeOptions(selected) {
  const types = feeScheduleOptions.feeTypes?.length ? feeScheduleOptions.feeTypes : defaultFeeTypes;
  return [
    `<option value="">Toàn dòng</option>`,
    ...types.map((item) => {
      const value = item.code || "";
      const label = [item.code, item.labelVi].filter(Boolean).join(" · ");
      return `<option value="${escapeAttr(value)}" ${value === selected ? "selected" : ""}>${escapeHtml(label || value)}</option>`;
    }),
  ].join("");
}

function renderFeeAdjustmentRows(rows = null) {
  const sourceRows = rows || collectFeeAdjustmentRows();
  const drafts = sourceRows.length ? sourceRows : [defaultFeeAdjustmentDraft()];
  feeAdjustmentRowsEl.innerHTML = drafts
    .map((row) => {
      const adjustmentType = headerKeyClient(row.adjustmentType || "discount");
      const feeTypeCode = headerKeyClient(row.feeTypeCode || "");
      return `
        <tr>
          <td><input data-fee-adjust-field="studentCode" value="${escapeAttr(row.studentCode || "")}" placeholder="S001" /></td>
          <td><select data-fee-adjust-field="adjustmentType">${feeAdjustmentTypeOptions(adjustmentType)}</select></td>
          <td><select data-fee-adjust-field="feeTypeCode">${feeAdjustmentFeeTypeOptions(feeTypeCode)}</select></td>
          <td><input data-fee-adjust-field="amount" type="number" min="0" step="1000" value="${Number(row.amount || 0)}" /></td>
          <td><input data-fee-adjust-field="reason" value="${escapeAttr(row.reason || "")}" placeholder="Lý do bắt buộc" /></td>
          <td><button class="icon-only" type="button" data-fee-adjust-remove title="Xóa điều chỉnh">${muiIcon("delete")}<span class="sr-only">Xóa</span></button></td>
        </tr>
      `;
    })
    .join("");
  feeAdjustmentRowsEl.querySelectorAll("input, select").forEach((input) => {
    input.addEventListener("input", updateFeeAdjustmentCount);
    input.addEventListener("change", updateFeeAdjustmentCount);
  });
  feeAdjustmentRowsEl.querySelectorAll("[data-fee-adjust-remove]").forEach((button) => {
    button.addEventListener("click", () => {
      button.closest("tr")?.remove();
      if (!feeAdjustmentRowsEl.querySelector("tr")) {
        renderFeeAdjustmentRows([defaultFeeAdjustmentDraft()]);
        return;
      }
      updateFeeAdjustmentCount();
    });
  });
  updateFeeAdjustmentCount();
}

function addFeeAdjustmentRow(row = defaultFeeAdjustmentDraft()) {
  const rows = collectFeeAdjustmentRows({ includeBlank: true }).filter(
    (item) => item.studentCode || item.feeTypeCode || item.amount || item.reason,
  );
  rows.push(row);
  renderFeeAdjustmentRows(rows);
}

function collectFeeAdjustmentRows(options = {}) {
  return [...feeAdjustmentRowsEl.querySelectorAll("tr")]
    .map((row) => ({
      studentCode: row.querySelector('[data-fee-adjust-field="studentCode"]')?.value.trim().toUpperCase() || "",
      adjustmentType: headerKeyClient(row.querySelector('[data-fee-adjust-field="adjustmentType"]')?.value || ""),
      feeTypeCode: headerKeyClient(row.querySelector('[data-fee-adjust-field="feeTypeCode"]')?.value || ""),
      amount: Number(row.querySelector('[data-fee-adjust-field="amount"]')?.value || 0),
      reason: row.querySelector('[data-fee-adjust-field="reason"]')?.value.trim() || "",
    }))
    .filter((item) => options.includeBlank || item.studentCode || item.feeTypeCode || item.amount || item.reason);
}

function updateFeeAdjustmentCount() {
  const count = collectFeeAdjustmentRows().length + parseFeeAdjustmentsCsv().length;
  feeAdjustmentCountEl.textContent = `${count} điều chỉnh`;
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
    adjustments: [...collectFeeAdjustmentRows(), ...parseFeeAdjustmentsCsv()],
  };
}

function openFeeScheduleDialog() {
  renderFeeAdjustmentRows();
  setFeeGuideStep("scope");
  openAppDialog({
    title: "Tạo/sửa bảng phí",
    kicker: "Production fees",
    icon: "price_change",
    nodes: [feeGuideStepsEl, document.querySelector(".fee-schedule-grid"), document.querySelector(".fee-schedule-body")],
    size: "xl",
    actions: [
      { label: "Đóng", icon: "close", onClick: closeAppDialog },
      { label: "Preview", icon: "visibility", onClick: previewFeeSchedule },
      { label: "Lưu bảng phí", icon: "save", variant: "primary", onClick: saveFeeSchedule, closeOnSuccess: true },
    ],
  });
}

function setFeeGuideStep(step) {
  feeGuideStepsEl?.querySelectorAll("[data-fee-guide-step]").forEach((button) => {
    button.classList.toggle("is-active", button.dataset.feeGuideStep === step);
  });
  const panel = document.querySelector(`[data-fee-guide-panel="${step}"]`);
  panel?.scrollIntoView({ block: "nearest", behavior: "smooth" });
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
  setFeeGuideStep("preview");
  setFeeScheduleStatus(data.issues?.length ? "Có lỗi" : "Đã preview", data.issues?.length ? "error" : "ready");
  return !data.issues?.length;
}

async function saveFeeSchedule() {
  const draft = collectFeeScheduleDraft();
  const confirmed = await confirmDialog({
    title: "Lưu bảng phí?",
    message: "Thao tác này lưu bảng phí và điều chỉnh theo học sinh; các điều chỉnh cần lý do để phục vụ audit.",
    confirmLabel: "Lưu bảng phí",
    confirmIcon: "save",
    danger: draft.adjustments.length > 0 || draft.status === "active",
    actor: true,
    auditNote: draft.adjustments.length ? "Mỗi adjustment phải có reason rõ ràng trước khi lưu." : "Operator và thời điểm lưu được ghi nhận.",
    details: [
      { label: "Kỳ thu", value: draft.periodCode || "-" },
      { label: "Status", value: draft.status || "draft" },
      { label: "Điều chỉnh", value: `${draft.adjustments.length} dòng` },
      { label: "Operator", value: draft.operatorName || currentActorLabel() },
    ],
  });
  if (!confirmed) return false;
  setFeeScheduleStatus("Đang lưu", "busy");
  const res = await fetch("/api/v1/fee-schedules/save", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(draft),
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
  setFeeGuideStep("preview");
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
    feeSchedulePreviewRowsEl.innerHTML = `<tr><td colspan="7" class="empty-cell">Chưa có preview</td></tr>`;
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
          <td><span class="tag ${row.billingRecipientReady ? "tag-ready" : "tag-warning"}">${row.billingRecipientReady ? "Sẵn sàng" : "Thiếu"}</span></td>
        </tr>
      `,
    )
    .join("");
  if (!rows.length) {
    feeSchedulePreviewRowsEl.innerHTML = `<tr><td colspan="7" class="empty-cell">Không có học sinh trong phạm vi này</td></tr>`;
  }
}

function renderFeeSchedules(schedules) {
  schedules = schedules || [];
  feeSchedulesData = schedules || [];
  feeScheduleCountEl.textContent = `${schedules.length} bảng phí`;
  feeSchedulesEl.innerHTML = schedules
    .map((schedule) => {
      const scope = schedule.className || (schedule.grade ? `Khối ${schedule.grade}` : "Toàn năm học");
      const period = [schedule.periodCode, schedule.month ? `T${schedule.month}` : ""].filter(Boolean).join(" · ");
      const updated = [schedule.updatedActor || "system", schedule.updatedAt ? formatDateTime(schedule.updatedAt) : ""].filter(Boolean).join(" · ");
      return `
        <div class="fee-schedule-list-item">
          <div>
            <strong>${escapeHtml(schedule.name || schedule.periodCode || "Bảng phí")}</strong>
            <span>${escapeHtml([schedule.schoolCode, schedule.schoolYearCode, scope, period].filter(Boolean).join(" · "))}</span>
            <small>${Number(schedule.itemCount || 0)} khoản · ${Number(schedule.studentCount || 0)} HS · preview ${formatMoney(schedule.previewTotal || 0)} · ${escapeHtml(updated || "-")}</small>
          </div>
          <div class="fee-schedule-list-side">
            <span class="tag">${escapeHtml(schedule.status || "draft")}</span>
            <span class="money">${formatMoney(schedule.itemTotal || 0)}</span>
            <span>${Number(schedule.adjustmentCount || 0)} điều chỉnh</span>
            <button type="button" data-fee-schedule-invoice="${escapeAttr(schedule.id)}">
              ${muiIcon("request_quote")}<span>Sinh hóa đơn</span>
            </button>
          </div>
        </div>
      `;
    })
    .join("");
  if (!schedules.length) {
    feeSchedulesEl.innerHTML = `<div class="empty-cell fee-list-empty">Chưa có bảng phí đã lưu</div>`;
  }
  feeSchedulesEl.querySelectorAll("[data-fee-schedule-invoice]").forEach((button) => {
    button.addEventListener("click", () => openInvoiceFromFeeSchedule(button.dataset.feeScheduleInvoice || ""));
  });
}

async function openInvoiceFromFeeSchedule(scheduleId) {
  const schedule = feeSchedulesData.find((item) => item.id === scheduleId);
  if (!schedule) {
    setFeeScheduleStatus("Không tìm thấy bảng phí", "error");
    return;
  }
  appContext.schoolId = schedule.schoolId || appContext.schoolId;
  appContext.schoolYearId = schedule.schoolYearId || appContext.schoolYearId;
  appContext.periodCode = schedule.periodCode || appContext.periodCode;
  appContext.month = schedule.month ? String(schedule.month) : appContext.month;
  renderAppContextControls();

  await activateTab("invoiceTab");
  await loadInvoices(true);
  invoiceScheduleEl.value = optionValueOrEmpty(invoiceScheduleEl, scheduleId);
  if (!invoiceScheduleEl.value) {
    setInvoiceStatus("Không tìm thấy bảng phí trong danh sách hóa đơn", "error");
    return;
  }
  await previewInvoices();
  if (!openInvoiceDialogBtn.hidden) {
    openInvoiceDialog();
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
  const params = new URLSearchParams();
  const schedule = selectedInvoiceSchedule();
  if (schedule?.schoolYearId) params.set("schoolYearId", schedule.schoolYearId);
  if (schedule?.classId) params.set("classId", schedule.classId);
  if (!schedule?.classId && schedule?.grade) params.set("grade", schedule.grade);
  if (schedule?.periodCode) params.set("periodCode", schedule.periodCode);
  const res = await fetch(`/api/v1/invoices${params.toString() ? `?${params.toString()}` : ""}`);
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

function selectedInvoiceSchedule() {
  return (invoiceOptions.schedules || []).find((schedule) => schedule.id === invoiceScheduleEl.value) || null;
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
  renderAppContextControls();
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
  setInvoiceWorkbenchStep("scope");
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

function setInvoiceWorkbenchStep(step) {
  invoiceWorkbenchStepsEl?.querySelectorAll("[data-invoice-step]").forEach((button) => {
    button.classList.toggle("is-active", button.dataset.invoiceStep === step);
  });
  const panel = document.querySelector(`[data-invoice-step-panel="${step}"]`);
  panel?.scrollIntoView({ block: "nearest", behavior: "smooth" });
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
  setInvoiceWorkbenchStep(data.issues?.length ? "issues" : "preview");
  setInvoiceStatus(data.issues?.length ? "Có lỗi" : "Đã preview", data.issues?.length ? "error" : "ready");
  return !data.issues?.length;
}

async function generateInvoices() {
  const request = collectInvoiceRequest();
  const schedule = selectedInvoiceSchedule();
  const confirmed = await confirmDialog({
    title: "Sinh hóa đơn?",
    message: "Thao tác này sẽ ghi dữ liệu invoice vào database.",
    confirmLabel: "Sinh hóa đơn",
    confirmIcon: "post_add",
    danger: true,
    actor: true,
    auditNote: "Preview lại trước khi sinh nếu scope, kỳ thu hoặc regenerate vừa thay đổi.",
    details: [
      { label: "Bảng phí", value: schedule?.name || schedule?.periodCode || request.feeScheduleId || "-" },
      { label: "Kỳ thu", value: schedule?.periodCode || appContext.periodCode || "-" },
      { label: "Regenerate", value: request.regenerate ? "Có" : "Không" },
      { label: "Tài khoản", value: request.bankAccount || "-" },
    ],
  });
  if (!confirmed) return false;
  setInvoiceStatus("Đang sinh hóa đơn", "busy");
  const res = await fetch("/api/v1/invoices/generate", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(request),
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
  setInvoiceWorkbenchStep(data.preview?.issues?.length ? "issues" : "generate");
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
    invoiceIssuePanelEl.innerHTML = "";
    invoicePreviewRowsEl.innerHTML = `<tr><td colspan="8" class="empty-cell">Chưa có preview</td></tr>`;
    return;
  }

  const summary = data.summary || {};
  invoicePreviewSummaryEl.className = `invoice-summary${issues.length ? " error" : ""}`;
  const summaryText = [
    `${summary.studentCount || rows.length || 0} học sinh`,
    `đã có ${summary.existingCount || 0}`,
    `sẵn sàng ${summary.readyCount || 0}`,
    `regen ${summary.regenerableCount || 0}`,
    `chặn ${summary.blockedCount || 0}`,
    `mặc định ${formatMoney(summary.baseAmount || 0)}`,
    `điều chỉnh ${formatMoney(summary.adjustmentAmount || 0)}`,
    `phải thu ${formatMoney(summary.totalAmount || 0)}`,
  ].join(" · ");
  invoicePreviewSummaryEl.innerHTML = `<div>${escapeHtml(summaryText)}</div>`;
  renderInvoiceIssues(issues, rows);

  invoicePreviewRowsEl.innerHTML = rows
    .map((row) => {
      const itemCount = Number(row.itemCount ?? (row.items || []).length);
      const adjustmentCount = Number(row.adjustmentCount ?? (row.adjustments || []).length);
      return `
        <tr>
          <td><strong>${escapeHtml(row.invoiceCode || "")}</strong><small>${escapeHtml(row.periodCode || "")}</small></td>
          <td>${escapeHtml(row.studentCode || "")}<small>${escapeHtml(row.studentName || "")}</small></td>
          <td>${escapeHtml(row.className || "")}<small>${escapeHtml(row.schoolYearCode || "")}</small></td>
          <td>${itemCount}<small>${adjustmentCount ? `${adjustmentCount} adjustment` : "snapshot"}</small></td>
          <td class="money">${formatMoney(row.totalAmount || 0)}</td>
          <td><span class="tag ${row.billingRecipientReady ? "tag-ready" : "tag-warning"}">${row.billingRecipientReady ? "Ready" : "Missing"}</span></td>
          <td><span class="tag ${invoiceGenerationTone(row.generationState)}">${escapeHtml(invoiceGenerationLabel(row.generationState))}</span></td>
          <td><span class="tag ${invoiceIssueTone(row.issueState)}">${escapeHtml(invoiceIssueLabel(row.issueState))}</span></td>
        </tr>
      `;
    })
    .join("");
  if (!rows.length) {
    invoicePreviewRowsEl.innerHTML = `<tr><td colspan="8" class="empty-cell">Không có hóa đơn trong preview</td></tr>`;
  }
}

function renderInvoiceIssues(issues = [], rows = []) {
  const rowIssues = rows
    .filter((row) => row.issueState && row.issueState !== "ready")
    .map((row) => ({
      type: row.issueState,
      studentCode: row.studentCode,
      message:
        row.issueState === "blocking"
          ? invoiceGenerationLabel(row.generationState)
          : row.billingRecipientReady
            ? invoiceIssueLabel(row.issueState)
            : "Thiếu billing recipient đang hoạt động",
    }));
  const allIssues = [...issues, ...rowIssues].slice(0, 80);
  if (!allIssues.length) {
    invoiceIssuePanelEl.innerHTML = `<div class="invoice-issue-empty">${muiIcon("check_circle")}<span>Không có blocking issue trong preview hiện tại.</span></div>`;
    return;
  }
  invoiceIssuePanelEl.innerHTML = `
    <div class="invoice-issue-list">
      ${allIssues
        .map(
          (issue) => `
            <div class="invoice-issue-item ${issue.type === "blocking" || issue.type === "cannot_regenerate_paid_invoice" ? "is-blocking" : ""}">
              <span class="tag ${invoiceIssueTone(issue.type === "cannot_regenerate_paid_invoice" ? "blocking" : issue.type)}">${escapeHtml(issue.type || "issue")}</span>
              <strong>${escapeHtml(issue.studentCode || "-")}</strong>
              <span>${escapeHtml(issue.message || "")}</span>
            </div>
          `,
        )
        .join("")}
    </div>
  `;
}

function invoiceGenerationLabel(state) {
  const labels = {
    ready_to_generate: "Ready to generate",
    ready_to_regenerate: "Ready to regenerate",
    already_generated: "Already generated",
    blocked_paid_regenerate: "Blocked paid regen",
  };
  return labels[state] || state || "Ready";
}

function invoiceGenerationTone(state) {
  if (state === "blocked_paid_regenerate") return "tag-danger";
  if (state === "already_generated") return "tag-info";
  if (state === "ready_to_generate" || state === "ready_to_regenerate") return "tag-ready";
  return "tag-warning";
}

function invoiceIssueLabel(state) {
  const labels = {
    ready: "Ready",
    warning: "Warning",
    blocking: "Blocking",
    review: "Review",
    open: "Open",
  };
  return labels[state] || state || "Ready";
}

function invoiceIssueTone(state) {
  if (state === "ready") return "tag-ready";
  if (state === "blocking" || state === "review") return "tag-danger";
  if (state === "warning" || state === "open") return "tag-warning";
  return "tag-info";
}

function renderInvoices(invoices = []) {
  invoicesData = invoices || [];
  if (!invoicesData.some((invoice) => invoice.id === selectedInvoiceId)) {
    selectedInvoiceId = invoicesData[0]?.id || "";
  }
  invoiceCountEl.textContent = `${invoices.length} hóa đơn`;
  invoiceRowsEl.innerHTML = invoices
    .map((invoice) => {
      const outstanding = Number(invoice.outstandingAmount ?? Math.max(Number(invoice.totalAmount || 0) - Number(invoice.paidAmount || 0), 0));
      const sentLabel = invoice.sentCount ? `Sent ${Number(invoice.sentCount || 0)}` : "Chưa gửi";
      const sentMeta = invoice.lastSentAt ? formatDateTime(invoice.lastSentAt) : "";
      const notificationAction = isInvoiceNotificationCandidate(invoice)
        ? `<button type="button" data-invoice-notify="${escapeAttr(invoice.id || "")}">${muiIcon("campaign")}<span>Notify</span></button>`
        : "";
      const paymentIntentAction = hasPermission("payment.create")
        ? `<button type="button" data-invoice-intent="${escapeAttr(invoice.id || "")}">${muiIcon("add_card")}<span>Intent</span></button>`
        : "";
      return `
        <tr data-invoice-row="${escapeAttr(invoice.id || "")}" class="${invoice.id === selectedInvoiceId ? "is-selected" : ""}">
          <td><strong>${escapeHtml(invoice.invoiceCode || "")}</strong><small>${escapeHtml(invoice.periodCode || "")}</small></td>
          <td>${escapeHtml(invoice.studentCode || "")}<small>${escapeHtml(invoice.studentName || "")}</small></td>
          <td>${escapeHtml(invoice.className || "")}<small>${escapeHtml(invoice.schoolYearCode || "")}</small></td>
          <td class="money">${formatMoney(invoice.totalAmount || 0)}</td>
          <td class="money">${formatMoney(invoice.paidAmount || 0)}</td>
          <td class="money">${formatMoney(outstanding)}</td>
          <td><span class="tag ${invoiceIssueTone(invoice.issueState)}">${escapeHtml(invoice.status || "unpaid")}</span></td>
          <td><span class="tag ${invoice.sentCount ? "tag-ready" : "tag-info"}">${escapeHtml(sentLabel)}</span>${sentMeta ? `<small>${escapeHtml(sentMeta)}</small>` : ""}</td>
          <td>
            <div class="invoice-state-stack">
              <span class="tag ${invoice.qrReady ? "tag-ready" : "tag-warning"}">${invoice.qrReady ? "QR" : "No QR"}</span>
              <span class="tag ${invoice.pdfReady ? "tag-ready" : "tag-warning"}">${invoice.pdfReady ? "PDF" : "No PDF"}</span>
            </div>
          </td>
          <td>
            <div class="invoice-actions">
              <button type="button" data-invoice-qr="${escapeAttr(invoice.id || "")}">${muiIcon("qr_code")}<span>QR</span></button>
              <a class="button-link" href="/api/v1/invoices/pdf?id=${encodeURIComponent(invoice.id || "")}" target="_blank" rel="noreferrer">${muiIcon("picture_as_pdf")}<span>PDF</span></a>
              ${paymentIntentAction}
              ${notificationAction}
            </div>
          </td>
        </tr>
      `;
    })
    .join("");
  if (!invoices.length) {
    invoiceRowsEl.innerHTML = `<tr><td colspan="10" class="empty-cell">Chưa có hóa đơn</td></tr>`;
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
  invoiceRowsEl.querySelectorAll("[data-invoice-intent]").forEach((button) => {
    button.addEventListener("click", () => openPaymentIntentFromInvoice(invoicesData.find((invoice) => invoice.id === button.dataset.invoiceIntent)));
  });
  invoiceRowsEl.querySelectorAll("[data-invoice-notify]").forEach((button) => {
    button.addEventListener("click", () => openNotificationFromInvoice(invoicesData.find((invoice) => invoice.id === button.dataset.invoiceNotify)));
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
  const invoice = invoicesData.find((item) => item.id === selectedInvoiceId);
  renderInvoiceDetail(invoice);
  if (selectedInvoiceId) {
    loadInvoiceDetail(selectedInvoiceId);
  }
  if (!options.keepQr) {
    invoicePaymentStatusEl.textContent = selectedInvoiceId ? "Đã chọn" : "Chưa chọn";
    invoicePaymentStatusEl.dataset.tone = selectedInvoiceId ? "ready" : "";
    invoicePaymentPreviewEl.className = "preview-empty";
    invoicePaymentPreviewEl.textContent = selectedInvoiceId ? "Bấm QR để xem payload thanh toán của hóa đơn này" : "Chưa chọn hóa đơn";
  }
}

async function loadInvoiceDetail(invoiceId) {
  if (!invoiceId) return null;
  if (invoiceDetailCache.has(invoiceId)) {
    if (selectedInvoiceId === invoiceId) renderInvoiceDetail(invoiceDetailCache.get(invoiceId));
    return invoiceDetailCache.get(invoiceId);
  }
  const res = await fetch(`/api/v1/invoices/detail?id=${encodeURIComponent(invoiceId)}`);
  const text = await res.text();
  let invoice = null;
  try {
    invoice = JSON.parse(text);
  } catch {
    invoice = null;
  }
  if (selectedInvoiceId !== invoiceId) return null;
  if (!res.ok || !invoice) {
    setInvoiceStatus(text || "Không tải được chi tiết hóa đơn", "error");
    return null;
  }
  invoiceDetailCache.set(invoiceId, invoice);
  renderInvoiceDetail(invoice);
  return invoice;
}

function renderInvoiceDetail(invoice) {
  if (!invoice) {
    invoiceDetailSummaryEl.innerHTML = `<div class="detail-placeholder">${muiIcon("receipt_long")}<span>Chọn một hóa đơn để xem tổng tiền, trạng thái và thao tác QR/PDF.</span></div>`;
    return;
  }
  const outstanding = Number(invoice.outstandingAmount ?? Math.max(Number(invoice.totalAmount || 0) - Number(invoice.paidAmount || 0), 0));
  const hasSnapshot = Array.isArray(invoice.items);
  const items = invoice.items || [];
  const adjustments = invoice.adjustments || [];
  const history = invoice.statusHistory || [];
  const notificationAction = isInvoiceNotificationCandidate(invoice)
    ? `<button type="button" data-detail-invoice-notify="${escapeAttr(invoice.id || "")}">${muiIcon("campaign")}<span>Mở thông báo</span></button>`
    : "";
  const paymentIntentAction = hasPermission("payment.create")
    ? `<button type="button" data-detail-invoice-intent="${escapeAttr(invoice.id || "")}">${muiIcon("add_card")}<span>Tạo intent</span></button>`
    : "";
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
    <div class="invoice-mini-grid">
      <div><strong>${Number(invoice.itemCount ?? items.length)}</strong><span>line item</span></div>
      <div><strong>${Number(invoice.adjustmentCount ?? adjustments.length)}</strong><span>adjustment</span></div>
      <div><strong>${Number(invoice.paymentIntentCount || 0)}</strong><span>intent</span></div>
      <div><strong>${Number(invoice.matchedPaymentCount || 0)}</strong><span>matched</span></div>
      <div><strong>${Number(invoice.sentCount || 0)}</strong><span>sent</span></div>
      <div><strong>${invoice.qrReady ? "Ready" : "Missing"}</strong><span>QR</span></div>
    </div>
    ${
      hasSnapshot
        ? `
          <div class="detail-section">
            <h3 class="detail-section-title">Line-item snapshot</h3>
            <ul class="detail-list invoice-detail-list">
              ${
                items.length
                  ? items
                      .map(
                        (item) => `
                          <li>
                            <strong>${escapeHtml(item.labelVi || item.feeTypeCode || "")}</strong>
                            <span>${escapeHtml(item.labelEn || item.feeTypeCode || "")}</span>
                            <small>${formatMoney(item.amount || 0)}</small>
                          </li>
                        `,
                      )
                      .join("")
                  : `<li><span>Không có line item</span></li>`
              }
            </ul>
          </div>
          <div class="detail-section">
            <h3 class="detail-section-title">Adjustments</h3>
            <ul class="detail-list invoice-detail-list">
              ${
                adjustments.length
                  ? adjustments
                      .map(
                        (item) => `
                          <li>
                            <strong>${escapeHtml(item.labelVi || item.adjustmentType || "")}</strong>
                            <span>${escapeHtml(item.reason || item.feeTypeCode || "")}</span>
                            <small>${formatMoney(item.delta || item.amount || 0)}</small>
                          </li>
                        `,
                      )
                      .join("")
                  : `<li><span>Không có điều chỉnh</span></li>`
              }
            </ul>
          </div>
          <div class="detail-section">
            <h3 class="detail-section-title">Status history</h3>
            <ul class="detail-list invoice-detail-list">
              ${
                history.length
                  ? history
                      .map(
                        (item) => `
                          <li>
                            <strong>${escapeHtml([item.fromStatus || "-", item.toStatus || "-"].join(" -> "))}</strong>
                            <span>${escapeHtml(item.reason || "")}</span>
                            <small>${escapeHtml(formatDateTime(item.createdAt))}</small>
                          </li>
                        `,
                      )
                      .join("")
                  : `<li><span>Chưa có status history</span></li>`
              }
            </ul>
          </div>
        `
        : `<div class="invoice-detail-loading">${muiIcon("hourglass_top")}<span>Đang tải snapshot bất biến, lịch sử thanh toán và log gửi.</span></div>`
    }
    <div class="detail-actions">
      <button type="button" data-detail-invoice-qr="${escapeAttr(invoice.id || "")}">${muiIcon("qr_code")}<span>Xem QR</span></button>
      <a class="button-link" href="/api/v1/invoices/pdf?id=${encodeURIComponent(invoice.id || "")}" target="_blank" rel="noreferrer">${muiIcon("picture_as_pdf")}<span>Mở PDF</span></a>
      ${paymentIntentAction}
      ${notificationAction}
    </div>
  `;
  invoiceDetailSummaryEl.querySelectorAll("[data-detail-invoice-qr]").forEach((button) => {
    button.addEventListener("click", () => loadInvoicePayment(button.dataset.detailInvoiceQr));
  });
  invoiceDetailSummaryEl.querySelectorAll("[data-detail-invoice-intent]").forEach((button) => {
    button.addEventListener("click", () => openPaymentIntentFromInvoice(invoiceDetailCache.get(button.dataset.detailInvoiceIntent) || invoice));
  });
  invoiceDetailSummaryEl.querySelectorAll("[data-detail-invoice-notify]").forEach((button) => {
    button.addEventListener("click", () => openNotificationFromInvoice(invoiceDetailCache.get(button.dataset.detailInvoiceNotify) || invoice));
  });
}

function isInvoiceNotificationCandidate(invoice) {
  if (!invoice) return false;
  const outstanding = Number(invoice.outstandingAmount ?? Math.max(Number(invoice.totalAmount || 0) - Number(invoice.paidAmount || 0), 0));
  return outstanding > 0 && ["unpaid", "partial", "manual_review"].includes(invoice.status || "unpaid");
}

async function openPaymentIntentFromInvoice(invoice) {
  if (!invoice?.id) return;
  if (!hasPermission("payment.create")) {
    setInvoiceStatus("Không đủ quyền tạo payment intent", "error");
    return;
  }
  await activateTab("reconciliationTab");
  await loadPaymentReconciliation(true);
  await createPaymentIntent(invoice.id, "manual_vietqr");
}

async function openNotificationFromInvoice(invoice) {
  if (!invoice?.id) return;
  if (!hasPermission("notification.view") && !hasPermission("notification.send")) {
    setInvoiceStatus("Không đủ quyền mở thông báo", "error");
    return;
  }
  await activateTab("notificationTab");
  await loadNotifications(true);
  notificationSchoolYearEl.value = optionValueOrEmpty(notificationSchoolYearEl, invoice.schoolYearId || "");
  renderNotificationGradeOptions();
  notificationGradeEl.value = optionValueOrEmpty(notificationGradeEl, invoice.grade || "");
  renderNotificationClassOptions();
  notificationClassEl.value = optionValueOrEmpty(notificationClassEl, invoice.classId || "");
  notificationPeriodEl.value = invoice.periodCode || "";
  notificationInvoiceStatusEl.value = optionValueOrEmpty(notificationInvoiceStatusEl, ["partial", "manual_review"].includes(invoice.status) ? "partial" : "unpaid");
  setInvoiceStatus("Đã mở bộ lọc thông báo theo hóa đơn", "ready");
  if (hasPermission("notification.send") && !openNotificationDialogBtn.hidden) {
    openNotificationDialog();
  }
}

function exportInvoiceCsv() {
  if (!hasPermission("report.export")) {
    setInvoiceStatus("Không đủ quyền export", "error");
    return;
  }
  const params = new URLSearchParams({ dataset: "invoices" });
  const schedule = selectedInvoiceSchedule();
  if (schedule?.schoolYearId) params.set("schoolYearId", schedule.schoolYearId);
  if (schedule?.classId) params.set("classId", schedule.classId);
  if (!schedule?.classId && schedule?.grade) params.set("grade", schedule.grade);
  if (schedule?.periodCode) params.set("periodCode", schedule.periodCode);
  const link = document.createElement("a");
  link.href = `/api/v1/admin/reports/export?${params.toString()}`;
  link.download = "";
  document.body.appendChild(link);
  link.click();
  link.remove();
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
  if (paymentReconSchoolFilterEl.value) {
    params.set("schoolId", paymentReconSchoolFilterEl.value);
  }
  if (paymentReconYearFilterEl.value) {
    params.set("schoolYearId", paymentReconYearFilterEl.value);
  }
  if (paymentReconGradeFilterEl.value) {
    params.set("grade", paymentReconGradeFilterEl.value);
  }
  if (paymentReconClassFilterEl.value) {
    params.set("classId", paymentReconClassFilterEl.value);
  }
  if (paymentReconPeriodFilterEl.value.trim()) {
    params.set("periodCode", paymentReconPeriodFilterEl.value.trim());
  }
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
  renderPaymentReconFilters(data || paymentReconciliationData);
  renderPaymentReconSummary(data?.summary || null);
  renderPaymentReconInvoices(data?.invoices || [], data?.intents || {});
  renderPaymentReconTransactions(data?.transactions || []);
  renderPaymentReconReviewQueue(data?.invoices || [], data?.transactions || []);
}

function renderPaymentReconFilters(data = paymentReconciliationData) {
  const selectedSchool = paymentReconSchoolFilterEl.value;
  const selectedYear = paymentReconYearFilterEl.value;
  const selectedGrade = paymentReconGradeFilterEl.value;
  const selectedClass = paymentReconClassFilterEl.value;
  const schools = data?.schools || [];
  const years = data?.schoolYears || [];
  const classes = data?.classes || [];

  paymentReconSchoolFilterEl.innerHTML = [
    `<option value="">Tất cả trường</option>`,
    ...schools.map((school) => {
      const label = [school.code, school.name && school.name !== school.code ? school.name : ""].filter(Boolean).join(" · ");
      return `<option value="${escapeAttr(school.id)}">${escapeHtml(label || school.id)}</option>`;
    }),
  ].join("");
  paymentReconSchoolFilterEl.value = optionValueOrEmpty(paymentReconSchoolFilterEl, selectedSchool);

  const activeSchool = paymentReconSchoolFilterEl.value;
  const scopedYears = years
    .filter((year) => !activeSchool || year.schoolId === activeSchool)
    .sort((a, b) => String(a.code || a.name || "").localeCompare(String(b.code || b.name || ""), "vi", { numeric: true }));
  paymentReconYearFilterEl.innerHTML = [
    `<option value="">Tất cả năm học</option>`,
    ...scopedYears.map((year) => {
      const label = [year.schoolCode, year.code || year.name].filter(Boolean).join(" · ");
      return `<option value="${escapeAttr(year.id)}">${escapeHtml(label || year.id)}</option>`;
    }),
  ].join("");
  paymentReconYearFilterEl.value = optionValueOrEmpty(paymentReconYearFilterEl, selectedYear);

  const activeYear = paymentReconYearFilterEl.value;
  const scopedClasses = classes.filter((item) => {
    return (!activeSchool || item.schoolId === activeSchool) && (!activeYear || item.schoolYearId === activeYear);
  });
  const grades = [...new Set(scopedClasses.map((item) => item.grade).filter(Boolean))].sort((a, b) => String(a).localeCompare(String(b), "vi", { numeric: true }));
  paymentReconGradeFilterEl.innerHTML = [
    `<option value="">Tất cả khối</option>`,
    ...grades.map((grade) => `<option value="${escapeAttr(grade)}">Khối ${escapeHtml(grade)}</option>`),
  ].join("");
  paymentReconGradeFilterEl.value = optionValueOrEmpty(paymentReconGradeFilterEl, selectedGrade);

  const activeGrade = paymentReconGradeFilterEl.value;
  const classOptions = scopedClasses
    .filter((item) => !activeGrade || item.grade === activeGrade)
    .sort((a, b) => String(a.name || "").localeCompare(String(b.name || ""), "vi", { numeric: true }));
  paymentReconClassFilterEl.innerHTML = [
    `<option value="">Tất cả lớp</option>`,
    ...classOptions.map((item) => {
      const label = [item.schoolCode, item.schoolYearCode, item.name].filter(Boolean).join(" · ");
      return `<option value="${escapeAttr(item.id)}">${escapeHtml(label || item.id)}</option>`;
    }),
  ].join("");
  paymentReconClassFilterEl.value = optionValueOrEmpty(paymentReconClassFilterEl, selectedClass);

  renderPaymentProviderFilter(data?.providers || paymentReconciliationData.providers || []);
  renderAppContextControls();
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
  const manualReview = Number(summary.manualReviewCount || 0);
  paymentReconSummaryEl.innerHTML = `
    <div><strong>${formatMoney(summary.totalReceivable || 0)}</strong><span>Receivable</span></div>
    <div><strong>${formatMoney(summary.totalCollected || 0)}</strong><span>Collected</span></div>
    <div><strong>${formatMoney(outstanding)}</strong><span>Outstanding</span></div>
    <div><strong>${formatPercent(summary.collectionRate || 0)}</strong><span>Collection rate</span></div>
    <div><strong>${Number(summary.unpaidCount || 0)}</strong><span>Unpaid</span></div>
    <div><strong>${Number(summary.partialCount || 0)}</strong><span>Partial</span></div>
    <div><strong>${Number(summary.paidCount || 0)}</strong><span>Paid</span></div>
    <div><strong>${Number(summary.overpaidCount || 0)}</strong><span>Overpaid</span></div>
    <div><strong>${manualReview}</strong><span>Manual review</span></div>
    <div><strong>${Number(summary.unmatchedCount || 0)}</strong><span>Unmatched</span></div>
  `;
}

function renderPaymentReconInvoices(invoices, intents) {
  paymentReconInvoiceCountEl.textContent = `${invoices.length} hóa đơn`;
  const hasPayOS = (paymentReconciliationData.providers || []).some((provider) => provider.code === "payos");
  const canWritePayments = hasPermission("payment.create");
  const canNotify = hasPermission("notification.view") || hasPermission("notification.send");
  paymentReconInvoiceRowsEl.innerHTML = invoices
    .map((invoice) => {
      const paid = Number(invoice.paidAmount || 0);
      const total = Number(invoice.totalAmount || 0);
      const outstanding = Number(invoice.outstandingAmount ?? Math.max(total - paid, 0));
      const intent = intents?.[invoice.id];
      const intentLabel = intent?.provider ? `${intent.provider}: ${intent.status}` : "";
      const matches = paymentReconMatchesForInvoice(invoice.id);
      const matchLabel = matches.length ? `${matches.length} match` : invoice.matchedPaymentCount ? `${Number(invoice.matchedPaymentCount || 0)} matched` : "Chưa match";
      const notifyAction = canNotify && isInvoiceNotificationCandidate(invoice) ? `<button type="button" data-recon-notify="${escapeAttr(invoice.id || "")}">${muiIcon("campaign")}<span>Notify</span></button>` : "";
      const paymentActions = canWritePayments
        ? `
          <button type="button" data-recon-intent="${escapeAttr(invoice.id || "")}" data-recon-provider="manual_vietqr">${muiIcon("qr_code")}<span>QR</span></button>
          ${hasPayOS ? `<button type="button" data-recon-intent="${escapeAttr(invoice.id || "")}" data-recon-provider="payos">${muiIcon("link")}<span>payOS</span></button>` : ""}
          <button type="button" data-recon-cash="${escapeAttr(invoice.id || "")}" data-recon-default-amount="${escapeAttr(outstanding || total)}">${muiIcon("payments")}<span>Cash</span></button>
        `
        : "";
      const actions = `
        <div class="invoice-actions">
          <button type="button" data-recon-detail="${escapeAttr(invoice.id || "")}">${muiIcon("manage_search")}<span>Detail</span></button>
          ${paymentActions}
          <a class="button-link" href="/api/v1/invoices/pdf?id=${encodeURIComponent(invoice.id || "")}" target="_blank" rel="noreferrer">${muiIcon("picture_as_pdf")}<span>PDF</span></a>
          ${notifyAction}
        </div>
      `;
      return `
        <tr data-recon-invoice-row="${escapeAttr(invoice.id || "")}">
          <td><strong>${escapeHtml(invoice.invoiceCode || "")}</strong>${intentLabel ? `<small>${escapeHtml(intentLabel)}</small>` : ""}</td>
          <td>${escapeHtml(invoice.studentCode || "")}<small>${escapeHtml(invoice.studentName || "")}</small></td>
          <td>${escapeHtml(invoice.className || "")}<small>${escapeHtml(invoice.periodCode || invoice.schoolYearCode || "")}</small></td>
          <td class="money">${formatMoney(total)}</td>
          <td class="money">${formatMoney(paid)}<small>Còn ${formatMoney(outstanding)}</small></td>
          <td><span class="tag ${paymentReconStatusTone(invoice.status)}">${escapeHtml(invoice.status || "unpaid")}</span><small>${escapeHtml(matchLabel)}</small></td>
          <td>${actions}</td>
        </tr>
      `;
    })
    .join("");
  if (!invoices.length) {
    paymentReconInvoiceRowsEl.innerHTML = `<tr><td colspan="7" class="empty-cell">Chưa có hóa đơn để đối soát</td></tr>`;
  }
  paymentReconInvoiceRowsEl.querySelectorAll("[data-recon-detail]").forEach((button) => {
    button.addEventListener("click", () => showPaymentReconInvoiceDetail(button.dataset.reconDetail, { focusDetail: true }));
  });
  paymentReconInvoiceRowsEl.querySelectorAll("[data-recon-intent]").forEach((button) => {
    button.addEventListener("click", () => createPaymentIntent(button.dataset.reconIntent, button.dataset.reconProvider));
  });
  paymentReconInvoiceRowsEl.querySelectorAll("[data-recon-cash]").forEach((button) => {
    button.addEventListener("click", () => recordManualCashReceipt(button.dataset.reconCash, Number(button.dataset.reconDefaultAmount || 0)));
  });
  paymentReconInvoiceRowsEl.querySelectorAll("[data-recon-notify]").forEach((button) => {
    button.addEventListener("click", () => openNotificationFromInvoice((paymentReconciliationData.invoices || []).find((item) => item.id === button.dataset.reconNotify)));
  });
  paymentReconInvoiceRowsEl.querySelectorAll("[data-recon-invoice-row]").forEach((row) => {
    row.classList.toggle("is-selected", paymentReconSelection.type === "invoice" && paymentReconSelection.id === row.dataset.reconInvoiceRow);
    row.addEventListener("click", (event) => {
      if (event.target.closest("button, a")) return;
      showPaymentReconInvoiceDetail(row.dataset.reconInvoiceRow);
    });
  });
}

function renderPaymentReconTransactions(transactions) {
  paymentReconTransactionCountEl.textContent = `${transactions.length} giao dịch`;
  paymentReconTransactionRowsEl.innerHTML = transactions
    .map(
      (transaction) => `
        <tr data-recon-transaction-row="${escapeAttr(transaction.id || "")}">
          <td><strong>${escapeHtml(transaction.provider || "")}</strong><small>${escapeHtml(transaction.providerTransactionId || transaction.referenceCode || "")}</small></td>
          <td>${escapeHtml(formatDateTime(transaction.transactionTime))}</td>
          <td class="money">${formatMoney(transaction.amount || 0)}</td>
          <td>${escapeHtml(transaction.accountNumber || "")}<small>${escapeHtml(transaction.bankName || "")}</small></td>
          <td>${escapeHtml(transaction.description || transaction.referenceCode || "")}</td>
          <td>${escapeHtml(transaction.matchReason || transaction.referenceCode || "-")}<small>${escapeHtml(paymentMatchMeta(transaction))}</small></td>
          <td><span class="tag ${paymentReconStatusTone(transaction.matchStatus || transaction.status)}">${escapeHtml(transaction.invoiceCode || transaction.status || "unmatched")}</span></td>
        </tr>
      `,
    )
    .join("");
  if (!transactions.length) {
    paymentReconTransactionRowsEl.innerHTML = `<tr><td colspan="7" class="empty-cell">Chưa có giao dịch vào</td></tr>`;
  }
  paymentReconTransactionRowsEl.querySelectorAll("[data-recon-transaction-row]").forEach((row) => {
    row.classList.toggle("is-selected", paymentReconSelection.type === "transaction" && paymentReconSelection.id === row.dataset.reconTransactionRow);
    row.addEventListener("click", () => {
      showPaymentReconTransactionDetail(row.dataset.reconTransactionRow);
    });
  });
}

function renderPaymentReconReviewQueue(invoices, transactions) {
  const invoiceItems = (invoices || [])
    .filter((invoice) => ["partial", "overpaid", "manual_review"].includes(invoice.status || ""))
    .map((invoice) => {
      const outstanding = Number(invoice.outstandingAmount ?? Math.max(Number(invoice.totalAmount || 0) - Number(invoice.paidAmount || 0), 0));
      return {
        type: "invoice",
        id: invoice.id || "",
        title: invoice.invoiceCode || "Invoice",
        status: invoice.status || "review",
        amount: outstanding,
        meta: [invoice.studentName, invoice.className, invoice.periodCode].filter(Boolean).join(" · "),
      };
    });
  const transactionItems = (transactions || [])
    .filter((transaction) => ["unmatched", "manual_review"].includes(transaction.status || ""))
    .map((transaction) => ({
      type: "transaction",
      id: transaction.id || "",
      title: transaction.providerTransactionId || transaction.referenceCode || transaction.provider || "Transaction",
      status: transaction.status || "review",
      amount: Number(transaction.amount || 0),
      meta: transaction.description || transaction.matchReason || "Chưa match hóa đơn",
    }));
  const items = [...invoiceItems, ...transactionItems];
  paymentReconReviewCountEl.textContent = `${items.length} mục`;
  if (!items.length) {
    paymentReconReviewRowsEl.innerHTML = `<div class="detail-placeholder">${muiIcon("task_alt")}<span>Không có mục cần xử lý trong scope hiện tại.</span></div>`;
    return;
  }
  paymentReconReviewRowsEl.innerHTML = items
    .map((item) => {
      const selected = paymentReconSelection.type === item.type && paymentReconSelection.id === item.id;
      const dataAttr = item.type === "invoice" ? `data-review-invoice="${escapeAttr(item.id)}"` : `data-review-transaction="${escapeAttr(item.id)}"`;
      return `
        <button class="reconciliation-review-item ${selected ? "is-selected" : ""}" type="button" ${dataAttr}>
          <span class="tag ${paymentReconStatusTone(item.status)}">${escapeHtml(item.type === "invoice" ? "Invoice" : "Transaction")}</span>
          <strong>${escapeHtml(item.title)}</strong>
          <span>${escapeHtml(item.meta || "-")}</span>
          <small>${escapeHtml(item.status)} · ${formatMoney(item.amount || 0)}</small>
        </button>
      `;
    })
    .join("");
  paymentReconReviewRowsEl.querySelectorAll("[data-review-invoice]").forEach((button) => {
    button.addEventListener("click", () => showPaymentReconInvoiceDetail(button.dataset.reviewInvoice, { focusDetail: true }));
  });
  paymentReconReviewRowsEl.querySelectorAll("[data-review-transaction]").forEach((button) => {
    button.addEventListener("click", () => showPaymentReconTransactionDetail(button.dataset.reviewTransaction, { focusDetail: true }));
  });
}

function showPaymentReconInvoiceDetail(invoiceId, options = {}) {
  const invoice = (paymentReconciliationData.invoices || []).find((item) => item.id === invoiceId);
  renderPaymentReconDetail(invoiceDetailTemplate(invoice), { type: "invoice", id: invoiceId || "" });
  if (options.focusDetail) {
    setPaymentReconWorkbenchStep("detail");
  }
}

function showPaymentReconTransactionDetail(transactionId, options = {}) {
  const transaction = (paymentReconciliationData.transactions || []).find((item) => item.id === transactionId);
  renderPaymentReconDetail(transactionDetailTemplate(transaction), { type: "transaction", id: transactionId || "" });
  if (options.focusDetail) {
    setPaymentReconWorkbenchStep("detail");
  }
}

function setPaymentReconWorkbenchStep(step) {
  paymentReconStepsEl?.querySelectorAll("[data-reconciliation-step]").forEach((button) => {
    button.classList.toggle("is-active", button.dataset.reconciliationStep === step);
  });
  const panel = document.querySelector(`[data-reconciliation-step-panel="${step}"]`);
  panel?.scrollIntoView({ block: "start", behavior: "smooth" });
}

function reloadPaymentReconFromScopeFilters({ rerenderFilters = false } = {}) {
  if (rerenderFilters) {
    renderPaymentReconFilters(paymentReconciliationData);
  }
  syncAppContextFromActiveTab("reconciliationTab");
  return loadPaymentReconciliation(true);
}

function paymentReconMatchesForInvoice(invoiceId) {
  return paymentReconciliationData.matches?.[invoiceId] || [];
}

function paymentReconStatusTone(status) {
  if (["paid", "matched", "completed", "active"].includes(status || "")) return "tag-ready";
  if (["partial", "pending", "unpaid"].includes(status || "")) return "tag-warning";
  if (["overpaid", "manual_review", "unmatched", "failed", "cancelled"].includes(status || "")) return "tag-danger";
  return "tag-info";
}

function paymentMatchMeta(item) {
  const parts = [];
  if (item.matchType) parts.push(item.matchType);
  if (item.matchScore) parts.push(`score ${Number(item.matchScore || 0)}`);
  if (item.amountApplied) parts.push(`applied ${formatMoney(item.amountApplied)}`);
  return parts.join(" · ");
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
  setPaymentReconWorkbenchStep("detail");
  paymentReconciliationLoaded = false;
  await loadPaymentReconciliation(true);
}

function cashReceiptDialog(defaultAmount) {
  return new Promise((resolve) => {
    let settled = false;
    const body = document.createElement("div");
    body.className = "dialog-form-grid";
    body.innerHTML = `
      <div class="dialog-guardrail">
        <strong>${muiIcon("verified_user")}Audit-bound action</strong>
        <span>Actor: ${escapeHtml(currentActorLabel())}. Lý do ghi nhận là bắt buộc và sẽ đi vào audit/payment ledger.</span>
      </div>
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
            if (!field("reason").value.trim()) {
              showDialogError("Nhập lý do ghi nhận để phục vụ audit");
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
  setPaymentReconWorkbenchStep("detail");
  setPaymentReconStatus("Đã ghi nhận", "ready");
  return true;
}

function renderPaymentReconDetail(html, selection = null) {
  if (selection) {
    paymentReconSelection = selection;
    updatePaymentReconActiveRows();
  }
  paymentReconDetailEl.innerHTML = html || "Chưa chọn hóa đơn hoặc giao dịch";
  bindPaymentReconDetailActions();
}

function updatePaymentReconActiveRows() {
  paymentReconInvoiceRowsEl.querySelectorAll("[data-recon-invoice-row]").forEach((row) => {
    row.classList.toggle("is-selected", paymentReconSelection.type === "invoice" && paymentReconSelection.id === row.dataset.reconInvoiceRow);
  });
  paymentReconTransactionRowsEl.querySelectorAll("[data-recon-transaction-row]").forEach((row) => {
    row.classList.toggle("is-selected", paymentReconSelection.type === "transaction" && paymentReconSelection.id === row.dataset.reconTransactionRow);
  });
  paymentReconReviewRowsEl?.querySelectorAll("[data-review-invoice]").forEach((row) => {
    row.classList.toggle("is-selected", paymentReconSelection.type === "invoice" && paymentReconSelection.id === row.dataset.reviewInvoice);
  });
  paymentReconReviewRowsEl?.querySelectorAll("[data-review-transaction]").forEach((row) => {
    row.classList.toggle("is-selected", paymentReconSelection.type === "transaction" && paymentReconSelection.id === row.dataset.reviewTransaction);
  });
}

function bindPaymentReconDetailActions() {
  paymentReconDetailEl.querySelectorAll("[data-recon-detail-intent]").forEach((button) => {
    button.addEventListener("click", () => createPaymentIntent(button.dataset.reconDetailIntent, button.dataset.reconDetailProvider || "manual_vietqr"));
  });
  paymentReconDetailEl.querySelectorAll("[data-recon-detail-cash]").forEach((button) => {
    button.addEventListener("click", () => recordManualCashReceipt(button.dataset.reconDetailCash, Number(button.dataset.reconDefaultAmount || 0)));
  });
  paymentReconDetailEl.querySelectorAll("[data-recon-detail-notify]").forEach((button) => {
    button.addEventListener("click", () => openNotificationFromInvoice((paymentReconciliationData.invoices || []).find((item) => item.id === button.dataset.reconDetailNotify)));
  });
}

function invoiceDetailTemplate(invoice) {
  if (!invoice) return "Không tìm thấy hóa đơn";
  const total = Number(invoice.totalAmount || 0);
  const paid = Number(invoice.paidAmount || 0);
  const outstanding = Number(invoice.outstandingAmount ?? Math.max(total - paid, 0));
  const intent = paymentReconciliationData.intents?.[invoice.id];
  const matches = paymentReconMatchesForInvoice(invoice.id);
  const hasPayOS = (paymentReconciliationData.providers || []).some((provider) => provider.code === "payos");
  const paymentActions = hasPermission("payment.create")
    ? `
      <button type="button" data-recon-detail-intent="${escapeAttr(invoice.id || "")}" data-recon-detail-provider="manual_vietqr">${muiIcon("qr_code")}<span>QR/Intent</span></button>
      ${hasPayOS ? `<button type="button" data-recon-detail-intent="${escapeAttr(invoice.id || "")}" data-recon-detail-provider="payos">${muiIcon("link")}<span>payOS</span></button>` : ""}
      <button type="button" data-recon-detail-cash="${escapeAttr(invoice.id || "")}" data-recon-default-amount="${escapeAttr(outstanding || total)}">${muiIcon("payments")}<span>Cash receipt</span></button>
    `
    : "";
  const notificationAction = (hasPermission("notification.view") || hasPermission("notification.send")) && isInvoiceNotificationCandidate(invoice) ? `<button type="button" data-recon-detail-notify="${escapeAttr(invoice.id || "")}">${muiIcon("campaign")}<span>Notify</span></button>` : "";
  const intentDetail = intent
    ? `
      <div class="reconciliation-detail-grid">
        <span>Provider</span><strong>${escapeHtml(intent.provider || "")}</strong>
        <span>Intent</span><strong>${escapeHtml(intent.intentCode || "")}</strong>
        <span>Reference</span><strong>${escapeHtml(intent.providerReference || "-")}</strong>
        <span>Status</span><strong>${escapeHtml(intent.status || "")}</strong>
        <span>Created</span><strong>${escapeHtml(formatDateTime(intent.createdAt))}</strong>
      </div>
      ${intent.paymentUrl ? `<a class="button-link" href="${escapeAttr(intent.paymentUrl)}" target="_blank" rel="noreferrer">${muiIcon("open_in_new")}<span>Mở link thanh toán</span></a>` : ""}
    `
    : `<div class="invoice-detail-loading">${muiIcon("add_card")}<span>Chưa có payment intent trong scope hiện tại.</span></div>`;
  const matchList = matches.length
    ? matches
        .map(
          (match) => `
            <li>
              <strong>${escapeHtml(match.provider || "")} · ${escapeHtml(match.providerTransactionId || match.transactionId || "")}</strong>
              <span>${escapeHtml(match.reason || match.matchType || "Matched")}</span>
              <small>${escapeHtml(match.status || "")} · score ${Number(match.score || 0)} · ${formatMoney(match.amountApplied || 0)} · ${escapeHtml(formatDateTime(match.createdAt))}</small>
            </li>
          `,
        )
        .join("")
    : `<li><span>Chưa có reconciliation match.</span></li>`;
  return `
    <div class="detail-hero">
      ${muiIcon("receipt_long")}
      <div>
        <strong>${escapeHtml(invoice.invoiceCode || "-")}</strong>
        <span>${escapeHtml(invoice.studentCode || "")} · ${escapeHtml(invoice.studentName || "")}</span>
      </div>
    </div>
    <div class="reconciliation-detail-grid">
      <span>Lớp / kỳ</span><strong>${escapeHtml(invoice.className || "")} · ${escapeHtml(invoice.periodCode || "")}</strong>
      <span>Phải thu</span><strong>${formatMoney(total)}</strong>
      <span>Đã thu</span><strong>${formatMoney(paid)}</strong>
      <span>Còn thiếu</span><strong>${formatMoney(outstanding)}</strong>
      <span>Status</span><strong>${escapeHtml(invoice.status || "")}</strong>
      <span>Tài khoản thu</span><strong>${escapeHtml([invoice.bankBin || "", invoice.bankAccount || ""].filter(Boolean).join(" / ") || "-")}</strong>
      <span>QR reference</span><strong>${escapeHtml([invoice.qrBillNumber || "", invoice.qrNote || ""].filter(Boolean).join(" · ") || "-")}</strong>
    </div>
    <div class="invoice-mini-grid">
      <div><strong>${Number(invoice.paymentIntentCount || (intent ? 1 : 0))}</strong><span>intent</span></div>
      <div><strong>${matches.length || Number(invoice.matchedPaymentCount || 0)}</strong><span>matched</span></div>
      <div><strong>${Number(invoice.sentCount || 0)}</strong><span>sent</span></div>
      <div><strong>${invoice.qrReady ? "Ready" : "Missing"}</strong><span>QR</span></div>
      <div><strong>${invoice.pdfReady ? "Ready" : "Missing"}</strong><span>PDF</span></div>
      <div><strong>${escapeHtml(invoice.dueDate || "-")}</strong><span>due</span></div>
    </div>
    <div class="detail-section">
      <h3 class="detail-section-title">Payment intent</h3>
      ${intentDetail}
    </div>
    <div class="detail-section">
      <h3 class="detail-section-title">Reconciliation matches</h3>
      <ul class="detail-list invoice-detail-list">${matchList}</ul>
    </div>
    <div class="detail-actions">
      ${paymentActions}
      <a class="button-link" href="/api/v1/invoices/pdf?id=${encodeURIComponent(invoice.id || "")}" target="_blank" rel="noreferrer">${muiIcon("picture_as_pdf")}<span>PDF</span></a>
      ${notificationAction}
    </div>
  `;
}

function transactionDetailTemplate(transaction) {
  if (!transaction) return "Không tìm thấy giao dịch";
  return `
    <div class="detail-hero">
      ${muiIcon("sync_alt")}
      <div>
        <strong>${escapeHtml(transaction.providerTransactionId || transaction.referenceCode || transaction.provider || "-")}</strong>
        <span>${escapeHtml(transaction.description || "Giao dịch vào")}</span>
      </div>
    </div>
    <div class="reconciliation-detail-grid">
      <span>Provider</span><strong>${escapeHtml(transaction.provider || "")}</strong>
      <span>Reference</span><strong>${escapeHtml(transaction.providerTransactionId || transaction.referenceCode || "")}</strong>
      <span>Hóa đơn</span><strong>${escapeHtml(transaction.invoiceCode || "Chưa match")}</strong>
      <span>Số tiền</span><strong>${formatMoney(transaction.amount || 0)}</strong>
      <span>Applied</span><strong>${formatMoney(transaction.amountApplied || 0)}</strong>
      <span>Thời gian</span><strong>${escapeHtml(formatDateTime(transaction.transactionTime))}</strong>
      <span>Tài khoản</span><strong>${escapeHtml([transaction.bankName || "", transaction.accountNumber || ""].filter(Boolean).join(" / ") || "-")}</strong>
      <span>Nội dung</span><strong>${escapeHtml(transaction.description || "")}</strong>
      <span>Match reason</span><strong>${escapeHtml(transaction.matchReason || "-")}</strong>
      <span>Match score</span><strong>${escapeHtml(paymentMatchMeta(transaction) || "-")}</strong>
      <span>Status</span><strong>${escapeHtml(transaction.matchStatus || transaction.status || "")}</strong>
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
        <span>Created</span><strong>${escapeHtml(formatDateTime(intent.createdAt))}</strong>
      </div>
      ${link}
      ${qr.vietqr ? `<textarea class="payload" readonly>${escapeHtml(qr.vietqr)}</textarea>` : ""}
    </div>
  `;
}

async function loadNotifications(force = false) {
  if (notificationLoaded && !force) return;
  setNotificationStatus("Đang tải", "busy");
  const cronRequest = hasPermission("email_cron.view") || hasPermission("email_cron.update") ? fetch("/api/v1/email/cron") : Promise.resolve(null);
  const [optionsRes, logsRes, cronRes] = await Promise.all([
    fetch("/api/v1/notifications/options"),
    fetch("/api/v1/notifications/logs?limit=50"),
    cronRequest,
  ]);
  const optionsText = await optionsRes.text();
  if (!optionsRes.ok) {
    setNotificationStatus(optionsText || "Không tải được notification", "error");
    renderNotificationControls();
    renderNotificationPreview(null);
    renderNotificationCampaigns([]);
    renderNotificationLogs([]);
    renderNotificationCronSnapshot(null);
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
  if (cronRes?.ok) {
    notificationCronData = await cronRes.json();
    renderNotificationCronSnapshot(notificationCronData);
  } else {
    renderNotificationCronSnapshot(null);
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
  renderAppContextControls();
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
    forceResend: notificationForceResendEl.value === "true",
  };
}

function openNotificationDialog() {
  setNotificationWorkbenchStep("target");
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

function setNotificationWorkbenchStep(step) {
  notificationWorkbenchStepsEl?.querySelectorAll("[data-notification-step]").forEach((button) => {
    button.classList.toggle("is-active", button.dataset.notificationStep === step);
  });
  const panel = document.querySelector(`[data-notification-step-panel="${step}"]`);
  panel?.scrollIntoView({ block: "nearest", behavior: "smooth" });
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
  setNotificationWorkbenchStep((data.issues || []).length ? "target" : "recipients");
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
  setNotificationWorkbenchStep("recipients");
  setNotificationStatus("Đã lưu", "ready");
  return data.campaign;
}

async function sendNotificationCampaign() {
  const input = collectNotificationInput();
  const confirmed = await confirmDialog({
    title: "Gửi campaign?",
    message: "Thao tác này sẽ gửi email thật qua provider hiện tại và ghi log theo từng invoice/recipient.",
    confirmLabel: "Gửi campaign",
    confirmIcon: "send",
    danger: true,
    actor: true,
    auditNote: "Send log được ghi theo campaign/template/invoice/recipient.",
    details: [
      { label: "Campaign", value: input.name || currentNotificationCampaignId || "-" },
      { label: "Type", value: input.campaignType || "-" },
      { label: "Template", value: notificationTemplateEl.selectedOptions?.[0]?.textContent || "-" },
      { label: "Recipients", value: `${notificationPreviewData?.recipients?.length || 0} trong preview hiện tại` },
    ],
  });
  if (!confirmed) return false;
  setNotificationStatus("Đang gửi", "busy");
  const sendInput = { ...input, confirmSend: true, dryRun: false };
  const res = await fetch("/api/v1/notifications/campaigns/send", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(sendInput),
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
  setNotificationWorkbenchStep("send");
  setNotificationStatus("Đã xử lý gửi", "ready");
  return true;
}

async function retrySelectedNotificationRecipients() {
  const recipientIds = Array.from(selectedNotificationRecipientIds);
  if (!currentNotificationCampaignId) {
    setNotificationStatus("Chọn hoặc lưu campaign trước khi retry", "error");
    setNotificationWorkbenchStep("send");
    return false;
  }
  if (!recipientIds.length) {
    setNotificationStatus("Chọn recipient cần retry", "error");
    setNotificationWorkbenchStep("recipients");
    return false;
  }
  const confirmed = await confirmDialog({
    title: "Retry recipient đã chọn?",
    message: `Thao tác này sẽ gửi lại email thật cho ${recipientIds.length} recipient đã chọn và bỏ qua trạng thái đã gửi trước đó.`,
    confirmLabel: "Retry selected",
    confirmIcon: "replay",
    danger: true,
    actor: true,
    auditNote: "Retry dùng forceResend=true và ghi delivery log mới.",
    details: [
      { label: "Campaign", value: currentNotificationCampaignId || "-" },
      { label: "Recipients", value: `${recipientIds.length} selected` },
    ],
  });
  if (!confirmed) return false;
  setNotificationStatus("Đang retry", "busy");
  const input = {
    ...collectNotificationInput(),
    campaignId: currentNotificationCampaignId,
    recipientIds,
    forceResend: true,
    confirmSend: true,
    dryRun: false,
  };
  const res = await fetch("/api/v1/notifications/campaigns/send", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
  const text = await res.text();
  if (!res.ok) {
    setNotificationStatus(text || "Không retry được recipient", "error");
    return false;
  }
  const data = JSON.parse(text);
  currentNotificationCampaignId = data.campaign?.id || currentNotificationCampaignId;
  renderNotificationResults(data);
  await loadNotificationLogs(currentNotificationCampaignId);
  notificationLoaded = false;
  setNotificationWorkbenchStep("send");
  setNotificationStatus("Đã retry selected", "ready");
  return true;
}

function renderNotificationPreview(data) {
  const issues = data?.issues || [];
  const recipients = data?.recipients || [];
  const summary = data?.summary || {};
  if (!data) {
    selectedNotificationRecipientKey = "";
    selectedNotificationRecipientIds = new Set();
    notificationRecipientCountEl.textContent = "0 email";
    notificationSummaryEl.textContent = "Chưa có preview";
    notificationRecipientsEl.innerHTML = `<tr><td colspan="8" class="empty-cell">Chưa có recipient</td></tr>`;
    renderNotificationEmailPreviewPlaceholder();
    return;
  }
  const validRecipientIds = new Set(recipients.map((item) => item.id).filter(Boolean));
  selectedNotificationRecipientIds = new Set(Array.from(selectedNotificationRecipientIds).filter((id) => validRecipientIds.has(id)));
  if (selectedNotificationRecipientKey && !recipients.some((item) => notificationRecipientKey(item) === selectedNotificationRecipientKey)) {
    selectedNotificationRecipientKey = "";
    renderNotificationEmailPreviewPlaceholder();
  }
  notificationRecipientCountEl.textContent = selectedNotificationRecipientIds.size
    ? `${recipients.length} email · ${selectedNotificationRecipientIds.size} chọn`
    : `${recipients.length} email`;
  if (issues.length) {
    notificationSummaryEl.innerHTML = issues.map((item) => `<div><strong>${escapeHtml(item.type)}</strong><span>${escapeHtml(item.message)}</span></div>`).join("");
  } else {
    notificationSummaryEl.innerHTML = `
      <div><strong>${Number(summary.invoiceCount || 0)}</strong><span>invoice</span></div>
      <div><strong>${Number(summary.recipientCount || 0)}</strong><span>recipient</span></div>
      <div><strong>${formatMoney(summary.totalAmount || 0)}</strong><span>phải thu</span></div>
      <div><strong>${formatMoney(summary.unpaidAmount || 0)}</strong><span>còn phải thu</span></div>
      <div><strong>${Number(summary.alreadySent || 0)}</strong><span>đã gửi trước</span></div>
      <div><strong>${Number(summary.qrMissingCount || 0)}</strong><span>thiếu QR</span></div>
      <div><strong>${Number(summary.errorCount || 0)}</strong><span>lỗi gửi</span></div>
      <div><strong>${Number(summary.retryEligibleCount || 0)}</strong><span>retry</span></div>
    `;
  }
  notificationRecipientsEl.innerHTML = recipients.length
    ? recipients.map(notificationRecipientRowTemplate).join("")
    : `<tr><td colspan="8" class="empty-cell">Không có recipient phù hợp</td></tr>`;
  notificationRecipientsEl.querySelectorAll("tr[data-recipient-key]").forEach((row) => {
    row.addEventListener("click", (event) => {
      if (event.target.closest("button, input, label, a, select, textarea")) return;
      const recipient = findNotificationRecipientByKey(row.dataset.recipientKey || "");
      if (recipient) selectNotificationRecipient(recipient);
    });
  });
  notificationRecipientsEl.querySelectorAll("[data-notification-recipient-id]").forEach((checkbox) => {
    checkbox.addEventListener("change", () => {
      const id = checkbox.dataset.notificationRecipientId || "";
      if (!id) return;
      if (checkbox.checked) {
        selectedNotificationRecipientIds.add(id);
      } else {
        selectedNotificationRecipientIds.delete(id);
      }
      renderNotificationPreview({ ...notificationPreviewData, recipients, summary, issues });
    });
  });
  notificationRecipientsEl.querySelectorAll("[data-notification-email-preview]").forEach((button) => {
    button.addEventListener("click", () => {
      const recipient = findNotificationRecipientByKey(button.dataset.notificationEmailPreview || "");
      if (recipient) selectNotificationRecipient(recipient);
    });
  });
}

function notificationRecipientRowTemplate(item) {
  const key = notificationRecipientKey(item);
  const status = notificationRecipientStatus(item);
  const tone = notificationRecipientTone(item);
  const checked = item.id && selectedNotificationRecipientIds.has(item.id);
  const selectedClass = selectedNotificationRecipientKey === key ? " class=\"is-selected\"" : "";
  const qrLabel = item.qrReady ? "QR ready" : "No QR";
  const sentLabel = item.sendCount ? `${Number(item.sendCount || 0)} sent` : "Chưa gửi";
  const retryLabel = item.retryEligible ? `<span class="tag tag-warning">${muiIcon("replay")}Retry</span>` : `<span class="tag tag-info">No retry</span>`;
  return `
    <tr data-recipient-key="${escapeAttr(key)}"${selectedClass}>
      <td>
        <input class="notification-recipient-check" type="checkbox" data-notification-recipient-id="${escapeAttr(item.id || "")}" ${checked ? "checked" : ""} ${item.id ? "" : "disabled"} aria-label="Chọn recipient" />
      </td>
      <td><strong>${escapeHtml(item.invoiceCode || "")}</strong><small>${escapeHtml(item.periodCode || "")}</small></td>
      <td>${escapeHtml(item.studentCode || "")}<small>${escapeHtml(item.studentName || "")} · ${escapeHtml(item.className || "")}</small></td>
      <td>${escapeHtml(item.recipientEmail || "")}<small>${escapeHtml(item.recipientName || "")}</small></td>
      <td>${formatMoney(item.outstandingAmount ?? item.amount ?? 0)}<small>Tổng ${formatMoney(item.amount || 0)} · đã thu ${formatMoney(item.paidAmount || 0)}</small></td>
      <td><span class="tag ${item.qrReady ? "tag-ready" : "tag-warning"}">${item.qrReady ? muiIcon("qr_code_2") : muiIcon("qr_code_scanner")}${escapeHtml(qrLabel)}</span></td>
      <td>
        <span class="tag ${item.sendCount ? "tag-ready" : "tag-info"}">${escapeHtml(sentLabel)}</span>
        <small>${escapeHtml(item.lastSentAt ? formatDateTime(item.lastSentAt) : item.lastLogStatus || "")}</small>
        ${retryLabel}
      </td>
      <td>
        <span class="status-pill" data-tone="${escapeAttr(tone)}">${escapeHtml(status.replaceAll("_", " "))}</span>
        ${item.lastError ? `<small>${escapeHtml(item.lastError)}</small>` : ""}
        <button class="icon-button notification-preview-email-btn" type="button" data-notification-email-preview="${escapeAttr(key)}" title="Preview email">${muiIcon("drafts")}<span class="sr-only">Preview email</span></button>
      </td>
    </tr>
  `;
}

function notificationRecipientKey(item) {
  return item.id || `${item.invoiceId || ""}|${String(item.recipientEmail || "").toLowerCase()}`;
}

function findNotificationRecipientByKey(key) {
  return (notificationPreviewData.recipients || []).find((item) => notificationRecipientKey(item) === key) || null;
}

function notificationRecipientStatus(item) {
  if (item.status) return item.status;
  if (item.alreadySent) return "already_sent";
  return item.invoiceStatus || "pending";
}

function notificationRecipientTone(item) {
  const status = notificationRecipientStatus(item);
  if (status === "error") return "error";
  if (status === "sent" || status === "dry_run") return "ready";
  if (item.retryEligible || status === "skipped") return "error";
  if (status === "already_sent") return "warning";
  return "ready";
}

function selectNotificationRecipient(recipient) {
  selectedNotificationRecipientKey = notificationRecipientKey(recipient);
  notificationRecipientsEl.querySelectorAll("tr[data-recipient-key]").forEach((row) => {
    row.classList.toggle("is-selected", row.dataset.recipientKey === selectedNotificationRecipientKey);
  });
  previewNotificationRecipientEmail(recipient);
}

function renderNotificationEmailPreviewPlaceholder(message = "Chọn recipient để xem email preview") {
  if (notificationEmailPreviewStatusEl) notificationEmailPreviewStatusEl.textContent = "Chưa chọn";
  if (notificationEmailPreviewFrameEl) {
    notificationEmailPreviewFrameEl.srcdoc = `<body style="margin:0;padding:18px;font-family:system-ui,sans-serif;color:#667168;background:#fff;">${escapeHtml(message)}</body>`;
  }
}

async function previewNotificationRecipientEmail(recipient) {
  if (!recipient) {
    renderNotificationEmailPreviewPlaceholder();
    return false;
  }
  notificationEmailPreviewStatusEl.textContent = "Đang render";
  setNotificationWorkbenchStep("email");
  const res = await fetch("/api/v1/notifications/campaigns/email-preview", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      ...collectNotificationInput(),
      campaignId: currentNotificationCampaignId,
      recipientId: recipient.id || "",
      invoiceId: recipient.invoiceId || "",
      recipientEmail: recipient.recipientEmail || "",
    }),
  });
  const text = await res.text();
  if (!res.ok) {
    notificationEmailPreviewStatusEl.textContent = "Lỗi";
    notificationEmailPreviewFrameEl.srcdoc = `<body style="margin:0;padding:18px;font-family:system-ui,sans-serif;color:#ba2c2c;background:#fff;">${escapeHtml(text || "Không render được email preview")}</body>`;
    setNotificationStatus(text || "Không render được email preview", "error");
    return false;
  }
  const data = JSON.parse(text);
  notificationEmailPreviewStatusEl.textContent = data.to || recipient.recipientEmail || "Preview";
  notificationEmailPreviewFrameEl.srcdoc = data.html || `<pre>${escapeHtml(data.text || data.subject || "Email preview")}</pre>`;
  setNotificationStatus("Email preview xong", "ready");
  return true;
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
  notificationForceResendEl.value = "false";
  setNotificationWorkbenchStep("send");
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
      if (!result) return recipient;
      return {
        ...recipient,
        status: result.status,
        lastError: result.error || "",
        lastLogStatus: result.status,
        alreadySent: result.status === "sent" ? true : recipient.alreadySent,
        sendCount: result.status === "sent" ? Number(recipient.sendCount || 0) + 1 : Number(recipient.sendCount || 0),
        retryEligible: result.status === "error" || result.status === "skipped",
      };
    }),
  };
  renderNotificationPreview(notificationPreviewData);
  renderNotificationLogs(logs);
  if (data.campaign) {
    const campaigns = [data.campaign, ...(notificationOptions.campaigns || []).filter((item) => item.id !== data.campaign.id)];
    renderNotificationCampaigns(campaigns);
  }
}

async function loadNotificationCronSnapshot() {
  if (!hasPermission("email_cron.view") && !hasPermission("email_cron.update")) {
    renderNotificationCronSnapshot(null);
    return null;
  }
  const res = await fetch("/api/v1/email/cron");
  if (!res.ok) {
    renderNotificationCronSnapshot(null);
    return null;
  }
  notificationCronData = await res.json();
  renderNotificationCronSnapshot(notificationCronData);
  return notificationCronData;
}

function renderNotificationCronSnapshot(data) {
  if (!notificationCronSnapshotEl) return;
  if (!data) {
    notificationCronSnapshotEl.textContent = hasPermission("email_cron.view") || hasPermission("email_cron.update") ? "Không tải được cron" : "Không có quyền xem cron";
    return;
  }
  const sentLast24h = data.sentLast24h ?? data.sentToday ?? 0;
  const metrics = [
    { label: "Cron", value: data.enabled ? "ACTIVE" : "PAUSED", tone: data.enabled ? "tag-ready" : "tag-warning" },
    { label: "Send time", value: data.sendTime || "08:00", tone: "tag-info" },
    { label: "Daily limit", value: Number(data.dailyLimit || 0), tone: "tag-info" },
    { label: "Queued", value: Number(data.queued || 0), tone: data.queued ? "tag-warning" : "tag-ready" },
    { label: "Sent", value: Number(data.sent || 0), tone: "tag-ready" },
    { label: "Errors", value: Number(data.errors || 0), tone: data.errors ? "tag-danger" : "tag-ready" },
    { label: "Sent 24h", value: sentLast24h, tone: "tag-info" },
  ];
  const recent = (data.lastResults || [])
    .slice(0, 6)
    .map((item) => {
      const message = item.error || item.resendId || item.messageId || item.status;
      return `<div><strong>${escapeHtml(item.status || "-")}</strong><span>${escapeHtml(item.email || "-")} · ${escapeHtml(item.studentName || "-")}</span><small>${escapeHtml(message || "")}</small></div>`;
    })
    .join("");
  notificationCronSnapshotEl.innerHTML = `
    <div class="notification-cron-metrics">
      ${metrics.map((item) => `<div><span>${escapeHtml(item.label)}</span><strong class="tag ${item.tone}">${escapeHtml(item.value)}</strong></div>`).join("")}
    </div>
    <div class="notification-cron-times">
      <span>Next ${escapeHtml(data.nextRunAt ? formatDateTime(data.nextRunAt) : "-")}</span>
      <span>Last ${escapeHtml(data.lastRunAt ? formatDateTime(data.lastRunAt) : "-")}</span>
    </div>
    <div class="notification-cron-results">${recent || "<span>Chưa có kết quả gần đây</span>"}</div>
  `;
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
    const rows = collectRows();
    const confirmed = await confirmDialog({
      title: "Gửi email thật?",
      message: "Thao tác này sẽ gửi email thật qua provider hiện tại cho các dòng đang có trong bảng.",
      confirmLabel: "Gửi email",
      confirmIcon: "send",
      danger: true,
      actor: true,
      auditNote: "Preview hoặc dry-run trước khi gửi thật; provider secrets không hiển thị trong UI.",
      details: [
        { label: "Provider", value: emailProviderEl.value || "gmail" },
        { label: "Recipients", value: `${rows.filter((row) => row.email).length}/${rows.length} dòng có email` },
        { label: "Template", value: emailTemplateEl.value || "-" },
      ],
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
  const confirmed = await confirmDialog({
    title: "Tắt cron gửi email?",
    message: "Thao tác này dừng scheduler cục bộ; các email chưa gửi sẽ vẫn nằm trong queue nhưng không tự chạy.",
    confirmLabel: "Tắt cron",
    confirmIcon: "event_busy",
    danger: true,
    actor: true,
    details: [
      { label: "Queue", value: cronQueueSummaryEl.value || "-" },
      { label: "Provider", value: emailProviderEl.value || "gmail" },
    ],
  });
  if (!confirmed) return false;
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
    actor: true,
    auditNote: "Chỉ chạy khi đã kiểm tra quota và queue; không dùng để test gửi thật.",
    details: [
      { label: "Queue", value: cronQueueSummaryEl.value || "-" },
      { label: "Daily limit", value: cronDailyLimitEl.value || "-" },
      { label: "Provider", value: emailProviderEl.value || "gmail" },
    ],
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
  notificationCronData = data;
  renderNotificationCronSnapshot(data);

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
  if (appDialogEl.dataset.busy === "true") return;
  if (event.target === appDialogEl) {
    closeAppDialog();
  }
});
appDialogEl.addEventListener("cancel", (event) => {
  if (appDialogEl.dataset.busy === "true") {
    event.preventDefault();
  }
});
appDialogEl.addEventListener("close", () => {
  const onClose = activeDialogOnClose;
  const trigger = lastDialogTrigger;
  lastDialogTrigger = null;
  activeDialogOnClose = null;
  delete appDialogEl.dataset.busy;
  appDialogEl.removeAttribute("aria-busy");
  appDialogEl.removeAttribute("aria-describedby");
  restoreDialogContent();
  if (onClose) onClose();
  if (trigger?.isConnected && !appDialogEl.open) {
    window.setTimeout(() => trigger.focus(), 0);
  }
});
document.addEventListener("keydown", activateFocusedInteractiveRow);
new MutationObserver(scheduleEnhanceInteractiveRows).observe(document.body, {
  attributes: true,
  attributeFilter: ["class"],
  childList: true,
  subtree: true,
});
scheduleEnhanceInteractiveRows();

tabButtons.forEach((button) => {
  button.addEventListener("click", () => activateTab(button.dataset.tabTarget));
});

appContextSchoolEl.addEventListener("change", applyAppContextToActiveTab);
appContextYearEl.addEventListener("change", applyAppContextToActiveTab);
appContextPeriodEl.addEventListener("change", applyAppContextToActiveTab);
appContextPeriodEl.addEventListener("input", scheduleApplyAppContext);
appContextMonthEl.addEventListener("change", applyAppContextToActiveTab);
appContextMonthEl.addEventListener("input", scheduleApplyAppContext);
adminWorkQueueEl.addEventListener("click", (event) => {
  const button = event.target.closest("[data-dashboard-action]");
  if (button) runDashboardAction(button.dataset.dashboardAction, button.dataset.readinessIssue || "");
});
adminQuickActionsEl.addEventListener("click", (event) => {
  const button = event.target.closest("[data-dashboard-action]");
  if (button) runDashboardAction(button.dataset.dashboardAction, button.dataset.readinessIssue || "");
});
operatorOnboardingEl?.addEventListener("click", (event) => {
  const button = event.target.closest("[data-dashboard-action]");
  if (button) runDashboardAction(button.dataset.dashboardAction, button.dataset.readinessIssue || "");
});
adminReadinessCenterEl.addEventListener("click", (event) => {
  const button = event.target.closest("[data-dashboard-action]");
  if (button) runDashboardAction(button.dataset.dashboardAction, button.dataset.readinessIssue || "");
});
adminReadinessSeverityEl.addEventListener("change", () => renderAdminReadiness(adminDashboardData?.readiness || null));
adminReadinessTypeEl.addEventListener("change", () => renderAdminReadiness(adminDashboardData?.readiness || null));

refreshAdminDashboardBtn.addEventListener("click", () => loadAdminDashboard(true));
adminDashboardSchoolEl.addEventListener("change", async () => {
  renderAdminFilters("dashboard");
  syncAppContextFromActiveTab("dashboardTab");
  await loadAdminDashboard(true);
});
adminDashboardYearEl.addEventListener("change", async () => {
  renderAdminFilters("dashboard");
  syncAppContextFromActiveTab("dashboardTab");
  await loadAdminDashboard(true);
});
adminDashboardGradeEl.addEventListener("change", async () => {
  renderAdminFilters("dashboard");
  await loadAdminDashboard(true);
});
adminDashboardClassEl.addEventListener("change", () => loadAdminDashboard(true));
adminDashboardPeriodEl.addEventListener("change", () => {
  syncAppContextFromActiveTab("dashboardTab");
  loadAdminDashboard(true);
});
adminDashboardMonthEl.addEventListener("change", () => {
  syncAppContextFromActiveTab("dashboardTab");
  loadAdminDashboard(true);
});
adminDashboardInvoiceStatusEl.addEventListener("change", () => loadAdminDashboard(true));

refreshAdminReportsBtn.addEventListener("click", () => loadAdminReports(true));
exportAdminReportClassesBtn.addEventListener("click", () => exportAdminReport("classes"));
exportAdminReportInvoicesBtn.addEventListener("click", () => exportAdminReport("invoices"));
exportAdminReportTransactionsBtn.addEventListener("click", () => exportAdminReport("transactions"));
adminReportsSchoolEl.addEventListener("change", async () => {
  renderAdminFilters("reports");
  syncAppContextFromActiveTab("reportsTab");
  await loadAdminReports(true);
});
adminReportsYearEl.addEventListener("change", async () => {
  renderAdminFilters("reports");
  syncAppContextFromActiveTab("reportsTab");
  await loadAdminReports(true);
});
adminReportsGradeEl.addEventListener("change", async () => {
  renderAdminFilters("reports");
  await loadAdminReports(true);
});
adminReportsClassEl.addEventListener("change", () => loadAdminReports(true));
adminReportsPeriodEl.addEventListener("change", () => {
  syncAppContextFromActiveTab("reportsTab");
  loadAdminReports(true);
});
adminReportsMonthEl.addEventListener("change", () => {
  syncAppContextFromActiveTab("reportsTab");
  loadAdminReports(true);
});
adminReportsInvoiceStatusEl.addEventListener("change", () => loadAdminReports(true));
adminReportsProviderEl.addEventListener("change", () => loadAdminReports(true));
refreshOperationsBtn.addEventListener("click", () => loadOperations(true));
operationSourceFilterEl.addEventListener("change", () => loadOperations(true));
operationLevelFilterEl.addEventListener("change", () => loadOperations(true));
operationNameFilterEl.addEventListener("change", () => loadOperations(true));
operationStatusFilterEl.addEventListener("change", () => loadOperations(true));
auditActionFilterEl.addEventListener("change", () => loadOperations(true));
operationEntityTypeFilterEl.addEventListener("change", () => loadOperations(true));
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
  syncAppContextFromActiveTab("masterDataTab");
  await loadMasterStudents();
});
masterSchoolYearFilterEl.addEventListener("change", async () => {
  renderMasterFilters();
  syncAppContextFromActiveTab("masterDataTab");
  await loadMasterStudents();
});
masterGradeFilterEl.addEventListener("change", async () => {
  renderMasterClassFilter();
  await loadMasterStudents();
});
masterClassFilterEl.addEventListener("change", loadMasterStudents);
masterBillingFilterEl.addEventListener("change", () => renderMasterStudents(masterStudentsRawData));
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
feeGuideStepsEl?.querySelectorAll("[data-fee-guide-step]").forEach((button) => {
  button.addEventListener("click", () => setFeeGuideStep(button.dataset.feeGuideStep || "scope"));
});
feeScheduleSchoolEl.addEventListener("change", async () => {
  renderFeeScheduleControls();
  syncAppContextFromActiveTab("feeTemplateTab");
  await loadFeeScheduleList();
});
feeScheduleYearEl.addEventListener("change", async () => {
  renderFeeScheduleControls();
  syncAppContextFromActiveTab("feeTemplateTab");
  await loadFeeScheduleList();
});
feeScheduleGradeEl.addEventListener("change", async () => {
  renderFeeScheduleClassFilter();
  await loadFeeScheduleList();
});
feeScheduleClassEl.addEventListener("change", loadFeeScheduleList);
feeSchedulePeriodEl.addEventListener("change", async () => {
  syncAppContextFromActiveTab("feeTemplateTab");
  await loadFeeScheduleList();
});
feeScheduleMonthEl.addEventListener("change", async () => {
  syncAppContextFromActiveTab("feeTemplateTab");
  await loadFeeScheduleList();
});
addFeeAdjustmentRowBtn.addEventListener("click", () => addFeeAdjustmentRow());
feeAdjustmentsCsvEl.addEventListener("input", () => {
  updateFeeAdjustmentCount();
});
previewFeeScheduleBtn.addEventListener("click", previewFeeSchedule);
saveFeeScheduleBtn.addEventListener("click", saveFeeSchedule);

refreshInvoicesBtn.addEventListener("click", () => loadInvoices(true));
openInvoiceDialogBtn.addEventListener("click", openInvoiceDialog);
invoiceWorkbenchStepsEl?.querySelectorAll("[data-invoice-step]").forEach((button) => {
  button.addEventListener("click", () => setInvoiceWorkbenchStep(button.dataset.invoiceStep || "scope"));
});
invoiceScheduleEl.addEventListener("change", async () => {
  invoiceDetailCache = new Map();
  await loadInvoiceList();
});
previewInvoicesBtn.addEventListener("click", previewInvoices);
generateInvoicesBtn.addEventListener("click", generateInvoices);
exportInvoiceCsvBtn.addEventListener("click", exportInvoiceCsv);
refreshPaymentReconBtn.addEventListener("click", () => loadPaymentReconciliation(true));
paymentReconStepsEl?.querySelectorAll("[data-reconciliation-step]").forEach((button) => {
  button.addEventListener("click", () => setPaymentReconWorkbenchStep(button.dataset.reconciliationStep || "scope"));
});
paymentReconSchoolFilterEl.addEventListener("change", () => reloadPaymentReconFromScopeFilters({ rerenderFilters: true }));
paymentReconYearFilterEl.addEventListener("change", () => reloadPaymentReconFromScopeFilters({ rerenderFilters: true }));
paymentReconGradeFilterEl.addEventListener("change", () => reloadPaymentReconFromScopeFilters({ rerenderFilters: true }));
paymentReconClassFilterEl.addEventListener("change", () => reloadPaymentReconFromScopeFilters());
paymentReconPeriodFilterEl.addEventListener("change", () => reloadPaymentReconFromScopeFilters());
paymentProviderFilterEl.addEventListener("change", () => loadPaymentReconciliation(true));
paymentInvoiceStatusFilterEl.addEventListener("change", () => loadPaymentReconciliation(true));
paymentTransactionStatusFilterEl.addEventListener("change", () => loadPaymentReconciliation(true));
refreshNotificationsBtn.addEventListener("click", () => loadNotifications(true));
openNotificationDialogBtn.addEventListener("click", openNotificationDialog);
notificationWorkbenchStepsEl?.querySelectorAll("[data-notification-step]").forEach((button) => {
  button.addEventListener("click", () => setNotificationWorkbenchStep(button.dataset.notificationStep || "target"));
});
notificationCampaignTypeEl.addEventListener("change", () => {
  currentNotificationCampaignId = "";
  const match = (notificationOptions.templates || []).find((item) => item.code === notificationCampaignTypeEl.value);
  if (match) notificationTemplateEl.value = match.id;
});
notificationSchoolYearEl.addEventListener("change", () => {
  renderNotificationGradeOptions();
  renderNotificationClassOptions();
  syncAppContextFromActiveTab("notificationTab");
});
notificationGradeEl.addEventListener("change", renderNotificationClassOptions);
notificationPeriodEl.addEventListener("change", () => syncAppContextFromActiveTab("notificationTab"));
previewNotificationsBtn.addEventListener("click", previewNotifications);
saveNotificationCampaignBtn.addEventListener("click", saveNotificationCampaign);
retryNotificationRecipientsBtn.addEventListener("click", retrySelectedNotificationRecipients);
sendNotificationCampaignBtn.addEventListener("click", sendNotificationCampaign);

saveEmailConfigBtn.addEventListener("click", saveEmailConfig);
openEmailConfigDialogBtn.addEventListener("click", openEmailConfigDialog);
openCronConfigDialogBtn.addEventListener("click", openCronConfigDialog);
openNotificationCronConfigBtn.addEventListener("click", openCronConfigDialog);
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
