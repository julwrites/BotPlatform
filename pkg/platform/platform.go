package platform

import "github.com/julwrites/BotPlatform/pkg/def"

type Platform interface {
	Translate(body []byte) (def.SessionData, error)
	Post(env def.SessionData) bool
}
