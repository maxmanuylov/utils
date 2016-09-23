package http_transport

import (
    "net/http"
)

func NewDefault() *http.Transport {
    defaultTransport := &http.Transport{}
    *defaultTransport = *http.DefaultTransport.(*http.Transport) // copying
    return defaultTransport
}
