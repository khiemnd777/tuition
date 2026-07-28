import { describe, expect, it } from "vitest";

import {
  GMAIL_SHEET_STORAGE_KEY,
  clearPersonalGmailSheetURL,
  loadPersonalGmailSheetURL,
  normalizePersonalGmailSheetURL,
  savePersonalGmailSheetURL,
} from "../src/gmail-sheet.js";

function memoryStorage() {
  const values = new Map();
  return {
    getItem: (key) => values.get(key) || null,
    setItem: (key, value) => values.set(key, value),
    removeItem: (key) => values.delete(key),
  };
}

describe("personal Gmail Sheet connection", () => {
  it("normalizes a user-owned Google Sheet edit link", () => {
    expect(normalizePersonalGmailSheetURL(
      "https://docs.google.com/spreadsheets/d/abc_DEF-123/edit?gid=456#gid=456",
    )).toBe("https://docs.google.com/spreadsheets/d/abc_DEF-123/edit");
  });

  it("rejects template copy and non-Google links", () => {
    expect(() => normalizePersonalGmailSheetURL(
      "https://docs.google.com/spreadsheets/d/template-id/copy",
    )).toThrow("link tạo bản sao");
    expect(() => normalizePersonalGmailSheetURL(
      "https://example.com/spreadsheets/d/abc/edit",
    )).toThrow("docs.google.com");
  });

  it("stores only the personal Sheet URL and can disconnect", () => {
    const storage = memoryStorage();
    const normalized = savePersonalGmailSheetURL(
      storage,
      "https://docs.google.com/spreadsheets/d/sheet-123/edit#gid=0",
    );

    expect(normalized).toBe("https://docs.google.com/spreadsheets/d/sheet-123/edit");
    expect(storage.getItem(GMAIL_SHEET_STORAGE_KEY)).toBe(normalized);
    expect(loadPersonalGmailSheetURL(storage)).toBe(normalized);

    clearPersonalGmailSheetURL(storage);
    expect(loadPersonalGmailSheetURL(storage)).toBe("");
  });
});
