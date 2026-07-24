import { describe, expect, it } from "vitest";

import {
  DEFAULT_EMAIL_TEMPLATE,
  buildEml,
  renderEmailTemplate,
  templateForGmailMerge,
} from "../src/email.js";

const item = {
  id: "row-001",
  studentName: "Nguyễn An",
  parentName: "Nguyễn Văn Bình",
  className: "3.02",
  bankName: "VietinBank",
  bankAccount: "0011001932418",
  email: "parent@example.com",
  amount: 120000,
  paymentItems: [{ label: "Học phí", labelEn: "Tuition", amount: 120000 }],
  billNumber: "SUN001",
  note: "HP NGUYEN AN",
};

describe("email template export", () => {
  it("renders escaped merge fields and trusted payment blocks", () => {
    const rendered = renderEmailTemplate(DEFAULT_EMAIL_TEMPLATE, {
      ...item,
      parentName: "<script>alert(1)</script>",
    }, { qrSrc: "data:image/png;base64,AAAA" });
    expect(rendered.subject).toContain("Nguyễn An");
    expect(rendered.html).not.toContain("<script>");
    expect(rendered.html).toContain("&lt;script&gt;");
    expect(rendered.html).toContain("data:image/png;base64,AAAA");
    expect(rendered.text).toContain("120.000 ₫");
  });

  it("builds an unsent MIME message whose CID matches the inline QR", () => {
    const rendered = renderEmailTemplate(DEFAULT_EMAIL_TEMPLATE, item, { qrSrc: "cid:qr-row-001" });
    const eml = buildEml({
      to: item.email,
      subject: rendered.subject,
      html: rendered.html,
      text: rendered.text,
      qrBase64: "QUJDRA==",
      contentId: "qr-row-001",
      qrFilename: "SUN001.png",
    });
    expect(eml).toContain("X-Unsent: 1");
    expect(eml).toContain("Content-ID: <qr-row-001>");
    expect(eml).toContain("Content-Disposition: inline; filename=\"SUN001.png\"");
  });

  it("converts supported fields to Gmail merge tags and flags per-row QR", () => {
    const gmail = templateForGmailMerge(DEFAULT_EMAIL_TEMPLATE);
    expect(gmail.html).toContain("@StudentName");
    expect(gmail.hasPerRecipientQR).toBe(true);
    expect(gmail.html).not.toContain("{{qr_image}}");
  });
});
