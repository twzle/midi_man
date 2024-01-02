package midi

import (
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"gitlab.com/gomidi/midi/v2"
	"log"
	"midi_manipulator/pkg/model"
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
			switch val.status.(type) {
			case model.NotePushed:
				val.status = model.NoteReleased{Device: md.name, KeyCode: int(key), Velocity: int(velocity)}
			case model.NoteHold:
				val.status = model.NoteReleasedAfterHold{Device: md.name, KeyCode: int(key), Velocity: int(velocity)}
			}
		}
	case msg.GetControlChange(&channel, &key, &velocity):
		// CONTROL PUSHED STATUS
		kctx := KeyContext{key: key, velocity: velocity, usedAt: time.Now(),
			status: model.ControlPushed{Device: md.name, KeyCode: int(key), Value: int(velocity)}}
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
			signal := model.NotePushed{Device: md.name, KeyCode: int(kctx.key), Velocity: int(kctx.velocity)}
			signalSequence = append(signalSequence, signal)
			// UPDATE KEY STATUS IN BUFFER
			kctx.status = signal
		case model.NotePushed:
			if time.Now().Sub(kctx.usedAt) >= md.holdDelta {
				signal := model.NoteHold{Device: md.name, KeyCode: int(kctx.key), Velocity: int(kctx.velocity)}
				signalSequence = append(signalSequence, signal)
				// UPDATE KEY STATUS IN BUFFER
				kctx.status = signal
			}
		case model.NoteReleased:
			signal := model.NoteReleased{Device: md.name, KeyCode: int(kctx.key),
				Velocity: int(kctx.velocity)}
			signalSequence = append(signalSequence, signal)
			// DELETE KEY FROM BUFFER
			delete(md.clickBuffer, kctx.key)
		case model.NoteReleasedAfterHold:
			signal := model.NoteReleasedAfterHold{Device: md.name, KeyCode: int(kctx.key),
				Velocity: int(kctx.velocity)}
			signalSequence = append(signalSequence, signal)
			// DELETE KEY FROM BUFFER
			delete(md.clickBuffer, kctx.key)
		case model.ControlPushed:
			signal := model.ControlPushed{Device: md.name, KeyCode: int(kctx.key),
				Value: int(kctx.velocity)}
			signalSequence = append(signalSequence, signal)
			// DELETE KEY FROM BUFFER
			delete(md.clickBuffer, kctx.key)
		}
	}
	return signalSequence
}

func (md *MidiDevice) listen(signals chan<- core.Signal) {
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
		default:
			for _, signal := range signalSequence {
				md.sendSignal(signals, signal)
			}
		}
	}
}
