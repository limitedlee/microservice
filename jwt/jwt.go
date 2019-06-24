package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
)

func JWT() gin.HandlerFunc  {
	return func(context *gin.Context) {

		if context.Request.RequestURI!="/health" {
			token:=context.GetHeader("Authorization")
			if token!="" {
				auths:=strings.Split(token," ")
				if len(auths)>1 {
					protocol:=auths[0]
					jwtToken:=auths[1]
					if strings.ToLower(protocol)=="bearer"  {
						sub,err:=getSubFromToken(jwtToken)
						if err!=nil {
							fmt.Println("token error:",err)
						}else{
							fmt.Println("111111111111:",sub)
						}
					}
				}
			}
		} else {
			context.Next()
		}

	}
}


// getSubFromToken 获取Token的主题（也可以更改获取其他值）
// 参数tokenStr指的是 从客户端传来的待验证Token
// 验证Token过程中，如果Token生成过程中，指定了iat与exp参数值，将会自动根据时间戳进行时间验证
func getSubFromToken(tokenStr string) (claims jwt.Claims,err error) {
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
		return claims,err
	}

	return claims, errors.New("Token无效或者无对应值")
}
