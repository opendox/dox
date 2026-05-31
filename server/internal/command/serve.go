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
 * @File    : serve.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-05-31
 * @Modified: 2026-05-31
 */

package command

import (
	"context"
	"errors"

	"github.com/opendox/dox/server/internal/bootstrap"
	"github.com/spf13/cobra"
)

// ServeWaitFunc waits until the booted server runtime should shut down.
type ServeWaitFunc func(context.Context, *bootstrap.Runtime) error

func newServeCommand(cfg Config) *cobra.Command {
	options := bootstrap.ConfigOptions{}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the Dox Web backend runtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			snapshot, err := bootstrap.LoadSetting(ctx, options)
			if err != nil {
				return err
			}

			runtime, err := bootstrap.Boot(ctx, snapshot, bootstrap.BootOptions{})
			if err != nil {
				return err
			}

			wait := cfg.ServeWait
			if wait == nil {
				wait = waitForContextCancellation
			}

			waitErr := wait(ctx, runtime)
			shutdownErr := runtime.Shutdown(context.Background())
			return errors.Join(waitErr, shutdownErr)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.ConfigDir, "config-dir", bootstrap.DefaultConfigDir, "Directory containing server configuration files")
	flags.StringVar(&options.Env, "env", bootstrap.DefaultConfigEnv, "Server runtime environment")
	flags.StringVar(&options.Format, "config-format", bootstrap.DefaultConfigFormat, "Server configuration file format")
	flags.StringVar(&options.EnvPrefix, "env-prefix", bootstrap.DefaultConfigEnvPrefix, "Environment variable prefix for server configuration overrides")
	flags.DurationVar(&options.Timeout, "config-timeout", 0, "Optional timeout for server configuration loading")

	return cmd
}

func waitForContextCancellation(ctx context.Context, runtime *bootstrap.Runtime) error {
	if ctx == nil {
		ctx = context.Background()
	}
	<-ctx.Done()
	return nil
}
