package Server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chaos-star/marvel/Etcd"
	"github.com/chaos-star/marvel/Utils"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
	"net"
)

type ServiceOption struct {
	Sd   *grpc.ServiceDesc
	Ss   interface{}
	Name string
}

type RpcServer struct {
	*grpc.Server
	etcd     *Etcd.Engine
	services []ServiceOption
	prefix   string
	port     int
}

type ConnBody struct {
	Op       int
	Addr     string
	Metadata interface{}
}

func Initialize(etcd *Etcd.Engine, port int, prefix string) *RpcServer {
	return &RpcServer{
		Server: grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				grpc_middleware.ChainUnaryServer(
					grpc_ctxtags.UnaryServerInterceptor(),
					grpc_opentracing.UnaryServerInterceptor(),
					grpc_recovery.UnaryServerInterceptor(),
				),
			),
		),
		etcd:     etcd,
		services: []ServiceOption{},
		prefix:   prefix,
		port:     port,
	}
}

func (s *RpcServer) Register(sd *grpc.ServiceDesc, ss interface{}) {
	var service ServiceOption
	service.Sd = sd
	service.Ss = ss
	s.services = append(s.services, service)
}

func (s *RpcServer) Load(services []ServiceOption) {
	s.services = append(s.services, services...)
}

func (s *RpcServer) Run() error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	if len(s.services) <= 0 {
		return errors.New("no service available")
	}
	ip := Utils.GetLocalIP()

	for _, item := range s.services {
		s.RegisterService(item.Sd, item.Ss)
		var conn ConnBody
		conn.Addr = fmt.Sprintf("%s:%d", ip, s.port)
		srvKey := fmt.Sprintf("%s%s/%s", s.prefix, item.Name, conn.Addr)
		body, _ := json.Marshal(conn)
		ers, err := s.etcd.RegisterService(srvKey, string(body), 10, true)
		if err != nil {
			return err
		}
		defer ers.Close()
	}

	go func() {
		fmt.Printf("Rpc Start Success! Listen At:%v \n", listen.Addr())
		s.print()
	}()
	err = s.Serve(listen)
	return nil
}

func (s *RpcServer) print() {
	//var content strings.Builder
	srvs := s.GetServiceInfo()
	if len(srvs) > 0 {
		for k, _ := range srvs {
			fmt.Println(k)
		}
	}
}
