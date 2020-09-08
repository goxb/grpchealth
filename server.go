package grpchealth

import (
	"context"
	"time"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// ServerOption is server option
type ServerOption func(*serverOptions)

// WithChecksFunc set checks function
func WithChecksFunc(checksFunc func() healthpb.HealthCheckResponse_ServingStatus) ServerOption {
	return func(o *serverOptions) {
		o.checksFunc = checksFunc
	}
}

// WithRegularInterval set checking health interval
func WithRegularInterval(d time.Duration) ServerOption {
	return func(o *serverOptions) {
		o.interval = d
	}
}

type serverOptions struct {
	checksFunc func() healthpb.HealthCheckResponse_ServingStatus
	interval   time.Duration
}

func applyServerOptions(opts ...ServerOption) *serverOptions {
	o := &serverOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Server implements service Health.
type Server struct {
	service  string
	hs       *health.Server
	checkFn  func() healthpb.HealthCheckResponse_ServingStatus
	interval time.Duration
}

// NewServer returns a health Server.
func NewServer(service string, opts ...ServerOption) *Server {
	o := applyServerOptions(opts...)

	//checkFn: fn,
	s := &Server{
		service:  service,
		hs:       health.NewServer(),
		checkFn:  o.checksFunc,
		interval: o.interval,
	}

	if s.interval > 0 {
		go s.checksAtRegularInterval()
	}

	return s
}

// Check implements `service Health`.
func (s *Server) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	// regular checks
	if s.interval > 0 {
		return s.hs.Check(ctx, in)
	}

	if s.checkFn != nil {
		servingStatus := s.checkFn()
		s.hs.SetServingStatus(s.service, servingStatus)
	}

	return s.hs.Check(ctx, in)
}

// Watch implements `service Health`.
func (s *Server) Watch(in *healthpb.HealthCheckRequest, stream healthpb.Health_WatchServer) error {
	return s.hs.Watch(in, stream)
}

// SetServingStatus is called when need to reset the serving status of a service
// or insert a new service entry into the statusMap.
func (s *Server) SetServingStatus(servingStatus healthpb.HealthCheckResponse_ServingStatus) {
	s.hs.SetServingStatus(s.service, servingStatus)
}

func (s *Server) checksAtRegularInterval() {
	for {
		servingStatus := s.checkFn()
		s.hs.SetServingStatus(s.service, servingStatus)
		time.Sleep(s.interval)
	}
}
