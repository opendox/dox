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

  @File    : .codex/agents/plugin-architecture.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 5. Plugin Architecture


Dox platform capabilities must enter the system through plugins.

A plugin is not just a code import, a frontend menu item, or a backend endpoint. A plugin is a complete platform capability definition that describes how an external platform, data source, business capability, or system extension integrates with Dox.

Plugins may represent:

- Amazon SP-API.
- Amazon Ads API.
- Lingxing API.
- Feishu API.
- SIF keyword data.
- SellerSprite-style market data.
- Amazon retail frontend collection.
- Future third-party platforms, internal platforms, or proprietary business capabilities.

### Plugin Definition

A plugin should describe the capability it brings as explicitly as possible. Plugin behavior must not be scattered across the Web, Scheduling, Collection, and Computation systems without a clear definition.

A plugin definition may include:

- Plugin identifier.
- Plugin name.
- Plugin version.
- Plugin type.
- Platform authentication method.
- Configuration fields.
- Permission points.
- Menu entries.
- Route entries.
- Task types.
- Collection capabilities.
- Computation capabilities.
- Data models.
- Data-layer write conventions.
- Alert capabilities.
- Notification capabilities.
- Report or page entries.
- Enablement conditions.
- Plugin dependencies.
- Health check method.

Plugin definitions should be as declarative as practical. The system should be able to understand from the plugin definition what configuration is required, what tasks can be produced, what pages can be enabled, what data can be written, and what alerts can be triggered.

### Plugin Enablement Flow

When a user configures and enables a platform plugin in the Web System, the system should activate related capabilities according to the plugin definition.

Enabling a plugin may include:

- The Web System stores platform credentials and plugin configuration.
- The Web System issues menus, routes, and permission points according to user, role, permission, and plugin state.
- The Scheduling System registers task types that the plugin can schedule.
- The Collection System enables or connects the plugin's collection capability.
- The Computation System enables or connects the plugin's computation capability.
- Related collection tasks, computation tasks, and alert rules become schedulable or configurable.

Users do not need to understand code, processes, workers, or deployment details. After users complete authentication and configuration in the Web System, the plugin capability should become available through the Dox system.

### Plugins And Frontend

The frontend is a shell for plugin capabilities, not the decision-maker for whether a plugin is available.

The frontend must not hard-code whether a platform is available. Menus, routes, permissions, plugin entries, and some page capabilities should be issued by the Web System according to user, role, permission, and plugin enablement state.

The frontend may provide:

- Shared layout.
- Shared page containers.
- Shared forms.
- Shared tables.
- Shared task views.
- Shared report views.
- Plugin configuration page hosting.
- Plugin analytics page hosting.

Whether a specific plugin is visible, which entries are shown, and which actions the user may perform must be decided by the backend.

### Plugins And Scheduling

The Scheduling System must understand plugin task capabilities.

After a plugin is enabled, the Scheduling System should know:

- Which task types the plugin provides.
- Which configuration each task requires.
- Which Collection capability executes each task.
- Whether a task may produce computation tasks.
- Whether a task can be scheduled periodically.
- Whether a task can be triggered manually.
- Whether a task has retry, timeout, priority, and concurrency limits.
- Which events or follow-up tasks may be triggered by task results.

Scheduling policy must not be hard-coded into the Collection System. The Collection System executes tasks. Task chains and next actions are decided by the Scheduling System according to plugin definitions and system policy.

### Plugins And Collection

The Collection System executes plugin collection capabilities.

A collection plugin should define:

- Which APIs it calls or which pages it collects.
- Which credentials it requires.
- Which runtime configuration it requires.
- Which `raw` data it produces.
- Which `ods` data it produces.
- How it reports task status.
- How it handles rate limits, retries, platform errors, and credential invalidation.
- Whether it supports incremental collection.
- Whether it supports batch collection.
- Whether it requires proxies, browsers, fingerprints, or other runtime resources.

Collection capability must be observable, retryable, and traceable. Collection must not be implemented as an unrecoverable black-box script.

### Plugins And Computation

The Computation System executes plugin-related computation capabilities.

A computation plugin should define:

- Which data layers it reads.
- Which data layers it writes.
- Which metrics it generates.
- Which reports it generates.
- Which alert inputs it generates.
- Whether it supports incremental computation.
- Whether it supports recomputation.
- Its input and output table relationships.
- How its computation tasks are scheduled.

Plugin computation capability must not be placed in the Web request path.

### Plugins And Open-source Boundary

Dox should consider an open-source base with private plugins.

The open-source base may provide:

- Plugin protocol.
- Plugin registration mechanism.
- Plugin configuration mechanism.
- Generic task scheduling model.
- Generic collection task model.
- Generic computation task model.
- Generic frontend plugin hosting capability.
- Example plugins.

Private plugins may contain:

- High-value platform integrations.
- Commercial data sources.
- Proprietary collection logic.
- Proprietary computation logic.
- Proprietary reports.
- Proprietary algorithms and business strategies.

Therefore, the core system must not hard-code private plugin capabilities into the base platform. The base platform must remain clear, extensible, and replaceable.

