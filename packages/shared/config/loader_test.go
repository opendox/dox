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
 * @File    : loader_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-25
 * @Modified: 2026-04-25
 */

package config

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

type loaderSetting struct {
	App  loaderAppSetting  `mapstructure:"app"`
	HTTP loaderHTTPSetting `mapstructure:"http"`
}

type loaderAppSetting struct {
	Name string `mapstructure:"name"`
}

type loaderHTTPSetting struct {
	Port    int           `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
}

func TestLoadReadsMergesDecodesAndFingerprints(t *testing.T) {
	dir := t.TempDir()
	basePath := filepath.Join(dir, "base.yaml")
	writeLoaderFixture(t, basePath, `
app:
  name: dox
http:
  port: 8080
  timeout: 5s
`)

	loader := NewLoader(LoaderConfig{
		Providers: map[ProviderKind]Provider{
			ProviderKindEnv: EnvProvider{Lookup: func() []string {
				return []string{
					"DOX_SERVER_APP_NAME=dox-env",
					"DOX_SERVER_HTTP_TIMEOUT=10s",
				}
			}},
		},
	})

	buildRequest := func(target any) Request {
		return Request{
			Runtime: "server",
			Env:     "dev",
			Target:  target,
			Sources: []Source{
				{
					Name:     "base",
					Kind:     ProviderKindFile,
					Parser:   ParserKindYAML,
					Location: basePath,
					Required: true,
					Priority: 10,
				},
				{
					Name:     "env",
					Kind:     ProviderKindEnv,
					Parser:   ParserKindNone,
					Location: "DOX_SERVER_",
					Priority: 100,
				},
			},
			Options: Options{UnknownKeyPolicy: UnknownKeyPolicyReject},
		}
	}

	var target loaderSetting
	result, err := loader.Load(context.Background(), buildRequest(&target))
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if target.App.Name != "dox-env" {
		t.Fatalf("expected env override for app.name, got %q", target.App.Name)
	}
	if target.HTTP.Port != 8080 {
		t.Fatalf("expected http.port from file, got %d", target.HTTP.Port)
	}
	if target.HTTP.Timeout != 10*time.Second {
		t.Fatalf("expected env override for http.timeout, got %s", target.HTTP.Timeout)
	}
	if result.Runtime != "server" || result.Env != "dev" {
		t.Fatalf("unexpected result identity: %+v", result)
	}
	if !reflect.DeepEqual(result.SourceNames, []string{"base", "env"}) {
		t.Fatalf("unexpected source names: %+v", result.SourceNames)
	}
	if !strings.HasPrefix(result.Fingerprint, "sha256:") {
		t.Fatalf("expected sha256 fingerprint, got %q", result.Fingerprint)
	}
	assertOverride(t, result.Diagnostics.Overrides, "app.name", "env", "base")
	assertOverride(t, result.Diagnostics.Overrides, "http.timeout", "env", "base")

	var targetAgain loaderSetting
	resultAgain, err := loader.Load(context.Background(), buildRequest(&targetAgain))
	if err != nil {
		t.Fatalf("load config again: %v", err)
	}
	if resultAgain.Fingerprint != result.Fingerprint {
		t.Fatalf("expected stable fingerprint, got %q and %q", result.Fingerprint, resultAgain.Fingerprint)
	}
}

func TestLoadKeepsOptionalSourceDiagnostics(t *testing.T) {
	dir := t.TempDir()
	basePath := filepath.Join(dir, "base.yaml")
	writeLoaderFixture(t, basePath, `
app:
  name: dox
`)

	var target loaderSetting
	result, err := Load(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{
			{
				Name:     "base",
				Kind:     ProviderKindFile,
				Parser:   ParserKindYAML,
				Location: basePath,
				Required: true,
				Priority: 10,
			},
			{
				Name:     "local",
				Kind:     ProviderKindFile,
				Parser:   ParserKindYAML,
				Location: filepath.Join(dir, "missing.yaml"),
				Required: false,
				Priority: 20,
			},
		},
	})
	if err != nil {
		t.Fatalf("load config with optional source: %v", err)
	}
	if target.App.Name != "dox" {
		t.Fatalf("expected base app.name, got %q", target.App.Name)
	}
	diagnostic := findSourceDiagnostic(result.Diagnostics.Sources, "local")
	if diagnostic == nil || !diagnostic.Skipped || diagnostic.Loaded {
		t.Fatalf("expected skipped optional diagnostic, got %+v", diagnostic)
	}
}

func TestLoadAllowsEmptySourcesWhenExplicit(t *testing.T) {
	var target map[string]any
	result, err := Load(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Options: Options{AllowEmptySources: true},
	})
	if err != nil {
		t.Fatalf("load empty sources: %v", err)
	}
	if len(target) != 0 {
		t.Fatalf("expected empty target map, got %+v", target)
	}
	if len(result.SourceNames) != 0 {
		t.Fatalf("expected no source names, got %+v", result.SourceNames)
	}
	if !strings.HasPrefix(result.Fingerprint, "sha256:") {
		t.Fatalf("expected sha256 fingerprint, got %q", result.Fingerprint)
	}
}

func TestLoadReturnsParseError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "base.yaml")
	writeLoaderFixture(t, path, "app: [")

	var target map[string]any
	_, err := Load(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{{
			Name:     "base",
			Kind:     ProviderKindFile,
			Parser:   ParserKindYAML,
			Location: path,
			Required: true,
		}},
	})

	if !IsKind(err, ErrorKindParse) {
		t.Fatalf("expected parse error, got %v", err)
	}
}

func TestLoadReturnsDecodeError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "base.yaml")
	writeLoaderFixture(t, path, `
app:
  name: dox
  extra: true
`)

	var target loaderSetting
	_, err := Load(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{{
			Name:     "base",
			Kind:     ProviderKindFile,
			Parser:   ParserKindYAML,
			Location: path,
			Required: true,
		}},
	})

	if !IsKind(err, ErrorKindDecode) {
		t.Fatalf("expected decode error, got %v", err)
	}
}

func TestLoadRejectsUnsupportedProviderKind(t *testing.T) {
	var target map[string]any
	_, err := Load(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{{
			Name:     "memory",
			Kind:     ProviderKind("memory"),
			Parser:   ParserKindNone,
			Location: "memory",
			Required: true,
		}},
	})

	if !IsKind(err, ErrorKindContract) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestLoadRejectsUnsupportedParserKind(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "base.jsonnet")
	writeLoaderFixture(t, path, "{}")

	var target map[string]any
	_, err := Load(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{{
			Name:     "base",
			Kind:     ProviderKindFile,
			Parser:   ParserKind("jsonnet"),
			Location: path,
			Required: true,
		}},
	})

	if !IsKind(err, ErrorKindContract) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestLoadUsesCustomProviderAndParser(t *testing.T) {
	customProvider := ProviderKind("memory")
	customParser := ParserKind("passthrough")
	loader := NewLoader(LoaderConfig{
		Providers: map[ProviderKind]Provider{
			customProvider: ProviderFunc(func(ctx context.Context, source Source) (*Payload, error) {
				return &Payload{
					Source: source,
					Values: map[string]any{
						"app": map[string]any{"name": "custom"},
					},
					Diagnostic: SourceDiagnostic{
						Name:   source.Name,
						Kind:   source.Kind,
						Loaded: true,
					},
				}, nil
			}),
		},
		Parsers: map[ParserKind]Parser{
			customParser: ParserFunc(func(ctx context.Context, payload Payload) (map[string]any, error) {
				return cloneStructuredMap(payload.Values), nil
			}),
		},
	})

	var target loaderSetting
	result, err := loader.Load(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{{
			Name:     "memory",
			Kind:     customProvider,
			Parser:   customParser,
			Location: "memory",
			Required: true,
		}},
	})
	if err != nil {
		t.Fatalf("load custom source: %v", err)
	}
	if target.App.Name != "custom" {
		t.Fatalf("expected custom app.name, got %q", target.App.Name)
	}
	if !reflect.DeepEqual(result.SourceNames, []string{"memory"}) {
		t.Fatalf("unexpected source names: %+v", result.SourceNames)
	}
}

func TestLoadHonorsTimeout(t *testing.T) {
	slowProvider := ProviderKind("slow")
	loader := NewLoader(LoaderConfig{
		Providers: map[ProviderKind]Provider{
			slowProvider: ProviderFunc(func(ctx context.Context, source Source) (*Payload, error) {
				<-ctx.Done()
				return nil, SourceError("ctx", "context is done", ctx.Err())
			}),
		},
	})

	var target map[string]any
	_, err := loader.Load(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{{
			Name:     "slow",
			Kind:     slowProvider,
			Parser:   ParserKindNone,
			Location: "memory",
			Required: true,
		}},
		Options: Options{Timeout: time.Nanosecond},
	})

	if !IsKind(err, ErrorKindSource) {
		t.Fatalf("expected source timeout error, got %v", err)
	}
}

func writeLoaderFixture(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
}

func findSourceDiagnostic(diagnostics []SourceDiagnostic, name string) *SourceDiagnostic {
	for index := range diagnostics {
		if diagnostics[index].Name == name {
			return &diagnostics[index]
		}
	}
	return nil
}
