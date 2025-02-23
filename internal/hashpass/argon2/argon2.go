package argon2

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"github.com/RinatZaynet/CouchFilmCritic/internal/hashpass"
	"github.com/RinatZaynet/CouchFilmCritic/internal/helpers/random"
	a2 "golang.org/x/crypto/argon2"
)

const (
	defaultSaltLen int    = 16
	defaultKeyLen  uint32 = 32
)

type Options struct {
	Time    uint32
	Memory  uint32
	Threads uint8
}

type Manager struct {
	Opt *Options
}

func NewManager(opt *Options) *Manager {
	return &Manager{opt}
}

func (m *Manager) HashingPassword(password []byte) (formatHash string) {
	salt := random.RandomSliceByte(defaultSaltLen)

	hash := a2.IDKey(password, salt, m.Opt.Time, m.Opt.Memory, m.Opt.Threads, defaultKeyLen)

	formatHast := fmt.Sprintf("$argon2id$v=%d$t=%d,m=%d,p=%d$%s$%s",
		a2.Version, m.Opt.Time, m.Opt.Memory, m.Opt.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return formatHast
}

func (m *Manager) CompareHashAndPassword(password []byte, formatHash string) error {
	const fn = "argon2.Manager.CompareHashAndPassword"

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
