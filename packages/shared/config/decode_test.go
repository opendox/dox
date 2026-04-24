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
 * @File    : decode_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-25
 * @Modified: 2026-04-25
 */

package config

import (
	"context"
	"errors"
	"testing"
	"time"
)

type decodeSetting struct {
	App  decodeAppSetting  `mapstructure:"app"`
	HTTP decodeHTTPSetting `mapstructure:"http"`
}

type decodeAppSetting struct {
	Name string `mapstructure:"name"`
}

type decodeHTTPSetting struct {
	Port    int           `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
}

func TestDecodeValuesDecodesStructTarget(t *testing.T) {
	var target decodeSetting

	err := DecodeValues(context.Background(), map[string]any{
		"app": map[string]any{
			"name": "dox",
		},
		"http": map[string]any{
			"port":    "8080",
			"timeout": "5s",
		},
	}, &target, Options{UnknownKeyPolicy: UnknownKeyPolicyReject})
	if err != nil {
		t.Fatalf("decode values: %v", err)
	}

	if target.App.Name != "dox" {
		t.Fatalf("expected app.name to be decoded, got %q", target.App.Name)
	}
	if target.HTTP.Port != 8080 {
		t.Fatalf("expected weak string port decode, got %d", target.HTTP.Port)
	}
	if target.HTTP.Timeout != 5*time.Second {
		t.Fatalf("expected duration decode, got %s", target.HTTP.Timeout)
	}
}

func TestDecodeValuesDecodesMapTarget(t *testing.T) {
	var target map[string]any

	err := DecodeValues(context.Background(), map[string]any{
		"app": map[string]any{
			"name": "dox",
		},
	}, &target, Options{})
	if err != nil {
		t.Fatalf("decode map target: %v", err)
	}

	app := assertMap(t, target["app"])
	if got := app["name"]; got != "dox" {
		t.Fatalf("expected app.name to be dox, got %v", got)
	}
}

func TestDecodeMergeResultDecodesValues(t *testing.T) {
	var target decodeSetting

	err := DecodeMergeResult(context.Background(), &MergeResult{
		Values: map[string]any{
			"app": map[string]any{"name": "dox"},
		},
	}, &target, Options{UnknownKeyPolicy: UnknownKeyPolicyAllow})
	if err != nil {
		t.Fatalf("decode merge result: %v", err)
	}
	if target.App.Name != "dox" {
		t.Fatalf("expected app.name to be dox, got %q", target.App.Name)
	}
}

func TestDecodeValuesRejectsInvalidTarget(t *testing.T) {
	err := DecodeValues(context.Background(), map[string]any{}, decodeSetting{}, Options{})

	if !IsKind(err, ErrorKindContract) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestDecodeValuesRejectsUnknownKeysByDefault(t *testing.T) {
	var target decodeSetting

	err := DecodeValues(context.Background(), map[string]any{
		"app": map[string]any{
			"name":  "dox",
			"extra": true,
		},
	}, &target, Options{})

	if !IsKind(err, ErrorKindDecode) {
		t.Fatalf("expected decode error, got %v", err)
	}
}

func TestDecodeValuesAllowsUnknownKeysWhenExplicit(t *testing.T) {
	var target decodeSetting

	err := DecodeValues(context.Background(), map[string]any{
		"app": map[string]any{
			"name":  "dox",
			"extra": true,
		},
	}, &target, Options{UnknownKeyPolicy: UnknownKeyPolicyAllow})
	if err != nil {
		t.Fatalf("expected unknown key allow to decode, got %v", err)
	}
	if target.App.Name != "dox" {
		t.Fatalf("expected app.name to be decoded, got %q", target.App.Name)
	}
}

func TestDecodeValuesReturnsDecodeErrorForInvalidFieldType(t *testing.T) {
	var target decodeSetting

	err := DecodeValues(context.Background(), map[string]any{
		"http": map[string]any{
			"port": "not-a-number",
		},
	}, &target, Options{UnknownKeyPolicy: UnknownKeyPolicyAllow})

	if !IsKind(err, ErrorKindDecode) {
		t.Fatalf("expected decode error, got %v", err)
	}
}

func TestDecodeMergeResultRejectsNilResult(t *testing.T) {
	var target decodeSetting

	err := DecodeMergeResult(context.Background(), nil, &target, Options{})

	if !IsKind(err, ErrorKindContract) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestDecodeErrorKindHelpers(t *testing.T) {
	err := DecodeError("target", "failed", errors.New("boom"))

	if !IsKind(err, ErrorKindDecode) {
		t.Fatalf("expected decode error kind, got %v", err)
	}
	if !errors.Is(err, &Error{Kind: ErrorKindDecode}) {
		t.Fatal("expected errors.Is to match decode error kind")
	}
}
