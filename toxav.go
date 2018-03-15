package tox

/*
#include <vpx/vpx_image.h>
#include <tox/tox.h>
#include <tox/toxav.h>

void callbackCallWrapperForC(ToxAV *toxAV, uint32_t friend_number, bool audio_enabled, bool video_enabled, void *user_data);
void callbackCallStateWrapperForC(ToxAV *toxAV, uint32_t friendNumber, uint32_t state, void* user_data);
void callbackAudioBitRateWrapperForC(ToxAV *toxAV, uint32_t friendNumber, uint32_t audioBitRate, void* user_data);
void callbackVideoBitRateWrapperForC(ToxAV *toxAV, uint32_t friendNumber, uint32_t videoBitRate, void* user_data);
void callbackAudioReceiveFrameWrapperForC(ToxAV *toxAV, uint32_t friendNumber, int16_t *pcm, size_t sample_count, uint8_t channels, uint32_t sampling_rate, void* user_data);
void callbackVideoReceiveFrameWrapperForC(ToxAV *toxAV, uint32_t friendNumber, uint16_t width, uint16_t height,
     uint8_t *y, uint8_t *u, uint8_t *v, int32_t ystride, int32_t ustride, int32_t vstride, void* user_data);

extern void i420_to_rgb(int width, int height, const uint8_t *y, const uint8_t *u, const uint8_t *v,
            int ystride, int ustride, int vstride, unsigned char *out);
extern void rgb_to_i420(unsigned char* rgb, vpx_image_t *img);
*/
import "C"
import (
	"errors"
	"unsafe"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

type cb_call_ftype func(friendNumber uint32, audioEnabled bool, videoEnabled bool)
type cb_call_state_ftype func(friendNumber uint32, state toxenums.TOXAV_FRIEND_CALL_STATE)
type cb_audio_bit_rate_ftype func(friendNumber uint32, audioBitRate uint32)
type cb_video_bit_rate_ftype func(friendNumber uint32, videoBitRate uint32)
type cb_audio_receive_frame_ftype func(friendNumber uint32, pcm []byte, sampleCount int, channels int, samplingRate int)
type cb_video_receive_frame_ftype func(friendNumber uint32, width uint16, height uint16, data []byte)

// TODO add av loop to send/recv
// TODO controlled by tox loop. Friend add/delete should block av loop
// TODO webrtc to work with browser?
type ToxAV struct {
	toxav *C.ToxAV

	// session datas
	out_image  []byte
	out_width  C.uint16_t
	out_hegith C.uint16_t
	in_image   *C.vpx_image_t
	in_width   C.uint16_t
	in_height  C.uint16_t

	// callbacks
	cb_call                cb_call_ftype
	cb_call_state          cb_call_state_ftype
	cb_audio_bit_rate      cb_audio_bit_rate_ftype
	cb_video_bit_rate      cb_video_bit_rate_ftype
	cb_audio_receive_frame cb_audio_receive_frame_ftype
	cb_video_receive_frame cb_video_receive_frame_ftype
}

func NewToxAV(t *Tox) (*ToxAV, error) {
	var cerr C.TOXAV_ERR_NEW
	toxav := C.toxav_new(t.toxcore, &cerr)
	if cerr != 0 {
		return nil, toxenums.TOXAV_ERR_NEW(cerr)
	}

	tav := &ToxAV{
		toxav: toxav,
	}
	cbAVUserDatas.set(toxav, tav)
	return tav, nil
}

func (this *ToxAV) Kill() {
	C.toxav_kill(this.toxav)
}

func (this *ToxAV) IterationInterval() uint64 {
	return uint64(C.toxav_iteration_interval(this.toxav))
}

func (this *ToxAV) Iterate() {
	C.toxav_iterate(this.toxav)
}

func (this *ToxAV) Call(friendNumber uint32, audioBitRate uint32, videoBitRate uint32) error {
	var cerr C.TOXAV_ERR_CALL
	C.toxav_call(this.toxav, C.uint32_t(friendNumber), C.uint32_t(audioBitRate), C.uint32_t(videoBitRate), &cerr)
	if cerr != 0 {
		return toxenums.TOXAV_ERR_CALL(cerr)
	}
	return nil
}

func (t *ToxAV) CToxAV() *C.ToxAV { return t.toxav }

var cbAVUserDatas = newUserDataAV()

//export callbackCallWrapperForC
func callbackCallWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, audioEnabled C.bool, videoEnabled C.bool, ud unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	this.cb_call(uint32(friendNumber), bool(audioEnabled), bool(videoEnabled))
}
func (this *ToxAV) CallbackCall(cbfn cb_call_ftype) {
	this.cb_call = cbfn
	C.toxav_callback_call(this.toxav, (*C.toxav_call_cb)(C.callbackCallWrapperForC), nil)
}

