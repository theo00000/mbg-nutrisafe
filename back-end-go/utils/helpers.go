package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

func GenerateSlug(name string) string {
	name = strings.ToLower(name)
	var result strings.Builder
	prevDot := false
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
			prevDot = false
		} else if !prevDot && result.Len() > 0 {
			result.WriteByte('.')
			prevDot = true
		}
	}
	s := strings.TrimRight(result.String(), ".")
	if s == "" {
		s = "akun"
	}
	return s
}

func GeneratePassword(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

// GeneratePasswordFromName membuat password dari nama entitas + 4 angka random
// Contoh: "SDN Cipta Karya" → "SDNCiptaKarya4821"
func GeneratePasswordFromName(name string) string {
	words := strings.Fields(name)
	var base strings.Builder
	for _, w := range words {
		if len(w) > 0 {
			base.WriteString(strings.ToUpper(string([]rune(w)[0])))
			if len([]rune(w)) > 1 {
				base.WriteString(string([]rune(w)[1:]))
			}
		}
	}

	// hilangkan karakter non-alphanumeric
	var clean strings.Builder
	for _, r := range base.String() {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			clean.WriteRune(r)
		}
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := fmt.Sprintf("%04d", rng.Intn(10000))
	return clean.String() + digits
}
