// server command: ./examples
// client command: ./examples -addr <server toxid>
package main

import (
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/TokTok/go-toxcore-c"
	"github.com/TokTok/go-toxcore-c/toxenums"
)

// tox-account
// server:
//{
//  "Address": "13478A827170C31B4F8BD4FA6E33FF25E0BFBDDC01A755A3124949DC8703E77B24886745508E",
//  "Secret": "0954F174F6DF774EC51BCD23A848DBF50B64B9F8ECC7E30C417A6A964543E192",
//  "Pubkey": "13478A827170C31B4F8BD4FA6E33FF25E0BFBDDC01A755A3124949DC8703E77B",
//  "Nospam": 612919109
//}

var (
	serverAddress = tox.MustDecodeAddress("13478A827170C31B4F8BD4FA6E33FF25E0BFBDDC01A755A3124949DC8703E77B24886745508E")
	serverSecret  = tox.MustDecodeSecret("0954F174F6DF774EC51BCD23A848DBF50B64B9F8ECC7E30C417A6A964543E192")
	serverPubkey  = tox.MustDecodePubkey("13478A827170C31B4F8BD4FA6E33FF25E0BFBDDC01A755A3124949DC8703E77B")

	serverNospam uint32 = 612919109
)

//client:
//{
//  "Address": "3739E8B71573C65A0DEBA5C9F381EF6BD6DA5D9A5E6C7C4C3B0C9567BED7387D3AFA6DC56EF1",
//  "Secret": "423E3CA623F2F6EA8C877758E5A9990A01B2E41081AA5224CA658AF9C93A3B37",
//  "Pubkey": "3739E8B71573C65A0DEBA5C9F381EF6BD6DA5D9A5E6C7C4C3B0C9567BED7387D",
//  "Nospam": 989490629
//}

var (
	clientAddress = tox.MustDecodeAddress("3739E8B71573C65A0DEBA5C9F381EF6BD6DA5D9A5E6C7C4C3B0C9567BED7387D3AFA6DC56EF1")
	clientSecret  = tox.MustDecodeSecret("423E3CA623F2F6EA8C877758E5A9990A01B2E41081AA5224CA658AF9C93A3B37")
	clientPubkey  = tox.MustDecodePubkey("3739E8B71573C65A0DEBA5C9F381EF6BD6DA5D9A5E6C7C4C3B0C9567BED7387D")

	clientNospam uint32 = 989490629
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}

var bootstrapNode = tox.BootstrapNode{
	Addr:    "205.185.116.116",
	Port:    33445,
	TcpPort: 33445,
	Pubkey:  *tox.MustDecodePubkey("A179B09749AC826FF01F37A9613F6B57118AE014D4196A0E1105A98F93A54702"),
}

var fnameServer = "./toxecho-server.data"
var fname = "./toxecho-client.data"
var debug = true
var nickPrefix = "EchoBot."
var statusText = "Send me text, file, audio, video."

var serverMode = flag.Bool("server", false, "server mode")

