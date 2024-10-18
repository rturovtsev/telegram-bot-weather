package model

type Response struct {
	Now       int64      `json:"now"`       // время сервера unix timestamp
	NowDt     string     `json:"now_dt"`    // время сервера в UTC
	Info      Info       `json:"info"`      // информация о населенном пункте
	Fact      Fact       `json:"fact"`      // текущая погода
	Forecasts []Forecast `json:"forecasts"` // прогноз погоды
}

type Info struct {
	N             bool    `json:"n"`
	URL           string  `json:"url"`             // URL населенного пункта на сайте
	Lat           float64 `json:"lat"`             // широта
	Lon           float64 `json:"lon"`             // долгота
	Tzinfo        Tzinfo  `json:"tzinfo"`          // информация о часовом поясе
	DefPressureMm int     `json:"def_pressure_mm"` // нормальное давление для данной координаты в мм рт.ст.
	DefPressurePa int     `json:"def_pressure_pa"` // нормальное давление для данной координаты в паскалях
	Zoom          int     `json:"zoom"`
	Nr            bool    `json:"nr"`
	Ns            bool    `json:"ns"`
	Nsr           bool    `json:"nsr"`
	P             bool    `json:"p"`
	F             bool    `json:"f"`
	H             bool    `json:"_h"`
}

type Tzinfo struct {
	Name   string `json:"name"`
	Abbr   string `json:"abbr"`
	Dst    bool   `json:"dst"`
	Offset int    `json:"offset"`
}

type Fact struct {
	Daytime      string  `json:"daytime"`  // светлое или темное время суток. d - светлое, n - темное
	ObsTime      int64   `json:"obs_time"` // время замера погодных данных в unix timestamp
	Season       string  `json:"season"`   // сезон. «summer» — лето, «autumn» — осень, «winter» — зима. «spring» — весна
	Source       string  `json:"source"`
	Uptime       int64   `json:"uptime"`
	Cloudness    float64 `json:"cloudness"`  // облачность. 0 — ясно. 0.25 — малооблачно. 0.5 — облачно с прояснениями. 0.75 — облачно с прояснениями. 1 — пасмурно.
	Condition    string  `json:"condition"`  // состояние погоды. Расшифровка: clear — ясно. partly-cloudy — малооблачно. cloudy — облачно с прояснениями. overcast — пасмурно. light-rain — небольшой дождь. rain — дождь. heavy-rain — сильный дождь. showers — ливень. wet-snow — дождь со снегом. light-snow — небольшой снег. snow — снег. snow-showers — снегопад. hail — град. thunderstorm — гроза. thunderstorm-with-rain — дождь с грозой. thunderstorm-with-hail — гроза с градом.
	FeelsLike    float64 `json:"feels_like"` // ощущаемая температура
	Humidity     float64 `json:"humidity"`   // влажность в процентах
	Icon         string  `json:"icon"`       // иконка погоды. Иконка доступна по адресу https://yastatic.net/weather/i/icons/funky/dark/<значение из поля icon>.svg.
	IsThunder    bool    `json:"is_thunder"` // признак того, что погода является грозой
	Polar        bool    `json:"polar"`      // признак того, что время суток, указанное в поле daytime является полярным
	PrecProb     float64 `json:"prec_prob"`
	PrecStrength float64 `json:"prec_strength"` // сила осадков 0 — без осадков. 0.25 — слабый дождь/слабый снег. 0.5 — дождь/снег. 0.75 — сильный дождь/сильный снег. 1 — сильный ливень/очень сильный снег.
	PrecType     float64 `json:"prec_type"`     // тип осадков. 0 — без осадков, 1 — дождь, 2 — дождь со снегом, 3 — снег, 4 — град.
	PressureMm   float64 `json:"pressure_mm"`   // давление в мм рт.ст.
	PressurePa   float64 `json:"pressure_pa"`   // давление в паскалях
	Temp         float64 `json:"temp"`          // температура
	UvIndex      float64 `json:"uv_index"`
	WindAngle    float64 `json:"wind_angle"`
	WindDir      string  `json:"wind_dir"`
	WindGust     float64 `json:"wind_gust"`  // скорость порывов ветра
	WindSpeed    float64 `json:"wind_speed"` // скорость ветра
}

