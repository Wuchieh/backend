package line

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"string_backend_0001/internal/conf"
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
