package user

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go-template-service/biz/service"
	"go-template-service/biz/utils"
	user "go-template-service/hertz_gen/basic/user"
)

// GetUserDetail .
// @router /users/currentUser [GET]
func GetUserDetail(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.UserReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewGetUserDetailService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
