package tox

import (
	"log"
)

func (t *Tox) ParseLosslessPacket(friendNumber uint32, data []byte) {
	tf, ok := t.friends[friendNumber]
	if !ok {
		log.Printf("Got TCP frame with unknown friend ID #%d\n", t.recvFrom)
		return
	}

	magic := ProtocolMagic((uint16(data[0]) << 8) | (uint16(data[1])))
	if magic != PROTOCOL_MAGIC {
		log.Printf("Received data frame with invalid protocol magic number 0x%x\n", magic)
		return
	}

	packetType := data[PacketTypeOffset]
	if packetType >= PacketTypeInvalid {
		log.Printf("Received data frame with invalid PacketType 0x%x\n", packetType)
		return
	}

	t.recvFrame = data
	t.recvFrom = friendNumber
	t.recvFriend = tf
	t.recvType = packetType
	t.recvSize = uint16(len(data))

	t.handle_frame()
}

// mostly imported from https://github.com/gjedeer/tuntox.git
func (t *Tox) handle_frame() {
	switch t.recvType {
	case PacketTypeStreamBigServer:
		t.handle_stream_data_frame(t.recvFriend.bigServer, t.recvFriend.FriendBig)
	case PacketTypeStreamLittleServer:
		t.handle_stream_data_frame(t.recvFriend.littleServer, !t.recvFriend.FriendBig)
	case PacketTypeStreamOpen:
		t.handle_stream_open_frame()
	case PacketTypeStreamReady:
		t.handle_stream_ready_frame()
	case PacketTypeStreamClose:
		t.handle_stream_close_frame()
	case PacketTypePing:
		t.handle_ping_frame()
	case PacketTypePong:
		t.handle_pong_frame()
	default:
		log.Printf("Got unknown tcp packet type 0x%x from friend %d\n", t.recvType, t.recvFrom)
	}
}

func (t *Tox) handle_stream_data_frame(stream *TcpStream, clientMode bool) {
	if t.recvSize < PacketStreamDataOffset {
		log.Printf("Declared data too small: %d bytes, excepted at least %d bytes\n", t.recvSize, PacketStreamDataOffset+1)
		return
	}

	if t.recvSize > TOX_MAX_CUSTOM_PACKET_SIZE {
		log.Printf("Declared data too big (attempted buffer overflow?): %d bytes, excepted at most %d bytes\n", t.recvSize, TOX_MAX_CUSTOM_PACKET_SIZE)
		return
	}

	size := (uint16(t.recvFrame[PacketStreamDataSizeOffset]) << 8) | (uint16(t.recvFrame[PacketStreamDataSizeOffset+1]))
	if size+PacketStreamDataOffset != t.recvSize {
		log.Printf("Received frame (attempted buffer overflow?): %d bytes, excepted %d bytes\n", size+PacketStreamDataOffset, t.recvSize)
		return
	}

	if t.recvFriend.waitingAck && clientMode {
		return
	}

	if stream == nil {
		return
	}

	data := t.recvFrame[PacketStreamDataOffset:]
	total := int(size)
	sent := 0
	for sent < total {
		n, err := stream.pipe.Write(data[sent:])
		if err != nil {
			stream.close_l()
			log.Printf("Could not write to pipe of friend #%d: %v\n", t.recvFrom, err)
			return
		}
		sent += n
	}
}

func (t *Tox) handle_stream_open_frame() {
	if t.recvSize != PacketStreamOpenReadySize {
		log.Printf("Got invalid stream open frame")
		return
	}

	tf := t.recvFriend
	if tf.FriendBig {
		if tf.littleServer != nil {
			tf.littleServer.close_local_l()
			tf.littleServer = nil
		}
	} else {
		if tf.bigServer != nil {
			tf.bigServer.close_local_l()
			tf.bigServer = nil
		}
	}

	t.tunnelAcceptMu.Lock()
	defer t.tunnelAcceptMu.Unlock()
	if t.tunnelAcceptClosed {
		return
	}

	t.bufStreamReadyFrameNoData[PacketStreamOpenReadySeqOffset] = t.recvFrame[PacketStreamOpenReadySeqOffset]
	e := t.FriendSendLosslessPacket_l(tf.FriendNumber, t.bufStreamReadyFrameNoData[:], false)
	if e != 0 {
		return
	}

	stream := tf.newTcpStream(true, !tf.FriendBig)
	if tf.FriendBig {
		tf.littleServer = stream
	} else {
		tf.bigServer = stream
	}
	t.tunnelAccept <- stream
}

func (t *Tox) handle_stream_ready_frame() {
	if t.recvSize != PacketStreamOpenReadySize {
		log.Printf("Got invalid stream ready frame")
		return
	}

	tf := t.recvFriend
	if tf.waitingAck && t.recvFrame[PacketStreamOpenReadySeqOffset] == tf.dialSeq {
		tf.waitingAck = false
	}
}

func (t *Tox) handle_stream_close_frame() {
	if t.recvSize != PacketStreamCloseSize {
		log.Printf("Got invalid stream close frame")
		return
	}

	tf := t.recvFriend
	switch t.recvFrame[PacketStreamCloseSize-1] {
	case 0:
		if tf.bigServer != nil {
			tf.bigServer.remoteClosed = true
			tf.bigServer.close_local_l()
			tf.bigServer = nil
		}
	case 1:
		if tf.littleServer != nil {
			tf.littleServer.remoteClosed = true
			tf.littleServer.close_local_l()
			tf.littleServer = nil
		}
	}
}
