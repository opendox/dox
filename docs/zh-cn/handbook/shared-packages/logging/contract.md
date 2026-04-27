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

  @File    : docs/zh-cn/handbook/shared-packages/logging/contract.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Logging 契约

Shared logging 契约定义调用方可以从 `packages/shared/logging` 依赖什么，哪些行为目前只是配置形状，以及哪些 error 属于包边界。

## Ownership Contract

| 区域 | Package Owns | Runtime Owns |
| --- | --- | --- |
| Logging vocabulary | Resource、correlation、event、node、tags、fields model。 | Runtime-specific naming policy 和 business event taxonomy。 |
| Business logging API | `Logger`、`Attr` 和 typed attribute constructors。 | 把 facade 传给 services 和 middleware。 |
| zap mapping | Level、encoder、config、console core、JSON/file core、lumberjack helpers。 | Output paths、lifecycle、deployment topology 和 sink ownership。 |
| OpenTelemetry SDK base | Resource mapping、propagator mapping、SDK provider construction、flush/shutdown helpers。 | Installing globals、exporter wiring 和 lifecycle integration。 |
| Config validation | Defaults 和 package validation errors。 | Runtime-specific configuration restrictions 和 environment-specific policy。 |

> [!NOTE]
> 这个包是共享基础设施。它不会自己为任何进程启动 logging。

## Config 契约

`Config` 是 root shared logging configuration shape。

| Group | 用途 |
| --- | --- |
| `Level` | Root Dox logging level。 |
| `Development` | Root development-mode switch，会进入 zap defaults。 |
| `Resource` | Logging-specific resource overrides，目前是 `service_version`。 |
| `Zap` | zap-facing config shape 和 encoder symbols。 |
| `Cores` | Console/file core declarations。 |
| `Buffering` | Future buffered writer behavior。 |
| `Shutdown` | Flush/shutdown timeout。 |
| `Redaction` | Sensitive key replacement policy shape。 |
| `OTel` | OpenTelemetry propagation、provider、exporter、batch settings。 |

`Config.Default` 先填默认值，再进行 validation。`Config.Validate` 检查受支持 enum values、required paths、duplicate core names、positive durations and sizes、redaction key shape、OTLP protocol shape 和 batch constraints。

## Default 契约

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
| OTLP exporter | 默认 disabled |

显式 disabled 的 pointer toggles 在 defaults 后仍保持 disabled。

## Validation Error 契约

包暴露：

- `FieldError`，包含 `Field` 和 `Reason`；
- `ValidationError`，包含 `Fields []FieldError`；
- 适合日志的紧凑 error string。

调用方应检查 `ValidationError.Fields`，不要解析 error string。

<details>
<summary>Validation fields 示例</summary>

```text
level: level is not supported
cores[1].name: core name must be unique
redaction.keys[0]: redaction key must not be empty
otel.traces.sampler.ratio: trace sampler ratio must be between 0 and 1
```

</details>

## 已实现与仅配置契约

| Setting | Validation | 当前 Runtime Effect |
| --- | --- | --- |
| `cores[].datasets` | 校验为 non-empty entries。 | 尚未用于 event dataset routing。 |
| `buffering.*` | enabled 时校验。 | 尚未安装 buffered writer。 |
| `redaction.*` | enabled 时校验。 | 尚未对 log field 或 value 执行 redaction。 |
| `cores[].rotation.driver=external/logrotate` | 合法 enum value。 | 当前 zap file sink 会拒绝。 |
| `otel.exporter.otlp.*` | Shape 会校验。 | Enabled OTLP exporter 会被 `NewOpenTelemetrySDKBase` 拒绝。 |
| `zap.sampling.hook_metrics` | 配置形状的一部分。 | 尚未安装 sampling hook metrics。 |
| `DefaultFilePathTemplate` | 作为默认字符串使用。 | Template variables 不会被渲染。 |

> [!WARNING]
> 仅配置契约应被视为 forward-compatible schema，而不是 active runtime behavior。

## Logger 契约

Business-facing code 使用 `Logger` 和 `Attr`。

Facade 保证：

- log level methods: `Debug`, `Info`, `Warn`, `Error`, `DPanic`, `Panic`, `Fatal`;
- logger derivation: `Named`, `With`;
- flush path: `Sync`;
- Dox-owned attributes，用于 resource、correlation、event、node、tags、fields、single fields 和 errors；
- 每次 write 时合并 context correlation。

Facade 刻意不在 method signatures 中暴露 `zap.Logger` 或 `zap.Field`。

## Context Correlation 契约

Correlation helpers 提供简单 overlay semantics：

- `ContextWithCorrelation` 存储 correlation value；nil context 会使用 `context.Background`。
- `ContextWithMergedCorrelation` 把 non-empty overlay fields 合并进已有 context correlation。
- `CorrelationFromContext` 对 nil 或缺失 context 返回 `(Correlation, false)`。
- `MergeCorrelation` 保留 base fields，除非 overlay fields 非空。

## 契约之外

以下行为当前不属于包契约：

- server/scheduler/collector/compute startup 时打开 runtime loggers；
- HTTP middleware correlation；
- task/job/plugin correlation injection；
- OpenTelemetry global provider installation；
- OTLP exporter construction；
- collector deployment examples；
- runtime-specific sink path rendering；
- multi-process rotation coordination；
- field-level redaction execution；
- dataset-based core routing。

## 相关页面

- [模型](model.md)
- [Runtime 边界](runtime-boundary.md)
- [函数与 API](functions.md)
