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
 * @File    : zap.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package logging

import (
	"fmt"
	"sort"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLevelEncoders = map[string]zapcore.LevelEncoder{
	"lowercase":       zapcore.LowercaseLevelEncoder,
	"lowercase_color": zapcore.LowercaseColorLevelEncoder,
	"capital":         zapcore.CapitalLevelEncoder,
	"capital_color":   zapcore.CapitalColorLevelEncoder,
}

var zapTimeEncoders = map[string]zapcore.TimeEncoder{
	"rfc3339nano": zapcore.RFC3339NanoTimeEncoder,
	"rfc3339":     zapcore.RFC3339TimeEncoder,
	"iso8601":     zapcore.ISO8601TimeEncoder,
	"epoch":       zapcore.EpochTimeEncoder,
	"millis":      zapcore.EpochMillisTimeEncoder,
	"nanos":       zapcore.EpochNanosTimeEncoder,
}

var zapDurationEncoders = map[string]zapcore.DurationEncoder{
	"seconds": zapcore.SecondsDurationEncoder,
	"millis":  zapcore.MillisDurationEncoder,
	"nanos":   zapcore.NanosDurationEncoder,
	"string":  zapcore.StringDurationEncoder,
}

var zapCallerEncoders = map[string]zapcore.CallerEncoder{
	"short": zapcore.ShortCallerEncoder,
	"full":  zapcore.FullCallerEncoder,
}

var zapNameEncoders = map[string]zapcore.NameEncoder{
	"full": zapcore.FullNameEncoder,
}

// ZapCoreBase carries zap primitives for runtime integrations.
//
// Business code should depend on Dox logging types and the future Dox logger
// API instead of using these zap primitives directly.
type ZapCoreBase struct {
	Config       zap.Config
	Level        zap.AtomicLevel
	Core         zapcore.Core
	EnabledCores int

	options []zap.Option
	close   func()
}

// Options returns zap options mapped from the Dox logging configuration.
func (b *ZapCoreBase) Options() []zap.Option {
	if b == nil || len(b.options) == 0 {
		return nil
	}
	options := make([]zap.Option, len(b.options))
	copy(options, b.options)
	return options
}

// Close releases sinks opened by NewZapCoreBase.
func (b *ZapCoreBase) Close() {
	if b == nil || b.close == nil {
		return
	}
	closeFn := b.close
	b.close = nil
	closeFn()
}

// NewZapLevel maps a Dox logging level to a zapcore level.
func NewZapLevel(level Level) (zapcore.Level, error) {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel, nil
	case LevelInfo:
		return zapcore.InfoLevel, nil
	case LevelWarn:
		return zapcore.WarnLevel, nil
	case LevelError:
		return zapcore.ErrorLevel, nil
	case LevelDPanic:
		return zapcore.DPanicLevel, nil
	case LevelPanic:
		return zapcore.PanicLevel, nil
	case LevelFatal:
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InvalidLevel, validationError("level", "level is not supported")
	}
}

// NewZapAtomicLevel creates a zap AtomicLevel from a Dox logging level.
func NewZapAtomicLevel(level Level) (zap.AtomicLevel, error) {
	zapLevel, err := NewZapLevel(level)
	if err != nil {
		return zap.AtomicLevel{}, err
	}
	return zap.NewAtomicLevelAt(zapLevel), nil
}

// NewZapEncoderConfig maps the Dox symbolic encoder config to zapcore.
func NewZapEncoderConfig(config EncoderConfig) (zapcore.EncoderConfig, error) {
	config.Default()
	v := validator{}
	config.validate(&v, "encoder_config")
	if err := v.err(); err != nil {
		return zapcore.EncoderConfig{}, err
	}

	return zapcore.EncoderConfig{
		MessageKey:       config.MessageKey,
		LevelKey:         config.LevelKey,
		TimeKey:          config.TimeKey,
		NameKey:          config.NameKey,
		CallerKey:        config.CallerKey,
		FunctionKey:      config.FunctionKey,
		StacktraceKey:    config.StacktraceKey,
		SkipLineEnding:   config.SkipLineEnding,
		LineEnding:       config.LineEnding,
		EncodeLevel:      zapLevelEncoders[config.LevelEncoder],
		EncodeTime:       zapTimeEncoders[config.TimeEncoder],
		EncodeDuration:   zapDurationEncoders[config.DurationEncoder],
		EncodeCaller:     zapCallerEncoders[config.CallerEncoder],
		EncodeName:       zapNameEncoders[config.NameEncoder],
		ConsoleSeparator: config.ConsoleSeparator,
	}, nil
}

