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

  @File    : docs/zh-cn/handbook/shared-packages/config/README.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Config 包手册

`packages/shared/config` 是 Dox 共享的配置加载 SDK。它为各个后端 runtime 提供一套显式流程：读取声明的配置源、解析 payload、合并值、解码到调用方拥有的目标对象，并返回可观察诊断信息。

这份手册面向开发者和编码 Agent。Web、Scheduling、Collection、Computation 以及后续 runtime 引用 `github.com/opendox/dox/packages/shared/config` 时，都应把这里当作包级契约。

## 目录

- [契约](contract.md)：`Request`、`Source`、`Options`、`Result`、诊断和错误语义。
- [管线](pipeline.md)：provider、parser、merge、decode、fingerprint 和 diagnostics 流程。
- [函数与 API](functions.md)：导出入口、接口、辅助函数和调用方责任。

## 包定位

这个包负责配置加载管线契约。它在工作开始前验证 API 使用方式，为一个显式 `Request` 执行加载流程，并返回调用方可以记录或检查的 `Result`。

这个包不负责 runtime 级 setting 校验。各 runtime 包在解码出原始配置值后，自己定义 setting 结构体、字段约束、默认策略和运行规则。

## 当前能力

当前实现包含：

- 本地文件 provider 和环境变量 provider；
- YAML、JSON、TOML 和 `none` parser；
- map 深度合并，scalar 和 slice 替换；
- 按 `Priority` 升序排列 source；
- merge 前展开环境变量 dotted key；
- 通过 `mapstructure` 解码到 struct pointer 或 map pointer；
- 默认拒绝未知 key；
- 基于合并后值的稳定 `sha256:` fingerprint；
- source diagnostics 和 override diagnostics；
- 自定义 provider、parser、merger、decoder 的扩展点。

`ProviderKindRemote` 作为命名 kind 已存在，但默认 loader 没有注册 remote provider。调用方要使用 remote 或其他自定义 source kind，必须先注册 provider。

## 当前非能力

这个包当前不实现：

- runtime 级 setting 校验；
- 文件监听或热重载；
- 默认 loader 中的 remote provider 读取；
- 默认文件路径发现；
- secret 加载；
- schema 生成；
- 对 `Options.RedactKeys` 的值脱敏执行。

`Options.RedactKeys` 目前只是 option 结构的一部分，当前管线不会把它应用到 values、diagnostics、errors 或 fingerprints。不要把它当作脱敏保证。

## 基本用法

```go
package runtime

import (
	"context"
	"time"

	sharedconfig "github.com/opendox/dox/packages/shared/config"
)

type Setting struct {
	App struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"app"`
	HTTP struct {
		Port    int           `mapstructure:"port"`
		Timeout time.Duration `mapstructure:"timeout"`
	} `mapstructure:"http"`
}

func LoadSetting(ctx context.Context, basePath string) (*Setting, *sharedconfig.Result, error) {
	var target Setting
	result, err := sharedconfig.Load(ctx, sharedconfig.Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []sharedconfig.Source{
			{
				Name:     "base",
				Kind:     sharedconfig.ProviderKindFile,
				Parser:   sharedconfig.ParserKindYAML,
				Location: basePath,
				Required: true,
				Priority: 10,
			},
			{
				Name:     "env",
				Kind:     sharedconfig.ProviderKindEnv,
				Parser:   sharedconfig.ParserKindNone,
				Location: "DOX_SERVER_",
				Required: false,
				Priority: 100,
			},
		},
		Options: sharedconfig.Options{
			UnknownKeyPolicy: sharedconfig.UnknownKeyPolicyReject,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	return &target, result, nil
}
```

这个例子先读取一个必需 YAML 文件，再用匹配前缀的环境变量作为更高优先级覆盖。调用方仍然负责 setting 默认值、领域校验和 runtime 启动行为。

## 调用方规则

调用方应该：

- 创建由 runtime 包自己拥有的 typed setting target；
- 显式声明每一个 source；
- 使用唯一 source name 和唯一 priority；
- 让低优先级 base 文件排在高优先级 override 之前；
- 环境变量 source 使用 `ParserKindNone`；
- 用 `IsKind` 检查 typed error；
- 需要运维可追踪性时，记录 `Result.SourceNames`、`Result.Fingerprint` 和 diagnostics。

调用方不应该：

- 依赖未声明的默认 source；
- 传入 nil context 或 nil target；
- 不注册 provider 就使用 `ProviderKindRemote`；
- 把 Koanf 或 mapstructure 暴露成 runtime 自己的公开契约；
- 假设 `RedactKeys` 会移除敏感值。

## 阅读顺序

实现新的 runtime 集成时，建议按顺序阅读：

1. [契约](contract.md)，理解合法 request 和失败分类。
2. [管线](pipeline.md)，理解值如何流动以及如何覆盖。
3. [函数与 API](functions.md)，选择入口函数和扩展点。
