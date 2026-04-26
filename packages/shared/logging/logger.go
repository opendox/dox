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
 * @File    : logger.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-27
 * @Modified: 2026-04-27
 */

package logging

import (
	"context"
	"errors"
	"sort"

	"go.uber.org/zap"
)

// Logger is the Dox-owned business logging facade.
type Logger interface {
	Debug(ctx context.Context, message string, attrs ...Attr)
	Info(ctx context.Context, message string, attrs ...Attr)
	Warn(ctx context.Context, message string, attrs ...Attr)
	Error(ctx context.Context, message string, attrs ...Attr)
	DPanic(ctx context.Context, message string, attrs ...Attr)
	Panic(ctx context.Context, message string, attrs ...Attr)
	Fatal(ctx context.Context, message string, attrs ...Attr)
	Named(name string) Logger
	With(attrs ...Attr) Logger
	Sync() error
}

// Attr attaches Dox logging model values to a logger or log call.
type Attr struct {
	apply func(*logValues)
}

// ResourceAttr attaches resource identity fields.
func ResourceAttr(resource Resource) Attr {
	return Attr{apply: func(values *logValues) {
		values.resource = mergeResource(values.resource, resource)
	}}
}

// CorrelationAttr attaches correlation fields.
func CorrelationAttr(correlation Correlation) Attr {
	return Attr{apply: func(values *logValues) {
		values.correlation = MergeCorrelation(values.correlation, correlation)
	}}
}

// EventAttr attaches observability event fields.
func EventAttr(event Event) Attr {
	return Attr{apply: func(values *logValues) {
		values.event = mergeEvent(values.event, event)
	}}
}

// NodeAttr attaches service-internal node fields.
func NodeAttr(node Node) Attr {
	return Attr{apply: func(values *logValues) {
		values.node = mergeNode(values.node, node)
	}}
}

// TagsAttr attaches low-cardinality business tags.
func TagsAttr(tags Tags) Attr {
	return Attr{apply: func(values *logValues) {
		if len(tags) == 0 {
			return
		}
		if values.tags == nil {
			values.tags = Tags{}
		}
		for key, value := range tags {
			if key == "" {
				continue
			}
			values.tags[key] = value
		}
	}}
}

// FieldsAttr attaches event facts and higher-cardinality fields.
func FieldsAttr(fields Fields) Attr {
	return Attr{apply: func(values *logValues) {
		if len(fields) == 0 {
			return
		}
		if values.fields == nil {
			values.fields = Fields{}
		}
		for key, value := range fields {
			if key == "" {
				continue
			}
			values.fields[key] = value
		}
	}}
}

// FieldAttr attaches one event fact under the Dox fields object.
func FieldAttr(key string, value any) Attr {
	return Attr{apply: func(values *logValues) {
		if key == "" {
			return
		}
		if values.fields == nil {
			values.fields = Fields{}
		}
		values.fields[key] = value
	}}
}

// ErrorAttr attaches an error to the log record.
func ErrorAttr(err error) Attr {
	return Attr{apply: func(values *logValues) {
		values.err = err
	}}
}

// NewLogger wraps an initialized zap core base with the Dox logger facade.
func NewLogger(base *ZapCoreBase, attrs ...Attr) (Logger, error) {
	if base == nil {
		return nil, errors.New("logging: zap core base must not be nil")
	}

	options := base.Options()
	options = append(options, zap.AddCallerSkip(1))

	logger := zap.New(base.Core, options...)
	values := logValues{}
	values.apply(attrs...)
	return doxLogger{
		logger: logger,
		values: values,
	}, nil
}

type doxLogger struct {
	logger *zap.Logger
	values logValues
}

func (l doxLogger) Debug(ctx context.Context, message string, attrs ...Attr) {
	l.write(ctx, LevelDebug, message, attrs...)
}

func (l doxLogger) Info(ctx context.Context, message string, attrs ...Attr) {
	l.write(ctx, LevelInfo, message, attrs...)
}

func (l doxLogger) Warn(ctx context.Context, message string, attrs ...Attr) {
	l.write(ctx, LevelWarn, message, attrs...)
}

func (l doxLogger) Error(ctx context.Context, message string, attrs ...Attr) {
	l.write(ctx, LevelError, message, attrs...)
}

func (l doxLogger) DPanic(ctx context.Context, message string, attrs ...Attr) {
	l.write(ctx, LevelDPanic, message, attrs...)
}

func (l doxLogger) Panic(ctx context.Context, message string, attrs ...Attr) {
	l.write(ctx, LevelPanic, message, attrs...)
}

func (l doxLogger) Fatal(ctx context.Context, message string, attrs ...Attr) {
	l.write(ctx, LevelFatal, message, attrs...)
}

func (l doxLogger) Named(name string) Logger {
	if l.logger == nil || name == "" {
		return l
	}
	l.logger = l.logger.Named(name)
	return l
}

