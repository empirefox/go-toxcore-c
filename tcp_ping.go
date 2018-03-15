package tox

import (
	"encoding/binary"
	"math"
	"time"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

func (t *Tox) setPingMultiple_l(data *pingMultipleData) {
	_, ok := t.pingMap_l[data.FriendNumber]
	if !ok {
		data.Result <- toxenums.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND
		return
	}

	if data.Multiple == 0 {
		t.pingMap_l[data.FriendNumber][1] = DefaultPingMultiple
	} else {
		t.pingMap_l[data.FriendNumber][1] = data.Multiple
	}
	data.Result <- nil
	return
}

// TODO refactor: move ping to TcpConn?
func (t *Tox) doTcpPing_l() {
	ms := uint32(time.Now().UnixNano() / int64(time.Millisecond))
	binary.BigEndian.PutUint32(t.pingFrameNoData[PROTOCOL_BUFFER_OFFSET:], ms)

	for fn, ns := range t.pingMap_l {
		if ns[1] < 0 {
			continue
		}

		if ns[0] == 0 {
			if ns[2] > PingMaxTryTimes {
				// close timeout
				t.closeTcpTunnel_l(fn)
				continue
			}

			data := sendTcpPacketData{
				FriendNumber: fn,
				Data:         t.pingFrameNoData[:],
				NoRetry:      true,
			}
			t.sendTcpPacket_l(&data)
			ns[2]++ // pings_from_last_pong

			// if err, check timeout now
			if data.err != 0 && ns[2] > PingMaxTryTimes {
				// close timeout
				t.closeTcpTunnel_l(fn)
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
	_, ok := t.validConnWithFrame()
	if !ok {
		return
	}
	if len(t.tcpFrame_l.Data) != PingPongTimestampSize {
		return
	}
	copy(t.pongFrameNoData[PROTOCOL_BUFFER_OFFSET:], t.tcpFrame_l.Data)
	t.sendTcpPacket_l(&sendTcpPacketData{
		FriendNumber: t.tcpFrame_l.FriendNumber,
		Data:         t.pongFrameNoData[:],
	})
}

func (t *Tox) handle_pong_frame() {
	c, ok := t.validConnWithFrame()
	if !ok {
		return
	}

	if len(t.tcpFrame_l.Data) != PingPongTimestampSize {
		return
	}

	t.pingMap_l[c.frame.FriendNumber][2] = 0 // pings_from_last_pong

	if t.cbTcpPong != nil {
		ms := int32(time.Now().UnixNano() / int64(time.Millisecond))
		ms -= int32(binary.BigEndian.Uint32(t.tcpFrame_l.Data))
		if ms < 0 {
			ms += math.MaxInt32
		}
		t.cbTcpPong(t.tcpFrame_l.FriendNumber, uint32(ms))
	}
}
