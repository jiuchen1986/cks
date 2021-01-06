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
	"go.uber.org/zap/zapcore"
)

// LogOptionType indicates configurable log options
type LogOptionType uint16

const (
	// LogLevelOpt is used to configure global log level
	LogLevelOpt LogOptionType = 1 << iota

	// EnableLogFileOpt is used to enable local log file.
	// Local log file is enabled only if
	// both EnableLogFileOpt and LogFilePathOpt exist.
	// Value of EnableLogFileOpt is ignored.
	// Note that currently no internal rotation supported
	EnableLogFileOpt

	// LogFilePathOpt is used to configure
	// the path for local log file
	LogFilePathOpt
)

// LogOption is used to configure global logger behaviors
type LogOption interface {
	// Type get the type of this option
	OptType() LogOptionType
	// ConfigLogger configure the logger according to this option
	ConfigLogger(LogOptionType, *zap.Config) error
}

type logLevelOption struct {
	Level zapcore.Level
}

// mapping from string to zap log level,
// update this map when new level is added
var logLevelMap map[string]zapcore.Level = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
}

// PrintAvailLogLevel returns a string listing all supported log level
// seperated by comma, e.g. "debug, info, warn, ..."
func PrintAvailLogLevel() string {
	s := ""
	for k := range logLevelMap {
		s = fmt.Sprintf("%s\"%s\", ", s, k)
	}

	return s[:len(s)-2]
}

// NewLogLevelOption returns a logLevelOption with specified log level
// and an error if exists
func NewLogLevelOption(lvl string) (LogOption, error) {
	if l, ok := logLevelMap[lvl]; ok {
		return &logLevelOption{Level: l}, nil
	}
	return nil, errors.Errorf("unsupported log level: %s, only support %s", lvl, PrintAvailLogLevel())
}

func (ll *logLevelOption) OptType() LogOptionType {
	return LogLevelOpt
}

func (ll *logLevelOption) ConfigLogger(opts LogOptionType, cfg *zap.Config) error {
	cfg.Level.SetLevel(ll.Level)
	return nil
}

type enableLogFileOption struct{}

// NewEnableLogFileOption returns an enableLogFileOption
// which enables local log file if file path is specified
func NewEnableLogFileOption() (LogOption, error) {
	return &enableLogFileOption{}, nil
}

func (elf *enableLogFileOption) OptType() LogOptionType {
	return EnableLogFileOpt
}

func (elf *enableLogFileOption) ConfigLogger(opts LogOptionType, cfg *zap.Config) error {
	// do nothing here because only logFilePathOption can
	// finally enable the local log file
	return nil
}

type logFilePathOption struct {
	FilePath string
}

// NewLogFilePathOption returns an logFilePathOption
// which enables local log file and configures the file path
func NewLogFilePathOption(p string) (LogOption, error) {
	return &logFilePathOption{FilePath: p}, nil
}

func (lfp *logFilePathOption) OptType() LogOptionType {
	return LogFilePathOpt
}

func (lfp *logFilePathOption) ConfigLogger(opts LogOptionType, cfg *zap.Config) error {
	enable := EnableLogFileOpt | LogFilePathOpt
	if (opts & enable) == enable {
		cfg.OutputPaths = append(cfg.OutputPaths, lfp.FilePath)
	}
	return nil
}
