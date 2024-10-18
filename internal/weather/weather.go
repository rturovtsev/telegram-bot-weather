package weather

import (
	"encoding/json"
	"fmt"
	"github.com/rturovtsev/telegram-bot-weather/internal/model"
	"io"
	"log"
	"net/http"
	"strconv"
)

const TverURL = "https://api.weather.yandex.ru/v2/forecast?lat=56.83270&lon=35.93039"
const MoscowURL = "https://api.weather.yandex.ru/v2/forecast?lat=55.70077&lon=37.360078"

func GetWeather(url string, token string) string {
	weather := makeRequest(url, token)

	factTemp := weather.Fact.Temp
	feelsLike := weather.Fact.FeelsLike
	condition := weather.Fact.Condition

	dayShort := weather.Forecasts[0].Parts.DayShort
	tempMin := dayShort.TempMin
	//tempMax := dayShort.TempMax
	tempMax := dayShort.Temp
	windSpeed := dayShort.WindSpeed
	windGust := dayShort.WindGust

	umbrellaNeeded, rainHours := checkForRain(weather.Forecasts[0].Hours)

	switch condition {
	case "clear":
		condition = "ясно"
	case "partly-cloudy":
		condition = "малооблачно"
	case "cloudy":
		condition = "облачно с прояснениями"
	case "overcast":
		condition = "пасмурно"
	case "light-rain":
		condition = "небольшой дождь"
	case "rain":
		condition = "дождь"
	case "heavy-rain":
		condition = "сильный дождь"
	case "showers":
		condition = "ливень"
	case "wet-snow":
		condition = "дождь со снегом"
	case "light-snow":
		condition = "небольшой снег"
	case "snow":
		condition = "снег"
	case "snow-showers":
		condition = "снегопад"
	case "hail":
		condition = "град"
	case "thunderstorm":
		condition = "гроза"
	case "thunderstorm-with-rain":
		condition = "дождь с грозой"
	case "thunderstorm-with-hail":
		condition = "гроза с градом"
	default:
		condition = "не понятно"
	}

	var umbrellaText string
	if !umbrellaNeeded {
		umbrellaText = "Зонт не нужен"
	} else {
		umbrellaText = fmt.Sprintf("Зонт нужен с %s с вероятностью %d%%", rainHours, int(dayShort.PrecProb))
	}

	text := fmt.Sprintf(
		"%s\nТемпература сейчас %.0f°C, %s, ощущается как %.0f°C\n"+
			"В течение дня температура от %.0f°C до %.0f°C\n"+
			"В течение дня ветер %.1f м/с, порывы до %.1f м/с",
		umbrellaText,
		factTemp,
		condition,
		feelsLike,
		tempMin,
		tempMax,
		windSpeed,
		windGust,
	)

	return text
}

func checkForRain(hours []model.Hour) (bool, string) {
	umbrellaNeeded := false
	rainHours := ""
	rainHoursStart := ""
	rainHoursEnd := ""

	for _, hour := range hours {
		hourInt, err := strconv.Atoi(hour.Hour)
		if err != nil || hourInt < 6 || hourInt > 21 {
			continue
		}

		if hour.PrecType != 0 {
			umbrellaNeeded = true
			if rainHoursStart == "" {
				rainHoursStart = hour.Hour + ":00"
			} else {
				rainHoursEnd = hour.Hour + ":00"
			}
		}
	}

	rainHours = fmt.Sprintf("%s до %s:00", rainHoursStart, rainHoursEnd)

	return umbrellaNeeded, rainHours
}

func makeRequest(url string, token string) (weatherResponse *model.Response) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-Yandex-Weather-Key", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Fatal(err)
	}

	return
}
