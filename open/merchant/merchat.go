package merchant

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lin-jim-leon/kuaishou/util"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	PickUrl       = "https://openapi.kwaixiaodian.com/open/distribution/selection/pick?appkey=%s&access_token=%s&method=open.distribution.selection.pick&param=%s&sign=%s&version=1&signMethod=MD5&timestamp=%d"
	Selectdetail  = "https://openapi.kwaixiaodian.com/open/distribution/query/selection/item/detail?appkey=%s&access_token=%s&method=open.distribution.query.selection.item.detail&param=%s&sign=%s&version=1&signMethod=MD5&timestamp=%d"
	Corderlisturl = "https://openapi.kwaixiaodian.com/open/distribution/cps/distributor/order/cursor/list?appkey=%s&access_token=%s&method=open.distribution.cps.distributor.order.cursor.list&param=%s&sign=%s&version=1&signMethod=MD5&timestamp=%d"
)

type Iteminfo struct {
	Result   int    `json:"result"`
	ItemId   int64  `json:"itemId"`
	ErrorMsg string `json:"errorMsg"`
}

type AddItemsres struct {
	Result   int        `json:"result"`
	Msg      string     `json:"msg"`
	ErrorMsg string     `json:"error_msg"`
	Code     string     `json:"code"`
	Data     []Iteminfo `json:"data"`
	SubMsg   string     `json:"sub_msg"`
	SubCode  string     `json:"sub_code"`
}

// generateSign 生成签名
func generateSign(appkey, signsecret, accesstoken, param string, timestamp int64, method string) string {
	// 将请求参数按照字典顺序排序
	params := map[string]string{
		"access_token": accesstoken,
		"appkey":       appkey,
		"method":       method,
		"param":        param,
		"signMethod":   "MD5",
		"timestamp":    fmt.Sprintf("%d", timestamp),
		"version":      "1",
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构造待签名字符串
	var signStr strings.Builder
	for _, k := range keys {
		signStr.WriteString(fmt.Sprintf("%s=%s&", k, params[k]))
	}
	signStr.WriteString("signSecret=" + signsecret)

	// 计算签名值
	hasher := md5.New()
	hasher.Write([]byte(signStr.String()))
	encrypted := hasher.Sum(nil)
	return hex.EncodeToString(encrypted)
}

// 商品加橱窗
func AddItemsToShelf(Appkey string, signsecret string, accesstoken string, itemidlist []string) (Adinfo AddItemsres, err error) {
	// 构造请求参数 param
	rawParam := fmt.Sprintf(`{"itemIds":[%s]}`, strings.Join(itemidlist, ","))
	// 获取当前时间戳（毫秒级）
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	method := "open.distribution.selection.pick"
	// 对请求参数进行签名
	sign := generateSign(Appkey, signsecret, accesstoken, rawParam, timestamp, method)
	// 构造完整的请求 URL
	encodedParam := url.QueryEscape(rawParam)
	uri := fmt.Sprintf(PickUrl, Appkey, accesstoken, encodedParam, sign, timestamp)

	// 发送 HTTP GET 请求
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return AddItemsres{}, err
	}
	// 解析响应数据
	var result AddItemsres
	err = json.Unmarshal(response, &result)
	if err != nil {
		return AddItemsres{}, err
	}

	// 检查是否有错误信息
	if result.ErrorMsg != "SUCCESS" {
		return AddItemsres{}, fmt.Errorf("AddItemsToShelf error: %s", result.ErrorMsg)
	}

	return result, nil
}

type SelsetionRes struct {
	Code      string                    `json:"code"`
	Msg       string                    `json:"msg"`
	SubCode   string                    `json:"sub_code"`
	SubMsg    string                    `json:"sub_msg"`
	Result    int                       `json:"result"`
	ErrorMsg  string                    `json:"error_msg"`
	ItemList  []DistributeSelectionItem `json:"itemList"`
	ShopTitle string                    `json:"shopTitle"`
}

