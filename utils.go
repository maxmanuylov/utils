package utils

func Close(closer io.Closer) {
    _ = closer.Close()
}
