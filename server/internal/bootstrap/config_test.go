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
 * @File    : config_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-25
 * @Modified: 2026-04-25
 */

package bootstrap

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	sharedconfig "github.com/opendox/dox/packages/shared/config"
)

func TestBuildConfigSourcesUsesDefaultConvention(t *testing.T) {
	options, err := normalizeConfigOptions(ConfigOptions{})
	if err != nil {
		t.Fatalf("normalize default options: %v", err)
	}

	sources := buildConfigSources(options)
	if len(sources) != 4 {
		t.Fatalf("expected four config sources, got %+v", sources)
	}

	expected := []sharedconfig.Source{
		{
			Name:     "base",
			Kind:     sharedconfig.ProviderKindFile,
			Parser:   sharedconfig.ParserKindYAML,
			Location: filepath.Join("configs", "base.yaml"),
			Required: true,
			Priority: 10,
		},
		{
			Name:     "environment",
			Kind:     sharedconfig.ProviderKindFile,
			Parser:   sharedconfig.ParserKindYAML,
			Location: filepath.Join("configs", "dev.yaml"),
			Required: false,
			Priority: 20,
		},
		{
			Name:     "local",
			Kind:     sharedconfig.ProviderKindFile,
			Parser:   sharedconfig.ParserKindYAML,
			Location: filepath.Join("configs", "local.yaml"),
			Required: false,
			Priority: 30,
		},
		{
			Name:     "env",
			Kind:     sharedconfig.ProviderKindEnv,
			Parser:   sharedconfig.ParserKindNone,
			Location: "DOX_SERVER_",
			Required: false,
			Priority: 100,
		},
	}
	if !reflect.DeepEqual(sources, expected) {
		t.Fatalf("unexpected default sources:\nwant: %+v\n got: %+v", expected, sources)
	}
}

func TestLoadConfigBuildsSnapshotAndAppliesOverrides(t *testing.T) {
	dir := t.TempDir()
	writeBootstrapFixture(t, filepath.Join(dir, "base.yaml"), `
app:
  name: base
http:
  timeout: 5s
`)
	writeBootstrapFixture(t, filepath.Join(dir, "dev.yaml"), `
app:
  environment: dev
http:
  timeout: 10s
`)

	prefix := "DOX_BOOTSTRAP_TEST_SNAPSHOT_"
	t.Setenv(prefix+"APP_NAME", "env")
	t.Setenv("DOX_ENV", "prod")
	t.Setenv("DOX_CONFIG_DIR", "/tmp/not-used")
	t.Setenv("DOX_CONFIG_FORMAT", "json")

	snapshot, err := LoadConfig(context.Background(), ConfigOptions{
		ConfigDir: dir,
		Env:       "dev",
		Format:    "yaml",
		EnvPrefix: prefix,
	})
	if err != nil {
		t.Fatalf("load config snapshot: %v", err)
	}

	if snapshot.Runtime != "server" || snapshot.Env != "dev" {
		t.Fatalf("unexpected snapshot identity: %+v", snapshot)
	}
	if !reflect.DeepEqual(snapshot.SourceNames, []string{"base", "environment", "local", "env"}) {
		t.Fatalf("unexpected source names: %+v", snapshot.SourceNames)
	}
	if !strings.HasPrefix(snapshot.Fingerprint, "sha256:") {
		t.Fatalf("expected sha256 fingerprint, got %q", snapshot.Fingerprint)
	}

	app := assertBootstrapMap(t, snapshot.Values["app"])
	if got := app["name"]; got != "env" {
		t.Fatalf("expected app.name from env override, got %v", got)
	}
	if got := app["environment"]; got != "dev" {
		t.Fatalf("expected app.environment from environment file, got %v", got)
	}
	http := assertBootstrapMap(t, snapshot.Values["http"])
	if got := http["timeout"]; got != "10s" {
		t.Fatalf("expected http.timeout from environment file, got %v", got)
	}
	if _, exists := snapshot.Values["env"]; exists {
		t.Fatalf("did not expect bootstrap control DOX_ENV in config values: %+v", snapshot.Values)
	}
	if _, exists := snapshot.Values["config"]; exists {
		t.Fatalf("did not expect bootstrap control DOX_CONFIG_* in config values: %+v", snapshot.Values)
	}

	localDiagnostic := findBootstrapDiagnostic(snapshot.Diagnostics.Sources, "local")
	if localDiagnostic == nil || !localDiagnostic.Skipped || localDiagnostic.Loaded {
		t.Fatalf("expected skipped local diagnostic, got %+v", localDiagnostic)
	}
	assertBootstrapOverride(t, snapshot.Diagnostics.Overrides, "app.name", "env", "base")
	assertBootstrapOverride(t, snapshot.Diagnostics.Overrides, "http.timeout", "environment", "base")
}

