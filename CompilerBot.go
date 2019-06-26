package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

func main(){
	if os.Getenv("COMPILER_BOT_TOKEN") == ""{
		fmt.Println("Required environment variable is not set (\"COMPILER_BOT_TOKEN\")")
		return
	}
	discord, err := discordgo.New()
	discord.Token = "Bot " + os.Getenv("COMPILER_BOT_TOKEN")
	if err != nil {
		fmt.Println(err)
		return
	}
	discord.AddHandler(onMessage)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	<-make(chan bool)
}
func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content,"c-compile") {

		//Compile something
	}
}