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

  @File    : .codex/agents/business-domain.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 8. Business Domain Principles


The business core of Dox is to build a sustainably expanding data flywheel and business analytics capability around the Amazon ecosystem.

Dox does not merely display data. Dox should help users form actionable business understanding around industries, products, keywords, ASINs, competitors, traffic structures, advertising data, market trends, and later supply-chain performance.

### Core Business Objects

Core business objects may include:

- Industry.
- Market.
- Store.
- Product.
- ASIN.
- Keyword.
- Traffic keyword.
- Search result.
- Search placement.
- Listing.
- Ranking list.
- Competitor.
- Advertising campaign.
- Advertising keyword.
- Platform plugin.
- Collection task.
- Computation task.
- Report.
- Alert.
- Supplier.
- Supplier product.

These objects do not all need to be fully implemented in the current project scope, but agents must understand that they are part of the long-term Dox domain model.

### Keyword And ASIN Data Flywheel

One core Dox business flywheel is continuous expansion through keywords and ASINs.

A typical flow:

1. A user provides owned products, target ASINs, an industry, or a market direction.
2. The system obtains related advertising traffic keywords or keywords through platforms such as SIF.
3. The system searches Amazon retail frontend pages with those keywords and collects search results, placement data, and related ASINs.
4. The system analyzes these ASINs for listings, rankings, traffic keywords, advertising keywords, competitor relationships, and market performance.
5. Newly discovered keywords and ASINs enter later collection and computation rounds.
6. The system gradually forms an industry-level network of keywords, ASINs, traffic structures, competitor structures, and market opportunities.

This flywheel must not expand blindly without limits. The system must control expansion through strategy, such as:

- Maximum expansion depth.
- Keyword quality.
- ASIN deduplication.
- Industry boundary.
- Collection budget.
- Priority.
- Allowlist and blocklist rules.
- Frequency control.
- Incremental strategy.
- Stop conditions.

These strategies may be built into the system and may also be adjustable by users in the Web System.

### External Data Sources

Dox will integrate multiple kinds of external data sources.

Typical data sources include:

- Amazon SP-API.
- Amazon Ads API.
- Amazon retail frontend search results.
- Amazon Best Sellers rankings.
- Listing pages.
- Search placement data.
- SIF keyword data.
- SellerSprite-style market data.
- Lingxing data.
- Feishu data.
- Future third-party platforms or internal systems.

Different data sources should enter the system through plugins. Do not hard-code platform-specific business logic into the core system.

### Business Data And System Data

Dox must distinguish business data from system runtime data.

Business data includes:

- ASIN performance.
- Keyword performance.
- Search results.
- Market trends.
- Competitor relationships.
- Advertising performance.
- Product performance.
- Report metrics.
- Business alert inputs.

System runtime data includes:

- Task status.
- Worker status.
- Queue status.
- Plugin health status.
- Scheduling status.
- Error logs.
- Retry records.
- System alerts.

These data categories may be related, but they must not be collapsed into the same model. Business users care about business changes. Developers and operators care about system runtime state.

### Supplier And Product Operations Direction

Suppliers, supplier products, fulfillment performance, logistics speed, product quality, post-launch market performance, and supplier scoring are important future business directions for Dox.

Long term, Dox may help sellers:

- Discover product opportunities through market data.
- Obtain product data through supplier plugins or supplier collaboration capabilities.
- Evaluate the market competitiveness of supplier products.
- Track supplier fulfillment performance.
- Evaluate suppliers by combining post-launch sales, advertising, keyword, competitor, and market performance.
- Quickly eliminate weak products or suppliers.
- Quickly increase investment in products and strategies that perform well.

However, Dox must not implement a complete supplier scoring system and must not make a supplier collaboration platform part of the default current scope.

Agents may leave extension space for supplier-related directions through naming, model boundaries, and documentation, but must not expand the active task into a supplier scoring system without an explicit issue.

### AI Business Capability Direction

AI is a future capability amplifier for Dox, not a replacement for system facts.

AI may be used to:

- Summarize market structure.
- Explain keyword and ASIN networks.
- Generate competitor analysis.
- Recommend keyword expansion directions.
- Explain business alerts.
- Generate supplier evaluation summaries.
- Recommend advertising directions.
- Analyze advertising data feedback.
- Provide next-step action suggestions to users.

AI output must distinguish:

- Facts.
- Inferences.
- Suggestions.

AI must not invent data, metrics, API fields, platform rules, or business results. AI suggestions should be traceable to data sources where practical.

Automatic advertising execution, automatic budget adjustment, automatic supplier decisions, and similar high-risk actions are future directions. Dox must not implement them by default. Any automatic business action capability must have permission control, audit logs, thresholds, human confirmation, or rollback mechanisms.

### Current Scope Control

Dox should prioritize:

- Web, Scheduling, Collection, and Computation system boundaries.
- Plugin mechanism.
- PostgreSQL-first data layering.
- Basic collection flow.
- Basic computation flow.
- Basic business presentation.
- Basic alerting and notification.
- GitHub workflow and engineering standards.

Dox must not overbuild long-term directions such as:

- Complete supplier scoring.
- Complete strict multi-tenant isolation.
- AI-driven autonomous business decisions.
- Automatic advertising execution.
- Spark, Flink, or Hadoop computation platforms.
- Complex data lake architecture.

