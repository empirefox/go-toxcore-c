package tox

import (
	"log"
)

func (t *Tox) ParseLosslessPacket(friendNumber uint32, data []byte) {
	total := uint16(len(data))
	if total < PROTOCOL_BUFFER_OFFSET {
		log.Printf("Received too short data frame - only %d bytes, at least %d expected\n", total, PROTOCOL_BUFFER_OFFSET)
		return
	}

	magic := ProtocolMagic((uint16(data[0]) << 8) | (uint16(data[1])))
	if magic != PROTOCOL_MAGIC {
		log.Printf("Received data frame with invalid protocol magic number 0x%x\n", magic)
		return
	}

	size := (uint16(data[FrameDataSizeOffset]) << 8) | (uint16(data[FrameDataSizeOffset+1])) + PROTOCOL_BUFFER_OFFSET
	if size != total {
		log.Printf("Received frame size (attempted buffer overflow?): %d bytes, excepted %d bytes\n", total, size)
		return
	}

	if size > PROTOCOL_MAX_PACKET_SIZE {
		log.Printf("Declared data length too big (attempted buffer overflow?): %d bytes, excepted at most %d bytes\n", size, PROTOCOL_MAX_PACKET_SIZE)
		return
	}

	t.tcpFrame_l = TcpFrame{
		FriendNumber: friendNumber,
		Magic:        magic,
		PacketType:   PacketType(data[FramePacketTypeOffset]),
		ConnID:       data[FrameConnIdOffset],
		Data:         data[PROTOCOL_BUFFER_OFFSET:],
	}
	if t.tcpFrame_l.PacketType >= PACKET_TYPE_INVALID {
		log.Printf("Received data frame with invalid PacketType 0x%x\n", t.tcpFrame_l.PacketType)
		return
	}

	t.handle_frame()
}

// mostly imported from https://github.com/gjedeer/tuntox.git
func (t *Tox) handle_frame() {
	switch t.tcpFrame_l.PacketType {
	case PACKET_TYPE_PING:
		t.handle_ping_frame()
	case PACKET_TYPE_PONG:
		t.handle_pong_frame()
	case PACKET_TYPE_TCP:
		t.handle_tcp_frame()
	case PACKET_TYPE_REQUESTTUNNEL:
		t.handle_request_tunnel_frame()
	case PACKET_TYPE_ACKTUNNEL:
		t.handle_acktunnel_frame()
	case PACKET_TYPE_TCP_FIN:
		t.handle_tcp_fin_frame()
	default:
		log.Printf("Got unknown tcp packet type 0x%x from friend %d\n", t.tcpFrame_l.PacketType, t.tcpFrame_l.FriendNumber)
	}
}

func (t *Tox) handle_tcp_frame() {
	c, ok := t.validConnWithFrame()
	if !ok {
		return
	}

	size := len(t.tcpFrame_l.Data)
	sent := 0
	for sent < size {
		n, err := c.pipe.Write(t.tcpFrame_l.Data[sent:])
		if err != nil {
			log.Printf("Could not write to pipe of friend #%d: %v\n", c.frame.FriendNumber, err)
			t.closeTcpTunnel_l(t.tcpFrame_l.FriendNumber)
			return
		}
		sent += n
	}
}

func (t *Tox) handle_request_tunnel_frame() {
	// close exist old conn
	c, ok := t.tunnels_l[t.tcpFrame_l.FriendNumber]
	if ok {
		if t.tcpFrame_l.ConnID == c.frame.ConnID || !c.server {
			return
		}
		// TODO really alow this?
		t.closeTcpTunnel_l(t.tcpFrame_l.FriendNumber)
	}

	// create server conn
	t.tunnelAcceptMu.Lock()
	if t.tunnelAcceptClosed {
		t.tunnelAcceptMu.Unlock()
		t.finFrameNoData[FrameConnIdOffset] = t.tcpFrame_l.ConnID
		data := sendTcpPacketData{
			FriendNumber: c.frame.FriendNumber,
			Data:         t.finFrameNoData[:],
		}
		t.sendTcpPacket_l(&data)
		return
	}
	t.tunnelAccept <- t.newTcpConn(t.tcpFrame_l.FriendNumber, t.tcpFrame_l.ConnID, true)
	t.tunnelAcceptMu.Unlock()

	// TODO send_tunnel_ack_frame
}

func (t *Tox) handle_acktunnel_frame() {
	log.Printf("ACK should not got here!")
}

func (t *Tox) handle_tcp_fin_frame() {
	_, ok := t.validConnWithFrame()
	if !ok {
		return
	}

	t.closeTcpTunnel_l(t.tcpFrame_l.FriendNumber)
}
