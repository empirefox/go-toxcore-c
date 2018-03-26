package tox

/*
#include <tox/tox.h>

void callbackFriendRequestWrapperForC(Tox *, uint8_t *, uint8_t *, uint16_t, void*);
void callbackFriendMessageWrapperForC(Tox *, uint32_t, int, uint8_t*, uint32_t, void*);
void callbackFriendNameWrapperForC(Tox *, uint32_t, uint8_t*, uint32_t, void*);
void callbackFriendStatusMessageWrapperForC(Tox *, uint32_t, uint8_t*, uint32_t, void*);
void callbackFriendStatusWrapperForC(Tox *, uint32_t, TOX_USER_STATUS, void*);
void callbackFriendConnectionStatusWrapperForC(Tox *, uint32_t, TOX_CONNECTION, void*);
void callbackFriendTypingWrapperForC(Tox *, uint32_t, uint8_t, void*);
void callbackFriendReadReceiptWrapperForC(Tox *, uint32_t, uint32_t, void*);
void callbackFriendLossyPacketWrapperForC(Tox *, uint32_t, uint8_t*, size_t, void*);
void callbackFriendLosslessPacketWrapperForC(Tox *, uint32_t, uint8_t*, size_t, void*);
void callbackSelfConnectionStatusWrapperForC(Tox *, TOX_CONNECTION, void*);
void callbackFileRecvControlWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number,
                                      TOX_FILE_CONTROL control, void *user_data);
void callbackFileRecvWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number, uint32_t kind,
                               uint64_t file_size, uint8_t *filename, size_t filename_length, void *user_data);
void callbackFileRecvChunkWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number, uint64_t position,
                                    uint8_t *data, size_t length, void *user_data);
void callbackFileChunkRequestWrapperForC(Tox *tox, uint32_t friend_number, uint32_t file_number, uint64_t position,
                                       size_t length, void *user_data);
*/
import "C"
import (
	"unsafe"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

// friend callback type
type cb_friend_request_ftype func(pubkey *[PUBLIC_KEY_SIZE]byte, message []byte)
type cb_friend_message_ftype func(friendNumber uint32, mtype toxenums.TOX_MESSAGE_TYPE, message []byte)
type cb_friend_name_ftype func(friendNumber uint32, newName string)
type cb_friend_status_message_ftype func(friendNumber uint32, newStatus string)
type cb_friend_status_ftype func(friendNumber uint32, status toxenums.TOX_USER_STATUS)
type cb_friend_connection_status_ftype func(friendNumber uint32, status toxenums.TOX_CONNECTION)
type cb_friend_typing_ftype func(friendNumber uint32, isTyping uint8)
type cb_friend_read_receipt_ftype func(friendNumber uint32, receipt uint32)
type cb_friend_lossy_packet_ftype func(friendNumber uint32, data []byte)
type cb_friend_lossless_packet_ftype func(friendNumber uint32, data []byte)

// self callback type
type cb_self_connection_status_ftype func(status toxenums.TOX_CONNECTION)

// file callback type
type cb_file_recv_control_ftype func(friendNumber uint32, fileNumber uint32, control toxenums.TOX_FILE_CONTROL)
type cb_file_recv_ftype func(friendNumber uint32, fileNumber uint32, kind toxenums.TOX_FILE_KIND, fileSize uint64, fileName []byte)
type cb_file_recv_chunk_ftype func(friendNumber uint32, fileNumber uint32, position uint64, data []byte)
type cb_file_chunk_request_ftype func(friend_number uint32, file_number uint32, position uint64, length int)

var cbUserDatas = newUserData()

//export callbackFriendRequestWrapperForC
func callbackFriendRequestWrapperForC(m *C.Tox, a0 *C.uint8_t, a1 *C.uint8_t, a2 C.uint16_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	var pubkey [PUBLIC_KEY_SIZE]byte
	copy(pubkey[:], (*[1 << 30]byte)(unsafe.Pointer(a0))[:])
	message := C.GoBytes(unsafe.Pointer(a1), C.int(a2))
	t.cb_friend_request(&pubkey, message)
}
func (t *Tox) CallbackFriendRequest(cbfn cb_friend_request_ftype) {
	t.cb_friend_request = cbfn
	C.tox_callback_friend_request(t.toxcore, (*C.tox_friend_request_cb)(C.callbackFriendRequestWrapperForC))
}

//export callbackFriendMessageWrapperForC
func callbackFriendMessageWrapperForC(m *C.Tox, a0 C.uint32_t, mtype C.int, a1 *C.uint8_t, a2 C.uint32_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	message := C.GoBytes(unsafe.Pointer(a1), (C.int)(a2))
	t.cb_friend_message(uint32(a0), toxenums.TOX_MESSAGE_TYPE(mtype), message)
}
func (t *Tox) CallbackFriendMessage(cbfn cb_friend_message_ftype) {
	t.cb_friend_message = cbfn
	C.tox_callback_friend_message(t.toxcore, (*C.tox_friend_message_cb)(C.callbackFriendMessageWrapperForC))
}

//export callbackFriendNameWrapperForC
func callbackFriendNameWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, a2 C.uint32_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	name := C.GoStringN((*C.char)((unsafe.Pointer)(a1)), C.int(a2))
	t.cb_friend_name(uint32(a0), name)
}
func (t *Tox) CallbackFriendName(cbfn cb_friend_name_ftype) {
	t.cb_friend_name = cbfn
	C.tox_callback_friend_name(t.toxcore, (*C.tox_friend_name_cb)(C.callbackFriendNameWrapperForC))
}

