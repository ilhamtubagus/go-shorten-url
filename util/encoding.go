package util

import (
	"math/big"
	"strings"
)

// Base62 characters
const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// EncodeBase62 Encode a number to base62
func EncodeBase62(num *big.Int) string {
	if num.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}

	base := big.NewInt(62)
	var encoded strings.Builder
	mod := new(big.Int)

	for num.Cmp(big.NewInt(0)) > 0 {
		num.DivMod(num, base, mod)
		encoded.WriteByte(base62Chars[mod.Int64()])
	}

	// Reverse the string since we construct it backwards
	runes := []rune(encoded.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
