package backlight

import (
	"log"
	"testing"
)

func BenchmarkTurnLight(b *testing.B) {
	decodedBacklightConfig, err := InitConfig("../../configs/backlight_config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, err = decodedBacklightConfig.TurnLight("MPD226", 2, "red", "on")
		if err != nil {
			log.Fatal(err)
		}
	}
}
