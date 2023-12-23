FROM golang:1.19-alpine as golang
LABEL org.opencontainers.image.source https://github.com/rssnyder/discord-bot

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /discord-bot

# FROM gcr.io/distroless/static-debian11

# COPY --from=golang /discord-bot .

ENTRYPOINT /discord-bot -token "$TOKEN" -nickname "$NICKNAME" -activity "$ACTIVITY" -status "${STATUS:-0}" -refresh "${REFRESH:-60}" -metrics "${METRICS:-:8080}"
