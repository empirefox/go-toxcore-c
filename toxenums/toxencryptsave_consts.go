//go:generate stringer -type=TOX_ERR_KEY_DERIVATION,TOX_ERR_ENCRYPTION,TOX_ERR_DECRYPTION,TOX_ERR_GET_SALT
package toxenums

import "fmt"

type TOX_ERR_KEY_DERIVATION int

const (
	TOX_ERR_KEY_DERIVATION_OK TOX_ERR_KEY_DERIVATION = iota
	TOX_ERR_KEY_DERIVATION_NULL
	TOX_ERR_KEY_DERIVATION_FAILED
)

type TOX_ERR_ENCRYPTION int

const (
	TOX_ERR_ENCRYPTION_OK TOX_ERR_ENCRYPTION = iota
	TOX_ERR_ENCRYPTION_NULL
	TOX_ERR_ENCRYPTION_KEY_DERIVATION_FAILED
	TOX_ERR_ENCRYPTION_FAILED
)

type TOX_ERR_DECRYPTION int

const (
	TOX_ERR_DECRYPTION_OK TOX_ERR_DECRYPTION = iota
	TOX_ERR_DECRYPTION_NULL
	TOX_ERR_DECRYPTION_INVALID_LENGTH
	TOX_ERR_DECRYPTION_BAD_FORMAT
	TOX_ERR_DECRYPTION_KEY_DERIVATION_FAILED
	TOX_ERR_DECRYPTION_FAILED
)

type TOX_ERR_GET_SALT int

const (
	TOX_ERR_GET_SALT_OK TOX_ERR_GET_SALT = iota
	TOX_ERR_GET_SALT_NULL
	TOX_ERR_GET_SALT_BAD_FORMAT
)

func (e TOX_ERR_KEY_DERIVATION) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_ENCRYPTION) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_DECRYPTION) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_GET_SALT) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
