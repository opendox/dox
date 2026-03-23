/*
 * Copyright © 2026 dox authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 *
 * version.go
 *
 * - Author   : Frost Leo <frostleo.dev@gmail.com>
 * - Created  : 2026-03-23
 * - Modified : 2026-03-23
 */

package version

import (
	"runtime"
)

var (
	_name       = "Dox"
	_version    = "0.1.0"
	_gitCommit  = "unknown"
	_gitBranch  = "unknown"
	_gitTag     = "unknown"
	_buildTime  = "unknown"
	_buildUser  = "unknown"
	_cgoEnabled = "0"
)

// info is the singleton, populated once at package init.
var info Info

// init Inject version information using -ldflags "-X ..."
func init() {
	info = Info{
		name:       _name,
		version:    _version,
		gitCommit:  _gitCommit,
		gitBranch:  _gitBranch,
		gitTag:     _gitTag,
		buildTime:  _buildTime,
		buildUser:  _buildUser,
		goVersion:  runtime.Version(),
		cgoEnabled: _cgoEnabled,
	}
}
