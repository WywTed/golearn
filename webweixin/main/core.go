package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/skip2/go-qrcode"
	"strings"
	"net/url"
	"net"
	"os"
	"html"
	"encoding/json"
)

const BASE_URL = "https://login.weixin.qq.com/"


func main() {
	requestListen()
}

func requestListen(){
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(":8080", nil)
}
func requestHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fmt.Fprintf(w, "hi i love you %s", html.EscapeString(r.URL.Path[1:]))
	if r.Method == "GET"{
		fmt.Println("method: ", r.Method)
		fmt.Println("username: ", r.Form["username"])
		fmt.Println("password: ", r.Form["password"])
		for k, v := range r.Form{
			fmt.Print("key: ", k, ";")
			fmt.Print("value: ", strings.Join(v, ""))
		}

	} else if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Printf("%s\n", result)

		if result == nil {
			fmt.Println("result is null")
			return
		}

		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		for k, v := range m {
			switch vv := v.(type) {
			case string:
				fmt.Println(k, "is string", vv)
			case int:
				fmt.Println(k, "is int", vv)
			case float64:
				fmt.Println(k, "is float64", vv)
			case []interface{}:
				fmt.Println(k, "is an array:")
				for i, u := range vv {
					fmt.Println(i, u)
				}
			default:
				fmt.Println(k, "is of a type i don't know how to handle")
			}
		}
		var s Serverslice;
		json.Unmarshal([]byte(result), &s)
		fmt.Println(s.ServersID)
		for i:= 0; i<len(s.Servers); i++{
			fmt.Println(s.Servers[i].ServerName + "  " + s.Servers[i].ServerIP)
		}
	}

}

type Server struct{
	ServerName string
	ServerIP string
}
type Serverslice struct{
	Servers []Server
	ServersID string
}


func server(){
	if os.Args[1] == "server" {
		server, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
			return;
		}
		for {
			con, err := server.Accept()
			if err != nil {
				panic(err)
			}
			con.Write([]byte("hello world"))
		}
	} else {
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			panic(err)
		}
		for {
			buf := make([]byte, 1024)
			if length, err := conn.Read(buf); err == nil{
				if length > 0 {
					buf[length] = 0
					fmt.Printf("%s", string(buf[0:length]))
				}
			}
		}
	}

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