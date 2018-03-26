package tox

//#include <tox/tox.h>
import "C"
import (
	"bytes"
	"encoding/hex"
	"log"
	"strings"
	"time"
	"unsafe"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

type (
	BootstrapNode struct {
		Addr    string
		Port    uint16
		TcpPort uint16
		Pubkey  [PUBLIC_KEY_SIZE]byte
	}
	BootstrapResult struct {
		Codes   []toxenums.TOX_ERR_BOOTSTRAP
		Success uint
		LastErr toxenums.TOX_ERR_BOOTSTRAP
	}
)

// Error return nil if at least one succeed.
func (r *BootstrapResult) Error() error {
	if r.Success == 0 {
		return r.LastErr
	}
	return nil
}

func (t *Tox) onFriendAdded_l(friendNumber uint32, pubkey *[PUBLIC_KEY_SIZE]byte) {
	t.friendIdToPk.Store(friendNumber, pubkey)
	t.friendPkToId.Store(*pubkey, friendNumber)

	friendBig := bytes.Compare(pubkey[:], t.Pubkey[:]) == 1
	var dialSeq byte = 1
	if friendBig {
		dialSeq = 0
	}
	pingValue := defaultPinMapValue
	tf := ToxFriend{
		FriendNumber: friendNumber,
		Pubkey:       *pubkey,
		FriendBig:    friendBig,

		tox:        t,
		dialSeq:    dialSeq,
		remoteAddr: addr(strings.ToUpper(hex.EncodeToString(pubkey[:]))),
		ping:       &pingValue,
	}
	t.friends[friendNumber] = &tf
}

func (t *Tox) onFriendDeleted_l(friendNumber uint32) {
	tf, ok := t.friends[friendNumber]
	if ok {
		tf.CloseStreams_l()
		delete(t.friends, friendNumber)
	}

	pubkey, ok := t.friendPkToId.Load(friendNumber)
	if ok {
		t.friendPkToId.Delete(*(pubkey.(*[PUBLIC_KEY_SIZE]byte)))
	}
	t.friendIdToPk.Delete(friendNumber)
}

func (t *Tox) BootstrapNode_l(node *BootstrapNode) toxenums.TOX_ERR_BOOTSTRAP {
	addrb := []byte(node.Addr)
	addr := (*C.char)(unsafe.Pointer(&addrb[0]))
	port := C.uint16_t(node.Port)
	tcp := C.uint16_t(node.TcpPort)
	cpubkey := (*C.uint8_t)(&node.Pubkey[0])

	var cerr C.TOX_ERR_BOOTSTRAP
	C.tox_bootstrap(t.toxcore, addr, port, cpubkey, &cerr)
	if cerr == 0 && tcp != 0 {
		C.tox_add_tcp_relay(t.toxcore, addr, tcp, cpubkey, &cerr)
	}

	return toxenums.TOX_ERR_BOOTSTRAP(cerr)
}

func (t *Tox) BootstrapNodes_l(nodes []BootstrapNode) *BootstrapResult {
	result := BootstrapResult{
		Codes: make([]toxenums.TOX_ERR_BOOTSTRAP, len(nodes)),
	}
	for i, node := range nodes {
		err := t.BootstrapNode_l(&node)
		result.Codes[i] = err
		if err == 0 {
			result.Success++
		} else {
			result.LastErr = err
		}
	}
	return &result
}

func (t *Tox) GetSavedata_l() []byte {
	savedata := make([]byte, C.tox_get_savedata_size(t.toxcore))
	C.tox_get_savedata(t.toxcore, (*C.uint8_t)(&savedata[0]))
	return savedata
}

func (t *Tox) SelfSetNospam_l(nospam uint32) {
	C.tox_self_set_nospam(t.toxcore, C.uint32_t(nospam))
}

func (t *Tox) SelfGetNospam_l() uint32 {
	return uint32(C.tox_self_get_nospam(t.toxcore))
}

func (t *Tox) SelfSetName_l(n string) error {
	name := []byte(n)
	name_size := C.size_t(len(name))

	var err error
	var cerr C.TOX_ERR_SET_INFO
	C.tox_self_set_name(t.toxcore, (*C.uint8_t)(&name[0]), name_size, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_SET_INFO(cerr)
	}
	return err
}

func (t *Tox) SelfGetName_l() string {
	name_size := C.tox_self_get_name_size(t.toxcore)
	if name_size == 0 {
		return ""
	}

	name := make([]byte, name_size)
	C.tox_self_get_name(t.toxcore, (*C.uint8_t)(&name[0]))
	return string(name)
}

func (t *Tox) SelfSetStatus_l(status toxenums.TOX_USER_STATUS) {
	C.tox_self_set_status(t.toxcore, C.TOX_USER_STATUS(status))
}

func (t *Tox) SelfGetStatus_l() toxenums.TOX_USER_STATUS {
	return toxenums.TOX_USER_STATUS(C.tox_self_get_status(t.toxcore))
}

func (t *Tox) SelfSetStatusMessage_l(message string) error {
	status := []byte(message)
	status_size := C.size_t(len(status))

	var err error
	var cerr C.TOX_ERR_SET_INFO
	C.tox_self_set_status_message(t.toxcore, (*C.uint8_t)(&status[0]), status_size, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_SET_INFO(cerr)
	}
	return err
}

func (t *Tox) SelfGetStatusMessage_l() string {
	message_size := C.tox_self_get_status_message_size(t.toxcore)
	if message_size == 0 {
		return ""
	}

	message := make([]byte, message_size)
	C.tox_self_get_status_message(t.toxcore, (*C.uint8_t)(&message[0]))
	return string(message)
}

func (t *Tox) FriendAdd_l(address *[ADDRESS_SIZE]byte, message []byte) (uint32, error) {
	caddress := (*C.uint8_t)(&address[0])
	cmessage := (*C.uint8_t)(&message[0])
	message_size := C.size_t(len(message))

	t.blockAv()
	defer t.unblockAv()

	var err error
	var cerr C.TOX_ERR_FRIEND_ADD
	friendNumber := uint32(C.tox_friend_add(t.toxcore, caddress, cmessage, message_size, &cerr))
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_ADD(cerr)
	} else {
		var pubkey [PUBLIC_KEY_SIZE]byte
		copy(pubkey[:], address[:])
		t.onFriendAdded_l(friendNumber, &pubkey)
	}

	return friendNumber, err
}

