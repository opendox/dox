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
 * @File    : root_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package command

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRootCommandMetadata(t *testing.T) {
	cmd := NewRootCommand(Config{})

	if cmd.Use != "dox-server" {
		t.Fatalf("expected root command use dox-server, got %q", cmd.Use)
	}
	if cmd.Short == "" {
		t.Fatal("expected root command short description")
	}
}

func TestVersionCommandPrintsSharedVersionInfo(t *testing.T) {
	var out bytes.Buffer

	err := Execute(context.Background(), []string{"version"}, Config{Out: &out})
	if err != nil {
		t.Fatalf("expected version command to succeed: %v", err)
	}

	output := out.String()
	for _, expected := range []string{
		"dox 0.1.0",
		"Go Version",
		"Git Commit",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected version output to contain %q, got:\n%s", expected, output)
		}
	}
}
