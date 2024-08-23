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

const (
	deviceDisconnectedCheckLabelFormat = "DEVICE_DISCONNECTED_%s"
)

// Representation of MIDI-device entity
type MidiDevice struct {
	name               string
	active             bool
	ports              MidiPorts
	clickBuffer        ClickBuffer
	holdDelta          time.Duration
	startupDelay       time.Duration
	reconnectInterval  time.Duration
	mutex              sync.Mutex
	stopReconnect      chan struct{}
	stopListen         chan struct{}
	stopBlinking       chan struct{}
	namespace          string
	connected          atomic.Bool
	controls           map[int]*Control
	signals            chan<- core.Signal
	logger             *zap.Logger
	conf               config.DeviceConfig
	reconnectedEvent   chan bool
	checkManager       core.CheckRegistry
	blinkingKeys       map[int]blinkingKey
	blinkingQueueMutex sync.Mutex
}

// Representation of blinking key entity
type blinkingKey struct {
	keyCode         int
	OnColorName     string
	OffColorName    string
	backlightConfig *backlight.DeviceBacklightConfig
}

// Representation of MIDI-ports entity
type MidiPorts struct {
	in  drivers.In
	out drivers.Out
}

// Function returns alias of MIDI-device
func (md *MidiDevice) GetAlias() string {
	return md.name
}

// Function executes command on MIDI-device
func (md *MidiDevice) ExecuteCommand(command model.MidiCommand, backlightConfig *backlight.DeviceBacklightConfig) error {
	md.mutex.Lock()
	defer md.mutex.Unlock()

	switch cmd := command.(type) {
	case model.TurnLightOnCommand:
		md.turnLightOn(cmd, backlightConfig)
	case model.TurnLightOffCommand:
		md.turnLightOff(cmd, backlightConfig)
	case model.SingleBlinkCommand:
		md.singleBlink(cmd, backlightConfig)
	case model.SingleReversedBlinkCommand:
		md.singleReversedBlink(cmd, backlightConfig)
	case model.SetActiveNamespaceCommand:
		md.setActiveNamespace(cmd, backlightConfig)
	case model.StartBlinkingCommand:
		md.blinkingQueueMutex.Lock()
		md.blinkingKeys[cmd.KeyCode] = blinkingKey{
			keyCode:         cmd.KeyCode,
			OnColorName:     cmd.OnColorName,
			OffColorName:    cmd.OffColorName,
			backlightConfig: backlightConfig,
		}
		md.blinkingQueueMutex.Unlock()
	case model.StopBlinkingCommand:
		md.blinkingQueueMutex.Lock()
		delete(md.blinkingKeys, cmd.KeyCode)
		md.blinkingQueueMutex.Unlock()
	default:
		md.logger.Warn("Unknown command", zap.Any("command", cmd))
	}
	return nil
}

// Function frees resources of current MIDI-device
func (md *MidiDevice) Stop() {
	close(md.stopListen)
	close(md.stopReconnect)
	close(md.stopBlinking)
}

// Function initialized working process for MIDI-device
func (md *MidiDevice) RunDevice(backlightConfig *backlight.DeviceBacklightConfig) {
	time.Sleep(md.startupDelay)
	go md.reconnect(backlightConfig)
	go md.listen()
	go md.blinking()
}

// Function initializes connection with MIDI-device connected through system
func (md *MidiDevice) initConnection(backlightConfig *backlight.DeviceBacklightConfig) error {
	md.mutex.Lock()
	defer md.mutex.Unlock()

	if err := md.connectDevice(); err != nil {
		return err
	}
	go md.startupIllumination(backlightConfig)
	md.clickBuffer = make(map[uint8]*KeyContext)
	md.applyControls(md.conf.Controls)
	return nil
}

