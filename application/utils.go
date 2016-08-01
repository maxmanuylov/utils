package application

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func WaitForTermination() {
    signalsChannel := make(chan os.Signal, 1)
    signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGTERM)
    <-signalsChannel
}

func Exit(errorMessage string) {
    fmt.Fprintln(os.Stderr, errorMessage)
    os.Exit(255)
}
