package manipulator

import (
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"log"
	"midi_manipulator/pkg/config"
	"midi_manipulator/pkg/signals"
	"sync"
	"time"
)

type MidiManipulator struct {
	device      MidiDevice
	clickBuffer ClickBuffer
	holdDelta   time.Duration
	mutex       sync.Mutex
}

type MidiDevice struct {
	name  string
	ports MidiPorts
}

type MidiPorts struct {
	in *drivers.In
}

func (mm *MidiManipulator) setHoldDelta(delta float64) {
	mm.holdDelta = time.Duration(float64(time.Second) * delta)
}

func (mm *MidiManipulator) setDeviceName(deviceName string) {
	mm.device.name = deviceName
}

func (mm *MidiManipulator) getPortsByDeviceName(deviceName string) drivers.In {
	inPort, err := midi.FindInPort(deviceName)
	if err != nil {
		log.Println("Input port was not found")
		return nil
	}

	return inPort
}

func (mm *MidiManipulator) applyConfiguration(config config.MIDIConfig) {
	mm.setDeviceName(config.DeviceName)
	mm.setHoldDelta(config.HoldDelta)
	mm.clickBuffer.buffer = make(map[int]KeyContext)
}

func (mm *MidiManipulator) connectDevice(inPort drivers.In) {
	port, err := midi.InPort(inPort.Number())
	if err != nil {
		return
	}

	mm.device.ports.in = &port
}

func (mm *MidiManipulator) listen(signals chan<- core.Signal, shutdown <-chan bool) {
	stop, err := midi.ListenTo(*mm.device.ports.in, mm.getMidiMessage, midi.UseSysEx())

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	for {
		signal_sequence := mm.messageToSignal()
		select {
		case <-shutdown:
			stop()
			return
		default:
			for _, signal := range signal_sequence {
				mm.sendSignal(signals, signal)
			}
		}
	}
}

func (mm *MidiManipulator) sendSignal(signals chan<- core.Signal, signal core.Signal) {
	if signal != nil {
		log.Printf("%s, %d\n", signal.Code(), signal)
		signals <- signal
	}
}

func (mm *MidiManipulator) getMidiMessage(msg midi.Message, timestamps int32) {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	var bt []byte
	var channel, key, velocity uint8
	switch {
	case msg.GetAfterTouch(&channel, &velocity):
		return
	case msg.GetSysEx(&bt):
		return
	case msg.GetNoteOn(&channel, &key, &velocity):
		// NIL STATUS
		kctx := KeyContext{int(key), int(velocity), time.Now(),
			nil}
		mm.clickBuffer.setKeyContext(int(key), kctx)
	case msg.GetNoteOff(&channel, &key, &velocity):
		// NOTE RELEASED STATUS
		val, ok := mm.clickBuffer.getKeyContext(int(key))
		if ok {
			val.status = signals.NoteReleased{int(key), int(velocity)}
			mm.clickBuffer.setKeyContext(int(key), val)
		}
	case msg.GetControlChange(&channel, &key, &velocity):
		// CONTROL PUSHED STATUS
		kctx := KeyContext{int(key), int(velocity), time.Now(),
			signals.ControlPushed{int(key), int(velocity)}}
		mm.clickBuffer.setKeyContext(int(key), kctx)
	}
}

func (mm *MidiManipulator) messageToSignal() []core.Signal {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	var signal_sequence []core.Signal
	for _, kctx := range mm.clickBuffer.buffer {
		switch kctx.status.(type) {
		case nil:
			signal := signals.NotePushed{kctx.key,
				kctx.velocity}
			signal_sequence = append(signal_sequence, signal)
			// UPDATE KEY STATUS IN BUFFER
			kctx.status = signal
			mm.clickBuffer.setKeyContext(kctx.key, kctx)
		case signals.NotePushed:
			if time.Now().Sub(kctx.usedAt) >= mm.holdDelta {
				signal := signals.NoteHold{kctx.key,
					kctx.velocity}
				signal_sequence = append(signal_sequence, signal)
				// UPDATE KEY STATUS IN BUFFER
				kctx.status = signal
				mm.clickBuffer.setKeyContext(kctx.key, kctx)
			}
		case signals.NoteHold:
		case signals.NoteReleased:
			{
				signal := signals.NoteReleased{kctx.key,
					kctx.velocity}
				signal_sequence = append(signal_sequence, signal)
				// DELETE KEY FROM BUFFER
				delete(mm.clickBuffer.buffer, kctx.key)
			}
		case signals.ControlPushed:
			{
				signal := signals.ControlPushed{kctx.key,
					kctx.velocity}
				signal_sequence = append(signal_sequence, signal)
				// DELETE KEY FROM BUFFER
				delete(mm.clickBuffer.buffer, kctx.key)
			}
		}
	}
	return signal_sequence
}

func (mm *MidiManipulator) Run(config config.MIDIConfig, signals chan<- core.Signal, shutdown <-chan bool) {
	inPort := mm.getPortsByDeviceName(config.DeviceName)
	if inPort == nil {
		return
	}

	mm.applyConfiguration(config)
	mm.connectDevice(inPort)
	mm.listen(signals, shutdown)
}
