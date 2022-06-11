package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var (
	token = ""

)

func ConnectToDiscord() {

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill) //, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from bot himself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// !Hello responde "Hello World!"
	if m.Content == "!Hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello World!")
		return
	}

	// !Listar imprime una lista de los canales
	if m.Content == "!Listar"{
		for _, guild := range s.State.Guilds{
			channels, _ := s.GuildChannels(guild.ID)
			i := 1
			for _, c := range channels{
				if (i <= 2){
					i++
					continue
				}
				mensaje := "" + strconv.Itoa(i-2) + " - " + string(c.Name)
				s.ChannelMessageSend(m.ChannelID, mensaje)
				i++
			}
		}
		return
	}
}
