//go:generate stringer -type=TOX_USER_STATUS,TOX_MESSAGE_TYPE,TOX_PROXY_TYPE,TOX_SAVEDATA_TYPE,TOX_LOG_LEVEL,TOX_ERR_OPTIONS_NEW,TOX_ERR_NEW,TOX_ERR_BOOTSTRAP,TOX_CONNECTION,TOX_ERR_SET_INFO,TOX_ERR_FRIEND_ADD,TOX_ERR_FRIEND_DELETE,TOX_ERR_FRIEND_BY_PUBLIC_KEY,TOX_ERR_FRIEND_GET_PUBLIC_KEY,TOX_ERR_FRIEND_GET_LAST_ONLINE,TOX_ERR_FRIEND_QUERY,TOX_ERR_SET_TYPING,TOX_ERR_FRIEND_SEND_MESSAGE,TOX_FILE_KIND,TOX_FILE_CONTROL,TOX_ERR_FILE_CONTROL,TOX_ERR_FILE_SEEK,TOX_ERR_FILE_GET,TOX_ERR_FILE_SEND,TOX_ERR_FILE_SEND_CHUNK,TOX_CONFERENCE_TYPE,TOX_ERR_CONFERENCE_NEW,TOX_ERR_CONFERENCE_DELETE,TOX_ERR_CONFERENCE_PEER_QUERY,TOX_ERR_CONFERENCE_INVITE,TOX_ERR_CONFERENCE_JOIN,TOX_ERR_CONFERENCE_SEND_MESSAGE,TOX_ERR_CONFERENCE_TITLE,TOX_ERR_CONFERENCE_GET_TYPE,TOX_ERR_FRIEND_CUSTOM_PACKET,TOX_ERR_GET_PORT
package toxenums

import "fmt"

type TOX_USER_STATUS int

const (
	TOX_USER_STATUS_NONE TOX_USER_STATUS = iota
	TOX_USER_STATUS_AWAY
	TOX_USER_STATUS_BUSY
)

type TOX_MESSAGE_TYPE int

const (
	TOX_MESSAGE_TYPE_NORMAL TOX_MESSAGE_TYPE = iota
	TOX_MESSAGE_TYPE_ACTION
)

type TOX_PROXY_TYPE int

const (
	TOX_PROXY_TYPE_NONE TOX_PROXY_TYPE = iota
	TOX_PROXY_TYPE_HTTP
	TOX_PROXY_TYPE_SOCKS5
)

type TOX_SAVEDATA_TYPE int

const (
	TOX_SAVEDATA_TYPE_NONE TOX_SAVEDATA_TYPE = iota
	TOX_SAVEDATA_TYPE_TOX_SAVE
	TOX_SAVEDATA_TYPE_SECRET_KEY
)

type TOX_LOG_LEVEL int

const (
	TOX_LOG_LEVEL_TRACE TOX_LOG_LEVEL = iota
	TOX_LOG_LEVEL_DEBUG
	TOX_LOG_LEVEL_INFO
	TOX_LOG_LEVEL_WARNING
	TOX_LOG_LEVEL_ERROR
)

type TOX_ERR_OPTIONS_NEW int

const (
	TOX_ERR_OPTIONS_NEW_OK TOX_ERR_OPTIONS_NEW = iota
	TOX_ERR_OPTIONS_NEW_MALLOC
)

type TOX_ERR_NEW int

const (
	TOX_ERR_NEW_OK TOX_ERR_NEW = iota
	TOX_ERR_NEW_NULL
	TOX_ERR_NEW_MALLOC
	TOX_ERR_NEW_PORT_ALLOC
	TOX_ERR_NEW_PROXY_BAD_TYPE
	TOX_ERR_NEW_PROXY_BAD_HOST
	TOX_ERR_NEW_PROXY_BAD_PORT
	TOX_ERR_NEW_PROXY_NOT_FOUND
	TOX_ERR_NEW_LOAD_ENCRYPTED
	TOX_ERR_NEW_LOAD_BAD_FORMAT
)

type TOX_ERR_BOOTSTRAP int

const (
	TOX_ERR_BOOTSTRAP_OK TOX_ERR_BOOTSTRAP = iota
	TOX_ERR_BOOTSTRAP_NULL
	TOX_ERR_BOOTSTRAP_BAD_HOST
	TOX_ERR_BOOTSTRAP_BAD_PORT
)

