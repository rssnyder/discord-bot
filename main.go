package main

import (
	"flag"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := flag.String("token", "", "discord bot token")
	nickname := flag.String("nickname", "", "bot nickname")
	activity := flag.String("activity", "", "bot activity")
	status := flag.Int("status", 0, "0: playing, 1: listening")
	refresh := flag.Int("refresh", 300, "seconds between refresh")
	flag.Parse()

	dg, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	guilds, err := dg.UserGuilds(100, "", "")
	if err != nil {
		log.Println(err)
		*nickname = ""
	}
	if len(guilds) == 0 {
		*nickname = ""
	}

	ticker := time.NewTicker(time.Duration(*refresh) * time.Second)

	for {
		select {
		case <-ticker.C:
			if *nickname != "" {
				for _, g := range guilds {
					err = dg.GuildMemberNickname(g.ID, "@me", *nickname)
					if err != nil {
						log.Println(err)
						continue
					} else {
						log.Printf("Set nickname in %s: %s\n", g.Name, *nickname)
					}
				}
			}
			if *activity != "" {
				switch *status {
				case 0:
					err = dg.UpdateGameStatus(0, *activity)
				case 1:
					err = dg.UpdateListeningStatus(*activity)
				}
				if err != nil {
					log.Printf("Unable to set activity: %s\n", err)
				} else {
					log.Printf("Set activity: %s\n", *activity)
				}
			}
		}
	}
}
