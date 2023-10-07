package comments

import "strings"

// Метод для валидации комментария
func IsValidComment(text string) bool {
	// Запрещенные слова или фразы
	forbiddenWords := []string{"qwerty", "йцукен", "zxvbnm"}

	// Приводим текст комментария к нижнему регистру для регистронезависимой проверки
	normalizedText := strings.ToLower(text)

	for _, word := range forbiddenWords {
		if strings.Contains(normalizedText, word) {
			return false
		}
	}
	return true
}
