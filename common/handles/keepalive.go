package handles

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/limitedlee/microservice/common/nacos"
	"io/ioutil"
	"net/http"
	"sync"
)

//全局变量 grpc 连接池map
var GrpcPool = make(map[string][]nacos.PoolUrl, 0)

func CheckHealthy(response http.ResponseWriter, request *http.Request) {
	_, _ = response.Write([]byte("ok"))
}

func ApiChangesPool(c echo.Context) error {
	changeData := make(map[string][]nacos.PoolUrl, 0)
	if err := c.Bind(changeData); err != nil {
		return c.String(200, "false")
	}
	ChangeGrpcPool(changeData)

	return c.String(http.StatusOK, "true")
}

func ChangesPool(response http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		_, _ = response.Write([]byte("false"))
	}
	changeData := make(map[string][]nacos.PoolUrl, 0)

	err = json.Unmarshal(body, &changeData)
	if err != nil {
		_, _ = response.Write([]byte("false"))
	}
	ChangeGrpcPool(changeData)

	_, _ = response.Write([]byte("true"))
}

var mutex sync.Mutex //定义一个锁的变量(互斥锁的关键字是Mutex，其是一个结构体，传参一定要传地址，否则就不对了)

func ChangeGrpcPool(changeData map[string][]nacos.PoolUrl) {
	if len(changeData) <= 0 {
		return
	}
	if len(GrpcPool) <= 0 {
		return
	}
	data := make(map[string][]nacos.PoolUrl, 0)
	for key := range GrpcPool {
		if len(changeData[key]) > 0 {
			data[key] = changeData[key]
		}
	}
	if len(data) <= 0 {
		return
	}
	mutex.Lock() //对共享变量操作之前先加锁
	GrpcPool = data
	mutex.Unlock() //对共享变量操作完毕在解锁，

}
