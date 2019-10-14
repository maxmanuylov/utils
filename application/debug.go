package application

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

const defaultLogsFolder = "/var/log"

func DumpStackTracesOnSigQuit(appName string) error {
	return DumpStackTracesOnSigQuitTo(filepath.Join(defaultLogsFolder, appName, "dumps"))
}

func DumpStackTracesOnSigQuitToUserHome(relativePath string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Failed to detect user home dir: %s", err.Error())
	}

	return DumpStackTracesOnSigQuitTo(filepath.Join(userHomeDir, relativePath))
}

func DumpStackTracesOnSigQuitTo(dumpsFolder string) error {
	if err := os.MkdirAll(dumpsFolder, os.FileMode(777)); err != nil {
		return fmt.Errorf("Failed to create dumps folder (%s): %v\n", dumpsFolder, err)
	}

	DoOnSigQuit(func() {
		DumpStackTraces(dumpsFolder)
	})

	return nil
}

func DoOnSigQuit(f func()) {
	signalsChannel := make(chan os.Signal, 10)
	notifyOnTermination(signalsChannel)
	signal.Notify(signalsChannel, syscall.SIGQUIT)

	go func() {
		for sig := range signalsChannel {
			if sig != syscall.SIGQUIT {
				break
			}
			go f()
		}
	}()
}

func DumpStackTraces(dumpsFolder string) {
	writers := make([]io.Writer, 1, 2)
	writers[0] = os.Stderr

	dumpFile, err := createDumpFile(dumpsFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create dump file: %v\n", err)
	} else {
		defer dumpFile.Close()
		defer dumpFile.Sync()
		writers = append(writers, dumpFile)
	}

	stackTraces := GetStackTraces()

	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "[============================== Goroutines Dump ==============================]")
	fmt.Fprintln(os.Stderr)

	io.MultiWriter(writers...).Write(stackTraces)

	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "[=============================================================================]")
	fmt.Fprintln(os.Stderr)
}

func createDumpFile(dumpsFolder string) (*os.File, error) {
	now := time.Now()
	timestamp := fmt.Sprintf("%s%03d", now.Format("20060102-150405"), now.Nanosecond()/1e6)
	pid := os.Getpid()

	for i := 1; i < 100; i++ {
		dumpFile, err := os.OpenFile(makeDumpFilePath(dumpsFolder, timestamp, pid, i), os.O_RDWR|os.O_CREATE|os.O_EXCL, os.FileMode(0666))
		if err != nil {
			if os.IsExist(err) {
				continue
			}
			return nil, err
		}
		return dumpFile, nil
	}

	return nil, fmt.Errorf("Failed to find unique file name for timestamp = %s", timestamp)
}

func makeDumpFilePath(dumpsFolder, timestamp string, pid, attempt int) string {
	var filename string
	if attempt == 1 {
		filename = fmt.Sprintf("%d-%s-goroutines.txt", pid, timestamp)
	} else {
		filename = fmt.Sprintf("%d-%s-goroutines.%d.txt", pid, timestamp, attempt)
	}
	return filepath.Join(dumpsFolder, filename)
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
