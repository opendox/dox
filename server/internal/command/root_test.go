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
 * @File    : root_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-05-31
 */

package command

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opendox/dox/server/internal/bootstrap"
)

func TestRootCommandMetadata(t *testing.T) {
	cmd := NewRootCommand(Config{})

	if cmd.Use != "dox-server" {
		t.Fatalf("expected root command use dox-server, got %q", cmd.Use)
	}
	if cmd.Short == "" {
		t.Fatal("expected root command short description")
	}
}

func TestRootCommandRegistersServeCommand(t *testing.T) {
	cmd := NewRootCommand(Config{})

	serve, _, err := cmd.Find([]string{"serve"})
	if err != nil {
		t.Fatalf("find serve command: %v", err)
	}
	if serve == nil || serve.Use != "serve" {
		t.Fatalf("expected serve command to be registered, got %#v", serve)
	}
}

func TestVersionCommandPrintsSharedVersionInfo(t *testing.T) {
	var out bytes.Buffer

	err := Execute(context.Background(), []string{"version"}, Config{Out: &out})
	if err != nil {
		t.Fatalf("expected version command to succeed: %v", err)
	}

	output := out.String()
	for _, expected := range []string{
		"dox 0.1.0",
		"Go Version",
		"Git Commit",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected version output to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestServeCommandLoadsSettingBootsAndStops(t *testing.T) {
	dir := t.TempDir()
	logDir := t.TempDir()
	writeCommandFixture(t, filepath.Join(dir, "base.yaml"), `
identity:
  service:
    name: api
logging:
  resource:
    service_version: command-test
  zap:
    disable_caller: true
    disable_stacktrace: true
    disable_error_verbose: true
    error_output_paths:
      - "`+filepath.Join(logDir, "errors.log")+`"
  cores:
    - name: service-file
      enabled: true
      type: file
      level: info
      encoding: json
      output_paths:
        - "`+filepath.Join(logDir, "${service.namespace}-${service.name}.jsonl")+`"
      datasets: ["*"]
      rotation:
        driver: none
`)

	var booted bool
	err := Execute(context.Background(), []string{
		"serve",
		"--config-dir", dir,
		"--env", "test",
		"--config-format", "yaml",
		"--env-prefix", "DOX_COMMAND_TEST_SERVE_SUCCESS_",
		"--config-timeout", "2s",
	}, Config{
		ServeWait: func(ctx context.Context, runtime *bootstrap.Runtime) error {
			booted = true
			if runtime == nil || runtime.Logger == nil {
				t.Fatalf("expected booted runtime with logger, got %#v", runtime)
			}
			if runtime.Resource.ServiceName != "api" || runtime.Resource.ServiceVersion != "command-test" {
				t.Fatalf("expected serve command to boot loaded setting, got %+v", runtime.Resource)
			}
			runtime.Logger.Info(ctx, "serve command ready")
			return nil
		},
	})
	if err != nil {
		t.Fatalf("expected serve command to succeed: %v", err)
	}
	if !booted {
		t.Fatal("expected serve wait hook to observe a booted runtime")
	}

	payload, err := os.ReadFile(filepath.Join(logDir, "dox-api.jsonl"))
	if err != nil {
		t.Fatalf("expected serve shutdown to flush runtime logs: %v", err)
	}
	if !strings.Contains(string(payload), `"message":"serve command ready"`) {
		t.Fatalf("expected serve runtime log to be flushed, got:\n%s", payload)
	}
}

func TestServeCommandReturnsConfigLoadError(t *testing.T) {
	err := Execute(context.Background(), []string{
		"serve",
		"--config-dir", t.TempDir(),
		"--env", "test",
		"--config-format", "yaml",
		"--env-prefix", "DOX_COMMAND_TEST_SERVE_MISSING_",
	}, Config{
		ServeWait: func(context.Context, *bootstrap.Runtime) error {
			t.Fatal("serve wait must not run when setting loading fails")
			return nil
		},
	})

	if err == nil {
		t.Fatal("expected serve command to return config load error")
	}
}

func TestWaitForContextCancellationReturnsAfterCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := waitForContextCancellation(ctx, nil); err != nil {
		t.Fatalf("expected canceled context to stop wait without an error, got %v", err)
	}
}

func writeCommandFixture(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(strings.TrimLeft(content, "\n")), 0o644); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
}
