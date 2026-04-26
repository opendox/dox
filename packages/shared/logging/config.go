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
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package logging

import (
	"errors"
	"time"
)

const (
	// DefaultFilePathTemplate is the default service JSONL log path template.
	DefaultFilePathTemplate = "logs/${service.namespace}-${service.name}.jsonl"
	// DefaultRedactionReplacement is the default replacement for sensitive values.
	DefaultRedactionReplacement = "[REDACTED]"
)

// Level identifies a logging level accepted by the shared logging contract.
type Level string

const (
	LevelDebug  Level = "debug"
	LevelInfo   Level = "info"
	LevelWarn   Level = "warn"
	LevelError  Level = "error"
	LevelDPanic Level = "dpanic"
	LevelPanic  Level = "panic"
	LevelFatal  Level = "fatal"
)

// IsValid reports whether l is a supported logging level.
func (l Level) IsValid() bool {
	switch l {
	case LevelDebug, LevelInfo, LevelWarn, LevelError, LevelDPanic, LevelPanic, LevelFatal:
		return true
	default:
		return false
	}
}

// Encoding identifies a configured log encoder format.
type Encoding string

const (
	EncodingConsole Encoding = "console"
	EncodingJSON    Encoding = "json"
)

// IsValid reports whether e is a supported log encoding.
func (e Encoding) IsValid() bool {
	switch e {
	case EncodingConsole, EncodingJSON:
		return true
	default:
		return false
	}
}

// CoreType identifies one configured zapcore-style output core.
type CoreType string

const (
	CoreTypeConsole CoreType = "console"
	CoreTypeFile    CoreType = "file"
)

// IsValid reports whether t is a supported core type.
func (t CoreType) IsValid() bool {
	switch t {
	case CoreTypeConsole, CoreTypeFile:
		return true
	default:
		return false
	}
}

// RotationDriver identifies a file rotation implementation strategy.
type RotationDriver string

const (
	RotationDriverLumberjack RotationDriver = "lumberjack"
	RotationDriverExternal   RotationDriver = "external"
	RotationDriverLogrotate  RotationDriver = "logrotate"
	RotationDriverNone       RotationDriver = "none"
)

// IsValid reports whether d is a supported rotation driver.
func (d RotationDriver) IsValid() bool {
	switch d {
	case RotationDriverLumberjack, RotationDriverExternal, RotationDriverLogrotate, RotationDriverNone:
		return true
	default:
		return false
	}
}

// OTLPProtocol identifies an OTLP exporter protocol.
type OTLPProtocol string

const (
	OTLPProtocolGRPC OTLPProtocol = "grpc"
	OTLPProtocolHTTP OTLPProtocol = "http"
)

// IsValid reports whether p is a supported OTLP protocol.
func (p OTLPProtocol) IsValid() bool {
	switch p {
	case OTLPProtocolGRPC, OTLPProtocolHTTP:
		return true
	default:
		return false
	}
}

// TraceSamplerType identifies an OpenTelemetry trace sampler shape.
type TraceSamplerType string

const (
	TraceSamplerAlwaysOn                TraceSamplerType = "always_on"
	TraceSamplerAlwaysOff               TraceSamplerType = "always_off"
	TraceSamplerTraceIDRatio            TraceSamplerType = "traceidratio"
	TraceSamplerParentBasedTraceIDRatio TraceSamplerType = "parentbased_traceidratio"
)

// IsValid reports whether t is a supported trace sampler type.
func (t TraceSamplerType) IsValid() bool {
	switch t {
	case TraceSamplerAlwaysOn, TraceSamplerAlwaysOff, TraceSamplerTraceIDRatio, TraceSamplerParentBasedTraceIDRatio:
		return true
	default:
		return false
	}
}

// Config is the shared logging configuration contract.
type Config struct {
	Level       Level               `json:"level" yaml:"level" mapstructure:"level"`
	Development bool                `json:"development" yaml:"development" mapstructure:"development"`
	Resource    ResourceConfig      `json:"resource" yaml:"resource" mapstructure:"resource"`
	Zap         ZapConfig           `json:"zap" yaml:"zap" mapstructure:"zap"`
	Cores       []CoreConfig        `json:"cores" yaml:"cores" mapstructure:"cores"`
	Buffering   BufferingConfig     `json:"buffering" yaml:"buffering" mapstructure:"buffering"`
	Shutdown    ShutdownConfig      `json:"shutdown" yaml:"shutdown" mapstructure:"shutdown"`
	Redaction   RedactionConfig     `json:"redaction" yaml:"redaction" mapstructure:"redaction"`
	OTel        OpenTelemetryConfig `json:"otel" yaml:"otel" mapstructure:"otel"`
}