func TestLoadConfigSkipsMissingOptionalSources(t *testing.T) {
	dir := t.TempDir()
	writeBootstrapFixture(t, filepath.Join(dir, "base.yaml"), `
app:
  name: base
`)

	snapshot, err := LoadConfig(context.Background(), ConfigOptions{
		ConfigDir: dir,
		Env:       "dev",
		Format:    "yaml",
		EnvPrefix: "DOX_BOOTSTRAP_TEST_EMPTY_18_",
	})
	if err != nil {
		t.Fatalf("load config snapshot: %v", err)
	}

	app := assertBootstrapMap(t, snapshot.Values["app"])
	if got := app["name"]; got != "base" {
		t.Fatalf("expected app.name from base file, got %v", got)
	}
	for _, name := range []string{"environment", "local", "env"} {
		diagnostic := findBootstrapDiagnostic(snapshot.Diagnostics.Sources, name)
		if diagnostic == nil || !diagnostic.Skipped || diagnostic.Loaded {
			t.Fatalf("expected skipped diagnostic for %s, got %+v", name, diagnostic)
		}
	}
}

func TestLoadConfigReturnsSourceErrorWhenBaseMissing(t *testing.T) {
	_, err := LoadConfig(context.Background(), ConfigOptions{
		ConfigDir: t.TempDir(),
		Env:       "dev",
		Format:    "yaml",
		EnvPrefix: "DOX_BOOTSTRAP_TEST_MISSING_BASE_",
	})

	if !sharedconfig.IsKind(err, sharedconfig.ErrorKindSource) {
		t.Fatalf("expected source error, got %v", err)
	}
}

func TestLoadConfigRejectsUnsupportedFormat(t *testing.T) {
	_, err := LoadConfig(context.Background(), ConfigOptions{
		ConfigDir: t.TempDir(),
		Env:       "dev",
		Format:    "ini",
		EnvPrefix: "DOX_BOOTSTRAP_TEST_FORMAT_",
	})

	if !sharedconfig.IsKind(err, sharedconfig.ErrorKindContract) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func writeBootstrapFixture(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(strings.TrimLeft(body, "\n")), 0o600); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
}

func assertBootstrapMap(t *testing.T, value any) map[string]any {
	t.Helper()
	values, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("expected map value, got %T: %v", value, value)
	}
	return values
}

func findBootstrapDiagnostic(diagnostics []sharedconfig.SourceDiagnostic, name string) *sharedconfig.SourceDiagnostic {
	for index := range diagnostics {
		if diagnostics[index].Name == name {
			return &diagnostics[index]
		}
	}
	return nil
}

func assertBootstrapOverride(t *testing.T, overrides []sharedconfig.MergeOverride, path string, source string, previous string) {
	t.Helper()
	for _, override := range overrides {
		if override.Path == path && override.Source == source && override.PreviousSource == previous {
			return
		}
	}
	t.Fatalf("expected override path=%s source=%s previous=%s in %+v", path, source, previous, overrides)
}
