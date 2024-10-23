package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rturovtsev/telegram-bot-weather/internal/chat"
	"github.com/rturovtsev/telegram-bot-weather/internal/handler"
	"log"
	"os"
	"strconv"
)

var chatIDs []int64
var chatFile string
var adminID int64
var err error
var bot *tgbotapi.BotAPI

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	yandexToken := os.Getenv("YANDEX_TOKEN")
	env := os.Getenv("ENV")
	adminIDEnv := os.Getenv("ADMIN_ID")

	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set")
	}
	if env == "dev" {
		chatFile = "chats.json"
	} else {
		chatFile = "/app/data/chats.json"
	}

	if adminIDEnv != "" {
		adminID, err = strconv.ParseInt(adminIDEnv, 10, 64)
		if err != nil {
			log.Printf("Error admin id convert: %v", err)
		}
	}

	chatIDs = chat.LoadChatIDs(chatFile)

	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = env == "dev"

	log.Printf("Authorized on account %s", bot.Self.UserName)

	handler.ScheduleMessage(bot, chatIDs, yandexToken)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && adminID != 0 {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "send":
					if update.Message.From.ID == adminID {
						handler.SendDailyMessage(bot, chatIDs, yandexToken)
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет разрешения на использование этой команды.")
						bot.Send(msg)
					}
				}
			}
		}

		if update.MyChatMember != nil {
			chatType := update.MyChatMember.Chat.Type

			if chatType == "channel" {
				if update.MyChatMember.NewChatMember.Status == "administrator" {
					chat.AddChatID(update.MyChatMember.Chat.ID, chatFile) // добавим ID канала, если бот стал администратором
					log.Printf("Добавление в канал [%s] %s", update.MyChatMember.From.UserName, update.MyChatMember.Chat.Title)
				} else if update.MyChatMember.NewChatMember.Status == "left" || update.MyChatMember.NewChatMember.Status == "kicked" || update.MyChatMember.NewChatMember.Status == "member" {
					chat.RemoveChatID(update.MyChatMember.Chat.ID, chatFile)
					log.Printf("Удаление из канала [%s] %s", update.MyChatMember.From.UserName, update.MyChatMember.Chat.Title)
				}
			} else if chatType == "private" {
				if update.MyChatMember.NewChatMember.Status == "left" || update.MyChatMember.NewChatMember.Status == "kicked" {
					chat.RemoveChatID(update.Message.Chat.ID, chatFile)
					log.Printf("Удаление из личного чата [%s] %s", update.Message.From.UserName, update.Message.Text)
				} else {
					chat.AddChatID(update.Message.Chat.ID, chatFile)
					log.Printf("Добавление в личный чат [%s] %s", update.Message.From.UserName, update.Message.Text)
				}
			} else if chatType == "group" || chatType == "supergroup" {
				if update.MyChatMember.NewChatMember.Status == "member" {
					chat.AddChatID(update.MyChatMember.Chat.ID, chatFile) // Добавить ID группы
					log.Printf("Добавление в группу [%s] %s", update.MyChatMember.From.UserName, update.MyChatMember.Chat.Title)
				} else if update.MyChatMember.NewChatMember.Status == "left" || update.MyChatMember.NewChatMember.Status == "kicked" {
					chat.RemoveChatID(update.MyChatMember.Chat.ID, chatFile) // Удалить ID группы
					log.Printf("Удаление из группы [%s] %s", update.MyChatMember.From.UserName, update.MyChatMember.Chat.Title)
				}
			} else {
				continue
			}
		}
	}
}
