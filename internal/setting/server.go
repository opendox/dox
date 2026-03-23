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
 * server.go
 *
 * - Author   : Frost Leo <frostleo.dev@gmail.com>
 * - Created  : 2026-03-23
 * - Modified : 2026-03-23
 */

package setting

import "time"

// Server is the top-level HTTP server configuration for Fiber.
type Server struct {
	Listen     Listen     `json:"listen" yaml:"listen" mapstructure:"listen"`
	Routing    Routing    `json:"routing" yaml:"routing" mapstructure:"routing"`
	Body       Body       `json:"body" yaml:"body" mapstructure:"body"`
	Timeout    Timeout    `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Buffer     Buffer     `json:"buffer" yaml:"buffer" mapstructure:"buffer"`
	Connection Connection `json:"connection" yaml:"connection" mapstructure:"connection"`
	Header     Header     `json:"header" yaml:"header" mapstructure:"header"`
	Proxy      Proxy      `json:"proxy" yaml:"proxy" mapstructure:"proxy"`
	TLS        TLS        `json:"tls" yaml:"tls" mapstructure:"tls"`
	Startup    Startup    `json:"startup" yaml:"startup" mapstructure:"startup"`
	View       View       `json:"view" yaml:"view" mapstructure:"view"`
}

// Listen defines basic listening parameters.
type Listen struct {
	Addr               string   `json:"addr" yaml:"addr" mapstructure:"addr"`
	ListenerNetwork    string   `json:"listener_network" yaml:"listener_network" mapstructure:"listener_network"`
	UnixSocketFileMode uint32   `json:"unix_socket_file_mode" yaml:"unix_socket_file_mode" mapstructure:"unix_socket_file_mode"`
	ServerHeader       string   `json:"server_header" yaml:"server_header" mapstructure:"server_header"`
	Concurrency        int      `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
	RequestMethods     []string `json:"request_methods" yaml:"request_methods" mapstructure:"request_methods"`
}

// Routing controls route matching behaviour.
type Routing struct {
	CaseSensitive           bool `json:"case_sensitive" yaml:"case_sensitive" mapstructure:"case_sensitive"`
	StrictRouting           bool `json:"strict_routing" yaml:"strict_routing" mapstructure:"strict_routing"`
	UnescapePath            bool `json:"unescape_path" yaml:"unescape_path" mapstructure:"unescape_path"`
	DisableHeadAutoRegister bool `json:"disable_head_auto_register" yaml:"disable_head_auto_register" mapstructure:"disable_head_auto_register"`
}

// Body governs request body handling.
type Body struct {
	BodyLimit                    int  `json:"body_limit" yaml:"body_limit" mapstructure:"body_limit"`
	MaxRanges                    int  `json:"max_ranges" yaml:"max_ranges" mapstructure:"max_ranges"`
	StreamRequestBody            bool `json:"stream_request_body" yaml:"stream_request_body" mapstructure:"stream_request_body"`
	DisablePreParseMultipartForm bool `json:"disable_pre_parse_multipart_form" yaml:"disable_pre_parse_multipart_form" mapstructure:"disable_pre_parse_multipart_form"`
	EnableSplittingOnParsers     bool `json:"enable_splitting_on_parsers" yaml:"enable_splitting_on_parsers" mapstructure:"enable_splitting_on_parsers"`
}

// Timeout defines I/O and lifecycle timeouts.
type Timeout struct {
	ReadTimeout     time.Duration `json:"read_timeout" yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout" yaml:"write_timeout" mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout" yaml:"idle_timeout" mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout" yaml:"shutdown_timeout" mapstructure:"shutdown_timeout"`
}

// Buffer sets per-connection buffer sizes.
type Buffer struct {
	ReadBufferSize  int `json:"read_buffer_size" yaml:"read_buffer_size" mapstructure:"read_buffer_size"`
	WriteBufferSize int `json:"write_buffer_size" yaml:"write_buffer_size" mapstructure:"write_buffer_size"`
}

// Connection holds connection-level behavioural flags.
type Connection struct {
	DisableKeepalive  bool `json:"disable_keepalive" yaml:"disable_keepalive" mapstructure:"disable_keepalive"`
	GETOnly           bool `json:"get_only" yaml:"get_only" mapstructure:"get_only"`
	ReduceMemoryUsage bool `json:"reduce_memory_usage" yaml:"reduce_memory_usage" mapstructure:"reduce_memory_usage"`
	Immutable         bool `json:"immutable" yaml:"immutable" mapstructure:"immutable"`
}

// Header configures response headers and normalisation.
type Header struct {
	DisableDefaultContentType bool              `json:"disable_default_content_type" yaml:"disable_default_content_type" mapstructure:"disable_default_content_type"`
	DisableDefaultDate        bool              `json:"disable_default_date" yaml:"disable_default_date" mapstructure:"disable_default_date"`
	DisableHeaderNormalizing  bool              `json:"disable_header_normalizing" yaml:"disable_header_normalizing" mapstructure:"disable_header_normalizing"`
	CompressedFileSuffixes    map[string]string `json:"compressed_file_suffixes" yaml:"compressed_file_suffixes" mapstructure:"compressed_file_suffixes"`
}

// Proxy configures trusted proxy and client IP extraction.
type Proxy struct {
	TrustProxy         bool     `json:"trust_proxy" yaml:"trust_proxy" mapstructure:"trust_proxy"`
	TrustedProxies     []string `json:"trusted_proxies" yaml:"trusted_proxies" mapstructure:"trusted_proxies"`
	ProxyHeader        string   `json:"proxy_header" yaml:"proxy_header" mapstructure:"proxy_header"`
	EnableIPValidation bool     `json:"enable_ip_validation" yaml:"enable_ip_validation" mapstructure:"enable_ip_validation"`
}

// TLS provides file-based TLS/mTLS configuration.
type TLS struct {
	CertFile       string `json:"cert_file" yaml:"cert_file" mapstructure:"cert_file"`
	CertKeyFile    string `json:"cert_key_file" yaml:"cert_key_file" mapstructure:"cert_key_file"`
	CertClientFile string `json:"cert_client_file" yaml:"cert_client_file" mapstructure:"cert_client_file"`
	TLSMinVersion  uint16 `json:"tls_min_version" yaml:"tls_min_version" mapstructure:"tls_min_version"`
}

// Startup controls process-level startup behaviour.
type Startup struct {
	EnablePrefork         bool `json:"enable_prefork" yaml:"enable_prefork" mapstructure:"enable_prefork"`
	DisableStartupMessage bool `json:"disable_startup_message" yaml:"disable_startup_message" mapstructure:"disable_startup_message"`
	EnablePrintRoutes     bool `json:"enable_print_routes" yaml:"enable_print_routes" mapstructure:"enable_print_routes"`
}

// View configures template rendering.
type View struct {
	ViewsLayout         string `json:"views_layout" yaml:"views_layout" mapstructure:"views_layout"`
	PassLocalsToViews   bool   `json:"pass_locals_to_views" yaml:"pass_locals_to_views" mapstructure:"pass_locals_to_views"`
	PassLocalsToContext bool   `json:"pass_locals_to_context" yaml:"pass_locals_to_context" mapstructure:"pass_locals_to_context"`
}
