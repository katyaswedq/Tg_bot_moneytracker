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

type expenseServiceMock struct{ mock.Mock }

func (m *expenseServiceMock) Add(ctx context.Context, in dto.ExpenseAddInput) (*dto.ExpenseAddOutput, error) {
	args := m.Called(ctx, in)
	out, _ := args.Get(0).(*dto.ExpenseAddOutput)
	return out, args.Error(1)
}

func TestAddExpenseHandler_Handle_NilMessageOrFrom(t *testing.T) {
	svc := &expenseServiceMock{}
	h := NewAddExpenseHandler(svc)

	res, err := h.Handle(context.Background(), tgbotapi.Update{})
	require.NoError(t, err)
	require.Nil(t, res)

	res, err = h.Handle(context.Background(), tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/add 1 Еда",
			Chat: &tgbotapi.Chat{ID: 1},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 4},
			},
		},
	})
	require.NoError(t, err)
	require.Nil(t, res)
}

func TestAddExpenseHandler_Handle_Usage_NoArgs(t *testing.T) {
	h := NewAddExpenseHandler(&expenseServiceMock{})

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/add",
			Chat: &tgbotapi.Chat{ID: 1},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 4},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "Использование")
}

func TestAddExpenseHandler_Handle_Usage_TooFewArgs(t *testing.T) {
	h := NewAddExpenseHandler(&expenseServiceMock{})

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/add 500",
			Chat: &tgbotapi.Chat{ID: 1},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 4},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "Использование")
}

func TestAddExpenseHandler_Handle_AmountNotNumber(t *testing.T) {
	h := NewAddExpenseHandler(&expenseServiceMock{})

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/add abc Еда Обед",
			Chat: &tgbotapi.Chat{ID: 1},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 4},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "сумма должна быть числом")
}

func TestAddExpenseHandler_Handle_InvalidAmountDomainError(t *testing.T) {
	svc := &expenseServiceMock{}
	h := NewAddExpenseHandler(svc)

	svc.On("Add", mock.Anything, mock.Anything).Return((*dto.ExpenseAddOutput)(nil), domainerr.ErrInvalidAmount)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/add -1 Еда Обед",
			Chat: &tgbotapi.Chat{ID: 1},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 4},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, "положительным")

	svc.AssertExpectations(t)
}

func TestAddExpenseHandler_Handle_CategoryNotFound(t *testing.T) {
	svc := &expenseServiceMock{}
	h := NewAddExpenseHandler(svc)

	svc.On("Add", mock.Anything, mock.Anything).Return((*dto.ExpenseAddOutput)(nil), domainerr.ErrCategoryNotFound)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/add 500 Еда Обед",
			Chat: &tgbotapi.Chat{ID: 1},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 4},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Contains(t, res.Text, `Категория "Еда" не найдена`)

	svc.AssertExpectations(t)
}

func TestAddExpenseHandler_Handle_Success_WithDescription(t *testing.T) {
	svc := &expenseServiceMock{}
	h := NewAddExpenseHandler(svc)

	id := uuid.New()
	desc := "Обед в кафе"

	svc.On("Add", mock.Anything, dto.ExpenseAddInput{
		TelegramID:   10,
		Amount:       500,
		CategoryName: "Еда",
		Description:  &desc,
	}).Return(&dto.ExpenseAddOutput{
		ID:          id,
		Amount:      500,
		Category:    "Еда",
		Description: &desc,
	}, nil)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/add 500 Еда Обед в кафе",
			Chat: &tgbotapi.Chat{ID: 7},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 4},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, int64(7), res.ChatID)
	require.Contains(t, res.Text, "✅ Расход добавлен!")
	require.Contains(t, res.Text, "Сумма: 500")
	require.Contains(t, res.Text, "Категория: Еда")
	require.Contains(t, res.Text, "Описание: Обед в кафе")

	svc.AssertExpectations(t)
}

func TestAddExpenseHandler_Handle_Success_NoDescription(t *testing.T) {
	svc := &expenseServiceMock{}
	h := NewAddExpenseHandler(svc)

	id := uuid.New()

	svc.On("Add", mock.Anything, dto.ExpenseAddInput{
		TelegramID:   10,
		Amount:       300,
		CategoryName: "Транспорт",
		Description:  nil,
	}).Return(&dto.ExpenseAddOutput{
		ID:          id,
		Amount:      300,
		Category:    "Транспорт",
		Description: nil,
	}, nil)

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/add 300 Транспорт",
			Chat: &tgbotapi.Chat{ID: 8},
			From: &tgbotapi.User{ID: 10},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 4},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, int64(8), res.ChatID)
	require.Contains(t, res.Text, "✅ Расход добавлен!")
	require.Contains(t, res.Text, "Сумма: 300")
	require.Contains(t, res.Text, "Категория: Транспорт")
	require.NotContains(t, res.Text, "Описание:")

	svc.AssertExpectations(t)
}
