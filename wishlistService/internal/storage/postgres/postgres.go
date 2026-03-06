package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"wishlistService/internal/config"
	"wishlistService/internal/domain/models"
	"wishlistService/internal/domain/utils"
	"wishlistService/internal/storage"

	"github.com/asgwg01/wishlists/pkg/types/price"
	"github.com/asgwg01/wishlists/pkg/types/trace"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

type Storage struct {
	log *slog.Logger
	db  *sql.DB
}

func NewStorage(log *slog.Logger, cfg config.StorageConfig) (*Storage, error) {
	const logPrefix = "postgres.Storage.NewStorage"
	l := log.With(
		slog.String("where", logPrefix),
	)

	connectionStr := "postgres://" + cfg.User + ":" + cfg.Password +
		"@" + cfg.Host + ":" + cfg.Port +
		"/" + cfg.DBName +
		"?sslmode=disable"

	l.Debug("Create new psql conn", slog.String("str", connectionStr))

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {

		return nil, fmt.Errorf("can't open storage %s", err)
	}

	l.Info("Postgres connected", slog.String("host", cfg.Host), slog.String("port", cfg.Port))

	return &Storage{log: log, db: db}, nil
}

// IWishlistStorage
func (s *Storage) CreateWishlist(ctx context.Context, wishlist models.Wishlist) (models.Wishlist, error) {
	const logPrefix = "postgres.Storage.CreateWishlist"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Creare wishlist", utils.WishlistToSlog(wishlist))

	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Can not create transaction", slog.String("err", err.Error()))
		return models.Wishlist{}, err
	} else { // transaction

		exists, err := s.checkWishlistExist(tx, wishlist.ID)
		if err != nil {
			tx.Rollback()
			return models.Wishlist{}, err
		}
		if exists {
			log.Info("Wishlist already exist", slog.String("wishlist_id", wishlist.ID.String()))
			return models.Wishlist{}, storage.ErrorWishlistAlreadyExist
		}

		query, err := tx.Prepare(`
		INSERT INTO
		wishlists(
			owner_id,
			title,
			description,
			is_public,
			created_at,
			updated_at
		)
		VALUES
		($1, $2, $3, $4, $5, $6)
		RETURNING id;
		`)

		if err != nil {
			tx.Rollback()
			return models.Wishlist{}, fmt.Errorf("INSERT error %s", err)
		}

		var newUuidStr string
		err = query.QueryRow(
			wishlist.OwnerID.String(),
			wishlist.Title,
			wishlist.Description,
			wishlist.IsPublic,
			wishlist.CreatedAt,
			wishlist.UpdatedAt,
		).Scan(&newUuidStr)
		if err != nil {
			tx.Rollback()
			log.Error("Error inserting wishlist", utils.WishlistToSlog(wishlist))
			return models.Wishlist{}, fmt.Errorf("INSERT error %s", err)
		}

		newUuid, err := uuid.Parse(newUuidStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return models.Wishlist{}, fmt.Errorf("Error parse uuid %w", err)
		}

		err = tx.Commit()
		if err != nil {
			log.Error("error commit transaction", slog.String("err", err.Error()))
			return models.Wishlist{}, fmt.Errorf("error commit transaction %s", err)
		}

		wishlist.ID = newUuid

		log.Info("Creare wishlist success!", utils.WishlistToSlog(wishlist))

	} // end transaction

	return wishlist, nil
}
func (s *Storage) GetWishlistByID(ctx context.Context, id uuid.UUID) (models.Wishlist, error) {
	const logPrefix = "postgres.Storage.GetWishlistByID"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Get wishlist", slog.String("id", id.String()))

	query, err := s.db.Prepare(`
	SELECT
		id,
		owner_id,
		title,
		description,
		is_public,
		created_at,
		updated_at
	FROM wishlists
	WHERE id = $1
	`)

	if err != nil {
		return models.Wishlist{}, fmt.Errorf("SELECT error %s", err)
	}

	result := models.Wishlist{}

	var queryUuidWlStr string
	var queryUuidOwnerStr string
	err = query.QueryRow(id.String()).Scan(
		&queryUuidWlStr,
		&queryUuidOwnerStr,
		&result.Title,
		&result.Description,
		&result.IsPublic,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("SELECT Query error", slog.String("err", storage.ErrorWishlistNotExist.Error()))
			return models.Wishlist{}, storage.ErrorWishlistNotExist
		} else {
			log.Error("SELECT Query error", slog.String("err", err.Error()))
			return models.Wishlist{}, err
		}
	}

	queryUuidWl, err := uuid.Parse(queryUuidWlStr)
	if err != nil {
		log.Error("Error parse uuid", slog.String("err", err.Error()))
		return models.Wishlist{}, fmt.Errorf("Error parse uuid %w", err)
	}
	result.ID = queryUuidWl

	queryUuidOwner, err := uuid.Parse(queryUuidOwnerStr)
	if err != nil {
		log.Error("Error parse uuid", slog.String("err", err.Error()))
		return models.Wishlist{}, fmt.Errorf("Error parse uuid %w", err)
	}
	result.OwnerID = queryUuidOwner

	return result, nil
}
func (s *Storage) GetWishlistsByOwnerID(ctx context.Context, ownerID uuid.UUID, includePrivate bool) ([]models.Wishlist, error) {
	const logPrefix = "postgres.Storage.GetWishlistsByOwnerID"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Get wishlist", slog.String("owner_id", ownerID.String()))

	query, err := s.db.Prepare(`
	SELECT
		id,
		owner_id,
		title,
		description,
		is_public,
		created_at,
		updated_at
	FROM wishlists
	WHERE owner_id = $1
	`)

	if err != nil {
		return []models.Wishlist{}, fmt.Errorf("SELECT error %s", err)
	}

	result := []models.Wishlist{}

	rows, err := query.Query(ownerID.String())
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error("rows.Close error", slog.String("err", closeErr.Error()))
		}
	}()
	if err != nil {
		log.Error("SELECT Query error", slog.String("err", err.Error()))
		return []models.Wishlist{}, err

	}

	for rows.Next() {
		wl := models.Wishlist{}

		var queryUuidWlStr string
		var queryUuidOwnerStr string
		err := rows.Scan(
			&queryUuidWlStr,
			&queryUuidOwnerStr,
			&wl.Title,
			&wl.Description,
			&wl.IsPublic,
			&wl.CreatedAt,
			&wl.UpdatedAt,
		)
		if err != nil {
			log.Error("SELECT Query error", slog.String("err", err.Error()))
			return []models.Wishlist{}, err
		}

		queryUuidWl, err := uuid.Parse(queryUuidWlStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return []models.Wishlist{}, fmt.Errorf("Error parse uuid %w", err)
		}
		wl.ID = queryUuidWl

		queryUuidOwner, err := uuid.Parse(queryUuidOwnerStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return []models.Wishlist{}, fmt.Errorf("Error parse uuid %w", err)
		}
		wl.OwnerID = queryUuidOwner

		if wl.IsPublic || includePrivate {
			result = append(result, wl)
		}
	}
	return result, nil
}
func (s *Storage) UpdateWishlist(ctx context.Context, wishlist models.Wishlist) (models.Wishlist, error) {
	const logPrefix = "postgres.Storage.UpdateWishlist"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Update wishlist", utils.WishlistToSlog(wishlist))

	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Can not create transaction", slog.String("err", err.Error()))
		return models.Wishlist{}, err
	} else { // transaction

		exists, err := s.checkWishlistExist(tx, wishlist.ID)
		if err != nil {
			tx.Rollback()
			return models.Wishlist{}, err
		}
		if !exists {
			log.Info("Wishlist is not exist", slog.String("wishlist_id", wishlist.ID.String()))
			return models.Wishlist{}, storage.ErrorWishlistNotExist
		}

		query, err := tx.Prepare(`
		UPDATE wishlists
		SET
			title = $2,
			description = $3,
			is_public = $4,
			updated_at = $5
		WHERE
			id = $1
		`)

		if err != nil {
			tx.Rollback()
			return models.Wishlist{}, fmt.Errorf("UPDATE error %s", err)
		}

		_, err = query.Exec(
			wishlist.ID.String(),
			wishlist.Title,
			wishlist.Description,
			wishlist.IsPublic,
			wishlist.UpdatedAt,
		)
		if err != nil {
			tx.Rollback()
			log.Error("Error update wishlist", utils.WishlistToSlog(wishlist))
			return models.Wishlist{}, fmt.Errorf("UPDATE error %s", err)
		}

		err = tx.Commit()
		if err != nil {
			log.Error("error commit transaction", slog.String("err", err.Error()))
			return models.Wishlist{}, fmt.Errorf("error commit transaction %s", err)
		}

		log.Info("Update wishlist success!", utils.WishlistToSlog(wishlist))

	} // end transaction

	return wishlist, nil
}
func (s *Storage) DeleteWishlist(ctx context.Context, id uuid.UUID) error {
	const logPrefix = "postgres.Storage.DeleteWishlist"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Delete wishlist", slog.String("id", id.String()))

	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Can not create transaction", slog.String("err", err.Error()))
		return err
	} else { // transaction

		exists, err := s.checkWishlistExist(tx, id)
		if err != nil {
			tx.Rollback()
			return err
		}
		if !exists {
			log.Info("Wishlist is not exist", slog.String("wishlist_id", id.String()))
			return storage.ErrorWishlistNotExist
		}
		query, err := tx.Prepare(`
		DELETE FROM wishlists
		WHERE id = $1
		`)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("DELETE error %s", err)
		}

		_, err = query.Exec(
			id.String(),
		)
		if err != nil {
			tx.Rollback()
			log.Error("Error delete wishlist", slog.String("wishlist_id", id.String()))
			return fmt.Errorf("DELETE error %s", err)
		}

		err = tx.Commit()
		if err != nil {
			log.Error("error commit transaction", slog.String("err", err.Error()))
			return fmt.Errorf("error commit transaction %w", err)
		}

		log.Info("Delete wishlist success!", slog.String("wishlist_id", id.String()))

	} // end transaction

	return nil
}
func (s *Storage) ListPublicWishlists(ctx context.Context, limit, offset int) ([]models.Wishlist, int, error) {
	const logPrefix = "postgres.Storage.ListPublicWishlists"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Puplic wishlists", slog.Int("limit", limit), slog.Int("offset", limit))

	// Получаем общее количество
	query, err := s.db.Prepare(`
	SELECT COUNT(*) 
	FROM wishlists 
	WHERE is_public = true
	`)
	if err != nil {
		return []models.Wishlist{}, 0, fmt.Errorf("SELECT error %s", err)
	}

	var total int
	err = query.QueryRow().Scan(
		&total,
	)
	if err != nil {
		log.Error("SELECT Query error", slog.String("err", err.Error()))
		return []models.Wishlist{}, 0, err
	}

	// Получаем данные с пагинацией
	query, err = s.db.Prepare(`
	SELECT
		id,
		owner_id,
		title,
		description,
		is_public,
		created_at,
		updated_at
	FROM wishlists
    WHERE 
		is_public = true
    ORDER BY 
		created_at DESC
    LIMIT $1 
	OFFSET $2`)

	rows, err := query.Query(limit, offset)
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error("rows.Close error", slog.String("err", closeErr.Error()))
		}
	}()

	if err != nil {
		log.Error("SELECT Query error", slog.String("err", err.Error()))
		return []models.Wishlist{}, 0, err

	}
	defer rows.Close()

	results := []models.Wishlist{}

	for rows.Next() {
		wl := models.Wishlist{}

		var queryUuidWlStr string
		var queryUuidOwnerStr string
		err := rows.Scan(
			&queryUuidWlStr,
			&queryUuidOwnerStr,
			&wl.Title,
			&wl.Description,
			&wl.IsPublic,
			&wl.CreatedAt,
			&wl.UpdatedAt,
		)
		if err != nil {
			log.Error("SELECT Query error", slog.String("err", err.Error()))
			return []models.Wishlist{}, 0, err
		}

		queryUuidWl, err := uuid.Parse(queryUuidWlStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return []models.Wishlist{}, 0, fmt.Errorf("Error parse uuid %w", err)
		}
		wl.ID = queryUuidWl

		queryUuidOwner, err := uuid.Parse(queryUuidOwnerStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return []models.Wishlist{}, 0, fmt.Errorf("Error parse uuid %w", err)
		}
		wl.OwnerID = queryUuidOwner

		results = append(results, wl)
	}

	return results, total, nil
}

