package metrics

import (
	"errors"
	"fmt"
)

// Константы для сообщений об ошибках.
const (
	errorUnknownMetricType       = "unknown metric type: %s"    // Ошибка: неизвестный тип метрики.
	errorInvalidCounterValueType = "invalid counter value type" // Ошибка: неверный тип значения счетчика.
)

// MemStorageInterface - интерфейс хранилища метрик.
type MemStorageInterface interface {
	Update(metric Metric) error // Обновляет значение метрики.
}

// MemStorage - реализация хранилища метрик в памяти.
type MemStorage struct {
	metrics map[MetricType]map[string]interface{} // metrics - отображение (map), где ключом является тип метрики, а значением - другое отображение, где ключом является имя метрики, а значением - значение метрики.
}

// NewMemStorage создает новый экземпляр MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: map[MetricType]map[string]interface{}{
			TypeGauge:   make(map[string]interface{}), // Инициализируем отображение для метрик типа Gauge.
			TypeCounter: make(map[string]interface{}), // Инициализируем отображение для метрик типа Counter.
		},
	}
}

// Update обновляет значение метрики в хранилище.
func (s *MemStorage) Update(metric Metric) error {

	switch metric.Type {
	case TypeGauge:
		// Для метрики типа Gauge просто присваиваем новое значение.
		s.metrics[TypeGauge][metric.Name] = metric.Value
	case TypeCounter:
		// Для метрики типа Counter нужно сначала получить текущее значение, а затем добавить к нему новое.
		currentValue, ok := s.metrics[TypeCounter][metric.Name]
		if !ok {
			// Если метрики с таким именем еще нет, считаем, что текущее значение равно 0.
			currentValue = int64(0)
		}
		currentValueInt, ok := currentValue.(int64)
		if !ok {
			// fmt.Printf("metrics[%s][%s] = %v + %v", TypeCounter, metric.Name, currentValueInt, newValueInt)
			return errors.New(errorInvalidCounterValueType)
		}
		newValueInt, ok := metric.Value.(int64)
		if !ok {
			fmt.Printf("metrics[%s][%s] = %v + %v", TypeCounter, metric.Name, currentValueInt, newValueInt)
			return errors.New(errorInvalidCounterValueType)
		}
		fmt.Printf("metrics[%s][%s] = %v + %v", TypeCounter, metric.Name, currentValueInt, newValueInt)
		s.metrics[TypeCounter][metric.Name] = currentValueInt + newValueInt // Добавляем новое значение к текущему.
	default:
		return fmt.Errorf(errorUnknownMetricType, metric.Type)
	}

	return nil
}
