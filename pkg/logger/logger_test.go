package logger_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	lgr "github.com/jiuchen1986/cks/pkg/logger"
)

var (
	filePathA string = "/tmp/eke.log"
	filePathB string = "/var/log/eke.log"
)

func clean(undo func()) {
	os.RemoveAll(filePathA)
	os.RemoveAll(filePathB)
	undo()
}

func TestLogLevel(t *testing.T) {
	warnOpt, er := lgr.NewLogLevelOption("warn")
	if er != nil {
		t.Fatal(er)
	}
	undo, _, err := lgr.InitLogger(warnOpt)
	if err != nil {
		t.Fatal(err)
	}
	defer clean(undo)

	logger := zap.L()

	ce := logger.Check(zapcore.InfoLevel, "info level msg")
	assert.Nil(t, ce, "Info level msg should be disabled.")
	ce = logger.Check(zapcore.DebugLevel, "debug level msg")
	assert.Nil(t, ce, "Debug level msg should be disabled.")
	ce = logger.Check(zapcore.WarnLevel, "warn level msg")
	assert.NotNil(t, ce, "Warn level msg should be enabled.")
}

func TestMultiLogLevel(t *testing.T) {
	var er error
	warnOpt, er := lgr.NewLogLevelOption("warn")
	if er != nil {
		t.Fatal(er)
	}
	infoOpt, er := lgr.NewLogLevelOption("info")
	if er != nil {
		t.Fatal(er)
	}
	debugOpt, er := lgr.NewLogLevelOption("debug")
	if er != nil {
		t.Fatal(er)
	}

	undo, _, err := lgr.InitLogger(debugOpt, infoOpt, warnOpt)
	if err != nil {
		t.Fatal(err)
	}
	defer clean(undo)

	logger := zap.L()

	ce := logger.Check(zapcore.InfoLevel, "info level msg")
	assert.Nil(t, ce, "Info level msg should be disabled.")
	ce = logger.Check(zapcore.DebugLevel, "debug level msg")
	assert.Nil(t, ce, "Debug level msg should be disabled.")
	ce = logger.Check(zapcore.WarnLevel, "warn level msg")
	assert.NotNil(t, ce, "Warn level msg should be enabled.")
}

func TestOutputPath(t *testing.T) {
	var er error
	var undo func()
	var cfg *zap.Config

	enableLogFileOpt, er := lgr.NewEnableLogFileOption()
	if er != nil {
		t.Fatal(er)
	}
	logFilePathOptA, er := lgr.NewLogFilePathOption(filePathA)
	if er != nil {
		t.Fatal(er)
	}

	undo, cfg, er = lgr.InitLogger(enableLogFileOpt)
	if er != nil {
		t.Fatal(er)
	}
	assert.Equal(t, []string{"stdout"}, cfg.OutputPaths, "Only stdout should be in output path "+
		"when only enable log file.")
	undo()

	undo, cfg, er = lgr.InitLogger(logFilePathOptA, enableLogFileOpt)
	if er != nil {
		t.Fatal(er)
	}
	assert.Equal(t, []string{"stdout", filePathA}, cfg.OutputPaths, "File path should be in output path "+
		"when log file is enalbed and file path is specified.")
	undo()

	undo, cfg, er = lgr.InitLogger(logFilePathOptA)
	if er != nil {
		t.Fatal(er)
	}
	assert.Equal(t, []string{"stdout"}, cfg.OutputPaths, "Only stdout should be in output path "+
		"when log file is not enalbed.")
	clean(undo)
}

func TestMultiOutputPath(t *testing.T) {
	var er error

	enableLogFileOpt, er := lgr.NewEnableLogFileOption()
	if er != nil {
		t.Fatal(er)
	}
	logFilePathOptA, er := lgr.NewLogFilePathOption(filePathA)
	if er != nil {
		t.Fatal(er)
	}
	logFilePathOptB, er := lgr.NewLogFilePathOption(filePathB)
	if er != nil {
		t.Fatal(er)
	}

	undo, cfg, er := lgr.InitLogger(enableLogFileOpt, logFilePathOptA, logFilePathOptB)
	if er != nil {
		t.Fatal(er)
	}
	defer clean(undo)

	assert.Equal(t, []string{"stdout", filePathB}, cfg.OutputPaths, "The file path of last option "+
		"should be in output path.")
}
