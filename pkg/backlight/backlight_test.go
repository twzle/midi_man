package backlight

import (
	"testing"
)

func BenchmarkTurnLight(b *testing.B) {
	_, decodedBacklightConfig, _ := InitConfig("../../configs/backlight_config.yaml")

	for i := 0; i < b.N; i++ {
		_, _ = decodedBacklightConfig.TurnLight("MPD226", 2, "red", "on")
	}
}
