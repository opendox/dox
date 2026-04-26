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
 * @File    : identity.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package setting

import (
	"errors"
	"strings"

	sharedsetting "github.com/opendox/dox/packages/shared/setting"
)

const (
	// ServerRuntime is the Dox runtime identity owned by the Web backend server.
	ServerRuntime = sharedsetting.RuntimeServer
)

var (
	// ErrInvalidServerRuntime reports a supported Dox runtime that is not server.
	ErrInvalidServerRuntime = errors.New("server setting: identity.system.runtime must be server")
)

// IdentityDefaultOptions carries bootstrap-derived identity defaults.
type IdentityDefaultOptions struct {
	Env string
}

// Identity groups the shared identity fragments used by the server runtime.
type Identity struct {
	Organization sharedsetting.Organization `json:"organization" yaml:"organization" mapstructure:"organization"`
	Application  sharedsetting.Application  `json:"application" yaml:"application" mapstructure:"application"`
	System       sharedsetting.System       `json:"system" yaml:"system" mapstructure:"system"`
	Service      sharedsetting.Service      `json:"service" yaml:"service" mapstructure:"service"`
	Deployment   sharedsetting.Deployment   `json:"deployment" yaml:"deployment" mapstructure:"deployment"`
}

// Default fills stable server identity defaults without bootstrap-derived values.
func (i *Identity) Default() error {
	return i.DefaultWithOptions(IdentityDefaultOptions{})
}

// DefaultWithOptions fills stable server identity defaults and optional bootstrap-derived values.
func (i *Identity) DefaultWithOptions(options IdentityDefaultOptions) error {
	if i == nil {
		return errors.New("server setting: identity must not be nil")
	}
	if err := i.Organization.Default(); err != nil {
		return err
	}
	if err := i.Application.Default(); err != nil {
		return err
	}
	if err := i.System.Default(); err != nil {
		return err
	}
	if i.System.Runtime == "" {
		i.System.Runtime = ServerRuntime
	}
	if err := i.Service.Default(i.Application, i.System); err != nil {
		return err
	}
	if i.Deployment.Env == "" {
		i.Deployment.Env = sharedsetting.Env(strings.TrimSpace(options.Env))
	}
	if err := i.Deployment.Default(); err != nil {
		return err
	}
	return nil
}

// Validate verifies shared identity fragments plus server-owned runtime rules.
func (i Identity) Validate() error {
	return errors.Join(
		i.Organization.Validate(),
		i.Application.Validate(),
		i.System.Validate(),
		i.validateServerRuntime(),
		i.Service.Validate(),
		i.Deployment.Validate(),
	)
}

func (i Identity) validateServerRuntime() error {
	if i.System.Runtime.IsValid() && i.System.Runtime != ServerRuntime {
		return ErrInvalidServerRuntime
	}
	return nil
}
