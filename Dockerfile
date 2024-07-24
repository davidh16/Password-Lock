FROM golang:1.22-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM base as local

ENV GIN_MODE=debug

RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN air init

CMD air --build.bin="dlv exec --accept-multiclient --headless --listen=:2345 --continue --api-version=2 ./tmp/main"

FROM base as build

RUN go build -o runnable .

FROM gcr.io/distroless/static:nonroot AS debug

ENV GIN_MODE=debug

WORKDIR /
COPY --from=build /app/runnable /app/runnable

ENTRYPOINT ["/app/runnable"]

FROM gcr.io/distroless/static:nonroot AS production

ENV GIN_MODE=release

WORKDIR /
COPY --from=build /app/runnable /app/runnable

ENTRYPOINT ["/app/runnable"]