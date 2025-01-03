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
        port = "1234"
    }

	// Handle endpoints
	video.HandleVideo(router)
	video.HandleSession(router)

	// log and start
	log.Printf("Server starting on port %s.", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
