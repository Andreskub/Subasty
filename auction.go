package main

import (
	"errors"
	"strconv"

	//"github.com/bwmarrin/discordgo"
)

type Auction struct {
	id              string
	userCreator     string
	userWinner      string
	initialAmount   int
	currentAmount   int
	productName     string
	openStatus      bool
}

func newAuction(productName string, initialAmount int, user_id string) (*Auction, error) {
	auctionName := "Auction - " + productName
	//auctionChannel, err := s.GuildChannelCreate(s.State.Guilds[0].ID, auctionName, discordgo.ChannelTypeGuildText)
	return &Auction{
		id:            auctionName, //auctionChannel.ID,
		userCreator:  user_id,
		userWinner:   "",
		initialAmount: initialAmount,
		currentAmount: initialAmount,
		productName:   productName,
		openStatus:    true,
	}, nil //err
}

func (a *Auction) bid(amount int, user_biding string) error {
	if a.openStatus == false{
		return errors.New("the auction has already been closed")
	}
	if amount <= a.currentAmount {
		return errors.New("input amount is lower than current amount")
	}
	a.currentAmount = amount
	a.userWinner = user_biding
	return nil
}

func (a *Auction) showInfo() string {
	message := ""
	message += "This auction is for " + a.productName + "\n"
	message += "It is currently at a price of: " + strconv.Itoa(a.currentAmount) + " and " + a.userWinner + " is winning!\n"
	return message
}

func (a *Auction) showInfoEnded() string{
	message := ""
	message += "The actual Auction has ended!\n"
	message += "The winner of the auction was: " + a.userWinner + " who bought it for: " + strconv.Itoa(a.currentAmount) + "\n"
	message += "¡Thank you for participating!"

	return message
}

func (a *Auction) getCreator() string{
	return a.userCreator
}

func (a *Auction) hasEnded() bool{
	return !a.openStatus
}

func (a *Auction) end() error{
	if a.openStatus == false{
		return errors.New("the auction has already been closed")
	}
	a.openStatus = false
	return nil
}