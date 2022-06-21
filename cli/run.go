package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var input []string
var outputFile string
var pageNumbers bool
var authorName bool
var doubleSpace bool

var cmd = &cobra.Command{
	Use:     "ciere",
	Version: "0.0.1",
	Short:   "Convert markdown into docx for submissions",
	Long:    "Convert markdown into docx for submissions in the format publishers like",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	cmd.PersistentFlags().StringArrayVarP(&input, "stories", "s", []string{}, "Files to submit (required)")
	cmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Supply a name for the output file")
	cmd.PersistentFlags().BoolVarP(&pageNumbers, "page", "p", false, "Put page numbers in header")
	cmd.PersistentFlags().BoolVarP(&authorName, "author", "a", false, "Put author name in header")
	cmd.PersistentFlags().BoolVarP(&doubleSpace, "double", "d", false, "Double-space the document")

	cmd.MarkPersistentFlagRequired("stories")

}

func Run() int {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "CLI error: '%s'\n", err)
		return 1
	}

	println(outputFile)
	return 0
}
