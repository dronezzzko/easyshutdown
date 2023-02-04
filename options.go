package easyshutdown

import (
	"log"
	"net/http"
	"time"

	opentrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

type Option func(sd *Shutdown)

func WithShutdownTimeout(d time.Duration) Option {
	return func(sd *Shutdown) {
		sd.shutdownTimeout = d
	}
}

func WithLogger(l *log.Logger) Option {
	return func(sd *Shutdown) {
		sd.logger = l
	}
}

func WithHTTPServer(srv *http.Server) Option {
	return func(sd *Shutdown) {
		sd.services.httpServer = srv
	}
}

func WithHTTPSServer(srv *http.Server) Option {
	return func(sd *Shutdown) {
		sd.services.httpsServer = srv
	}
}

func WithGrpcServer(srv *grpc.Server) Option {
	return func(sd *Shutdown) {
		sd.services.grpcServer = srv
	}
}

func WithTracerProvider(t *opentrace.TracerProvider) Option {
	return func(sd *Shutdown) {
		sd.services.tracerProvider = t
	}
}
