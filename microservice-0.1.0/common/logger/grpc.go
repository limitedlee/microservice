package logger

import (
	"context"
	"fmt"
	"github.com/limitedlee/microservice/common"
	"log"
	"runtime"
	"time"

	"github.com/limitedlee/microservice/common/config"
	"google.golang.org/grpc"
)

var logGrpcUrl string

func init() {
	//获取项目配置中的数据
	var err error
	logGrpcUrl, err = config.Get("LogGrpc")
	if err != nil {
		log.Printf("get config info fail: %v", err)
	}
}
func Error(in ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)

	pc2, _, _, _ := runtime.Caller(1)
	f2 := runtime.FuncForPC(pc2)

	go writeLog(fmt.Sprintf("%v => %v", f.Name(), f2.Name()), "ERROR", in)
}

func Info(in ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)

	pc2, _, _, _ := runtime.Caller(1)
	f2 := runtime.FuncForPC(pc2)

	go writeLog(fmt.Sprintf("%v => %v", f.Name(), f2.Name()), "INFO", in)
}
func Fatal(in ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)

	pc2, _, _, _ := runtime.Caller(1)
	f2 := runtime.FuncForPC(pc2)
	go writeLog(fmt.Sprintf("%v => %v", f.Name(), f2.Name()), "FATAL", in)
}
func Warn(in ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)

	pc2, _, _, _ := runtime.Caller(1)
	f2 := runtime.FuncForPC(pc2)

	go writeLog(fmt.Sprintf("%v => %v", f.Name(), f2.Name()), "WARN", in)
}

func writeLog(logger string, level string, in []interface{}) (r *Reply) {
	for {
		err := func() (err error) {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()

			var conn *grpc.ClientConn
			conn, err = grpc.Dial(logGrpcUrl, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := NewLogClient(conn)

			ctx, cf := context.WithTimeout(context.Background(), time.Second*60)
			defer cf()

			m := LogRequest{}
			m.Logger = logger
			m.Appid = common.PbConfig.Grpc.Appid

			for _, v := range in {
				switch v.(type) {
				case string:
					m.Message += v.(string)
				case runtime.Error:
					m.Exception = v.(runtime.Error).Error()
				case error:
					m.Exception = v.(error).Error()
				default:
					m.Exception = fmt.Sprintf("%T - ", v) + fmt.Sprintf("%v", v)
				}
			}

			switch level {
			case "ERROR":
				r, err = client.Error(ctx, &m)
			case "INFO":
				r, err = client.Info(ctx, &m)
			case "WARN":
				r, err = client.Warn(ctx, &m)
			case "FATAL":
				r, err = client.Fatal(ctx, &m)
			}
			return
		}()

		if err != nil {
			log.Printf("%v\n", err)
			break
			//time.Sleep(time.Second * 3)
		} else {
			break
		}
	}
	return r
}
