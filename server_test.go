package grpchealth

import (
	"testing"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Make sure the service implementation complies with the proto definition.
func TestRegister(t *testing.T) {
	s := grpc.NewServer()
	hs := NewServer("demo", WithChecksFunc(checksFunc))
	healthpb.RegisterHealthServer(s, hs)
	s.Stop()
}

func checksFunc() healthpb.HealthCheckResponse_ServingStatus {
	// TODO: check your service
	return healthpb.HealthCheckResponse_SERVING
}
