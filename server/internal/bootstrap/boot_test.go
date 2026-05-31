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
 * @File    : boot_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-05-31
 * @Modified: 2026-05-31
 */

package bootstrap

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	sharedconfig "github.com/opendox/dox/packages/shared/config"
	sharedlogging "github.com/opendox/dox/packages/shared/logging"
	sharedsetting "github.com/opendox/dox/packages/shared/setting"
	serversetting "github.com/opendox/dox/server/internal/setting"
)

func TestBootBuildsLoggingRuntimeFromValidatedSetting(t *testing.T) {
	tempDir := t.TempDir()
	snapshot := newBootTestSnapshot(t, filepath.Join(tempDir, "${service.namespace}-${service.name}.jsonl"))

	runtime, err := Boot(context.Background(), snapshot, BootOptions{})
	if err != nil {
		t.Fatalf("boot runtime: %v", err)
	}

	if runtime.Snapshot != snapshot {
		t.Fatal("expected runtime to retain the setting snapshot")
	}
	if runtime.Logger == nil {
		t.Fatal("expected runtime logger to be initialized")
	}
	if runtime.Resource.ServiceNamespace != "dox" || runtime.Resource.ServiceName != "api" {
		t.Fatalf("expected logging resource from identity, got %+v", runtime.Resource)
	}
	if runtime.Resource.ServiceVersion != "test-version" {
		t.Fatalf("expected logging resource version override, got %+v", runtime.Resource)
	}
	if runtime.Resource.DeploymentEnvironment != "test" ||
		runtime.Resource.DoxOrganization != "opendox" ||
		runtime.Resource.DoxRuntime != "server" {
		t.Fatalf("expected logging resource deployment identity, got %+v", runtime.Resource)
	}

	runtime.Logger.Info(context.Background(), "boot logger ready")
	if err := runtime.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown runtime: %v", err)
	}

	logPath := filepath.Join(tempDir, "dox-api.jsonl")
	payload, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read rendered log file %s: %v", logPath, err)
	}
	for _, expected := range []string{
		`"message":"boot logger ready"`,
		`"service.namespace":"dox"`,
		`"service.name":"api"`,
		`"service.version":"test-version"`,
		`"deployment.environment.name":"test"`,
		`"dox.runtime":"server"`,
	} {
		if !strings.Contains(string(payload), expected) {
			t.Fatalf("expected log payload to contain %s, got:\n%s", expected, payload)
		}
	}
}

func TestBootRejectsNilSnapshot(t *testing.T) {
	_, err := Boot(context.Background(), nil, BootOptions{})

	if !sharedconfig.IsKind(err, sharedconfig.ErrorKindContract) {
		t.Fatalf("expected contract error for nil snapshot, got %v", err)
	}
}

func TestBootRejectsUnsupportedOpenTelemetryGlobalInstallOption(t *testing.T) {
	snapshot := newBootTestSnapshot(t, filepath.Join(t.TempDir(), "service.jsonl"))

	_, err := Boot(context.Background(), snapshot, BootOptions{
		InstallOpenTelemetryGlobals: true,
	})

	if !sharedconfig.IsKind(err, sharedconfig.ErrorKindContract) {
		t.Fatalf("expected contract error for unsupported boot option, got %v", err)
	}
}

func TestBootReturnsOpenTelemetrySetupError(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "service.jsonl")
	snapshot := newBootTestSnapshot(t, logPath)
	snapshot.Setting.Logging.OTel.Enabled = boolPtr(true)
	snapshot.Setting.Logging.OTel.Exporter.OTLP.Enabled = true
	snapshot.Setting.Logging.OTel.Exporter.OTLP.Endpoint = "http://127.0.0.1:4318"

	_, err := Boot(context.Background(), snapshot, BootOptions{})

	if err == nil {
		t.Fatal("expected boot to fail when unsupported OTLP exporter setup is requested")
	}
	if !strings.Contains(err.Error(), "OpenTelemetry") {
		t.Fatalf("expected OpenTelemetry setup error, got %v", err)
	}

	if removeErr := os.Remove(logPath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
		t.Fatalf("expected boot rollback to release logging file sink, got remove error: %v", removeErr)
	}
}

func TestRuntimeLifecycleShutdownRunsHooksInReverseOrderAndAggregatesErrors(t *testing.T) {
	firstErr := errors.New("first cleanup failed")
	secondErr := errors.New("second cleanup failed")
	var order []string

	lifecycle := runtimeLifecycle{}
	lifecycle.register("first", func(context.Context) error {
		order = append(order, "first")
		return firstErr
	})
	lifecycle.register("second", func(context.Context) error {
		order = append(order, "second")
		return secondErr
	})

	err := lifecycle.shutdown(context.Background())

	if !reflect.DeepEqual(order, []string{"second", "first"}) {
		t.Fatalf("expected reverse shutdown order, got %#v", order)
	}
	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("expected shutdown error to preserve cleanup failures, got %v", err)
	}

	order = nil
	if err := lifecycle.shutdown(context.Background()); err != nil {
		t.Fatalf("expected second shutdown to be idempotent, got %v", err)
	}
	if len(order) != 0 {
		t.Fatalf("expected second shutdown to skip already-run hooks, got %#v", order)
	}
}

func newBootTestSnapshot(t *testing.T, outputPath string) *SettingSnapshot {
	t.Helper()

	setting := serversetting.Setting{
		Identity: serversetting.Identity{
			Service: sharedsetting.Service{
				Name:       "api",
				InstanceID: "server-01",
			},
			Deployment: sharedsetting.Deployment{
				Region:       "us-east-1",
				Zone:         "use1-az1",
				Cluster:      "dox-test",
				K8sNamespace: "dox-system",
			},
		},
		Logging: sharedlogging.Config{
			Resource: sharedlogging.ResourceConfig{
				ServiceVersion: "test-version",
			},
			Zap: sharedlogging.ZapConfig{
				DisableCaller:       true,
				DisableStacktrace:   true,
				DisableErrorVerbose: true,
				ErrorOutputPaths:    []string{filepath.Join(filepath.Dir(outputPath), "errors.log")},
			},
			Cores: []sharedlogging.CoreConfig{
				{
					Name:        "service-file",
					Enabled:     boolPtr(true),
					Type:        sharedlogging.CoreTypeFile,
					Level:       sharedlogging.LevelInfo,
					Encoding:    sharedlogging.EncodingJSON,
					OutputPaths: []string{outputPath},
					Datasets:    []string{"*"},
					Rotation: sharedlogging.RotationConfig{
						Driver: sharedlogging.RotationDriverNone,
					},
				},
			},
		},
	}
	if err := setting.DefaultWithOptions(serversetting.DefaultOptions{Env: "test"}); err != nil {
		t.Fatalf("default test setting: %v", err)
	}
	if err := setting.Validate(); err != nil {
		t.Fatalf("validate test setting: %v", err)
	}

	return &SettingSnapshot{
		Runtime: configRuntime,
		Env:     "test",
		Setting: setting,
	}
}

func boolPtr(value bool) *bool {
	return &value
}
