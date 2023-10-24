package main

import (
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"sync"
	"time"
)

// TODO: Вынести в конфиги
var delta = time.Second
var ctx = KeyContext{}
var mutex sync.Mutex

func connectDevice() (in drivers.In, out drivers.Out) {
	inPort, err := midi.FindInPort("MPD226")
	if err != nil {
		fmt.Printf("can't find")
		return
	}

	outPort, err := midi.FindOutPort("MPD226")
	if err != nil {
		fmt.Printf("can't find")
		return
	}

	in, _ = midi.InPort(inPort.Number())
	out, _ = midi.OutPort(outPort.Number())

	return in, out
}

func startupIllumination(out drivers.Out) {
	// AKAI MPD226 DIV'S
	for i := 0; i < 4; i++ {
		msg := midi.Message{177, byte(i), 127}
		time.Sleep(time.Millisecond * 50)
		out.Send(msg)
	}

	for i := 0; i < 4; i++ {
		msg := midi.Message{177, byte(i), 0}
		time.Sleep(time.Millisecond * 50)
		out.Send(msg)
	}

	// AKAI MPD226 PADS
	for i := 60; i < 88; i++ {
		msg := midi.Message{145, byte(i), 4}
		time.Sleep(time.Millisecond * 50)
		out.Send(msg)
	}

	for i := 60; i < 88; i++ {
		msg := midi.Message{129, byte(i), 2}
		time.Sleep(time.Millisecond * 50)
		out.Send(msg)
	}
}

func listen(in drivers.In, signals chan<- core.Signal, done <-chan bool) {
	_, err := midi.ListenTo(in, getMidiMessage, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	for {
		signal := messageToSignal()
		select {
		case <-done:
			return
		default:
			sendSignal(signals, signal)
		}
	}
}

func sendSignal(signals chan<- core.Signal, signal core.Signal) {
	if signal != nil {
		fmt.Println(signal.Code(), signal)
		signals <- signal
	}
}

func getMidiMessage(msg midi.Message, timestamps int32) {
	mutex.Lock()
	defer mutex.Unlock()
	var bt []byte
	var channel, key, velocity uint8
	switch {
	case msg.GetAfterTouch(&channel, &velocity):
		return
	case msg.GetSysEx(&bt):
		return
	case msg.GetNoteOn(&channel, &key, &velocity):
		// Бан одновременного нажатия множества клавиш
		if !ctx.getPreviousKey().isActive() {
			ctx.setCurrentKey(MidiKey{int(key), int(velocity), time.Now(),
				NotePushed{int(key), int(velocity)}})
			_ = ctx.getCurrentKey().getStatus()
		}
	case msg.GetNoteOff(&channel, &key, &velocity):
		if int(key) == ctx.getCurrentKey().key {
			ctx.getCurrentKey().setStatus(NoteReleased{int(key), int(velocity)})
		}
		_ = ctx.getCurrentKey().getStatus()
	case msg.GetControlChange(&channel, &key, &velocity):
		// Бан одновременного нажатия множества клавиш
		if !ctx.getPreviousKey().isActive() {
			ctx.setCurrentKey(MidiKey{int(key), int(velocity), time.Now(),
				ControlPushed{int(key), int(velocity)}})
		}
	default:
		fmt.Println(msg)
	}
}

func messageToSignal() core.Signal {
	mutex.Lock()
	defer mutex.Unlock()
	var signal core.Signal
	if !ctx.compareKeys() {
		ctx.setPreviousKey(MidiKey{ctx.currentKey.key,
			ctx.currentKey.velocity,
			ctx.currentKey.usedAt,
			nil})
	}
	switch ctx.currentKey.status.(type) {
	case NotePushed:
		if time.Now().Sub(ctx.currentKey.usedAt) >= delta {
			ctx.getCurrentKey().setStatus(NoteHold{ctx.currentKey.key, ctx.currentKey.velocity})
			return nil
		}
		if !ctx.compareStatuses() {
			ctx.getPreviousKey().setStatus(ctx.currentKey.status)
			signal = NotePushed{ctx.currentKey.key,
				ctx.currentKey.velocity}
			return signal
		}
	case NoteHold:
		if !ctx.compareStatuses() {
			ctx.getPreviousKey().setStatus(ctx.currentKey.status)
			signal = NoteHold{ctx.currentKey.key,
				ctx.currentKey.velocity}
			return signal
		}
	case NoteReleased:
		if !ctx.compareStatuses() {
			ctx.getPreviousKey().setStatus(ctx.currentKey.status)
			signal = NoteReleased{ctx.currentKey.key,
				ctx.currentKey.velocity}
			return signal
		}
	case ControlPushed:
		if !ctx.compareStatuses() {
			ctx.getPreviousKey().setStatus(ctx.currentKey.status)
			signal = ControlPushed{ctx.currentKey.key,
				ctx.currentKey.velocity}
			return signal
		}
	}
	return nil
}

func Run(signals chan<- core.Signal, done <-chan bool) {
	defer midi.CloseDriver()
	in, out := connectDevice()
	defer in.Close()
	defer out.Close()
	startupIllumination(out)
	listen(in, signals, done)
}