type Forecast struct {
	Date      string `json:"date"`       // Дата прогноза в формате ГГГГ-ММ-ДД
	DateTs    int64  `json:"date_ts"`    // Дата прогноза в формате Unixtime.
	Week      int    `json:"week"`       // Порядковый номер недели.
	Sunrise   string `json:"sunrise"`    // Время окончания восхода Солнца, локальное время
	Sunset    string `json:"sunset"`     // Время начала заката Солнца, локальное время
	RiseBegin string `json:"rise_begin"` // Время начала восхода Солнца, локальное время
	SetEnd    string `json:"set_end"`    // Время окончания заката Солнца, локальное время
	MoonCode  int    `json:"moon_code"`  // Код фазы Луны. Возможные значения: 0 — полнолуние. 1-3 — убывающая Луна. 4 — последняя четверть. 5-7 — убывающая Луна. 8 — новолуние. 9-11 — растущая Луна. 12 — первая четверть. 13-15 — растущая Луна.
	MoonText  string `json:"moon_text"`  // Текстовый код для фазы Луны. Возможные значения: moon-code-0 — полнолуние. moon-code-1 — убывающая луна. moon-code-2 — убывающая луна. moon-code-3 — убывающая луна. moon-code-4 — последняя четверть. moon-code-5 — убывающая луна. moon-code-6 — убывающая луна. moon-code-7 — убывающая луна. moon-code-8 — новолуние. moon-code-9 — растущая луна. moon-code-10 — растущая луна. moon-code-11 — растущая луна. moon-code-12 — первая четверть. moon-code-13 — растущая луна. moon-code-14 — растущая луна. moon-code-15 — растущая луна.
	Parts     Parts  `json:"parts"`      // Прогнозы по времени суток и 12-часовые прогнозы
	Hours     []Hour `json:"hours"`      // Прогноз на 24 часа.
}

type Parts struct {
	Day        PartsDay `json:"day"`         // прогноз на день.
	DayShort   PartsDay `json:"day_short"`   // 12-часовой прогноз на день.
	Evening    PartsDay `json:"evening"`     // прогноз на вечер
	Morning    PartsDay `json:"morning"`     // прогноз на утро
	Night      PartsDay `json:"night"`       // прогноз на ночь.
	NightShort PartsDay `json:"night_short"` // прогноз на ночь, для которого исключены поля temp_min и temp_max, в поле temp указывается минимальная температура за ночной период
}

type PartsDay struct {
	Daytime      string  `json:"daytime"` // Светлое или темное время суток. Возможные значения: «d» — светлое время суток. «n» — темное время суток.
	Source       string  `json:"_source"`
	Cloudness    float64 `json:"cloudness"`     // Облачность. Возможные значения: 0 — ясно. 0.25 — малооблачно. 0.5 — облачно с прояснениями. 0.75 — облачно с прояснениями. 1 — пасмурно.
	Condition    string  `json:"condition"`     // Код расшифровки погодного описания. Возможные значения: clear — ясно. partly-cloudy — малооблачно. cloudy — облачно с прояснениями. overcast — пасмурно. light-rain — небольшой дождь. rain — дождь. heavy-rain — сильный дождь. showers — ливень. wet-snow — дождь со снегом. light-snow — небольшой снег. snow — снег. snow-showers — снегопад. hail — град. thunderstorm — гроза. thunderstorm-with-rain — дождь с грозой. thunderstorm-with-hail — гроза с градом.
	FreshSnowMm  float64 `json:"fresh_snow_mm"` // Количество свежего снега (в мм). Вычисляется на основе значения поля prec_mm.
	Humidity     float64 `json:"humidity"`      // Влажность воздуха (в процентах)
	Icon         string  `json:"icon"`          // Код иконки погоды. Иконка доступна по адресу https://yastatic.net/weather/i/icons/funky/dark/<значение из поля icon>.svg.
	Polar        bool    `json:"polar"`         // Признак того, что время суток, указанное в поле daytime, является полярным
	PrecMm       float64 `json:"prec_mm"`       // Прогнозируемое количество осадков (в мм).
	PrecPeriod   float64 `json:"prec_period"`   // Прогнозируемый период осадков (в минутах).
	PrecProb     float64 `json:"prec_prob"`     // Вероятность выпадения осадков (в процентах).
	PrecStrength float64 `json:"prec_strength"` // Сила осадков. Возможные значения: 0 — без осадков. 0.25 — слабый дождь/слабый снег. 0.5 — дождь/снег. 0.75 — сильный дождь/сильный снег. 1 — сильный ливень/очень сильный снег.
	PrecType     float64 `json:"prec_type"`     // Тип осадков. Возможные значения: 0 — без осадков. 1 — дождь. 2 — дождь со снегом. 3 — снег.
	TempAvg      float64 `json:"temp_avg"`      // Средняя температура для времени суток (°C).
	Temp         float64 `json:"temp"`          // Видимо ошибка, Максимальная температура для времени суток (°C).
	TempMax      float64 `json:"temp_max"`      // Максимальная температура для времени суток (°C).
	TempMin      float64 `json:"temp_min"`      // Минимальная температура для времени суток (°C)
	FeelsLike    float64 `json:"feels_like"`    // ощущаемая температура
	UvIndex      float64 `json:"uv_index"`      // Ультрафиолетовый индекс
	WindAngle    float64 `json:"wind_angle"`
	WindDir      string  `json:"wind_dir"`   // Направление ветра. Возможные значения: «nw» — северо-западное. «n» — северное. «ne» — северо-восточное. «e» — восточное. «se» — юго-восточное. «s» — южное. «sw» — юго-западное. «w» — западное. «c» — штиль.
	WindGust     float64 `json:"wind_gust"`  // Скорость порывов ветра (в м/с)
	WindSpeed    float64 `json:"wind_speed"` // Скорость ветра (в м/с).
}

