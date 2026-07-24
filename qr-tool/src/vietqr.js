import { BANK_BY_BIN } from "./banks.js";

const NAPAS_AID = "A000000727";
const NAPAS_SERVICE_ACCOUNT = "QRIBFTTA";

export function onlyDigits(value = "") {
  return String(value).replace(/\D/g, "");
}

export function cleanAccount(value = "") {
  return String(value).trim().replace(/[^A-Za-z0-9]/g, "");
}

export function ascii(value = "") {
  return String(value)
    .replace(/[ĐÐ]/g, "D")
    .replace(/đ/g, "d")
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "");
}

export function cleanANS(value = "", maxLength = 25) {
  const cleaned = ascii(String(value).trim())
    .replace(/[^A-Za-z0-9 $%*+\-./:]/g, "")
    .replace(/\s+/g, " ")
    .trim();
  return maxLength > 0 ? cleaned.slice(0, maxLength) : cleaned;
}

export function parseAmount(value) {
  if (typeof value === "number") return Number.isFinite(value) ? Math.trunc(value) : 0;
  const normalized = String(value ?? "").trim().replace(/[,\.\s_]/g, "");
  if (!normalized) return 0;
  const amount = Number.parseInt(normalized, 10);
  return Number.isFinite(amount) ? amount : 0;
}

export function formatVND(value) {
  return `${new Intl.NumberFormat("vi-VN").format(Number(value) || 0)} ₫`;
}

export function crc16CCITTFalse(value) {
  let crc = 0xffff;
  for (const byte of new TextEncoder().encode(value)) {
    crc ^= byte << 8;
    for (let bit = 0; bit < 8; bit += 1) {
      crc = crc & 0x8000 ? ((crc << 1) ^ 0x1021) & 0xffff : (crc << 1) & 0xffff;
    }
  }
  return crc.toString(16).toUpperCase().padStart(4, "0");
}

function crc32IEEE(value) {
  let crc = 0xffffffff;
  for (const byte of new TextEncoder().encode(value)) {
    crc ^= byte;
    for (let bit = 0; bit < 8; bit += 1) {
      crc = (crc >>> 1) ^ (crc & 1 ? 0xedb88320 : 0);
    }
  }
  return (crc ^ 0xffffffff) >>> 0;
}

function tlv(id, value) {
  if (!/^\d{2}$/.test(id)) throw new Error(`invalid TLV id ${id}`);
  const text = String(value);
  if (!text) throw new Error(`empty TLV value for ${id}`);
  if (text.length > 99) throw new Error(`TLV ${id} length exceeds 99`);
  return `${id}${String(text.length).padStart(2, "0")}${text}`;
}

export function generateVietQR(request) {
  const bankBin = onlyDigits(request.bankBin);
  const accountNumber = cleanAccount(request.accountNumber);
  const amount = parseAmount(request.amount);
  const billNumber = cleanANS(request.billNumber, 25);
  const purpose = cleanANS(request.purpose, 25);

  if (bankBin.length !== 6) throw new Error("BIN ngân hàng phải có 6 chữ số");
  if (!BANK_BY_BIN.has(bankBin)) throw new Error("BIN ngân hàng không nằm trong danh sách VietQR");
  if (!accountNumber) throw new Error("Số tài khoản là bắt buộc");
  if (accountNumber.length > 19) throw new Error("Số tài khoản tối đa 19 ký tự");
  if (amount < 0) throw new Error("Số tiền không được âm");

  const beneficiary = tlv("01", tlv("00", bankBin) + tlv("01", accountNumber));
  const merchantAccount = tlv("38", tlv("00", NAPAS_AID) + beneficiary + tlv("02", NAPAS_SERVICE_ACCOUNT));
  const dynamic = Boolean(request.dynamic) || amount > 0;

  let payload = tlv("00", "01") + tlv("01", dynamic ? "12" : "11") + merchantAccount + tlv("53", "704");
  if (amount > 0) payload += tlv("54", String(amount));
  payload += tlv("58", "VN");

  let additional = "";
  if (billNumber) additional += tlv("01", billNumber);
  if (purpose) additional += tlv("08", purpose);
  if (additional) payload += tlv("62", additional);

  const withoutCRC = `${payload}6304`;
  return withoutCRC + crc16CCITTFalse(withoutCRC);
}

export function defaultBillNumber(row) {
  const parts = [
    row.studentName,
    row.parentName,
    row.bankBin,
    row.bankAccount,
    String(row.amount || 0),
    row.note,
  ];
  return `SUN${crc32IEEE(parts.join("|")).toString(16).toUpperCase().padStart(8, "0")}`;
}

function cleanPaymentItems(items = []) {
  return items
    .map((item) => {
      const label = String(item.label || item.labelEn || "").trim();
      const labelEn = String(item.labelEn || item.label || "").trim();
      return { label, labelEn, amount: parseAmount(item.amount) };
    })
    .filter((item) => item.label || item.labelEn || item.amount);
}

export function cleanPaymentRow(input = {}) {
  const row = {
    ...input,
    id: String(input.id || "").trim(),
    studentName: String(input.studentName || "").trim(),
    parentName: String(input.parentName || "").trim(),
    className: String(input.className || "").trim(),
    bankBin: onlyDigits(input.bankBin),
    bankAccount: cleanAccount(input.bankAccount),
    email: String(input.email || "").trim(),
    amount: parseAmount(input.amount),
    paymentItems: cleanPaymentItems(input.paymentItems),
    billNumber: cleanANS(input.billNumber, 25),
    note: cleanANS(input.note, 25),
    importErrors: Array.isArray(input.importErrors) ? input.importErrors : [],
  };
  if (row.paymentItems.length > 0) {
    row.amount = row.paymentItems.reduce((sum, item) => sum + item.amount, 0);
  }
  if (!row.note) row.note = cleanANS(`HP ${row.studentName}`, 25);
  if (!row.billNumber) row.billNumber = defaultBillNumber(row);
  return row;
}

export function validatePaymentRow(input) {
  const row = cleanPaymentRow(input);
  const errors = [...row.importErrors];
  if (row.bankBin.length !== 6) errors.push("BIN ngân hàng phải có 6 chữ số");
  else if (!BANK_BY_BIN.has(row.bankBin)) errors.push("BIN ngân hàng không nằm trong danh sách VietQR");
  if (!row.bankAccount) errors.push("Số tài khoản là bắt buộc");
  if (row.bankAccount.length > 19) errors.push("Số tài khoản tối đa 19 ký tự");
  if (row.amount < 0) errors.push("Số tiền không được âm");
  row.paymentItems.forEach((item, index) => {
    if (item.amount < 0) errors.push(`Khoản phí ${index + 1} không được âm`);
  });
  if (row.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(row.email)) errors.push("Email không hợp lệ");
  return errors;
}

export function buildQRItem(input) {
  const row = cleanPaymentRow(input);
  const errors = validatePaymentRow(row);
  const bank = BANK_BY_BIN.get(row.bankBin);
  if (errors.length > 0) return { ...row, bankName: bank?.shortName || "", vietqr: "", qrData: "", errors };
  try {
    const vietqr = generateVietQR({
      bankBin: row.bankBin,
      accountNumber: row.bankAccount,
      amount: row.amount,
      billNumber: row.billNumber,
      purpose: row.note,
      dynamic: row.amount > 0,
    });
    return { ...row, bankName: bank.shortName, vietqr, qrData: "", errors: [] };
  } catch (error) {
    return { ...row, bankName: bank?.shortName || "", vietqr: "", qrData: "", errors: [error.message] };
  }
}
