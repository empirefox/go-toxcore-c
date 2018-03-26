package tox

//#include <tox/tox.h>
import "C"
import (
	"time"
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

		case fn := <-t.chLoopRequest:
			fn()

		case <-pingTicker.C:
			t.doTcpPing_l()

		case <-t.stop:
			close(t.stopped)
			return
		}
	}
}

func (t *Tox) blockAv()   {}
func (t *Tox) unblockAv() {}
