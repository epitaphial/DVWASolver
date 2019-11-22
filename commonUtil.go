package main

import (
    "fmt"
    "io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func DVWADebugger(debugInfo string){
	fmt.Println(debugInfo)
}

// 获取页面response报文，以及cookie
func GetContents(requrl string) (string ,string,error) {
    resp,err := http.Get(requrl)
    if err != nil {
        return "","",err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return "", "",fmt.Errorf("get content failed status code is %d ",resp.StatusCode)
    }

	cookie := resp.Cookies()[0].Value

    bytes,err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "" , "",nil
    }
    return string(bytes),string(cookie),nil
}

// 解析报文，拿取csrf token
func GetToken(contentLogin string) string{
	pattern := `<input type='hidden' name='user_token' value='(.*)' />`
	re := regexp.MustCompile(pattern)
	tokenList := re.FindAllStringSubmatch(contentLogin,-1)
	var tk string
	for _,token := range tokenList{
		tk = token[1]
	}
	return tk
}

// 注册cookie
func Postlogin(urlDVWALogin string,loginToken string,cookieLogin string) string{
	client := &http.Client{}
	params:=url.Values{}
	params.Add("username","admin")
	params.Add("password","password")
	params.Add("Login","Login")
	params.Add("user_token",loginToken)
	parmsenc := ioutil.NopCloser(strings.NewReader(params.Encode()))
	req, err := http.NewRequest("POST",urlDVWALogin, parmsenc)
	if err != nil{
		return ""
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:70.0) Gecko/20100101 Firefox/70.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "close")
	req.Header.Set("Cookie","PHPSESSID="+cookieLogin+"; security=low")
	resp, err := client.Do(req)
	if err != nil{
		return ""
	}
	defer resp.Body.Close()	
	return "ok"
}

//初始化
func InitDvwaUrl(urlDVWA string) string{
	urlDVWALogin := urlDVWA + "login.php"
	contentLogin,cookie,err := GetContents(urlDVWALogin)
	if err != nil {
		return ""
	}
	tokenLogin := GetToken(contentLogin)
	status := Postlogin(urlDVWALogin,tokenLogin,cookie)
	if status == "" {
		return status
	}else{
		return cookie
	}
}