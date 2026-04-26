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
 * @File    : validate.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package logging

import (
	"strconv"
	"strings"
)

var levelEncoderNames = map[string]struct{}{
	"lowercase":       {},
	"lowercase_color": {},
	"capital":         {},
	"capital_color":   {},
}

var timeEncoderNames = map[string]struct{}{
	"rfc3339nano": {},
	"rfc3339":     {},
	"iso8601":     {},
	"epoch":       {},
	"millis":      {},
	"nanos":       {},
}

var durationEncoderNames = map[string]struct{}{
	"seconds": {},
	"millis":  {},
	"nanos":   {},
	"string":  {},
}

var callerEncoderNames = map[string]struct{}{
	"short": {},
	"full":  {},
}

var nameEncoderNames = map[string]struct{}{
	"full": {},
}

// FieldError describes one logging configuration validation failure.
type FieldError struct {
	Field  string
	Reason string
}

// ValidationError describes one or more logging validation failures.
type ValidationError struct {
	Fields []FieldError
}

// Error returns a compact validation failure message.
func (e *ValidationError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if len(e.Fields) == 0 {
		return "logging validation failed"
	}

	parts := make([]string, 0, len(e.Fields))
	for _, field := range e.Fields {
		if field.Reason == "" {
			parts = append(parts, field.Field)
			continue
		}
		parts = append(parts, field.Field+": "+field.Reason)
	}
	return "logging validation failed: " + strings.Join(parts, ", ")
}

// Validate verifies the shared logging configuration contract.
func (c Config) Validate() error {
	v := validator{}
	v.level("level", c.Level)
	c.Zap.validate(&v)
	v.cores(c.Cores)
	c.Buffering.validate(&v)
	c.Shutdown.validate(&v)
	c.Redaction.validate(&v)
	c.OTel.validate(&v)
	return v.err()
}

func (c ZapConfig) validate(v *validator) {
	v.level("zap.level", c.Level)
	v.encoding("zap.encoding", c.Encoding)
	c.EncoderConfig.validate(v, "zap.encoder_config")
	if len(c.OutputPaths) == 0 {
		v.add("zap.output_paths", "at least one output path is required")
	}
	for index, path := range c.OutputPaths {
		if strings.TrimSpace(path) == "" {
			v.add(indexField("zap.output_paths", index), "output path must not be empty")
		}
	}
	if len(c.ErrorOutputPaths) == 0 {
		v.add("zap.error_output_paths", "at least one error output path is required")
	}
	for index, path := range c.ErrorOutputPaths {
		if strings.TrimSpace(path) == "" {
			v.add(indexField("zap.error_output_paths", index), "error output path must not be empty")
		}
	}
	if c.Sampling.Enabled {
		if c.Sampling.Initial <= 0 {
			v.add("zap.sampling.initial", "initial must be positive when sampling is enabled")
		}
		if c.Sampling.Thereafter <= 0 {
			v.add("zap.sampling.thereafter", "thereafter must be positive when sampling is enabled")
		}
	}
}

func (c EncoderConfig) validate(v *validator, field string) {
	v.symbol(field+".level_encoder", c.LevelEncoder, levelEncoderNames)
	v.symbol(field+".time_encoder", c.TimeEncoder, timeEncoderNames)
	v.symbol(field+".duration_encoder", c.DurationEncoder, durationEncoderNames)
	v.symbol(field+".caller_encoder", c.CallerEncoder, callerEncoderNames)
	v.symbol(field+".name_encoder", c.NameEncoder, nameEncoderNames)
}

func (c BufferingConfig) validate(v *validator) {
	if !boolValue(c.Enabled) {
		return
	}
	if c.SizeBytes <= 0 {
		v.add("buffering.size_bytes", "size must be positive when buffering is enabled")
	}
	if c.FlushInterval <= 0 {
		v.add("buffering.flush_interval", "flush interval must be positive when buffering is enabled")
	}
}

func (c ShutdownConfig) validate(v *validator) {
	if c.Timeout <= 0 {
		v.add("shutdown.timeout", "timeout must be positive")
	}
}

func (c RedactionConfig) validate(v *validator) {
	if !boolValue(c.Enabled) {
		return
	}
	if strings.TrimSpace(c.Replacement) == "" {
		v.add("redaction.replacement", "replacement must not be empty when redaction is enabled")
	}
	if len(c.Keys) == 0 {
		v.add("redaction.keys", "at least one redaction key is required when redaction is enabled")
	}
	for index, key := range c.Keys {
		if strings.TrimSpace(key) == "" {
			v.add(indexField("redaction.keys", index), "redaction key must not be empty")
		}
	}
}

func (c OpenTelemetryConfig) validate(v *validator) {
	if !boolValue(c.Enabled) {
		return
	}
	c.Traces.validate(v)
	c.Exporter.validate(v)
	c.Batch.validate(v)
}

