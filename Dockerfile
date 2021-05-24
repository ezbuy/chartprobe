FROM golang:1.16 AS BUILDER

WORKDIR /go/src/github.com/ezbuy/chartprobe
COPY . .

RUN CGO_ENABLED=0 go build -o chartprobe main.go

FROM alpine:latest

WORKDIR /app/chartprobe

COPY --from=BUILDER /go/src/github.com/ezbuy/chartprobe/chartprobe /app/chartprobe/chartprobe
