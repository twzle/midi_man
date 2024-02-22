package midi

import (
	"errors"
	"fmt"
	"midi_manipulator/pkg/backlight"
	"midi_manipulator/pkg/config"
	"midi_manipulator/pkg/model"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.miem.hse.ru/hubman/hubman-lib/core"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"go.uber.org/zap"
)

type MidiDevice struct {
	name              string
	active            bool
	ports             MidiPorts
	clickBuffer       ClickBuffer
	holdDelta         time.Duration
	startupDelay      time.Duration
	reconnectInterval time.Duration
	mutex             sync.Mutex
	stopReconnect     chan struct{}
	stopListen        chan struct{}
	namespace         string
	connected         atomic.Bool
	controls          map[byte]*Control
	signals           chan<- core.Signal
	logger            *zap.Logger
	conf              config.DeviceConfig
	reconnectedEvent  chan bool
}

type MidiPorts struct {
	in  *drivers.In
	out *drivers.Out
}

func (md *MidiDevice) GetAlias() string {
	return md.name
}

func (md *MidiDevice) ExecuteCommand(command model.MidiCommand, backlightConfig *backlight.DecodedDeviceBacklightConfig) error {
	md.mutex.Lock()
	defer md.mutex.Unlock()

	switch v := command.(type) {
	case model.TurnLightOnCommand:
		md.turnLightOn(command.(model.TurnLightOnCommand), backlightConfig)
	case model.TurnLightOffCommand:
		md.turnLightOff(command.(model.TurnLightOffCommand), backlightConfig)
	case model.SingleBlinkCommand:
		md.singleBlink(command.(model.SingleBlinkCommand), backlightConfig)
	case model.SingleReversedBlinkCommand:
		md.singleReversedBlink(command.(model.SingleReversedBlinkCommand), backlightConfig)
	case model.SetActiveNamespaceCommand:
		md.setActiveNamespace(command.(model.SetActiveNamespaceCommand), backlightConfig)
	default:
		md.logger.Warn("Unknown command", zap.Any("command", v))
	}
	return nil
}

func (md *MidiDevice) Stop() {
	md.stopListen <- struct{}{}
	close(md.stopListen)
	md.stopReconnect <- struct{}{}
	close(md.stopReconnect)
}

func (md *MidiDevice) RunDevice(backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	time.Sleep(md.startupDelay)
	go md.reconnect(backlightConfig)
	go md.listen()
}

func (md *MidiDevice) initConnection(backlightConfig *backlight.DecodedDeviceBacklightConfig) error {
	md.mutex.Lock()
	defer md.mutex.Unlock()

	if err := md.connectDevice(); err != nil {
		return err
	}
	md.startupIllumination(backlightConfig)
	md.clickBuffer = make(map[uint8]*KeyContext)
	md.applyControls(md.conf.Controls)
	return nil
}

func (md *MidiDevice) reconnect(backlightConfig *backlight.DecodedDeviceBacklightConfig) {
	ticker := time.NewTicker(md.reconnectInterval)
	for {
		select {
		case <-md.stopReconnect:
			ticker.Stop()
			return
		case <-ticker.C:
			connected := HasDeviceWithName(md.name, midi.GetInPorts()) &&
				HasDeviceWithName(md.name, midi.GetOutPorts())
			if md.connected.Load() {
				if !connected {
					md.logger.Warn("Device disconnected")
					md.reconnectedEvent <- connected
				}
			} else {
				if connected {
					err := md.initConnection(backlightConfig)
					if err != nil {
						md.logger.Warn("Unable to connect device", zap.Error(err))
					} else {
						md.reconnectedEvent <- connected
					}
				} else {
					md.logger.Debug("No hardware connection to device")
				}
			}
		}
	}
}

func (md *MidiDevice) connectDevice() error {
	if inErr := md.connectInPort(); inErr != nil {
		return fmt.Errorf("connection of device failed, inPort:\"{%w}\"", inErr)
	}
	if outErr := md.connectOutPort(); outErr != nil {
		return fmt.Errorf("connection of device failed, outPort:\"{%w}\"", outErr)
	}
	return nil
}

func (md *MidiDevice) connectOutPort() error {
	portNum := -1
	for _, inPort := range midi.GetOutPorts() {
		md.logger.Debug("Found out midi port", zap.String("name", inPort.String()), zap.Int("num", inPort.Number()))
		if strings.Contains(inPort.String(), md.name) || strings.Contains(md.name, inPort.String()) {
			md.logger.Info("Matched midi out port", zap.String("name", inPort.String()), zap.Int("num", inPort.Number()))
			portNum = inPort.Number()
			break
		}
	}
	if portNum == -1 {
		return errors.New("not found midi out port")
	}
	port, err := midi.OutPort(portNum)
	if err != nil {
		return err
	}

	md.ports.out = &port
	return nil
}

func (md *MidiDevice) connectInPort() error {
	portNum := -1
	for _, inPort := range midi.GetInPorts() {
		md.logger.Debug("Found in midi port", zap.String("name", inPort.String()), zap.Int("num", inPort.Number()))
		if strings.Contains(inPort.String(), md.name) || strings.Contains(md.name, inPort.String()) {
			portNum = inPort.Number()
			md.logger.Info("Matched midi in port", zap.String("name", inPort.String()), zap.Int("num", inPort.Number()))
			break
		}
	}
	if portNum == -1 {
		return errors.New("not found midi in port")
	}
	port, err := midi.InPort(portNum)
	if err != nil {
		return err
	}

	md.ports.in = &port
	return nil
}

func (md *MidiDevice) applyConfiguration(
	deviceConfig config.DeviceConfig,
	signals chan<- core.Signal,
	logger *zap.Logger,
) {
	md.conf = deviceConfig
	md.name = deviceConfig.DeviceName
	md.active = deviceConfig.Active
	md.holdDelta = time.Duration(deviceConfig.HoldDelta) * time.Millisecond
	md.startupDelay = time.Duration(deviceConfig.StartupDelay) * time.Millisecond
	md.reconnectInterval = time.Duration(deviceConfig.ReconnectInterval) * time.Millisecond
	md.clickBuffer = make(map[uint8]*KeyContext)
	md.stopListen = make(chan struct{})
	md.stopReconnect = make(chan struct{})
	md.reconnectedEvent = make(chan bool)
	md.namespace = deviceConfig.Namespace
	md.signals = signals
	md.logger = logger.With(zap.String("alias", md.name))
	md.applyControls(deviceConfig.Controls)
}

func (md *MidiDevice) applyControls(controls config.Controls) {
	md.controls = make(map[byte]*Control)
	for _, controlKey := range controls.Keys {
		control := Control{
			Key:              controlKey,
			Rotate:           controls.Rotate,
			ValueRange:       controls.ValueRange,
			InitialValue:     controls.InitialValue,
			DecrementTrigger: controls.Triggers.Decrement,
			IncrementTrigger: controls.Triggers.Increment,
		}
		md.controls[controlKey] = &control
	}
}

func NewDevice(deviceConfig config.DeviceConfig, signals chan<- core.Signal, logger *zap.Logger) *MidiDevice {
	midiDevice := MidiDevice{}
	midiDevice.applyConfiguration(deviceConfig, signals, logger)

	return &midiDevice
}

func (md *MidiDevice) sendNamespaceChangedSignal(signals chan<- core.Signal, oldNamespace string, newNamespace string) {
	signal := model.NamespaceChanged{
		Device:       md.name,
		OldNamespace: oldNamespace,
		NewNamespace: newNamespace,
	}
	signals <- signal
}
