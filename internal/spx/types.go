package spx

import "time"

// Event представляет одно событие профилирования SPX
type Event struct {
	FunctionID int       `json:"function_id"`
	EventType  int       `json:"event_type"` // 1 = enter, 0 = exit
	Time       int64     `json:"time"`       // микросекунды
	Memory     int64     `json:"memory"`     // байты
}

// Profile содержит полный профиль SPX
type Profile struct {
	Events    []Event           `json:"events"`
	Functions map[int]string    `json:"functions"` // ID -> имя функции
}

// CallInfo информация о вызове функции
type CallInfo struct {
	FunctionID   int           `json:"function_id"`
	Parent       *int          `json:"parent"`        // ID родительской функции
	Duration     time.Duration `json:"duration"`      // длительность выполнения
	MemoryDelta  int64         `json:"memory_delta"`  // изменение памяти
	StartTime    int64         `json:"start_time"`
	EndTime      int64         `json:"end_time"`
	Children     []int         `json:"children"`      // ID дочерних функций
}

// CallGraph представляет граф вызовов
type CallGraph struct {
	Nodes map[int]*CallNode `json:"nodes"`
	Edges []*CallEdge       `json:"edges"`
	Root  int               `json:"root"`
}

// CallNode узел в графе вызовов
type CallNode struct {
	FunctionID      int           `json:"function_id"`
	Name            string        `json:"name"`
	TotalDuration   time.Duration `json:"total_duration"`
	SelfDuration    time.Duration `json:"self_duration"`
	CallCount       int           `json:"call_count"`
	TotalMemory     int64         `json:"total_memory"`
	Percentage      float64       `json:"percentage"`     // процент от общего времени
	SelfPercentage  float64       `json:"self_percentage"`
}

// CallEdge ребро в графе вызовов
type CallEdge struct {
	From           int           `json:"from"`
	To             int           `json:"to"`
	CallCount      int           `json:"call_count"`
	TotalDuration  time.Duration `json:"total_duration"`
	Percentage     float64       `json:"percentage"`
}

// FunctionStats статистика функции
type FunctionStats struct {
	TotalDuration time.Duration
	SelfDuration  time.Duration
	CallCount     int
	TotalMemory   int64
	MaxDuration   time.Duration
}
