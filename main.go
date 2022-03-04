package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	Store = &MemStore{}
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

	Store.LoadDefaults([]Picture{
		{URL: "https://i.ytimg.com/vi/ByH9LuSILxU/maxresdefault.jpg"},
		{URL: "https://m.media-amazon.com/images/I/71bcwxa4FmL._AC_SX425_.jpg"},
	})
}

type DataStore interface {
	GetRandomCat() Picture
}

type Picture struct {
	URL string
}

type MemStore struct {
	Defaults []Picture
	Len      int
}

func (s *MemStore) LoadDefaults(pics []Picture) {
	s.Defaults = pics
	s.Len = len(pics)
}

func (s *MemStore) GetRandomCat() Picture {
	i := rand.Intn(s.Len)
	return s.Defaults[i]
}

func main() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalln("error creating Discord session,", err)
	}

	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!cat" {
		s.ChannelMessageSend(m.ChannelID, Store.GetRandomCat().URL)
	}

	if strings.Contains(m.Content, "sad") {
		s.ChannelMessageSend(m.ChannelID, Store.GetRandomCat().URL)
	}
}
