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

var (
	oauthConfigMap = map[string]func() *oauth2.Config{
		Google:  google.GetOauth2Config,
		Discord: discord.GetOauth2Config,
		Line:    line.GetOauth2Config,
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

type loginResp struct {
	ID       string `json:"id"`
	Redirect string `json:"redirect"`
}

func login(c *gin.Context) {
	platform := c.Param("platform")
	p, ok := oauthConfigMap[platform]

	if !ok {
		c.JSON(pkg.CreateResponse(http.StatusBadRequest, "invalid platform"))
		return
	}

	web := conf.Conf.Web
	cfg := *p()
	cfg.RedirectURL = fmt.Sprintf("http://%s:%d/api/oauth/callback", web.Host, web.Port)

	s := newStore()
	s.Platform = platform
	s.State = 1

	c.JSON(pkg.CreateSuccessResponse(loginResp{
		ID:       s.ID,
		Redirect: cfg.AuthCodeURL(s.ID),
	}))
}

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
	web := conf.Conf.Web
	cfg := *p()
	cfg.RedirectURL = fmt.Sprintf("http://%s:%d/api/oauth/callback", web.Host, web.Port)

	var data any
	var err error

	switch Store.Platform {
	case Google:
		data, err = google.GetData(code, &cfg)
	case Discord:
		data, err = discord.GetData(code, &cfg)
	case Line:
		data, err = line.GetData(code, &cfg)
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
