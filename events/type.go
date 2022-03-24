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

// Supported event type. When users call vefaas.Start(), we actually started a http web server. Incoming events/requests
// are simply just http request. We use headers (X-Faas-Event-Type) to determine whether this is a regular http request
// or a event trigger (timer trigger, tos trigger, etc).
//
// Not intended to be used by regular vefaas-golang-runtime users. FaaS Native http users (running native http apps on vefaas)
// may want to use the value comparing with the X-Faas-Event-Type header, to determine if their request is coming from
// vefaas triggers.
const (
	EventTypeAny = "any"

	// EventTypeHTTP represents http request from vefaas provided HTTP triggers (or other guys calling you with mesh http egress)
	EventTypeHTTP = "http"

	// EventTypeCloudEvent represents trigger events (timer/kafka/rocketmq/tos/abase_binlog, etc)
	EventTypeCloudEvent = "cloudevent"
)