type TOX_CONNECTION int

const (
	TOX_CONNECTION_NONE TOX_CONNECTION = iota
	TOX_CONNECTION_TCP
	TOX_CONNECTION_UDP
)

type TOX_ERR_SET_INFO int

const (
	TOX_ERR_SET_INFO_OK TOX_ERR_SET_INFO = iota
	TOX_ERR_SET_INFO_NULL
	TOX_ERR_SET_INFO_TOO_LONG
)

type TOX_ERR_FRIEND_ADD int

const (
	TOX_ERR_FRIEND_ADD_OK TOX_ERR_FRIEND_ADD = iota
	TOX_ERR_FRIEND_ADD_NULL
	TOX_ERR_FRIEND_ADD_TOO_LONG
	TOX_ERR_FRIEND_ADD_NO_MESSAGE
	TOX_ERR_FRIEND_ADD_OWN_KEY
	TOX_ERR_FRIEND_ADD_ALREADY_SENT
	TOX_ERR_FRIEND_ADD_BAD_CHECKSUM
	TOX_ERR_FRIEND_ADD_SET_NEW_NOSPAM
	TOX_ERR_FRIEND_ADD_MALLOC
)

type TOX_ERR_FRIEND_DELETE int

const (
	TOX_ERR_FRIEND_DELETE_OK TOX_ERR_FRIEND_DELETE = iota
	TOX_ERR_FRIEND_DELETE_FRIEND_NOT_FOUND
)

type TOX_ERR_FRIEND_BY_PUBLIC_KEY int

const (
	TOX_ERR_FRIEND_BY_PUBLIC_KEY_OK TOX_ERR_FRIEND_BY_PUBLIC_KEY = iota
	TOX_ERR_FRIEND_BY_PUBLIC_KEY_NULL
	TOX_ERR_FRIEND_BY_PUBLIC_KEY_NOT_FOUND
)

type TOX_ERR_FRIEND_GET_PUBLIC_KEY int

const (
	TOX_ERR_FRIEND_GET_PUBLIC_KEY_OK TOX_ERR_FRIEND_GET_PUBLIC_KEY = iota
	TOX_ERR_FRIEND_GET_PUBLIC_KEY_FRIEND_NOT_FOUND
)

type TOX_ERR_FRIEND_GET_LAST_ONLINE int

const (
	TOX_ERR_FRIEND_GET_LAST_ONLINE_OK TOX_ERR_FRIEND_GET_LAST_ONLINE = iota
	TOX_ERR_FRIEND_GET_LAST_ONLINE_FRIEND_NOT_FOUND
)

type TOX_ERR_FRIEND_QUERY int

const (
	TOX_ERR_FRIEND_QUERY_OK TOX_ERR_FRIEND_QUERY = iota
	TOX_ERR_FRIEND_QUERY_NULL
	TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND
)

type TOX_ERR_SET_TYPING int

const (
	TOX_ERR_SET_TYPING_OK TOX_ERR_SET_TYPING = iota
	TOX_ERR_SET_TYPING_FRIEND_NOT_FOUND
)

type TOX_ERR_FRIEND_SEND_MESSAGE int

const (
	TOX_ERR_FRIEND_SEND_MESSAGE_OK TOX_ERR_FRIEND_SEND_MESSAGE = iota
	TOX_ERR_FRIEND_SEND_MESSAGE_NULL
	TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_FOUND
	TOX_ERR_FRIEND_SEND_MESSAGE_FRIEND_NOT_CONNECTED
	TOX_ERR_FRIEND_SEND_MESSAGE_SENDQ
	TOX_ERR_FRIEND_SEND_MESSAGE_TOO_LONG
	TOX_ERR_FRIEND_SEND_MESSAGE_EMPTY
)

type TOX_FILE_KIND int

const (
	TOX_FILE_KIND_DATA TOX_FILE_KIND = iota
	TOX_FILE_KIND_AVATAR
)

type TOX_FILE_CONTROL int

const (
	TOX_FILE_CONTROL_RESUME TOX_FILE_CONTROL = iota
	TOX_FILE_CONTROL_PAUSE
	TOX_FILE_CONTROL_CANCEL
)

