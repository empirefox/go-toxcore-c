package tox

//#include <tox/tox.h>
import "C"
import (
	"log"
	"time"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

type (
	sendTcpPacketData struct {
		FriendNumber uint32
		Data         []byte
		NoRetry      bool
		Result       chan error

		err toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET
	}

	PingMultipleData struct {
		FriendNumber uint32
		Multiple     int8
		Result       chan error
	}
)

func (t *Tox) sendTcpPacket_l(data *sendTcpPacketData) {
	fn := C.uint32_t(data.FriendNumber)
	data_size := C.size_t(len(data.Data))

	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET
	var e toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET

	// import from tuntox send_frame
	i := time.Duration(1)
	j := time.Duration(0)
	try := 0
	for i < 65 { // 65->2667ms 33->651ms per packet max 17->155ms
		try++

		C.tox_friend_send_lossless_packet(t.toxcore, fn, (*C.uint8_t)(&data.Data[0]), data_size, &cerr)
		e = toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET(cerr)
		switch e {
		case toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET_OK:
			goto end
		case toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET_SENDQ:
		case toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_CONNECTED:
			log.Printf("[%d] Failed to send packet to friend %d (Friend gone)\n", i, data.FriendNumber)
			goto end
		default:
			log.Printf("[%d] Failed to send packet to friend %d (err: %v)\n", i, data.FriendNumber, e)
		}

		if t.inToxIterate || data.NoRetry {
			goto end
		}

		i = i << 1
		for j = 0; j < i; j++ {
			C.tox_iterate(t.toxcore, nil)
			time.Sleep(j * time.Millisecond)
		}
	}

end:
	if e != 0 {
		if data.Result != nil {
			data.Result <- e
		} else {
			data.err = e
		}
		return
	}
	if try > 1 {
		log.Printf("Packet succeeded at try %d (friend %d)\n", try, data.FriendNumber)
	}
	if data.Result != nil {
		data.Result <- nil
	}
}
