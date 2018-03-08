package tox

/*
#include <tox/toxencryptsave.h>
*/
import "C"
import (
	"github.com/TokTok/go-toxcore-c/toxenums"
)

const PASS_KEY_LENGTH = int(C.TOX_PASS_KEY_LENGTH)
const PASS_ENCRYPTION_EXTRA_LENGTH = int(C.TOX_PASS_ENCRYPTION_EXTRA_LENGTH)

type ToxPassKey struct {
	cpk *C.Tox_Pass_Key
}

func (this *ToxPassKey) Free() {
	C.tox_pass_key_free(this.cpk)
}

func Derive(passphrase []byte) (*ToxPassKey, error) {
	this := &ToxPassKey{}

	passphrase_ := (*C.uint8_t)(&passphrase[0])

	var cerr C.TOX_ERR_KEY_DERIVATION
	this.cpk = C.tox_pass_key_derive(passphrase_, C.size_t(len(passphrase)), &cerr)
	if cerr != 0 {
		return nil, toxenums.TOX_ERR_KEY_DERIVATION(cerr)
	}
	return this, nil
}

func DeriveWithSalt(passphrase []byte, salt []byte) (*ToxPassKey, error) {
	this := &ToxPassKey{}

	passphrase_ := (*C.uint8_t)(&passphrase[0])
	salt_ := (*C.uint8_t)(&salt[0])

	var cerr C.TOX_ERR_KEY_DERIVATION
	this.cpk = C.tox_pass_key_derive_with_salt(passphrase_, C.size_t(len(passphrase)), salt_, &cerr)
	if cerr != 0 {
		return nil, toxenums.TOX_ERR_KEY_DERIVATION(cerr)
	}
	return this, nil
}

func (this *ToxPassKey) Encrypt(plaintext []byte) ([]byte, error) {
	ciphertext := make([]byte, len(plaintext)+PASS_ENCRYPTION_EXTRA_LENGTH)
	ciphertext_ := (*C.uint8_t)(&ciphertext[0])
	plaintext_ := (*C.uint8_t)(&plaintext[0])

	var cerr C.TOX_ERR_ENCRYPTION
	C.tox_pass_key_encrypt(this.cpk, plaintext_, C.size_t(len(plaintext)), ciphertext_, &cerr)
	if cerr != 0 {
		return nil, toxenums.TOX_ERR_ENCRYPTION(cerr)
	}
	return ciphertext, nil
}

func (this *ToxPassKey) Decrypt(ciphertext []byte) ([]byte, error) {
	ciphertext_ := (*C.uint8_t)(&ciphertext[0])
	plaintext := make([]byte, len(ciphertext)-PASS_ENCRYPTION_EXTRA_LENGTH)
	plaintext_ := (*C.uint8_t)(&plaintext[0])

	var cerr C.TOX_ERR_DECRYPTION
	C.tox_pass_key_decrypt(this.cpk, ciphertext_, C.size_t(len(ciphertext)), plaintext_, &cerr)
	if cerr != 0 {
		return nil, toxenums.TOX_ERR_DECRYPTION(cerr)
	}
	return plaintext, nil
}

func GetSalt(ciphertext []byte) ([]byte, error) {
	ciphertext_ := (*C.uint8_t)(&ciphertext[0])
	salt := make([]byte, int(C.TOX_PASS_SALT_LENGTH))
	salt_ := (*C.uint8_t)(&salt[0])

	var cerr C.TOX_ERR_GET_SALT
	C.tox_get_salt(ciphertext_, salt_, &cerr)
	if cerr != 0 {
		return nil, toxenums.TOX_ERR_GET_SALT(cerr)
	}
	return salt, nil
}

func IsDataEncrypted(data []byte) bool {
	data_ := (*C.uint8_t)(&data[0])
	return bool(C.tox_is_data_encrypted(data_))
}

func PassEncrypt(plaintext []byte, passphrase []byte) ([]byte, error) {
	ciphertext := make([]byte, len(plaintext)+PASS_ENCRYPTION_EXTRA_LENGTH)
	ciphertext_ := (*C.uint8_t)(&ciphertext[0])
	plaintext_ := (*C.uint8_t)(&plaintext[0])
	passphrase_ := (*C.uint8_t)(&passphrase[0])

	var cerr C.TOX_ERR_ENCRYPTION
	C.tox_pass_encrypt(plaintext_, C.size_t(len(plaintext)), passphrase_, C.size_t(len(passphrase)), ciphertext_, &cerr)
	if cerr != 0 {
		return nil, toxenums.TOX_ERR_ENCRYPTION(cerr)
	}
	return ciphertext, nil
}

func PassDecrypt(ciphertext []byte, passphrase []byte) ([]byte, error) {
	ciphertext_ := (*C.uint8_t)(&ciphertext[0])
	plaintext := make([]byte, len(ciphertext)-PASS_ENCRYPTION_EXTRA_LENGTH)
	plaintext_ := (*C.uint8_t)(&plaintext[0])
	passphrase_ := (*C.uint8_t)(&plaintext[0])

	var cerr C.TOX_ERR_DECRYPTION
	C.tox_pass_decrypt(ciphertext_, C.size_t(len(ciphertext)), passphrase_, C.size_t(len(passphrase)), plaintext_, &cerr)
	if cerr != 0 {
		return nil, toxenums.TOX_ERR_DECRYPTION(cerr)
	}
	return plaintext, nil
}
