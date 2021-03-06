// Code generated by "stringer -type=TOX_ERR_KEY_DERIVATION,TOX_ERR_ENCRYPTION,TOX_ERR_DECRYPTION,TOX_ERR_GET_SALT"; DO NOT EDIT.

package toxenums

import "fmt"

const _TOX_ERR_KEY_DERIVATION_name = "TOX_ERR_KEY_DERIVATION_OKTOX_ERR_KEY_DERIVATION_NULLTOX_ERR_KEY_DERIVATION_FAILED"

var _TOX_ERR_KEY_DERIVATION_index = [...]uint8{0, 25, 52, 81}

func (i TOX_ERR_KEY_DERIVATION) String() string {
	if i < 0 || i >= TOX_ERR_KEY_DERIVATION(len(_TOX_ERR_KEY_DERIVATION_index)-1) {
		return fmt.Sprintf("TOX_ERR_KEY_DERIVATION(%d)", i)
	}
	return _TOX_ERR_KEY_DERIVATION_name[_TOX_ERR_KEY_DERIVATION_index[i]:_TOX_ERR_KEY_DERIVATION_index[i+1]]
}

const _TOX_ERR_ENCRYPTION_name = "TOX_ERR_ENCRYPTION_OKTOX_ERR_ENCRYPTION_NULLTOX_ERR_ENCRYPTION_KEY_DERIVATION_FAILEDTOX_ERR_ENCRYPTION_FAILED"

var _TOX_ERR_ENCRYPTION_index = [...]uint8{0, 21, 44, 84, 109}

func (i TOX_ERR_ENCRYPTION) String() string {
	if i < 0 || i >= TOX_ERR_ENCRYPTION(len(_TOX_ERR_ENCRYPTION_index)-1) {
		return fmt.Sprintf("TOX_ERR_ENCRYPTION(%d)", i)
	}
	return _TOX_ERR_ENCRYPTION_name[_TOX_ERR_ENCRYPTION_index[i]:_TOX_ERR_ENCRYPTION_index[i+1]]
}

const _TOX_ERR_DECRYPTION_name = "TOX_ERR_DECRYPTION_OKTOX_ERR_DECRYPTION_NULLTOX_ERR_DECRYPTION_INVALID_LENGTHTOX_ERR_DECRYPTION_BAD_FORMATTOX_ERR_DECRYPTION_KEY_DERIVATION_FAILEDTOX_ERR_DECRYPTION_FAILED"

var _TOX_ERR_DECRYPTION_index = [...]uint8{0, 21, 44, 77, 106, 146, 171}

func (i TOX_ERR_DECRYPTION) String() string {
	if i < 0 || i >= TOX_ERR_DECRYPTION(len(_TOX_ERR_DECRYPTION_index)-1) {
		return fmt.Sprintf("TOX_ERR_DECRYPTION(%d)", i)
	}
	return _TOX_ERR_DECRYPTION_name[_TOX_ERR_DECRYPTION_index[i]:_TOX_ERR_DECRYPTION_index[i+1]]
}

const _TOX_ERR_GET_SALT_name = "TOX_ERR_GET_SALT_OKTOX_ERR_GET_SALT_NULLTOX_ERR_GET_SALT_BAD_FORMAT"

var _TOX_ERR_GET_SALT_index = [...]uint8{0, 19, 40, 67}

func (i TOX_ERR_GET_SALT) String() string {
	if i < 0 || i >= TOX_ERR_GET_SALT(len(_TOX_ERR_GET_SALT_index)-1) {
		return fmt.Sprintf("TOX_ERR_GET_SALT(%d)", i)
	}
	return _TOX_ERR_GET_SALT_name[_TOX_ERR_GET_SALT_index[i]:_TOX_ERR_GET_SALT_index[i+1]]
}
