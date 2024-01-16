package logic

import (
	"context"
	"fmt"

	"github.com/wagfog/wag_zero_demo/application/user/rpc/internal/svc"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByMobileLogic {
	return &FindByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByMobileLogic) FindByMobile(in *service.FindByMobileRequest) (*service.FindByMobileResponse, error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if user == nil {
		return &service.FindByMobileResponse{}, nil
	}
	if err != nil {
		fmt.Printf("Find user mobile %s ERROR:%s", in.Mobile, err.Error())
		return nil, err
	}

	return &service.FindByMobileResponse{
		UserId:   user.Id,
		Username: user.Username,
		Mobile:   user.Mobile,
		Avatar:   user.Avatar,
	}, nil
}
