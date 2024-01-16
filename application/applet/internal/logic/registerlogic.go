package logic

import (
	"context"
	"errors"
	"strings"

	"github.com/wagfog/wag_zero_demo/application/applet/internal/code"
	"github.com/wagfog/wag_zero_demo/application/applet/internal/svc"
	"github.com/wagfog/wag_zero_demo/application/applet/internal/types"
	"github.com/wagfog/wag_zero_demo/application/user/rpc/user"
	"github.com/wagfog/wag_zero_demo/pkg/encrypt"
	"github.com/wagfog/wag_zero_demo/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	prefixActivation = "biz#activation#%s"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// todo: add your logic here and delete this line
	//去除空格
	req.Name = strings.TrimSpace(req.Name)
	req.Mobile = strings.TrimSpace(req.Mobile)
	if len(req.Mobile) == 0 {
		return nil, code.RegisterMobileEmpty
	}
	req.Password = strings.TrimSpace(req.Password)
	if len(req.Password) == 0 {
		return nil, code.RegisterPasswdEmpty
	} else {
		req.Password = encrypt.EncPassword(req.Password)
	}
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}
	//检查验证码是否合法
	err = checkVerificationCode(l.svcCtx.BizRedis, req.Mobile, req.VerificationCode)
	if err != nil {
		logx.Errorf("checkVerificationCode error: %v", err)
		return nil, err
	}

	//对电话号码进行加密
	mobile, err := encrypt.EncMoblie(req.Mobile)
	if err != nil {
		logx.Errorf("EncMobile mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	//通过电话号码查找user
	u, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &user.FindByMobileRequest{
		Mobile: mobile,
	})
	if err != nil {
		logx.Errorf("FindByMobile error: %v", err)
		return nil, err
	}
	//是否存在该用户
	if u != nil && u.UserId > 0 {
		return nil, code.MobileHasRegistered
	}
	//发起注册请求
	regRet, err := l.svcCtx.UserRPC.Register(l.ctx, &user.RegisterRequest{
		Username: req.Name,
		Mobile:   mobile,
	})
	if err != nil {
		logx.Errorf("Register error: %v", err)
		return nil, err
	}

	token, err := jwt.BuildToken(jwt.TokenOptions{
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		Fileds: map[string]interface{}{
			"userId": regRet.UserId,
		},
	})
	if err != nil {
		logx.Errorf("Register error: %v", err)
		return nil, err
	}

	_ = delActivationCache(req.Mobile, req.VerificationCode, l.svcCtx.BizRedis)

	return &types.RegisterResponse{
		UserId: regRet.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExipre,
		},
	}, nil
}

func checkVerificationCode(rds *redis.Redis, mobile, code string) error {
	cacheCode, err := gerActivationCache(mobile, rds)
	if err != nil {
		logx.Errorf("Register error: %v", err)
		return err
	}
	if cacheCode == "" {
		return errors.New("verificationCode expired")
	}
	if cacheCode != code {
		return errors.New("verification code failed")
	}

	return nil
}

// cacheCode,err := getAc
