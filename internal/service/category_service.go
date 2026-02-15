package service

import (
	"context"
	"strings"

	"tgfin/internal/domain/usecase"
	"tgfin/internal/service/dto"
)

type CategoryService struct {
	tx       user.UnitOfWork
	userRepo user.UserRepository
	catRepo  user.CategoryRepository
}

func NewCategoryService(tx user.UnitOfWork, userRepo user.UserRepository, catRepo user.CategoryRepository) *CategoryService {
	return &CategoryService{
		tx:       tx,
		userRepo: userRepo,
		catRepo:  catRepo,
	}
}

func (s *CategoryService) Add(ctx context.Context, in dto.CategoryAddInput) (*dto.CategoryAddOutput, error) {
	name := strings.TrimSpace(in.Name)

	var out *dto.CategoryAddOutput

	err := s.tx.WithinTx(ctx, func(ctx context.Context, exec user.Exec) error {
		u, err := s.userRepo.GetByTelegramID(ctx, exec, in.TelegramID)
		if err != nil {
			return err
		}

		cat, err := s.catRepo.Create(ctx, exec, u.ID, name)
		if err != nil {
			return err
		}

		out = &dto.CategoryAddOutput{
			ID:   cat.ID,
			Name: cat.Name,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *CategoryService) List(ctx context.Context, in dto.CategoryListInput) (*dto.CategoryListOutput, error) {
	var out *dto.CategoryListOutput

	err := s.tx.WithinTx(ctx, func(ctx context.Context, exec user.Exec) error {
		u, err := s.userRepo.GetByTelegramID(ctx, exec, in.TelegramID)
		if err != nil {
			return err
		}

		items, err := s.catRepo.ListByUser(ctx, exec, u.ID)
		if err != nil {
			return err
		}

		res := make([]dto.CategoryListItem, 0, len(items))
		for _, c := range items {
			res = append(res, dto.CategoryListItem{
				ID:   c.ID,
				Name: c.Name,
			})
		}

		out = &dto.CategoryListOutput{Items: res}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s *CategoryService) Delete(ctx context.Context, in dto.CategoryDeleteInput) error {
	return s.tx.WithinTx(ctx, func(ctx context.Context, exec user.Exec) error{
		u, err := s.userRepo.GetByTelegramID(ctx, exec, in.TelegramID)
		if err != nil {
			return err
		}
		return s.catRepo.Delete(ctx, exec, u.ID, in.CategoryID)
	})
}
