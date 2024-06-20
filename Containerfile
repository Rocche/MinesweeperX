FROM docker.io/golang:1.22.4-alpine3.20

ADD *.go go.mod .
ADD templates /go/templates

RUN go build .

ENTRYPOINT ["./MinesweeperX"]
