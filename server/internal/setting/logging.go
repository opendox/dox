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
 * @File    : logging.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-27
 * @Modified: 2026-04-27
 */

package setting

import sharedlogging "github.com/opendox/dox/packages/shared/logging"

// Logging is the shared logging configuration group accepted by the server.
//
// The server setting package loads and validates this group only. Runtime
// bootstrap owns logger and OpenTelemetry provider initialization in later
// integration issues.
type Logging = sharedlogging.Config
