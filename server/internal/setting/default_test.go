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
 * @File    : default_test.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-28
 * @Modified: 2026-04-28
 */

package setting

import (
	"errors"
	"reflect"
	"testing"
)

func TestDefaultGroupsAppliesDefaultsInOrder(t *testing.T) {
	var calls []string
	first := defaultTestGroup{name: "first", calls: &calls}
	second := defaultTestGroup{name: "second", calls: &calls}

	if err := DefaultGroups(&first, &second); err != nil {
		t.Fatalf("default groups: %v", err)
	}
	if !reflect.DeepEqual(calls, []string{"first", "second"}) {
		t.Fatalf("unexpected default order: %+v", calls)
	}
}

func TestDefaultGroupsStopsOnFirstFailure(t *testing.T) {
	expected := errors.New("default failed")
	var calls []string
	first := defaultTestGroup{name: "first", calls: &calls}
	second := defaultTestGroup{name: "second", calls: &calls, err: expected}
	third := defaultTestGroup{name: "third", calls: &calls}

	if err := DefaultGroups(&first, &second, &third); !errors.Is(err, expected) {
		t.Fatalf("expected default failure, got %v", err)
	}
	if !reflect.DeepEqual(calls, []string{"first", "second"}) {
		t.Fatalf("expected defaulting to stop on second group, got %+v", calls)
	}
}

func TestDefaultGroupsRejectsNilGroup(t *testing.T) {
	if err := DefaultGroups(nil); !errors.Is(err, ErrNilDefaultable) {
		t.Fatalf("expected nil defaultable error, got %v", err)
	}

	var typedNil *defaultTestGroup
	if err := DefaultGroups(typedNil); !errors.Is(err, ErrNilDefaultable) {
		t.Fatalf("expected typed nil defaultable error, got %v", err)
	}
}

type defaultTestGroup struct {
	name  string
	calls *[]string
	err   error
}

func (g *defaultTestGroup) Default() error {
	if g.calls != nil {
		*g.calls = append(*g.calls, g.name)
	}
	return g.err
}
