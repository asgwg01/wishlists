package client

import (
	"context"
	wishlistv1 "pkg/proto/wishlists/v1"
	"pkg/types/trace"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	//"google.golang.org/grpc/metadata"
)

type WishlistClient struct {
	client wishlistv1.WishlistServiceClient
	conn   *grpc.ClientConn
}

func NewWishlistClient(addr string, port string) (*WishlistClient, error) {
	conn, err := grpc.NewClient(
		addr+":"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	if err != nil {
		return nil, err
	}

	return &WishlistClient{
		client: wishlistv1.NewWishlistServiceClient(conn),
		conn:   conn,
	}, nil
}

// addUserToContext добавляет информацию о пользователе в metadata
// func (c *WishlistClient) addUserToContext(ctx context.Context, userID, userEmail, userName string) context.Context {
// 	md := metadata.Pairs(
// 		"user-id", userID,
// 		"user-email", userEmail,
// 		"user-name", userName,
// 	)
// 	return metadata.NewOutgoingContext(ctx, md)
// }

// Wishlist methods
func (c *WishlistClient) CreateWishlist(ctx context.Context, userID, userEmail, userName, title, description string, isPublic bool) (*wishlistv1.WishlistResponce, error) {
	ctx = trace.InjectIntoGRPC(ctx)

	//ctx = c.addUserToContext(ctx, userID, userEmail, userName)
	return c.client.CreateWishlist(ctx, &wishlistv1.CreateWishlistRequest{
		UserId:      userID,
		Title:       title,
		Description: description,
		IsPublic:    isPublic,
	})
}

func (c *WishlistClient) GetWishlist(ctx context.Context, userID, wishlistID string) (*wishlistv1.WishlistResponce, error) {
	ctx = trace.InjectIntoGRPC(ctx)

	return c.client.GetWishlist(ctx, &wishlistv1.GetWishlistRequest{
		Id:     wishlistID,
		UserId: userID,
	})
}

func (c *WishlistClient) GetUserWishlists(ctx context.Context, userID, requestingUserID string) (*wishlistv1.WishlistListResponce, error) {
	ctx = trace.InjectIntoGRPC(ctx)

	return c.client.GetUserWishlists(ctx, &wishlistv1.GetUserWishlistsRequest{
		UserId:           userID,
		RequestingUserId: requestingUserID,
	})
}

func (c *WishlistClient) UpdateWishlist(ctx context.Context, userID, wishlistID, title, description string, isPublic bool) (*wishlistv1.WishlistResponce, error) {
	ctx = trace.InjectIntoGRPC(ctx)

	return c.client.UpdateWishlist(ctx, &wishlistv1.UpdateWishlistRequest{
		Id:          wishlistID,
		UserId:      userID,
		Title:       title,
		Description: description,
		IsPublic:    isPublic,
	})
}

func (c *WishlistClient) DeleteWishlist(ctx context.Context, userID, wishlistID string) error {
	ctx = trace.InjectIntoGRPC(ctx)

	_, err := c.client.DeleteWishlist(ctx, &wishlistv1.DeleteWishlistRequest{
		Id:     wishlistID,
		UserId: userID,
	})
	return err
}

func (c *WishlistClient) ListPublicWishlists(ctx context.Context, page, pageSize int32) (*wishlistv1.WishlistListResponce, error) {
	ctx = trace.InjectIntoGRPC(ctx)

	return c.client.ListPublicWishlists(ctx, &wishlistv1.ListPublicWishlistsRequest{
		Page:     page,
		PageSize: pageSize,
	})
}

// Item methods
func (c *WishlistClient) AddItem(ctx context.Context, userID, userEmail, userName, wishlistID, name, description, imageURL, productURL string, price int64) (*wishlistv1.ItemResponce, error) {
	//ctx = c.addUserToContext(ctx, userID, userEmail, userName)
	ctx = trace.InjectIntoGRPC(ctx)
	return c.client.AddItem(ctx, &wishlistv1.AddItemRequest{
		WishlistId:  wishlistID,
		UserId:      userID,
		Name:        name,
		Description: description,
		ImageUrl:    imageURL,
		ProductUrl:  productURL,
		Price:       price,
	})
}

func (c *WishlistClient) GetItem(ctx context.Context, userID, itemID string) (*wishlistv1.ItemResponce, error) {
	ctx = trace.InjectIntoGRPC(ctx)
	return c.client.GetItem(ctx, &wishlistv1.GetItemRequest{
		Id:     itemID,
		UserId: userID,
	})
}

func (c *WishlistClient) ListItems(ctx context.Context, userID, wishlistID string, page, pageSize int32) (*wishlistv1.ItemListResponce, error) {
	ctx = trace.InjectIntoGRPC(ctx)
	return c.client.ListItems(ctx, &wishlistv1.ListItemsRequest{
		WishlistId: wishlistID,
		UserId:     userID,
		Page:       page,
		PageSize:   pageSize,
	})
}

func (c *WishlistClient) UpdateItem(ctx context.Context, userID, userEmail, userName, itemID, name, description, imageURL, productURL string, price int64) (*wishlistv1.ItemResponce, error) {
	//ctx = c.addUserToContext(ctx, userID, userEmail, userName)
	ctx = trace.InjectIntoGRPC(ctx)
	return c.client.UpdateItem(ctx, &wishlistv1.UpdateItemRequest{
		Id:          itemID,
		UserId:      userID,
		Name:        name,
		Description: description,
		ImageUrl:    imageURL,
		ProductUrl:  productURL,
		Price:       price,
	})
}

func (c *WishlistClient) DeleteItem(ctx context.Context, userID, itemID string) error {
	ctx = trace.InjectIntoGRPC(ctx)
	_, err := c.client.DeleteItem(ctx, &wishlistv1.DeleteItemRequest{
		Id:     itemID,
		UserId: userID,
	})
	return err
}

// Booking methods
func (c *WishlistClient) BookItem(ctx context.Context, userID, userEmail, userName, itemID string) (*wishlistv1.ItemResponce, error) {
	//ctx = c.addUserToContext(ctx, userID, userEmail, userName)
	ctx = trace.InjectIntoGRPC(ctx)
	return c.client.BookItem(ctx, &wishlistv1.BookItemRequest{
		ItemId: itemID,
		UserId: userID,
	})
}

func (c *WishlistClient) UnbookItem(ctx context.Context, userID, userEmail, userName, itemID string) (*wishlistv1.ItemResponce, error) {
	//ctx = c.addUserToContext(ctx, userID, userEmail, userName)
	ctx = trace.InjectIntoGRPC(ctx)
	return c.client.UnbookItem(ctx, &wishlistv1.UnbookItemRequest{
		ItemId: itemID,
		UserId: userID,
	})
}

func (c *WishlistClient) GetUserBookings(ctx context.Context, userID string) (*wishlistv1.BookingListResponce, error) {
	ctx = trace.InjectIntoGRPC(ctx)
	return c.client.GetBookings(ctx, &wishlistv1.GetBookingsRequest{
		UserId: userID,
	})
}

func (c *WishlistClient) Close() error {
	return c.conn.Close()
}
