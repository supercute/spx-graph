package spx

import "time"

type Event struct {
	FunctionID int   `json:"function_id"`
	EventType  int   `json:"event_type"` // 1 = enter, 0 = exit
	Time       int64 `json:"time"`       // микросекунды
	Memory     int64 `json:"memory"`     // байты
}

type Profile struct {
	Events    []Event        `json:"events"`
	Functions map[int]string `json:"functions"` // ID -> имя функции
}

type CallInfo struct {
	FunctionID  int           `json:"function_id"`
	Parent      *int          `json:"parent"`       // ID родительской функции
	Duration    time.Duration `json:"duration"`     // длительность выполнения
	MemoryDelta int64         `json:"memory_delta"` // изменение памяти
	StartTime   int64         `json:"start_time"`
	EndTime     int64         `json:"end_time"`
	Children    []int         `json:"children"` // ID дочерних функций
}

type CallGraph struct {
	Nodes map[int]*CallNode `json:"nodes"`
	Edges []*CallEdge       `json:"edges"`
	Root  int               `json:"root"`
}

type CallNode struct {
	FunctionID     int           `json:"function_id"`
	Name           string        `json:"name"`
	TotalDuration  time.Duration `json:"total_duration"`
	SelfDuration   time.Duration `json:"self_duration"`
	CallCount      int           `json:"call_count"`
	TotalMemory    int64         `json:"total_memory"`
	Percentage     float64       `json:"percentage"` // процент от общего времени
	SelfPercentage float64       `json:"self_percentage"`
}

type CallEdge struct {
	From          int           `json:"from"`
	To            int           `json:"to"`
	CallCount     int           `json:"call_count"`
	TotalDuration time.Duration `json:"total_duration"`
	Percentage    float64       `json:"percentage"`
}

type FunctionStats struct {
	TotalDuration time.Duration
	SelfDuration  time.Duration
	CallCount     int
	TotalMemory   int64
	MaxDuration   time.Duration
}
