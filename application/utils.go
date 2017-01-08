package application

import (
    "fmt"
    "os"
)

func Exit(errorMessage string) {
    fmt.Fprintln(os.Stderr, errorMessage)
    os.Exit(255)
}
