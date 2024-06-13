FROM golang:1.22-alpine as base

LABEL org.opencontainers.image.source = "https://github.com/rssnyder/discord-bot"

RUN apk --update add ca-certificates

WORKDIR $GOPATH/src/discord-bot/app/

COPY . .

RUN go mod download
RUN go mod verify

RUN go build -o /discord-bot .

FROM scratch

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

COPY --from=base /discord-bot .

CMD ["./discord-bot"]
