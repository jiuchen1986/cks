/*
Copyright © 2020 Xin Chen <devops.chen@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package error

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/jiuchen1986/cks/pkg/utils"
)

// HandleErr handles error and executes a series of clean functions might need.
// Those clean functions should be without input, and any output are ignored here
type HandleErr func(error, ...func())

// ExitOnErr is a global HandleErr function
// that will finally cause the whole program exit.
// ExitOnErr is set to detailExitOnErr initially,
// and might be set to other HandleErr in UpdateErrHandling
var ExitOnErr HandleErr

func init() {
	ExitOnErr = HandleErr(detailExitOnErr)
}

// a HandleErr that prints out error and causes program exit
func simpleExitOnErr(err error, cleanFuncs ...func()) {
	if err != nil {
		utils.Printf("exit on fatal error: %s\n", err)
		for _, f := range cleanFuncs {
			if f != nil {
				f()
			}
		}
		os.Exit(1)
	}
}

// a HandleErr that prints out error details and causes program exit
func detailExitOnErr(err error, cleanFuncs ...func()) {
	if err != nil {

		s := err.Error()
		// only give stack trace info for the root error if the error is wrapped
		for ; errors.Unwrap(err) != nil; err = errors.Unwrap(err) {
		}
		utils.Printf("exit on fatal error: %+v\n", errors.WithMessage(err, s))
		for _, f := range cleanFuncs {
			if f != nil {
				f()
			}
		}
		os.Exit(1)
	}
}

// update this map when new type is added
var exitOnErrMap map[string]HandleErr = map[string]HandleErr{
	"simple": simpleExitOnErr,
	"detail": detailExitOnErr,
}

// PrintAvailExitOnErr returns a string listing all supported types
// of ExitOnErr seperated by comma, e.g. "simple, detail, stack"
func PrintAvailExitOnErr() string {
	s := ""
	for k := range exitOnErrMap {
		s = fmt.Sprintf("%s\"%s\", ", s, k)
	}

	return s[:len(s)-2]
}

// UpdateErrHandling updates configuration for error handling.
// Currently only change ExitOnErr is supported
func UpdateErrHandling(t string) error {
	utils.Println("start to init error handling.")

	if h, ok := exitOnErrMap[t]; ok {
		utils.Printf("%s information is enabled in ExitOnErr.", t)
		ExitOnErr = h
		return nil
	}

	return errors.Errorf("unknown type of ExitOnErr: %s, only support %s", t, PrintAvailExitOnErr())
}
