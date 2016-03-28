package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reusee/vviccommon"
)

type Api struct {
}

func NewApi() (*Api, error) {
	return new(Api), nil
}

type RespCommon struct {
	Ok bool `json:"ok"`
}

type PingReq struct {
	Greetings string
}
type PingResp struct {
	RespCommon
	Echo string
}

func (a *Api) Ping(req *PingReq, resp *PingResp) error {
	resp.Ok = true
	resp.Echo = req.Greetings
	return nil
}

type GetInfoReq struct {
	Id int
}

type GetInfoResp struct {
	RespCommon
	Title        string
	Price        string
	Id           int
	Attrs        [][]string
	ThemeImages  []string
	DetailImages []string
}

func (a *Api) GetInfo(req *GetInfoReq, resp *GetInfoResp) (err error) {
	defer ct(&err)
	pt("get info of %d\n", req.Id)
	pagePath := fmt.Sprintf("http://www.vvic.com/api/item/%d", req.Id)
	httpResp, err := http.Get(pagePath)
	ce(err, "get")
	defer httpResp.Body.Close()
	if httpResp.StatusCode == 400 {
		return nil
	}
	var data struct {
		Code int
		Data struct {
			Upload_num       int
			Color            string
			Is_tx            int
			Discount_price   string
			Index_img_url    string
			Is_df            int
			Title            string // 标题
			Tid              int64
			Price            string
			Color_pics       string
			Id               int // 货号
			Art_no           string
			Imgs             string // 图片
			Support_services string // 退现 实拍 代发
			Discount_value   float64
			Shop_name        string
			Discount_type    string
			Attrs            string // 属性
			Is_sp            int
			Shop_id          int
			Size             string // 尺寸
			Bname            string // 市场名
			Up_time          string
			Bid              int
			Tcid             string
			Status           int
			Cid              string
			Desc             string // 描述html
		}
	}
	err = json.NewDecoder(httpResp.Body).Decode(&data)
	ce(err, "decode")

	for _, imgPath := range strings.Split(data.Data.Imgs, ",") {
		if !strings.HasPrefix(imgPath, "http:") {
			imgPath = "http:" + imgPath
		}
		resp.ThemeImages = append(resp.ThemeImages, imgPath)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data.Data.Desc))
	ce(err, "goquery doc")
	doc.Find("img").Each(func(i int, se *goquery.Selection) {
		imgSrc, _ := se.Attr("src")
		resp.DetailImages = append(resp.DetailImages, imgSrc)
	})

	resp.Title = vviccommon.TidyTitle(data.Data.Title)
	resp.Price = data.Data.Discount_price
	resp.Id = data.Data.Id

	attrs := map[string]string{}
	for _, attr := range strings.Split(data.Data.Attrs, ",") {
		parts := strings.SplitN(attr, ":", 2)
		attrs[parts[0]] = parts[1]
	}
	attrKeys := []string{
		"风格",
		"裙长",
		"版型",
		"领型",
		"袖型",
		"元素",
		"颜色",
		"尺码",
		"图案",
		"适用",
		"组合",
		"款式",
		"袖长",
		"腰型",
		"门襟",
		"裙型",
		"质地",
	}
loop_key:
	for _, key := range attrKeys {
		for attrKey, attr := range attrs {
			if strings.Contains(attrKey, key) {
				resp.Attrs = append(resp.Attrs, []string{
					attrKey, attr,
				})
				delete(attrs, attrKey)
				continue loop_key
			}
		}
	}
	for attrKey, attr := range attrs {
		resp.Attrs = append(resp.Attrs, []string{
			attrKey, attr,
		})
	}

	resp.Ok = true
	return nil
}
