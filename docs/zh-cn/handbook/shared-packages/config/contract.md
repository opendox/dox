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

  @File    : docs/zh-cn/handbook/shared-packages/config/contract.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Config 契约

config 包的契约围绕一个显式 `Request` 展开。调用方提供 runtime 标识、environment 标识、target pointer、source 描述列表和 options。loader 会在读取 source 之前验证这些契约。

## Request 契约

| 字段 | 是否必需 | 语义 |
| --- | --- | --- |
| `Runtime` | 是 | 小写 runtime 名称。必须以字母开头，可包含小写字母、数字和连字符。 |
| `Env` | 是 | 小写环境名称。命名规则和 `Runtime` 相同。 |
| `Target` | 是 | 指向 struct 或 map 的非 nil pointer。调用方拥有 target 类型和后续校验。 |
| `Sources` | 通常是 | 要读取的 source 列表。除非 `Options.AllowEmptySources` 为 true，否则空 source 会被拒绝。 |
| `Options` | 否 | 控制管线行为。空 options 会被规范化为包默认值。 |

这个包只验证 request 形状。它不校验端口、业务开关、租户策略、部署规则等领域 setting 值。

## Source 契约

每个 `Source` 声明一次 provider 读取。

| 字段 | 是否必需 | 语义 |
| --- | --- | --- |
| `Name` | 是 | 唯一 source 名称。必须以字母开头，可包含小写字母、数字、连字符、下划线和点。 |
| `Kind` | 是 | Provider kind。内置 kind 是 `file` 和 `env`；自定义名称可通过验证，但加载前必须注册。 |
| `Parser` | 是 | Parser kind。内置 kind 是 `yaml`、`json`、`toml` 和 `none`；自定义名称可通过验证，但加载前必须注册。 |
| `Location` | 是 | Provider 相关位置。文件 source 是文件系统路径；环境变量 source 是必需前缀。 |
| `Required` | 否 | 必需 source 不可读或缺失会失败。可选缺失 file/env source 会被跳过并留下 diagnostics。 |
| `Priority` | 否 | 非负 merge 优先级。数字越小越早 merge；后续高优先级 source 覆盖前面的值。priority 必须唯一。 |
| `Options` | 否 | Provider 级 string options。当前内置 provider 不消费 source options。 |

Provider 与 parser 的兼容性属于 request 契约：

- `ProviderKindEnv` 必须使用 `ParserKindNone`。
- `ProviderKindFile` 必须使用非 `ParserKindNone` 的 parser。
- 自定义 provider 和 parser kind 必须通过命名校验，并在 load 前注册到 `LoaderConfig`。

## 内置 Kind

| Kind | 值 | 默认注册 | 说明 |
| --- | --- | --- | --- |
| `ProviderKindFile` | `file` | 是 | 读取本地文件为 raw bytes。 |
| `ProviderKindEnv` | `env` | 是 | 按前缀读取环境变量，并返回 structured dotted keys。 |
| `ProviderKindRemote` | `remote` | 否 | 为未来或自定义场景命名。默认 loader 不读取 remote source。 |
| `ParserKindNone` | `none` | 是 | 不做 raw parse，直接返回 provider values。环境变量 source 使用它。 |
| `ParserKindYAML` | `yaml` | 是 | 解析 YAML object payload。 |
| `ParserKindJSON` | `json` | 是 | 解析 JSON object payload，并保留 JSON number。 |
| `ParserKindTOML` | `toml` | 是 | 解析 TOML object payload。 |

## Option 契约

| 字段 | 默认值 | 语义 |
| --- | --- | --- |
| `AllowEmptySources` | `false` | 为 true 时允许空 source 列表。 |
| `MergeStrategy` | `deep_replace` | 当前只支持 `MergeStrategyDeepReplace`。 |
| `UnknownKeyPolicy` | `reject` | `reject` 遇到未使用 key 会失败；`allow` 会忽略未知 key。 |
| `Timeout` | `0` | 为正数时，对整个 load 操作应用一个 context timeout。负数非法。 |
| `RedactKeys` | 空 | 当前实现接受这个字段，但不执行脱敏。 |

不支持的 merge strategy 和 unknown-key policy 会返回 `ErrorKindContract`。

## Result 契约

只有成功 load 并 decode 后才会返回 `Result`。

| 字段 | 语义 |
| --- | --- |
| `Runtime` | 从 request 复制。 |
| `Env` | 从 request 复制。 |
| `SourceNames` | 按 priority 排序后的 source 名称。被跳过的可选 source 仍会出现。 |
| `Fingerprint` | Decode 成功后，对 merge 后 value map 计算得到的 `sha256:` fingerprint。 |
| `Diagnostics` | source 参与情况和 override diagnostics。 |

Fingerprint 代表 merge 后 structured values。它不是 secret-safe digest 契约，因为当前实现不会应用 `RedactKeys`。

## Diagnostics 契约

`Diagnostics.Sources` 记录每个 source 如何参与 merge。

| 字段 | 语义 |
| --- | --- |
| `Name` | Source 名称。 |
| `Kind` | Provider kind。 |
| `Required` | Source 的 required 标记。 |
| `Loaded` | provider 加载到数据，或 merge 推断出已有 values 时为 true。 |
| `Skipped` | 可选缺失 source 被跳过时为 true。 |
| `Message` | skip 或 provider 状态的人类可读说明。 |

`Diagnostics.Overrides` 记录值替换事件。

| 字段 | 语义 |
| --- | --- |
| `Path` | 被替换的 dot path。 |
| `Source` | 提供替换值的 source。 |
| `PreviousSource` | 之前提供旧值的 source。 |

嵌套 map 会 deep merge。scalar 和 slice 会被替换。当后续 source 改变之前 source 拥有的路径时，会记录 override diagnostics。

## Error 契约

包内分类错误都使用 `*config.Error`。

| Kind | 含义 |
| --- | --- |
| `ErrorKindContract` | 调用方违反 API 或 request 规则。 |
| `ErrorKindSource` | Provider 无法读取声明的 source。 |
| `ErrorKindParse` | Parser 无法转换 source payload。 |
| `ErrorKindMerge` | Parsed source values 无法合并。 |
| `ErrorKindDecode` | Merged values 无法解码到 target。 |

使用 `IsKind(err, kind)` 或 `errors.Is(err, &config.Error{Kind: kind})` 判断错误类型，不要解析错误字符串。错误字符串用于人类阅读，不应作为机器契约。

## 扩展契约

包通过 `LoaderConfig` 支持自定义管线组件：

- `Providers` 注册额外 `ProviderKind` handler，或覆盖内置 provider。
- `Parsers` 注册额外 `ParserKind` handler，或覆盖内置 parser。
- `Merger` 替换默认 deep-replace merger。
- `Decoder` 替换默认 mapstructure decoder。

自定义组件必须保持调用方期望的公共语义。默认 loader 仍会在调用自定义 provider 和 parser 前验证 request 形状。

## 契约之外

以下行为不属于这个包的契约：

- runtime 级 setting 默认值和校验；
- 从进程参数、部署清单或服务发现中选择 source list；
- remote read，除非调用方注册并拥有 provider；
- hot reload 和 watcher 生命周期；
- secret 解析；
- 对返回 values 或 diagnostics 做脱敏；
- 把 Koanf 或 mapstructure 行为保证成 runtime 公开 API。
