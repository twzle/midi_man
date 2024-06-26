package midi

import (
	"midi_manipulator/pkg/model"
	"time"

	"git.miem.hse.ru/hubman/hubman-lib/core"
	"gitlab.com/gomidi/midi/v2"
	"go.uber.org/zap"
)

// Function writes signals to provided channel for device entity
func (md *MidiDevice) sendSignals(signals []core.Signal) {
	for _, signal := range signals {
		if signal != nil {
			md.logger.Debug("Received signal from MIDI device", zap.String("signal", signal.Code()), zap.Any("payload", signal))
			md.signals <- signal
		}
	}
}

// Function processing signals from MIDI-device
func (md *MidiDevice) processMidiMessage(msg midi.Message, _ int32) {
	md.mutex.Lock()
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
		velocity, valid := md.handleControls(int(key), int(velocity))
		if valid {
			kctx := KeyContext{key: key, velocity: uint8(velocity), usedAt: time.Now(),
				status: model.ControlPushed{Device: md.name, KeyCode: int(key), Value: int(velocity)}}
			md.clickBuffer.SetKeyContext(key, kctx)
		}
	}
	md.mutex.Unlock()

	md.sendSignals(md.messageToSignal())
}


// Function converts message to hubman-compatible signal
func (md *MidiDevice) messageToSignal() []core.Signal {
	md.mutex.Lock()
	defer md.mutex.Unlock()

	var signalSequence []core.Signal
	for _, kctx := range md.clickBuffer {
		switch kctx.status.(type) {
		case nil:
			signal := model.NotePushed{
				Device:    md.name,
				KeyCode:   int(kctx.key),
				Velocity:  int(kctx.velocity),
				Namespace: md.namespace,
			}
			signalSequence = append(signalSequence, signal)
			// UPDATE KEY STATUS IN BUFFER
			kctx.status = signal
		case model.NotePushed:
			if time.Since(kctx.usedAt) >= md.holdDelta {
				signal := model.NoteHold{
					Device:    md.name,
					KeyCode:   int(kctx.key),
					Velocity:  int(kctx.velocity),
					Namespace: md.namespace,
				}
				signalSequence = append(signalSequence, signal)
				// UPDATE KEY STATUS IN BUFFER
				kctx.status = signal
			}
		case model.NoteReleased:
			signal := model.NoteReleased{
				Device:    md.name,
				KeyCode:   int(kctx.key),
				Velocity:  int(kctx.velocity),
				Namespace: md.namespace,
			}
			signalSequence = append(signalSequence, signal)
			// DELETE KEY FROM BUFFER
			delete(md.clickBuffer, kctx.key)
		case model.NoteReleasedAfterHold:
			signal := model.NoteReleasedAfterHold{
				Device:    md.name,
				KeyCode:   int(kctx.key),
				Velocity:  int(kctx.velocity),
				Namespace: md.namespace,
			}
			signalSequence = append(signalSequence, signal)
			// DELETE KEY FROM BUFFER
			delete(md.clickBuffer, kctx.key)
		case model.ControlPushed:
			signal := model.ControlPushed{
				Device:    md.name,
				KeyCode:   int(kctx.key),
				Value:     int(kctx.velocity),
				Namespace: md.namespace,
			}
			signalSequence = append(signalSequence, signal)
			// DELETE KEY FROM BUFFER
			delete(md.clickBuffer, kctx.key)
		}
	}
	return signalSequence
}

// Function listening singals from single MIDI-device
func (md *MidiDevice) listen() {
	stopMidiListener := func() {}
	for {
		select {
		case <-md.stopListen:
			stopMidiListener()
			return
		case connected := <-md.reconnectedEvent:
			md.connected.Store(connected)
			if !connected {
				continue
			}
			var err error
			stopMidiListener, err = midi.ListenTo(md.ports.in, md.processMidiMessage, midi.UseSysEx())
			if err != nil {
				md.logger.Warn("error in init listen", zap.Error(err))
			}
		}
	}
}
