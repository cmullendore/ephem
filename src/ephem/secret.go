package ephem

import (
	"crypto/sha512"
	"encoding/base64"
)

type Secret struct {
	Path *string
	Data *[]byte
}

func getHashBase64(s *[]byte) *string {
	h := sha512.New()
	h.Write([]byte(*s))
	hashed := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return &hashed
}
