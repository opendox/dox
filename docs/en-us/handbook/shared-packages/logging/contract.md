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

  @File    : docs/en-us/handbook/shared-packages/logging/contract.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Logging Contract

The shared logging contract defines what callers can rely on from `packages/shared/logging`, which behavior is only configuration shape today, and which errors are part of the package boundary.

## Ownership Contract

| Area | Package Owns | Runtime Owns |
| --- | --- | --- |
| Logging vocabulary | Resource, correlation, event, node, tags, and fields model. | Runtime-specific naming policy and business event taxonomy. |
| Business logging API | `Logger`, `Attr`, and typed attribute constructors. | Passing the facade to services and middleware. |
| zap mapping | Level, encoder, config, console core, JSON/file core, lumberjack helpers. | Choosing output paths, lifecycle, deployment topology, and sink ownership. |
| OpenTelemetry SDK base | Resource mapping, propagator mapping, SDK provider construction, flush/shutdown helpers. | Installing globals, exporter wiring, and lifecycle integration. |
| Config validation | Defaults and package validation errors. | Runtime-specific configuration restrictions and environment-specific policy. |

> [!NOTE]
> The package is a shared foundation. It deliberately does not start logging for any process on its own.

## Config Contract

`Config` is the root shared logging configuration shape.

| Group | Purpose |
| --- | --- |
| `Level` | Root Dox logging level. |
| `Development` | Root development-mode switch that feeds zap defaults. |
| `Resource` | Logging-specific resource overrides, currently `service_version`. |
| `Zap` | zap-facing config shape and encoder symbols. |
| `Cores` | Console/file core declarations. |
| `Buffering` | Future buffered writer behavior. |
| `Shutdown` | Flush/shutdown timeout. |
| `Redaction` | Sensitive key replacement policy shape. |
| `OTel` | OpenTelemetry propagation, provider, exporter, and batch settings. |

`Config.Default` fills defaults before validation. `Config.Validate` checks supported enum values, required paths, duplicate core names, positive durations and sizes, redaction key shape, OTLP protocol shape, and batch constraints.

## Default Contract

| Field | Default |
| --- | --- |
| Root level | `info` |
| zap encoding | `json` |
| zap output paths | `stdout` |
| zap error output paths | `stderr` |
| encoder message key | `message` |
| encoder level key | `severity_text` |
| encoder time key | `timestamp` |
| buffering | enabled, `262144` bytes, `1s` flush interval |
| shutdown timeout | `5s` |
| redaction | enabled, `[REDACTED]`, default sensitive keys |
| OpenTelemetry | root/traces/metrics/logs enabled, trace context and baggage enabled |
| OTLP exporter | disabled by default |

Explicit disabled pointer toggles remain disabled after defaults are applied.

## Validation Error Contract

The package exposes:

- `FieldError`, with `Field` and `Reason`;
- `ValidationError`, with `Fields []FieldError`;
- compact error strings for logs.

Callers should inspect `ValidationError.Fields` instead of parsing error strings.

<details>
<summary>Example validation fields</summary>

```text
level: level is not supported
cores[1].name: core name must be unique
redaction.keys[0]: redaction key must not be empty
otel.traces.sampler.ratio: trace sampler ratio must be between 0 and 1
```

</details>

## Implemented Versus Contract-Only Settings

| Setting | Validation | Runtime Effect Today |
| --- | --- | --- |
| `cores[].datasets` | Validated as non-empty entries. | Not used for event dataset routing yet. |
| `buffering.*` | Validated when enabled. | No buffered writer is installed yet. |
| `redaction.*` | Validated when enabled. | No log field or value redaction is applied yet. |
| `cores[].rotation.driver=external/logrotate` | Valid enum value. | Rejected by the current zap file sink. |
| `otel.exporter.otlp.*` | Shape is validated. | Enabled OTLP exporter is rejected by `NewOpenTelemetrySDKBase`. |
| `zap.sampling.hook_metrics` | Part of config shape. | No sampling hook metrics are installed yet. |
| `DefaultFilePathTemplate` | Used as default string. | Template variables are not rendered. |

> [!WARNING]
> Treat contract-only settings as forward-compatible schema, not as active runtime behavior.

## Logger Contract

Business-facing code uses `Logger` and `Attr`.

The facade guarantees:

- log level methods: `Debug`, `Info`, `Warn`, `Error`, `DPanic`, `Panic`, `Fatal`;
- logger derivation: `Named`, `With`;
- flush path: `Sync`;
- Dox-owned attributes for resource, correlation, event, node, tags, fields, single fields, and errors;
- context correlation merge on every write.

The facade deliberately does not expose `zap.Logger` or `zap.Field` in method signatures.

## Context Correlation Contract

Correlation helpers provide simple overlay semantics:

- `ContextWithCorrelation` stores a correlation value and accepts nil context by using `context.Background`.
- `ContextWithMergedCorrelation` merges non-empty overlay fields into existing context correlation.
- `CorrelationFromContext` returns `(Correlation, false)` for nil or missing context.
- `MergeCorrelation` keeps base fields unless overlay fields are non-empty.

## Out of Contract

The following behavior is outside the package contract today:

- opening runtime loggers during server/scheduler/collector/compute startup;
- HTTP middleware correlation;
- task/job/plugin correlation injection;
- OpenTelemetry global provider installation;
- OTLP exporter construction;
- collector deployment examples;
- runtime-specific sink path rendering;
- multi-process rotation coordination;
- field-level redaction execution;
- dataset-based core routing.

## Related Pages

- [Model](model.md)
- [Runtime Boundary](runtime-boundary.md)
- [Functions and API](functions.md)
