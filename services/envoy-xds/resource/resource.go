// Copyright 2018 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// Package resource creates test xDS resources
package resouce

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	//als "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v2"
	//alf "github.com/envoyproxy/go-control-plane/envoy/config/filter/accesslog/v2"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	tcp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/tcp_proxy/v2"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"time"
)

const (
	//localhost = "127.0.0.1"
	localhost = "0.0.0.0"

	// XdsCluster is the cluster name for the control server (used by non-ADS set-up)
	XdsCluster = "xds_cluster"

	// Ads mode for resources: one aggregated xDS service
	Ads = "ads"

	// Xds mode for resources: individual xDS services
	Xds = "xds"

	// Rest mode for resources: polling using Fetch
	Rest = "rest"
)

var (
	// RefreshDelay for the polling config source
	RefreshDelay = 500 * time.Millisecond
)

// MakeEndpoint creates a endpoint on a given port.
func MakeEndpoint(clusterName string, addr string, port uint32) *v2.ClusterLoadAssignment {
	return &v2.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []endpoint.LocalityLbEndpoints{{
			LbEndpoints: []endpoint.LbEndpoint{{
				Endpoint: &endpoint.Endpoint{
					Address: &core.Address{
						Address: &core.Address_SocketAddress{
							SocketAddress: &core.SocketAddress{
								Protocol: core.TCP,
								Address:  addr,
								PortSpecifier: &core.SocketAddress_PortValue{
									PortValue: port,
								},
							},
						},
					},
				},
			}},
		}},
	}
}

// MakeCluster creates a cluster using either ADS or EDS.
// func MakeCluster(mode string, clusterName string, codeType v2.HttpConnectionManager_CodecType) *v2.Cluster {
func MakeCluster(mode string, clusterName string) *v2.Cluster {
	var edsSource *core.ConfigSource
	switch mode {
	case Ads:
		edsSource = &core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_Ads{
				Ads: &core.AggregatedConfigSource{},
			},
		}
	case Xds:
		edsSource = &core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
				ApiConfigSource: &core.ApiConfigSource{
					ApiType: core.ApiConfigSource_GRPC,
					GrpcServices: []*core.GrpcService{{
						TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: XdsCluster},
						},
					}},
				},
			},
		}
	case Rest:
		edsSource = &core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
				ApiConfigSource: &core.ApiConfigSource{
					ApiType:      core.ApiConfigSource_REST,
					ClusterNames: []string{XdsCluster},
					RefreshDelay: &RefreshDelay,
				},
			},
		}
	}

	return &v2.Cluster{
		Name:           clusterName,
		ConnectTimeout: 5 * time.Second,
		Type:           v2.Cluster_EDS,
		EdsClusterConfig: &v2.Cluster_EdsClusterConfig{
			EdsConfig: edsSource,
		},
		// CommonHttpProtocolOptions: &core.HttpProtocolOptions{},
		Http2ProtocolOptions: &core.Http2ProtocolOptions{},
		LbPolicy:             v2.Cluster_ROUND_ROBIN,
	}
}

func MakeClusterHTTP1(mode string, clusterName string) *v2.Cluster {
	var edsSource *core.ConfigSource
	switch mode {
	case Ads:
		edsSource = &core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_Ads{
				Ads: &core.AggregatedConfigSource{},
			},
		}
	case Xds:
		edsSource = &core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
				ApiConfigSource: &core.ApiConfigSource{
					ApiType: core.ApiConfigSource_GRPC,
					GrpcServices: []*core.GrpcService{{
						TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: XdsCluster},
						},
					}},
				},
			},
		}
	case Rest:
		edsSource = &core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
				ApiConfigSource: &core.ApiConfigSource{
					ApiType:      core.ApiConfigSource_REST,
					ClusterNames: []string{XdsCluster},
					RefreshDelay: &RefreshDelay,
				},
			},
		}
	}

	return &v2.Cluster{
		Name:           clusterName,
		ConnectTimeout: 5 * time.Second,
		Type:           v2.Cluster_EDS,
		EdsClusterConfig: &v2.Cluster_EdsClusterConfig{
			EdsConfig: edsSource,
		},
		// CommonHttpProtocolOptions: &core.HttpProtocolOptions{},
		// Http2ProtocolOptions: &core.Http2ProtocolOptions{},
		LbPolicy:             v2.Cluster_ROUND_ROBIN,
	}
}


func MakeRouteRule(clusterName string, prefix string, decorator string, header_matcher []*route.HeaderMatcher) route.Route {
	return route.Route{
		Match: route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Prefix{
				Prefix: prefix,
			},
			Headers: header_matcher,
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterName,
				},
			},
		},

		Decorator: &route.Decorator{
			Operation: decorator, //"checkAvailability"
		},
	}
}

func MakeRoute(routeName string, routeRules []route.Route) *v2.RouteConfiguration {
	return &v2.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []route.VirtualHost{{
			Name:    routeName,
			Domains: []string{"*"},
			Routes:  routeRules,
		}},
	}
}

//func MakeRoute(routeName, clusterName string, prefix string, decorator string, header_matcher []*route.HeaderMatcher) *v2.RouteConfiguration {
//return &v2.RouteConfiguration{
//Name: routeName,
//VirtualHosts: []route.VirtualHost{{
//Name:    routeName,
//Domains: []string{"*"},
//Routes: []route.Route{{
//Match: route.RouteMatch{
//PathSpecifier: &route.RouteMatch_Prefix{
//Prefix: prefix,
//},
//Headers: header_matcher,
//},
//Action: &route.Route_Route{
//Route: &route.RouteAction{
//ClusterSpecifier: &route.RouteAction_Cluster{
//Cluster: clusterName,
//},
//},
//},

