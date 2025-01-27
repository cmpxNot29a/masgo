package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/cmpxNot29a/masgo/internal/metrics"
	"github.com/cmpxNot29a/masgo/internal/storage"
)

// Константы для путей и сообщений об ошибках.
const (
	pathSeparator     = "/"      // Разделитель пути в URL.
	pathUpdateSegment = "update" // Сегмент пути "update"

	errorOnlyPostAllowed = "only POST requests are allowed" // Ошибка: разрешены только POST-запросы.
	// errorInvalidRequestFormat = "invalid request format"         // Ошибка: неверный формат запроса.
	errorNotFaundPage        = "page not found"          // Ошибка: страница не найдена.
	errorMetricNameRequired  = "metric name is required" // Ошибка: требуется имя метрики.
	errorInvalidGaugeValue   = "invalid gauge value"     // Ошибка: неверное значение gauge.
	errorInvalidCounterValue = "invalid counter value"   // Ошибка: неверное значение counter.
	errorInvalidMetricType   = "invalid metric type"     // Ошибка: неверный тип метрики.
)

// validateURL проверяет корректность URL запроса.
func ValidateURL(r *http.Request) error {
	// Проверяем, что используется метод POST.
	if r.Method != http.MethodPost {
		return errors.New(errorOnlyPostAllowed)
	}

	// Удаляем начальные и конечные слеши и разбиваем URL на части.
	path := strings.Trim(r.URL.Path, pathSeparator)
	parts := strings.Split(path, pathSeparator)

	// Проверяем, что URL имеет формат /update/{type}/{name}/{value},
	// то есть 4 части после удаления начальных и конечных слешей.
	if len(parts) != 4 {
		return errors.New(errorNotFaundPage)
	}

	// Проверяем, что первая часть URL - "update".
	if parts[0] != pathUpdateSegment {
		return errors.New(errorNotFaundPage)
	}

	// Проверяем наличие имени метрики (третья часть URL).
	if parts[2] == "" {
		return errors.New(errorMetricNameRequired)
	}

	return nil
}

// validateMetric проверяет корректность типа и значения метрики.
func ValidateMetric(metricTypeStr, metricName, metricValueStr string) (metrics.MetricType, interface{}, error) {
	// Проверяем, что имя метрики не пустое.
	if metricName == "" {
		return "", nil, errors.New(errorMetricNameRequired)
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
			return "", nil, errors.New(errorInvalidGaugeValue)
		}
	case metrics.TypeCounter:
		// Парсим значение как int64.
		metricValue, err = strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			return "", nil, errors.New(errorInvalidCounterValue)
		}
	default:
		return "", nil, errors.New(errorInvalidMetricType)
	}

	return metricType, metricValue, nil
}

// validateAndUpdateRequest валидирует запрос и, в случае успеха, обновляет хранилище метрик.
func ValidateAndUpdateRequest(r *http.Request, storage storage.Repository) error {
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

func UpdateHandler(storage storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := ValidateAndUpdateRequest(r, storage)
		if err != nil {
			// Обрабатываем ошибки валидации.
			switch err.Error() {
			case errorOnlyPostAllowed:
				// Метод отличается от POST.
				http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			case errorNotFaundPage:
				// Неправильный формат запроса.
				http.Error(w, err.Error(), http.StatusNotFound)
			case errorMetricNameRequired:
				// Не указано имя метрики.
				http.Error(w, err.Error(), http.StatusNotFound)
			case errorInvalidGaugeValue, errorInvalidCounterValue, errorInvalidMetricType:
				// Ошибка в типе метрики или значении.
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				// Непредвиденная ошибка.
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		// Отправляем ответ с кодом 200 OK в случае успеха.
		w.WriteHeader(http.StatusOK)
	}
}
