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

type ProtocolMagic uint16

const (
	PacketTypeStreamOpen byte = iota
	PacketTypeStreamReady
	PacketTypeStreamClose
	PacketTypeStreamBigServer
	PacketTypeStreamLittleServer
	PacketTypePing
	PacketTypePong
	PacketTypeInvalid
)

type TcpStreamState uint8

const (
	TcpStreamStateConnected TcpStreamState = iota
	TcpStreamStateHalfReady
	TcpStreamStateReady
	TcpStreamStateHalfClose
	TcpStreamStateClosed
	TcpStreamStateLost
)

const (
	PROTOCOL_MAGIC_V1 ProtocolMagic = 0xa26a // 0b1010001001101010
	PROTOCOL_MAGIC    ProtocolMagic = PROTOCOL_MAGIC_V1

	PROTOCOL_MAGIC_HIGH = byte(PROTOCOL_MAGIC >> 8)
	PROTOCOL_MAGIC_LOW  = byte(PROTOCOL_MAGIC & 0xff)

	TOX_MAX_CUSTOM_PACKET_SIZE = 1373

	PingMaxTryTimes     = 2
	DefaultPingUnit     = time.Second
	DefaultPingMultiple = 4

	PacketTypeOffset               = 2
	PacketMinSize                  = 4
	PacketPingPongSize             = 7
	PacketPingPongTimeOffset       = 3
	PacketStreamOpenReadySize      = 4
	PacketStreamOpenReadySeqOffset = 3
	PacketStreamCloseSize          = 4
	PacketStreamDataSizeOffset     = 3
	PacketStreamDataOffset         = 5
)

var (
	defaultPinMapValue = [3]int8{0, -1, 0}
	// shoud add 4 byte timestamp, total length = 10
	// change header if const changes
	pingFrameNoData = [PacketPingPongSize]byte{
		0xa2, 0x6a, // PROTOCOL_MAGIC_V1
		PacketTypePing,         // PACKET_TYPE_PING
		0x00, 0x00, 0x00, 0x00, // uint32 timestamp
	}
	pongFrameNoData = [PacketPingPongSize]byte{
		0xa2, 0x6a, // PROTOCOL_MAGIC_V1
		PacketTypePong,         // PACKET_TYPE_PONG
		0x00, 0x00, 0x00, 0x00, // uint32 timestamp, len = PingPongTimestampSize
	}
	streamOpenFrameNoData = [PacketStreamOpenReadySize]byte{
		0xa2, 0x6a, // PROTOCOL_MAGIC_V1
		PacketTypeStreamOpen, // PACKET_TYPE_REQUESTTUNNEL
		0x00,                 // conn id
	}
	streamReadyFrameNoData = [PacketStreamOpenReadySize]byte{
		0xa2, 0x6a, // PROTOCOL_MAGIC_V1
		PacketTypeStreamReady, // PACKET_TYPE_REQUESTTUNNEL
		0x00, // conn id
	}
	streamCloseFrameNoData = [PacketStreamCloseSize]byte{
		0xa2, 0x6a, // PROTOCOL_MAGIC_V1
		PacketTypeStreamClose, // PACKET_TYPE_TCP_FIN
		0x00, // close tag bigServer(0)/littleServer(1)
	}
)
