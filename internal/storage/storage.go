package storage

import (
	"github.com/cmpxNot29a/masgo/internal/metrics"
)

// Repository - интерфейс для работы с хранилищем метрик.
type Repository interface {
	Update(metric metrics.Metric) error
}
