FROM golang:alpine

WORKDIR /app

ADD . .
RUN go mod download

RUN go build -o /main.go

CMD [ "/main" ]