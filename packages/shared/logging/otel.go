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
 * @File    : otel.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-27
 * @Modified: 2026-04-27
 */

package logging

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// OpenTelemetrySDKBase carries OpenTelemetry SDK primitives for runtime
// integrations.
//
// Business code should depend on Dox logging types and Logger instead of using
// these OpenTelemetry primitives directly.
type OpenTelemetrySDKBase struct {
	Resource        *sdkresource.Resource
	Propagator      propagation.TextMapPropagator
	TracerProvider  *sdktrace.TracerProvider
	MeterProvider   *sdkmetric.MeterProvider
	LoggerProvider  *sdklog.LoggerProvider
	ShutdownTimeout time.Duration
}

// ForceFlush flushes enabled OpenTelemetry providers using the configured
// logging shutdown timeout.
func (b *OpenTelemetrySDKBase) ForceFlush(ctx context.Context) error {
	if b == nil {
		return nil
	}

	ctx, cancel := b.withShutdownTimeout(ctx)
	defer cancel()

	var err error
	if b.TracerProvider != nil {
		err = errors.Join(err, b.TracerProvider.ForceFlush(ctx))
	}
	if b.MeterProvider != nil {
		err = errors.Join(err, b.MeterProvider.ForceFlush(ctx))
	}
	if b.LoggerProvider != nil {
		err = errors.Join(err, b.LoggerProvider.ForceFlush(ctx))
	}
	return err
}

// Shutdown shuts down enabled OpenTelemetry providers using the configured
// logging shutdown timeout.
func (b *OpenTelemetrySDKBase) Shutdown(ctx context.Context) error {
	if b == nil {
		return nil
	}

	ctx, cancel := b.withShutdownTimeout(ctx)
	defer cancel()

	var err error
	if b.TracerProvider != nil {
		err = errors.Join(err, b.TracerProvider.Shutdown(ctx))
	}
	if b.MeterProvider != nil {
		err = errors.Join(err, b.MeterProvider.Shutdown(ctx))
	}
	if b.LoggerProvider != nil {
		err = errors.Join(err, b.LoggerProvider.Shutdown(ctx))
	}
	return err
}

// NewOpenTelemetrySDKBase maps Dox logging config to OpenTelemetry SDK
// providers without installing any runtime globals.
func NewOpenTelemetrySDKBase(config Config, resource Resource) (*OpenTelemetrySDKBase, error) {
	normalized, err := normalizeOpenTelemetryConfig(config)
	if err != nil {
		return nil, err
	}

	otelResource, err := NewOpenTelemetryResource(resource, normalized.Resource)
	if err != nil {
		return nil, err
	}

	base := &OpenTelemetrySDKBase{
		Resource:        otelResource,
		ShutdownTimeout: normalized.Shutdown.Timeout,
	}

	if !boolValue(normalized.OTel.Enabled) {
		base.Propagator = propagation.NewCompositeTextMapPropagator()
		return base, nil
	}

	base.Propagator = NewOpenTelemetryPropagator(normalized.OTel.Propagation)

	if boolValue(normalized.OTel.Traces.Enabled) {
		sampler, err := NewOpenTelemetryTraceSampler(normalized.OTel.Traces.Sampler)
		if err != nil {
			return nil, err
		}
		base.TracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithResource(otelResource),
			sdktrace.WithSampler(sampler),
		)
	}

	if boolValue(normalized.OTel.Metrics.Enabled) {
		base.MeterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(otelResource),
		)
	}

	if boolValue(normalized.OTel.Logs.Enabled) {
		base.LoggerProvider = sdklog.NewLoggerProvider(
			sdklog.WithResource(otelResource),
		)
	}

	return base, nil
}

// NewOpenTelemetryResource maps the Dox resource model to OpenTelemetry
// resource attributes and merges them over the OpenTelemetry SDK defaults.
func NewOpenTelemetryResource(model Resource, config ResourceConfig) (*sdkresource.Resource, error) {
	attrs := openTelemetryResourceAttributes(model, config)
	doxResource := sdkresource.NewSchemaless(attrs...)
	otelResource, err := sdkresource.Merge(sdkresource.Default(), doxResource)
	if err != nil {
		return nil, fmt.Errorf("logging: build OpenTelemetry resource: %w", err)
	}
	return otelResource, nil
}