//export callbackFriendStatusMessageWrapperForC
func callbackFriendStatusMessageWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, a2 C.uint32_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	statusText := C.GoStringN((*C.char)(unsafe.Pointer(a1)), C.int(a2))
	t.cb_friend_status_message(uint32(a0), statusText)
}
func (t *Tox) CallbackFriendStatusMessage(cbfn cb_friend_status_message_ftype) {
	t.cb_friend_status_message = cbfn
	C.tox_callback_friend_status_message(t.toxcore, (*C.tox_friend_status_message_cb)(C.callbackFriendStatusMessageWrapperForC))
}

//export callbackFriendStatusWrapperForC
func callbackFriendStatusWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.TOX_USER_STATUS, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	t.cb_friend_status(uint32(a0), toxenums.TOX_USER_STATUS(a1))
}
func (t *Tox) CallbackFriendStatus(cbfn cb_friend_status_ftype) {
	t.cb_friend_status = cbfn
	C.tox_callback_friend_status(t.toxcore, (*C.tox_friend_status_cb)(C.callbackFriendStatusWrapperForC))
}

//export callbackFriendConnectionStatusWrapperForC
func callbackFriendConnectionStatusWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.TOX_CONNECTION, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	t.cb_friend_connection_status(uint32(a0), toxenums.TOX_CONNECTION(a1))
}
func (t *Tox) CallbackFriendConnectionStatus(cbfn cb_friend_connection_status_ftype) {
	t.cb_friend_connection_status = cbfn
	C.tox_callback_friend_connection_status(t.toxcore, (*C.tox_friend_connection_status_cb)(C.callbackFriendConnectionStatusWrapperForC))
}

//export callbackFriendTypingWrapperForC
func callbackFriendTypingWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint8_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	t.cb_friend_typing(uint32(a0), uint8(a1))
}
func (t *Tox) CallbackFriendTyping(cbfn cb_friend_typing_ftype) {
	t.cb_friend_typing = cbfn
	C.tox_callback_friend_typing(t.toxcore, (*C.tox_friend_typing_cb)(C.callbackFriendTypingWrapperForC))
}

//export callbackFriendReadReceiptWrapperForC
func callbackFriendReadReceiptWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint32_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	t.cb_friend_read_receipt(uint32(a0), uint32(a1))
}
func (t *Tox) CallbackFriendReadReceipt(cbfn cb_friend_read_receipt_ftype) {
	t.cb_friend_read_receipt = cbfn
	C.tox_callback_friend_read_receipt(t.toxcore, (*C.tox_friend_read_receipt_cb)(C.callbackFriendReadReceiptWrapperForC))
}

//export callbackFriendLossyPacketWrapperForC
func callbackFriendLossyPacketWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, len C.size_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	msg := C.GoBytes(unsafe.Pointer(a1), C.int(len))
	t.cb_friend_lossy_packet(uint32(a0), msg)
}
func (t *Tox) CallbackFriendLossyPacket(cbfn cb_friend_lossy_packet_ftype) {
	t.cb_friend_lossy_packet = cbfn
	C.tox_callback_friend_lossy_packet(t.toxcore, (*C.tox_friend_lossy_packet_cb)(C.callbackFriendLossyPacketWrapperForC))
}

//export callbackFriendLosslessPacketWrapperForC
func callbackFriendLosslessPacketWrapperForC(m *C.Tox, a0 C.uint32_t, a1 *C.uint8_t, len C.size_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	msg := C.GoBytes(unsafe.Pointer(a1), C.int(len))
	t.cb_friend_lossless_packet(uint32(a0), msg)
}
func (t *Tox) CallbackFriendLosslessPacket(cbfn cb_friend_lossless_packet_ftype) {
	t.cb_friend_lossless_packet = cbfn
	C.tox_callback_friend_lossless_packet(t.toxcore, (*C.tox_friend_lossless_packet_cb)(C.callbackFriendLosslessPacketWrapperForC))
}

