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
 * @File    : application.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-25
 * @Modified: 2026-04-25
 */

package setting

import "errors"

const (
	// DefaultApplicationName is the shared Dox application family name.
	DefaultApplicationName = "dox"
	// DefaultApplicationNamespace is the shared Dox service namespace.
	DefaultApplicationNamespace = "opendox"
)

// Runtime identifies one Dox runtime system.
type Runtime string

const (
	// RuntimeServer identifies the Web backend runtime.
	RuntimeServer Runtime = "server"
	// RuntimeScheduler identifies the scheduling runtime.
	RuntimeScheduler Runtime = "scheduler"
	// RuntimeCollector identifies the collection runtime.
	RuntimeCollector Runtime = "collector"
	// RuntimeCompute identifies the computation runtime.
	RuntimeCompute Runtime = "compute"
)

// IsValid reports whether r is one of the supported Dox runtime systems.
func (r Runtime) IsValid() bool {
	switch r {
	case RuntimeServer, RuntimeScheduler, RuntimeCollector, RuntimeCompute:
		return true
	default:
		return false
	}
}

// Env identifies the deployment environment for a Dox runtime.
type Env string

const (
	// EnvDev identifies a development environment.
	EnvDev Env = "dev"
	// EnvTest identifies a test environment.
	EnvTest Env = "test"
	// EnvStaging identifies a staging environment.
	EnvStaging Env = "staging"
	// EnvProd identifies a production environment.
	EnvProd Env = "prod"
)

// IsValid reports whether e is one of the supported Dox deployment environments.
func (e Env) IsValid() bool {
	switch e {
	case EnvDev, EnvTest, EnvStaging, EnvProd:
		return true
	default:
		return false
	}
}

// Application describes the shared identity of a Dox runtime process.
type Application struct {
	Name      string  `json:"name" yaml:"name" mapstructure:"name" validate:"required,dox_kebab"`
	Namespace string  `json:"namespace" yaml:"namespace" mapstructure:"namespace" validate:"required,dox_kebab"`
	Runtime   Runtime `json:"runtime" yaml:"runtime" mapstructure:"runtime" validate:"required,dox_runtime"`
	Service   string  `json:"service" yaml:"service" mapstructure:"service" validate:"required,dox_kebab"`
	Env       Env     `json:"env" yaml:"env" mapstructure:"env" validate:"required,dox_env"`
}

// Default fills stable shared application identity defaults.
func (a *Application) Default() error {
	if a == nil {
		return errors.New("setting: application must not be nil")
	}
	if a.Name == "" {
		a.Name = DefaultApplicationName
	}
	if a.Namespace == "" {
		a.Namespace = DefaultApplicationNamespace
	}
	if a.Env == "" {
		a.Env = EnvDev
	}
	if a.Service == "" && a.Name != "" && a.Runtime != "" {
		a.Service = a.Name + "-" + string(a.Runtime)
	}
	return nil
}

// Validate verifies that application identity fields use supported Dox values.
func (a Application) Validate() error {
	return Validate(a)
}
