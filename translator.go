package botmultiplexer

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Translator component to handle translation of a HTTP payload into a
// consistent format, and to translate that format back into a HTTP payload
// for posting

func TranslateToProps(req *http.Request, env *SessionData) bool {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("Error occurred reading http request: %v", err)
		return false
	}
	log.Printf("Request body: %s", strings.ReplaceAll(string(reqBody), "\n", ""))

	translated := false

	switch env.Type {
	case TYPE_TELEGRAM:
		translated = TelegramTranslate(reqBody, env)
	}

	if translated {
		return translated
	}

	return false
}

func PostFromProps(env *SessionData) bool {
	switch env.Type {
	case TYPE_TELEGRAM:
		return PostTelegram(env)
	}

	return false
}