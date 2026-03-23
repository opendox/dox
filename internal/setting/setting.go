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
 * setting.go
 *
 * - Author   : Frost Leo <frostleo.dev@gmail.com>
 * - Created  : 2026-03-23
 * - Modified : 2026-03-23
 */

package setting

// Setting Configuration sets required for each component
type Setting struct {
	Application Application `json:"application" yaml:"application" mapstructure:"application"`
	Cache       Cache       `json:"cache" yaml:"cache" mapstructure:"cache"`
	Database    Database    `json:"database" yaml:"database" mapstructure:"database"`
	Server      Server      `json:"server" yaml:"server" mapstructure:"server"`
}
