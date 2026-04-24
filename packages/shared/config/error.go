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
 * @File    : error.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-24
 * @Modified: 2026-04-24
 */

package config

import (
	"errors"
	"fmt"
)

// ErrorKind classifies loader failures without tying callers to string parsing.
type ErrorKind string

const (
	// ErrorKindContract means the caller violated the loader API contract.
	ErrorKindContract ErrorKind = "contract"
	// ErrorKindSource is reserved for future provider read failures.
	ErrorKindSource ErrorKind = "source"
	// ErrorKindDecode is reserved for future decode failures.
	ErrorKindDecode ErrorKind = "decode"
)

// Error describes a typed configuration loading error.
type Error struct {
	Kind   ErrorKind
	Field  string
	Reason string
	Err    error
}

// Error returns a stable human-readable message for the configuration error.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	message := fmt.Sprintf("config %s error", e.Kind)
	if e.Field != "" {
		message += " at " + e.Field
	}
	if e.Reason != "" {
		message += ": " + e.Reason
	}
	if e.Err != nil {
		message += ": " + e.Err.Error()
	}
	return message
}

// Unwrap returns the wrapped error when one is present.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Is reports whether the target error has the same configuration error kind.
func (e *Error) Is(target error) bool {
	var other *Error
	if !errors.As(target, &other) {
		return false
	}
	return e != nil && other != nil && e.Kind == other.Kind
}

// ContractError creates an error for invalid loader API usage.
func ContractError(field string, reason string) error {
	return &Error{Kind: ErrorKindContract, Field: field, Reason: reason}
}

// SourceError creates an error for future provider read failures.
func SourceError(field string, reason string, err error) error {
	return &Error{Kind: ErrorKindSource, Field: field, Reason: reason, Err: err}
}

// DecodeError creates an error for future decode failures.
func DecodeError(field string, reason string, err error) error {
	return &Error{Kind: ErrorKindDecode, Field: field, Reason: reason, Err: err}
}

// IsKind reports whether err contains a configuration error of the given kind.
func IsKind(err error, kind ErrorKind) bool {
	var cfgErr *Error
	if !errors.As(err, &cfgErr) {
		return false
	}
	return cfgErr.Kind == kind
}
