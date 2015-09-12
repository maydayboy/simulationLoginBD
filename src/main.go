package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"runtime"
	"strings"
)

var gCurCookies []*http.Cookie
var gCurCookiesJar *cookiejar.Jar
var gLog *log.Logger

func init() {
	gCurCookies = nil
	gCurCookiesJar, _ = cookiejar.New(nil)

}

//Get Current Path filename
func getCurFileName() string {
	_, fileName, _, _ := runtime.Caller(0)
	var fileNameWithSuffix string
	fileNameWithSuffix = path.Base(fileName)
	var fileSuffix string
	fileSuffix = path.Ext(fileNameWithSuffix)
	var fileNameOnly string
	fileNameOnly = strings.TrimSuffix(fileNameWithSuffix, fileSuffix)
	return fileNameOnly
}
func newLogger() (file *os.File, gLogger *log.Logger) {
	var fileNameOnly string
	fileNameOnly = getCurFileName()
	var logFileName = fileNameOnly + ".log"
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("open File error =", err)
		os.Exit(-1)
	}

	logger := log.New(logFile, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)

	return logFile, logger
}
func getRespHtml(url string) string {
	gLog.Println("getResphtml,url=", url)
	var resphtml string
	httpClient := &http.Client{
		CheckRedirect: nil,
		Jar:           gCurCookiesJar,
	}
	httpreq, err := http.NewRequest("GET", url, nil)
	httpresp, err := httpClient.Do(httpreq)
	if err != nil {
		gLog.Println(err)
	}
	defer httpresp.Body.Close()
	body, errReadAll := ioutil.ReadAll(httpresp.Body)
	if errReadAll != nil {
		gLog.Println("ioutil.ReadAll(httpresp.Body) err=", err)
	}
	gCurCookies = gCurCookiesJar.Cookies(httpreq.URL)
	resphtml = string(body)
	return resphtml
}

func dbgPrinCurCookies() {
	var cookieNum = len(gCurCookies)
	for i := 0; i < cookieNum; i++ {
		curCk := gCurCookies[i]
		gLog.Printf("========Cookie[%d]", i)
		gLog.Printf("Name\t=%s", curCk.Name)
		gLog.Printf("Value\t=%s", curCk.Value)
		gLog.Printf("Path\t=%s", curCk.Path)
		gLog.Printf("Domain\t=%s", curCk.Domain)
		gLog.Printf("Expires\t=%s", curCk.Expires)
		gLog.Printf("RawExpires=%s", curCk.RawExpires)
		gLog.Printf("MaxAge\t=%d", curCk.MaxAge)
		gLog.Printf("Secure\t=%t", curCk.Secure)
		gLog.Printf("HttpOnly=%t", curCk.HttpOnly)
		gLog.Printf("Raw\t=%s", curCk.Raw)
		gLog.Printf("Unparsed=%s", curCk.Unparsed)
	}
}
func main() {
	logFile, gLogger := newLogger()
	gLog = gLogger
	defer logFile.Close()
	fmt.Println(getRespHtml("http://www.baidu.com"))
	dbgPrinCurCookies()
	gLog.Println("hello world")

}
