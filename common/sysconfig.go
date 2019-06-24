package common

type SystemConfig struct {
	//服务英文名
	Name 						string
	//服务中文名
	DisplayName 				string
	//版本号
	Version 					string
	//本地地址,用于回调
	LocalAddress 				string
	//Consul服务地址
	ServiceDiscoveryAddress 	string
	//数据库链接字符串
	DatabaseConnectionString 	string
	//缓存链接字符串
	CacheConnectionString		string
}
