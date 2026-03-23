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
 * database.go
 *
 * - Author   : Frost Leo <frostleo.dev@gmail.com>
 * - Created  : 2026-03-23
 * - Modified : 2026-03-23
 */

package setting

import "time"

// Database holds all named database source configurations.
type Database struct {
	Sources map[string]DBSource `json:"sources" yaml:"sources" mapstructure:"sources"`
}

// DBSource represents a single database connection configuration.
type DBSource struct {
	Driver         string        `json:"driver" yaml:"driver" mapstructure:"driver"`    // "postgres" | "mysql"
	Net            string        `json:"net" yaml:"net" mapstructure:"net"`             // "tcp" | "tcp6" | "unix", default: "tcp"
	Host           string        `json:"host" yaml:"host" mapstructure:"host"`          // hostname, IP, or unix socket path
	Port           int           `json:"port" yaml:"port" mapstructure:"port"`          // default: 5432 (pg) / 3306 (mysql)
	DBName         string        `json:"db_name" yaml:"db_name" mapstructure:"db_name"` // database name
	User           string        `json:"user" yaml:"user" mapstructure:"user"`
	Password       string        `json:"password" yaml:"password" mapstructure:"password"`
	ConnectTimeout time.Duration `json:"connect_timeout" yaml:"connect_timeout" mapstructure:"connect_timeout"`
	ReadTimeout    time.Duration `json:"read_timeout" yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout" yaml:"write_timeout" mapstructure:"write_timeout"`
	TLS            DBTLS         `json:"tls" yaml:"tls" mapstructure:"tls"`
	Pool           DBPool        `json:"pool" yaml:"pool" mapstructure:"pool"`

	Postgres *PostgresOptions `json:"postgres,omitempty" yaml:"postgres,omitempty" mapstructure:"postgres"` // only when driver=postgres
	MySQL    *MySQLOptions    `json:"mysql,omitempty" yaml:"mysql,omitempty" mapstructure:"mysql"`          // only when driver=mysql
}

// DBTLS provides TLS/SSL configuration for database connections.
type DBTLS struct {
	Mode               string `json:"mode" yaml:"mode" mapstructure:"mode"`                                                 // pg: disable/allow/prefer/require/verify-ca/verify-full; mysql: false/true/skip-verify/preferred
	CACertFile         string `json:"ca_cert_file" yaml:"ca_cert_file" mapstructure:"ca_cert_file"`                         // CA certificate path (PEM)
	CertFile           string `json:"cert_file" yaml:"cert_file" mapstructure:"cert_file"`                                  // client certificate path (PEM) for mTLS
	KeyFile            string `json:"key_file" yaml:"key_file" mapstructure:"key_file"`                                     // client private key path (PEM) for mTLS
	ServerName         string `json:"server_name" yaml:"server_name" mapstructure:"server_name"`                            // override TLS server name verification
	InsecureSkipVerify bool   `json:"insecure_skip_verify" yaml:"insecure_skip_verify" mapstructure:"insecure_skip_verify"` // skip server cert verification
	MinVersion         string `json:"min_version" yaml:"min_version" mapstructure:"min_version"`                            // "1.0" | "1.1" | "1.2" | "1.3", default: "1.2"
}

// DBPool configures the database/sql connection pool.
type DBPool struct {
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `json:"max_open_conns" yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time" mapstructure:"conn_max_idle_time"`
}

// PostgresOptions holds PostgreSQL-specific configuration via pgx.
type PostgresOptions struct {
	RuntimeParams            map[string]string  `json:"runtime_params" yaml:"runtime_params" mapstructure:"runtime_params"`                                     // session params: search_path, application_name, timezone, etc.
	SSLNegotiation           string             `json:"ssl_negotiation" yaml:"ssl_negotiation" mapstructure:"ssl_negotiation"`                                  // "postgres" (default) | "direct"
	Fallbacks                []PostgresFallback `json:"fallbacks,omitempty" yaml:"fallbacks,omitempty" mapstructure:"fallbacks"`                                // HA failover targets
	DefaultQueryExecMode     string             `json:"default_query_exec_mode" yaml:"default_query_exec_mode" mapstructure:"default_query_exec_mode"`          // cache_statement/cache_describe/describe_exec/exec/simple_protocol
	StatementCacheCapacity   int                `json:"statement_cache_capacity" yaml:"statement_cache_capacity" mapstructure:"statement_cache_capacity"`       // default: 512
	DescriptionCacheCapacity int                `json:"description_cache_capacity" yaml:"description_cache_capacity" mapstructure:"description_cache_capacity"` // default: 512
	MinProtocolVersion       string             `json:"min_protocol_version" yaml:"min_protocol_version" mapstructure:"min_protocol_version"`                   // "3.0" | "3.2" | "latest"
	MaxProtocolVersion       string             `json:"max_protocol_version" yaml:"max_protocol_version" mapstructure:"max_protocol_version"`                   // "3.0" | "3.2" | "latest"
	ChannelBinding           string             `json:"channel_binding" yaml:"channel_binding" mapstructure:"channel_binding"`                                  // "disable" | "prefer" | "require"
	KerberosSrvName          string             `json:"kerberos_srv_name" yaml:"kerberos_srv_name" mapstructure:"kerberos_srv_name"`
	KerberosSpn              string             `json:"kerberos_spn" yaml:"kerberos_spn" mapstructure:"kerberos_spn"`
}

