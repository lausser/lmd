package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kdar/factorlog"
)

const logFormat = "[%{Date} %{Time}][%{Severity}][%{File}:%{Line}] %{Message}"
const logColors = "%{Color \"yellow\" \"WARN\"}%{Color \"red\" \"ERROR\"}"
const logColorReset = "%{Color \"reset\"}"

var log *factorlog.FactorLog

func InitLogging(conf *Config) {
	var logFormatter factorlog.Formatter
	var targetWriter io.Writer
	var err error
	if conf.LogFile == "" {
		logFormatter = factorlog.NewStdFormatter(logColors + logFormat + logColorReset)
		targetWriter = os.Stdout
	} else if strings.ToLower(conf.LogFile) == "stderr" {
		logFormatter = factorlog.NewStdFormatter(logColors + logFormat + logColorReset)
		targetWriter = os.Stderr
	} else {
		logFormatter = factorlog.NewStdFormatter(logFormat)
		if _, err = os.Stat(conf.LogFile); err != nil {
			targetWriter, err = os.Create(conf.LogFile)
		} else {
			targetWriter, err = os.Open(conf.LogFile)
		}
	}
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %s", err.Error()))
	}
	log = factorlog.New(targetWriter, logFormatter)
	var LogLevel = "Warn"
	if conf.LogLevel != "" {
		LogLevel = conf.LogLevel
	}
	if strings.ToLower(LogLevel) == "off" {
		log.SetMinMaxSeverity(factorlog.StringToSeverity("PANIC"), factorlog.StringToSeverity("PANIC"))
	} else {
		log.SetMinMaxSeverity(factorlog.StringToSeverity(strings.ToUpper(LogLevel)), factorlog.StringToSeverity("PANIC"))
	}
}
