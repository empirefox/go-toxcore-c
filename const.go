package tox

/*
#include <tox/tox.h>
*/
import "C"

const (
	PUBLIC_KEY_SIZE = 32
	SECRET_KEY_SIZE = 32
	ADDRESS_SIZE    = 38
	HASH_LENGTH     = 32
	FILE_ID_LENGTH  = 32
)

var (
	MAX_NAME_LENGTH           int = int(C.tox_max_name_length())
	MAX_STATUS_MESSAGE_LENGTH int = int(C.tox_max_status_message_length())
	MAX_FRIEND_REQUEST_LENGTH int = int(C.tox_max_friend_request_length())
	MAX_MESSAGE_LENGTH        int = int(C.tox_max_message_length())
	MAX_CUSTOM_PACKET_SIZE    int = int(C.tox_max_custom_packet_size())
	MAX_FILENAME_LENGTH       int = int(C.tox_max_filename_length())
)
