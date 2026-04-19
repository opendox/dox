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

  @File    : .codex/agents/collaboration-standard.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 2. Maintainer And Collaboration Standard


The maintainer, Leo, is an experienced architect and full-stack engineer. He understands product requirements, system architecture, backend design, frontend implementation, testing, verification, and long-term maintenance cost.

Agents must not treat the maintainer as a non-technical user. Superficially working code with poor structure is unacceptable. Placeholder logic, fake production readiness, vague abstractions, hidden risks, and low-quality shortcuts are unacceptable.

Agent work must be optimized for long-term maintainability, not for making the current response look complete. Implementations must be clear, reviewable, testable, evolvable, and consistent with Dox system boundaries.

Agents must be explicit about uncertainty, technical tradeoffs, incomplete parts, and verification results. Do not claim work is complete when it is not. Do not present experimental code as production-ready implementation.

Agents should avoid over-explaining basic concepts. When communicating with the maintainer, prefer engineering facts, design tradeoffs, risks, and executable plans.

Agents should:

- Read existing code before changing it.
- Keep changes scoped to the active issue.
- Avoid unrelated refactors.
- Ask questions only when blocked by a real architectural decision.
- Document important design decisions.
- Verify changes with tests, builds, screenshots, or clear manual checks.

