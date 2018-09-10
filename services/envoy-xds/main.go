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

// Package main contains the test driver for testing xDS manually.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	//"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/gogo/protobuf/proto"
    resource "envoy-xds/resource"
    consul "envoy-xds/consul"
	//"github.com/envoyproxy/go-control-plane/pkg/test"
	//"github.com/envoyproxy/go-control-plane/pkg/test/resource"
)

var (
	debug bool

	port         uint
	gatewayPort  uint
	upstreamPort uint
	basePort     uint
	alsPort      uint

	delay    time.Duration
	requests int
	updates  int

	mode          string
	clusters      int
	httpListeners int
	tcpListeners  int

	nodeID string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Use debug logging")
	flag.UintVar(&port, "port", 18000, "Management server port")
	flag.UintVar(&gatewayPort, "gateway", 18001, "Management server port for HTTP gateway")
	flag.UintVar(&upstreamPort, "upstream", 18080, "Upstream HTTP/1.1 port")
	flag.UintVar(&basePort, "base", 9000, "Listener port")
	flag.UintVar(&alsPort, "als", 18090, "Accesslog server port")
	flag.DurationVar(&delay, "delay", 500*time.Millisecond, "Interval between request batch retries")
	flag.IntVar(&requests, "r", 5, "Number of requests between snapshot updates")
	flag.IntVar(&updates, "u", 3, "Number of snapshot updates")
	flag.StringVar(&mode, "xds", resource.Ads, "Management server type (ads, xds, rest)")
	flag.IntVar(&clusters, "clusters", 4, "Number of clusters")
	flag.IntVar(&httpListeners, "http", 2, "Number of HTTP listeners (and RDS configs)")
	flag.IntVar(&tcpListeners, "tcp", 2, "Number of TCP pass-through listeners")
	flag.StringVar(&nodeID, "nodeID", "test-id", "Node ID")
}

// main returns code 1 if any of the batches failed to pass all requests
func main() {

	flag.Parse()
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	//ctx, cancel := context.WithTimeout(context.Background(), 5000 * time.Second)
	//defer cancel()
	ctx := context.Background()

	// start upstream
	//go test.RunHTTP(ctx, upstreamPort)

	// create a cache
	signal := make(chan struct{})
	cb := &callbacks{signal: signal}
	config := cache.NewSnapshotCache(mode == resource.Ads, Hasher{}, logger{})
	srv := server.NewServer(config, cb)
	//als := &main.AccessLogService{}

	// create a test snapshot
    //snapshots := resource.TestSnapshot{
    //Xds:              mode,
    //UpstreamPort:     uint32(upstreamPort),
    //BasePort:         uint32(basePort),
    //NumClusters:      clusters,
    //NumHTTPListeners: httpListeners,
    //NumTCPListeners:  tcpListeners,
    //}

	//生成快照（缓存）,多个版本
    //for i := 0; i < updates; i++ {
        //snapshots.Version = fmt.Sprintf("v%d", i)
        //log.WithFields(log.Fields{"version": snapshots.Version}).Info("update snapshot")

        //snapshot := snapshots.Generate()
        //if err := snapshot.Consistent(); err != nil {
            //log.Errorf("snapshot inconsistency: %+v", snapshot)
        //}

        //err := config.SetSnapshot(nodeID, snapshot)
        //if err != nil {
            //log.Errorf("snapshot error %q for %+v", err, snapshot)
            //os.Exit(1)
        //}
    //}

	// start the xDS server
	//go RunAccessLogServer(ctx, als, alsPort)
	go RunManagementServer(ctx, srv, port)
	//go test.RunManagementServer(ctx, srv, port)
	go RunManagementGateway(ctx, srv, gatewayPort)
	//go test.RunManagementGateway(ctx, srv, gatewayPort)


    consul_api, err := consul.Init()
    if err != nil {
        log.Errorf(err.Error())
    }
	//定时产生配置
    FetchDataFromConsulServer := func(config cache.SnapshotCache, nodeID string) {
        //snapshots := resource.TestSnapshot{
            //Xds:              mode,
            //UpstreamPort:     12345,
            //BasePort:         9876,
            //NumClusters:      4,
            //NumHTTPListeners: 3,
            //NumTCPListeners:  3,
        //}

        for {
            var clusters []cache.Resource
            var endpoints []cache.Resource

            version := fmt.Sprintf("version:%d", time.Now().Minute())
            //拉取roomsvr信息
            services, err := consul_api.GetService("roomsvr")
            if err != nil {
                log.Errorf("consul_api:GetService(%s) failed, %s", "roomsvr", err)
            } else{
                for _, x := range services {
                    log.Infof("roomsvr==> version:%s {name:%s addr:%s tags:%s port:%d}", version, x.Name, x.Address, x.Tags[0], x.Port)
                    clusters = append(clusters, resource.MakeCluster(resource.Xds, x.Name))
                    endpoints = append(endpoints, resource.MakeEndpoint(x.Name, x.Address, uint32(x.Port)))
                }
            }

            //拉取zonesvr信息
            services, err = consul_api.GetService("zonesvr")
            if err != nil {
                log.Errorf("consul_api:GetService(%s) failed, %s", "zonesvr", err)
            } else{
                for _, x := range services {
                    log.Infof("zonesvr==> version:%s {name:%s addr:%s tags:%s port:%d}", 
                        version, x.Name, x.Address, x.Tags[0], x.Port)

                    clusters = append(clusters, resource.MakeCluster(resource.Xds, x.Name))
                    endpoints = append(endpoints, resource.MakeEndpoint(x.Name, x.Address, uint32(x.Port)))
                }
            }


            //构造一个listener
            listeners := make([]cache.Resource, 1)
			listeners[0] = resource.MakeHTTPListener(resource.Xds, "for-front-proxy", 80, "local_route")

            snapshot := cache.NewSnapshot(version, endpoints, clusters, nil, listeners)
            err = config.SetSnapshot(nodeID, snapshot)
            if err != nil {
                log.Errorf("snapshot error %q for %+v", err, snapshot)
            }

            time.Sleep(time.Second * 15)
        }

        //for i := 0; i < updates; i++ {
            //snapshots.Version = fmt.Sprintf("v%d", i)
            //log.WithFields(log.Fields{"version": snapshots.Version}).Info("update snapshot")

            //snapshot := snapshots.Generate()
            //if err := snapshot.Consistent(); err != nil {
                //log.Errorf("snapshot inconsistency: %+v", snapshot)
            //}

            //err := config.SetSnapshot(nodeID, snapshot)
            //if err != nil {
                //log.Errorf("snapshot error %q for %+v", err, snapshot)
                //os.Exit(1)
            //}
            //time.Sleep(10 * time.Second)
        //}
        //log.Infof("FetchDataFromConsulServer end...")
    }
    go FetchDataFromConsulServer(config, nodeID)

	log.Infof("waiting for the first request...")
	<-signal
	//// pass is true if all requests succeed at least once in a run
	//pass := false
	//for j := 0; j < requests; j++ {
	//ok, failed := callEcho()
	//if failed == 0 && !pass {
	//pass = true
	//}
	//log.WithFields(log.Fields{"batch": j, "ok": ok, "failed": failed, "pass": pass}).Info("request batch")
	//select {
	//case <-time.After(delay):
	//case <-ctx.Done():
	//return
	//}
	//}

	//als.Dump(func(s string) { log.Debug(s) })
	//cb.Report()

	//if !pass {
	//log.Errorf("failed all requests in a run %d", i)
	//os.Exit(1)
	//}
	//}

	log.Infof("main loop exit, mode:%s", mode)
}

