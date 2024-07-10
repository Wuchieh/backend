package line

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"string_backend_0001/internal/conf"
	"string_backend_0001/internal/pkg"
)

var (
	cfg *oauth2.Config
)

const (
	STATE        = "state"
	Oauth2APIUrl = "https://api.line.me/v2/profile"
)

func NewGoogleOAuthConfig() *oauth2.Config {
	lineConf := conf.Conf.LineOauth
	return &oauth2.Config{
		ClientID:     lineConf.ClientID,
		ClientSecret: lineConf.ClientSecret,
		RedirectURL:  lineConf.RedirectURL,
		Scopes:       lineConf.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
			TokenURL: "https://api.line.me/oauth2/v2.1/token",
			//DeviceAuthURL: "https://oauth2.googleapis.com/device/code",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}

func Router(r *gin.RouterGroup) {
	cfg = NewGoogleOAuthConfig()
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

func getUserDataFromLine(code string) (*Profile, error) {
	token, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, Oauth2APIUrl, nil)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer func() { response.Body.Close() }()

	var o Profile
	err = pkg.RespUnmarshal(response, &o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}
