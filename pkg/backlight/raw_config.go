package backlight

type RawColor struct {
	ColorName string `json:"color_name" yaml:"color_name"`
	Payload   string `json:"payload" yaml:"payload"`
}

type RawColorSpace struct {
	Id  int        `json:"color_space_id" yaml:"color_space_id"`
	On  []RawColor `json:"on" yaml:"on"`
	Off []RawColor `json:"off" yaml:"off"`
}

type RawStatus struct {
	Type          string `json:"type" yaml:"type"`
	FallbackColor string `json:"fallback_color" yaml:"fallback_color"`
	Bytes         string `json:"bytes" yaml:"bytes"`
}

type RawKeyBacklightStatuses struct {
	On  RawStatus `json:"on" yaml:"on"`
	Off RawStatus `json:"off" yaml:"off"`
}

type RawKeyBacklight struct {
	KeyRange          [2]byte                 `json:"key_range" yaml:"key_range"`
	ColorSpace        int                     `json:"color_space" yaml:"color_space"`
	BacklightStatuses RawKeyBacklightStatuses `json:"statuses" yaml:"statuses"`
	KeyNumberShift    int                     `json:"key_number_shift" yaml:"key_number_shift"`
}

type RawDeviceBacklightConfig struct {
	DeviceName          string            `json:"device_name" yaml:"device_name"`
	BacklightTimeOffset int               `json:"backlight_time_offset" yaml:"backlight_time_offset"`
	ColorSpaces         []RawColorSpace   `json:"color_spaces" yaml:"color_spaces"`
	KeyboardBacklight   []RawKeyBacklight `json:"keyboard_backlight" yaml:"keyboard_backlight"`
}

type RawBacklightConfig struct {
	DeviceBacklightConfigurations []RawDeviceBacklightConfig `json:"device_light_configuration" yaml:"device_light_configuration"`
}
