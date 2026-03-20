<!--
Copyright (c) 2026 dox authors.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

.github/SECURITY.md

- Author   : Frost Leo <frostleo.dev@gmail.com>
- Created  : 2026-03-20
- Modified : 2026-03-20
-->

# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| latest  | Yes                |

## Reporting a Vulnerability

We take the security of dox seriously. If you discover a security vulnerability, please report it responsibly.

### How to Report

1. **Do not** open a public issue for security vulnerabilities
2. Email your findings to frostleo.dev@gmail.com
3. Include as much detail as possible:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### What to Expect

- **Acknowledgment**: We will acknowledge receipt of your report within 48 hours
- **Assessment**: We will assess the vulnerability and determine its severity
- **Updates**: We will keep you informed of our progress
- **Resolution**: We aim to resolve critical vulnerabilities within 7 days
- **Credit**: We will credit you in the security advisory (unless you prefer to remain anonymous)

### Scope

The following are in scope for security reports:

- Authentication and authorization issues
- Data exposure vulnerabilities
- SQL injection, XSS, CSRF
- Remote code execution
- Privilege escalation
- Cryptographic issues

### Out of Scope

- Denial of service attacks
- Social engineering
- Physical security
- Issues in dependencies (report to the respective project)

## Security Best Practices

When deploying dox:

1. Always use HTTPS in production
2. Keep your instance updated to the latest version
3. Use strong, unique passwords
4. Enable two-factor authentication when available
5. Regularly backup your data
6. Follow the principle of least privilege for database access
7. Review and rotate API keys periodically

## Disclosure Policy

We follow a coordinated disclosure process:

1. Reporter submits vulnerability
2. We validate and assess the issue
3. We develop and test a fix
4. We release the fix
5. We publish a security advisory
6. Reporter may publish their findings after the advisory

Thank you for helping keep dox and its users safe.
