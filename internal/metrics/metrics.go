package metrics

// MetricType - тип метрики (gauge или counter).
type MetricType string

// Константы для типов метрик.
const (
	TypeGauge   MetricType = "gauge"   // Тип метрики gauge (число с плавающей точкой).
	TypeCounter MetricType = "counter" // Тип метрики counter (целое число).
)

// Metric - структура, представляющая метрику.
type Metric struct {
	Type  MetricType  // Тип метрики.
	Name  string      // Имя метрики.
	Value interface{} // Значение метрики (может быть float64 или int64).
}
