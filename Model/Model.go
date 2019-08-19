package Model

import(
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Download interface{
	DownloadFile()
}

type SerialDownload struct{
	Urls []string
}

type ConDownload struct{
	Urls []string
}

type Input struct{
	Type string   `json:"type"`
	Urls []string `json:"urls"`
}

type Response struct{
	Id string `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime time.Time `json:"end_time"`
	Status string `json:"status"`
	DownloadType string `json:"downloadType"`
	Files map[string]string
}

func (s SerialDownload)DownloadFile(urlVsAdd map[string]string,uuid string){
	for _,link := range s.Urls{
		down(link,urlVsAdd ,uuid)
	}
}
func (c ConDownload)DownloadFile(urlVsAdd map[string]string,uuid string){
	for _,link := range c.Urls{
		down(link,urlVsAdd ,uuid)
	}
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

func down(site string, urlVsAdd map[string]string, uuid string){
	j:= generateUUID()
	// don't worry about errors
	response, e := http.Get(site)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()
	_ = os.MkdirAll("/tmp/"+uuid,os.ModePerm)
	//open a file for writing
	file, err := os.Create("/tmp/"+uuid+"/"+j+".jpg")
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
	urlVsAdd[site]= "/tmp/"+uuid+"/"+j+".jpg"
}