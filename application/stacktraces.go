package application

import (
	"runtime"
)

var stackTracesDumpFunc = func() []byte { return GetStackTraces() }

func DumpStackTracesOnSigQuit(appName string) (CancelFunc, error) {
	return DumpOnSigQuit(appName, "goroutines", stackTracesDumpFunc)
}

func DumpStackTracesOnSigQuitToUserHome(appName string) (CancelFunc, error) {
	return DumpOnSigQuitToUserHome(appName, "goroutines", stackTracesDumpFunc)
}

func DumpStackTracesOnSigQuitTo(dumpsFolder string) (CancelFunc, error) {
	return DumpOnSigQuitTo(dumpsFolder, "goroutines", stackTracesDumpFunc)
}

func GetStackTraces() []byte {
	buf := make([]byte, 1024+1024*runtime.NumGoroutine())

	for i := 0; i < 10; i++ {
		if i != 0 {
			buf = make([]byte, 2*len(buf))
		}

		if n := runtime.Stack(buf, true); n < len(buf) {
			return buf[:n]
		}
	}

	return buf
}
