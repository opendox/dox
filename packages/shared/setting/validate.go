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
 * @Created : 2026-04-25
 * @Modified: 2026-04-25
 */

package setting

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var kebabPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)
var identifierPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9._-]*[a-z0-9])?$`)

var validateOnce sync.Once
var validateInstance *validator.Validate

// FieldError describes one setting validation failure without exposing the
// third-party validator error type.
type FieldError struct {
	Field string
	Rule  string
}

// ValidationError describes one or more setting validation failures.
type ValidationError struct {
	Fields []FieldError
}

// Error returns a compact validation failure message.
func (e *ValidationError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if len(e.Fields) == 0 {
		return "setting validation failed"
	}

	parts := make([]string, 0, len(e.Fields))
	for _, field := range e.Fields {
		if field.Rule == "" {
			parts = append(parts, field.Field)
			continue
		}
		parts = append(parts, field.Field+"="+field.Rule)
	}
	return "setting validation failed: " + strings.Join(parts, ", ")
}

// Validate verifies setting structs with Dox-owned validation rules.
func Validate(value any) error {
	engine := settingValidator()
	if err := engine.Struct(value); err != nil {
		return convertValidationError(err)
	}
	return nil
}

func settingValidator() *validator.Validate {
	validateOnce.Do(func() {
		validateInstance = validator.New(validator.WithRequiredStructEnabled())
		validateInstance.RegisterTagNameFunc(settingFieldName)
		mustRegisterValidation("dox_kebab", validateKebab)
		mustRegisterValidation("dox_identifier", validateIdentifier)
		mustRegisterValidation("dox_runtime", validateRuntime)
		mustRegisterValidation("dox_env", validateEnv)
	})
	return validateInstance
}

func settingFieldName(field reflect.StructField) string {
	name := field.Tag.Get("mapstructure")
	if name == "" {
		name = field.Tag.Get("json")
	}
	if name == "" || name == "-" {
		return field.Name
	}
	name, _, _ = strings.Cut(name, ",")
	if name == "" {
		return field.Name
	}
	return name
}

func mustRegisterValidation(tag string, fn validator.Func) {
	if err := validateInstance.RegisterValidation(tag, fn); err != nil {
		panic(fmt.Sprintf("setting: register %s validation: %v", tag, err))
	}
}

func validateKebab(level validator.FieldLevel) bool {
	return kebabPattern.MatchString(level.Field().String())
}

func validateIdentifier(level validator.FieldLevel) bool {
	return identifierPattern.MatchString(level.Field().String())
}

func validateRuntime(level validator.FieldLevel) bool {
	return Runtime(level.Field().String()).IsValid()
}

func validateEnv(level validator.FieldLevel) bool {
	return Env(level.Field().String()).IsValid()
}

func convertValidationError(err error) error {
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return &ValidationError{Fields: []FieldError{{Field: "value", Rule: "invalid"}}}
	}

	fields := make([]FieldError, 0, len(validationErrors))
	for _, field := range validationErrors {
		fields = append(fields, FieldError{
			Field: field.Namespace(),
			Rule:  field.Tag(),
		})
	}
	return &ValidationError{Fields: fields}
}
