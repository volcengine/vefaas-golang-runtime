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
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/volcengine/vefaas-golang-runtime/events"
)

var (
	contextType       = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorType         = reflect.TypeOf((*error)(nil)).Elem()
	httpRequestType   = reflect.TypeOf(&events.HTTPRequest{})
	eventResponseType = reflect.TypeOf(&events.EventResponse{})
	cloudEventType    = reflect.TypeOf(&events.CloudEvent{})
)

type (
	// httpFunctionHandler handles regular http request, like api gateway trigger.
	httpFunctionHandler = func(context.Context, *events.HTTPRequest) (*events.EventResponse, error)

	// cloudeventFunctionHandler handles CloudEvent request, like timer event,
	// tos event, kafka event, etc.
	cloudeventFunctionHandler = func(context.Context, *events.CloudEvent) (*events.EventResponse, error)

	// anyFunctionHandler handles requests of any type, it might be used when
	// you want to handle both regular http request and other event trigger requests
	// like timer, tos, kafka, etc.
	anyFunctionHandler = func(context.Context, interface{}) (*events.EventResponse, error)
)

// validateHandler validates and creates the base function handler, which will do
// basic payload unmarshaling before deferring to handlerSymbol.
//
// If provided handlerSymbol is not valid, the returned handler will be a function
// that just report the validation error.
func validateHandler(handlerSymbol interface{}) (eventType string, functionHandler interface{}) {
	var err error
	if handlerSymbol == nil {
		err = errors.New("expected a handler function, but got nil")
		fmt.Fprintln(os.Stderr, err)
		eventType = events.EventTypeAny
		functionHandler = errorHandler(err)
		return
	}

	// Vaidate kind.
	handlerType := reflect.TypeOf(handlerSymbol)
	if handlerType.Kind() != reflect.Func {
		err = fmt.Errorf("expected handler kind: %s, but got: %s", reflect.Func, handlerType.Kind())
		fmt.Fprintln(os.Stderr, err)
		eventType = events.EventTypeAny
		functionHandler = errorHandler(err)
		return
	}

	// Validate arguments.
	eventType, err = validateHandlerArguments(handlerType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		eventType = events.EventTypeAny
		functionHandler = errorHandler(err)
		return
	}

	// Validate return values.
	err = validateHandlerReturnValues(handlerType, eventType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		eventType = events.EventTypeAny
		functionHandler = errorHandler(err)
		return
	}

	functionHandler = handlerSymbol

	return
}

func validateHandlerArguments(handler reflect.Type) (eventType string, err error) {
	if handler.NumIn() != 2 {
		err = fmt.Errorf("handler should take two arguments, but got %d", handler.NumIn())
		return
	}

	if !handler.In(0).Implements(contextType) {
		err = fmt.Errorf("the first argument of handler does not implement context.Context")
		return
	}
	if handler.In(1).Kind() == reflect.Interface {
		eventType = events.EventTypeAny
		return
	}
	if handler.In(1) == httpRequestType {
		eventType = events.EventTypeHTTP
		return
	}
	if handler.In(1) == cloudEventType {
		eventType = events.EventTypeCloudEvent
		return
	}

	err = fmt.Errorf("the second argument of handler should be one of "+
		"(*events.HTTPRequest, *events.CloudEvent, interface{}), but got %s", handler.In(1))
	return
}

func validateHandlerReturnValues(handler reflect.Type, eventType string) error {
	if handler.NumOut() != 2 {
		return fmt.Errorf("handler should return two values, but got %d", handler.NumOut())
	}
	if handler.Out(0) != eventResponseType {
		return fmt.Errorf("the first return value of handler should be events.EventResponse, but got %s", handler.Out(0))
	}
	if !handler.Out(1).Implements(errorType) {
		return fmt.Errorf("the second return value of handler should implement error, but got %s", handler.Out(1))
	}

	return nil
}

func errorHandler(err error) anyFunctionHandler {
	return func(ctx context.Context, payload interface{}) (*events.EventResponse, error) {
		return nil, err
	}
}
