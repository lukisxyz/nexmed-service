package auth

import (
	"bytes"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/argon2"
)

type HashSalt struct {
	Hash []byte
}

type Argon2idHash struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

var argon2IdHash *Argon2idHash

func SetHashParam(time, saltLen uint32, memory uint32, threads uint8, keyLen uint32) {
	argon2IdHash = &Argon2idHash{
		time:    time,
		saltLen: saltLen,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
	}
}

func generateSalt(length uint32) ([]byte, error) {
    salt := make([]byte, length)
    _, err := rand.Read(salt)
    if err != nil {
        return nil, err
    }
    return salt, nil
}

func randomSecret(length uint32) ([]byte, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func generateHash(password, salt []byte) (*HashSalt, error) {
	var err error
	if len(salt) == 0 {
		salt, err = randomSecret(argon2IdHash.saltLen)
	}
	if err != nil {
		return nil, err
	}
	hash := argon2.IDKey(
		password,
		salt,
		argon2IdHash.time,
		argon2IdHash.memory,
		argon2IdHash.threads,
		argon2IdHash.keyLen,
	)
	return &HashSalt{Hash: hash}, nil
}

func compare(hash, salt, password []byte) error {
	hashSalt, err := generateHash(password, salt)
	if err != nil {
		return err
	}
	if !bytes.Equal(hash, hashSalt.Hash) {
		return errors.New("hash doesn't match")
	}
	return nil
}