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

  @File    : .codex/agents/project-identity.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 1. Project Identity And Scope


Dox is a plugin-oriented data collection, scheduling, computation, and business analytics platform for the Amazon ecosystem.

Dox is not a generic admin dashboard, a simple SP-API wrapper, a crawler-only tool, or a passive data display panel. Its core purpose is to let users configure platform credentials, industry context, and business rules, then enable platform plugins that drive data collection, task scheduling, offline computation, business aggregation, result presentation, and alerting.

Dox must be organized around four core systems:

- Web System
- Scheduling System
- Collection System
- Computation System

The Web System is responsible for user interaction, authentication, authorization, plugin configuration, menu and route issuance, business configuration, and result presentation.

The Scheduling System is responsible for generating tasks from Web configuration and system policies, then scheduling the Collection System and Computation System to execute those tasks.

The Collection System is responsible for calling external APIs, running crawlers, integrating third-party platforms, collecting raw data, and producing normalized source data.

The Computation System is responsible for reading lower-layer data from PostgreSQL, then performing cleaning, transformation, aggregation, metric calculation, and generation of application-facing data for Web presentation.

The current project scope is PostgreSQL-first. PostgreSQL is the primary storage system. Dox should implement a lightweight warehouse-style layering model inside PostgreSQL, such as `raw`, `ods`, `dwd`, `dws`, and `ads`. Web APIs should read application-facing upper-layer data, not scan raw collection data or perform heavy aggregation in request handlers.

Dox must not introduce Spark, Flink, Hadoop, or a complex data lake stack as default implementation choices. These systems may remain part of the long-term architecture vocabulary, but they are not part of the current project scope.

Dox must not implement a complete supplier scoring system, strict multi-tenant isolation, AI-driven autonomous business decisions, or automatic advertising execution by default. These are important future directions. The current architecture may leave room for them through naming, interface boundaries, and documentation, but it must not expand the implementation scope around them.

The focus of Dox is to establish clear system boundaries, a plugin mechanism, scheduling flow, collection flow, computation flow, PostgreSQL data layering, GitHub workflow, and a sustainable engineering foundation.

