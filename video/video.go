package video

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleVideo(router *mux.Router) {
	router.HandleFunc("/video/{sessionID}/{memberID}", OptionsHandler).Methods("OPTIONS")

	router.HandleFunc("/video/{sessionID}/{memberID}", postVideo).Methods("POST")
	router.HandleFunc("/video/{sessionID}/{memberID}", getVideo).Methods("GET")
}

type Video struct {
	Data string `json:"video"`
}

var videoMap = make(map[string]map[string]Video)

func postVideo(w http.ResponseWriter, r *http.Request) {
	// find sessionID and memberID as urlparams
	urlParams := mux.Vars(r)
	sessionID := urlParams["sessionID"]
	memberID := urlParams["memberID"]

	// decode a json string as videodata
	var videoData Video
	if err := json.NewDecoder(r.Body).Decode(&videoData); err != nil {
		errstring := "Invalid input" + videoData.Data
		log.Println(w, errstring, http.StatusBadRequest)
		return
	}

	// if this is the first frame sent to this session, initialize it's map
	if videoMap[sessionID] == nil {
		videoMap[sessionID] = make(map[string]Video)
	}

	// add the videodata frame
	videoMap[sessionID][memberID] = videoData

	w.WriteHeader(http.StatusCreated)
} 

func getVideo(w http.ResponseWriter, r *http.Request) {	
	urlParams := mux.Vars(r)
	sessionID := urlParams["sessionID"]
	memberID := urlParams["memberID"]
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videoMap[sessionID][memberID])
}

// non-endpoint related funcs

func clearVideoData(hostID string, memberID string) {
	session, sessionExists := videoMap[hostID]
	if !sessionExists {
		log.Println("Looks like the session doesn't exist in videoMap")
		return 
	}

	memberVideo, memberExists := session[memberID]
	if !memberExists {
		log.Println("Looks like the member doesn't exist in videoMap")
		return
	}

	memberVideo.Data = ""
	session[memberID] = memberVideo
	videoMap[hostID] = session
	log.Println(videoMap[hostID][memberID].Data)
}