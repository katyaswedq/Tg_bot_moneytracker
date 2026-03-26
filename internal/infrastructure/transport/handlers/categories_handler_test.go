package handlers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"tgfin/internal/service/dto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type categoryServiceMock struct{ mock.Mock }

func (m *categoryServiceMock) Add(ctx context.Context, in dto.CategoryAddInput) (*dto.CategoryAddOutput, error) {
	args := m.Called(ctx, in)
	out, _ := args.Get(0).(*dto.CategoryAddOutput)
	return out, args.Error(1)
}

func (m *categoryServiceMock) List(ctx context.Context, in dto.CategoryListInput) (*dto.CategoryListOutput, error) {
	args := m.Called(ctx, in)
	out, _ := args.Get(0).(*dto.CategoryListOutput)
	return out, args.Error(1)
}

func (m *categoryServiceMock) Delete(ctx context.Context, in dto.CategoryDeleteInput) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func TestCategoriesHandler_Handle_NilMessageOrFrom(t *testing.T) {
	svc := &categoryServiceMock{}
	h := NewCategoriesHandler(svc)

	res, err := h.Handle(context.Background(), tgbotapi.Update{})
	require.NoError(t, err)
	require.Nil(t, res)

	res, err = h.Handle(context.Background(), tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/categories",
			Chat: &tgbotapi.Chat{ID: 1},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 11},
			},
		},
	})
	require.NoError(t, err)
	require.Nil(t, res)
}

func TestCategoriesHandler_Handle_EmptyList(t *testing.T) {
	svc := &categoryServiceMock{}
	h := NewCategoriesHandler(svc)

	svc.On("List", mock.Anything, dto.CategoryListInput{TelegramID: 10}).
		Return(&dto.CategoryListOutput{Items: nil}, nil)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/categories",
			Chat: &tgbotapi.Chat{ID: 123},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 11},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, int64(123), res.ChatID)
	require.Contains(t, res.Text, "пока нет категорий")

	svc.AssertExpectations(t)
}

func TestCategoriesHandler_Handle_WithItems(t *testing.T) {
	svc := &categoryServiceMock{}
	h := NewCategoriesHandler(svc)

	id1 := uuid.New()
	id2 := uuid.New()

	svc.On("List", mock.Anything, dto.CategoryListInput{TelegramID: 10}).
		Return(&dto.CategoryListOutput{
			Items: []dto.CategoryListItem{
				{ID: id1, Name: "Еда"},
				{ID: id2, Name: "Транспорт"},
			},
		}, nil)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/categories",
			Chat: &tgbotapi.Chat{ID: 555},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 11},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, int64(555), res.ChatID)
	require.Contains(t, res.Text, "Ваши категории")
	require.Contains(t, res.Text, "Еда")
	require.Contains(t, res.Text, id1.String())
	require.Contains(t, res.Text, "Транспорт")
	require.Contains(t, res.Text, id2.String())
	require.Contains(t, res.Text, "/category delete")

	svc.AssertExpectations(t)
}
