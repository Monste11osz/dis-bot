package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"
)

type WeatherInfo struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func commandWeather(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Content) <= len("!weather") {
		s.ChannelMessageSend(m.ChannelID, "Неверный запрос!\n Вот пример: `!weather Moscow`")
	}

	city := m.Content[len("!weather "):]
	weatherInfo, err := weatherProc(city)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Ошибка обработки данных!")
	}
	response := fmt.Sprintf("Weather in %s:\nTemperature: %.2f°C\nDescription: %s", city, weatherInfo.Main.Temp, weatherInfo.Weather[0].Description)

	embed := &discordgo.MessageEmbed{
		Description: response,
		Color:       0x0000FF,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func weatherProc(city string) (*WeatherInfo, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=a49a69e6d0bb51878e4056cdc3fd5cc4&units=metric", city)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weatherInfo WeatherInfo
	if err := json.NewDecoder(resp.Body).Decode(&weatherInfo); err != nil {
		return nil, err
	}
	return &weatherInfo, nil
}
