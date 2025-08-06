package idgen_test

import (
	"strings"
	"testing"

	"github.com/GAKiknadze/payment_service/internal/idgen"
)

func TestGenerateShortID_Length(t *testing.T) {
	testCases := []struct {
		name   string
		length int
	}{
		{"Zero length", 0},
		{"Single character", 1},
		{"Typical length", 8},
		{"Long ID", 32},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id := idgen.GenerateShortID(tc.length)
			if len(id) != tc.length {
				t.Errorf("expected length %d, got %d", tc.length, len(id))
			}
		})
	}
}

func TestGenerateShortID_ValidCharacters(t *testing.T) {
	id := idgen.GenerateShortID(1000) // Генерируем длинный ID для надежности проверки

	for i := 0; i < len(id); i++ {
		char := id[i]
		if idx := idgen.Alphabet[strings.IndexByte(idgen.Alphabet, char)]; idx != char {
			t.Errorf("invalid character '%c' found in ID", char)
		}
	}
}

func TestGenerateShortID_Uniqueness(t *testing.T) {
	const (
		count = 1000
		size  = 16
	)
	ids := make(map[string]struct{}, count)

	for i := 0; i < count; i++ {
		id := idgen.GenerateShortID(size)
		if _, exists := ids[id]; exists {
			t.Fatalf("duplicate ID at iteration %d: %s", i, id)
		}
		ids[id] = struct{}{}
	}

	if len(ids) != count {
		t.Errorf("unexpected collisions: %d unique IDs out of %d", len(ids), count)
	}
}

func TestGenerateShortID_PanicOnNegative(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for negative length")
		} else if msg, ok := r.(string); !ok || !strings.Contains(msg, "negative length") {
			t.Errorf("unexpected panic message: %v", r)
		}
	}()
	idgen.GenerateShortID(-1)
}

func TestGenerateShortID_ZeroLength(t *testing.T) {
	if id := idgen.GenerateShortID(0); id != "" {
		t.Errorf("expected empty string for zero length, got '%s'", id)
	}
}
