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
 * @File    : zap_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewZapLevelMapsDoxLevels(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected zapcore.Level
	}{
		{name: "debug", level: LevelDebug, expected: zapcore.DebugLevel},
		{name: "info", level: LevelInfo, expected: zapcore.InfoLevel},
		{name: "warn", level: LevelWarn, expected: zapcore.WarnLevel},
		{name: "error", level: LevelError, expected: zapcore.ErrorLevel},
		{name: "dpanic", level: LevelDPanic, expected: zapcore.DPanicLevel},
		{name: "panic", level: LevelPanic, expected: zapcore.PanicLevel},
		{name: "fatal", level: LevelFatal, expected: zapcore.FatalLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := NewZapLevel(tt.level)
			if err != nil {
				t.Fatalf("map zap level: %v", err)
			}
			if level != tt.expected {
				t.Fatalf("expected zap level %s, got %s", tt.expected, level)
			}
		})
	}

	if _, err := NewZapLevel(Level("trace")); !hasValidationField(err, "level") {
		t.Fatalf("expected package validation error for invalid level, got %v", err)
	}
}

func TestNewZapAtomicLevelMapsDoxLevel(t *testing.T) {
	level, err := NewZapAtomicLevel(LevelWarn)
	if err != nil {
		t.Fatalf("map zap atomic level: %v", err)
	}
	if level.Level() != zapcore.WarnLevel {
		t.Fatalf("expected warn atomic level, got %s", level.Level())
	}
	level.SetLevel(zapcore.ErrorLevel)
	if level.Enabled(zapcore.WarnLevel) {
		t.Fatal("expected updated atomic level to reject warn")
	}
}

func TestNewZapEncoderConfigMapsSymbols(t *testing.T) {
	config := EncoderConfig{
		MessageKey:       "message",
		LevelKey:         "severity_text",
		TimeKey:          "timestamp",
		NameKey:          "logger",
		CallerKey:        "caller",
		FunctionKey:      "function",
		StacktraceKey:    "stacktrace",
		LineEnding:       "\n",
		LevelEncoder:     "lowercase",
		TimeEncoder:      "iso8601",
		DurationEncoder:  "seconds",
		CallerEncoder:    "full",
		NameEncoder:      "full",
		ConsoleSeparator: " | ",
	}

	encoderConfig, err := NewZapEncoderConfig(config)
	if err != nil {
		t.Fatalf("map zap encoder config: %v", err)
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	buffer, err := encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.WarnLevel,
		Time:       time.Date(2026, 4, 26, 10, 11, 12, 0, time.UTC),
		LoggerName: "iam",
		Caller: zapcore.EntryCaller{
			Defined: true,
			File:    "/srv/dox/auth.go",
			Line:    42,
		},
		Message: "login rejected",
	}, nil)
	if err != nil {
		t.Fatalf("encode entry: %v", err)
	}
	defer buffer.Free()

	text := buffer.String()
	for _, expected := range []string{
		`"severity_text":"warn"`,
		`"message":"login rejected"`,
		`"logger":"iam"`,
		`"caller":"/srv/dox/auth.go:42"`,
		`"timestamp":"2026-04-26T10:11:12.000Z"`,
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("expected encoded entry to contain %s, got %s", expected, text)
		}
	}
}

