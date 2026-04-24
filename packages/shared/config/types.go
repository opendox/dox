/**
 * dox
 * Copyright (C) 2026  OpenDox
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * @File    : types.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-25
 */

package config

import (
	"context"
	"time"
)

// ProviderKind identifies the provider implementation that should read a source.
type ProviderKind string

const (
	// ProviderKindFile identifies the local file provider.
	ProviderKindFile ProviderKind = "file"
	// ProviderKindEnv identifies the environment variable provider.
	ProviderKindEnv ProviderKind = "env"
	// ProviderKindRemote identifies a future remote configuration provider.
	ProviderKindRemote ProviderKind = "remote"
)

// ParserKind identifies the parser implementation for source payloads.
type ParserKind string

const (
	// ParserKindNone means the provider already returns structured values.
	ParserKindNone ParserKind = "none"
	// ParserKindYAML identifies the YAML parser.
	ParserKindYAML ParserKind = "yaml"
	// ParserKindJSON identifies the JSON parser.
	ParserKindJSON ParserKind = "json"
	// ParserKindTOML identifies the TOML parser.
	ParserKindTOML ParserKind = "toml"
)

// MergeStrategy defines how values from a later source override earlier values.
type MergeStrategy string

const (
	// MergeStrategyDeepReplace deep-merges maps and replaces scalar and slice values.
	MergeStrategyDeepReplace MergeStrategy = "deep_replace"
)

// UnknownKeyPolicy defines how the decoder should handle unknown fields.
type UnknownKeyPolicy string

const (
	// UnknownKeyPolicyAllow permits keys that are not represented by the target type.
	UnknownKeyPolicyAllow UnknownKeyPolicy = "allow"
	// UnknownKeyPolicyReject rejects keys that are not represented by the target type.
	UnknownKeyPolicyReject UnknownKeyPolicy = "reject"
)

// Source describes one configuration source for provider reads.
type Source struct {
	Name     string
	Kind     ProviderKind
	Parser   ParserKind
	Location string
	Required bool
	Priority int
	Options  map[string]string
}

// Payload is the raw provider output that later pipeline stages parse, merge, and decode.
type Payload struct {
	Source     Source
	Raw        []byte
	Values     map[string]any
	Metadata   map[string]string
	Diagnostic SourceDiagnostic
}

// Provider reads one configuration source and returns a payload for later stages.
type Provider interface {
	Read(ctx context.Context, source Source) (*Payload, error)
}

// Parser converts provider bytes into structured values for later merge stages.
type Parser interface {
	Parse(ctx context.Context, payload Payload) (map[string]any, error)
}

// ParsedSource is one provider source after parser output is available.
type ParsedSource struct {
	Source     Source
	Values     map[string]any
	Diagnostic SourceDiagnostic
}

// Merger combines parsed sources into one structured value map.
type Merger interface {
	Merge(ctx context.Context, sources []ParsedSource, options Options) (*MergeResult, error)
}

// Decoder copies merged configuration values into a caller-provided target.
type Decoder interface {
	Decode(ctx context.Context, values map[string]any, target any, options Options) error
}

// MergeResult describes the output of a parsed source merge operation.
type MergeResult struct {
	Values      map[string]any
	SourceNames []string
	Diagnostics Diagnostics
}

// Request describes one explicit configuration loading operation.
type Request struct {
	Runtime string
	Env     string
	Target  any
	Sources []Source
	Options Options
}

// Options controls loading, merge, decode, and diagnostics behavior.
type Options struct {
	AllowEmptySources bool
	MergeStrategy     MergeStrategy
	UnknownKeyPolicy  UnknownKeyPolicy
	Timeout           time.Duration
	RedactKeys        []string
}

// Result describes the output shape of a completed load.
type Result struct {
	Runtime     string
	Env         string
	SourceNames []string
	Fingerprint string
	Diagnostics Diagnostics
}

// Diagnostics records source and merge metadata for operational review.
type Diagnostics struct {
	Sources   []SourceDiagnostic
	Overrides []MergeOverride
}

// SourceDiagnostic describes how one source participated in a load operation.
type SourceDiagnostic struct {
	Name     string
	Kind     ProviderKind
	Required bool
	Loaded   bool
	Skipped  bool
	Message  string
}

// MergeOverride records which source replaced a previously merged value.
type MergeOverride struct {
	Path           string
	Source         string
	PreviousSource string
}
