package spx

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Analyzer struct {
	profile *Profile
}

func NewAnalyzer(profile *Profile) *Analyzer {
	return &Analyzer{profile: profile}
}

// BuildCallGraph build call graph from profile events
func (a *Analyzer) BuildCallGraph() *CallGraph {
	callInfos := a.buildCallTree()
	stats := a.calculateStats(callInfos)

	return a.createCallGraph(callInfos, stats)
}

// buildCallTree builds a call tree from events
func (a *Analyzer) buildCallTree() []CallInfo {
	var stack []CallInfo
	var callInfos []CallInfo

	for _, event := range a.profile.Events {
		funcID := event.FunctionID

		if event.EventType == 1 { // start
			var parent *int
			if len(stack) > 0 {
				parentID := stack[len(stack)-1].FunctionID
				parent = &parentID
			}

			callInfo := CallInfo{
				FunctionID: funcID,
				Parent:     parent,
				StartTime:  event.Time,
				Children:   make([]int, 0),
			}

			stack = append(stack, callInfo)

		} else if event.EventType == 0 { // end
			for i := len(stack) - 1; i >= 0; i-- {
				if stack[i].FunctionID == funcID {
					callInfo := stack[i]
					callInfo.EndTime = event.Time
					callInfo.Duration = time.Duration(event.Time-callInfo.StartTime) * time.Microsecond
					callInfo.MemoryDelta = event.Memory

					// Add to parent as child
					if callInfo.Parent != nil {
						for j := i - 1; j >= 0; j-- {
							if stack[j].FunctionID == *callInfo.Parent {
								stack[j].Children = append(stack[j].Children, funcID)
								break
							}
						}
					}

					callInfos = append(callInfos, callInfo)

					// Remove from stack
					stack = append(stack[:i], stack[i+1:]...)
					break
				}
			}
		}
	}

	return callInfos
}

// calculateStats calculate statistics for each function
func (a *Analyzer) calculateStats(callInfos []CallInfo) map[int]*FunctionStats {
	stats := make(map[int]*FunctionStats)

	for _, call := range callInfos {
		funcID := call.FunctionID

		if stats[funcID] == nil {
			stats[funcID] = &FunctionStats{}
		}

		stats[funcID].TotalDuration += call.Duration
		stats[funcID].CallCount++
		stats[funcID].TotalMemory += call.MemoryDelta

		if call.Duration > stats[funcID].MaxDuration {
			stats[funcID].MaxDuration = call.Duration
		}

		// Calculate self-time (time excluding child calls)
		selfTime := call.Duration
		for _, childID := range call.Children {
			for _, childCall := range callInfos {
				if childCall.FunctionID == childID &&
					childCall.Parent != nil &&
					*childCall.Parent == funcID &&
					childCall.StartTime >= call.StartTime &&
					childCall.EndTime <= call.EndTime {
					selfTime -= childCall.Duration
				}
			}
		}
		if selfTime > 0 {
			stats[funcID].SelfDuration += selfTime
		}
	}

	return stats
}

// createCallGraph creates the call graph structure
func (a *Analyzer) createCallGraph(callInfos []CallInfo, stats map[int]*FunctionStats) *CallGraph {
	nodes := make(map[int]*CallNode)

	// Calculate total execution time
	var totalTime time.Duration
	for _, stat := range stats {
		if stat.TotalDuration > totalTime {
			totalTime = stat.TotalDuration
		}
	}

	// Create nodes
	for funcID, stat := range stats {
		funcName := a.profile.Functions[funcID]
		if funcName == "" {
			funcName = fmt.Sprintf("func_%d", funcID)
		}

		// Split long names
		if len(funcName) > 50 {
			parts := strings.Split(funcName, "/")
			if len(parts) > 1 {
				funcName = ".../" + parts[len(parts)-1]
			}
		}

		percentage := 0.0
		selfPercentage := 0.0
		if totalTime > 0 {
			percentage = float64(stat.TotalDuration) / float64(totalTime) * 100
			selfPercentage = float64(stat.SelfDuration) / float64(totalTime) * 100
		}

		nodes[funcID] = &CallNode{
			FunctionID:     funcID,
			Name:           funcName,
			TotalDuration:  stat.TotalDuration,
			SelfDuration:   stat.SelfDuration,
			CallCount:      stat.CallCount,
			TotalMemory:    stat.TotalMemory,
			Percentage:     percentage,
			SelfPercentage: selfPercentage,
		}
	}

	// Create edges
	edgeStats := make(map[string]struct {
		count    int
		duration time.Duration
	})

	for _, call := range callInfos {
		if call.Parent != nil {
			key := fmt.Sprintf("%d->%d", *call.Parent, call.FunctionID)
			edge := edgeStats[key]
			edge.count++
			edge.duration += call.Duration
			edgeStats[key] = edge
		}
	}

	var edges []*CallEdge
	for key, edge := range edgeStats {
		parts := strings.Split(key, "->")
		if len(parts) != 2 {
			continue
		}

		from, _ := strconv.Atoi(parts[0])
		to, _ := strconv.Atoi(parts[1])

		percentage := 0.0
		if totalTime > 0 {
			percentage = float64(edge.duration) / float64(totalTime) * 100
		}

		edges = append(edges, &CallEdge{
			From:          from,
			To:            to,
			CallCount:     edge.count,
			TotalDuration: edge.duration,
			Percentage:    percentage,
		})
	}

	// Find root node
	root := 0
	for _, call := range callInfos {
		if call.Parent == nil {
			root = call.FunctionID
			break
		}
	}

	return &CallGraph{
		Nodes: nodes,
		Edges: edges,
		Root:  root,
	}
}
