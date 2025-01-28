package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
	"net/http/httputil"
)

// Константы для настройки агента.
const (
	serverAddress         = "http://localhost:8080" // Адрес сервера
	pollInterval          = 2 * time.Second         // Интервал сбора метрик
	reportInterval        = 10 * time.Second        // Интервал отправки метрик
	contentTypeTextPlain  = "text/plain"
	endpointUpdateGauge   = "/update/gauge/"
	endpointUpdateCounter = "/update/counter/"
	unsupportedFieldMsg   = "Unsupported field type for metric %s (type: %s)\n"
	errorSendingMetric    = "Error sending metric %s: %v\n"
	errorSendingCounter   = "Error sending %s: %v\n"
)

// Константы для имен метрик.
const (
	metricAlloc         = "Alloc"
	metricBuckHashSys   = "BuckHashSys"
	metricFrees         = "Frees"
	metricGCCPUFraction = "GCCPUFraction"
	metricGCSys         = "GCSys"
	metricHeapAlloc     = "HeapAlloc"
	metricHeapIdle      = "HeapIdle"
	metricHeapInuse     = "HeapInuse"
	metricHeapObjects   = "HeapObjects"
	metricHeapReleased  = "HeapReleased"
	metricHeapSys       = "HeapSys"
	metricLastGC        = "LastGC"
	metricLookups       = "Lookups"
	metricMCacheInuse   = "MCacheInuse"
	metricMCacheSys     = "MCacheSys"
	metricMSpanInuse    = "MSpanInuse"
	metricMSpanSys      = "MSpanSys"
	metricMallocs       = "Mallocs"
	metricNextGC        = "NextGC"
	metricNumForcedGC   = "NumForcedGC"
	metricNumGC         = "NumGC"
	metricOtherSys      = "OtherSys"
	metricPauseTotalNs  = "PauseTotalNs"
	metricStackInuse    = "StackInuse"
	metricStackSys      = "StackSys"
	metricSys           = "Sys"
	metricTotalAlloc    = "TotalAlloc"
	metricRandomValue   = "RandomValue"
	metricPollCount     = "PollCount"
)

// Метрики, которые будем собирать из пакета runtime.
var runtimeMetrics = []string{
	metricAlloc, metricBuckHashSys, metricFrees, metricGCCPUFraction, metricGCSys,
	metricHeapAlloc, metricHeapIdle, metricHeapInuse, metricHeapObjects, metricHeapReleased,
	metricHeapSys, metricLastGC, metricLookups, metricMCacheInuse, metricMCacheSys,
	metricMSpanInuse, metricMSpanSys, metricMallocs, metricNextGC, metricNumForcedGC,
	metricNumGC, metricOtherSys, metricPauseTotalNs, metricStackInuse, metricStackSys,
	metricSys, metricTotalAlloc,
}

// pollCount - счетчик обновлений метрик.
var pollCount int64

func collectMetrics() map[string]float64 {
	metrics := make(map[string]float64)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Используем рефлексию для динамического доступа к полям MemStats.
	memStatsValue := reflect.ValueOf(memStats)
	for _, metricName := range runtimeMetrics {
		field := memStatsValue.FieldByName(metricName)
		if field.IsValid() {
			switch field.Kind() {
			case reflect.Uint64:
				metrics[metricName] = float64(field.Uint()) // Преобразуем uint64 в float64
			case reflect.Uint32:
				metrics[metricName] = float64(field.Uint()) // Преобразуем uint32 в float64
			case reflect.Float64:
				metrics[metricName] = field.Float() // Преобразуем float64
			default:
				fmt.Printf(unsupportedFieldMsg, metricName, field.Kind())
			}
		}
	}

	// Добавляем дополнительные метрики.
	metrics[metricRandomValue] = rand.Float64() // Произвольное значение
	pollCount++                                 // Увеличиваем счетчик обновлений

	return metrics
}

// sendMetrics отправляет метрики на сервер.
func sendMetrics(metrics map[string]float64) {
	for metricName, metricValue := range metrics {
		url := fmt.Sprintf("%s%s%s/%f", serverAddress, endpointUpdateGauge, metricName, metricValue)
		resp, err := http.Post(url, contentTypeTextPlain, bytes.NewBuffer([]byte{}))
		if err != nil {
			fmt.Printf(errorSendingMetric, metricName, err)
			continue
		}
		b, _ := httputil.DumpResponse(resp, true)
		println(url)
		println(string(b))

	}

	// Отправляем счетчик PollCount.
	url := fmt.Sprintf("%s%s%s/%d", serverAddress, endpointUpdateCounter, metricPollCount, pollCount)
	resp, err := http.Post(url, contentTypeTextPlain, bytes.NewBuffer([]byte{}))
	if err != nil {
		fmt.Printf(errorSendingCounter, metricPollCount, err)
		return
	}
	resp.Body.Close()
}

func main() {
	// Запускаем сбор метрик в отдельной горутине.
	go func() {
		for {
			metrics := collectMetrics()
			sendMetrics(metrics)
			time.Sleep(reportInterval)
		}
	}()

	for {
		time.Sleep(pollInterval)
	}
}