func TestNewZapEncoderConfigSupportsAllConfiguredSymbols(t *testing.T) {
	for _, name := range []string{"lowercase", "lowercase_color", "capital", "capital_color"} {
		config := EncoderConfig{LevelEncoder: name}
		if _, err := NewZapEncoderConfig(config); err != nil {
			t.Fatalf("expected level encoder %q to be supported: %v", name, err)
		}
	}
	for _, name := range []string{"rfc3339nano", "rfc3339", "iso8601", "epoch", "millis", "nanos"} {
		config := EncoderConfig{TimeEncoder: name}
		if _, err := NewZapEncoderConfig(config); err != nil {
			t.Fatalf("expected time encoder %q to be supported: %v", name, err)
		}
	}
	for _, name := range []string{"seconds", "millis", "nanos", "string"} {
		config := EncoderConfig{DurationEncoder: name}
		if _, err := NewZapEncoderConfig(config); err != nil {
			t.Fatalf("expected duration encoder %q to be supported: %v", name, err)
		}
	}
	for _, name := range []string{"short", "full"} {
		config := EncoderConfig{CallerEncoder: name}
		if _, err := NewZapEncoderConfig(config); err != nil {
			t.Fatalf("expected caller encoder %q to be supported: %v", name, err)
		}
	}

	if _, err := NewZapEncoderConfig(EncoderConfig{LevelEncoder: "shout"}); !hasValidationField(err, "encoder_config.level_encoder") {
		t.Fatalf("expected package validation error for invalid encoder, got %v", err)
	}
}

