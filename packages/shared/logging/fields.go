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
 * @File    : fields.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package logging

const (
	// FieldServiceNamespace is the OpenTelemetry service namespace field.
	FieldServiceNamespace = "service.namespace"
	// FieldServiceName is the OpenTelemetry service name field.
	FieldServiceName = "service.name"
	// FieldServiceInstanceID is the OpenTelemetry service instance field.
	FieldServiceInstanceID = "service.instance.id"
	// FieldServiceVersion is the OpenTelemetry service version field.
	FieldServiceVersion = "service.version"
	// FieldDeploymentEnvironmentName identifies the deployment environment.
	FieldDeploymentEnvironmentName = "deployment.environment.name"
	// FieldCloudRegion identifies the cloud region.
	FieldCloudRegion = "cloud.region"
	// FieldCloudAvailabilityZone identifies the cloud availability zone.
	FieldCloudAvailabilityZone = "cloud.availability_zone"
	// FieldK8sClusterName identifies the Kubernetes cluster.
	FieldK8sClusterName = "k8s.cluster.name"
	// FieldK8sNamespaceName identifies the Kubernetes namespace.
	FieldK8sNamespaceName = "k8s.namespace.name"
	// FieldDoxOrganization identifies the Dox owning organization.
	FieldDoxOrganization = "dox.organization"
	// FieldDoxApplication identifies the Dox application family.
	FieldDoxApplication = "dox.application"
	// FieldDoxRuntime identifies the Dox runtime system.
	FieldDoxRuntime = "dox.runtime"
)

const (
	// FieldTraceID identifies an OpenTelemetry trace.
	FieldTraceID = "trace_id"
	// FieldSpanID identifies an OpenTelemetry span.
	FieldSpanID = "span_id"
	// FieldTraceFlags carries OpenTelemetry trace flags.
	FieldTraceFlags = "trace_flags"
	// FieldRequestID identifies an HTTP or API request.
	FieldRequestID = "request_id"
	// FieldCorrelationID identifies a Dox business execution chain.
	FieldCorrelationID = "correlation_id"
	// FieldJobID identifies a Dox job.
	FieldJobID = "job_id"
	// FieldTaskID identifies a Dox task.
	FieldTaskID = "task_id"
	// FieldWorkflowID identifies a Dox workflow.
	FieldWorkflowID = "workflow_id"
	// FieldPluginID identifies a Dox plugin.
	FieldPluginID = "plugin_id"
	// FieldPluginRunID identifies one plugin execution.
	FieldPluginRunID = "plugin_run_id"
)

const (
	// FieldEventName identifies the observability event name.
	FieldEventName = "event.name"
	// FieldEventDataset identifies the observability event dataset.
	FieldEventDataset = "event.dataset"
	// FieldEventCategory identifies the observability event category.
	FieldEventCategory = "event.category"
	// FieldEventType identifies the observability event type.
	FieldEventType = "event.type"
	// FieldEventAction identifies the observability event action.
	FieldEventAction = "event.action"
	// FieldEventOutcome identifies the observability event outcome.
	FieldEventOutcome = "event.outcome"
)

const (
	// FieldComponent identifies the service-internal component.
	FieldComponent = "component"
	// FieldOperation identifies the current operation.
	FieldOperation = "operation"
	// FieldTags stores low-cardinality business node labels.
	FieldTags = "tags"
	// FieldFields stores event facts and higher-cardinality details.
	FieldFields = "fields"
)
