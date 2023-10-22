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

type MidiKey struct {
	key      int
	velocity int
	usedAt   time.Time
	status   core.Signal
}

var delta = time.Second
var currentUsedKey = MidiKey{-1, -1, time.Time{}, nil}
var previousUsedKey = MidiKey{-1, -1, time.Time{}, nil}
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

func listen(in drivers.In, signals chan<- core.Signal) {
	_, err := midi.ListenTo(in, getMidiMessage, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	for {
		messageToSignal(signals)
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
		currentUsedKey = MidiKey{int(key), int(velocity), time.Now(),
			NotePushed{int(key), int(velocity)}}
		_ = currentUsedKey.status
	case msg.GetNoteOff(&channel, &key, &velocity):
		if int(key) == currentUsedKey.key {
			currentUsedKey.status = NoteReleased{int(key), int(velocity)}
		}
		_ = currentUsedKey.status
	case msg.GetControlChange(&channel, &key, &velocity):
		currentUsedKey = MidiKey{int(key), int(velocity), time.Now(),
			ControlPushed{int(key), int(velocity)}}
	default:
		fmt.Println(msg)
	}
}

func messageToSignal(signals chan<- core.Signal) {
	mutex.Lock()
	defer mutex.Unlock()
	if previousUsedKey.key != currentUsedKey.key {
		previousUsedKey.status = nil
		previousUsedKey.key = currentUsedKey.key
	}
	switch currentUsedKey.status.(type) {
	case NotePushed:
		if time.Now().Sub(currentUsedKey.usedAt) >= delta {
			currentUsedKey.status = NoteHold{currentUsedKey.key,
				currentUsedKey.velocity}
			return
		}
		signal := currentUsedKey.status
		if previousUsedKey.status != currentUsedKey.status {
			signals <- signal
			previousUsedKey.status = currentUsedKey.status
		}
	case NoteHold:
		signal := NoteHold{currentUsedKey.key,
			currentUsedKey.velocity}
		if previousUsedKey.status != currentUsedKey.status {
			signals <- signal
			previousUsedKey.status = currentUsedKey.status
		}
	case NoteReleased:
		signal := NoteReleased{currentUsedKey.key,
			currentUsedKey.velocity}
		if previousUsedKey.status != currentUsedKey.status {
			signals <- signal
			previousUsedKey.status = currentUsedKey.status
		}
	case ControlPushed:
		signal := ControlPushed{currentUsedKey.key,
			currentUsedKey.velocity}
		if previousUsedKey.status != currentUsedKey.status {
			signals <- signal
			previousUsedKey.status = currentUsedKey.status
		}
	}
}

func Run(signals chan<- core.Signal) {
	defer midi.CloseDriver()
	in, out := connectDevice()
	defer in.Close()
	defer out.Close()
	startupIllumination(out)
	listen(in, signals)
}
