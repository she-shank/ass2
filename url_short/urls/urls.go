package urls

import (
	"crypto/sha256"
	"regexp"
)

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Allows only "http" and "https" protocols.
const regex = `^(https?):\/\/(?:www\.)?((?:[a-zA-Z0-9\-]+\.)+[a-zA-Z]{2,})(?::(\d{1,5}))?(\/\S*)?$`

var len = 6

// Collisions can happen, where two different URLs generate the same hash.
func Shorten(url string) (string, error) {
	hash := sha256.Sum256([]byte(url))

	// 6 bytes (48 bits) produce a decimal number that typically converts
	// to a Base62 string of approximately 7 characters.
	shortHash := hash[:len]
	decimalValue := 0
	for _, b := range shortHash {
		decimalValue = (decimalValue << 8) + int(b)
	}

	// Using Base62 instead of Base64 to avoid the use of characters that
	// can be confused with each other (e.g., 0 and O) and/or are bad for URLs (_).
	var base62Value string
	for decimalValue > 0 {
		remainder := decimalValue % 62
		base62Value = string(chars[remainder]) + base62Value
		decimalValue /= 62
	}
	return base62Value, nil
}

func ValidateURL(incURL string) bool {
	match, _ := regexp.MatchString(regex, incURL)
	return match
}

// func StripDomain(incURL string) string {
// 	parsedURL, err := url.Parse(incURL)
// 	if err != nil {
// 		return ""
// 	}
// 	return parsedURL.Hostname()
// }
