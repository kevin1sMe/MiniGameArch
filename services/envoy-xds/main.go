package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/gogo/protobuf/proto"
    //"github.com/onrik/logrus/filename"
    resource "envoy-xds/resource"
    consul "envoy-xds/consul"
)

var (
	debug bool

	port         uint
	gatewayPort  uint
	//alsPort      uint

    refreshDelay uint

	mode          string
	nodeID string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Use debug logging")
	flag.UintVar(&port, "port", 18000, "Management server port")
	flag.UintVar(&gatewayPort, "gateway", 18001, "Management server port for HTTP gateway")
	flag.UintVar(&refreshDelay, "refresh", 15, "Refresh services delay(second)")
	//flag.UintVar(&alsPort, "als", 18090, "Accesslog server port")
	flag.StringVar(&mode, "xds", resource.Ads, "Management server type (ads, xds, rest)")
	flag.StringVar(&nodeID, "nodeID", "test-id", "Node ID")
}

func main() {

	flag.Parse()
	if debug {
		log.SetLevel(log.DebugLevel)
	}
    //TODO 给日志加文件名及行号。。
    //log.SetFormatter(&log.TextFormatter{
     //TimestampFormat{})
    //}
    //log.AddHook(filename.NewHook(log.GetLevel()))

	//ctx, cancel := context.WithTimeout(context.Background(), 5000 * time.Second)
	//defer cancel()
	ctx := context.Background()

	// create a cache
	signal := make(chan struct{})
	cb := &callbacks{signal: signal}
	config := cache.NewSnapshotCache(mode == resource.Ads, Hasher{}, logger{})
	srv := server.NewServer(config, cb)
	//als := &main.AccessLogService{}

	// start the xDS server
	//go RunAccessLogServer(ctx, als, alsPort)
	go RunManagementServer(ctx, srv, port)
	go RunManagementGateway(ctx, srv, gatewayPort)

    go FetchDataFromConsulServer(config, refreshDelay, nodeID)

	log.Infof("waiting for the first request...")
	<-signal

	log.Infof("main loop exit, mode:%s", mode)
}

//日志接口
//----------------------------------------------------------------------------------------
type logger struct{}

func (logger logger) Infof(format string, args ...interface{}) {
    log.Debugf(format, args...)
}
func (logger logger) Errorf(format string, args ...interface{}) {
    log.Errorf(format, args...)
}

//这里需要实现 go-control-plane/pkg/server中定义的Callbacks, 并且定义了一些扩展字段
//----------------------------------------------------------------------------------------
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
}

func (cb *callbacks) OnStreamResponse(id int64, req *v2.DiscoveryRequest, rsp *v2.DiscoveryResponse) {
	log.Debugf("OnStreamResponse enter, id:%d, ==> req:%s <=== rsp:%s ==>", id, proto.MarshalTextString(req), proto.MarshalTextString(rsp))
}

func (cb *callbacks) OnFetchRequest(req *v2.DiscoveryRequest) {
	log.Debugf("OnFetchRequest:%s", proto.MarshalTextString(req))
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.fetches++
}

func (cb *callbacks) OnFetchResponse(req *v2.DiscoveryRequest, rsp *v2.DiscoveryResponse) {
    log.Debugf("OnFetchResponse, req:[%s] rsp:[%s]", proto.MarshalTextString(req), proto.MarshalTextString(rsp))
}

//定时产生配置--for test
func FetchDataFromConsulServer(config cache.SnapshotCache, delay uint, nodeID string) {

    //初始化consul api
    consul_api, err := consul.Init()
    if err != nil {
        log.Errorf(err.Error())
    }

    //定时拉取服务
    for {
        var clusters []cache.Resource
        var endpoints []cache.Resource
        var routers []cache.Resource


        all_services, err := consul_api.GetAllServices()
        if err != nil {
            log.Errorf(err.Error())
            time.Sleep( 5* time.Second)
            continue
        }

        version := fmt.Sprintf("version:%d", time.Now().Minute())

        for _, s := range all_services {
            services, err := consul_api.GetService(s.Name)
            if err != nil {
                log.Errorf("consul_api:GetService(%s) failed, %s", s.Name, err)
            } else{
                for _, x := range services {
                    clusters = append(clusters, resource.MakeCluster(resource.Xds, x.Name))
                    endpoints = append(endpoints, resource.MakeEndpoint(x.Name, x.Address, uint32(x.Port)))
                }

                routers = append(routers, resource.MakeRoute("local_route", s.Name))
            }
        }
            //拉取zonesvr信息
            //services, err = consul_api.GetService("zonesvr")
            //if err != nil {
            //log.Errorf("consul_api:GetService(%s) failed, %s", "zonesvr", err)
            //} else{
            //for _, x := range services {
            //log.Infof("zonesvr==> version:%s {name:%s addr:%s tags:%s port:%d}", 
            //version, x.Name, x.Address, x.Tags[0], x.Port)

            //clusters = append(clusters, resource.MakeCluster(resource.Xds, x.Name))
            //endpoints = append(endpoints, resource.MakeEndpoint(x.Name, x.Address, uint32(x.Port)))
            //}
            //}


            //构造一个listener
            listeners := make([]cache.Resource, 1)
            listeners[0] = resource.MakeHTTPListener(resource.Xds, "for-front-proxy", 80, "local_route")

            snapshot := cache.NewSnapshot(version, endpoints, clusters, routers, listeners)
            err = config.SetSnapshot(nodeID, snapshot)
            if err != nil {
                log.Errorf("snapshot error %q for %+v", err, snapshot)
            }


            time.Sleep(time.Duration(delay) * time.Second)
        }
}
