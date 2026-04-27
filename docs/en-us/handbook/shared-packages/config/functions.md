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

  @File    : docs/en-us/handbook/shared-packages/config/functions.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Config Functions and API

The exported API surface is grouped by loader entry points, pipeline extension points, helpers, and caller obligations.

## Load Entry Points

| API | Purpose | Use When |
| --- | --- | --- |
| `Load(ctx, req)` | Runs a request through a new default loader. | The runtime only needs built-in file/env providers and YAML/JSON/TOML/none parsers. |
| `NewDefaultLoader()` | Creates a default loader with built-in components. | The runtime wants a reusable loader instance without custom components. |
| `NewLoader(LoaderConfig)` | Creates a loader with built-ins plus overrides. | The runtime needs custom providers, parsers, merger, or decoder. |
| `(*DefaultLoader).Load(ctx, req)` | Runs one request through a configured loader. | The runtime already owns a loader instance. |

Use `Load` for simple local integrations. Use `NewLoader` when adding custom source kinds, such as a memory provider in tests or a remote provider owned by a runtime.

## Loader Configuration

| API | Purpose |
| --- | --- |
| `Loader` | Interface with `Load(ctx, req) (*Result, error)`. |
| `DefaultLoader` | Built-in implementation that orchestrates providers, parsers, merger, decoder, and fingerprinting. |
| `LoaderConfig` | Registers provider/parser maps and optional custom merger/decoder. |

`LoaderConfig.Providers` and `LoaderConfig.Parsers` are additive over built-ins. If a map contains a built-in key, the supplied component replaces that built-in component for the created loader.

## Provider API

| API | Purpose |
| --- | --- |
| `Provider` | Interface with `Read(ctx, source) (*Payload, error)`. |
| `ProviderFunc` | Adapter that lets a function implement `Provider`. |
| `FileProvider` | Built-in provider for local file reads. |
| `EnvProvider` | Built-in provider for environment variable prefix reads. |
| `ProviderKindFile` | Built-in `file` provider kind. |
| `ProviderKindEnv` | Built-in `env` provider kind. |
| `ProviderKindRemote` | Named `remote` kind without a default registered provider. |

Provider implementations should return `SourceError` for read failures and `ContractError` for provider/source mismatches. They should not merge, decode, or validate runtime-specific settings.

## Parser API

| API | Purpose |
| --- | --- |
| `Parser` | Interface with `Parse(ctx, payload) (map[string]any, error)`. |
| `ParserFunc` | Adapter that lets a function implement `Parser`. |
| `NoneParser` | Returns provider values for structured payloads. |
| `YAMLParser` | Parses YAML object payloads. |
| `JSONParser` | Parses JSON object payloads. |
| `TOMLParser` | Parses TOML object payloads. |
| `BuiltinParser(kind)` | Returns a built-in parser and a boolean. |
| `ParsePayload(ctx, payload)` | Parses a payload with its declared built-in parser. |

Parser implementations should return `ParseError` for malformed source data. They should return `ContractError` when the payload does not match the parser implementation.

## Merge API

| API | Purpose |
| --- | --- |
| `Merger` | Interface with `Merge(ctx, sources, options) (*MergeResult, error)`. |
| `DeepReplaceMerger` | Built-in merger. Deep-merges maps and replaces scalars/slices. |
| `MergeParsedSources(ctx, sources, options)` | Convenience helper using `DeepReplaceMerger`. |
| `ParsedSourceFromPayload(payload, values)` | Binds parsed values back to provider metadata. |

The built-in merger sorts by ascending priority. It records ordered source names, source diagnostics, and override diagnostics.

## Decode API

| API | Purpose |
| --- | --- |
| `Decoder` | Interface with `Decode(ctx, values, target, options) error`. |
| `MapstructureDecoder` | Built-in mapstructure-backed decoder. |
| `DecodeValues(ctx, values, target, options)` | Convenience helper using `MapstructureDecoder`. |
| `DecodeMergeResult(ctx, result, target, options)` | Decodes the values from a merge result. |

Decode helpers are useful in tests or in custom pipelines. The full loader should usually be preferred by runtimes because it preserves diagnostics and fingerprinting.

## Validation API

| API | Purpose |
| --- | --- |
| `ValidateLoadRequest(ctx, req)` | Validates request shape before loading. |

Use validation directly when a runtime wants to fail fast before constructing a loader or before adding runtime-specific checks. Validation does not prove custom providers or parsers are registered.

## Error API

| API | Purpose |
| --- | --- |
| `Error` | Typed package error with `Kind`, `Field`, `Reason`, and wrapped `Err`. |
| `ErrorKindContract` | API or request contract failure. |
| `ErrorKindSource` | Source read failure. |
| `ErrorKindParse` | Source parse failure. |
| `ErrorKindMerge` | Merge failure. |
| `ErrorKindDecode` | Decode failure. |
| `ContractError` | Creates a contract error. |
| `SourceError` | Creates a source error. |
| `ParseError` | Creates a parse error. |
| `MergeError` | Creates a merge error. |
| `DecodeError` | Creates a decode error. |
| `IsKind(err, kind)` | Reports whether an error contains a config error of the requested kind. |

Callers should branch on error kind rather than string contents.

## Data Types

| Type | Purpose |
| --- | --- |
| `Request` | One load operation. |
| `Options` | Load, merge, decode, and diagnostics behavior. |
| `Source` | One provider source descriptor. |
| `Payload` | Provider output before parser conversion. |
| `ParsedSource` | Parsed source values with source metadata. |
| `MergeResult` | Merged values, source names, and diagnostics. |
| `Result` | Successful load output after decode and fingerprinting. |
| `Diagnostics` | Source diagnostics and override diagnostics. |
| `SourceDiagnostic` | One source participation record. |
| `MergeOverride` | One override record. |

These types are plain Go structs. The package does not hide source shape behind global state.

## Common Recipes

### Default local load

Use `Load(ctx, req)` with file and env sources. This is the normal path for local runtime startup.

### Test-only memory source

Use `NewLoader(LoaderConfig{Providers: ..., Parsers: ...})` with custom kinds and `ProviderFunc` / `ParserFunc`. This keeps tests away from filesystem and process environment state.

### Pre-merge inspection

Use providers, parsers, and `MergeParsedSources` directly only when building a custom runtime pipeline. Most consumers should not need this.

### Decode-only tests

Use `DecodeValues` or `DecodeMergeResult` when testing target struct tags, duration conversion, or unknown-key behavior without source reads.

## Caller Obligations

Every consumer must own:

- setting target type definition;
- runtime-specific defaults;
- domain validation after decode;
- source list selection;
- environment prefix policy;
- secret and redaction policy;
- startup logging policy;
- any hot reload, remote configuration, or watcher lifecycle.

The shared package gives a deterministic loading primitive. It does not replace runtime configuration design.

## Related Pages

- [Shared config package manual](README.md)
- [Shared config contract](contract.md)
- [Shared config pipeline](pipeline.md)
