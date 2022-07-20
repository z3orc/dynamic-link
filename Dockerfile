FROM golang:alpine

WORKDIR /app

ADD . .
RUN go mod download

RUN go build ./main.go

CMD [ "./main" ]