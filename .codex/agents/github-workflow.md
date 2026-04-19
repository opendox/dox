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

  @File    : .codex/agents/github-workflow.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 9. GitHub Workflow


Dox development must follow an issue-first and pull-request-first workflow.

Except for very small local experiments or explicitly requested maintenance operations, features, fixes, refactors, documentation, CI, build changes, and architecture changes should start from a GitHub issue and land through a pull request.

### Issues

Each clear feature, fix, maintenance task, or documentation task should have a corresponding issue.

Issues should describe:

- Background.
- Goal.
- Scope.
- Non-goals.
- Related modules.
- Acceptance criteria.
- Risks or notes.
- Related links.

Issues must not be only short titles. Agents must understand the issue scope and acceptance criteria before implementation.

When asked to create issues, agents should use project issue templates and keep titles and labels consistent.

### Branches

Each implementation task should create a branch from the related issue.

Recommended branch format:

- `codex/gh-<issue-number>-<short-name>`

Examples:

- `codex/gh-12-governance-foundation`
- `codex/gh-18-postgres-data-layering`
- `codex/gh-24-plugin-registry`

Branch scope must match the issue. Do not mix unrelated features, refactors, or formatting into one branch.

### Commits

Commit messages should bind to the issue and use a clear type.

Recommended format:

- `gh-<issue-number> <type>: <summary>`

Allowed types:

- `feat`
- `fix`
- `docs`
- `refactor`
- `test`
- `style`
- `perf`
- `build`
- `ci`
- `chore`

Examples:

- `gh-12 docs: Add governance workflow`
- `gh-18 feat: Add raw data table schema`
- `gh-24 test: Add plugin registry tests`
- `gh-31 ci: Add frontend type-check workflow`

Commit messages must be semantically clear. Avoid vague messages such as:

- `update`
- `fix`
- `add feature`
- `change something`
- `wip`

If a PR contains several meaningful steps, it may contain multiple commits. Each commit should explain why it exists.

Formal commits must include a meaningful body. Do not create commits that contain only a subject line unless the maintainer explicitly asks for a throwaway local checkpoint.

Commit bodies should include:

- A short paragraph describing the purpose of the change.
- A concise bullet list of important changes.
- `Date: YYYY-MM-DD`.
- `Refs: #<issue-number>` or `Closes: #<issue-number>`.
- `Signed-off-by: Frost Leo <frostleo.dev@gmail.com>` when committing as the maintainer.

When a usable GPG or SSH signing key is available, formal commits must be cryptographically signed. Do not silently fall back to unsigned commits. If signing fails, stop and report the signing failure instead of pushing.

Recommended body format:

```text
<purpose sentence>

- <important change>
- <important change>

Date: YYYY-MM-DD
Refs: #<issue-number>
Signed-off-by: Frost Leo <frostleo.dev@gmail.com>
```

### Pull Requests

All formal changes should land through pull requests.

Recommended PR title format:

- `gh-<issue-number> [<type>]: <summary>`

Examples:

- `gh-12 [docs]: Initialize governance workflow`
- `gh-18 [feat]: Add PostgreSQL raw data layer`
- `gh-24 [refactor]: Introduce plugin registry boundary`

PR bodies should include:

- `Description`: what this PR does.
- `Related Links`: linked issues, documentation, or discussions.
- `Implementation Notes`: key implementation choices.
- `Verification`: tests, builds, screenshots, or manual checks performed.
- `Risks`: remaining risks or uncovered areas.

PRs must link related issues, such as:

- `Closes #12`
- `Refs #18`

UI PRs should include screenshots or Playwright verification notes where practical.

Backend PRs should describe test commands, database migration impact, and configuration impact.

Architecture PRs should describe design tradeoffs and non-goals.

### PR Scope Control

PRs must be small and reviewable.

A PR should solve one issue or one tightly related group of problems.

Agents must not mix the following into a PR without an explicit reason:

- Unrelated refactors.
- Unrelated formatting.
- Unrelated dependency upgrades.
- Large unrelated documentation edits.
- Unrequested architecture changes.
- Feature additions unrelated to the issue.

If implementation reveals a new problem, prefer creating a follow-up issue instead of expanding the current PR.

### CHANGELOG

When requested by the user or required by the issue, PRs should update `CHANGELOG.md`.

The changelog should record changes meaningful to users, developers, or system behavior. It should not list every internal detail.

### Labels And Templates

Dox should maintain project-specific labels, issue templates, and a PR template.

When creating issues or PRs, agents should use the correct type, area, and priority labels where possible.

Do not rely on GitHub default labels as the long-term governance system.