func (t *Tox) FriendAddNorequest_l(pubkey *[PUBLIC_KEY_SIZE]byte) (uint32, error) {
	cpubkey := (*C.uint8_t)(&pubkey[0])

	t.blockAv()
	defer t.unblockAv()

	var err error
	var cerr C.TOX_ERR_FRIEND_ADD
	friendNumber := uint32(C.tox_friend_add_norequest(t.toxcore, cpubkey, &cerr))
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_ADD(cerr)
	} else {
		pk := *pubkey
		t.onFriendAdded_l(friendNumber, &pk)
	}
	return friendNumber, err
}

func (t *Tox) FriendDelete_l(friendNumber uint32) error {
	fn := C.uint32_t(friendNumber)

	t.blockAv()
	defer t.unblockAv()

	var err error
	var cerr C.TOX_ERR_FRIEND_DELETE
	C.tox_friend_delete(t.toxcore, fn, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_DELETE(cerr)
	} else {
		t.onFriendDeleted_l(friendNumber)
	}
	return err
}

func (t *Tox) FriendSendLosslessPacket_l(friendNumber uint32, data []byte, noRetry bool) toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET {
	fn := C.uint32_t(friendNumber)
	cdata := (*C.uint8_t)(&data[0])
	cdata_size := C.size_t(len(data))

	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET
	var e toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET

	// import from tuntox send_frame
	i := time.Duration(1)
	j := time.Duration(0)
	try := 0
	for i < 65 { // 65->2667ms 33->651ms per packet max 17->155ms
		try++

		C.tox_friend_send_lossless_packet(t.toxcore, fn, cdata, cdata_size, &cerr)
		e = toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET(cerr)
		switch e {
		case toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET_OK:
			goto end
		case toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET_SENDQ:
		case toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET_FRIEND_NOT_CONNECTED:
			log.Printf("[%d] Failed to send packet to friend %d (Friend gone)\n", i, friendNumber)
			goto end
		default:
			log.Printf("[%d] Failed to send packet to friend %d (err: %v)\n", i, friendNumber, e)
		}

		if t.inToxIterate || noRetry {
			goto end
		}

		i = i << 1
		for j = 0; j < i; j++ {
			C.tox_iterate(t.toxcore, nil)
			time.Sleep(j * time.Millisecond)
		}
	}

end:
	if e == 0 && try > 1 {
		log.Printf("Packet succeeded at try %d (friend %d)\n", try, friendNumber)
	}
	return e
}

