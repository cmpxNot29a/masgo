package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cmpxNot29a/masgo/internal/metrics"
)

// Константы для путей и сообщений об ошибках.
const (
	pathSeparator = "/" // Разделитель пути в URL.

	errorOnlyPostAllowed      = "only POST requests are allowed" // Ошибка: разрешены только POST-запросы.
	errorInvalidRequestFormat = "invalid request format"         // Ошибка: неверный формат запроса.
	errorMetricNameRequired   = "metric name is required"        // Ошибка: требуется имя метрики.
	errorInvalidGaugeValue    = "invalid gauge value"            // Ошибка: неверное значение gauge.
	errorInvalidCounterValue  = "invalid counter value"          // Ошибка: неверное значение counter.
	errorInvalidMetricType    = "invalid metric type"            // Ошибка: неверный тип метрики.
)

// validateURL проверяет корректность URL запроса.
func ValidateURL(r *http.Request) error {
	// Проверяем, что используется метод POST.
	if r.Method != http.MethodPost {
		return fmt.Errorf(errorOnlyPostAllowed)
	}

	// Разбиваем URL на части по разделителю "/".
	parts := strings.Split(r.URL.Path, pathSeparator)
	// Ожидаем, что URL будет иметь формат /update/{type}/{name}/{value}, то есть 5 частей, включая пустую в начале из-за /
	if len(parts) != 5 {
		return fmt.Errorf(errorInvalidRequestFormat)
	}

	return nil
}

// validateMetric проверяет корректность типа и значения метрики.
func ValidateMetric(metricTypeStr, metricName, metricValueStr string) (metrics.MetricType, interface{}, error) {
	// Проверяем, что имя метрики не пустое.
	if metricName == "" {
		return "", nil, fmt.Errorf(errorMetricNameRequired)
	}

	// Преобразуем строку в тип метрики.
	metricType := metrics.MetricType(metricTypeStr)
	// Объявляем переменную для хранения значения метрики.
	var metricValue interface{}
	var err error

	// В зависимости от типа метрики парсим значение.
	switch metricType {
	case metrics.TypeGauge:
		// Парсим значение как float64.
		metricValue, err = strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			return "", nil, fmt.Errorf(errorInvalidGaugeValue)
		}
	case metrics.TypeCounter:
		// Парсим значение как int64.
		metricValue, err = strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			return "", nil, fmt.Errorf(errorInvalidCounterValue)
		}
	default:
		return "", nil, fmt.Errorf(errorInvalidMetricType)
	}

	return metricType, metricValue, nil
}

// validateAndUpdateRequest валидирует запрос и, в случае успеха, обновляет хранилище метрик.
func ValidateAndUpdateRequest(r *http.Request, storage metrics.MemStorageInterface) error {
	// Проверяем корректность URL.
	err := ValidateURL(r)
	if err != nil {
		return err
	}

	// Разбиваем URL на части.
	parts := strings.Split(r.URL.Path, pathSeparator)
	// Проверяем корректность метрики.
	metricType, metricValue, err := ValidateMetric(parts[2], parts[3], parts[4])
	if err != nil {
		return err
	}

	// Создаем объект метрики.
	metric := metrics.Metric{
		Type:  metricType,
		Name:  parts[3],
		Value: metricValue,
	}

	// Обновляем метрику в хранилище.
	err = storage.Update(metric)
	if err != nil {
		return err
	}

	return nil
}

// UpdateHandler обрабатывает HTTP-запросы на обновление метрик.
func UpdateHandler(storage metrics.MemStorageInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Валидируем запрос и обновляем хранилище.
		err := ValidateAndUpdateRequest(r, storage)
		if err != nil {
			// Обрабатываем ошибки валидации.
			switch err.Error() {
			case errorOnlyPostAllowed:
				http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			case errorInvalidRequestFormat:
				http.Error(w, err.Error(), http.StatusBadRequest)
			case errorMetricNameRequired:
				http.Error(w, err.Error(), http.StatusNotFound)
			case errorInvalidGaugeValue, errorInvalidCounterValue, errorInvalidMetricType:
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		// Отправляем ответ с кодом 200 OK.
		w.WriteHeader(http.StatusOK)
	}
}
