package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	user "go-template-service/hertz_gen/basic/user"
	common "go-template-service/hertz_gen/common"
)

type GetUserDetailService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserDetailService(Context context.Context, RequestContext *app.RequestContext) *GetUserDetailService {
	return &GetUserDetailService{RequestContext: RequestContext, Context: Context}
}

func (h *GetUserDetailService) Run(req *user.UserReq) (resp *user.UserResp, err error) {
	cu, exist := h.RequestContext.Get("currentUser")
	if !exist {
		h.RequestContext.JSON(consts.StatusBadRequest, "No user detail found.")
		return
	}

	ccu := cu.(common.User)
	res := &user.UserResp{
		User:    &ccu,
		Success: true,
	}

	return res, nil
}
