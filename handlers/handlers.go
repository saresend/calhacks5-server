package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/saresend/calhacks5-server/state"
)

var upgrader = websocket.Upgrader{}

var connections = make(map[string]*websocket.Conn)

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	println("Attempting new socket connection...")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println(err)
		return
	}
	user := new(state.User)

	if err = conn.ReadJSON(user); err != nil {
		println(err)
		return
	}
	connections[user.GetUsername()] = conn
	user.AddUser()

}

type UserUpdate struct {
	Username       string `json:"username"`
	SelectedPrompt int    `json:"selectedPrompt"`
}

type Update struct {
	CurrentUser     UserUpdate `json:"currentUser"`
	NextUser        UserUpdate `json:"nextUser"`
	TerminationTime int64      `json:"terminationTime"`
	Upvotes         int        `json:"upvotes"`
	Downvotes       int        `json:"downvotes"`
}

func Broadcast() {
	newState := state.GetState()
	var newCurrentUser UserUpdate
	if newState.CurrentUser != nil {
		newCurrentUser = UserUpdate{newState.CurrentUser.Username, newState.CurrentUser.SelectedPrompt}
	}
	var newNextUser UserUpdate
	if newState.NextUser != nil {
		newNextUser = UserUpdate{newState.NextUser.Username, newState.NextUser.SelectedPrompt}
	}
	update := Update{newCurrentUser, newNextUser,
		newState.TerminationTime, newState.Upvotes, newState.Downvotes}
	// b, _ := json.Marshal(newState)
	// os.Stdout.Write(b)
	// println()
	for id, conn := range connections {
		if err := conn.WriteJSON(update); err != nil {
			delete(connections, id)
			state.RemoveUser(id)
		}
	}
}

func GetPrompts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(state.GetPrompts())
}

type SetPromptsRequest struct {
	Username string
	Prompts  []int
}

func SetPrompts(w http.ResponseWriter, r *http.Request) {
	request := new(SetPromptsRequest)

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return
	}
	state.SetPrompts(request.Username, request.Prompts)
}

type VoteRequest struct {
	Username    string
	CurrentUser string
	upvote      bool
}

func Vote(w http.ResponseWriter, r *http.Request) {
	request := new(VoteRequest)

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return
	}
	state.MakeVote(request.upvote)
}
