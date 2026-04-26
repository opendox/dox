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
 * @File    : otel_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-27
 * @Modified: 2026-04-27
 */

package logging

import (
	"context"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestNewOpenTelemetryResourceMapsDoxResource(t *testing.T) {
	resource, err := NewOpenTelemetryResource(Resource{
		ServiceNamespace:      "dox",
		ServiceName:           "iam",
		ServiceInstanceID:     "server-01",
		ServiceVersion:        "0.1.0",
		DeploymentEnvironment: "prod",
		CloudRegion:           "us-east-1",
		CloudAvailabilityZone: "us-east-1a",
		K8sClusterName:        "dox-prod",
		K8sNamespaceName:      "dox-system",
		DoxOrganization:       "opendox",
		DoxApplication:        "dox",
		DoxRuntime:            "server",
	}, ResourceConfig{
		ServiceVersion: "0.2.0",
	})
	if err != nil {
		t.Fatalf("map OpenTelemetry resource: %v", err)
	}

	assertResourceString(t, resource, FieldServiceNamespace, "dox")
	assertResourceString(t, resource, FieldServiceName, "iam")
	assertResourceString(t, resource, FieldServiceInstanceID, "server-01")
	assertResourceString(t, resource, FieldServiceVersion, "0.2.0")
	assertResourceString(t, resource, FieldDeploymentEnvironmentName, "prod")
	assertResourceString(t, resource, FieldCloudRegion, "us-east-1")
	assertResourceString(t, resource, FieldCloudAvailabilityZone, "us-east-1a")
	assertResourceString(t, resource, FieldK8sClusterName, "dox-prod")
	assertResourceString(t, resource, FieldK8sNamespaceName, "dox-system")
	assertResourceString(t, resource, FieldDoxOrganization, "opendox")
	assertResourceString(t, resource, FieldDoxApplication, "dox")
	assertResourceString(t, resource, FieldDoxRuntime, "server")
}

func TestNewOpenTelemetryPropagatorMapsConfiguredFields(t *testing.T) {
	propagator := NewOpenTelemetryPropagator(OpenTelemetryPropagation{})
	assertContains(t, propagator.Fields(), "traceparent")
	assertContains(t, propagator.Fields(), "baggage")

	propagator = NewOpenTelemetryPropagator(OpenTelemetryPropagation{
		TraceContext: boolPtr(true),
		Baggage:      boolPtr(false),
	})
	assertContains(t, propagator.Fields(), "traceparent")
	assertNotContains(t, propagator.Fields(), "baggage")

	propagator = NewOpenTelemetryPropagator(OpenTelemetryPropagation{
		TraceContext: boolPtr(false),
		Baggage:      boolPtr(false),
	})
	if fields := propagator.Fields(); len(fields) != 0 {
		t.Fatalf("expected no-op propagator fields, got %#v", fields)
	}
}

func TestNewOpenTelemetryTraceSamplerMapsConfiguredSampler(t *testing.T) {
	tests := []struct {
		name          string
		config        TraceSamplerConfig
		decision      sdktrace.SamplingDecision
		checkDecision bool
		descriptor    string
	}{
		{
			name: "always on",
			config: TraceSamplerConfig{
				Type: TraceSamplerAlwaysOn,
			},
			decision:      sdktrace.RecordAndSample,
			checkDecision: true,
			descriptor:    "AlwaysOn",
		},
		{
			name: "always off",
			config: TraceSamplerConfig{
				Type: TraceSamplerAlwaysOff,
			},
			decision:      sdktrace.Drop,
			checkDecision: true,
			descriptor:    "AlwaysOff",
		},
		{
			name: "trace ratio",
			config: TraceSamplerConfig{
				Type:  TraceSamplerTraceIDRatio,
				Ratio: 0.5,
			},
			descriptor: "TraceIDRatioBased",
		},
		{
			name: "parent based trace ratio",
			config: TraceSamplerConfig{
				Type:  TraceSamplerParentBasedTraceIDRatio,
				Ratio: 0.5,
			},
			descriptor: "ParentBased",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler, err := NewOpenTelemetryTraceSampler(tt.config)
			if err != nil {
				t.Fatalf("map OpenTelemetry trace sampler: %v", err)
			}
			if !strings.Contains(sampler.Description(), tt.descriptor) {
				t.Fatalf("expected sampler description to contain %q, got %q", tt.descriptor, sampler.Description())
			}
			if !tt.checkDecision {
				return
			}
			result := sampler.ShouldSample(sdktrace.SamplingParameters{
				ParentContext: context.Background(),
				TraceID:       trace.TraceID{1, 2, 3},
				Name:          "test",
			})
			if result.Decision != tt.decision {
				t.Fatalf("expected decision %v, got %v", tt.decision, result.Decision)
			}
		})
	}
}

func TestNewOpenTelemetryTraceSamplerRejectsInvalidConfig(t *testing.T) {
	if _, err := NewOpenTelemetryTraceSampler(TraceSamplerConfig{Type: TraceSamplerType("tail")}); !hasValidationField(err, "otel.traces.sampler.type") {
		t.Fatalf("expected invalid sampler type validation error, got %v", err)
	}
	if _, err := NewOpenTelemetryTraceSampler(TraceSamplerConfig{Type: TraceSamplerTraceIDRatio, Ratio: 2}); !hasValidationField(err, "otel.traces.sampler.ratio") {
		t.Fatalf("expected invalid sampler ratio validation error, got %v", err)
	}
}

