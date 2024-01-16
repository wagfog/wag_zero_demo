package main

import (
	"flag"
	"fmt"

	"github.com/wagfog/wag_zero_demo/application/user/rpc/internal/config"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/internal/server"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/internal/svc"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/service"
	"github.com/wagfog/wag_zero_demo/pkg/interceptor"

	"github.com/zeromicro/go-zero/core/conf"
	cs "github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		service.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == cs.DevMode || c.Mode == cs.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	//在server拦截器中，我们把自定义的错误转换成RPC的错误，是不能直接转的
	// 我们需要把业务自定义的错误存到grpc status的detail中
	s.AddUnaryInterceptors(interceptor.ServerErrorInterceptor())

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
