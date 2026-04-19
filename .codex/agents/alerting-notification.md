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

  @File    : .codex/agents/alerting-notification.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 7. Alerting And Notification Boundaries


Dox must distinguish system task alerts from business metric alerts.

These alert categories care about different problems, target different audiences, use different triggering strategies, and require different presentation. They must not be collapsed into a single undifferentiated alert model just because both are called alerts.

### System Task Alerts

System task alerts focus on system runtime state, task execution state, collection stability, computation stability, and platform integration health.

System task alerts primarily target:

- Developers.
- Operators.
- System administrators.
- Authorized advanced users.

System task alerts may include:

- Collection workers offline.
- Computation workers offline.
- Collection task failures.
- Computation task failures.
- Retry exhaustion.
- Task timeouts.
- Queue backlog.
- Platform credential invalidation.
- Third-party platform rate limits.
- Crawler abnormalities.
- Data write failures.
- `raw`, `ods`, `dwd`, `dws`, or `ads` generation failures.
- Plugin health check failures.
- Scheduling System abnormalities.

The goal of system task alerts is to help maintainers detect and recover system problems quickly.

System task alerts should preserve enough technical context, such as:

- `task_id`
- `plugin_id`
- `platform_id`
- `worker_id`
- `queue_name`
- `error_code`
- `error_message`
- `retry_count`
- `last_attempt_at`
- `trace_id` or `correlation_id`
- linked raw data or task records

### Business Metric Alerts

Business metric alerts focus on business changes, market changes, product performance, competitor changes, and operational opportunities.

Business metric alerts primarily target:

- Business users.
- Operators.
- Sellers.
- Managers.
- Decision-makers.

Business metric alerts may include:

- Keyword ranking changes.
- ASIN ranking changes.
- Search placement changes.
- Competitor changes.
- Market landscape changes.
- Traffic structure abnormalities.
- Listing performance abnormalities.
- Advertising keyword opportunities.
- Advertising data feedback abnormalities.
- Product performance abnormalities.
- Supplier-chain related risk inputs.
- Business metrics reaching thresholds.

The goal of business metric alerts is not to tell users where the system is broken. It is to tell users what happened in the business, why it matters, and what action may be useful.

Business metric alerts should include where practical:

- Business object.
- Metric name.
- Current value.
- Comparison value.
- Change magnitude.
- Time window.
- Data source.
- Triggering rule.
- Impact scope.
- Suggested action or next-step entry.

### Notification Channels

System task alerts and business metric alerts may share lower-level notification channels, such as:

- In-app notifications.
- Feishu.
- Email.
- Webhooks.
- AI bots.
- Future client channels.

However, the two alert categories must remain separate in model, permission, subscription, template, priority, and presentation.

The same notification channel may carry different alert types, but notification content must clearly identify the alert category. Business users should not receive large volumes of low-level system errors by default, and developers must not miss critical system failures.

### Notification Channel Plugins

Notification channels themselves should be pluginized.

In-app notifications, Feishu, email, webhooks, AI bots, and future clients should be treated as notification channel plugins instead of hard-coded branches inside the notification system.

Notification channel plugins can be introduced, configured, enabled, disabled, and removed independently.

A notification channel plugin should define:

- Channel identifier.
- Channel name.
- Channel type.
- Authentication method.
- Configuration fields.
- Supported message types.
- Supported template capabilities.
- Sending limits.
- Retry policy.
- Failure handling.
- Health check method.
- Permission requirements.
- Whether it supports interactive messages.
- Whether it supports AI-generated content.

The alerting system decides whether a notification is needed, who receives it, what content is sent, and what priority it has.

Notification channel plugins only send notifications to the corresponding client. They must not decide business alert rules or system alert rules.

The same alert may be sent to different notification channels according to user subscriptions, role, severity, and plugin configuration.

Notification channel plugins must follow the plugin architecture principles. They are part of Dox platform capabilities, but they must not be tightly coupled to the core alert model.

### Alert Subscriptions And Permissions

Alerts must be controlled by permissions and subscription rules.

System task alerts must not be sent to ordinary business users by default.

Business metric alerts must not be sent to all system administrators by default.

Users or roles should be able to subscribe, with proper authorization, to:

- Plugin-level alerts.
- Platform-level alerts.
- Task-level alerts.
- Business-object-level alerts.
- Metric-level alerts.
- Severity-level alerts.

### Alerts And Computation

Business metric alerts should usually be triggered from results produced by the Computation System.

The Web request path must not perform heavy computation to decide business alerts. Alerts that require statistics, aggregation, trend detection, or cross-period comparison should be produced as alert inputs by the Computation System and then evaluated by the alerting system.

System task alerts may be produced directly by the Scheduling System, Collection System, Computation System, or plugin health checks.

### AI And Alerts

AI may be used to explain alerts, summarize impact, generate suggested actions, or help users understand abnormalities.

AI must not replace alert facts. Alerts must have clear data sources, triggering rules, and business context before AI explains them. AI output must distinguish facts, inferences, and suggestions.

High-risk suggestions must not be executed automatically unless the user explicitly authorizes them and the system has permission control, audit logs, and rollback capability.

