package Controller

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/ndhaka007/FileDownloadManager/Model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	i = 0
	mp = make(map[string] Model.DownloadFile)
)

func down(site string, j int){
	// don't worry about errors
	response, e := http.Get(site)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create("/tmp/"+strconv.Itoa(j)+".jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Success!")
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

func HomePage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Ok","200")
	fmt.Fprintf(w,"Ok")
	fmt.Println("Endpoint Hit: homePage")
}


func Download(w http.ResponseWriter, r *http.Request) {
	//to generate uuid
	uuid := generateUUID()
	resp := Model.Response{uuid}
	w.Header().Set("Content-Type", "application/json")
	ret,_:=json.Marshal(resp)

	//read payload
	reqBody, _ := ioutil.ReadAll(r.Body)
	//convert to json
	var file Model.DownloadFile
	_ = json.Unmarshal(reqBody, &file)
	if file.Type == "serial" {
		for _, link := range file.Urls {
			down(link, i)
			i++
		}
		w.Write(ret)
	}	else{
		w.Write(ret)

		for _, link := range file.Urls {
			down(link, i)
			i++
		}
	}
	mp[uuid] = file

}

func parseURL(url string) string{
	var s = ""
	for i:=len(url)-1;url[i]!='/';i--{
		s=string(url[i])+s
	}
	return s
}

func Status(w http.ResponseWriter, r *http.Request){
	id := parseURL(r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	ret,_:=json.Marshal(mp[id])
	w.Write(ret)

}