package golog

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/limitedlee/microservice/golog/config"
	"github.com/limitedlee/microservice/golog/proto"
	"google.golang.org/grpc"
)

var client proto.LogClient

func init() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	conn, err := grpc.Dial(config.App.Grpc.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	client = proto.NewLogClient(conn)
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

	writeLog(fmt.Sprintf("%v => %v", f.Name(), f2.Name()), "ERROR", in)
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

	writeLog(fmt.Sprintf("%v => %v", f.Name(), f2.Name()), "INFO", in)
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
	writeLog(fmt.Sprintf("%v => %v", f.Name(), f2.Name()), "FATAL", in)
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

	writeLog(fmt.Sprintf("%v => %v", f.Name(), f2.Name()), "WARN", in)
}

func writeLog(logger string, level string, in []interface{}) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		m := proto.LogRequest{}
		m.Logger = logger
		m.Appid = config.App.Grpc.Appid

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

		ctx, _ := context.WithTimeout(context.Background(), time.Second*60)

		var r *proto.Reply
		var err2 error
		switch level {
		case "ERROR":
			r, err2 = client.Error(ctx, &m)
		case "INFO":
			r, err2 = client.Info(ctx, &m)
		case "WARN":
			r, err2 = client.Warn(ctx, &m)
		case "FATAL":
			r, err2 = client.Fatal(ctx, &m)
		}
		if err2 != nil {
			log.Fatalf("could not greet: %v", err2)
		}
		log.Println(r)
	}()
}
