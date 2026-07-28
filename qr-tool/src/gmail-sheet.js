export const GMAIL_SHEET_STORAGE_KEY = "dekisugi.gmail-sheet-url.v1";

export function normalizePersonalGmailSheetURL(value) {
  let url;
  try {
    url = new URL(String(value || "").trim());
  } catch {
    throw new Error("Link Google Sheet không hợp lệ");
  }

  if (url.protocol !== "https:" || url.hostname !== "docs.google.com") {
    throw new Error("Hãy dùng link Google Sheet bắt đầu bằng https://docs.google.com");
  }

  const match = url.pathname.match(/^\/spreadsheets\/d\/([A-Za-z0-9_-]+)(?:\/|$)/);
  if (!match) throw new Error("Không tìm thấy mã Google Sheet trong link này");
  if (/\/copy\/?$/.test(url.pathname)) {
    throw new Error("Đây là link tạo bản sao. Hãy mở Sheet của bạn rồi copy link có chữ /edit");
  }

  return `https://docs.google.com/spreadsheets/d/${match[1]}/edit`;
}

export function loadPersonalGmailSheetURL(storage) {
  try {
    const value = storage?.getItem(GMAIL_SHEET_STORAGE_KEY);
    return value ? normalizePersonalGmailSheetURL(value) : "";
  } catch {
    return "";
  }
}

export function savePersonalGmailSheetURL(storage, value) {
  const normalized = normalizePersonalGmailSheetURL(value);
  if (!storage) throw new Error("Trình duyệt không cho phép ghi nhớ Sheet. Hãy kiểm tra chế độ riêng tư");
  try {
    storage.setItem(GMAIL_SHEET_STORAGE_KEY, normalized);
  } catch {
    throw new Error("Trình duyệt không cho phép ghi nhớ Sheet. Hãy kiểm tra chế độ riêng tư");
  }
  return normalized;
}

export function clearPersonalGmailSheetURL(storage) {
  try {
    storage?.removeItem(GMAIL_SHEET_STORAGE_KEY);
  } catch {
    // The connection is already cleared for this page even if browser storage is unavailable.
  }
}
