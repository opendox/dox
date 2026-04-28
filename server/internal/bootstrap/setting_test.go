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
 * @File    : setting_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-28
 * @Modified: 2026-04-28
 */

package bootstrap

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	sharedconfig "github.com/opendox/dox/packages/shared/config"
	sharedlogging "github.com/opendox/dox/packages/shared/logging"
	sharedsetting "github.com/opendox/dox/packages/shared/setting"
	serversetting "github.com/opendox/dox/server/internal/setting"
)

func TestLoadSettingBuildsValidatedSnapshot(t *testing.T) {
	dir := t.TempDir()
	writeBootstrapFixture(t, filepath.Join(dir, "base.yaml"), `
identity:
  service:
    name: api
logging:
  level: debug
`)

	snapshot, err := LoadSetting(context.Background(), ConfigOptions{
		ConfigDir: dir,
		Env:       "prod",
		Format:    "yaml",
		EnvPrefix: "DOX_BOOTSTRAP_TEST_SETTING_SUCCESS_",
	})
	if err != nil {
		t.Fatalf("load setting snapshot: %v", err)
	}

	if snapshot.Runtime != configRuntime || snapshot.Env != "prod" {
		t.Fatalf("unexpected setting snapshot identity: %+v", snapshot)
	}
	if !strings.HasPrefix(snapshot.Fingerprint, "sha256:") {
		t.Fatalf("expected sha256 fingerprint, got %q", snapshot.Fingerprint)
	}
	if len(snapshot.SourceNames) != 4 {
		t.Fatalf("expected source names to be preserved, got %+v", snapshot.SourceNames)
	}
	if len(snapshot.Diagnostics.Sources) != 4 {
		t.Fatalf("expected diagnostics to be preserved, got %+v", snapshot.Diagnostics)
	}
	if snapshot.Setting.Identity.System.Runtime != serversetting.ServerRuntime {
		t.Fatalf("expected server runtime default, got %q", snapshot.Setting.Identity.System.Runtime)
	}
	if snapshot.Setting.Identity.Service.Name != "api" {
		t.Fatalf("expected decoded service name, got %q", snapshot.Setting.Identity.Service.Name)
	}
	if snapshot.Setting.Identity.Deployment.Env != sharedsetting.EnvProd {
		t.Fatalf("expected deployment env from config options, got %q", snapshot.Setting.Identity.Deployment.Env)
	}
	if snapshot.Setting.Logging.Level != sharedlogging.LevelDebug {
		t.Fatalf("expected decoded logging level debug, got %q", snapshot.Setting.Logging.Level)
	}
	if len(snapshot.Setting.Logging.Cores) == 0 {
		t.Fatal("expected logging defaults to be applied")
	}
}

func TestLoadConfigAllowsUnknownRawKeysAndLoadSettingRejectsThem(t *testing.T) {
	dir := t.TempDir()
	writeBootstrapFixture(t, filepath.Join(dir, "base.yaml"), `
identity: {}
unexpected_group:
  enabled: true
`)

	rawSnapshot, err := LoadConfig(context.Background(), ConfigOptions{
		ConfigDir: dir,
		Env:       "dev",
		Format:    "yaml",
		EnvPrefix: "DOX_BOOTSTRAP_TEST_SETTING_UNKNOWN_RAW_",
	})
	if err != nil {
		t.Fatalf("load raw config snapshot: %v", err)
	}
	if _, exists := rawSnapshot.Values["unexpected_group"]; !exists {
		t.Fatalf("expected raw snapshot to keep unknown key, got %+v", rawSnapshot.Values)
	}

	_, err = LoadSetting(context.Background(), ConfigOptions{
		ConfigDir: dir,
		Env:       "dev",
		Format:    "yaml",
		EnvPrefix: "DOX_BOOTSTRAP_TEST_SETTING_UNKNOWN_TYPED_",
	})
	if !sharedconfig.IsKind(err, sharedconfig.ErrorKindDecode) {
		t.Fatalf("expected typed setting decode error for unknown key, got %v", err)
	}
	if !strings.Contains(err.Error(), "unexpected_group") {
		t.Fatalf("expected unknown key to be identified, got %v", err)
	}
}

func TestLoadSettingReturnsValidationErrorForInvalidSetting(t *testing.T) {
	dir := t.TempDir()
	writeBootstrapFixture(t, filepath.Join(dir, "base.yaml"), `
identity:
  system:
    runtime: scheduler
`)

	_, err := LoadSetting(context.Background(), ConfigOptions{
		ConfigDir: dir,
		Env:       "dev",
		Format:    "yaml",
		EnvPrefix: "DOX_BOOTSTRAP_TEST_SETTING_INVALID_",
	})
	if !errors.Is(err, serversetting.ErrInvalidServerRuntime) {
		t.Fatalf("expected invalid server runtime validation error, got %v", err)
	}
}

func TestLoadSettingPropagatesSourceError(t *testing.T) {
	_, err := LoadSetting(context.Background(), ConfigOptions{
		ConfigDir: t.TempDir(),
		Env:       "dev",
		Format:    "yaml",
		EnvPrefix: "DOX_BOOTSTRAP_TEST_SETTING_MISSING_BASE_",
	})

	if !sharedconfig.IsKind(err, sharedconfig.ErrorKindSource) {
		t.Fatalf("expected source error, got %v", err)
	}
}

func TestAssembleSettingRejectsNilSnapshot(t *testing.T) {
	_, err := AssembleSetting(context.Background(), nil)

	if !sharedconfig.IsKind(err, sharedconfig.ErrorKindContract) {
		t.Fatalf("expected contract error, got %v", err)
	}
}
