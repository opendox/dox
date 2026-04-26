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
 * @Modified: 2026-04-26
 */

package setting

import (
	"errors"
	"testing"
)

func TestIdentityFragmentsApplyConservativeDefaults(t *testing.T) {
	organization := Organization{}
	application := Application{}
	system := System{Runtime: RuntimeServer}
	service := Service{}
	deployment := Deployment{}

	for _, item := range []struct {
		name      string
		defaulter func() error
	}{
		{"organization", organization.Default},
		{"application", application.Default},
		{"system", system.Default},
		{"service", func() error {
			return service.Default(application, system)
		}},
		{"deployment", deployment.Default},
	} {
		if err := item.defaulter(); err != nil {
			t.Fatalf("default %s: %v", item.name, err)
		}
	}

	if organization.Name != "opendox" {
		t.Fatalf("expected organization default opendox, got %q", organization.Name)
	}
	if application.Name != "dox" {
		t.Fatalf("expected application default dox, got %q", application.Name)
	}
	if system.Runtime != RuntimeServer {
		t.Fatalf("expected explicit runtime to remain server, got %q", system.Runtime)
	}
	if service.Namespace != "dox" {
		t.Fatalf("expected service namespace from application name, got %q", service.Namespace)
	}
	if service.Name != "server" {
		t.Fatalf("expected service name from system runtime, got %q", service.Name)
	}
	if deployment.Env != EnvDev {
		t.Fatalf("expected deployment env default dev, got %q", deployment.Env)
	}

	for name, validator := range map[string]func() error{
		"organization": organization.Validate,
		"application":  application.Validate,
		"system":       system.Validate,
		"service":      service.Validate,
		"deployment":   deployment.Validate,
	} {
		if err := validator(); err != nil {
			t.Fatalf("validate %s: %v", name, err)
		}
	}
}

func TestSystemDefaultDoesNotInventRuntime(t *testing.T) {
	system := System{}

	if err := system.Default(); err != nil {
		t.Fatalf("default system: %v", err)
	}
	if system.Runtime != "" {
		t.Fatalf("expected runtime to remain empty, got %q", system.Runtime)
	}
	if err := system.Validate(); !hasValidationField(err, "System.runtime", "required") {
		t.Fatalf("expected required runtime validation error, got %v", err)
	}
}

func TestServiceDefaultDoesNotInventNameWithoutRuntime(t *testing.T) {
	application := Application{}
	if err := application.Default(); err != nil {
		t.Fatalf("default application: %v", err)
	}

	service := Service{}
	if err := service.Default(application, System{}); err != nil {
		t.Fatalf("default service: %v", err)
	}
	if service.Namespace != "dox" {
		t.Fatalf("expected service namespace from application name, got %q", service.Namespace)
	}
	if service.Name != "" {
		t.Fatalf("expected service name to remain empty without runtime, got %q", service.Name)
	}
	if err := service.Validate(); !hasValidationField(err, "Service.name", "required") {
		t.Fatalf("expected required service name validation error, got %v", err)
	}
}

func TestServiceDefaultPreservesExplicitName(t *testing.T) {
	application := Application{Name: "dox"}
	system := System{Runtime: RuntimeServer}
	service := Service{Name: "iam"}

	if err := service.Default(application, system); err != nil {
		t.Fatalf("default service: %v", err)
	}
	if service.Namespace != "dox" {
		t.Fatalf("expected service namespace from application name, got %q", service.Namespace)
	}
	if service.Name != "iam" {
		t.Fatalf("expected explicit service name to remain, got %q", service.Name)
	}
	if err := service.Validate(); err != nil {
		t.Fatalf("validate service: %v", err)
	}
}

func TestSystemValidateRejectsInvalidRuntime(t *testing.T) {
	system := System{Runtime: Runtime("worker")}

	if err := system.Validate(); !hasValidationField(err, "System.runtime", "dox_runtime") {
		t.Fatalf("expected invalid runtime validation error, got %v", err)
	}
}

func TestDeploymentValidateRejectsInvalidEnv(t *testing.T) {
	deployment := Deployment{Env: Env("production")}

	if err := deployment.Validate(); !hasValidationField(err, "Deployment.env", "dox_env") {
		t.Fatalf("expected invalid env validation error, got %v", err)
	}
}

func TestServiceValidateRejectsInvalidNamespaceAndName(t *testing.T) {
	service := Service{
		Namespace:  "Dox",
		Name:       "iam-service-",
		InstanceID: "server-pod-1",
	}

	err := service.Validate()
	if !hasValidationField(err, "Service.namespace", "dox_kebab") {
		t.Fatalf("expected invalid namespace validation error, got %v", err)
	}
	if !hasValidationField(err, "Service.name", "dox_kebab") {
		t.Fatalf("expected invalid service name validation error, got %v", err)
	}
}

func TestServicesMayShareInstanceID(t *testing.T) {
	services := []Service{
		{Namespace: "dox", Name: "iam", InstanceID: "server-pod-abc"},
		{Namespace: "dox", Name: "audit", InstanceID: "server-pod-abc"},
	}

	for _, service := range services {
		if err := service.Validate(); err != nil {
			t.Fatalf("validate service %+v: %v", service, err)
		}
	}
}

func TestServiceValidateRejectsInvalidInstanceID(t *testing.T) {
	service := Service{
		Namespace:  "dox",
		Name:       "iam",
		InstanceID: "server-pod-",
	}

	if err := service.Validate(); !hasValidationField(err, "Service.instance_id", "dox_identifier") {
		t.Fatalf("expected invalid instance_id validation error, got %v", err)
	}
}

func TestOrganizationValidateRejectsInvalidGovernanceIdentifier(t *testing.T) {
	organization := Organization{
		Name:       "opendox",
		Owner:      "Platform Team",
		CostCenter: "dox-core",
		Project:    "dox",
	}

	if err := organization.Validate(); !hasValidationField(err, "Organization.owner", "dox_identifier") {
		t.Fatalf("expected invalid owner validation error, got %v", err)
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
