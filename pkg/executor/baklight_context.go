package executor

type BacklightContext struct {
	backlightBuffer []KeyBacklight
}

type KeyBacklight struct {
	KeyCode  float64
	IsActive bool
}
