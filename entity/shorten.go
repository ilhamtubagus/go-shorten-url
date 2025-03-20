package entity

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ilhamtubagus/go-shorten-url/util"
	"math/big"
	"os"
)

var (
	encodeBase62 = util.EncodeBase62
)

type ShortenedURL struct {
	ShortCode    string `json:"shortCode" bson:"shortCode"`
	OriginalURL  string `json:"originalURL" bson:"originalURL"`
	ShortenedURL string `json:"shortenedURL,omitempty" bson:",omitempty"`
}

func (s *ShortenedURL) GenerateShortCode(salt ...string) string {
	var plain string
	if len(salt) > 0 {
		plain = fmt.Sprintf("%s%s", salt[0], s.OriginalURL)
	} else {
		plain = s.OriginalURL
	}

	hash := md5.Sum([]byte(plain))
	hashHex := hex.EncodeToString(hash[:])

	// Get first 6 bytes (12 hex characters)
	first6BytesHex := hashHex[:12]

	decimalValue := new(big.Int)
	decimalValue.SetString(first6BytesHex, 16)

	shortCode := encodeBase62(decimalValue)
	s.ShortCode = shortCode

	return shortCode
}

func (s *ShortenedURL) GenerateShortenedURL() error {
	if s.ShortCode == "" {
		return fmt.Errorf("short code not specified")
	}

	host := fmt.Sprintf("%s:%s", os.Getenv("SERVICE_HOST"), os.Getenv("SERVICE_PORT"))
	s.ShortenedURL = fmt.Sprintf("%s/%s", host, s.ShortCode)

	return nil
}
