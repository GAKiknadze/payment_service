package idgen

import (
	"crypto/rand"
	"math/big"
)

const (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// GenerateShortID генерирует короткий идентификатор заданной длины
// Использует криптографически безопасный генератор
// Пример: для length=8 -> "A1B2C3D4"
func GenerateShortID(length int) string {
	result := make([]byte, length)
	alphabetSize := big.NewInt(int64(len(alphabet)))

	for i := 0; i < length; i++ {
		// Получаем безопасное случайное число в диапазоне алфавита
		num, err := rand.Int(rand.Reader, alphabetSize)
		if err != nil {
			// В реальном приложении здесь должна быть обработка ошибки
			return fallbackShortID(length)
		}
		result[i] = alphabet[num.Int64()]
	}

	return string(result)
}

// fallbackShortID используется при неудаче с crypto/rand
func fallbackShortID(length int) string {
	// В реальном приложении здесь должна быть более надежная реализация
	// или panic с логированием критической ошибки
	result := make([]byte, length)
	for i := range result {
		result[i] = 'X'
	}
	return string(result)
}
