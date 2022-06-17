package main

import (
	"errors"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type Auction struct {
	id            string
	initialAmount int
	currentAmount int
	productName   string
}

func newAuction(productName string, initialAmount int, s *discordgo.Session) (*Auction, error) {
	auctionName := "Auction - " + productName
	//auctionChannel, err := s.GuildChannelCreate(s.State.Guilds[0].ID, auctionName, discordgo.ChannelTypeGuildText)

	return &Auction{
		id:            auctionName, //auctionChannel.ID,
		initialAmount: initialAmount,
		currentAmount: initialAmount,
		productName:   productName,
	}, nil //err
}

func (a *Auction) bid(amount int) error {

	if amount <= a.currentAmount {
		return errors.New("input amount is lower than current amount")
	}
	a.currentAmount = amount
	return nil
}

func (a *Auction) showInfo() string {
	return strconv.Itoa(a.currentAmount)
}
