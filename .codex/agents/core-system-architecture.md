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

  @File    : .codex/agents/core-system-architecture.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 3. Core System Architecture


Dox must be organized around four core systems:

- Web System
- Scheduling System
- Collection System
- Computation System

These systems may live in the same monorepo, but they must be treated as separate runtime systems. They should be able to start independently, deploy independently, scale independently, and must not assume that they run on the same machine, in the same process, or in the same deployment unit.

The control relationship must stay clear:

- The Web System controls the Scheduling System.
- The Scheduling System controls the Collection System.
- The Scheduling System controls the Computation System.
- The Collection System does not control the Scheduling System.
- The Computation System does not control the Scheduling System.
- The Web System must not directly perform the responsibilities of the Collection System or the Computation System.

### Web System

The Web System is the management and analytics platform for Dox, as well as the main user interaction entry point.

The Web System is responsible for:

- Users, authentication, authorization, and IAM.
- Plugin configuration and platform credential management.
- Business configuration management.
- Menu, route, permission, and plugin entry issuance for the frontend.
- CRUD capabilities for core business objects.
- Management of collection tasks, scheduling tasks, computation tasks, and plugin status.
- Presentation of business results, analytics reports, task status, and system status.
- Receiving user operations and handing asynchronous work to the Scheduling System.
- Querying business data, analytics data, task data, and necessary lower-layer collection data.
- Providing controlled lower-layer data access for debugging, validation, manual review, and operations management.

The Web System may access uncleaned or lower-layer data, but this access must be controlled. Such access must have a clear use case, authorization boundary, pagination, filtering, indexing, and performance limits.

The Web System must not perform large-scale long-running computation in request handlers. It must not directly perform industry-level aggregation, full-table scans, long SQL statistics jobs, bulk transformations, or complex report generation in the request path.

When a user operation requires heavy computation, batch processing, aggregation, or long-running execution, the Web System should create a task or event and hand it to the Scheduling System.

### Scheduling System

The Scheduling System is the task orchestration center of Dox.

The Scheduling System is responsible for:

- Generating tasks from Web System configuration, plugin enablement state, and system policies.
- Managing the lifecycle of collection tasks and computation tasks.
- Dispatching collection tasks to the Collection System.
- Dispatching computation tasks to the Computation System.
- Managing task status, priority, retry, failure, timeout, and scheduling policy.
- Listening to collection results, computation results, and system events.
- Returning task status and result summaries to the Web System.
- Coordinating collection processes and computation processes that may run on different machines.

The Scheduling System may contain Bridge capabilities. The role of the Bridge is to connect Web, collection, and computation through explicit task and event protocols instead of direct coupling.

The Scheduling System owns task decisions. The Collection System and Computation System execute tasks; they do not own global scheduling policy.

### Collection System

The Collection System is the data collection execution layer of Dox.

The Collection System is responsible for:

- Executing collection tasks dispatched by the Scheduling System.
- Calling external APIs such as Amazon SP-API, Amazon Ads API, Lingxing API, Feishu API, and similar platforms.
- Collecting Amazon retail frontend data such as search results, Best Sellers rankings, listings, and search placement data.
- Integrating external data platforms such as SIF keyword data and SellerSprite-style market data.
- Performing necessary basic cleaning, normalization, and task status reporting.
- Writing raw data or normalized source data.
- Sending collection result events back to the Scheduling System or message queue.

The Collection System must not decide the next global collection strategy. It may return collection results and data that can be used for later expansion, but whether to continue expanding keywords, ASINs, or task chains must be decided by the Scheduling System according to policy.

### Computation System

The Computation System is the execution layer for long-running computation, offline processing, aggregation, and report generation.

The Computation System is responsible for:

- Executing computation tasks dispatched by the Scheduling System.
- Reading lower-layer or detail-level data from PostgreSQL.
- Cleaning, transforming, aggregating, and calculating metrics from collected data.
- Generating upper-layer data for Web presentation.
- Generating reports, summary tables, trend data, business metrics, and alert inputs.
- Writing computation results back to PostgreSQL into layers suitable for business querying and analytics presentation.
- Sending computation result events back to the Scheduling System or Notification System.

The core boundary of the Computation System is to absorb long-running computation that does not belong in the Web request path. The Web System may query data and perform CRUD operations, but computations that require long execution time, large scans, large-scale aggregation, or bulk transformation should be handled by the Computation System.

