package line

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	IdToken             = "id_token"
	ExpiresIn           = "expires_in"
	Scope               = "scope"
	STATE               = "state"
	ProfileUrl          = "https://api.line.me/v2/profile"
	TokenUrl            = "https://api.line.me/oauth2/v2.1/token"
	VerifyUrl           = "https://api.line.me/oauth2/v2.1/verify"
	RevokeUrl           = "https://api.line.me/oauth2/v2.1/revoke"
	DeauthorizeUrl      = "https://api.line.me/user/v1/deauthorize"
	UserinfoUrl         = "https://api.line.me/oauth2/v2.1/userinfo"
	FriendshipStatusUrl = "https://api.line.me/friendship/v1/status"
	AuthURL             = "https://access.line.me/oauth2/v2.1/authorize"
)

var (
	Endpoint = oauth2.Endpoint{
		AuthURL:   AuthURL,
		TokenURL:  TokenUrl,
		AuthStyle: oauth2.AuthStyleInParams,
		//DeviceAuthURL: "https://oauth2.googleapis.com/device/code",
	}
)

type IssueAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type AccessTokenValidity struct {
	Scope     string `json:"scope"`
	ClientId  string `json:"client_id"`
	ExpiresIn int    `json:"expires_in"`
}

type RefreshAccessToken struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyIDTokenOption struct {
	Nonce  string `json:"nonce,omitempty"`
	UserId string `json:"user_id,omitempty"`
}

type VerifyIDToken struct {
	Iss     string   `json:"iss"`
	Sub     string   `json:"sub"`
	Aud     string   `json:"aud"`
	Exp     int      `json:"exp"`
	Iat     int      `json:"iat"`
	Nonce   string   `json:"nonce"`
	Amr     []string `json:"amr"`
	Name    string   `json:"name"`
	Picture string   `json:"picture"`
	Email   string   `json:"email"`
}

type FriendFlag struct {
	FriendFlag bool `json:"friendFlag"`
}

type Profile struct {
	UserId        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureUrl    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

type Oauth2 struct {
	AccessToken  string    `json:"access_token"`
	ExpiresIn    int       `json:"expires_in"`
	IdToken      string    `json:"id_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	Scope        string    `json:"scope"`
	TokenType    string    `json:"token_type"`

	code    string
	profile *Profile
}

func CreateOauth2(token *oauth2.Token, code string) *Oauth2 {
	base := &Oauth2{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		TokenType:    token.TokenType,
		code:         code,
	}

	if idToken, ok := token.Extra(IdToken).(string); ok {
		base.IdToken = idToken
	}

	if expiresIn, ok := token.Extra(ExpiresIn).(int); ok {
		base.ExpiresIn = expiresIn
	}

	if scope, ok := token.Extra(Scope).(string); ok {
		base.Scope = scope
	}

	return base
}

func (o *Oauth2) GetProfile() (*Profile, error) {
	if o.profile != nil {
		return o.profile, nil
	}

	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, ProfileUrl, nil)
	req.Header.Set("Authorization", "Bearer "+o.AccessToken)

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer func() { response.Body.Close() }()

	var p Profile
	err = RespUnmarshal(response, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (o *Oauth2) IssueAccessToken(cfg *oauth2.Config, codeVerifier ...string) (*IssueAccessToken, error) {
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", o.code)
	data.Add("redirect_uri", cfg.AuthCodeURL(STATE))
	data.Add("client_id", cfg.ClientID)
	data.Add("client_secret", cfg.ClientSecret)
	if len(codeVerifier) > 0 {
		data.Add("code_verifier", codeVerifier[0])
	}

	req, err := http.NewRequest(http.MethodPost, TokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { resp.Body.Close() }()

	var ia IssueAccessToken
	err = RespUnmarshal(resp, &ia)
	if err != nil {
		return nil, err
	}

	return &ia, nil
}

func (o *Oauth2) VerifyAccessTokenValidity() (*AccessTokenValidity, error) {
	client := http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s?access_token=%s", VerifyUrl, o.AccessToken))
	if err != nil {
		return nil, err
	}

	defer func() { resp.Body.Close() }()
	var a AccessTokenValidity
	err = RespUnmarshal(resp, &a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (o *Oauth2) RefreshAccessToken() (*RefreshAccessToken, error) {
	data := url.Values{}
	data.Add("grant_type", "refresh_token")
	data.Add("refresh_token", o.RefreshToken)

	req, err := http.NewRequest(http.MethodPost, TokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var r RefreshAccessToken
	err = RespUnmarshal(resp, &r)
	if err != nil {
		return nil, err
	}

	o.AccessToken = r.AccessToken
	o.ExpiresIn = r.ExpiresIn
	o.RefreshToken = r.RefreshToken
	o.Scope = r.Scope

	return &r, nil
}

func (o *Oauth2) RevokeAccessToken(clientId string, clientSecret ...string) bool {
	data := url.Values{}
	data.Add("access_token", o.AccessToken)
	data.Add("client_id", clientId)
	if len(clientSecret) > 0 {
		data.Add("code_verifier", clientSecret[0])
	}

	req, err := http.NewRequest(http.MethodPost, RevokeUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func() { resp.Body.Close() }()

	return resp.StatusCode == http.StatusOK
}

func (o *Oauth2) DeauthorizeYourAppToWhichTheUserHasGrantedPermissions(userAccessToken string) bool {
	var data map[string]string
	data["userAccessToken"] = userAccessToken

	jsonData, err := json.Marshal(data)
	if err != nil {
		return false
	}

	req, err := http.NewRequest(http.MethodPost, DeauthorizeUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.AccessToken)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func() { resp.Body.Close() }()

	return resp.StatusCode == http.StatusNoContent
}

func (o *Oauth2) VerifyIDToken(clientId string, option ...VerifyIDTokenOption) (*VerifyIDToken, error) {
	if o.IdToken == "" {
		return nil, errors.New("not get id token")
	}

	data := url.Values{}
	data.Add("id_token", o.IdToken)
	data.Add("client_id", clientId)
	if len(option) > 0 {
		if option[0].Nonce != "" {
			data.Add("nonce", option[0].Nonce)
		}
		if option[0].UserId != "" {
			data.Add("user_id", option[0].UserId)
		}
	}

	req, err := http.NewRequest(http.MethodPost, VerifyUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var lineErr RespError
		err = RespUnmarshal(resp, &lineErr)
		return nil, &lineErr
	}

	var v VerifyIDToken
	err = RespUnmarshal(resp, &v)
	if err != nil {
		return nil, err
	}

	o.profile = &Profile{
		UserId:      v.Sub,
		DisplayName: v.Name,
		PictureUrl:  v.Picture,
	}

	return &v, nil
}

func (o *Oauth2) GetUserInformation() (*Profile, error) {
	req, err := http.NewRequest(http.MethodGet, UserinfoUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+o.AccessToken)
	client := http.Client{}
	do, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { do.Body.Close() }()

	var p Profile
	err = RespUnmarshal(do, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (o *Oauth2) GetUserProfile() (*Profile, error) {
	return o.GetProfile()
}

func (o *Oauth2) GetFriendshipStatus() (*FriendFlag, error) {
	req, _ := http.NewRequest(http.MethodGet, FriendshipStatusUrl, nil)
	req.Header.Set("Authorization", "Bearer "+o.AccessToken)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { resp.Body.Close() }()

	var f FriendFlag
	err = RespUnmarshal(resp, &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
