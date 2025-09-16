package spx

import "time"

type Event struct {
	FunctionID int   `json:"function_id"`
	EventType  int   `json:"event_type"` // 1 = enter, 0 = exit
	Time       int64 `json:"time"`       // microseconds
	Memory     int64 `json:"memory"`     // bytes
}

type Profile struct {
	Events    []Event        `json:"events"`
	Functions map[int]string `json:"functions"` // ID -> function name
}

type CallInfo struct {
	FunctionID  int           `json:"function_id"`
	Parent      *int          `json:"parent"`       // parent function ID
	Duration    time.Duration `json:"duration"`     // execution duration
	MemoryDelta int64         `json:"memory_delta"` // memory change
	StartTime   int64         `json:"start_time"`
	EndTime     int64         `json:"end_time"`
	Children    []int         `json:"children"` // child function IDs
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
	Percentage     float64       `json:"percentage"` // percentage of total time
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