// IItemStorage
func (s *Storage) CreateItem(ctx context.Context, item models.Item) (models.Item, error) {
	const logPrefix = "postgres.Storage.CreateItem"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Creare item", utils.ItemToSlog(item))

	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Can not create transaction", slog.String("err", err.Error()))
		return models.Item{}, err
	} else { // transaction

		exists, err := s.checkItemExist(tx, item.ID)
		if err != nil {
			tx.Rollback()
			return models.Item{}, err
		}
		if exists {
			log.Info("Item already exist", slog.String("item_id", item.ID.String()))
			return models.Item{}, storage.ErrorItemAlreadyExist
		}

		exists, err = s.checkWishlistExist(tx, item.WishlistID)
		if err != nil {
			tx.Rollback()
			return models.Item{}, err
		}
		if !exists {
			log.Info("wishlist is not exist", slog.String("wishlist_id", item.WishlistID.String()))
			return models.Item{}, storage.ErrorWishlistNotExist
		}

		query, err := tx.Prepare(`
		INSERT INTO
		items(
			wishlist_id,
			name,
			description,
			image_url,
			product_url,
			price,
			created_at,
			updated_at
		)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;
		`)

		if err != nil {
			tx.Rollback()
			return models.Item{}, fmt.Errorf("INSERT error %s", err)
		}

		var newUuidStr string
		err = query.QueryRow(
			item.WishlistID.String(),
			item.Name,
			item.Description,
			item.ImageURL,
			item.ProductURL,
			item.Price.FullPriceKopecks(),
			item.CreatedAt,
			item.UpdatedAt,
		).Scan(&newUuidStr)
		if err != nil {
			tx.Rollback()
			log.Error("Error inserting item", utils.ItemToSlog(item))
			return models.Item{}, fmt.Errorf("INSERT error %s", err)
		}

		newUuid, err := uuid.Parse(newUuidStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return models.Item{}, fmt.Errorf("Error parse uuid %w", err)
		}

		err = tx.Commit()
		if err != nil {
			log.Error("error commit transaction", slog.String("err", err.Error()))
			return models.Item{}, fmt.Errorf("error commit transaction %s", err)
		}

		item.ID = newUuid

		log.Info("Creare item success!", utils.ItemToSlog(item))

	} // end transaction

	return item, nil
}
func (s *Storage) GetItemByID(ctx context.Context, id uuid.UUID) (models.Item, error) {
	const logPrefix = "postgres.Storage.GetItemByID"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Get item", slog.String("id", id.String()))

	query, err := s.db.Prepare(`
	SELECT
		id,
		wishlist_id,
		name,
		description,
		image_url,
		product_url,
		price,
		booked_by,
		booked_at,
		created_at,
		updated_at
	FROM items
	WHERE id = $1
	`)
	if err != nil {
		return models.Item{}, fmt.Errorf("SELECT error %s", err)
	}

	result := models.Item{}

	var idStr string
	var wlIdStr string
	var booked_by *string
	//var booked_at *time.Time
	var priceInt int
	err = query.QueryRow(id.String()).Scan(
		&idStr,
		&wlIdStr,
		&result.Name,
		&result.Description,
		&result.ImageURL,
		&result.ProductURL,
		&priceInt,
		&booked_by,
		&result.BookedAt,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("SELECT Query error", slog.String("err", storage.ErrorItemNotExist.Error()))
			return models.Item{}, storage.ErrorItemNotExist
		} else {
			log.Error("SELECT Query error", slog.String("err", err.Error()))
			return models.Item{}, err
		}
	}

	itemId, err := uuid.Parse(idStr)
	if err != nil {
		log.Error("Error parse uuid", slog.String("err", err.Error()))
		return models.Item{}, fmt.Errorf("Error parse uuid %w", err)
	}
	result.ID = itemId

	wlId, err := uuid.Parse(wlIdStr)
	if err != nil {
		log.Error("Error parse uuid", slog.String("err", err.Error()))
		return models.Item{}, fmt.Errorf("Error parse uuid %w", err)
	}
	result.WishlistID = wlId

	if booked_by != nil {
		bById, err := uuid.Parse(*booked_by)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return models.Item{}, fmt.Errorf("Error parse uuid %w", err)
		}
		result.BookedBy = &bById
	}

	result.Price = price.Price{FullPrice: uint(priceInt)}

	return result, nil
}
func (s *Storage) GetItemsByWishlistID(ctx context.Context, wishlistID uuid.UUID, limit, offset int) ([]models.Item, int, error) {
	const logPrefix = "postgres.Storage.GetItemsByWishlistID"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Get items", slog.String("wishlist_id", wishlistID.String()), slog.Int("limit", limit), slog.Int("offset", limit))

	// Получаем общее количество
	query, err := s.db.Prepare(`
	SELECT COUNT(*) 
	FROM items 
	WHERE wishlist_id = $1
	`)
	if err != nil {
		return []models.Item{}, 0, fmt.Errorf("SELECT error %s", err)
	}

	var total int
	err = query.QueryRow(wishlistID.String()).Scan(
		&total,
	)
	if err != nil {
		log.Error("SELECT Query error", slog.String("err", err.Error()))
		return []models.Item{}, 0, err
	}

	// Получаем данные с пагинацией
	query, err = s.db.Prepare(`
	SELECT
		id,
		wishlist_id,
		name,
		description,
		image_url,
		product_url,
		price,
		booked_by,
		booked_at,
		created_at,
		updated_at
	FROM items
    WHERE 
		wishlist_id = $1
    ORDER BY 
		created_at DESC
    LIMIT $2 
	OFFSET $3`)

	rows, err := query.Query(wishlistID.String(), limit, offset)
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error("rows.Close error", slog.String("err", closeErr.Error()))
		}
	}()

	if err != nil {
		log.Error("SELECT Query error", slog.String("err", err.Error()))
		return []models.Item{}, 0, err

	}
	defer rows.Close()

	results := []models.Item{}

	for rows.Next() {
		item := models.Item{}

		var idStr string
		var wlIdStr string
		var booked_by *string
		var priceInt int
		err = rows.Scan(
			&idStr,
			&wlIdStr,
			&item.Name,
			&item.Description,
			&item.ImageURL,
			&item.ProductURL,
			&priceInt,
			&booked_by,
			&item.BookedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("SELECT Query error", slog.String("err", storage.ErrorItemNotExist.Error()))
				return []models.Item{}, 0, storage.ErrorItemNotExist
			} else {
				log.Error("SELECT Query error", slog.String("err", err.Error()))
				return []models.Item{}, 0, err
			}
		}

		itemId, err := uuid.Parse(idStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return []models.Item{}, 0, fmt.Errorf("Error parse uuid %w", err)
		}
		item.ID = itemId

		wlId, err := uuid.Parse(wlIdStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return []models.Item{}, 0, fmt.Errorf("Error parse uuid %w", err)
		}
		item.WishlistID = wlId

		if booked_by != nil {
			bById, err := uuid.Parse(*booked_by)
			if err != nil {
				log.Error("Error parse uuid", slog.String("err", err.Error()))
				return []models.Item{}, 0, fmt.Errorf("Error parse uuid %w", err)
			}
			item.BookedBy = &bById
		}

		item.Price = price.Price{FullPrice: uint(priceInt)}

		results = append(results, item)
	}

	return results, total, nil
}

