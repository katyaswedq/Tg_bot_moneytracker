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

type catUowMock struct{ mock.Mock }

func (m *catUowMock) WithinTx(ctx context.Context, fn func(ctx context.Context, exec user.Exec) error) error {
	m.Called(ctx)
	return fn(ctx, nil)
}

type catUserRepoMock struct{ mock.Mock }

func (m *catUserRepoMock) Create(ctx context.Context, exec user.Exec, u *models.User) error {
	args := m.Called(ctx, exec, u)
	return args.Error(0)
}

func (m *catUserRepoMock) GetByTelegramID(ctx context.Context, exec user.Exec, telegramID int64) (*models.User, error) {
	args := m.Called(ctx, exec, telegramID)
	u, _ := args.Get(0).(*models.User)
	return u, args.Error(1)
}

type catCategoryRepoMock struct{ mock.Mock }

func (m *catCategoryRepoMock) CreateDefault(ctx context.Context, exec user.Exec, userID uuid.UUID) error {
	args := m.Called(ctx, exec, userID)
	return args.Error(0)
}

func (m *catCategoryRepoMock) Create(ctx context.Context, exec user.Exec, userID uuid.UUID, name string) (*models.Category, error) {
	args := m.Called(ctx, exec, userID, name)
	c, _ := args.Get(0).(*models.Category)
	return c, args.Error(1)
}

func (m *catCategoryRepoMock) ListByUser(ctx context.Context, exec user.Exec, userID uuid.UUID) ([]*models.Category, error) {
	args := m.Called(ctx, exec, userID)
	items, _ := args.Get(0).([]*models.Category)
	return items, args.Error(1)
}

func (m *catCategoryRepoMock) Delete(ctx context.Context, exec user.Exec, userID uuid.UUID, categoryID uuid.UUID) error {
	args := m.Called(ctx, exec, userID, categoryID)
	return args.Error(0)
}

func (m *catCategoryRepoMock) GetByName(ctx context.Context, exec user.Exec, userID uuid.UUID, name string) (*models.Category, error) {
	args := m.Called(ctx, exec, userID, name)
	c, _ := args.Get(0).(*models.Category)
	return c, args.Error(1)
}

func TestCategoryService_Add_Success_TrimsName(t *testing.T) {
	ctx := context.Background()

	uow := &catUowMock{}
	ur := &catUserRepoMock{}
	cr := &catCategoryRepoMock{}

	uow.On("WithinTx", mock.Anything).Return(nil)

	uID := uuid.New()
	ur.On("GetByTelegramID", mock.Anything, mock.Anything, int64(5)).
		Return(&models.User{ID: uID, TelegramID: 5}, nil)

	catID := uuid.New()
	cr.On("Create", mock.Anything, mock.Anything, uID, "Спорт").
		Return(&models.Category{ID: catID, UserId: uID, Name: "Спорт"}, nil)

	s := NewCategoryService(uow, ur, cr)

	out, err := s.Add(ctx, dto.CategoryAddInput{
		TelegramID: 5,
		Name:       "  Спорт  ",
	})
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, catID, out.ID)
	require.Equal(t, "Спорт", out.Name)

	uow.AssertExpectations(t)
	ur.AssertExpectations(t)
	cr.AssertExpectations(t)
}

func TestCategoryService_Add_Duplicate(t *testing.T) {
	ctx := context.Background()

	uow := &catUowMock{}
	ur := &catUserRepoMock{}
	cr := &catCategoryRepoMock{}

	uow.On("WithinTx", mock.Anything).Return(nil)

	uID := uuid.New()
	ur.On("GetByTelegramID", mock.Anything, mock.Anything, int64(6)).
		Return(&models.User{ID: uID, TelegramID: 6}, nil)

	cr.On("Create", mock.Anything, mock.Anything, uID, "Еда").
		Return((*models.Category)(nil), domainerr.ErrCategoryAlreadyExists)

	s := NewCategoryService(uow, ur, cr)

	out, err := s.Add(ctx, dto.CategoryAddInput{
		TelegramID: 6,
		Name:       "Еда",
	})
	require.ErrorIs(t, err, domainerr.ErrCategoryAlreadyExists)
	require.Nil(t, out)

	uow.AssertExpectations(t)
	ur.AssertExpectations(t)
	cr.AssertExpectations(t)
}

func TestCategoryService_List_Success(t *testing.T) {
	ctx := context.Background()

	uow := &catUowMock{}
	ur := &catUserRepoMock{}
	cr := &catCategoryRepoMock{}

	uow.On("WithinTx", mock.Anything).Return(nil)

	uID := uuid.New()
	ur.On("GetByTelegramID", mock.Anything, mock.Anything, int64(7)).
		Return(&models.User{ID: uID, TelegramID: 7}, nil)

	c1 := &models.Category{ID: uuid.New(), UserId: uID, Name: "Еда"}
	c2 := &models.Category{ID: uuid.New(), UserId: uID, Name: "Транспорт"}

	cr.On("ListByUser", mock.Anything, mock.Anything, uID).
		Return([]*models.Category{c1, c2}, nil)

	s := NewCategoryService(uow, ur, cr)

	out, err := s.List(ctx, dto.CategoryListInput{TelegramID: 7})
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Len(t, out.Items, 2)
	require.Equal(t, c1.ID, out.Items[0].ID)
	require.Equal(t, "Еда", out.Items[0].Name)
	require.Equal(t, c2.ID, out.Items[1].ID)
	require.Equal(t, "Транспорт", out.Items[1].Name)

	uow.AssertExpectations(t)
	ur.AssertExpectations(t)
	cr.AssertExpectations(t)
}

func TestCategoryService_Delete_NotFound(t *testing.T) {
	ctx := context.Background()

	uow := &catUowMock{}
	ur := &catUserRepoMock{}
	cr := &catCategoryRepoMock{}

	uow.On("WithinTx", mock.Anything).Return(nil)

	uID := uuid.New()
	ur.On("GetByTelegramID", mock.Anything, mock.Anything, int64(8)).
		Return(&models.User{ID: uID, TelegramID: 8}, nil)

	cr.On("Delete", mock.Anything, mock.Anything, uID, mock.AnythingOfType("uuid.UUID")).
		Return(domainerr.ErrCategoryNotFound)

	s := NewCategoryService(uow, ur, cr)

	err := s.Delete(ctx, dto.CategoryDeleteInput{
		TelegramID: 8,
		CategoryID: uuid.New(),
	})
	require.ErrorIs(t, err, domainerr.ErrCategoryNotFound)

	uow.AssertExpectations(t)
	ur.AssertExpectations(t)
	cr.AssertExpectations(t)
}

func TestCategoryService_Delete_Success(t *testing.T) {
	ctx := context.Background()

	uow := &catUowMock{}
	ur := &catUserRepoMock{}
	cr := &catCategoryRepoMock{}

	uow.On("WithinTx", mock.Anything).Return(nil)

	uID := uuid.New()
	ur.On("GetByTelegramID", mock.Anything, mock.Anything, int64(9)).
		Return(&models.User{ID: uID, TelegramID: 9}, nil)

	catID := uuid.New()
	cr.On("Delete", mock.Anything, mock.Anything, uID, catID).
		Return(nil)

	s := NewCategoryService(uow, ur, cr)

	err := s.Delete(ctx, dto.CategoryDeleteInput{
		TelegramID: 9,
		CategoryID: catID,
	})
	require.NoError(t, err)

	uow.AssertExpectations(t)
	ur.AssertExpectations(t)
	cr.AssertExpectations(t)
}
