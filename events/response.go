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

// EventResponse represents the returned response from a vefaas-golang-runtime handler.
type EventResponse struct {
	// StatusCode is supposed to be a valid HTTP Status code
	StatusCode int

	// Headers contains customized header returned from function.
	Headers map[string]string

	// Body just body, no surprise
	Body []byte
}
