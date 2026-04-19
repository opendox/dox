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

  @File    : .codex/agents/frontend-engineering.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 12. Frontend Engineering And Design Rules

Dox frontend is the interaction shell of the management and analytics platform, and it is also the hosting layer for plugin capabilities.

The frontend is not a generic admin template or a simple collection of CRUD pages. It must support Dox management, analytics, configuration, task observation, data insight, and plugin-oriented business capabilities.

## Technology Boundaries

Current frontend technology direction includes:

- Vue 3.
- Vite.
- TypeScript.
- Pinia.
- Vue Router.
- Ant Design Vue.
- Figma as a design reference source.
- Playwright as a browser verification tool.

Do not replace the core frontend stack without an issue and a clear reason.

Do not casually introduce new UI frameworks, state management libraries, routing systems, CSS frameworks, or large visualization libraries.

## File Splitting And Module Boundaries

Frontend code must keep clear file splitting.

Do not put pages, requests, state, forms, tables, charts, business rules, and complex interactions all into one `.vue` file.

Recommended splitting:

- `views/`: page entry points and page-level composition.
- `components/`: page-local business components.
- `components/common/` or the project-defined directory: cross-page reusable components.
- `composables/`: reusable state, interactions, and business flows.
- `api/`: request wrappers.
- `stores/`: cross-page state.
- `types/`: type definitions.
- `utils/`: pure function utilities.
- `constants/`: constants and enums.
- `adapters/`: backend data to frontend view-model transformations.
- `plugins/`: frontend plugins or extension capabilities.

Page components should organize layout and flow. They should not contain too many business implementation details.

## Vue Rules

Vue code should prefer:

- Composition API.
- `<script setup lang="ts">`.
- Explicit props, emits, and types.
- Reusable composables.
- Clear component boundaries.

Avoid:

- Options API.
- Implicit `any`.
- Oversized single-file components.
- Putting business requests, state management, and complex presentation logic into one component.
- Over-abstracted generic components.
- Meaningless wrapper components.

## Reusable Component System

Dox frontend should build reusable components instead of reimplementing tables, charts, filters, and status displays on every page.

Prefer building component families such as:

- Data tables.
- Chart containers.
- Metric cards.
- Filters.
- Time range pickers.
- Task lists.
- Task details.
- Queue status views.
- Plugin configuration forms.
- Plugin status cards.
- Alert lists.
- Report layouts.
- Empty states.
- Error states.
- Loading states.
- Permission-aware rendering.
- Detail drawers.
- Confirmation dialogs.

Reusable components should carry variation through props, slots, configuration objects, or composables. Avoid copy-pasting page code.

Do not over-abstract only for the sake of reuse. Abstract only when a component reduces real duplication, expresses a stable business pattern, or improves consistency.

## Frontend Plugin Shell

The frontend is a shell for plugin capabilities, not the decision-maker for whether a plugin is available.

Menus, routes, permissions, plugin entries, and some page capabilities should be issued by the Web System according to user, permission, role, and plugin enablement state.

The frontend must not hard-code that a platform is always visible.

The frontend should provide shared capabilities:

- App layout.
- Navigation.
- Plugin entry hosting.
- Plugin configuration page containers.
- Task list and task detail views.
- Report and dashboard hosting.
- Permission-aware rendering.
- Loading, empty, error, and degraded states.

## Frontend HTTP And Interceptor Plugins

Frontend HTTP clients, request interceptors, response interceptors, and cross-cutting request capabilities should also be pluginized.

Request-flow capabilities must not become unmaintainable hard-coded `if/else` branches.

Frontend HTTP plugins may include:

- Authentication token injection.
- CSRF handling.
- Signing.
- Encryption.
- Request deduplication.
- Request cache.
- Retry.
- Rate limiting.
- Error normalization.
- Logging.
- Metrics collection.
- Mocking.
- Reserved tenant or workspace headers.
- Plugin platform-related header injection.

Each interceptor plugin should have a single responsibility, explicit ordering, enable/disable controls, and tests where practical.

Do not put business page logic into global request interceptors.

## Figma Usage

Figma is the source for design language and visual references.

When a task involves UI, components, pages, design systems, or visual fidelity, agents should ask whether a Figma file, node, or reference design exists.

If the user provides a Figma file or node, read the Figma context before implementation.

Code returned from Figma is usually reference code only. Do not copy React/Tailwind output directly into the Vue project.

Extract from Figma:

- Colors.
- Typography.
- Spacing.
- Radius.
- Shadows.
- Blur.
- Component structure.
- Interaction states.
- Responsive intent.

Then implement using Dox Vue, CSS, and component patterns.

## Apple-inspired Design Language

Dox frontend should use an Apple-inspired professional analytics design language.

This means:

- Restrained.
- Precise.
- Clear.
- Premium but not flashy.
- Rhythmic spacing.
- Clear typography hierarchy.
- Soft and meaningful motion.
- Native-app-like component feel.
- Professional data-tool information density for analytics pages.

Do not copy Apple trademarks, official pages, proprietary icons, official copy, or protected assets.

Dox should learn from Apple and macOS design language:

- Clear information hierarchy.
- Precise spacing.
- Soft depth.
- Pill buttons.
- Lightweight glass feel.
- Precise shadows.
- Stable radius system.
- Natural interaction feedback.

The final goal is a professional analytics platform, not an Apple website clone.

## Analytics Interface Rules

Dox is a professional analytics platform. Data pages must not only be spacious and pretty.

Analytics interfaces should balance:

- High information density.
- Clear comparison.
- Fast filtering.
- Stable tables.
- Trend charts.
- Abnormality hints.
- Status markers.
- Traceable sources.
- Explainable metrics.
- Task and data time windows.

Avoid analytics pages that have:

- Only large cards and no details.
- Pretty charts without metric definitions.
- Metrics without units.
- Metrics without time windows.
- Metrics without data sources.
- Status colors without semantics.
- Tables that cannot be scanned.

## Playwright Verification

Important UI changes should use Playwright or browser screenshots where practical.

Verification should cover:

- Desktop viewport.
- Mobile viewport.
- Console errors.
- Whether the page loads.
- Whether key interactions work.
- Whether text overflows.
- Whether UI elements overlap.
- Whether important states are visible.

If Playwright or browser verification cannot run, explain why in the PR.
