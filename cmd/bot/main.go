package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rturovtsev/telegram-bot-weather/internal/chat"
	"github.com/rturovtsev/telegram-bot-weather/internal/handler"
	"log"
	"os"
)

var chatIDs []int64
var chatFile string

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	env := os.Getenv("ENV")

	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set")
	}
	if env == "dev" {
		chatFile = "chats.json"
	} else {
		chatFile = "/app/data/chats.json"
	}

	chatIDs = chat.LoadChatIDs(chatFile)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = env == "dev"

	log.Printf("Authorized on account %s", bot.Self.UserName)

	handler.ScheduleMessage(bot, chatIDs)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chat.AddChatID(update.Message.Chat.ID, chatFile)
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		} else if update.MyChatMember != nil {
			if update.MyChatMember.NewChatMember.Status == "administrator" {
				if update.MyChatMember.Chat.Type == "channel" {
					chat.AddChatID(update.MyChatMember.Chat.ID, chatFile) // добавим ID канала, если бот стал администратором
					log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				}
			}
		} else {
			continue
		}
	}
}
