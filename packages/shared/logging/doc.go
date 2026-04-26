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
 * @File    : doc.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

// Package logging defines the shared Dox logging model and configuration
// contract.
//
// This package is the first-stage observability vocabulary for Dox backend
// runtimes. It defines resource identity, correlation, observability event,
// node, tag, field, and logging configuration types. It also maps the Dox
// logging configuration to zap, zapcore, and lumberjack primitives for runtime
// integration. It does not initialize OpenTelemetry SDK providers or the Dox
// business logger API. Those are follow-up runtime integration milestones.
package logging
