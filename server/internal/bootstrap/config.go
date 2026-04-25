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
 * @File    : config.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-25
 * @Modified: 2026-04-25
 */

package bootstrap

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	sharedconfig "github.com/opendox/dox/packages/shared/config"
)

const (
	// DefaultConfigDir is the default directory for server configuration files.
	DefaultConfigDir = "configs"
	// DefaultConfigEnv is the default server runtime environment name.
	DefaultConfigEnv = "dev"
	// DefaultConfigFormat is the default server configuration file format.
	DefaultConfigFormat = "yaml"
	// DefaultConfigEnvPrefix is the default prefix for server configuration value overrides.
	DefaultConfigEnvPrefix = "DOX_SERVER_"

	configRuntime = "server"
)

const (
	configBasePriority        = 10
	configEnvironmentPriority = 20
	configLocalPriority       = 30
	configEnvPriority         = 100
)

// ConfigOptions defines startup inputs for server configuration loading.
type ConfigOptions struct {
	ConfigDir string
	Env       string
	Format    string
	EnvPrefix string
	Timeout   time.Duration
}

// ConfigSnapshot captures the merged server configuration values loaded during bootstrap.
type ConfigSnapshot struct {
	Runtime     string
	Env         string
	Values      map[string]any
	SourceNames []string
	Fingerprint string
	Diagnostics sharedconfig.Diagnostics
}

// LoadConfig loads the server startup configuration snapshot.
//
// Source order is:
//   - configs/base.<format>, required baseline
//   - configs/<env>.<format>, optional environment override
//   - configs/local.<format>, optional local override
//   - environment variables with EnvPrefix, optional final override
//
// The snapshot intentionally uses a generic map target. Concrete resource
// settings belong to later server-owned setting packages.
func LoadConfig(ctx context.Context, options ConfigOptions) (*ConfigSnapshot, error) {
	normalized, err := normalizeConfigOptions(options)
	if err != nil {
		return nil, err
	}

	values := map[string]any{}
	result, err := sharedconfig.Load(ctx, sharedconfig.Request{
		Runtime: configRuntime,
		Env:     normalized.Env,
		Target:  &values,
		Sources: buildConfigSources(normalized),
		Options: sharedconfig.Options{
			Timeout:          normalized.Timeout,
			UnknownKeyPolicy: sharedconfig.UnknownKeyPolicyAllow,
		},
	})
	if err != nil {
		return nil, err
	}

	return &ConfigSnapshot{
		Runtime:     result.Runtime,
		Env:         result.Env,
		Values:      cloneConfigValues(values),
		SourceNames: append([]string(nil), result.SourceNames...),
		Fingerprint: result.Fingerprint,
		Diagnostics: result.Diagnostics,
	}, nil
}

func normalizeConfigOptions(options ConfigOptions) (ConfigOptions, error) {
	options.ConfigDir = strings.TrimSpace(options.ConfigDir)
	if options.ConfigDir == "" {
		options.ConfigDir = DefaultConfigDir
	}

	options.Env = strings.TrimSpace(options.Env)
	if options.Env == "" {
		options.Env = DefaultConfigEnv
	}

	format := strings.TrimSpace(options.Format)
	format = strings.TrimPrefix(format, ".")
	format = strings.ToLower(format)
	if format == "" {
		format = DefaultConfigFormat
	}
	if _, err := parserForFormat(format); err != nil {
		return options, err
	}
	options.Format = format

	options.EnvPrefix = strings.TrimSpace(options.EnvPrefix)
	if options.EnvPrefix == "" {
		options.EnvPrefix = DefaultConfigEnvPrefix
	}

	return options, nil
}

func buildConfigSources(options ConfigOptions) []sharedconfig.Source {
	parser, _ := parserForFormat(options.Format)
	return []sharedconfig.Source{
		{
			Name:     "base",
			Kind:     sharedconfig.ProviderKindFile,
			Parser:   parser,
			Location: filepath.Join(options.ConfigDir, "base."+options.Format),
			Required: true,
			Priority: configBasePriority,
		},
		{
			Name:     "environment",
			Kind:     sharedconfig.ProviderKindFile,
			Parser:   parser,
			Location: filepath.Join(options.ConfigDir, options.Env+"."+options.Format),
			Required: false,
			Priority: configEnvironmentPriority,
		},
		{
			Name:     "local",
			Kind:     sharedconfig.ProviderKindFile,
			Parser:   parser,
			Location: filepath.Join(options.ConfigDir, "local."+options.Format),
			Required: false,
			Priority: configLocalPriority,
		},
		{
			Name:     "env",
			Kind:     sharedconfig.ProviderKindEnv,
			Parser:   sharedconfig.ParserKindNone,
			Location: options.EnvPrefix,
			Required: false,
			Priority: configEnvPriority,
		},
	}
}

func parserForFormat(format string) (sharedconfig.ParserKind, error) {
	switch format {
	case "yaml", "yml":
		return sharedconfig.ParserKindYAML, nil
	case "json":
		return sharedconfig.ParserKindJSON, nil
	case "toml":
		return sharedconfig.ParserKindTOML, nil
	default:
		return "", sharedconfig.ContractError("options.format", "config format is not supported")
	}
}

func cloneConfigValues(input map[string]any) map[string]any {
	if len(input) == 0 {
		return map[string]any{}
	}
	values := make(map[string]any, len(input))
	for key, value := range input {
		values[key] = cloneConfigValue(value)
	}
	return values
}

func cloneConfigValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return cloneConfigValues(typed)
	case []any:
		values := make([]any, len(typed))
		for index, entry := range typed {
			values[index] = cloneConfigValue(entry)
		}
		return values
	default:
		return value
	}
}
