package video

import (
	"encoding/json"
	"net/http"
	"strconv"

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

var videoMap = make(map[int]map[int]Video)

func postVideo(w http.ResponseWriter, r *http.Request) {
	// find sessionID and memberID as urlparams
	urlParams := mux.Vars(r)
	sessionID, err := strconv.Atoi(urlParams["sessionID"])
	if err != nil {
		http.Error(w, "Wrong url param", http.StatusBadRequest)
		return
	}
	memberID, err := strconv.Atoi(urlParams["memberID"])
	if err != nil {
		http.Error(w, "Wrong url param", http.StatusBadRequest)
		return
	}

	// decode a json string as videodata
	var videoData Video
	if err := json.NewDecoder(r.Body).Decode(&videoData); err != nil {
		errstring := "Invalid input" + videoData.Data
		http.Error(w, errstring, http.StatusBadRequest)
		return
	}

	// if this is the first frame sent to this session, initialize it's map
	if videoMap[sessionID] == nil {
		videoMap[sessionID] = make(map[int]Video)
	}

	// add the videodata frame
	videoMap[sessionID][memberID] = videoData

	w.WriteHeader(http.StatusCreated)
} 

func getVideo(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	sessionID, err := strconv.Atoi(urlParams["sessionID"])
	if err != nil {
		http.Error(w, "Wrong url param", http.StatusBadRequest)
		return
	}
	memberID, err := strconv.Atoi(urlParams["memberID"])
	if err != nil {
		http.Error(w, "Wrong url param", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videoMap[sessionID][memberID])
}