// NewOpenTelemetryPropagator maps Dox propagation settings to an OpenTelemetry
// composite TextMapPropagator.
func NewOpenTelemetryPropagator(config OpenTelemetryPropagation) propagation.TextMapPropagator {
	config.Default()

	propagators := make([]propagation.TextMapPropagator, 0, 2)
	if boolValue(config.TraceContext) {
		propagators = append(propagators, propagation.TraceContext{})
	}
	if boolValue(config.Baggage) {
		propagators = append(propagators, propagation.Baggage{})
	}
	return propagation.NewCompositeTextMapPropagator(propagators...)
}

// NewOpenTelemetryTraceSampler maps a Dox trace sampler config to an
// OpenTelemetry trace sampler.
func NewOpenTelemetryTraceSampler(config TraceSamplerConfig) (sdktrace.Sampler, error) {
	config.Default()

	if !config.Type.IsValid() {
		return nil, validationError("otel.traces.sampler.type", "trace sampler type is not supported")
	}
	if config.Ratio < 0 || config.Ratio > 1 {
		return nil, validationError("otel.traces.sampler.ratio", "trace sampler ratio must be between 0 and 1")
	}

	switch config.Type {
	case TraceSamplerAlwaysOn:
		return sdktrace.AlwaysSample(), nil
	case TraceSamplerAlwaysOff:
		return sdktrace.NeverSample(), nil
	case TraceSamplerTraceIDRatio:
		return sdktrace.TraceIDRatioBased(config.Ratio), nil
	case TraceSamplerParentBasedTraceIDRatio:
		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(config.Ratio)), nil
	default:
		return nil, validationError("otel.traces.sampler.type", "trace sampler type is not supported")
	}
}

func normalizeOpenTelemetryConfig(config Config) (Config, error) {
	if err := config.Default(); err != nil {
		return Config{}, err
	}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}
	if boolValue(config.OTel.Enabled) && config.OTel.Exporter.OTLP.Enabled {
		return Config{}, validationError("otel.exporter.otlp.enabled", "OTLP exporter setup is not supported by the SDK base")
	}
	return config, nil
}

func openTelemetryResourceAttributes(model Resource, config ResourceConfig) []attribute.KeyValue {
	serviceVersion := model.ServiceVersion
	if strings.TrimSpace(config.ServiceVersion) != "" {
		serviceVersion = config.ServiceVersion
	}

	attrs := make([]attribute.KeyValue, 0, 12)
	attrs = appendStringAttribute(attrs, FieldServiceNamespace, model.ServiceNamespace)
	attrs = appendStringAttribute(attrs, FieldServiceName, model.ServiceName)
	attrs = appendStringAttribute(attrs, FieldServiceInstanceID, model.ServiceInstanceID)
	attrs = appendStringAttribute(attrs, FieldServiceVersion, serviceVersion)
	attrs = appendStringAttribute(attrs, FieldDeploymentEnvironmentName, model.DeploymentEnvironment)
	attrs = appendStringAttribute(attrs, FieldCloudRegion, model.CloudRegion)
	attrs = appendStringAttribute(attrs, FieldCloudAvailabilityZone, model.CloudAvailabilityZone)
	attrs = appendStringAttribute(attrs, FieldK8sClusterName, model.K8sClusterName)
	attrs = appendStringAttribute(attrs, FieldK8sNamespaceName, model.K8sNamespaceName)
	attrs = appendStringAttribute(attrs, FieldDoxOrganization, model.DoxOrganization)
	attrs = appendStringAttribute(attrs, FieldDoxApplication, model.DoxApplication)
	attrs = appendStringAttribute(attrs, FieldDoxRuntime, model.DoxRuntime)
	return attrs
}

func appendStringAttribute(attrs []attribute.KeyValue, key string, value string) []attribute.KeyValue {
	if strings.TrimSpace(value) == "" {
		return attrs
	}
	return append(attrs, attribute.String(key, value))
}

func (b *OpenTelemetrySDKBase) withShutdownTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	if b.ShutdownTimeout <= 0 {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, b.ShutdownTimeout)
}
