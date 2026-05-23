---
name: abcsun-frontend-ui
description: Frontend UI workflow for the ABC SUN embedded vanilla web app. Use when changing web/index.html, web/styles.css, web/app.js, tabs, table editing, QR preview, email config UI, cron UI, client-side API calls, or responsive layout.
---

# ABC SUN Frontend UI

Use this for browser-facing changes in `web/*`.

## Architecture

- No frontend build step exists; `web/*` is embedded by Go via `//go:embed`.
- `web/index.html` owns structure.
- `web/styles.css` owns the operational dashboard look.
- `web/app.js` owns state, rendering, API calls, and event handlers.
- Server data contracts are in Go handlers; update Go tests when client behavior depends on API shape.

## Design Rules

- Keep the app as a dense operational tool, not a landing page.
- Match the existing restrained palette, 6-8px radii, Material Symbols icons, tabs, tables, and panels.
- Prefer compact controls with clear labels. Do not add marketing copy or tutorial text.
- Keep fee editing, QR preview, email preview, and cron controls scannable.
- Avoid layout shifts when rows, fee items, status pills, or buttons update.
- Keep mobile usable: tables may scroll horizontally, but controls must not overlap.

## Implementation Workflow

1. Trace the current UI state path in `web/app.js` before editing.
2. If API fields change, update Go structs/handlers/tests and the JS together.
3. Escape user-controlled values with existing helpers such as `escapeHtml` and `escapeAttr`.
4. Preserve confirmation prompts for real email sends and cron runs.
5. Keep provider-specific UI visibility in sync with backend validation rules.

## Verification

Run:

```sh
go test ./...
```

For UI changes:

1. Start the app with `go run .` or `PORT=8081 go run .`.
2. Open the local URL.
3. Exercise the changed tab or workflow.
4. Check browser console/network errors when browser tooling is available.

For email UI changes, validate preview or dry-run only unless the user explicitly asks for a real send.
