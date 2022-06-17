package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token      = ""
	auctions   map[string]*Auction
	botActions = []string{
		"!hello to say hi",
		"!list to list all bid's channels",
		"!auction <itemName> <initialAmount> to create a new auction",
		"!bid <amount> to bid in auction",
		"!info to show bid information",
		"!help to show all commands",
	}
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
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt) //, os.Kill)
	<-sc

	dg.Close()
}

func messageControler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from bot himself
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(strings.ToLower(m.Content), "!") {
		return
	}

	switch {
	case m.Author.ID == s.State.User.ID || !strings.HasPrefix(strings.ToLower(m.Content), "!"):
		return
	// !hello responds "Hello World!"
	case strings.ToLower(m.Content) == "!hello":
		{
			s.ChannelMessageSend(m.ChannelID, "Hello World!")
			return
		}
	// !list prints all channels
	case strings.ToLower(m.Content) == "!list":
		{
			go handleListChannels(s, m)
			return
		}
	// !auction <item name> <initial amount>
	case strings.HasPrefix(strings.ToLower(m.Content), "!auction"):
		{
			go handlecreateAuction(s, m)
			return
		}
	// !bid <amount>
	case strings.HasPrefix(strings.ToLower(m.Content), "!bid"):
		{
			go handleBid(s, m)
			return
		}
	// !info
	case strings.HasPrefix(strings.ToLower(m.Content), "!info"):
		{
			go handleInfo(s, m)
			return
		}
	// !help to show all commands
	case strings.HasPrefix(strings.ToLower(m.Content), "!help"):
		{
			message := "List of commands:"
			lenght := len(botActions)

			for i := 1; i <= lenght; i++ {
				message += "\n" + strconv.Itoa(i) + ")  " + botActions[i-1]
			}

			s.ChannelMessageSend(m.ChannelID, message)
			return
		}
	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command. Try !help for commands info")
	}

	// !hello responds "Hello World!"
	// if strings.ToLower(m.Content) == "!hello" {
	// 	type GuildChannelCreateData struct {
	// 		Name                 string                 `json:"name"`
	// 		Type                 ChannelType            `json:"type"`
	// 		Topic                string                 `json:"topic,omitempty"`
	// 		Bitrate              int                    `json:"bitrate,omitempty"`
	// 		UserLimit            int                    `json:"user_limit,omitempty"`
	// 		RateLimitPerUser     int                    `json:"rate_limit_per_user,omitempty"`
	// 		Position             int                    `json:"position,omitempty"`
	// 		PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	// 		ParentID             string                 `json:"parent_id,omitempty"`
	// 		NSFW                 bool                   `json:"nsfw,omitempty"`
	// 	}
	// 	s.GuildChannelCreate(&GuildChannelCreateData{
	// 		Name: "Channel",
	// 		Type: "text",
	// 	}, err)
	// 	s.ChannelMessageSend(m.ChannelID, "Hello World!")
	// 	return
	// }
}

func handleListChannels(s *discordgo.Session, m *discordgo.MessageCreate) {

	for _, guild := range s.State.Guilds {
		channels, _ := s.GuildChannels(guild.ID)
		i := 1
		for _, c := range channels {
			if i <= 2 {
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
	if len(splits) != 3 || strings.ToLower(splits[0]) != "!auction" {
		s.ChannelMessageSend(m.ChannelID, "Error: wrong usage of !auction. Use !auction <item name> <initial amount> to create a new auction")
		return
	}

	//user := m.Author.Username
	itemName := splits[1]
	initialAmount, _ := strconv.Atoi(splits[2])

	auction, err := newAuction(itemName, initialAmount, s)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error while creating the Auction")
		return
	}
	auctions[m.ChannelID] = auction
}

func handleBid(s *discordgo.Session, m *discordgo.MessageCreate) {

	splits := strings.Split(m.Content, " ")
	if len(splits) != 2 || strings.ToLower(splits[0]) != "!bid" {
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
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: The amount must be greater than the current amount")
	}
	fmt.Println(auction)
	//return
}

func handleInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.ToLower(m.Content) != "!info" {
		s.ChannelMessageSend(m.ChannelID, "Error: Wrong use of !info")
		return
	}
	auction := auctions[m.ChannelID]
	info := auction.showInfo()
	s.ChannelMessageSend(m.ChannelID, info)
	//return
}
