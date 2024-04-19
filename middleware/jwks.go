package middleware

import (
	"context"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	u "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-jwt/jwt/v5"
	utils "go-template-service/biz/utils"
	"go-template-service/conf"
	"net/http"
	"strings"
)

func JWKSMiddleware() app.HandlerFunc {
	// your code...
	return func(ctx context.Context, c *app.RequestContext) {
		// your code...
		jwksURL := conf.GetConf().Idp.JwksURL

		jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
		if err != nil {
			hlog.Errorf("Failed to create JWKs from URL. Error: %v", err)
			utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Get jwt access token from Authorization bearer token
		bearerToken := c.Request.Header.Get("Authorization")
		if bearerToken == "" {
			hlog.Errorf("No access token found.")
			c.JSON(consts.StatusUnauthorized, u.H{
				"message": "No access token found.",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		accessToken := strings.Split(bearerToken, " ")[1]

		// Parse the JWT and verify it
		token, err := jwt.Parse(accessToken, jwks.Keyfunc)
		if err != nil {
			hlog.Errorf("Failed to parse JWT. Error: %v", err)
			utils.SendErrResponse(ctx, c, consts.StatusUnauthorized, err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check if the token is valid
		if !token.Valid {
			hlog.Errorf("Invalid JWT.")
			c.JSON(consts.StatusUnauthorized, u.H{
				"message": "Invalid JWT.",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// store token claim into request context
		c.Set("currentUser", token.Claims)

		c.Next(ctx)
	}
}
