package weather

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rturovtsev/telegram-bot-weather/internal/model"
)

func TestGetWeather(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   model.Response
		expectedOutput string
	}{
		{
			name: "Sunny day with rain later",
			mockResponse: model.Response{
				Fact: model.Fact{
					Temp:       10,
					FeelsLike:  7,
					Condition:  "clear",
					PressureMm: 740,
				},
				Forecasts: []model.Forecast{
					{
						Parts: model.Parts{
							DayShort: model.PartsDay{
								TempMin:   5,
								Temp:      15,
								WindSpeed: 2.0,
								WindGust:  2.6,
								PrecProb:  30,
								PrecMm:    1.0,
							},
						},
						Hours: []model.Hour{
							{Hour: "12", PrecType: 0},
							{Hour: "13", PrecType: 1},
							{Hour: "14", PrecType: 1},
							{Hour: "15", PrecType: 0},
						},
					},
				},
			},
			expectedOutput: "Зонт нужен с 13:00 до 14:00 с вероятностью 30%\n" +
				"Температура сейчас 10°C, ясно, ощущается как 7°C\n" +
				"Давление 740 мм рт. ст.\n" +
				"В течение дня температура от 5°C до 15°C\n" +
				"Ветер 2.0 м/с, порывы до 2.6 м/с",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				respBytes, _ := json.Marshal(tt.mockResponse)
				w.Write(respBytes)
			}))
			defer server.Close()

			// Call GetWeather
			result := GetWeather(server.URL, "test_token")

			// Assert
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}

func TestCheckForRainBlocks(t *testing.T) {
	tests := []struct {
		name           string
		hours          []model.Hour
		expectedNeeded bool
		expectedHours  string
	}{
		{
			name: "Rain 1",
			hours: []model.Hour{
				{Hour: "8", PrecType: 0},
				{Hour: "9", PrecType: 0},
				{Hour: "10", PrecType: 0},
				{Hour: "11", PrecType: 0},
				{Hour: "12", PrecType: 0},
				{Hour: "13", PrecType: 0},
				{Hour: "14", PrecType: 0},
				{Hour: "15", PrecType: 0},
			},
			expectedNeeded: false,
			expectedHours:  "",
		},
		{
			name: "Rain 2",
			hours: []model.Hour{
				{Hour: "8", PrecType: 1},
				{Hour: "9", PrecType: 1},
				{Hour: "10", PrecType: 0},
				{Hour: "11", PrecType: 0},
				{Hour: "12", PrecType: 1},
				{Hour: "13", PrecType: 1},
				{Hour: "14", PrecType: 1},
				{Hour: "15", PrecType: 0},
			},
			expectedNeeded: true,
			expectedHours:  "с 8:00 до 9:00 и с 12:00 до 14:00",
		},
		{
			name: "Rain 3",
			hours: []model.Hour{
				{Hour: "8", PrecType: 1},
				{Hour: "9", PrecType: 0},
				{Hour: "10", PrecType: 0},
				{Hour: "11", PrecType: 0},
				{Hour: "12", PrecType: 0},
				{Hour: "13", PrecType: 0},
				{Hour: "14", PrecType: 0},
				{Hour: "15", PrecType: 0},
			},
			expectedNeeded: true,
			expectedHours:  "в 8:00",
		},
		{
			name: "Rain 4",
			hours: []model.Hour{
				{Hour: "8", PrecType: 1},
				{Hour: "9", PrecType: 0},
				{Hour: "10", PrecType: 0},
				{Hour: "11", PrecType: 0},
				{Hour: "12", PrecType: 0},
				{Hour: "13", PrecType: 0},
				{Hour: "14", PrecType: 0},
				{Hour: "15", PrecType: 1},
			},
			expectedNeeded: true,
			expectedHours:  "в 8:00 и в 15:00",
		},
		{
			name: "Rain 5",
			hours: []model.Hour{
				{Hour: "8", PrecType: 1},
				{Hour: "9", PrecType: 0},
				{Hour: "10", PrecType: 0},
				{Hour: "11", PrecType: 1},
				{Hour: "12", PrecType: 1},
				{Hour: "13", PrecType: 0},
				{Hour: "14", PrecType: 0},
				{Hour: "15", PrecType: 1},
			},
			expectedNeeded: true,
			expectedHours:  "в 8:00, с 11:00 до 12:00 и в 15:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call checkForRainBlocks
			umbrellaNeeded, rainHoursStr := checkForRainBlocks(tt.hours)

			// Assert
			assert.Equal(t, tt.expectedNeeded, umbrellaNeeded)
			assert.Equal(t, tt.expectedHours, rainHoursStr)
		})
	}
}
