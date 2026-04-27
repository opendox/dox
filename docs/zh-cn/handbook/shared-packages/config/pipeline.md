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

  @File    : docs/zh-cn/handbook/shared-packages/config/pipeline.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Config 管线

默认 loader 会把一个 request 送入固定的本地管线：

```text
Request
  -> validate request and normalize options
  -> apply request timeout when configured
  -> provider read for each source
  -> parser parse for each payload
  -> merge parsed sources by priority
  -> decode merged values into target
  -> fingerprint merged values
  -> return Result
```

每个阶段都会接收 request context。如果 `Options.Timeout` 为正数，loader 会在 provider read 开始前用该 timeout 包装 context。

## 阶段 1：Validation

`ValidateLoadRequest` 在读取 source 前执行。它检查：

- context 非 nil；
- `Runtime` 和 `Env` 合法；
- target 是指向 struct 或 map 的非 nil pointer；
- option 值受支持；
- 除非 `AllowEmptySources` 为 true，否则 source 必须存在；
- source name 唯一；
- source priority 唯一；
- source name、provider kind、parser kind、location 和 priority 合法；
- 内置 file/env provider 与 parser 的兼容性。

Validation 允许符合命名规则的自定义 provider 和 parser kind。若 load 时没有注册对应组件，后续仍会以 `ErrorKindContract` 失败。

## 阶段 2：Provider Read

Provider 返回 `Payload`。Payload 可以包含 `Raw` bytes、structured `Values`、provider `Metadata` 和 `SourceDiagnostic`。

### File Provider

`FileProvider` 从本地文件系统读取 `Source.Location`。

- 必需文件缺失返回 `ErrorKindSource`。
- 可选文件缺失返回 skipped payload。
- location 指向目录会返回 `ErrorKindSource`。
- 成功读取时填充 `Raw`、文件路径 metadata 和 loaded diagnostic。

File provider 不解析 bytes，只负责读取 source。

### Environment Provider

`EnvProvider` 优先从 `Lookup` 读取环境变量条目；未提供时读取 `os.Environ`。

它使用 `Source.Location` 作为前缀过滤条目。匹配的 key 会这样规范化：

- 移除 source prefix；
- 裁剪首尾下划线；
- 转小写；
- 把双下划线替换成点；
- 把剩余下划线替换成点。

示例：

| 环境变量 key | Prefix | 规范化 key |
| --- | --- | --- |
| `DOX_SERVER_APP_NAME` | `DOX_SERVER_` | `app.name` |
| `DOX_SERVER_HTTP__ADDRESS` | `DOX_SERVER_` | `http.address` |

没有匹配变量时，required env source 返回 `ErrorKindSource`；optional env source 被跳过并留下 diagnostics。

## 阶段 3：Parser

Parser 把一个 provider payload 转换成 `map[string]any`。

| Parser | 输入 | 行为 |
| --- | --- | --- |
| `NoneParser` | `Payload.Values` | 返回 shallow copy。用于 env source。 |
| `YAMLParser` | `Payload.Raw` | 解析 YAML，并要求 root 是 object。 |
| `JSONParser` | `Payload.Raw` | 解析 JSON，拒绝 trailing data，保留 number 为 `json.Number`，并要求 root 是 object。 |
| `TOMLParser` | `Payload.Raw` | 解析 TOML，并要求 root 是 object。 |

被跳过的 optional payload 会解析为空 map。格式错误返回 `ErrorKindParse`。Parser 实现与 payload kind 不匹配返回 `ErrorKindContract`。

## 阶段 4：Merge

`DeepReplaceMerger` 按 `Source.Priority` 升序排列 parsed sources。优先级数字越小越早 merge。后续高优先级 source 覆盖之前的值。

Merge 行为：

- 嵌套 map deep merge；
- scalar 替换之前的 scalar；
- slice 替换之前的 slice；
- skipped source 会出现在 `SourceNames` 和 source diagnostics 中，但不会改变 values；
- 重复 source name 或 priority 返回 `ErrorKindContract`；
- 使用 `ParserKindNone` 的 env source 会在 merge 前把 dotted keys 展开成嵌套 map。

示例：

```text
base priority 10:
  app.name = dox
  http.timeout = 5s

env priority 100:
  app.name = dox-env
  http.timeout = 10s

merged:
  app.name = dox-env
  http.timeout = 10s
```

结果会为 `app.name` 和 `http.timeout` 记录 override diagnostics，表示 `env` 替换了 `base`。

冲突的环境变量 dotted keys 会返回 `ErrorKindMerge`。例如同一个 env source 不能同时提供 `app=dox` 和 `app.name=dox`，因为同一路径不能既是 scalar 又是 map。

## 阶段 5：Decode

`MapstructureDecoder` 把 merged values 解码到调用方 target。

Decode 行为：

- target 必须是指向 struct 或 map 的非 nil pointer；
- struct tag 使用 `mapstructure`；
- 启用 weak type conversion；
- 启用 string 到 `time.Duration` 的转换；
- decode 前会 zero target fields；
- 默认拒绝未知 key；
- `UnknownKeyPolicyAllow` 允许未知 key。

Decode 失败返回 `ErrorKindDecode`。target 非法返回 `ErrorKindContract`。

## 阶段 6：Fingerprint 和 Result

Decode 成功后，loader 会根据 merged structured value map 计算稳定的 `sha256:` fingerprint。Result 包含：

- request runtime；
- request environment；
- 排序后的 source names；
- fingerprint；
- source diagnostics；
- merge override diagnostics。

Fingerprint 代表 merged values，而不是 decoded target。如果 decode 失败，不会返回 result。

## 运维备注

Source name 应体现用途，例如 `base`、`local`、`secrets` 或 `env`。Priority 建议留出间隔，例如 `10`、`20`、`100`，方便之后插入新的 override 层。

为了保持确定性：

- 避免重复 priority；
- 文件格式保持 object root；
- 环境变量 prefix 按 runtime 隔离；
- 启动成功后记录 source names 和 fingerprint；
- 排查覆盖行为时记录 diagnostics。

不要记录 raw merged values，除非调用方已经应用自己的脱敏策略。

## 相关页面

- [Shared config 包手册](README.md)
- [Shared config 契约](contract.md)
- [Shared config 函数与 API](functions.md)
