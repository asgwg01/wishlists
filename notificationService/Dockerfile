FROM golang:bookworm

WORKDIR /appnotificationService

COPY . .

RUN go mod tidy

RUN go build -o .bin/notificationService cmd/notificationService/main.go

CMD [ ".bin/notificationService"]