import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

describe("payment recipient workflow", () => {
  it("places a mandatory recipient step before field mapping", () => {
    const html = readFileSync(new URL("../index.html", import.meta.url), "utf8");
    const recipientStep = html.indexOf('id="recipientSection"');
    const mappingStep = html.indexOf('id="mappingSection"');

    expect(recipientStep).toBeGreaterThan(-1);
    expect(mappingStep).toBeGreaterThan(recipientStep);
    expect(html).toContain('value="shared"');
    expect(html).toContain('value="per_row"');
    expect(html).toContain('id="confirmRecipientSetup"');
    expect(html).toContain('id="confirmRecipientExport"');
    expect(html).toContain("Tiền sẽ chuyển về đâu?");
  });
});
