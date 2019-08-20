package Controller

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/ndhaka007/FileDownloadManager/Model"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	Mp = make(map[string]*Model.Response)
)

//API to check health
func HomePage(w http.ResponseWriter, r *http.Request){
	if r.Method =="GET" {
		w.Header().Set("Ok", "200")
		_, _ = fmt.Fprintf(w, "Ok")
	}
}

//API to download files
func Download(w http.ResponseWriter, r *http.Request) {
	if r.Method =="POST" {
		//to generate uuid for request
		uuid := generateUUID()
		//map to store url and there respective downloaded images
		urlVsAdd := make(map[string]string)
		//request time
		startTime := time.Now()
		//read payload
		reqBody, _ := ioutil.ReadAll(r.Body)
		//convert to json
		var file Model.Input
		_ = json.Unmarshal(reqBody, &file)
		if file.Type == "serial" {
			//serial download
			se := Model.SerialDownload{Urls: file.Urls}
			e := se.DownloadFile(urlVsAdd, uuid)

			//save the data in response and map and return the id
			resp := &Model.Response{Id: uuid, StartTime: startTime, EndTime: time.Now(), Status: "successful", DownloadType: file.Type, Files: urlVsAdd}
			if e != nil {
				resp.Status = "failed"
			}
			w.Header().Set("Content-Type", "application/json")
			ret, _ := json.Marshal(resp.Id)
			Mp[uuid] = resp
			_, _ = w.Write(ret)
		} else {
			//concurrent download
			//save the data in response and map and return the id
			resp := &Model.Response{Id: uuid, StartTime: startTime, EndTime: time.Now(), Status: "queue", DownloadType: file.Type, Files: urlVsAdd}
			w.Header().Set("Content-Type", "application/json")
			ret, _ := json.Marshal(resp.Id)
			Mp[uuid] = resp
			_, _ = w.Write(ret)

			co := Model.ConDownload{Urls: file.Urls}
			e := co.DownloadFile(urlVsAdd, uuid, resp)
			if e != nil {
				resp.Status = "failed"
			}
			Mp[uuid] = resp
		}
	}
}

//API to check status
func Status(w http.ResponseWriter, r *http.Request){
	if r.Method =="GET" {
		id := parseURL(r.URL.Path)
		if _, va := Mp[id]; !va {
			w.Header().Set("id not found", "404")
			_, _ = fmt.Fprintf(w, "404 Id not Found")
		} else {
			w.Header().Set("Content-Type", "application/json")
			ret, _ := json.Marshal(Mp[id])
			_, _ = w.Write(ret)
		}
	}
}

//function to generate UUID
func generateUUID() string{
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

//function to parse URL and extract id for status
func parseURL(url string) string{
	var s = ""
	for i:=len(url)-1;url[i]!='/';i--{
		s=string(url[i])+s
	}
	return s
}