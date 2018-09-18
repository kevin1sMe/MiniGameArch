# 最小游戏架构  
这个分支使用 consul 作为服务发现与注册框架。envoy 作为服务之间透明的网络代理。grpc 作为服务之间的通信方式。

## consul 特性介绍
Consul 是 HashiCorp 公司推出的开源工具，用于实现分布式系统的服务发现与配置。与其他分布式服务注册与发现的方案，Consul 的方案更“一站式”，内置了服务注册与发现框 架、分布一致性协议实现、健康检查、Key/Value存储、多数据中心方案，不再需要依赖其他工具（比如ZooKeeper等）。使用起来也较 为简单。Consul使用Go语言编写，因此具有天然可移植性(支持Linux、windows和Mac OS X)；安装包仅包含一个可执行文件，方便部署，与 Docker 等轻量级容器可无缝配合 。

## 项目架构
项目中包含一个 zonesvr 和一个 roomsvr。


## 使用  
* 安装好docker及docker-compose等基础套件。
* 使用`docker-compose up --build`来构建并启动服务。
* 访问zonesvr服务。`curl http://localhost:8000/zone/hello`
* 服务追踪：可访问: `http://localhost:16686/`
* 查看consul服务状态，可访问: `http://localhost:8500/ui/dc1/services`
* 发起一个rpc调用: `curl 'localhost:8000/zone/room/1'`
* 查看http接口状态码 `curl -I -m 10 -o /dev/null -s -w %{http_code} url:port`
