package bot

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"time"
)

var chatIDs []int64

const chatFile = "chats.json"

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	env := os.Getenv("ENV")

	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = env == "dev"

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}

func saveChatIDs() {
	file, err := os.Create(chatFile)
	if err != nil {
		log.Panicln("Ошибка при сохранении chat ID:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(chatIDs); err != nil {
		log.Panicln("Ошибка при кодировании chat ID:", err)
	}
}

func loadChatIDs() {
	file, err := os.Open(chatFile)
	if os.IsNotExist(err) {
		return // Файл не существует, ничего страшного
	}
	if err != nil {
		log.Panicln("Ошибка при открытии файла:", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&chatIDs); err != nil {
		log.Panicln("Ошибка при декодировании chat ID:", err)
	}
}

func addChatID(chatID int64) {
	for _, id := range chatIDs {
		if id == chatID {
			return // ChatID уже существует, ничего не делать
		}
	}
	chatIDs = append(chatIDs, chatID)
	saveChatIDs()
}

func scheduleMessage(bot *tgbotapi.BotAPI) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, now.Location())
			if next.Before(now) {
				next = next.Add(24 * time.Hour)
			}
			t := time.NewTimer(next.Sub(now))
			<-t.C

			for _, chatID := range chatIDs {
				message := tgbotapi.NewMessage(chatID, "Доброе утро!")
				bot.Send(message)

				photo := tgbotapi.FileURL("https://xras.ru/upload_test/files/fc3_REL0.png")
				photoMessage := tgbotapi.NewPhoto(chatID, photo)
				photoMessage.Caption = "Прогноз магнитных бурь на три дня"
				bot.Send(photoMessage)
			}
		}
	}()
}
