package application

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

type DumpFunc func() []byte
type CancelFunc func()

func DumpOnSigQuit(appName, dumpKind string, dumpFunc DumpFunc) (CancelFunc, error) {
	return DumpOnSigQuitTo(filepath.Join("/var/log", appName, "dumps"), dumpKind, dumpFunc)
}

func DumpOnSigQuitToUserHome(appName, dumpKind string, dumpFunc DumpFunc) (CancelFunc, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Failed to detect user home dir: %s", err.Error())
	}

	return DumpOnSigQuitTo(filepath.Join(userHomeDir, ".dumps", appName), dumpKind, dumpFunc)
}

func DumpOnSigQuitTo(dumpsFolder, dumpKind string, dumpFunc DumpFunc) (CancelFunc, error) {
	if err := CreateDumpsFolder(dumpsFolder); err != nil {
		return nil, err
	}

	return DoOnSigQuit(func() {
		if err := Dump(dumpsFolder, dumpKind, dumpFunc); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
		}
	}), nil
}

func CreateDumpsFolder(dumpsFolder string) error {
	if err := os.MkdirAll(dumpsFolder, os.FileMode(0777)); err != nil {
		return fmt.Errorf("Failed to create dumps folder (%s): %v\n", dumpsFolder, err)
	}
	return nil
}

func DoOnSigQuit(f func()) CancelFunc {
	signalsChannel := make(chan os.Signal, 10)
	notifyOnTermination(signalsChannel)
	signal.Notify(signalsChannel, sigQuit)

	go func() {
		for sig := range signalsChannel {
			if sig != sigQuit {
				break
			}
			go f()
		}
	}()

	return func() {
		signal.Stop(signalsChannel)
		close(signalsChannel)
	}
}

func Dump(dumpsFolder, dumpKind string, dumpFunc DumpFunc) (err error) {
	writers := make([]io.Writer, 1, 2)
	writers[0] = os.Stderr

	dumpFile, _err := createDumpFile(dumpsFolder, dumpKind)
	if _err != nil {
		err = fmt.Errorf("Failed to create dump file: %v\n", _err)
	} else {
		defer func() {
			_ = dumpFile.Sync()
			_ = dumpFile.Close()
		}()
		writers = append(writers, dumpFile)
	}

	dump := dumpFunc()

	_, _ = fmt.Fprintln(os.Stderr)
	_, _ = fmt.Fprintf(os.Stderr, "[============================== Dump[%s] ==============================]\n", dumpKind)
	_, _ = fmt.Fprintln(os.Stderr)

	_, _ = io.MultiWriter(writers...).Write(dump)

	_, _ = fmt.Fprintln(os.Stderr)
	_, _ = fmt.Fprintln(os.Stderr, "[=============================================================================]")
	_, _ = fmt.Fprintln(os.Stderr)

	return
}

func createDumpFile(dumpsFolder, dumpKind string) (*os.File, error) {
	now := time.Now()
	timestamp := fmt.Sprintf("%s%03d", now.Format("20060102-150405"), now.Nanosecond()/1e6)
	pid := os.Getpid()

	for i := 1; i < 100; i++ {
		dumpFile, err := os.OpenFile(makeDumpFilePath(dumpsFolder, dumpKind, timestamp, pid, i), os.O_RDWR|os.O_CREATE|os.O_EXCL, os.FileMode(0666))
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

func makeDumpFilePath(dumpsFolder, dumpKind, timestamp string, pid, attempt int) string {
	var filename string
	if attempt == 1 {
		filename = fmt.Sprintf("%d-%s-%s.txt", pid, timestamp, dumpKind)
	} else {
		filename = fmt.Sprintf("%d-%s-%s.%d.txt", pid, timestamp, dumpKind, attempt)
	}
	return filepath.Join(dumpsFolder, filename)
}
