package line

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"string_backend_0001/internal/conf"
	"string_backend_0001/internal/logger"
	"string_backend_0001/pkg"
	"string_backend_0001/sdk/line"
)

var (
	cfg *oauth2.Config
)

const (
	STATE = "state"
)

func NewOAuthConfig() *oauth2.Config {
	lineConf := conf.Conf.LineOauth
	return &oauth2.Config{
		ClientID:     lineConf.ClientID,
		ClientSecret: lineConf.ClientSecret,
		RedirectURL:  lineConf.RedirectURL,
		Scopes:       lineConf.Scopes,
		Endpoint:     line.Endpoint,
	}
}

func GetOauth2Config() *oauth2.Config {
	return cfg
}

func GetData(code string, config ...*oauth2.Config) (*line.Profile, error) {
	if len(config) == 0 {
		return getUserDataFromLine(code)
	} else {
		return getUserDataFromLine(code, config[0])
	}
}

func Router(r *gin.RouterGroup) {
	cfg = NewOAuthConfig()
	r.GET("/callback", callback)
	r.GET("/login", login)
}

func callback(c *gin.Context) {
	callbackErr := c.Query("error")
	if callbackErr != "" {
		handler, err := line.ErrHandler(c.FullPath())
		if err != nil {
			logger.Error("%+v", err)
			c.JSON(pkg.CreateResponse(http.StatusInternalServerError, "undefined error"))
		} else {
			switch handler.Type {
			case line.ErrInvalidRequest:

			case line.ErrAccessDenied:
				c.JSON(pkg.CreateResponse(http.StatusUnauthorized, "登入失敗, 用戶拒絕登入"))
			case line.ErrUnsupportedResponseType:
				logger.Error("line oauth callback error type: %s, description:%s", handler.Type, handler.Description)
				c.JSON(pkg.CreateResponse(http.StatusInternalServerError, "請聯繫管理員"))
			case line.ErrInvalidScope:
				logger.Error("line oauth callback error type: %s, description:%s", handler.Type, handler.Description)
				c.JSON(pkg.CreateResponse(http.StatusInternalServerError, "請聯繫管理員"))
			case line.ErrServerError:
				logger.Error("line oauth callback error type: %s, description:%s", handler.Type, handler.Description)
				c.JSON(pkg.CreateResponse(http.StatusInternalServerError, "line 登入伺服器錯誤"))
			default:
				logger.Error("line oauth callback error type: %s, description:%s", handler.Type, handler.Description)
				c.JSON(pkg.CreateResponse(http.StatusInternalServerError, "發生未知錯誤"))
			}
		}
		return
	}

	state := c.Query(STATE)
	if state != STATE {
		c.JSON(pkg.CreateResponse(http.StatusUnauthorized, "invalid csrf token"))
		return
	}

	code := c.Query("code")
	userInfo, err := getUserDataFromLine(code)
	if err != nil {
		c.JSON(pkg.CreateResponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.Set(pkg.Line, userInfo)
	c.Next()
	if c.IsAborted() {
		return
	}
	c.JSON(pkg.CreateSuccessResponse(userInfo))
}

func login(c *gin.Context) {
	c.Redirect(http.StatusFound, cfg.AuthCodeURL(STATE))
}

func getUserDataFromLine(code string, config ...*oauth2.Config) (*line.Profile, error) {
	cfg := func() *oauth2.Config {
		if len(config) > 0 {
			return config[0]
		} else {
			return cfg
		}
	}()

	token, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	Oauth2 := line.CreateOauth2(token, code)

	return Oauth2.GetProfile()
}
