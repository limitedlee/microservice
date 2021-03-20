package micro

import (
	"github.com/limitedlee/microservice/common/net"
)

type Options struct {
	serviceType ServiceType
	name        string //服务名称
	port        int
	runMode     RunMode
}

type Option func(opts *Options)

func SetServiceType(serviceType ServiceType) Option {
	return func(opts *Options) {
		opts.serviceType = serviceType
	}
}

//设置服务名称
func SetServiceName(name string) Option {
	return func(opts *Options) {
		opts.name = name
	}
}

//设置允许模式
func SetRunMode(mode RunMode) Option {
	return func(opts *Options) {
		opts.runMode = mode

		if opts.runMode == ByHost || opts.runMode == ByDocker {
			opts.port = net.GetAvailablePort()
		} else if opts.runMode == ByK8s {
			opts.port = 8899
		} else {
			panic("undefined run mode")
		}
	}
}
