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
 * @File    : setting.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-26
 * @Modified: 2026-04-26
 */

package setting

import "errors"

// DefaultOptions carries bootstrap-derived values that can seed server settings.
type DefaultOptions struct {
	Env string
}

// Setting is the root configuration aggregate for the Dox Web backend runtime.
type Setting struct {
	Identity Identity `json:"identity" yaml:"identity" mapstructure:"identity"`
	Logging  Logging  `json:"logging" yaml:"logging" mapstructure:"logging"`
}

// Default fills stable server defaults without bootstrap-derived values.
func (s *Setting) Default() error {
	return s.DefaultWithOptions(DefaultOptions{})
}

// DefaultWithOptions fills stable server defaults and optional bootstrap-derived values.
func (s *Setting) DefaultWithOptions(options DefaultOptions) error {
	if s == nil {
		return errors.New("server setting: setting must not be nil")
	}
	if err := s.Identity.DefaultWithOptions(IdentityDefaultOptions{
		Env: options.Env,
	}); err != nil {
		return err
	}
	if err := s.Logging.Default(); err != nil {
		return err
	}
	return nil
}

// Validate verifies the concrete server configuration aggregate.
func (s Setting) Validate() error {
	return errors.Join(
		s.Identity.Validate(),
		s.Logging.Validate(),
	)
}
