package consul

import (
	"fmt"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/limitedlee/microservice/common"
	"log"
	"strconv"
	"strings"
)

func RegisterService(sysConfig *common.SystemConfig) {
	config := consulApi.DefaultConfig()
	config.Address = sysConfig.ServiceDiscoveryAddress

	client, err := consulApi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error:", err)
	}

	ip := strings.Split(sysConfig.LocalAddress, ":")[0]
	port, _ := strconv.Atoi(strings.Split(sysConfig.LocalAddress, ":")[1])

	registration := new(consulApi.AgentServiceRegistration)
	registration.ID = sysConfig.Name + "-" + strings.ReplaceAll(sysConfig.LocalAddress, ":", "-")
	registration.Name = sysConfig.Name
	registration.Address = ip
	registration.Port = port
	registration.Check = &consulApi.AgentServiceCheck{
		DeregisterCriticalServiceAfter: "5s",
		HTTP:                           fmt.Sprintf("http://%s/health", sysConfig.LocalAddress),
		Interval:                       "2s",
		Timeout:                        "1s"}

	registration.Tags = make([]string, 1)
	registration.Tags[0] = sysConfig.DisplayName

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error:", err)
	}

}
