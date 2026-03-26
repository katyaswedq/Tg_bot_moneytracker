package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"tgfin/internal/service/dto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type startServiceMock struct{ mock.Mock }

func (m *startServiceMock) Start(ctx context.Context, in dto.StartInput) error {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func TestStartHandler_Handle_NilMessageOrFrom(t *testing.T) {
	svc := &startServiceMock{}
	h := NewStartHandler(svc)

	res, err := h.Handle(context.Background(), tgbotapi.Update{})
	require.NoError(t, err)
	require.Nil(t, res)

	res, err = h.Handle(context.Background(), tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/start",
			Chat: &tgbotapi.Chat{ID: 1},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 6},
			},
		},
	})
	require.NoError(t, err)
	require.Nil(t, res)
}

func TestStartHandler_Handle_OK_WithUsername(t *testing.T) {
	svc := &startServiceMock{}
	h := NewStartHandler(svc)

	username := "myuser"

	svc.On("Start", mock.Anything, dto.StartInput{
		TelegramID: 10,
		UserName:   &username,
		FirstName:  "Ann",
	}).Return(nil)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/start",
			Chat: &tgbotapi.Chat{ID: 123},
			From: &tgbotapi.User{
				ID:        10,
				UserName:  username,
				FirstName: "Ann",
			},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 6},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, int64(123), res.ChatID)
	require.Contains(t, res.Text, "Добро пожаловать")
	require.Contains(t, res.Text, "Созданы базовые категории")

	svc.AssertExpectations(t)
}

func TestStartHandler_Handle_OK_WithoutUsername(t *testing.T) {
	svc := &startServiceMock{}
	h := NewStartHandler(svc)

	svc.On("Start", mock.Anything, dto.StartInput{
		TelegramID: 10,
		UserName:   nil,
		FirstName:  "Ann",
	}).Return(nil)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/start",
			Chat: &tgbotapi.Chat{ID: 123},
			From: &tgbotapi.User{
				ID:        10,
				UserName:  "",
				FirstName: "Ann",
			},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 6},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, int64(123), res.ChatID)
	require.Contains(t, res.Text, "✅ Вы зарегистрированы")

	svc.AssertExpectations(t)
}

func TestStartHandler_Handle_ServiceError_ReturnsUserMessageAndError(t *testing.T) {
	svc := &startServiceMock{}
	h := NewStartHandler(svc)

	svc.On("Start", mock.Anything, mock.Anything).Return(assertErrStart{})

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/start",
			Chat: &tgbotapi.Chat{ID: 111},
			From: &tgbotapi.User{ID: 10, FirstName: "Ann"},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 6},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.Error(t, err)
	require.NotNil(t, res)
	require.Equal(t, int64(111), res.ChatID)
	require.Contains(t, res.Text, "Произошла ошибка")

	svc.AssertExpectations(t)
}

type assertErrStart struct{}

func (assertErrStart) Error() string { return "boom" }
