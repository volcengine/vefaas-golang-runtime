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
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/volcengine/vefaas-golang-runtime/events"
	"github.com/volcengine/vefaas-golang-runtime/version"
)

var (
	listenPort           = "5000"
	requestTimeoutSecond = 900
)

const (
	startServerExitCode = 170
)

// Start stats vefaas runtime server with provided handler to handle incoming
// requests.
//
// Currently, the supported handler signatures are:
//
// - func(context.Context, *events.HTTPRequest) (*events.EventResponse, error)
// for handling regular http requests, like those from api gateway trigger.
//
// - func(context.Context, *events.CloudEvent) (*events.EventResponse, error)
// for handling CloudEvent request, like those from timer trigger, tos trigger,
// kafka trigger, etc.
//
// - func(context.Context, interface{}) (*events.EventResponse, error)
// for handling requests of any type, especially for those business that
// handle both regular http requests and CloudEvent requests, and the developer
// can use type assertion to distinguish and process them.
func Start(handler interface{}) {
	StartWithInitializer(handler, nil)
}

// StartWithInitializer starts vefaas runtime server with provided
// handler and initializer.
//
// See Start for the supported handler signatures.
//
// The functionality of initializer is to do the initialization work before your
// runtime server can handle any incoming request, like setup the connection to
// database, setup the http/rpc client for downstream services, etc.
// Currently the supported initializer signatures are:
// - func(context.Context) error
func StartWithInitializer(handler interface{}, initializer interface{}) {
	rand.Seed(time.Now().UTC().UnixNano())

	// Validate handler.
	eventType, functionHandler := validateHandler(handler)

	// Validate initializer.
	functionInitializer := validateInitializer(initializer)

	// Initialize metadata.
	if s := os.Getenv("_FAAS_FUNC_TIMEOUT"); s != "" {
		if tmp, err := strconv.Atoi(s); err == nil {
			requestTimeoutSecond = tmp
		}
	}

	// Start http server.
	if s := os.Getenv("_FAAS_RUNTIME_PORT"); s != "" {
		listenPort = s
	}
	listener, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen port with tcp, %v.\n", err)
		os.Exit(startServerExitCode)
	}
	defer listener.Close()
	server := &http.Server{
		Handler: functionServer{
			initializer: functionInitializer,
			handleFunc:  buildHandler(eventType, functionHandler),
		},
	}

	go func() {
		err := server.Serve(listener)
		if err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server exited unexpectedly, %v.\n", err)
			os.Exit(startServerExitCode)
		}
	}()

	// Graceful shutdown.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	// Shutdown http server to close open listeners and idle connections,
	// wait active connections to return to idle and then shut down,
	// and refuse further connections and requests.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(requestTimeoutSecond)*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Shutdown http server error, %v.\n", err)
	}
}

func buildHandler(eventType string, handler interface{}) func(rw http.ResponseWriter, rq *http.Request) {
	switch eventType {
	case events.EventTypeHTTP:
		return handleHttpEvent(handler)
	case events.EventTypeCloudEvent:
		return handleCloudEvent(handler)
	case events.EventTypeAny:
		return handleAnyEvent(handler)
	default:
		return func(rw http.ResponseWriter, rq *http.Request) {
			rw.WriteHeader(http.StatusBadRequest)
		}
	}
}

type functionServer struct {
	initializer interface{}
	handleFunc  func(http.ResponseWriter, *http.Request)
}

func (s functionServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Faas-Internal-Request") == "true" {
		switch r.URL.Path {
		case "/v1/initialize":
			switch r.Method {
			case http.MethodPost:
				err := initializeFunction(context.Background(), s.initializer)
				if err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
				} else {
					rw.WriteHeader(http.StatusOK)
				}
			default:
				rw.WriteHeader(http.StatusMethodNotAllowed)
			}
		case "/v1/version":
			switch r.Method {
			case http.MethodGet:
				rw.WriteHeader(http.StatusOK)
				_, _ = rw.Write([]byte(version.Version))
			default:
				rw.WriteHeader(http.StatusMethodNotAllowed)
			}
		default:
			rw.WriteHeader(http.StatusNotFound)
		}

		return
	}

	s.handleFunc(rw, r)
}
