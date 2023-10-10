FROM golang:1.20-alpine as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN air init

COPY . .

CMD air --build.bin="dlv exec --accept-multiclient --headless --listen=:2345 --continue --api-version=2 ./tmp/main"


