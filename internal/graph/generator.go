package graph

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"math"
	"os"
	"strings"
	"time"

	"github.com/goccy/go-graphviz"
	"github.com/supercute/spx-graph/internal/spx"
)

type Generator struct {
	callGraph *spx.CallGraph
	functions map[int]string
}

func NewGenerator(callGraph *spx.CallGraph, functions map[int]string) *Generator {
	return &Generator{
		callGraph: callGraph,
		functions: functions,
	}
}

func (g *Generator) GenerateSVG() (string, error) {
	ctx := context.Background()

	gv, err := graphviz.New(ctx)
	if err != nil {
		return "", err
	}
	defer gv.Close()

	graph, err := gv.Graph()
	if err != nil {
		return "", err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			// ignore error on close
		}
	}()

	// Graph settings
	graph.Set("rankdir", "TB")
	graph.Set("nodesep", "0.5")
	graph.Set("ranksep", "0.75")
	graph.Set("bgcolor", "white")
	graph.Set("fontname", "Arial")

	// Find max percentage for normalization
	maxPercentage := 0.0
	maxSelfPercentage := 0.0
	for _, node := range g.callGraph.Nodes {
		if node.Percentage > maxPercentage {
			maxPercentage = node.Percentage
		}
		if node.SelfPercentage > maxSelfPercentage {
			maxSelfPercentage = node.SelfPercentage
		}
	}

	// Create nodes
	nodeMap := make(map[int]*graphviz.Node)
	for id, node := range g.callGraph.Nodes {
		// Skipping nodes with very low time
		if node.Percentage < 0.1 {
			continue
		}

		nodeName := fmt.Sprintf("n%d", id)
		n, err := graph.CreateNodeByName(nodeName)
		if err != nil {
			continue
		}

		label := g.createNodeLabel(node)

		// Calculate size and color for node
		nodeSize := g.calculateNodeSize(node, maxSelfPercentage)
		nodeColor := g.calculateNodeColor(node)

		// Set attributes
		n.SetLabel(label)
		n.SetShape("box")
		n.SetStyle("filled")
		n.SetFillColor(nodeColor)
		n.SetFontSize(10.0)
		n.SetWidth(nodeSize)
		n.SetHeight(nodeSize * 0.6)

		nodeMap[id] = n
	}

	// Create adges
	maxEdgePercentage := 0.0
	for _, edge := range g.callGraph.Edges {
		if edge.Percentage > maxEdgePercentage {
			maxEdgePercentage = edge.Percentage
		}
	}

	for _, edge := range g.callGraph.Edges {
		fromNode := nodeMap[edge.From]
		toNode := nodeMap[edge.To]

		if fromNode == nil || toNode == nil {
			continue
		}

		// Skipping low edges
		if edge.Percentage < 0.1 {
			continue
		}

		edgeName := fmt.Sprintf("e_%d_%d", edge.From, edge.To)
		e, err := graph.CreateEdgeByName(edgeName, fromNode, toNode)
		if err != nil {
			continue
		}

		// Calculate size and color for edge
		edgeWidth := g.calculateEdgeWidth(edge, maxEdgePercentage)
		edgeColor := g.calculateEdgeColor(edge)

		e.SetPenWidth(edgeWidth)
		e.SetColor(edgeColor)
		e.SetLabel(fmt.Sprintf("%.1f%%\\n%d calls", edge.Percentage, edge.CallCount))
		e.SetFontSize(8.0)
	}

	var buf bytes.Buffer
	if err := gv.Render(ctx, graph, graphviz.SVG, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *Generator) formatFunctionName(name string) string {
	// Сокращаем длинные пути
	if len(name) > 40 {
		if strings.Contains(name, "/") {
			parts := strings.Split(name, "/")
			if len(parts) > 2 {
				return ".../" + strings.Join(parts[len(parts)-2:], "/")
			}
		}
		// Если это длинное имя класса/метода
		if len(name) > 40 {
			return name[:37] + "..."
		}
	}
	return name
}

