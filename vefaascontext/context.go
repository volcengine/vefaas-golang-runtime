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
	"net/http"
)

const (
	requestIdHeader       = "X-Faas-Request-Id"
	accessKeyIdHeader     = "X-Faas-Access-Key-Id"
	secretAccessKeyHeader = "X-Faas-Secret-Access-Key"
	sessionTokenHeader    = "X-Faas-Session-Token"
)

type contextKey int

const (
	requestIdContextKey contextKey = iota
	accessKeyIdContextKey
	secretAccessKeyContextKey
	sessionTokenContextKey
)

// WithRequestIdContext stores request id into context.
func WithRequestIdContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestIdContextKey, req.Header.Get(requestIdHeader))
}

// RequestIdFromContext retrieves request id from context.
func RequestIdFromContext(ctx context.Context) (id string) {
	if ctx != nil {
		id, _ = ctx.Value(requestIdContextKey).(string)
	}

	return
}

// WithAccessKeyIdContext stores access key id into context.
func WithAccessKeyIdContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, accessKeyIdContextKey, req.Header.Get(accessKeyIdHeader))
}

// AccessKeyIdFromContext retrieves access key id from context.
func AccessKeyIdFromContext(ctx context.Context) (id string) {
	if ctx != nil {
		id, _ = ctx.Value(accessKeyIdContextKey).(string)
	}

	return
}

// WithSecretAccessKeyContext stores secret access key into context.
func WithSecretAccessKeyContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, secretAccessKeyContextKey, req.Header.Get(secretAccessKeyHeader))
}

// SecretAccessKeyFromContext retrieves secret access key from context.
func SecretAccessKeyFromContext(ctx context.Context) (id string) {
	if ctx != nil {
		id, _ = ctx.Value(secretAccessKeyContextKey).(string)
	}

	return
}

// WithSessionTokenContext stores session token into context.
func WithSessionTokenContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, sessionTokenContextKey, req.Header.Get(sessionTokenHeader))
}

// SessionTokenFromContext retrieves session token from context.
func SessionTokenFromContext(ctx context.Context) (id string) {
	if ctx != nil {
		id, _ = ctx.Value(sessionTokenContextKey).(string)
	}

	return
}
