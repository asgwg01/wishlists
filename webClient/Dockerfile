FROM golang:bookworm

WORKDIR /webClient

COPY . .

RUN go mod tidy

RUN go build -o .bin/webClient cmd/webClient/main.go

CMD [ ".bin/webClient"]