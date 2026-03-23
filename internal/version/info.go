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
 * info.go
 *
 * - Author   : Frost Leo <frostleo.dev@gmail.com>
 * - Created  : 2026-03-23
 * - Modified : 2026-03-23
 */

package version

import (
	"fmt"
	"strings"
)

// Info Version build information
type Info struct {
	name       string
	version    string
	gitCommit  string
	gitBranch  string
	gitTag     string
	buildTime  string
	buildUser  string
	goVersion  string
	cgoEnabled string
}

// Read-only collection
func (i Info) Name() string       { return i.name }
func (i Info) Version() string    { return i.version }
func (i Info) GitCommit() string  { return i.gitCommit }
func (i Info) GitBranch() string  { return i.gitBranch }
func (i Info) GitTag() string     { return i.gitTag }
func (i Info) BuildTime() string  { return i.buildTime }
func (i Info) BuildUser() string  { return i.buildUser }
func (i Info) GoVersion() string  { return i.goVersion }
func (i Info) CGOEnabled() string { return i.cgoEnabled }

// String
func (i Info) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Name: %s\n", i.name))
	b.WriteString(fmt.Sprintf("Version: %s\n", i.version))
	b.WriteString(fmt.Sprintf("GitCommit: %s\n", i.gitCommit))
	b.WriteString(fmt.Sprintf("GitBranch: %s\n", i.gitBranch))
	b.WriteString(fmt.Sprintf("GitTag: %s\n", i.gitTag))
	b.WriteString(fmt.Sprintf("BuildTime: %s\n", i.buildTime))
	b.WriteString(fmt.Sprintf("BuildUser: %s\n", i.buildUser))
	b.WriteString(fmt.Sprintf("GoVersion: %s\n", i.goVersion))
	b.WriteString(fmt.Sprintf("CGOEnabled: %s\n", i.cgoEnabled))

	return b.String()
}
