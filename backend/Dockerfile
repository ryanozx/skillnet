FROM golang:alpine

WORKDIR /usr/src/app/backend

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy

EXPOSE 8080