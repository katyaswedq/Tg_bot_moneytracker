package service

import (
	"context"
	"errors"
	"strings"

	domainerr "tgfin/internal/domain/error"
	"tgfin/internal/domain/models"
	"tgfin/internal/domain/usecase"
	"tgfin/internal/service/dto"
)

type ExpenseService struct {
	tx          user.UnitOfWork
	userRepo    user.UserRepository
	categoryRepo user.CategoryRepository
	expenseRepo user.ExpenseRepository
}

func NewExpenseService(tx user.UnitOfWork, userRepo user.UserRepository, categoryRepo user.CategoryRepository, expenseRepo user.ExpenseRepository) *ExpenseService {
	return &ExpenseService{
		tx:           tx,
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
		expenseRepo:  expenseRepo,
	}
}

func (s *ExpenseService) Add(ctx context.Context, in dto.ExpenseAddInput) (*dto.ExpenseAddOutput, error) {
	if in.Amount <= 0 {
		return nil, domainerr.ErrInvalidAmount
	}

	categoryName := strings.TrimSpace(in.CategoryName)
	if categoryName == "" {
		return nil, domainerr.ErrCategoryNotFound
	}

	if in.Description != nil {
		d := strings.TrimSpace(*in.Description)
		if d == "" {
			in.Description = nil
		} else {
			in.Description = &d
		}
	}

	var out *dto.ExpenseAddOutput

	err := s.tx.WithinTx(ctx, func(ctx context.Context, exec user.Exec) error {
		u, err := s.userRepo.GetByTelegramID(ctx, exec, in.TelegramID)
		if err != nil {
			return err
		}
		cat, err := s.categoryRepo.GetByName(ctx, exec, u.ID, categoryName)
		if err != nil {
			return err
		}

		e := models.NewExpense(u.ID, cat.ID, in.Amount, in.Description)
		if err := s.expenseRepo.Create(ctx, exec, e); err != nil {
			return err
		}

		out = &dto.ExpenseAddOutput{
			ID:          e.ID,
			Amount:      e.Amount,
			Category:    cat.Name,
			Description: e.Description,
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, domainerr.ErrCategoryNotFound) {
			return nil, domainerr.ErrCategoryNotFound
		}
		return nil, err
	}

	return out, nil
}
