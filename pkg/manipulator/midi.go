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
	device    MidiDevice
	keyCtx    KeyContext
	holdDelta time.Duration
	mutex     sync.Mutex
}

type MidiDevice struct {
	name  string
	ports MidiPorts
}

type MidiPorts struct {
	in  drivers.In
	out drivers.Out
}

func (mm *MidiManipulator) setHoldDelta(delta float64) {
	mm.holdDelta = time.Duration(float64(time.Second) * delta)
}

func (mm *MidiManipulator) setDeviceName(deviceName string) {
	mm.device.name = deviceName
}

func (mm *MidiManipulator) getPortsByDeviceName(deviceName string) (drivers.In, drivers.Out) {
	inPort, err := midi.FindInPort(deviceName)
	if err != nil {
		log.Println("Input port was not found")
		return nil, nil
	}

	outPort, err := midi.FindOutPort(deviceName)
	if err != nil {
		log.Println("Output port was not found")
		return nil, nil
	}

	return inPort, outPort
}

func (mm *MidiManipulator) applyConfiguration(config config.MIDIConfig) {
	mm.setDeviceName(config.DeviceName)
	mm.setHoldDelta(config.HoldDelta)
}

func (mm *MidiManipulator) connectDevice(inPort drivers.In) {
	mm.device.ports.in, _ = midi.InPort(inPort.Number())
}

func (mm *MidiManipulator) listen(signals chan<- core.Signal, shutdown <-chan bool) {
	stop, err := midi.ListenTo(mm.device.ports.in, mm.getMidiMessage, midi.UseSysEx())

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}

	for {
		signal := mm.messageToSignal()
		select {
		case <-shutdown:
			stop()
			return
		default:
			mm.sendSignal(signals, signal)
		}
	}
}

func (mm *MidiManipulator) sendSignal(signals chan<- core.Signal, signal core.Signal) {
	if signal != nil {
		log.Printf("%s, %f\n", signal.Code(), signal)
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
		// Бан одновременного нажатия множества клавиш
		if !mm.keyCtx.isPreviousKeyActive() {
			mm.keyCtx.setCurrentKey(MidiKey{float64(key), float64(velocity), time.Now(),
				signals.NotePushed{float64(key), float64(velocity)}})
		}
	case msg.GetNoteOff(&channel, &key, &velocity):
		if float64(key) == mm.keyCtx.currentKey.getKeyCode() {
			mm.keyCtx.currentKey.setStatus(signals.NoteReleased{float64(key), float64(velocity)})
		}
	case msg.GetControlChange(&channel, &key, &velocity):
		// Бан одновременного нажатия множества клавиш
		if !mm.keyCtx.isPreviousKeyActive() {
			mm.keyCtx.setCurrentKey(MidiKey{float64(key), float64(velocity), time.Now(),
				signals.ControlPushed{float64(key), float64(velocity)}})
		}
	}
}

func (mm *MidiManipulator) messageToSignal() core.Signal {
	mm.mutex.Lock()
	defer mm.mutex.Unlock()
	var signal core.Signal
	if !mm.keyCtx.compareKeys() {
		mm.keyCtx.setPreviousKey(MidiKey{mm.keyCtx.currentKey.getKeyCode(),
			mm.keyCtx.currentKey.getVelocity(),
			mm.keyCtx.currentKey.getUsedAt(),
			nil})
	}
	switch mm.keyCtx.currentKey.status.(type) {
	case signals.NotePushed:
		if time.Now().Sub(mm.keyCtx.currentKey.getUsedAt()) >= mm.holdDelta {
			mm.keyCtx.currentKey.setStatus(signals.NoteHold{mm.keyCtx.currentKey.getKeyCode(),
				mm.keyCtx.currentKey.getVelocity()})
			return nil
		}
		if !mm.keyCtx.compareStatuses() {
			mm.keyCtx.previousKey.setStatus(mm.keyCtx.currentKey.getStatus())
			signal = signals.NotePushed{float64(mm.keyCtx.currentKey.getKeyCode()),
				mm.keyCtx.currentKey.getVelocity()}
			return signal
		}
	case signals.NoteHold:
		if !mm.keyCtx.compareStatuses() {
			mm.keyCtx.previousKey.setStatus(mm.keyCtx.currentKey.getStatus())
			signal = signals.NoteHold{mm.keyCtx.currentKey.getKeyCode(),
				mm.keyCtx.currentKey.getVelocity()}
			return signal
		}
	case signals.NoteReleased:
		if !mm.keyCtx.compareStatuses() {
			mm.keyCtx.previousKey.setStatus(mm.keyCtx.currentKey.getStatus())
			signal = signals.NoteReleased{mm.keyCtx.currentKey.getKeyCode(),
				mm.keyCtx.currentKey.getVelocity()}
			return signal
		}
	case signals.ControlPushed:
		if !mm.keyCtx.compareStatuses() {
			mm.keyCtx.previousKey.setStatus(mm.keyCtx.currentKey.getStatus())
			signal = signals.ControlPushed{mm.keyCtx.currentKey.getKeyCode(),
				mm.keyCtx.currentKey.getVelocity()}
			return signal
		}
	}
	return nil
}

func (mm *MidiManipulator) Run(config config.MIDIConfig, signals chan<- core.Signal, shutdown <-chan bool) {
	defer midi.CloseDriver()
	inPort, outPort := mm.getPortsByDeviceName(config.DeviceName)
	if inPort == nil || outPort == nil {
		return
	}

	mm.applyConfiguration(config)
	mm.connectDevice(inPort)
	defer mm.device.ports.in.Close()
	defer mm.device.ports.out.Close()
	mm.listen(signals, shutdown)
}
