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
// EDA event, node, tag, field, and logging configuration types. It does not
// initialize zap, lumberjack, OpenTelemetry SDK providers, runtime loggers, or
// concrete sinks. Those are follow-up runtime integration milestones.
package logging
