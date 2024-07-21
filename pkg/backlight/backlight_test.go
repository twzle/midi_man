package backlight

import (
	"log"
	"testing"
)

// Function checks the performance of backlight configuration initialization
func BenchmarkTurnLight(b *testing.B) {
	decodedBacklightConfig, err := InitConfig("../../configs/backlight_config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = decodedBacklightConfig.TurnLight("Arduino", 2, "red", "on")
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Function checks the performance of backlight configuration initialization
func BenchmarkReadConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := InitConfig("../../configs/backlight_config.yaml") // O(N)
		if err != nil {
			log.Fatal(err)
		}
	}
}
