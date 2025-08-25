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

package vefaas

import (
	"fmt"
	"net/http"
	"time"

	"github.com/volcengine/vefaas-golang-runtime/events"
	"github.com/volcengine/vefaas-golang-runtime/utils"
	"github.com/volcengine/vefaas-golang-runtime/vefaascontext"
)

func handleHttpEvent(handler interface{}) func(rw http.ResponseWriter, rq *http.Request) {
	functionHandler := handler.(httpFunctionHandler)
	return func(rw http.ResponseWriter, rq *http.Request) {
		defer utils.RecoverFunc(rw, nil)

		eventType := rq.Header.Get("X-Faas-Event-Type")
		if eventType != "" && eventType != events.EventTypeHTTP {
			utils.SetInvalidEventTypeHeader(rw, eventType, events.EventTypeHTTP)
			return
		}

		ctx := rq.Context()
		ctx = vefaascontext.WithRequestIdContext(ctx, rq)
		ctx = vefaascontext.WithAccessKeyIdContext(ctx, rq)
		ctx = vefaascontext.WithSecretAccessKeyContext(ctx, rq)
		ctx = vefaascontext.WithSessionTokenContext(ctx, rq)

		rawBody, err := utils.RawBodyFromHttpRequest(rq)
		if err != nil {
			return
		}

		remoteAddr := rq.RemoteAddr
		remoteIP := rq.Header.Get("X-Real-Ip")
		remotePort := rq.Header.Get("X-Real-Port")
		if remoteIP != "" && remotePort != "" {
			remoteAddr = fmt.Sprintf("%s:%s", remoteIP, remotePort)
		}

		req := &events.HTTPRequest{
			HTTPMethod:            rq.Method,
			Path:                  rq.URL.Path,
			RemoteAddr:            remoteAddr,
			PathParameters:        make(map[string]string),
			QueryStringParameters: make(map[string]string),
			Headers:               make(map[string]string),
			Body:                  rawBody,
		}
		utils.SetHttpParamsAndHeaders(req, rq)

		startTime := time.Now()
		resp, err := functionHandler(ctx, req)
		utils.SetExecutionDurationHeader(rw, startTime)
		if err != nil {
			utils.SetFunctionExecutionErrorHeader(rw, err)
			return
		}

		if resp == nil {
			utils.SetFunctionNoResponseErrorHeader(rw)
			return
		}

		if resp.Headers != nil {
			for k, v := range resp.Headers {
				rw.Header().Set(k, v)
			}
		}
		// In case user set this response header, we need to rewrite it here.
		utils.SetExecutionDurationHeader(rw, startTime)

		if resp.StatusCode != 0 {
			rw.WriteHeader(resp.StatusCode)
		}
		if resp.Body != nil {
			_, _ = rw.Write(resp.Body)
		}
	}
}