type TOX_ERR_FILE_CONTROL int

const (
	TOX_ERR_FILE_CONTROL_OK TOX_ERR_FILE_CONTROL = iota
	TOX_ERR_FILE_CONTROL_FRIEND_NOT_FOUND
	TOX_ERR_FILE_CONTROL_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_CONTROL_NOT_FOUND
	TOX_ERR_FILE_CONTROL_NOT_PAUSED
	TOX_ERR_FILE_CONTROL_DENIED
	TOX_ERR_FILE_CONTROL_ALREADY_PAUSED
	TOX_ERR_FILE_CONTROL_SENDQ
)

type TOX_ERR_FILE_SEEK int

const (
	TOX_ERR_FILE_SEEK_OK TOX_ERR_FILE_SEEK = iota
	TOX_ERR_FILE_SEEK_FRIEND_NOT_FOUND
	TOX_ERR_FILE_SEEK_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_SEEK_NOT_FOUND
	TOX_ERR_FILE_SEEK_DENIED
	TOX_ERR_FILE_SEEK_INVALID_POSITION
	TOX_ERR_FILE_SEEK_SENDQ
)

type TOX_ERR_FILE_GET int

const (
	TOX_ERR_FILE_GET_OK TOX_ERR_FILE_GET = iota
	TOX_ERR_FILE_GET_NULL
	TOX_ERR_FILE_GET_FRIEND_NOT_FOUND
	TOX_ERR_FILE_GET_NOT_FOUND
)

type TOX_ERR_FILE_SEND int

const (
	TOX_ERR_FILE_SEND_OK TOX_ERR_FILE_SEND = iota
	TOX_ERR_FILE_SEND_NULL
	TOX_ERR_FILE_SEND_FRIEND_NOT_FOUND
	TOX_ERR_FILE_SEND_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_SEND_NAME_TOO_LONG
	TOX_ERR_FILE_SEND_TOO_MANY
)

type TOX_ERR_FILE_SEND_CHUNK int

const (
	TOX_ERR_FILE_SEND_CHUNK_OK TOX_ERR_FILE_SEND_CHUNK = iota
	TOX_ERR_FILE_SEND_CHUNK_NULL
	TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_FOUND
	TOX_ERR_FILE_SEND_CHUNK_FRIEND_NOT_CONNECTED
	TOX_ERR_FILE_SEND_CHUNK_NOT_FOUND
	TOX_ERR_FILE_SEND_CHUNK_NOT_TRANSFERRING
	TOX_ERR_FILE_SEND_CHUNK_INVALID_LENGTH
	TOX_ERR_FILE_SEND_CHUNK_SENDQ
	TOX_ERR_FILE_SEND_CHUNK_WRONG_POSITION
)

type TOX_CONFERENCE_TYPE int

const (
	TOX_CONFERENCE_TYPE_TEXT TOX_CONFERENCE_TYPE = iota
	TOX_CONFERENCE_TYPE_AV
)

type TOX_ERR_CONFERENCE_NEW int

const (
	TOX_ERR_CONFERENCE_NEW_OK TOX_ERR_CONFERENCE_NEW = iota
	TOX_ERR_CONFERENCE_NEW_INIT
)

type TOX_ERR_CONFERENCE_DELETE int

const (
	TOX_ERR_CONFERENCE_DELETE_OK TOX_ERR_CONFERENCE_DELETE = iota
	TOX_ERR_CONFERENCE_DELETE_CONFERENCE_NOT_FOUND
)

type TOX_ERR_CONFERENCE_PEER_QUERY int

const (
	TOX_ERR_CONFERENCE_PEER_QUERY_OK TOX_ERR_CONFERENCE_PEER_QUERY = iota
	TOX_ERR_CONFERENCE_PEER_QUERY_CONFERENCE_NOT_FOUND
	TOX_ERR_CONFERENCE_PEER_QUERY_PEER_NOT_FOUND
	TOX_ERR_CONFERENCE_PEER_QUERY_NO_CONNECTION
)

type TOX_ERR_CONFERENCE_INVITE int

const (
	TOX_ERR_CONFERENCE_INVITE_OK TOX_ERR_CONFERENCE_INVITE = iota
	TOX_ERR_CONFERENCE_INVITE_CONFERENCE_NOT_FOUND
	TOX_ERR_CONFERENCE_INVITE_FAIL_SEND
)

