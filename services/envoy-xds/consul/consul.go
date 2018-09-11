package consul

import (
    "fmt"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
    "encoding/json"
	log "github.com/sirupsen/logrus"
)

type Consul struct {
	client *api.Client
}

type ConsulService struct {
	Name    string `json:"name"`
	Address string `json:"addr"`
    Tags    []string `json:"tags"`
	Port    int  `json:"port"`
}

func Init() (*Consul, error) {
    config := api.DefaultConfig()
    client, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("consul client init failed %s", err)
    }

    return &Consul{client: client}, nil
}


//获得单个service的详细状态
func (c *Consul) GetService(serviceName string) ([]ConsulService, error) {
    log.Infof("GetService(%s)", serviceName)
	serviceAddressesPorts := []ConsulService{}
	// get consul service addresses and ports
	addresses, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get consul service")
	}

    log.Infof("GetService(%s), rsp address size:%d", serviceName, len(addresses))

	for _, addr := range addresses {
		// append service addresses and ports
        srv := ConsulService{
			Name:    addr.Service.Service,
			Address: addr.Node.Address,
            Tags : addr.Service.Tags,
			Port:    addr.Service.Port}

        str, err := json.Marshal(&srv)
        if err != nil {
		    return nil, errors.Wrap(err, "json marshal failed")
        }
        log.Infof("GetService(%s), %s", serviceName, str)
		serviceAddressesPorts = append(serviceAddressesPorts, srv)
	}

	return serviceAddressesPorts, nil
}


//获得所有注册的service
func (c *Consul) GetAllServices() ([]ConsulService, error) {
    services,  _, err := c.client.Catalog().Services(nil)
    if err != nil {
        log.Errorf("GetAllServices() failed, %s", err)
        return nil, errors.New("catalog get Services failed")
    }

    log.Infof("GetAllServices(), rsp size:%d", len(services))

	var srv []ConsulService

    for s, values := range services {
        if s == "consul" {
            log.Infof("skip services[consul]")
            continue
        }
        for  i, v := range values {
            log.Infof("services: k[%s] ==> v[%d] = %s", s, i, v)
        }
        srv = append(srv, ConsulService { Name: s })
        //serviceAddressesPorts = append(serviceAddressesPorts, ConsulService{
			//Name:    addr.Service.Service,
			//Address: addr.Node.Address,
            //Tags : addr.Service.Tags, 
			//Port:    addr.Service.Port})
            //log.Infof("GetService(%s), {name:%s tags:%s addres:%s port:%d}",
                //serviceName,  addr.Service.Service, addr.Service.Tags[0],  addr.Node.Address, addr.Service.Port)

    }

	return srv, nil
}
