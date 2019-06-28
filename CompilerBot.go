package main

import (
	"./compiler"
	"./util"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var(
	supportingExtensoins = []string{"java","go","js","py"}
)

func main(){
	if os.Getenv("COMPILER_BOT_TOKEN") == ""{
		log.Fatal("Required environment variable is not set (\"COMPILER_BOT_TOKEN\")")
		return
	}
	err := exec.Command("docker").Run()
	if err != nil{
		log.Fatal("Couldn't find Docker command!")
	}
	er := exec.Command("docker","ps").Run()
	if er != nil{
		log.Fatal("You don't have permission to connect to Docker daemon or Docker daemon isn't running!")
	}
	log.Println("Updating docker images... (it may takes few minutes)")
	compiler.Build("java")
	compiler.Build("python")
	compiler.Build("nodejs")
	compiler.Build("golang")
	discord, err := discordgo.New()
	discord.Token = "Bot " + os.Getenv("COMPILER_BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}
	discord.AddHandler(onMessage)
	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Started CompilerBot")
	<-make(chan bool)
}
func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "c-compile" {
		if len(m.Attachments) == 0{
			embed := &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{},
				Color:  0xFF0000,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:   "Usage",
						Value:  "Send a file with comment \"c-compile\"!",
						Inline: false,
					},
					&discordgo.MessageEmbedField{
						Name:   "Environment",
						Value:  "Java: [latest](https://hub.docker.com/_/openjdk)\n" +
							"GoLang: [latest](https://hub.docker.com/_/golang)\n" +
							"Python: [latest v3](https://hub.docker.com/_/python)\n" +
							"Node.js: [latest](https://hub.docker.com/_/node)",
						Inline: false,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}
		var splittedName = strings.Split(m.Attachments[0].Filename,".")
		var extension = splittedName[len(splittedName)-1]
		if !util.ArrayContains(supportingExtensoins,extension){
			embed := &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{},
				Color:  0xFF0000,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:   "Unsupported file type",
						Value:  "Sorry, but file type \""+extension+"\" isn't supported yet.",
						Inline: false,
					},
					&discordgo.MessageEmbedField{
						Name:   "Supported file type list",
						Value:  ".java\n"+
							".js\n"+
							".py\n"+
							".go",
						Inline: false,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			return
		}
		newUUID,err := downloadFile(m.Attachments[0].Filename,m.Attachments[0].URL)
		if err != nil {
			embed := &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{},
				Color:  0xFF0000,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:   "Something went wrong",
						Value:  err.Error(),
						Inline: false,
					},
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			log.Println(err)
			return
		}
		log.Println("Downloaded file \""+m.Attachments[0].Filename+"\" with UUID \""+newUUID+"\"")
		switch extension {
		case "java":
			compiler.Run(newUUID,"java",m.ChannelID,s)
			break
		case "go":
			compiler.Run(newUUID,"golang",m.ChannelID,s)
			break
		case "js":
			compiler.Run(newUUID,"nodejs",m.ChannelID,s)
			break
		case "py":
			compiler.Run(newUUID,"python",m.ChannelID,s)
			break

		}
	}
}

func downloadFile(fileName,URL string) (string,error){
	resp, err := http.Get(URL)
	if err != nil {
		return "",err
	}
	defer resp.Body.Close()
	var newUUID = uuid.New().String()
	error := os.Mkdir(newUUID,os.FileMode(0755))
	if error != nil {
		return "",error
	}
	out, err := os.Create(newUUID+"/"+fileName)
	if err != nil {
		return "",err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "",err
	}
	out.Close()
	return newUUID,nil
}