#!/usr/bin/env bash
cd "$(dirname "$0")"
mockgen -destination=mocks/mock_balancer.go -package=mocks google.golang.org/grpc/balancer ClientConn,SubConn
mockgen -destination=mocks/mock_stream.go -package=mocks google.golang.org/grpc ClientStream
