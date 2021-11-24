package jwt

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/limitedlee/microservice/common/config"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
)

func init() {
	PublicKeyString, err := config.Get("PublicKey")
	PrivateKeyString, err := config.Get("PrivateKey")

	if err != nil {
		log.Println(err.Error())
	}

	publicKeyByte, _ := ioutil.ReadFile(PublicKeyString)
	privateKeyByte, _ := ioutil.ReadFile(PrivateKeyString)

	PublicKey, _ = jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)
	PrivateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)
}

// getSubFromToken 获取Token的主题（也可以更改获取其他值）
// 参数tokenStr指的是 从客户端传来的待验证Token
// 验证Token过程中，如果Token生成过程中，指定了iat与exp参数值，将会自动根据时间戳进行时间验证
func getSubFromToken(tokenStr string) (claims jwt.Claims, err error) {
	// 基于公钥验证Token合法性
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 基于JWT的第一部分中的alg字段值进行一次验证
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("验证Token的加密类型错误")
		}
		return PublicKey, nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		//return claims["sub"].(string), nil
		return claims, err
	}

	return claims, errors.New("Token无效或者无对应值")
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
	}

	i := strings.LastIndex(path, "/")
	if i < 0 {
		return "", errors.New(`Can't find "/" or "\".`)
	}

	return string(path[0 : i+1]), nil
}