func (l doxLogger) With(attrs ...Attr) Logger {
	l.values.apply(attrs...)
	return l
}

func (l doxLogger) Sync() error {
	if l.logger == nil {
		return nil
	}
	return l.logger.Sync()
}

func (l doxLogger) write(ctx context.Context, level Level, message string, attrs ...Attr) {
	if l.logger == nil {
		return
	}

	values := l.values.clone()
	if correlation, ok := CorrelationFromContext(ctx); ok {
		values.correlation = MergeCorrelation(values.correlation, correlation)
	}
	values.apply(attrs...)
	fields := values.zapFields()

	switch level {
	case LevelDebug:
		l.logger.Debug(message, fields...)
	case LevelInfo:
		l.logger.Info(message, fields...)
	case LevelWarn:
		l.logger.Warn(message, fields...)
	case LevelError:
		l.logger.Error(message, fields...)
	case LevelDPanic:
		l.logger.DPanic(message, fields...)
	case LevelPanic:
		l.logger.Panic(message, fields...)
	case LevelFatal:
		l.logger.Fatal(message, fields...)
	}
}

type logValues struct {
	resource    Resource
	correlation Correlation
	event       Event
	node        Node
	tags        Tags
	fields      Fields
	err         error
}

func (v *logValues) apply(attrs ...Attr) {
	for _, attr := range attrs {
		if attr.apply != nil {
			attr.apply(v)
		}
	}
}

func (v logValues) clone() logValues {
	return logValues{
		resource:    v.resource,
		correlation: v.correlation,
		event:       v.event,
		node:        v.node,
		tags:        copyTags(v.tags),
		fields:      copyFields(v.fields),
		err:         v.err,
	}
}

func (v logValues) zapFields() []zap.Field {
	fields := make([]zap.Field, 0, 32)
	fields = appendResourceFields(fields, v.resource)
	fields = appendCorrelationFields(fields, v.correlation)
	fields = appendEventFields(fields, v.event)
	fields = appendNodeFields(fields, v.node)
	if len(v.tags) > 0 {
		fields = append(fields, zap.Any(FieldTags, sortedTags(v.tags)))
	}
	if len(v.fields) > 0 {
		fields = append(fields, zap.Any(FieldFields, sortedFields(v.fields)))
	}
	if v.err != nil {
		fields = append(fields, zap.Error(v.err))
	}
	return fields
}

func appendResourceFields(fields []zap.Field, resource Resource) []zap.Field {
	if resource.ServiceNamespace != "" {
		fields = append(fields, zap.String(FieldServiceNamespace, resource.ServiceNamespace))
	}
	if resource.ServiceName != "" {
		fields = append(fields, zap.String(FieldServiceName, resource.ServiceName))
	}
	if resource.ServiceInstanceID != "" {
		fields = append(fields, zap.String(FieldServiceInstanceID, resource.ServiceInstanceID))
	}
	if resource.ServiceVersion != "" {
		fields = append(fields, zap.String(FieldServiceVersion, resource.ServiceVersion))
	}
	if resource.DeploymentEnvironment != "" {
		fields = append(fields, zap.String(FieldDeploymentEnvironmentName, resource.DeploymentEnvironment))
	}
	if resource.CloudRegion != "" {
		fields = append(fields, zap.String(FieldCloudRegion, resource.CloudRegion))
	}
	if resource.CloudAvailabilityZone != "" {
		fields = append(fields, zap.String(FieldCloudAvailabilityZone, resource.CloudAvailabilityZone))
	}
	if resource.K8sClusterName != "" {
		fields = append(fields, zap.String(FieldK8sClusterName, resource.K8sClusterName))
	}
	if resource.K8sNamespaceName != "" {
		fields = append(fields, zap.String(FieldK8sNamespaceName, resource.K8sNamespaceName))
	}
	if resource.DoxOrganization != "" {
		fields = append(fields, zap.String(FieldDoxOrganization, resource.DoxOrganization))
	}
	if resource.DoxApplication != "" {
		fields = append(fields, zap.String(FieldDoxApplication, resource.DoxApplication))
	}
	if resource.DoxRuntime != "" {
		fields = append(fields, zap.String(FieldDoxRuntime, resource.DoxRuntime))
	}
	return fields
}

