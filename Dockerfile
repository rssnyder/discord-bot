FROM golang:1.19-alpine
LABEL org.opencontainers.image.source https://github.com/rssnyder/discord-bot

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /discord-bot

ENTRYPOINT /discord-bot -token "$TOKEN" -nickname "${NICKNAME:-I am a bot}" -activity "${ACTIVITY:-A special bot}" -status "${STATUS:-0}" -refresh "${REFRESH:-30}" -metrics "${METRICS:-:8080}"