func (this *ToxAV) Answer(friendNumber uint32, audioBitRate uint32, videoBitRate uint32) error {
	var cerr C.TOXAV_ERR_ANSWER
	C.toxav_answer(this.toxav, C.uint32_t(friendNumber), C.uint32_t(audioBitRate), C.uint32_t(videoBitRate), &cerr)
	if cerr != 0 {
		return toxenums.TOXAV_ERR_ANSWER(cerr)
	}
	return nil
}

//export callbackCallStateWrapperForC
func callbackCallStateWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, state C.uint32_t, ud unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	this.cb_call_state(uint32(friendNumber), toxenums.TOXAV_FRIEND_CALL_STATE(state))
}
func (this *ToxAV) CallbackCallState(cbfn cb_call_state_ftype) {
	this.cb_call_state = cbfn
	C.toxav_callback_call_state(this.toxav, (*C.toxav_call_state_cb)(C.callbackCallStateWrapperForC), nil)
}

func (this *ToxAV) CallControl(friendNumber uint32, control int) error {
	var cerr C.TOXAV_ERR_CALL_CONTROL
	C.toxav_call_control(this.toxav, C.uint32_t(friendNumber), C.TOXAV_CALL_CONTROL(control), &cerr)
	if cerr != 0 {
		return toxenums.TOXAV_ERR_CALL_CONTROL(cerr)
	}
	return nil
}

func (this *ToxAV) AudioSetBitRate(friendNumber uint32, audioBitRate uint32) error {
	var cerr C.TOXAV_ERR_BIT_RATE_SET
	C.toxav_audio_set_bit_rate(this.toxav, C.uint32_t(friendNumber), C.uint32_t(audioBitRate), &cerr)
	if cerr != 0 {
		return toxenums.TOXAV_ERR_BIT_RATE_SET(cerr)
	}
	return nil
}

func (this *ToxAV) VideoSetBitRate(friendNumber uint32, videoBitRate uint32) error {
	var cerr C.TOXAV_ERR_BIT_RATE_SET
	C.toxav_video_set_bit_rate(this.toxav, C.uint32_t(friendNumber), C.uint32_t(videoBitRate), &cerr)
	if cerr != 0 {
		return toxenums.TOXAV_ERR_BIT_RATE_SET(cerr)
	}
	return nil
}

//export callbackAudioBitRateWrapperForC
func callbackAudioBitRateWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, audioBitRate C.uint32_t, ud unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	this.cb_audio_bit_rate(uint32(friendNumber), uint32(audioBitRate))
}
func (this *ToxAV) CallbackAudioBitRate(cbfn cb_audio_bit_rate_ftype) {
	this.cb_audio_bit_rate = cbfn
	C.toxav_callback_audio_bit_rate(this.toxav, (*C.toxav_audio_bit_rate_cb)(C.callbackAudioBitRateWrapperForC), nil)
}

//export callbackVideoBitRateWrapperForC
func callbackVideoBitRateWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, videoBitRate C.uint32_t, ud unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	this.cb_video_bit_rate(uint32(friendNumber), uint32(videoBitRate))
}
func (this *ToxAV) CallbackVideoBitRate(cbfn cb_video_bit_rate_ftype) {
	this.cb_video_bit_rate = cbfn
	C.toxav_callback_video_bit_rate(this.toxav, (*C.toxav_video_bit_rate_cb)(C.callbackVideoBitRateWrapperForC), nil)
}

func (this *ToxAV) AudioSendFrame(friendNumber uint32, pcm []byte, sampleCount int, channels int, samplingRate int) error {
	pcm_ := (*C.int16_t)(unsafe.Pointer(&pcm[0]))
	var cerr C.TOXAV_ERR_SEND_FRAME
	C.toxav_audio_send_frame(this.toxav, C.uint32_t(friendNumber), pcm_, C.size_t(sampleCount), C.uint8_t(channels), C.uint32_t(samplingRate), &cerr)
	if cerr != 0 {
		return toxenums.TOXAV_ERR_SEND_FRAME(cerr)
	}
	return nil
}

