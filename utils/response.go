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
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

func RecoverFunc(rw http.ResponseWriter, callback func()) {
	if err := recover(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "panic: %v\n%s", err, debug.Stack())
		rw.Header().Set(
			"X-Faas-Response-Error-Code", "function_panic",
		)
		rw.Header().Set(
			"X-Faas-Response-Error-Message",
			"Function panic, please check log for more details.",
		)
		rw.WriteHeader(http.StatusInternalServerError)

		if callback != nil {
			callback()
		}
	}
}

func RawBodyFromHttpRequest(r *http.Request) ([]byte, error) {
	var rawBody bytes.Buffer
	if r.Body != nil {
		var err error
		_, err = io.Copy(&rawBody, r.Body)
		if err != nil {
			return nil, err
		}
	}

	return rawBody.Bytes(), nil
}

func SetInvalidCloudEventHeader(rw http.ResponseWriter, err error) {
	rw.Header().Set(
		"X-Faas-Response-Error-Code", "invalid_cloud_event",
	)
	rw.Header().Set(
		"X-Faas-Response-Error-Message",
		fmt.Sprintf(`The request is not valid cloudevent message, %v.`, err),
	)
	rw.WriteHeader(http.StatusBadRequest)
}

func SetFunctionExecutionErrorHeader(rw http.ResponseWriter, err error) {
	rw.Header().Set(
		"X-Faas-Response-Error-Code", "function_execution_error",
	)
	rw.Header().Set(
		"X-Faas-Response-Error-Message",
		fmt.Sprintf(`Function returns error, %v.`, err),
	)
	rw.WriteHeader(http.StatusInternalServerError)
}

func SetFunctionNoResponseErrorHeader(rw http.ResponseWriter) {
	rw.Header().Set(
		"X-Faas-Response-Error-Code", "function_no_response",
	)
	rw.Header().Set(
		"X-Faas-Response-Error-Message",
		"No response was returned from function.",
	)
	rw.WriteHeader(http.StatusInternalServerError)
}

func SetInvalidEventTypeHeader(rw http.ResponseWriter, t string, expectedT ...string) {
	rw.Header().Set(
		"X-Faas-Response-Error-Code", "invalid_event_type",
	)
	rw.Header().Set(
		"X-Faas-Response-Error-Message",
		fmt.Sprintf(`The request event type "%s" is not acceptable, expected type "%v".`, t, expectedT),
	)
	rw.WriteHeader(http.StatusBadRequest)
}

func SetExecutionDurationHeader(rw http.ResponseWriter, startTime time.Time) {
	rw.Header().Set("X-Faas-Execution-Duration",
		fmt.Sprintf("%.2f", float64(time.Since(startTime).Nanoseconds())/float64(time.Millisecond)))
}
