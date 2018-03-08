package tox

import (
	"bytes"
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
