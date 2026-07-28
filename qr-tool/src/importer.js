import * as XLSX from "xlsx";

import { cleanPaymentRow, parseAmount } from "./vietqr.js";

export const MAX_ROWS = 500;

export const IMPORT_FIELD_GROUPS = [
  {
    key: "student_school",
    label: "Học sinh & trường",
    fields: [
      { key: "student_code", label: "Mã học sinh" },
      { key: "student_name", label: "Tên học sinh" },
      { key: "school_name", label: "Tên trường" },
      { key: "cohort", label: "Niên khóa" },
      { key: "year", label: "Năm" },
      { key: "class_name", label: "Lớp" },
    ],
  },
  {
    key: "parent",
    label: "Phụ huynh",
    fields: [
      { key: "parent_name", label: "Tên phụ huynh" },
      { key: "email", label: "Email" },
    ],
  },
  {
    key: "payment",
    label: "Thanh toán",
    fields: [
      { key: "bank_bin", label: "BIN ngân hàng", required: true },
      { key: "bank_account", label: "Số tài khoản", required: true },
      { key: "amount", label: "Số tiền" },
      { key: "bill_number", label: "Mã hóa đơn" },
      { key: "note", label: "Nội dung chuyển khoản" },
    ],
  },
  {
    key: "fees",
    label: "Khoản phí / nâng cao",
    fields: [
      { key: "fee_item", label: "Khoản phí — lấy tên cột", repeatable: true },
      { key: "payment_items", label: "Danh sách khoản phí JSON" },
    ],
  },
];

export const IMPORT_FIELDS = IMPORT_FIELD_GROUPS.flatMap((group) => group.fields);

const LABEL_BY_KEY = new Map(IMPORT_FIELDS.map((field) => [field.key, field.label]));

const ALIASES = {
  student_code: ["student_code", "student_id", "student_no", "ma_hoc_sinh", "ma_hs"],
  student_name: ["student_name", "student", "ten_hoc_sinh", "hoc_sinh", "ho_va_ten", "ho_ten", "ten_hs"],
  school_name: ["school_name", "school", "ten_truong", "truong"],
  cohort: ["cohort", "nien_khoa", "khoa_hoc"],
  year: ["year", "nam", "nam_hoc", "grade", "grade_name", "khoi"],
  parent_name: ["parent_name", "parent", "ten_phu_huynh", "phu_huynh", "ten_ba_me", "ba_me", "ten_bo_me"],
  class_name: ["class_name", "class", "lop", "ten_lop", "lop_hoc"],
  bank_bin: ["bank_bin", "bin", "ma_bin", "bank_code", "bin_ngan_hang", "ma_ngan_hang"],
  bank_account: ["bank_account", "account", "account_number", "tai_khoan_ngan_hang", "stk", "so_tai_khoan"],
  email: ["email", "mail", "email_phu_huynh"],
  amount: ["amount", "so_tien", "hoc_phi", "tong_phi", "tong_hoc_phi", "so_tien_thanh_toan"],
  payment_items: ["payment_items", "items", "fee_items", "danh_sach_khoan_phi", "danh_sach_khoan_phi_json"],
  fee_item: [
    "fee_item", "khoan_phi", "khoan_thu",
    "tuition_april", "hoc_phi_thang_04", "hoc_phi_t4",
    "shuttle_april", "phi_xe_thang_04", "phi_xe_t4",
    "tuition_may", "hoc_phi_thang_05", "hoc_phi_t5",
    "health_insurance", "bao_hiem_y_te",
    "uniform_fee", "dong_phuc",
    "international_material", "sach_ctqt", "sach",
    "previous_fees", "cac_khoan_phi_thang_truoc", "phi_thang_truoc",
  ],
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
  const recipientMode = defaults.mode === "per_row" ? "per_row" : "shared";

  return table.records.map((record, index) => {
    const parsedItems = jsonPaymentItems(valueFor(table, record, mapping, "payment_items"));
    const columnItems = mappedFeeItems(table, record, mapping);
    return cleanPaymentRow({
      id: `row-${String(index + 1).padStart(3, "0")}`,
      sourceRow: index + 2,
      studentCode: valueFor(table, record, mapping, "student_code"),
      studentName: valueFor(table, record, mapping, "student_name"),
      schoolName: valueFor(table, record, mapping, "school_name"),
      cohort: valueFor(table, record, mapping, "cohort"),
      year: valueFor(table, record, mapping, "year"),
      parentName: valueFor(table, record, mapping, "parent_name"),
      className: valueFor(table, record, mapping, "class_name"),
      bankBin: recipientMode === "per_row" ? valueFor(table, record, mapping, "bank_bin") : defaults.bankBin,
      bankAccount: recipientMode === "per_row" ? valueFor(table, record, mapping, "bank_account") : defaults.bankAccount,
      accountName: recipientMode === "shared" ? defaults.accountName : "",
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
