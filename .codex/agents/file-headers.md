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

  @File    : .codex/agents/file-headers.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 10. File Headers And Source Metadata


Dox source files must keep consistent file headers.

File headers express copyright, license, file purpose, authorship, and creation/modification metadata. They are part of project governance and long-term maintenance, not decoration.

New source files must include file headers unless the file type or generated-file workflow makes that inappropriate.

Existing file header style must be preserved when editing existing files.

Do not batch-modify file headers across unrelated files for a small change.

Do not reformat file headers into a style inconsistent with the project.

Do not manually edit generated code only to add file headers unless the generator supports it.

### Go File Header

Go source files should use this style:

```go
/**
 * dox
 * Copyright (C) 2026  OpenDox
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * @File    : file_name.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : YYYY-MM-DD
 * @Modified: YYYY-MM-DD
 */

package example
```

### TypeScript File Header

TypeScript source files should use this style:

```ts
/**
 * dox
 * Copyright (C) 2026  OpenDox
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * @File    : file_name.ts
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : YYYY-MM-DD
 * @Modified: YYYY-MM-DD
 */
```

### Vue File Header

Vue single-file components should use this style:

```vue
<!--
  dox
  Copyright (C) 2026  OpenDox

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.

  @File    : ComponentName.vue
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : YYYY-MM-DD
  @Modified: YYYY-MM-DD
-->
```

