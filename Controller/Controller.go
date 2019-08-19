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
	Mp = make(map[string] Model.Response)
)

func HomePage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Ok","200")
	fmt.Fprintf(w,"Ok")
	fmt.Println("Endpoint Hit: homePage")
}

func Download(w http.ResponseWriter, r *http.Request) {
	//to generate uuid
	uuid := generateUUID()

	urlVsAdd := make(map[string]string)
	t :=time.Now()
	//read payload
	reqBody, _ := ioutil.ReadAll(r.Body)
	//convert to json
	var file Model.Input
	_ = json.Unmarshal(reqBody, &file)
	if file.Type == "serial" {
		se := Model.SerialDownload{Urls: file.Urls}

		se.DownloadFile(urlVsAdd,uuid)

		resp := Model.Response{Id: uuid, StartTime: t, EndTime: time.Now(), Status: "successful", DownloadType: file.Type, Files: urlVsAdd}
		w.Header().Set("Content-Type", "application/json")
		ret,_:=json.Marshal(resp.Id)
		Mp[uuid] = resp
		w.Write(ret)
	}	else{
		resp := Model.Response{Id: uuid, StartTime: t, EndTime: time.Now(), Status: "queue", DownloadType: file.Type, Files: urlVsAdd}
		w.Header().Set("Content-Type", "application/json")
		ret,_:=json.Marshal(resp.Id)
		Mp[uuid] = resp
		w.Write(ret)
		 co := Model.ConDownload{Urls: file.Urls}

		 co.DownloadFile(urlVsAdd,uuid)

		resp.Status = "successful"
		Mp[uuid] = resp
	}
}

func Status(w http.ResponseWriter, r *http.Request){
	id := parseURL(r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	ret,_:=json.Marshal(Mp[id])
	w.Write(ret)
}

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

func parseURL(url string) string{
	var s = ""
	for i:=len(url)-1;url[i]!='/';i--{
		s=string(url[i])+s
	}
	return s
}