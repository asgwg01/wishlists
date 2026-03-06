FROM golang:bookworm

WORKDIR /appApiGateway

COPY . .

# RUN mkdir .bin
RUN go mod tidy
# RUN go build -o .bin/apiGateway cmd/apiGateway/main.go
RUN go build -v -o main ./cmd/apiGateway

CMD ["./main"]
# CMD [ ".bin/apiGateway" ]