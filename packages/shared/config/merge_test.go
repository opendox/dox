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
 * @File    : merge_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package config

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestMergeParsedSourcesDeepMergesByPriority(t *testing.T) {
	result, err := MergeParsedSources(context.Background(), []ParsedSource{
		parsedFileSource("local", 20, map[string]any{
			"app": map[string]any{
				"debug": true,
			},
			"http": map[string]any{
				"timeout": "10s",
			},
		}),
		parsedFileSource("base", 10, map[string]any{
			"app": map[string]any{
				"name": "dox",
			},
			"http": map[string]any{
				"address": "127.0.0.1:8080",
				"timeout": "5s",
			},
		}),
	}, Options{})
	if err != nil {
		t.Fatalf("merge parsed sources: %v", err)
	}

	if !reflect.DeepEqual(result.SourceNames, []string{"base", "local"}) {
		t.Fatalf("unexpected source names: %+v", result.SourceNames)
	}
	app := assertMap(t, result.Values["app"])
	if got := app["name"]; got != "dox" {
		t.Fatalf("expected app.name from base, got %v", got)
	}
	if got := app["debug"]; got != true {
		t.Fatalf("expected app.debug from local, got %v", got)
	}
	http := assertMap(t, result.Values["http"])
	if got := http["address"]; got != "127.0.0.1:8080" {
		t.Fatalf("expected http.address from base, got %v", got)
	}
	if got := http["timeout"]; got != "10s" {
		t.Fatalf("expected http.timeout override from local, got %v", got)
	}
	assertOverride(t, result.Diagnostics.Overrides, "http.timeout", "local", "base")
}

func TestMergeParsedSourcesReplacesScalarsAndSlices(t *testing.T) {
	result, err := MergeParsedSources(context.Background(), []ParsedSource{
		parsedFileSource("base", 10, map[string]any{
			"mode":     "prod",
			"features": []any{"base", "audit"},
		}),
		parsedFileSource("local", 20, map[string]any{
			"mode":     "dev",
			"features": []any{"local"},
		}),
	}, Options{})
	if err != nil {
		t.Fatalf("merge parsed sources: %v", err)
	}

	if got := result.Values["mode"]; got != "dev" {
		t.Fatalf("expected scalar override, got %v", got)
	}
	if got := result.Values["features"]; !reflect.DeepEqual(got, []any{"local"}) {
		t.Fatalf("expected slice replacement, got %+v", got)
	}
	assertOverride(t, result.Diagnostics.Overrides, "mode", "local", "base")
	assertOverride(t, result.Diagnostics.Overrides, "features", "local", "base")
}

func TestMergeParsedSourcesExpandsEnvironmentDottedKeys(t *testing.T) {
	result, err := MergeParsedSources(context.Background(), []ParsedSource{
		parsedFileSource("base", 10, map[string]any{
			"app": map[string]any{
				"name": "dox",
			},
			"http": map[string]any{
				"address": "127.0.0.1:8080",
			},
		}),
		parsedEnvSource("env", 100, map[string]any{
			"app.name":     "dox-env",
			"http.address": "0.0.0.0:8080",
		}),
	}, Options{})
	if err != nil {
		t.Fatalf("merge parsed sources: %v", err)
	}

	app := assertMap(t, result.Values["app"])
	if got := app["name"]; got != "dox-env" {
		t.Fatalf("expected env app.name override, got %v", got)
	}
	http := assertMap(t, result.Values["http"])
	if got := http["address"]; got != "0.0.0.0:8080" {
		t.Fatalf("expected env http.address override, got %v", got)
	}
	assertOverride(t, result.Diagnostics.Overrides, "app.name", "env", "base")
	assertOverride(t, result.Diagnostics.Overrides, "http.address", "env", "base")
}

func TestMergeParsedSourcesKeepsSkippedSourceDiagnostics(t *testing.T) {
	result, err := MergeParsedSources(context.Background(), []ParsedSource{
		parsedFileSource("base", 10, map[string]any{"app": map[string]any{"name": "dox"}}),
		skippedFileSource("local", 20),
	}, Options{})
	if err != nil {
		t.Fatalf("merge parsed sources: %v", err)
	}

	if !reflect.DeepEqual(result.SourceNames, []string{"base", "local"}) {
		t.Fatalf("unexpected source names: %+v", result.SourceNames)
	}
	if len(result.Diagnostics.Sources) != 2 {
		t.Fatalf("expected two source diagnostics, got %+v", result.Diagnostics.Sources)
	}
	if !result.Diagnostics.Sources[1].Skipped {
		t.Fatalf("expected skipped diagnostic, got %+v", result.Diagnostics.Sources[1])
	}
	app := assertMap(t, result.Values["app"])
	if got := app["name"]; got != "dox" {
		t.Fatalf("expected skipped source not to override values, got %v", got)
	}
}

func TestMergeParsedSourcesRejectsDuplicatePriorities(t *testing.T) {
	_, err := MergeParsedSources(context.Background(), []ParsedSource{
		parsedFileSource("base", 10, map[string]any{}),
		parsedFileSource("local", 10, map[string]any{}),
	}, Options{})

	if !IsKind(err, ErrorKindContract) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestMergeParsedSourcesRejectsConflictingEnvironmentDottedKeys(t *testing.T) {
	_, err := MergeParsedSources(context.Background(), []ParsedSource{
		parsedEnvSource("env", 10, map[string]any{
			"app":      "dox",
			"app.name": "dox",
		}),
	}, Options{})

	if !IsKind(err, ErrorKindMerge) {
		t.Fatalf("expected merge error, got %v", err)
	}
}

func TestMergeErrorKindHelpers(t *testing.T) {
	err := MergeError("source.base.values", "failed", errors.New("boom"))

	if !IsKind(err, ErrorKindMerge) {
		t.Fatalf("expected merge error kind, got %v", err)
	}
	if !errors.Is(err, &Error{Kind: ErrorKindMerge}) {
		t.Fatal("expected errors.Is to match merge error kind")
	}
}

func parsedFileSource(name string, priority int, values map[string]any) ParsedSource {
	return ParsedSource{
		Source: Source{
			Name:     name,
			Kind:     ProviderKindFile,
			Parser:   ParserKindYAML,
			Location: "configs/" + name + ".yaml",
			Required: true,
			Priority: priority,
		},
		Values: values,
		Diagnostic: SourceDiagnostic{
			Name:     name,
			Kind:     ProviderKindFile,
			Required: true,
			Loaded:   true,
		},
	}
}

func parsedEnvSource(name string, priority int, values map[string]any) ParsedSource {
	return ParsedSource{
		Source: Source{
			Name:     name,
			Kind:     ProviderKindEnv,
			Parser:   ParserKindNone,
			Location: "DOX_SERVER_",
			Priority: priority,
		},
		Values: values,
		Diagnostic: SourceDiagnostic{
			Name:   name,
			Kind:   ProviderKindEnv,
			Loaded: true,
		},
	}
}

func skippedFileSource(name string, priority int) ParsedSource {
	source := parsedFileSource(name, priority, nil)
	source.Source.Required = false
	source.Diagnostic = SourceDiagnostic{
		Name:    name,
		Kind:    ProviderKindFile,
		Skipped: true,
		Message: "optional file source does not exist",
	}
	return source
}

func assertOverride(t *testing.T, overrides []MergeOverride, path string, source string, previousSource string) {
	t.Helper()
	for _, override := range overrides {
		if override.Path == path && override.Source == source && override.PreviousSource == previousSource {
			return
		}
	}
	t.Fatalf("expected override %s %s <- %s in %+v", path, source, previousSource, overrides)
}
