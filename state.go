package main

import (
	"math/rand"
	"time"
)

// DB struct
type User struct {
	username       string
	prompts        []int
	selectedPrompt int
}

type Prompt struct {
	id          int
	description string
}

type PublicState struct {
	currentUser     *User
	terminationTime int64
	nextUser        *User
}

var publicState PublicState
var users map[string]*User
var activeUsers map[string]*User
var streamingUsers map[string]*User

func Init() {
	publicState = PublicState{nil, time.Now().Unix(), nil}
	users = make(map[string]*User)
}

func (user *User) AddUser() {
	user.prompts = make([]int, 0)
	user.selectedPrompt = -1
	users[user.username] = user
}

func (user *User) GetUsername() string {
	return user.username
}

// RemoveUser removes user from all queuess
func RemoveUser(username string) {
	delete(activeUsers, username)
	delete(streamingUsers, username)
	if publicState.nextUser == users[username] {
		publicState.nextUser = nil
	}
}

// SetPrompts chooses which prompts a user wants to broadcast on
func SetPrompts(username string, prompts []int) {
	users[username].prompts = prompts
	// remove user from queue if they deselect all prompts
	if len(prompts) == 0 {
		delete(streamingUsers, username)
		if publicState.nextUser == users[username] {
			publicState.nextUser = nil
		}
	} else {
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
	println("hello")
	// if time on screen expired
	if publicState.terminationTime < time.Now().Unix() {
		if publicState.nextUser == nil {
			publicState.nextUser = publicState.currentUser
		}
		publicState.currentUser = publicState.nextUser
		publicState.nextUser = nil
		publicState.terminationTime = time.Now().Unix() + 15
	}
	// select next user
	if publicState.nextUser == nil {
		for k, v := range streamingUsers {
			if activeUsers[k] != nil && len(v.prompts) > 0 {
				publicState.nextUser = v
				v.selectedPrompt = v.prompts[rand.Int()%len(v.prompts)]
				delete(streamingUsers, k)
				break
			}
		}
	}

	// send new state to all
	Broadcast(publicState)
}
