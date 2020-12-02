package nacos

type ConfigRequest struct {
	Tenant  string `json:"tenant"`  //租户信息，对应 Nacos 的命名空间ID字段
	DataId  string `json:"dataId"`  //配置 ID
	Group   string `json:"group"`   //配置分组
	Content string `json:"content"` //配置内容
	Type    string `json:"type"`    //配置类型
}

type InitConfigRequest struct {
	NamespaceId string `json:"namespaceId"` //租户信息，对应 Nacos 的命名空间ID字段
	ServerName  string `json:"serverName"`  //
	GroupName   string `json:"groupName"`   //
	Ip          string `json:"ip"`          //
	Port        int    `json:"port"`        //
}
