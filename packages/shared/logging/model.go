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
 * @File    : model.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package logging

// Resource describes the service instance that produced telemetry.
type Resource struct {
	ServiceNamespace      string `json:"service.namespace,omitempty" yaml:"service.namespace,omitempty" mapstructure:"service.namespace"`
	ServiceName           string `json:"service.name,omitempty" yaml:"service.name,omitempty" mapstructure:"service.name"`
	ServiceInstanceID     string `json:"service.instance.id,omitempty" yaml:"service.instance.id,omitempty" mapstructure:"service.instance.id"`
	ServiceVersion        string `json:"service.version,omitempty" yaml:"service.version,omitempty" mapstructure:"service.version"`
	DeploymentEnvironment string `json:"deployment.environment.name,omitempty" yaml:"deployment.environment.name,omitempty" mapstructure:"deployment.environment.name"`
	CloudRegion           string `json:"cloud.region,omitempty" yaml:"cloud.region,omitempty" mapstructure:"cloud.region"`
	CloudAvailabilityZone string `json:"cloud.availability_zone,omitempty" yaml:"cloud.availability_zone,omitempty" mapstructure:"cloud.availability_zone"`
	K8sClusterName        string `json:"k8s.cluster.name,omitempty" yaml:"k8s.cluster.name,omitempty" mapstructure:"k8s.cluster.name"`
	K8sNamespaceName      string `json:"k8s.namespace.name,omitempty" yaml:"k8s.namespace.name,omitempty" mapstructure:"k8s.namespace.name"`
	DoxOrganization       string `json:"dox.organization,omitempty" yaml:"dox.organization,omitempty" mapstructure:"dox.organization"`
	DoxApplication        string `json:"dox.application,omitempty" yaml:"dox.application,omitempty" mapstructure:"dox.application"`
	DoxRuntime            string `json:"dox.runtime,omitempty" yaml:"dox.runtime,omitempty" mapstructure:"dox.runtime"`
}

// Correlation connects one request, task, plugin run, or event-driven chain.
type Correlation struct {
	TraceID       string `json:"trace_id,omitempty" yaml:"trace_id,omitempty" mapstructure:"trace_id"`
	SpanID        string `json:"span_id,omitempty" yaml:"span_id,omitempty" mapstructure:"span_id"`
	TraceFlags    string `json:"trace_flags,omitempty" yaml:"trace_flags,omitempty" mapstructure:"trace_flags"`
	RequestID     string `json:"request_id,omitempty" yaml:"request_id,omitempty" mapstructure:"request_id"`
	CorrelationID string `json:"correlation_id,omitempty" yaml:"correlation_id,omitempty" mapstructure:"correlation_id"`
	JobID         string `json:"job_id,omitempty" yaml:"job_id,omitempty" mapstructure:"job_id"`
	TaskID        string `json:"task_id,omitempty" yaml:"task_id,omitempty" mapstructure:"task_id"`
	WorkflowID    string `json:"workflow_id,omitempty" yaml:"workflow_id,omitempty" mapstructure:"workflow_id"`
	PluginID      string `json:"plugin_id,omitempty" yaml:"plugin_id,omitempty" mapstructure:"plugin_id"`
	PluginRunID   string `json:"plugin_run_id,omitempty" yaml:"plugin_run_id,omitempty" mapstructure:"plugin_run_id"`
}

// Event describes the observability event represented by a log record.
type Event struct {
	Name     string `json:"event.name,omitempty" yaml:"event.name,omitempty" mapstructure:"event.name"`
	Dataset  string `json:"event.dataset,omitempty" yaml:"event.dataset,omitempty" mapstructure:"event.dataset"`
	Category string `json:"event.category,omitempty" yaml:"event.category,omitempty" mapstructure:"event.category"`
	Type     string `json:"event.type,omitempty" yaml:"event.type,omitempty" mapstructure:"event.type"`
	Action   string `json:"event.action,omitempty" yaml:"event.action,omitempty" mapstructure:"event.action"`
	Outcome  string `json:"event.outcome,omitempty" yaml:"event.outcome,omitempty" mapstructure:"event.outcome"`
}

// Node describes where inside a service an event happened.
type Node struct {
	Component string `json:"component,omitempty" yaml:"component,omitempty" mapstructure:"component"`
	Operation string `json:"operation,omitempty" yaml:"operation,omitempty" mapstructure:"operation"`
}

// Tags are low-cardinality business node labels.
type Tags map[string]string

// Fields are event facts and higher-cardinality details.
type Fields map[string]any