func TestNewOpenTelemetrySDKBaseRespectsProviderToggles(t *testing.T) {
	base, err := NewOpenTelemetrySDKBase(Config{
		Shutdown: ShutdownConfig{Timeout: 150 * time.Millisecond},
		OTel: OpenTelemetryConfig{
			Traces: OpenTelemetryTraces{
				Enabled: boolPtr(true),
			},
			Metrics: OpenTelemetrySignal{
				Enabled: boolPtr(false),
			},
			Logs: OpenTelemetrySignal{
				Enabled: boolPtr(true),
			},
		},
	}, Resource{
		ServiceNamespace: "dox",
		ServiceName:      "iam",
		DoxRuntime:       "server",
	})
	if err != nil {
		t.Fatalf("build OpenTelemetry SDK base: %v", err)
	}

	if base.Resource == nil || base.Propagator == nil {
		t.Fatalf("expected resource and propagator, got %+v", base)
	}
	if base.TracerProvider == nil {
		t.Fatal("expected traces provider to be enabled")
	}
	if base.MeterProvider != nil {
		t.Fatal("expected metrics provider to be disabled")
	}
	if base.LoggerProvider == nil {
		t.Fatal("expected logs provider to be enabled")
	}
	if base.ShutdownTimeout != 150*time.Millisecond {
		t.Fatalf("expected shutdown timeout to map, got %s", base.ShutdownTimeout)
	}
	if err := base.ForceFlush(nil); err != nil {
		t.Fatalf("force flush OpenTelemetry providers: %v", err)
	}
	if err := base.Shutdown(nil); err != nil {
		t.Fatalf("shutdown OpenTelemetry providers: %v", err)
	}
}

func TestNewOpenTelemetrySDKBaseRespectsRootDisabledToggle(t *testing.T) {
	base, err := NewOpenTelemetrySDKBase(Config{
		OTel: OpenTelemetryConfig{
			Enabled: boolPtr(false),
		},
	}, Resource{
		ServiceName: "iam",
	})
	if err != nil {
		t.Fatalf("build disabled OpenTelemetry SDK base: %v", err)
	}

	if base.Resource == nil || base.Propagator == nil {
		t.Fatalf("expected disabled base to still expose resource and propagator, got %+v", base)
	}
	if base.TracerProvider != nil || base.MeterProvider != nil || base.LoggerProvider != nil {
		t.Fatalf("expected disabled OpenTelemetry providers, got %+v", base)
	}
	if fields := base.Propagator.Fields(); len(fields) != 0 {
		t.Fatalf("expected disabled OpenTelemetry propagator to be no-op, got %#v", fields)
	}
}

func TestNewOpenTelemetrySDKBaseRejectsEnabledOTLPExporter(t *testing.T) {
	_, err := NewOpenTelemetrySDKBase(Config{
		OTel: OpenTelemetryConfig{
			Exporter: OpenTelemetryExporter{
				OTLP: OTLPExporterConfig{
					Enabled: true,
				},
			},
		},
	}, Resource{})
	if !hasValidationField(err, "otel.exporter.otlp.enabled") {
		t.Fatalf("expected unsupported OTLP exporter validation error, got %v", err)
	}
}

func TestNewOpenTelemetrySDKBaseDoesNotInstallGlobals(t *testing.T) {
	beforeTracer := reflect.TypeOf(otel.GetTracerProvider()).String()
	beforeMeter := reflect.TypeOf(otel.GetMeterProvider()).String()
	beforePropagatorFields := append([]string(nil), otel.GetTextMapPropagator().Fields()...)

	base, err := NewOpenTelemetrySDKBase(Config{}, Resource{
		ServiceName: "iam",
	})
	if err != nil {
		t.Fatalf("build OpenTelemetry SDK base: %v", err)
	}
	t.Cleanup(func() {
		if err := base.Shutdown(context.Background()); err != nil {
			t.Fatalf("shutdown OpenTelemetry SDK base: %v", err)
		}
	})

	if got := reflect.TypeOf(otel.GetTracerProvider()).String(); got != beforeTracer {
		t.Fatalf("expected global tracer provider type %q to stay unchanged, got %q", beforeTracer, got)
	}
	if got := reflect.TypeOf(otel.GetMeterProvider()).String(); got != beforeMeter {
		t.Fatalf("expected global meter provider type %q to stay unchanged, got %q", beforeMeter, got)
	}
	if got := otel.GetTextMapPropagator().Fields(); !slices.Equal(got, beforePropagatorFields) {
		t.Fatalf("expected global propagator fields %#v to stay unchanged, got %#v", beforePropagatorFields, got)
	}
}

func assertResourceString(t *testing.T, resource interface {
	Set() *attribute.Set
}, key string, expected string) {
	t.Helper()

	value, ok := resource.Set().Value(attribute.Key(key))
	if !ok {
		t.Fatalf("expected resource attribute %q", key)
	}
	if got := value.AsString(); got != expected {
		t.Fatalf("expected resource attribute %q to be %q, got %q", key, expected, got)
	}
}

func assertContains(t *testing.T, values []string, expected string) {
	t.Helper()
	if contains(values, expected) {
		return
	}
	t.Fatalf("expected %#v to contain %q", values, expected)
}

func assertNotContains(t *testing.T, values []string, expected string) {
	t.Helper()
	if !contains(values, expected) {
		return
	}
	t.Fatalf("expected %#v not to contain %q", values, expected)
}

var _ propagation.TextMapPropagator = NewOpenTelemetryPropagator(OpenTelemetryPropagation{})
