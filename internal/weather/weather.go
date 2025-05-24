package weather

import (
	"encoding/json"
	"fmt"
	"github.com/rturovtsev/telegram-bot-weather/internal/model"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const TverURL = "https://api.weather.yandex.ru/v2/forecast?lat=56.83270&lon=35.93039"
const MoscowURL = "https://api.weather.yandex.ru/v2/forecast?lat=55.70077&lon=37.360078"
const AntaliaURL = "https://api.weather.yandex.ru/v2/forecast?lat=36.88590&lon=30.67414"

const PrecTypeRain = 1
const PrecTypeSnowRain = 2

func GetWeather(url string, token string) string {
	weather := makeRequest(url, token)

	factTemp := weather.Fact.Temp
	factPressureMm := weather.Fact.PressureMm
	feelsLike := weather.Fact.FeelsLike
	condition := weather.Fact.Condition

	dayShort := weather.Forecasts[0].Parts.DayShort
	tempMin := dayShort.TempMin
	tempMax := dayShort.Temp
	windSpeed := dayShort.WindSpeed
	windGust := dayShort.WindGust
	uvIndex := dayShort.UvIndex

	umbrellaNeeded, rainHours := checkForRainBlocks(weather.Forecasts[0].Hours)

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
		umbrellaText = fmt.Sprintf("Зонт нужен %s", rainHours)
		if int(dayShort.PrecProb) > 0 {
			umbrellaText = fmt.Sprintf("%s с вероятностью %d%%", umbrellaText, int(dayShort.PrecProb))
		}
	}
	var uvText string
	if uvIndex >= 4 && uvIndex < 6 {
		uvText = "\nРекомендуется использование SPF крема"
	} else if uvIndex >= 6 {
		uvText = "\nИспользование SPF 50+ крема <strong>обязательно</strong>"
	}

	text := fmt.Sprintf(
		"%s\n"+
			"Температура сейчас %.0f°C, %s, ощущается как %.0f°C\n"+
			"Давление %.0f мм рт. ст.\n"+
			"В течение дня температура от %.0f°C до %.0f°C\n"+
			"Ветер %.1f м/с, порывы до %.1f м/с"+
			"%s",
		umbrellaText,
		factTemp,
		condition,
		feelsLike,
		factPressureMm,
		tempMin,
		tempMax,
		windSpeed,
		windGust,
		uvText,
	)

	return text
}

func checkForRainBlocks(hours []model.Hour) (bool, string) {
	umbrellaNeeded := false
	var rainHours []string
	var rainHoursArs [][]int
	var tmpHour []int

	for _, hour := range hours {
		hourInt, err := strconv.Atoi(hour.Hour)
		if err != nil || hourInt < 6 || hourInt > 21 {
			continue
		}

		if hour.PrecType == PrecTypeRain || hour.PrecType == PrecTypeSnowRain {
			umbrellaNeeded = true
			tmpHour = append(tmpHour, hourInt)
		} else {
			if len(tmpHour) > 0 {
				rainHoursArs = append(rainHoursArs, tmpHour)
				tmpHour = nil
			}
		}
	}

	if tmpHour != nil {
		rainHoursArs = append(rainHoursArs, tmpHour)
	}

	for _, item := range rainHoursArs {
		if len(item) == 1 {
			rainHours = append(rainHours, fmt.Sprintf("в %d:00", item[0]))
		} else {
			rainHours = append(rainHours, fmt.Sprintf("с %d:00 до %d:00", item[0], item[len(item)-1]))
		}
	}

	rainHoursStr := joinWithCommaAnd(rainHours)

	return umbrellaNeeded, rainHoursStr
}

func joinWithCommaAnd(strArr []string) string {
	if len(strArr) == 0 {
		return ""
	}

	if len(strArr) == 1 {
		return strArr[0]
	}

	if len(strArr) == 2 {
		return strArr[0] + " и " + strArr[1]
	}

	return strings.Join(strArr[:len(strArr)-1], ", ") + " и " + strArr[len(strArr)-1]
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
