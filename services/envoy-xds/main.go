package main

import (
    "context"
    "flag"
    "fmt"
    log "github.com/sirupsen/logrus"
    "strings"
    "sync"
    "time"

    "envoy-xds/consul"
    "github.com/envoyproxy/go-control-plane/envoy/api/v2"
    "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
    "github.com/envoyproxy/go-control-plane/pkg/cache"
    "github.com/envoyproxy/go-control-plane/pkg/server"
    "github.com/gogo/protobuf/proto"
    // "github.com/onrik/logrus/filename"
    resource "envoy-xds/resource"
)

var (
	debug bool

	port         uint
	gatewayPort  uint
	//alsPort      uint

    refreshDelay uint

	mode          string
	// nodeID string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Use debug logging")
	flag.UintVar(&port, "port", 18000, "Management server port")
	flag.UintVar(&gatewayPort, "gateway", 18001, "Management server port for HTTP gateway")
	flag.UintVar(&refreshDelay, "refresh", 15, "Refresh services delay(second)")
	//flag.UintVar(&alsPort, "als", 18090, "Accesslog server port")
	flag.StringVar(&mode, "xds", resource.Ads, "Management server type (ads, xds, rest)")
	// flag.StringVar(&nodeID, "nodeID", "test-id", "Node ID")
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
	cb := &callbacks{signal: signal, nodeList: make(map[string]string, 128) }
	config := cache.NewSnapshotCache(mode == resource.Ads, Hasher{}, logger{})
	srv := server.NewServer(config, cb)
	//als := &main.AccessLogService{}

	// start the xDS server
	//go RunAccessLogServer(ctx, als, alsPort)
	go RunManagementServer(ctx, srv, port)
	go RunManagementGateway(ctx, srv, gatewayPort)

    go FetchDataFromConsulServer(config, refreshDelay, cb)

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

	nodeList map[string]string
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
	//记录版本
    cb.nodeList[req.Node.GetCluster()] = req.VersionInfo
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
func FetchDataFromConsulServer(config cache.SnapshotCache, delay uint, cb *callbacks) {
    //初始化consul api
    consul_api, err := consul.Init()
    if err != nil {
        log.Errorf(err.Error())
    }

    //定时拉取服务
    for {


        all_services, err := consul_api.GetAllServices()
        if err != nil {
            log.Errorf(err.Error())
            time.Sleep( 5* time.Second)
            continue
        }

        version := fmt.Sprintf("version:%d", time.Now().Minute())

        //遍历所有的node, 为他们创建各有不同的配置
        for nodeName, nodeVersion := range cb.nodeList {
            log.Infof("nodeName:%s nodeVersion:%s", nodeName, nodeVersion)
            var clusters []cache.Resource
            var endpoints []cache.Resource
            var routers []cache.Resource
            //var listeners []cache.Resource

            var routeRules []route.Route

            for _, s := range all_services {
                services, err := consul_api.GetService(s.Name)
                if err != nil {
                    log.Errorf("consul_api:GetService(%s) failed, %s", s.Name, err)
                } else{
                    for _, x := range services {
                        if x.Name == "local_zonesvr" {
                            clusters = append(clusters, resource.MakeClusterHTTP1(resource.Xds, x.Name))//, v22.HTTP1))
                        }else {
                            clusters = append(clusters, resource.MakeCluster(resource.Xds, x.Name))//, v22.HTTP2))
                        }
                        endpoints = append(endpoints, resource.MakeEndpoint(x.Name, x.Address, uint32(x.Port)))
                    }

                    log.Infof("routeRules, s.Name:%s", s.Name)
                    //TODO 这些逻辑抽出去封装起来。不同的node需要不同的规则和返回, 不同的service也有不一样的路由规则
                    if !strings.Contains(s.Name, "local_") {
                        if s.Name == nodeName {
                            //替换掉所有原本到自己的route rules, 比如roomsvr上的envoy不用负责转发到roomsvr，而应该是到local_roomsvr
                            log.Infof("routeRules, s.Name:%s , nodeName:%s add to LocalServiceRouteRules", s.Name, nodeName)
                            routeRules = append(routeRules, BuildLocalServiceRouteRules(s.Name)...)
                        } else {
                                log.Infof("routeRules, s.Name:%s , add to ServiceRouteRules", s.Name)
                                routeRules = append(routeRules, BuildServiceRouteRules(s.Name)...)
                        }
                    }
                }
            }

            //为了HTTP的测试协议适配的转发通道
            if nodeName == "front-proxy" {
                routeRules = append(routeRules, BuildZoneServiceRouteRules("zonesvr")...)
            } else if nodeName == "zonesvr" {
                routeRules = append(routeRules, BuildZoneServiceRouteRules("local_zonesvr")...)
            }

            //----------------------
            route_cfg_name := "envoy." + nodeName + ".minigame"
            tmpRouters := append(routers, resource.MakeRoute(route_cfg_name, routeRules))

            //构造一个listener
            listeners := make([]cache.Resource, 1)
            listeners[0] = resource.MakeHTTPListener(resource.Xds, nodeName, 80, route_cfg_name)

            log.Infof("NewSnapshot(node[%s] ==> version:%s, endpoint sz:%d, clusters sz:%d, routers sz:%d listeners sz:%d)",
                nodeName, version, len(endpoints), len(clusters), len(tmpRouters), len(listeners))
            snapshot := cache.NewSnapshot(version, endpoints, clusters, tmpRouters, listeners)
            err = config.SetSnapshot(nodeName, snapshot)
            if err != nil {
                log.Errorf("snapshot error %q for %+v", err, snapshot)
            }
        }

        time.Sleep(time.Duration(delay) * time.Second)
    }
}

func BuildZoneServiceRouteRules(serviceName string) []route.Route{
    var routeRules []route.Route
    var prefix string
    var decorator string
    var header_matcher []*route.HeaderMatcher

    //对于/zone的路由，是现阶段的临时方案
    prefix = "/zone"
    decorator = "checkAvailability"
    routeRules = append(routeRules, resource.MakeRouteRule(serviceName, prefix, decorator, header_matcher))

    return routeRules
}

func BuildServiceRouteRules(serviceName string) []route.Route{
    var routeRules []route.Route
    var prefix string
    var decorator string
    var header_matcher []*route.HeaderMatcher
    decorator = "checkAvailability"

    prefix = "/"
    header_matcher = append(header_matcher, &route.HeaderMatcher{
        Name: "content-type",
        HeaderMatchSpecifier: &route.HeaderMatcher_ExactMatch{
            ExactMatch: "application/grpc",
        },
    },)
    header_matcher = append(header_matcher, &route.HeaderMatcher{
        Name: "server-family",
        HeaderMatchSpecifier: &route.HeaderMatcher_ExactMatch{
            ExactMatch: serviceName,
        },
    },)
    routeRules = append(routeRules, resource.MakeRouteRule(serviceName, prefix, decorator, header_matcher))
    return routeRules
}

func BuildLocalServiceRouteRules(serviceName string) []route.Route{
    var routeRules []route.Route
    var prefix string
    var decorator string
    var header_matcher []*route.HeaderMatcher
    decorator = "checkAvailability"

    prefix = "/"
    header_matcher = append(header_matcher, &route.HeaderMatcher{
        Name: "content-type",
        HeaderMatchSpecifier: &route.HeaderMatcher_ExactMatch{
            ExactMatch: "application/grpc",
        },
    },)
    header_matcher = append(header_matcher, &route.HeaderMatcher{
        Name: "server-family",
        HeaderMatchSpecifier: &route.HeaderMatcher_ExactMatch{
            ExactMatch: serviceName,
        },
    },)

    localServiceName := "local_" + serviceName
    routeRules = append(routeRules, resource.MakeRouteRule(localServiceName, prefix, decorator, header_matcher))
    return routeRules
}
