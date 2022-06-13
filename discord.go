package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"strings"
	
	"github.com/bwmarrin/discordgo"
)

var (
	token = ""
	auctions map[string]*Auction
)

func ConnectToDiscord() {
	auctions = make(map[string]*Auction)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(messageControler)
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

func messageControler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from bot himself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// !hello responds "Hello World!"
	if strings.ToLower(m.Content) == "!hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello World!")
		return
	}

	// !list prints all channels
	if strings.ToLower(m.Content) == "!list"{
		handleListChannels(s,m)
		return
	}

	// !auction <item name> <initial amount>
	if strings.HasPrefix(strings.ToLower(m.Content), "!auction"){
		handlecreateAuction(s,m)
		return
	}

	// !bid <amount>
	if strings.HasPrefix(strings.ToLower(m.Content), "!bid"){
		handleBid(s,m)
		return
	}
}

func handleListChannels(s *discordgo.Session, m *discordgo.MessageCreate){
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
}


func handlecreateAuction(s *discordgo.Session, m *discordgo.MessageCreate) {

	splits := strings.Split(m.Content, " ")
	if (len(splits) != 3 || strings.ToLower(splits[0]) != "!auction") {
		s.ChannelMessageSend(m.ChannelID, "Error: wrong usage of !auction. Use !auction <item name> <initial amount> to create a new auction")
		return
	}

	//user := m.Author.Username
	itemName := splits[1]
	initialAmount, _ := strconv.Atoi(splits[2])

	auction, err := newAuction(itemName, initialAmount, s)
	if (err != nil){
		s.ChannelMessageSend(m.ChannelID, "Error while creating the Auction")
		return
	}
	auctions[auction.id] = auction
}

func handleBid(s *discordgo.Session, m *discordgo.MessageCreate) {
	
	splits := strings.Split(m.Content, " ")
	if (len(splits) != 2 || strings.ToLower(splits[0]) != "!bid") {
		s.ChannelMessageSend(m.ChannelID, "Error: Wrong usage of !bid. Use !bid <amount> to bid")
		return
	}
	
	if _, err := strconv.Atoi(splits[1]); err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: The input is not a number")
		return
	}
	//user := m.Author.Username
	amount, _ := strconv.Atoi(splits[1])

	auction := auctions[m.ChannelID]
	err := auction.bid(amount)
	if (err != nil){
		s.ChannelMessageSend(m.ChannelID, "Error: The amount must be greater than the current amount")
	}
	fmt.Println(auction)
	return
}