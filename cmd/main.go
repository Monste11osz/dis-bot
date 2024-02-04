package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	CommandPrefix = "!"
	token         = "your token"
)

var polls = make(map[string]Poll)

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error create ", err)
		return
	}
	dg.AddHandler(messageCreate)
	err = dg.Open()
	if err != nil {
		fmt.Println("Error running ", err)
		return
	}
	fmt.Println("Bot is running now ")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	content := strings.TrimSpace(m.Content)
	if strings.HasPrefix(content, CommandPrefix) {
		content = strings.TrimPrefix(content, CommandPrefix)
		command := strings.Fields(content)

		fmt.Println("Command:", command[0])
		switch command[0] {
		case "help":
			commandHelp(s, m)
		case "poll":
			commandPoll(s, m, command[1:])
		case "weather":
			commandWeather(s, m)
		case "info":
			s.ChannelMessageSend(m.ChannelID, "Приветствую тебя мой дорогой пользователь. Я бот написанный на языке Golang.\nМой ассортимент команд не столь обширный, но чутка информативный! Чтобы узнать о них, напиши `!help` ")
		case "translate":
			commandTranslate(s, m, command[1:])
		default:
			s.ChannelMessageSend(m.ChannelID, "Привет. Для пользованием бота, ознакомьтесь со справкой, а именно: напишите `!info`")

		}
	}
}

func commandHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Список команд",
		Description: "Вот список доступных команд:",
		Color:       0x0000FF,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "!poll", Value: "Опрос для пользователей. Для пользователей сервера, вы можете создавать опросы на разные темы. \n Пример: `Любимая марка? | Мерс Бмв`", Inline: false},
			{Name: "!info", Value: "Краткая информация, навигация!", Inline: false},
			{Name: "!weather", Value: "Информация о погоде. Написав по примеру запрос, вы можете узнать погоду в любой точке земного шара.\n Пример: `!weather Astana`", Inline: false},
			{Name: "!translate", Value: "Данной командой вы запустите перевод текста на удобный вам язык.\n Пример: `!translate ru Your english is very nice!`", Inline: false},
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
