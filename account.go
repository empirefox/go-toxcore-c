package tox

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"strings"

	"golang.org/x/crypto/curve25519"
)

var ErrKeyLen = errors.New("key len error")

type Account struct {
	Address [ADDRESS_SIZE]byte
	Secret  [SECRET_KEY_SIZE]byte
	Pubkey  [PUBLIC_KEY_SIZE]byte
	Nospam  uint32
}

func GenerateAccount(randReader io.Reader, nospam uint32) (*Account, error) {
	var scalar [32]byte
	if _, err := io.ReadFull(randReader, scalar[:]); err != nil {
		return nil, err
	}
	return NewAccount(&scalar, nospam), nil
}

func NewAccountFrom(secret string, nospam uint32) (*Account, error) {
	var scalar [32]byte

	sk, err := hex.DecodeString(strings.ToLower(secret))
	if err != nil {
		return nil, err
	}

	if copy(scalar[:], sk) != 32 {
		return nil, ErrKeyLen
	}

	return NewAccount(&scalar, nospam), nil
}

func NewAccount(secret *[SECRET_KEY_SIZE]byte, nospam uint32) *Account {
	var toxid [ADDRESS_SIZE]byte

	var pubkey [PUBLIC_KEY_SIZE]byte
	curve25519.ScalarBaseMult(&pubkey, secret)
	copy(toxid[:PUBLIC_KEY_SIZE], pubkey[:])

	binary.BigEndian.PutUint32(toxid[PUBLIC_KEY_SIZE:], nospam)

	checksum := toxid[36:]
	for i := 0; i < 36; i++ {
		checksum[i%2] ^= toxid[i]
	}

	return &Account{
		Address: toxid,
		Secret:  *secret,
		Pubkey:  pubkey,
		Nospam:  nospam,
	}
}

type HumanAccount struct {
	Address string
	Secret  string
	Pubkey  string
	Nospam  uint32
}

func (a *Account) HumanReadable() *HumanAccount {
	return &HumanAccount{
		Address: strings.ToUpper(hex.EncodeToString(a.Address[:])),
		Secret:  strings.ToUpper(hex.EncodeToString(a.Secret[:])),
		Pubkey:  strings.ToUpper(hex.EncodeToString(a.Pubkey[:])),
		Nospam:  a.Nospam,
	}
}
