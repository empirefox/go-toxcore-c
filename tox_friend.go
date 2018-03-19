package tox

import (
	"errors"
	"log"
)

var (
	ErrDialAlreadyLocked = errors.New("Dial already locked")
)

type ToxFriend struct {
	FriendNumber uint32
	Address      [ADDRESS_SIZE]byte
	Pubkey       [PUBLIC_KEY_SIZE]byte
	FriendBig    bool

	tox          *Tox
	bigServer    *TcpStream
	littleServer *TcpStream
	dialLocked   bool

	// A(0) B(1)
	dialSeq    byte
	waitingAck bool

	remoteAddr addr
	ping       *[3]int8 // count, trigger, pings_from_last_pong
}

func (tf *ToxFriend) SyncDail() (c *TcpStream, err error) {
	t := tf.tox
	done := make(chan struct{}, 1)
	t.DoInLoop(func() {
		c, err = tf.SyncDail_l()
		close(done)
	})
	<-done
	return
}

func (tf *ToxFriend) SyncDail_l() (c *TcpStream, err error) {
	if err = tf.lockDial(); err != nil {
		return
	}

	t := tf.tox
	tf.dialSeq += 2
	t.bufStreamOpenFrameNoData[PacketStreamOpenReadySeqOffset] = tf.dialSeq
	data := sendTcpPacketData{
		FriendNumber: tf.FriendNumber,
		Data:         t.bufStreamOpenFrameNoData[:],
	}
	t.sendTcpPacket_l(&data)
	if data.err != 0 {
		tf.unlockDial()
		err = data.err
		return
	}

	tf.waitingAck = true
	c = tf.newTcpStream(false, !tf.FriendBig)
	if tf.FriendBig {
		tf.bigServer = c
	} else {
		tf.littleServer = c
	}
	return
}

func (tf *ToxFriend) CloseStreams_l() {
	if tf.littleServer != nil {
		tf.littleServer.remoteClosed = true
		tf.littleServer.close_local_l()
		tf.littleServer = nil
	}
	if tf.bigServer != nil {
		tf.bigServer.remoteClosed = true
		tf.bigServer.close_local_l()
		tf.bigServer = nil
	}
}

func (tf *ToxFriend) lockDial() error {
	if tf.dialLocked {
		return ErrDialAlreadyLocked
	}
	tf.dialLocked = true
	return nil
}

func (tf *ToxFriend) unlockDial() {
	if !tf.dialLocked {
		log.Fatalln("unlock the unlocked dial")
	}
	tf.dialLocked = false
}
