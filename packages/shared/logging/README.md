<!--
  dox
  Copyright (C) 2026  OpenDox

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program. If not, see <http://www.gnu.org/licenses/>.

  @File    : packages/shared/logging/README.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-26
  @Modified: 2026-04-26
-->

# Shared Logging Model Contract

`packages/shared/logging` defines the shared Dox logging vocabulary and configuration contract.

This package maps the shared Dox logging configuration to zap and zapcore primitives for runtime integrations. It does not make zap the business logging API, and it does not initialize lumberjack, OpenTelemetry SDK providers, or the Dox logger facade.

## Boundary

The package owns stable names and configuration shapes for:

- resource identity;
- correlation identity;
- observability event classification;
- service-internal node fields;
- low-cardinality business tags;
- event facts and higher-cardinality fields;
- zap-facing, rotation, buffering, redaction, and OpenTelemetry configuration shapes.

The package must not:

- expose `*zap.Logger` or `zap.Field` as a business API;
- import lumberjack or OpenTelemetry SDK packages in the zap core base;
- implement file rotation, OTLP, or async queues in the zap core base;
- wire logging into server, scheduler, collector, compute, IAM, or HTTP middleware.

## Zap Core Base

The zap core base provides implementation helpers for runtime bootstrap code:

- Dox `Level` to `zapcore.Level` and `zap.AtomicLevel` mapping;
- symbolic encoder mapping for level, time, duration, caller, and logger name encoders;
- `zap.Config` mapping with sampling disabled unless explicitly enabled;
- enabled console and JSON core construction through zap output paths;
- zap options for development mode, caller, stacktrace, error output, and initial fields.

`DisableErrorVerbose` is applied by the core base so `zap.Error` fields keep the basic error string without adding zap's extra `errorVerbose` field.

File core declarations currently use zap's basic output path support. Rotation fields stay in the configuration contract for the follow-up lumberjack v2 sink issue.

## Resource

Resource fields answer who produced telemetry:

```text
service.namespace
service.name
service.instance.id
service.version
deployment.environment.name
cloud.region
cloud.availability_zone
k8s.cluster.name
k8s.namespace.name
dox.organization
dox.application
dox.runtime
```

`dox.runtime` identifies one Dox runtime: `server`, `scheduler`, `collector`, or `compute`.

`service.name` identifies a service capability and is not the same as runtime. For example, the `server` runtime can host the `iam` service.

## Correlation

Correlation fields connect one request, task, plugin execution, or cross-runtime chain:

```text
trace_id
span_id
trace_flags
request_id
correlation_id
job_id
task_id
workflow_id
plugin_id
plugin_run_id
```

`trace_id`, `span_id`, and `trace_flags` align with OpenTelemetry. `correlation_id` is Dox-owned and must survive across request, task, event, and plugin boundaries.

## Observability Events

`event.*` describes what the log record observes:

```text
event.name
event.dataset
event.category
event.type
event.action
event.outcome
```

Example:

```text
event.name = iam.login.rejected
event.dataset = dox.iam.security
event.category = authentication
event.type = denied
event.action = login
event.outcome = failure
```

There is no first-class `channel` field. Dataset, category, type, action, and outcome carry that classification.

## Tags Versus Fields

`tags` are low-cardinality business labels declared by the current node. They are suitable for dropdown filters, grouping, and alert routing.

Good tags:

```text
risk_level
login_method
credential_type
reject_reason
provider
queue
worker_pool
retryable
rate_limit_bucket
```

Do not put these into tags:

```text
runtime
service
env
region
component
operation
trace_id
request_id
correlation_id
user_id
account
duration_ms
error_message
```

Those values are resource fields, node fields, correlation fields, or event facts.

## IAM Login Rejected Example

```json
{
  "service.namespace": "dox",
  "service.name": "iam",
  "service.instance.id": "server-01",
  "service.version": "0.1.0",
  "dox.runtime": "server",
  "deployment.environment.name": "prod",
  "trace_id": "trace_001",
  "span_id": "span_auth_001",
  "request_id": "req_001",
  "correlation_id": "corr_001",
  "event.name": "iam.login.rejected",
  "event.dataset": "dox.iam.security",
  "event.category": "authentication",
  "event.type": "denied",
  "event.action": "login",
  "event.outcome": "failure",
  "component": "auth_service",
  "operation": "verify_credential",
  "tags": {
    "risk_level": "medium",
    "login_method": "password",
    "credential_type": "password",
    "reject_reason": "invalid_password"
  },
  "fields": {
    "account": "alice@example.com",
    "tenant_id": "tenant_a",
    "client_ip": "203.0.113.10",
    "failed_attempts": 3
  }
}
```

## Configuration Shape

The configuration contract includes:

- root logging settings;
- zap-facing configuration shape;
- per-core declarations for future console and file sinks;
- rotation configuration for future lumberjack v2 mapping;
- buffering and shutdown settings;
- redaction policy;
- OpenTelemetry propagation, traces, metrics, logs, OTLP exporter, and batch settings.

The default core declarations are:

```yaml
cores:
  - name: console
    enabled: true
    type: console
    encoding: console
    output_paths: ["stdout"]
    datasets: ["*"]
  - name: service-file
    enabled: true
    type: file
    encoding: json
    output_paths: ["logs/${service.namespace}-${service.name}.jsonl"]
    datasets: ["*"]
    rotation:
      driver: lumberjack
      enabled: true
      max_size_mb: 100
      max_backups: 10
      max_age_days: 14
      compress: true
      local_time: true
```

## Follow-Up Work

Separate issues should implement:

- JSONL file sink and lumberjack v2 rotation;
- Dox logger API and context correlation;
- OpenTelemetry SDK base;
- server setting and bootstrap integration;
- HTTP correlation middleware;
- IAM login chain sample;
- scheduler, collector, and compute integration;
- Filebeat, Fluent Bit, and OpenTelemetry Collector examples.
