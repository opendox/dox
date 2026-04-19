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

  @File    : CONTRIBUTING.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# Contributing

Thank you for your interest in Dox. This project is building a plugin-oriented Amazon product performance intelligence platform, so contributions work best when the scope is clear and the system boundary is respected.

## Before You Start

Please open or comment on an issue before starting significant work. A good issue should explain:

- the problem or opportunity;
- the expected outcome;
- affected systems such as Web, Scheduling, Collection, Computation, plugins, frontend, or data;
- what is out of scope;
- acceptance criteria;
- risks, migration impact, or security concerns.

Small typo fixes may go directly to a pull request, but features, fixes, refactors, documentation changes, CI changes, and architecture work should be issue-first.

## Development Scope

Dox is organized around four systems:

- Web controls configuration, identity, authorization, plugin activation, routes, menus, CRUD workflows, and analysis presentation.
- Scheduling turns configuration and policy into tasks, queues, retries, priorities, and dispatch.
- Collection integrates APIs, crawlers, and external platforms to capture source data.
- Computation reads PostgreSQL data layers and builds application-facing metrics, reports, and alert inputs.

Please keep changes aligned with these boundaries. Avoid mixing unrelated refactors, formatting, feature additions, and dependency upgrades in the same pull request.

## Branches And Commits

Use a branch name that points back to the issue when possible:

```text
codex/gh-<issue-number>-<short-name>
```

Commit messages should be clear and typed:

```text
gh-<issue-number> <type>: <summary>
```

Common types are `feat`, `fix`, `docs`, `refactor`, `test`, `style`, `perf`, `build`, `ci`, and `chore`.

Formal commits should include a body, not only a subject line. The body should explain the purpose, list the important changes, include the date, link the issue, and include a sign-off:

```text
gh-<issue-number> <type>: <summary>

Describe why this change exists.

- List the important changes.
- Keep each item reviewable.

Date: YYYY-MM-DD
Refs: #<issue-number>
Signed-off-by: Frost Leo <frostleo.dev@gmail.com>
```

Use `Closes: #<issue-number>` when the commit or pull request completes the issue. Use `Refs: #<issue-number>` when the change is related but does not close it. Commits should be GPG or SSH signed when a signing key is available.

## Pull Requests

A pull request should be small enough to review carefully. Please include:

- what changed;
- the related issue or discussion;
- implementation notes;
- verification steps;
- known risks or follow-up work.

Backend PRs should mention tests, database impact, configuration impact, and migration impact when relevant. Frontend PRs should include screenshots or Playwright verification notes when visual behavior changes.

## Security And Secrets

Never include real credentials, tokens, cookies, private keys, access keys, platform responses containing secrets, or sensitive customer/business data in issues, pull requests, commits, logs, screenshots, or examples.

If you find a security issue, do not open a public issue. Follow [SECURITY.md](SECURITY.md).

## Documentation

Update documentation when a change affects architecture, configuration, APIs, database schema, task or queue behavior, plugin behavior, frontend contracts, security behavior, or development workflow.

Documentation should describe what is true today. Future directions are useful, but they must be clearly marked as future directions.
