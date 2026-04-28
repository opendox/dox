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
 * @File    : validate_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-28
 * @Modified: 2026-04-28
 */

package setting

import (
	"errors"
	"testing"
)

func TestValidateGroupsJoinsValidationFailures(t *testing.T) {
	firstErr := errors.New("first invalid")
	secondErr := errors.New("second invalid")
	first := validateTestGroup{err: firstErr}
	second := validateTestGroup{err: secondErr}

	err := ValidateGroups(first, second)
	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("expected joined validation errors, got %v", err)
	}
}

func TestValidateGroupsReportsNilGroups(t *testing.T) {
	if err := ValidateGroups(nil); !errors.Is(err, ErrNilValidatable) {
		t.Fatalf("expected nil validatable error, got %v", err)
	}

	var typedNil *validateTestGroup
	if err := ValidateGroups(typedNil); !errors.Is(err, ErrNilValidatable) {
		t.Fatalf("expected typed nil validatable error, got %v", err)
	}
}

type validateTestGroup struct {
	err error
}

func (g validateTestGroup) Validate() error {
	return g.err
}
