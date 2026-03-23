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
 * cache.go
 *
 * - Author   : Frost Leo <frostleo.dev@gmail.com>
 * - Created  : 2026-03-23
 * - Modified : 2026-03-23
 */

package setting

import "time"

// Cache holds all named cache source configurations.
type Cache struct {
	Sources map[string]CacheSource `json:"sources" yaml:"sources" mapstructure:"sources"`
}

// CacheSource represents a single cache connection configuration.
type CacheSource struct {
	Driver string        `json:"driver" yaml:"driver" mapstructure:"driver"` // "redis"
	Prefix string        `json:"prefix" yaml:"prefix" mapstructure:"prefix"` // key namespace prefix
	TTL    time.Duration `json:"ttl" yaml:"ttl" mapstructure:"ttl"`          // default TTL for cache entries

	Redis *RedisOptions `json:"redis,omitempty" yaml:"redis,omitempty" mapstructure:"redis"`
}

// RedisOptions maps to go-redis/v9 Options.
type RedisOptions struct {
	Network               string        `json:"network" yaml:"network" mapstructure:"network"`    // "tcp" | "unix", default: "tcp"
	Addr                  string        `json:"addr" yaml:"addr" mapstructure:"addr"`             // "host:port"
	Username              string        `json:"username" yaml:"username" mapstructure:"username"` // Redis 6.0+ ACL
	Password              string        `json:"password" yaml:"password" mapstructure:"password"`
	DB                    int           `json:"db" yaml:"db" mapstructure:"db"`                                                       // database number
	Protocol              int           `json:"protocol" yaml:"protocol" mapstructure:"protocol"`                                     // 2 or 3, default: 3
	ClientName            string        `json:"client_name" yaml:"client_name" mapstructure:"client_name"`                            // CLIENT SETNAME
	MaxRetries            int           `json:"max_retries" yaml:"max_retries" mapstructure:"max_retries"`                            // default: 3, -1 disables
	MinRetryBackoff       time.Duration `json:"min_retry_backoff" yaml:"min_retry_backoff" mapstructure:"min_retry_backoff"`          // default: 8ms, -1 disables
	MaxRetryBackoff       time.Duration `json:"max_retry_backoff" yaml:"max_retry_backoff" mapstructure:"max_retry_backoff"`          // default: 512ms, -1 disables
	DialTimeout           time.Duration `json:"dial_timeout" yaml:"dial_timeout" mapstructure:"dial_timeout"`                         // default: 5s
	ReadTimeout           time.Duration `json:"read_timeout" yaml:"read_timeout" mapstructure:"read_timeout"`                         // default: 3s, -1 no timeout, -2 disable deadline
	WriteTimeout          time.Duration `json:"write_timeout" yaml:"write_timeout" mapstructure:"write_timeout"`                      // default: 3s, -1 no timeout, -2 disable deadline
	DialerRetries         int           `json:"dialer_retries" yaml:"dialer_retries" mapstructure:"dialer_retries"`                   // default: 5
	DialerRetryTimeout    time.Duration `json:"dialer_retry_timeout" yaml:"dialer_retry_timeout" mapstructure:"dialer_retry_timeout"` // default: 100ms
	ContextTimeoutEnabled bool          `json:"context_timeout_enabled" yaml:"context_timeout_enabled" mapstructure:"context_timeout_enabled"`
	ReadBufferSize        int           `json:"read_buffer_size" yaml:"read_buffer_size" mapstructure:"read_buffer_size"`    // default: 32KiB
	WriteBufferSize       int           `json:"write_buffer_size" yaml:"write_buffer_size" mapstructure:"write_buffer_size"` // default: 32KiB
	PoolFIFO              bool          `json:"pool_fifo" yaml:"pool_fifo" mapstructure:"pool_fifo"`                         // default: false (LIFO)
	PoolSize              int           `json:"pool_size" yaml:"pool_size" mapstructure:"pool_size"`                         // default: 10 * GOMAXPROCS
	PoolTimeout           time.Duration `json:"pool_timeout" yaml:"pool_timeout" mapstructure:"pool_timeout"`                // default: ReadTimeout + 1s
	MinIdleConns          int           `json:"min_idle_conns" yaml:"min_idle_conns" mapstructure:"min_idle_conns"`          // default: 0
	MaxIdleConns          int           `json:"max_idle_conns" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`          // default: 0
	MaxActiveConns        int           `json:"max_active_conns" yaml:"max_active_conns" mapstructure:"max_active_conns"`    // default: 0 (unlimited)
	MaxConcurrentDials    int           `json:"max_concurrent_dials" yaml:"max_concurrent_dials" mapstructure:"max_concurrent_dials"`
	ConnMaxIdleTime       time.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time" mapstructure:"conn_max_idle_time"`                   // default: 30m
	ConnMaxLifetime       time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`                      // default: 0 (no limit)
	ConnMaxLifetimeJitter time.Duration `json:"conn_max_lifetime_jitter" yaml:"conn_max_lifetime_jitter" mapstructure:"conn_max_lifetime_jitter"` // default: 0
	TLS                   *RedisTLS     `json:"tls,omitempty" yaml:"tls,omitempty" mapstructure:"tls"`
	DisableIdentity       bool          `json:"disable_identity" yaml:"disable_identity" mapstructure:"disable_identity"` // disable CLIENT SETINFO
	IdentitySuffix        string        `json:"identity_suffix" yaml:"identity_suffix" mapstructure:"identity_suffix"`
	FailingTimeoutSeconds int           `json:"failing_timeout_seconds" yaml:"failing_timeout_seconds" mapstructure:"failing_timeout_seconds"` // default: 15
}

// RedisTLS provides TLS configuration for Redis connections.
type RedisTLS struct {
	Enabled            bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	CACertFile         string `json:"ca_cert_file" yaml:"ca_cert_file" mapstructure:"ca_cert_file"`
	CertFile           string `json:"cert_file" yaml:"cert_file" mapstructure:"cert_file"`       // client cert for mTLS
	KeyFile            string `json:"key_file" yaml:"key_file" mapstructure:"key_file"`          // client key for mTLS
	ServerName         string `json:"server_name" yaml:"server_name" mapstructure:"server_name"` // override SNI
	InsecureSkipVerify bool   `json:"insecure_skip_verify" yaml:"insecure_skip_verify" mapstructure:"insecure_skip_verify"`
	MinVersion         string `json:"min_version" yaml:"min_version" mapstructure:"min_version"` // "1.2" | "1.3", default: "1.2"
}

func (c *Cache) Validate() error {
	// TODO: validate driver type, redis addr format, TTL bounds, etc.
	return nil
}
