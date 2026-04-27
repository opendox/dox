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

  @File    : docs/zh-cn/handbook/shared-packages/setting/functions.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# 第 3 章：Shared Setting 函数与 API

| 上一章 | 上级 | 下一章 |
| --- | --- | --- |
| [第 2 章：模型](model.md) | [Shared setting 包](README.md) | 结束 |

> [!NOTE]
> 本章是这个包的 API reference。如果某个行为依赖模型语义，相关模型章节会在正文中链接。

## 常量

| API | Value | 用途 |
| --- | --- | --- |
| `DefaultOrganizationName` | `opendox` | 默认 organization identity。 |
| `DefaultApplicationName` | `dox` | 默认 product 或 application family identity。 |

## Runtime API

| API | 用途 |
| --- | --- |
| `type Runtime string` | Dox runtime identity value。 |
| `RuntimeServer` | Web backend runtime。 |
| `RuntimeScheduler` | Scheduling runtime。 |
| `RuntimeCollector` | Collection runtime。 |
| `RuntimeCompute` | Computation runtime。 |
| `(Runtime).IsValid() bool` | 对支持的 runtime value 返回 true。 |

Runtime validation 是 shared 的，但 runtime selection 不是。比如 server package 可以要求 `RuntimeServer`，scheduler package 应该要求 `RuntimeScheduler`。

## Env API

| API | 用途 |
| --- | --- |
| `type Env string` | Deployment environment value。 |
| `EnvDev` | Development environment。 |
| `EnvTest` | Test environment。 |
| `EnvStaging` | Staging environment。 |
| `EnvProd` | Production environment。 |
| `(Env).IsValid() bool` | 对支持的 environment value 返回 true。 |

## Fragment API

| Fragment | Default Method | Validate Method |
| --- | --- | --- |
| `Organization` | `(*Organization).Default() error` | `(Organization).Validate() error` |
| `Application` | `(*Application).Default() error` | `(Application).Validate() error` |
| `System` | `(*System).Default() error` | `(System).Validate() error` |
| `Service` | `(*Service).Default(application, system) error` | `(Service).Validate() error` |
| `Deployment` | `(*Deployment).Default() error` | `(Deployment).Validate() error` |

Default methods 在 nil receiver 上会返回 error。Validate methods 调用 package-level `Validate` helper。

## Default Method 行为

| Method | 行为 |
| --- | --- |
| `Organization.Default` | 把空 `Name` 设置为 `DefaultOrganizationName`。 |
| `Application.Default` | 把空 `Name` 设置为 `DefaultApplicationName`。 |
| `System.Default` | 只检查 receiver，不设置 `Runtime`。 |
| `Service.Default` | 从 `application.Name` 设置空 `Namespace`；runtime 已知时从 `system.Runtime` 设置空 `Name`。 |
| `Deployment.Default` | 把空 `Env` 设置为 `EnvDev`。 |

<details>
<summary>示例：runtime aggregate 中的 default 顺序</summary>

```go
func (i *Identity) Default() error {
	if err := i.Organization.Default(); err != nil {
		return err
	}
	if err := i.Application.Default(); err != nil {
		return err
	}
	if err := i.System.Default(); err != nil {
		return err
	}
	if i.System.Runtime == "" {
		i.System.Runtime = sharedsetting.RuntimeServer
	}
	if err := i.Service.Default(i.Application, i.System); err != nil {
		return err
	}
	return i.Deployment.Default()
}
```

设置 `RuntimeServer` 的那一行由 runtime package 自己拥有。

</details>

## Validation API

| API | 用途 |
| --- | --- |
| `Validate(value any) error` | 使用已注册的 Dox validation tags 校验 struct。 |
| `FieldError` | 带 `Field` 和 `Rule` 的机器可读字段失败。 |
| `ValidationError` | 字段失败列表。 |
| `(*ValidationError).Error() string` | 适合日志的紧凑 error message。 |

Validation engine 初始化一次，并注册：

- `dox_kebab`;
- `dox_identifier`;
- `dox_runtime`;
- `dox_env`.

字段名优先来自 `mapstructure` tags。这能让 validation output 与 decoded configuration keys 对齐。

## 调用方 Recipes

### 组合 runtime identity group

定义 runtime-owned group，用 shared fragments 作为字段或嵌入。不要把 `packages/shared/setting` 暴露成整个 runtime setting aggregate。

### 添加 runtime-specific validation

把 shared validation 和 runtime-specific validation 合并：

```go
func (i Identity) Validate() error {
	return errors.Join(
		i.Organization.Validate(),
		i.Application.Validate(),
		i.System.Validate(),
		i.validateServerRuntime(),
		i.Service.Validate(),
		i.Deployment.Validate(),
	)
}
```

### 保持 runtime bootstrap 分离

Bootstrap 可以提供 env、instance ID 或 region 等 seed values。Shared package 不应该直接读取 flags、environment variables、files 或 process metadata。

## 调用方责任

每个 runtime consumer 必须自己拥有：

- root `Setting` aggregate；
- group composition；
- runtime identity selection；
- runtime-specific defaults；
- 依赖多个 fragments 组合关系的 validation；
- bootstrap-derived seed values；
- 更严格 runtime behavior 的文档。

## API 稳定性备注

- Shared fragments 设计目标是在 Dox runtimes 之间复用。
- Validation tags 是 public package behavior 的一部分。
- 第三方 validator 实现不会作为 caller-facing error type 暴露。
- 新增 runtime 或 environment value 会改变 shared validation semantics，应视为契约变更。

## 导航

| 上一章 | 上级 | 下一章 |
| --- | --- | --- |
| [第 2 章：模型](model.md) | [Shared setting 包](README.md) | 结束 |
