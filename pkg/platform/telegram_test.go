package platform

import (
	"encoding/json"
	"testing"

	"github.com/julwrites/BotPlatform/pkg/def"
)

func GenerateTestData() []byte {
	var req TelegramRequest
	req.Message.Id = 9876
	req.Message.Chat.Id = 4567
	req.Message.Text = "/Command Text"
	req.Message.Sender.Id = 1234
	req.Message.Sender.Firstname = "First"
	req.Message.Sender.Lastname = "Last"
	req.Message.Sender.Username = "User"

	data, _ := json.Marshal(req)

	return data
}

func TestTelegramTranslate(t *testing.T) {
	data := GenerateTestData()

	// Create a telegram instance (secrets don't matter for translate)
	tg := NewTelegram("token", "admin")
	env, err := tg.Translate(data)

	if err != nil {
		t.Errorf("Failed TestTelegramTranslate: %v", err)
	}

	if env.Type != def.TYPE_TELEGRAM {
		t.Errorf("Failed TestTelegramTranslate, Type is wrong")
	}
	if env.Msg.Command != "Command" {
		t.Errorf("Failed TestTelegramTranslate, Msg Command is wrong")
	}
	if env.Msg.Message != "Text" {
		t.Errorf("Failed TestTelegramTranslate, Msg Text is wrong")
	}
	if env.Msg.Id != "9876" {
		t.Errorf("Failed TestTelegramTranslate, Msg ID is wrong")
	}
	if env.Channel != "4567" {
		t.Errorf("Failed TestTelegramTranslate, Channel ID is wrong")
	}
	if env.User.Id != "1234" {
		t.Errorf("Failed TestTelegramTranslate, User ID is wrong")
	}
	if env.User.Firstname != "First" {
		t.Errorf("Failed TestTelegramTranslate, User First name is wrong")
	}
	if env.User.Lastname != "Last" {
		t.Errorf("Failed TestTelegramTranslate, User Last name is wrong")
	}
	if env.User.Username != "User" {
		t.Errorf("Failed TestTelegramTranslate, User Username is wrong")
	}
}

func TestPrepTelegramMessage(t *testing.T) {
	var post TelegramPost

	post.Text = "Text"
	post.ParseMode = def.TELEGRAM_PARSE_MODE_MD
	post.Id = "1234"
	post.ReplyId = "4567"

	{
		var env def.SessionData
		env.Res.ParseMode = def.TELEGRAM_PARSE_MODE_HTML
		if env.Res.ParseMode != "HTML" {
			t.Errorf("Failed TestPrepTelegramMessage ParseMode check")
		}

		postHTML := post
		postHTML.ParseMode = def.TELEGRAM_PARSE_MODE_HTML
		data := PrepTelegramMessage(postHTML, env)
		var result TelegramPost
		json.Unmarshal(data, &result)
		if result.ParseMode != "HTML" {
			t.Errorf("Failed TestPrepTelegramMessage HTML ParseMode")
		}
	}

	{
		var env def.SessionData
		env.Res.Affordances.Remove = true

		data := PrepTelegramMessage(post, env)
		var remove TelegramRemovePost
		error := json.Unmarshal(data, &remove)

		if error != nil {
			t.Errorf("Failed TestPrepTelegramMessage RemovePost unmarshal JSON")
		}
		if remove.Markup.Remove != true {
			t.Errorf("Failed TestPrepTelegramMessage RemovePost")
		}
	}
	{
		var env def.SessionData
		var options []def.Option
		options = append(options, def.Option{Link: "Link1", Text: "OptionText1"})
		options = append(options, def.Option{Link: "Link2", Text: "OptionText2"})
		options = append(options, def.Option{Link: "Link3", Text: "OptionText3"})
		env.Res.Affordances.Options = options

		{
			data := PrepTelegramMessage(post, env)
			var reply TelegramReplyPost
			error := json.Unmarshal(data, &reply)

			if error != nil {
				t.Errorf("Failed TestPrepTelegramMessage ReplyPost unmarshal JSON")
			}
			if len(reply.Markup.Keyboard) == 0 {
				t.Errorf("Failed TestPrepTelegramMessage ReplyPost keyboard options")
			}
		}
		{
			env.Res.Affordances.Inline = true

			data := PrepTelegramMessage(post, env)
			var inline TelegramInlinePost
			error := json.Unmarshal(data, &inline)

			if error != nil {
				t.Errorf("Failed TestPrepTelegramMessage InlinePost unmarshal JSON")
			}
			if len(inline.Markup.Keyboard) == 0 {
				t.Errorf("Failed TestPrepTelegramMessage InlinePost keyboard options")
			}
		}
	}
}
