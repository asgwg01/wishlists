FROM golang:bookworm

WORKDIR /appWishlistService

COPY . .

RUN go mod tidy
RUN make wishlist-wishlistService-build

CMD [ "make", "wishlist-wishlistService-start"]