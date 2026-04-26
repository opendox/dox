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
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package setting

import (
	"context"
	"testing"

	sharedconfig "github.com/opendox/dox/packages/shared/config"
	sharedsetting "github.com/opendox/dox/packages/shared/setting"
)

func TestSettingDefaultAppliesIdentityDefaults(t *testing.T) {
	setting := Setting{}

	if err := setting.DefaultWithOptions(DefaultOptions{Env: "test"}); err != nil {
		t.Fatalf("default setting: %v", err)
	}
	if setting.Identity.System.Runtime != ServerRuntime {
		t.Fatalf("expected server runtime, got %q", setting.Identity.System.Runtime)
	}
	if setting.Identity.Deployment.Env != sharedsetting.EnvTest {
		t.Fatalf("expected deployment env from options, got %q", setting.Identity.Deployment.Env)
	}
	if err := setting.Validate(); err != nil {
		t.Fatalf("validate setting: %v", err)
	}
}

func TestSettingDecodeValuesSupportsNestedIdentity(t *testing.T) {
	values := map[string]any{
		"identity": map[string]any{
			"organization": map[string]any{
				"name":        "opendox",
				"owner":       "platform",
				"cost_center": "dox-core",
			},
			"application": map[string]any{
				"name": "dox",
			},
			"service": map[string]any{
				"name":        "iam",
				"instance_id": "server-pod-1",
			},
			"deployment": map[string]any{
				"region":        "us-east-1",
				"zone":          "us-east-1a",
				"cluster":       "dox-prod-1",
				"k8s_namespace": "dox-prod",
			},
		},
	}

	setting := Setting{}
	if err := sharedconfig.DecodeValues(context.Background(), values, &setting, sharedconfig.Options{}); err != nil {
		t.Fatalf("decode setting: %v", err)
	}
	if err := setting.DefaultWithOptions(DefaultOptions{Env: "prod"}); err != nil {
		t.Fatalf("default setting: %v", err)
	}
	if err := setting.Validate(); err != nil {
		t.Fatalf("validate setting: %v", err)
	}

	if setting.Identity.System.Runtime != ServerRuntime {
		t.Fatalf("expected server runtime default, got %q", setting.Identity.System.Runtime)
	}
	if setting.Identity.Service.Namespace != "dox" {
		t.Fatalf("expected service namespace from application, got %q", setting.Identity.Service.Namespace)
	}
	if setting.Identity.Service.Name != "iam" {
		t.Fatalf("expected explicit service name, got %q", setting.Identity.Service.Name)
	}
	if setting.Identity.Deployment.Env != sharedsetting.EnvProd {
		t.Fatalf("expected deployment env from bootstrap option, got %q", setting.Identity.Deployment.Env)
	}
}

func TestSettingValidateRejectsInvalidDecodedIdentity(t *testing.T) {
	values := map[string]any{
		"identity": map[string]any{
			"system": map[string]any{
				"runtime": "server",
			},
			"service": map[string]any{
				"namespace": "dox",
				"name":      "iam-",
			},
		},
	}

	setting := Setting{}
	if err := sharedconfig.DecodeValues(context.Background(), values, &setting, sharedconfig.Options{}); err != nil {
		t.Fatalf("decode setting: %v", err)
	}
	if err := setting.Default(); err != nil {
		t.Fatalf("default setting: %v", err)
	}
	if err := setting.Validate(); !hasSharedValidationField(err, "Service.name", "dox_kebab") {
		t.Fatalf("expected invalid decoded service name, got %v", err)
	}
}

func TestSettingValidateRejectsInvalidBootstrapEnv(t *testing.T) {
	setting := Setting{}

	if err := setting.DefaultWithOptions(DefaultOptions{Env: "production"}); err != nil {
		t.Fatalf("default setting: %v", err)
	}
	if err := setting.Validate(); !hasSharedValidationField(err, "Deployment.env", "dox_env") {
		t.Fatalf("expected invalid bootstrap env validation error, got %v", err)
	}
}

func TestSettingDefaultRejectsNilReceiver(t *testing.T) {
	var setting *Setting

	if err := setting.Default(); err == nil {
		t.Fatal("expected nil setting default error")
	}
}
