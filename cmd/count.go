/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"io"
	"log"

	"github.com/lucblassel/fastago/pkg/seqs"
	"github.com/spf13/cobra"
)

// countCmd represents the count command
var countCmd = &cobra.Command{
	Use:   "count",
	Short: "count the number of sequences in a fasta file",
	RunE: func(cmd *cobra.Command, args []string) error {
		records := make(chan seqs.SeqRecord)
		errs := make(chan error)
		go seqs.ReadFastaRecords(inputReader, records, errs)

		count := 0

		for records != nil && errs != nil {
			select {
			case <-records:
				count++
			case err := <-errs:
				if err != nil {
					return err
				}
				fmt.Fprintln(outputWriter, count)
				return nil

			}
		}

		return nil
	},
}

func init() {
	statsCmd.AddCommand(countCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// countCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// countCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func countSeqs(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	count := 0

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 && line[0] == '>' {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error in counting sequences: ", err)
	}

	fmt.Fprintln(output, count)
}
