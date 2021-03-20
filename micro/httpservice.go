package micro

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type HttpService struct {
	Opts             Options
	Service          *echo.Echo
	ServiceDiscovery func(dsAddress, namespaceId string)
	groups           map[string]*Group
}

func (hs *HttpService) New() interface{} {
	hs.Service = echo.New()
	return hs.Service
}

func (hs *HttpService) Init(groups map[string]*Group) {
	for groupPrefix, group := range groups {
		if groupPrefix != "/" {
			e := hs.Service.Group(fmt.Sprintf("/%s%s", hs.Opts.name, groupPrefix))
			for _, route := range group.routes {
				if route.verb == GET {
					e.GET(route.prefix, echo.HandlerFunc(route.handleFunc))
				}
				if route.verb == PUT {
					e.PUT(route.prefix, echo.HandlerFunc(route.handleFunc))
				}
				if route.verb == POST {
					e.POST(route.prefix, echo.HandlerFunc(route.handleFunc))
				}
				if route.verb == DELETE {
					e.DELETE(route.prefix, echo.HandlerFunc(route.handleFunc))
				}
			}
		} else {
			for _, route := range group.routes {
				path := fmt.Sprintf("/%s%s", hs.Opts.name, route.prefix)
				if route.verb == GET {
					hs.Service.GET(path, echo.HandlerFunc(route.handleFunc))
				}
				if route.verb == PUT {
					hs.Service.PUT(path, echo.HandlerFunc(route.handleFunc))
				}
				if route.verb == POST {
					hs.Service.POST(path, echo.HandlerFunc(route.handleFunc))
				}
				if route.verb == DELETE {
					hs.Service.DELETE(path, echo.HandlerFunc(route.handleFunc))
				}
			}
		}
	}
}

func (hs *HttpService) Run() {
	address := fmt.Sprintf(":%d", hs.Opts.port)
	hs.Service.Start(address)
}
