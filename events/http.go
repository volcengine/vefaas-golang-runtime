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

package events

// HTTPRequest is a wrapper around net/http request,
// it's useful if users want to build web api using vefaas-golang-runtime handler.
type HTTPRequest struct {
	// Method
	HTTPMethod string

	// path, like /abc
	Path string

	// ip:port
	RemoteAddr string

	// parsed path params
	PathParameters map[string]string

	// parsed url query string params
	QueryStringParameters map[string]string

	// http headers
	Headers map[string]string

	// body in raw bytes
	Body []byte
}