func (t *Tox) FriendSendLossyPacket_l(friendNumber uint32, data []byte) error {
	fn := C.uint32_t(friendNumber)
	data_size := C.size_t(len(data))

	var err error
	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET
	C.tox_friend_send_lossy_packet(t.toxcore, fn, (*C.uint8_t)(&data[0]), data_size, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET(cerr)
	}
	return err
}

func (t *Tox) FriendSendMessage_l(friendNumber uint32, mtype toxenums.TOX_MESSAGE_TYPE, message []byte) (uint32, error) {
	fn := C.uint32_t(friendNumber)
	cmtype := C.TOX_MESSAGE_TYPE(mtype)
	cmessage := (*C.uint8_t)(&message[0])
	cmessage_size := C.size_t(len(message))

	var err error
	var cerr C.TOX_ERR_FRIEND_SEND_MESSAGE
	messageId := uint32(C.tox_friend_send_message(t.toxcore, fn, cmtype, cmessage, cmessage_size, &cerr))
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_SEND_MESSAGE(cerr)
	}
	return messageId, err
}

func (t *Tox) FriendGetName_l(friendNumber uint32) (string, error) {
	fn := C.uint32_t(friendNumber)

	var cerr C.TOX_ERR_FRIEND_QUERY
	name_size := C.tox_friend_get_name_size(t.toxcore, fn, &cerr)
	if cerr != 0 {
		return "", toxenums.TOX_ERR_FRIEND_QUERY(cerr)
	}

	if name_size == 0 {
		return "", nil
	}

	name := make([]byte, name_size)
	C.tox_friend_get_name(t.toxcore, fn, (*C.uint8_t)(&name[0]), &cerr)
	if cerr != 0 {
		return "", toxenums.TOX_ERR_FRIEND_QUERY(cerr)
	}
	return string(name), nil
}

func (t *Tox) FriendGetStatusMessage_l(friendNumber uint32) (string, error) {
	fn := C.uint32_t(friendNumber)

	var cerr C.TOX_ERR_FRIEND_QUERY
	message_size := C.tox_friend_get_status_message_size(t.toxcore, fn, &cerr)
	if cerr != 0 {
		return "", toxenums.TOX_ERR_FRIEND_QUERY(cerr)
	}

	if message_size == 0 {
		return "", nil
	}

	message := make([]byte, message_size)
	C.tox_friend_get_status_message(t.toxcore, fn, (*C.uint8_t)(&message[0]), &cerr)
	if cerr != 0 {
		return "", toxenums.TOX_ERR_FRIEND_QUERY(cerr)
	}
	return string(message), nil
}