// ResourceConfig carries logging-specific resource overrides.
type ResourceConfig struct {
	ServiceVersion string `json:"service_version" yaml:"service_version" mapstructure:"service_version"`
}

// ZapConfig mirrors the zap.Config shape without importing zap.
type ZapConfig struct {
	Level               Level          `json:"level" yaml:"level" mapstructure:"level"`
	Development         bool           `json:"development" yaml:"development" mapstructure:"development"`
	DisableCaller       bool           `json:"disable_caller" yaml:"disable_caller" mapstructure:"disable_caller"`
	DisableStacktrace   bool           `json:"disable_stacktrace" yaml:"disable_stacktrace" mapstructure:"disable_stacktrace"`
	DisableErrorVerbose bool           `json:"disable_error_verbose" yaml:"disable_error_verbose" mapstructure:"disable_error_verbose"`
	Encoding            Encoding       `json:"encoding" yaml:"encoding" mapstructure:"encoding"`
	EncoderConfig       EncoderConfig  `json:"encoder_config" yaml:"encoder_config" mapstructure:"encoder_config"`
	OutputPaths         []string       `json:"output_paths" yaml:"output_paths" mapstructure:"output_paths"`
	ErrorOutputPaths    []string       `json:"error_output_paths" yaml:"error_output_paths" mapstructure:"error_output_paths"`
	InitialFields       map[string]any `json:"initial_fields" yaml:"initial_fields" mapstructure:"initial_fields"`
	Sampling            SamplingConfig `json:"sampling" yaml:"sampling" mapstructure:"sampling"`
}

// EncoderConfig mirrors the zapcore.EncoderConfig shape with symbolic encoders.
type EncoderConfig struct {
	MessageKey       string `json:"message_key" yaml:"message_key" mapstructure:"message_key"`
	LevelKey         string `json:"level_key" yaml:"level_key" mapstructure:"level_key"`
	TimeKey          string `json:"time_key" yaml:"time_key" mapstructure:"time_key"`
	NameKey          string `json:"name_key" yaml:"name_key" mapstructure:"name_key"`
	CallerKey        string `json:"caller_key" yaml:"caller_key" mapstructure:"caller_key"`
	FunctionKey      string `json:"function_key" yaml:"function_key" mapstructure:"function_key"`
	StacktraceKey    string `json:"stacktrace_key" yaml:"stacktrace_key" mapstructure:"stacktrace_key"`
	SkipLineEnding   bool   `json:"skip_line_ending" yaml:"skip_line_ending" mapstructure:"skip_line_ending"`
	LineEnding       string `json:"line_ending" yaml:"line_ending" mapstructure:"line_ending"`
	LevelEncoder     string `json:"level_encoder" yaml:"level_encoder" mapstructure:"level_encoder"`
	TimeEncoder      string `json:"time_encoder" yaml:"time_encoder" mapstructure:"time_encoder"`
	DurationEncoder  string `json:"duration_encoder" yaml:"duration_encoder" mapstructure:"duration_encoder"`
	CallerEncoder    string `json:"caller_encoder" yaml:"caller_encoder" mapstructure:"caller_encoder"`
	NameEncoder      string `json:"name_encoder" yaml:"name_encoder" mapstructure:"name_encoder"`
	ConsoleSeparator string `json:"console_separator" yaml:"console_separator" mapstructure:"console_separator"`
}

// SamplingConfig mirrors zap sampling configuration without importing zap.
type SamplingConfig struct {
	Enabled     bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Initial     int  `json:"initial" yaml:"initial" mapstructure:"initial"`
	Thereafter  int  `json:"thereafter" yaml:"thereafter" mapstructure:"thereafter"`
	HookMetrics bool `json:"hook_metrics" yaml:"hook_metrics" mapstructure:"hook_metrics"`
}

