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
 * @File    : validate.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package config

import (
	"context"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var runtimeNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
var envNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
var sourceNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-_.]*$`)
var kindNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-_.]*$`)

// ValidateLoadRequest verifies that a loader request fails fast on API misuse.
func ValidateLoadRequest(ctx context.Context, req Request) error {
	if ctx == nil {
		return ContractError("ctx", "context must not be nil")
	}
	if err := validateRuntime(req.Runtime); err != nil {
		return err
	}
	if err := validateEnv(req.Env); err != nil {
		return err
	}
	if err := validateTarget(req.Target); err != nil {
		return err
	}
	if err := validateOptions(req.Options); err != nil {
		return err
	}
	if !req.Options.AllowEmptySources && len(req.Sources) == 0 {
		return ContractError("sources", "at least one source is required")
	}

	seenNames := map[string]struct{}{}
	seenPriorities := map[int]struct{}{}
	for index, source := range req.Sources {
		if err := validateSource(index, source); err != nil {
			return err
		}
		sourceName := strings.TrimSpace(source.Name)
		if _, exists := seenNames[sourceName]; exists {
			return ContractError(sourceField(index, "name"), "source name must be unique")
		}
		seenNames[sourceName] = struct{}{}
		if _, exists := seenPriorities[source.Priority]; exists {
			return ContractError(sourceField(index, "priority"), "source priority must be unique")
		}
		seenPriorities[source.Priority] = struct{}{}
	}
	return nil
}

func validateRuntime(runtime string) error {
	runtime = strings.TrimSpace(runtime)
	if runtime == "" {
		return ContractError("runtime", "runtime is required")
	}
	if !runtimeNamePattern.MatchString(runtime) {
		return ContractError("runtime", "runtime must use lowercase letters, digits, or hyphens and start with a letter")
	}
	return nil
}

func validateEnv(env string) error {
	env = strings.TrimSpace(env)
	if env == "" {
		return ContractError("env", "environment is required")
	}
	if !envNamePattern.MatchString(env) {
		return ContractError("env", "environment must use lowercase letters, digits, or hyphens and start with a letter")
	}
	return nil
}

func validateTarget(target any) error {
	if target == nil {
		return ContractError("target", "target must not be nil")
	}
	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return ContractError("target", "target must be a non-nil pointer")
	}
	targetKind := value.Elem().Kind()
	switch targetKind {
	case reflect.Struct, reflect.Map:
	default:
		return ContractError("target", "target must point to a struct or map")
	}
	return nil
}

func validateOptions(options Options) error {
	_, err := normalizeOptions(options)
	return err
}

func normalizeOptions(options Options) (Options, error) {
	if options.MergeStrategy == "" {
		options.MergeStrategy = MergeStrategyDeepReplace
	}
	if options.MergeStrategy != MergeStrategyDeepReplace {
		return options, ContractError("options.merge_strategy", "merge strategy is not supported")
	}
	if options.UnknownKeyPolicy == "" {
		options.UnknownKeyPolicy = UnknownKeyPolicyReject
	}
	switch options.UnknownKeyPolicy {
	case UnknownKeyPolicyAllow, UnknownKeyPolicyReject:
	default:
		return options, ContractError("options.unknown_key_policy", "unknown key policy is not supported")
	}
	if options.Timeout < 0 {
		return options, ContractError("options.timeout", "timeout must not be negative")
	}
	return options, nil
}

func validateSource(index int, source Source) error {
	return validateSourceFields(source, func(field string) string {
		return sourceField(index, field)
	})
}

func validateProviderSource(source Source) error {
	return validateSourceFields(source, func(field string) string {
		return "source." + field
	})
}

func validateSourceFields(source Source, fieldName func(string) string) error {
	sourceName := strings.TrimSpace(source.Name)
	if sourceName == "" {
		return ContractError(fieldName("name"), "source name is required")
	}
	if !sourceNamePattern.MatchString(sourceName) {
		return ContractError(fieldName("name"), "source name must use lowercase letters, digits, hyphens, underscores, or dots and start with a letter")
	}
	kind := strings.TrimSpace(string(source.Kind))
	if kind == "" {
		return ContractError(fieldName("kind"), "provider kind is required")
	}
	if !kindNamePattern.MatchString(kind) {
		return ContractError(fieldName("kind"), "provider kind must use lowercase letters, digits, hyphens, underscores, or dots and start with a letter")
	}
	parser := strings.TrimSpace(string(source.Parser))
	if parser == "" {
		return ContractError(fieldName("parser"), "parser kind is required")
	}
	if !kindNamePattern.MatchString(parser) {
		return ContractError(fieldName("parser"), "parser kind must use lowercase letters, digits, hyphens, underscores, or dots and start with a letter")
	}
	if source.Kind == ProviderKindEnv && source.Parser != ParserKindNone {
		return ContractError(fieldName("parser"), "environment sources must use the none parser")
	}
	if source.Kind == ProviderKindFile && source.Parser == ParserKindNone {
		return ContractError(fieldName("parser"), "file sources must declare a parser")
	}
	if strings.TrimSpace(source.Location) == "" {
		return ContractError(fieldName("location"), "source location is required")
	}
	if source.Priority < 0 {
		return ContractError(fieldName("priority"), "priority must not be negative")
	}
	return nil
}

func sourceField(index int, field string) string {
	return "sources[" + strconv.Itoa(index) + "]." + field
}
