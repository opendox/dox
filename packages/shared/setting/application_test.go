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
 * @File    : application_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-25
 * @Modified: 2026-04-25
 */

package setting

import (
	"errors"
	"testing"
)

func TestApplicationDefaultDerivesSharedIdentity(t *testing.T) {
	app := Application{Runtime: RuntimeServer}

	if err := app.Default(); err != nil {
		t.Fatalf("default application: %v", err)
	}

	if app.Name != "dox" {
		t.Fatalf("expected default name dox, got %q", app.Name)
	}
	if app.Namespace != "opendox" {
		t.Fatalf("expected default namespace opendox, got %q", app.Namespace)
	}
	if app.Env != EnvDev {
		t.Fatalf("expected default env dev, got %q", app.Env)
	}
	if app.Runtime != RuntimeServer {
		t.Fatalf("expected runtime to remain server, got %q", app.Runtime)
	}
	if app.Service != "dox-server" {
		t.Fatalf("expected derived service dox-server, got %q", app.Service)
	}
	if err := app.Validate(); err != nil {
		t.Fatalf("validate defaulted application: %v", err)
	}
}

func TestApplicationDefaultDoesNotInventRuntime(t *testing.T) {
	app := Application{}

	if err := app.Default(); err != nil {
		t.Fatalf("default application: %v", err)
	}
	if app.Runtime != "" {
		t.Fatalf("expected runtime to remain empty, got %q", app.Runtime)
	}
	if app.Service != "" {
		t.Fatalf("expected service to remain empty without runtime, got %q", app.Service)
	}
	if err := app.Validate(); !hasValidationField(err, "Application.runtime", "required") {
		t.Fatalf("expected required runtime validation error, got %v", err)
	}
}

func TestApplicationDefaultPreservesExplicitService(t *testing.T) {
	app := Application{
		Name:      "dox",
		Namespace: "opendox",
		Runtime:   RuntimeCollector,
		Service:   "amazon-collector",
		Env:       EnvProd,
	}

	if err := app.Default(); err != nil {
		t.Fatalf("default application: %v", err)
	}
	if app.Service != "amazon-collector" {
		t.Fatalf("expected explicit service to remain, got %q", app.Service)
	}
	if err := app.Validate(); err != nil {
		t.Fatalf("validate application: %v", err)
	}
}

func TestApplicationValidateRejectsInvalidRuntime(t *testing.T) {
	app := validApplication()
	app.Runtime = Runtime("worker")

	if err := app.Validate(); !hasValidationField(err, "Application.runtime", "dox_runtime") {
		t.Fatalf("expected invalid runtime validation error, got %v", err)
	}
}

func TestApplicationValidateRejectsInvalidEnv(t *testing.T) {
	app := validApplication()
	app.Env = Env("production")

	if err := app.Validate(); !hasValidationField(err, "Application.env", "dox_env") {
		t.Fatalf("expected invalid env validation error, got %v", err)
	}
}

func TestApplicationValidateRejectsInvalidIdentifierSyntax(t *testing.T) {
	app := validApplication()
	app.Service = "Dox Server"

	if err := app.Validate(); !hasValidationField(err, "Application.service", "dox_kebab") {
		t.Fatalf("expected invalid service validation error, got %v", err)
	}
}

func TestApplicationValidateRejectsTrailingHyphen(t *testing.T) {
	app := validApplication()
	app.Service = "dox-server-"

	if err := app.Validate(); !hasValidationField(err, "Application.service", "dox_kebab") {
		t.Fatalf("expected invalid service validation error, got %v", err)
	}
}

func TestRuntimeAndEnvValidity(t *testing.T) {
	if !RuntimeScheduler.IsValid() || !RuntimeCollector.IsValid() || !RuntimeCompute.IsValid() {
		t.Fatal("expected supported runtimes to be valid")
	}
	if Runtime("queue").IsValid() {
		t.Fatal("did not expect unsupported runtime to be valid")
	}
	if !EnvTest.IsValid() || !EnvStaging.IsValid() || !EnvProd.IsValid() {
		t.Fatal("expected supported envs to be valid")
	}
	if Env("stage").IsValid() {
		t.Fatal("did not expect unsupported env to be valid")
	}
}

func validApplication() Application {
	return Application{
		Name:      "dox",
		Namespace: "opendox",
		Runtime:   RuntimeServer,
		Service:   "dox-server",
		Env:       EnvDev,
	}
}

func hasValidationField(err error, field string, rule string) bool {
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		return false
	}
	for _, item := range validationErr.Fields {
		if item.Field == field && item.Rule == rule {
			return true
		}
	}
	return false
}
