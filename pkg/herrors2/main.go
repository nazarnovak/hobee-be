package herrors

import (
	"errors"
	"fmt"
	"runtime"
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

func (e Error) String() string {
	if e.error == nil {
		return fmt.Sprintf("%s: No errors\n", e.prettyStack)
	}

	return fmt.Sprintf("%s: %s", e.error.Error(), e.prettyStack)
}

//type stackFrame struct {
//	Filename string
//	Method   string
//	Line     int
//}

func Wrap(err error, keyVals ...interface{}) error {
	e := Error{}

	if err != nil {
		e.error = err
	}

	e.keyVals = mapKeyVals(err, keyVals...)
fmt.Println("In wrap:", getCallerPrettyStack())
	e.prettyStack += getCallerPrettyStack()

	return e
}

func New(msg string, keyVals ...interface{}) error {
	e := Error{}

	e.error = errors.New(msg)
	e.keyVals = mapKeyVals(e, keyVals)
fmt.Println("In new keyVals:", mapKeyVals(e, keyVals))
fmt.Println("In new:", getCallerPrettyStack())
	e.prettyStack = getCallerPrettyStack()

	return e
}

func getCallerPrettyStack() string {
	// Skip is set to 2 because this function + the caller of this function which is inside this package
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return fmt.Sprintf("Unknown caller\n")
	}

	fn := runtime.FuncForPC(pc)
	fnName := ""
	if fn != nil {
		fnName = fn.Name()
	}

	return fmt.Sprintf("%s:%d\t%s()\n", file, line, fnName)
}

//func WrappedNew(msg string, keyVals ...interface{}) error {
//	return Wrap(New(msg, keyVals...))
//}
//
//func (e Error) PrettyStack() string {
//	return e.prettyStack
//}

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

//func handlePrettyStack(err error) string {
//	type prettyStackGetter interface {
//		PrettyStack() string
//	}
//
//	// Check if the error we received satisfies hobee errors
//	herr, ok := err.(prettyStackGetter)
//	if !ok {
//		return buildPrettyStack()
//	}
//
//	return herr.PrettyStack()
//}
//
//func ShortenFilePath(s string) string {
//	idx := strings.Index(s, "/src/pkg/")
//	if idx != -1 {
//		return s[idx+5:]
//	}
//	for _, pattern := range knownFilePathPatterns {
//		idx = strings.Index(s, pattern)
//		if idx != -1 {
//			return s[idx:]
//		}
//	}
//	return s
//}
//
//func functionName(pc uintptr) string {
//	fn := runtime.FuncForPC(pc)
//	if fn == nil {
//		return "???"
//	}
//	name := fn.Name()
//	end := strings.LastIndex(name, string(os.PathSeparator))
//	return name[end+1:]
//}
//
//func buildPrettyStack() string {
//	var stackLines []string
//
//	sfs := make([]stackFrame, 0)
//
//	// NOTE: 128 is just to protect against the possibility of
//	// an infinite loop. In reality the loop will exit much sooner.
//	// By starting on nr 3 we avoid the frames introduced in this file.
//	for i := 0; i <= 128; i++ {
//		pc, file, line, ok := runtime.Caller(i)
//		if !ok {
//			// End of the stack
//			break
//		}
//		file = ShortenFilePath(file)
//		fnc := functionName(pc)
//
//		if strings.HasPrefix(fnc, "herrors") {
//			// Stack from herrors. Not very useful! :)
//			continue
//		}
//
//		if fnc == "runtime.main" || fnc == "http.HandlerFunc.ServeHTTP" {
//			// We stop here. Including this and the following lines
//			// is never helpful.
//			break
//		}
//		sfs = append(sfs, stackFrame{file, functionName(pc), line})
//	}
//
//	for _, sf := range sfs {
//		stackLines = append(stackLines, fmt.Sprintf("%s:%d\t%s()", sf.Filename, sf.Line, sf.Method))
//	}
//
//	return strings.Join(stackLines, "\n")
//}

// main -> f1 -> f2

// log(3) -> wrap(2) -> new(1)
// 3 -> show time [level] error\nstacktrace?
// 2 -> wrap - Take existing error, and add a stacktrace to it with the current line. + keyvals
// 1 -> new - Set the error with the message provided to the New(), and add a single stack trace line where New() was
// called. + keyvals

// Wrap existing error
// log(2) -> wrap(1)
// 2 -> show time [level] error\nstacktrace
// 1 -> wrap - Take existing error, and add a stacktrace to it with the current line. + keyvals

// 1) Have the error by either passing a string to New or Wrap an existing error
// 2) Pass the error along if it's a Wrap()
// 3) When using New() or Wrap() - you always add a key/vals map, and override it on the next level
// 4) Start the stacktrace from the New() string slice with [0] index, then add more stack trace strings
// 5) In the end, format it as Time [Level] Error message\n Stacktrace concatenated with \n
