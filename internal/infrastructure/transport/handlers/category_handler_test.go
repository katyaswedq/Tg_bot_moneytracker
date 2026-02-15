package handlers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainerr "tgfin/internal/domain/error"
	"tgfin/internal/service/dto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type categoryServiceMock2 struct{ mock.Mock }

func (m *categoryServiceMock2) Add(ctx context.Context, in dto.CategoryAddInput) (*dto.CategoryAddOutput, error) {
	args := m.Called(ctx, in)
	out, _ := args.Get(0).(*dto.CategoryAddOutput)
	return out, args.Error(1)
}

func (m *categoryServiceMock2) List(ctx context.Context, in dto.CategoryListInput) (*dto.CategoryListOutput, error) {
	args := m.Called(ctx, in)
	out, _ := args.Get(0).(*dto.CategoryListOutput)
	return out, args.Error(1)
}

func (m *categoryServiceMock2) Delete(ctx context.Context, in dto.CategoryDeleteInput) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func TestCategoryHandler_Handle_NoArgs_ShowsUsage(t *testing.T) {
	svc := &categoryServiceMock2{}
	h := NewCategoryHandler(svc)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/category",
			Chat: &tgbotapi.Chat{ID: 1},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 9},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "Использование")
}

func TestCategoryHandler_Handle_Add_NoName(t *testing.T) {
	svc := &categoryServiceMock2{}
	h := NewCategoryHandler(svc)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/category add",
			Chat: &tgbotapi.Chat{ID: 1},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 9},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "Укажите название")
}

func TestCategoryHandler_Handle_Add_Success_NameWithSpaces(t *testing.T) {
	svc := &categoryServiceMock2{}
	h := NewCategoryHandler(svc)

	id := uuid.New()

	svc.On("Add", mock.Anything, dto.CategoryAddInput{
		TelegramID: 10,
		Name:       "Супермаркет у дома",
	}).Return(&dto.CategoryAddOutput{
		ID:   id,
		Name: "Супермаркет у дома",
	}, nil)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/category add Супермаркет у дома",
			Chat: &tgbotapi.Chat{ID: 123},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 9},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, int64(123), res.ChatID)
	require.Contains(t, res.Text, "✅ Категория создана")
	require.Contains(t, res.Text, "Супермаркет у дома")
	require.Contains(t, res.Text, id.String())

	svc.AssertExpectations(t)
}

func TestCategoryHandler_Handle_Add_Duplicate(t *testing.T) {
	svc := &categoryServiceMock2{}
	h := NewCategoryHandler(svc)

	svc.On("Add", mock.Anything, dto.CategoryAddInput{
		TelegramID: 10,
		Name:       "Еда",
	}).Return((*dto.CategoryAddOutput)(nil), domainerr.ErrCategoryAlreadyExists)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/category add Еда",
			Chat: &tgbotapi.Chat{ID: 123},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 9},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "уже существует")

	svc.AssertExpectations(t)
}

func TestCategoryHandler_Handle_Delete_NoID(t *testing.T) {
	svc := &categoryServiceMock2{}
	h := NewCategoryHandler(svc)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/category delete",
			Chat: &tgbotapi.Chat{ID: 1},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 9},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "Укажите ID")
}

func TestCategoryHandler_Handle_Delete_BadUUID(t *testing.T) {
	svc := &categoryServiceMock2{}
	h := NewCategoryHandler(svc)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/category delete not-a-uuid",
			Chat: &tgbotapi.Chat{ID: 1},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 9},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "Неверный формат ID")
}

func TestCategoryHandler_Handle_Delete_NotFound(t *testing.T) {
	svc := &categoryServiceMock2{}
	h := NewCategoryHandler(svc)

	id := uuid.New()

	svc.On("Delete", mock.Anything, dto.CategoryDeleteInput{
		TelegramID: 10,
		CategoryID: id,
	}).Return(domainerr.ErrCategoryNotFound)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/category delete " + id.String(),
			Chat: &tgbotapi.Chat{ID: 2},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 9},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "не найдена")

	svc.AssertExpectations(t)
}

func TestCategoryHandler_Handle_Delete_Success(t *testing.T) {
	svc := &categoryServiceMock2{}
	h := NewCategoryHandler(svc)

	id := uuid.New()

	svc.On("Delete", mock.Anything, dto.CategoryDeleteInput{
		TelegramID: 10,
		CategoryID: id,
	}).Return(nil)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/category delete " + id.String(),
			Chat: &tgbotapi.Chat{ID: 2},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 9},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "✅ Категория удалена")

	svc.AssertExpectations(t)
}
