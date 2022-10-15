FROM golang:1.19.2-alpine

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./server ./server
COPY ./server/web/tmpl ./web/tmpl

RUN go build -v -o main server/cmd/main.go

CMD [ "./main" ]