func (t *Tox) FriendGetStatus_l(friendNumber uint32) (toxenums.TOX_USER_STATUS, error) {
	fn := C.uint32_t(friendNumber)

	var err error
	var cerr C.TOX_ERR_FRIEND_QUERY
	status := toxenums.TOX_USER_STATUS(C.tox_friend_get_status(t.toxcore, fn, &cerr))
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_QUERY(cerr)
	}
	return status, err
}

func (t *Tox) FriendGetLastOnline_l(friendNumber uint32) (uint64, error) {
	fn := C.uint32_t(friendNumber)

	var err error
	var cerr C.TOX_ERR_FRIEND_GET_LAST_ONLINE
	unixTime := uint64(C.tox_friend_get_last_online(t.toxcore, fn, &cerr))
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_GET_LAST_ONLINE(cerr)
	}
	return unixTime, err
}

func (t *Tox) SelfSetTyping_l(friendNumber uint32, typing bool) error {
	fn := C.uint32_t(friendNumber)
	ctyping := C._Bool(typing)

	var err error
	var cerr C.TOX_ERR_SET_TYPING
	C.tox_self_set_typing(t.toxcore, fn, ctyping, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_SET_TYPING(cerr)
	}
	return err
}

// file send

func (t *Tox) FileControl_l(friendNumber, fileNumber uint32, control toxenums.TOX_FILE_CONTROL) error {
	fn := C.uint32_t(friendNumber)
	file_number := C.uint32_t(fileNumber)

	var err error
	var cerr C.TOX_ERR_FILE_CONTROL
	C.tox_file_control(t.toxcore, fn, file_number, C.TOX_FILE_CONTROL(control), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FILE_CONTROL(cerr)
	}
	return err
}

func (t *Tox) FileSend_l(friendNumber uint32, kind toxenums.TOX_FILE_KIND, fileSize uint64, fileId *[FILE_ID_LENGTH]byte, fileName []byte) (uint32, error) {
	fn := C.uint32_t(friendNumber)

	var file_id *C.uint8_t
	if fileId != nil {
		file_id = (*C.uint8_t)(&fileId[0])
	}

	var err error
	var cerr C.TOX_ERR_FILE_SEND
	r := C.tox_file_send(t.toxcore, fn, C.uint32_t(kind), C.uint64_t(fileSize), file_id, (*C.uint8_t)(&fileName[0]), C.size_t(len(fileName)), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FILE_SEND(cerr)
	}
	return uint32(r), err
}

func (t *Tox) FileSendChunk_l(friendNumber uint32, fileNumber uint32, position uint64, data []byte) error {
	fn := C.uint32_t(friendNumber)
	file_number := C.uint32_t(fileNumber)

	var err error
	var cerr C.TOX_ERR_FILE_SEND_CHUNK
	C.tox_file_send_chunk(t.toxcore, fn, file_number, C.uint64_t(position), (*C.uint8_t)(&data[0]), C.size_t(len(data)), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FILE_SEND_CHUNK(cerr)
	}
	return err
}

func (t *Tox) FileSeek_l(friendNumber uint32, fileNumber uint32, position uint64) error {
	fn := C.uint32_t(friendNumber)
	file_number := C.uint32_t(fileNumber)

	var err error
	var cerr C.TOX_ERR_FILE_SEEK
	C.tox_file_seek(t.toxcore, fn, file_number, C.uint64_t(position), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FILE_SEEK(cerr)
	}
	return err
}

func (t *Tox) FileGetFileId_l(friendNumber uint32, fileNumber uint32) (*[FILE_ID_LENGTH]byte, error) {
	var fileId [FILE_ID_LENGTH]byte
	fn := C.uint32_t(friendNumber)
	file_number := C.uint32_t(fileNumber)

	var err error
	var cerr C.TOX_ERR_FILE_GET
	C.tox_file_get_file_id(t.toxcore, fn, file_number, (*C.uint8_t)(&fileId[0]), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FILE_GET(cerr)
	}
	return &fileId, err
}
