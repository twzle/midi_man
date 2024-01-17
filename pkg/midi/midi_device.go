package midi

import (
	"fmt"
	"git.miem.hse.ru/hubman/hubman-lib/core"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"log"
	"midi_manipulator/pkg/backlight"
	"midi_manipulator/pkg/config"
	"midi_manipulator/pkg/model"
	"sync"
	"time"
)

type MidiDevice struct {
	name        string
	active      bool
	ports       MidiPorts
	clickBuffer ClickBuffer
	holdDelta   time.Duration
	mutex       sync.Mutex
	stop        chan bool
	namespace   string
	controls    map[byte]*Control
}

type MidiPorts struct {
	in  *drivers.In
	out *drivers.Out
}

func (md *MidiDevice) GetAlias() string {
	return md.name
}

func (md *MidiDevice) ExecuteCommand(command model.MidiCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) error {
	switch v := command.(type) {
	case model.TurnLightOnCommand:
		md.turnLightOn(command.(model.TurnLightOnCommand), backlightConfig)
	case model.TurnLightOffCommand:
		md.turnLightOff(command.(model.TurnLightOffCommand), backlightConfig)
	case model.SingleBlinkCommand:
		md.singleBlink(command.(model.SingleBlinkCommand), backlightConfig)
	case model.SingleReversedBlinkCommand:
		md.singleReversedBlink(command.(model.SingleReversedBlinkCommand), backlightConfig)
	default:
		fmt.Printf("Unknown command with type: \"%T\"\n", v)
	}
	return nil
}

func (md *MidiDevice) StopDevice() error {
	md.stop <- true
	return nil
}

func (md *MidiDevice) RunDevice(signals chan<- core.Signal, backlightConfig *backlight.DecodedDeviceBacklightConfig) error {
	go md.startupIllumination(backlightConfig)
	go md.listen(signals)
	return nil
}

func (md *MidiDevice) connectDevice() error {
	var err error
	in_err := md.connectInPort()
	out_err := md.connectOutPort()

	if in_err != nil || out_err != nil {
		err = fmt.Errorf("connection of device \"{%s}\" failed", md.name)
	}
	return err
}

func (md *MidiDevice) connectOutPort() error {
	port, err := midi.FindOutPort(md.name)
	if err != nil {
		log.Printf("Output port named {%s} was not found\n", md.name)
		return err
	}

	port, err = midi.OutPort(port.Number())
	if err != nil {
		log.Printf("Output port #{%d} was not found\n", port.Number())
		return err
	}

	md.ports.out = &port
	return nil
}

func (md *MidiDevice) connectInPort() error {
	port, err := midi.FindInPort(md.name)
	if err != nil {
		log.Printf("Input port named {%s} was not found", md.name)
		return err
	}

	port, err = midi.InPort(port.Number())
	if err != nil {
		log.Printf("Input port #{%d} was not found\n", port.Number())
		return err
	}

	md.ports.in = &port
	return nil
}

func (md *MidiDevice) applyConfiguration(deviceConfig config.DeviceConfig) {
	md.name = deviceConfig.DeviceName
	md.active = deviceConfig.Active
	md.holdDelta = time.Duration(float64(time.Millisecond) * deviceConfig.HoldDelta)
	md.clickBuffer = make(map[uint8]*KeyContext)
	md.stop = make(chan bool, 1)
	md.namespace = deviceConfig.Namespace
	md.controls = make(map[byte]*Control)
	md.applyControls(deviceConfig.Controls)
}

func (md *MidiDevice) applyControls(controls config.Controls) {
	for _, controlKey := range controls.Keys {
		control := Control{
			Key: controlKey, 
			Rotate: controls.Rotate, 
			ValueRange: controls.ValueRange, 
			InitialValue: controls.InitialValue, 
			DecrementTrigger: controls.Triggers.Decrement, 
			IncrementTrigger: controls.Triggers.Increment,
		}
		md.controls[controlKey] = &control
	}
}

func (md *MidiDevice) updateConfiguration(config config.DeviceConfig) {
	md.active = config.Active
	md.holdDelta = time.Duration(float64(time.Millisecond) * config.HoldDelta)
}

func NewDevice(deviceConfig config.DeviceConfig) (*MidiDevice, error) {
	midiDevice := MidiDevice{}
	midiDevice.applyConfiguration(deviceConfig)

	err := midiDevice.connectDevice()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &midiDevice, nil
}
