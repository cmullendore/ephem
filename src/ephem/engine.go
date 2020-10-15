package ephem

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"

	"github.com/cmullendore/ephem/src/config"
)

// Engine is the primary engine object, containing the appropriate
// configuration, cache, and log objects necessary to execute any
// given intbound request.
type Engine struct {
	Configuration *config.Ephem
	Cache         *Cache
}

// CreateEngine creates a new request handler and initializes the
// internal objects, returning a complete and useable object.
func CreateEngine() *Engine {

	var c = config.LoadEphemConfig()

	return &Engine{
		Configuration: c,
		Cache:         NewCache(c),
	}
}

// SaveItem persists and item into the local cache. Note that it is
// up to the cache layer to persist the item all the way to storage.
func (h *Engine) SaveItem(path *string, item *[]byte) *error {
	var key []byte
	var err error
	if h.Configuration.Aes256KeySource == "URL" {

		rawBytes, err5 := base64.RawURLEncoding.DecodeString(*path)
		if err5 != nil {
			log.Println(err)
			return &err
		}
		if len(rawBytes) < 32 {
			err = errors.New("Insufficient Key Length")
		}
		key = rawBytes[:32]
	} else if h.Configuration.Aes256Key != "" {
		key, err = base64.RawURLEncoding.DecodeString(h.Configuration.Aes256Key)
	}

	if err != nil {
		log.Println(err)
		return &err
	}

	enc, err2 := encryptToString(key, item)
	if err2 != nil {
		return err2
	}

	var pathBytes = []byte(*path)
	if err := h.Cache.SaveItem(getHashBase64(&pathBytes), enc, h.Configuration.MaxAgeSeconds, 0); err != nil {
		return err
	}

	return nil
}

// GetItem retrieves an item from the local cache. If the item does not
// exist in the cache it will be retrieved from the database. If it does
// not exist in the database, this will return an error.
func (h *Engine) GetItem(path *string) (*[]byte, *error) {
	var key []byte
	var err error

	if h.Configuration.Aes256KeySource == "URL" {

		rawBytes, err5 := base64.RawURLEncoding.DecodeString(*path)
		if err5 != nil {
			log.Println(err)
		}
		if len(rawBytes) < 32 {
			err = errors.New("Insufficient Key Length")
		}
		key = rawBytes[:32]
	} else if h.Configuration.Aes256Key != "" {
		key, err = base64.RawURLEncoding.DecodeString(h.Configuration.Aes256Key)
	}

	if err != nil {
		log.Println(err)
	}

	var pathBytes = []byte(*path)

	item, err1 := h.Cache.GetItem(getHashBase64(&pathBytes))
	if err1 != nil {
		return nil, err1
	}
	if item == nil {
		return nil, nil
	}

	dec, err3 := decryptToBytes(key, item)
	if err3 != nil {
		return nil, err3
	}

	return dec, nil
}

// DeleteItem purges an item from the local cache, which will then purge
// it from the database immediately, not waiting for the configured
// automatic purge settings.
func (h *Engine) DeleteItem(path *string) *error {
	var pathBytes = []byte(*path)

	if err := h.Cache.DeleteItem(getHashBase64(&pathBytes)); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func encryptToString(key []byte, secret *[]byte) (*string, *error) {
	c, err := newAesCipher(key)
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, &err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, &err
	}

	str := base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, *secret, nil))
	return &str, nil
}

func decryptToBytes(key []byte, secret *string) (*[]byte, *error) {
	c, err1 := newAesCipher(key)
	if err1 != nil {
		return nil, &err1
	}

	encBytes, err2 := base64.StdEncoding.DecodeString(*secret)
	if err2 != nil {
		return nil, &err2
	}

	gcmDecrypt, err3 := cipher.NewGCM(c)
	if err3 != nil {
		return nil, &err3
	}

	nonceSize := gcmDecrypt.NonceSize()

	nonce, encryptedMessage := encBytes[:nonceSize], encBytes[nonceSize:]

	decBytes, err4 := gcmDecrypt.Open(nil, []byte(nonce), []byte(encryptedMessage), nil)
	if err4 != nil {
		return nil, &err4
	}

	return &decBytes, nil
}

func newAesCipher(key []byte) (cipher.Block, error) {
	if len(key) < 32 {
		return nil, errors.New("Insufficient key length")
	}

	cipher, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}

	return cipher, nil
}
