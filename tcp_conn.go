package tox

import (
	"encoding/hex"
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"
)

var ErrNoDeadline = os.ErrNoDeadline

type TcpConn struct {
	t      *Tox
	server bool

	// TODO try io.Pipe?
	pipe       pipe
	remoteAddr addr
	buf        [PROTOCOL_MAX_PACKET_SIZE]byte
	frame      TcpFrame
	result     chan error
	closed     atomic.Value
}

func (t *Tox) newTcpConn(friendNumber uint32, connid byte, server bool) *TcpConn {
	c := TcpConn{
		t:      t,
		server: server,
		pipe:   pipe{b: new(dataBuffer)},
		frame: TcpFrame{
			FriendNumber: friendNumber,
			Magic:        PROTOCOL_MAGIC,
			PacketType:   PACKET_TYPE_TCP,
			ConnID:       connid,
		},
		result: make(chan error, 1),
	}
	c.frame.WriteProtocol(c.buf[:])
	pubkey, ok := t.FriendGetPublicKey(friendNumber)
	if ok {
		c.remoteAddr = addr(hex.EncodeToString(pubkey[:]))
	} else {
		c.remoteAddr = addr("unknown remote pubkey")
	}
	c.closed.Store(false)

	t.tunnels_l[friendNumber] = &c
	return &c
}

func (c *TcpConn) Close() error {
	if c.closed.Load().(bool) {
		return nil
	}
	result := make(chan error, 1)
	c.t.DoInLoop(&CloseTcpTunnelData{
		FriendNumber: c.frame.FriendNumber,
		Result:       result,
	})
	return <-result
}

// Write write directly to underline implement
func (c *TcpConn) Write(p []byte) (n int, err error) {
	for {
		dataSize := copy(c.buf[PROTOCOL_BUFFER_OFFSET:], p[n:])
		c.buf[FrameDataSizeOffset] = byte(dataSize >> 8)
		c.buf[FrameDataSizeOffset+1] = byte(dataSize)
		c.t.DoInLoop(&sendTcpPacketData{
			FriendNumber: c.frame.FriendNumber,
			Data:         c.buf[:PROTOCOL_BUFFER_OFFSET+dataSize],
			Result:       c.result,
		})
		err = <-c.result
		if err != nil {
			return
		}
		n += dataSize
		if dataSize < READ_BUFFER_SIZE {
			return
		}
	}
	return
}

func (c *TcpConn) Read(p []byte) (n int, err error) { return c.pipe.Read(p) }
func (c *TcpConn) LocalAddr() net.Addr              { return &c.t.localAddr }
func (c *TcpConn) RemoteAddr() net.Addr             { return &c.remoteAddr }

func (c *TcpConn) SetDeadline(t time.Time) error {
	log.Println("SetDeadline should not be called")
	return ErrNoDeadline
}
func (c *TcpConn) SetReadDeadline(t time.Time) error {
	log.Println("SetReadDeadline should not be called")
	return ErrNoDeadline
}
func (c *TcpConn) SetWriteDeadline(t time.Time) error {
	log.Println("SetWriteDeadline should not be called")
	return ErrNoDeadline
}

type addr string

func (a *addr) Network() string { return "tox" }
func (a *addr) String() string  { return string(*a) }
