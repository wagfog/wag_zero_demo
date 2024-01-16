package svc

import (
	"github.com/wagfog/wag_zero_demo/application/applet/internal/config"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/user"
	"github.com/wagfog/wag_zero_demo/pkg/interceptor"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	UserRPC  user.User
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {

	userRpc := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptor.ClientErrorInterceptor()))

	return &ServiceContext{
		Config:   c,
		UserRPC:  user.NewUser(userRpc),
		BizRedis: redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
	}
}
