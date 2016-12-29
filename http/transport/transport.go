package http_transport

import (
    "context"
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

    originalDialContext := transport.DialContext

    transport.DialContext = func(ctx context.Context, _, _ string) (net.Conn, error) {
        return originalDialContext(ctx, "unix", socketFile)
    }

    return transport
}
