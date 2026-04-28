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
 * @File    : setting.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-28
 * @Modified: 2026-04-28
 */

package bootstrap

import (
	"context"
	"fmt"

	sharedconfig "github.com/opendox/dox/packages/shared/config"
	serversetting "github.com/opendox/dox/server/internal/setting"
)

// SettingSnapshot captures typed, defaulted, and validated server settings.
type SettingSnapshot struct {
	Runtime     string
	Env         string
	Setting     serversetting.Setting
	SourceNames []string
	Fingerprint string
	Diagnostics sharedconfig.Diagnostics
}

// LoadSetting loads server configuration sources and assembles typed settings.
func LoadSetting(ctx context.Context, options ConfigOptions) (*SettingSnapshot, error) {
	configSnapshot, err := LoadConfig(ctx, options)
	if err != nil {
		return nil, err
	}
	return AssembleSetting(ctx, configSnapshot)
}

// AssembleSetting decodes, defaults, and validates settings from a config snapshot.
func AssembleSetting(ctx context.Context, snapshot *ConfigSnapshot) (*SettingSnapshot, error) {
	if snapshot == nil {
		return nil, sharedconfig.ContractError("snapshot", "config snapshot must not be nil")
	}

	setting := serversetting.Setting{}
	if err := sharedconfig.DecodeValues(ctx, snapshot.Values, &setting, sharedconfig.Options{
		UnknownKeyPolicy: sharedconfig.UnknownKeyPolicyReject,
	}); err != nil {
		return nil, err
	}
	if err := setting.DefaultWithOptions(serversetting.DefaultOptions{Env: snapshot.Env}); err != nil {
		return nil, fmt.Errorf("server setting default failed: %w", err)
	}
	if err := setting.Validate(); err != nil {
		return nil, fmt.Errorf("server setting validation failed: %w", err)
	}

	return &SettingSnapshot{
		Runtime:     snapshot.Runtime,
		Env:         snapshot.Env,
		Setting:     setting,
		SourceNames: append([]string(nil), snapshot.SourceNames...),
		Fingerprint: snapshot.Fingerprint,
		Diagnostics: cloneConfigDiagnostics(snapshot.Diagnostics),
	}, nil
}

func cloneConfigDiagnostics(input sharedconfig.Diagnostics) sharedconfig.Diagnostics {
	return sharedconfig.Diagnostics{
		Sources:   append([]sharedconfig.SourceDiagnostic(nil), input.Sources...),
		Overrides: append([]sharedconfig.MergeOverride(nil), input.Overrides...),
	}
}
