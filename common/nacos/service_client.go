package nacos

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/limitedlee/microservice/common/config"
	"github.com/lsls907/nacos-sdk-go/clients"
	"github.com/lsls907/nacos-sdk-go/common/constant"
	"github.com/lsls907/nacos-sdk-go/vo"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ConfigsUrl = "v1/cs/configs"

	//nacos 全局配置名称
	MasterConfigName ="MasterPool"
)

//Register with default cluster and group
//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
func RegisterServiceInstance(param vo.RegisterInstanceParam) {
	ipAddr, _ := config.Get("nacos-addr")
	port, _ := config.Get("nacos-port")
	namespaceId, _ := config.Get("nacos-namespace-id")
	intPort, _ := strconv.Atoi(port)
	sc := []constant.ServerConfig{
		{
			IpAddr: ipAddr,
			Port:   uint64(intPort),
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         namespaceId,      //namespace id
		TimeoutMs:           5000,             // 请求Nacos服务端的超时时间，默认是10000ms
		NotLoadCacheAtStart: true,             // 在启动的时候不读取缓存在CacheDir的service信息
		LogDir:              "/tmp/nacos/log", // 日志存储路径
		//CacheDir:            "/tmp/nacos/cache",
		RotateTime: "24h",  // 日志轮转周期，比如：30m, 1h, 24h, 默认是24h
		MaxAge:     3,      // 日志最大文件数，默认3
		LogLevel:   "info", // 日志默认级别，值必须是：debug,info,warn,error，默认值是info
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)

	}
	success, _ := client.RegisterInstance(param)
	fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)

	//初始化配置信息
	initConfigs(InitConfigRequest{
		NamespaceId: namespaceId,
		ServerName:  MasterConfigName,
		GroupName:   param.GroupName,
		Ip:          param.Ip,
		Port:        intPort,
	})
}

// Get preferred outbound ip of this machine
func GetOutboundIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}


//init nacos configs
func initConfigs(params InitConfigRequest) {
	configInfo := ConfigRequest{
		Tenant: params.NamespaceId,
		DataId: params.ServerName,
		Group:  params.GroupName,
	}
	poolMap := make(map[string][]string, 0)
	v := getConfigs(configInfo)
	if len(v) > 0 {
		_ = json.Unmarshal([]byte(v), poolMap)
	}
	if len(poolMap) > 0 {
		if len(poolMap[params.ServerName]) > 0 {
			ips := poolMap[params.ServerName]
			ips = append(ips, fmt.Sprintf("%s:%d", params.Ip, params.Port))
			poolMap[params.ServerName] = ips
		} else {
			ips := []string{fmt.Sprintf("%s:%d", params.Ip, params.Port)}
			poolMap[params.ServerName] = ips
		}
	}
	content, err := json.Marshal(poolMap)
	if err != nil {
		fmt.Print("initConfigs fmap to json error:", err)
		return
	}
	configInfo.Content = string(content)

	flag := setConfigs(configInfo)
	if !flag {
		fmt.Printf("InitConfigs,param:%+v,result:%+v \n\n", configInfo, flag)
	}
}

func setConfigs(configInfo ConfigRequest) bool {

	nacosOpenApiUrl, _ := config.Get("nacos-url")
	data := fmt.Sprintf("?dataId=%s&group=%s&tenant=%s&content=%s", configInfo.DataId,
		configInfo.Group, configInfo.Tenant, configInfo.Content)
	var url = nacosOpenApiUrl + ConfigsUrl + data

	val, err := httpPost(url, nil, "")
	if err != nil {
		panic(err)
		return false
	}
	if string(val) == "true" {
		return true
	}
	return false

}

//get nacos configs
func getConfigs(configInfo ConfigRequest) string {
	nacosOpenApiUrl, _ := config.Get("nacos-url")
	data := fmt.Sprintf("?dataId=%s&group=%s&tenant=%s", configInfo.DataId, configInfo.Group, configInfo.Tenant)
	var url = nacosOpenApiUrl + ConfigsUrl + data
	val, err := httpGet(url, "")
	if err == nil {
		//panic(err)
		return ""
	}

	if string(val) == "config data not exist" {
		return ""
	}
	return string(val)
}

//url ：
//dataJson : 数据对象转化成json字符串
func httpPost(url string, dataJson []byte, Headers string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	client := http.Client{Timeout: time.Second * 60, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	req.Header.Set("Content-Type", "application/json")

	if len(Headers) > 0 {
		strList := strings.Split(Headers, "&")
		for i := 0; i < len(strList); i++ {
			v := strList[i]
			valueList := strings.Split(v, "|")
			req.Header.Set(valueList[0], valueList[1])
		}
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
		return []byte(""), err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

//url ：
func httpGet(url string, headers string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	client := http.Client{Timeout: time.Second * 60, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	if len(headers) > 0 {
		strList := strings.Split(headers, "&")
		for i := 0; i < len(strList); i++ {
			v := strList[i]
			valueList := strings.Split(v, "|")
			req.Header.Set(valueList[0], valueList[1])
		}
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
		return []byte(""), err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body, err
}
