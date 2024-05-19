package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/lin-jim-leon/kuaishou/util"
)

const (
	accessTokenUrl  = "https://open.kuaishou.com/oauth2/access_token?grant_type=authorization_code&app_id=%s&app_secret=%s&code=%s"
	refreshTokenURL = "https://open.kuaishou.com/oauth2/refresh_token?grant_type=refresh_token&app_id=%s&app_secret=%s&refresh_token=%s"
)

// accesstoken信息
type AccessTokenRes struct {
	Result                int      `json:"result"`
	AccessToken           string   `json:"access_token"`
	ExpiresIn             int      `json:"expires_in"`
	RefreshToken          string   `json:"refresh_token"`
	RefreshTokenExpiresIn int      `json:"refresh_token_expires_in"`
	OpenId                string   `json:"open_id"`
	Scopes                []string `json:"scopes"`
	ErrorMsg              string   `json:"error_msg"`
}

// GetAccessToken 通过网页授权的code 换取access_token
func GetAccessToken(ClientKey string, ClientSecret string, code string) (accessToken AccessTokenRes, err error) {
	uri := fmt.Sprintf(accessTokenUrl, ClientKey, ClientSecret, code)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return AccessTokenRes{}, err
	}
	var result AccessTokenRes
	err = json.Unmarshal(response, &result)
	if err != nil {
		return AccessTokenRes{}, err
	}
	if len(result.ErrorMsg) > 0 {
		return AccessTokenRes{}, fmt.Errorf("GetAccessToken error: error_msg=%s", result.ErrorMsg)
	}
	return result, nil
}

type RefreshTokenRes struct {
	Result                int      `json:"result"`
	AccessToken           string   `json:"access_token"`
	ExpiresIn             int      `json:"expires_in"`
	RefreshToken          string   `json:"refresh_token"`
	RefreshTokenExpiresIn int      `json:"refresh_token_expires_in"`
	Scopes                []string `json:"scopes"`
	ErrorMsg              string   `json:"error_msg"`
}

// RefreshAccessToken 刷新AccessToken.
// 当access_token过期（过期时间2天）后，可以通过该接口使用refresh_token（过期时间180天）进行刷新
func RefreshAccessToken(refreshkey string, clientkey string, clientsecret string) (accessToken RefreshTokenRes, err error) {
	uri := fmt.Sprintf(refreshTokenURL, clientkey, clientsecret, refreshkey)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return RefreshTokenRes{}, err
	}
	var result RefreshTokenRes
	err = json.Unmarshal(response, &result)
	if err != nil {
		return RefreshTokenRes{}, err
	}
	if len(result.ErrorMsg) > 0 {
		return RefreshTokenRes{}, fmt.Errorf("RefreshAccessToken error: error_msg=%s", result.ErrorMsg)
	}
	return result, nil
}
