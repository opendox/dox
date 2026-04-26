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
 * @File    : config_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package logging

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	sharedconfig "github.com/opendox/dox/packages/shared/config"
	"go.yaml.in/yaml/v3"
)

func TestConfigDefaultAppliesSharedContract(t *testing.T) {
	cfg := Config{}

	if err := cfg.Default(); err != nil {
		t.Fatalf("default config: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("validate default config: %v", err)
	}

	if cfg.Level != LevelInfo {
		t.Fatalf("expected root level info, got %q", cfg.Level)
	}
	if cfg.Zap.Encoding != EncodingJSON {
		t.Fatalf("expected zap encoding json, got %q", cfg.Zap.Encoding)
	}
	if cfg.Zap.EncoderConfig.MessageKey != "message" {
		t.Fatalf("expected message key default, got %q", cfg.Zap.EncoderConfig.MessageKey)
	}
	if len(cfg.Cores) != 2 {
		t.Fatalf("expected two default cores, got %d", len(cfg.Cores))
	}
	assertCore(t, cfg.Cores[0], "console", CoreTypeConsole, EncodingConsole, []string{"stdout"})
	assertCore(t, cfg.Cores[1], "service-file", CoreTypeFile, EncodingJSON, []string{DefaultFilePathTemplate})
	if cfg.Cores[1].Rotation.Driver != RotationDriverLumberjack {
		t.Fatalf("expected lumberjack rotation, got %q", cfg.Cores[1].Rotation.Driver)
	}
	if cfg.Buffering.FlushInterval != time.Second {
		t.Fatalf("expected buffering flush interval 1s, got %s", cfg.Buffering.FlushInterval)
	}
	if cfg.Shutdown.Timeout != 5*time.Second {
		t.Fatalf("expected shutdown timeout 5s, got %s", cfg.Shutdown.Timeout)
	}
	if !boolValue(cfg.Redaction.Enabled) || cfg.Redaction.Replacement != DefaultRedactionReplacement {
		t.Fatalf("expected redaction defaults, got %+v", cfg.Redaction)
	}
	if !contains(cfg.Redaction.Keys, "authorization") || !contains(cfg.Redaction.Keys, "client_secret") {
		t.Fatalf("expected sensitive redaction keys, got %#v", cfg.Redaction.Keys)
	}
	if !boolValue(cfg.OTel.Propagation.TraceContext) || !boolValue(cfg.OTel.Propagation.Baggage) {
		t.Fatalf("expected OpenTelemetry propagation defaults, got %+v", cfg.OTel.Propagation)
	}
	if cfg.OTel.Exporter.OTLP.Enabled {
		t.Fatal("expected OTLP exporter to be disabled by default")
	}
}

func TestFieldNameConstants(t *testing.T) {
	assertEqual(t, FieldServiceName, "service.name")
	assertEqual(t, FieldServiceInstanceID, "service.instance.id")
	assertEqual(t, FieldDeploymentEnvironmentName, "deployment.environment.name")
	assertEqual(t, FieldTraceID, "trace_id")
	assertEqual(t, FieldCorrelationID, "correlation_id")
	assertEqual(t, FieldEventName, "event.name")
	assertEqual(t, FieldEventDataset, "event.dataset")
	assertEqual(t, FieldComponent, "component")
	assertEqual(t, FieldOperation, "operation")
	assertEqual(t, FieldTags, "tags")
	assertEqual(t, FieldFields, "fields")
}

func TestConfigDefaultPreservesExplicitDisabledToggles(t *testing.T) {
	cfg := Config{
		Cores: []CoreConfig{
			{
				Name:        "console",
				Enabled:     boolPtr(false),
				Type:        CoreTypeConsole,
				Level:       LevelInfo,
				Encoding:    EncodingConsole,
				OutputPaths: []string{"stdout"},
				Datasets:    []string{"*"},
			},
			{
				Name:        "service-file",
				Enabled:     boolPtr(false),
				Type:        CoreTypeFile,
				Level:       LevelInfo,
				Encoding:    EncodingJSON,
				OutputPaths: []string{DefaultFilePathTemplate},
				Datasets:    []string{"*"},
				Rotation: RotationConfig{
					Enabled:   boolPtr(false),
					Compress:  boolPtr(false),
					LocalTime: boolPtr(false),
				},
			},
		},
		Buffering: BufferingConfig{
			Enabled: boolPtr(false),
		},
		Redaction: RedactionConfig{
			Enabled: boolPtr(false),
		},
		OTel: OpenTelemetryConfig{
			Enabled: boolPtr(false),
			Propagation: OpenTelemetryPropagation{
				TraceContext: boolPtr(false),
				Baggage:      boolPtr(false),
			},
			Traces: OpenTelemetryTraces{
				Enabled: boolPtr(false),
			},
			Metrics: OpenTelemetrySignal{
				Enabled: boolPtr(false),
			},
			Logs: OpenTelemetrySignal{
				Enabled: boolPtr(false),
			},
			Exporter: OpenTelemetryExporter{
				OTLP: OTLPExporterConfig{
					Insecure: boolPtr(false),
				},
			},
		},
	}

	if err := cfg.Default(); err != nil {
		t.Fatalf("default config: %v", err)
	}

	if boolValue(cfg.Cores[0].Enabled) || boolValue(cfg.Cores[1].Enabled) {
		t.Fatalf("expected explicit disabled cores to stay disabled, got %+v", cfg.Cores)
	}
	if boolValue(cfg.Cores[1].Rotation.Enabled) {
		t.Fatalf("expected explicit disabled rotation to stay disabled, got %+v", cfg.Cores[1].Rotation)
	}
	if boolValue(cfg.Cores[1].Rotation.Compress) || boolValue(cfg.Cores[1].Rotation.LocalTime) {
		t.Fatalf("expected explicit disabled rotation options to stay disabled, got %+v", cfg.Cores[1].Rotation)
	}
	if boolValue(cfg.Buffering.Enabled) {
		t.Fatalf("expected explicit disabled buffering to stay disabled, got %+v", cfg.Buffering)
	}
	if boolValue(cfg.Redaction.Enabled) {
		t.Fatalf("expected explicit disabled redaction to stay disabled, got %+v", cfg.Redaction)
	}
	if boolValue(cfg.OTel.Enabled) ||
		boolValue(cfg.OTel.Propagation.TraceContext) ||
		boolValue(cfg.OTel.Propagation.Baggage) ||
		boolValue(cfg.OTel.Traces.Enabled) ||
		boolValue(cfg.OTel.Metrics.Enabled) ||
		boolValue(cfg.OTel.Logs.Enabled) ||
		boolValue(cfg.OTel.Exporter.OTLP.Insecure) {
		t.Fatalf("expected explicit disabled OpenTelemetry toggles to stay disabled, got %+v", cfg.OTel)
	}
}

func TestConfigValidateRejectsInvalidValues(t *testing.T) {
	cfg := Config{}
	if err := cfg.Default(); err != nil {
		t.Fatalf("default config: %v", err)
	}

	cfg.Level = Level("verbose")
	cfg.Zap.Encoding = Encoding("plain")
	cfg.Zap.EncoderConfig.TimeEncoder = "clock"
	cfg.Cores[0].Name = cfg.Cores[1].Name
	cfg.Cores[0].Datasets = []string{""}
	cfg.Cores[1].Rotation.MaxSizeMB = 0
	cfg.Buffering.SizeBytes = 0
	cfg.Shutdown.Timeout = 0
	cfg.Redaction.Keys = []string{""}
	cfg.OTel.Exporter.OTLP.Enabled = true
	cfg.OTel.Exporter.OTLP.Endpoint = ""
	cfg.OTel.Traces.Sampler.Ratio = 2

	err := cfg.Validate()
	for _, field := range []string{
		"level",
		"zap.encoding",
		"zap.encoder_config.time_encoder",
		"cores[1].name",
		"cores[0].datasets[0]",
		"cores[1].rotation.max_size_mb",
		"buffering.size_bytes",
		"shutdown.timeout",
		"redaction.keys[0]",
		"otel.exporter.otlp.endpoint",
		"otel.traces.sampler.ratio",
	} {
		if !hasValidationField(err, field) {
			t.Fatalf("expected validation field %s, got %v", field, err)
		}
	}
}

func TestConfigDecodesWithMapstructure(t *testing.T) {
	var cfg Config

	err := sharedconfig.DecodeValues(context.Background(), map[string]any{
		"level": "debug",
		"zap": map[string]any{
			"level":    "warn",
			"encoding": "json",
			"encoder_config": map[string]any{
				"message_key":      "msg",
				"level_encoder":    "lowercase",
				"time_encoder":     "iso8601",
				"duration_encoder": "seconds",
				"caller_encoder":   "full",
				"name_encoder":     "full",
			},
			"output_paths":       []any{"stdout"},
			"error_output_paths": []any{"stderr"},
		},
		"cores": []any{
			map[string]any{
				"name":         "console",
				"enabled":      true,
				"type":         "console",
				"level":        "debug",
				"encoding":     "console",
				"output_paths": []any{"stdout"},
				"datasets":     []any{"*"},
			},
		},
		"buffering": map[string]any{
			"enabled":        true,
			"size_bytes":     4096,
			"flush_interval": "2s",
		},
		"shutdown": map[string]any{
			"timeout": "3s",
		},
		"redaction": map[string]any{
			"enabled":     true,
			"replacement": "[redacted]",
			"keys":        []any{"token"},
		},
		"otel": map[string]any{
			"enabled": true,
			"propagation": map[string]any{
				"trace_context": true,
				"baggage":       true,
			},
			"traces": map[string]any{
				"enabled": true,
				"sampler": map[string]any{
					"type":  "traceidratio",
					"ratio": 0.5,
				},
			},
			"metrics": map[string]any{"enabled": true},
			"logs":    map[string]any{"enabled": true},
			"exporter": map[string]any{
				"otlp": map[string]any{
					"enabled":  false,
					"endpoint": "collector:4317",
					"protocol": "grpc",
					"timeout":  "5s",
				},
			},
			"batch": map[string]any{
				"max_queue_size":        64,
				"schedule_delay":        "1s",
				"export_timeout":        "2s",
				"max_export_batch_size": 16,
			},
		},
	}, &cfg, sharedconfig.Options{})
	if err != nil {
		t.Fatalf("decode logging config: %v", err)
	}
	if cfg.Level != LevelDebug {
		t.Fatalf("expected decoded level debug, got %q", cfg.Level)
	}
	if cfg.Zap.EncoderConfig.MessageKey != "msg" {
		t.Fatalf("expected decoded message key msg, got %q", cfg.Zap.EncoderConfig.MessageKey)
	}
	if cfg.Buffering.FlushInterval != 2*time.Second {
		t.Fatalf("expected decoded buffering duration, got %s", cfg.Buffering.FlushInterval)
	}
	if cfg.OTel.Traces.Sampler.Ratio != 0.5 {
		t.Fatalf("expected decoded sampler ratio 0.5, got %f", cfg.OTel.Traces.Sampler.Ratio)
	}
}

func TestConfigDecodePreservesExplicitDisabledToggles(t *testing.T) {
	var cfg Config

	err := sharedconfig.DecodeValues(context.Background(), map[string]any{
		"cores": []any{
			map[string]any{
				"name":         "console",
				"enabled":      false,
				"type":         "console",
				"level":        "info",
				"encoding":     "console",
				"output_paths": []any{"stdout"},
				"datasets":     []any{"*"},
			},
			map[string]any{
				"name":         "service-file",
				"enabled":      false,
				"type":         "file",
				"level":        "info",
				"encoding":     "json",
				"output_paths": []any{DefaultFilePathTemplate},
				"datasets":     []any{"*"},
				"rotation": map[string]any{
					"enabled":    false,
					"compress":   false,
					"local_time": false,
				},
			},
		},
		"buffering": map[string]any{
			"enabled": false,
		},
		"redaction": map[string]any{
			"enabled": false,
		},
		"otel": map[string]any{
			"enabled": false,
			"propagation": map[string]any{
				"trace_context": false,
				"baggage":       false,
			},
			"traces": map[string]any{
				"enabled": false,
			},
			"metrics": map[string]any{
				"enabled": false,
			},
			"logs": map[string]any{
				"enabled": false,
			},
			"exporter": map[string]any{
				"otlp": map[string]any{
					"insecure": false,
				},
			},
		},
	}, &cfg, sharedconfig.Options{})
	if err != nil {
		t.Fatalf("decode logging config: %v", err)
	}
	if err := cfg.Default(); err != nil {
		t.Fatalf("default logging config: %v", err)
	}

	if boolValue(cfg.Cores[0].Enabled) || boolValue(cfg.Cores[1].Enabled) {
		t.Fatalf("expected decoded disabled cores to stay disabled, got %+v", cfg.Cores)
	}
	if boolValue(cfg.Cores[1].Rotation.Enabled) ||
		boolValue(cfg.Cores[1].Rotation.Compress) ||
		boolValue(cfg.Cores[1].Rotation.LocalTime) {
		t.Fatalf("expected decoded disabled rotation settings to stay disabled, got %+v", cfg.Cores[1].Rotation)
	}
	if boolValue(cfg.Buffering.Enabled) || boolValue(cfg.Redaction.Enabled) {
		t.Fatalf("expected decoded disabled buffering/redaction to stay disabled, got buffering=%+v redaction=%+v", cfg.Buffering, cfg.Redaction)
	}
	if boolValue(cfg.OTel.Enabled) ||
		boolValue(cfg.OTel.Propagation.TraceContext) ||
		boolValue(cfg.OTel.Propagation.Baggage) ||
		boolValue(cfg.OTel.Traces.Enabled) ||
		boolValue(cfg.OTel.Metrics.Enabled) ||
		boolValue(cfg.OTel.Logs.Enabled) ||
		boolValue(cfg.OTel.Exporter.OTLP.Insecure) {
		t.Fatalf("expected decoded disabled OpenTelemetry toggles to stay disabled, got %+v", cfg.OTel)
	}
}

func TestConfigJSONAndYAMLTags(t *testing.T) {
	cfg := Config{}
	if err := cfg.Default(); err != nil {
		t.Fatalf("default config: %v", err)
	}

	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	jsonText := string(jsonBytes)
	for _, item := range []string{`"encoder_config"`, `"output_paths"`, `"error_output_paths"`, `"max_size_mb"`} {
		if !strings.Contains(jsonText, item) {
			t.Fatalf("expected JSON to contain %s, got %s", item, jsonText)
		}
	}

	yamlBytes, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal yaml: %v", err)
	}
	yamlText := string(yamlBytes)
	for _, item := range []string{"encoder_config:", "output_paths:", "error_output_paths:", "max_size_mb:"} {
		if !strings.Contains(yamlText, item) {
			t.Fatalf("expected YAML to contain %s, got %s", item, yamlText)
		}
	}
}

