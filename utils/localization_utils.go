package utils

import "log"

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
