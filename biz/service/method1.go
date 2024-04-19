package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-jwt/jwt/v5"

	"github.com/cloudwego/hertz/pkg/app"
	hello "go-template-service/hertz_gen/cwgo/http/hello"
)

type Method1Service struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewMethod1Service(Context context.Context, RequestContext *app.RequestContext) *Method1Service {
	return &Method1Service{RequestContext: RequestContext, Context: Context}
}

func (h *Method1Service) Run(req *hello.HelloReq) (resp *hello.HelloResp, err error) {
	// get userDetail from token
	userDetail, exist := h.RequestContext.Get("currentUser")
	if !exist {
		h.RequestContext.JSON(consts.StatusBadRequest, "No user detail found.")
		return
	}

	// return userDetail name to client
	claim := userDetail.(jwt.MapClaims)
	sub, _ := claim["name"].(string)
	res := &hello.HelloResp{
		RespBody: "Hello " + sub,
	}
	return res, nil
}
