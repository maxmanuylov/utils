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

func Switch(apps map[string]func()) {
    if len(os.Args) < 2 {
        Exit("Application is not specified")
    }

    appName := os.Args[1]
    os.Args = append(os.Args[:1], os.Args[2:]...)

    app, ok := apps[appName]
    if !ok {
        Exit(fmt.Sprintf("Unknown application: %s", appName))
    }

	DumpStackTracesOnSigQuit(appName)

	app()
}

func Exit(errorMessage string) {
    fmt.Fprintln(os.Stderr, errorMessage)
    os.Exit(255)
}
