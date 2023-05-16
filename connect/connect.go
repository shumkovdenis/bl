package connect

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
)

func NewInsecureClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
			// Don't forget timeouts!
		},
	}
}

func WithHandlerOptions(interceptors ...connect.Interceptor) connect.HandlerOption {
	return connect.WithHandlerOptions(
		connect.WithInterceptors(
			InjectTraceContext(),
			InjectTraceContextLogger(),
		),
		connect.WithInterceptors(interceptors...),
	)
}

func WithClientOptions(interceptors ...connect.Interceptor) connect.ClientOption {
	return connect.WithClientOptions(
		connect.WithGRPC(),
		connect.WithInterceptors(interceptors...),
	)
}
