package tox

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"unsafe"
)

func safeptr(b []byte) unsafe.Pointer {
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return unsafe.Pointer(h.Data)
}

func DecodeAddress(addr string) (*[ADDRESS_SIZE]byte, error) {
	var addrb [ADDRESS_SIZE]byte
	n, err := hex.Decode(addrb[:], bytes.ToLower([]byte(addr)))
	if err != nil {
		return nil, err
	}
	if n != ADDRESS_SIZE {
		return nil, fmt.Errorf("Tox address bytes len should be %d, but got %d", ADDRESS_SIZE, n)
	}
	return &addrb, nil
}

func MustDecodeAddress(address string) *[ADDRESS_SIZE]byte {
	addressb, err := DecodeAddress(address)
	if err != nil {
		panic(err)
	}
	return addressb
}

func ToPubkey(address *[ADDRESS_SIZE]byte) *[PUBLIC_KEY_SIZE]byte {
	var pubkey [PUBLIC_KEY_SIZE]byte
	copy(pubkey[:], address[:])
	return &pubkey
}

func DecodePubkey(pubkey string) (*[PUBLIC_KEY_SIZE]byte, error) {
	var pubkeyb [PUBLIC_KEY_SIZE]byte
	n, err := hex.Decode(pubkeyb[:], bytes.ToLower([]byte(pubkey)))
	if err != nil {
		return nil, err
	}
	if n != PUBLIC_KEY_SIZE {
		return nil, fmt.Errorf("Tox public key bytes len should be %d, but got %d", PUBLIC_KEY_SIZE, n)
	}
	return &pubkeyb, nil
}

func MustDecodePubkey(pubkey string) *[PUBLIC_KEY_SIZE]byte {
	return MustDecodeSecret(pubkey)
}

func MustDecodeSecret(secret string) *[SECRET_KEY_SIZE]byte {
	b, err := DecodePubkey(secret)
	if err != nil {
		panic(err)
	}
	return b
}

func WriteSavedata(fname string, savedata []byte) error {
	if _, err := os.Stat(fname); err != nil {
		err = ioutil.WriteFile(fname, savedata, 0755)
		if err != nil {
			return err
		}
	} else {
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		if bytes.Compare(data, savedata) != 0 {
			err := ioutil.WriteFile(fname, savedata, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
