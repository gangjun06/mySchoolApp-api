package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/osang-school/backend/internal/conf"
)

var (
	dg *discordgo.Session
)

func Init() {
	var err error
	dg, err = discordgo.New("Bot " + conf.Discord().Token)
	if err != nil {
		log.Fatalln(err)
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	if err := dg.Open(); err != nil {
		log.Fatalln(err)
	}

}

func Close() {
	dg.Close()
}

func ready(s *discordgo.Session, m *discordgo.Ready) {
	s.UpdateGameStatus(0, "오상중학교 학생회 봇")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
}

func SendEmbed(channelID string, embed *discordgo.MessageEmbed) {
	dg.ChannelMessageSendEmbed(channelID, embed)
}
