import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

describe("connected Gmail Sheet recovery", () => {
  it("shows a visible upgrade action before the collapsed management area", () => {
    const html = readFileSync(new URL("../index.html", import.meta.url), "utf8");
    const returningFlow = html.indexOf('id="gmailReturningFlow"');
    const upgradeAction = html.indexOf('id="downloadGmailUpgradeConnected"');
    const upgradeSteps = html.indexOf('id="gmailUpgradeSteps"');
    const managementArea = html.indexOf("Đổi Sheet, ngắt kết nối hoặc tải lại bộ cài");

    expect(returningFlow).toBeGreaterThan(-1);
    expect(upgradeAction).toBeGreaterThan(returningFlow);
    expect(upgradeSteps).toBeGreaterThan(upgradeAction);
    expect(managementArea).toBeGreaterThan(upgradeSteps);
    expect(html).toContain("Không thấy “0. Nhập dữ liệu mới” trong Sheet?");
  });
});
