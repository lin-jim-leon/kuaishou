package kuaishou

import (
	"github.com/lin-jim-leon/kuaishou/open/merchant"
	"github.com/lin-jim-leon/kuaishou/open/oauth"
	"github.com/lin-jim-leon/kuaishou/open/user"
)

// GetAccessToken 通过网页授权的code 换取access_token
func GetAccessToken(ClientKey string, ClientSecret string, code string) (accessToken oauth.AccessTokenRes, err error) {
	return oauth.GetAccessToken(ClientKey, ClientSecret, code)
}

// RefreshAccessToken 刷新AccessToken.
// 当access_token过期（过期时间2天）后，可以通过该接口使用refresh_token（过期时间180天）进行刷新
func RefreshAccessToken(refreshkey string, clientkey string, clientsecret string) (accessToken oauth.RefreshTokenRes, err error) {
	return oauth.RefreshAccessToken(refreshkey, clientkey, clientsecret)
}

// GetUserinfo 获取用户信息
func GetUserinfo(appid string, accesstoken string) (info user.Userresponse, err error) {
	return user.GetUserinfo(appid, accesstoken)
}

// 商品加橱窗
func AddItemsToShelf(Appkey string, signsecret string, accesstoken string, itemidlist []string) (Adinfo merchant.AddItemsres, err error) {
	return merchant.AddItemsToShelf(Appkey, accesstoken, signsecret, itemidlist)
}

// 查询选品详情
func Queryselectiondetail(Appkey string, signsecret string, accesstoken string, itemidlist []string) (Selectitem merchant.SelsetionRes, err error) {
	return merchant.Queryselectiondetail(Appkey, signsecret, accesstoken, itemidlist)
}

// 订单列表（游标）
func Corderlist(Appkey string, signsecret string, accesstoken string, cpsOrderStatus int, pageSize int, beginTime int64, endTime int64) (corlist merchant.OrderlistRes, err error) {
	return merchant.Corderlist(Appkey, signsecret, accesstoken, cpsOrderStatus, pageSize, beginTime, endTime)
}
