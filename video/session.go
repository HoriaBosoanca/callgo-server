package video

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

// so that it is encoded as json object
type Name struct {
	Name string `json:"name"`
}

type Participant struct {
	Name string `json:"name"`
	ID string `json:"ID"`
}

type Session struct {
	Host Participant `json:"host"`
	Members []Participant `json:"members"`
}

var sessionsMap map[string]Session = make(map[string]Session)

func HandleSession(router *mux.Router) {
	router.HandleFunc("/session", OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/session/{id}", OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/session/{hostID}/{memberID}", OptionsHandler).Methods("OPTIONS")
	
	router.HandleFunc("/session", createSession).Methods("POST")
	router.HandleFunc("/session/{id}", addMember).Methods("POST")
	router.HandleFunc("/session/{id}", getSessionByID).Methods("GET")
	router.HandleFunc("/session/{hostID}/{memberID}", leaveSession).Methods("DELETE")
}

func createSession(w http.ResponseWriter, r *http.Request) {
	hostID := uuid.New().String()
	var name Name
	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		http.Error(w, "Invalid name", http.StatusBadRequest)
	}
	sessionsMap[hostID] = Session{Host: Participant{Name: name.Name, ID: hostID}}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(hostID)
}

func getSessionByID(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	hostID := urlParams["id"]
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sessionsMap[hostID])
}

func addMember(w http.ResponseWriter, r* http.Request) {
	var name Name
	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		http.Error(w, "Invalid name", http.StatusBadRequest)
	}

	urlParams := mux.Vars(r)
	hostID := urlParams["id"]

	session, exists := sessionsMap[hostID]
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
        return
	}
	memberID := uuid.New().String()
	session.Members = append(sessionsMap[hostID].Members, Participant{Name: name.Name, ID: memberID})
	sessionsMap[hostID] = session

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(memberID)
}

func leaveSession(w http.ResponseWriter, r *http.Request) {
	hostID := mux.Vars(r)["hostID"]
	memberID := mux.Vars(r)["memberID"]

	session, exists := sessionsMap[hostID]
	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	memberIndex := -1
	for i, member := range session.Members {
		if member.ID == memberID {
			memberIndex = i
			break
		}
	}

	if memberIndex == -1 {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	session.Members = append(session.Members[:memberIndex], session.Members[memberIndex+1:]...)
	sessionsMap[hostID] = session

	clearVideoData(hostID, memberID)

	w.WriteHeader(http.StatusNoContent)
}
