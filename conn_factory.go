package go_pool

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"io"
	"math/rand"
	"time"
)

func GetHostsFuncFromConsul(consulClient *api.Client, ServiceId string) GetHostsFunc {
	return func() (hosts []string, err error) {
		services, _, err := consulClient.Health().Service(ServiceId, "", true, nil)
		if err != nil {
			return
		}
		for _, service := range services {
			host := service.Service.Address
			if host == "" {
				host = service.Node.Address
			}
			host = fmt.Sprintf("%s:%d", host, service.Service.Port)
			hosts = append(hosts, host)
		}

		return
	}
}

type GetHostsFunc func() (hosts []string, err error)

func ConnFactoryGrpc(getHostsFunc GetHostsFunc) ConnFactory {
	return func() (closer io.Closer, err error) {

		hosts, err := getHostsFunc()
		if err != nil {
			return
		}

		hostsCount := len(hosts)
		if hostsCount == 0 {
			err = ErrNoAliveHost
			return
		}

		rand.Seed(time.Now().Unix())
		index := rand.Intn(hostsCount)
		host := hosts[index]
		closer, err = grpc.Dial(host, grpc.WithInsecure())

		return
	}
}
