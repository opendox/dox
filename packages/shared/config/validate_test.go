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
 * @File    : validate_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package config

import (
	"context"
	"errors"
	"testing"
)

type testSetting struct {
	Name string
}

func TestValidateLoadRequestAcceptsValidContract(t *testing.T) {
	var target testSetting

	err := ValidateLoadRequest(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{
			{
				Name:     "base",
				Kind:     ProviderKindFile,
				Parser:   ParserKindYAML,
				Location: "configs/base.yaml",
				Required: true,
				Priority: 10,
			},
			{
				Name:     "env",
				Kind:     ProviderKindEnv,
				Parser:   ParserKindNone,
				Location: "DOX_SERVER_",
				Required: false,
				Priority: 100,
			},
		},
		Options: Options{
			MergeStrategy:    MergeStrategyDeepReplace,
			UnknownKeyPolicy: UnknownKeyPolicyReject,
		},
	})
	if err != nil {
		t.Fatalf("expected valid request, got %v", err)
	}
}

func TestValidateLoadRequestRejectsNilContext(t *testing.T) {
	var target testSetting

	err := ValidateLoadRequest(nil, Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{{Name: "base", Kind: ProviderKindFile, Parser: ParserKindYAML, Location: "configs/base.yaml"}},
	})

	assertContractError(t, err)
}

func TestValidateLoadRequestRejectsInvalidTarget(t *testing.T) {
	err := ValidateLoadRequest(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  testSetting{},
		Sources: []Source{{Name: "base", Kind: ProviderKindFile, Parser: ParserKindYAML, Location: "configs/base.yaml"}},
	})

	assertContractError(t, err)
}

func TestValidateLoadRequestRejectsEmptySourcesByDefault(t *testing.T) {
	var target testSetting

	err := ValidateLoadRequest(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
	})

	assertContractError(t, err)
}

func TestValidateLoadRequestAllowsEmptySourcesWhenExplicit(t *testing.T) {
	var target testSetting

	err := ValidateLoadRequest(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Options: Options{AllowEmptySources: true},
	})
	if err != nil {
		t.Fatalf("expected empty sources to be allowed, got %v", err)
	}
}

func TestValidateLoadRequestRejectsDuplicateSourceNames(t *testing.T) {
	var target testSetting

	err := ValidateLoadRequest(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{
			{Name: "base", Kind: ProviderKindFile, Parser: ParserKindYAML, Location: "configs/base.yaml"},
			{Name: "base", Kind: ProviderKindFile, Parser: ParserKindYAML, Location: "configs/local.yaml", Priority: 1},
		},
	})

	assertContractError(t, err)
}

func TestValidateLoadRequestRejectsInvalidProviderParserPair(t *testing.T) {
	var target testSetting

	err := ValidateLoadRequest(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{
			{Name: "env", Kind: ProviderKindEnv, Parser: ParserKindYAML, Location: "DOX_SERVER_"},
		},
	})

	assertContractError(t, err)
}

func TestValidateLoadRequestAllowsCustomProviderAndParserKinds(t *testing.T) {
	var target testSetting

	err := ValidateLoadRequest(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{
			{Name: "remote-main", Kind: ProviderKind("consul"), Parser: ParserKind("jsonnet"), Location: "dox/server/dev"},
		},
	})
	if err != nil {
		t.Fatalf("expected custom provider and parser kinds to be accepted, got %v", err)
	}
}

func TestValidateLoadRequestRejectsDuplicatePriorities(t *testing.T) {
	var target testSetting

	err := ValidateLoadRequest(context.Background(), Request{
		Runtime: "server",
		Env:     "dev",
		Target:  &target,
		Sources: []Source{
			{Name: "base", Kind: ProviderKindFile, Parser: ParserKindYAML, Location: "configs/base.yaml", Priority: 10},
			{Name: "local", Kind: ProviderKindFile, Parser: ParserKindYAML, Location: "configs/local.yaml", Priority: 10},
		},
	})

	assertContractError(t, err)
}

func TestErrorKindHelpers(t *testing.T) {
	err := SourceError("source", "failed", errors.New("boom"))

	if !IsKind(err, ErrorKindSource) {
		t.Fatalf("expected source error kind, got %v", err)
	}
	if IsKind(err, ErrorKindContract) {
		t.Fatal("did not expect contract error kind")
	}
	if !errors.Is(err, &Error{Kind: ErrorKindSource}) {
		t.Fatal("expected errors.Is to match source error kind")
	}
}

func assertContractError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected contract error, got nil")
	}
	if !IsKind(err, ErrorKindContract) {
		t.Fatalf("expected contract error, got %v", err)
	}
}
