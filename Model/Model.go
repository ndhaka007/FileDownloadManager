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
	DownloadFile()error
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

func (s SerialDownload)DownloadFile(urlVsAdd map[string]string,uuid string)error{
	_ = os.MkdirAll("/tmp/"+uuid,os.ModePerm)
	for _,link := range s.Urls{
		if _, ok := urlVsAdd[link]; !ok {
			j := generateUUID()
			// don't worry about errors
			response, e := http.Get(link)
			if e != nil {
				fmt.Println(e)
				return e
			}
			defer response.Body.Close()

			//open a file for writing
			file, err := os.Create("/tmp/" + uuid + "/" + j)
			if err != nil {
				//log.Fatal(err)
				return err
			}
			defer file.Close()

			// Use io.Copy to just dump the response body to the file. This supports huge files
			_, err = io.Copy(file, response.Body)
			if err != nil {
				//log.Fatal(err)
				return err
			}
			fmt.Println("Success!")
			urlVsAdd[link] = "/tmp/" + uuid + "/" + j
		}
	}
	return nil
}


func (c ConDownload)DownloadFile(urlVsAdd map[string]string,uuid string, resp *Response)error{
	concurrency := 6
	req := make(chan string, concurrency)
	set := make(map[string]bool)
	for _, i := range c.Urls{
		set[i] = true
	}
	for i:=0;i<concurrency;i++{
		go concurrent(req,len(set), 0, urlVsAdd ,uuid, resp)
	}
	go func() {
		for link := range set {
				req <- link
		}
		fmt.Println("z")
	}()
	return nil
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

func condown(site string, urlVsAdd map[string]string, uuid string){
	j:= generateUUID()
	// don't worry about errors
	response, e := http.Get(site)
	if e != nil {
		fmt.Println(e)
	}
	defer response.Body.Close()
	_ = os.MkdirAll("/tmp/"+uuid,os.ModePerm)
	//open a file for writing
	file, err := os.Create("/tmp/"+uuid+"/"+j)
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
	urlVsAdd[site]= "/tmp/"+uuid+"/"+j
}

func concurrent(reqChan chan string,total int, comp int,urlVsAdd map[string]string,uuid string, resp *Response) {
	for {
		select {
		case link, ok := <-reqChan:
			if !ok{
				return
			}
			condown(link, urlVsAdd, uuid)
		}
		if len(urlVsAdd)==total {
			close(reqChan)
			resp.Status = "successful"
			resp.EndTime = time.Now()
			return
		}
	}
}