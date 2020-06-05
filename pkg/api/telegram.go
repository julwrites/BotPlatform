package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julwrites/BotMultiplexer/pkg/def"
)

// Classes

type TelegramSender struct {
	Id        int    `json:"id"`
	Bot       bool   `json:"is_bot"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Username  string `json:"username"`
	Language  string `json:"langauge_code"`
}

type TelegramChat struct {
	Id        int    `json:"id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type TelegramMessage struct {
	Sender TelegramSender `json:"from"`
	Chat   TelegramChat   `json:"chat"`
	Text   string         `json:"text"`
	Id     int            `json:"message_id"`
}

type TelegramRequest struct {
	Message TelegramMessage `json:"message"`
}

type TelegramPost struct {
	Id        string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
	ReplyId   string `json:"reply_to_message_id"`
}

type InlineButton struct {
	Text string `json:"text"`
	Url  string `json:"url"`
}

type InlineMarkup struct {
	Keyboard [][]InlineButton `json:"inline_keyboard"`
}

type TelegramInlinePost struct {
	TelegramPost
	Markup InlineMarkup `json:"reply_markup"`
}

type KeyButton struct {
	Text string `json:"text"`
}

type ReplyMarkup struct {
	Keyboard  [][]KeyButton `json:"keyboard"`
	Resize    bool          `json:"resize_keyboard"`
	Once      bool          `json:"one_time_keyboard"`
	Selective bool          `json:"selective"`
}

type TelegramReplyPost struct {
	TelegramPost
	Markup ReplyMarkup `json:"reply_markup"`
}

type RemoveMarkup struct {
	Remove    bool `json:"remove_keyboard"`
	Selective bool `json:"selective"`
}

type TelegramRemovePost struct {
	TelegramPost
	Markup RemoveMarkup `json:"reply_markup"`
}

// Translate to properties

func TelegramTranslate(body []byte) def.SessionData {
	log.Printf("Parsing Telegram message")

	var env def.SessionData
	env.Type = def.TYPE_TELEGRAM

	var data TelegramRequest
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Failed to unmarshal request body: %v", err)
		return env
	}

	env.User.Firstname = data.Message.Sender.Firstname
	env.User.Lastname = data.Message.Sender.Lastname
	env.User.Username = data.Message.Sender.Username
	env.User.Id = strconv.Itoa(data.Message.Sender.Id)

	// TODO: Implement support for groups
	// env.User.Type = data.Message.Sender.Title

	log.Printf("User: %s %s | %s : %s", env.User.Firstname, env.User.Lastname, env.User.Username, env.User.Id)

	tokens := strings.Split(data.Message.Text, " ")
	if strings.Index(tokens[0], "/") == 0 {
		env.Msg.Command = string((tokens[0])[1:])                                // Get the first token and strip off the prefix
		data.Message.Text = strings.Replace(data.Message.Text, tokens[0], "", 1) // Replace the command
	}
	env.Msg.Message = data.Message.Text
	env.Msg.Id = strconv.Itoa(data.Message.Id)

	env.Channel = strconv.Itoa(data.Message.Chat.Id)

	log.Printf("Message: %s | %s", env.Msg.Command, env.Msg.Message)

	return env
}

// Translate to Telegram

func HasOptions(env def.SessionData) bool {
	return len(env.Res.Affordances.Options) > 0 || env.Res.Affordances.Remove
}

