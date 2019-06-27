package compiler

import (
	"bytes"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func Build(lang string){
	downloadFiles(lang)
	err := exec.Command("docker", "build" ,"-t", "compile/"+lang,"./"+lang).Run()
	if err != nil{
		log.Fatal("Failed to build "+lang+" image!")
		return
	}
}

func Run(uuid,lang,channelId string,discord *discordgo.Session){
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0xFF0000,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Compiling",
				Value:  "We're compiling your program and it'll be executed!",
				Inline: false,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	discord.ChannelMessageSendEmbed(channelId, embed)
	log.Println("Running file with UUID \""+uuid+"\"")
	path, _ := os.Getwd()
	src, _ := os.Open(lang+"/compile.sh")
	defer src.Close()
	dst, _ := os.Create(uuid+"/compile.sh")
	io.Copy(dst, src)
	dst.Close()
	os.Chmod(uuid+"/compile.sh",os.FileMode(0777))
	commandArgs := []string{"run","--rm","-i","-v", path+"/"+uuid+":/build","--network=none","--memory=128MB","compile/"+lang}
	cmd := exec.Command("docker",commandArgs...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Start()
	if err != nil{
		log.Println(err)
		log.Println("Failed to run "+lang+" image!")
		return
	}
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	timeout := time.After(time.Minute)
	select {
	case <-timeout:
		cmd.Process.Kill()
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Color:  0xFF0000,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Command execution timed out",
					Value:  "Sorry, but CompilerBot has execution timeout to prevent malicious programs.",
					Inline: false,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		discord.ChannelMessageSendEmbed(channelId, embed)
	case _ = <-done:
		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Color:  0x00FF00,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Output",
					Value:  stdout.String(),
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "Errors",
					Value:  stderr.String(),
					Inline: false,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		discord.ChannelMessageSendEmbed(channelId, embed)
	}
	os.RemoveAll(uuid)
}

func downloadFiles(lang string){
	os.Mkdir(lang, os.FileMode(0755))
	downloadAFile(lang,"compile.sh","https://ry0tak.github.io/Files/docker/"+lang+"/compile.sh")
	downloadAFile(lang,"Dockerfile","https://ry0tak.github.io/Files/docker/"+lang+"/Dockerfile")
}

func downloadAFile(folder,filename,url string){
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	out, err := os.Create(folder+"/"+filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
}