package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/skip2/go-qrcode"
	"strings"
	"net/url"
)

const BASE_URL = "https://login.weixin.qq.com/"


func main() {
	getQruuid()
}

func a(){
	c, err := qrcode.Encode("http://www.baidu.com", qrcode.Medium, 256)
	if err != nil {
		panic(err)
	}
	print(string(c))

}

func getQruuid(){
	params := url.Values{}
	params.Add("redirect_url", "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage")
	params.Add("fun", "new")
	params.Add("lang", "zh_CN")
	params.Add("appid", "wx782c26e4c19acffb")
	params.Add("_", string(1525426049365))

	resp, err := http.NewRequest("GET", BASE_URL, strings.NewReader(params.Encode()))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))
}

func getQuuid(){
	//u, _ := url.Parse(BASE_URL)
}