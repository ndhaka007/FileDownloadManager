package Model

import(
	"time"
)

type DownloadFile struct{
	Type string   `json:"type"`
	Urls []string `json:"urls"`
}

type Response struct{
	Id string `json:"id"`
	Start_time Time `json:start_time`
	End_time Time `json:end_time`
	Status string `json:"status"`
	Download_type string `json:"SERIAL"`
	Files map[string]string
}