type DistributeSelectionItem struct {
	ItemID                      int64           `json:"itemId"`
	ZKFinalPrice                int64           `json:"zkFinalPrice"`
	ItemImgURL                  string          `json:"itemImgUrl"`
	ItemPrice                   int64           `json:"itemPrice"`
	SKUList                     []DistributeSKU `json:"skuList"`
	SoldCountThirtyDays         int64           `json:"soldCountThirtyDays"`
	ShopTitle                   string          `json:"shopTitle"`
	CommissionRate              int             `json:"commissionRate"`
	BrandName                   string          `json:"brandName"`
	ItemTitle                   string          `json:"itemTitle"`
	ProfitAmount                int64           `json:"profitAmount"`
	ShopScore                   string          `json:"shopScore"`
	ShopStar                    string          `json:"shopStar"`
	ItemDesc                    string          `json:"itemDesc"`
	ShopType                    int             `json:"shopType"`
	CategoryID                  int             `json:"categoryId"`
	ItemGalleryURLs             []string        `json:"itemGalleryUrls"`
	ItemDescURLs                []string        `json:"itemDescUrls"`
	MerchantSoldCountThirtyDays int64           `json:"merchantSoldCountThirtyDays"`
	ShopID                      int64           `json:"shopId"`
}

type DistributeSKU struct {
	SKUPrice      int64  `json:"skuPrice"`
	SKUStock      int64  `json:"skuStock"`
	Specification string `json:"specification"`
	SKUID         int64  `json:"skuId"`
}

func Queryselectiondetail(Appkey string, signsecret string, accesstoken string, itemidlist []string) (Selectitem SelsetionRes, err error) {
	// 构造请求参数 param
	rawParam := fmt.Sprintf(`{"itemId":[%s]}`, strings.Join(itemidlist, ","))
	// 获取当前时间戳（毫秒级）
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	method := "open.distribution.query.selection.item.detail"
	// 对请求参数进行签名
	sign := generateSign(Appkey, signsecret, accesstoken, rawParam, timestamp, method)
	// 构造完整的请求 URL
	encodedParam := url.QueryEscape(rawParam)
	uri := fmt.Sprintf(Selectdetail, Appkey, accesstoken, encodedParam, sign, timestamp)
	fmt.Println(uri)
	// 发送 HTTP GET 请求
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return SelsetionRes{}, err
	}
	fmt.Println(string(response))
	// 解析响应数据
	var result SelsetionRes
	err = json.Unmarshal(response, &result)
	if err != nil {
		return SelsetionRes{}, err
	}

	// 检查是否有错误信息
	if result.ErrorMsg != "SUCCESS" {
		return SelsetionRes{}, fmt.Errorf("Queryselectiondetail error: %s", result.ErrorMsg)
	}

	return result, nil
}

type OrderlistRes struct {
	Code     string      `json:"code"`      // 主返回码
	Msg      string      `json:"msg"`       // 主返回信息
	SubCode  string      `json:"sub_code"`  // 子返回码
	SubMsg   string      `json:"sub_msg"`   // 子返回信息
	Result   int         `json:"result"`    // 返回码
	ErrorMsg string      `json:"error_msg"` // 返回错误信息
	Data     OrderCursor `json:"data"`      // 返回数据
}

type OrderCursor struct {
	Cursor     string      `json:"pcursor"`   // 分销订单位点游标(请求透传, "nomore"标识后续无数据)
	OrderViews []OrderView `json:"orderView"` // 订单信息
}

