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
 * @File    : deployment.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-25
 * @Modified: 2026-04-25
 */

package setting

import "errors"

// Deployment describes where a Dox runtime process is deployed.
type Deployment struct {
	Region     string `json:"region" yaml:"region" mapstructure:"region" validate:"omitempty,dox_identifier"`
	Zone       string `json:"zone" yaml:"zone" mapstructure:"zone" validate:"omitempty,dox_identifier"`
	Cluster    string `json:"cluster" yaml:"cluster" mapstructure:"cluster" validate:"omitempty,dox_identifier"`
	InstanceID string `json:"instance_id" yaml:"instance_id" mapstructure:"instance_id" validate:"omitempty,dox_identifier"`
}

// Default is currently a no-op because deployment identity is environment-specific.
func (d *Deployment) Default() error {
	if d == nil {
		return errors.New("setting: deployment must not be nil")
	}
	return nil
}

// Validate verifies that non-empty deployment identity fields use stable syntax.
func (d Deployment) Validate() error {
	return Validate(d)
}
