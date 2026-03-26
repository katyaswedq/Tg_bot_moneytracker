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

type startUowMock struct{ mock.Mock }

func (m *startUowMock) WithinTx(ctx context.Context, fn func(ctx context.Context, exec user.Exec) error) error {
	m.Called(ctx)
	return fn(ctx, nil)
}

type startUserRepoMock struct{ mock.Mock }

func (m *startUserRepoMock) Create(ctx context.Context, exec user.Exec, u *models.User) error {
	args := m.Called(ctx, exec, u)
	return args.Error(0)
}

func (m *startUserRepoMock) GetByTelegramID(ctx context.Context, exec user.Exec, telegramID int64) (*models.User, error) {
	args := m.Called(ctx, exec, telegramID)
	u, _ := args.Get(0).(*models.User)
	return u, args.Error(1)
}

type startCategoryRepoMock struct{ mock.Mock }

func (m *startCategoryRepoMock) CreateDefault(ctx context.Context, exec user.Exec, userID uuid.UUID) error {
	args := m.Called(ctx, exec, userID)
	return args.Error(0)
}

func (m *startCategoryRepoMock) Create(ctx context.Context, exec user.Exec, userID uuid.UUID, name string) (*models.Category, error) {
	args := m.Called(ctx, exec, userID, name)
	c, _ := args.Get(0).(*models.Category)
	return c, args.Error(1)
}

func (m *startCategoryRepoMock) ListByUser(ctx context.Context, exec user.Exec, userID uuid.UUID) ([]*models.Category, error) {
	args := m.Called(ctx, exec, userID)
	items, _ := args.Get(0).([]*models.Category)
	return items, args.Error(1)
}

func (m *startCategoryRepoMock) Delete(ctx context.Context, exec user.Exec, userID uuid.UUID, categoryID uuid.UUID) error {
	args := m.Called(ctx, exec, userID, categoryID)
	return args.Error(0)
}

func (m *startCategoryRepoMock) GetByName(ctx context.Context, exec user.Exec, userID uuid.UUID, name string) (*models.Category, error) {
	args := m.Called(ctx, exec, userID, name)
	c, _ := args.Get(0).(*models.Category)
	return c, args.Error(1)
}

func TestStartService_Start_NewUser_CreatesDefaults(t *testing.T) {
	ctx := context.Background()

	uow := &startUowMock{}
	ur := &startUserRepoMock{}
	cr := &startCategoryRepoMock{}

	uow.On("WithinTx", mock.Anything).Return(nil)

	ur.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*models.User")).
		Run(func(args mock.Arguments) {
			u := args.Get(2).(*models.User)
			u.ID = uuid.New()
		}).
		Return(nil)

	cr.On("CreateDefault", mock.Anything, mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil)

	s := NewStartService(uow, ur, cr)

	err := s.Start(ctx, dto.StartInput{
		TelegramID: 1,
		UserName:   nil,
		FirstName:  "A",
	})
	require.NoError(t, err)

	uow.AssertExpectations(t)
	ur.AssertExpectations(t)
	cr.AssertExpectations(t)
}

func TestStartService_Start_UserAlreadyExists_DoesNothing(t *testing.T) {
	ctx := context.Background()

	uow := &startUowMock{}
	ur := &startUserRepoMock{}
	cr := &startCategoryRepoMock{}

	uow.On("WithinTx", mock.Anything).Return(nil)

	ur.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*models.User")).
		Return(domainerr.ErrUserAlreadyExists)

	s := NewStartService(uow, ur, cr)

	err := s.Start(ctx, dto.StartInput{
		TelegramID: 2,
		UserName:   nil,
		FirstName:  "B",
	})
	require.NoError(t, err)

	cr.AssertNotCalled(t, "CreateDefault", mock.Anything, mock.Anything, mock.Anything)

	uow.AssertExpectations(t)
	ur.AssertExpectations(t)
	cr.AssertExpectations(t)
}

func TestStartService_Start_UserCreateError_ReturnsError(t *testing.T) {
	ctx := context.Background()

	uow := &startUowMock{}
	ur := &startUserRepoMock{}
	cr := &startCategoryRepoMock{}

	uow.On("WithinTx", mock.Anything).Return(nil)

	ur.On("Create", mock.Anything, mock.Anything, mock.AnythingOfType("*models.User")).
		Return(assertErr{})

	s := NewStartService(uow, ur, cr)

	err := s.Start(ctx, dto.StartInput{
		TelegramID: 3,
		UserName:   nil,
		FirstName:  "C",
	})
	require.Error(t, err)

	uow.AssertExpectations(t)
	ur.AssertExpectations(t)
	cr.AssertExpectations(t)
}

type assertErr struct{}

func (assertErr) Error() string { return "boom" }