func PrepTelegramMessage(base TelegramPost, env def.SessionData) []byte {
	var data []byte
	var jsonErr error

	if HasOptions(env) {
		if env.Res.Affordances.Remove {
			var message TelegramRemovePost
			message.TelegramPost = base
			message.Markup.Remove = true
			message.Markup.Selective = true
			data, jsonErr = json.Marshal(message)
			log.Printf("Post with Affordance Removal command")
		} else if len(env.Res.Affordances.Options) > 0 {
			if env.Res.Affordances.Inline {
				var buttons []InlineButton
				for i := 0; i < len(env.Res.Affordances.Options); i++ {
					buttons = append(buttons, InlineButton{env.Res.Affordances.Options[i].Text, env.Res.Affordances.Options[i].Link})
				}
				var markup InlineMarkup
				markup.Keyboard = append([][]InlineButton{}, buttons)
				var message TelegramInlinePost
				message.TelegramPost = base
				message.Markup = markup
				data, jsonErr = json.Marshal(message)
				log.Printf("Post with Inline Affordance command")
			} else {
				var buttons []KeyButton
				for i := 0; i < len(env.Res.Affordances.Options); i++ {
					buttons = append(buttons, KeyButton{env.Res.Affordances.Options[i].Text})
				}
				var markup ReplyMarkup
				markup.Keyboard = append([][]KeyButton{}, buttons)
				var message TelegramReplyPost
				message.TelegramPost = base
				message.Markup = markup
				data, jsonErr = json.Marshal(message)
				log.Printf("Post with Keyboard Affordance command")
			}
		}
	} else {
		data, jsonErr = json.Marshal(base)
	}

	if jsonErr != nil {
		log.Printf("Error occurred during conversion to JSON: %v", jsonErr)
		return nil
	}

	return data
}

func PostTelegram(env def.SessionData) bool {
	endpoint := "https://api.telegram.org/bot" + env.Secrets.TELEGRAM_ID + "/sendMessage"
	header := "application/json;charset=utf-8"

	text := Format(env.Res.Message, TelegramNormal, TelegramBold, TelegramItalics, TelegramSuperscript)

	chunks := Split(text, 4000)

	var base TelegramPost
	base.Id = env.User.Id
	base.ParseMode = def.TELEGRAM_PARSE_MODE
	base.ReplyId = env.Msg.Id

	for _, chunk := range chunks {
		base.Text = chunk
		data := PrepTelegramMessage(base, env)

		buffer := bytes.NewBuffer(data)
		res, postErr := http.Post(endpoint, header, buffer)
		if postErr != nil {
			log.Printf("Error occurred during post: %v", postErr)
			return false
		}

		log.Printf("Posted message, response %v", res)
	}

	return true
}

// Formatting methods

func TelegramNormal(str string) string {
	str = strings.ReplaceAll(str, "[", "\\[")
	str = strings.ReplaceAll(str, "]", "\\]")
	str = strings.ReplaceAll(str, "(", "\\(")
	str = strings.ReplaceAll(str, ")", "\\)")
	str = strings.ReplaceAll(str, "~", "\\~")
	str = strings.ReplaceAll(str, ">", "\\>")
	str = strings.ReplaceAll(str, "#", "\\#")
	str = strings.ReplaceAll(str, "+", "\\+")
	str = strings.ReplaceAll(str, "-", "\\-")
	str = strings.ReplaceAll(str, "=", "\\=")
	str = strings.ReplaceAll(str, "|", "\\|")
	str = strings.ReplaceAll(str, "{", "\\{")
	str = strings.ReplaceAll(str, "}", "\\}")
	str = strings.ReplaceAll(str, ".", "\\.")
	str = strings.ReplaceAll(str, "!", "\\!")

	return str
}

func TelegramBold(str string) string {
	return fmt.Sprintf("*%s*", str)
}

func TelegramItalics(str string) string {
	return fmt.Sprintf("_%s_", str)
}

func TelegramSuperscript(str string) string {
	var out string

	for _, c := range str {
		switch string(c) {
		case "0":
			out = out + "\u2070"
			break
		case "1":
			out = out + "\u00b9"
			break
		case "2":
			out = out + "\u00b2"
			break
		case "3":
			out = out + "\u00b3"
			break
		case "4":
			out = out + "\u2074"
			break
		case "5":
			out = out + "\u2075"
			break
		case "6":
			out = out + "\u2076"
			break
		case "7":
			out = out + "\u2077"
			break
		case "8":
			out = out + "\u2078"
			break
		case "9":
			out = out + "\u2079"
			break
		}
	}

	return out
}