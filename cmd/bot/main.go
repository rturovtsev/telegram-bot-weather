package main

import (
	"encoding/json"
	"github.com/fogleman/gg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
)

var chatIDs []int64

const chatFile = "chats.json"

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	env := os.Getenv("ENV")

	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	loadChatIDs()

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
		if update.Message == nil {
			continue
		}
		addChatID(update.Message.Chat.ID)

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}

func downloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func addBackgroundToImage(src image.Image) image.Image {
	const (
		width  = 725 // ширина нового изображения с фоном
		height = 400 // высота нового изображения с фоном
	)

	dc := gg.NewContext(width, height)
	dc.SetRGB(0, 0, 0) // черный цвет
	dc.Clear()

	x := (width - src.Bounds().Dx()) / 2
	y := (height - src.Bounds().Dy()) / 2

	dc.DrawImage(src, x, y)

	return dc.Image()
}

func saveChatIDs() {
	file, err := os.Create(chatFile)
	if err != nil {
		log.Panicln("Ошибка при сохранении chat ID:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(chatIDs); err != nil {
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
	if err = decoder.Decode(&chatIDs); err != nil {
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
		/*for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		t := time.NewTimer(next.Sub(now))
		<-t.C*/

		for _, chatID := range chatIDs {
			srcImage, err := downloadImage("https://xras.ru/upload_test/files/fc3_REL0.png")
			if err != nil {
				log.Println("Ошибка загрузки изображения:", err)
				continue
			}

			editedImage := addBackgroundToImage(srcImage)

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
		/*}*/
	}()
}
