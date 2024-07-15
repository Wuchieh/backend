package discord

import (
	"context"
	"golang.org/x/oauth2"
)

var (
	Endpoint = oauth2.Endpoint{
		AuthURL:   "https://discord.com/oauth2/authorize",
		TokenURL:  "https://discord.com/api/oauth2/token",
		AuthStyle: oauth2.AuthStyleInParams,
	}
)

const (
	RevokeUrl = "https://discord.com/api/oauth2/token/revoke"
)

type User struct {
	AccentColor int    `json:"accent_color"`
	Avatar      string `json:"avatar"`
	//AvatarDecorationData interface{} `json:"avatar_decoration_data"`
	Banner      string `json:"banner,omitempty"`
	BannerColor string `json:"banner_color"`
	//Clan          interface{} `json:"clan"`
	Discriminator string `json:"discriminator"`
	Email         string `json:"email"`
	Flags         int    `json:"flags"`
	GlobalName    string `json:"global_name"`
	Id            string `json:"id"`
	Locale        string `json:"locale"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	PremiumType   int    `json:"premium_type"`
	PublicFlags   int    `json:"public_flags"`
	Username      string `json:"username"`
	Verified      bool   `json:"verified"`
}

type Oauth struct {
	oauth *oauth2.Config

	Token *oauth2.Token
}

func NewOauth(o *oauth2.Config) *Oauth {
	return &Oauth{
		oauth: o,
	}
}

func (o *Oauth) GetUserData() (*User, error) {
	resp, err := o.oauth.Client(context.Background(), o.Token).Get("https://discord.com/api/users/@me")
	//resp, err := o.oauth.Client(context.Background(), o.Token).Get("https://discord.com/api/oauth2/@me")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var u User
	err = RespUnmarshal(resp, &u)

	return &u, err
}

func (o *Oauth) Exchange(code string) error {
	token, err := o.oauth.Exchange(context.Background(), code)
	if err != nil {
		return err
	}
	o.Token = token
	return nil
}
