package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/YoshihikoAbe/extdrm"
	"github.com/spf13/cobra"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump SOURCE DESTINATION PRESET",
	Short: "Decrypt the contents of an encrypted extdrm filesystem",
	Args:  cobra.MinimumNArgs(3),
	Run:   runDump,
}

func init() {
	rootCmd.AddCommand(dumpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dumpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dumpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	dumpCmd.Flags().Int("workers", 0, "Number of workers. Specify a value less than one, and the number of logical CPUs available to the process will be used")
}

func runDump(cmd *cobra.Command, args []string) {
	workers, _ := cmd.Flags().GetInt("workers")
	if workers < 1 {
		workers = runtime.NumCPU()
	}

	src := args[0]
	dest := args[1]
	filename := args[2]

	data, err := os.ReadFile(filename)
	if err != nil {
		fatal(err)
	}
	config := &extdrm.DrmConfig{}
	if err := json.Unmarshal(data, config); err != nil {
		fatal(err)
	}

	start := time.Now()
	ch, err := extdrm.ReadFS(*config, src)
	if err != nil {
		fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			for {
				file, ok := <-ch
				if !ok {
					wg.Done()
					return
				}

				func(file extdrm.DrmFile) {
					defer file.Close()

					dir, _ := path.Split(file.Path)
					if err := os.MkdirAll(path.Join(dest, dir), 0777); err != nil {
						fmt.Println(err)
						return
					}

					out, err := os.Create(path.Join(dest, file.Path))
					if err != nil {
						fmt.Println(err)
						return
					}
					defer out.Close()

					if _, err := io.Copy(out, file); err != nil {
						fmt.Println(err)
					}
				}(file)
			}
		}()
	}
	wg.Wait()
	fmt.Println("time elapsed:", time.Since(start))
}
