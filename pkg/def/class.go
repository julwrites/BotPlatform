package def

// Struct definitions for bot

type UserData struct {
	Firstname string
	Lastname  string
	Username  string
	Id        string
	Type      string // Group/Individual
	Action    string // Current action if any
	Config    string
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
	Type    string
	Channel string
	User    UserData
	Msg     MessageData
	Res     ResponseData
	Props   map[string]interface{}
}