//export callbackSelfConnectionStatusWrapperForC
func callbackSelfConnectionStatusWrapperForC(m *C.Tox, status C.TOX_CONNECTION, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	t.cb_self_connection_status(toxenums.TOX_CONNECTION(status))
}
func (t *Tox) CallbackSelfConnectionStatus(cbfn cb_self_connection_status_ftype) {
	t.cb_self_connection_status = cbfn
	C.tox_callback_self_connection_status(t.toxcore, (*C.tox_self_connection_status_cb)(C.callbackSelfConnectionStatusWrapperForC))
}

//export callbackFileRecvControlWrapperForC
func callbackFileRecvControlWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t, control C.TOX_FILE_CONTROL, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	t.cb_file_recv_control(uint32(friendNumber), uint32(fileNumber), toxenums.TOX_FILE_CONTROL(control))
}
func (t *Tox) CallbackFileRecvControl(cbfn cb_file_recv_control_ftype) {
	t.cb_file_recv_control = cbfn
	C.tox_callback_file_recv_control(t.toxcore, (*C.tox_file_recv_control_cb)(C.callbackFileRecvControlWrapperForC))
}

//export callbackFileRecvWrapperForC
func callbackFileRecvWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t, kind C.uint32_t, fileSize C.uint64_t, fileName *C.uint8_t, fileNameLength C.size_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	fileName_ := C.GoBytes(unsafe.Pointer(fileName), C.int(fileNameLength))
	t.cb_file_recv(uint32(friendNumber), uint32(fileNumber), toxenums.TOX_FILE_KIND(kind), uint64(fileSize), fileName_)
}
func (t *Tox) CallbackFileRecv(cbfn cb_file_recv_ftype) {
	t.cb_file_recv = cbfn
	C.tox_callback_file_recv(t.toxcore, (*C.tox_file_recv_cb)(C.callbackFileRecvWrapperForC))
}

//export callbackFileRecvChunkWrapperForC
func callbackFileRecvChunkWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t, position C.uint64_t, data *C.uint8_t, length C.size_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	data_ := C.GoBytes((unsafe.Pointer)(data), C.int(length))
	t.cb_file_recv_chunk(uint32(friendNumber), uint32(fileNumber), uint64(position), data_)
}
func (t *Tox) CallbackFileRecvChunk(cbfn cb_file_recv_chunk_ftype) {
	t.cb_file_recv_chunk = cbfn
	C.tox_callback_file_recv_chunk(t.toxcore, (*C.tox_file_recv_chunk_cb)(C.callbackFileRecvChunkWrapperForC))
}

//export callbackFileChunkRequestWrapperForC
func callbackFileChunkRequestWrapperForC(m *C.Tox, friendNumber C.uint32_t, fileNumber C.uint32_t, position C.uint64_t, length C.size_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	t.cb_file_chunk_request(uint32(friendNumber), uint32(fileNumber), uint64(position), int(length))
}
func (t *Tox) CallbackFileChunkRequest(cbfn cb_file_chunk_request_ftype) {
	t.cb_file_chunk_request = cbfn
	C.tox_callback_file_chunk_request(t.toxcore, (*C.tox_file_chunk_request_cb)(C.callbackFileChunkRequestWrapperForC))
}

// TODO return error when full
func (t *Tox) DoInLoop(fn func()) {
	t.chLoopRequest <- fn
}

func (t *Tox) SelfGetFriendList() map[uint32]*[PUBLIC_KEY_SIZE]byte {
	list := make(map[uint32]*[PUBLIC_KEY_SIZE]byte, 32)
	t.friendIdToPk.Range(func(key, value interface{}) bool {
		list[key.(uint32)] = value.(*[PUBLIC_KEY_SIZE]byte)
		return true
	})
	return list
}

func (t *Tox) FriendByPublicKey(pubkey *[PUBLIC_KEY_SIZE]byte) (uint32, bool) {
	r, ok := t.friendPkToId.Load(*pubkey)
	if !ok {
		return 0, false
	}
	return r.(uint32), true
}

func (t *Tox) FriendGetPublicKey(friendNumber uint32) (*[PUBLIC_KEY_SIZE]byte, bool) {
	r, ok := t.friendIdToPk.Load(friendNumber)
	if !ok {
		return nil, false
	}
	return r.(*[PUBLIC_KEY_SIZE]byte), true
}

func (t *Tox) FriendExists(friendNumber uint32) bool {
	_, ok := t.friendIdToPk.Load(friendNumber)
	return ok
}

func ToxHash(data []byte) []byte {
	_hash := make([]byte, C.TOX_HASH_LENGTH)
	var _datalen = C.size_t(len(data))
	C.tox_hash((*C.uint8_t)(&_hash[0]), (*C.uint8_t)(&data[0]), _datalen)
	return _hash
}

func (t *Tox) StopNotifier() <-chan struct{} { return t.stop }
func (t *Tox) CTox() *C.Tox                  { return t.toxcore }