func appendCorrelationFields(fields []zap.Field, correlation Correlation) []zap.Field {
	if correlation.TraceID != "" {
		fields = append(fields, zap.String(FieldTraceID, correlation.TraceID))
	}
	if correlation.SpanID != "" {
		fields = append(fields, zap.String(FieldSpanID, correlation.SpanID))
	}
	if correlation.TraceFlags != "" {
		fields = append(fields, zap.String(FieldTraceFlags, correlation.TraceFlags))
	}
	if correlation.RequestID != "" {
		fields = append(fields, zap.String(FieldRequestID, correlation.RequestID))
	}
	if correlation.CorrelationID != "" {
		fields = append(fields, zap.String(FieldCorrelationID, correlation.CorrelationID))
	}
	if correlation.JobID != "" {
		fields = append(fields, zap.String(FieldJobID, correlation.JobID))
	}
	if correlation.TaskID != "" {
		fields = append(fields, zap.String(FieldTaskID, correlation.TaskID))
	}
	if correlation.WorkflowID != "" {
		fields = append(fields, zap.String(FieldWorkflowID, correlation.WorkflowID))
	}
	if correlation.PluginID != "" {
		fields = append(fields, zap.String(FieldPluginID, correlation.PluginID))
	}
	if correlation.PluginRunID != "" {
		fields = append(fields, zap.String(FieldPluginRunID, correlation.PluginRunID))
	}
	return fields
}

func appendEventFields(fields []zap.Field, event Event) []zap.Field {
	if event.Name != "" {
		fields = append(fields, zap.String(FieldEventName, event.Name))
	}
	if event.Dataset != "" {
		fields = append(fields, zap.String(FieldEventDataset, event.Dataset))
	}
	if event.Category != "" {
		fields = append(fields, zap.String(FieldEventCategory, event.Category))
	}
	if event.Type != "" {
		fields = append(fields, zap.String(FieldEventType, event.Type))
	}
	if event.Action != "" {
		fields = append(fields, zap.String(FieldEventAction, event.Action))
	}
	if event.Outcome != "" {
		fields = append(fields, zap.String(FieldEventOutcome, event.Outcome))
	}
	return fields
}

func appendNodeFields(fields []zap.Field, node Node) []zap.Field {
	if node.Component != "" {
		fields = append(fields, zap.String(FieldComponent, node.Component))
	}
	if node.Operation != "" {
		fields = append(fields, zap.String(FieldOperation, node.Operation))
	}
	return fields
}

func mergeResource(base Resource, overlay Resource) Resource {
	if overlay.ServiceNamespace != "" {
		base.ServiceNamespace = overlay.ServiceNamespace
	}
	if overlay.ServiceName != "" {
		base.ServiceName = overlay.ServiceName
	}
	if overlay.ServiceInstanceID != "" {
		base.ServiceInstanceID = overlay.ServiceInstanceID
	}
	if overlay.ServiceVersion != "" {
		base.ServiceVersion = overlay.ServiceVersion
	}
	if overlay.DeploymentEnvironment != "" {
		base.DeploymentEnvironment = overlay.DeploymentEnvironment
	}
	if overlay.CloudRegion != "" {
		base.CloudRegion = overlay.CloudRegion
	}
	if overlay.CloudAvailabilityZone != "" {
		base.CloudAvailabilityZone = overlay.CloudAvailabilityZone
	}
	if overlay.K8sClusterName != "" {
		base.K8sClusterName = overlay.K8sClusterName
	}
	if overlay.K8sNamespaceName != "" {
		base.K8sNamespaceName = overlay.K8sNamespaceName
	}
	if overlay.DoxOrganization != "" {
		base.DoxOrganization = overlay.DoxOrganization
	}
	if overlay.DoxApplication != "" {
		base.DoxApplication = overlay.DoxApplication
	}
	if overlay.DoxRuntime != "" {
		base.DoxRuntime = overlay.DoxRuntime
	}
	return base
}

func mergeEvent(base Event, overlay Event) Event {
	if overlay.Name != "" {
		base.Name = overlay.Name
	}
	if overlay.Dataset != "" {
		base.Dataset = overlay.Dataset
	}
	if overlay.Category != "" {
		base.Category = overlay.Category
	}
	if overlay.Type != "" {
		base.Type = overlay.Type
	}
	if overlay.Action != "" {
		base.Action = overlay.Action
	}
	if overlay.Outcome != "" {
		base.Outcome = overlay.Outcome
	}
	return base
}

func mergeNode(base Node, overlay Node) Node {
	if overlay.Component != "" {
		base.Component = overlay.Component
	}
	if overlay.Operation != "" {
		base.Operation = overlay.Operation
	}
	return base
}

func copyTags(tags Tags) Tags {
	if tags == nil {
		return nil
	}
	copied := make(Tags, len(tags))
	for key, value := range tags {
		copied[key] = value
	}
	return copied
}

func copyFields(fields Fields) Fields {
	if fields == nil {
		return nil
	}
	copied := make(Fields, len(fields))
	for key, value := range fields {
		copied[key] = value
	}
	return copied
}

func sortedTags(tags Tags) Tags {
	if len(tags) == 0 {
		return nil
	}
	keys := make([]string, 0, len(tags))
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sorted := make(Tags, len(tags))
	for _, key := range keys {
		sorted[key] = tags[key]
	}
	return sorted
}

func sortedFields(fields Fields) Fields {
	if len(fields) == 0 {
		return nil
	}
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sorted := make(Fields, len(fields))
	for _, key := range keys {
		sorted[key] = fields[key]
	}
	return sorted
}
