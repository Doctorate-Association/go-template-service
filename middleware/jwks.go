package middleware

import (
	"context"
	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"
	u "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-jwt/jwt/v5"
	"go-template-service/biz/dal/redis"
	"go-template-service/biz/models"
	"go-template-service/biz/utils"
	"go-template-service/conf"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var kf keyfunc.Keyfunc

// JWKSMiddleware  middleware to verify JWT token
func JWKSMiddleware() app.HandlerFunc {
	// your code...
	return func(ctx context.Context, c *app.RequestContext) {
		// your code...
		if kf == nil {
			nkf, err := createNewKeyFunc(ctx)
			if err != nil {
				hlog.Errorf("Failed to create JWKs from URL. Error: %v", err)
				utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			kf = nkf
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

		val, err := redis.RedisClient.Get(ctx, accessToken).Result()

		if err == nil {

			var cu models.User
			err = json.Unmarshal([]byte(val), &cu)
			if err != nil {
				hlog.Errorf("Failed to unmarshal user info. Error: %v", err)
				utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			hlog.Infof("Get user  %s info from redis from", cu.DisplayName)
			c.Set("currentUser", cu)
			c.Next(ctx)
			return
		}

		// Parse the JWT and verify it
		token, err := jwt.Parse(accessToken, kf.Keyfunc)
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

		cu := models.User{
			ID:          token.Claims.(jwt.MapClaims)["id"].(string),
			Name:        token.Claims.(jwt.MapClaims)["name"].(string),
			DisplayName: token.Claims.(jwt.MapClaims)["displayName"].(string),
			Avatar:      token.Claims.(jwt.MapClaims)["avatar"].(string),
			Email:       token.Claims.(jwt.MapClaims)["email"].(string),
			Phone:       token.Claims.(jwt.MapClaims)["phone"].(string),
			Exp:         token.Claims.(jwt.MapClaims)["exp"].(float64),
			Nvb:         token.Claims.(jwt.MapClaims)["nbf"].(float64),
			Iat:         token.Claims.(jwt.MapClaims)["iat"].(float64),
			Jti:         token.Claims.(jwt.MapClaims)["jti"].(string),
		}
		// store token claim into request context
		c.Set("currentUser", cu)

		j, err := json.Marshal(cu)
		if err != nil {
			hlog.Errorf("Failed to marshal user info. Error: %v", err)
			utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		//Store user into redis
		err = redis.RedisClient.Set(ctx, accessToken, j, getTimeDiff(token.Claims)).Err()
		if err != nil {
			hlog.Errorf("Failed to set user info to redis. Error: %v", err)
			utils.SendErrResponse(ctx, c, consts.StatusInternalServerError, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Next(ctx)
	}
}

func createNewKeyFunc(ctx context.Context) (keyfunc.Keyfunc, error) {
	jwksURL := conf.GetConf().Idp.JwksURL

	given := jwkset.NewMemoryStorage()

	remoteJWKSets := make(map[string]jwkset.Storage)

	ur, err := url.ParseRequestURI(jwksURL)

	jwksetHTTPStorageOptions := jwkset.HTTPClientStorageOptions{
		Client:                    http.DefaultClient, // Could be replaced with a custom client.
		Ctx:                       ctx,                // Used to end background refresh goroutine.
		HTTPExpectedStatus:        http.StatusOK,
		HTTPMethod:                http.MethodGet,
		HTTPTimeout:               10 * time.Second, // Timeout for the HTTP request.
		NoErrorReturnFirstHTTPReq: true,             // Create storage regardless if the first HTTP request fails.
		RefreshErrorHandler: func(ctx context.Context, err error) {
			hlog.CtxErrorf(ctx, "Failed to refresh HTTP JWK Set from remote HTTP resource from url:%s Error: %v", ur.String(), err)
		},
		RefreshInterval: time.Hour, // How often to refresh the JWK Set
		Storage:         nil,
	}

	store, err := jwkset.NewStorageFromHTTP(ur, jwksetHTTPStorageOptions)
	if err != nil {
		log.Fatalf("Failed to create HTTP client storage for %q: %s", ur, err)
	}
	remoteJWKSets[ur.String()] = store

	jwksetHTTPClientOptions := jwkset.HTTPClientOptions{
		Given:          given,
		HTTPURLs:       remoteJWKSets,
		PrioritizeHTTP: false,
	}
	combined, err := jwkset.NewHTTPClient(jwksetHTTPClientOptions)

	keyfuncOptions := keyfunc.Options{
		Ctx:          ctx,
		Storage:      combined,
		UseWhitelist: []jwkset.USE{jwkset.UseSig},
	}
	if err != nil {
		log.Fatalf("Failed to create HTTP client storage: %s", err)
	}

	return keyfunc.New(keyfuncOptions)
}

func getTimeDiff(claims jwt.Claims) time.Duration {
	exp := claims.(jwt.MapClaims)["exp"].(float64)
	expTime := time.Unix(int64(exp), 0)
	now := time.Now()
	return expTime.Sub(now)
}
