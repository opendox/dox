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
 * @File    : provider_file.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package config

import (
	"context"
	"errors"
	"os"
)

// FileProvider reads local configuration files as raw bytes.
type FileProvider struct{}

// Read reads the file source location and returns raw bytes for later parsing.
func (p FileProvider) Read(ctx context.Context, source Source) (*Payload, error) {
	if ctx == nil {
		return nil, ContractError("ctx", "context must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, SourceError("ctx", "context is done", err)
	}
	if err := validateProviderSource(source); err != nil {
		return nil, err
	}
	if source.Kind != ProviderKindFile {
		return nil, ContractError("source.kind", "file provider requires file source kind")
	}

	info, err := os.Stat(source.Location)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && !source.Required {
			return skippedPayload(source, "optional file source does not exist"), nil
		}
		return nil, SourceError("source.location", "file source is not readable", err)
	}
	if info.IsDir() {
		return nil, SourceError("source.location", "file source points to a directory", nil)
	}

	body, err := os.ReadFile(source.Location)
	if err != nil {
		return nil, SourceError("source.location", "file source read failed", err)
	}

	return &Payload{
		Source:   source,
		Raw:      body,
		Metadata: sourceMetadata(source),
		Diagnostic: SourceDiagnostic{
			Name:     source.Name,
			Kind:     source.Kind,
			Required: source.Required,
			Loaded:   true,
		},
	}, nil
}

func skippedPayload(source Source, message string) *Payload {
	return &Payload{
		Source:   source,
		Metadata: sourceMetadata(source),
		Diagnostic: SourceDiagnostic{
			Name:     source.Name,
			Kind:     source.Kind,
			Required: source.Required,
			Skipped:  true,
			Message:  message,
		},
	}
}

func sourceMetadata(source Source) map[string]string {
	switch source.Kind {
	case ProviderKindFile:
		return map[string]string{"path": source.Location}
	case ProviderKindEnv:
		return map[string]string{"prefix": source.Location, "count": "0"}
	default:
		return map[string]string{"location": source.Location}
	}
}
