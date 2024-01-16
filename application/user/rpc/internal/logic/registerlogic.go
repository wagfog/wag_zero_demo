package logic

import (
	"context"
	"time"

	"github.com/wagfog/wag_zero_demo/application/user/rpc/internal/code"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/internal/model"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/internal/svc"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 期望的是在RPC中也可以通过New方法定义业务错误码，但直接定义返回的话其实是有问题的，因为gRPC识别不了该业务错误码，错误返回给applet-api后最终会被转换成ServerErr：
func (l *RegisterLogic) Register(in *service.RegisterRequest) (*service.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	if len(in.Username) == 0 {
		return nil, &code.RegisterNameEmpty
	}

	res, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username:   in.Username,
		Avatar:     in.Avatar,
		Mobile:     in.Mobile,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	if err != nil {
		logx.Errorf("Register req: %v error: %v", in, err)
		return nil, err
	}
	userId, err := res.LastInsertId()
	if err != nil {
		logx.Errorf("LastInsertId error: %v", err)
		return nil, err
	}

	return &service.RegisterResponse{UserId: userId}, nil
}
