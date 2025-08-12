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

package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cloudevents/sdk-go/v2/binding"
	cloudevents_v2 "github.com/cloudevents/sdk-go/v2/event"
	cehttp_v2 "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/volcengine/vefaas-golang-runtime/events"
)

func SetHttpParamsAndHeaders(req *events.HTTPRequest, rq *http.Request) {
	if len(rq.URL.Query()) > 0 {
		for k, v := range rq.URL.Query() {
			req.QueryStringParameters[k] = strings.Join(v, ",")
		}
	}
	if len(rq.Header) > 0 {
		for k, v := range rq.Header {
			req.Headers[k] = strings.Join(v, ",")
		}
	}
}

func GetBatchEventsVersionFromRequest(rq *http.Request) string {
	batchEventsVersion := rq.Header.Get("X-Bytefaas-Batch-Events-Version")
	if batchEventsVersion == "" {
		batchEventsVersion = rq.Header.Get("X-Faas-Batch-Events-Version")
	}
	return batchEventsVersion
}

func TransformRawBodyToBatchCloudEvents(ctx context.Context, batchEventsVersion string, rawBody []byte) ([]*events.CloudEvent, error) {
	var es []*events.CloudEvent
	if batchEventsVersion == "1.0" {
		eventsV2, err := transformBytesToCloudEventsV2(ctx, rawBody)
		if err != nil {
			return nil, err
		}
		for _, event := range eventsV2 {
			es = append(es, &events.CloudEvent{Event: event})
		}
		return es, nil
	} else {
		return nil, fmt.Errorf("batch events version %s is not supported", batchEventsVersion)
	}
}

func transformBytesToCloudEventsV2(ctx context.Context, bytesData []byte) ([]*cloudevents_v2.Event, error) {
	buf := bytes.NewBuffer(bytesData)

	totalNumBytes := buf.Next(4)
	totalNum := transformBytesToInt32(totalNumBytes)

	multiEvents := make([]*cloudevents_v2.Event, 0)

	var i uint32
	for i = 0; i < totalNum; i++ {
		headerLenBytes := buf.Next(4)
		headerLen := transformBytesToInt32(headerLenBytes)
		headerBytes := buf.Next(int(headerLen))
		header := transformBytesToHeader(headerBytes)

		bodyLenBytes := buf.Next(4)
		bodyLen := transformBytesToInt32(bodyLenBytes)
		body := buf.Next(int(bodyLen))

		event, err := transformHeadersAndBodyToCloudEventV2(ctx, header, body)
		if err != nil {
			return nil, err
		}

		multiEvents = append(multiEvents, event)
	}

	return multiEvents, nil
}

func transformBytesToInt32(bytesData []byte) uint32 {
	var v uint32
	v = uint32(bytesData[3])
	v = v | (uint32(bytesData[2]) << 8)
	v = v | (uint32(bytesData[1]) << 16)
	v = v | (uint32(bytesData[0]) << 24)
	return v
}

func transformBytesToHeader(bytesData []byte) http.Header {
	header := http.Header{}
	mulHeaderBytes := bytes.Split(bytesData, []byte("\r\n"))
	for _, singleHeaderBytes := range mulHeaderBytes {
		if len(singleHeaderBytes) == 0 {
			continue
		}
		kv := bytes.Split(singleHeaderBytes, []byte(":"))
		if len(kv) == 2 {
			k := string(kv[0])
			vv := string(kv[1])
			for _, v := range strings.Split(vv, ",") {
				header.Add(k, string(v))
			}
		}
	}
	return header
}

func transformHeadersAndBodyToCloudEventV2(ctx context.Context, headers http.Header, body []byte) (event *cloudevents_v2.Event, err error) {
	readerCloser := io.NopCloser(bytes.NewReader(body))
	msg := cehttp_v2.NewMessage(headers, readerCloser)
	event, err = binding.ToEvent(ctx, msg)
	if err != nil {
		err = fmt.Errorf("failed to parse cloudevent in batch request, err: %v", err)
	}
	return
}