func main() {
	flag.Parse()

	if *serverMode {
		fname = fnameServer
	}

	var savedata *[32]byte
	var nospam uint32
	if *serverMode {
		savedata = serverSecret
		nospam = serverNospam
	} else {
		savedata = clientSecret
		nospam = clientNospam
	}

	t, err := tox.NewTox(&tox.ToxOptions{
		Savedata_type:           toxenums.TOX_SAVEDATA_TYPE_SECRET_KEY,
		Savedata_data:           savedata[:],
		Tcp_port:                33445,
		NospamIfSecretType:      nospam,
		ProxyToNoneIfErr:        true,
		AutoTcpPortIfErr:        true,
		DisableTcpPortIfAutoErr: true,
		PingUnit:                time.Second,
	})
	if err != nil {
		log.Println("NewTox", err)
	}

	berr := t.BootstrapNode_l(&bootstrapNode)
	if debug && berr != 0 {
		log.Println("Bootstrap:", berr)
	}

	if debug {
		log.Printf("keys: secret:%X pubkey:%X\n", t.Secret, t.Pubkey)
	}
	log.Printf("toxid: %X\n", t.Address)

	defaultName := t.SelfGetName_l()
	humanName := nickPrefix + hex.EncodeToString(t.Address[:])[:5]
	t.SelfSetName_l(humanName)
	humanName = t.SelfGetName_l()
	if debug {
		log.Println("nickName:", humanName, "defaultName:", defaultName, err)
	}

	t.SelfSetStatusMessage_l(statusText)
	if debug {
		log.Println(statusText, t.SelfGetStatusMessage_l())
	}

	sd := t.GetSavedata_l()
	if debug {
		log.Println("savedata", len(sd), t)
	}

	// add friend norequest
	fv := t.SelfGetFriendList()
	if debug {
		log.Println("add friends:", len(fv))
		for fn, pk := range fv {
			log.Printf("friend: %d %X\n", fn, *pk)
		}
	}

	workerCh := make(chan interface{}, 32)
	workerSavedata := make(chan []byte, 1)
	workerSavedata <- nil

	// callbacks
	t.CallbackSelfConnectionStatus(func(status toxenums.TOX_CONNECTION) {
		if debug {
			log.Println("on self conn status:", status)
		}
		if status != toxenums.TOX_CONNECTION_NONE && !*serverMode {
			_, err := t.FriendAdd_l(serverAddress, []byte("Hi! I am a tunnel client."))
			if err != nil {
				log.Println("FriendAdd_l", err)
			}
		}
	})
	t.CallbackFriendRequest(func(friendId *[tox.PUBLIC_KEY_SIZE]byte, message []byte) {
		log.Printf("%X: %s\n", friendId, message)
		num, err := t.FriendAddNorequest_l(friendId)
		if debug {
			log.Println("on friend request:", num, err)
		}
		if num < 100000 {
			workerSavedata <- t.GetSavedata_l()
		}
	})
	t.CallbackFriendMessage(func(friendNumber uint32, mtype toxenums.TOX_MESSAGE_TYPE, message []byte) {
		if debug {
			log.Printf("on friend message: %d %s\n", friendNumber, message)
		}
		n, err := t.FriendSendMessage_l(friendNumber, mtype, append([]byte("RE: "), message...))
		if err != nil {
			log.Println(n, err)
		}
	})
	t.CallbackFriendConnectionStatus(func(friendNumber uint32, status toxenums.TOX_CONNECTION) {
		if debug {
			friendId, err := t.FriendGetPublicKey(friendNumber)
			log.Printf("on friend connection status: %d %v %X %v\n", friendNumber, status, friendId, err)
		}
		if !*serverMode {
			c, err := t.Dial_l(friendNumber)
			if err != nil {
				log.Printf("dail tcp to friend failed: %d %v\n", friendNumber, err)
				log.Println("if not a dup dial, save the pubkey and try Dial later from out side of the callbacks")
				return
			}
			go onClientConn(c)
		}
	})
	t.CallbackFriendStatus(func(friendNumber uint32, status toxenums.TOX_USER_STATUS) {
		if debug {
			friendId, err := t.FriendGetPublicKey(friendNumber)
			log.Printf("on friend status: %d %v %X %v\n", friendNumber, status, friendId, err)
		}
	})
	t.CallbackFriendStatusMessage(func(friendNumber uint32, statusText string) {
		if debug {
			friendId, err := t.FriendGetPublicKey(friendNumber)
			log.Printf("on friend status text: %d %v %X %v\n", friendNumber, statusText, friendId, err)
		}
	})

	t.CallbackFileRecvControl(func(friendNumber uint32, fileNumber uint32, control toxenums.TOX_FILE_CONTROL) {
		if debug {
			friendId, err := t.FriendGetPublicKey(friendNumber)
			log.Printf("on recv file control: %d %d %v %X %v\n", friendNumber, fileNumber, control, friendId, err)
		}
		// worker's job now.
		workerCh <- &tox.FileControlData{
			FriendNumber: friendNumber,
			FileNumber:   fileNumber,
			Control:      control,
			Result:       make(chan error, 1),
		}
	})
	t.CallbackFileRecv(func(friendNumber uint32, fileNumber uint32, kind toxenums.TOX_FILE_KIND, fileSize uint64, fileName []byte) {
		if debug {
			friendId, err := t.FriendGetPublicKey(friendNumber)
			log.Printf("on recv file: %d %d %v %d %s %X %v\n", friendNumber, fileNumber, kind, fileSize, fileName, friendId, err)
		}
		if fileSize > 1024*1024*1024 {
			// good guy
		}
		// worker's job now.
		toWorkerData := make(chan *tox.FileSendResult, 1)
		toWorkerData <- &tox.FileSendResult{FileNumber: fileNumber}
		workerCh <- &tox.FileSendData{
			FriendNumber: friendNumber,
			Kind:         kind,
			FileSize:     fileSize,
			FileId:       nil,
			FileName:     append([]byte("RE_"), fileName...),
			Result:       toWorkerData,
		}
	})
	t.CallbackFileRecvChunk(func(friendNumber uint32, fileNumber uint32, position uint64, data []byte) {
		friendId, err := t.FriendGetPublicKey(friendNumber)
		if len(data) == 0 {
			if debug {
				log.Printf("recv file finished: %d %d %X %v\n", friendNumber, fileNumber, friendId, err)
			}
			return
		}
		// worker's job now.
		workerCh <- &tox.FileSendChunkData{
			FriendNumber: friendNumber,
			FileNumber:   fileNumber,
			Position:     position,
			Data:         data,
			Result:       nil,
		}
	})
	t.CallbackFileChunkRequest(func(friendNumber uint32, fileNumber uint32, position uint64, length int) {
		friendId, err := t.FriendGetPublicKey(friendNumber)
		if length == 0 {
			if debug {
				log.Printf("send file finished: %d %d %X %v\n", friendNumber, fileNumber, friendId, err)
			}
			return
		}
		// worker's job now.
		workerCh <- &FileChunkRequestData{friendNumber, fileNumber, position, length}
	})

	t.CallbackFriendLosslessPacket(t.ParseLosslessPacket)
	t.CallbackTcpPong(func(friendNumber uint32, ms uint32) {
		log.Println("Ping", friendNumber, ms)
	})

	// audio/video
	av, err := tox.NewToxAV(t)
	if err != nil {
		log.Println(err, av)
	}
	if av == nil {
	}
	av.CallbackCall(func(friendNumber uint32, audioEnabled bool, videoEnabled bool) {
		if debug {
			log.Println("oncall:", friendNumber, audioEnabled, videoEnabled)
		}
		var audioBitRate uint32 = 48
		var videoBitRate uint32 = 64
		err := av.Answer(friendNumber, audioBitRate, videoBitRate)
		if err != nil {
			log.Println(err)
		}
	})
	av.CallbackCallState(func(friendNumber uint32, state toxenums.TOXAV_FRIEND_CALL_STATE) {
		if debug {
			log.Println("on call state:", friendNumber, state)
		}
	})
	av.CallbackAudioReceiveFrame(func(friendNumber uint32, pcm []byte, sampleCount int, channels int, samplingRate int) {
		if debug {
			if rand.Int()%23 == 3 {
				log.Println("on recv audio frame:", friendNumber, len(pcm), sampleCount, channels, samplingRate)
			}
		}
		err := av.AudioSendFrame(friendNumber, pcm, sampleCount, channels, samplingRate)
		if err != nil {
			log.Println(err)
		}
	})
	av.CallbackVideoReceiveFrame(func(friendNumber uint32, width uint16, height uint16, frames []byte) {
		if debug {
			if rand.Int()%45 == 3 {
				log.Println("on recv video frame:", friendNumber, width, height, len(frames))
			}
		}
		err := av.VideoSendFrame(friendNumber, width, height, frames)
		if err != nil {
			log.Println(err)
		}
	})

	// toxav loops
	go func() {
		shutdown := false
		loopc := 0
		itval := uint64(0)
		for !shutdown {
			iv := av.IterationInterval()
			if iv != itval {
				// wtf
				if iv-itval > 20 || itval-iv > 20 {
					log.Println("av itval changed:", itval, iv, iv-itval, itval-iv)
				}
				itval = iv
			}

			av.Iterate()
			loopc += 1
			time.Sleep(1000 * 50 * time.Microsecond)
		}

		av.Kill()
		t.Stop()
		t.Kill()
	}()

	// toxcore loops
	go t.Run()

	// tcp tunnel server
	if *serverMode {
		go func() {
			for {
				c, err := t.Accept()
				if err != nil {
					return
				}
				go func() {
					io.Copy(c, c)
					c.Close()
				}()
			}
		}()
	}

	// worker loop below

	// worker stuff
	var recvFiles = make(map[uint64]uint32)
	var sendFiles = make(map[uint64]uint32)
	var sendDatas = make(map[string][]byte)
	var chunkReqs = make([]string, 0, 8)
	trySendChunk := func(friendNumber uint32, fileNumber uint32, position uint64) {
		sentKeys := make(map[string]bool)
		for _, reqkey := range chunkReqs {
			lst := strings.Split(reqkey, "_")
			pos, err := strconv.ParseUint(lst[2], 10, 64)
			if err != nil {
			}
			if data, ok := sendDatas[reqkey]; ok {
				req := &tox.FileSendChunkData{
					FriendNumber: friendNumber,
					FileNumber:   fileNumber,
					Position:     pos,
					Data:         data,
					Result:       make(chan error, 1),
				}
				t.DoInLoop(req)
				if err = <-req.Result; err != nil {
					if terr, ok := err.(toxenums.TOX_ERR_FILE_SEND_CHUNK); ok && terr > 6 {
					} else {
						log.Println("file send chunk error:", err, reqkey)
					}
					break
				} else {
					delete(sendDatas, reqkey)
					sentKeys[reqkey] = true
				}
			}
		}
		leftChunkReqs := make([]string, 0)
		for _, reqkey := range chunkReqs {
			if _, ok := sentKeys[reqkey]; !ok {
				leftChunkReqs = append(leftChunkReqs, reqkey)
			}
		}
		chunkReqs = leftChunkReqs
	}
	// worker stuff end

	for {
		select {
		case idata := <-workerCh:
			switch data := idata.(type) {
			case *tox.FileSendData:
				result := <-data.Result
				fileNumber := result.FileNumber

				t.DoInLoop(data)
				result = <-data.Result
				if result.Error != nil {
					log.Println("RE file:", data.FileName, result.Error)
					continue
				}

				recvFiles[uint64(data.FriendNumber)<<32|uint64(fileNumber)] = result.FileNumber
				sendFiles[uint64(data.FriendNumber)<<32|uint64(result.FileNumber)] = fileNumber

			case *tox.FileControlData:
				key := uint64(data.FriendNumber)<<32 | uint64(data.FileNumber)
				fno, ok := sendFiles[key]
				if !ok {
					log.Println("RE control FileNumber not found:", data.FileNumber)
					continue
				}
				data.FileNumber = fno
				t.DoInLoop(data)
				if err := <-data.Result; err != nil {
					log.Println("RE control:", data.FileNumber, err)
				}

			case *tox.FileSendChunkData:
				reFileNumber := recvFiles[uint64(data.FriendNumber)<<32|uint64(data.FileNumber)]
				key := makekey(data.FriendNumber, reFileNumber, data.Position)
				sendDatas[key] = data.Data
				trySendChunk(data.FriendNumber, reFileNumber, data.Position)

			case *FileChunkRequestData:
				if data.Length == 0 {
					origFileNumber := sendFiles[uint64(data.FriendNumber)<<32|uint64(data.FileNumber)]
					delete(sendFiles, uint64(data.FriendNumber)<<32|uint64(data.FileNumber))
					delete(recvFiles, uint64(data.FriendNumber)<<32|uint64(origFileNumber))
				} else {
					key := makekey(data.FriendNumber, data.FileNumber, data.Position)
					chunkReqs = append(chunkReqs, key)
					trySendChunk(data.FriendNumber, data.FileNumber, data.Position)
				}

			default:
				log.Fatalln("should not go here")
			}

		case liveData := <-workerSavedata:
			if liveData == nil {
				data := make(tox.GetSavedataData, 1)
				t.DoInLoop(data)
				liveData = <-data
			}
			err = tox.WriteSavedata(fname, liveData)
			if debug {
				log.Println("savedata write:", err)
			}
		}
	}
}

type FileChunkRequestData struct {
	FriendNumber uint32
	FileNumber   uint32
	Position     uint64
	Length       int
}

func onClientConn(c net.Conn) {
	defer c.Close()

	go func() {
		buf := make([]byte, 16)
		for {
			_, err := c.Read(buf)
			if err != nil {
				log.Println("client reade done")
				return
			}
			log.Println("\nserver:", string(buf))
		}
	}()

	var i uint32
	var s [4]byte
	for {
		binary.BigEndian.PutUint32(s[:], i)
		msg := fmt.Sprintf("message %X", s)
		i++
		log.Println("\nclient:", msg, len(msg))
		_, err := c.Write([]byte(msg))
		if err != nil {
			log.Println("client write done")
			return
		}
		time.Sleep(time.Second * 2)
	}
}

func makekey(no uint32, a0 interface{}, a1 interface{}) string {
	return fmt.Sprintf("%d_%v_%v", no, a0, a1)
}
