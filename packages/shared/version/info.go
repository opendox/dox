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
 * @File    : info.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

const unknownValue = "unknown"

var (
	buildName       = "dox"
	buildVersion    = "0.1.0"
	buildTime       = unknownValue
	buildUser       = unknownValue
	buildCGOEnabled = unknownValue
	buildGitCommit  = unknownValue
	buildGitBranch  = unknownValue
	buildGitTag     = unknownValue
	buildGitDirty   = unknownValue
)

// Info describes the build and source metadata of a Dox backend binary.
type Info struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	BuildTime  string `json:"build_time"`
	BuildUser  string `json:"build_user"`
	GoVersion  string `json:"go_version"`
	GOOS       string `json:"goos"`
	GOARCH     string `json:"goarch"`
	CGOEnabled string `json:"cgo_enabled"`
	GitCommit  string `json:"git_commit"`
	GitBranch  string `json:"git_branch"`
	GitTag     string `json:"git_tag"`
	GitDirty   string `json:"git_dirty"`
}

// GetInfo returns normalized version metadata for the current binary.
func GetInfo() Info {
	info := Info{
		Name:       stableValue(buildName, "dox"),
		Version:    stableValue(buildVersion, "0.1.0"),
		BuildTime:  stableValue(buildTime, unknownValue),
		BuildUser:  stableValue(buildUser, unknownValue),
		GoVersion:  runtime.Version(),
		GOOS:       runtime.GOOS,
		GOARCH:     runtime.GOARCH,
		CGOEnabled: stableValue(buildCGOEnabled, unknownValue),
		GitCommit:  stableValue(buildGitCommit, unknownValue),
		GitBranch:  stableValue(buildGitBranch, unknownValue),
		GitTag:     stableValue(buildGitTag, unknownValue),
		GitDirty:   stableValue(buildGitDirty, unknownValue),
	}
	applyBuildInfoFallback(&info)
	return info
}

// String returns a stable multi-line representation for CLI output.
func (i Info) String() string {
	return fmt.Sprintf(`%s %s
  Build Time  : %s
  Build User  : %s
  Go Version  : %s
  GOOS        : %s
  GOARCH      : %s
  CGO Enabled : %s
  Git Commit  : %s
  Git Branch  : %s
  Git Tag     : %s
  Git Dirty   : %s
`, i.Name, i.Version,
		i.BuildTime,
		i.BuildUser,
		i.GoVersion,
		i.GOOS,
		i.GOARCH,
		i.CGOEnabled,
		i.GitCommit,
		i.GitBranch,
		i.GitTag,
		i.GitDirty,
	)
}

// Short returns a one-line representation for concise logs and diagnostics.
func (i Info) Short() string {
	commit := i.GitCommit
	if !hasValue(commit) {
		commit = unknownValue
	}
	if len(commit) > 12 {
		commit = commit[:12]
	}

	state := "clean"
	if i.IsDirty() {
		state = "dirty"
	}
	if !hasValue(i.GitDirty) {
		state = unknownValue
	}

	return fmt.Sprintf("%s %s (%s, %s)", i.Name, i.Version, commit, state)
}

// Fields returns version metadata as string fields for structured logging.
func (i Info) Fields() map[string]string {
	return map[string]string{
		"name":        i.Name,
		"version":     i.Version,
		"build_time":  i.BuildTime,
		"build_user":  i.BuildUser,
		"go_version":  i.GoVersion,
		"goos":        i.GOOS,
		"goarch":      i.GOARCH,
		"cgo_enabled": i.CGOEnabled,
		"git_commit":  i.GitCommit,
		"git_branch":  i.GitBranch,
		"git_tag":     i.GitTag,
		"git_dirty":   i.GitDirty,
	}
}

// IsDirty reports whether the build was produced from a modified working tree.
func (i Info) IsDirty() bool {
	switch strings.ToLower(strings.TrimSpace(i.GitDirty)) {
	case "1", "true", "yes", "y", "dirty", "modified":
		return true
	default:
		return false
	}
}

func applyBuildInfoFallback(info *Info) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	if !hasValue(info.Version) && hasValue(buildInfo.Main.Version) && buildInfo.Main.Version != "(devel)" {
		info.Version = buildInfo.Main.Version
	}

	settings := map[string]string{}
	for _, setting := range buildInfo.Settings {
		settings[setting.Key] = setting.Value
	}
	if !hasValue(info.GitCommit) {
		info.GitCommit = stableValue(settings["vcs.revision"], info.GitCommit)
	}
	if !hasValue(info.BuildTime) {
		info.BuildTime = stableValue(settings["vcs.time"], info.BuildTime)
	}
	if !hasValue(info.GitDirty) {
		info.GitDirty = stableValue(settings["vcs.modified"], info.GitDirty)
	}
}

func stableValue(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func hasValue(value string) bool {
	value = strings.TrimSpace(value)
	return value != "" && !strings.EqualFold(value, unknownValue)
}
