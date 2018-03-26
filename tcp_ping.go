package tox

import (
	"encoding/binary"
	"log"
	"math"
	"time"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

func (t *Tox) SetPingMultiple_l(friendNumber uint32, multiple int8) error {
	tf, ok := t.friends[friendNumber]
	if !ok {
		return toxenums.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND
	}

	if multiple == 0 {
		tf.ping[1] = DefaultPingMultiple
	} else {
		tf.ping[1] = multiple
	}
	return nil
}

// TODO refactor: move ping to TcpConn?
func (t *Tox) doTcpPing_l() {
	ms := uint32(time.Now().UnixNano() / int64(time.Millisecond))
	binary.BigEndian.PutUint32(t.bufPingFrameNoData[PacketPingPongTimeOffset:], ms)

	for fn, tf := range t.friends {
		ns := tf.ping
		if ns[1] < 0 {
			continue
		}

		if ns[0] == 0 {
			if ns[2] > PingMaxTryTimes {
				tf.CloseStreams_l()
				continue
			}

			e := t.FriendSendLosslessPacket_l(fn, t.bufPingFrameNoData[:], true)
			ns[2]++ // pings_from_last_pong

			// if err, check timeout now
			if e != 0 && ns[2] > PingMaxTryTimes {
				tf.CloseStreams_l()
				continue
			}
		}
		if ns[0] < ns[1] {
			ns[0]++
		} else {
			ns[0] = 0
		}
	}
}

func (t *Tox) handle_ping_frame() {
	if t.recvSize != PacketPingPongSize {
		log.Printf("Got invalid ping frame")
		return
	}

	copy(t.bufPongFrameNoData[PacketPingPongTimeOffset:], t.recvFrame[PacketPingPongTimeOffset:])
	t.FriendSendLosslessPacket_l(t.recvFrom, t.bufPongFrameNoData[:], false)
}

func (t *Tox) handle_pong_frame() {
	if t.recvSize != PacketPingPongSize {
		log.Printf("Got invalid pong frame")
		return
	}

	t.recvFriend.ping[2] = 0 // pings_from_last_pong

	if t.cbTcpPong != nil {
		ms := int32(time.Now().UnixNano() / int64(time.Millisecond))
		ms -= int32(binary.BigEndian.Uint32(t.recvFrame[PacketPingPongTimeOffset:]))
		if ms < 0 {
			ms += math.MaxInt32
		}
		t.cbTcpPong(t.recvFrom, uint32(ms))
	}
}