// PostgresFallback defines an alternative host/port for connection failover.
type PostgresFallback struct {
	Host    string `json:"host" yaml:"host" mapstructure:"host"`
	Port    int    `json:"port" yaml:"port" mapstructure:"port"`
	TLSMode string `json:"tls_mode" yaml:"tls_mode" mapstructure:"tls_mode"` // overrides parent DBTLS.Mode for this fallback
}

// MySQLOptions holds MySQL-specific configuration via go-sql-driver/mysql.
type MySQLOptions struct {
	Charset                  string            `json:"charset" yaml:"charset" mapstructure:"charset"`                                  // e.g. "utf8mb4"
	Collation                string            `json:"collation" yaml:"collation" mapstructure:"collation"`                            // e.g. "utf8mb4_unicode_ci"
	Loc                      string            `json:"loc" yaml:"loc" mapstructure:"loc"`                                              // time.Location name, default: "UTC"
	ParseTime                *bool             `json:"parse_time" yaml:"parse_time" mapstructure:"parse_time"`                         // parse DATE/DATETIME to time.Time, default: true
	TimeTruncate             time.Duration     `json:"time_truncate" yaml:"time_truncate" mapstructure:"time_truncate"`                // truncate time.Time precision
	MaxAllowedPacket         int               `json:"max_allowed_packet" yaml:"max_allowed_packet" mapstructure:"max_allowed_packet"` // default: 64MiB, 0=auto-fetch
	InterpolateParams        bool              `json:"interpolate_params" yaml:"interpolate_params" mapstructure:"interpolate_params"` // client-side placeholder interpolation
	MultiStatements          bool              `json:"multi_statements" yaml:"multi_statements" mapstructure:"multi_statements"`
	AllowCleartextPasswords  bool              `json:"allow_cleartext_passwords" yaml:"allow_cleartext_passwords" mapstructure:"allow_cleartext_passwords"`
	AllowFallbackToPlaintext bool              `json:"allow_fallback_to_plaintext" yaml:"allow_fallback_to_plaintext" mapstructure:"allow_fallback_to_plaintext"` // like --ssl-mode=PREFERRED
	AllowNativePasswords     *bool             `json:"allow_native_passwords" yaml:"allow_native_passwords" mapstructure:"allow_native_passwords"`                // default: true
	AllowOldPasswords        bool              `json:"allow_old_passwords" yaml:"allow_old_passwords" mapstructure:"allow_old_passwords"`
	ServerPubKey             string            `json:"server_pub_key" yaml:"server_pub_key" mapstructure:"server_pub_key"`                // registered RSA public key name
	CheckConnLiveness        *bool             `json:"check_conn_liveness" yaml:"check_conn_liveness" mapstructure:"check_conn_liveness"` // default: true
	ClientFoundRows          bool              `json:"client_found_rows" yaml:"client_found_rows" mapstructure:"client_found_rows"`       // UPDATE returns matched rows
	ColumnsWithAlias         bool              `json:"columns_with_alias" yaml:"columns_with_alias" mapstructure:"columns_with_alias"`
	RejectReadOnly           bool              `json:"reject_read_only" yaml:"reject_read_only" mapstructure:"reject_read_only"`                // reject read-only connections (AWS Aurora failover)
	ConnectionAttributes     string            `json:"connection_attributes" yaml:"connection_attributes" mapstructure:"connection_attributes"` // "key1:val1,key2:val2"
	Params                   map[string]string `json:"params" yaml:"params" mapstructure:"params"`                                              // session system variables
}
