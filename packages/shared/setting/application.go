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
 * @Modified: 2026-04-26
 */

package setting

import "errors"

const (
	// DefaultOrganizationName is the default Dox organization identity.
	DefaultOrganizationName = "opendox"
	// DefaultApplicationName is the shared Dox product or application family name.
	DefaultApplicationName = "dox"
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

// Organization describes shared ownership and governance identity.
type Organization struct {
	Name       string `json:"name" yaml:"name" mapstructure:"name" validate:"required,dox_identifier"`
	Owner      string `json:"owner" yaml:"owner" mapstructure:"owner" validate:"omitempty,dox_identifier"`
	CostCenter string `json:"cost_center" yaml:"cost_center" mapstructure:"cost_center" validate:"omitempty,dox_identifier"`
	Project    string `json:"project" yaml:"project" mapstructure:"project" validate:"omitempty,dox_identifier"`
}

// Default fills stable shared organization defaults.
func (o *Organization) Default() error {
	if o == nil {
		return errors.New("setting: organization must not be nil")
	}
	if o.Name == "" {
		o.Name = DefaultOrganizationName
	}
	return nil
}

// Validate verifies that organization identity fields use stable syntax.
func (o Organization) Validate() error {
	return Validate(o)
}

// Application describes the product or application family identity.
type Application struct {
	Name string `json:"name" yaml:"name" mapstructure:"name" validate:"required,dox_kebab"`
}

// Default fills stable shared application defaults.
func (a *Application) Default() error {
	if a == nil {
		return errors.New("setting: application must not be nil")
	}
	if a.Name == "" {
		a.Name = DefaultApplicationName
	}
	return nil
}

// Validate verifies that application identity fields use stable syntax.
func (a Application) Validate() error {
	return Validate(a)
}

// System describes the Dox core system identity for a runtime.
type System struct {
	Runtime Runtime `json:"runtime" yaml:"runtime" mapstructure:"runtime" validate:"required,dox_runtime"`
}

// Default is currently a no-op because Dox runtime identity must be explicit.
func (s *System) Default() error {
	if s == nil {
		return errors.New("setting: system must not be nil")
	}
	return nil
}

// Validate verifies that system identity fields use supported Dox values.
func (s System) Validate() error {
	return Validate(s)
}

// Service describes one logical service identity.
type Service struct {
	Namespace  string `json:"namespace" yaml:"namespace" mapstructure:"namespace" validate:"required,dox_kebab"`
	Name       string `json:"name" yaml:"name" mapstructure:"name" validate:"required,dox_kebab"`
	InstanceID string `json:"instance_id" yaml:"instance_id" mapstructure:"instance_id" validate:"omitempty,dox_identifier"`
}

// Default fills service identity defaults from the application and system identity.
func (s *Service) Default(application Application, system System) error {
	if s == nil {
		return errors.New("setting: service must not be nil")
	}
	if s.Namespace == "" {
		s.Namespace = application.Name
	}
	if s.Name == "" && system.Runtime != "" {
		s.Name = string(system.Runtime)
	}
	return nil
}

// Validate verifies that service identity fields use stable syntax.
func (s Service) Validate() error {
	return Validate(s)
}
