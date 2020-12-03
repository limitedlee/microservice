package handles

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/limitedlee/microservice/common/nacos"

	"io/ioutil"
	"net/http"
)

//全局变量 grpc 连接池map

func CheckHealthy(response http.ResponseWriter, request *http.Request) {
	_, _ = response.Write([]byte("ok"))
}

func ApiChangesPool(c echo.Context) error {
	changeData := make(map[string][]nacos.PoolUrl, 0)
	if err := c.Bind(changeData); err != nil {
		return c.String(200, "false")
	}
	changeGrpcPool(changeData)

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
	changeGrpcPool(changeData)

	_, _ = response.Write([]byte("true"))
}

func changeGrpcPool(changeData map[string][]nacos.PoolUrl) {
	if len(changeData) <= 0 {
		return
	}
	if len(nacos.GrpcPool) <= 0 {
		return
	}
	data := make(map[string][]nacos.PoolUrl, 0)
	for key := range nacos.GrpcPool {
		if len(changeData[key]) > 0 {
			data[key] = changeData[key]
		}
	}
	if len(data) <= 0 {
		return
	}
	nacos.Mutex.Lock() //对共享变量操作之前先加锁
	nacos.GrpcPool = data
	nacos.Mutex.Unlock() //对共享变量操作完毕在解锁，

}
