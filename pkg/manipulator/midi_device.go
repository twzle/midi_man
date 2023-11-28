package manipulator

import (
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"log"
	"midi_manipulator/pkg/config"
	"midi_manipulator/pkg/signals"
	"sync"
	"time"
)

type MidiDevice struct {
	name        string
	ports       MidiPort
	clickBuffer ClickBuffer
	holdDelta   time.Duration
	mutex       sync.Mutex
}

type MidiPort struct {
	in *drivers.In
}

func (md *MidiDevice) setHoldDelta(delta float64) {
	md.holdDelta = time.Duration(float64(time.Second) * delta)
}

func (md *MidiDevice) setDeviceName(deviceName string) {
	md.name = deviceName
}

func (md *MidiDevice) setClickBuffer() {
	md.clickBuffer = make(map[uint8]*KeyContext)
}

func (md *MidiDevice) applyConfiguration(config config.MidiConfig) {
	md.setDeviceName(config.DeviceName)
	md.setHoldDelta(config.HoldDelta)
	md.setClickBuffer()
}

func (md *MidiDevice) sendSignal(signals chan<- core.Signal, signal core.Signal) {
	if signal != nil {
		log.Println(signal.Code(), signal)
		signals <- signal
	}
}

func (md *MidiDevice) getMidiMessage(msg midi.Message, timestamps int32) {
	md.mutex.Lock()
	defer md.mutex.Unlock()
	var channel, key, velocity uint8
	switch {
	case msg.GetNoteOn(&channel, &key, &velocity):
		// NIL STATUS
		kctx := KeyContext{key, velocity, time.Now(),
			nil}
		md.clickBuffer.SetKeyContext(key, kctx)
	case msg.GetNoteOff(&channel, &key, &velocity):
		// NOTE RELEASED STATUS
		val, ok := md.clickBuffer.GetKeyContext(key)
		if ok {
			val.status = signals.NoteReleased{md.name, int(key), int(velocity)}
			//mm.clickBuffer.SetKeyContext(key, val)
		}
	case msg.GetControlChange(&channel, &key, &velocity):
		// CONTROL PUSHED STATUS
		kctx := KeyContext{key, velocity, time.Now(),
			signals.ControlPushed{md.name, int(key), int(velocity)}}
		md.clickBuffer.SetKeyContext(key, kctx)
	}
}

func (md *MidiDevice) messageToSignal() []core.Signal {
	md.mutex.Lock()
	defer md.mutex.Unlock()
	var signalSequence []core.Signal
	for _, kctx := range md.clickBuffer {
		switch kctx.status.(type) {
		case nil:
			signal := signals.NotePushed{md.name, int(kctx.key), int(kctx.velocity)}
			signalSequence = append(signalSequence, signal)
			// UPDATE KEY STATUS IN BUFFER
			kctx.status = signal
		case signals.NotePushed:
			if time.Now().Sub(kctx.usedAt) >= md.holdDelta {
				signal := signals.NoteHold{md.name, int(kctx.key), int(kctx.velocity)}
				signalSequence = append(signalSequence, signal)
				// UPDATE KEY STATUS IN BUFFER
				kctx.status = signal
			}
		case signals.NoteReleased:
			signal := signals.NoteReleased{md.name, int(kctx.key),
				int(kctx.velocity)}
			signalSequence = append(signalSequence, signal)
			// DELETE KEY FROM BUFFER
			delete(md.clickBuffer, kctx.key)
		case signals.ControlPushed:
			signal := signals.ControlPushed{md.name, int(kctx.key),
				int(kctx.velocity)}
			signalSequence = append(signalSequence, signal)
			// DELETE KEY FROM BUFFER
			delete(md.clickBuffer, kctx.key)
		}
	}
	return signalSequence
}
