FROM go1.16 AS Builder

WORKDIR /go/src/github.com/ezbuy/chartprobe
COPY . .

RUN CGO_ENABLED=0 go build -o chartprobe main.go

FROM alpine:latest

WORKDIR /app/chartprobe

COPY --from=GOBUILDER /go/src/github.com/ezbuy/chartprobe/chartprobe /app/chartprobe/chartprobe