func TestNewZapConfigMapsDoxZapConfig(t *testing.T) {
	config := Config{
		Level:       LevelDebug,
		Development: true,
		Zap: ZapConfig{
			Level:               LevelWarn,
			Development:         true,
			DisableCaller:       true,
			DisableStacktrace:   true,
			DisableErrorVerbose: true,
			Encoding:            EncodingConsole,
			EncoderConfig: EncoderConfig{
				MessageKey:      "msg",
				LevelEncoder:    "capital",
				TimeEncoder:     "rfc3339",
				DurationEncoder: "string",
				CallerEncoder:   "short",
				NameEncoder:     "full",
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			InitialFields: map[string]any{
				"service.name": "iam",
			},
			Sampling: SamplingConfig{
				Enabled:    true,
				Initial:    3,
				Thereafter: 4,
			},
		},
	}

	zapConfig, err := NewZapConfig(config)
	if err != nil {
		t.Fatalf("map zap config: %v", err)
	}

	if zapConfig.Level.Level() != zapcore.WarnLevel {
		t.Fatalf("expected warn level, got %s", zapConfig.Level.Level())
	}
	if !zapConfig.Development || !zapConfig.DisableCaller || !zapConfig.DisableStacktrace {
		t.Fatalf("expected development/caller/stacktrace settings, got %+v", zapConfig)
	}
	if zapConfig.Encoding != "console" {
		t.Fatalf("expected console encoding, got %q", zapConfig.Encoding)
	}
	if zapConfig.OutputPaths[0] != "stdout" || zapConfig.ErrorOutputPaths[0] != "stderr" {
		t.Fatalf("expected output paths to map, got outputs=%#v errors=%#v", zapConfig.OutputPaths, zapConfig.ErrorOutputPaths)
	}
	if zapConfig.InitialFields["service.name"] != "iam" {
		t.Fatalf("expected initial fields to map, got %#v", zapConfig.InitialFields)
	}
	if zapConfig.Sampling == nil || zapConfig.Sampling.Initial != 3 || zapConfig.Sampling.Thereafter != 4 {
		t.Fatalf("expected sampling config to map, got %+v", zapConfig.Sampling)
	}
}

func TestNewZapConfigKeepsSamplingDisabledByDefault(t *testing.T) {
	zapConfig, err := NewZapConfig(Config{})
	if err != nil {
		t.Fatalf("map default zap config: %v", err)
	}
	if zapConfig.Sampling != nil {
		t.Fatalf("expected nil zap sampling config by default, got %+v", zapConfig.Sampling)
	}
}

func TestNewZapCoreBaseWritesEnabledConsoleAndJSONCores(t *testing.T) {
	tempDir := t.TempDir()
	consolePath := filepath.Join(tempDir, "console.log")
	jsonPath := filepath.Join(tempDir, "service.jsonl")

	base, err := NewZapCoreBase(Config{
		Level: LevelDebug,
		Zap: ZapConfig{
			Level:             LevelDebug,
			DisableCaller:     true,
			DisableStacktrace: true,
			ErrorOutputPaths:  []string{filepath.Join(tempDir, "errors.log")},
			InitialFields: map[string]any{
				"service.name": "iam",
			},
		},
		Cores: []CoreConfig{
			{
				Name:        "console",
				Enabled:     boolPtr(true),
				Type:        CoreTypeConsole,
				Level:       LevelDebug,
				Encoding:    EncodingConsole,
				OutputPaths: []string{consolePath},
				Datasets:    []string{"*"},
			},
			{
				Name:        "service-file",
				Enabled:     boolPtr(true),
				Type:        CoreTypeFile,
				Level:       LevelInfo,
				Encoding:    EncodingJSON,
				OutputPaths: []string{jsonPath},
				Datasets:    []string{"*"},
			},
		},
	})
	if err != nil {
		t.Fatalf("build zap core base: %v", err)
	}
	t.Cleanup(base.Close)

	if base.EnabledCores != 2 {
		t.Fatalf("expected two enabled cores, got %d", base.EnabledCores)
	}

	logger := zap.New(base.Core, base.Options()...)
	logger.Debug("debug-only")
	logger.Info("login rejected", zap.String("credential_type", "password"))
	_ = logger.Sync()
	base.Close()

	consoleText := readFile(t, consolePath)
	if !strings.Contains(consoleText, "debug-only") || !strings.Contains(consoleText, "login rejected") {
		t.Fatalf("expected console core to receive debug and info logs, got %s", consoleText)
	}

	jsonText := readFile(t, jsonPath)
	if strings.Contains(jsonText, "debug-only") {
		t.Fatalf("expected JSON core level to drop debug log, got %s", jsonText)
	}
	for _, expected := range []string{
		`"message":"login rejected"`,
		`"credential_type":"password"`,
		`"service.name":"iam"`,
	} {
		if !strings.Contains(jsonText, expected) {
			t.Fatalf("expected JSON output to contain %s, got %s", expected, jsonText)
		}
	}
}

func TestNewZapCoreBaseSkipsDisabledCores(t *testing.T) {
	base, err := NewZapCoreBase(Config{
		Zap: ZapConfig{
			DisableCaller:     true,
			DisableStacktrace: true,
			ErrorOutputPaths:  []string{"stderr"},
		},
		Cores: []CoreConfig{
			{
				Name:     "disabled",
				Enabled:  boolPtr(false),
				Type:     CoreTypeConsole,
				Level:    LevelInfo,
				Encoding: EncodingConsole,
				Datasets: []string{"*"},
			},
		},
	})
	if err != nil {
		t.Fatalf("build zap core base: %v", err)
	}
	t.Cleanup(base.Close)

	if base.EnabledCores != 0 {
		t.Fatalf("expected no enabled cores, got %d", base.EnabledCores)
	}
	if base.Core.Enabled(zapcore.InfoLevel) {
		t.Fatal("expected no-op core to reject info")
	}
}

func TestNewLumberjackLoggerMapsRotationConfig(t *testing.T) {
	coreConfig := CoreConfig{
		Name:        "service-file",
		Enabled:     boolPtr(true),
		Type:        CoreTypeFile,
		Level:       LevelInfo,
		Encoding:    EncodingJSON,
		OutputPaths: []string{"/var/log/dox/service.jsonl"},
		Datasets:    []string{"*"},
		Rotation: RotationConfig{
			Driver:     RotationDriverLumberjack,
			Enabled:    boolPtr(true),
			MaxSizeMB:  64,
			MaxBackups: 7,
			MaxAgeDays: 5,
			Compress:   boolPtr(true),
			LocalTime:  boolPtr(true),
		},
	}

	logger, err := newLumberjackLogger(coreConfig, "cores[0]")
	if err != nil {
		t.Fatalf("map lumberjack logger: %v", err)
	}

	if logger.Filename != "/var/log/dox/service.jsonl" ||
		logger.MaxSize != 64 ||
		logger.MaxBackups != 7 ||
		logger.MaxAge != 5 ||
		!logger.Compress ||
		!logger.LocalTime {
		t.Fatalf("expected lumberjack config to map, got %+v", logger)
	}
}

func TestNewZapCoreBaseWritesLumberjackJSONLFileCore(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "service.jsonl")

	base, err := NewZapCoreBase(Config{
		Zap: ZapConfig{
			DisableCaller:     true,
			DisableStacktrace: true,
			ErrorOutputPaths:  []string{filepath.Join(tempDir, "errors.log")},
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
					Driver:     RotationDriverLumberjack,
					Enabled:    boolPtr(true),
					MaxSizeMB:  1,
					MaxBackups: 2,
					MaxAgeDays: 3,
					Compress:   boolPtr(false),
					LocalTime:  boolPtr(false),
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("build zap core base: %v", err)
	}
	t.Cleanup(base.Close)

	logger := zap.New(base.Core, base.Options()...)
	logger.Info("login accepted", zap.String("event.name", "iam.login.accepted"))
	_ = logger.Sync()
	base.Close()

	entries := readJSONLines(t, jsonPath)
	if len(entries) != 1 {
		t.Fatalf("expected one JSONL entry, got %d: %#v", len(entries), entries)
	}
	if entries[0]["message"] != "login accepted" || entries[0]["event.name"] != "iam.login.accepted" {
		t.Fatalf("expected JSONL fields to be written, got %#v", entries[0])
	}
	if base.close != nil {
		t.Fatal("expected ZapCoreBase.Close to clear the close function")
	}
}

func TestNewZapCoreBaseWritesNoRotationFileCore(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "service.jsonl")

	base, err := NewZapCoreBase(Config{
		Zap: ZapConfig{
			DisableCaller:     true,
			DisableStacktrace: true,
			ErrorOutputPaths:  []string{filepath.Join(tempDir, "errors.log")},
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

	logger := zap.New(base.Core, base.Options()...)
	logger.Info("no rotation")
	_ = logger.Sync()
	base.Close()

	entries := readJSONLines(t, jsonPath)
	if len(entries) != 1 || entries[0]["message"] != "no rotation" {
		t.Fatalf("expected no-rotation JSONL entry, got %#v", entries)
	}
}

func TestNewZapCoreBaseRejectsUnsupportedFileRotationDriver(t *testing.T) {
	tempDir := t.TempDir()

	_, err := NewZapCoreBase(Config{
		Zap: ZapConfig{
			DisableCaller:     true,
			DisableStacktrace: true,
			ErrorOutputPaths:  []string{filepath.Join(tempDir, "errors.log")},
		},
		Cores: []CoreConfig{
			{
				Name:        "service-file",
				Enabled:     boolPtr(true),
				Type:        CoreTypeFile,
				Level:       LevelInfo,
				Encoding:    EncodingJSON,
				OutputPaths: []string{filepath.Join(tempDir, "service.jsonl")},
				Datasets:    []string{"*"},
				Rotation: RotationConfig{
					Driver: RotationDriverExternal,
				},
			},
		},
	})
	if !hasValidationField(err, "cores[0].rotation.driver") {
		t.Fatalf("expected unsupported rotation driver validation error, got %v", err)
	}
}

func TestNewZapCoreBaseRejectsInvalidLumberjackOutputPaths(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		outputPaths []string
		field       string
	}{
		{
			name:        "empty",
			outputPaths: []string{""},
			field:       "cores[0].output_paths[0]",
		},
		{
			name:        "multiple",
			outputPaths: []string{filepath.Join(tempDir, "first.jsonl"), filepath.Join(tempDir, "second.jsonl")},
			field:       "cores[0].output_paths",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewZapCoreBase(Config{
				Zap: ZapConfig{
					DisableCaller:     true,
					DisableStacktrace: true,
					ErrorOutputPaths:  []string{filepath.Join(tempDir, "errors.log")},
				},
				Cores: []CoreConfig{
					{
						Name:        "service-file",
						Enabled:     boolPtr(true),
						Type:        CoreTypeFile,
						Level:       LevelInfo,
						Encoding:    EncodingJSON,
						OutputPaths: tt.outputPaths,
						Datasets:    []string{"*"},
						Rotation: RotationConfig{
							Driver:     RotationDriverLumberjack,
							Enabled:    boolPtr(true),
							MaxSizeMB:  1,
							MaxBackups: 1,
							MaxAgeDays: 1,
							Compress:   boolPtr(false),
							LocalTime:  boolPtr(false),
						},
					},
				},
			})
			if !hasValidationField(err, tt.field) {
				t.Fatalf("expected validation field %s, got %v", tt.field, err)
			}
		})
	}
}

func TestNewZapCoreBaseSuppressesVerboseErrors(t *testing.T) {
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
			},
		},
	})
	if err != nil {
		t.Fatalf("build zap core base: %v", err)
	}
	t.Cleanup(base.Close)

	logger := zap.New(base.Core, base.Options()...)
	logger.Error("failed", zap.Error(verboseTestError{}))
	_ = logger.Sync()
	base.Close()

	jsonText := readFile(t, jsonPath)
	if strings.Contains(jsonText, "errorVerbose") {
		t.Fatalf("expected verbose error field to be suppressed, got %s", jsonText)
	}
	if !strings.Contains(jsonText, `"error":"basic"`) {
		t.Fatalf("expected basic error field, got %s", jsonText)
	}
}

