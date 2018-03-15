package tox

//#include <stdlib.h>
//#include <tox/tox.h>
import "C"
import (
	"sync"
	"time"
	"unsafe"

	"github.com/TokTok/go-toxcore-c/toxenums"
	"github.com/phayes/freeport"
)

// Tox method end with _l should be used inside of callbacks or before Run() called.
// All xxx_l can do the same with DoInLoop((*)xxxData{}) out side of callbacks.
type Tox struct {
	opts    *ToxOptions
	toxcore *C.Tox

	// uint32 -> *[PUBLIC_KEY_SIZE]byte
	friendIdToPk sync.Map

	// [PUBLIC_KEY_SIZE]byte -> uint32
	friendPkToId sync.Map

	// init with create
	Address [ADDRESS_SIZE]byte
	Pubkey  [PUBLIC_KEY_SIZE]byte
	Secret  [SECRET_KEY_SIZE]byte

	cb_friend_request           cb_friend_request_ftype
	cb_friend_message           cb_friend_message_ftype
	cb_friend_name              cb_friend_name_ftype
	cb_friend_status_message    cb_friend_status_message_ftype
	cb_friend_status            cb_friend_status_ftype
	cb_friend_connection_status cb_friend_connection_status_ftype
	cb_friend_typing            cb_friend_typing_ftype
	cb_friend_read_receipt      cb_friend_read_receipt_ftype
	cb_friend_lossy_packet      cb_friend_lossy_packet_ftype
	cb_friend_lossless_packet   cb_friend_lossless_packet_ftype
	cb_self_connection_status   cb_self_connection_status_ftype

	cb_conference_invite            cb_conference_invite_ftype
	cb_conference_message           cb_conference_message_ftype
	cb_conference_title             cb_conference_title_ftype
	cb_conference_peer_name         cb_conference_peer_name_ftype
	cb_conference_peer_list_changed cb_conference_peer_list_changed_ftype

	cb_file_recv_control  cb_file_recv_control_ftype
	cb_file_recv          cb_file_recv_ftype
	cb_file_recv_chunk    cb_file_recv_chunk_ftype
	cb_file_chunk_request cb_file_chunk_request_ftype

	inToxIterate bool

	cbTcpPong     CallbackTcpPongFn
	cbPostIterate []CallbackPostIterateOnceFn

	pingFrameNoData          [PacketPingSize]byte
	pongFrameNoData          [PacketPingSize]byte
	tunnelRequestFrameNoData [PROTOCOL_BUFFER_OFFSET]byte
	finFrameNoData           [PROTOCOL_BUFFER_OFFSET]byte

	// this is buf of recv frame
	tcpFrame_l TcpFrame

	localAddr addr
	pingUnit  time.Duration
	pingMap_l map[uint32]*[3]int8 // count, trigger, pings_from_last_pong

	// TODO split to 2 map, add tag to protocol header,
	// we can support both client/server mode with the same peer.
	tunnels_l   map[uint32]*TcpConn
	tunnelids_l map[uint32]byte

	tunnelAccept       chan *TcpConn
	tunnelAcceptMu     sync.Mutex
	tunnelAcceptClosed bool

	chLoopRequest chan interface{}

	stopOnce sync.Once
	stop     chan struct{}
	stopped  chan struct{}
	killOnce sync.Once
	killed   chan struct{}
}

