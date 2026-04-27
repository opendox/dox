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

  @File    : docs/en-us/handbook/shared-packages/config/README.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Config Package Handbook

`packages/shared/config` is the shared Dox configuration loading SDK. It gives each backend runtime an explicit way to read declared sources, parse source payloads, merge values, decode into caller-owned targets, and return operational diagnostics.

This handbook is written for developers and coding agents. Treat it as the package-level contract for consumers in Web, Scheduling, Collection, Computation, and any later runtime that imports `github.com/opendox/dox/packages/shared/config`.

## Table of Contents

- [Contract](contract.md): request, source, option, result, diagnostics, and error semantics.
- [Pipeline](pipeline.md): provider, parser, merge, decode, fingerprint, and diagnostics flow.
- [Functions and API](functions.md): exported entry points, interfaces, helpers, and caller obligations.

## Package Position

The package owns the loading pipeline contract. It validates API usage before work starts, executes the pipeline for one explicit `Request`, and returns a `Result` that callers can log or inspect.

The package does not own runtime-specific setting validation. Runtime packages define their own setting structs, field constraints, default policy, and operational rules after this package has decoded raw configuration values.

## Current Capability

The current implementation includes:

- local file and environment providers;
- YAML, JSON, TOML, and `none` parsers;
- deep map merge with scalar and slice replacement;
- source ordering by ascending `Priority`;
- environment dotted-key expansion before merge;
- decode into struct pointers or map pointers through `mapstructure`;
- reject-by-default unknown key handling;
- stable `sha256:` fingerprints over merged values;
- source diagnostics and override diagnostics;
- extension points for custom providers, parsers, mergers, and decoders.

`ProviderKindRemote` exists as a named kind, but no remote provider is registered by the default loader. Callers must register a provider before using remote or other custom source kinds.

## Current Non-capability

The package currently does not implement:

- runtime-specific setting validation;
- file watching or hot reload;
- remote provider reads in the default loader;
- default file path discovery;
- secret loading;
- schema generation;
- value redaction enforcement for `Options.RedactKeys`.

`Options.RedactKeys` is accepted as part of the option shape, but the current pipeline does not apply redaction to values, diagnostics, errors, or fingerprints. Do not treat it as a sanitization guarantee.

## Basic Usage

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

This example reads a required YAML file first, then applies matching environment variables as higher-priority overrides. The caller still owns setting defaults, domain validation, and runtime startup behavior.

## Consumer Rules

Consumers should:

- create a typed setting target owned by the runtime package;
- declare every source explicitly;
- use unique source names and unique priorities;
- keep lower-priority base files before higher-priority overrides;
- use `ParserKindNone` for environment sources;
- inspect typed errors with `IsKind`;
- persist or log `Result.SourceNames`, `Result.Fingerprint`, and diagnostics when operational traceability is needed.

Consumers should not:

- rely on undeclared default sources;
- pass a nil context or nil target;
- use `ProviderKindRemote` without registering a provider;
- expose Koanf or mapstructure as part of their own public runtime contract;
- assume `RedactKeys` removes sensitive values.

## Reading Order

Read these pages in order when implementing a new runtime integration:

1. [Contract](contract.md) to understand valid requests and failure classes.
2. [Pipeline](pipeline.md) to understand how values move and override each other.
3. [Functions and API](functions.md) to choose the entry point and extension points.