func (s *Storage) UpdateItem(ctx context.Context, item models.Item) (models.Item, error) {
	const logPrefix = "postgres.Storage.UpdateItem"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Update item", utils.ItemToSlog(item))

	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Can not create transaction", slog.String("err", err.Error()))
		return models.Item{}, err
	} else { // transaction

		exists, err := s.checkItemExist(tx, item.ID)
		if err != nil {
			tx.Rollback()
			return models.Item{}, err
		}
		if !exists {
			log.Info("Item not exist", slog.String("item_id", item.ID.String()))
			return models.Item{}, storage.ErrorItemNotExist
		}

		query, err := tx.Prepare(`
		UPDATE items
		SET
			name = $2,
			description = $3,
			image_url = $4,
			product_url = $5,
			price = $6,
			updated_at = $7
		WHERE
			id = $1
		`)
		if err != nil {
			tx.Rollback()
			return models.Item{}, fmt.Errorf("UPDATE error %s", err)
		}

		_, err = query.Exec(
			item.ID.String(),
			item.Name,
			item.Description,
			item.ImageURL,
			item.ProductURL,
			item.Price.FullPriceKopecks(),
			item.UpdatedAt,
		)
		err = tx.Commit()
		if err != nil {
			log.Error("error commit transaction", slog.String("err", err.Error()))
			return models.Item{}, fmt.Errorf("error commit transaction %s", err)
		}

		log.Info("Update item success!", utils.ItemToSlog(item))

	} // end transaction

	return item, nil
}
func (s *Storage) DeleteItem(ctx context.Context, id uuid.UUID) error {
	const logPrefix = "postgres.Storage.DeleteItem"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Delete item", slog.String("id", id.String()))

	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Can not create transaction", slog.String("err", err.Error()))
		return err
	} else { // transaction

		exists, err := s.checkItemExist(tx, id)
		if err != nil {
			tx.Rollback()
			return err
		}
		if !exists {
			log.Info("Item is not exist", slog.String("item_id", id.String()))
			return storage.ErrorWishlistNotExist
		}
		query, err := tx.Prepare(`
		DELETE FROM items
		WHERE id = $1
		`)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("DELETE error %s", err)
		}

		_, err = query.Exec(
			id.String(),
		)
		if err != nil {
			tx.Rollback()
			log.Error("Error delete item", slog.String("item_id", id.String()))
			return fmt.Errorf("DELETE error %s", err)
		}

		err = tx.Commit()
		if err != nil {
			log.Error("error commit transaction", slog.String("err", err.Error()))
			return fmt.Errorf("error commit transaction %w", err)
		}

		log.Info("Delete item success!", slog.String("item_id", id.String()))

	} // end transaction

	return nil
}

