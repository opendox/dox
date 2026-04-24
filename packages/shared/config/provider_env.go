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
 * @File    : provider_env.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package config

import (
	"context"
	"os"
	"sort"
	"strconv"
	"strings"
)

// EnvProvider reads environment variables by explicit prefix.
type EnvProvider struct {
	Lookup func() []string
}

// Read filters environment variables by source location and returns structured values.
func (p EnvProvider) Read(ctx context.Context, source Source) (*Payload, error) {
	if ctx == nil {
		return nil, ContractError("ctx", "context must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, SourceError("ctx", "context is done", err)
	}
	if err := validateProviderSource(source); err != nil {
		return nil, err
	}
	if source.Kind != ProviderKindEnv {
		return nil, ContractError("source.kind", "env provider requires env source kind")
	}

	var entries []string
	if p.Lookup != nil {
		entries = p.Lookup()
	} else {
		entries = os.Environ()
	}
	sort.Strings(entries)

	values := map[string]any{}
	prefix := source.Location
	for _, entry := range entries {
		key, value, ok := strings.Cut(entry, "=")
		if !ok || !strings.HasPrefix(key, prefix) {
			continue
		}
		normalized := normalizeEnvKey(strings.TrimPrefix(key, prefix))
		if normalized == "" {
			continue
		}
		values[normalized] = value
	}
	if len(values) == 0 {
		if !source.Required {
			return skippedPayload(source, "optional environment source did not match any variables"), nil
		}
		return nil, SourceError("source.location", "environment source did not match any variables", nil)
	}

	return &Payload{
		Source: source,
		Values: values,
		Metadata: map[string]string{
			"prefix": source.Location,
			"count":  strconv.Itoa(len(values)),
		},
		Diagnostic: SourceDiagnostic{
			Name:     source.Name,
			Kind:     source.Kind,
			Required: source.Required,
			Loaded:   true,
		},
	}, nil
}

func normalizeEnvKey(key string) string {
	key = strings.Trim(key, "_")
	if key == "" {
		return ""
	}
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "__", ".")
	key = strings.ReplaceAll(key, "_", ".")
	return key
}
