import * as XLSX from "xlsx";

import { cleanPaymentRow, parseAmount } from "./vietqr.js";

export const MAX_ROWS = 500;

const FIELD_DEFINITIONS = [
  ["student_name", "Người thanh toán / Học sinh"],
  ["parent_name", "Phụ huynh"],
  ["class_name", "Lớp"],
  ["bank_bin", "BIN ngân hàng"],
  ["bank_account", "Số tài khoản"],
  ["email", "Email"],
  ["amount", "Số tiền"],
  ["payment_items", "Danh sách khoản phí JSON"],
  ["fee_item", "Khoản phí (lấy tên cột)"],
  ["tuition_april", "Học phí tháng 04"],
  ["shuttle_april", "Phí xe tháng 04"],
  ["tuition_may", "Học phí tháng 05"],
  ["health_insurance", "Bảo hiểm y tế"],
  ["uniform_fee", "Đồng phục"],
  ["international_material", "Sách CTQT"],
  ["previous_fees", "Các khoản phí trước"],
  ["bill_number", "Mã hóa đơn / Tham chiếu"],
  ["note", "Nội dung chuyển khoản"],
];

export const IMPORT_FIELDS = FIELD_DEFINITIONS.map(([key, label]) => ({
  key,
  label,
  required: key === "bank_bin" || key === "bank_account",
}));

const LABEL_BY_KEY = new Map(FIELD_DEFINITIONS);

const ALIASES = {
  student_name: ["student_name", "student", "ten_hoc_sinh", "hoc_sinh", "ho_va_ten", "ho_ten", "ten_hs"],
  parent_name: ["parent_name", "parent", "ten_phu_huynh", "phu_huynh", "ten_ba_me", "ba_me", "ten_bo_me"],
  class_name: ["class_name", "class", "lop", "ten_lop", "lop_hoc"],
  bank_bin: ["bank_bin", "bin", "ma_bin", "bank_code", "bin_ngan_hang", "ma_ngan_hang"],
  bank_account: ["bank_account", "account", "account_number", "tai_khoan_ngan_hang", "stk", "so_tai_khoan"],
  email: ["email", "mail", "email_phu_huynh"],
  amount: ["amount", "so_tien", "hoc_phi", "tong_phi", "tong_hoc_phi", "so_tien_thanh_toan"],
  payment_items: ["payment_items", "items", "fee_items", "khoan_phi"],
  tuition_april: ["tuition_april", "hoc_phi_thang_04", "hoc_phi_t4"],
  shuttle_april: ["shuttle_april", "phi_xe_thang_04", "phi_xe_t4"],
  tuition_may: ["tuition_may", "hoc_phi_thang_05", "hoc_phi_t5"],
  health_insurance: ["health_insurance", "bao_hiem_y_te"],
  uniform_fee: ["uniform_fee", "dong_phuc"],
  international_material: ["international_material", "sach_ctqt", "sach"],
  previous_fees: ["previous_fees", "cac_khoan_phi_thang_truoc", "phi_thang_truoc"],
  bill_number: ["bill_number", "bill", "invoice", "invoice_number", "ma_hoa_don", "ma_tham_chieu", "ma_hd"],
  note: ["note", "noi_dung", "memo", "noi_dung_chuyen_khoan", "dien_giai"],
};

const FEE_FIELDS = {
  tuition_april: ["Tiền học phí Tháng 04", "Tuition fees for April"],
  shuttle_april: ["Phí xe đưa rước Tháng 04", "Shuttle fees for April"],
  tuition_may: ["Tiền học phí Tháng 05", "Tuition fees for May"],
  health_insurance: ["Bảo hiểm y tế", "Health Insurance"],
  uniform_fee: ["Đồng phục", "Uniform fee"],
  international_material: ["Sách CTQT", "International material"],
  previous_fees: ["Các khoản phí tháng trước", "Previous month's fees"],
};

export function normalizeHeader(value = "") {
  return String(value)
    .trim()
    .replace(/^\uFEFF/, "")
    .replace(/[ĐÐ]/g, "D")
    .replace(/đ/g, "d")
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "_")
    .replace(/^_+|_+$/g, "");
}

export function suggestMapping(headers) {
  const result = {};
  for (const header of headers) {
    const key = normalizeHeader(header);
    result[header] = Object.entries(ALIASES).find(([, aliases]) => aliases.includes(key))?.[0] || "";
  }
  return result;
}

