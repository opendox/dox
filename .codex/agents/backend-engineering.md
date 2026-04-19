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

  @File    : .codex/agents/backend-engineering.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 11. Backend Engineering Rules


Dox backend work must prioritize clarity, maintainability, testability, and long-term evolution.

Backend implementation should follow idiomatic Go. Do not introduce unnecessary abstractions only to make the architecture look complex. Architecture must serve Dox system boundaries, plugin mechanisms, scheduling flow, collection flow, computation flow, and PostgreSQL-first data layering.

### Technology Boundaries

Current backend technology direction includes:

- Go.
- Fiber as the HTTP framework.
- Ent as the schema and data access foundation.
- PostgreSQL as the primary storage system.
- Redis for cache, locks, rate limits, short-lived state, or session-like capabilities.
- Cobra for CLI and executable commands.
- Zap or the project-defined logging system for structured logs.
- Viper or the project-defined configuration system for configuration loading.

Do not replace core technology choices without an issue and a clear reason.

Do not casually introduce large frameworks, ORMs, message queue clients, configuration systems, or logging systems.

### Module Boundaries

Backend code should keep clear module boundaries.

Recommended boundaries include:

- `cmd`: CLI commands and process entry points.
- `app`: business modules.
- `internal`: internal infrastructure, system capabilities, and non-public implementation.
- `pkg`: public packages only when they truly need external reuse.
- `configs`: configuration files and templates.
- `docs`: backend-related documentation.
- `internal/ent/schema`: Ent schema definitions.
- generated `internal/ent` code: must not be edited manually.

Business capability should live in explicit modules. Do not scatter business logic into `cmd`, middleware, or infrastructure packages.

### Four-system Backend Boundaries

Backend implementation must respect Web, Scheduling, Collection, and Computation responsibilities.

Web-related code must not directly execute collection work or heavy computation.

Scheduling-related code should manage task orchestration, task state, retries, priority, and scheduling policy.

Collection-related code should execute collection tasks, call platform APIs, run crawlers, or integrate external data sources.

Computation-related code should handle batch processing, aggregation, report generation, and upper-layer data generation for Web use.

Even if these systems share code or processes early, their responsibilities must not be mixed.

### Ent And Database

Ent schemas are one source of truth for the data model.

When changing database models, modify schemas under `internal/ent/schema` first, then update generated Ent code through the generation workflow.

Do not manually edit generated Ent code.

Database design must consider:

- Primary keys and unique keys.
- Idempotent writes.
- Deduplication.
- Indexes.
- Future partitioning.
- Data source.
- Task ID.
- Plugin ID.
- Platform ID.
- Time fields.
- Data layer.
- Audit fields.
- Future migration and extension.

Migrations must be managed through explicit migration workflows. Do not depend on implicit database changes.

### PostgreSQL Query Principles

PostgreSQL queries should be explicit, explainable, and optimizable.

Avoid:

- Unbounded full-table scans.
- Long-running aggregation in Web request handlers.
- N+1 queries.
- Large-range filters without indexes.
- Implicit loading of large relationship graphs.
- Complex business computation inside request handlers.

Complex statistics, aggregation, and batch processing belong in the Computation System.

The Web System may query lower-layer data, but it must have clear filters, pagination, indexes, and authorization boundaries.

### Redis Usage Principles

Redis may be used for:

- Cache.
- Distributed locks.
- Rate limits.
- Short-lived task state.
- Short-lived session state.
- Deduplication windows.
- Temporary counters.

Redis must not become an untraceable long-term source of business facts.

Important business facts must eventually be stored in PostgreSQL or another explicit durable storage system.

Redis usage must consider:

- Key naming.
- TTL.
- Idempotency.
- Lock release.
- Failure recovery.
- Cache breakdown.
- Cache consistency.
- Multi-instance concurrency.

### Redis Key Naming

Redis keys must use clear, stable, and consistent naming.

Recommended format:

- `dox:<app>:<module>:<feature>:<identifier>`

Examples:

- `dox:iam:auth:session:<session_id>`
- `dox:iam:auth:mfa:<challenge_id>`
- `dox:scheduling:task:lock:<task_id>`
- `dox:collection:amazon:rate_limit:<credential_id>`
- `dox:plugin:feishu:token:<credential_id>`

Redis keys should make it clear:

- That the key belongs to Dox.
- Which app or system owns it.
- Which major module owns it.
- Which feature owns it.
- Which business object or identifier it belongs to.

Avoid vague keys such as:

- `token`
- `cache`
- `lock`
- `user:1`
- `tmp`

Redis key design must consider:

- TTL.
- Naming collisions.
- Environment isolation.
- Future multi-tenant extension.
- Deletion strategy.
- Debuggability.
- Bulk scan risk.
- Whether sensitive information is included.

Redis keys must not contain raw secrets, tokens, passwords, complete credentials, or other sensitive values. Use internal IDs, hashes, or safe references when referencing sensitive objects.

### Context And Cancellation

All operations that may involve databases, Redis, external APIs, queues, files, networks, or long-running execution should use `context.Context`.

Do not lose cancellation signals in background goroutines.

Long-running tasks should support timeout, cancellation, and status reporting.

### Error Handling

Error handling must be explicit.

Do not swallow errors.

Do not use `panic` for normal business errors.

Errors should preserve enough context, such as task, plugin, platform, operation, input identifier, and downstream error information.

Errors exposed to users and errors written to internal logs should be handled differently. Do not write sensitive credentials, tokens, secrets, or raw credentials into logs or frontend responses.

### Logging And Observability

Backend logs must be structured.

Logs should support observation by type, system, app, and module so they can later integrate with ELK or other observability platforms.

Logs should include where practical:

- `log_type`: such as `request`, `task`, `collection`, `computation`, `auth`, `audit`, `security`, `plugin`, `queue`, or `system`.
- `system`: such as `web`, `scheduling`, `collection`, or `computation`.
- `app`: such as `iam`, `plugin`, `amazon`, or `notification`.
- `module`: major module.
- `feature`: feature or capability.
- `request_id` or `correlation_id`.
- `task_id`.
- `plugin_id`.
- `platform_id`.
- `worker_id`.
- `user_id` or operator.
- `operation`.
- `duration`.
- `result`.
- `error`.

The backend should support splitting logs by app or module. For example, IAM logs may be written to an IAM log file, scheduling logs to a scheduling log file, and collection logs to a collection log file.

The goal of log splitting is better debugging and better later ELK integration, not excessive fragmentation. Log splitting strategy should stay consistent and configurable.

Security, audit, authentication, task, and collection error logs should have stable fields for search, aggregation, and alerting.

Logs must not leak sensitive information.

Collection, Scheduling, and Computation flows must preserve enough context for debugging.

### Testing

Backend changes should add tests according to risk.

Prioritize tests for:

- Business rules.
- Data transformations.
- Plugin definitions.
- Scheduling state machines.
- Idempotency logic.
- Error handling.
- Configuration loading.
- Permission checks.
- Data-layer transformations.

Logic that can be covered by unit tests should not rely only on manual verification.

Database tests must use explicit test databases or isolation strategies. Tests must not accidentally connect to production or real sensitive environments.

### Generated Code And Tooling Code

Generated code must not be edited manually.

If generated code has a problem, fix the schema, template, generator, or generation workflow.

Tooling code, scripts, and one-off migrations must clearly state their purpose. They must not be mixed into the main business path without a clear reason.

