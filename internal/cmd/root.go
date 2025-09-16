package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/supercute/spx-graph/internal/graph"
	"github.com/supercute/spx-graph/internal/server"
	"github.com/supercute/spx-graph/internal/spx"
	"time"
)

var (
	inputFile  string
	outputFile string
	port       int
)

var rootCmd = &cobra.Command{
	Use:   "spx-graph",
	Short: "SPX profile graph visualizer",
	Long: `Visualize SPX profiling data as interactive call graphs, similar to go tool pprof.

Examples:
  spx-graph --file profile.txt.gz
  spx-graph --file profile.txt.gz -o result.html
  spx-graph --file profile.txt.gz --port 9090`,
	RunE: runGraph,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "file", "f", "", "SPX profile file (.txt or .txt.gz)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output HTML file (default: start server)")
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Server port")
	err := rootCmd.MarkFlagRequired("file")
	if err != nil {
		return
	}
}

func runGraph(cmd *cobra.Command, args []string) error {
	// Parse spx file
	profile, err := spx.ParseProfile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to parse profile: %w", err)
	}

	fmt.Printf("Parsed %d events, %d functions\n",
		len(profile.Events), len(profile.Functions))

	// Analyze call three
	fmt.Printf("Analyze and build graph...\n")
	callGraph := buildGraph(profile)

	fmt.Printf("Build call graph with %d nodes, %d edges\n",
		len(callGraph.Nodes), len(callGraph.Edges))

	// Generate graph
	generator := graph.NewGenerator(callGraph, profile.Functions)

	if outputFile != "" {
		fmt.Printf("Saving HTML to: %s\n", outputFile)
		return generator.SaveHTML(outputFile)
	} else {
		srv := server.New(generator, port)
		fmt.Printf("Starting server at http://localhost:%d\n", port)
		fmt.Println("Press Ctrl+C to stop")
		return srv.Start()
	}
}

func buildGraph(profile *spx.Profile) *spx.CallGraph {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Printf("Build time is %v\n", elapsed)
	}()
	analyzer := spx.NewAnalyzer(profile)
	callGraph := analyzer.BuildCallGraph()
	return callGraph
}
