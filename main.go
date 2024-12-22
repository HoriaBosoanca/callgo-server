package main

import (
	"log"
	"os"
	"net/http"

	"Callgo/video"
	
	"github.com/gorilla/mux"
)

func main() {
	// init
	router := mux.NewRouter()
	router.Use(video.EnableCORS)
	
	// get port
	port := os.Getenv("CALLGOPORT")
    if port == "" {
        port = "8080"
    }

	// Handle endpoints
	video.HandleVideo(router)
	video.HandleSession(router)

	certFile := "/etc/letsencrypt/live/horia.live/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/horia.live/privkey.pem"

	// log and start
	log.Printf("Server starting on port %s.", port)
	log.Fatal(http.ListenAndServeTLS(":"+port, certFile, keyFile, router))
}
