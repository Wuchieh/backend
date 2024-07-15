package discord

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
	"string_backend_0001/internal/conf"
	"string_backend_0001/internal/pkg"
	"string_backend_0001/sdk/discord"
)

var (
	cfg *oauth2.Config
)

const (
	STATE = "state"
)

func NewOAuthConfig() *oauth2.Config {
	Conf := conf.Conf.DiscordOauth
	return &oauth2.Config{
		ClientID:     Conf.ClientID,
		ClientSecret: Conf.ClientSecret,
		RedirectURL:  Conf.RedirectURL,
		Scopes:       Conf.Scopes,
		Endpoint:     discord.Endpoint,
	}
}

func Router(r *gin.RouterGroup) {
	cfg = NewOAuthConfig()
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
	userInfo, err := getUserData(code)
	if err != nil {
		c.JSON(pkg.CreateResponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.Set(pkg.Discord, userInfo)
	c.Next()
	if c.IsAborted() {
		return
	}
	c.JSON(pkg.CreateSuccessResponse(userInfo))
}

func getUserData(code string) (*discord.User, error) {
	oauth := discord.NewOauth(cfg)
	err := oauth.Exchange(code)
	if err != nil {
		return nil, err
	}
	data, err := oauth.GetUserData()
	return data, err
}
