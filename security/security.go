package security

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/nafisalfiani/ketson-go-lib/log"
)

type Interface interface {
	HashPassword(ctx context.Context, password string) (string, error)
	CompareHashPassword(ctx context.Context, hashPassword, password string) bool
}

type Config struct {
	SecretKey string `env:"SECRET_KEY"`
}

type security struct {
	log log.Interface
	cfg Config
}

func Init(log log.Interface, cfg Config) Interface {
	return &security{
		log: log,
		cfg: cfg,
	}
}

func (s *security) HashPassword(ctx context.Context, password string) (string, error) {
	computedHash := hmac.New(sha256.New, []byte(s.cfg.SecretKey))
	n, err := computedHash.Write([]byte(password))
	if err != nil {
		return "", err
	}
	s.log.Debug(ctx, fmt.Sprintf("%v bytes written", n))

	return hex.EncodeToString(computedHash.Sum(nil)), nil
}

func (s *security) CompareHashPassword(ctx context.Context, hashPassword, password string) bool {
	hashedPassword, err := s.HashPassword(ctx, password)
	if err != nil {
		return false
	}

	return hashPassword == hashedPassword
}