func (this *ToxAV) VideoSendFrame(friendNumber uint32, width uint16, height uint16, data []byte) error {
	if this.in_image != nil && (uint16(this.in_width) != width || uint16(this.in_height) != height) {
		C.vpx_img_free(this.in_image)
		this.in_image = nil
	}

	if this.in_image == nil {
		this.in_width = C.uint16_t(width)
		this.in_height = C.uint16_t(height)
		this.in_image = C.vpx_img_alloc(nil, C.VPX_IMG_FMT_I420, C.uint(this.in_width), C.uint(this.in_height), 1)
	}

	C.rgb_to_i420((*C.uchar)(unsafe.Pointer(&data[0])), this.in_image)

	var cerr C.TOXAV_ERR_SEND_FRAME
	C.toxav_video_send_frame(this.toxav, C.uint32_t(friendNumber), C.uint16_t(width), C.uint16_t(height),
		(*C.uint8_t)(this.in_image.planes[0]),
		(*C.uint8_t)(this.in_image.planes[1]),
		(*C.uint8_t)(this.in_image.planes[2]),
		&cerr)
	if cerr != 0 {
		return toxenums.TOXAV_ERR_SEND_FRAME(cerr)
	}
	return nil
}

//export callbackAudioReceiveFrameWrapperForC
func callbackAudioReceiveFrameWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, pcm *C.int16_t, sampleCount C.size_t, channels C.uint8_t, samplingRate C.uint32_t, ud unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)
	length := sampleCount * C.size_t(channels) * 2
	pcm_p := unsafe.Pointer(pcm)
	pcm_b := C.GoBytes(pcm_p, C.int(length))
	this.cb_audio_receive_frame(uint32(friendNumber), pcm_b, int(sampleCount), int(channels), int(samplingRate))
}
func (this *ToxAV) CallbackAudioReceiveFrame(cbfn cb_audio_receive_frame_ftype) {
	this.cb_audio_receive_frame = cbfn
	C.toxav_callback_audio_receive_frame(this.toxav, (*C.toxav_audio_receive_frame_cb)(C.callbackAudioReceiveFrameWrapperForC), nil)
}

//export callbackVideoReceiveFrameWrapperForC
func callbackVideoReceiveFrameWrapperForC(m *C.ToxAV, friendNumber C.uint32_t, width C.uint16_t, height C.uint16_t, y *C.uint8_t, u *C.uint8_t, v *C.uint8_t, ystride C.int32_t, ustride C.int32_t, vstride C.int32_t, ud unsafe.Pointer) {
	var this = cbAVUserDatas.get(m)

	if this.out_image != nil && (this.out_width != width || this.out_hegith != height) {
		this.out_image = nil
	}

	var buf_size int = int(width) * int(height) * 3

	if this.out_image == nil {
		this.out_width = width
		this.out_hegith = height
		this.out_image = make([]byte, buf_size, buf_size)
	}

	out := unsafe.Pointer(&(this.out_image[0]))
	C.i420_to_rgb(C.int(width), C.int(height), y, u, v, C.int(ystride), C.int(ustride), C.int(vstride), (*C.uchar)(out))

	this.cb_video_receive_frame(uint32(friendNumber), uint16(width), uint16(height), this.out_image)
}
func (this *ToxAV) CallbackVideoReceiveFrame(cbfn cb_video_receive_frame_ftype) {
	this.cb_video_receive_frame = cbfn
	C.toxav_callback_video_receive_frame(this.toxav, (*C.toxav_video_receive_frame_cb)(C.callbackVideoReceiveFrameWrapperForC), nil)
}

// TODO
// toxav_add_av_groupchat
// toxav_join_av_groupchat
// toxav_group_send_audio

func (this *Tox) AddAVGroupChat() int {
	r := C.toxav_add_av_groupchat(this.toxcore, nil, nil)
	return int(r)
}

func (this *Tox) JoinAVGroupChat(friendNumber uint32, cookie []byte) (int, error) {
	var _fn = C.uint32_t(friendNumber)
	var _data = (*C.uint8_t)(&cookie[0])
	var _length = C.uint16_t(len(cookie))

	// TODO nil => real
	r := C.toxav_join_av_groupchat(this.toxcore, _fn, _data, _length, nil, nil)
	if int(r) == -1 {
		return int(r), errors.New("Join av group chat failed")
	}
	return int(r), nil
}
