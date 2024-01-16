package logic

import (
	"context"
	"fmt"

	"github.com/wagfog/wag_zero_demo/application/user/rpc/internal/svc"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/service"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *service.FindByIdRequest) (*service.FindByIdrResponse, error) {
	// todo: add your logic here and delete this line
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		fmt.Printf("Find user id %s ERROR:%s", in.UserId, err.Error())
		return nil, err
	}

	return &service.FindByIdrResponse{
		UserId:   u.Id,
		Username: u.Username,
		Mobile:   u.Mobile,
		Avatar:   u.Avatar,
	}, nil

}
