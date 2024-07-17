package google

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	Oauth2 "google.golang.org/api/oauth2/v2"
	"net/http"
	"string_backend_0001/internal/conf"
	"string_backend_0001/pkg"
)

var (
	cfg *oauth2.Config
)

const (
	STATE        = "state"
	Oauth2APIUrl = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

func NewGoogleOAuthConfig() *oauth2.Config {
	googleConf := conf.Conf.GoogleOauth
	return &oauth2.Config{
		ClientID:     googleConf.ClientID,
		ClientSecret: googleConf.ClientSecret,
		RedirectURL:  googleConf.RedirectURL,
		Scopes:       googleConf.Scopes,
		Endpoint:     google.Endpoint,
	}
}

func GetOauth2Config() *oauth2.Config {
	return cfg
}

func GetData(code string, config ...*oauth2.Config) (*Oauth2.Userinfo, error) {
	if len(config) == 0 {
		return getUserDataFromGoogle(code)
	} else {
		return getUserDataFromGoogle(code, config[0])
	}
}

func Router(r *gin.RouterGroup) {
	cfg = NewGoogleOAuthConfig()
	r.GET("/callback", callback)
	r.GET("/login", login)
}

func login(c *gin.Context) {
	c.Redirect(http.StatusFound, cfg.AuthCodeURL(STATE))
}

func callback(c *gin.Context) {
	state := c.Query(STATE)
	if state != STATE {
		c.JSON(pkg.CreateResponse(http.StatusUnauthorized, "invalid csrf token"))
		return
	}

	code := c.Query("code")
	userInfo, err := getUserDataFromGoogle(code)
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

func getUserDataFromGoogle(code string, config ...*oauth2.Config) (*Oauth2.Userinfo, error) {
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

	response, err := http.Get(Oauth2APIUrl + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer func() { response.Body.Close() }()

	var o Oauth2.Userinfo
	err = pkg.RespUnmarshal(response, &o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}
