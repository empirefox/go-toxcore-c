package tox

import (
	"context"
	"net"
	"os"
)

// Accept waits for and returns the next connection to the listener.
func (t *Tox) Accept() (net.Conn, error) {
	c, ok := <-t.tunnelAccept
	if ok {
		return c, nil
	}
	return nil, os.ErrClosed
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (t *Tox) Close() error {
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
func (t *Tox) Dial_l(friendNumber uint32) (net.Conn, error) {
	return t.CreateTcpTunnel_l(friendNumber)
}

func (t *Tox) Dial(ctx context.Context, pubkey *[PUBLIC_KEY_SIZE]byte) (net.Conn, error) {
	resultCh := make(chan *CreateTcpTunnelResult, 1)
	t.DoInLoop(&CreateTcpTunnelByPublicKeyData{
		Pubkey: pubkey,
		Result: resultCh,
	})
	select {
	case result := <-resultCh:
		return result.Conn, result.Error
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
