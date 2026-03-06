package dowork

import (
	"context"
	"log/slog"
	"time"
)

type WorkService struct {
	log *slog.Logger
	//storage    ISomeStorage

}

type IServiceWorkSome interface {
	WorkSome(ctx context.Context,
		str string,
	) (ok bool, err error)
}

// New return a new instance of the Auth service (service layer)
func New(
	log *slog.Logger,
	//storage ISomeStorage,
) *WorkService {
	return &WorkService{
		log: log,
		//storage:    storage,
	}
}

func (a *WorkService) WorkSome(
	ctx context.Context,
	str string,
) (bool, error) {
	const logPrefix = "auth.WorkSome"
	log := a.log.With(
		slog.String("where", logPrefix),
		slog.String("str", str),
	)

	log.Info("Start do something")

	log.Info("Do something...")
	time.Sleep(1 * time.Second)

	log.Info("End do something")

	return true, nil
}
