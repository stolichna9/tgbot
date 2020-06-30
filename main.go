package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tkanos/gonfig"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	if runtime.GOOS == "windows" {
		err = gonfig.GetConf(dir+"\\config.json", &config)
	} else {
		err = gonfig.GetConf(dir+"/config.json", &config)
	}

	menunav := "mainmenu" // current menu position

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	userCollection := client.Database("store").Collection("user")

	// Connect to a telegram bot
	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)

	if err != nil {
		log.Panic(err)
	}

	log.Println("Authorized on telegram account " + bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		reply := ""
		switch menunav {
		case "mainmenu":
			switch update.Message.Command() {
			case "start", "mainmenu":
				reply = "You are in main menu."
				buttons := make([][]tgbotapi.KeyboardButton, 2)
				buttons[0] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/items"), tgbotapi.NewKeyboardButton("/myaccaunt"))
				buttons[1] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/mainmenu"))
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
					Keyboard:        buttons,
					ResizeKeyboard:  true,
					OneTimeKeyboard: true,
					Selective:       false,
				}
				bot.Send(msg)
				menunav = "mainmenu"
			case "myaccaunt":
				var user User
				filterUser := bson.D{{"username", update.Message.From.UserName}, {"telegramID", update.Message.From.ID}}
				err := userCollection.FindOne(context.TODO(), filterUser).Decode(&user)
				if err != nil {
					if err.Error() != "mongo: no documents in result" {
						log.Fatal(err)
					} else {
						_, err = userCollection.InsertOne(context.TODO(), bson.D{
							{"username", update.Message.From.UserName},
							{"telegramID", update.Message.From.ID},
							{"balance", 0},
						})
					}
				}
				reply = "Your ID: " + strconv.Itoa(update.Message.From.ID) + "\n"
				reply += "Your Username: @" + update.Message.From.UserName + "\n"
				reply += "Your balance: " + strconv.Itoa(user.Balance)
				buttons := make([][]tgbotapi.KeyboardButton, 2)
				buttons[0] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/deposit"))
				buttons[1] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/mainmenu"))
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
					Keyboard:        buttons,
					ResizeKeyboard:  true,
					OneTimeKeyboard: true,
					Selective:       false,
				}
				bot.Send(msg)
				menunav = "myaccauntmenu"
			}
		case "myaccauntmenu":
			switch update.Message.Command() {
			case "deposit":
				reply = "You are in deposit menu."
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
				menunav = "depositmenu"
			case "mainmenu":
				reply = "You are in main menu."
				buttons := make([][]tgbotapi.KeyboardButton, 2)
				buttons[0] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/items"), tgbotapi.NewKeyboardButton("/myaccaunt"))
				buttons[1] = tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/mainmenu"))
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
					Keyboard:        buttons,
					ResizeKeyboard:  true,
					OneTimeKeyboard: true,
					Selective:       false,
				}
				bot.Send(msg)
				menunav = "mainmenu"
			}

		}

	}

}