// NewZapConfig maps the Dox logging configuration to zap.Config.
func NewZapConfig(config Config) (zap.Config, error) {
	normalized, err := normalizeZapConfig(config)
	if err != nil {
		return zap.Config{}, err
	}
	return newZapConfigFromNormalized(normalized)
}

// NewZapCoreBase builds enabled zap cores and options from the Dox config.
func NewZapCoreBase(config Config) (*ZapCoreBase, error) {
	normalized, err := normalizeZapConfig(config)
	if err != nil {
		return nil, err
	}

	zapConfig, err := newZapConfigFromNormalized(normalized)
	if err != nil {
		return nil, err
	}

	errSink, closeErrSink, err := zap.Open(zapConfig.ErrorOutputPaths...)
	if err != nil {
		return nil, fmt.Errorf("logging: open zap error output paths: %w", err)
	}

	core, enabledCores, closeCore, err := newZapCoreFromNormalized(normalized, zapConfig.Level)
	if err != nil {
		closeErrSink()
		return nil, err
	}

	return &ZapCoreBase{
		Config:       zapConfig,
		Level:        zapConfig.Level,
		Core:         core,
		EnabledCores: enabledCores,
		options:      newZapOptions(normalized.Zap, errSink),
		close: func() {
			closeCore()
			closeErrSink()
		},
	}, nil
}

