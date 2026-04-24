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
 * @File    : parser_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package config

import (
	"context"
	"encoding/json"
	"testing"
)

func TestYAMLParserParsesObjectPayload(t *testing.T) {
	values, err := YAMLParser{}.Parse(context.Background(), filePayload(ParserKindYAML, []byte(`
app:
  name: dox
http:
  address: 127.0.0.1:8080
features:
  - config
  - parser
`)))
	if err != nil {
		t.Fatalf("parse yaml payload: %v", err)
	}

	app := assertMap(t, values["app"])
	if got := app["name"]; got != "dox" {
		t.Fatalf("expected app.name to be dox, got %v", got)
	}
	http := assertMap(t, values["http"])
	if got := http["address"]; got != "127.0.0.1:8080" {
		t.Fatalf("expected http.address, got %v", got)
	}
}

func TestJSONParserParsesObjectPayload(t *testing.T) {
	values, err := JSONParser{}.Parse(context.Background(), filePayload(ParserKindJSON, []byte(`{
	"app": {"name": "dox"},
	"http": {"port": 8080}
}`)))
	if err != nil {
		t.Fatalf("parse json payload: %v", err)
	}

	app := assertMap(t, values["app"])
	if got := app["name"]; got != "dox" {
		t.Fatalf("expected app.name to be dox, got %v", got)
	}
	http := assertMap(t, values["http"])
	if got := http["port"]; got != json.Number("8080") {
		t.Fatalf("expected http.port to preserve json number, got %T %v", got, got)
	}
}

func TestTOMLParserParsesObjectPayload(t *testing.T) {
	values, err := TOMLParser{}.Parse(context.Background(), filePayload(ParserKindTOML, []byte(`
[app]
name = "dox"

[http]
port = 8080
`)))
	if err != nil {
		t.Fatalf("parse toml payload: %v", err)
	}

	app := assertMap(t, values["app"])
	if got := app["name"]; got != "dox" {
		t.Fatalf("expected app.name to be dox, got %v", got)
	}
	http := assertMap(t, values["http"])
	if got := http["port"]; got != int64(8080) {
		t.Fatalf("expected http.port to be int64, got %T %v", got, got)
	}
}

func TestNoneParserReturnsProviderValues(t *testing.T) {
	payload := Payload{
		Source: Source{
			Name:     "env",
			Kind:     ProviderKindEnv,
			Parser:   ParserKindNone,
			Location: "DOX_SERVER_",
		},
		Values: map[string]any{"app.name": "dox"},
		Diagnostic: SourceDiagnostic{
			Name:   "env",
			Kind:   ProviderKindEnv,
			Loaded: true,
		},
	}

	values, err := NoneParser{}.Parse(context.Background(), payload)
	if err != nil {
		t.Fatalf("parse none payload: %v", err)
	}
	values["app.name"] = "changed"
	if got := payload.Values["app.name"]; got != "dox" {
		t.Fatalf("expected parser result to be a shallow copy, original got %v", got)
	}
}

func TestParsePayloadUsesDeclaredBuiltinParser(t *testing.T) {
	values, err := ParsePayload(context.Background(), filePayload(ParserKindJSON, []byte(`{"app":{"name":"dox"}}`)))
	if err != nil {
		t.Fatalf("parse payload: %v", err)
	}

	app := assertMap(t, values["app"])
	if got := app["name"]; got != "dox" {
		t.Fatalf("expected app.name to be dox, got %v", got)
	}
}

func TestParsersRejectMalformedInput(t *testing.T) {
	tests := []struct {
		name    string
		parser  Parser
		payload Payload
	}{
		{
			name:    "yaml",
			parser:  YAMLParser{},
			payload: filePayload(ParserKindYAML, []byte("app: [")),
		},
		{
			name:    "json",
			parser:  JSONParser{},
			payload: filePayload(ParserKindJSON, []byte(`{"app":`)),
		},
		{
			name:    "toml",
			parser:  TOMLParser{},
			payload: filePayload(ParserKindTOML, []byte("[app")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.parser.Parse(context.Background(), tt.payload)
			if !IsKind(err, ErrorKindParse) {
				t.Fatalf("expected parse error, got %v", err)
			}
		})
	}
}

func TestParserRejectsNonObjectPayload(t *testing.T) {
	_, err := JSONParser{}.Parse(context.Background(), filePayload(ParserKindJSON, []byte(`["dox"]`)))

	if !IsKind(err, ErrorKindParse) {
		t.Fatalf("expected parse error, got %v", err)
	}
}

func TestParserRejectsMismatchedKind(t *testing.T) {
	_, err := JSONParser{}.Parse(context.Background(), filePayload(ParserKindYAML, []byte(`{"app":{"name":"dox"}}`)))

	if !IsKind(err, ErrorKindContract) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestParsePayloadRejectsUnsupportedParserKind(t *testing.T) {
	_, err := ParsePayload(context.Background(), filePayload(ParserKind("jsonnet"), []byte(`{}`)))

	if !IsKind(err, ErrorKindContract) {
		t.Fatalf("expected contract error, got %v", err)
	}
}

func TestParserSkipsOptionalPayload(t *testing.T) {
	payload := filePayload(ParserKindYAML, nil)
	payload.Diagnostic = SourceDiagnostic{
		Name:    "local",
		Kind:    ProviderKindFile,
		Skipped: true,
	}

	values, err := YAMLParser{}.Parse(context.Background(), payload)
	if err != nil {
		t.Fatalf("expected skipped payload to parse as empty values, got %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("expected empty skipped values, got %+v", values)
	}
}

func filePayload(parser ParserKind, raw []byte) Payload {
	return Payload{
		Source: Source{
			Name:     "base",
			Kind:     ProviderKindFile,
			Parser:   parser,
			Location: "configs/base",
			Required: true,
		},
		Raw: raw,
		Diagnostic: SourceDiagnostic{
			Name:     "base",
			Kind:     ProviderKindFile,
			Required: true,
			Loaded:   true,
		},
	}
}

func assertMap(t *testing.T, value any) map[string]any {
	t.Helper()
	values, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", value)
	}
	return values
}
