package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestHelpHandler_Handle_NilMessage(t *testing.T) {
	h := NewHelpHandler()

	res, err := h.Handle(context.Background(), tgbotapi.Update{})
	require.NoError(t, err)
	require.Nil(t, res)
}

func TestHelpHandler_Handle_OK(t *testing.T) {
	h := NewHelpHandler()

	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/help",
			Chat: &tgbotapi.Chat{ID: 123},
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 5},
			},
		},
	}

	res, err := h.Handle(context.Background(), upd)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, int64(123), res.ChatID)
	require.Contains(t, res.Text, "/add")
	require.Contains(t, res.Text, "/categories")
	require.Contains(t, res.Text, "/category add")
	require.Contains(t, res.Text, "/start")
	require.Contains(t, res.Text, "/help")
}
