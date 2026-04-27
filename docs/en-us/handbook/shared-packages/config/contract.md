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

  @File    : docs/en-us/handbook/shared-packages/config/contract.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Config Contract

The config package contract is centered on one explicit `Request`. A caller provides runtime identity, environment identity, a target pointer, ordered source descriptors, and options. The loader validates that contract before reading sources.

## Request Contract

| Field | Required | Semantics |
| --- | --- | --- |
| `Runtime` | Yes | Lowercase runtime name. It must start with a letter and may contain lowercase letters, digits, and hyphens. |
| `Env` | Yes | Lowercase environment name. It follows the same naming rule as `Runtime`. |
| `Target` | Yes | Non-nil pointer to a struct or map. The caller owns the target type and later validation. |
| `Sources` | Usually | Source list to read. Empty sources are rejected unless `Options.AllowEmptySources` is true. |
| `Options` | No | Pipeline behavior. Empty options are normalized to package defaults. |

The package validates request shape only. It does not validate domain-specific setting values such as ports, business toggles, tenant policy, or runtime deployment rules.

## Source Contract

Each `Source` declares one provider read.

| Field | Required | Semantics |
| --- | --- | --- |
| `Name` | Yes | Unique source name. It must start with a letter and may contain lowercase letters, digits, hyphens, underscores, and dots. |
| `Kind` | Yes | Provider kind. Built-in kinds are `file` and `env`; custom names are accepted by validation but must be registered before loading. |
| `Parser` | Yes | Parser kind. Built-in kinds are `yaml`, `json`, `toml`, and `none`; custom names are accepted by validation but must be registered before loading. |
| `Location` | Yes | Provider-specific location. For files this is a filesystem path. For env sources this is the required variable prefix. |
| `Required` | No | Required sources fail when unreadable or missing. Optional missing file/env sources are skipped with diagnostics. |
| `Priority` | No | Non-negative merge priority. Lower numbers merge first; later higher-priority sources override earlier values. Priorities must be unique. |
| `Options` | No | Provider-specific string options. The current built-in providers do not consume source options. |

Provider and parser compatibility is part of the request contract:

- `ProviderKindEnv` must use `ParserKindNone`.
- `ProviderKindFile` must use a parser other than `ParserKindNone`.
- Custom provider and parser kinds must pass naming validation and must be registered in `LoaderConfig` before load time.

## Built-in Kinds

| Kind | Value | Registered by default | Notes |
| --- | --- | --- | --- |
| `ProviderKindFile` | `file` | Yes | Reads local files as raw bytes. |
| `ProviderKindEnv` | `env` | Yes | Reads environment variables by prefix and returns structured dotted keys. |
| `ProviderKindRemote` | `remote` | No | Named for future or custom use. The default loader does not read remote sources. |
| `ParserKindNone` | `none` | Yes | Returns provider values without raw parsing. Used by env sources. |
| `ParserKindYAML` | `yaml` | Yes | Parses YAML object payloads. |
| `ParserKindJSON` | `json` | Yes | Parses JSON object payloads and preserves JSON numbers. |
| `ParserKindTOML` | `toml` | Yes | Parses TOML object payloads. |

## Option Contract

| Field | Default | Semantics |
| --- | --- | --- |
| `AllowEmptySources` | `false` | Allows an empty source list when true. |
| `MergeStrategy` | `deep_replace` | Only `MergeStrategyDeepReplace` is currently supported. |
| `UnknownKeyPolicy` | `reject` | `reject` fails on unused decoded keys. `allow` ignores unknown keys. |
| `Timeout` | `0` | When positive, applies one context timeout to the full load operation. Negative durations are invalid. |
| `RedactKeys` | Empty | Accepted but not enforced by the current implementation. |

Unsupported merge strategies and unknown-key policies fail with `ErrorKindContract`.

## Result Contract

`Result` is returned only after a successful load and decode.

| Field | Semantics |
| --- | --- |
| `Runtime` | Copied from the request. |
| `Env` | Copied from the request. |
| `SourceNames` | Source names after priority ordering. Skipped optional sources are still represented. |
| `Fingerprint` | `sha256:` fingerprint computed from the merged value map after decode succeeds. |
| `Diagnostics` | Source participation and override diagnostics. |

The fingerprint is computed from the merged structured values. It is not a secret-safe digest contract, because the current implementation does not apply `RedactKeys`.

## Diagnostics Contract

`Diagnostics.Sources` records how each source participated in the merge.

| Field | Semantics |
| --- | --- |
| `Name` | Source name. |
| `Kind` | Provider kind. |
| `Required` | Required flag from the source. |
| `Loaded` | True when provider data was loaded or merge inferred loaded values. |
| `Skipped` | True when an optional missing source was skipped. |
| `Message` | Human-readable context for skip or provider status. |

`Diagnostics.Overrides` records value replacement events.

| Field | Semantics |
| --- | --- |
| `Path` | Dot path that was replaced. |
| `Source` | Source that supplied the replacing value. |
| `PreviousSource` | Earlier source that supplied the previous value. |

Nested maps are deep-merged. Scalars and slices are replaced. Override diagnostics are recorded when a later source changes a previously owned path.

## Error Contract

All package-classified errors use `*config.Error`.

| Kind | Meaning |
| --- | --- |
| `ErrorKindContract` | The caller violated API or request rules. |
| `ErrorKindSource` | A provider could not read a declared source. |
| `ErrorKindParse` | A parser could not convert source payload data. |
| `ErrorKindMerge` | Parsed source values could not be merged. |
| `ErrorKindDecode` | Merged values could not be decoded into the target. |

Use `IsKind(err, kind)` or `errors.Is(err, &config.Error{Kind: kind})` instead of parsing error strings. Error strings are human-readable but should not be treated as machine contracts.

## Extension Contract

The package supports custom pipeline components through `LoaderConfig`:

- `Providers` registers additional `ProviderKind` handlers or overrides built-ins.
- `Parsers` registers additional `ParserKind` handlers or overrides built-ins.
- `Merger` replaces the default deep-replace merger.
- `Decoder` replaces the default mapstructure decoder.

Custom components must preserve the public semantics expected by the caller. The default loader will still validate request shape before invoking custom providers and parsers.

## Out of Contract

The following behavior is outside this package contract:

- runtime-specific setting defaults and validation;
- selecting source lists from process flags, deployment manifests, or service discovery;
- remote reads unless a caller registers and owns the provider;
- hot reload and watcher lifecycle;
- secret resolution;
- redacting returned values or diagnostics;
- guaranteeing Koanf or mapstructure behavior as a public runtime API.
