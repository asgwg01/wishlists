package handler

import (
	"context"
	"errors"
	wishlistv1 "pkg/proto/wishlists/v1"
	"pkg/types/price"
	"wishlistService/internal/domain/types"
	itemservice "wishlistService/internal/services/itemService"
	"wishlistService/internal/services/wishlistService"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var (
	ErrorEmailExist = errors.New("email already exist")
)

// GrpcHandler реализует сгенерированный интерфейс
type GrpcHandler struct {
	wishlistv1.UnimplementedWishlistServiceServer
	wishlistService wishlistService.IWshlistService
	itemService     itemservice.IItemService
}

func NewGrpcHandler(
	wishlistService wishlistService.IWshlistService,
	itemService itemservice.IItemService,
) *GrpcHandler {
	return &GrpcHandler{
		wishlistService: wishlistService,
		itemService:     itemService,
	}
}

func (h *GrpcHandler) CreateWishlist(ctx context.Context, req *wishlistv1.CreateWishlistRequest) (*wishlistv1.WishlistResponce, error) {
	err := validateCreateWishlistRequest(req)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	wishlist, err := h.wishlistService.Create(ctx, userID, req.GetTitle(), req.GetDescription(), req.GetIsPublic())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &wishlistv1.WishlistResponce{
		Id:          wishlist.ID.String(),
		OwnerId:     wishlist.OwnerID.String(),
		Title:       wishlist.Title,
		Description: wishlist.Description,
		IsPublic:    wishlist.IsPublic,
		CreatedAt:   wishlist.CreatedAt.Unix(),
		UpdatedAt:   wishlist.CreatedAt.Unix(),
	}, nil
}
func (h *GrpcHandler) GetWishlist(ctx context.Context, req *wishlistv1.GetWishlistRequest) (*wishlistv1.WishlistResponce, error) {
	err := validateGetWishlistRequest(req)
	if err != nil {
		return nil, err
	}
	wlID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	wishlist, err := h.wishlistService.Get(ctx, wlID, userID)
	if err != nil {
		switch err {
		case types.ErrorWishlistNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &wishlistv1.WishlistResponce{
		Id:          wishlist.ID.String(),
		OwnerId:     wishlist.OwnerID.String(),
		Title:       wishlist.Title,
		Description: wishlist.Description,
		IsPublic:    wishlist.IsPublic,
		CreatedAt:   wishlist.CreatedAt.Unix(),
		UpdatedAt:   wishlist.CreatedAt.Unix(),
	}, nil
}
func (h *GrpcHandler) GetUserWishlists(ctx context.Context, req *wishlistv1.GetUserWishlistsRequest) (*wishlistv1.WishlistListResponce, error) {
	err := validateGetUserWishlistsRequest(req)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	requestingUserId, err := uuid.Parse(req.GetRequestingUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	wishlists, err := h.wishlistService.GetUserWishlists(ctx, userID, requestingUserId)
	if err != nil {
		switch err {
		case types.ErrorWishlistNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	responceLists := make([]*wishlistv1.WishlistResponce, len(wishlists))
	for i, wl := range wishlists {
		responceLists[i] = &wishlistv1.WishlistResponce{
			Id:          wl.ID.String(),
			OwnerId:     wl.OwnerID.String(),
			Title:       wl.Title,
			Description: wl.Description,
			IsPublic:    wl.IsPublic,
			CreatedAt:   wl.CreatedAt.Unix(),
			UpdatedAt:   wl.CreatedAt.Unix(),
		}
	}

	return &wishlistv1.WishlistListResponce{
		Wishlists:  responceLists,
		TotalCount: int32(len(responceLists)),
	}, nil
}
func (h *GrpcHandler) UpdateWishlist(ctx context.Context, req *wishlistv1.UpdateWishlistRequest) (*wishlistv1.WishlistResponce, error) {
	err := validateUpdateWishlistRequest(req)
	if err != nil {
		return nil, err
	}
	wlID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	wishlist, err := h.wishlistService.Update(ctx, wlID, userID, req.GetTitle(), req.GetDescription(), req.GetIsPublic())
	if err != nil {
		switch err {
		case types.ErrorWishlistNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &wishlistv1.WishlistResponce{
		Id:          wishlist.ID.String(),
		OwnerId:     wishlist.OwnerID.String(),
		Title:       wishlist.Title,
		Description: wishlist.Description,
		IsPublic:    wishlist.IsPublic,
		CreatedAt:   wishlist.CreatedAt.Unix(),
		UpdatedAt:   wishlist.UpdatedAt.Unix(),
	}, nil
}
func (h *GrpcHandler) DeleteWishlist(ctx context.Context, req *wishlistv1.DeleteWishlistRequest) (*emptypb.Empty, error) {
	err := validateDeleteWishlistRequest(req)
	if err != nil {
		return nil, err
	}
	wlID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	err = h.wishlistService.Delete(ctx, wlID, userID)
	if err != nil {
		switch err {
		case types.ErrorWishlistNotFound:
			return &emptypb.Empty{}, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorAccessDenied:
			return &emptypb.Empty{}, status.Error(codes.PermissionDenied, "access denied")
		default:
			return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}
func (h *GrpcHandler) ListPublicWishlists(ctx context.Context, req *wishlistv1.ListPublicWishlistsRequest) (*wishlistv1.WishlistListResponce, error) {

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	wishlists, total, err := h.wishlistService.PublicWishlists(ctx, page, pageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responceLists := make([]*wishlistv1.WishlistResponce, len(wishlists))
	for i, wl := range wishlists {
		responceLists[i] = &wishlistv1.WishlistResponce{
			Id:          wl.ID.String(),
			OwnerId:     wl.OwnerID.String(),
			Title:       wl.Title,
			Description: wl.Description,
			IsPublic:    wl.IsPublic,
			CreatedAt:   wl.CreatedAt.Unix(),
			UpdatedAt:   wl.UpdatedAt.Unix(),
		}
	}

	return &wishlistv1.WishlistListResponce{
		Wishlists:  responceLists,
		TotalCount: int32(total),
	}, nil
}
func (h *GrpcHandler) AddItem(ctx context.Context, req *wishlistv1.AddItemRequest) (*wishlistv1.ItemResponce, error) {
	err := validateAddItemRequest(req)
	if err != nil {
		return nil, err
	}

	wlID, err := uuid.Parse(req.GetWishlistId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid wishlist_id format")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	item, err := h.itemService.Add(ctx, wlID, userID, req.GetName(), req.GetDescription(),
		req.GetImageUrl(), req.GetProductUrl(), price.Price{FullPrice: uint(req.GetPrice())})
	if err != nil {
		switch err {
		case types.ErrorWishlistNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	responce := wishlistv1.ItemResponce{
		Id:          item.ID.String(),
		WishlistId:  item.WishlistID.String(),
		Name:        item.Name,
		Description: item.Description,
		ImageUrl:    item.ImageURL,
		ProductUrl:  item.ProductURL,
		Price:       int64(item.Price.FullPriceKopecks()),
		CreatedAt:   item.CreatedAt.Unix(),
		UpdatedAt:   item.UpdatedAt.Unix(),
	}

	if item.BookedBy != nil {
		bby := item.BookedBy.String()
		responce.BookedBy = &bby
	}
	if item.BookedAt != nil {
		bat := item.BookedAt.Unix()
		responce.BookedAt = &bat
	}

	return &responce, nil
}
func (h *GrpcHandler) GetItem(ctx context.Context, req *wishlistv1.GetItemRequest) (*wishlistv1.ItemResponce, error) {
	err := validateGetItemRequest(req)
	if err != nil {
		return nil, err
	}

	itemID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item_id format")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	item, err := h.itemService.Get(ctx, itemID, userID)
	if err != nil {
		switch err {
		case types.ErrorItemNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	responce := wishlistv1.ItemResponce{
		Id:          item.ID.String(),
		WishlistId:  item.WishlistID.String(),
		Name:        item.Name,
		Description: item.Description,
		ImageUrl:    item.ImageURL,
		ProductUrl:  item.ProductURL,
		Price:       int64(item.Price.FullPriceKopecks()),
		CreatedAt:   item.CreatedAt.Unix(),
		UpdatedAt:   item.UpdatedAt.Unix(),
	}

	if item.BookedBy != nil {
		bby := item.BookedBy.String()
		responce.BookedBy = &bby
	}
	if item.BookedAt != nil {
		bat := item.BookedAt.Unix()
		responce.BookedAt = &bat
	}

	return &responce, nil
}
func (h *GrpcHandler) UpdateItem(ctx context.Context, req *wishlistv1.UpdateItemRequest) (*wishlistv1.ItemResponce, error) {
	err := validateUpdateItemRequest(req)
	if err != nil {
		return nil, err
	}

	itemID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item_id format")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	item, err := h.itemService.Update(ctx, itemID, userID, req.GetName(),
		req.GetDescription(),
		req.GetImageUrl(),
		req.GetProductUrl(),
		price.Price{FullPrice: uint(req.GetPrice())},
	)
	if err != nil {
		switch err {
		case types.ErrorItemNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	responce := wishlistv1.ItemResponce{
		Id:          item.ID.String(),
		WishlistId:  item.WishlistID.String(),
		Name:        item.Name,
		Description: item.Description,
		ImageUrl:    item.ImageURL,
		ProductUrl:  item.ProductURL,
		Price:       int64(item.Price.FullPriceKopecks()),
		CreatedAt:   item.CreatedAt.Unix(),
		UpdatedAt:   item.UpdatedAt.Unix(),
	}

	if item.BookedBy != nil {
		bby := item.BookedBy.String()
		responce.BookedBy = &bby
	}
	if item.BookedAt != nil {
		bat := item.BookedAt.Unix()
		responce.BookedAt = &bat
	}

	return &responce, nil
}
func (h *GrpcHandler) DeleteItem(ctx context.Context, req *wishlistv1.DeleteItemRequest) (*emptypb.Empty, error) {
	err := validateDeleteItemRequest(req)
	if err != nil {
		return nil, err
	}

	itemID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid item_id format")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	err = h.itemService.Delete(ctx, itemID, userID)
	if err != nil {
		switch err {
		case types.ErrorItemNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}
func (h *GrpcHandler) ListItems(ctx context.Context, req *wishlistv1.ListItemsRequest) (*wishlistv1.ItemListResponce, error) {
	err := validateListItemsRequest(req)
	if err != nil {
		return nil, err
	}

	wlID, err := uuid.Parse(req.GetWishlistId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid wishlist_id format")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	// Пагинация
	page := 1
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	items, total, err := h.itemService.List(ctx, wlID, userID, page, pageSize)
	if err != nil {
		switch err {
		case types.ErrorWishlistNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	responceLists := make([]*wishlistv1.ItemResponce, len(items))
	for i, item := range items {
		responce := wishlistv1.ItemResponce{
			Id:          item.ID.String(),
			WishlistId:  item.WishlistID.String(),
			Name:        item.Name,
			Description: item.Description,
			ImageUrl:    item.ImageURL,
			ProductUrl:  item.ProductURL,
			Price:       int64(item.Price.FullPriceKopecks()),
			CreatedAt:   item.CreatedAt.Unix(),
			UpdatedAt:   item.UpdatedAt.Unix(),
		}

		if item.BookedBy != nil {
			bby := item.BookedBy.String()
			responce.BookedBy = &bby
		}
		if item.BookedAt != nil {
			bat := item.BookedAt.Unix()
			responce.BookedAt = &bat
		}
		responceLists[i] = &responce
	}

	return &wishlistv1.ItemListResponce{
		Items:      responceLists,
		TotalCount: int32(total),
	}, nil
}
func (h *GrpcHandler) BookItem(ctx context.Context, req *wishlistv1.BookItemRequest) (*wishlistv1.ItemResponce, error) {
	err := validateBookItemRequest(req)
	if err != nil {
		return nil, err
	}

	itemID, err := uuid.Parse(req.GetItemId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid wishlist_id format")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	item, err := h.itemService.Book(ctx, itemID, userID)
	if err != nil {
		switch err {
		case types.ErrorItemNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorItemAlreadyBooked:
			return nil, status.Error(codes.FailedPrecondition, "item already booked")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	responce := wishlistv1.ItemResponce{
		Id:          item.ID.String(),
		WishlistId:  item.WishlistID.String(),
		Name:        item.Name,
		Description: item.Description,
		ImageUrl:    item.ImageURL,
		ProductUrl:  item.ProductURL,
		Price:       int64(item.Price.FullPriceKopecks()),
		CreatedAt:   item.CreatedAt.Unix(),
		UpdatedAt:   item.UpdatedAt.Unix(),
	}

	if item.BookedBy != nil {
		bby := item.BookedBy.String()
		responce.BookedBy = &bby
	}
	if item.BookedAt != nil {
		bat := item.BookedAt.Unix()
		responce.BookedAt = &bat
	}

	return &responce, nil
}
func (h *GrpcHandler) UnbookItem(ctx context.Context, req *wishlistv1.UnbookItemRequest) (*wishlistv1.ItemResponce, error) {
	err := validateUnbookItemRequest(req)
	if err != nil {
		return nil, err
	}

	itemID, err := uuid.Parse(req.GetItemId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid wishlist_id format")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	item, err := h.itemService.Unbook(ctx, itemID, userID)
	if err != nil {
		switch err {
		case types.ErrorItemNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorItemNotBooked:
			return nil, status.Error(codes.FailedPrecondition, "item not booked")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	responce := wishlistv1.ItemResponce{
		Id:          item.ID.String(),
		WishlistId:  item.WishlistID.String(),
		Name:        item.Name,
		Description: item.Description,
		ImageUrl:    item.ImageURL,
		ProductUrl:  item.ProductURL,
		Price:       int64(item.Price.FullPriceKopecks()),
		CreatedAt:   item.CreatedAt.Unix(),
		UpdatedAt:   item.UpdatedAt.Unix(),
	}

	if item.BookedBy != nil {
		bby := item.BookedBy.String()
		responce.BookedBy = &bby
	}
	if item.BookedAt != nil {
		bat := item.BookedAt.Unix()
		responce.BookedAt = &bat
	}

	return &responce, nil
}

func (h *GrpcHandler) GetBookings(ctx context.Context, req *wishlistv1.GetBookingsRequest) (*wishlistv1.BookingListResponce, error) {
	err := validateGetBookingsRequest(req)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	items, err := h.itemService.GetUserBookings(ctx, userId, page, pageSize)
	if err != nil {
		switch err {
		case types.ErrorItemNotFound:
			return nil, status.Error(codes.NotFound, "wishlist not found")
		case types.ErrorItemNotBooked:
			return nil, status.Error(codes.FailedPrecondition, "item not booked")
		case types.ErrorAccessDenied:
			return nil, status.Error(codes.PermissionDenied, "access denied")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// Преобразуем Item в Booking для ответа
	bookings := make([]*wishlistv1.Booking, len(items))
	for i, item := range items {
		bookings[i] = &wishlistv1.Booking{
			ItemId:     item.ID.String(),
			WishlistId: item.WishlistID.String(),
			UserId:     item.BookedBy.String(),
			BookedAt:   item.BookedAt.Unix(),
			ItemName:   item.Name,
		}
	}

	return &wishlistv1.BookingListResponce{
		Bookings: bookings,
	}, nil
}

func validateCreateWishlistRequest(req *wishlistv1.CreateWishlistRequest) error {
	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.GetTitle() == "" {
		return status.Error(codes.InvalidArgument, "title is required")
	}

	return nil
}

func validateGetWishlistRequest(req *wishlistv1.GetWishlistRequest) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "id is required")
	}

	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}

func validateGetUserWishlistsRequest(req *wishlistv1.GetUserWishlistsRequest) error {
	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.GetRequestingUserId() == "" {
		return status.Error(codes.InvalidArgument, "requestingUserId is required")
	}
	return nil
}

func validateUpdateWishlistRequest(req *wishlistv1.UpdateWishlistRequest) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "id is required")
	}

	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.GetTitle() == "" {
		return status.Error(codes.InvalidArgument, "title is required")
	}

	return nil
}

func validateDeleteWishlistRequest(req *wishlistv1.DeleteWishlistRequest) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "id is required")
	}

	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}

func validateAddItemRequest(req *wishlistv1.AddItemRequest) error {
	if req.GetWishlistId() == "" {
		return status.Error(codes.InvalidArgument, "wishlist_id is required")
	}
	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	return nil
}

func validateGetItemRequest(req *wishlistv1.GetItemRequest) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "id is required")
	}

	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}

func validateUpdateItemRequest(req *wishlistv1.UpdateItemRequest) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "id is required")
	}

	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	return nil
}

func validateDeleteItemRequest(req *wishlistv1.DeleteItemRequest) error {
	if req.GetId() == "" {
		return status.Error(codes.InvalidArgument, "id is required")
	}

	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}

func validateListItemsRequest(req *wishlistv1.ListItemsRequest) error {
	if req.GetWishlistId() == "" {
		return status.Error(codes.InvalidArgument, "wishlist_id is required")
	}

	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}

func validateBookItemRequest(req *wishlistv1.BookItemRequest) error {
	if req.GetItemId() == "" {
		return status.Error(codes.InvalidArgument, "item_id is required")
	}

	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}

func validateUnbookItemRequest(req *wishlistv1.UnbookItemRequest) error {
	if req.GetItemId() == "" {
		return status.Error(codes.InvalidArgument, "item_id is required")
	}

	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}

func validateGetBookingsRequest(req *wishlistv1.GetBookingsRequest) error {
	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