func normalizeZapConfig(config Config) (Config, error) {
	if err := config.Default(); err != nil {
		return Config{}, err
	}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

func newZapConfigFromNormalized(config Config) (zap.Config, error) {
	level, err := NewZapAtomicLevel(config.Zap.Level)
	if err != nil {
		return zap.Config{}, err
	}
	encoderConfig, err := NewZapEncoderConfig(config.Zap.EncoderConfig)
	if err != nil {
		return zap.Config{}, err
	}

	return zap.Config{
		Level:             level,
		Development:       config.Zap.Development,
		DisableCaller:     config.Zap.DisableCaller,
		DisableStacktrace: config.Zap.DisableStacktrace,
		Sampling:          newZapSamplingConfig(config.Zap.Sampling),
		Encoding:          string(config.Zap.Encoding),
		EncoderConfig:     encoderConfig,
		OutputPaths:       copyStrings(config.Zap.OutputPaths),
		ErrorOutputPaths:  copyStrings(config.Zap.ErrorOutputPaths),
		InitialFields:     copyInitialFields(config.Zap.InitialFields),
	}, nil
}

func newZapCoreFromNormalized(config Config, rootLevel zap.AtomicLevel) (zapcore.Core, int, func(), error) {
	encoderConfig, err := NewZapEncoderConfig(config.Zap.EncoderConfig)
	if err != nil {
		return nil, 0, nil, err
	}

	cores := make([]zapcore.Core, 0, len(config.Cores))
	closeFns := make([]func(), 0, len(config.Cores))
	for _, coreConfig := range config.Cores {
		if !boolValue(coreConfig.Enabled) {
			continue
		}

		encoder, err := newZapEncoder(coreConfig.Encoding, encoderConfig)
		if err != nil {
			closeZapSinks(closeFns)
			return nil, 0, nil, fmt.Errorf("logging: build zap core %q: %w", coreConfig.Name, err)
		}
		writer, closeWriter, err := zap.Open(coreConfig.OutputPaths...)
		if err != nil {
			closeZapSinks(closeFns)
			return nil, 0, nil, fmt.Errorf("logging: open zap core %q output paths: %w", coreConfig.Name, err)
		}
		closeFns = append(closeFns, closeWriter)

		coreLevel, err := NewZapLevel(coreConfig.Level)
		if err != nil {
			closeZapSinks(closeFns)
			return nil, 0, nil, fmt.Errorf("logging: build zap core %q: %w", coreConfig.Name, err)
		}
		cores = append(cores, zapcore.NewCore(encoder, writer, zapCoreLevelEnabler(rootLevel, coreLevel)))
	}

	core := zapcore.NewTee(cores...)
	if config.Zap.DisableErrorVerbose {
		core = noErrorVerboseCore{Core: core}
	}
	if config.Zap.Sampling.Enabled {
		core = zapcore.NewSamplerWithOptions(
			core,
			time.Second,
			config.Zap.Sampling.Initial,
			config.Zap.Sampling.Thereafter,
		)
	}

	return core, len(cores), func() {
		closeZapSinks(closeFns)
	}, nil
}

func newZapEncoder(encoding Encoding, config zapcore.EncoderConfig) (zapcore.Encoder, error) {
	switch encoding {
	case EncodingConsole:
		return zapcore.NewConsoleEncoder(config), nil
	case EncodingJSON:
		return zapcore.NewJSONEncoder(config), nil
	default:
		return nil, validationError("encoding", "encoding is not supported")
	}
}

func newZapSamplingConfig(config SamplingConfig) *zap.SamplingConfig {
	if !config.Enabled {
		return nil
	}
	return &zap.SamplingConfig{
		Initial:    config.Initial,
		Thereafter: config.Thereafter,
	}
}

func newZapOptions(config ZapConfig, errSink zapcore.WriteSyncer) []zap.Option {
	options := []zap.Option{zap.ErrorOutput(errSink)}
	if config.Development {
		options = append(options, zap.Development())
	}
	if !config.DisableCaller {
		options = append(options, zap.AddCaller())
	}
	if !config.DisableStacktrace {
		stackLevel := zapcore.ErrorLevel
		if config.Development {
			stackLevel = zapcore.WarnLevel
		}
		options = append(options, zap.AddStacktrace(stackLevel))
	}
	if fields := newZapInitialFields(config.InitialFields); len(fields) > 0 {
		options = append(options, zap.Fields(fields...))
	}
	return options
}

func newZapInitialFields(fields map[string]any) []zap.Field {
	if len(fields) == 0 {
		return nil
	}

	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	zapFields := make([]zap.Field, 0, len(fields))
	for _, key := range keys {
		zapFields = append(zapFields, zap.Any(key, fields[key]))
	}
	return zapFields
}

func zapCoreLevelEnabler(rootLevel zap.AtomicLevel, coreLevel zapcore.Level) zapcore.LevelEnabler {
	return zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return rootLevel.Enabled(level) && coreLevel.Enabled(level)
	})
}

func copyStrings(values []string) []string {
	if values == nil {
		return nil
	}
	copied := make([]string, len(values))
	copy(copied, values)
	return copied
}

func copyInitialFields(values map[string]any) map[string]any {
	if values == nil {
		return nil
	}
	copied := make(map[string]any, len(values))
	for key, value := range values {
		copied[key] = value
	}
	return copied
}

func closeZapSinks(closeFns []func()) {
	for _, closeFn := range closeFns {
		closeFn()
	}
}

func validationError(field string, reason string) error {
	return &ValidationError{
		Fields: []FieldError{
			{
				Field:  field,
				Reason: reason,
			},
		},
	}
}

type noErrorVerboseCore struct {
	zapcore.Core
}

func (c noErrorVerboseCore) With(fields []zapcore.Field) zapcore.Core {
	return noErrorVerboseCore{Core: c.Core.With(disableErrorVerboseFields(fields))}
}

func (c noErrorVerboseCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

func (c noErrorVerboseCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	return c.Core.Write(entry, disableErrorVerboseFields(fields))
}

func disableErrorVerboseFields(fields []zapcore.Field) []zapcore.Field {
	var copied []zapcore.Field
	for index, field := range fields {
		if field.Type != zapcore.ErrorType {
			continue
		}
		err, ok := field.Interface.(error)
		if !ok || err == nil {
			continue
		}
		if copied == nil {
			copied = make([]zapcore.Field, len(fields))
			copy(copied, fields)
		}
		copied[index] = zapcore.Field{
			Key:    field.Key,
			Type:   zapcore.StringType,
			String: err.Error(),
		}
	}
	if copied == nil {
		return fields
	}
	return copied
}
