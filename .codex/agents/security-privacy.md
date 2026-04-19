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

  @File    : .codex/agents/security-privacy.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 13. Security And Privacy Rules

Dox must treat security, privacy, and credential protection as first-class concerns from the beginning.

Dox will handle user identities, platform credentials, third-party API credentials, collection data, business data, task data, and potentially commercially sensitive market data. Agents must not treat security as a later patch.

## IAM And Authentication

The Web System must provide clear authentication and authorization boundaries.

Authentication and authorization capabilities should consider:

- User login.
- Session management.
- Token management.
- MFA.
- OAuth and third-party login.
- Roles.
- Permission points.
- Menu permissions.
- API permissions.
- Plugin permissions.
- Task permissions.
- Audit logs.

Authentication logic must stay clear. Do not collapse login, authorization, menus, plugin enablement, and task access into one tangled model.

## Credential Management

Dox may store many kinds of external platform credentials, such as:

- Amazon SP-API credentials.
- Amazon Ads API credentials.
- Lingxing API credentials.
- Feishu API credentials.
- SIF or other platform credentials.
- Proxy, crawler, or browser runtime credentials.
- Notification channel credentials.

Credentials must be treated as sensitive data.

Credential rules:

- Do not write plaintext credentials to logs.
- Do not return plaintext credentials to the frontend.
- Do not put plaintext credentials into Redis keys.
- Do not put plaintext credentials into task names, event names, or error messages.
- Do not commit plaintext credentials to Git.
- Do not use real credentials in test data.
- Do not expose credentials in screenshots, PR descriptions, or issues.

When referencing credentials, use `credential_id`, internal IDs, hashes, or safe references.

## Configuration And Secrets

Sensitive configuration must support injection through environment variables or secure configuration sources.

Configuration templates may contain field structure, but must not contain real secrets.

Local configuration files, production configuration files, tokens, secrets, private keys, cookies, refresh tokens, and similar data must be protected by `.gitignore` or secure workflows.

When agents create configuration templates, use obvious placeholders such as:

- `<CHANGE_ME>`
- `<REDACTED>`
- `<YOUR_SECRET>`
- `<YOUR_CLIENT_ID>`

Do not generate fake secrets that look like real secrets.

## Log Redaction

Logs must be redacted.

Do not log:

- Passwords.
- Tokens.
- Refresh tokens.
- Access tokens.
- API secrets.
- Client secrets.
- Private keys.
- Cookies.
- Authorization headers.
- Raw platform credentials.
- Information that may compromise account security.

Error logs should preserve debugging context, but sensitive fields must be hidden or redacted.

Prefer shared redaction utilities over hand-written redaction in each module.

## Authorization And Data Access

Dox authorization is not only page authorization.

Authorization includes:

- API authorization.
- Plugin authorization.
- Platform credential authorization.
- Task viewing authorization.
- Task operation authorization.
- Queue management authorization.
- Lower-layer data access authorization.
- Alert subscription authorization.
- Notification channel configuration authorization.

The Web System may access lower-layer data, but lower-layer data access must be authorized.

High-risk operations such as queue adjustment, task cancellation, task re-run, task reordering, plugin enablement, and credential modification must have explicit authorization and audit records.

## Audit Logs

Key operations must write audit logs.

Audited operations may include:

- Login.
- Logout.
- MFA success or failure.
- User permission changes.
- Role changes.
- Plugin enablement or disablement.
- Platform credential creation, update, or deletion.
- Manual collection task trigger.
- Task cancellation.
- Task re-run.
- Queue adjustment.
- Computation task changes.
- Alert rule changes.
- Notification channel changes.
- Viewing or exporting sensitive data.

Audit logs should record:

- Operator.
- Operation type.
- Operation target.
- Operation time.
- Source IP or device information.
- Result.
- Failure reason.
- `correlation_id` or `request_id`.

Audit logs must not leak secrets.

## Plugin Security

Plugin capabilities must be controlled by authorization and configuration.

A plugin must not become visible to all users only because it is installed or exists in code.

Plugin enablement must go through Web System configuration and authorization.

When a plugin executes tasks, it must use an explicit `credential_id` and task context. It must not implicitly read global credentials.

Plugin errors must not leak platform credentials or sensitive responses.

## Collection And External Platform Security

The Collection System must respect platform authentication, rate limits, and failure handling boundaries.

Collection tasks must consider:

- Rate limits.
- Retry.
- Credential invalidation.
- Account risk.
- Proxy risk.
- Task cancellation.
- Duplicate data.
- Secure raw data storage.
- External platform error redaction.

Do not expose complete sensitive external platform responses directly to ordinary users.

## AI Security

AI must not access or output unauthorized sensitive data.

AI must not invent credentials, API fields, platform rules, business data, or security conclusions.

AI-generated suggestions must distinguish facts, inferences, and suggestions.

When AI participates in high-risk actions, the system must provide authorization, audit logs, human confirmation, or rollback mechanisms.

## Issue And PR Security

Issues, PRs, commit messages, screenshots, and log snippets must not contain sensitive information.

If the user pastes a secret, agents should warn the user to revoke or rotate that secret instead of continuing to spread it.