func (c OpenTelemetryTraces) validate(v *validator) {
	if !boolValue(c.Enabled) {
		return
	}
	if !c.Sampler.Type.IsValid() {
		v.add("otel.traces.sampler.type", "trace sampler type is not supported")
	}
	if c.Sampler.Ratio < 0 || c.Sampler.Ratio > 1 {
		v.add("otel.traces.sampler.ratio", "trace sampler ratio must be between 0 and 1")
	}
}

func (c OpenTelemetryExporter) validate(v *validator) {
	c.OTLP.validate(v)
}

func (c OTLPExporterConfig) validate(v *validator) {
	if !c.Protocol.IsValid() {
		v.add("otel.exporter.otlp.protocol", "OTLP protocol is not supported")
	}
	if c.Enabled && strings.TrimSpace(c.Endpoint) == "" {
		v.add("otel.exporter.otlp.endpoint", "endpoint is required when OTLP exporter is enabled")
	}
	if c.Enabled && c.Timeout <= 0 {
		v.add("otel.exporter.otlp.timeout", "timeout must be positive when OTLP exporter is enabled")
	}
	for key := range c.Headers {
		if strings.TrimSpace(key) == "" {
			v.add("otel.exporter.otlp.headers", "header keys must not be empty")
		}
	}
}

func (c OpenTelemetryBatch) validate(v *validator) {
	if c.MaxQueueSize <= 0 {
		v.add("otel.batch.max_queue_size", "max queue size must be positive")
	}
	if c.ScheduleDelay <= 0 {
		v.add("otel.batch.schedule_delay", "schedule delay must be positive")
	}
	if c.ExportTimeout <= 0 {
		v.add("otel.batch.export_timeout", "export timeout must be positive")
	}
	if c.MaxExportBatchSize <= 0 {
		v.add("otel.batch.max_export_batch_size", "max export batch size must be positive")
	}
	if c.MaxExportBatchSize > c.MaxQueueSize {
		v.add("otel.batch.max_export_batch_size", "max export batch size must not exceed max queue size")
	}
}

type validator struct {
	fields []FieldError
}

func (v *validator) cores(cores []CoreConfig) {
	if len(cores) == 0 {
		v.add("cores", "at least one core is required")
		return
	}
	seen := map[string]struct{}{}
	for index, core := range cores {
		field := indexField("cores", index)
		name := strings.TrimSpace(core.Name)
		if name == "" {
			v.add(field+".name", "core name is required")
		} else if _, exists := seen[name]; exists {
			v.add(field+".name", "core name must be unique")
		} else {
			seen[name] = struct{}{}
		}

		if !core.Type.IsValid() {
			v.add(field+".type", "core type is not supported")
		}
		v.level(field+".level", core.Level)
		v.encoding(field+".encoding", core.Encoding)
		if boolValue(core.Enabled) && len(core.OutputPaths) == 0 {
			v.add(field+".output_paths", "enabled core requires at least one output path")
		}
		for pathIndex, path := range core.OutputPaths {
			if strings.TrimSpace(path) == "" {
				v.add(indexField(field+".output_paths", pathIndex), "output path must not be empty")
			}
		}
		if len(core.Datasets) == 0 {
			v.add(field+".datasets", "at least one dataset routing entry is required")
		}
		for datasetIndex, dataset := range core.Datasets {
			if strings.TrimSpace(dataset) == "" {
				v.add(indexField(field+".datasets", datasetIndex), "dataset routing entry must not be empty")
			}
		}
		if core.Type == CoreTypeFile {
			core.Rotation.validate(v, field+".rotation")
		}
	}
}

func (c RotationConfig) validate(v *validator, field string) {
	if !c.Driver.IsValid() {
		v.add(field+".driver", "rotation driver is not supported")
	}
	if !boolValue(c.Enabled) {
		return
	}
	if c.Driver == RotationDriverLumberjack {
		if c.MaxSizeMB <= 0 {
			v.add(field+".max_size_mb", "max size must be positive")
		}
		if c.MaxBackups < 0 {
			v.add(field+".max_backups", "max backups must not be negative")
		}
		if c.MaxAgeDays < 0 {
			v.add(field+".max_age_days", "max age must not be negative")
		}
	}
}

func (v *validator) level(field string, level Level) {
	if !level.IsValid() {
		v.add(field, "level is not supported")
	}
}

func (v *validator) encoding(field string, encoding Encoding) {
	if !encoding.IsValid() {
		v.add(field, "encoding is not supported")
	}
}

func (v *validator) symbol(field string, value string, allowed map[string]struct{}) {
	if strings.TrimSpace(value) == "" {
		v.add(field, "encoder name is required")
		return
	}
	if _, ok := allowed[value]; !ok {
		v.add(field, "encoder name is not supported")
	}
}

func (v *validator) add(field string, reason string) {
	v.fields = append(v.fields, FieldError{Field: field, Reason: reason})
}

func (v *validator) err() error {
	if len(v.fields) == 0 {
		return nil
	}
	return &ValidationError{Fields: v.fields}
}

func indexField(field string, index int) string {
	return field + "[" + strconv.Itoa(index) + "]"
}
