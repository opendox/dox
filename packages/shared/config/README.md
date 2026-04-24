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

  @File    : packages/shared/config/README.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-24
  @Modified: 2026-04-24
-->

# Shared Config Loader Contract

`packages/shared/config` defines the shared runtime configuration loading contract for Dox backend runtimes.

This package is a loader SDK. Callers must pass an explicit request, target, source list, and options. The package validates API usage and fails fast when the request contract is invalid.

## Boundary

The config package validates loader contract rules:

- request shape
- target pointer requirements
- source descriptors
- provider and parser naming
- built-in provider and parser compatibility
- option consistency
- future error categories

The config package does not validate runtime-specific setting values. That belongs to each runtime setting package, such as `server/internal/setting`.

## Current Scope

This milestone only defines the contract. It does not implement file loading, environment loading, parsing, merging, decoding, or remote configuration providers.
