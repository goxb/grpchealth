# grpchealth
it is a package that implement grpc health probe v1

### Example
```
    healthServer := grpchealth.NewServer("demo",
		health.WithChecksFunc(checkHealth),
		//health.WithRegularInterval(time.Minute),
	)

    func checksFunc() healthpb.HealthCheckResponse_ServingStatus {
	// TODO: check your service
	return healthpb.HealthCheckResponse_SERVING
}
```
