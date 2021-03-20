package micro

import (
	"github.com/limitedlee/microservice/common/config"
	"github.com/limitedlee/microservice/common/handles"
	"github.com/limitedlee/microservice/common/nacos"
	"github.com/lsls907/nacos-sdk-go/vo"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
	"strings"
)

type MicService struct {
	GrpcServer *grpc.Server
	Routes     map[string]func(http.ResponseWriter, *http.Request)
}

func (m *MicService) NewServer() {
	m.GrpcServer = grpc.NewServer(grpc.UnaryInterceptor(filter))
}

func (m *MicService) Start(serviceName string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handles.CheckHealthy)

	mux.HandleFunc("/pool/change", handles.ChangesPool)
	for key, route := range m.Routes {
		mux.HandleFunc(key, route) //websocket
	}
	baseUrl, _ := config.Get("BaseUrl")
	items := strings.Split(baseUrl, ":")
	//addr := fmt.Sprintf(":%v", items[len(items)-1])

	if len(items) <= 0 {
		panic("Please define the port，example(:7065)")
	}
	intNum, _ := strconv.Atoi(items[1])
	port := uint64(intNum)

	nacos.RegisterServiceInstance(vo.RegisterInstanceParam{
		Ip:          nacos.GetOutboundIp(),
		Port:        port,
		ServiceName: serviceName,
		Weight:      10,
		ClusterName: "DEFAULT",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		GroupName:   "DEFAULT_GROUP",
	})
	//http.ListenAndServe(addr, grpcHandleFunc(m.GrpcServer, mux))
}

//func grpcHandleFunc(grpcServer *grpc.Server, otherHander http.Handler) http.Handler {
//	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
//			grpcServer.ServeHTTP(w, r)
//		} else {
//			otherHander.ServeHTTP(w, r)
//		}
//	}), &http2.Server{})
//}
////
//func filter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
//	defer func() {
//		if e := recover(); e != nil {
//			debug.PrintStack()
//			err = status.Errorf(codes.Internal, "Panic err: %v", e)
//		}
//	}()
//
//	md, ok := metadata.FromIncomingContext(ctx)
//	if !ok {
//		return nil, grpc.Errorf(codes.Unauthenticated, "valid token required.")
//	}
//	jwtToken, ok := md["authorization"]
//
//	if jwtToken != nil {
//		if !ok {
//			return nil, grpc.Errorf(codes.Unauthenticated, "valid token required.")
//		}
//
//		index := strings.Index(jwtToken[0], " ")
//		count := strings.Count(jwtToken[0], "")
//		token := jwtToken[0][index+1 : count-1]
//
//		_, err := validateToken(token, jw.PublicKey)
//		if err != nil {
//			return nil, grpc.Errorf(codes.Unauthenticated, fmt.Sprintf("valid token required.%v", err))
//		}
//	}
//
//	return handler(ctx, req)
//}
//
//func validateToken(token string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
//	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
//		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
//			log.Printf("Unexpected signing method: %v", t.Header["alg"])
//			return nil, fmt.Errorf("invalid token")
//		}
//		return publicKey, nil
//	})
//	if err == nil && jwtToken.Valid {
//		return jwtToken, nil
//	}
//	return nil, err
//}
