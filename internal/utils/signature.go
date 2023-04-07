package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

func GenerateSignature(apiSecret string) (string, int64) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	message := strconv.FormatInt(timestamp, 10) + "stream"
	hash := hmac.New(sha256.New, []byte(apiSecret))

	hash.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), timestamp
}
