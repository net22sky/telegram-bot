package utils

import (
	"github.com/net22sky/telegram-bot/db/services"
	"log"
)

// GetLocalizedString получает локализованную строку по ключу.
func GetLocalizedString(l map[string]interface{}, key string) string {
	if value, exists := l[key]; exists {
		if strValue, ok := value.(string); ok {
			return strValue
		}
		log.Printf("Ошибка преобразования строки для ключа %s", key)
	}
	return "Строка не найдена"
}

// getUserLanguage возвращает язык пользователя.
// Параметры:
// - userID: ID пользователя.
// - userService: Сервис для работы с пользователями.
// - locales: Локализованные строки для разных языков.
// Возвращает:
// - Язык пользователя (по умолчанию "ru").
func GetUserLanguage(userID int64, userService *services.UserService) string {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		log.Printf("Ошибка при получении пользователя: %v", err)
		return "en"
	}
	if user.Language != "" {
		return user.Language
	}
	return "ru"
}
