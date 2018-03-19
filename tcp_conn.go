package tox

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	ErrNoDeadline = os.ErrNoDeadline
)

type TcpStream struct {
	tf        *ToxFriend
	server    bool
	big       bool
	bigServer bool
	closeTag  byte

	pipe pipe

	buf [TOX_MAX_CUSTOM_PACKET_SIZE]byte

	result       chan error
	closed       bool
	remoteClosed bool
}

func (tf *ToxFriend) newTcpStream(server, big bool) *TcpStream {
	bigServer := server && big || !server && !big
	var closeTag byte = 1
	packetType := PacketTypeStreamLittleServer
	if bigServer {
		closeTag = 0
		packetType = PacketTypeStreamBigServer
	}
	c := TcpStream{
		tf:        tf,
		server:    server,
		big:       big,
		bigServer: bigServer,
		closeTag:  closeTag,
		pipe:      pipe{b: new(dataBuffer)},
		buf:       [TOX_MAX_CUSTOM_PACKET_SIZE]byte{PROTOCOL_MAGIC_HIGH, PROTOCOL_MAGIC_LOW, packetType},
		result:    make(chan error, 1),
	}

	return &c
}

func (c *TcpStream) Close() (err error) {
	done := make(chan struct{}, 1)
	c.tf.tox.DoInLoop(func() {
		err = c.close_l()
		close(done)
	})
	<-done
	return
}

func (c *TcpStream) close_l() (err error) {
	if err = c.close_local_l(); err != nil {
		return
	}
	t := c.tf.tox
	if !c.remoteClosed {
		t.bufStreamCloseFrameNoData[PacketStreamCloseSize-1] = c.closeTag
		data := sendTcpPacketData{
			FriendNumber: c.tf.FriendNumber,
			Data:         t.bufStreamCloseFrameNoData[:],
		}
		t.sendTcpPacket_l(&data)
		if data.err != 0 {
			err = data.err
		}
	}
	return
}

func (c *TcpStream) close_local_l() error {
	if c.closed {
		return io.EOF
	}
	c.closed = true
	c.pipe.CloseWithError(io.EOF)
	if c.bigServer {
		c.tf.bigServer = nil
	} else {
		c.tf.littleServer = nil
	}
	if !c.server {
		c.tf.unlockDial()
	}
	return nil
}

// Write write directly to underline implement
func (c *TcpStream) Write(p []byte) (n int, err error) {
	total := len(p)
	for n < total {
		dataSize := copy(c.buf[PacketStreamDataOffset:], p[n:])
		c.buf[PacketStreamDataSizeOffset] = byte(dataSize >> 8)
		c.buf[PacketStreamDataSizeOffset+1] = byte(dataSize)
		c.tf.tox.DoInLoop(func() {
			if c.closed {
				c.result <- io.EOF
				return
			}

			data := sendTcpPacketData{
				FriendNumber: c.tf.FriendNumber,
				Data:         c.buf[:PacketStreamDataOffset+dataSize],
				Result:       c.result,
			}
			c.tf.tox.sendTcpPacket_l(&data)
		})
		err = <-c.result
		if err != nil {
			return
		}
		n += dataSize
	}
	return
}

func (c *TcpStream) Read(p []byte) (n int, err error) { return c.pipe.Read(p) }
func (c *TcpStream) LocalAddr() net.Addr              { return &c.tf.tox.localAddr }
func (c *TcpStream) RemoteAddr() net.Addr             { return &c.tf.remoteAddr }

func (c *TcpStream) SetDeadline(t time.Time) error {
	log.Println("SetDeadline should not be called")
	return ErrNoDeadline
}
func (c *TcpStream) SetReadDeadline(t time.Time) error {
	log.Println("SetReadDeadline should not be called")
	return ErrNoDeadline
}
func (c *TcpStream) SetWriteDeadline(t time.Time) error {
	log.Println("SetWriteDeadline should not be called")
	return ErrNoDeadline
}

type addr string

func (a *addr) Network() string { return "tox" }
func (a *addr) String() string  { return string(*a) }
