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

  @File    : docs/en-us/handbook/shared-packages/setting/functions.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-27
  @Modified: 2026-04-27
-->

# Shared Setting Functions and API

The shared setting API surface includes enum constants, fragment methods, validation helpers, error types, and caller obligations for composing `packages/shared/setting` into runtime-owned aggregates.

## Constants

| API | Value | Purpose |
| --- | --- | --- |
| `DefaultOrganizationName` | `opendox` | Default organization identity. |
| `DefaultApplicationName` | `dox` | Default product or application family identity. |

## Runtime API

| API | Purpose |
| --- | --- |
| `type Runtime string` | Dox runtime identity value. |
| `RuntimeServer` | Web backend runtime. |
| `RuntimeScheduler` | Scheduling runtime. |
| `RuntimeCollector` | Collection runtime. |
| `RuntimeCompute` | Computation runtime. |
| `(Runtime).IsValid() bool` | Returns true for supported runtime values. |

Runtime validation is shared, but runtime selection is not. For example, the server package may require `RuntimeServer`, while a scheduler package should require `RuntimeScheduler`.

## Env API

| API | Purpose |
| --- | --- |
| `type Env string` | Deployment environment value. |
| `EnvDev` | Development environment. |
| `EnvTest` | Test environment. |
| `EnvStaging` | Staging environment. |
| `EnvProd` | Production environment. |
| `(Env).IsValid() bool` | Returns true for supported environment values. |

## Fragment API

| Fragment | Default Method | Validate Method |
| --- | --- | --- |
| `Organization` | `(*Organization).Default() error` | `(Organization).Validate() error` |
| `Application` | `(*Application).Default() error` | `(Application).Validate() error` |
| `System` | `(*System).Default() error` | `(System).Validate() error` |
| `Service` | `(*Service).Default(application, system) error` | `(Service).Validate() error` |
| `Deployment` | `(*Deployment).Default() error` | `(Deployment).Validate() error` |

Default methods return an error when called on a nil receiver. Validate methods call the package-level `Validate` helper.

## Default Method Behavior

| Method | Behavior |
| --- | --- |
| `Organization.Default` | Sets empty `Name` to `DefaultOrganizationName`. |
| `Application.Default` | Sets empty `Name` to `DefaultApplicationName`. |
| `System.Default` | Checks receiver only. It does not set `Runtime`. |
| `Service.Default` | Sets empty `Namespace` from `application.Name`; sets empty `Name` from `system.Runtime` when runtime is known. |
| `Deployment.Default` | Sets empty `Env` to `EnvDev`. |

<details>
<summary>Example: default order in a runtime aggregate</summary>

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

The runtime package owns the line that sets `RuntimeServer`.

</details>

## Validation API

| API | Purpose |
| --- | --- |
| `Validate(value any) error` | Validates a struct with registered Dox validation tags. |
| `FieldError` | Machine-readable field failure with `Field` and `Rule`. |
| `ValidationError` | List of field failures. |
| `(*ValidationError).Error() string` | Compact log-friendly error message. |

The validation engine is initialized once and registers:

- `dox_kebab`;
- `dox_identifier`;
- `dox_runtime`;
- `dox_env`.

Field names are derived from `mapstructure` tags first. That keeps validation output aligned with decoded configuration keys.

## Caller Recipes

### Compose a runtime identity group

Define a runtime-owned group that embeds or contains shared fragments. Do not expose `packages/shared/setting` as the entire runtime setting aggregate.

### Add runtime-specific validation

Join shared validation with runtime-specific validation:

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

### Keep runtime bootstrap separate

Bootstrap may provide seed values such as env, instance ID, or region. The shared package should not read flags, environment variables, files, or process metadata directly.

## Caller Obligations

Every runtime consumer must own:

- the root `Setting` aggregate;
- group composition;
- runtime identity selection;
- runtime-specific defaults;
- validation that depends on combinations across fragments;
- bootstrap-derived seed values;
- documentation for stricter runtime behavior.

## API Stability Notes

- Shared fragments are intended to be reused across Dox runtimes.
- Validation tags are part of the public package behavior.
- The third-party validator implementation is not exposed as the caller-facing error type.
- Adding a new runtime or environment value changes shared validation semantics and should be treated as a contract change.

## Related Pages

- [Shared setting package manual](README.md)
- [Shared setting contract](contract.md)
- [Shared setting model](model.md)
