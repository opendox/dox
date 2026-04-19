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

  @File    : SECURITY.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# Security Policy

Dox is expected to handle identities, platform credentials, task data, Amazon ecosystem data, third-party API responses, and commercially sensitive business data. Security reports are taken seriously.

## Reporting A Vulnerability

Please do not report security vulnerabilities through public GitHub issues.

Use one of these private channels instead:

- GitHub Security Advisory: `https://github.com/opendox/dox/security/advisories/new`
- Email: `frostleo.dev@gmail.com`

Please include as much of the following as you can safely provide:

- vulnerability type;
- affected files, APIs, plugins, or workflows;
- affected branch, tag, commit, or version;
- reproduction steps;
- expected impact;
- proof of concept, if appropriate;
- any configuration needed to reproduce the issue.

Do not include real production credentials or sensitive customer/business data in the report. Use redacted examples when possible.

## Sensitive Data Guidelines

Do not publish:

- passwords, tokens, cookies, private keys, OAuth secrets, refresh tokens, or API secrets;
- Amazon SP-API credentials, Ads API credentials, Lingxing credentials, Feishu credentials, crawler credentials, or notification channel secrets;
- authorization headers or raw platform credential payloads;
- logs or screenshots that reveal secrets;
- lower-layer raw business data that should not be public.

When referencing sensitive resources, use safe identifiers such as `credential_id`, internal IDs, hashes, or redacted placeholders.

## Expected Response

The maintainers will review security reports as soon as practical and may ask for additional details. Critical reports that can lead to credential exposure, unauthorized access, data leakage, or task manipulation will be prioritized.

Public disclosure should be coordinated with the maintainers so users have time to apply fixes.
