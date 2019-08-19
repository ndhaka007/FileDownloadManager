package Routes

import(
	"github.com/ndhaka007/FileDownloadManager/Controller"
	"log"
	"net/http"
)

func HandleRequests() {
	//to check health status
	http.HandleFunc("/health", Controller.HomePage)
	//to download images
	http.HandleFunc("/downloads", Controller.Download)
	//to check status
	http.HandleFunc("/downloads/", Controller.Status)

	log.Fatal(http.ListenAndServe(":8081", nil))
}
