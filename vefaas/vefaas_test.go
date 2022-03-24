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
	"encoding/json"
	"fmt"
	"log"

	"github.com/volcengine/vefaas-golang-runtime/events"
	"github.com/volcengine/vefaas-golang-runtime/vefaascontext"
)

// ExampleStart demonstrated how to start a simple vefaas function serving http request.
func ExampleStart() {
	// Define your handler.
	handler := func(ctx context.Context, r *events.HTTPRequest) (*events.EventResponse, error) {
		log.Printf("request id: %v", vefaascontext.RequestIdFromContext(ctx))
		log.Printf("request headers: %v", r.Headers)

		body, _ := json.Marshal(map[string]string{"message": "Hello veFaaS!"})
		return &events.EventResponse{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: body,
		}, nil
	}

	// Start your vefaas function =D.
	Start(handler)
}

// ExampleStartWithInitializer shows how to start a simple vefaas function with some initialization logic.
func ExampleStartWithInitializer() {
	// First define your handler.
	handler := func(ctx context.Context, r *events.HTTPRequest) (*events.EventResponse, error) {
		body, _ := json.Marshal(map[string]string{"message": "Hello veFaaS!"})
		return &events.EventResponse{
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: body,
		}, nil
	}

	// Then we have a dummy initializer.
	initializer := func(ctx context.Context) error {
		fmt.Println("init function...")
		// Some initialization work, like setup http client, create database connections, etc.
		fmt.Println("init function done")
		return nil
	}

	// Start your vefaas function =D.
	StartWithInitializer(handler, initializer)
}
