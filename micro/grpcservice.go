package micro

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/dgrijalva/jwt-go"
	jw "github.com/limitedlee/microservice/common/jwt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcService struct {
	mux              *http.ServeMux
	Opts             Options
	Service          *grpc.Server
	ServiceDiscovery func(dsAddress, namespaceId string)
}

func (gs *GrpcService) New() interface{} {
	maxSize := 50 * 1024 * 1024
	opts:=[]grpc.ServerOption{
		grpc.UnaryInterceptor(filter),
		grpc.MaxRecvMsgSize(maxSize),
		grpc.MaxSendMsgSize(maxSize),
	}
	gs.Service = grpc.NewServer(opts...)
	return gs.Service
}

func (gs *GrpcService) Init(groups map[string]*Group) {
	path := fmt.Sprintf("/%s/ping", gs.Opts.name)

	gs.mux = http.NewServeMux()
	gs.mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("ok"))
	})
}

func (gs *GrpcService) Run() {
	address := fmt.Sprintf(":%d", gs.Opts.port)
	http.ListenAndServe(address, grpcHandleFunc(gs.Service, gs.mux))
}

func grpcHandleFunc(grpcServer *grpc.Server, otherHander http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHander.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func filter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "valid token required.")
	}
	jwtToken, ok := md["authorization"]

	if jwtToken != nil {
		if !ok {
			return nil, grpc.Errorf(codes.Unauthenticated, "valid token required.")
		}

		index := strings.Index(jwtToken[0], " ")
		count := strings.Count(jwtToken[0], "")
		token := jwtToken[0][index+1 : count-1]

		_, err := validateToken(token, jw.PublicKey)
		if err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, fmt.Sprintf("valid token required.%v", err))
		}
	}

	return handler(ctx, req)
}

func validateToken(token string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			//log.Printf("Unexpected signing method: %v", t.Header["alg"])
			return nil, fmt.Errorf("invalid token")
		}
		return publicKey, nil
	})
	if err == nil && jwtToken.Valid {
		return jwtToken, nil
	}
	return nil, err
}
