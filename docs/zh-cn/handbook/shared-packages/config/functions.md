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

  @File    : docs/zh-cn/handbook/shared-packages/config/functions.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Config 函数与 API

本页列出包消费者可以依赖的导出 API。它不是 Go documentation 的替代品，而是解释每个导出符号在包契约中的位置。

## Load 入口

| API | 用途 | 使用场景 |
| --- | --- | --- |
| `Load(ctx, req)` | 用新的默认 loader 执行一个 request。 | Runtime 只需要内置 file/env provider 和 YAML/JSON/TOML/none parser。 |
| `NewDefaultLoader()` | 创建包含内置组件的默认 loader。 | Runtime 需要可复用 loader 实例，但不需要自定义组件。 |
| `NewLoader(LoaderConfig)` | 创建带内置组件和覆盖项的 loader。 | Runtime 需要自定义 provider、parser、merger 或 decoder。 |
| `(*DefaultLoader).Load(ctx, req)` | 用已配置 loader 执行一个 request。 | Runtime 已拥有 loader 实例。 |

简单本地集成使用 `Load`。需要新增自定义 source kind 时使用 `NewLoader`，例如测试中的 memory provider，或由 runtime 自己拥有的 remote provider。

## Loader 配置

| API | 用途 |
| --- | --- |
| `Loader` | 拥有 `Load(ctx, req) (*Result, error)` 的接口。 |
| `DefaultLoader` | 内置实现，编排 provider、parser、merger、decoder 和 fingerprint。 |
| `LoaderConfig` | 注册 provider/parser map，以及可选自定义 merger/decoder。 |

`LoaderConfig.Providers` 和 `LoaderConfig.Parsers` 会叠加到内置组件上。如果 map 中包含内置 key，传入组件会替换该 loader 中的内置组件。

## Provider API

| API | 用途 |
| --- | --- |
| `Provider` | 拥有 `Read(ctx, source) (*Payload, error)` 的接口。 |
| `ProviderFunc` | 让函数实现 `Provider` 的 adapter。 |
| `FileProvider` | 读取本地文件的内置 provider。 |
| `EnvProvider` | 按环境变量前缀读取的内置 provider。 |
| `ProviderKindFile` | 内置 `file` provider kind。 |
| `ProviderKindEnv` | 内置 `env` provider kind。 |
| `ProviderKindRemote` | 命名为 `remote` 的 kind，默认没有注册 provider。 |

Provider 实现应使用 `SourceError` 表达读取失败，使用 `ContractError` 表达 provider/source 不匹配。Provider 不应该 merge、decode 或校验 runtime 级 setting。

## Parser API

| API | 用途 |
| --- | --- |
| `Parser` | 拥有 `Parse(ctx, payload) (map[string]any, error)` 的接口。 |
| `ParserFunc` | 让函数实现 `Parser` 的 adapter。 |
| `NoneParser` | 为 structured payload 返回 provider values。 |
| `YAMLParser` | 解析 YAML object payload。 |
| `JSONParser` | 解析 JSON object payload。 |
| `TOMLParser` | 解析 TOML object payload。 |
| `BuiltinParser(kind)` | 返回内置 parser 和布尔值。 |
| `ParsePayload(ctx, payload)` | 使用 payload 声明的内置 parser 解析 payload。 |

Parser 实现应使用 `ParseError` 表达 source data 格式错误。当 payload 与 parser 实现不匹配时，应返回 `ContractError`。

## Merge API

| API | 用途 |
| --- | --- |
| `Merger` | 拥有 `Merge(ctx, sources, options) (*MergeResult, error)` 的接口。 |
| `DeepReplaceMerger` | 内置 merger。deep-merge maps，并替换 scalars/slices。 |
| `MergeParsedSources(ctx, sources, options)` | 使用 `DeepReplaceMerger` 的便捷函数。 |
| `ParsedSourceFromPayload(payload, values)` | 把 parsed values 与 provider metadata 重新绑定。 |

