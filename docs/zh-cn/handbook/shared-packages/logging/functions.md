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

  @File    : docs/zh-cn/handbook/shared-packages/logging/functions.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Logging 函数与 API

Shared logging API surface 定义 `packages/shared/logging` 暴露的 exported types、constructors、helpers、constants 和 caller-facing behavior。

## Constants 和 Enums

| API | 用途 |
| --- | --- |
| `DefaultFilePathTemplate` | 默认 JSONL file path string。当前不会由本包渲染。 |
| `DefaultRedactionReplacement` | Redaction config 的默认 replacement string。 |
| `Level` and level constants | Dox logging levels: debug, info, warn, error, dpanic, panic, fatal。 |
| `Encoding` and encoding constants | 支持的 encodings: console and json。 |
| `CoreType` and core constants | 支持的 core types: console and file。 |
| `RotationDriver` and driver constants | Rotation strategy values: lumberjack, external, logrotate, none。 |
| `OTLPProtocol` and protocol constants | OTLP protocol values: grpc and http。 |
| `TraceSamplerType` and sampler constants | OpenTelemetry sampler shape values。 |

每个 enum type 都有 `IsValid` method。

## Config API

| API | 用途 |
| --- | --- |
| `Config.Default() error` | 应用 shared logging defaults。 |
| `Config.Validate() error` | 校验 shared logging configuration contract。 |
| `DefaultRedactionKeys() []string` | 返回默认 sensitive key list。 |

大多数 nested config structs 都有 `Default` methods，由 `Config.Default` 使用；调用方通常只需要调用 root config defaults。

## Model API

| Type | 用途 |
| --- | --- |
| `Resource` | Telemetry producer identity。 |
| `Correlation` | Request/task/workflow/plugin correlation identity。 |
| `Event` | Observability event classification。 |
| `Node` | Service-internal component and operation。 |
| `Tags` | Low-cardinality business labels。 |
| `Fields` | Higher-cardinality event facts。 |

`fields.go` 中的 field name constants 定义稳定 output keys。

## Logger API

| API | 用途 |
| --- | --- |
| `Logger` | Business-facing logging facade。 |
| `Attr` | Opaque attribute application type。 |
| `NewLogger(base, attrs...)` | 把 `ZapCoreBase` 包装成 Dox logger facade。 |
| `ResourceAttr` | 添加 resource fields。 |
| `CorrelationAttr` | 添加 correlation fields。 |
| `EventAttr` | 添加 event fields。 |
| `NodeAttr` | 添加 component and operation fields。 |
| `TagsAttr` | 添加 tags object entries。 |
| `FieldsAttr` | 添加 fields object entries。 |
| `FieldAttr` | 添加一个 fields object entry。 |
| `ErrorAttr` | 添加一个 error field。 |

`NewLogger` 会拒绝 nil `ZapCoreBase`。Logger facade method signatures 不暴露 zap types。

## Context API

| API | 用途 |
| --- | --- |
| `ContextWithCorrelation(ctx, correlation)` | 在 context 上存储 correlation。 |
| `ContextWithMergedCorrelation(ctx, correlation)` | 把 non-empty fields 合并进 context correlation。 |
| `CorrelationFromContext(ctx)` | 取出 correlation 和 boolean。 |
| `MergeCorrelation(base, overlay)` | 把 non-empty overlay fields 应用到 base。 |

Helper functions 会安全处理 nil context。

## Zap API

| API | 用途 |
| --- | --- |
| `ZapCoreBase` | Runtime-owned zap primitives and close lifecycle。 |
| `NewZapLevel(level)` | 把 Dox level 映射到 zapcore level。 |
| `NewZapAtomicLevel(level)` | 创建 zap atomic level。 |
| `NewZapEncoderConfig(config)` | 把 symbolic encoder config 映射到 zapcore encoder config。 |
| `NewZapConfig(config)` | 把 Dox config 映射到 zap config。 |
| `NewZapCoreBase(config)` | 构建 enabled zap cores 和 zap options。 |
| `(*ZapCoreBase).Options()` | 返回复制后的 zap options。 |
| `(*ZapCoreBase).Close()` | 关闭已打开 sinks，一次性执行。 |

Runtime bootstrap 应拥有 `ZapCoreBase` lifecycle，不应把 zap primitives 暴露给 business code。

## OpenTelemetry API

| API | 用途 |
| --- | --- |
| `OpenTelemetrySDKBase` | Runtime-owned OpenTelemetry SDK primitives。 |
| `NewOpenTelemetrySDKBase(config, resource)` | 构建 resource、propagator 和 enabled providers。 |
| `NewOpenTelemetryResource(model, config)` | 把 Dox resource 映射到 SDK resource。 |
| `NewOpenTelemetryPropagator(config)` | 构建 trace context/baggage propagator。 |
| `NewOpenTelemetryTraceSampler(config)` | 构建 SDK trace sampler。 |
| `(*OpenTelemetrySDKBase).ForceFlush(ctx)` | Flush enabled SDK providers。 |
| `(*OpenTelemetrySDKBase).Shutdown(ctx)` | Shutdown enabled SDK providers。 |

SDK base 返回 provider objects。它不安装 globals 或 exporters。

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

这个 sketch 省略 runtime policy，例如 path rendering、global OpenTelemetry installation 和 shutdown orchestration。

## 相关页面

- [契约](contract.md)
- [模型](model.md)
- [Runtime 边界](runtime-boundary.md)
