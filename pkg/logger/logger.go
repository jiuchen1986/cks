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
package logger

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/jiuchen1986/cks/pkg/utils"
)

// InitLogger initialize the zap global logger with log options
// and return an function undo the changes of the zap global logger,
// which should be call before the whole program exits.
// Supported levels include debug, info, warn, error, panic.
// The return config is in purpose of test, it should be removed until
// a better way to test is found
func InitLogger(options ...LogOption) (func(), *zap.Config, error) {
	utils.Println("start to init log system.")

	cfg := zap.NewDevelopmentConfig()

	// use stdout as default output path
	cfg.OutputPaths = []string{"stdout"}

	// disable the stack trace as it's developed by ourselves
	// see Error function below for details
	cfg.DisableStacktrace = true

	opts, bitmap := tidyOptsAndCalBitmap(options)

	// apply options
	for _, v := range opts {
		if err := v.ConfigLogger(bitmap, &cfg); err != nil {
			return nil, nil, errors.Wrap(err, "failed to apply options")
		}
	}

	logger, err := cfg.Build(zap.WithCaller(true), zap.AddCallerSkip(1))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build global logger")
	}

	return zap.ReplaceGlobals(logger), &cfg, nil
}

// if same options are provided, make sure only the last one takes effect
// also bitmap is calculated
func tidyOptsAndCalBitmap(opts []LogOption) (map[LogOptionType]LogOption, LogOptionType) {
	result := map[LogOptionType]LogOption{}
	bitmap := LogOptionType(0)
	for _, opt := range opts {
		result[opt.OptType()] = opt
		bitmap = bitmap | opt.OptType()
	}
	return result, bitmap
}

// StructuredLogger is an interface wraps the zap logger
type StructuredLogger interface {
	Sync() error
	Info(string, ...map[string]string)
	Warn(string, ...map[string]string)
	Error(string, ...error)
	Debug(string, ...map[string]string)
	Panic(string, ...error)
}

type wrapLogger struct {
	logger *zap.Logger
}

// GetGlobalStructuredLogger returns a wrapped global zap logger
func GetGlobalStructuredLogger() StructuredLogger {
	return &wrapLogger{logger: zap.L()}
}

// Sync is wrapped sync of zap logger
func (w *wrapLogger) Sync() error {
	return w.logger.Sync()
}

// Info is wrapped Info of zap logger
// with input of the msg and a map for desired fields
// where keys and values in the map are both string
func (w *wrapLogger) Info(msg string, fields ...map[string]string) {
	fs := mapToStringFields(fields)
	w.logger.Info(msg, fs...)
}

// Warn is wrapped Warn of zap logger
// with input of the msg and a map for desired fields
// where keys and values in the map are both string
func (w *wrapLogger) Warn(msg string, fields ...map[string]string) {
	fs := mapToStringFields(fields)
	w.logger.Warn(msg, fs...)
}

// Error is wrapped Error of zap logger
// with input of the msg and a slice of error
func (w *wrapLogger) Error(msg string, fields ...error) {
	fs := sliceToErrorFields(fields)
	w.logger.Error(msg, fs...)
}

// Debug is wrapped Debug of zap logger
// with input of the msg and a map for desired fields
// where keys and values in the map are both string
func (w *wrapLogger) Debug(msg string, fields ...map[string]string) {
	fs := mapToStringFields(fields)
	w.logger.Debug(msg, fs...)
}

// Panic is wrapped Panic of zap logger
// with input of the msg and a slice of error
func (w *wrapLogger) Panic(msg string, fields ...error) {
	fs := sliceToErrorFields(fields)
	w.logger.Panic(msg, fs...)
}

func mapToStringFields(fields []map[string]string) []zap.Field {
	fs := []zap.Field{}
	for _, m := range fields {
		for k, v := range m {
			fs = append(fs, zap.String(k, v))
		}
	}
	return fs
}

func sliceToErrorFields(fields []error) []zap.Field {
	fs := []zap.Field{}
	for _, e := range fields {
		fs = append(fs, zap.Error(e))
	}
	return fs
}

// SugaredLogger is an interface wraps the zap sugared logger
type SugaredLogger interface {
	Sync() error
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Debug(...interface{})
	Panic(...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Debugf(string, ...interface{})
	Panicf(string, ...interface{})
	UnwrappedStackErrorf(error, string, ...interface{})
}

type wrapSugaredLogger struct {
	logger *zap.SugaredLogger
}

// GetGlobalLogger returns a wrapped global zap sugared logger
func GetGlobalLogger() SugaredLogger {
	return &wrapSugaredLogger{logger: zap.S()}
}

// Sync wraps up the Sync of zap.SugaredLogger
func (w *wrapSugaredLogger) Sync() error {
	return w.logger.Sync()
}

// Info wraps up the Info of zap.SugaredLogger
func (w *wrapSugaredLogger) Info(args ...interface{}) {
	w.logger.Info(args...)
}

// Warn wraps up the Warn of zap.SugaredLogger
func (w *wrapSugaredLogger) Warn(args ...interface{}) {
	w.logger.Warn(args...)
}

// Error wraps up the Error of zap.SugaredLogger
func (w *wrapSugaredLogger) Error(args ...interface{}) {
	w.logger.Error(args...)
}

// Debug wraps up the Debug of zap.SugaredLogger
func (w *wrapSugaredLogger) Debug(args ...interface{}) {
	w.logger.Debug(args...)
}

// Panic wraps up the Panic of zap.SugaredLogger
func (w *wrapSugaredLogger) Panic(args ...interface{}) {
	w.logger.Panic(args...)
}

// Infof wraps up the Infof of zap.SugaredLogger
func (w *wrapSugaredLogger) Infof(template string, args ...interface{}) {
	w.logger.Infof(template, args...)
}

// Warnf wraps up the Warnf of zap.SugaredLogger
func (w *wrapSugaredLogger) Warnf(template string, args ...interface{}) {
	w.logger.Warnf(template, args...)
}

// Errorf wraps up the Errorf of zap.SugaredLogger
func (w *wrapSugaredLogger) Errorf(template string, args ...interface{}) {
	w.logger.Errorf(template, args...)
}

// Debugf wraps up the Debugf of zap.SugaredLogger
func (w *wrapSugaredLogger) Debugf(template string, args ...interface{}) {
	w.logger.Infof(template, args...)
}

// Panicf wraps up the Panicf of zap.SugaredLogger
func (w *wrapSugaredLogger) Panicf(template string, args ...interface{}) {
	w.logger.Panicf(template, args...)
}

// UnwrappedStackErrorf wraps up the Errorf of zap.SugaredLogger by unwrapping
// the input error and output the stack trace info of the unwrapped error.
// Note that user should not use "%+v" to print out any detailed information of any errors
// in template. Otherwise, please directly use the Errorf
func (w *wrapSugaredLogger) UnwrappedStackErrorf(err error, template string, args ...interface{}) {
	// only give stack trace info for the root error if the error is wrapped
	for ; errors.Unwrap(err) != nil; err = errors.Unwrap(err) {
	}
	w.logger.Errorf(template, args...)
	fmt.Printf("%+v\n", err)
}
