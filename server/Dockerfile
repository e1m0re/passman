FROM golang:1.23

ENV HTTP_SERVER_PORT=80
ENV HTTPS_SERVER_PORT=443

WORKDIR /go/server

COPY ./server/ ./

WORKDIR /go/server

RUN go build -o server ./cmd/server/main.go

EXPOSE ${HTTP_SERVER_PORT} ${HTTPS_SERVER_PORT}

ENTRYPOINT ./server