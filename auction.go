package main

import (
	"github.com/bwmarrin/discordgo"
	"errors"
)

type Auction struct {
	id string
	initialAmount int
	currentAmount int
	productName string
}

func newAuction(productName string, initialAmount int, s *discordgo.Session) (*Auction, error){
	auctionName := "Auction - " + productName
	auctionChannel, err := s.GuildChannelCreate(s.State.Guilds[0].ID, auctionName, discordgo.ChannelTypeGuildText)
	
	return &Auction{
		id: auctionChannel.ID,
		initialAmount: initialAmount,
		currentAmount: initialAmount,
		productName: productName,
	}, err
}

func (a *Auction) bid(amount int)(error){

	if amount <= a.currentAmount {
		return errors.New("Amount lower than current")
	}
	a.currentAmount = amount
	return nil
}