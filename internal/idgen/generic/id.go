package generic

import (
	"strings"

	"github.com/GAKiknadze/payment_service/internal/idgen"
)

// Конфигурация для ID
type IdConfig struct {
	Prefix string
	Err    error
}

// Интерфейс для получения конфигурации
type IDConfiger interface {
	Config() IdConfig
}

// Обобщенный тип ID
type ID[T any] string

// Методы для обобщенного типа
func (id ID[T]) String() string {
	return string(id)
}

func (id ID[T]) Equals(other ID[T]) bool {
	return id == other
}

// Функция создания ID
func NewID[T IDConfiger](id string) (ID[T], error) {
	var zero T
	cfg := zero.Config()
	if !idgen.ValidatePrefixedID(id, cfg.Prefix, 8) {
		return "", cfg.Err
	}
	return ID[T](strings.ToUpper(id)), nil
}

// Функция генерации ID
func GenerateID[T IDConfiger]() ID[T] {
	var zero T
	cfg := zero.Config()
	return ID[T](idgen.GeneratePrefixedID(cfg.Prefix, 8))
}
