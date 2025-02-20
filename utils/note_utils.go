package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/net22sky/telegram-bot/db/services"
	"log"
	"strconv"
	"strings"
)

// DeleteNote удаляет заметку по её ID, если она принадлежит указанному пользователю.
func DeleteNote(bot *tgbotapi.BotAPI, message *tgbotapi.Message, l map[string]interface{}, noteService *services.NoteService) {
	parts := strings.SplitN(message.Text, " ", 2)
	if len(parts) < 2 || parts[1] == "" {
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "unknown_command"))
		return
	}

	noteIDStr := parts[1]
	userID := message.From.ID
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil || noteID <= 0 {
		log.Printf("Ошибка при преобразовании ID заметки: %v", err)
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "invalid_note_id"))
		return
	}

	err = noteService.DeleteNoteByID(uint(noteID), int64(userID))
	if err != nil {
		log.Printf("Ошибка при удалении заметки: %v", err)
		SendMessage(bot, message.Chat.ID, GetLocalizedString(l, "note_deletion_error"))
		return
	}

	SendMessage(bot, message.Chat.ID, fmt.Sprintf(GetLocalizedString(l, "note_deleted"), noteID))
}
