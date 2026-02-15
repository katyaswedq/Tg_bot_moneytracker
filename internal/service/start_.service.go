package service

import (
	"context"
	"errors"
	"time"

	domainerr "tgfin/internal/domain/error"
	"tgfin/internal/domain/models"
	"tgfin/internal/domain/usecase"
	"tgfin/internal/service/dto"
)

type StartService struct {
	tx       user.UnitOfWork
	userRepo user.UserRepository
	catRepo  user.CategoryRepository
}

func NewStartService(tx user.UnitOfWork, userRepo user.UserRepository, catRepo user.CategoryRepository) *StartService {
	return &StartService{
		tx: tx,
		userRepo: userRepo,
		catRepo: catRepo,
	}
}

func (s *StartService) Start(ctx context.Context, in dto.StartInput) error {
	return s.tx.WithinTx(ctx, func(ctx context.Context, exec user.Exec) error {
		u := &models.User{
			TelegramID: in.TelegramID,
			UserName:   in.UserName,
			FirstName:  in.FirstName,
			CreatedAt:  time.Now(),
		}

		if err := s.userRepo.Create(ctx, exec, u); err != nil{
			if errors.Is(err, domainerr.ErrUserAlreadyExists){
				return nil
			}
			return err
		}

		return s.catRepo.CreateDefault(ctx, exec, u.ID)
	})
}