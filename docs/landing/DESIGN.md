# DEKISUGI Landing Page Design

## Purpose

The Landing Page is the public pre-auth entry point for DEKISUGI Finance Hub. It should help a school operator understand the product, review available subscription plans, and choose either Login or Enrollment without entering the authenticated admin workspace first.

## Audience

- School owners or finance operators evaluating DEKISUGI.
- Existing tenant users returning to Login.
- Internal operations users who need a public intake path for new schools.

## Page Structure

### Header

- Brand: `DEKISUGI Finance Hub`.
- Navigation anchors: overview, subscriptions, contact.
- Primary actions: Login and Enrollment.
- Header must remain compact and readable on mobile.

### Welcome / Greeting

- First viewport must clearly identify `DEKISUGI Finance Hub`.
- Supporting copy should describe school payment operations: invoices, VietQR, reconciliation, parent notifications, and tenant subscriptions.
- Use a real app-rendered visual asset when possible, such as the public QR PNG endpoint or an HTML/CSS operations artifact that previews the product workflow.
- The hero should leave a hint of following content visible where viewport height allows.

### Operational Areas

Show concise cards for:

- Tuition collection and VietQR.
- Payment reconciliation.
- Parent notification/email.
- Deployment setup and subscription operations.

### Subscriptions

- Subscription cards must render from real system data, not static mock plans.
- Source endpoint for the public page: `GET /api/v1/public/subscription-plans`.
- Protected operator endpoint remains `GET /api/v1/subscriptions/plans`.
- If plans are missing or the database is unavailable, show an empty/error state.
- Do not invent pricing, quotas, or plan names not returned by the API.

### Access

- Login and Enrollment live in the landing access section.
- Preserve existing form IDs in `web/index.html` because `web/app.js` wires behavior by selector.
- Enrollment submits intake to `/api/v1/intake`; the internal operations team later creates the tenant, owner, provider config, and subscription.
- Customer-facing copy should not expose internal control-plane role names; use phrases like `đội ngũ phụ trách`, `bộ phận hỗ trợ`, or `cấu hình hệ thống`.

### Footer

- Show contact:
  - `khiemnd777@gmail.com`
  - `0974322365`
- Show: `Powered by KNA Software - knasoftware.com`.
- Keep Terms and Privacy links available through the auth form or footer area.

## Interaction Rules

- Hover transforms should make cards/sections feel active:
  - Cards may use `translateY(-4px to -6px)` with a subtle shadow.
  - Section hover may use a smaller transform and must not create text overlap.
- Do not use native `window.alert`, `window.confirm`, or `window.prompt`.
- Do not hide Login behind a marketing funnel. Login must remain a top-level action.
- Enrollment must be reachable from the header and hero.

## Visual Style

- Use the existing CSS variables in `web/styles.css`.
- Use Material Symbols icons already loaded in `web/index.html`.
- Cards should use 8px radius or less.
- Keep the palette balanced with restrained green, blue, amber, and neutral accents.
- Avoid one-note palettes, decorative gradient blobs, and nested cards.

## Responsive Rules

- Desktop: header, hero, section grids, and auth form should be scannable without overlap.
- Tablet: section grids may collapse to two columns.
- Mobile: cards, proof chips, hero actions, and auth controls collapse to one column.
- Button text may wrap; controls must keep stable dimensions and remain tappable.

## Accessibility

- Use semantic landmarks: header, main, section, footer.
- Anchor navigation must target existing section IDs.
- Dynamic subscription plan status should use `role="status"` or an equivalent live region.
- Form labels must remain visible.
- Rendered API values must be escaped with `escapeHtml` or `escapeAttr`.

## Docker Runtime

- Landing runs through the `landing` compose service with `container_name: finance_hub_landing`.
- The container is an nginx reverse proxy to `api:18080`; it should not run a second Go app instance.
- `make restart` should start `api`, `admin`, and `landing`.
- Default public URL: `http://localhost:18182`.
