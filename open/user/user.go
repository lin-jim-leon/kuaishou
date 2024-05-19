package user

import (
	"encoding/json"
	"fmt"
	"github.com/lin-jim-leon/kuaishou/util"
)

const (
	UserinfoUrl = "https://open.kuaishou.com/openapi/user_info?app_id=%s&access_token=%s"
)

type Info struct {
	Name    string `json:"name"`
	Sex     string `json:"sex"`
	Fan     int    `json:"fan"`
	Follow  int    `json:"follow"`
	Head    string `json:"head"`
	BigHead string `json:"bigHead"`
	City    string `json:"city"`
}

type Userresponse struct {
	Result   int    `json:"result"`
	ErrorMsg string `json:"error_msg"`
	Data     Info   `json:"user_info"`
}

// GetUserinfo 获取用户信息
func GetUserinfo(appid string, accesstoken string) (info Userresponse, err error) {
	uri := fmt.Sprintf(UserinfoUrl, appid, accesstoken)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return Userresponse{}, err
	}
	var result Userresponse
	err = json.Unmarshal(response, &result)
	if err != nil {
		return Userresponse{}, err
	}
	if len(result.ErrorMsg) > 0 {
		return Userresponse{}, fmt.Errorf("GetUserinfo error: error_msg=%s", result.ErrorMsg)
	}
	return result, nil
}
