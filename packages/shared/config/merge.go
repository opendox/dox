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
 * @File    : merge.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package config

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/knadh/koanf/providers/confmap"
	koanf "github.com/knadh/koanf/v2"
)

// DeepReplaceMerger deep-merges maps and replaces scalar or slice values.
type DeepReplaceMerger struct{}

// Merge combines parsed sources by ascending priority.
func (m DeepReplaceMerger) Merge(ctx context.Context, sources []ParsedSource, options Options) (*MergeResult, error) {
	if ctx == nil {
		return nil, ContractError("ctx", "context must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, MergeError("ctx", "context is done", err)
	}
	if err := validateOptions(options); err != nil {
		return nil, err
	}

	ordered, err := orderParsedSources(sources)
	if err != nil {
		return nil, err
	}

	engine := koanf.New(".")
	result := &MergeResult{}
	owners := map[string]string{}
	for _, source := range ordered {
		if err := ctx.Err(); err != nil {
			return nil, MergeError("ctx", "context is done", err)
		}
		result.SourceNames = append(result.SourceNames, source.Source.Name)
		result.Diagnostics.Sources = append(result.Diagnostics.Sources, mergeSourceDiagnostic(source))
		if source.Diagnostic.Skipped {
			continue
		}

		values, err := mergeSourceValues(source)
		if err != nil {
			return nil, err
		}
		incoming, err := loadSourceKoanf(source, values)
		if err != nil {
			return nil, err
		}
		incomingValues := incoming.Raw()
		overrides := collectMergeOverrides(engine.Raw(), owners, nil, incomingValues, source.Source.Name)
		if err := engine.Merge(incoming); err != nil {
			return nil, MergeError(sourceMergeField(source.Source, "values"), "koanf merge failed", err)
		}
		setValueOwners(owners, nil, incomingValues, source.Source.Name)
		result.Diagnostics.Overrides = append(result.Diagnostics.Overrides, overrides...)
	}
	result.Values = cloneStructuredMap(engine.Raw())
	return result, nil
}

// MergeParsedSources merges parsed source values with the built-in deep replace merger.
func MergeParsedSources(ctx context.Context, sources []ParsedSource, options Options) (*MergeResult, error) {
	return DeepReplaceMerger{}.Merge(ctx, sources, options)
}

// ParsedSourceFromPayload binds parsed values back to the provider payload source metadata.
func ParsedSourceFromPayload(payload Payload, values map[string]any) ParsedSource {
	return ParsedSource{
		Source:     payload.Source,
		Values:     cloneStructuredMap(values),
		Diagnostic: payload.Diagnostic,
	}
}

func orderParsedSources(sources []ParsedSource) ([]ParsedSource, error) {
	ordered := make([]ParsedSource, len(sources))
	copy(ordered, sources)
	sort.SliceStable(ordered, func(i int, j int) bool {
		return ordered[i].Source.Priority < ordered[j].Source.Priority
	})

	seenNames := map[string]struct{}{}
	seenPriorities := map[int]struct{}{}
	for index, source := range ordered {
		if err := validateProviderSource(source.Source); err != nil {
			return nil, err
		}
		if _, exists := seenNames[source.Source.Name]; exists {
			return nil, ContractError(parsedSourceField(index, "name"), "source name must be unique")
		}
		seenNames[source.Source.Name] = struct{}{}
		if _, exists := seenPriorities[source.Source.Priority]; exists {
			return nil, ContractError(parsedSourceField(index, "priority"), "source priority must be unique")
		}
		seenPriorities[source.Source.Priority] = struct{}{}
	}
	return ordered, nil
}

func mergeSourceValues(source ParsedSource) (map[string]any, error) {
	values := cloneStructuredMap(source.Values)
	if source.Source.Kind != ProviderKindEnv || source.Source.Parser != ParserKindNone {
		return values, nil
	}
	expanded, err := expandDottedValues(values)
	if err != nil {
		return nil, MergeError(sourceMergeField(source.Source, "values"), "environment source values are not expandable", err)
	}
	return expanded, nil
}

func loadSourceKoanf(source ParsedSource, values map[string]any) (*koanf.Koanf, error) {
	sourceConfig := koanf.New(".")
	if err := sourceConfig.Load(confmap.Provider(values, ""), nil); err != nil {
		return nil, MergeError(sourceMergeField(source.Source, "values"), "koanf source load failed", err)
	}
	return sourceConfig, nil
}

func collectMergeOverrides(current map[string]any, owners map[string]string, path []string, incoming map[string]any, sourceName string) []MergeOverride {
	var overrides []MergeOverride
	keys := sortedMapKeys(incoming)
	for _, key := range keys {
		nextPath := appendPath(path, key)
		incomingValue := incoming[key]
		existingValue, exists := current[key]

		incomingMap, incomingIsMap := incomingValue.(map[string]any)
		existingMap, existingIsMap := existingValue.(map[string]any)
		if exists && incomingIsMap && existingIsMap {
			nestedOverrides := collectMergeOverrides(existingMap, owners, nextPath, incomingMap, sourceName)
			overrides = append(overrides, nestedOverrides...)
			continue
		}

		if exists && !reflect.DeepEqual(existingValue, incomingValue) {
			previousSource := ownerForPath(owners, nextPath)
			if previousSource != "" && previousSource != sourceName {
				overrides = append(overrides, MergeOverride{
					Path:           strings.Join(nextPath, "."),
					Source:         sourceName,
					PreviousSource: previousSource,
				})
			}
		}
	}
	return overrides
}

func expandDottedValues(values map[string]any) (map[string]any, error) {
	expanded := map[string]any{}
	keys := sortedMapKeys(values)
	for _, key := range keys {
		parts := strings.Split(key, ".")
		for _, part := range parts {
			if strings.TrimSpace(part) == "" {
				return nil, fmt.Errorf("empty path segment in %q", key)
			}
		}
		if err := setDottedValue(expanded, parts, cloneStructuredValue(values[key])); err != nil {
			return nil, err
		}
	}
	return expanded, nil
}

func setDottedValue(values map[string]any, path []string, value any) error {
	if len(path) == 1 {
		if existing, exists := values[path[0]]; exists && !reflect.DeepEqual(existing, value) {
			return fmt.Errorf("conflicting value for %q", strings.Join(path, "."))
		}
		values[path[0]] = value
		return nil
	}
	head := path[0]
	next, exists := values[head]
	if !exists {
		child := map[string]any{}
		values[head] = child
		return setDottedValue(child, path[1:], value)
	}
	child, ok := next.(map[string]any)
	if !ok {
		return fmt.Errorf("conflicting scalar at %q", head)
	}
	return setDottedValue(child, path[1:], value)
}

func mergeSourceDiagnostic(source ParsedSource) SourceDiagnostic {
	diagnostic := source.Diagnostic
	if diagnostic.Name == "" {
		diagnostic.Name = source.Source.Name
	}
	if diagnostic.Kind == "" {
		diagnostic.Kind = source.Source.Kind
	}
	if !diagnostic.Required {
		diagnostic.Required = source.Source.Required
	}
	if !diagnostic.Loaded && !diagnostic.Skipped && len(source.Values) > 0 {
		diagnostic.Loaded = true
	}
	return diagnostic
}

func ownerForPath(owners map[string]string, path []string) string {
	for len(path) > 0 {
		if owner := owners[strings.Join(path, ".")]; owner != "" {
			return owner
		}
		path = path[:len(path)-1]
	}
	return ""
}

func setValueOwners(owners map[string]string, path []string, value any, sourceName string) {
	joined := strings.Join(path, ".")
	owners[joined] = sourceName
	if values, ok := value.(map[string]any); ok {
		for _, key := range sortedMapKeys(values) {
			setValueOwners(owners, appendPath(path, key), values[key], sourceName)
		}
	}
}

func cloneStructuredMap(input map[string]any) map[string]any {
	if len(input) == 0 {
		return map[string]any{}
	}
	values := make(map[string]any, len(input))
	for key, value := range input {
		values[key] = cloneStructuredValue(value)
	}
	return values
}

func cloneStructuredValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return cloneStructuredMap(typed)
	case []any:
		values := make([]any, len(typed))
		for index, entry := range typed {
			values[index] = cloneStructuredValue(entry)
		}
		return values
	default:
		return value
	}
}

func sortedMapKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func appendPath(path []string, key string) []string {
	next := make([]string, 0, len(path)+1)
	next = append(next, path...)
	next = append(next, key)
	return next
}

func parsedSourceField(index int, field string) string {
	return "parsed_sources[" + fmt.Sprint(index) + "].source." + field
}

func sourceMergeField(source Source, field string) string {
	if source.Name == "" {
		return "source." + field
	}
	return "source." + source.Name + "." + field
}
