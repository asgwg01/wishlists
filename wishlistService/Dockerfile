FROM golang:bookworm

WORKDIR /appWishlistService

COPY . .

RUN go mod tidy

RUN go build -o .bin/wishlistService cmd/wishlistService/main.go

CMD [ ".bin/wishlistService"]