内置 merger 按 priority 升序排序。它记录 ordered source names、source diagnostics 和 override diagnostics。

## Decode API

| API | 用途 |
| --- | --- |
| `Decoder` | 拥有 `Decode(ctx, values, target, options) error` 的接口。 |
| `MapstructureDecoder` | 基于 mapstructure 的内置 decoder。 |
| `DecodeValues(ctx, values, target, options)` | 使用 `MapstructureDecoder` 的便捷函数。 |
| `DecodeMergeResult(ctx, result, target, options)` | 解码 merge result 中的 values。 |

Decode helper 适合测试或自定义管线。Runtime 通常应优先使用完整 loader，因为它会保留 diagnostics 和 fingerprint。

## Validation API

| API | 用途 |
| --- | --- |
| `ValidateLoadRequest(ctx, req)` | 加载前验证 request shape。 |

当 runtime 想在构造 loader 前或运行 runtime-specific checks 前 fail fast，可以直接使用 validation。Validation 不保证自定义 provider 或 parser 已注册。

## Error API

| API | 用途 |
| --- | --- |
| `Error` | 带有 `Kind`、`Field`、`Reason` 和 wrapped `Err` 的 typed package error。 |
| `ErrorKindContract` | API 或 request 契约错误。 |
| `ErrorKindSource` | Source read 失败。 |
| `ErrorKindParse` | Source parse 失败。 |
| `ErrorKindMerge` | Merge 失败。 |
| `ErrorKindDecode` | Decode 失败。 |
| `ContractError` | 创建 contract error。 |
| `SourceError` | 创建 source error。 |
| `ParseError` | 创建 parse error。 |
| `MergeError` | 创建 merge error。 |
| `DecodeError` | 创建 decode error。 |
| `IsKind(err, kind)` | 判断 error 中是否包含指定 kind 的 config error。 |

调用方应该基于 error kind 分支，而不是解析字符串。

## 数据类型

| Type | 用途 |
| --- | --- |
| `Request` | 一次 load 操作。 |
| `Options` | Load、merge、decode 和 diagnostics 行为。 |
| `Source` | 一个 provider source 描述。 |
| `Payload` | Parser 转换前的 provider 输出。 |
| `ParsedSource` | 带 source metadata 的 parsed source values。 |
| `MergeResult` | Merged values、source names 和 diagnostics。 |
| `Result` | Decode 和 fingerprint 成功后的 load 输出。 |
| `Diagnostics` | Source diagnostics 和 override diagnostics。 |
| `SourceDiagnostic` | 一个 source 参与记录。 |
| `MergeOverride` | 一个 override 记录。 |

这些类型都是普通 Go struct。包不会把 source shape 隐藏在全局状态后面。

## 常见用法

### 默认本地加载

使用 `Load(ctx, req)` 搭配 file 和 env sources。这是本地 runtime 启动的常规路径。

### 测试用 memory source

使用 `NewLoader(LoaderConfig{Providers: ..., Parsers: ...})` 搭配自定义 kind 和 `ProviderFunc` / `ParserFunc`。这样测试不需要依赖文件系统或进程环境状态。

### Merge 前检查

只有在构建自定义 runtime 管线时，才需要直接使用 providers、parsers 和 `MergeParsedSources`。大多数消费者不需要这样做。

### Decode-only 测试

测试 target struct tags、duration conversion 或 unknown-key 行为时，可以使用 `DecodeValues` 或 `DecodeMergeResult`，不需要读取 source。

## 调用方责任

每个消费者必须自己负责：

- setting target type 定义；
- runtime-specific defaults；
- decode 后的领域校验；
- source list 选择；
- 环境变量 prefix 策略；
- secret 和 redaction 策略；
- 启动日志策略；
- hot reload、remote configuration 或 watcher 生命周期。

共享包提供的是确定性的配置加载原语。它不替代 runtime 配置设计。
