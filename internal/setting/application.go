/*
 * Copyright © 2026 dox authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 *
 * application.go
 *
 * - Author   : Frost Leo <frostleo.dev@gmail.com>
 * - Created  : 2026-03-23
 * - Modified : 2026-03-23
 */

package setting

import (
	"errors"
	
	"github.com/opendox/dox/internal/version"
)

// Application basic information configuration
type Application struct {
	Name        string `json:"name" yaml:"name" mapstructure:"name"`
	Version     string `json:"version" yaml:"version" mapstructure:"version"`
	Description string `json:"description" yaml:"description" mapstructure:"description"`
}

func (a *Application) Default() error {
	if a.Name == "" {
		a.Name = "Dox"
	}
	if a.Version == "" {
		a.Version = version.GetInfo().Version()
	}
	if a.Description == "" {
		a.Description = "Enterprise-grade Amazon product analytics platform"
	}
	return nil
}

func (a *Application) Validate() error {
	if a.Name == "" {
		return errors.New("application: name is required")
	}
	if a.Version == "" {
		return errors.New("application: version is required")
	}
	return nil
}
