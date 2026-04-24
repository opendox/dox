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
 * @File    : loader.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-25
 */

package config

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Loader defines the configuration loading entrypoint contract.
type Loader interface {
	Load(ctx context.Context, req Request) (*Result, error)
}

// ProviderFunc adapts a function into a provider implementation.
type ProviderFunc func(ctx context.Context, source Source) (*Payload, error)

// Read calls f with the provided source.
func (f ProviderFunc) Read(ctx context.Context, source Source) (*Payload, error) {
	if f == nil {
		return nil, ContractError("provider", "provider function must not be nil")
	}
	return f(ctx, source)
}

// LoaderConfig customizes the default loader with additional pipeline components.
type LoaderConfig struct {
	Providers map[ProviderKind]Provider
	Parsers   map[ParserKind]Parser
	Merger    Merger
	Decoder   Decoder
}

// DefaultLoader orchestrates provider, parser, merge, and decode stages.
type DefaultLoader struct {
	providers map[ProviderKind]Provider
	parsers   map[ParserKind]Parser
	merger    Merger
	decoder   Decoder
}

// NewLoader creates a loader with built-in components plus caller overrides.
func NewLoader(config LoaderConfig) *DefaultLoader {
	providers := builtinProviders()
	for kind, provider := range config.Providers {
		providers[kind] = provider
	}

	parsers := builtinParsers()
	for kind, parser := range config.Parsers {
		parsers[kind] = parser
	}

	merger := config.Merger
	if merger == nil {
		merger = DeepReplaceMerger{}
	}
	decoder := config.Decoder
	if decoder == nil {
		decoder = MapstructureDecoder{}
	}

	return &DefaultLoader{
		providers: providers,
		parsers:   parsers,
		merger:    merger,
		decoder:   decoder,
	}
}

// NewDefaultLoader creates a loader with all built-in local components.
func NewDefaultLoader() *DefaultLoader {
	return NewLoader(LoaderConfig{})
}

// Load runs a request through the built-in local loader.
func Load(ctx context.Context, req Request) (*Result, error) {
	return NewDefaultLoader().Load(ctx, req)
}

// Load validates the request, loads sources, merges values, decodes the target, and returns diagnostics.
func (l *DefaultLoader) Load(ctx context.Context, req Request) (*Result, error) {
	if l == nil {
		return nil, ContractError("loader", "loader must not be nil")
	}
	if err := ValidateLoadRequest(ctx, req); err != nil {
		return nil, err
	}
	options, err := normalizeOptions(req.Options)
	if err != nil {
		return nil, err
	}

	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()
	}

	parsedSources := make([]ParsedSource, 0, len(req.Sources))
	providers := l.effectiveProviders()
	parsers := l.effectiveParsers()
	for index, source := range req.Sources {
		provider, ok := providers[source.Kind]
		if !ok || provider == nil {
			return nil, ContractError(sourceField(index, "kind"), "provider kind is not registered")
		}
		payload, err := provider.Read(ctx, source)
		if err != nil {
			return nil, err
		}
		if payload == nil {
			return nil, SourceError(sourceField(index, "kind"), "provider returned nil payload", nil)
		}
		payload.Source = source

		parser, ok := parsers[source.Parser]
		if !ok || parser == nil {
			return nil, ContractError(sourceField(index, "parser"), "parser kind is not registered")
		}
		values, err := parser.Parse(ctx, *payload)
		if err != nil {
			return nil, err
		}
		parsedSources = append(parsedSources, ParsedSourceFromPayload(*payload, values))
	}

	merger := l.merger
	if merger == nil {
		merger = DeepReplaceMerger{}
	}
	mergeResult, err := merger.Merge(ctx, parsedSources, options)
	if err != nil {
		return nil, err
	}

	decoder := l.decoder
	if decoder == nil {
		decoder = MapstructureDecoder{}
	}
	if err := decoder.Decode(ctx, mergeResult.Values, req.Target, options); err != nil {
		return nil, err
	}

	fingerprint, err := fingerprintValues(mergeResult.Values)
	if err != nil {
		return nil, MergeError("result.values", "fingerprint serialization failed", err)
	}

	return &Result{
		Runtime:     req.Runtime,
		Env:         req.Env,
		SourceNames: append([]string(nil), mergeResult.SourceNames...),
		Fingerprint: fingerprint,
		Diagnostics: mergeResult.Diagnostics,
	}, nil
}

func (l *DefaultLoader) effectiveProviders() map[ProviderKind]Provider {
	if len(l.providers) == 0 {
		return builtinProviders()
	}
	return l.providers
}

func (l *DefaultLoader) effectiveParsers() map[ParserKind]Parser {
	if len(l.parsers) == 0 {
		return builtinParsers()
	}
	return l.parsers
}

func builtinProviders() map[ProviderKind]Provider {
	return map[ProviderKind]Provider{
		ProviderKindFile: FileProvider{},
		ProviderKindEnv:  EnvProvider{},
	}
}

func builtinParsers() map[ParserKind]Parser {
	return map[ParserKind]Parser{
		ParserKindNone: NoneParser{},
		ParserKindYAML: YAMLParser{},
		ParserKindJSON: JSONParser{},
		ParserKindTOML: TOMLParser{},
	}
}

func fingerprintValues(values map[string]any) (string, error) {
	body, err := json.Marshal(cloneStructuredMap(values))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(body)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