//Decorator: & route.Decorator{
//Operation: decorator,   //"checkAvailability"
//},

//}},
//}},
//}
//}

// MakeHTTPListener creates a listener using either ADS or RDS for the route.
func MakeHTTPListener(mode string, listenerName string, port uint32, route string) *v2.Listener {
	// data source configuration
	rdsSource := core.ConfigSource{}
	switch mode {
	case Ads:
		rdsSource.ConfigSourceSpecifier = &core.ConfigSource_Ads{
			Ads: &core.AggregatedConfigSource{},
		}
	case Xds:
		rdsSource.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
			ApiConfigSource: &core.ApiConfigSource{
				ApiType: core.ApiConfigSource_GRPC,
				GrpcServices: []*core.GrpcService{{
					TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
						EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: XdsCluster},
					},
				}},
			},
		}
	case Rest:
		rdsSource.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
			ApiConfigSource: &core.ApiConfigSource{
				ApiType:      core.ApiConfigSource_REST,
				ClusterNames: []string{XdsCluster},
				RefreshDelay: &RefreshDelay,
			},
		}
	}

	// access log service configuration
	//alsConfig := &als.HttpGrpcAccessLogConfig{
	//CommonConfig: &als.CommonGrpcAccessLogConfig{
	//LogName: "echo",
	//GrpcService: &core.GrpcService{
	//TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
	//EnvoyGrpc: &core.GrpcService_EnvoyGrpc{
	//ClusterName: XdsCluster,
	//},
	//},
	//},
	//},
	//}
	//alsConfigPbst, err := util.MessageToStruct(alsConfig)
	//if err != nil {
	//panic(err)
	//}

	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType: hcm.AUTO,
		//StatPrefix: "http",
		StatPrefix: "ingress_http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    rdsSource,
				RouteConfigName: route,
			},
		},
		HttpFilters: []*hcm.HttpFilter{
			// {
			// Name: util.GRPCHTTP1Bridge,
		//}, {
		{
			Name: util.Router,
		}},
		Tracing: &hcm.HttpConnectionManager_Tracing{
			OperationName: hcm.EGRESS,
		},
		//AccessLog: []*alf.AccessLog{{
		//Name:   util.HTTPGRPCAccessLog,
		//Config: alsConfigPbst,
		//}},
	}
	pbst, err := util.MessageToStruct(manager)
	if err != nil {
		panic(err)
	}

	return &v2.Listener{
		Name: listenerName,
		Address: core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.TCP,
					Address:  localhost,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: port,
					},
				},
			},
		},
		FilterChains: []listener.FilterChain{{
			Filters: []listener.Filter{{
				Name:   util.HTTPConnectionManager,
				Config: pbst,
			}},
		}},
	}
}

// MakeTCPListener creates a TCP listener for a cluster.
func MakeTCPListener(listenerName string, port uint32, clusterName string) *v2.Listener {
	// TCP filter configuration
	config := &tcp.TcpProxy{
		StatPrefix: "tcp",
		Cluster:    clusterName,
	}
	pbst, err := util.MessageToStruct(config)
	if err != nil {
		panic(err)
	}
	return &v2.Listener{
		Name: listenerName,
		Address: core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.TCP,
					Address:  localhost,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: port,
					},
				},
			},
		},
		FilterChains: []listener.FilterChain{{
			Filters: []listener.Filter{{
				Name:   util.TCPProxy,
				Config: pbst,
			}},
		}},
	}
}

// Generate produces a snapshot from the parameters.
//func (ts TestSnapshot) Generate() cache.Snapshot {
//clusters := make([]cache.Resource, ts.NumClusters)
//endpoints := make([]cache.Resource, ts.NumClusters)
//for i := 0; i < ts.NumClusters; i++ {
//name := fmt.Sprintf("cluster-%s-%d", ts.Version, i)
//clusters[i] = MakeCluster(ts.Xds, name)
//endpoints[i] = MakeEndpoint(name, localhost, ts.UpstreamPort)
//}

//routes := make([]cache.Resource, ts.NumHTTPListeners)
//for i := 0; i < ts.NumHTTPListeners; i++ {
//name := fmt.Sprintf("route-%s-%d", ts.Version, i)
//routes[i] = MakeRoute(name, cache.GetResourceName(clusters[i%ts.NumClusters]))
//}

//total := ts.NumHTTPListeners + ts.NumTCPListeners
//listeners := make([]cache.Resource, total)
//for i := 0; i < total; i++ {
//port := ts.BasePort + uint32(i)
//// listener name must be same since ports are shared and previous listener is drained
//name := fmt.Sprintf("listener-%d", port)
//if i < ts.NumHTTPListeners {
//listeners[i] = MakeHTTPListener(ts.Xds, name, port, cache.GetResourceName(routes[i]))
//} else {
//listeners[i] = MakeTCPListener(name, port, cache.GetResourceName(clusters[i%ts.NumClusters]))
//}
//}

//return cache.NewSnapshot(ts.Version, endpoints, clusters, routes, listeners)
//}
