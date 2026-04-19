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

  @File    : AGENTS.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# Dox Agent Rules

Dox uses progressive agent rules. Keep this file as the repository-level agent entry point. Detailed rules live under `.codex/agents/` and should be loaded only when relevant.

There is no required `.codex/AGENTS.md` mirror. The root `AGENTS.md` is the stable entry point; `.codex/agents/` holds the detailed rule modules.

## Always Read

Before doing any Dox work, agents must read:

- `.codex/agents/project-identity.md`
- `.codex/agents/collaboration-standard.md`
- `.codex/agents/file-headers.md`

These documents define what Dox is, how to collaborate with the maintainer, and the mandatory copyright/source header requirements for new source files.

## Read When Relevant

Read task-specific rules before editing related areas:

- Core system boundaries: `.codex/agents/core-system-architecture.md`
- PostgreSQL schemas, queries, data layering, or computation inputs/outputs: `.codex/agents/postgresql-data-layering.md`
- Plugin work: `.codex/agents/plugin-architecture.md`
- Scheduling, queues, tasks, events, workers, or idempotency: `.codex/agents/event-queue-architecture.md`
- Alerts, notifications, notification clients, or notification plugins: `.codex/agents/alerting-notification.md`
- Business domain modeling, Amazon data, keywords, ASINs, market data, suppliers, or AI business behavior: `.codex/agents/business-domain.md`
- GitHub issues, branches, commits, pull requests, labels, templates, or changelog work: `.codex/agents/github-workflow.md`
- Go backend, Fiber, Ent, PostgreSQL, Redis, logging, configuration, migration, or backend tests: `.codex/agents/backend-engineering.md`
- Frontend UI, Vue, Vite, Pinia, Figma, Playwright, frontend plugins, reusable components, or frontend HTTP work: `.codex/agents/frontend-engineering.md`
- Security, IAM, credentials, secrets, authorization, audit logs, external platform credentials, or AI safety work: `.codex/agents/security-privacy.md`
- Tests, CI, builds, Playwright verification, database migration checks, or PR verification notes: `.codex/agents/testing-verification.md`
- README files, architecture docs, API docs, database docs, frontend docs, operations docs, changelog entries, or other project documentation: `.codex/agents/documentation.md`

If a task touches multiple areas, read every relevant document before editing.

## Project Documentation

Project documentation belongs under `docs/` when project documentation is explicitly requested. It is community-facing and must not be mixed with coding-agent process rules.

Documentation may be bilingual when that helps both Chinese-speaking maintainers and international contributors. Prefer Chinese for business context and domain reasoning, and English for public summaries, API terms, code identifiers, and ecosystem-facing descriptions.

Do not put coding-agent process rules into community documentation unless the maintainer explicitly asks for it.

## Repository Bootstrap Discipline

Do not create application source trees, README files, assets, GitHub workflows, documentation folders, or implementation scaffolding just because the architecture mentions them.

When the task is about agent rules, limit changes to `AGENTS.md` and `.codex/agents/` unless the maintainer explicitly asks for repository scaffolding or community documentation.

## Rule Precedence

- User instructions for the current task take priority when they are explicit.
- These agent rules define the default project behavior when the user does not override it.
- Do not silently ignore a rule. If a rule conflicts with the task, explain the conflict and use the user's latest instruction as the deciding context.

## Scope Control

Keep changes scoped to the active issue or request. Do not expand implementation into long-term directions unless the user explicitly asks for that work.
