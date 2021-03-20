package micro

import (
	"github.com/labstack/echo/v4"
	"github.com/limitedlee/microservice/common/nacos"
	"net/http"
)

type MicroService interface {
	New() interface{}
	Init(groups map[string]*Group)
	Run()
}

type Service struct {
	service          MicroService
	opts             Options
	serviceDiscovery func()
	groups           map[string]*Group
	Instance         interface{}
}

type HandleFunc func(c echo.Context) error

//创建服务
func NewService(opts ...Option) *Service {
	baseService := new(Service)

	thisOpts := Options{}

	for _, opt := range opts {
		opt(&thisOpts)
	}

	if thisOpts.name == "" {
		panic("undefined service name")
	}

	if thisOpts.port == 0 {
		panic("undefined port")
	}

	if thisOpts.serviceType == Http {
		var httpService = new(HttpService)
		httpService.Opts = thisOpts
		baseService.service = httpService
	} else if thisOpts.serviceType == GRPC {
		var grpcService = new(GrpcService)
		grpcService.Opts = thisOpts
		baseService.service = grpcService
	} else {
		panic("undefined service type")
	}

	baseService.Instance = baseService.service.New()

	route := Route{
		verb:       GET,
		prefix:     "/ping",
		handleFunc: CheckHealthy,
	}

	group := &Group{}
	group.routes = make([]Route, 0, 999)
	group.routes = append(group.routes, route)

	baseService.groups = make(map[string]*Group)
	baseService.groups["/"] = group
	baseService.opts = thisOpts

	return baseService
}

func (s *Service) Get(prefix string, f HandleFunc) {
	if s.opts.serviceType != Http {
		panic("unsupported grpc defined route")
	}

	route := Route{
		prefix: prefix,
		verb:   GET,
	}
	route.handleFunc = f

	defaultGroupName := "/"
	if len(s.groups[defaultGroupName].routes) == 0 {
		s.groups[defaultGroupName].routes = make([]Route, 0, 999)
	}
	s.groups[defaultGroupName].routes = append(s.groups[defaultGroupName].routes, route)
}

func (s *Service) POST(prefix string, f HandleFunc) {
	if s.opts.serviceType != Http {
		panic("unsupported grpc defined route")
	}

	route := Route{
		prefix: prefix,
		verb:   POST,
	}
	route.handleFunc = f

	defaultGroupName := "/"
	if len(s.groups[defaultGroupName].routes) == 0 {
		s.groups[defaultGroupName].routes = make([]Route, 0, 999)
	}
	s.groups[defaultGroupName].routes = append(s.groups[defaultGroupName].routes, route)
}

func (s *Service) PUT(prefix string, f HandleFunc) {
	if s.opts.serviceType != Http {
		panic("unsupported grpc defined route")
	}

	route := Route{
		prefix: prefix,
		verb:   PUT,
	}
	route.handleFunc = f

	defaultGroupName := "/"
	if len(s.groups[defaultGroupName].routes) == 0 {
		s.groups[defaultGroupName].routes = make([]Route, 0, 999)
	}
	s.groups[defaultGroupName].routes = append(s.groups[defaultGroupName].routes, route)
}

func (s *Service) DELETE(prefix string, f HandleFunc) {
	if s.opts.serviceType != Http {
		panic("unsupported grpc defined route")
	}

	route := Route{
		prefix: prefix,
		verb:   DELETE,
	}
	route.handleFunc = f

	defaultGroupName := "/"
	if len(s.groups[defaultGroupName].routes) == 0 {
		s.groups[defaultGroupName].routes = make([]Route, 0, 999)
	}
	s.groups[defaultGroupName].routes = append(s.groups[defaultGroupName].routes, route)
}

//注册服务
func (bs *Service) RegisterService(dsAddress, namespaceId string) {
	bs.serviceDiscovery = func() {
		nacos.InitDiscovery(dsAddress, namespaceId)
		go nacos.SubServices()
		////TODO:Watch services change and Subscribe
	}

	bs.serviceDiscovery()
}

func (bs *Service) Init() {
	bs.service.Init(bs.groups)
}

func (bs *Service) Run() {
	bs.service.Run()
}

func CheckHealthy(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
