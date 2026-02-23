package platform

import (
	"testing"
)

func BenchmarkTelegramPreprocessing(b *testing.B) {
	input := "This is a test message with some special characters: []()~>#+-=|{}.!"
	for i := 0; i < b.N; i++ {
		TelegramPreprocessing(input)
	}
}
