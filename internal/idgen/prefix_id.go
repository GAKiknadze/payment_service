package idgen

import "strings"

// GeneratePrefixedID генерирует идентификатор с префиксом
// Пример: GeneratePrefixedID("TAR", 8) -> "TAR-ABCD1234"
func GeneratePrefixedID(prefix string, length int) string {
	// Нормализуем префикс: заглавные буквы, удаляем недопустимые символы
	cleanPrefix := strings.ToUpper(prefix)
	cleanPrefix = strings.Map(func(r rune) rune {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, cleanPrefix)

	// Гарантируем, что префикс не пустой
	if cleanPrefix == "" {
		cleanPrefix = "ID"
	}

	return cleanPrefix + "-" + GenerateShortID(length)
}
