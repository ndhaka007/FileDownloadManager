package Model

import(
	"crypto/rand"
	"fmt"
	"log"
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

func (s SerialDownload)DownloadFile(){

}
func (c ConDownload)DownloadFile(){

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