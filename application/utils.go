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

func DumpStackTracesOnSigQuit(appName string) {
    DumpStackTracesOnSigQuitTo(filepath.Join(defaultLogsFolder, appName, "dumps"))
}

func DumpStackTracesOnSigQuitTo(dumpsFolder string) {
    if err := os.MkdirAll(dumpsFolder, os.FileMode(777)); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create dumps folder (%s): %v\n", dumpsFolder, err)
    }

    signalsChannel := make(chan os.Signal, 10)
    notifyOnTermination(signalsChannel)
    signal.Notify(signalsChannel, syscall.SIGQUIT)

    go func() {
        for sig := range signalsChannel {
            if sig != syscall.SIGQUIT {
                break
            }
            go dumpStackTraces(dumpsFolder)
        }
    }()
}

func WaitForTermination() {
    signalsChannel := make(chan os.Signal, 1)
    notifyOnTermination(signalsChannel)
    <-signalsChannel
}

func Exit(errorMessage string) {
    fmt.Fprintln(os.Stderr, errorMessage)
    os.Exit(255)
}

func notifyOnTermination(c chan<- os.Signal) {
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
}

func dumpStackTraces(dumpsFolder string) {
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

    stackTraces := getStackTraces()

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
    timestamp := fmt.Sprintf("%s%03d", now.Format("20060102-150405"), now.Nanosecond() / 1e6)
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

func getStackTraces() []byte {
    buf := make([]byte, 1024 + 1024 * runtime.NumGoroutine())

    for i := 0; i < 10; i++ {
        if i != 0 {
            buf = make([]byte, 2 * len(buf))
        }

        if n := runtime.Stack(buf, true); n < len(buf) {
            return buf[:n]
        }
    }

    return buf
}