func (s *Storage) UpdateItemBooking(ctx context.Context, itemID uuid.UUID, bookedBy *uuid.UUID, bookedAt *time.Time) (models.Item, error) {
	const logPrefix = "postgres.Storage.UpdateItemBooking"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Update item booking", slog.String("item_id", itemID.String()), slog.Any("booked_by", bookedBy), slog.Any("booked_at", bookedAt))

	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Can not create transaction", slog.String("err", err.Error()))
		return models.Item{}, err
	} else { // transaction

		exists, err := s.checkItemExist(tx, itemID)
		if err != nil {
			tx.Rollback()
			return models.Item{}, err
		}
		if !exists {
			log.Info("Item not exist", slog.String("item_id", itemID.String()))
			return models.Item{}, storage.ErrorItemNotExist
		}

		query, err := tx.Prepare(`
		UPDATE items
		SET
			booked_by = $2,
			booked_at = $3,
			updated_at = $4
		WHERE
			id = $1
		RETURNING
			id,
			wishlist_id,
			name,
			description,
			image_url,
			product_url,
			price,
			booked_by,
			booked_at,
			created_at,
			updated_at;
		`)
		if err != nil {
			tx.Rollback()
			return models.Item{}, fmt.Errorf("UPDATE error %s", err)
		}

		result := models.Item{}

		var idStr string
		var wlIdStr string
		var booked_by *string
		//var booked_at *time.Time
		var priceInt int
		err = query.QueryRow(
			itemID.String(),
			bookedBy,
			bookedAt,
			time.Now(),
		).Scan(
			&idStr,
			&wlIdStr,
			&result.Name,
			&result.Description,
			&result.ImageURL,
			&result.ProductURL,
			&priceInt,
			&booked_by,
			&result.BookedAt,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("UPDATE Query error", slog.String("err", storage.ErrorItemNotExist.Error()))
				return models.Item{}, storage.ErrorItemNotExist
			} else {
				log.Error("UPDATE Query error", slog.String("err", err.Error()))
				return models.Item{}, err
			}
		}

		itemId, err := uuid.Parse(idStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return models.Item{}, fmt.Errorf("Error parse uuid %w", err)
		}
		result.ID = itemId

		wlId, err := uuid.Parse(wlIdStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return models.Item{}, fmt.Errorf("Error parse uuid %w", err)
		}
		result.WishlistID = wlId

		if booked_by != nil {
			bById, err := uuid.Parse(*booked_by)
			if err != nil {
				log.Error("Error parse uuid", slog.String("err", err.Error()))
				return models.Item{}, fmt.Errorf("Error parse uuid %w", err)
			}
			result.BookedBy = &bById
		}

		result.Price = price.Price{FullPrice: uint(priceInt)}

		err = tx.Commit()
		if err != nil {
			log.Error("error commit transaction", slog.String("err", err.Error()))
			return models.Item{}, fmt.Errorf("error commit transaction %s", err)
		}

		log.Info("Update item success!", slog.String("item_id", itemID.String()))
		return result, nil

	} // end transaction
}

