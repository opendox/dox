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
 * @File    : identity_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package setting

import (
	"errors"
	"testing"

	sharedsetting "github.com/opendox/dox/packages/shared/setting"
)

func TestIdentityDefaultAppliesServerIdentity(t *testing.T) {
	identity := Identity{}

	if err := identity.DefaultWithOptions(IdentityDefaultOptions{Env: "staging"}); err != nil {
		t.Fatalf("default identity: %v", err)
	}

	if identity.Organization.Name != sharedsetting.DefaultOrganizationName {
		t.Fatalf("expected default organization, got %q", identity.Organization.Name)
	}
	if identity.Application.Name != sharedsetting.DefaultApplicationName {
		t.Fatalf("expected default application, got %q", identity.Application.Name)
	}
	if identity.System.Runtime != ServerRuntime {
		t.Fatalf("expected server runtime, got %q", identity.System.Runtime)
	}
	if identity.Service.Namespace != sharedsetting.DefaultApplicationName {
		t.Fatalf("expected service namespace from application, got %q", identity.Service.Namespace)
	}
	if identity.Service.Name != string(ServerRuntime) {
		t.Fatalf("expected service name from runtime, got %q", identity.Service.Name)
	}
	if identity.Deployment.Env != sharedsetting.EnvStaging {
		t.Fatalf("expected deployment env from options, got %q", identity.Deployment.Env)
	}
	if err := identity.Validate(); err != nil {
		t.Fatalf("validate identity: %v", err)
	}
}

func TestIdentityDefaultPreservesExplicitValues(t *testing.T) {
	identity := Identity{
		Application: sharedsetting.Application{Name: "dox"},
		System:      sharedsetting.System{Runtime: ServerRuntime},
		Service: sharedsetting.Service{
			Namespace:  "platform",
			Name:       "iam",
			InstanceID: "server-pod-1",
		},
		Deployment: sharedsetting.Deployment{Env: sharedsetting.EnvProd},
	}

	if err := identity.DefaultWithOptions(IdentityDefaultOptions{Env: "staging"}); err != nil {
		t.Fatalf("default identity: %v", err)
	}

	if identity.Service.Namespace != "platform" {
		t.Fatalf("expected explicit service namespace to remain, got %q", identity.Service.Namespace)
	}
	if identity.Service.Name != "iam" {
		t.Fatalf("expected explicit service name to remain, got %q", identity.Service.Name)
	}
	if identity.Service.InstanceID != "server-pod-1" {
		t.Fatalf("expected explicit service instance id to remain, got %q", identity.Service.InstanceID)
	}
	if identity.Deployment.Env != sharedsetting.EnvProd {
		t.Fatalf("expected explicit deployment env to remain, got %q", identity.Deployment.Env)
	}
	if err := identity.Validate(); err != nil {
		t.Fatalf("validate identity: %v", err)
	}
}

func TestIdentityDefaultFallsBackToSharedDeploymentEnv(t *testing.T) {
	identity := Identity{}

	if err := identity.Default(); err != nil {
		t.Fatalf("default identity: %v", err)
	}
	if identity.Deployment.Env != sharedsetting.EnvDev {
		t.Fatalf("expected shared deployment env default, got %q", identity.Deployment.Env)
	}
}

func TestIdentityValidateRejectsNonServerRuntime(t *testing.T) {
	identity := validIdentity()
	identity.System.Runtime = sharedsetting.RuntimeScheduler

	if err := identity.Validate(); !errors.Is(err, ErrInvalidServerRuntime) {
		t.Fatalf("expected invalid server runtime error, got %v", err)
	}
}

func TestIdentityValidateRejectsSharedIdentityFailures(t *testing.T) {
	identity := validIdentity()
	identity.Service.Name = "Dox Server"

	if err := identity.Validate(); !hasSharedValidationField(err, "Service.name", "dox_kebab") {
		t.Fatalf("expected shared service validation error, got %v", err)
	}
}

func TestIdentityDefaultRejectsNilReceiver(t *testing.T) {
	var identity *Identity

	if err := identity.Default(); err == nil {
		t.Fatal("expected nil identity default error")
	}
}

func validIdentity() Identity {
	return Identity{
		Organization: sharedsetting.Organization{Name: sharedsetting.DefaultOrganizationName},
		Application:  sharedsetting.Application{Name: sharedsetting.DefaultApplicationName},
		System:       sharedsetting.System{Runtime: ServerRuntime},
		Service: sharedsetting.Service{
			Namespace: sharedsetting.DefaultApplicationName,
			Name:      string(ServerRuntime),
		},
		Deployment: sharedsetting.Deployment{Env: sharedsetting.EnvDev},
	}
}

func hasSharedValidationField(err error, field string, rule string) bool {
	var validationErr *sharedsetting.ValidationError
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
