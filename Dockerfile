FROM golang:1.22

WORKDIR /go/src/builder

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /go/bin/app

RUN rm -rf /go/src/builder

CMD ["/go/bin/app"]