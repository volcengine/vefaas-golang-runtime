/*
 * Copyright 2022 Volcengine
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vefaascontext

import (
	"context"
)

type contextKey int

const (
	requestIdContextKey contextKey = iota
)

// WithRequestIdContext stores request id into context.
func WithRequestIdContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIdContextKey, id)
}

// RequestIdFromContext retrieves request id from context.
func RequestIdFromContext(ctx context.Context) (id string) {
	if ctx != nil {
		id, _ = ctx.Value(requestIdContextKey).(string)
	}

	return
}