// CoreConfig declares one future zapcore output core.
type CoreConfig struct {
	Name        string         `json:"name" yaml:"name" mapstructure:"name"`
	Enabled     *bool          `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Type        CoreType       `json:"type" yaml:"type" mapstructure:"type"`
	Level       Level          `json:"level" yaml:"level" mapstructure:"level"`
	Encoding    Encoding       `json:"encoding" yaml:"encoding" mapstructure:"encoding"`
	OutputPaths []string       `json:"output_paths" yaml:"output_paths" mapstructure:"output_paths"`
	Datasets    []string       `json:"datasets" yaml:"datasets" mapstructure:"datasets"`
	Rotation    RotationConfig `json:"rotation" yaml:"rotation" mapstructure:"rotation"`
}

// RotationConfig declares a future file rotation mapping.
type RotationConfig struct {
	Driver     RotationDriver `json:"driver" yaml:"driver" mapstructure:"driver"`
	Enabled    *bool          `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	MaxSizeMB  int            `json:"max_size_mb" yaml:"max_size_mb" mapstructure:"max_size_mb"`
	MaxBackups int            `json:"max_backups" yaml:"max_backups" mapstructure:"max_backups"`
	MaxAgeDays int            `json:"max_age_days" yaml:"max_age_days" mapstructure:"max_age_days"`
	Compress   *bool          `json:"compress" yaml:"compress" mapstructure:"compress"`
	LocalTime  *bool          `json:"local_time" yaml:"local_time" mapstructure:"local_time"`
}

// BufferingConfig declares future buffered writer behavior.
type BufferingConfig struct {
	Enabled       *bool         `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	SizeBytes     int           `json:"size_bytes" yaml:"size_bytes" mapstructure:"size_bytes"`
	FlushInterval time.Duration `json:"flush_interval" yaml:"flush_interval" mapstructure:"flush_interval"`
}

// ShutdownConfig declares logging shutdown behavior.
type ShutdownConfig struct {
	Timeout time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

// RedactionConfig declares sensitive key replacement policy.
type RedactionConfig struct {
	Enabled     *bool    `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Replacement string   `json:"replacement" yaml:"replacement" mapstructure:"replacement"`
	Keys        []string `json:"keys" yaml:"keys" mapstructure:"keys"`
}

// OpenTelemetryConfig declares future OpenTelemetry SDK integration settings.
type OpenTelemetryConfig struct {
	Enabled     *bool                    `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Propagation OpenTelemetryPropagation `json:"propagation" yaml:"propagation" mapstructure:"propagation"`
	Traces      OpenTelemetryTraces      `json:"traces" yaml:"traces" mapstructure:"traces"`
	Metrics     OpenTelemetrySignal      `json:"metrics" yaml:"metrics" mapstructure:"metrics"`
	Logs        OpenTelemetrySignal      `json:"logs" yaml:"logs" mapstructure:"logs"`
	Exporter    OpenTelemetryExporter    `json:"exporter" yaml:"exporter" mapstructure:"exporter"`
	Batch       OpenTelemetryBatch       `json:"batch" yaml:"batch" mapstructure:"batch"`
}

// OpenTelemetryPropagation declares propagation settings.
type OpenTelemetryPropagation struct {
	TraceContext *bool `json:"trace_context" yaml:"trace_context" mapstructure:"trace_context"`
	Baggage      *bool `json:"baggage" yaml:"baggage" mapstructure:"baggage"`
}

// OpenTelemetryTraces declares trace settings.
type OpenTelemetryTraces struct {
	Enabled *bool              `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Sampler TraceSamplerConfig `json:"sampler" yaml:"sampler" mapstructure:"sampler"`
}

// TraceSamplerConfig declares trace sampler settings.
type TraceSamplerConfig struct {
	Type  TraceSamplerType `json:"type" yaml:"type" mapstructure:"type"`
	Ratio float64          `json:"ratio" yaml:"ratio" mapstructure:"ratio"`
}

// OpenTelemetrySignal declares a simple OpenTelemetry signal toggle.
type OpenTelemetrySignal struct {
	Enabled *bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
}

// OpenTelemetryExporter declares exporter settings.
type OpenTelemetryExporter struct {
	OTLP OTLPExporterConfig `json:"otlp" yaml:"otlp" mapstructure:"otlp"`
}

