FROM golang:1.22-alpine as base

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

FROM base as local

ENV GIN_MODE=debug

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN air init

CMD air --build.bin="dlv exec --accept-multiclient --headless --listen=:2345 --continue --api-version=2 ./tmp/main"

FROM base as debug

ENV GIN_MODE=debug

CMD ["go", "run", "main.go"]

FROM base as production

ENV GIN_MODE=release

CMD ["go", "run", "main.go"]