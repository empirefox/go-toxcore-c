package tox

//#include <tox/tox.h>
import "C"
import (
	"bytes"
	"encoding/hex"
	"strings"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

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

// send packet

func (t *Tox) friendSendLossyPacket_l(data *FriendSendLossyPacketData) {
	fn := C.uint32_t(data.FriendNumber)
	data_size := C.size_t(len(data.Data))

	var err error
	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET
	C.tox_friend_send_lossy_packet(t.toxcore, fn, (*C.uint8_t)(&data.Data[0]), data_size, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET(cerr)
	}

	data.Result <- err
}

func (t *Tox) friendSendLosslessPacket_l(data *FriendSendLosslessPacketData) {
	fn := C.uint32_t(data.FriendNumber)
	data_size := C.size_t(len(data.Data))

	var err error
	var cerr C.TOX_ERR_FRIEND_CUSTOM_PACKET
	C.tox_friend_send_lossless_packet(t.toxcore, fn, (*C.uint8_t)(&data.Data[0]), data_size, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_CUSTOM_PACKET(cerr)
	}

	data.Result <- err
}

func (t *Tox) friendSendMessage_l(data *FriendSendMessageData) {
	messageId, err := t.FriendSendMessage_l(data.FriendNumber, data.Type, data.Message)
	data.Result <- &FriendSendMessageResult{
		MessageId: messageId,
		Error:     err,
	}
}

func (t *Tox) bootstrapNodes_l(data *BootstrapData) {
	data.Result <- t.BootstrapNodes_l(data.Nodes)
}

func (t *Tox) getSavedata_l(data GetSavedataData) {
	data <- t.GetSavedata_l()
}

func (t *Tox) friendAdd_l(data *FriendAddData) {
	friendNumber, err := t.FriendAdd_l(data.Address, data.Message)
	data.Result <- &FriendAddResult{
		FriendNumber: friendNumber,
		Error:        err,
	}
}

func (t *Tox) friendAddNorequest_l(data *FriendAddNorequestData) {
	friendNumber, err := t.FriendAddNorequest_l(data.Pubkey)
	data.Result <- &FriendAddResult{
		FriendNumber: friendNumber,
		Error:        err,
	}
}

func (t *Tox) friendDelete_l(data *FriendDeleteData) {
	data.Result <- t.FriendDelete_l(data.FriendNumber)
}

func (t *Tox) selfSetName_l(data *SelfSetNameData) {
	data.Result <- t.SelfSetName_l(data.Name)
}

func (t *Tox) selfGetName_l(data SelfGetNameData) {
	data <- t.SelfGetName_l()
}

func (t *Tox) selfSetStatusMessage_l(data *SelfSetStatusMessageData) {
	data.Result <- t.SelfSetStatusMessage_l(data.Message)
}

func (t *Tox) selfGetStatusMessage_l(data SelfGetStatusMessageData) {
	data <- t.SelfGetStatusMessage_l()
}

func (t *Tox) selfSetStatus_l(data SelfSetStatusData) {
	C.tox_self_set_status(t.toxcore, C.TOX_USER_STATUS(data))
}

func (t *Tox) selfGetStatus_l(data SelfGetStatusData) {
	data <- toxenums.TOX_USER_STATUS(C.tox_self_get_status(t.toxcore))
}

func (t *Tox) selfSetNospam_l(data SelfSetNospamData) {
	C.tox_self_set_nospam(t.toxcore, C.uint32_t(data))
}

func (t *Tox) selfGetNospam_l(data SelfGetNospamData) {
	data <- uint32(C.tox_self_get_nospam(t.toxcore))
}

func (t *Tox) friendGetName_l(data *FriendGetNameData) {
	fn := C.uint32_t(data.FriendNumber)

	var cerr C.TOX_ERR_FRIEND_QUERY
	name_size := C.tox_friend_get_name_size(t.toxcore, fn, &cerr)
	if cerr != 0 {
		data.Result <- &FriendGetNameResult{Error: toxenums.TOX_ERR_FRIEND_QUERY(cerr)}
		return
	}

	if name_size == 0 {
		data.Result <- new(FriendGetNameResult)
		return
	}

	name := make([]byte, name_size)
	C.tox_friend_get_name(t.toxcore, fn, (*C.uint8_t)(&name[0]), &cerr)
	if cerr != 0 {
		data.Result <- &FriendGetNameResult{Error: toxenums.TOX_ERR_FRIEND_QUERY(cerr)}
		return
	}

	data.Result <- &FriendGetNameResult{Name: string(name)}
}

func (t *Tox) friendGetStatusMessage_l(data *FriendGetStatusMessageData) {
	fn := C.uint32_t(data.FriendNumber)

	var cerr C.TOX_ERR_FRIEND_QUERY
	message_size := C.tox_friend_get_status_message_size(t.toxcore, fn, &cerr)
	if cerr != 0 {
		data.Result <- &FriendGetStatusMessageResult{Error: toxenums.TOX_ERR_FRIEND_QUERY(cerr)}
		return
	}

	if message_size == 0 {
		data.Result <- new(FriendGetStatusMessageResult)
		return
	}

	message := make([]byte, message_size)
	C.tox_friend_get_status_message(t.toxcore, fn, (*C.uint8_t)(&message[0]), &cerr)
	if cerr != 0 {
		data.Result <- &FriendGetStatusMessageResult{Error: toxenums.TOX_ERR_FRIEND_QUERY(cerr)}
		return
	}

	data.Result <- &FriendGetStatusMessageResult{Message: string(message)}
}

func (t *Tox) friendGetStatus_l(data *FriendGetStatusData) {
	fn := C.uint32_t(data.FriendNumber)

	var err error
	var cerr C.TOX_ERR_FRIEND_QUERY
	status := toxenums.TOX_USER_STATUS(C.tox_friend_get_status(t.toxcore, fn, &cerr))
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_QUERY(cerr)
	}

	data.Result <- &FriendGetStatusResult{
		Status: status,
		Error:  err,
	}
}

