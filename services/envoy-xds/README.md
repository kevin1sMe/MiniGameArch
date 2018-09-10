### 实现envoy的动态配置服务

https://www.envoyproxy.io/docs/envoy/latest/api-v2/api/v2/discovery.proto#discoveryresponse

### 编译
`make`

### 基于go-control-plane框架
### 基于consul获取服务列表


### 注意
* 建议使用go1.11版本构建  
* 需要使用modules模块, 你可以这么使用  
    * 在此目录下已经有go.mod了，要下载依赖的包，使用go mod tidy。必须做一次
    * 若需要把下载的包拷贝到此目录，使用 go mod vendor。 此步可选


