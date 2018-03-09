package tox

/*
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <tox/tox.h>

void callbackConferenceInviteWrapperForC(Tox*, uint32_t, TOX_CONFERENCE_TYPE, uint8_t *, size_t, void *);
void callbackConferenceMessageWrapperForC(Tox *, uint32_t, uint32_t, TOX_MESSAGE_TYPE, uint8_t *, size_t, void *);
// void callbackConferenceActionWrapperForC(Tox*, uint32_t, uint32_t, uint8_t*, size_t, void*);

void callbackConferenceTitleWrapperForC(Tox*, uint32_t, uint32_t, uint8_t*, size_t, void*);
void callbackConferencePeerNameWrapperForC(Tox*, uint32_t, uint32_t, uint8_t*, size_t, void*);
void callbackConferencePeerListChangedWrapperForC(Tox*, uint32_t, void*);

// fix nouse compile warning
static inline __attribute__((__unused__)) void fixnousetoxgroup() {
}

*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

// conference callback type
type cb_conference_invite_ftype func(this *Tox, friendNumber uint32, itype toxenums.TOX_CONFERENCE_TYPE, cookie []byte, userData interface{})
type cb_conference_message_ftype func(this *Tox, groupNumber uint32, peerNumber uint32, message []byte, userData interface{})

type cb_conference_action_ftype func(this *Tox, groupNumber uint32, peerNumber uint32, action []byte, userData interface{})
type cb_conference_title_ftype func(this *Tox, groupNumber uint32, peerNumber uint32, title string, userData interface{})
type cb_conference_peer_name_ftype func(this *Tox, groupNumber uint32, peerNumber uint32, name string, userData interface{})
type cb_conference_peer_list_changed_ftype func(this *Tox, groupNumber uint32, userData interface{})

// tox_callback_conference_***

//export callbackConferenceInviteWrapperForC
func callbackConferenceInviteWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.TOX_CONFERENCE_TYPE, a2 *C.uint8_t, a3 C.size_t, a4 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_conference_invites {
		cbfn := *(*cb_conference_invite_ftype)(cbfni)
		cookie := C.GoBytes((unsafe.Pointer)(a2), C.int(a3))
		this.putcbevts(func() { cbfn(this, uint32(a0), toxenums.TOX_CONFERENCE_TYPE(a1), cookie, ud) })
	}
}

func (this *Tox) CallbackConferenceInvite(cbfn cb_conference_invite_ftype, userData interface{}) {
	this.CallbackConferenceInviteAdd(cbfn, userData)
}
func (this *Tox) CallbackConferenceInviteAdd(cbfn cb_conference_invite_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_conference_invites[cbfnp]; ok {
		return
	}
	this.cb_conference_invites[cbfnp] = userData

	C.tox_callback_conference_invite(this.toxcore, (*C.tox_conference_invite_cb)(C.callbackConferenceInviteWrapperForC))
}

//export callbackConferenceMessageWrapperForC
func callbackConferenceMessageWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint32_t, mtype C.TOX_MESSAGE_TYPE, a2 *C.uint8_t, a3 C.size_t, a4 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	if toxenums.TOX_MESSAGE_TYPE(mtype) == toxenums.TOX_MESSAGE_TYPE_NORMAL {
		for cbfni, ud := range this.cb_conference_messages {
			cbfn := *(*cb_conference_message_ftype)(cbfni)
			message := C.GoBytes((unsafe.Pointer)(a2), C.int(a3))
			this.putcbevts(func() { cbfn(this, uint32(a0), uint32(a1), message, ud) })
		}
	} else {
		for cbfni, ud := range this.cb_conference_actions {
			cbfn := *(*cb_conference_action_ftype)(cbfni)
			message := C.GoBytes((unsafe.Pointer)(a2), C.int(a3))
			this.putcbevts(func() { cbfn(this, uint32(a0), uint32(a1), message, ud) })
		}
	}
}

func (this *Tox) CallbackConferenceMessage(cbfn cb_conference_message_ftype, userData interface{}) {
	this.CallbackConferenceMessageAdd(cbfn, userData)
}
func (this *Tox) CallbackConferenceMessageAdd(cbfn cb_conference_message_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_conference_messages[cbfnp]; ok {
		return
	}
	this.cb_conference_messages[cbfnp] = userData

	if !this.cb_conference_message_setted {
		this.cb_conference_message_setted = true

		C.tox_callback_conference_message(this.toxcore, (*C.tox_conference_message_cb)(C.callbackConferenceMessageWrapperForC))
	}
}

func (this *Tox) CallbackConferenceAction(cbfn cb_conference_action_ftype, userData interface{}) {
	this.CallbackConferenceActionAdd(cbfn, userData)
}
func (this *Tox) CallbackConferenceActionAdd(cbfn cb_conference_action_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_conference_actions[cbfnp]; ok {
		return
	}
	this.cb_conference_actions[cbfnp] = userData

	if !this.cb_conference_message_setted {
		this.cb_conference_message_setted = true
		C.tox_callback_conference_message(this.toxcore, (*C.tox_conference_message_cb)(C.callbackConferenceMessageWrapperForC))
	}
}

//export callbackConferenceTitleWrapperForC
func callbackConferenceTitleWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint32_t, a2 *C.uint8_t, a3 C.size_t, a4 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_conference_titles {
		cbfn := *(*cb_conference_title_ftype)(cbfni)
		title := C.GoStringN((*C.char)((unsafe.Pointer)(a2)), C.int(a3))
		this.putcbevts(func() { cbfn(this, uint32(a0), uint32(a1), title, ud) })
	}
}

func (this *Tox) CallbackConferenceTitle(cbfn cb_conference_title_ftype, userData interface{}) {
	this.CallbackConferenceTitleAdd(cbfn, userData)
}
func (this *Tox) CallbackConferenceTitleAdd(cbfn cb_conference_title_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_conference_titles[cbfnp]; ok {
		return
	}
	this.cb_conference_titles[cbfnp] = userData

	C.tox_callback_conference_title(this.toxcore, (*C.tox_conference_title_cb)(C.callbackConferenceTitleWrapperForC))
}

//export callbackConferencePeerNameWrapperForC
func callbackConferencePeerNameWrapperForC(m *C.Tox, a0 C.uint32_t, a1 C.uint32_t, a2 *C.uint8_t, a3 C.size_t, a4 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_conference_peer_names {
		cbfn := *(*cb_conference_peer_name_ftype)(cbfni)
		peer_name := C.GoStringN((*C.char)((unsafe.Pointer)(a2)), C.int(a3))
		this.putcbevts(func() { cbfn(this, uint32(a0), uint32(a1), peer_name, ud) })
	}
}

func (this *Tox) CallbackConferencePeerName(cbfn cb_conference_peer_name_ftype, userData interface{}) {
	this.CallbackConferencePeerNameAdd(cbfn, userData)
}
func (this *Tox) CallbackConferencePeerNameAdd(cbfn cb_conference_peer_name_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_conference_peer_names[cbfnp]; ok {
		return
	}
	this.cb_conference_peer_names[cbfnp] = userData

	C.tox_callback_conference_peer_name(this.toxcore, (*C.tox_conference_peer_name_cb)(C.callbackConferencePeerNameWrapperForC))
}

//export callbackConferencePeerListChangedWrapperForC
func callbackConferencePeerListChangedWrapperForC(m *C.Tox, a0 C.uint32_t, a1 unsafe.Pointer) {
	var this = cbUserDatas.get(m)
	for cbfni, ud := range this.cb_conference_peer_list_changeds {
		cbfn := *(*cb_conference_peer_list_changed_ftype)(cbfni)
		this.putcbevts(func() { cbfn(this, uint32(a0), ud) })
	}
}

func (this *Tox) CallbackConferencePeerListChanged(cbfn cb_conference_peer_list_changed_ftype, userData interface{}) {
	this.CallbackConferencePeerListChangedAdd(cbfn, userData)
}
func (this *Tox) CallbackConferencePeerListChangedAdd(cbfn cb_conference_peer_list_changed_ftype, userData interface{}) {
	cbfnp := (unsafe.Pointer)(&cbfn)
	if _, ok := this.cb_conference_peer_list_changeds[cbfnp]; ok {
		return
	}
	this.cb_conference_peer_list_changeds[cbfnp] = userData

	C.tox_callback_conference_peer_list_changed(this.toxcore, (*C.tox_conference_peer_list_changed_cb)(C.callbackConferencePeerListChangedWrapperForC))
}

// methods tox_conference_*
func (this *Tox) ConferenceNew() (uint32, error) {
	this.lock()
	defer this.unlock()

	var cerr C.TOX_ERR_CONFERENCE_NEW
	r := C.tox_conference_new(this.toxcore, &cerr)
	if cerr != 0 {
		return uint32(r), toxenums.TOX_ERR_CONFERENCE_NEW(cerr)
	}
	return uint32(r), nil
}

func (this *Tox) ConferenceDelete(groupNumber uint32) error {
	this.lock()

	var _gn = C.uint32_t(groupNumber)
	var cerr C.TOX_ERR_CONFERENCE_DELETE
	C.tox_conference_delete(this.toxcore, _gn, &cerr)
	if cerr != 0 {
		this.unlock()
		return toxenums.TOX_ERR_CONFERENCE_DELETE(cerr)
	}
	this.unlock()
	return nil
}

func (this *Tox) ConferencePeerGetName(groupNumber uint32, peerNumber uint32) (string, error) {
	var _gn = C.uint32_t(groupNumber)
	var _pn = C.uint32_t(peerNumber)
	_name := make([]byte, MAX_NAME_LENGTH)

	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	C.tox_conference_peer_get_name(this.toxcore, _gn, _pn, (*C.uint8_t)(&_name[0]), &cerr)
	if cerr != 0 {
		return "", toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}

	return C.GoString((*C.char)(safeptr(_name[:]))), nil
}

func (this *Tox) ConferencePeerGetPublicKey(groupNumber uint32, peerNumber uint32) (*[PUBLIC_KEY_SIZE]byte, error) {
	var _gn = C.uint32_t(groupNumber)
	var _pn = C.uint32_t(peerNumber)
	var pubkey [PUBLIC_KEY_SIZE]byte

	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	C.tox_conference_peer_get_public_key(this.toxcore, _gn, _pn, (*C.uint8_t)(&pubkey[0]), &cerr)
	if cerr != 0 {
		return nil, toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}

	return &pubkey, nil
}

func (this *Tox) ConferenceInvite(friendNumber uint32, groupNumber uint32) error {
	this.lock()
	defer this.unlock()

	var _fn = C.uint32_t(friendNumber)
	var _gn = C.uint32_t(groupNumber)

	// if give a friendNumber which not exists,
	// the tox_invite_friend has a strange behaive: cause other tox_* call failed
	// and the call will return true, but only strange thing accurs
	// so just precheck the friendNumber and then go
	if !this.FriendExists(friendNumber) {
		return fmt.Errorf("friend not exists: %d", friendNumber)
	}

	var cerr C.TOX_ERR_CONFERENCE_INVITE
	C.tox_conference_invite(this.toxcore, _fn, _gn, &cerr)
	if cerr != 0 {
		return toxenums.TOX_ERR_CONFERENCE_INVITE(cerr)
	}
	return nil
}

func (this *Tox) ConferenceJoin(friendNumber uint32, cookie []byte) (uint32, error) {
	if cookie == nil {
		return 0, toxenums.TOX_ERR_CONFERENCE_JOIN_INVALID_LENGTH
	}

	this.lock()
	defer this.unlock()

	var _fn = C.uint32_t(friendNumber)
	var _length = C.size_t(len(cookie))

	var cerr C.TOX_ERR_CONFERENCE_JOIN
	r := C.tox_conference_join(this.toxcore, _fn, (*C.uint8_t)(&cookie[0]), _length, &cerr)
	if cerr != 0 {
		return uint32(r), toxenums.TOX_ERR_CONFERENCE_JOIN(cerr)
	}
	return uint32(r), nil
}

func (this *Tox) ConferenceSendMessage(groupNumber uint32, mtype toxenums.TOX_MESSAGE_TYPE, message []byte) error {
	this.lock()
	defer this.unlock()

	var _gn = C.uint32_t(groupNumber)
	var _length = C.size_t(len(message))

	switch mtype {
	case toxenums.TOX_MESSAGE_TYPE_NORMAL:
	case toxenums.TOX_MESSAGE_TYPE_ACTION:
	default:
		return fmt.Errorf("Invalid message type: %v", mtype)
	}

	var cerr C.TOX_ERR_CONFERENCE_SEND_MESSAGE
	C.tox_conference_send_message(this.toxcore, _gn, C.TOX_MESSAGE_TYPE(mtype), (*C.uint8_t)(&message[0]), _length, &cerr)
	if cerr != 0 {
		return toxenums.TOX_ERR_CONFERENCE_SEND_MESSAGE(cerr)
	}
	return nil
}

func (this *Tox) ConferenceSetTitle(groupNumber uint32, title string) error {
	this.lock()
	defer this.unlock()

	var _gn = C.uint32_t(groupNumber)
	var _title = []byte(title)
	var _length = C.size_t(len(title))

	var cerr C.TOX_ERR_CONFERENCE_TITLE
	C.tox_conference_set_title(this.toxcore, _gn, (*C.uint8_t)(&_title[0]), _length, &cerr)
	if cerr != 0 {
		return toxenums.TOX_ERR_CONFERENCE_TITLE(cerr)
	}
	return nil
}

func (this *Tox) ConferenceGetTitle(groupNumber uint32) (string, error) {
	var _gn = C.uint32_t(groupNumber)
	_title := make([]byte, MAX_NAME_LENGTH)

	var cerr C.TOX_ERR_CONFERENCE_TITLE
	C.tox_conference_get_title(this.toxcore, _gn, (*C.uint8_t)(&_title[0]), &cerr)
	if cerr != 0 {
		return "", toxenums.TOX_ERR_CONFERENCE_TITLE(cerr)
	}
	return C.GoString((*C.char)(safeptr(_title[:]))), nil
}

func (this *Tox) ConferencePeerNumberIsOurs(groupNumber uint32, peerNumber uint32) (bool, error) {
	var _gn = C.uint32_t(groupNumber)
	var _pn = C.uint32_t(peerNumber)

	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	r := C.tox_conference_peer_number_is_ours(this.toxcore, _gn, _pn, &cerr)
	if cerr != 0 {
		return false, toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}
	return bool(r), nil
}

func (this *Tox) ConferencePeerCount(groupNumber uint32) (uint32, error) {
	var _gn = C.uint32_t(groupNumber)

	var cerr C.TOX_ERR_CONFERENCE_PEER_QUERY
	r := C.tox_conference_peer_count(this.toxcore, _gn, &cerr)
	if cerr != 0 {
		return 0, toxenums.TOX_ERR_CONFERENCE_PEER_QUERY(cerr)
	}
	return uint32(r), nil
}

func (this *Tox) ConferenceGetChatlistSize() uint32 {
	r := C.tox_conference_get_chatlist_size(this.toxcore)
	return uint32(r)
}

func (this *Tox) ConferenceGetChatlist() []uint32 {
	var sz uint32 = this.ConferenceGetChatlistSize()
	vec := make([]uint32, sz)
	if sz == 0 {
		return vec
	}

	vec_p := unsafe.Pointer(&vec[0])
	C.tox_conference_get_chatlist(this.toxcore, (*C.uint32_t)(vec_p))
	return vec
}

func (this *Tox) ConferenceGetType(groupNumber uint32) (t toxenums.TOX_CONFERENCE_TYPE, err error) {
	var _gn = C.uint32_t(groupNumber)

	var cerr C.TOX_ERR_CONFERENCE_GET_TYPE
	t = toxenums.TOX_CONFERENCE_TYPE(C.tox_conference_get_type(this.toxcore, _gn, &cerr))
	if cerr != 0 {
		err = toxenums.TOX_ERR_CONFERENCE_GET_TYPE(cerr)
	}
	return
}
