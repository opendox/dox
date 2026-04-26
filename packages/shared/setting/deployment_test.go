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
 * @File    : deployment_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-25
 * @Modified: 2026-04-26
 */

package setting

import "testing"

func TestDeploymentAllowsEmptyIdentity(t *testing.T) {
	deployment := Deployment{}

	if err := deployment.Default(); err != nil {
		t.Fatalf("default deployment: %v", err)
	}
	if deployment.Env != EnvDev {
		t.Fatalf("expected default env dev, got %q", deployment.Env)
	}
	if err := deployment.Validate(); err != nil {
		t.Fatalf("validate empty deployment: %v", err)
	}
}

func TestDeploymentValidateAcceptsStableIdentifiers(t *testing.T) {
	deployment := Deployment{
		Env:          EnvProd,
		Region:       "us-east-1",
		Zone:         "us-east-1a",
		Cluster:      "dox-prod-1",
		K8sNamespace: "dox-prod",
	}

	if err := deployment.Validate(); err != nil {
		t.Fatalf("validate deployment: %v", err)
	}
}

func TestDeploymentValidateRejectsInvalidIdentifier(t *testing.T) {
	deployment := Deployment{
		Env:          EnvProd,
		Region:       "us-east-1",
		Zone:         "us-east-1a",
		Cluster:      "Dox Prod",
		K8sNamespace: "dox-prod",
	}

	if err := deployment.Validate(); !hasValidationField(err, "Deployment.cluster", "dox_identifier") {
		t.Fatalf("expected invalid cluster validation error, got %v", err)
	}
}

func TestDeploymentValidateRejectsTrailingNamespaceSeparator(t *testing.T) {
	deployment := Deployment{
		Env:          EnvProd,
		K8sNamespace: "dox-prod-",
	}

	if err := deployment.Validate(); !hasValidationField(err, "Deployment.k8s_namespace", "dox_identifier") {
		t.Fatalf("expected invalid k8s_namespace validation error, got %v", err)
	}
}
