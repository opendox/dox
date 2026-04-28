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
 * @File    : validate.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-28
 * @Modified: 2026-04-28
 */

package setting

import "errors"

var (
	// ErrNilValidatable reports a nil setting group in a validation call.
	ErrNilValidatable = errors.New("server setting: validatable group must not be nil")
)

// Validatable is implemented by setting groups that can validate their values.
type Validatable interface {
	Validate() error
}

// ValidateGroups validates groups and joins all reported failures.
func ValidateGroups(groups ...Validatable) error {
	errs := make([]error, 0, len(groups))
	for _, group := range groups {
		if isNilInterface(group) {
			errs = append(errs, ErrNilValidatable)
			continue
		}
		errs = append(errs, group.Validate())
	}
	return errors.Join(errs...)
}
