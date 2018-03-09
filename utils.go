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

var toxdebug = false

func SetDebug(debug bool) {
	toxdebug = debug
}

var loglevel = 0

func SetLogLevel(level int) {
	loglevel = level
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

func FileExist(fname string) bool {
	_, err := os.Stat(fname)
	if err != nil {
		return false
	}
	return true
}

func (this *Tox) WriteSavedata(fname string) error {
	if !FileExist(fname) {
		err := ioutil.WriteFile(fname, this.GetSavedata(), 0755)
		if err != nil {
			return err
		}
	} else {
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		liveData := this.GetSavedata()
		if bytes.Compare(data, liveData) != 0 {
			err := ioutil.WriteFile(fname, this.GetSavedata(), 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (this *Tox) LoadSavedata(fname string) ([]byte, error) {
	return ioutil.ReadFile(fname)
}

func LoadSavedata(fname string) ([]byte, error) {
	return ioutil.ReadFile(fname)
}
