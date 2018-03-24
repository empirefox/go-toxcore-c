package tox

//#include <tox/tox.h>
import "C"
import (
	"time"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

func (t *Tox) Run() {
	timer := time.NewTimer(0)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	pingUnit := t.pingUnit
	if pingUnit < DefaultPingUnit {
		pingUnit = DefaultPingUnit
	}
	pingTicker := time.NewTicker(pingUnit)
	defer pingTicker.Stop()

	for {
		select {
		case <-timer.C:
			t.inToxIterate = true
			C.tox_iterate(t.toxcore, nil)
			t.inToxIterate = false

			ms := time.Duration(C.tox_iteration_interval(t.toxcore)) * time.Millisecond
			if t.cbPostIterate != nil {
				for _, cb := range t.cbPostIterate {
					ms -= cb()
				}
				t.cbPostIterate = nil
			}
			timer.Reset(ms)

		case idata := <-t.chLoopRequest:
			t.doInLoop(idata)

		case <-pingTicker.C:
			t.doTcpPing_l()

		case <-t.stop:
			close(t.stopped)
			return
		}
	}
}

func (t *Tox) doInLoop(idata interface{}) {
	switch data := idata.(type) {
	case *sendTcpPacketData:
		t.sendTcpPacket_l(data)

	case *FriendSendLossyPacketData:
		t.friendSendLossyPacket_l(data)
	case *FriendSendLosslessPacketData:
		t.friendSendLosslessPacket_l(data)

	case *FriendSendMessageData:
		t.friendSendMessage_l(data)

	case *BootstrapData:
		t.bootstrapNodes_l(data)
	case GetSavedataData:
		t.getSavedata_l(data)

		// Friend add/delete should block av loop
		// TODO add block av logic
	case *FriendAddData:
		t.blockAv()
		t.friendAdd_l(data)
		t.unblockAv()
	case *FriendAddNorequestData:
		t.blockAv()
		t.friendAddNorequest_l(data)
		t.unblockAv()
	case *FriendDeleteData:
		t.blockAv()
		t.friendDelete_l(data)
		t.unblockAv()

	case *SelfSetNameData:
		t.selfSetName_l(data)
	case SelfGetNameData:
		t.selfGetName_l(data)
	case *SelfSetStatusMessageData:
		t.selfSetStatusMessage_l(data)
	case SelfGetStatusMessageData:
		t.selfGetStatusMessage_l(data)
	case SelfSetStatusData:
		t.selfSetStatus_l(data)
	case SelfGetStatusData:
		t.selfGetStatus_l(data)
	case SelfSetNospamData:
		t.selfSetNospam_l(data)
	case SelfGetNospamData:
		t.selfGetNospam_l(data)

	case *FriendGetNameData:
		t.friendGetName_l(data)
	case *FriendGetStatusMessageData:
		t.friendGetStatusMessage_l(data)
	case *FriendGetStatusData:
		t.friendGetStatus_l(data)
	case *FriendGetLastOnlineData:
		t.friendGetLastOnline_l(data)

	case *SelfSetTypingData:
		t.selfSetTyping_l(data)

	case *FileControlData:
		t.fileControl_l(data)
	case *FileSendData:
		t.fileSend_l(data)
	case *FileSendChunkData:
		t.fileSendChunk_l(data)
	case *FileSeekData:
		t.fileSeek_l(data)
	case *FileGetFileIdData:
		t.fileGetFileId_l(data)

	case ConferenceNewData:
		t.conferenceNew_l(data)
	case *ConferenceDeleteData:
		t.conferenceDelete_l(data)
	case *ConferencePeerGetNameData:
		t.conferencePeerGetName_l(data)
	case *ConferencePeerGetPublicKeyData:
		t.conferencePeerGetPublicKey_l(data)
	case *ConferenceInviteData:
		t.conferenceInvite_l(data)
	case *ConferenceJoinData:
		t.conferenceJoin_l(data)
	case *ConferenceSendMessageData:
		t.conferenceSendMessage_l(data)
	case *ConferenceSetTitleData:
		t.conferenceSetTitle_l(data)
	case *ConferenceGetTitleData:
		t.conferenceGetTitle_l(data)
	case *ConferencePeerNumberIsOursData:
		t.conferencePeerNumberIsOurs_l(data)
	case *ConferencePeerCountData:
		t.conferencePeerCount_l(data)
	case ConferenceGetChatlistData:
		t.conferenceGetChatlist_l(data)
	case *ConferenceGetTypeData:
		t.conferenceGetType_l(data)

	case *PingMultipleData:
		t.setPingMultiple_l(data)

	case IterateThenData:
		if _, ok := data.Do.(IterateThenData); ok {
			panic("IterateThenData.Do should not be *IterateThenData")
		}
		C.tox_iterate(t.toxcore, nil)
		t.doInLoop(data.Do)
	case func():
		data()

	default:
		panic("should not go here")
	}
}

func (t *Tox) blockAv()   {}
func (t *Tox) unblockAv() {}

type (
	FriendSendLossyPacketData struct {
		FriendNumber uint32
		Data         []byte
		Result       chan error
	}
	FriendSendLosslessPacketData struct {
		FriendNumber uint32
		Data         []byte
		Result       chan error
	}
)

type (
	BootstrapNode struct {
		Addr    string
		Port    uint16
		TcpPort uint16
		Pubkey  [PUBLIC_KEY_SIZE]byte
	}
	BootstrapData struct {
		Nodes  []BootstrapNode
		Result chan *BootstrapResult
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

type (
	GetSavedataData chan []byte
)

type (
	FriendAddData struct {
		Address *[ADDRESS_SIZE]byte
		Message []byte
		Result  chan *FriendAddResult
	}
	FriendAddNorequestData struct {
		Pubkey *[PUBLIC_KEY_SIZE]byte
		Result chan *FriendAddResult
	}
	FriendAddResult struct {
		FriendNumber uint32
		Error        error
	}

	FriendDeleteData struct {
		FriendNumber uint32
		Result       chan error
	}

	FriendSendMessageData struct {
		FriendNumber uint32
		Type         toxenums.TOX_MESSAGE_TYPE
		Message      []byte
		Result       chan *FriendSendMessageResult
	}
	FriendSendMessageResult struct {
		MessageId uint32
		Error     error
	}
)

type (
	SelfSetNameData struct {
		Name   string
		Result chan error
	}
	SelfGetNameData chan string

	SelfSetStatusMessageData struct {
		Message string
		Result  chan error
	}
	SelfGetStatusMessageData chan string

	SelfSetStatusData toxenums.TOX_USER_STATUS
	SelfGetStatusData chan toxenums.TOX_USER_STATUS

	SelfSetNospamData uint32
	SelfGetNospamData chan uint32
)

type (
	FriendGetNameData struct {
		FriendNumber uint32
		Result       chan *FriendGetNameResult
	}
	FriendGetNameResult struct {
		Name  string
		Error error
	}

	FriendGetStatusMessageData struct {
		FriendNumber uint32
		Result       chan *FriendGetStatusMessageResult
	}
	FriendGetStatusMessageResult struct {
		Message string
		Error   error
	}

	FriendGetStatusData struct {
		FriendNumber uint32
		Result       chan *FriendGetStatusResult
	}
	FriendGetStatusResult struct {
		Status toxenums.TOX_USER_STATUS
		Error  error
	}

	FriendGetLastOnlineData struct {
		FriendNumber uint32
		Result       chan *FriendGetLastOnlineResult
	}
	FriendGetLastOnlineResult struct {
		Unix  uint64
		Error error
	}
)

type (
	SelfSetTypingData struct {
		FriendNumber uint32
		Typing       bool
		Result       chan error
	}
)

type (
	FileControlData struct {
		FriendNumber uint32
		FileNumber   uint32
		Control      toxenums.TOX_FILE_CONTROL
		Result       chan error
	}

	FileSendData struct {
		FriendNumber uint32
		Kind         toxenums.TOX_FILE_KIND
		FileSize     uint64
		FileId       *[FILE_ID_LENGTH]byte
		FileName     []byte
		Result       chan *FileSendResult
	}
	FileSendResult struct {
		FileNumber uint32
		Error      error
	}
	FileSendChunkData struct {
		FriendNumber uint32
		FileNumber   uint32
		Position     uint64
		Data         []byte
		Result       chan error
	}
	FileSeekData struct {
		FriendNumber uint32
		FileNumber   uint32
		Position     uint64
		Result       chan error
	}
	FileGetFileIdData struct {
		FriendNumber uint32
		FileNumber   uint32
		Result       chan *FileGetFileIdResult
	}
	FileGetFileIdResult struct {
		FileId *[FILE_ID_LENGTH]byte
		Error  error
	}
)

type (
	ConferenceNewData   chan *ConferenceNewResult
	ConferenceNewResult struct {
		ConferenceNumber uint32
		Error            error
	}

	ConferenceDeleteData struct {
		ConferenceNumber uint32
		Result           chan error
	}

	ConferencePeerGetNameData struct {
		ConferenceNumber uint32
		PeerNumber       uint32
		Result           chan *ConferencePeerGetNameResult
	}
	ConferencePeerGetNameResult struct {
		Name  string
		Error error
	}

	ConferencePeerGetPublicKeyData struct {
		ConferenceNumber uint32
		PeerNumber       uint32
		Result           chan *ConferencePeerGetPublicKeyResult
	}
	ConferencePeerGetPublicKeyResult struct {
		Pubkey *[PUBLIC_KEY_SIZE]byte
		Error  error
	}

	ConferenceInviteData struct {
		ConferenceNumber uint32
		FriendNumber     uint32
		Result           chan error
	}

	ConferenceJoinData struct {
		FriendNumber uint32
		Cookie       []byte
		Result       chan *ConferenceJoinResult
	}
	ConferenceJoinResult struct {
		ConferenceNumber uint32
		Error            error
	}

	ConferenceSendMessageData struct {
		ConferenceNumber uint32
		Type             toxenums.TOX_MESSAGE_TYPE
		Message          []byte
		Result           chan error
	}

	ConferenceSetTitleData struct {
		ConferenceNumber uint32
		Title            string
		Result           chan error
	}
	ConferenceGetTitleData struct {
		ConferenceNumber uint32
		Result           chan *ConferenceGetTitleResult
	}
	ConferenceGetTitleResult struct {
		Title string
		Error error
	}

	ConferencePeerNumberIsOursData struct {
		ConferenceNumber uint32
		PeerNumber       uint32
		Result           chan *ConferencePeerNumberIsOursResult
	}
	ConferencePeerNumberIsOursResult struct {
		Is    bool
		Error error
	}

	ConferencePeerCountData struct {
		ConferenceNumber uint32
		Result           chan *ConferencePeerCountResult
	}
	ConferencePeerCountResult struct {
		Count uint32
		Error error
	}

	ConferenceGetChatlistData chan []uint32

	ConferenceGetTypeData struct {
		ConferenceNumber uint32
		Result           chan *ConferenceGetTypeResult
	}
	ConferenceGetTypeResult struct {
		Type  toxenums.TOX_CONFERENCE_TYPE
		Error error
	}
)

type (
	IterateThenData struct {
		Do interface{}
	}
)
