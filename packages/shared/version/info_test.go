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
 * @File    : info_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetInfoReturnsStableDefaults(t *testing.T) {
	info := GetInfo()

	if info.Name != "dox" {
		t.Fatalf("expected default name dox, got %q", info.Name)
	}
	if info.Version != "0.1.0" {
		t.Fatalf("expected default version 0.1.0, got %q", info.Version)
	}
	if info.GoVersion != runtime.Version() {
		t.Fatalf("expected Go version %q, got %q", runtime.Version(), info.GoVersion)
	}
	if info.GOOS != runtime.GOOS {
		t.Fatalf("expected GOOS %q, got %q", runtime.GOOS, info.GOOS)
	}
	if info.GOARCH != runtime.GOARCH {
		t.Fatalf("expected GOARCH %q, got %q", runtime.GOARCH, info.GOARCH)
	}
	if info.Fields()["name"] != "dox" {
		t.Fatal("expected fields to include name")
	}
}

func TestStringIncludesStructuredBuildMetadata(t *testing.T) {
	info := Info{
		Name:       "dox-server",
		Version:    "1.2.3",
		BuildTime:  "2026-04-24T00:00:00Z",
		BuildUser:  "builder",
		GoVersion:  "go1.25.6",
		GOOS:       "linux",
		GOARCH:     "amd64",
		CGOEnabled: "false",
		GitCommit:  "abcdef1234567890",
		GitBranch:  "main",
		GitTag:     "v1.2.3",
		GitDirty:   "false",
	}

	output := info.String()
	for _, expected := range []string{
		"dox-server 1.2.3",
		"Build Time  : 2026-04-24T00:00:00Z",
		"GOOS        : linux",
		"Git Commit  : abcdef1234567890",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestShortOutputTruncatesCommitAndReportsState(t *testing.T) {
	info := Info{
		Name:      "dox-scheduler",
		Version:   "0.1.0",
		GitCommit: "abcdef1234567890",
		GitDirty:  "true",
	}

	got := info.Short()
	want := "dox-scheduler 0.1.0 (abcdef123456, dirty)"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestIsDirtyNormalizesCommonValues(t *testing.T) {
	dirtyValues := []string{"1", "true", "yes", "dirty", "modified"}
	for _, value := range dirtyValues {
		if !((Info{GitDirty: value}).IsDirty()) {
			t.Fatalf("expected %q to be dirty", value)
		}
	}

	cleanValues := []string{"", "unknown", "UNKNOWN", "0", "false", "clean"}
	for _, value := range cleanValues {
		if (Info{GitDirty: value}).IsDirty() {
			t.Fatalf("expected %q to be clean", value)
		}
	}
}
