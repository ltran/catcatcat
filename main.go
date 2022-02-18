package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("Starting the bot")
	discord, err := discordgo.New("Bot " + "authentication token")

	if err != nil {
		panic("OH NO!")
	}
	fmt.Println(discord)
}
