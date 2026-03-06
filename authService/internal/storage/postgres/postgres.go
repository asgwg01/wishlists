package postgres

import (
	"authService/internal/config"
	"authService/internal/domain/models"
	"authService/internal/domain/utils"
	"authService/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

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

func (s *Storage) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	const logPrefix = "postgres.Storage.CreateUser"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Creare user", utils.UserToSlog(user))

	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Can not create transaction", slog.String("err", err.Error()))
		return models.User{}, err
	} else { // transaction

		exists, err := s.checkUserExist(tx, user)
		if err != nil {
			tx.Rollback()
			return models.User{}, err
		}
		if exists {
			log.Info("User with email is exist", slog.String("email", user.Email))
			return models.User{}, storage.ErrorEmailExist
		}

		query, err := tx.Prepare(`
		INSERT INTO
		users(
			email,
			name,
			password_hash,
			create_at,
			update_at
		)
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING uuid;
		`)

		if err != nil {
			tx.Rollback()
			return models.User{}, fmt.Errorf("INSERT error %s", err)
		}

		var newUuidStr string
		err = query.QueryRow(
			user.Email,
			user.Name,
			string(user.PasswordHash),
			user.CreateAt,
			user.UpdateAt,
		).Scan(&newUuidStr)
		if err != nil {
			tx.Rollback()
			log.Error("Error inserting user", utils.UserToSlog(user))
			return models.User{}, fmt.Errorf("INSERT error %s", err)
		}

		newUuid, err := uuid.Parse(newUuidStr)
		if err != nil {
			log.Error("Error parse uuid", slog.String("err", err.Error()))
			return models.User{}, fmt.Errorf("Error parse uuid %w", err)
		}

		err = tx.Commit()
		if err != nil {
			log.Error("error commit transaction", slog.String("err", err.Error()))
			return models.User{}, fmt.Errorf("error commit transaction %s", err)
		}

		log.Info("Creare user success!", utils.UserToSlog(user))

		user.ID = newUuid
	} // end transaction

	return user, nil
}
func (s *Storage) GetUserByID(ctx context.Context, id uuid.UUID) (models.User, error) {
	const logPrefix = "postgres.Storage.GetUserByID"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Get user", slog.String("uuid", id.String()))

	query, err := s.db.Prepare(`
	SELECT
		uuid,
		email,
		name,
		password_hash,
		create_at,
		update_at
	FROM users
	WHERE uuid = $1
	`)

	if err != nil {
		return models.User{}, fmt.Errorf("SELECT error %s", err)
	}

	result := models.User{}

	var queryUuidStr string
	err = query.QueryRow(id.String()).Scan(
		&queryUuidStr,
		&result.Email,
		&result.Name,
		&result.PasswordHash,
		&result.CreateAt,
		&result.UpdateAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("SELECT Query error", slog.String("err", storage.ErrorUserNotExist.Error()))
			return models.User{}, storage.ErrorUserNotExist
		} else {
			log.Error("SELECT Query error", slog.String("err", err.Error()))
			return models.User{}, err
		}
	}

	queryUuid, err := uuid.Parse(queryUuidStr)
	if err != nil {
		log.Error("Error parse uuid", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("Error parse uuid %w", err)
	}

	result.ID = queryUuid

	return result, nil
}
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const logPrefix = "postgres.Storage.GetUserByEmail"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Get user", slog.String("email", email))

	query, err := s.db.Prepare(`
	SELECT
		uuid,
		email,
		name,
		password_hash,
		create_at,
		update_at
	FROM users
	WHERE email = $1
	`)

	if err != nil {
		return models.User{}, fmt.Errorf("SELECT error %s", err)
	}

	result := models.User{}

	var queryUuidStr string
	err = query.QueryRow(email).Scan(
		&queryUuidStr,
		&result.Email,
		&result.Name,
		&result.PasswordHash,
		&result.CreateAt,
		&result.UpdateAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("SELECT Query error", slog.String("err", storage.ErrorUserNotExist.Error()))
			return models.User{}, storage.ErrorUserNotExist
		} else {
			log.Error("SELECT Query error", slog.String("err", err.Error()))
			return models.User{}, err
		}
	}

	queryUuid, err := uuid.Parse(queryUuidStr)
	if err != nil {
		log.Error("Error parse uuid", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("Error parse uuid %w", err)
	}

	result.ID = queryUuid

	return result, nil
}
func (s *Storage) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	const logPrefix = "postgres.Storage.UpdateUser"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Update user", utils.UserToSlog(user))

	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Can not create transaction", slog.String("err", err.Error()))
		return models.User{}, err
	} else { // transaction

		exists, err := s.checkUserExist(tx, user)
		if err != nil {
			tx.Rollback()
			return models.User{}, err
		}
		if !exists {
			log.Info("User with email not exist", slog.String("email", user.Email))
			return models.User{}, storage.ErrorUserNotExist
		}

		// Update Subs
		query, err := tx.Prepare(`
		UPDATE users
		SET 
		email = $2,
		name = $3,
		password_hash = $4,
		create_at = $5,
		update_at = $6
		WHERE uuid = $1;
		`)
		if err != nil {
			tx.Rollback()
			return models.User{}, fmt.Errorf("UPDATE error %s", err)
		}

		_, err = query.Exec(
			user.ID.String(),
			user.Email,
			user.Name,
			string(user.PasswordHash),
			user.CreateAt,
			user.UpdateAt,
		)
		if err != nil {
			tx.Rollback()
			log.Error("Error update user", utils.UserToSlog(user))
			return models.User{}, fmt.Errorf("UPDATE error %s", err)
		}

		err = tx.Commit()
		if err != nil {
			log.Error("error commit transaction", slog.String("err", err.Error()))
			return models.User{}, fmt.Errorf("error commit transaction %s", err)
		}

		log.Info("Update user success!", utils.UserToSlog(user))

	} // end transaction

	return user, nil
}
func (s *Storage) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const logPrefix = "postgres.Storage.DeleteUser"
	log := s.log.With(

		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)
	log.Info("Delete user", slog.String("uuid", id.String()))

	tx, err := s.db.Begin()
	if err != nil {
		log.Error("Can not create transaction", slog.String("err", err.Error()))
		return err
	} else { // transaction

		// check exist
		query, err := tx.Prepare(`
		SELECT EXISTS ( 
			SELECT 1
			FROM users
			WHERE uuid=$1
		);`)
		if err != nil {
			tx.Rollback()
			log.Error("Error check subscription exists", slog.String("uuid", id.String()))
			return storage.ErrorUserNotExist
		}

		var exists bool
		err = query.QueryRow(id.String()).Scan(&exists)
		if err != nil {
			tx.Rollback()
			log.Error("Error subscription check exists", slog.String("uuid", id.String()))
			return err
		}
		if !exists {
			tx.Rollback()
			log.Info("subscription is not exist", slog.String("uuid", id.String()))
			return storage.ErrorUserNotExist
		}

		// delete
		query, err = tx.Prepare(`
		DELETE FROM users
		WHERE uuid=$1;
		`)
		if err != nil {
			tx.Rollback()
			log.Error("Error delete subscription", slog.String("uuid", id.String()))
			return fmt.Errorf("DELETE error %s", err)
		}

		_, err = query.Exec(id.String())
		if err != nil {
			tx.Rollback()
			log.Error("Error delete subscription", slog.String("uuid", id.String()))
			return err
		}

		err = tx.Commit()
		if err != nil {
			log.Error("error commit transaction", slog.String("err", err.Error()))
			return fmt.Errorf("error commit transaction %s", err)
		}

		log.Info("Delete subscription success!", slog.String("uuid", id.String()))
	} // transaction end

	return nil
}

func (s *Storage) checkUserExist(tx *sql.Tx, user models.User) (bool, error) {
	const logPrefix = "postgres.Storage.checkUserExist"
	log := s.log.With(
		slog.String("where", logPrefix),
	)
	log.Info("check exist user", utils.UserToSlog(user))

	// Check User
	query, err := tx.Prepare(`
		SELECT EXISTS ( 
			SELECT 1
			FROM users
			WHERE email=$1
		);`)
	if err != nil {
		return false, fmt.Errorf("SELECT error %s", err)
	}

	var exists bool
	err = query.QueryRow(user.Email).Scan(&exists)
	if err != nil {
		log.Error("Error check user exists", slog.String("email", user.Email))
		return false, err
	}

	return exists, nil
}
