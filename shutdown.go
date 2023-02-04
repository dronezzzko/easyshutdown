package easyshutdown

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	opentrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

const (
	// Kubernetes (rolling update) doesn't wait until a pod is out of rotation before sending SIGTERM,
	// and external LB could still route traffic to a non-existing pod resulting in a surge of 50x API errors.
	// It's recommended to wait for a while before terminating the program; see for details:
	// https://github.com/kubernetes-retired/contrib/issues/1140.
	preShutdownDelay = 3 * time.Second

	defaultTimeoutDuration = 3 * time.Second
)

type supportedServices struct {
	httpServer     *http.Server
	httpsServer    *http.Server
	grpcServer     *grpc.Server
	tracerProvider *opentrace.TracerProvider
}

type Shutdown struct {
	logger          *log.Logger
	shutdownTimeout time.Duration
	services        *supportedServices
}

func NewShutdown(options ...Option) (*Shutdown, error) {
	srv := &Shutdown{
		shutdownTimeout: defaultTimeoutDuration,
		logger:          log.New(os.Stderr, "easyshutdown ", log.LstdFlags),
		services:        &supportedServices{},
	}

	for _, o := range options {
		o(srv)
	}

	return srv, nil
}

func (s *Shutdown) Graceful() {
	stopCh := signals()
	<-stopCh

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)

	defer cancel()

	s.logger.Println("Shutting down HTTP/HTTPS server")

	// Wait for Kubernetes readiness probe to remove this instance from the LB.
	time.Sleep(preShutdownDelay)

	// stop OpenTelemetry tracer provider
	if s.services.tracerProvider != nil {
		if err := s.services.tracerProvider.Shutdown(ctx); err != nil {
			s.logger.Printf("stopping tracer provider: failed: %s", err.Error())
		}
	}

	if s.services.grpcServer != nil {
		s.logger.Println("Shutting down GRPC server")
		s.services.grpcServer.GracefulStop()
	}

	if s.services.httpServer != nil {
		if err := s.services.httpServer.Shutdown(ctx); err != nil {
			s.logger.Printf("HTTP server graceful shutdown: failed: %s", err.Error())
		}
	}

	if s.services.httpsServer != nil {
		if err := s.services.httpsServer.Shutdown(ctx); err != nil {
			s.logger.Printf("HTTPS server graceful shutdown: failed: %s", err.Error())
		}
	}
}
