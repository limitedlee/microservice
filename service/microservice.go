package service

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/limitedlee/microservice/jwt"
	"github.com/limitedlee/microservice/rsa"
	"log"
	"strconv"
)

type Microservice struct {
	Service struct {
		//服务英文名
		Name string
		//服务中文名
		DisplayName string
		//版本号
		Version string
		Host    string
		Port    int
	}

	//Consul服务地址
	Consul struct {
		Host string
		Port int
	}

	//数据库链接字符串
	Mysql struct {
		Host    string
		Port    int
		User    string
		Pwd     string
		Default string
	}

	//缓存链接字符串
	Cache struct {
		Host    string
		Port    int
		Pwd     string
		Default int
	}

	Rsa struct {
		PublicKey  string
		PrivateKey string
		Issuer     string
	}

	Route *gin.Engine
}

func (m *Microservice) LoadConfig() (service Microservice) {
	_, err := toml.DecodeFile("appsetting.toml", &m)
	if err != nil {
		log.Fatal(err)
	}
	return *m
}

func (m *Microservice) registerService() {
	consulConfig := consulApi.DefaultConfig()
	consulConfig.Address = fmt.Sprintf("%s:%d", m.Consul.Host, m.Consul.Port)
	client, err := consulApi.NewClient(consulConfig)
	if err != nil {
		log.Fatal("consul client error:", err)
	}

	sConfig := m.Service
	registration := new(consulApi.AgentServiceRegistration)
	registration.ID = sConfig.Name + "-" + sConfig.Host + "-" + strconv.Itoa(sConfig.Port)
	registration.Name = sConfig.Name
	registration.Address = m.Service.Host
	registration.Port = m.Service.Port
	registration.Check = &consulApi.AgentServiceCheck{
		DeregisterCriticalServiceAfter: "5s",
		HTTP:                           fmt.Sprintf("http://%s:%d/health", sConfig.Host, sConfig.Port),
		Interval:                       "2s",
		Timeout:                        "1s"}

	registration.Tags = make([]string, 1)
	registration.Tags[0] = sConfig.DisplayName

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error:", err)
	}
}

func (m *Microservice) InitService() (r *gin.Engine) {
	m.LoadConfig()
	m.registerService()

	publickey, _ := rsa.LoadRsaKey(m.Rsa.PublicKey, m.Rsa.PrivateKey)
	jwt.PublicKey = publickey

	m.Route = gin.Default()
	m.Route.Use(jwt.JWT())
	m.Route.GET("/health", func(context *gin.Context) {
		context.String(200, "ok")
	})

	return m.Route
}

func (m *Microservice) Run() {
	_ = m.Route.Run(fmt.Sprintf(":%d", m.Service.Port))
}
