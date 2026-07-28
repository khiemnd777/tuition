/**
 * DEKISUGI Gmail Sender — code gắn với Google Sheet mẫu.
 * Dữ liệu chỉ nằm trong bản sao Google Sheet của người dùng.
 * @OnlyCurrentDoc
 */

var DEKISUGI_CONFIG = Object.freeze({
  SHEET_NAME: "EMAILS",
  MAX_ROWS: 500,
  MAX_CELL_LENGTH: 49000,
  MAX_PER_RUN: 90,
  RESERVED_QUOTA: 10,
  SEND_DELAY_MS: 400,
  QR_CONTENT_ID: "dekisugi_qr",
  SCHEDULE_HANDLER: "runScheduledBatch"
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
    .addItem("Mở bảng điều khiển", "showSidebar")
    .addToUi();
}

function setup() {
  ensureWelcomeSheet_();
  ensureSheet_();
  formatSheet_();
  MailApp.getRemainingDailyQuota();
  onOpen();
  showSidebar();
}

function showSidebar() {
  var output = HtmlService.createHtmlOutputFromFile("Sidebar")
    .setTitle("DEKISUGI Email");
  SpreadsheetApp.getUi().showSidebar(output);
}

function getSidebarState() {
  return getSidebarState_();
}

function importDataset(payload) {
  var lock = LockService.getDocumentLock();
  if (!lock.tryLock(5000)) throw new Error("Một thao tác khác đang chạy. Vui lòng thử lại sau vài giây.");

  try {
    validateDataset_(payload);
    removeScheduleTrigger_();

    var sheet = ensureSheet_();
    if (sheet.getFilter()) sheet.getFilter().remove();
    sheet.clear();
    sheet.showColumns(1, sheet.getMaxColumns());
    sheet.getRange(1, 1, 1, DEKISUGI_HEADERS.length).setValues([DEKISUGI_HEADERS]);

    if (payload.rows.length) {
      var values = payload.rows.map(function (row) {
        return DEKISUGI_HEADERS.map(function (header) {
          return safeSheetText_(row[header]);
        });
      });
      var range = sheet.getRange(2, 1, values.length, DEKISUGI_HEADERS.length);
      range.setNumberFormat("@");
      range.setValues(values);
    }

    var properties = PropertiesService.getDocumentProperties();
    properties.deleteAllProperties();
    properties.setProperties({
      DATA_IMPORTED_AT: new Date().toISOString(),
      DATA_EXPORTED_AT: String(payload.exportedAt || ""),
      DATA_ROW_COUNT: String(payload.rows.length)
    });

    formatSheet_();
    validateRows_();
    SpreadsheetApp.flush();
    return withMessage_(getSidebarState_(), "Đã nhập " + payload.rows.length + " dòng vào Google Sheet.");
  } finally {
    lock.releaseLock();
  }
}

function sendTestEmail(recipient) {
  recipient = String(recipient || "").trim();
  if (!isEmail_(recipient)) throw new Error("Email nhận bản thử không hợp lệ.");

  var lock = LockService.getDocumentLock();
  if (!lock.tryLock(5000)) throw new Error("Một thao tác khác đang chạy. Vui lòng thử lại sau vài giây.");
  try {
    validateRows_();
    var context = getContext_();
    var pending = pendingRowIndexes_(context);
    if (!pending.length) throw new Error("Không có dòng READY để gửi thử.");
    if (MailApp.getRemainingDailyQuota() < 1) throw new Error("Tài khoản đã hết quota gửi email hôm nay.");

    var record = recordForRow_(context, pending[0]);
    sendRecord_(record, recipient);
    PropertiesService.getDocumentProperties().setProperties({
      TEST_SENT_AT: new Date().toISOString(),
      TEST_RECIPIENT: recipient
    });
    return withMessage_(getSidebarState_(), "Đã gửi một email thử tới " + recipient + ". Danh sách thật chưa được gửi.");
  } finally {
    lock.releaseLock();
  }
}

function sendPendingNow() {
  requireSuccessfulTest_();
  var result = sendBatch_();
  return withMessage_(getSidebarState_(), batchMessage_(result));
}

function enableDailySchedule(hour) {
  requireSuccessfulTest_();
  hour = Number(hour);
  if (!Number.isInteger(hour) || hour < 0 || hour > 23) throw new Error("Giờ gửi phải nằm trong khoảng 00–23.");

  validateRows_();
  var context = getContext_();
  if (!pendingRowIndexes_(context).length) throw new Error("Không còn dòng READY để đặt lịch gửi.");

  removeScheduleTrigger_();
  var trigger = ScriptApp.newTrigger(DEKISUGI_CONFIG.SCHEDULE_HANDLER)
    .timeBased()
    .atHour(hour)
    .everyDays(1)
    .inTimezone(Session.getScriptTimeZone())
    .create();
  PropertiesService.getDocumentProperties().setProperties({
    SCHEDULE_TRIGGER_ID: trigger.getUniqueId(),
    SCHEDULE_HOUR: String(hour),
    SCHEDULE_ENABLED_AT: new Date().toISOString()
  });
  return withMessage_(getSidebarState_(), "Đã bật lịch gửi mỗi ngày khoảng " + pad2_(hour) + ":00.");
}

function disableDailySchedule() {
  removeScheduleTrigger_();
  return withMessage_(getSidebarState_(), "Đã tạm dừng lịch gửi tự động.");
}

function runScheduledBatch() {
  try {
    var result = sendBatch_();
    var properties = PropertiesService.getDocumentProperties();
    properties.setProperties({
      LAST_SCHEDULE_RUN_AT: new Date().toISOString(),
      LAST_SCHEDULE_RESULT: batchMessage_(result)
    });
    if (result.pending === 0) {
      removeScheduleTrigger_();
      properties.setProperty("SCHEDULE_COMPLETED_AT", new Date().toISOString());
    }
  } catch (error) {
    PropertiesService.getDocumentProperties().setProperties({
      LAST_SCHEDULE_RUN_AT: new Date().toISOString(),
      LAST_SCHEDULE_ERROR: cleanError_(error)
    });
    console.error(error);
  }
}

function sendBatch_() {
  var lock = LockService.getDocumentLock();
  if (!lock.tryLock(5000)) throw new Error("Một lượt gửi khác đang chạy. Vui lòng thử lại sau.");
  try {
    validateRows_();
    var context = getContext_();
    var pending = pendingRowIndexes_(context);
    if (!pending.length) return { sent: 0, failed: 0, pending: 0 };

    var remainingQuota = MailApp.getRemainingDailyQuota();
    var safeQuota = Math.max(0, Math.min(
      DEKISUGI_CONFIG.MAX_PER_RUN,
      remainingQuota - DEKISUGI_CONFIG.RESERVED_QUOTA
    ));
    if (safeQuota < 1) throw new Error("Quota còn lại không đủ. Công cụ luôn chừa 10 email cho nhu cầu gửi thông thường.");

    var sendCount = Math.min(safeQuota, pending.length);
    var sent = 0;
    var failed = 0;
    for (var index = 0; index < sendCount; index += 1) {
      var rowIndex = pending[index];
      var record = recordForRow_(context, rowIndex);
      writeCell_(context, rowIndex, "Status", "SENDING");
      writeCell_(context, rowIndex, "Error", "");
      SpreadsheetApp.flush();
      try {
        sendRecord_(record, record.Email);
        writeCell_(context, rowIndex, "Status", "SENT");
        writeCell_(context, rowIndex, "SentAt", new Date());
        sent += 1;
      } catch (error) {
        var message = cleanError_(error);
        writeCell_(context, rowIndex, "Status", "ERROR");
        writeCell_(context, rowIndex, "Error", message);
        failed += 1;
        if (/quota|limit|too many|service invoked/i.test(message)) break;
      }
      SpreadsheetApp.flush();
      Utilities.sleep(DEKISUGI_CONFIG.SEND_DELAY_MS);
    }
    return {
      sent: sent,
      failed: failed,
      pending: pendingRowIndexes_(getContext_()).length
    };
  } finally {
    lock.releaseLock();
  }
}

function getSidebarState_() {
  var properties = PropertiesService.getDocumentProperties();
  var counts = { total: 0, ready: 0, sent: 0, errors: 0, skipped: 0, sending: 0 };
  var hasData = false;
  try {
    var context = getContext_();
    hasData = context.values.length > 1;
    for (var rowIndex = 1; rowIndex < context.values.length; rowIndex += 1) {
      counts.total += 1;
      var status = String(recordForRow_(context, rowIndex).Status || "").toUpperCase();
      if (status === "READY") counts.ready += 1;
      else if (status === "SENT") counts.sent += 1;
      else if (status === "SKIP") counts.skipped += 1;
      else if (status === "SENDING") counts.sending += 1;
      else counts.errors += 1;
    }
  } catch (error) {
    hasData = false;
  }

  var scheduleHour = properties.getProperty("SCHEDULE_HOUR");
  var triggerActive = scheduleTriggerExists_();
  return {
    hasData: hasData,
    counts: counts,
    quota: MailApp.getRemainingDailyQuota(),
    testSent: Boolean(properties.getProperty("TEST_SENT_AT")),
    testRecipient: properties.getProperty("TEST_RECIPIENT") || "",
    testSentAt: properties.getProperty("TEST_SENT_AT") || "",
    scheduleActive: triggerActive,
    scheduleHour: scheduleHour === null ? 8 : Number(scheduleHour),
    scheduleLabel: triggerActive ? "Mỗi ngày khoảng " + pad2_(scheduleHour) + ":00" : "Chưa bật",
    lastScheduleRunAt: properties.getProperty("LAST_SCHEDULE_RUN_AT") || "",
    lastScheduleResult: properties.getProperty("LAST_SCHEDULE_RESULT") || "",
    lastScheduleError: properties.getProperty("LAST_SCHEDULE_ERROR") || ""
  };
}

function validateDataset_(payload) {
  if (!payload || payload.kind !== "dekisugi.gmail-data" || Number(payload.schemaVersion) !== 1) {
    throw new Error("Đây không phải file dữ liệu Gmail do DEKISUGI xuất.");
  }
  if (!Array.isArray(payload.rows)) throw new Error("File dữ liệu không có danh sách email.");
  if (payload.rows.length < 1) throw new Error("File dữ liệu không có dòng nào để nhập.");
  if (payload.rows.length > DEKISUGI_CONFIG.MAX_ROWS) throw new Error("Chỉ nhận tối đa 500 dòng mỗi lần import.");

  payload.rows.forEach(function (row, rowIndex) {
    if (!row || typeof row !== "object" || Array.isArray(row)) {
      throw new Error("Dòng " + (rowIndex + 1) + " không hợp lệ.");
    }
    DEKISUGI_HEADERS.forEach(function (header) {
      var value = row[header] === undefined || row[header] === null ? "" : String(row[header]);
      if (value.length > DEKISUGI_CONFIG.MAX_CELL_LENGTH) {
        throw new Error("Dòng " + (rowIndex + 1) + ", cột " + header + " vượt giới hạn Google Sheets.");
      }
    });
  });
}

function validateRows_() {
  var context = getContext_();
  for (var rowIndex = 1; rowIndex < context.values.length; rowIndex += 1) {
    var record = recordForRow_(context, rowIndex);
    if (record.Status === "SENT" || record.Status === "SENDING") continue;
    var errors = validateRecord_(record);
    var skipped = String(record.Send || "").toUpperCase() === "NO" && !errors.length;
    writeCell_(context, rowIndex, "Status", errors.length ? "ERROR" : (skipped ? "SKIP" : "READY"));
    writeCell_(context, rowIndex, "Error", errors.join("; "));
  }
  SpreadsheetApp.flush();
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

function ensureSheet_() {
  var spreadsheet = SpreadsheetApp.getActive();
  var sheet = spreadsheet.getSheetByName(DEKISUGI_CONFIG.SHEET_NAME);
  if (!sheet) sheet = spreadsheet.insertSheet(DEKISUGI_CONFIG.SHEET_NAME);
  if (sheet.getMaxColumns() < DEKISUGI_HEADERS.length) {
    sheet.insertColumnsAfter(sheet.getMaxColumns(), DEKISUGI_HEADERS.length - sheet.getMaxColumns());
  }
  if (!sheet.getLastRow()) sheet.getRange(1, 1, 1, DEKISUGI_HEADERS.length).setValues([DEKISUGI_HEADERS]);
  return sheet;
}

function ensureWelcomeSheet_() {
  var spreadsheet = SpreadsheetApp.getActive();
  var sheet = spreadsheet.getSheetByName("BAT_DAU");
  if (!sheet) sheet = spreadsheet.insertSheet("BAT_DAU", 0);
  sheet.getRange("A1:D12").breakApart();
  sheet.clear();
  sheet.getRange("A1:D1").merge().setValue("DEKISUGI EMAIL — BẮT ĐẦU TỪ ĐÂY");
  sheet.getRange("A3:D8").setValues([
    ["Bước", "Bạn bấm ở đâu", "Bạn cần làm gì", "Kết quả đúng"],
    ["1", "Thanh menu phía trên", "Chọn DEKISUGI Email → Mở bảng điều khiển", "Sidebar xuất hiện bên phải"],
    ["2", "Sidebar · Nhập danh sách", "Chọn file DEKISUGI_GMAIL_DATA.json vừa tải", "Thấy số dòng Sẵn sàng và Lỗi"],
    ["3", "Sidebar · Gửi một email thử", "Nhập email của chính bạn và bấm Gửi 1 email thử", "Email thử tới inbox; danh sách thật chưa gửi"],
    ["4", "Kiểm tra inbox", "Mở email thử, kiểm tra học sinh, số tiền và QR", "Nội dung và QR đều đúng"],
    ["5", "Sidebar · Gửi thật hoặc đặt lịch", "Chọn Gửi các email đang chờ hoặc Bật lịch", "Dòng thành công chuyển thành SENT"]
  ]);
  sheet.getRange("A10:D10").merge().setValue("Lần đầu sử dụng, Google sẽ yêu cầu bạn cấp quyền cho chính bản sao Sheet này. Chỉ tiếp tục khi địa chỉ là script.google.com và file do bạn vừa tạo bản sao.");
  sheet.getRange("A12:D12").merge().setValue("Không gửi hoặc chia sẻ Google Sheet này cho phụ huynh vì Sheet chứa email và dữ liệu thanh toán.");
  sheet.getRange("A1:D1").setBackground("#2b50a1").setFontColor("#ffffff").setFontWeight("bold").setFontSize(14);
  sheet.getRange("A3:D3").setBackground("#e8eefb").setFontWeight("bold");
  sheet.getRange("A10:D10").setBackground("#fff7e6").setFontColor("#7a5200");
  sheet.getRange("A12:D12").setBackground("#fff1f0").setFontColor("#9d261d").setFontWeight("bold");
  sheet.getRange("A1:D12").setWrap(true).setVerticalAlignment("middle");
  sheet.setColumnWidth(1, 70);
  sheet.setColumnWidth(2, 210);
  sheet.setColumnWidth(3, 330);
  sheet.setColumnWidth(4, 260);
  sheet.setRowHeight(1, 34);
  spreadsheet.setActiveSheet(sheet);
  spreadsheet.moveActiveSheet(1);
  return sheet;
}

function getContext_() {
  var sheet = SpreadsheetApp.getActive().getSheetByName(DEKISUGI_CONFIG.SHEET_NAME);
  if (!sheet || sheet.getLastRow() < 1) throw new Error("Chưa có dữ liệu. Hãy import file DEKISUGI_GMAIL_DATA.json trước.");
  var values = sheet.getDataRange().getDisplayValues();
  var index = {};
  values[0].forEach(function (header, columnIndex) { index[String(header).trim()] = columnIndex; });
  DEKISUGI_HEADERS.forEach(function (header) {
    if (index[header] === undefined) throw new Error("Thiếu cột bắt buộc: " + header);
  });
  return { sheet: sheet, values: values, index: index };
}

function formatSheet_() {
  var sheet = ensureSheet_();
  sheet.setFrozenRows(1);
  sheet.getRange(1, 1, 1, DEKISUGI_HEADERS.length)
    .setBackground("#2b50a1")
    .setFontColor("#ffffff")
    .setFontWeight("bold");
  if (!sheet.getFilter()) {
    sheet.getRange(1, 1, Math.max(1, sheet.getLastRow()), DEKISUGI_HEADERS.length).createFilter();
  }
  ["HtmlBody", "TextBody", "QrBase64", "QrFilename"].forEach(function (header) {
    var column = DEKISUGI_HEADERS.indexOf(header) + 1;
    if (!sheet.isColumnHiddenByUser(column)) sheet.hideColumns(column);
  });
  sheet.autoResizeColumns(1, 13);
  sheet.setColumnWidth(3, 220);
  sheet.setColumnWidth(13, 320);
  sheet.setColumnWidth(18, 90);
  sheet.setColumnWidth(20, 300);
}

function pendingRowIndexes_(context) {
  var rows = [];
  for (var rowIndex = 1; rowIndex < context.values.length; rowIndex += 1) {
    var record = recordForRow_(context, rowIndex);
    if (record.Status === "READY" && String(record.Send).toUpperCase() !== "NO") rows.push(rowIndex);
  }
  return rows;
}

function recordForRow_(context, rowIndex) {
  var record = {};
  Object.keys(context.index).forEach(function (header) {
    record[header] = context.values[rowIndex][context.index[header]];
  });
  return record;
}

function writeCell_(context, rowIndex, header, value) {
  context.sheet.getRange(rowIndex + 1, context.index[header] + 1).setValue(value);
  context.values[rowIndex][context.index[header]] = String(value || "");
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

function requireSuccessfulTest_() {
  if (!PropertiesService.getDocumentProperties().getProperty("TEST_SENT_AT")) {
    throw new Error("Hãy gửi một email thử và kiểm tra nội dung trước khi gửi thật hoặc bật lịch.");
  }
}

function removeScheduleTrigger_() {
  ScriptApp.getProjectTriggers().forEach(function (trigger) {
    if (trigger.getHandlerFunction() === DEKISUGI_CONFIG.SCHEDULE_HANDLER) ScriptApp.deleteTrigger(trigger);
  });
  var properties = PropertiesService.getDocumentProperties();
  ["SCHEDULE_TRIGGER_ID", "SCHEDULE_HOUR", "SCHEDULE_ENABLED_AT"].forEach(function (key) {
    properties.deleteProperty(key);
  });
}

function scheduleTriggerExists_() {
  return ScriptApp.getProjectTriggers().some(function (trigger) {
    return trigger.getHandlerFunction() === DEKISUGI_CONFIG.SCHEDULE_HANDLER;
  });
}

function safeSheetText_(value) {
  var text = value === undefined || value === null ? "" : String(value);
  return /^[=+\-@]/.test(text) ? "'" + text : text;
}

function isEmail_(value) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(value || "").trim());
}

function cleanError_(error) {
  return String(error && error.message ? error.message : error || "Lỗi chưa xác định")
    .replace(/[\r\n]+/g, " ")
    .slice(0, 500);
}

function batchMessage_(result) {
  return "Đã gửi: " + result.sent + " · Lỗi: " + result.failed + " · Còn chờ: " + result.pending;
}

function withMessage_(state, message) {
  state.message = message;
  return state;
}

function pad2_(value) {
  value = Number(value);
  return value < 10 ? "0" + value : String(value);
}
