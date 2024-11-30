package video

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Session struct {
	Host int `json:"host"`
	Members []int `json:"members"`
}

var sessionsMap map[int]Session = make(map[int]Session)

func HandleSession(router *mux.Router) {
	router.HandleFunc("/session", OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/session/{id}", OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/session", createSession).Methods("POST")
	router.HandleFunc("/session/{id}", getSessionByID).Methods("GET")
	
	router.HandleFunc("/member", OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/member/{id}", OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/member/{id}", addMember).Methods("POST")
}

var maxID = 0
func makeID() int {
	maxID++
	return maxID
}

func createSession(w http.ResponseWriter, r *http.Request) {
	hostID := makeID()
	sessionsMap[hostID] = Session{Host: hostID}

	for _, v := range(sessionsMap) {
		fmt.Println("Host: ", v.Host, "Members: ", v.Members)
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(hostID)
}

func getSessionByID(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	hostID, err := strconv.Atoi(urlParams["id"])
	if err != nil {
		http.Error(w, "Wrong url param", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sessionsMap[hostID])
}

func addMember(w http.ResponseWriter, r* http.Request) {
	urlParams := mux.Vars(r)
	hostID, err := strconv.Atoi(urlParams["id"])
	if err != nil {
		http.Error(w, "Wrong url param", http.StatusBadRequest)
		return
	}

	session, exists := sessionsMap[hostID]
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
        return
	}
	memberID := makeID()
	session.Members = append(sessionsMap[hostID].Members, memberID)
	sessionsMap[hostID] = session

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(memberID)
}
