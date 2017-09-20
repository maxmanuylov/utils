package application

import (
    "fmt"
    "os"
    "strings"
)

func Run(app func() error) {
    if err := app(); err != nil {
        Exit(err.Error())
    }
}

func Switch(apps map[string]func()) {
    doSwitch("Application", apps)
}

func SwitchCommand(apps map[string]func()) {
    doSwitch("Command", apps)
}

func doSwitch(kind string, apps map[string]func()) {
    if len(os.Args) < 2 {
        Exit(fmt.Sprintf("%s is not specified", kind))
    }

    appName := os.Args[1]
    os.Args = append(os.Args[:1], os.Args[2:]...)

    app, ok := apps[appName]
    if !ok {
        lowerKind := strings.ToLower(kind)

        appNames := make([]string, 0, len(apps))
        for app := range apps {
            appNames = append(appNames, app)
        }

        Exit(fmt.Sprintf("Unknown %s: %s; available %ss are \"%s\"", lowerKind, appName, lowerKind, strings.Join(appNames, "\", \"")))
    }

	DumpStackTracesOnSigQuit(appName)

	app()
}

func Exit(errorMessage string) {
    fmt.Fprintln(os.Stderr, errorMessage)
    os.Exit(255)
}
