package user

import (
	"app_server/model"
	"app_server/pkg/cfg"
	"app_server/pkg/db"
	"app_server/pkg/fn"
	"app_server/pkg/httpc"
	"app_server/pkg/jwt"
	"app_server/proto/user"
	"app_server/service/auth"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	connect "connectrpc.com/connect"
)

type UserService struct{}

func (s *UserService) WxUserLogin(ctx context.Context, connectReq *connect.Request[user.WxUserLoginRequest]) (*connect.Response[user.WxUserLoginResponse], error) {
	appName := connectReq.Msg.App
	appId, appSecret := cfg.Viper().GetString("wechat.mp."+appName+".app_id"), cfg.Viper().GetString("wechat.mp."+appName+".app_secret")

	if appId == "" || appSecret == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("app %s not found", appName))
	}

	// 调用微信code2session接口
	url := "https://api.weixin.qq.com/sns/jscode2session"

	// 创建HTTP客户端
	httpResp, err := httpc.Client().R().
		SetQueryParams(map[string]string{
			"appid":      appId,
			"secret":     appSecret,
			"js_code":    connectReq.Msg.Code,
			"grant_type": "authorization_code",
		}).
		Get(url)

	if err != nil {
		slog.ErrorContext(ctx, "call wechat api failed", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("call wechat api failed: %v", err))
	}

	var wxResp WxLoginResponse
	if err := json.Unmarshal(httpResp.Bytes(), &wxResp); err != nil {
		slog.ErrorContext(ctx, "unmarshal wechat response failed", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unmarshal wechat response failed: %v", err))
	}

	if wxResp.ErrCode != 0 {
		slog.ErrorContext(ctx, "wechat error", "error", fmt.Errorf("wechat error: %d, %s", wxResp.ErrCode, wxResp.ErrMsg))
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("wechat error: %d, %s", wxResp.ErrCode, wxResp.ErrMsg))
	}

	// 查询或创建用户
	u, err := WxLoginFirstOrCreate(ctx, wxResp.OpenID, wxResp.UnionID)
	if err != nil {
		slog.ErrorContext(ctx, "create user failed", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("create user failed: %v", err))
	}

	// 生成token
	token, err := jwt.Get().GenerateToken(fn.Itoa(u.ID), time.Hour*24*365)
	if err != nil {
		slog.ErrorContext(ctx, "generate token failed", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("generate token failed: %v", err))
	}

	return connect.NewResponse(&user.WxUserLoginResponse{
		Token: token,
	}), nil
}

func (s *UserService) PhoneLogin(ctx context.Context, connectReq *connect.Request[user.PhoneLoginRequest]) (*connect.Response[user.PhoneLoginResponse], error) {
	phone := connectReq.Msg.Phone
	verificationCode := connectReq.Msg.VerificationCode

	// 参数验证
	if phone == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("手机号不能为空"))
	}
	if verificationCode == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("验证码不能为空"))
	}

	// 验证手机号格式（简单验证）
	if len(phone) != 11 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("手机号格式不正确"))
	}

	// 验证验证码（暂时使用魔法验证码"1234"）
	if verificationCode != "1234" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("验证码错误"))
	}

	// 查询或创建用户
	u, err := PhoneLoginFirstOrCreate(ctx, phone)
	if err != nil {
		slog.ErrorContext(ctx, "phone login create user failed", "error", err, "phone", phone)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("登录失败: %v", err))
	}

	// 生成token
	token, err := jwt.Get().GenerateToken(fn.Itoa(u.ID), time.Hour*24*365)
	if err != nil {
		slog.ErrorContext(ctx, "generate token failed", "error", err, "phone", phone)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("生成token失败: %v", err))
	}

	slog.InfoContext(ctx, "phone login success", "phone", phone, "user_id", u.ID)

	return connect.NewResponse(&user.PhoneLoginResponse{
		Token: token,
	}), nil
}

func (s *UserService) GetUserProfile(ctx context.Context, connectReq *connect.Request[user.GetUserProfileRequest]) (*connect.Response[user.GetUserProfileResponse], error) {
	// 从上下文获取用户ID
	userID, err := auth.ParseUserID(connectReq.Header().Get("Authorization"))
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("用户未登录"))
	}

	// 查询用户信息
	var userModel model.User
	if err := db.GetDB().First(&userModel, userID).Error; err != nil {
		slog.ErrorContext(ctx, "查询用户信息失败", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("查询用户信息失败: %v", err))
	}

	// 构建响应
	response := &user.GetUserProfileResponse{
		User: userModel.ToProto(),
	}

	// 如果有ProfileID，查询用户档案
	if userModel.ProfileID > 0 {
		var profileModel model.Profile
		if err := db.GetDB().First(&profileModel, userModel.ProfileID).Error; err != nil {
			slog.WarnContext(ctx, "查询用户档案失败", "error", err)
		} else {
			response.Profile = profileModel.ToProto()
		}
	}

	return connect.NewResponse(response), nil
}

// WxLoginResponse 微信登录返回结构
type WxLoginResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}
