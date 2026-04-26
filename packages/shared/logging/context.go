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
 * @File    : context.go
 * @Author  : Frost Leo <frostleo.dev@gmail.com>
 * @Created : 2026-04-27
 * @Modified: 2026-04-27
 */

package logging

import "context"

type correlationContextKey struct{}

// ContextWithCorrelation stores correlation values on ctx.
func ContextWithCorrelation(ctx context.Context, correlation Correlation) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, correlationContextKey{}, correlation)
}

// ContextWithMergedCorrelation merges correlation values into ctx.
func ContextWithMergedCorrelation(ctx context.Context, correlation Correlation) context.Context {
	base, _ := CorrelationFromContext(ctx)
	return ContextWithCorrelation(ctx, MergeCorrelation(base, correlation))
}

// CorrelationFromContext retrieves correlation values from ctx.
func CorrelationFromContext(ctx context.Context) (Correlation, bool) {
	if ctx == nil {
		return Correlation{}, false
	}
	correlation, ok := ctx.Value(correlationContextKey{}).(Correlation)
	return correlation, ok
}

// MergeCorrelation returns base with non-empty overlay values applied.
func MergeCorrelation(base Correlation, overlay Correlation) Correlation {
	if overlay.TraceID != "" {
		base.TraceID = overlay.TraceID
	}
	if overlay.SpanID != "" {
		base.SpanID = overlay.SpanID
	}
	if overlay.TraceFlags != "" {
		base.TraceFlags = overlay.TraceFlags
	}
	if overlay.RequestID != "" {
		base.RequestID = overlay.RequestID
	}
	if overlay.CorrelationID != "" {
		base.CorrelationID = overlay.CorrelationID
	}
	if overlay.JobID != "" {
		base.JobID = overlay.JobID
	}
	if overlay.TaskID != "" {
		base.TaskID = overlay.TaskID
	}
	if overlay.WorkflowID != "" {
		base.WorkflowID = overlay.WorkflowID
	}
	if overlay.PluginID != "" {
		base.PluginID = overlay.PluginID
	}
	if overlay.PluginRunID != "" {
		base.PluginRunID = overlay.PluginRunID
	}
	return base
}
