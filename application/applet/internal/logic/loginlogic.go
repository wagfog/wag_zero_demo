package logic

import (
	"context"
	"strings"

	"github.com/wagfog/wag_zero_demo/application/applet/internal/code"
	"github.com/wagfog/wag_zero_demo/application/applet/internal/svc"
	"github.com/wagfog/wag_zero_demo/application/applet/internal/types"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/service"
	"github.com/wagfog/wag_zero_demo/pkg/encrypt"
	"github.com/wagfog/wag_zero_demo/pkg/jwt"
	"github.com/wagfog/wag_zero_demo/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// todo: add your logic here and delete this line
	req.Mobile = strings.TrimSpace(req.Mobile)
	if len(req.Mobile) == 0 {
		return nil, code.LoginMobileEmpty
	}
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}

	err = checkVerificationCode(l.svcCtx.BizRedis, req.Mobile, req.VerificationCode)
	if err != nil {
		// logx.Errorf("Register error: %v", err)
		return nil, err
	}

	Mobile, err := encrypt.EncMoblie(req.Mobile)
	if err != nil {
		logx.Errorf("EncMobile mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}

	u, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &service.FindByMobileRequest{
		Mobile: Mobile,
	})
	if err != nil {
		logx.Errorf("Find mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	if u == nil || u.UserId == 0 {
		return nil, xcode.AccessDenied
	}

	token, err := jwt.BuildToken(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fileds: map[string]interface{}{
			"userId": u.UserId,
		},
	})
	if err != nil {
		return nil, err
	}

	_ = delActivationCache(req.Mobile, req.VerificationCode, l.svcCtx.BizRedis)

	return &types.LoginResponse{
		UserId: u.UserId,
		Token: types.Token{
			AccessExpire: token.AccessExipre,
			AccessToken:  token.AccessToken,
		},
	}, nil
}
