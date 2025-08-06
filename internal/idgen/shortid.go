package idgen

import (
	"crypto/rand"
	"math/big"
)

const (
	Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var alphabetSize = big.NewInt(int64(len(Alphabet)))

// GenerateShortID генерирует короткий идентификатор заданной длины
// Использует криптографически безопасный генератор
// Пример: для length=8 -> "A1B2C3D4"
func GenerateShortID(length int) string {
	if length < 0 {
		panic("idgen: negative length requested")
	}
	id := make([]byte, length)
	for i := range id {
		num, err := rand.Int(rand.Reader, alphabetSize)
		if err != nil {
			panic("idgen: failure: " + err.Error())
		}
		id[i] = Alphabet[num.Int64()]
	}
	return string(id)
}
