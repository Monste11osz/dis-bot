package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"sync"
	"time"
)

type Poll struct {
	Question string
	Options  []string
	live     bool
}

var pollDeleteMutex sync.Mutex
var pollMessageID string

func commandPoll(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

	mes_str := strings.Join(args, " ")
	parts := strings.SplitN(mes_str, "|", 2)
	if len(parts) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Неверный формат. Пример: 'Вопрос | Вариант Вариант Вариант'.")
		return
	}

	question := strings.TrimSpace(parts[0])
	options := strings.Split(strings.TrimSpace(parts[1]), " ")

	poll := &Poll{
		Question: question,
		Options:  options,
		live:     false,
	}
	pollMessageID = m.ID
	go func() {
		time.Sleep(1 * time.Minute)
		poll.live = true
		if poll.live {
			startPollDeletion(s, m.ChannelID, pollMessageID)
		}
	}()
	sendPollMessage(s, m.ChannelID, poll)
}

func sendPollMessage(s *discordgo.Session, channelID string, poll *Poll) {
	embed := &discordgo.MessageEmbed{
		Title:       "Опрос",
		Description: poll.Question,
		Color:       0x0000ff,
	}

	for i, option := range poll.Options {
		name := fmt.Sprintf("%d. %s", i+1, option)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   name,
			Inline: true,
		})
	}

	message, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		fmt.Println("Ошибка отправки сообщения:", err)
		return
	}

	// Добавляем реакции к сообщению опроса
	for i := range poll.Options {
		emoji := emojiCode(i + 1)
		s.MessageReactionAdd(channelID, message.ID, emoji)
	}
}

func emojiCode(number int) string {
	// Возвращает эмодзи-код для числа
	return fmt.Sprintf("%d️⃣", number)
}

func startPollDeletion(s *discordgo.Session, channelID string, messageID string) {
	pollDeleteMutex.Lock()
	defer pollDeleteMutex.Unlock()

	// Проверка на совпадение идентификатора
	if pollMessageID == messageID {
		// Удаляем сообщение с опросом
		err := s.ChannelMessageDelete(channelID, messageID)
		if err != nil {
			fmt.Println("Ошибка удаления сообщения:", err)
			return
		}
		// Сбрасываем идентификатор сообщения
		pollMessageID = ""
		fmt.Println("Опрос удален.")
	}
}
