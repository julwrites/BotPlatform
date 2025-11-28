package def

import "github.com/julwrites/BotPlatform/pkg/secrets"

// Struct definitions for bot

type UserData struct {
	Firstname string `datastore:""`
	Lastname  string `datastore:""`
	Username  string `datastore:""`
	Id        string `datastore:""`
	Type      string `datastore:""` // Group/Individual
	Action    string `datastore:""` // Current action if any
	Config    string `datastore:""`
}

type MessageData struct {
	Id      string
	Chat    string
	Command string
	Message string
}

type Option struct {
	Text string
	Link string
}

type ResponseOptions struct {
	Inline  bool
	Options []Option
	Remove  bool
}

type ResponseData struct {
	Message     string
	ParseMode   string
	Affordances ResponseOptions
}

type SessionData struct {
	Secrets      secrets.SecretsData
	Type         string
	Channel      string
	User         UserData
	Msg          MessageData
	Res          ResponseData
	ResourcePath string
}
