package tox

import (
	"errors"
	"net"
	"os"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

// Accept waits for and returns the next connection to the listener.
func (t *Tox) Accept() (net.Conn, error) {
	c, ok := <-t.tunnelAccept
	if ok {
		return c, nil
	}
	return nil, os.ErrClosed
}

// Close Close implememts net.Listener and will do nothing. Kill tox will close the
// listener.
func (t *Tox) Close() error {
	return errors.New("Close tox listener will do nothing, kill tox please")
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (t *Tox) close() error {
	t.tunnelAcceptMu.Lock()
	defer t.tunnelAcceptMu.Unlock()

	if t.tunnelAcceptClosed {
		return nil
	}
	t.tunnelAcceptClosed = true
	close(t.tunnelAccept)
	return nil
}

// Addr returns the listener's network address.
func (t *Tox) Addr() net.Addr { return &t.localAddr }

// Dial_l dail from callbacks and it will not auto retry. It will not block the
// queue. If failed, save the pubkey then try Dial out side of callbacks later.
// Usage see https://github.com/empirefox/hybrid/blob/master/hybridtox/dc_tox.go#L118
func (t *Tox) Dial_l(friendNumber uint32) (net.Conn, error) {
	tf, ok := t.friends[friendNumber]
	if !ok {
		return nil, toxenums.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND
	}
	return tf.SyncDail_l()
}
