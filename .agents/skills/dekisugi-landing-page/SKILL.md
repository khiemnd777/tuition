---
name: dekisugi-landing-page
description: Landing Page workflow for the DEKISUGI public pre-auth experience. Use when changing the public landing screen, header, hero/welcome, subscriptions section, Login/Enrollment navigation, footer, hover transforms, public subscription plan API wiring, responsive landing layout, or finance_hub_landing Docker surface.
---

# DEKISUGI Landing Page

Use this for changes to the public pre-auth Landing Page and its Docker surface.

## Ownership

- `web/index.html`: public landing structure plus existing Login, Enrollment, Bootstrap, and Password Reset forms.
- `web/styles.css`: landing layout, visual system, hover transforms, responsive behavior.
- `web/app.js`: landing navigation, public subscription plan loading, enrollment plan select, auth form mode switching.
- `subscription.go` and `rbac.go`: read-only public subscription plan endpoint when the landing page needs real plan data.
- `docker-compose.yml`, `docker/landing.nginx.conf`, and `Makefile`: `finance_hub_landing` container and local Docker access.
- `docs/landing/DESIGN.md`: source of truth for landing design intent.

## Design Rules

- The first screen is the usable public landing experience, not a marketing-only splash.
- Keep the brand/product visible in the first viewport: `DEKISUGI Finance Hub`.
- Include clear areas: Header, Welcome/Greeting, operational areas, Subscriptions, Login/Enrollment access, Footer.
- Use Material Symbols already loaded by the app; avoid custom inline SVG icons unless necessary.
- Use hover transforms to highlight sections/cards, but keep text stable and avoid layout shift.
- Keep cards at 8px radius or less, matching the repo UI rules.
- Keep the palette balanced. Do not let the landing become one-note green, beige, purple, or dark slate.
- Text must fit on mobile and desktop without overlap. Use responsive grids and wrapping buttons.

## Data Rules

- Do not mock subscription plans. Render plans from Platform Admin data through the real API.
- Public landing may read plans through `GET /api/v1/public/subscription-plans`.
- Keep `GET /api/v1/subscriptions/plans` protected by `subscription.view`.
- If plan data is unavailable, show empty/error state; do not hard-code Standard or Trial as live data.
- Enrollment submits intake to `/api/v1/intake`. Platform Admin creates tenants/subscriptions later.

## Implementation Workflow

1. Read `docs/landing/DESIGN.md` and the current landing code path in `web/index.html`, `web/styles.css`, and `web/app.js`.
2. Preserve existing form IDs and auth flows unless the user explicitly asks for a contract change.
3. Escape rendered API values with `escapeHtml`/`escapeAttr`.
4. Keep subscription public API read-only and route-classified as public in `rbac.go`.
5. Update tests when changing RBAC or API behavior.
6. Update README/Makefile when Docker runtime behavior changes.

## Verification

Run focused checks based on the change:

```sh
go test ./...
node --check web/app.js
docker compose config
```

For UI/browser verification, follow `AGENTS.md`: only run browser verification after the user explicitly approves that validation step.
