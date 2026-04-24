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

The current package defines the loader contract, local provider primitives, and local parser primitives. It includes local file providers, environment variable providers, and YAML, JSON, and TOML parsers.

It does not implement multi-source merge behavior, decode behavior, server runtime integration, or remote configuration providers.

## Provider Scope

Local providers are responsible for reading source payloads and preserving source metadata for later pipeline stages.

Providers may:

- read required or optional local files;
- read environment variables by explicit prefix;
- return raw bytes or structured key-value payloads;
- attach source diagnostics for later operational review.

Providers must not:

- merge multiple sources;
- decode values into runtime setting structs;
- validate runtime-specific setting values;
- hide required source failures.

## Parser Scope

Parsers are responsible for converting one provider payload into structured values for later pipeline stages.

Parsers may:

- parse YAML, JSON, and TOML file payloads;
- pass through structured environment provider values through the none parser;
- report malformed payloads with typed parse errors;
- preserve generic map-based values for later merge and decode stages.

Parsers must not:

- merge multiple parsed payloads;
- decode values into runtime setting structs;
- validate runtime-specific setting values;
- hide malformed required source payloads.

Merge behavior, decode behavior, server runtime integration, and remote configuration providers are separate follow-up milestones.
