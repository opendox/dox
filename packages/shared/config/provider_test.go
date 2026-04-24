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
 * @File    : provider_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileProviderReadsRequiredFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "base.yaml")
	if err := os.WriteFile(path, []byte("app:\n  name: dox\n"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	payload, err := FileProvider{}.Read(context.Background(), Source{
		Name:     "base",
		Kind:     ProviderKindFile,
		Parser:   ParserKindYAML,
		Location: path,
		Required: true,
	})
	if err != nil {
		t.Fatalf("read file source: %v", err)
	}
	if string(payload.Raw) != "app:\n  name: dox\n" {
		t.Fatalf("unexpected raw payload: %q", payload.Raw)
	}
	if !payload.Diagnostic.Loaded || payload.Diagnostic.Skipped {
		t.Fatalf("unexpected diagnostic: %+v", payload.Diagnostic)
	}
}

func TestFileProviderRejectsMissingRequiredFile(t *testing.T) {
	_, err := FileProvider{}.Read(context.Background(), Source{
		Name:     "base",
		Kind:     ProviderKindFile,
		Parser:   ParserKindYAML,
		Location: filepath.Join(t.TempDir(), "missing.yaml"),
		Required: true,
	})

	if !IsKind(err, ErrorKindSource) {
		t.Fatalf("expected source error, got %v", err)
	}
}

func TestFileProviderSkipsMissingOptionalFile(t *testing.T) {
	payload, err := FileProvider{}.Read(context.Background(), Source{
		Name:     "local",
		Kind:     ProviderKindFile,
		Parser:   ParserKindYAML,
		Location: filepath.Join(t.TempDir(), "missing.yaml"),
		Required: false,
	})
	if err != nil {
		t.Fatalf("expected missing optional file to be skipped, got %v", err)
	}
	if !payload.Diagnostic.Skipped || payload.Diagnostic.Loaded {
		t.Fatalf("unexpected diagnostic: %+v", payload.Diagnostic)
	}
}

func TestEnvProviderFiltersByPrefix(t *testing.T) {
	payload, err := EnvProvider{Lookup: func() []string {
		return []string{
			"DOX_SERVER_APP_NAME=dox",
			"DOX_SERVER_HTTP__ADDRESS=127.0.0.1:8080",
			"OTHER_VALUE=ignored",
		}
	}}.Read(context.Background(), Source{
		Name:     "env",
		Kind:     ProviderKindEnv,
		Parser:   ParserKindNone,
		Location: "DOX_SERVER_",
	})
	if err != nil {
		t.Fatalf("read env source: %v", err)
	}
	if got := payload.Values["app.name"]; got != "dox" {
		t.Fatalf("expected app.name to be dox, got %v", got)
	}
	if got := payload.Values["http.address"]; got != "127.0.0.1:8080" {
		t.Fatalf("expected http.address, got %v", got)
	}
	if _, exists := payload.Values["other.value"]; exists {
		t.Fatal("did not expect unrelated env value")
	}
}

func TestEnvProviderRejectsMissingRequiredValues(t *testing.T) {
	_, err := EnvProvider{Lookup: func() []string {
		return []string{"OTHER_VALUE=ignored"}
	}}.Read(context.Background(), Source{
		Name:     "env",
		Kind:     ProviderKindEnv,
		Parser:   ParserKindNone,
		Location: "DOX_SERVER_",
		Required: true,
	})

	if !IsKind(err, ErrorKindSource) {
		t.Fatalf("expected source error, got %v", err)
	}
}

func TestEnvProviderSkipsMissingOptionalValues(t *testing.T) {
	payload, err := EnvProvider{Lookup: func() []string {
		return []string{"OTHER_VALUE=ignored"}
	}}.Read(context.Background(), Source{
		Name:     "env",
		Kind:     ProviderKindEnv,
		Parser:   ParserKindNone,
		Location: "DOX_SERVER_",
		Required: false,
	})
	if err != nil {
		t.Fatalf("expected missing optional environment values to be skipped, got %v", err)
	}
	if !payload.Diagnostic.Skipped || payload.Diagnostic.Loaded {
		t.Fatalf("unexpected diagnostic: %+v", payload.Diagnostic)
	}
	if got := payload.Metadata["prefix"]; got != "DOX_SERVER_" {
		t.Fatalf("expected prefix metadata, got %q", got)
	}
}

func TestProvidersRejectInvalidSourceContracts(t *testing.T) {
	_, err := FileProvider{}.Read(context.Background(), Source{
		Name:     "env",
		Kind:     ProviderKindEnv,
		Parser:   ParserKindNone,
		Location: "DOX_SERVER_",
	})
	if !IsKind(err, ErrorKindContract) {
		t.Fatalf("expected file provider contract error, got %v", err)
	}

	_, err = EnvProvider{}.Read(context.Background(), Source{
		Name:     "base",
		Kind:     ProviderKindFile,
		Parser:   ParserKindYAML,
		Location: "configs/base.yaml",
	})
	if !IsKind(err, ErrorKindContract) {
		t.Fatalf("expected env provider contract error, got %v", err)
	}
}
