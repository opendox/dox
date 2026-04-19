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

  @File    : .codex/agents/documentation.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# Documentation Rules

This document defines how Dox project documentation should be written and maintained.

## Language

Project documentation may be written in Chinese when it is meant for the maintainer, internal contributors, business discussion, architecture alignment, or implementation planning.

Use English when the document is primarily consumed by code-generation agents, public open-source users, package metadata, exported API contracts, or ecosystem tooling that expects English.

Mixed-language documentation is acceptable when it improves clarity:

- Keep code identifiers, table names, API paths, event names, queue names, config keys, and error codes in their original English form.
- Do not translate proper technical identifiers into vague Chinese phrases.
- Prefer Chinese explanations for business context, operational reasoning, domain background, and decision records when those documents are for the project team.
- Prefer English for templates or rule files that must be followed by many different coding agents.

## Purpose

Documentation exists to preserve long-term engineering context. It should help future contributors understand:

- what the system does;
- why an architecture or business decision was made;
- which boundaries must not be crossed;
- how to run, test, debug, deploy, and operate the system;
- how a feature should be used and extended.

Do not treat documentation as marketing copy. It is part of the engineering system.

## Documentation Types

Use the appropriate location for the type of knowledge being recorded:

- `README.md`: project entry point, positioning, quick start, and important links.
- `AGENTS.md`: root agent rule index and progressive loading entry point.
- `.codex/agents/`: detailed coding-agent rules. This is not community-facing project documentation.
- `docs/assets/`: public static assets for README files, project documentation, and GitHub Pages.
- `docs/en-us/`: English project documentation.
- `docs/zh-cn/`: Simplified Chinese project documentation.
- `docs/{language}/architecture/`: system architecture, four-system boundaries, plugin model, event queue design, data layering, and major decisions.
- `docs/{language}/domain/`: business domain concepts, Amazon data sources, keyword/ASIN flywheel, supplier direction, and domain glossary.
- `docs/{language}/development/`: local setup, dependency installation, tests, builds, debugging, and developer workflow.
- `docs/{language}/api/`: API contracts, authorization, request/response examples, error codes, pagination, sorting, and filtering.
- `docs/{language}/database/`: data model, migrations, PostgreSQL layering, table design, indexes, and data lifecycle.
- `docs/{language}/frontend/`: design language, Figma sources, component structure, reusable component contracts, page structure, and frontend plugin conventions.
- `docs/{language}/operations/`: deployment, configuration, queue operations, logs, alerts, task operations, and operational playbooks.
- `.github/`: GitHub-only configuration such as issue templates, pull request templates, and workflows. Do not use `.github/` as the primary location for README or GitHub Pages assets.
- `CHANGELOG.md`: meaningful behavior, API, migration, configuration, security, and plugin protocol changes.

Do not create these directories just because they are listed here. Create them only when the maintainer asks for that documentation, workflow, asset, or project structure.

## When To Update Documentation

Update documentation in the same change when code changes affect:

- user-visible behavior;
- architecture boundaries;
- service responsibilities;
- plugin protocols;
- configuration;
- API contracts;
- database schema or data lifecycle;
- task, queue, scheduling, or event flows;
- frontend interaction or reusable component contracts;
- alerting, notification, permission, or audit behavior;
- local development, testing, build, deployment, or operations workflow.

If an issue requires documentation, the pull request must include the corresponding documentation update. If no documentation update is needed, the pull request may state why.

## Architecture Documents

Architecture documents should explain decisions, not only describe the final shape.

Prefer this structure when writing substantial architecture notes:

- Background
- Goals
- Non-goals
- Current decision
- Alternatives considered
- Why this decision was chosen
- Impact on Web, Scheduling, Collection, and Computation
- Operational risks
- Future directions, clearly marked as future directions

Do not describe future ideas as already implemented.

## API Documents

API documentation should include:

- path and method;
- authentication and permission requirements;
- request parameters and body;
- response shape;
- error codes;
- pagination, sorting, and filtering behavior;
- related task or event behavior;
- security notes when credentials, tenant data, or external platform data are involved.

Keep examples realistic and aligned with actual code.

## Database Documents

Database and data-layer documentation should include:

- table or model purpose;
- data layer such as `raw`, `ods`, `dwd`, `dws`, or `ads`;
- write owner;
- primary readers;
- important indexes;
- idempotency and deduplication rules;
- lifecycle or retention expectations;
- upstream and downstream relationships.

Do not document PostgreSQL as a temporary placeholder. PostgreSQL is the primary database for Dox unless the user explicitly changes that direction.

## Frontend Documents

Frontend documentation should make reusable UI work easier over time.

Important frontend documentation should cover:

- Figma source links and node references when a design source exists;
- design language decisions;
- reusable component contracts;
- table, chart, form, navigation, and shell component patterns;
- page structure;
- frontend plugin hosting rules;
- frontend request plugin and interceptor conventions;
- Playwright verification expectations.

When writing reusable component documentation, include props, events, slots, examples, and expected states when relevant.

## Changelog

Record meaningful changes, not every internal implementation detail.

Changelog entries should cover:

- new features;
- behavior changes;
- breaking changes;
- important fixes;
- migrations;
- configuration changes;
- security changes;
- API changes;
- plugin protocol changes.

## Writing Style

Write clearly and directly. Prefer accurate engineering language over decorative wording.

Avoid:

- vague descriptions that hide constraints;
- promises that do not match the code;
- describing unimplemented future directions as current behavior;
- duplicating the same explanation across many files;
- writing docs that only repeat function names without explaining intent or boundaries.

Issues and pull requests are also project documentation. Keep important context in durable issues, pull requests, or docs instead of leaving it only in chat history.