// OTLPExporterConfig declares future OTLP exporter settings.
type OTLPExporterConfig struct {
	Enabled     bool              `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Endpoint    string            `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`
	Protocol    OTLPProtocol      `json:"protocol" yaml:"protocol" mapstructure:"protocol"`
	Insecure    *bool             `json:"insecure" yaml:"insecure" mapstructure:"insecure"`
	Headers     map[string]string `json:"headers" yaml:"headers" mapstructure:"headers"`
	Compression string            `json:"compression" yaml:"compression" mapstructure:"compression"`
	Timeout     time.Duration     `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

// OpenTelemetryBatch declares batch processor settings.
type OpenTelemetryBatch struct {
	MaxQueueSize       int           `json:"max_queue_size" yaml:"max_queue_size" mapstructure:"max_queue_size"`
	ScheduleDelay      time.Duration `json:"schedule_delay" yaml:"schedule_delay" mapstructure:"schedule_delay"`
	ExportTimeout      time.Duration `json:"export_timeout" yaml:"export_timeout" mapstructure:"export_timeout"`
	MaxExportBatchSize int           `json:"max_export_batch_size" yaml:"max_export_batch_size" mapstructure:"max_export_batch_size"`
}

// Default fills stable shared logging defaults.
func (c *Config) Default() error {
	if c == nil {
		return errors.New("logging: config must not be nil")
	}
	if c.Level == "" {
		c.Level = LevelInfo
	}
	c.Zap.Default(c.Level, c.Development)
	if len(c.Cores) == 0 {
		c.Cores = []CoreConfig{
			defaultConsoleCore(c.Level),
			defaultServiceFileCore(c.Level),
		}
	} else {
		for index := range c.Cores {
			c.Cores[index].Default(c.Level)
		}
	}
	c.Buffering.Default()
	c.Shutdown.Default()
	c.Redaction.Default()
	c.OTel.Default()
	return nil
}

// Default fills zap-facing defaults.
func (c *ZapConfig) Default(level Level, development bool) {
	if c.Level == "" {
		c.Level = level
	}
	c.Development = c.Development || development
	if c.Encoding == "" {
		c.Encoding = EncodingJSON
	}
	c.EncoderConfig.Default()
	if len(c.OutputPaths) == 0 {
		c.OutputPaths = []string{"stdout"}
	}
	if len(c.ErrorOutputPaths) == 0 {
		c.ErrorOutputPaths = []string{"stderr"}
	}
	if c.InitialFields == nil {
		c.InitialFields = map[string]any{}
	}
	c.Sampling.Default()
}

// Default fills encoder defaults.
func (c *EncoderConfig) Default() {
	if c.MessageKey == "" {
		c.MessageKey = "message"
	}
	if c.LevelKey == "" {
		c.LevelKey = "severity_text"
	}
	if c.TimeKey == "" {
		c.TimeKey = "timestamp"
	}
	if c.NameKey == "" {
		c.NameKey = "logger"
	}
	if c.CallerKey == "" {
		c.CallerKey = "caller"
	}
	if c.FunctionKey == "" {
		c.FunctionKey = "function"
	}
	if c.StacktraceKey == "" {
		c.StacktraceKey = "stacktrace"
	}
	if c.LineEnding == "" {
		c.LineEnding = "\n"
	}
	if c.LevelEncoder == "" {
		c.LevelEncoder = "capital"
	}
	if c.TimeEncoder == "" {
		c.TimeEncoder = "rfc3339nano"
	}
	if c.DurationEncoder == "" {
		c.DurationEncoder = "millis"
	}
	if c.CallerEncoder == "" {
		c.CallerEncoder = "short"
	}
	if c.NameEncoder == "" {
		c.NameEncoder = "full"
	}
	if c.ConsoleSeparator == "" {
		c.ConsoleSeparator = "\t"
	}
}

// Default fills sampling defaults without enabling sampling by itself.
func (c *SamplingConfig) Default() {
	if c.Initial == 0 {
		c.Initial = 100
	}
	if c.Thereafter == 0 {
		c.Thereafter = 100
	}
}

// Default fills core defaults.
func (c *CoreConfig) Default(level Level) {
	defaultBool(&c.Enabled, true)
	if c.Level == "" {
		c.Level = level
	}
	if c.Type == CoreTypeConsole {
		if c.Encoding == "" {
			c.Encoding = EncodingConsole
		}
		if len(c.OutputPaths) == 0 {
			c.OutputPaths = []string{"stdout"}
		}
	}
	if c.Type == CoreTypeFile {
		if c.Encoding == "" {
			c.Encoding = EncodingJSON
		}
		if len(c.OutputPaths) == 0 {
			c.OutputPaths = []string{DefaultFilePathTemplate}
		}
		c.Rotation.Default()
	}
	if len(c.Datasets) == 0 {
		c.Datasets = []string{"*"}
	}
}

// Default fills rotation defaults.
func (c *RotationConfig) Default() {
	if c.Driver == "" {
		c.Driver = RotationDriverLumberjack
	}
	if c.Driver == RotationDriverLumberjack {
		if c.MaxSizeMB == 0 {
			c.MaxSizeMB = 100
		}
		if c.MaxBackups == 0 {
			c.MaxBackups = 10
		}
		if c.MaxAgeDays == 0 {
			c.MaxAgeDays = 14
		}
		defaultBool(&c.Compress, true)
		defaultBool(&c.LocalTime, true)
	}
	if c.Driver != RotationDriverNone {
		defaultBool(&c.Enabled, true)
		return
	}
	defaultBool(&c.Enabled, false)
}

// Default fills buffering defaults.
func (c *BufferingConfig) Default() {
	if c.SizeBytes == 0 {
		c.SizeBytes = 262144
	}
	if c.FlushInterval == 0 {
		c.FlushInterval = time.Second
	}
	defaultBool(&c.Enabled, true)
}

// Default fills shutdown defaults.
func (c *ShutdownConfig) Default() {
	if c.Timeout == 0 {
		c.Timeout = 5 * time.Second
	}
}

// Default fills redaction defaults.
func (c *RedactionConfig) Default() {
	defaultBool(&c.Enabled, true)
	if c.Replacement == "" {
		c.Replacement = DefaultRedactionReplacement
	}
	if len(c.Keys) == 0 {
		c.Keys = DefaultRedactionKeys()
	}
}

// Default fills OpenTelemetry configuration defaults.
func (c *OpenTelemetryConfig) Default() {
	defaultBool(&c.Enabled, true)
	c.Propagation.Default()
	c.Traces.Default()
	defaultBool(&c.Metrics.Enabled, true)
	defaultBool(&c.Logs.Enabled, true)
	c.Exporter.Default()
	c.Batch.Default()
}

// Default fills OpenTelemetry propagation defaults.
func (c *OpenTelemetryPropagation) Default() {
	defaultBool(&c.TraceContext, true)
	defaultBool(&c.Baggage, true)
}

// Default fills OpenTelemetry trace defaults.
func (c *OpenTelemetryTraces) Default() {
	defaultBool(&c.Enabled, true)
	c.Sampler.Default()
}

// Default fills trace sampler defaults.
func (c *TraceSamplerConfig) Default() {
	if c.Type == "" {
		c.Type = TraceSamplerParentBasedTraceIDRatio
	}
	if c.Ratio == 0 {
		c.Ratio = 1
	}
}

// Default fills exporter defaults.
func (c *OpenTelemetryExporter) Default() {
	c.OTLP.Default()
}

// Default fills OTLP exporter defaults.
func (c *OTLPExporterConfig) Default() {
	if c.Endpoint == "" {
		c.Endpoint = "localhost:4317"
	}
	if c.Protocol == "" {
		c.Protocol = OTLPProtocolGRPC
	}
	defaultBool(&c.Insecure, true)
	if c.Headers == nil {
		c.Headers = map[string]string{}
	}
	if c.Compression == "" {
		c.Compression = "gzip"
	}
	if c.Timeout == 0 {
		c.Timeout = 10 * time.Second
	}
}

// Default fills OpenTelemetry batch defaults.
func (c *OpenTelemetryBatch) Default() {
	if c.MaxQueueSize == 0 {
		c.MaxQueueSize = 2048
	}
	if c.ScheduleDelay == 0 {
		c.ScheduleDelay = 5 * time.Second
	}
	if c.ExportTimeout == 0 {
		c.ExportTimeout = 30 * time.Second
	}
	if c.MaxExportBatchSize == 0 {
		c.MaxExportBatchSize = 512
	}
}

// DefaultRedactionKeys returns the default sensitive key list.
func DefaultRedactionKeys() []string {
	return []string{
		"authorization",
		"cookie",
		"set-cookie",
		"token",
		"access_token",
		"refresh_token",
		"password",
		"passwd",
		"secret",
		"private_key",
		"credential",
		"api_key",
		"client_secret",
		"session",
	}
}

func defaultConsoleCore(level Level) CoreConfig {
	return CoreConfig{
		Name:        "console",
		Enabled:     boolPtr(true),
		Type:        CoreTypeConsole,
		Level:       level,
		Encoding:    EncodingConsole,
		OutputPaths: []string{"stdout"},
		Datasets:    []string{"*"},
	}
}

func defaultServiceFileCore(level Level) CoreConfig {
	return CoreConfig{
		Name:        "service-file",
		Enabled:     boolPtr(true),
		Type:        CoreTypeFile,
		Level:       level,
		Encoding:    EncodingJSON,
		OutputPaths: []string{DefaultFilePathTemplate},
		Datasets:    []string{"*"},
		Rotation: RotationConfig{
			Driver:     RotationDriverLumberjack,
			Enabled:    boolPtr(true),
			MaxSizeMB:  100,
			MaxBackups: 10,
			MaxAgeDays: 14,
			Compress:   boolPtr(true),
			LocalTime:  boolPtr(true),
		},
	}
}

func defaultBool(target **bool, value bool) {
	if *target == nil {
		*target = boolPtr(value)
	}
}

func boolPtr(value bool) *bool {
	copied := value
	return &copied
}

func boolValue(value *bool) bool {
	return value != nil && *value
}
