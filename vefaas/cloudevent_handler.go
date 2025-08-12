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
	"net/http"
	"time"

	"github.com/cloudevents/sdk-go/v2/binding"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/volcengine/vefaas-golang-runtime/events"
	"github.com/volcengine/vefaas-golang-runtime/utils"
	"github.com/volcengine/vefaas-golang-runtime/vefaascontext"
)

func handleCloudEvent(handler interface{}) func(rw http.ResponseWriter, rq *http.Request) {
	functionHandler := handler.(cloudeventFunctionHandler)
	return func(rw http.ResponseWriter, rq *http.Request) {
		defer utils.RecoverFunc(rw, nil)

		eventType := rq.Header.Get("X-Faas-Event-Type")
		if eventType != events.EventTypeCloudEvent && eventType != events.EventTypeMultiCloudEvent {
			utils.SetInvalidEventTypeHeader(rw, eventType, events.EventTypeCloudEvent, events.EventTypeMultiCloudEvent)
			return
		}

		ctx := rq.Context()
		ctx = vefaascontext.WithRequestIdContext(ctx, rq.Header.Get("X-Faas-Request-Id"))

		if eventType == events.EventTypeCloudEvent {
			msg := cehttp.NewMessageFromHttpRequest(rq)
			event, err := binding.ToEvent(ctx, msg)
			if err != nil {
				utils.SetInvalidCloudEventHeader(rw, err)
				return
			}

			startTime := time.Now()
			resp, err := functionHandler(ctx, &events.CloudEvent{Event: event})

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
			// In case user set this response header, we rewrite it here.
			utils.SetExecutionDurationHeader(rw, startTime)

			if resp.StatusCode != 0 {
				rw.WriteHeader(resp.StatusCode)
			}
			if resp.Body != nil {
				_, _ = rw.Write(resp.Body)
			}
			return
		}

		// eventType == events.EventTypeMultiCloudEvent
		// if the handler only can handle one event, and request is multi event, will do for-loop
		rawBody, err := utils.RawBodyFromHttpRequest(rq)
		if err != nil {
			return
		}
		batchVersion := utils.GetBatchEventsVersionFromRequest(rq)
		multiEvents, err := utils.TransformRawBodyToBatchCloudEvents(ctx, batchVersion, rawBody)
		if err != nil || len(multiEvents) == 0 {
			utils.SetInvalidCloudEventHeader(rw, err)
			return
		}

		var res *events.EventResponse
		existInvalidStatus := false
		invalidStatus := 200

		startTime := time.Now()
		for _, faasCloudEvent := range multiEvents {
			res, err = functionHandler(ctx, faasCloudEvent)
			utils.SetExecutionDurationHeader(rw, startTime)

			if err != nil {
				utils.SetFunctionExecutionErrorHeader(rw, err)
				return
			}
			if res == nil {
				utils.SetFunctionNoResponseErrorHeader(rw)
				return
			}

			//record invalid status
			//but still call func to process other msgs
			if res.StatusCode >= 400 {
				existInvalidStatus = true
				invalidStatus = res.StatusCode
			}
		}

		//will return last resp as the batch response
		if res.Headers != nil {
			for k, v := range res.Headers {
				rw.Header().Set(k, v)
			}
		}
		// In case user set this response header, we rewrite it here.
		utils.SetExecutionDurationHeader(rw, startTime)
		if existInvalidStatus {
			rw.WriteHeader(invalidStatus)
		} else if res.StatusCode != 0 {
			rw.WriteHeader(res.StatusCode)
		}
		if res.Body != nil {
			_, _ = rw.Write(res.Body)
		}
	}
}
