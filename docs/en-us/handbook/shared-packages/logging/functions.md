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

  @File    : docs/en-us/handbook/shared-packages/logging/functions.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Logging Functions and API

The shared logging API surface defines the exported types, constructors, helpers, constants, and caller-facing behavior available from `packages/shared/logging`.

## Constants and Enums

| API | Purpose |
| --- | --- |
| `DefaultFilePathTemplate` | Default JSONL file path string. Currently not rendered by the package. |
| `DefaultRedactionReplacement` | Default replacement string for redaction config. |
| `Level` and level constants | Dox logging levels: debug, info, warn, error, dpanic, panic, fatal. |
| `Encoding` and encoding constants | Supported encodings: console and json. |
| `CoreType` and core constants | Supported core types: console and file. |
| `RotationDriver` and driver constants | Rotation strategy values: lumberjack, external, logrotate, none. |
| `OTLPProtocol` and protocol constants | OTLP protocol values: grpc and http. |
| `TraceSamplerType` and sampler constants | OpenTelemetry sampler shape values. |

Every enum type has an `IsValid` method.

## Config API

| API | Purpose |
| --- | --- |
| `Config.Default() error` | Applies shared logging defaults. |
| `Config.Validate() error` | Validates the shared logging configuration contract. |
| `DefaultRedactionKeys() []string` | Returns the default sensitive key list. |

Most nested config structs have `Default` methods used by `Config.Default`; callers normally call defaults on the root config.

## Model API

| Type | Purpose |
| --- | --- |
| `Resource` | Telemetry producer identity. |
| `Correlation` | Request/task/workflow/plugin correlation identity. |
| `Event` | Observability event classification. |
| `Node` | Service-internal component and operation. |
| `Tags` | Low-cardinality business labels. |
| `Fields` | Higher-cardinality event facts. |

Field name constants in `fields.go` define the stable output keys.

## Logger API

| API | Purpose |
| --- | --- |
| `Logger` | Business-facing logging facade. |
| `Attr` | Opaque attribute application type. |
| `NewLogger(base, attrs...)` | Wraps a `ZapCoreBase` in the Dox logger facade. |
| `ResourceAttr` | Adds resource fields. |
| `CorrelationAttr` | Adds correlation fields. |
| `EventAttr` | Adds event fields. |
| `NodeAttr` | Adds component and operation fields. |
| `TagsAttr` | Adds tags object entries. |
| `FieldsAttr` | Adds fields object entries. |
| `FieldAttr` | Adds one fields object entry. |
| `ErrorAttr` | Adds one error field. |

`NewLogger` rejects a nil `ZapCoreBase`. Logger facade method signatures do not expose zap types.

## Context API

| API | Purpose |
| --- | --- |
| `ContextWithCorrelation(ctx, correlation)` | Stores correlation on context. |
| `ContextWithMergedCorrelation(ctx, correlation)` | Merges non-empty fields into context correlation. |
| `CorrelationFromContext(ctx)` | Retrieves correlation and a boolean. |
| `MergeCorrelation(base, overlay)` | Applies non-empty overlay fields to base. |

Nil context is handled safely by the helper functions.

## Zap API

| API | Purpose |
| --- | --- |
| `ZapCoreBase` | Runtime-owned zap primitives and close lifecycle. |
| `NewZapLevel(level)` | Maps Dox level to zapcore level. |
| `NewZapAtomicLevel(level)` | Creates zap atomic level. |
| `NewZapEncoderConfig(config)` | Maps symbolic encoder config to zapcore encoder config. |
| `NewZapConfig(config)` | Maps Dox config to zap config. |
| `NewZapCoreBase(config)` | Builds enabled zap cores and zap options. |
| `(*ZapCoreBase).Options()` | Returns copied zap options. |
| `(*ZapCoreBase).Close()` | Closes opened sinks once. |

Runtime bootstrap should own `ZapCoreBase` lifecycle and should avoid exposing zap primitives to business code.

## OpenTelemetry API

| API | Purpose |
| --- | --- |
| `OpenTelemetrySDKBase` | Runtime-owned OpenTelemetry SDK primitives. |
| `NewOpenTelemetrySDKBase(config, resource)` | Builds resource, propagator, and enabled providers. |
| `NewOpenTelemetryResource(model, config)` | Maps Dox resource to SDK resource. |
| `NewOpenTelemetryPropagator(config)` | Builds trace context/baggage propagator. |
| `NewOpenTelemetryTraceSampler(config)` | Builds SDK trace sampler. |
| `(*OpenTelemetrySDKBase).ForceFlush(ctx)` | Flushes enabled SDK providers. |
| `(*OpenTelemetrySDKBase).Shutdown(ctx)` | Shuts down enabled SDK providers. |

The SDK base returns provider objects. It does not install globals or exporters.

## Runtime Integration Sketch

```go
cfg := logging.Config{}
if err := cfg.Default(); err != nil {
	return err
}
if err := cfg.Validate(); err != nil {
	return err
}

zapBase, err := logging.NewZapCoreBase(cfg)
if err != nil {
	return err
}
defer zapBase.Close()

logger, err := logging.NewLogger(zapBase, logging.ResourceAttr(resource))
if err != nil {
	return err
}
defer logger.Sync()
```

This sketch omits runtime policy such as path rendering, global OpenTelemetry installation, and shutdown orchestration.

## Related Pages

- [Contract](contract.md)
- [Model](model.md)
- [Runtime Boundary](runtime-boundary.md)
