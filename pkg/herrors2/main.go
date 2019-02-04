package herrors

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	knownFilePathPatterns = []string{
		"github.com/",
		"bitbucket.org/",
		"code.google.com/",
	}
)

type Error struct {
	error
	keyVals     map[string]interface{}
	prettyStack string
}

type stackFrame struct {
	Filename string
	Method   string
	Line     int
}

func Wrap(err error, keyVals ...interface{}) error {
	e := Error{}

	if err != nil {
		e.error = err
	}

	e.keyVals = mapKeyVals(err, keyVals...)
	e.prettyStack = handlePrettyStack(err)

	return e
}

func New(msg string, keyVals ...interface{}) error {
	e := Error{}

	e.error = errors.New(msg)
	e.keyVals = mapKeyVals(e, keyVals...)
	e.prettyStack = handlePrettyStack(e.error)

	return e
}

func WrappedNew(msg string, keyVals ...interface{}) error {
	return Wrap(New(msg, keyVals...))
}

func (e Error) PrettyStack() string {
	return e.prettyStack
}

func (e Error) KeyVals() map[string]interface{} {
	if e.keyVals == nil {
		return map[string]interface{}{}
	}

	return e.keyVals
}

func mapKeyVals(err error, keyVals ...interface{}) map[string]interface{} {
	kvs := map[string]interface{}{}

	type keyValueGetter interface {
		KeyVals() map[string]interface{}
	}

	if herr, ok := err.(keyValueGetter); ok {
		kvs = herr.KeyVals()
	}

	if len(keyVals) == 0 {
		return kvs
	}

	if len(keyVals)%2 == 1 {
		keyVals = append(keyVals, "<MISSING>")
	}

	l := len(keyVals)

	for i := 0; i < l; i += 2 {
		k := fmt.Sprintf("%s", keyVals[i])
		kvs[k] = keyVals[i+1]
	}

	return kvs
}

func handlePrettyStack(err error) string {
	type prettyStackGetter interface {
		PrettyStack() string
	}

	// Check if the error we received satisfies hobee errors
	herr, ok := err.(prettyStackGetter)
	if !ok {
		return buildPrettyStack()
	}

	return herr.PrettyStack()
}

func ShortenFilePath(s string) string {
	idx := strings.Index(s, "/src/pkg/")
	if idx != -1 {
		return s[idx+5:]
	}
	for _, pattern := range knownFilePathPatterns {
		idx = strings.Index(s, pattern)
		if idx != -1 {
			return s[idx:]
		}
	}
	return s
}

func functionName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "???"
	}
	name := fn.Name()
	end := strings.LastIndex(name, string(os.PathSeparator))
	return name[end+1:]
}

func buildPrettyStack() string {
	var stackLines []string

	sfs := make([]stackFrame, 0)

	// NOTE: 128 is just to protect against the possibility of
	// an infinite loop. In reality the loop will exit much sooner.
	// By starting on nr 3 we avoid the frames introduced in this file.
	for i := 0; i <= 128; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			// End of the stack
			break
		}
		file = ShortenFilePath(file)
		fnc := functionName(pc)

		if strings.HasPrefix(fnc, "herrors") {
			// Stack from herrors. Not very useful! :)
			continue
		}

		if fnc == "runtime.main" || fnc == "http.HandlerFunc.ServeHTTP" {
			// We stop here. Including this and the following lines
			// is never helpful.
			break
		}
		sfs = append(sfs, stackFrame{file, functionName(pc), line})
	}

	for _, sf := range sfs {
		stackLines = append(stackLines, fmt.Sprintf("%s:%d\t%s()", sf.Filename, sf.Line, sf.Method))
	}

	return strings.Join(stackLines, "\n")
}