func TestModelTypesRepresentObservabilityEvents(t *testing.T) {
	event := Event{
		Name:     "iam.login.rejected",
		Dataset:  "dox.iam.security",
		Category: "authentication",
		Type:     "denied",
		Action:   "login",
		Outcome:  "failure",
	}

	if event.Name != "iam.login.rejected" || event.Dataset != "dox.iam.security" {
		t.Fatalf("expected observability event to describe the log record, got %+v", event)
	}
}

func assertCore(t *testing.T, core CoreConfig, name string, coreType CoreType, encoding Encoding, outputPaths []string) {
	t.Helper()
	if core.Name != name {
		t.Fatalf("expected core name %q, got %q", name, core.Name)
	}
	if !boolValue(core.Enabled) {
		t.Fatalf("expected core %q to be enabled", name)
	}
	if core.Type != coreType {
		t.Fatalf("expected core %q type %q, got %q", name, coreType, core.Type)
	}
	if core.Encoding != encoding {
		t.Fatalf("expected core %q encoding %q, got %q", name, encoding, core.Encoding)
	}
	if len(core.OutputPaths) != len(outputPaths) {
		t.Fatalf("expected core %q output paths %#v, got %#v", name, outputPaths, core.OutputPaths)
	}
	for index, expected := range outputPaths {
		if core.OutputPaths[index] != expected {
			t.Fatalf("expected core %q output path %d to be %q, got %q", name, index, expected, core.OutputPaths[index])
		}
	}
}

func assertEqual(t *testing.T, got string, expected string) {
	t.Helper()
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func contains(items []string, expected string) bool {
	for _, item := range items {
		if item == expected {
			return true
		}
	}
	return false
}

func hasValidationField(err error, field string) bool {
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		return false
	}
	for _, item := range validationErr.Fields {
		if item.Field == field {
			return true
		}
	}
	return false
}
