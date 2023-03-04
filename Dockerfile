FROM golang:1.19-alpine as build

WORKDIR /tmp/tail-time

COPY cmd/tail-time ./cmd/tail-time
COPY internal ./internal
COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go build -o build/tail-time cmd/tail-time/main.go

FROM alpine:3

COPY --from=build /tmp/tail-time/build/tail-time /opt/tail-time

RUN adduser -D tail-time
USER tail-time:tail-time

ENTRYPOINT ["/opt/tail-time"]
