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
 * @File    : decode.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-25
 * @Modified: 2026-04-25
 */

package config

import (
	"context"

	"github.com/go-viper/mapstructure/v2"
)

// MapstructureDecoder decodes merged values into caller-owned targets.
type MapstructureDecoder struct{}

// Decode copies merged values into target according to decode options.
func (d MapstructureDecoder) Decode(ctx context.Context, values map[string]any, target any, options Options) error {
	if ctx == nil {
		return ContractError("ctx", "context must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return DecodeError("ctx", "context is done", err)
	}
	if err := validateTarget(target); err != nil {
		return err
	}
	normalizedOptions, err := normalizeOptions(options)
	if err != nil {
		return err
	}

	metadata := &mapstructure.Metadata{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		ErrorUnused:      normalizedOptions.UnknownKeyPolicy == UnknownKeyPolicyReject,
		Metadata:         metadata,
		Result:           target,
		TagName:          "mapstructure",
		WeaklyTypedInput: true,
		ZeroFields:       true,
	})
	if err != nil {
		return DecodeError("target", "decoder initialization failed", err)
	}
	if err := decoder.Decode(cloneStructuredMap(values)); err != nil {
		return DecodeError("target", "decode failed", err)
	}
	return nil
}

// DecodeValues decodes merged values with the built-in decoder.
func DecodeValues(ctx context.Context, values map[string]any, target any, options Options) error {
	return MapstructureDecoder{}.Decode(ctx, values, target, options)
}

// DecodeMergeResult decodes a merge result with the built-in decoder.
func DecodeMergeResult(ctx context.Context, result *MergeResult, target any, options Options) error {
	if result == nil {
		return ContractError("result", "merge result must not be nil")
	}
	return DecodeValues(ctx, result.Values, target, options)
}