func (s *Storage) GetBookedItemsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Item, error) {
	const logPrefix = "postgres.Storage.GetBookedItemsByUser"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Get booked items", slog.String("user_id", userID.String()))

	query, err := s.db.Prepare(`
	SELECT
		id,
		wishlist_id,
		name,
		description,
		image_url,
		product_url,
		price,
		booked_by,
		booked_at,
		created_at,
		updated_at
	FROM items
    WHERE 
		booked_by = $1
    ORDER BY 
		created_at DESC
	LIMIT $2 
	OFFSET $3
    `)

	rows, err := query.Query(userID.String(), limit, offset)
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Error("rows.Close error", slog.String("err", closeErr.Error()))
		}
	}()

	if err != nil {
		log.Error("SELECT Query error", slog.String("err", err.Error()))
		return []models.Item{}, err

	}
	defer rows.Close()

	results := []models.Item{}

	for rows.Next() {
		item := models.Item{}

		var idStr string
		var wlIdStr string
		var booked_by *string
		var priceInt int
		err = rows.Scan(
			&idStr,
			&wlIdStr,
			&item.Name,
			&item.Description,
			&item.ImageURL,
			&item.ProductURL,
			&priceInt,
			&booked_by,
			&item.BookedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("SELECT Query error", slog.String("err", storage.ErrorItemNotExist.Error()))
				return []models.Item{}, storage.ErrorItemNotExist
			} else {
				log.Error("SELECT Query error", slog.String("err", err.Error()))
				return []models.Item{}, err
			}
		}

		itemId, err := uuid.Parse(idStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return []models.Item{}, fmt.Errorf("Error parse uuid %w", err)
		}
		item.ID = itemId

		wlId, err := uuid.Parse(wlIdStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return []models.Item{}, fmt.Errorf("Error parse uuid %w", err)
		}
		item.WishlistID = wlId

		if booked_by != nil {
			bById, err := uuid.Parse(*booked_by)
			if err != nil {
				log.Error("Error parse uuid", slog.String("err", err.Error()))
				return []models.Item{}, fmt.Errorf("Error parse uuid %w", err)
			}
			item.BookedBy = &bById
		}

		item.Price = price.Price{FullPrice: uint(priceInt)}

		results = append(results, item)
	}

	return results, nil
}

