package state

import (
	"math/rand"
	"time"
)

// DB struct
type User struct {
	Username       string
	prompts        []int
	SelectedPrompt int
}

type Prompt struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type PublicState struct {
	CurrentUser     *User
	TerminationTime int64
	NextUser        *User
}

var publicState PublicState
var users = make(map[string]*User)
var activeUsers = make(map[string]*User)
var streamingUsers = make(map[string]*User)

func Init() {
	publicState = PublicState{nil, time.Now().Unix(), nil}
}

func (user *User) AddUser() {
	user.prompts = make([]int, 0)
	user.SelectedPrompt = -1
	users[user.Username] = user
	activeUsers[user.Username] = user
	println("Added " + user.Username)
}

func (user *User) GetUsername() string {
	return user.Username
}

// RemoveUser removes user from all queuess
func RemoveUser(username string) {
	delete(activeUsers, username)
	delete(streamingUsers, username)
	if publicState.NextUser == users[username] {
		publicState.NextUser = nil
	}
}

// SetPrompts chooses which prompts a user wants to broadcast on
func SetPrompts(username string, prompts []int) {
	if users[username] == nil {
		println("Setting prompts for invalid user")
		return
	}
	users[username].prompts = prompts
	// remove user from queue if they deselect all prompts
	if len(prompts) == 0 {
		delete(streamingUsers, username)
		if publicState.NextUser == users[username] {
			publicState.NextUser = nil
		}
	} else {
		println("Adding to streaming users")
		streamingUsers[username] = users[username]
	}
}

// GetPrompts returns available prompts
func GetPrompts() []Prompt {
	prompts := []Prompt{
		{0, "Dare 1"},
		{1, "Dare 2"},
		{2, "Dare 3"},
		{3, "Dare 4"},
		{4, "Dare 5"},
		{5, "Dare 6"},
		{6, "Dare 7"},
		{7, "Dare 8"},
		{8, "Dare 9"}}
	return prompts
}

// UpdateState updates state every second (triggered by tick) and updates the current and next user if needed
// TODO: add/subtract time based on votes
func UpdateState() {
	// if time on screen expired
	if publicState.TerminationTime < time.Now().Unix() {
		println("Going to next user")
		if publicState.NextUser == nil {
			publicState.NextUser = publicState.CurrentUser
		}
		publicState.CurrentUser = publicState.NextUser
		publicState.NextUser = nil
		publicState.TerminationTime = time.Now().Unix() + 15
	}
	// if there is no one on screen right now
	if publicState.CurrentUser == nil && len(streamingUsers) > 0 {
		println("Trying to choose current user...")
		for k, v := range streamingUsers {
			println(k + "?")
			if activeUsers[k] != nil && len(v.prompts) > 0 {
				publicState.CurrentUser = v
				v.SelectedPrompt = v.prompts[rand.Int()%len(v.prompts)]
				delete(streamingUsers, k)
				break
			}
		}
	}
	// select next user
	if publicState.NextUser == nil && len(streamingUsers) > 0 {
		println("Trying to choose next user...")
		for k, v := range streamingUsers {
			println(k + "?")
			if activeUsers[k] != nil && len(v.prompts) > 0 {
				publicState.NextUser = v
				v.SelectedPrompt = v.prompts[rand.Int()%len(v.prompts)]
				delete(streamingUsers, k)
				break
			}
		}
	}
}

func GetState() PublicState {
	return publicState
}
