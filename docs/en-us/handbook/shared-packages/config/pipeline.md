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

  @File    : docs/en-us/handbook/shared-packages/config/pipeline.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Config Pipeline

The default loader runs one request through a fixed local pipeline:

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

Every stage receives the request context. If `Options.Timeout` is positive, the loader wraps the context with that timeout before provider reads begin.

## Stage 1: Validation

`ValidateLoadRequest` runs before source reads. It checks:

- non-nil context;
- valid `Runtime` and `Env`;
- non-nil pointer target to a struct or map;
- supported option values;
- source presence unless `AllowEmptySources` is true;
- unique source names;
- unique source priorities;
- valid source names, provider kinds, parser kinds, locations, and priorities;
- built-in file/env provider and parser compatibility.

Validation accepts custom provider and parser kind names when they follow naming rules. Load will still fail later with `ErrorKindContract` if no matching registered component exists.

## Stage 2: Provider Read

The provider returns a `Payload`. A payload may contain `Raw` bytes, structured `Values`, provider `Metadata`, and a `SourceDiagnostic`.

### File Provider

`FileProvider` reads `Source.Location` from the local filesystem.

- Required missing files fail with `ErrorKindSource`.
- Optional missing files return a skipped payload.
- Directories fail with `ErrorKindSource`.
- Successful reads populate `Raw`, file path metadata, and a loaded diagnostic.

The file provider does not parse bytes. It only reads the source.

### Environment Provider

`EnvProvider` reads environment entries from `Lookup` when provided, otherwise from `os.Environ`.

It filters entries by `Source.Location` as a prefix. Matching keys are normalized by:

- removing the source prefix;
- trimming leading and trailing underscores;
- lowercasing;
- replacing double underscores with dots;
- replacing remaining underscores with dots.

Examples:

| Environment key | Prefix | Normalized key |
| --- | --- | --- |
| `DOX_SERVER_APP_NAME` | `DOX_SERVER_` | `app.name` |
| `DOX_SERVER_HTTP__ADDRESS` | `DOX_SERVER_` | `http.address` |

When no matching variables exist, a required env source fails with `ErrorKindSource`; an optional env source is skipped with diagnostics.

## Stage 3: Parser

The parser converts one provider payload into `map[string]any`.

| Parser | Input | Behavior |
| --- | --- | --- |
| `NoneParser` | `Payload.Values` | Returns a shallow copy. Used for env sources. |
| `YAMLParser` | `Payload.Raw` | Parses YAML and requires an object root. |
| `JSONParser` | `Payload.Raw` | Parses JSON, rejects trailing data, preserves numbers as `json.Number`, and requires an object root. |
| `TOMLParser` | `Payload.Raw` | Parses TOML and requires an object root. |

Skipped optional payloads parse as an empty map. Malformed payloads fail with `ErrorKindParse`. Parser implementation mismatches fail with `ErrorKindContract`.

## Stage 4: Merge

`DeepReplaceMerger` sorts parsed sources by ascending `Source.Priority`. Lower-priority numbers merge first. Higher-priority sources override previous values.

Merge behavior:

- nested maps are deep-merged;
- scalars replace previous scalars;
- slices replace previous slices;
- skipped sources appear in `SourceNames` and source diagnostics but do not change values;
- duplicate source names or priorities fail with `ErrorKindContract`;
- environment sources using `ParserKindNone` expand dotted keys into nested maps before merging.

Example:

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

The result records override diagnostics for `app.name` and `http.timeout`, with `env` replacing `base`.

Conflicting environment dotted keys fail with `ErrorKindMerge`. For example, one env source cannot provide both `app=dox` and `app.name=dox` because the same path would need to be both scalar and map.

## Stage 5: Decode

`MapstructureDecoder` decodes merged values into the caller target.

Decode behavior:

- target must be a non-nil pointer to a struct or map;
- struct tags use `mapstructure`;
- weak type conversion is enabled;
- string-to-`time.Duration` conversion is enabled;
- target fields are zeroed before decode;
- unknown keys are rejected by default;
- `UnknownKeyPolicyAllow` permits unknown keys.

Decode failures return `ErrorKindDecode`. Invalid targets return `ErrorKindContract`.

## Stage 6: Fingerprint and Result

After decode succeeds, the loader computes a stable `sha256:` fingerprint from the merged structured value map. The result includes:

- request runtime;
- request environment;
- ordered source names;
- fingerprint;
- source diagnostics;
- merge override diagnostics.

The fingerprint represents merged values, not the decoded target. If decode fails, no result is returned.

## Operational Notes

Use source names that reveal intent, such as `base`, `local`, `secrets`, or `env`. Keep priorities sparse, such as `10`, `20`, and `100`, so new override layers can be inserted later.

For deterministic operations:

- avoid duplicate priorities;
- keep file formats object-rooted;
- keep env variable prefixes runtime-specific;
- log source names and fingerprint after successful startup;
- log diagnostics when investigating override behavior.

Do not log raw merged values unless the caller has already applied its own redaction policy.

## Related Pages

- [Shared config package manual](README.md)
- [Shared config contract](contract.md)
- [Shared config functions and API](functions.md)
