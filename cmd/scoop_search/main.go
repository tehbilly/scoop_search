package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tehbilly/scoop_search"
	"os"
	"strings"
)

var (
	inMem      bool
	numResults int
	outputJson bool
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "scoop-search",
	Short: "Search scoop manifests, but more speedier!",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Don't want to be chatty _and_ output JSON
		if outputJson && verbose {
			verbose = false
		}

		var searcher *scoop_search.ScoopIndex

		if inMem {
			ims, err := scoop_search.NewSearcher(scoop_search.SearcherOpts{Verbose: verbose})
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Unable to open searcher: %s\n", err)
				return err
			}
			searcher = ims
		} else {
			ps, err := scoop_search.NewPersistentSearcher(scoop_search.SearcherOpts{Verbose: verbose})
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Unable to open searcher: %s\n", err)
				return err
			}
			searcher = ps
		}

		results, err := searcher.Search(scoop_search.SearchOptions{
			Verbose: verbose,
			Query:   strings.Join(args, " "),
			Num:     numResults,
		})
		if err != nil {
			panic(err)
		}

		var r renderer

		switch outputJson {
		case true:
			r = &jsonRenderer{}
		case false:
			r = newTableRenderer()
		}

		return r.Render(results)
	},
	Args: cobra.MinimumNArgs(1),
}

func main() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Output operational (debug) messages")
	rootCmd.PersistentFlags().BoolVarP(&inMem, "in-mem", "m", true, "Use in-memory search index")
	rootCmd.PersistentFlags().BoolVarP(&outputJson, "json", "j", false, "Output JSON instead of table")
	rootCmd.PersistentFlags().IntVarP(&numResults, "num-results", "n", 10, "Number of results to show")

	_ = rootCmd.Execute()
}
