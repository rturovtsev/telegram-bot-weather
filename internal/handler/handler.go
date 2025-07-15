package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rturovtsev/telegram-bot-weather/internal/weather"
	"time"
)

func ScheduleMessage(bot *tgbotapi.BotAPI, chatIDs []int64, token string) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, now.Location())
			if next.Before(now) {
				next = next.Add(24 * time.Hour)
			}
			t := time.NewTimer(next.Sub(now))
			<-t.C

			SendDailyMessage(bot, chatIDs, token)
		}
	}()
}

func SendDailyMessage(bot *tgbotapi.BotAPI, chatIDs []int64, token string) {
	for _, chatID := range chatIDs {
		/*srcImage, err := images.DownloadImage("https://xras.ru/upload_test/files/fc3_REL0.png")
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
		}*/
		//txt := "Прогноз магнитных бурь на три дня\n\n"
		//txt += "================\n\n"
		txt := "<strong>Тверь</strong>\n"
		txt += weather.GetWeather(weather.TverURL, token)
		

		targetDate := time.Date(2025, time.May, 31, 0, 0, 0, 0, time.UTC)
		kislovodsTargetDate := time.Date(2025, time.July, 26, 0, 0, 0, 0, time.UTC)
    	currentTime := time.Now()
		if currentTime.Before(targetDate) {
			txt += "\n\n<strong>Анталия</strong>\n"
			txt += weather.GetWeather(weather.AntaliaURL, token)
		} else {
			txt += "\n\n<strong>Сколково</strong>\n"
			txt += weather.GetWeather(weather.MoscowURL, token)
		}

		if currentTime.Before(kislovodsTargetDate) {
			txt += "\n\n<strong>Кисловодск</strong>\n"
			txt += weather.GetWeather(weather.KislovodsURL, token)
		}

		msg := tgbotapi.NewMessage(chatID, txt)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.DisableNotification = true
		bot.Send(msg)

		/*photo := tgbotapi.FilePath(file.Name())
		photoMessage := tgbotapi.NewPhoto(chatID, photo)
		photoMessage.Caption = txt
		photoMessage.ParseMode = tgbotapi.ModeHTML
		photoMessage.DisableNotification = true
		bot.Send(photoMessage)
		os.Remove(file.Name())*/
	}
}
