package oauth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"net/http"
	"string_backend_0001/internal/conf"
	"string_backend_0001/internal/logger"
	"string_backend_0001/internal/web/oauth/discord"
	"string_backend_0001/internal/web/oauth/google"
	"string_backend_0001/internal/web/oauth/line"
	"string_backend_0001/pkg"
	"time"
)

type store struct {
	ID       string
	Platform string
	Data     any

	/*
		0 準備中
		1 登入中
		2 登入成功
		10 登入失敗
	*/
	State int

	C chan struct{}
}

func (s *store) SetData(data any) {
	defer func() {
		s.C <- struct{}{}
	}()
	s.Data = data
}

func (s *store) LoginSuccess(data any) {
	s.State = 2
	s.SetData(data)
}

func (s *store) LoginFail(data any) {
	s.State = 10
	s.SetData(data)
}

func newStore() *store {
	id := uuid.NewString()
	s := &store{
		ID: id,
		C:  make(chan struct{}, 1),
	}
	storeMap[id] = s
	return s
}

const (
	Google  = "google"
	Discord = "discord"
	Line    = "line"
)

var getOauth = func(s func() *oauth2.Config) func() *oauth2.Config {
	var o *oauth2.Config
	return func() *oauth2.Config {
		if o == nil {
			o = new(oauth2.Config)
			*o = *s()
			web := conf.Conf.Web
			o.RedirectURL = fmt.Sprintf("http://%s:%d/api/oauth/callback", web.Host, web.Port)
		}
		return o
	}
}

var (
	oauthConfigMap = map[string]func() *oauth2.Config{
		Google:  getOauth(google.GetOauth2Config),
		Discord: getOauth(discord.GetOauth2Config),
		Line:    getOauth(line.GetOauth2Config),
	}

	storeMap = make(map[string]*store)

	ErrInvalidPlatform = errors.New("invalid platform")
)

func Router(r *gin.RouterGroup) {
	r.Use(cors)
	r.GET("/login/:platform", login)
	r.GET("/data/:id", getData)
	r.GET("/callback", callback)
}

func cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

type LoginResp struct {
	ID       string `json:"id"`
	Redirect string `json:"redirect"`
}

// @summary get oauth2 login data
// @description platform = google | line | discord
// @tags OAuth2
// @id OAuth2Login
// @accept json
// @produce json
// @param platform path string true "google | line | discord"
// @success 200 {object} pkg.Response{data=LoginResp}
// @router /api/oauth/login/{platform} [get]
func login(c *gin.Context) {
	platform := c.Param("platform")
	p, ok := oauthConfigMap[platform]

	if !ok {
		c.JSON(pkg.CreateResponse(http.StatusBadRequest, "invalid platform"))
		return
	}

	cfg := p()

	s := newStore()
	s.Platform = platform
	s.State = 1

	c.JSON(pkg.CreateSuccessResponse(LoginResp{
		ID:       s.ID,
		Redirect: cfg.AuthCodeURL(s.ID),
	}))
}

// @summary get oauth2 user data
// @description id from OAuth2Login
// @tags OAuth2
// @id OAuth2GetData
// @accept json
// @produce json
// @param id path string true "id"
// @success 200 {object} pkg.Response{data=any}
// @router /api/oauth/data/{id} [get]
func getData(c *gin.Context) {
	id := c.Param("id")
	Store, ok := storeMap[id]

	if !ok {
		c.JSON(pkg.CreateResponse(http.StatusBadRequest, "invalid id"))
		return
	}

	if Store.Data != nil {
		c.JSON(pkg.CreateSuccessResponse(Store.Data))
		return
	}

	select {
	case <-Store.C:
		if Store.State == 2 {
			c.JSON(pkg.CreateSuccessResponse(Store.Data))
		} else {
			c.JSON(pkg.CreateResponse(http.StatusBadRequest, "login fail"))
		}
		delete(storeMap, Store.ID)
	case <-time.After(time.Second * 30):
		c.JSON(pkg.CreateResponse(http.StatusBadRequest, "timeout"))
	}
}

func callback(c *gin.Context) {
	state := c.Query("state")

	code := c.Query("code")
	Store, ok := storeMap[state]

	if !ok {
		c.JSON(pkg.CreateResponse(http.StatusBadRequest, "invalid state"))
		return
	}

	p := oauthConfigMap[Store.Platform]
	cfg := p()

	var data any
	var err error

	switch Store.Platform {
	case Google:
		data, err = google.GetData(code, cfg)
	case Discord:
		data, err = discord.GetData(code, cfg)
	case Line:
		data, err = line.GetData(code, cfg)
	default:
		err = ErrInvalidPlatform
	}

	if err != nil {
		if !errors.Is(err, ErrInvalidPlatform) {
			logger.Error(err.Error())
		}
		Store.LoginFail(data)
		c.JSON(pkg.CreateResponse(http.StatusBadRequest, "get user data fail", err.Error()))
		return
	}

	Store.LoginSuccess(data)

	c.JSON(pkg.CreateSuccessResponse())
}
