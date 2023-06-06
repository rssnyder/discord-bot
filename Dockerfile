FROM golang:1.19-alpine
LABEL org.opencontainers.image.source https://github.com/rssnyder/discord-bot

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /discord-bot

ENTRYPOINT /discord-bot -token "$TOKEN" -nickname "$NICKNAME" -activity "$ACTIVITY" -status "$STATUS" -refresh "$REFRESH"
