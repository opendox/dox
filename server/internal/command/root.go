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
 * @File    : root.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package command

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Config defines external command streams for the server CLI.
type Config struct {
	Out    io.Writer
	ErrOut io.Writer
}

// Execute runs the server CLI with explicit args and stream configuration.
func Execute(ctx context.Context, args []string, cfg Config) error {
	cmd := NewRootCommand(cfg)
	cmd.SetArgs(args)
	return cmd.ExecuteContext(ctx)
}

// NewRootCommand builds the root command for the Dox Web backend runtime.
func NewRootCommand(cfg Config) *cobra.Command {
	out := cfg.Out
	if out == nil {
		out = os.Stdout
	}
	errOut := cfg.ErrOut
	if errOut == nil {
		errOut = os.Stderr
	}

	cmd := &cobra.Command{
		Use:           "dox-server",
		Short:         "Dox Web backend server",
		Long:          "dox-server is the CLI entrypoint for the Dox Web backend runtime.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetOut(out)
	cmd.SetErr(errOut)
	cmd.AddCommand(newVersionCommand())
	return cmd
}
