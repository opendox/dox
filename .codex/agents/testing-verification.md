<!--
  dox
  Copyright (C) 2026  OpenDox

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program. If not, see <http://www.gnu.org/licenses/>.

  @File    : .codex/agents/testing-verification.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 14. Testing And Verification Rules

Any formal Dox change must consider testing and verification.

Testing and verification are not ceremonial. They prove that the change matches the issue scope, does not break core paths, and has clearly described risk.

## General Principles

Agents must explain verification results when completing changes.

Verification may include:

- Unit tests.
- Integration tests.
- Type checking.
- Build.
- Database migration checks.
- Playwright browser verification.
- Figma comparison.
- Manual checks.
- Summaries of relevant logs or command output.

If a verification step cannot run, explain why. Do not pretend tests were run.

## Backend Verification

Go backend changes should run tests according to change scope.

Prefer the smallest meaningful verification:

- When changing one package, run that package's tests first.
- When changing shared infrastructure, run related package tests.
- When changing cross-module behavior, run broader tests.
- When changing public interfaces, configuration, database, scheduling, collection, or computation flows, consider `go test ./...`.

Backend verification may include:

- `go test ./...`
- Targeted package `go test`
- Compilation checks
- CLI command checks
- Configuration loading checks
- Database migration checks
- Log field checks
- Error path checks

Tests involving PostgreSQL, Redis, message queues, or external APIs must state environment requirements and isolation strategy.

Tests must not accidentally connect to production, real sensitive environments, or real third-party platforms.

## Frontend Verification

Frontend changes should run verification according to change scope.

Frontend verification may include:

- Type checking.
- Build.
- Page startup.
- Playwright browser verification.
- Desktop viewport screenshots.
- Mobile viewport screenshots.
- Console error checks.
- Key interaction checks.
- Text overflow and UI overlap checks.

Recommended commands include:

- `pnpm type-check`
- `pnpm build`
- `pnpm dev`

If a change affects UI, layout, interaction, or responsive behavior, use Playwright or screenshots where practical.

## Playwright Verification

Playwright is for real browser verification. It is not decoration.

Use Playwright for:

- New pages.
- Important components.
- Login pages.
- Plugin configuration pages.
- Data tables.
- Report pages.
- Dialogs.
- Mobile layout.
- Figma fidelity.
- Complex interactions.

Playwright verification should check where practical:

- Whether the page opens.
- Whether console errors exist.
- Whether key text exists.
- Whether key buttons are clickable.
- Whether forms accept input.
- Whether UI elements overlap.
- Whether text overflows.
- Whether desktop and mobile viewports work.

## Database And Migration Verification

Database-related changes require extra care.

Changes involving schema, migrations, data layers, indexes, partitioning, task tables, plugin tables, or business fact tables should explain:

- Which tables changed.
- Which fields changed.
- Whether existing data is affected.
- Whether migration is required.
- Whether rollback risk exists.
- Whether Web, Scheduling, Collection, or Computation systems are affected.

Migration files must be reviewable and must not depend on implicit database state.

## Scheduling, Collection, And Computation Verification

Changes involving Scheduling, Collection, or Computation should verify where practical:

- Whether tasks can be created.
- Whether task state transitions are correct.
- Whether tasks are retryable.
- Whether failure paths are recorded.
- Whether idempotency logic holds.
- Whether events include required fields.
- Whether logs include context such as `task_id`, `plugin_id`, and `correlation_id`.
- Whether the Web System can view task state.
- Whether the queue or task model can be operated.

If real queues or real external platforms are unavailable, use mocks, fakes, dry runs, or local verification and explain the behavior.

## Figma Verification

When implementing from Figma, explain:

- Which Figma file or node was used.
- Which design information was extracted.
- Which parts were adapted for engineering reasons.
- Whether screenshots were verified.
- Whether any intentional deviations from Figma exist.

Do not treat Figma-generated React/Tailwind code as final implementation.

## PR Verification Notes

PRs must include verification notes.

Recommended format:

```md
## Verification

- [x] `go test ./...`
- [x] `pnpm type-check`
- [x] `pnpm build`
- [x] Playwright checked `/auth/login` at desktop and mobile sizes
- [ ] Database migration tested against local PostgreSQL

Notes:
- PostgreSQL integration tests were not run because local database credentials are not configured.
```

Verification that was not run must be stated clearly. Do not omit it.
