FROM golang:bookworm

WORKDIR /appApiGateway

COPY . .

RUN go mod tidy

RUN go build -o .bin/apiGateway cmd/apiGateway/main.go

CMD [ ".bin/apiGateway"]