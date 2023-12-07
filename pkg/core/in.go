package core

import (
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"gitlab.com/gomidi/midi/v2"
	"log"
	"midi_manipulator/pkg/utils"
	"time"
)

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
			val.status = utils.NoteReleased{md.name, int(key), int(velocity)}
			//mm.clickBuffer.SetKeyContext(key, val)
		}
	case msg.GetControlChange(&channel, &key, &velocity):
		// CONTROL PUSHED STATUS
		kctx := KeyContext{key, velocity, time.Now(),
			utils.ControlPushed{md.name, int(key), int(velocity)}}
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
			signal := utils.NotePushed{md.name, int(kctx.key), int(kctx.velocity)}
			signalSequence = append(signalSequence, signal)
			// UPDATE KEY STATUS IN BUFFER
			kctx.status = signal
		case utils.NotePushed:
			if time.Now().Sub(kctx.usedAt) >= md.holdDelta {
				signal := utils.NoteHold{md.name, int(kctx.key), int(kctx.velocity)}
				signalSequence = append(signalSequence, signal)
				// UPDATE KEY STATUS IN BUFFER
				kctx.status = signal
			}
		case utils.NoteReleased:
			signal := utils.NoteReleased{md.name, int(kctx.key),
				int(kctx.velocity)}
			signalSequence = append(signalSequence, signal)
			// DELETE KEY FROM BUFFER
			delete(md.clickBuffer, kctx.key)
		case utils.ControlPushed:
			signal := utils.ControlPushed{md.name, int(kctx.key),
				int(kctx.velocity)}
			signalSequence = append(signalSequence, signal)
			// DELETE KEY FROM BUFFER
			delete(md.clickBuffer, kctx.key)
		}
	}
	return signalSequence
}

func (md *MidiDevice) listen(signals chan<- core.Signal, shutdown <-chan bool) {
	stop, err := midi.ListenTo(*md.ports.in, md.getMidiMessage, midi.UseSysEx())

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	for {
		signalSequence := md.messageToSignal()
		select {
		case <-md.stop:
			stop()
			return
		case <-shutdown:
			stop()
			return
		default:
			for _, signal := range signalSequence {
				md.sendSignal(signals, signal)
			}
		}
	}
}
