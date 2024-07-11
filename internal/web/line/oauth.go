package line

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"string_backend_0001/internal/conf"
	"string_backend_0001/internal/pkg"
	"string_backend_0001/internal/web/line/sdk"
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
		Endpoint:     sdk.Endpoint,
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

	c.Set(pkg.Google, userInfo)
	c.Next()
	if c.IsAborted() {
		return
	}
	c.JSON(pkg.CreateSuccessResponse(userInfo))
}

func login(c *gin.Context) {
	c.Redirect(http.StatusFound, cfg.AuthCodeURL(STATE))
}
func getUserDataFromLine(code string) (*sdk.Profile, error) {
	token, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	Oauth2 := sdk.CreateOauth2(token, code)

	return Oauth2.GetProfile()
}
