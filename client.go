package grpchealth

import (
	"context"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Client implements client Health.
type Client struct {
	hc healthpb.HealthClient
}

// NewClient return a new client
func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		hc: healthpb.NewHealthClient(conn),
	}
}

// Check return service's status
func (c *Client) Check(ctx context.Context, service string) (healthpb.HealthCheckResponse_ServingStatus, error) {
	reply, err := c.hc.Check(ctx, &healthpb.HealthCheckRequest{Service: service})
	if err != nil {
		return healthpb.HealthCheckResponse_NOT_SERVING, err
	}

	return reply.Status, nil
}
