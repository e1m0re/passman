FROM golang:1.23

ENV GRPC_SERVER_PORT=5000

WORKDIR /go/server

COPY ./ ./

WORKDIR /go/server

RUN go build -o server ./cmd/server/main.go

EXPOSE ${GRPC_SERVER_PORT}

VOLUME /passman_data

ENTRYPOINT ./server