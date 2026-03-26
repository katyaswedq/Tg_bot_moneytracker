package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainerr "tgfin/internal/domain/error"
	"tgfin/internal/domain/models"
	"tgfin/internal/domain/usecase"
	"tgfin/internal/service/dto"
)

type uowMock struct{ mock.Mock }

func (m *uowMock) WithinTx(ctx context.Context, fn func(ctx context.Context, exec user.Exec) error) error {
	m.Called(ctx)
	return fn(ctx, nil)
}

type userRepoMock struct{ mock.Mock }

func (m *userRepoMock) Create(ctx context.Context, exec user.Exec, u *models.User) error {
	args := m.Called(ctx, exec, u)
	return args.Error(0)
}

func (m *userRepoMock) GetByTelegramID(ctx context.Context, exec user.Exec, telegramID int64) (*models.User, error) {
	args := m.Called(ctx, exec, telegramID)
	u, _ := args.Get(0).(*models.User)
	return u, args.Error(1)
}

type categoryRepoMock struct{ mock.Mock }

func (m *categoryRepoMock) CreateDefault(ctx context.Context, exec user.Exec, userID uuid.UUID) error {
	args := m.Called(ctx, exec, userID)
	return args.Error(0)
}

func (m *categoryRepoMock) Create(ctx context.Context, exec user.Exec, userID uuid.UUID, name string) (*models.Category, error) {
	args := m.Called(ctx, exec, userID, name)
	c, _ := args.Get(0).(*models.Category)
	return c, args.Error(1)
}

func (m *categoryRepoMock) ListByUser(ctx context.Context, exec user.Exec, userID uuid.UUID) ([]*models.Category, error) {
	args := m.Called(ctx, exec, userID)
	items, _ := args.Get(0).([]*models.Category)
	return items, args.Error(1)
}

func (m *categoryRepoMock) Delete(ctx context.Context, exec user.Exec, userID uuid.UUID, categoryID uuid.UUID) error {
	args := m.Called(ctx, exec, userID, categoryID)
	return args.Error(0)
}

func (m *categoryRepoMock) GetByName(ctx context.Context, exec user.Exec, userID uuid.UUID, name string) (*models.Category, error) {
	args := m.Called(ctx, exec, userID, name)
	c, _ := args.Get(0).(*models.Category)
	return c, args.Error(1)
}

type expenseRepoMock struct{ mock.Mock }

func (m *expenseRepoMock) Create(ctx context.Context, exec user.Exec, e *models.Expense) error {
	args := m.Called(ctx, exec, e)
	return args.Error(0)
}

func TestExpenseService_Add_InvalidAmount(t *testing.T) {
	ctx := context.Background()

	s := NewExpenseService(&uowMock{}, &userRepoMock{}, &categoryRepoMock{}, &expenseRepoMock{})

	out, err := s.Add(ctx, dto.ExpenseAddInput{
		TelegramID:   1,
		Amount:       0,
		CategoryName: "Еда",
	})

	require.ErrorIs(t, err, domainerr.ErrInvalidAmount)
	require.Nil(t, out)
}

func TestExpenseService_Add_CategoryNotFound(t *testing.T) {
	ctx := context.Background()

	uow := &uowMock{}
	ur := &userRepoMock{}
	cr := &categoryRepoMock{}
	er := &expenseRepoMock{}

	uow.On("WithinTx", mock.Anything).Return(nil)

	uID := uuid.New()
	ur.On("GetByTelegramID", mock.Anything, mock.Anything, int64(10)).
		Return(&models.User{ID: uID, TelegramID: 10}, nil)

	cr.On("GetByName", mock.Anything, mock.Anything, uID, "Еда").
		Return((*models.Category)(nil), domainerr.ErrCategoryNotFound)

	s := NewExpenseService(uow, ur, cr, er)

	out, err := s.Add(ctx, dto.ExpenseAddInput{
		TelegramID:   10,
		Amount:       500,
		CategoryName: "Еда",
		Description:  nil,
	})

	require.ErrorIs(t, err, domainerr.ErrCategoryNotFound)
	require.Nil(t, out)

	uow.AssertExpectations(t)
	ur.AssertExpectations(t)
	cr.AssertExpectations(t)
	er.AssertExpectations(t)
}

func TestExpenseService_Add_Success(t *testing.T) {
	ctx := context.Background()

	uow := &uowMock{}
	ur := &userRepoMock{}
	cr := &categoryRepoMock{}
	er := &expenseRepoMock{}

	uow.On("WithinTx", mock.Anything).Return(nil)

	uID := uuid.New()
	ur.On("GetByTelegramID", mock.Anything, mock.Anything, int64(77)).
		Return(&models.User{ID: uID, TelegramID: 77}, nil)

	catID := uuid.New()
	cr.On("GetByName", mock.Anything, mock.Anything, uID, "Еда").
		Return(&models.Category{ID: catID, UserId: uID, Name: "Еда"}, nil)

	desc := "Обед"
	er.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*models.Expense")).
		Run(func(args mock.Arguments) {
			e := args.Get(2).(*models.Expense)
			e.ID = uuid.New()
		}).
		Return(nil)

	s := NewExpenseService(uow, ur, cr, er)

	out, err := s.Add(ctx, dto.ExpenseAddInput{
		TelegramID:   77,
		Amount:       500,
		CategoryName: "Еда",
		Description:  &desc,
	})

	require.NoError(t, err)
	require.NotNil(t, out)
	require.NotEqual(t, uuid.Nil, out.ID)
	require.Equal(t, int64(500), out.Amount)
	require.Equal(t, "Еда", out.Category)
	require.NotNil(t, out.Description)
	require.Equal(t, "Обед", *out.Description)

	uow.AssertExpectations(t)
	ur.AssertExpectations(t)
	cr.AssertExpectations(t)
	er.AssertExpectations(t)
}