// callEcho calls upstream echo service on all listener ports and returns an error
// if any of the listeners returned an error.
func callEcho() (int, int) {
	total := httpListeners + tcpListeners
	ok, failed := 0, 0
	ch := make(chan error, total)

	// spawn requests
	for i := 0; i < total; i++ {
		go func(i int) {
			client := http.Client{
				Timeout: 100 * time.Millisecond,
			}
			req, err := client.Get(fmt.Sprintf("http://localhost:%d", basePort+uint(i)))
			if err != nil {
				ch <- err
				return
			}
			defer req.Body.Close()
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				ch <- err
				return
			}
			if string(body) != Hello {
				ch <- fmt.Errorf("unexpected return %q", string(body))
				return
			}
			ch <- nil
		}(i)
	}

	for {
		out := <-ch
		if out == nil {
			ok++
		} else {
			failed++
		}
		if ok+failed == total {
			return ok, failed
		}
	}
}

type logger struct{}

func (logger logger) Infof(format string, args ...interface{}) {
	log.Debugf(format, args...)
}
func (logger logger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

type callbacks struct {
	signal   chan struct{}
	fetches  int
	requests int
	mu       sync.Mutex
}

func (cb *callbacks) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	log.WithFields(log.Fields{"fetches": cb.fetches, "requests": cb.requests}).Info("server callbacks")
}
func (cb *callbacks) OnStreamOpen(id int64, typ string) {
	log.Debugf("stream id:%d open for typ:%s", id, typ)
}
func (cb *callbacks) OnStreamClosed(id int64) {
	log.Debugf("stream id:%d closed", id)
}
func (cb *callbacks) OnStreamRequest(id int64, req *v2.DiscoveryRequest) {
	log.Debugf("OnStreamRequest enter, id:%d, req:%s", id, proto.MarshalTextString(req))
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.requests++

	//if cb.signal != nil {
	//log.Debugf("cb.signal close")
	//close(cb.signal)
	//cb.signal = nil
	//}
}
func (cb *callbacks) OnStreamResponse(id int64, req *v2.DiscoveryRequest, rsp *v2.DiscoveryResponse) {
	log.Debugf("OnStreamResponse enter, id:%d, ==> req:%s <=== rsp:%s ==>", id, proto.MarshalTextString(req), proto.MarshalTextString(rsp))
}
func (cb *callbacks) OnFetchRequest(req *v2.DiscoveryRequest) {
	log.Debugf("OnFetchRequest:%s", proto.MarshalTextString(req))
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.fetches++
	//if cb.signal != nil {
	//log.Debugf("cb.signal close")
	//close(cb.signal)
	//cb.signal = nil
	//}
}
func (cb *callbacks) OnFetchResponse(*v2.DiscoveryRequest, *v2.DiscoveryResponse) {
	log.Debugf("OnFetchResponse")
}
