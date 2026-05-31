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
 * @File    : lifecycle.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-05-31
 * @Modified: 2026-05-31
 */

package bootstrap

import (
	"context"
	"errors"
	"fmt"
)

// ShutdownFunc releases one runtime resource.
type ShutdownFunc func(context.Context) error

type shutdownHook struct {
	name string
	fn   ShutdownFunc
}

type runtimeLifecycle struct {
	hooks []shutdownHook
}

func (l *runtimeLifecycle) register(name string, fn ShutdownFunc) {
	if l == nil || fn == nil {
		return
	}
	l.hooks = append(l.hooks, shutdownHook{
		name: name,
		fn:   fn,
	})
}

func (l *runtimeLifecycle) shutdown(ctx context.Context) error {
	if l == nil || len(l.hooks) == 0 {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	hooks := l.hooks
	l.hooks = nil

	var err error
	for index := len(hooks) - 1; index >= 0; index-- {
		hook := hooks[index]
		if shutdownErr := hook.fn(ctx); shutdownErr != nil {
			err = errors.Join(err, fmt.Errorf("%s: %w", hook.name, shutdownErr))
		}
	}
	return err
}
