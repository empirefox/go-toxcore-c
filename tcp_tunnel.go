package tox

//#include <tox/tox.h>
import "C"
import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

var (
	ErrTunnelAlreadyCreated = errors.New("Tunnel already created")
)

type TcpFrame struct {
	FriendNumber uint32
	Magic        ProtocolMagic
	PacketType   PacketType
	ConnID       byte
	Data         []byte
}

func (f *TcpFrame) WriteProtocol(buf []byte) {
	buf[0], buf[1] = byte(f.Magic>>8), byte(f.Magic)
	buf[2], buf[3] = byte(f.PacketType), f.ConnID
}

type (
	CreateTcpTunnelData struct {
		FriendNumber uint32
		Result       chan *CreateTcpTunnelResult
	}
	CreateTcpTunnelByPublicKeyData struct {
		Pubkey *[PUBLIC_KEY_SIZE]byte
		Result chan *CreateTcpTunnelResult
	}
	CreateTcpTunnelResult struct {
		Conn  *TcpConn
		Error error
	}

	CloseTcpTunnelData struct {
		FriendNumber uint32
		Result       chan error
	}

	sendTcpPacketData struct {
		FriendNumber uint32
		Data         []byte
		NoRetry      bool
		Result       chan error

		err toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET
	}

	pingMultipleData struct {
		FriendNumber uint32
		Multiple     int8
		Result       chan error
	}
)

func (t *Tox) createTcpTunnel_l(data *CreateTcpTunnelData) {
	c, err := t.CreateTcpTunnel_l(data.FriendNumber)
	data.Result <- &CreateTcpTunnelResult{
		Conn:  c,
		Error: err,
	}
}

func (t *Tox) createTcpTunnelByPublicKey_l(data *CreateTcpTunnelByPublicKeyData) {
	friendNumber, ok := t.FriendByPublicKey(data.Pubkey)
	if !ok {
		data.Result <- &CreateTcpTunnelResult{Error: toxenums.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND}
		return
	}

	c, err := t.CreateTcpTunnel_l(friendNumber)
	data.Result <- &CreateTcpTunnelResult{
		Conn:  c,
		Error: err,
	}
}

func (t *Tox) CreateTcpTunnel_l(friendNumber uint32) (*TcpConn, error) {
	if !t.FriendExists(friendNumber) {
		return nil, toxenums.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND
	}

	_, ok := t.tunnels_l[friendNumber]
	if ok {
		return nil, ErrTunnelAlreadyCreated
	}

	connid := t.tunnelids_l[friendNumber]
	t.tunnelids_l[friendNumber] += 1

	t.tunnelRequestFrameNoData[FrameConnIdOffset] = connid
	data := sendTcpPacketData{
		FriendNumber: friendNumber,
		Data:         t.tunnelRequestFrameNoData[:],
	}
	t.sendTcpPacket_l(&data)
	if data.err != 0 {
		log.Printf("Fail to send packet to friend #%d with connid: %d\n", friendNumber, connid)
		return nil, data.err
	}
	// send_tunnel_request_packet
	log.Printf("Sending packet to friend #%d with connid: %d\n", friendNumber, connid)

	return t.newTcpConn(friendNumber, connid, false), nil
}

func (t *Tox) CloseTcpTunnel_l(data *CloseTcpTunnelData) {
	data.Result <- t.closeTcpTunnel_l(data.FriendNumber)
}

func (t *Tox) closeTcpTunnel_l(friendNumber uint32) error {
	c, ok := t.tunnels_l[friendNumber]
	if !ok {
		return toxenums.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND
	}

	t.finFrameNoData[FrameConnIdOffset] = c.frame.ConnID
	data := sendTcpPacketData{
		FriendNumber: c.frame.FriendNumber,
		Data:         t.finFrameNoData[:],
	}
	t.sendTcpPacket_l(&data)

	delete(t.tunnels_l, c.frame.FriendNumber)
	c.pipe.CloseWithError(io.EOF)
	c.closed.Store(true)
	return data.err
}

func (t *Tox) sendTcpPacket_l(data *sendTcpPacketData) {
	fn := C.uint32_t(data.FriendNumber)
	data_size := C.size_t(len(data.Data))

	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET
	var e toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET

	// import from tuntox send_frame
	i := time.Duration(1)
	j := time.Duration(0)
	try := 0
	for i < 33 { // 33->651ms per packet max 17->155ms
		try++

		C.tox_friend_send_lossless_packet(t.toxcore, fn, (*C.uint8_t)(&data.Data[0]), data_size, &cerr)
		e = toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET(cerr)
		switch e {
		case toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET_OK:
			goto end
		case toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET_SENDQ:
			log.Printf("[%d] Failed to send packet to friend %d (Packet queue is full)\n", i, data.FriendNumber)
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

func (t *Tox) validConnWithFrame() (c *TcpConn, ok bool) {
	c, ok = t.tunnels_l[t.tcpFrame_l.FriendNumber]
	if !ok {
		log.Printf("Got TCP frame with unknown friend ID #%d\n", t.tcpFrame_l.FriendNumber)
		return
	}

	if t.tcpFrame_l.ConnID != c.frame.ConnID {
		ok = false
		log.Printf("Got TCP frame with unknown tunnel ID %d of friend #%d\n", t.tcpFrame_l.ConnID, t.tcpFrame_l.FriendNumber)
		return
	}

	// should not ever happen bellow
	if t.tcpFrame_l.FriendNumber != c.frame.FriendNumber {
		ok = false
		log.Printf("Friend #%d tried to send packet to a tunnel which belongs to #%d\n", t.tcpFrame_l.FriendNumber, c.frame.FriendNumber)
		return
	}
	return
}
