package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var token, nickname, activity, status, refresh, metrics *string
var statusCode, refreshSec int
var err error
var nicknameUpdates, activityUpdates prometheus.Counter

func init() {
	token = flag.String("token", getEnv("TOKEN", ""), "discord bot token")
	nickname = flag.String("nickname", getEnv("NICKNAME", ""), "bot nickname")
	activity = flag.String("activity", getEnv("ACTIVITY", ""), "bot activity")
	status = flag.String("status", getEnv("STATUS", "2"), "0: playing, 1: listening, 2: watching")
	refresh = flag.String("refresh", getEnv("REFRESH", "120"), "seconds between refresh")
	metrics = flag.String("metrics", getEnv("METRICS", ":8080"), "address for prometheus metric serving")
	flag.Parse()

	if statusCode, err = strconv.Atoi(*status); err != nil {
		log.Println(err)
		statusCode = 2
	}
	if refreshSec, err = strconv.Atoi(*refresh); err != nil {
		log.Println(err)
		refreshSec = 120
	}

	nicknameUpdates = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "nickname_updates",
			Help: "Number of times discord nickname has been updated",
		},
	)
	activityUpdates = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "activity_updates",
			Help: "Number of times discord activity has been updated",
		},
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(nicknameUpdates)
	reg.MustRegister(activityUpdates)
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

	ticker := time.NewTicker(time.Duration(refreshSec) * time.Second)

	defer ticker.Stop()
	for ; true; <-ticker.C {
		if *nickname != "" {
			for _, g := range guilds {
				err = dg.GuildMemberNickname(g.ID, "@me", *nickname)
				if err != nil {
					log.Println(err)
					continue
				} else {
					log.Printf("Set nickname in %s: %s\n", g.Name, *nickname)
					nicknameUpdates.Inc()
				}
			}
		}
		if *activity != "" {
			err = setActivity(dg, *activity, statusCode)
			if err != nil {
				log.Printf("Unable to set activity: %s\n", err)
			} else {
				log.Printf("Set activity: %s\n", *activity)
				activityUpdates.Inc()
			}
		}
	}
}

func setActivity(session *discordgo.Session, text string, code int) (err error) {
	switch code {
	case 0:
		err = session.UpdateGameStatus(0, text)
	case 1:
		err = session.UpdateListeningStatus(text)
	case 2:
		err = session.UpdateWatchStatus(0, text)
	}

	return
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}