type TOX_ERR_CONFERENCE_JOIN int

const (
	TOX_ERR_CONFERENCE_JOIN_OK TOX_ERR_CONFERENCE_JOIN = iota
	TOX_ERR_CONFERENCE_JOIN_INVALID_LENGTH
	TOX_ERR_CONFERENCE_JOIN_WRONG_TYPE
	TOX_ERR_CONFERENCE_JOIN_FRIEND_NOT_FOUND
	TOX_ERR_CONFERENCE_JOIN_DUPLICATE
	TOX_ERR_CONFERENCE_JOIN_INIT_FAIL
	TOX_ERR_CONFERENCE_JOIN_FAIL_SEND
)

type TOX_ERR_CONFERENCE_SEND_MESSAGE int

const (
	TOX_ERR_CONFERENCE_SEND_MESSAGE_OK TOX_ERR_CONFERENCE_SEND_MESSAGE = iota
	TOX_ERR_CONFERENCE_SEND_MESSAGE_CONFERENCE_NOT_FOUND
	TOX_ERR_CONFERENCE_SEND_MESSAGE_TOO_LONG
	TOX_ERR_CONFERENCE_SEND_MESSAGE_NO_CONNECTION
	TOX_ERR_CONFERENCE_SEND_MESSAGE_FAIL_SEND
)

type TOX_ERR_CONFERENCE_TITLE int

const (
	TOX_ERR_CONFERENCE_TITLE_OK TOX_ERR_CONFERENCE_TITLE = iota
	TOX_ERR_CONFERENCE_TITLE_CONFERENCE_NOT_FOUND
	TOX_ERR_CONFERENCE_TITLE_INVALID_LENGTH
	TOX_ERR_CONFERENCE_TITLE_FAIL_SEND
)

type TOX_ERR_CONFERENCE_GET_TYPE int

const (
	TOX_ERR_CONFERENCE_GET_TYPE_OK TOX_ERR_CONFERENCE_GET_TYPE = iota
	TOX_ERR_CONFERENCE_GET_TYPE_CONFERENCE_NOT_FOUND
)

type TOX_ERR_FRIEND_CUSTOM_PACKET int

const (
	TOX_ERR_FRIEND_CUSTOM_PACKET_OK TOX_ERR_FRIEND_CUSTOM_PACKET = iota
	TOX_ERR_FRIEND_CUSTOM_PACKET_NULL
	TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_FOUND
	TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_CONNECTED
	TOX_ERR_FRIEND_CUSTOM_PACKET_INVALID
	TOX_ERR_FRIEND_CUSTOM_PACKET_EMPTY
	TOX_ERR_FRIEND_CUSTOM_PACKET_TOO_LONG
	TOX_ERR_FRIEND_CUSTOM_PACKET_SENDQ
)

type TOX_ERR_GET_PORT int

const (
	TOX_ERR_GET_PORT_OK TOX_ERR_GET_PORT = iota
	TOX_ERR_GET_PORT_NOT_BOUND
)

func (e TOX_ERR_OPTIONS_NEW) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_NEW) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_BOOTSTRAP) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_SET_INFO) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FRIEND_ADD) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FRIEND_DELETE) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FRIEND_BY_PUBLIC_KEY) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FRIEND_GET_PUBLIC_KEY) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FRIEND_GET_LAST_ONLINE) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FRIEND_QUERY) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_SET_TYPING) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FRIEND_SEND_MESSAGE) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FILE_CONTROL) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FILE_SEEK) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FILE_GET) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FILE_SEND) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FILE_SEND_CHUNK) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_CONFERENCE_NEW) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_CONFERENCE_DELETE) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_CONFERENCE_PEER_QUERY) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_CONFERENCE_INVITE) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_CONFERENCE_JOIN) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_CONFERENCE_SEND_MESSAGE) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_CONFERENCE_TITLE) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_CONFERENCE_GET_TYPE) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_FRIEND_CUSTOM_PACKET) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
func (e TOX_ERR_GET_PORT) Error() string {
	return fmt.Sprintf("tox err: %s (%d)", (interface{})(e).(fmt.Stringer).String(), e)
}
