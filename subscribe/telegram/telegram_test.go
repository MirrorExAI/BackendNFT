package telegram

import (
	"context"
	"os"
	"os/signal"
	"testing"
)

func TestTelegramBot_SendUserApproval(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	bot, _ := NewTelegramBot("6595532482:AAHTyurqRhRYFvw5h7WgiDty8wfSeq0LZko")
	bot.SendUserApproval(ctx)
	defer cancel()
}
