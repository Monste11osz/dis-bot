package main

import (
	"cloud.google.com/go/translate"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"strings"
)

func commandTranslate(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(m.Content) < len("!translate") {
		s.ChannelMessageSend(m.ChannelID, "Неверный запрос!\n Вот пример: `!translate ru Your english is very nice!`")
		return
	}
	txtTranslate := strings.Join(args[1:], " ")
	txtLanguage := language.Make(args[0])
	txtReady, err := translates(txtTranslate, txtLanguage)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Ошибка перевода!. Попробуйте повторить.")
		return
	}
	response := fmt.Sprintf("Translate: %s", txtReady)
	embed := &discordgo.MessageEmbed{
		Description: response,
		Color:       0x0000FF,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func translates(txtT string, txtL language.Tag) (string, error) {
	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey("AIzaSyDfxzeRhu6osIHMJVMgZWdvTEJ9PqV9mNQ"))
	if err != nil {
		return "", err
	}
	defer client.Close()

	translations, err := client.Translate(ctx, []string{txtT}, txtL, nil)
	if err != nil {
		return "", err
	}

	return translations[0].Text, nil
}
