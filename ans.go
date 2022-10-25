package webrtcvad

import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// #cgo CFLAGS: -I.
// #include <stdio.h>
// #include <stdlib.h>
// #include <math.h>
// #include <stdint.h>
// #include "cgo/ns/noise_suppression.h"
// #include "cgo/ns/dr_mp3.h"
// #include "cgo/ns/timing.h"
// #include "cgo/ns/dr_wav.h"
import "C"

func NewAns(sampleRate int) (*Ans, error) {
	var inst *C.NsHandle
	inst = C.WebRtcNs_Create()
	// ret := C.WebRtcVad_Create(&inst)
	// if ret != 0 {
	// 	return nil, errors.New("failed to create VAD")
	// }

	vad := &Ans{inst}
	runtime.SetFinalizer(vad, freeAns)

	ret := C.WebRtcNs_Init(inst, C.uint(sampleRate))
	if ret != 0 {
		return nil, fmt.Errorf("default mode could not be set:%v", ret)
	}

	return vad, nil
}

func freeAns(ans *Ans) {
	C.WebRtcNs_Free(ans.inst)
}

type Ans struct {
	inst *C.NsHandle
}

func (v *Ans) SetMode(mode int) error {
	ret := C.WebRtcNs_set_policy(v.inst, C.int(mode))
	if ret != 0 {
		return errors.New("mode could not be set")
	}
	return nil
}

func (v *Ans) Process(fs int, audioFrame []byte) (activeVoice bool, err error) {
	if len(audioFrame)%2 != 0 {
		return false, errors.New("audio frames must be 16bit little endian unsigned integers")
	}

	audioFramePtr := (*C.int16_t)(unsafe.Pointer(&audioFrame[0]))
	// frameLen := C.ulong(len(audioFrame) / 2)

	outFrame := make([]byte, len(audioFrame))
	outFramePtr := (*C.int16_t)(unsafe.Pointer(&outFrame[0]))
	C.WebRtcNs_Process(v.inst, &audioFramePtr, C.ulong(1), &outFramePtr)

	return false, nil
}
