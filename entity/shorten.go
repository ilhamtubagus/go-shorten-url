package entity

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/ilhamtubagus/go-shorten-url/util"
	"math/big"
)

var (
	encodeBase62 = util.EncodeBase62
)

type ShortenedURL struct {
	ShortenedURL string `json:"shortenedURL"`
	OriginalURL  string `json:"originalURL"`
}

func (s *ShortenedURL) GenerateShortCode() string {
	hash := md5.Sum([]byte(s.OriginalURL))
	hashHex := hex.EncodeToString(hash[:])

	// Get first 6 bytes (12 hex characters)
	first6BytesHex := hashHex[:12]

	decimalValue := new(big.Int)
	decimalValue.SetString(first6BytesHex, 16)

	shortCode := encodeBase62(decimalValue)
	s.ShortenedURL = shortCode

	return shortCode
}
