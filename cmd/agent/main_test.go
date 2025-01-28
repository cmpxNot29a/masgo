package main

import (
	"testing"
)

func Test_collectMetrics(t *testing.T) {
	tests := []struct {
		name       string
		wantKeys   []string // Ожидаемые ключи в метриках
		checkValue bool     // Нужно ли проверять конкретные значения (например, > 0)
	}{
		{
			name: "Basic runtime metrics collected",
			wantKeys: []string{
				"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
				"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
				"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
				"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
				"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
				"Sys", "TotalAlloc", "RandomValue",
			},
			checkValue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collectMetrics()

			// Проверка на наличие всех ожидаемых метрик.
			for _, key := range tt.wantKeys {
				if _, exists := got[key]; !exists {
					t.Errorf("Expected metric %s not found in result", key)
				}
			}

			// Если требуется, проверяем значения.
			if tt.checkValue {
				for key, value := range got {
					if value < 0 {
						t.Errorf("Metric %s has invalid value: %f", key, value)
					}
				}
			}
		})
	}
}
