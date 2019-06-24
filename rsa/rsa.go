package rsa

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func LoadRsaKey(pubicKeyFileName,privateKeyFileName string) (publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) {
	publicKeyByte, err := ioutil.ReadFile(pubicKeyFileName)
	if err != nil {
		log.Println(err.Error())
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyByte)

	privateKeyByte, err := ioutil.ReadFile(privateKeyFileName)
	if err != nil {
		log.Println(err.Error())
	}
	privateKey, _ = jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)


	path,_:=GetCurrentPath()

	fmt.Println(path)

	return publicKey, privateKey
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
