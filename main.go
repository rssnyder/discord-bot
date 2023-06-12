package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var token, nickname, activity, metrics *string
var status, refresh *int
var updates prometheus.Counter

func init() {
	token = flag.String("token", "", "discord bot token")
	nickname = flag.String("nickname", "", "bot nickname")
	activity = flag.String("activity", "", "bot activity")
	status = flag.Int("status", 0, "0: playing, 1: listening, 2: watching")
	refresh = flag.Int("refresh", 300, "seconds between refresh")
	metrics = flag.String("metrics", ":8080", "address for prometheus metric serving")
	flag.Parse()

	updates = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "updates",
			Help: "Number of times discord has been updated",
		},
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(updates)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	go func() {
		log.Fatal(http.ListenAndServe(*metrics, nil))
	}()
}

func main() {
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
						updates.Inc()
					}
				}
			}
			if *activity != "" {
				switch *status {
				case 0:
					err = dg.UpdateGameStatus(0, *activity)
				case 1:
					err = dg.UpdateListeningStatus(*activity)
				case 2:
					err = dg.UpdateWatchStatus(0, *activity)
				}
				if err != nil {
					log.Printf("Unable to set activity: %s\n", err)
				} else {
					log.Printf("Set activity: %s\n", *activity)
					updates.Inc()
				}
			}
		}
	}
}
