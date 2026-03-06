FROM golang:bookworm

WORKDIR /appAuthService

COPY . .

RUN go mod tidy
RUN make wishlist-authService-build

CMD [ "make", "wishlist-authService-start"]