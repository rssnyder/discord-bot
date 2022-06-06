# discord-bot

a simple discord bot that just changes its nickname and/or status

![image](https://user-images.githubusercontent.com/7338312/172267762-4a725451-ac86-4f81-aa3a-6ddd88e7967c.png)

![GitHub all releases](https://img.shields.io/github/downloads/rssnyder/discord-bot/total?style=flat-square)

```text
Usage of ./discord-bot:
  -activity string
        bot activity
  -nickname string
        bot nickname
  -refresh int
        seconds between refresh (default 300)
  -status int
        0: playing, 1: listening
  -token string
        discord bot token
```

make sure the bot has "change nickname" permissions in the server if using that feature

## docker

### command line

```shell
docker run -e "TOKEN=XXX..XXX" -e "NICKNAME=some nickname" -e "ACTIVITY=some activity" -e "STATUS=0" -e "REFRESH=5" ghcr.io/rssnyder/discord-bot
```

### docker compose

```yaml
---
version: "3"
services:
  discord-bot:
    image: ghcr.io/rssnyder/discord-bot
    environment:
      TOKEN: XXX..XXX
      NICKNAME: some nickname
      ACTIVITY: some activity
      STATUS: 0
      REFRESH: 5
```

## command line

### download binary

grab a download link from [here](https://github.com/rssnyder/discord-bot/releases).
```shell
curl -L https://github.com/rssnyder/discord-bot/releases/download/v<version>/discord-bot_<version>_<os>_<arch>.tar.gz -o discord-bot.tar.gz
tar zxf discord-bot.tar.gz
```

### run

```shell
./discord-bot -token "XXX..XXX" -nickname "some nickname" -activity "some activity" -status "0" -refresh "5"
```
