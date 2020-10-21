package user

import (
	"github.com/1024casts/fastim/handler"
	"github.com/1024casts/fastim/internal/model"
	"github.com/1024casts/fastim/internal/service"
	"github.com/1024casts/fastim/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
	"github.com/gin-gonic/gin"
)

// PhoneLogin 手机登录接口
// @Summary 用户登录接口
// @Description 仅限手机登录
// @Tags 用户
// @Produce  json
// @Param req body PhoneLoginCredentials true "phone"
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6Ik"}}"
// @Router /users/login [post]
func PhoneLogin(c *gin.Context) {
	log.Info("Phone Login function called.")

	// Binding the data with the u struct.
	var req PhoneLoginCredentials
	if err := c.Bind(&req); err != nil {
		log.Warnf("phone login bind param err: %v", err)
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}

	log.Infof("req %#v", req)
	// check param
	if req.Phone == 0 || req.VerifyCode == 0 {
		log.Warn("phone login bind param is empty")
		handler.SendResponse(c, errno.ErrParam, nil)
		return
	}

	// 验证校验码
	if !service.VCodeService.CheckLoginVCode(req.Phone, req.VerifyCode) {
		handler.SendResponse(c, errno.ErrVerifyCode, nil)
		return
	}

	// 登录
	t, err := service.NewUserService().PhoneLogin(c, req.Phone, req.VerifyCode)
	if err != nil {
		handler.SendResponse(c, errno.ErrVerifyCode, nil)
		return
	}

	handler.SendResponse(c, nil, model.Token{
		Token: t,
	})
}