// Function contains logic of reconnect to MIDI-device through system with specified time interval
func (md *MidiDevice) reconnect(backlightConfig *backlight.DeviceBacklightConfig) {
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

					check := core.NewCheck(
						fmt.Sprintf(deviceDisconnectedCheckLabelFormat, md.name),
						"device was disconnected",
					)
					md.checkManager.RegisterFail(check)
				}
			} else {
				if connected {
					err := md.initConnection(backlightConfig)
					if err != nil {
						md.logger.Warn("Unable to connect device", zap.Error(err))

						connCheck := core.NewCheck(
							fmt.Sprintf(deviceDisconnectedCheckLabelFormat, md.name),
							err.Error(),
						)
						md.checkManager.RegisterFail(connCheck)
					} else {
						md.reconnectedEvent <- connected

						connCheck := core.NewCheck(
							fmt.Sprintf(deviceDisconnectedCheckLabelFormat, md.name),
							"",
						)
						md.checkManager.RegisterSuccess(connCheck)
					}
				} else {
					msg := "No hardware connection to device"
					md.logger.Debug(msg)

					hardwareCheck := core.NewCheck(
						fmt.Sprintf(deviceDisconnectedCheckLabelFormat, md.name),
						msg,
					)
					md.checkManager.RegisterFail(hardwareCheck)
				}
			}
		}
	}
}

// Function contains logic of connection to IN and OUT ports of MIDI-device through system
func (md *MidiDevice) connectDevice() error {
	if inErr := md.connectInPort(); inErr != nil {
		return fmt.Errorf("connection of device failed, inPort:\"{%w}\"", inErr)
	}
	if outErr := md.connectOutPort(); outErr != nil {
		return fmt.Errorf("connection of device failed, outPort:\"{%w}\"", outErr)
	}
	return nil
}

// Function contains logic of connection to OUT port of MIDI-device through system
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

	md.ports.out = port
	return nil
}

// Function contains logic of connection to IN port of MIDI-device through system
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

	md.ports.in = port
	return nil
}


// Function apllying configurtaion to MIDI-device entity
func (md *MidiDevice) applyConfiguration(
	deviceConfig config.DeviceConfig,
	signals chan<- core.Signal,
	logger *zap.Logger,
	checkManager core.CheckRegistry,
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
	md.stopBlinking = make(chan struct{})
	md.reconnectedEvent = make(chan bool)
	md.blinkingKeys = make(map[int]blinkingKey)
	md.namespace = deviceConfig.Namespace
	md.signals = signals
	md.logger = logger.With(zap.String("alias", md.name))
	md.checkManager = checkManager
	md.applyControls(deviceConfig.Controls)
}


// Function applies configuration of controls to MIDI-device entity
func (md *MidiDevice) applyControls(controlsList []config.Controls) {
	md.controls = make(map[int]*Control)
	for _, controls := range controlsList {
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
}

// Function initializes midi device entity with values
func NewDevice(
	deviceConfig config.DeviceConfig,
	signals chan<- core.Signal,
	logger *zap.Logger,
	checkManager core.CheckRegistry,
) *MidiDevice {
	midiDevice := MidiDevice{}
	midiDevice.applyConfiguration(deviceConfig, signals, logger, checkManager)

	return &midiDevice
}

// Function sends callback signals after namespace is changed
func (md *MidiDevice) sendNamespaceChangedSignal(signals chan<- core.Signal, oldNamespace string, newNamespace string) {
	signal := model.NamespaceChanged{
		Device:       md.name,
		OldNamespace: oldNamespace,
		NewNamespace: newNamespace,
	}
	signals <- signal
}

// Function contains logic of continuous blinking for single component of MIDI-device
func (md *MidiDevice) blinking() {
	period := time.Duration(md.conf.BlinkingPeriodMS) * time.Millisecond
	if period == 0 {
		period = 1 * time.Second
	}
	ticker := time.NewTicker(period)
	isOn := true
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			md.blinkingQueueMutex.Lock()
			for _, b := range md.blinkingKeys {
				color := b.OffColorName
				if isOn {
					color = b.OnColorName
				}
				md.turnLightOn(
					model.TurnLightOnCommand{
						KeyCode:     b.keyCode,
						ColorName:   color,
						DeviceAlias: md.name,
					}, b.backlightConfig,
				)
			}
			isOn = !isOn
			md.blinkingQueueMutex.Unlock()
		case <-md.stopBlinking:
			return
		}
	}
}
