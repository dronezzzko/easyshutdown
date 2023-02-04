package easyshutdown

import (
	"log"
	"net/http"
	"time"

	opentrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

type Option func(sd *Shutdown)

func (s *Shutdown) WithShutdownTimeout(d time.Duration) Option {
	return func(sd *Shutdown) {
		sd.shutdownTimeout = d
	}
}

func (s *Shutdown) WithLogger(l *log.Logger) Option {
	return func(sd *Shutdown) {
		sd.logger = l
	}
}

func (s *Shutdown) WithHTTPServer(srv *http.Server) Option {
	return func(sd *Shutdown) {
		sd.services.httpServer = srv
	}
}

func (s *Shutdown) WithHTTPSServer(srv *http.Server) Option {
	return func(sd *Shutdown) {
		sd.services.httpsServer = srv
	}
}

func (s *Shutdown) WithGrpcServer(srv *grpc.Server) Option {
	return func(sd *Shutdown) {
		sd.services.grpcServer = srv
	}
}

func (s *Shutdown) WithTracerProvider(t *opentrace.TracerProvider) Option {
	return func(sd *Shutdown) {
		sd.services.tracerProvider = t
	}
}
