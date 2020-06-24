package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("YOUR_TOKEN")
	if err != nil {
		log.Panic(err)
	}
	log.Println("Authorized on account " + bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		reply := ""
		switch update.Message.Command() {
		case "start":
			reply = "choose any command in bellow"
		case "info":
			reply = "you choose command info"
		case "info2":
			reply = "you choose command info2"
		case "currencies":
			reply = "USD/RUB __/__"
		}
		buttons := make([][]tgbotapi.KeyboardButton, 1)
		buttons[0] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/info"), tgbotapi.NewKeyboardButton("/info2"))
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
			Keyboard:        buttons,
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
			Selective:       false,
		}
		bot.Send(msg)
	}

}
