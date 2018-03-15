package tox

//#include <tox/tox.h>
import "C"
import (
	"time"
)

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

// tcp tunnel
type ProtocolMagic uint16
type PacketType uint8

const (
	PACKET_TYPE_PONG PacketType = iota
	PACKET_TYPE_PING
	PACKET_TYPE_REQUESTTUNNEL
	PACKET_TYPE_ACKTUNNEL
	PACKET_TYPE_TCP
	PACKET_TYPE_TCP_FIN
	PACKET_TYPE_INVALID
)

const (
	PROTOCOL_MAGIC_V1 ProtocolMagic = 0xa26a
	PROTOCOL_MAGIC    ProtocolMagic = PROTOCOL_MAGIC_V1

	PROTOCOL_MAGIC_HIGH = byte(PROTOCOL_MAGIC >> 8)
	PROTOCOL_MAGIC_LOW  = byte(PROTOCOL_MAGIC & 0xff)

	TOX_MAX_CUSTOM_PACKET_SIZE = 1373
	PROTOCOL_BUFFER_OFFSET     = 6
	READ_BUFFER_SIZE           = TOX_MAX_CUSTOM_PACKET_SIZE - PROTOCOL_BUFFER_OFFSET
	PROTOCOL_MAX_PACKET_SIZE   = (READ_BUFFER_SIZE + PROTOCOL_BUFFER_OFFSET)

	PingMaxTryTimes       = 2
	FramePacketTypeOffset = 2
	FrameConnIdOffset     = 3
	FrameDataSizeOffset   = 4
	PacketPingSize        = 10
	DefaultPingUnit       = time.Second
	DefaultPingMultiple   = 4
	PingPongTimestampSize = 4
	DefaultTunnelQueueLen = 1024
)

var (
	defaultPinMapValue = [3]int8{0, -1, 0}
	// shoud add 4 byte timestamp, total length = 10
	// change header if const changes
	pingFrameNoData = [PacketPingSize]byte{
		0xa2, 0x6a, // PROTOCOL_MAGIC_V1
		0x01,       // PACKET_TYPE_PING
		0x00,       // conn id
		0x00, 0x04, // 4 len
		0x00, 0x00, 0x00, 0x00, // uint32 timestamp
	}
	pongFrameNoData = [PacketPingSize]byte{
		0xa2, 0x6a, // PROTOCOL_MAGIC_V1
		0x00,       // PACKET_TYPE_PONG
		0x00,       // conn id
		0x00, 0x04, // 4 len
		0x00, 0x00, 0x00, 0x00, // uint32 timestamp, len = PingPongTimestampSize
	}
	tunnelRequestFrameNoData = [PROTOCOL_BUFFER_OFFSET]byte{
		0xa2, 0x6a, // PROTOCOL_MAGIC_V1
		0x02,       // PACKET_TYPE_REQUESTTUNNEL
		0x00,       // conn id
		0x00, 0x00, // 0 len
	}
	finFrameNoData = [PROTOCOL_BUFFER_OFFSET]byte{
		0xa2, 0x6a, // PROTOCOL_MAGIC_V1
		0x05,       // PACKET_TYPE_TCP_FIN
		0x00,       // conn id
		0x00, 0x00, // 0 len
	}
)
