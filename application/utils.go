package application

import (
    "fmt"
    "os"
)

func Run(app func() error) {
    if err := app(); err != nil {
        Exit(err.Error())
    }
}

func Exit(errorMessage string) {
    fmt.Fprintln(os.Stderr, errorMessage)
    os.Exit(255)
}
