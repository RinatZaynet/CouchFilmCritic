package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	hashpass "github.com/RinatZaynet/CouchFilmCritic/internal/hashingPassword"
	a2 "golang.org/x/crypto/argon2"
)

const (
	defaultSaltLen uint32 = 16
	defaultKeyLen  uint32 = 32
)

type Options struct {
	Time    uint32
	Memory  uint32
	Threads uint8
}

type ManagerArgon2 struct {
	Opt *Options
}

func NewManagerArgon2(opt *Options) *ManagerArgon2 {
	return &ManagerArgon2{opt}
}

func (man *ManagerArgon2) HashingPassword(password []byte) (formatHash string, err error) {
	const fn = "argon2.ManagerArgon2.HasingPassword"

	salt := make([]byte, defaultSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	hash := a2.IDKey(password, salt, man.Opt.Time, man.Opt.Memory, man.Opt.Threads, defaultKeyLen)

	formatHast := fmt.Sprintf("$argon2id$v=%d$t=%d,m=%d,p=%d$%s$%s",
		a2.Version, man.Opt.Time, man.Opt.Memory, man.Opt.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return formatHast, nil
}

func (mng *ManagerArgon2) CompareHashAndPassword(password []byte, formatHash string) error {
	const fn = "argon2.ManagerArgon2.CompareHashAndPassword"

	p, err := parse(formatHash)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	compareHash := a2.IDKey(password, p.salt, p.time, p.memory, p.threads, p.keyLen)

	if subtle.ConstantTimeCompare(p.hash, compareHash) != 1 {
		return hashpass.ErrMismatchesTypes
	}

	return nil
}