func TestNewZapCoreBaseSamplingDefaultAndExplicitEnable(t *testing.T) {
	tests := []struct {
		name          string
		sampling      SamplingConfig
		expectedCount int
	}{
		{name: "default disabled", sampling: SamplingConfig{}, expectedCount: 3},
		{name: "explicit enabled", sampling: SamplingConfig{Enabled: true, Initial: 1, Thereafter: 100}, expectedCount: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			jsonPath := filepath.Join(tempDir, "service.jsonl")

			base, err := NewZapCoreBase(Config{
				Zap: ZapConfig{
					DisableCaller:     true,
					DisableStacktrace: true,
					ErrorOutputPaths:  []string{filepath.Join(tempDir, "errors.log")},
					Sampling:          tt.sampling,
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
					},
				},
			})
			if err != nil {
				t.Fatalf("build zap core base: %v", err)
			}
			t.Cleanup(base.Close)

			logger := zap.New(base.Core, base.Options()...)
			for index := 0; index < 3; index++ {
				logger.Info("sampled")
			}
			_ = logger.Sync()
			base.Close()

			jsonText := readFile(t, jsonPath)
			if count := strings.Count(jsonText, `"message":"sampled"`); count != tt.expectedCount {
				t.Fatalf("expected %d sampled log entries, got %d in %s", tt.expectedCount, count, jsonText)
			}
		})
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func readJSONLines(t *testing.T, path string) []map[string]any {
	t.Helper()
	text := strings.TrimSpace(readFile(t, path))
	if text == "" {
		t.Fatalf("expected %s to contain JSONL entries", path)
	}

	lines := strings.Split(text, "\n")
	entries := make([]map[string]any, 0, len(lines))
	for index, line := range lines {
		entry := map[string]any{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("parse JSONL line %d from %s: %v\n%s", index, path, err, line)
		}
		entries = append(entries, entry)
	}
	return entries
}

type verboseTestError struct{}

func (verboseTestError) Error() string {
	return "basic"
}

func (verboseTestError) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			_, _ = fmt.Fprint(state, "verbose")
			return
		}
		_, _ = fmt.Fprint(state, "basic")
	case 's':
		_, _ = fmt.Fprint(state, "basic")
	case 'q':
		_, _ = fmt.Fprint(state, `"basic"`)
	}
}
