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

  @File    : .codex/agents/postgresql-data-layering.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 4. PostgreSQL-first Data Layering


The current Dox project scope must stay PostgreSQL-first.

PostgreSQL is the primary storage system for the current implementation scope. Do not introduce Spark, Flink, Hadoop, or a complex data lake stack as default implementation choices. These systems may remain part of the long-term architecture vocabulary, but they must not become default dependencies for the current implementation.

Dox should implement a lightweight warehouse-style layering model inside PostgreSQL. Clear data layers must separate collection, normalization, cleaning, computation, aggregation, and presentation concerns.

Recommended data layers:

- `raw`: raw collection data.
- `ods`: normalized operational source data.
- `dwd`: cleaned detail-level business facts.
- `dws`: subject-level summaries and aggregated metrics.
- `ads`: application-facing presentation data.

### raw Layer

The `raw` layer stores data as close to the original source as practical.

The `raw` layer may include:

- Original third-party API responses.
- Original crawler results.
- Raw JSON payloads.
- Raw HTML, snapshots, screenshots, or references to them.
- Collection batch information.
- Collection task IDs.
- Source platform identifiers.
- Original error information.
- Ingestion time and source event time.

The purpose of the `raw` layer is to preserve factual source material for traceability, replay, re-parsing, and debugging.

The `raw` layer may become large. The Web System may inspect `raw` data in controlled scenarios such as task debugging, sample inspection, parser validation, and operations review, but it must not perform large scans or aggregations over `raw` data in request handlers.

### ods Layer

The `ods` layer stores source data after basic normalization.

The `ods` layer may include:

- Normalized ASIN records.
- Normalized keyword records.
- Normalized search result records.
- Normalized listing fields.
- Normalized platform response fields.
- Normalized task results.
- Deduplicated source records.
- References back to the `raw` layer.

The purpose of the `ods` layer is to convert data from different platforms, formats, and collection methods into a relatively consistent source data structure.

The Collection System may write to the `ods` layer, but complex cleaning and business aggregation should not be performed inside collection processes.

### dwd Layer

The `dwd` layer stores cleaned detail-level business facts.

The `dwd` layer may include:

- Keyword search result facts.
- ASIN ranking facts.
- Listing snapshot facts.
- Advertising keyword facts.
- Market detail facts.
- Collection batch facts.
- Relationship facts between products, keywords, and ASINs.

The purpose of the `dwd` layer is to produce reliable business detail data that can be consumed by the Computation System.

The `dwd` layer should minimize source-system format differences and express Dox's own business model.

### dws Layer

The `dws` layer stores subject-level summary data and aggregated metrics.

The `dws` layer may include:

- Keyword-level summaries.
- ASIN-level summaries.
- Market-level summaries.
- Competitor-level summaries.
- Traffic structure summaries.
- Time trend summaries.
- Search result coverage summaries.
- Task execution quality summaries.

The purpose of the `dws` layer is to provide stable aggregated data for reports, trends, alerts, and application-facing presentation.

### ads Layer

The `ads` layer stores application-facing data for Web presentation.

The `ads` layer may include:

- Page card data.
- Dashboard data.
- Report results.
- Ranking results.
- Trend chart results.
- Alert input results.
- Aggregated views directly consumed by the frontend.

The purpose of the `ads` layer is to let Web APIs read and present data quickly and reliably.

The Web System should prefer `ads` or `dws` data for pages, reports, and analytics views.

### Read And Write Boundaries

The Collection System primarily writes `raw` and `ods` data.

The Computation System primarily reads `raw`, `ods`, or `dwd` data and writes `dwd`, `dws`, and `ads` data.

The Web System primarily reads `ads` and `dws` data, but it may read `raw`, `ods`, or `dwd` data in controlled scenarios such as task debugging, data validation, manual review, collection sample inspection, and operations management.

The Web System must not perform large scans, long-running aggregation, or bulk transformation over `raw`, `ods`, or `dwd` data in request handlers. Such work belongs in the Computation System.

### Table Design Principles

PostgreSQL table design should consider from the beginning:

- Data source.
- Collection batch.
- Task ID.
- Platform identifier.
- Time fields.
- Deduplication keys.
- Idempotent writes.
- Indexing strategy.
- Partitioning strategy.
- Data retention strategy.
- Traceability to upstream `raw` data.
- Computation relationship to downstream `dws` and `ads` data.

Do not put all data into a single mixed-purpose table for short-term convenience. Different data layers must have clear purposes and lifecycles.

### Computation Principles

The Computation System should generate upper-layer data through batch-oriented tasks.

Dox may use Go, SQL, scheduled jobs, dispatched tasks, and PostgreSQL capabilities for computation. It must not default to Spark, Flink, or Hadoop.

If a computation may later move to a dedicated computation system, preserve that migration path through clear task boundaries, input/output tables, and documentation instead of introducing a complex stack too early.

