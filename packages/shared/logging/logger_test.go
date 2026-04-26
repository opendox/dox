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
 * @File    : logger_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-27
 * @Modified: 2026-04-27
 */

package logging

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewLoggerWritesDoxModelFieldsAndContextCorrelation(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "service.jsonl")

	base, err := NewZapCoreBase(Config{
		Zap: ZapConfig{
			DisableCaller:       true,
			DisableStacktrace:   true,
			DisableErrorVerbose: true,
			ErrorOutputPaths:    []string{filepath.Join(tempDir, "errors.log")},
		},
		Cores: []CoreConfig{
			{
				Name:        "service-file",
				Enabled:     boolPtr(true),
				Type:        CoreTypeFile,
				Level:       LevelInfo,
				Encoding:    EncodingJSON,
				OutputPaths: []string{jsonPath},
				Datasets:    []string{"*"},
				Rotation: RotationConfig{
					Driver: RotationDriverNone,
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("build zap core base: %v", err)
	}
	t.Cleanup(base.Close)

	logger, err := NewLogger(
		base,
		ResourceAttr(Resource{
			ServiceNamespace:      "dox",
			ServiceName:           "iam",
			ServiceInstanceID:     "server-01",
			ServiceVersion:        "0.1.0",
			DeploymentEnvironment: "test",
			DoxRuntime:            "server",
		}),
		NodeAttr(Node{
			Component: "auth_service",
			Operation: "verify_credential",
		}),
		TagsAttr(Tags{
			"risk_level":   "medium",
			"login_method": "password",
		}),
	)
	if err != nil {
		t.Fatalf("build logger facade: %v", err)
	}

	ctx := ContextWithMergedCorrelation(nil, Correlation{
		TraceID:       "trace_ctx",
		RequestID:     "req_ctx",
		CorrelationID: "corr_ctx",
	})
	ctx = ContextWithMergedCorrelation(ctx, Correlation{
		SpanID: "span_ctx",
	})

	logger = logger.Named("iam.auth").With(EventAttr(Event{
		Dataset:  "dox.iam.security",
		Category: "authentication",
		Action:   "login",
	}))
	logger.Info(
		ctx,
		"login rejected",
		EventAttr(Event{
			Name:    "iam.login.rejected",
			Type:    "denied",
			Outcome: "failure",
		}),
		CorrelationAttr(Correlation{
			CorrelationID: "corr_call",
		}),
		FieldsAttr(Fields{
			"account": "alice@example.com",
		}),
		FieldAttr("tenant_id", "tenant_a"),
		ErrorAttr(errors.New("invalid password")),
	)
	if err := logger.Sync(); err != nil {
		t.Fatalf("sync logger: %v", err)
	}
	base.Close()

	entries := readJSONLines(t, jsonPath)
	if len(entries) != 1 {
		t.Fatalf("expected one log entry, got %#v", entries)
	}
	entry := entries[0]

	for key, expected := range map[string]string{
		"logger":                       "iam.auth",
		"message":                      "login rejected",
		FieldServiceNamespace:          "dox",
		FieldServiceName:               "iam",
		FieldServiceInstanceID:         "server-01",
		FieldServiceVersion:            "0.1.0",
		FieldDeploymentEnvironmentName: "test",
		FieldDoxRuntime:                "server",
		FieldTraceID:                   "trace_ctx",
		FieldSpanID:                    "span_ctx",
		FieldRequestID:                 "req_ctx",
		FieldCorrelationID:             "corr_call",
		FieldEventName:                 "iam.login.rejected",
		FieldEventDataset:              "dox.iam.security",
		FieldEventCategory:             "authentication",
		FieldEventType:                 "denied",
		FieldEventAction:               "login",
		FieldEventOutcome:              "failure",
		FieldComponent:                 "auth_service",
		FieldOperation:                 "verify_credential",
		"error":                        "invalid password",
	} {
		if entry[key] != expected {
			t.Fatalf("expected %s=%q, got %#v in %#v", key, expected, entry[key], entry)
		}
	}
	if _, exists := entry["errorVerbose"]; exists {
		t.Fatalf("expected disable_error_verbose to suppress errorVerbose, got %#v", entry)
	}

	tags, ok := entry[FieldTags].(map[string]any)
	if !ok {
		t.Fatalf("expected tags object, got %#v", entry[FieldTags])
	}
	if tags["risk_level"] != "medium" || tags["login_method"] != "password" {
		t.Fatalf("expected tags to map, got %#v", tags)
	}

	fields, ok := entry[FieldFields].(map[string]any)
	if !ok {
		t.Fatalf("expected fields object, got %#v", entry[FieldFields])
	}
	if fields["account"] != "alice@example.com" || fields["tenant_id"] != "tenant_a" {
		t.Fatalf("expected fields to map, got %#v", fields)
	}
}

func TestContextCorrelationHelpersStoreMergeAndRetrieve(t *testing.T) {
	if _, ok := CorrelationFromContext(nil); ok {
		t.Fatal("expected nil context to have no correlation")
	}

	ctx := ContextWithCorrelation(nil, Correlation{
		TraceID:   "trace_1",
		RequestID: "req_1",
	})
	correlation, ok := CorrelationFromContext(ctx)
	if !ok {
		t.Fatal("expected stored correlation")
	}
	if correlation.TraceID != "trace_1" || correlation.RequestID != "req_1" {
		t.Fatalf("expected stored correlation, got %+v", correlation)
	}

	ctx = ContextWithMergedCorrelation(ctx, Correlation{
		RequestID: "req_2",
		TaskID:    "task_1",
	})
	correlation, ok = CorrelationFromContext(ctx)
	if !ok {
		t.Fatal("expected merged correlation")
	}
	if correlation.TraceID != "trace_1" || correlation.RequestID != "req_2" || correlation.TaskID != "task_1" {
		t.Fatalf("expected merged correlation, got %+v", correlation)
	}

	merged := MergeCorrelation(Correlation{
		TraceID: "trace_base",
		JobID:   "job_base",
	}, Correlation{
		TraceID:     "trace_overlay",
		PluginRunID: "plugin_run_1",
	})
	if merged.TraceID != "trace_overlay" || merged.JobID != "job_base" || merged.PluginRunID != "plugin_run_1" {
		t.Fatalf("expected explicit merge semantics, got %+v", merged)
	}
}

func TestNewLoggerRejectsNilCoreBase(t *testing.T) {
	if _, err := NewLogger(nil); err == nil {
		t.Fatal("expected nil zap core base to be rejected")
	}
}

func TestLoggerFacadeSignaturesDoNotExposeZapTypes(t *testing.T) {
	loggerType := reflect.TypeOf((*Logger)(nil)).Elem()
	for methodIndex := 0; methodIndex < loggerType.NumMethod(); methodIndex++ {
		method := loggerType.Method(methodIndex)
		for inputIndex := 0; inputIndex < method.Type.NumIn(); inputIndex++ {
			assertNoZapType(t, method.Name, method.Type.In(inputIndex))
		}
		for outputIndex := 0; outputIndex < method.Type.NumOut(); outputIndex++ {
			assertNoZapType(t, method.Name, method.Type.Out(outputIndex))
		}
	}
}

func assertNoZapType(t *testing.T, method string, typ reflect.Type) {
	t.Helper()
	if typ == nil {
		return
	}
	if strings.HasPrefix(typ.PkgPath(), "go.uber.org/zap") {
		t.Fatalf("logger method %s exposes zap type %s", method, typ)
	}
	switch typ.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Array:
		assertNoZapType(t, method, typ.Elem())
	case reflect.Map:
		assertNoZapType(t, method, typ.Key())
		assertNoZapType(t, method, typ.Elem())
	case reflect.Func:
		for index := 0; index < typ.NumIn(); index++ {
			assertNoZapType(t, method, typ.In(index))
		}
		for index := 0; index < typ.NumOut(); index++ {
			assertNoZapType(t, method, typ.Out(index))
		}
	}
}
