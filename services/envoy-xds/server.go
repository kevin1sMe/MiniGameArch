package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	//accesslog "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
)

// Hasher returns node ID as an ID
type Hasher struct {
}

// ID function
func (h Hasher) ID(node *core.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Id
}

// RunAccessLogServer starts an accesslog service.
//func RunAccessLogServer(ctx context.Context, als *AccessLogService, port uint) {
//grpcServer := grpc.NewServer()
//lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
//if err != nil {
//log.WithError(err).Fatal("failed to listen")
//}

//accesslog.RegisterAccessLogServiceServer(grpcServer, als)
//log.WithFields(log.Fields{"port": port}).Info("access log server listening")

//go func() {
//if err = grpcServer.Serve(lis); err != nil {
//log.Error(err)
//}
//}()
//<-ctx.Done()

//grpcServer.GracefulStop()
//}

const grpcMaxConcurrentStreams = 1000000

// RunManagementServer starts an xDS server at the given port.
func RunManagementServer(ctx context.Context, server xds.Server, port uint) {
	// gRPC golang library sets a very small upper bound for the number gRPC/h2
	// streams over a single TCP connection. If a proxy multiplexes requests over
	// a single connection to the management server, then it might lead to
	// availability problems.
	log.Debug(fmt.Sprintf("RunManagementServer enter"))
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions, grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams))
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}

	// register services
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	v2.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	v2.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	v2.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	v2.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	log.WithFields(log.Fields{"port": port}).Info("management server listening")
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Error("goroutine RunManagementServer exit, err:%s", err)
		}
	}()
	log.Debug(fmt.Sprintf("grpc svr start"))
	<-ctx.Done()

	grpcServer.GracefulStop()
	log.Debug(fmt.Sprintf("grpc stopped"))
}

// RunManagementGateway starts an HTTP gateway to an xDS server.
func RunManagementGateway(ctx context.Context, srv xds.Server, port uint) {
	log.WithFields(log.Fields{"port": port}).Info("gateway listening HTTP/1.1")
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: &xds.HTTPGateway{Server: srv}}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()
	log.Debug(fmt.Sprintf("gateway start"))

	<-ctx.Done()
	if err := server.Shutdown(ctx); err != nil {
		log.Error(fmt.Sprintf("shutdown?%s", err))
	}
	log.Debug(fmt.Sprintf("gateway stopped"))
}