func (g *Generator) createNodeLabel(node *spx.CallNode) string {
	name := g.formatFunctionName(node.Name)

	return fmt.Sprintf("%s\n%.1f%% (%.1f%%)\n%s",
		name,
		node.SelfPercentage,
		node.Percentage,
		g.formatDuration(node.SelfDuration))
}

func (g *Generator) formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.0fμs", float64(d.Nanoseconds())/1000)
	} else if d < time.Second {
		return fmt.Sprintf("%.1fms", float64(d.Nanoseconds())/1000000)
	} else {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

func (g *Generator) calculateNodeSize(node *spx.CallNode, maxPercentage float64) float64 {
	minSize := 0.5
	maxSize := 3.0

	if maxPercentage == 0 {
		return minSize
	}

	ratio := node.SelfPercentage / maxPercentage
	return minSize + (maxSize-minSize)*ratio
}

func (g *Generator) calculateNodeColor(node *spx.CallNode) string {
	return g.getHeatColor(node.SelfPercentage)
}

func (g *Generator) calculateEdgeWidth(edge *spx.CallEdge, maxPercentage float64) float64 {
	minWidth := 1.0
	maxWidth := 8.0

	if maxPercentage == 0 {
		return minWidth
	}

	ratio := edge.Percentage / maxPercentage
	return minWidth + (maxWidth-minWidth)*ratio
}

func (g *Generator) calculateEdgeColor(edge *spx.CallEdge) string {
	return g.getEdgeColor(edge.Percentage)
}

// dotColor returns a color for the given score (between -1.0 and 1.0)
// logic from go pprof
// see: https://github.com/google/pprof/blob/main/internal/graph/dotgraph.go
func (g *Generator) dotColor(score float64, isBackground bool) string {
	const shift = 0.7
	const bgSaturation = 0.1
	const bgValue = 0.93
	const fgSaturation = 1.0
	const fgValue = 0.7

	var saturation, value float64
	if isBackground {
		saturation = bgSaturation
		value = bgValue
	} else {
		saturation = fgSaturation
		value = fgValue
	}

	score = math.Max(-1.0, math.Min(1.0, score))

	if math.Abs(score) < 0.2 {
		saturation *= math.Abs(score) / 0.2
	}

	if score > 0.0 {
		score = math.Pow(score, (1.0 - shift))
	}
	if score < 0.0 {
		score = -math.Pow(-score, (1.0 - shift))
	}

	var r, gg, b float64
	if score < 0.0 {
		gg = value
		r = value * (1 + saturation*score)
	} else {
		r = value
		gg = value * (1 - saturation*score)
	}
	b = value * (1 - saturation)

	return fmt.Sprintf("#%02x%02x%02x", uint8(r*255.0), uint8(gg*255.0), uint8(b*255.0))
}

func (g *Generator) getHeatColor(percentage float64) string {
	score := percentage / 100.0

	return g.dotColor(score, true) // true for background color
}

func (g *Generator) getEdgeColor(percentage float64) string {
	score := percentage / 100.0

	return g.dotColor(score, false) // false для for foreground color
}

func (g *Generator) SaveHTML(filename string) error {
	svg, err := g.GenerateSVG()
	if err != nil {
		return fmt.Errorf("failed to generate SVG: %w", err)
	}

	html := g.generateHTML(svg)

	return os.WriteFile(filename, []byte(html), 0644)
}

func (g *Generator) GenerateHTML(svg string) string {
	return g.generateHTML(svg)
}

func (g *Generator) generateHTML(svg string) string {
	tmpl := g.getHTMLTemplate()

	t := template.Must(template.New("html").Parse(tmpl))

	nodeCount := len(g.callGraph.Nodes)
	edgeCount := len(g.callGraph.Edges)
	totalCalls := 0
	for _, edge := range g.callGraph.Edges {
		totalCalls += edge.CallCount
	}

	data := struct {
		SVG        template.HTML
		NodeCount  int
		EdgeCount  int
		TotalCalls int
	}{
		SVG:        template.HTML(svg),
		NodeCount:  nodeCount,
		EdgeCount:  edgeCount,
		TotalCalls: totalCalls,
	}

	var buf bytes.Buffer
	t.Execute(&buf, data)
	return buf.String()
}
