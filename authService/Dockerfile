FROM golang:bookworm

WORKDIR /appAuthService

COPY . .

RUN go mod tidy

RUN go build -o .bin/authService cmd/authService/main.go

CMD [ ".bin/authService"]