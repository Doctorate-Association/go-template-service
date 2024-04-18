package middleware

import (
	"context"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/golang-jwt/jwt/v5"
	"go-template-service/biz/utils"
	"go-template-service/conf"
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
			utils.SendErrResponse(ctx, c, 500, err)
			return
		}

		// Get jwt access token from Authorization bearer token
		bearerToken := c.Request.Header.Get("Authorization")
		if bearerToken == "" {
			hlog.Errorf("No access token found.")
			c.JSON(401, "No access token found.")
			return
		}

		accessToken := strings.Split(bearerToken, " ")[1]

		// Parse the JWT and verify it
		token, err := jwt.Parse(accessToken, jwks.Keyfunc)
		if err != nil {
			hlog.Errorf("Failed to parse JWT. Error: %v", err)
			utils.SendErrResponse(ctx, c, 401, err)
			return
		}

		// Check if the token is valid
		if !token.Valid {
			hlog.Errorf("Invalid JWT.")
			c.JSON(401, "Invalid JWT.")
			return
		}

		// store token claim into request context
		c.Set("currentUser", token.Claims)

		c.Next(ctx)
	}
}
