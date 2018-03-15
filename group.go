package tox

/*
#include <tox/tox.h>

void callbackConferenceInviteWrapperForC(Tox*, uint32_t, TOX_CONFERENCE_TYPE, uint8_t *, size_t, void *);
void callbackConferenceMessageWrapperForC(Tox *, uint32_t, uint32_t, TOX_MESSAGE_TYPE, uint8_t *, size_t, void *);
void callbackConferenceTitleWrapperForC(Tox*, uint32_t, uint32_t, uint8_t*, size_t, void*);
void callbackConferencePeerNameWrapperForC(Tox*, uint32_t, uint32_t, uint8_t*, size_t, void*);
void callbackConferencePeerListChangedWrapperForC(Tox*, uint32_t, void*);
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

// conference callback type
type cb_conference_invite_ftype func(friendNumber uint32, itype toxenums.TOX_CONFERENCE_TYPE, cookie []byte)
type cb_conference_message_ftype func(groupNumber uint32, peerNumber uint32, mtype toxenums.TOX_MESSAGE_TYPE, message []byte)
type cb_conference_title_ftype func(groupNumber uint32, peerNumber uint32, title string)
type cb_conference_peer_name_ftype func(groupNumber uint32, peerNumber uint32, name string)
type cb_conference_peer_list_changed_ftype func(groupNumber uint32)

// tox_callback_conference_***

//export callbackConferenceInviteWrapperForC
func callbackConferenceInviteWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.TOX_CONFERENCE_TYPE, a2 *C.uint8_t, a3 C.size_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	cookie := C.GoBytes((unsafe.Pointer)(a2), C.int(a3))
	t.cb_conference_invite(uint32(a0), toxenums.TOX_CONFERENCE_TYPE(a1), cookie)
}
func (t *Tox) CallbackConferenceInvite(cbfn cb_conference_invite_ftype) {
	t.cb_conference_invite = cbfn
	C.tox_callback_conference_invite(t.toxcore, (*C.tox_conference_invite_cb)(C.callbackConferenceInviteWrapperForC))
}

//export callbackConferenceMessageWrapperForC
func callbackConferenceMessageWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint32_t, mtype C.TOX_MESSAGE_TYPE, a2 *C.uint8_t, a3 C.size_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	message := C.GoBytes((unsafe.Pointer)(a2), C.int(a3))
	t.cb_conference_message(uint32(a0), uint32(a1), toxenums.TOX_MESSAGE_TYPE(mtype), message)
}
func (t *Tox) CallbackConferenceMessage(cbfn cb_conference_message_ftype) {
	t.cb_conference_message = cbfn
	C.tox_callback_conference_message(t.toxcore, (*C.tox_conference_message_cb)(C.callbackConferenceMessageWrapperForC))
}

//export callbackConferenceTitleWrapperForC
func callbackConferenceTitleWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint32_t, a2 *C.uint8_t, a3 C.size_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	title := C.GoStringN((*C.char)((unsafe.Pointer)(a2)), C.int(a3))
	t.cb_conference_title(uint32(a0), uint32(a1), title)
}
func (t *Tox) CallbackConferenceTitle(cbfn cb_conference_title_ftype) {
	t.cb_conference_title = cbfn
	C.tox_callback_conference_title(t.toxcore, (*C.tox_conference_title_cb)(C.callbackConferenceTitleWrapperForC))
}

//export callbackConferencePeerNameWrapperForC
func callbackConferencePeerNameWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint32_t, a2 *C.uint8_t, a3 C.size_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	peer_name := C.GoStringN((*C.char)((unsafe.Pointer)(a2)), C.int(a3))
	t.cb_conference_peer_name(uint32(a0), uint32(a1), peer_name)
}
func (t *Tox) CallbackConferencePeerName(cbfn cb_conference_peer_name_ftype) {
	t.cb_conference_peer_name = cbfn
	C.tox_callback_conference_peer_name(t.toxcore, (*C.tox_conference_peer_name_cb)(C.callbackConferencePeerNameWrapperForC))
}

//export callbackConferencePeerListChangedWrapperForC
func callbackConferencePeerListChangedWrapperForC(m *C.Tox, a0 C.uint32_t, ud unsafe.Pointer) {
	var t = cbUserDatas.get(m)
	t.cb_conference_peer_list_changed(uint32(a0))
}
func (t *Tox) CallbackConferencePeerListChanged(cbfn cb_conference_peer_list_changed_ftype) {
	t.cb_conference_peer_list_changed = cbfn
	C.tox_callback_conference_peer_list_changed(t.toxcore, (*C.tox_conference_peer_list_changed_cb)(C.callbackConferencePeerListChangedWrapperForC))
}

// methods tox_conference_*

func (t *Tox) conferenceNew_l(data ConferenceNewData) {
	var err error
	var cerr C.TOX_ERR_CONFERENCE_NEW
	r := C.tox_conference_new(t.toxcore, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_NEW(cerr)
	}

	data <- &ConferenceNewResult{
		ConferenceNumber: uint32(r),
		Error:            err,
	}
}

func (t *Tox) conferenceDelete_l(data *ConferenceDeleteData) {
	cn := C.uint32_t(data.ConferenceNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_DELETE
	C.tox_conference_delete(t.toxcore, cn, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_DELETE(cerr)
	}

	data.Result <- err
}

func (t *Tox) conferencePeerGetName_l(data *ConferencePeerGetNameData) {
	cn := C.uint32_t(data.ConferenceNumber)
	pn := C.uint32_t(data.PeerNumber)

	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	size := C.tox_conference_peer_get_name_size(t.toxcore, cn, pn, &cerr)
	if cerr != 0 {
		data.Result <- &ConferencePeerGetNameResult{Error: toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)}
		return
	}

	name := make([]byte, size)
	C.tox_conference_peer_get_name(t.toxcore, cn, pn, (*C.uint8_t)(&name[0]), &cerr)
	if cerr != 0 {
		data.Result <- &ConferencePeerGetNameResult{Error: toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)}
		return
	}

	data.Result <- &ConferencePeerGetNameResult{
		Name: string(name),
	}
}

func (t *Tox) conferencePeerGetPublicKey_l(data *ConferencePeerGetPublicKeyData) {
	cn := C.uint32_t(data.ConferenceNumber)
	pn := C.uint32_t(data.PeerNumber)

	var pubkey [PUBLIC_KEY_SIZE]byte
	var err error
	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	C.tox_conference_peer_get_public_key(t.toxcore, cn, pn, (*C.uint8_t)(&pubkey[0]), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}

	data.Result <- &ConferencePeerGetPublicKeyResult{
		Pubkey: &pubkey,
		Error:  err,
	}
}

func (t *Tox) conferenceInvite_l(data *ConferenceInviteData) {
	// if give a friendNumber which not exists,
	// the tox_invite_friend has a strange behaive: cause other tox_* call failed
	// and the call will return true, but only strange thing accurs
	// so just precheck the friendNumber and then go
	if !t.FriendExists(data.FriendNumber) {
		data.Result <- toxenums.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND
		return
	}

	cn := C.uint32_t(data.ConferenceNumber)
	fn := C.uint32_t(data.FriendNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_INVITE
	C.tox_conference_invite(t.toxcore, cn, fn, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_INVITE(cerr)
	}
	data.Result <- err
}

func (t *Tox) conferenceJoin_l(data *ConferenceJoinData) {
	if data.Cookie == nil {
		data.Result <- &ConferenceJoinResult{Error: toxenums.TOX_ERR_CONFERENCE_JOIN_INVALID_LENGTH}
		return
	}

	fn := C.uint32_t(data.FriendNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_JOIN
	r := C.tox_conference_join(t.toxcore, fn, (*C.uint8_t)(&data.Cookie[0]), C.size_t(len(data.Cookie)), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_JOIN(cerr)
	}

	data.Result <- &ConferenceJoinResult{
		ConferenceNumber: uint32(r),
		Error:            err,
	}
}

func (t *Tox) conferenceSendMessage_l(data *ConferenceSendMessageData) {
	switch data.Type {
	case toxenums.TOX_MESSAGE_TYPE_NORMAL:
	case toxenums.TOX_MESSAGE_TYPE_ACTION:
	default:
		data.Result <- fmt.Errorf("Invalid tox conference message type: %v", data.Type)
		return
	}

	cn := C.uint32_t(data.ConferenceNumber)
	message := (*C.uint8_t)(&data.Message[0])
	message_size := C.size_t(len(data.Message))

	var err error
	var cerr C.TOX_ERR_CONFERENCE_SEND_MESSAGE
	C.tox_conference_send_message(t.toxcore, cn, C.TOX_MESSAGE_TYPE(data.Type), message, message_size, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_SEND_MESSAGE(cerr)
	}

	data.Result <- err
}

func (t *Tox) conferenceSetTitle_l(data *ConferenceSetTitleData) {
	cn := C.uint32_t(data.ConferenceNumber)
	title := []byte(data.Title)
	title_size := C.size_t(len(title))

	var err error
	var cerr C.TOX_ERR_CONFERENCE_TITLE
	C.tox_conference_set_title(t.toxcore, cn, (*C.uint8_t)(&title[0]), title_size, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_TITLE(cerr)
	}

	data.Result <- err
}

func (t *Tox) conferenceGetTitle_l(data *ConferenceGetTitleData) {
	cn := C.uint32_t(data.ConferenceNumber)

	var cerr C.TOX_ERR_CONFERENCE_TITLE
	size := C.tox_conference_get_title_size(t.toxcore, cn, &cerr)
	if cerr != 0 {
		data.Result <- &ConferenceGetTitleResult{Error: toxenums.TOX_ERR_CONFERENCE_TITLE(cerr)}
		return
	}

	title := make([]byte, size)
	C.tox_conference_get_title(t.toxcore, cn, (*C.uint8_t)(&title[0]), &cerr)
	if cerr != 0 {
		data.Result <- &ConferenceGetTitleResult{Error: toxenums.TOX_ERR_CONFERENCE_TITLE(cerr)}
		return
	}

	data.Result <- &ConferenceGetTitleResult{Title: string(title)}
}

func (t *Tox) conferencePeerNumberIsOurs_l(data *ConferencePeerNumberIsOursData) {
	cn := C.uint32_t(data.ConferenceNumber)
	pn := C.uint32_t(data.PeerNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	r := C.tox_conference_peer_number_is_ours(t.toxcore, cn, pn, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}

	data.Result <- &ConferencePeerNumberIsOursResult{
		Is:    bool(r),
		Error: err,
	}
}

func (t *Tox) conferencePeerCount_l(data *ConferencePeerCountData) {
	cn := C.uint32_t(data.ConferenceNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	r := C.tox_conference_peer_count(t.toxcore, cn, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}

	data.Result <- &ConferencePeerCountResult{
		Count: uint32(r),
		Error: err,
	}
}

func (t *Tox) conferenceGetChatlist_l(data ConferenceGetChatlistData) {
	size := C.tox_conference_get_chatlist_size(t.toxcore)
	if size == 0 {
		data <- nil
		return
	}

	list := make([]uint32, size)
	C.tox_conference_get_chatlist(t.toxcore, (*C.uint32_t)(unsafe.Pointer(&list[0])))
	data <- list
}

func (t *Tox) conferenceGetType_l(data *ConferenceGetTypeData) {
	cn := C.uint32_t(data.ConferenceNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_GET_TYPE
	r := toxenums.TOX_CONFERENCE_TYPE(C.tox_conference_get_type(t.toxcore, cn, &cerr))
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_GET_TYPE(cerr)
	}

	data.Result <- &ConferenceGetTypeResult{
		Type:  toxenums.TOX_CONFERENCE_TYPE(r),
		Error: err,
	}
}
