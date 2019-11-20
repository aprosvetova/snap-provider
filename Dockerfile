FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
RUN adduser -D -g '' appuser

WORKDIR $GOPATH/src/github.com/aprosvetova/snap-provider
COPY . .

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/app .


FROM alpine

RUN apk update && apk add --no-cache ffmpeg

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/app /go/bin/app

USER appuser
ENTRYPOINT ["/go/bin/app"]