package tox

/*
#include <tox/tox.h>
extern void toxCallbackLog(Tox*, TOX_LOG_LEVEL, char*, uint32_t, char*, char*);
*/
import "C"
import (
	"time"
	"unsafe"

	"github.com/TokTok/go-toxcore-c/toxenums"
)

type ToxOptions struct {
	Ipv6_enabled            bool
	Udp_enabled             bool
	Proxy_type              toxenums.TOX_PROXY_TYPE
	Proxy_host              string
	Proxy_port              uint16
	Savedata_type           toxenums.TOX_SAVEDATA_TYPE
	Savedata_data           []byte
	Tcp_port                uint16
	Local_discovery_enabled bool
	Start_port              uint16
	End_port                uint16
	Hole_punching_enabled   bool
	LogCallback             func(t *Tox, level toxenums.TOX_LOG_LEVEL, file string, line uint32, fname string, msg string)

	// additions
	NospamIfSecretType      uint32
	Decrypt                 func(ciphertext []byte) ([]byte, error)
	ProxyToNoneIfErr        bool
	AutoTcpPortIfErr        bool
	DisableTcpPortIfAutoErr bool
	PingUnit                time.Duration
}

func NewToxOptions() *ToxOptions {
	toxopts := C.tox_options_new(nil)
	defer C.tox_options_free(toxopts)

	opts := new(ToxOptions)
	opts.Ipv6_enabled = bool(C.tox_options_get_ipv6_enabled(toxopts))
	opts.Udp_enabled = bool(C.tox_options_get_udp_enabled(toxopts))
	opts.Proxy_type = toxenums.TOX_PROXY_TYPE(C.tox_options_get_proxy_type(toxopts))
	opts.Proxy_port = uint16(C.tox_options_get_proxy_port(toxopts))
	opts.Tcp_port = uint16(C.tox_options_get_tcp_port(toxopts))
	opts.Local_discovery_enabled = bool(C.tox_options_get_local_discovery_enabled(toxopts))
	opts.Start_port = uint16(C.tox_options_get_start_port(toxopts))
	opts.End_port = uint16(C.tox_options_get_end_port(toxopts))
	opts.Hole_punching_enabled = bool(C.tox_options_get_hole_punching_enabled(toxopts))

	return opts
}

func (this *ToxOptions) toCToxOptions() *C.struct_Tox_Options {
	toxopts := C.tox_options_new(nil)
	C.tox_options_default(toxopts)
	C.tox_options_set_ipv6_enabled(toxopts, (C._Bool)(this.Ipv6_enabled))
	C.tox_options_set_udp_enabled(toxopts, (C._Bool)(this.Udp_enabled))

	if this.Savedata_data != nil {
		C.tox_options_set_savedata_data(toxopts, (*C.uint8_t)(&this.Savedata_data[0]), C.size_t(len(this.Savedata_data)))
		C.tox_options_set_savedata_type(toxopts, C.TOX_SAVEDATA_TYPE(this.Savedata_type))
	}
	C.tox_options_set_tcp_port(toxopts, (C.uint16_t)(this.Tcp_port))

	C.tox_options_set_proxy_type(toxopts, C.TOX_PROXY_TYPE(this.Proxy_type))
	C.tox_options_set_proxy_port(toxopts, C.uint16_t(this.Proxy_port))
	if len(this.Proxy_host) > 0 {
		C.tox_options_set_proxy_host(toxopts, C.CString(this.Proxy_host))
	}

	C.tox_options_set_local_discovery_enabled(toxopts, C._Bool(this.Local_discovery_enabled))
	C.tox_options_set_start_port(toxopts, C.uint16_t(this.Start_port))
	C.tox_options_set_end_port(toxopts, C.uint16_t(this.End_port))
	C.tox_options_set_hole_punching_enabled(toxopts, C._Bool(this.Hole_punching_enabled))

	C.tox_options_set_log_callback(toxopts, (*C.tox_log_cb)((unsafe.Pointer)(C.toxCallbackLog)))

	return toxopts
}

//export toxCallbackLog
func toxCallbackLog(ctox *C.Tox, level C.TOX_LOG_LEVEL, file *C.char, line C.uint32_t, fname *C.char, msg *C.char) {
	t := cbUserDatas.get(ctox)
	if t != nil && t.opts != nil && t.opts.LogCallback != nil {
		t.opts.LogCallback(t, toxenums.TOX_LOG_LEVEL(level), C.GoString(file), uint32(line), C.GoString(fname), C.GoString(msg))
	}
}
