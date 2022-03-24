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
	"os"
	"reflect"
	"runtime/debug"
)

var functionInitialized bool

type initializer = func(context.Context) error

// validateInitializer validates and creates function initializer, which is in charge
// of the function initialization.
//
// If provided initializerSymbol is not valid, the returned initializer will be a
// function that just report the validation error.
func validateInitializer(initializerSymbol interface{}) (initializer interface{}) {
	if initializerSymbol == nil {
		return
	}

	// Vaidate kind.
	initializerType := reflect.TypeOf(initializerSymbol)
	if initializerType.Kind() != reflect.Func {
		err := fmt.Errorf("expected initializer kind: %s, but got: %s", reflect.Func, initializerType.Kind())
		fmt.Fprintln(os.Stderr, err)
		initializer = initializerErrorFunc(err)
		return
	}

	// Validate arguments.
	err := validateInitializerArguments(initializerType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		initializer = initializerErrorFunc(err)
		return
	}

	// Validate return values.
	err = validateInitializerReturnValues(initializerType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		initializer = initializerErrorFunc(err)
		return
	}

	initializer = initializerSymbol

	return
}

func validateInitializerArguments(initializerType reflect.Type) error {
	if initializerType.NumIn() != 1 {
		return fmt.Errorf("initializer should take one argument (context.Context), but got %d", initializerType.NumIn())
	}
	if !initializerType.In(0).Implements(contextType) {
		return fmt.Errorf("the argument of initializer should implement context.Context, but got %s", initializerType.In(0))
	}

	return nil
}

func validateInitializerReturnValues(initializerType reflect.Type) error {
	if initializerType.NumOut() != 1 {
		return fmt.Errorf("initializer should return one value (error), but got %d", initializerType.NumOut())
	}
	if !initializerType.Out(0).Implements(errorType) {
		return fmt.Errorf("the return value of initializer should implement error, but got %s", initializerType.Out(0))
	}

	return nil
}

func initializerErrorFunc(err error) initializer {
	return func(ctx context.Context) error {
		return err
	}
}

func initializeFunction(ctx context.Context, initializerFunc interface{}) (err error) {
	// No initializer provided.
	if initializerFunc == nil {
		return
	}

	// Return directly if the initializer has been executed successfully.
	if functionInitialized {
		return
	}

	defer func() {
		// Recover from panic if exists.
		if errR := recover(); errR != nil {
			err = fmt.Errorf("panic while initializing function: %v\n%s", errR, string(debug.Stack()))
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	err = initializerFunc.(initializer)(ctx)
	if err != nil {
		err = fmt.Errorf("failed to initialize function, %v", err)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Set function as initialized.
	functionInitialized = true

	return
}
