package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rturovtsev/telegram-bot-weather/internal/chat"
	"github.com/rturovtsev/telegram-bot-weather/internal/images"
	"image/png"
	"log"
	"os"
	"time"
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

	scheduleMessage(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chat.AddChatID(update.Message.Chat.ID, chatFile)
		} else if update.MyChatMember != nil {
			if update.MyChatMember.NewChatMember.Status == "administrator" {
				if update.MyChatMember.Chat.Type == "channel" {
					chat.AddChatID(update.MyChatMember.Chat.ID, chatFile) // добавим ID канала, если бот стал администратором
				}
			}
		} else {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}

func scheduleMessage(bot *tgbotapi.BotAPI) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
			if next.Before(now) {
				next = next.Add(24 * time.Hour)
			}
			t := time.NewTimer(next.Sub(now))
			<-t.C

			for _, chatID := range chatIDs {
				srcImage, err := images.DownloadImage("https://xras.ru/upload_test/files/fc3_REL0.png")
				if err != nil {
					log.Println("Ошибка загрузки изображения:", err)
					continue
				}

				editedImage := images.AddBackgroundToImage(srcImage)

				file, err := os.Create("edited_image.png")
				if err != nil {
					log.Println("Ошибка создания файла:", err)
					continue
				}
				defer file.Close()

				err = png.Encode(file, editedImage)
				if err != nil {
					log.Println("Ошибка при сохранении изображения:", err)
					continue
				}

				/*message := tgbotapi.NewMessage(chatID, "Доброе утро!")
				bot.Send(message)*/

				photo := tgbotapi.FilePath(file.Name())
				photoMessage := tgbotapi.NewPhoto(chatID, photo)
				photoMessage.Caption = "Прогноз магнитных бурь на три дня"
				bot.Send(photoMessage)
				os.Remove(file.Name())
			}
		}
	}()
}
