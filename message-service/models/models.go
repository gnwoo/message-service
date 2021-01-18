package models

type Message struct {
	Sender string
	Channel string // Message channel or control channel
	Session string // A session could be a 1-1 conversation and a group chat
	Body string // Message body
	DispatchTime int
}
