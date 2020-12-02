package handles

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
)

//全局变量 grpc 连接池map
var GrpcPool =make(map[string]string, 0)


func CheckHealthy(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("ok"))
}


func ApiChangesPool(c echo.Context) error {



	return c.String(http.StatusOK, "v2.3 ....")
}

func ChangesPool(response http.ResponseWriter, request *http.Request) {


	body, _ := ioutil.ReadAll(request.Body)
	fmt.Println(body)

	//fmap:=make(map[string]map[string]string,0)


	response.Write([]byte("ok,test"))
}