type Hour struct {
	Hour         string  `json:"hour"`          // Значение часа, для которого дается прогноз (0-23), локальное время.
	HourTs       int64   `json:"hour_ts"`       // Время прогноза в Unixtime.
	Cloudness    float64 `json:"cloudness"`     // Облачность. Возможные значения: 0 — ясно. 0.25 — малооблачно. 0.5 — облачно с прояснениями. 0.75 — облачно с прояснениями. 1 — пасмурно.
	Condition    string  `json:"condition"`     // Код расшифровки погодного описания. Возможные значения: clear — ясно. partly-cloudy — малооблачно. cloudy — облачно с прояснениями. overcast — пасмурно. light-rain — небольшой дождь. rain — дождь. heavy-rain — сильный дождь. showers — ливень. wet-snow — дождь со снегом. light-snow — небольшой снег. snow — снег. snow-showers — снегопад. hail — град. thunderstorm — гроза. thunderstorm-with-rain — дождь с грозой. thunderstorm-with-hail — гроза с градом.
	FeelsLike    float64 `json:"feels_like"`    // ощущаемая температура
	Humidity     float64 `json:"humidity"`      // Влажность воздуха (в процентах)
	Icon         string  `json:"icon"`          // Код иконки погоды. Иконка доступна по адресу https://yastatic.net/weather/i/icons/funky/dark/<значение из поля icon>.svg.
	IsThunder    bool    `json:"is_thunder"`    // признак того, что погода является грозой
	PrecPeriod   float64 `json:"prec_period"`   // Прогнозируемый период осадков (в минутах).
	PrecStrength float64 `json:"prec_strength"` // Сила осадков. Возможные значения: 0 — без осадков. 0.25 — слабый дождь/слабый снег. 0.5 — дождь/снег. 0.75 — сильный дождь/сильный снег. 1 — сильный ливень/очень сильный снег.
	PrecType     float64 `json:"prec_type"`     // Тип осадков. Возможные значения: 0 — без осадков. 1 — дождь. 2 — дождь со снегом. 3 — снег.
	Temp         float64 `json:"temp"`          // Максимальная дневная или минимальная ночная температура (°C)
	UvIndex      float64 `json:"uv_index"`      // Ультрафиолетовый индекс
	WindAngle    float64 `json:"wind_angle"`
	WindDir      string  `json:"wind_dir"`   // Направление ветра. Возможные значения: «nw» — северо-западное. «n» — северное. «ne» — северо-восточное. «e» — восточное. «se» — юго-восточное. «s» — южное. «sw» — юго-западное. «w» — западное. «c» — штиль.
	WindGust     float64 `json:"wind_gust"`  // Скорость порывов ветра (в м/с)
	WindSpeed    float64 `json:"wind_speed"` // Скорость ветра (в м/с).
}