export function validateMapping(mapping) {
  const errors = [];
  const sourcesByTarget = new Map();
  for (const [source, target] of Object.entries(mapping || {})) {
    if (!target || target === "fee_item") continue;
    const sources = sourcesByTarget.get(target) || [];
    sources.push(source);
    sourcesByTarget.set(target, sources);
  }
  for (const [target, sources] of sourcesByTarget.entries()) {
    if (sources.length > 1) errors.push(`Trường ${LABEL_BY_KEY.get(target) || target} đang được map từ nhiều cột`);
  }
  return errors;
}

function uniqueHeaders(values) {
  const seen = new Map();
  return values.map((value, index) => {
    const base = String(value ?? "").trim() || `Cột ${index + 1}`;
    const count = (seen.get(base) || 0) + 1;
    seen.set(base, count);
    return count === 1 ? base : `${base} (${count})`;
  });
}

export async function readSpreadsheet(file) {
  const data = await file.arrayBuffer();
  const workbook = XLSX.read(data, { type: "array", cellDates: false, raw: false });
  const firstSheetName = workbook.SheetNames[0];
  if (!firstSheetName) throw new Error("File không có worksheet");
  const matrix = XLSX.utils.sheet_to_json(workbook.Sheets[firstSheetName], {
    header: 1,
    defval: "",
    raw: false,
    blankrows: false,
  });
  if (matrix.length === 0) throw new Error("File không có dữ liệu");
  const headers = uniqueHeaders(matrix[0]);
  const records = matrix.slice(1)
    .map((record) => headers.map((_, index) => String(record[index] ?? "").trim()))
    .filter((record) => record.some(Boolean));
  if (records.length > MAX_ROWS) throw new Error(`Tối đa ${MAX_ROWS} dòng mỗi lần export`);
  return { name: file.name, sheetName: firstSheetName, headers, records };
}

function valueFor(table, record, mapping, target) {
  const sourceIndex = table.headers.findIndex((header) => mapping[header] === target);
  return sourceIndex >= 0 ? record[sourceIndex] : "";
}

function mappedFeeItems(table, record, mapping) {
  const items = [];
  table.headers.forEach((header, index) => {
    const target = mapping[header];
    if (target === "fee_item" && record[index] !== "") {
      items.push({ label: header, labelEn: header, amount: parseAmount(record[index]) });
    }
  });
  Object.entries(FEE_FIELDS).forEach(([target, [label, labelEn]]) => {
    const raw = valueFor(table, record, mapping, target);
    if (raw !== "") items.push({ label, labelEn, amount: parseAmount(raw) });
  });
  return items;
}

function jsonPaymentItems(raw) {
  if (!raw) return { items: [], errors: [] };
  try {
    const items = JSON.parse(raw);
    if (!Array.isArray(items)) throw new Error("not array");
    return { items, errors: [] };
  } catch {
    return { items: [], errors: ["Danh sách khoản phí JSON không hợp lệ"] };
  }
}

export function buildPaymentRows(table, mapping, defaults = {}) {
  const mappingErrors = validateMapping(mapping);
  if (mappingErrors.length > 0) throw new Error(mappingErrors.join("; "));

  return table.records.map((record, index) => {
    const parsedItems = jsonPaymentItems(valueFor(table, record, mapping, "payment_items"));
    const columnItems = mappedFeeItems(table, record, mapping);
    return cleanPaymentRow({
      id: `row-${String(index + 1).padStart(3, "0")}`,
      sourceRow: index + 2,
      studentName: valueFor(table, record, mapping, "student_name"),
      parentName: valueFor(table, record, mapping, "parent_name"),
      className: valueFor(table, record, mapping, "class_name"),
      bankBin: valueFor(table, record, mapping, "bank_bin") || defaults.bankBin,
      bankAccount: valueFor(table, record, mapping, "bank_account") || defaults.bankAccount,
      email: valueFor(table, record, mapping, "email"),
      amount: valueFor(table, record, mapping, "amount"),
      paymentItems: parsedItems.items.length > 0 ? parsedItems.items : columnItems,
      billNumber: valueFor(table, record, mapping, "bill_number"),
      note: valueFor(table, record, mapping, "note"),
      importErrors: parsedItems.errors,
    });
  });
}

export function previewValues(table, header, limit = 3) {
  const index = table.headers.indexOf(header);
  if (index < 0) return "";
  return table.records.map((row) => row[index]).filter(Boolean).slice(0, limit).join(" · ");
}
