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

func (t *Tox) ConferenceNew_l() (conferenceNumber uint32, err error) {
	var cerr C.TOX_ERR_CONFERENCE_NEW
	r := C.tox_conference_new(t.toxcore, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_NEW(cerr)
	}
	return uint32(r), err
}

func (t *Tox) ConferenceDelete_l(conferenceNumber uint32) error {
	cn := C.uint32_t(conferenceNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_DELETE
	C.tox_conference_delete(t.toxcore, cn, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_DELETE(cerr)
	}
	return err
}

func (t *Tox) ConferencePeerGetName_l(conferenceNumber, peerNumber uint32) (string, error) {
	cn := C.uint32_t(conferenceNumber)
	pn := C.uint32_t(peerNumber)

	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	size := C.tox_conference_peer_get_name_size(t.toxcore, cn, pn, &cerr)
	if cerr != 0 {
		return "", toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}

	name := make([]byte, size)
	C.tox_conference_peer_get_name(t.toxcore, cn, pn, (*C.uint8_t)(&name[0]), &cerr)
	if cerr != 0 {
		return "", toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}

	return string(name), nil
}

func (t *Tox) ConferencePeerGetPublicKey_l(conferenceNumber, peerNumber uint32) (*[PUBLIC_KEY_SIZE]byte, error) {
	cn := C.uint32_t(conferenceNumber)
	pn := C.uint32_t(peerNumber)

	var pubkey [PUBLIC_KEY_SIZE]byte
	var err error
	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	C.tox_conference_peer_get_public_key(t.toxcore, cn, pn, (*C.uint8_t)(&pubkey[0]), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}
	return &pubkey, err
}

func (t *Tox) ConferenceInvite_l(conferenceNumber, friendNumber uint32) error {
	// if give a friendNumber which not exists,
	// the tox_invite_friend has a strange behaive: cause other tox_* call failed
	// and the call will return true, but only strange thing accurs
	// so just precheck the friendNumber and then go
	if !t.FriendExists(friendNumber) {
		return toxenums.TOX_ERR_FRIEND_QUERY_FRIEND_NOT_FOUND
	}

	cn := C.uint32_t(conferenceNumber)
	fn := C.uint32_t(friendNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_INVITE
	C.tox_conference_invite(t.toxcore, cn, fn, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_INVITE(cerr)
	}
	return err
}

func (t *Tox) ConferenceJoin_l(friendNumber uint32, cookie []byte) (uint32, error) {
	if cookie == nil {
		return 0, toxenums.TOX_ERR_CONFERENCE_JOIN_INVALID_LENGTH
	}

	fn := C.uint32_t(friendNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_JOIN
	r := C.tox_conference_join(t.toxcore, fn, (*C.uint8_t)(&cookie[0]), C.size_t(len(cookie)), &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_JOIN(cerr)
	}
	return uint32(r), err
}

func (t *Tox) ConferenceSendMessage_l(conferenceNumber uint32, typ toxenums.TOX_MESSAGE_TYPE, message []byte) error {
	switch typ {
	case toxenums.TOX_MESSAGE_TYPE_NORMAL:
	case toxenums.TOX_MESSAGE_TYPE_ACTION:
	default:
		return fmt.Errorf("Invalid tox conference message type: %v", typ)
	}

	cn := C.uint32_t(conferenceNumber)
	cmessage := (*C.uint8_t)(&message[0])
	cmessage_size := C.size_t(len(message))

	var err error
	var cerr C.TOX_ERR_CONFERENCE_SEND_MESSAGE
	C.tox_conference_send_message(t.toxcore, cn, C.TOX_MESSAGE_TYPE(typ), cmessage, cmessage_size, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_SEND_MESSAGE(cerr)
	}
	return err
}

func (t *Tox) ConferenceSetTitle_l(conferenceNumber uint32, title string) error {
	cn := C.uint32_t(conferenceNumber)
	ctitle := []byte(title)
	ctitle_size := C.size_t(len(ctitle))

	var err error
	var cerr C.TOX_ERR_CONFERENCE_TITLE
	C.tox_conference_set_title(t.toxcore, cn, (*C.uint8_t)(&ctitle[0]), ctitle_size, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_TITLE(cerr)
	}
	return err
}

func (t *Tox) ConferenceGetTitle_l(conferenceNumber uint32) (string, error) {
	cn := C.uint32_t(conferenceNumber)

	var cerr C.TOX_ERR_CONFERENCE_TITLE
	size := C.tox_conference_get_title_size(t.toxcore, cn, &cerr)
	if cerr != 0 {
		return "", toxenums.TOX_ERR_CONFERENCE_TITLE(cerr)
	}

	title := make([]byte, size)
	C.tox_conference_get_title(t.toxcore, cn, (*C.uint8_t)(&title[0]), &cerr)
	if cerr != 0 {
		return "", toxenums.TOX_ERR_CONFERENCE_TITLE(cerr)
	}

	return string(title), nil
}

func (t *Tox) ConferencePeerNumberIsOurs_l(conferenceNumber, peerNumber uint32) (bool, error) {
	cn := C.uint32_t(conferenceNumber)
	pn := C.uint32_t(peerNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	r := C.tox_conference_peer_number_is_ours(t.toxcore, cn, pn, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}
	return bool(r), err
}

func (t *Tox) ConferencePeerCount_l(conferenceNumber uint32) (uint32, error) {
	cn := C.uint32_t(conferenceNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	r := C.tox_conference_peer_count(t.toxcore, cn, &cerr)
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}
	return uint32(r), err
}

func (t *Tox) ConferenceGetChatlist_l() []uint32 {
	size := C.tox_conference_get_chatlist_size(t.toxcore)
	if size == 0 {
		return nil
	}

	list := make([]uint32, size)
	C.tox_conference_get_chatlist(t.toxcore, (*C.uint32_t)(unsafe.Pointer(&list[0])))
	return list
}

func (t *Tox) ConferenceGetType_l(conferenceNumber uint32) (toxenums.TOX_CONFERENCE_TYPE, error) {
	cn := C.uint32_t(conferenceNumber)

	var err error
	var cerr C.TOX_ERR_CONFERENCE_GET_TYPE
	r := toxenums.TOX_CONFERENCE_TYPE(C.tox_conference_get_type(t.toxcore, cn, &cerr))
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_GET_TYPE(cerr)
	}
	return toxenums.TOX_CONFERENCE_TYPE(r), err
}
