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
 * @File    : parser.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package config

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	"go.yaml.in/yaml/v3"
)

// ParserFunc adapts a function into a parser implementation.
type ParserFunc func(ctx context.Context, payload Payload) (map[string]any, error)

// Parse calls f with the provided payload.
func (f ParserFunc) Parse(ctx context.Context, payload Payload) (map[string]any, error) {
	return f(ctx, payload)
}

// NoneParser returns structured values already produced by a provider.
type NoneParser struct{}

// Parse returns provider values for sources that do not require raw parsing.
func (p NoneParser) Parse(ctx context.Context, payload Payload) (map[string]any, error) {
	if err := validateParserPayload(ctx, payload, ParserKindNone); err != nil {
		return nil, err
	}
	return cloneMap(payload.Values), nil
}

// JSONParser parses JSON object payloads.
type JSONParser struct{}

// Parse converts raw JSON bytes into structured values.
func (p JSONParser) Parse(ctx context.Context, payload Payload) (map[string]any, error) {
	if err := validateParserPayload(ctx, payload, ParserKindJSON); err != nil {
		return nil, err
	}
	var root any
	decoder := json.NewDecoder(bytes.NewReader(payload.Raw))
	decoder.UseNumber()
	if err := decoder.Decode(&root); err != nil {
		return nil, ParseError(payloadField(payload, "raw"), "json parser failed", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return nil, ParseError(payloadField(payload, "raw"), "json parser found trailing data", nil)
		}
		return nil, ParseError(payloadField(payload, "raw"), "json parser failed", err)
	}
	return parsedObject(payload, root)
}

// YAMLParser parses YAML object payloads.
type YAMLParser struct{}

// Parse converts raw YAML bytes into structured values.
func (p YAMLParser) Parse(ctx context.Context, payload Payload) (map[string]any, error) {
	if err := validateParserPayload(ctx, payload, ParserKindYAML); err != nil {
		return nil, err
	}
	var root any
	if err := yaml.Unmarshal(payload.Raw, &root); err != nil {
		return nil, ParseError(payloadField(payload, "raw"), "yaml parser failed", err)
	}
	return parsedObject(payload, root)
}

// TOMLParser parses TOML object payloads.
type TOMLParser struct{}

// Parse converts raw TOML bytes into structured values.
func (p TOMLParser) Parse(ctx context.Context, payload Payload) (map[string]any, error) {
	if err := validateParserPayload(ctx, payload, ParserKindTOML); err != nil {
		return nil, err
	}
	var root map[string]any
	if err := toml.Unmarshal(payload.Raw, &root); err != nil {
		return nil, ParseError(payloadField(payload, "raw"), "toml parser failed", err)
	}
	return parsedObject(payload, root)
}

// BuiltinParser returns a parser for the built-in parser kind.
func BuiltinParser(kind ParserKind) (Parser, bool) {
	switch kind {
	case ParserKindNone:
		return NoneParser{}, true
	case ParserKindJSON:
		return JSONParser{}, true
	case ParserKindYAML:
		return YAMLParser{}, true
	case ParserKindTOML:
		return TOMLParser{}, true
	default:
		return nil, false
	}
}

// ParsePayload parses a payload with its declared built-in parser.
func ParsePayload(ctx context.Context, payload Payload) (map[string]any, error) {
	parser, ok := BuiltinParser(payload.Source.Parser)
	if !ok {
		return nil, ContractError(payloadField(payload, "parser"), "parser kind is not supported")
	}
	return parser.Parse(ctx, payload)
}

func validateParserPayload(ctx context.Context, payload Payload, expected ParserKind) error {
	if ctx == nil {
		return ContractError("ctx", "context must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return ParseError("ctx", "context is done", err)
	}
	if err := validateProviderSource(payload.Source); err != nil {
		return err
	}
	if payload.Source.Parser != expected {
		return ContractError(payloadField(payload, "parser"), "parser kind does not match parser implementation")
	}
	if payload.Diagnostic.Skipped {
		return nil
	}
	if expected == ParserKindNone {
		return nil
	}
	return nil
}

func parsedObject(payload Payload, root any) (map[string]any, error) {
	if payload.Diagnostic.Skipped || root == nil {
		return map[string]any{}, nil
	}
	normalized, err := normalizeParsedValue(root)
	if err != nil {
		return nil, ParseError(payloadField(payload, "raw"), "parsed value contains unsupported object keys", err)
	}
	values, ok := normalized.(map[string]any)
	if !ok {
		return nil, ParseError(payloadField(payload, "raw"), "parsed payload must be an object", nil)
	}
	return values, nil
}

func normalizeParsedValue(value any) (any, error) {
	switch typed := value.(type) {
	case map[string]any:
		return normalizeStringMap(typed)
	case map[any]any:
		values := make(map[string]any, len(typed))
		for key, entry := range typed {
			name, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("non-string key %v", key)
			}
			normalized, err := normalizeParsedValue(entry)
			if err != nil {
				return nil, err
			}
			values[name] = normalized
		}
		return values, nil
	case []any:
		values := make([]any, len(typed))
		for index, entry := range typed {
			normalized, err := normalizeParsedValue(entry)
			if err != nil {
				return nil, err
			}
			values[index] = normalized
		}
		return values, nil
	default:
		return value, nil
	}
}

func normalizeStringMap(input map[string]any) (map[string]any, error) {
	values := make(map[string]any, len(input))
	for key, entry := range input {
		normalized, err := normalizeParsedValue(entry)
		if err != nil {
			return nil, err
		}
		values[key] = normalized
	}
	return values, nil
}

func cloneMap(input map[string]any) map[string]any {
	if len(input) == 0 {
		return map[string]any{}
	}
	values := make(map[string]any, len(input))
	for key, value := range input {
		values[key] = value
	}
	return values
}

func payloadField(payload Payload, field string) string {
	field = strings.TrimSpace(field)
	if payload.Source.Name == "" {
		return "payload." + field
	}
	return "source." + payload.Source.Name + "." + field
}
