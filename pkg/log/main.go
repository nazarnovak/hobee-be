package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"hobee-be/pkg/herrors"
	"strings"
)

const (
	console = "stderr"

	levelCritical = "CRITICAL"
	levelError    = "ERROR"
	levelWarning  = "WARNING"
	levelInfo     = "INFO"
)

var (
	canLogTo  = []string{console}
	levels    = []string{levelCritical, levelError, levelWarning, levelInfo}
	logTo     io.Writer
	logOutput = ""
)

type prettyStacker interface {
	PrettyStack() string
}

type keyValer interface {
	KeyVals() map[string]interface{}
}

func Init(to string) error {
	if to == "" {
		return herrors.New("Specify where to log errors")
	}
	// TODO: What if you take this off? Still good logs?
	// Better log messages for debugging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	for _, v := range canLogTo {
		if to == v {
			switch v {
			case console:
				logTo = os.Stderr
				logOutput = console
			}
		}
	}

	log.SetOutput(logTo)

	if logTo == nil {
		return herrors.New(fmt.Sprintf("Cannot log errors to: %s", to))
	}

	return nil
}

func logConsoleError(level string, err error) {
	fmt.Printf("%s [%s] %s\n", time.Now().UTC().Format("06-01-02 15:04:05.999"), level, err)
}

func logKeyVals(err error) {
	herr, ok := err.(keyValer)
	if !ok {
		return
	}

	keyVals := herr.KeyVals()
	if len(keyVals) == 0 {
		return
	}

	out := []string{}

	for k, v := range keyVals {
		out = append(out, fmt.Sprintf("%#v: %+v", k, v))
	}

	fmt.Println(strings.Join(out, ", "))
}

func logPrettyStack(err error) {
	herr, ok := err.(prettyStacker)
	if !ok {
		return
	}

	prettyStack := herr.PrettyStack()
	if prettyStack == "" {
		return
	}

	fmt.Println(prettyStack)
}

func logError(err error) {
	fmt.Printf("")
}

func Critical(ctx context.Context, err error) {
	switch logOutput {
	case console:
		logConsoleError(levelCritical, err)
		logKeyVals(err)
		logPrettyStack(err)
	}
}

func Error(ctx context.Context, err error) {
	switch logOutput {
	case console:
		logConsoleError(levelError, err)
		logKeyVals(err)
		logPrettyStack(err)
	}
}

func Warning(ctx context.Context, err error) {
	switch logOutput {
	case console:
		logConsoleError(levelWarning, err)
		logKeyVals(err)
		logPrettyStack(err)
	}
}

func Info(ctx context.Context, err error) {
	switch logOutput {
	case console:
		logConsoleError(levelInfo, err)
		logKeyVals(err)
		logPrettyStack(err)
	}
}
