package backlight

// Representation of color deserealized from backlight configuration before optimization
type RawColor struct {
	ColorName string `json:"color_name" yaml:"color_name"`
	Payload   string `json:"payload" yaml:"payload"`
}

// Representation of color space deserealized from backlight configuration before optimization
type RawColorSpace struct {
	Id  int        `json:"color_space_id" yaml:"color_space_id"`
	On  []RawColor `json:"on" yaml:"on"`
	Off []RawColor `json:"off" yaml:"off"`
}

// Representation of status deserealized from backlight configuration before optimization
type RawStatus struct {
	Type          string `json:"type" yaml:"type"`
	FallbackColor string `json:"fallback_color" yaml:"fallback_color"`
	Bytes         string `json:"bytes" yaml:"bytes"`
}

// Representation of key backlight statuses deserealized from backlight configuration before optimization
type RawKeyBacklightStatuses struct {
	On  RawStatus `json:"on" yaml:"on"`
	Off RawStatus `json:"off" yaml:"off"`
}

// Representation of key backlight deserealized from backlight configuration before optimization
type RawKeyBacklight struct {
	KeyRange          [2]byte                 `json:"key_range" yaml:"key_range"`
	ColorSpace        int                     `json:"color_space" yaml:"color_space"`
	BacklightStatuses RawKeyBacklightStatuses `json:"statuses" yaml:"statuses"`
	KeyNumberShift    int                     `json:"key_number_shift" yaml:"key_number_shift"`
}

// Representation of device backlight configuration deserealized from backlight configuration before optimization
type RawDeviceBacklightConfig struct {
	DeviceName          string            `json:"device_name" yaml:"device_name"`
	BacklightTimeOffset int               `json:"backlight_time_offset" yaml:"backlight_time_offset"`
	ColorSpaces         []RawColorSpace   `json:"color_spaces" yaml:"color_spaces"`
	KeyboardBacklight   []RawKeyBacklight `json:"keyboard_backlight" yaml:"keyboard_backlight"`
}

// Representation of backlight configuration deserealized from backlight configuration before optimization
type RawBacklightConfig struct {
	DeviceBacklightConfigurations []RawDeviceBacklightConfig `json:"device_light_configuration" yaml:"device_light_configuration"`
}