func (t *Tox) friendGetLastOnline_l(data *FriendGetLastOnlineData) {
	fn := C.uint32_t(data.FriendNumber)

	var err error
	var cerr C.TOX_ERR_FRIEND_GET_LAST_ONLINE
	unixTime := uint64(C.tox_friend_get_last_online(t.toxcore, fn, &cerr))
	if cerr != 0 {
		err = toxenums.TOX_ERR_FRIEND_GET_LAST_ONLINE(cerr)
	}

	data.Result <- &FriendGetLastOnlineResult{
		Unix:  unixTime,
		Error: err,
	}
}

func (t *Tox) selfSetTyping_l(data *SelfSetTypingData) {
	fn := C.uint32_t(data.FriendNumber)
	typing := C._Bool(data.Typing)

	var err error
	var cerr C.TOX_ERR_SET_TYPING
	C.tox_self_set_typing(t.toxcore, fn, typing, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_SET_TYPING(cerr)
	}

	data.Result <- err
}

// file send

func (t *Tox) fileControl_l(data *FileControlData) {
	fn := C.uint32_t(data.FriendNumber)
	fileNumber := C.uint32_t(data.FileNumber)

	var err error
	var cerr C.TOX_ERR_FILE_CONTROL
	C.tox_file_control(t.toxcore, fn, fileNumber, C.TOX_FILE_CONTROL(data.Control), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FILE_CONTROL(cerr)
	}

	data.Result <- err
}

func (t *Tox) fileSend_l(data *FileSendData) {
	fn := C.uint32_t(data.FriendNumber)

	var fileId *C.uint8_t
	if data.FileId != nil {
		fileId = (*C.uint8_t)(&data.FileId[0])
	}

	var err error
	var cerr C.TOX_ERR_FILE_SEND
	r := C.tox_file_send(t.toxcore, fn, C.uint32_t(data.Kind), C.uint64_t(data.FileSize),
		fileId, (*C.uint8_t)(&data.FileName[0]), C.size_t(len(data.FileName)), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FILE_SEND(cerr)
	}

	data.Result <- &FileSendResult{
		FileNumber: uint32(r),
		Error:      err,
	}
}

func (t *Tox) fileSendChunk_l(data *FileSendChunkData) {
	fn := C.uint32_t(data.FriendNumber)
	fileNumber := C.uint32_t(data.FileNumber)

	var err error
	var cerr C.TOX_ERR_FILE_SEND_CHUNK
	C.tox_file_send_chunk(t.toxcore, fn, fileNumber, C.uint64_t(data.Position),
		(*C.uint8_t)(&data.Data[0]), C.size_t(len(data.Data)), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FILE_SEND_CHUNK(cerr)
	}

	data.Result <- err
}

func (t *Tox) fileSeek_l(data *FileSeekData) {
	fn := C.uint32_t(data.FriendNumber)
	fileNumber := C.uint32_t(data.FileNumber)

	var err error
	var cerr C.TOX_ERR_FILE_SEEK
	C.tox_file_seek(t.toxcore, fn, fileNumber, C.uint64_t(data.Position), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FILE_SEEK(cerr)
	}

	data.Result <- err
}

func (t *Tox) fileGetFileId_l(data *FileGetFileIdData) {
	var fileId [FILE_ID_LENGTH]byte
	fn := C.uint32_t(data.FriendNumber)
	fileNumber := C.uint32_t(data.FileNumber)

	var err error
	var cerr C.TOX_ERR_FILE_GET
	C.tox_file_get_file_id(t.toxcore, fn, fileNumber, (*C.uint8_t)(&fileId[0]), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_FILE_GET(cerr)
	}

	data.Result <- &FileGetFileIdResult{
		FileId: &fileId,
		Error:  err,
	}
}
