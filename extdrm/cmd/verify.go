package cmd

import (
	"encoding/json"
	"os"

	"github.com/YoshihikoAbe/extdrm"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify ROOT METADATA",
	Short: "Verify the integrity of a filesystem dump",
	Args:  cobra.MinimumNArgs(2),
	Run:   runVerify,
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runVerify(cmd *cobra.Command, args []string) {
	data, err := os.ReadFile(args[1])
	if err != nil {
		fatal(err)
	}
	meta := &extdrm.Metadata{}
	if err := json.Unmarshal(data, meta); err != nil {
		fatal(err)
	}

	result, err := extdrm.VerifyFS(args[0], meta)
	if err != nil {
		fatal(err)
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fatal(err)
	}
	os.Stdout.Write(out)
}
