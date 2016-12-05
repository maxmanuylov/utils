package http_transport

import (
    "net"
    "net/http"
)

func NewDefault() *http.Transport {
    defaultTransport := &http.Transport{}
    *defaultTransport = *http.DefaultTransport.(*http.Transport) // copying
    return defaultTransport
}

func NewUnix(socketFile string) *http.Transport {
    transport := NewDefault()

    originalDial := transport.Dial

    transport.Dial = func(_, _ string) (net.Conn, error) {
        return originalDial("unix", socketFile)
    }

    return transport
}
