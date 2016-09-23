package tls_transport

import (
    "crypto/tls"
    "crypto/x509"
    "errors"
    "github.com/maxmanuylov/utils/http/transport"
    "net/http"
)

func New(caCert, clientCert, clientKey []byte) (*http.Transport, error) {
    tlsConfig, err := newTLSConfig(caCert, clientCert, clientKey)
    if err != nil {
        return nil, err
    }

    transport := http_transport.NewDefault()
    transport.TLSClientConfig = tlsConfig

    return transport, nil
}

func newTLSConfig(caCert, clientCert, clientKey []byte) (*tls.Config, error) {
    caPool := x509.NewCertPool()
    if !caPool.AppendCertsFromPEM(caCert) {
        return nil, errors.New("No CA certificate found")
    }

    certificate, err := tls.X509KeyPair(clientCert, clientKey)
    if err != nil {
        return nil, err
    }

    return &tls.Config{
        Certificates: []tls.Certificate{certificate},
        RootCAs: caPool,
    }, nil
}
