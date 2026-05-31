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

  @File    : server/internal/setting/README.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-26
  @Modified: 2026-05-31
-->

# Server Setting

`server/internal/setting` owns the concrete configuration aggregate for the Dox Web backend runtime.

Shared packages provide reusable configuration fragments. This package decides how the Web backend server runtime composes those fragments, which defaults are server-specific, and which validation rules are stricter than the shared fragment rules. Scheduler, collector, and computation runtimes should own separate setting aggregates when they are introduced.

## Boundaries

- `server/internal/bootstrap` loads source snapshots from files and environment variables, assembles typed server settings, and boots runtime resources from validated snapshots.
- `packages/shared/config` provides source loading, merging, and decoding primitives.
- `packages/shared/setting` defines reusable setting fragments.
- `packages/shared/logging` defines the shared logging model and runtime helper configuration.
- `server/internal/setting` defines the server runtime aggregate and group-level semantics.

Bootstrap should not own concrete HTTP, database, identity, logging, or security setting structs. It coordinates source loading, decode, defaulting, validation, and diagnostics preservation, while this package owns the setting groups and their semantics.

The expected assembly order is:

1. `server/internal/bootstrap` builds config sources from startup options.
2. `packages/shared/config` loads and merges raw values into a `map[string]any` snapshot.
3. `server/internal/bootstrap` decodes the snapshot into `Setting` with unknown keys rejected.
4. `server/internal/setting` applies defaults and validates group semantics.
5. `server/internal/bootstrap.Boot` receives the validated `SettingSnapshot` and constructs runtime resources.

## File Convention

Use one file per configuration group:

- `setting.go` defines the root `Setting` aggregate.
- `default.go` defines the defaulting contract shared by setting groups.
- `validate.go` defines the validation contract shared by setting groups.
- `identity.go` defines the identity group.
- `logging.go` defines the server logging configuration group backed by shared logging config.
- Future `database.go`, `http.go`, `security.go`, and similar files should define their own focused groups.

The root aggregate should compose groups instead of flattening every field:

```go
type Setting struct {
    Identity Identity `json:"identity" yaml:"identity" mapstructure:"identity"`
    Logging  Logging  `json:"logging" yaml:"logging" mapstructure:"logging"`
    Database Database `json:"database" yaml:"database" mapstructure:"database"`
}
```

Callers should pass narrow group settings to subsystems instead of passing the full root setting everywhere.

## Defaults And Validation

Setting groups that can fill stable defaults should implement:

```go
type Defaultable interface {
    Default() error
}
```

Setting groups that can validate their final values should implement:

```go
type Validatable interface {
    Validate() error
}
```

`DefaultGroups` applies defaults in order and stops on the first error. Use explicit root-level calls when a group needs bootstrap-derived seed values, such as identity using the bootstrap env.

`ValidateGroups` validates groups and joins reported errors so startup can report every invalid group found in one pass.

Future server-owned HTTP, database, Redis, RabbitMQ producer, security, plugin host, and similar setting groups should plug into these default and validation contracts before server runtime bootstrap consumes them. Scheduler, collector, and computation runtime settings should not be added to this package only because they share deployment infrastructure. The contracts do not define those concrete groups, construct clients, open network listeners, resolve secrets, or start runtime resources.

## Identity

The identity group composes shared identity fragments:

- `Organization`: ownership and governance metadata.
- `Application`: product or application family.
- `System`: Dox runtime identity.
- `Service`: logical service identity.
- `Deployment`: deployment environment and location.

The server package defaults `System.Runtime` to `server`. That default does not belong in `packages/shared/setting`, because scheduler, collector, and compute runtimes must own their own runtime identity.

`Deployment.Env` may be seeded from the bootstrap environment when the final server setting is created. If no seed or explicit value is provided, it falls back to the shared deployment default.

## Logging

The logging group is backed by `packages/shared/logging.Config`.

Server settings own loading, defaulting, and validation of this group. Runtime bootstrap constructs zap cores, the Dox logger facade, and the OpenTelemetry SDK base from the validated config.

This package must not open logging sinks, create log files, install OpenTelemetry globals, or wire HTTP/server modules to logging.
