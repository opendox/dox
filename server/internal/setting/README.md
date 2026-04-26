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
  @Modified: 2026-04-26
-->

# Server Setting

`server/internal/setting` owns the concrete configuration aggregate for the Dox Web backend runtime.

Shared packages provide reusable configuration fragments. This package decides how the server runtime composes those fragments, which defaults are server-specific, and which validation rules are stricter than the shared fragment rules.

## Boundaries

- `server/internal/bootstrap` loads source snapshots from files and environment variables.
- `packages/shared/config` provides source loading, merging, and decoding primitives.
- `packages/shared/setting` defines reusable setting fragments.
- `server/internal/setting` defines the server runtime aggregate and group-level semantics.

Bootstrap should not own concrete HTTP, database, identity, logging, or security setting structs. It should continue to provide loaded values and diagnostics.

## File Convention

Use one file per configuration group:

- `setting.go` defines the root `Setting` aggregate.
- `identity.go` defines the identity group.
- Future `database.go`, `http.go`, `logging.go`, `security.go`, and similar files should define their own focused groups.

The root aggregate should compose groups instead of flattening every field:

```go
type Setting struct {
    Identity Identity `json:"identity" yaml:"identity" mapstructure:"identity"`
    Database Database `json:"database" yaml:"database" mapstructure:"database"`
}
```

Callers should pass narrow group settings to subsystems instead of passing the full root setting everywhere.

## Identity

The identity group composes shared identity fragments:

- `Organization`: ownership and governance metadata.
- `Application`: product or application family.
- `System`: Dox runtime identity.
- `Service`: logical service identity.
- `Deployment`: deployment environment and location.

The server package defaults `System.Runtime` to `server`. That default does not belong in `packages/shared/setting`, because scheduler, collector, and compute runtimes must own their own runtime identity.

`Deployment.Env` may be seeded from the bootstrap environment when the final server setting is created. If no seed or explicit value is provided, it falls back to the shared deployment default.
