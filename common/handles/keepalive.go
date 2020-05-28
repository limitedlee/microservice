package handles

import (
	"net/http"
)

func CheckHealthy(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("ok"))
}
