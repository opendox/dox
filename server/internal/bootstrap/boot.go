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
 * @File    : boot.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-05-31
 * @Modified: 2026-05-31
 */

package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"syscall"

	sharedconfig "github.com/opendox/dox/packages/shared/config"
	sharedlogging "github.com/opendox/dox/packages/shared/logging"
)

// BootOptions carries runtime boot policy toggles.
type BootOptions struct {
	// InstallOpenTelemetryGlobals is reserved for a later runtime integration.
	InstallOpenTelemetryGlobals bool
}

// Runtime owns resources constructed for the server process.
type Runtime struct {
	Snapshot *SettingSnapshot
	Resource sharedlogging.Resource
	Logger   sharedlogging.Logger

	zapCore   *sharedlogging.ZapCoreBase
	otel      *sharedlogging.OpenTelemetrySDKBase
	lifecycle runtimeLifecycle
}

// Boot constructs server runtime resources from a validated setting snapshot.
func Boot(ctx context.Context, snapshot *SettingSnapshot, options BootOptions) (*Runtime, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if snapshot == nil {
		return nil, sharedconfig.ContractError("snapshot", "setting snapshot must not be nil")
	}
	if err := options.validate(); err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	resource := loggingResourceFromSnapshot(snapshot)
	loggingConfig := renderLoggingConfig(snapshot.Setting.Logging, resource)
	if err := prepareLoggingPaths(loggingConfig); err != nil {
		return nil, fmt.Errorf("server boot: prepare logging paths: %w", err)
	}

	runtime := &Runtime{
		Snapshot: snapshot,
		Resource: resource,
	}

	zapCore, err := sharedlogging.NewZapCoreBase(loggingConfig)
	if err != nil {
		return nil, fmt.Errorf("server boot: initialize zap core base: %w", err)
	}
	runtime.zapCore = zapCore
	runtime.lifecycle.register("logging.zap_core", func(context.Context) error {
		zapCore.Close()
		return nil
	})

	logger, err := sharedlogging.NewLogger(zapCore, sharedlogging.ResourceAttr(resource))
	if err != nil {
		return nil, runtime.rollback(ctx, fmt.Errorf("server boot: initialize logger: %w", err))
	}
	runtime.Logger = logger
	runtime.lifecycle.register("logging.logger", func(context.Context) error {
		return normalizeLoggerSyncError(logger.Sync())
	})

	otelBase, err := sharedlogging.NewOpenTelemetrySDKBase(loggingConfig, resource)
	if err != nil {
		return nil, runtime.rollback(ctx, fmt.Errorf("server boot: initialize OpenTelemetry SDK base: %w", err))
	}
	runtime.otel = otelBase
	runtime.lifecycle.register("logging.otel", func(ctx context.Context) error {
		return otelBase.Shutdown(ctx)
	})

	return runtime, nil
}

// Shutdown releases runtime resources in reverse construction order.
func (r *Runtime) Shutdown(ctx context.Context) error {
	if r == nil {
		return nil
	}
	return r.lifecycle.shutdown(ctx)
}

func (o BootOptions) validate() error {
	if o.InstallOpenTelemetryGlobals {
		return sharedconfig.ContractError(
			"options.install_open_telemetry_globals",
			"installing OpenTelemetry globals is not supported by server boot",
		)
	}
	return nil
}

func (r *Runtime) rollback(ctx context.Context, err error) error {
	if rollbackErr := r.Shutdown(ctx); rollbackErr != nil {
		return errors.Join(err, fmt.Errorf("server boot rollback failed: %w", rollbackErr))
	}
	return err
}

func normalizeLoggerSyncError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, syscall.EINVAL) {
		return nil
	}
	return err
}