func NewTox(opts *ToxOptions) (*Tox, error) {
	if opts == nil {
		opts = NewToxOptions()
	}
	toxopts := opts.toCToxOptions()
	defer func() {
		if opts.Proxy_host != "" {
			C.free(unsafe.Pointer(C.tox_options_get_proxy_host(toxopts)))
		}
		C.tox_options_free(toxopts)
	}()

	var cerr C.TOX_ERR_NEW
	decrypt := opts.Decrypt
	toxcore := C.tox_new(toxopts, &cerr)
	for cerr != 0 {
		switch err := toxenums.TOX_ERR_NEW(cerr); err {
		case toxenums.TOX_ERR_NEW_PORT_ALLOC:
			if opts.Tcp_port == 0 {
				return nil, err
			}
			if !opts.AutoTcpPortIfErr {
				return nil, err
			}
			port, ferr := freeport.GetFreePort()
			if ferr != nil {
				if !opts.DisableTcpPortIfAutoErr {
					return nil, ferr
				}
				port = 0
			}
			opts.Tcp_port = uint16(port)
			C.tox_options_set_tcp_port(toxopts, (C.uint16_t)(opts.Tcp_port))
		case toxenums.TOX_ERR_NEW_PROXY_BAD_TYPE,
			toxenums.TOX_ERR_NEW_PROXY_BAD_HOST,
			toxenums.TOX_ERR_NEW_PROXY_BAD_PORT,
			toxenums.TOX_ERR_NEW_PROXY_NOT_FOUND:
			if opts.Proxy_type == toxenums.TOX_PROXY_TYPE_NONE {
				return nil, err
			}
			if !opts.ProxyToNoneIfErr {
				return nil, err
			}
			opts.Proxy_type = toxenums.TOX_PROXY_TYPE_NONE
			C.tox_options_set_proxy_type(toxopts, C.TOX_PROXY_TYPE_NONE)
		case toxenums.TOX_ERR_NEW_LOAD_ENCRYPTED:
			if decrypt == nil {
				return nil, err
			}
			data, derr := decrypt(opts.Savedata_data)
			if derr != nil {
				return nil, derr
			}
			decrypt = nil
			opts.Savedata_data = data
			C.tox_options_set_savedata_data(toxopts, (*C.uint8_t)(&data[0]), C.size_t(len(data)))
		default:
			return nil, err
		}
		toxcore = C.tox_new(toxopts, &cerr)
	}

	// TODO make chan len configurable
	t := Tox{
		opts:    opts,
		toxcore: toxcore,

		pingFrameNoData:          pingFrameNoData,
		pongFrameNoData:          pongFrameNoData,
		tunnelRequestFrameNoData: tunnelRequestFrameNoData,
		finFrameNoData:           finFrameNoData,

		pingUnit:    opts.PingUnit,
		pingMap_l:   make(map[uint32]*[3]int8),
		tunnels_l:   make(map[uint32]*TcpConn),
		tunnelids_l: make(map[uint32]byte),

		tunnelAccept: make(chan *TcpConn, 16),

		chLoopRequest: make(chan interface{}, 1024),
		stop:          make(chan struct{}),
		stopped:       make(chan struct{}),
		killed:        make(chan struct{}),
	}

	if opts.Savedata_type == toxenums.TOX_SAVEDATA_TYPE_SECRET_KEY {
		t.SelfSetNospam_l(opts.NospamIfSecretType)
	}

	C.tox_self_get_address(toxcore, (*C.uint8_t)(&t.Address[0]))
	C.tox_self_get_public_key(toxcore, (*C.uint8_t)(&t.Pubkey[0]))
	C.tox_self_get_secret_key(toxcore, (*C.uint8_t)(&t.Secret[0]))

	size := C.tox_self_get_friend_list_size(toxcore)
	if size != 0 {
		list := make([]uint32, size)
		C.tox_self_get_friend_list(t.toxcore, (*C.uint32_t)(&list[0]))
		for _, friendNumber := range list {
			var pubkey [PUBLIC_KEY_SIZE]byte
			C.tox_friend_get_public_key(toxcore, C.uint32_t(friendNumber), (*C.uint8_t)(&pubkey[0]), nil)
			t.onFriendAdded_l(friendNumber, &pubkey)
		}
	}

	t.localAddr = addr(t.Pubkey[:])

	cbUserDatas.set(toxcore, &t)
	return &t, nil
}

// Kill only used before Run. If Run started, use StopAndKill.
func (t *Tox) Kill() {
	t.killOnce.Do(func() {
		C.tox_kill(t.toxcore)
		t.toxcore = nil
		close(t.killed)
	})
}

// StopAndKill only used after Run. If Run not started, use Kill.
func (t *Tox) StopAndKill() {
	t.stopOnce.Do(func() { close(t.stop) })
	<-t.stopped
}

// CallbackPostIterate heavey work must not be done here

type CallbackPostIterateOnceFn func() time.Duration
type CallbackTcpPongFn func(friendNumber uint32, ms uint32)

// mainly used to delete friend from callbacks
func (t *Tox) CallbackPostIterateOnce_l(cb CallbackPostIterateOnceFn) {
	t.cbPostIterate = append(t.cbPostIterate, cb)
}
func (t *Tox) CallbackTcpPong(cb CallbackTcpPongFn) { t.cbTcpPong = cb }