type OrderView struct {
	DistributorID         int64                           `json:"distributorId"`         // 分销推广者用户ID
	OrderID               int64                           `json:"oid"`                   // 订单ID
	CPSOrderStatus        int                             `json:"cpsOrderStatus"`        // 分销订单状态
	OrderCreateTime       int64                           `json:"orderCreateTime"`       // 订单创建时间
	PayTime               int64                           `json:"payTime"`               // 订单支付时间
	OrderTradeAmount      int64                           `json:"orderTradeAmount"`      // 订单交易总金额
	CPSOrderProductViews  []CPSOrderProductView           `json:"cpsOrderProductView"`   // 商品信息
	CreateTime            int64                           `json:"createTime"`            // 创建时间
	UpdateTime            int64                           `json:"updateTime"`            // 更新时间
	SettlementSuccessTime int64                           `json:"settlementSuccessTime"` // 结算时间
	SettlementAmount      int64                           `json:"settlementAmount"`      // 结算金额
	SendTime              int64                           `json:"sendTime"`              // 订单发货时间
	SendStatus            int                             `json:"sendStatus"`            // 订单发货状态
	RecvTime              int64                           `json:"recvTime"`              // 订单收货时间
	BuyerOpenID           string                          `json:"buyerOpenId"`           // 买家唯一识别ID
	BaseAmount            int64                           `json:"baseAmount"`            // 计佣基数
	ShareRateStr          string                          `json:"shareRateStr"`          // 分成比例
	SettlementBizType     int                             `json:"settlementBizType"`     // 订单业务结算类型
	OrderRefundInfo       []CPSDistributorOrderRefundInfo `json:"orderRefundInfo"`       // 退款信息
	OrderChannel          string                          `json:"orderChannel"`          // 订单渠道标识
}

type CPSOrderProductView struct {
	OrderID              int64  `json:"oid"`                  // 订单ID
	ItemID               int64  `json:"itemId"`               // 商品ID
	ItemTitle            string `json:"itemTitle"`            // 商品标题
	ItemPrice            int64  `json:"itemPrice"`            // 商品单价快照
	EstimatedIncome      int64  `json:"estimatedIncome"`      // 预估收入
	CommissionRate       int    `json:"commissionRate"`       // 佣金比率
	ServiceRate          int    `json:"serviceRate"`          // 平台服务费率
	ServiceAmount        int64  `json:"serviceAmount"`        // 平台服务费
	CPSPID               string `json:"cpsPid"`               // 推广位id
	SellerID             int64  `json:"sellerId"`             // 商家Id
	SellerNickName       string `json:"sellerNickName"`       // 商家昵称快照
	Num                  int    `json:"num"`                  // 商品数量
	StepCondition        int    `json:"stepCondition"`        // 阶梯佣金条件
	StepCommissionRate   int    `json:"stepCommissionRate"`   // 阶梯佣金比率
	StepCommissionAmount int64  `json:"stepCommissionAmount"` // 阶梯佣金金额
	ServiceIncome        int64  `json:"serviceIncome"`        // 接单服务收入
	ExcitationIncome     int64  `json:"excitationIncome"`     // 奖励收入
}

type CPSDistributorOrderRefundInfo struct {
	StartRefundTime int64  `json:"startRefundTime"` // 发起退款时间
	EndRefundTime   int64  `json:"endRefundTime"`   // 完成退款时间
	RefundFee       int64  `json:"refundFee"`       // 退款金额
	RefundStatus    string `json:"refundStatus"`    // 退款状态
}

// 订单列表（游标）
func Corderlist(Appkey string, signsecret string, accesstoken string, cpsOrderStatus int, pageSize int, beginTime int64, endTime int64, pcursor string) (corlist OrderlistRes, err error) {
	// 构造请求参数 param
	rawParam := fmt.Sprintf(`{"cpsOrderStatus":%d,"pageSize":%d,"sortType":1,"queryType":1,"beginTime":%d,"endTime":%d,"pcursor":"%s"}`, cpsOrderStatus, pageSize, beginTime, endTime, pcursor)
	// 获取当前时间戳（毫秒级）
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	method := "open.distribution.cps.distributor.order.cursor.list"
	// 对请求参数进行签名
	sign := generateSign(Appkey, signsecret, accesstoken, rawParam, timestamp, method)
	// 构造完整的请求 URL
	encodedParam := url.QueryEscape(rawParam)
	uri := fmt.Sprintf(Corderlisturl, Appkey, accesstoken, encodedParam, sign, timestamp)

	// 发送 HTTP GET 请求
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return OrderlistRes{}, err
	}
	// 解析响应数据
	var result OrderlistRes
	err = json.Unmarshal(response, &result)
	if err != nil {
		return OrderlistRes{}, err
	}

	// 检查是否有错误信息
	if result.ErrorMsg != "SUCCESS" {
		return OrderlistRes{}, fmt.Errorf("Corderlist error: %s", result.ErrorMsg)
	}

	return result, nil
}
