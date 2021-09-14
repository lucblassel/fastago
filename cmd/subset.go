/*
Copyright © 2021 LUC BLASSEL

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/lucblassel/fastago/pkg/seqs"
	"github.com/spf13/cobra"
	"os"
)

var namesFile string

// subsetCmd represents the subset command
var subsetCmd = &cobra.Command{
	Use:   "subset",
	Short: "Subset sequences by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		var names map[string]bool
		var err error

		if namesFile != "" {
			names, err = readNames(namesFile)
			if err != nil {
				return err
			}
		} else {
			names = make(map[string]bool)
			for _, name := range args {
				names[name] = true
			}
		}

		records := make(chan seqs.SeqRecord)
		errs := make(chan error)

		go seqs.ReadFastaRecords(inputReader, records, errs)

		for records != nil && errs != nil {
			select {
			case record := <-records:
				if names[record.Name] {
					output, err := record.Seq.FormatSeq(outputLineWidth)
					_, err = fmt.Fprintf(outputWriter, ">%s\n%s\n", record.Name, output)
					if err != nil {
						return err
					}
				}
			case err := <-errs:
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(subsetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subsetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	subsetCmd.Flags().StringVarP(&namesFile, "names", "n", "", "file containing the names of sequences to keep. One name by line")
}

func readNames(filename string) (map[string]bool, error) {
	names := make(map[string]bool)
	input, err := os.Open(filename)
	defer input.Close()

	if err != nil {
		return names, err
	}
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		name := scanner.Text()
		names[name] = true
	}

	if err := scanner.Err(); err != nil {
		return names, err
	}

	return names, nil
}
