package consul

import (
    "fmt"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Consul struct {
	client *api.Client
}

type ConsulService struct {
	Name    string
	Address string
    Tags    []string
	Port    int
}

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
		serviceAddressesPorts = append(serviceAddressesPorts, ConsulService{
			Name:    addr.Service.Service,
			Address: addr.Node.Address,
            Tags : addr.Service.Tags, 
			Port:    addr.Service.Port})
            log.Infof("GetService(%s), {name:%s tags:%s addres:%s port:%d}",
                serviceName,  addr.Service.Service, addr.Service.Tags[0],  addr.Node.Address, addr.Service.Port)
	}

	return serviceAddressesPorts, nil
}


func Init() (*Consul, error) {
    config := api.DefaultConfig()
    client, err := api.NewClient(config)
    if err != nil {
        return nil, fmt.Errorf("consul client init failed %s", err)
    }

    return &Consul{client: client}, nil
}