func (s *Storage) checkWishlistExist(tx *sql.Tx, wlID uuid.UUID) (bool, error) {
	const logPrefix = "postgres.Storage.checkWishlistExist"
	log := s.log.With(
		slog.String("where", logPrefix),
	)
	log.Info("check exist wishlist", slog.String("wishlist_id", wlID.String()))

	query, err := tx.Prepare(`
		SELECT EXISTS ( 
			SELECT 1
			FROM wishlists
			WHERE id=$1
		);`)
	if err != nil {
		return false, fmt.Errorf("SELECT error %s", err)
	}

	var exists bool
	err = query.QueryRow(wlID.String()).Scan(&exists)
	if err != nil {
		log.Error("Error check wishlist exists", slog.String("id", wlID.String()))
		return false, err
	}

	return exists, nil
}

func (s *Storage) checkItemExist(tx *sql.Tx, itemID uuid.UUID) (bool, error) {
	const logPrefix = "postgres.Storage.checkWishlistExist"
	log := s.log.With(
		slog.String("where", logPrefix),
	)
	log.Info("check exist wishlist", slog.String("item_id", itemID.String()))

	query, err := tx.Prepare(`
		SELECT EXISTS ( 
			SELECT 1
			FROM items
			WHERE id=$1
		);`)
	if err != nil {
		return false, fmt.Errorf("SELECT error %s", err)
	}

	var exists bool
	err = query.QueryRow(itemID.String()).Scan(&exists)
	if err != nil {
		log.Error("Error check wishlist exists", slog.String("id", itemID.String()))
		return false, err
	}

	return exists, nil
}
