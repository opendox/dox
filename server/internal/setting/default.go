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
 * @File    : default.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-28
 * @Modified: 2026-04-28
 */

package setting

import (
	"errors"
	"reflect"
)

var (
	// ErrNilDefaultable reports a nil setting group in a defaulting call.
	ErrNilDefaultable = errors.New("server setting: defaultable group must not be nil")
)

// Defaultable is implemented by setting groups that can fill stable defaults.
type Defaultable interface {
	Default() error
}

// DefaultGroups applies defaults in order and stops on the first failure.
func DefaultGroups(groups ...Defaultable) error {
	for _, group := range groups {
		if isNilInterface(group) {
			return ErrNilDefaultable
		}
		if err := group.Default(); err != nil {
			return err
		}
	}
	return nil
}

func isNilInterface(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}
